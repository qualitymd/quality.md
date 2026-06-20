---
type: Design Doc
title: Update command and improvements design
description: How qualitymd renames upgrade to update with apply-by-default, self-applies managed standalone, gates on readiness, surfaces release notes, adds a cached ambient update notice without blocking ordinary commands, and renames the paired /quality skill mode.
tags: [cli, update, upgrade, install, versioning, skill, design]
timestamp: 2026-06-20T00:00:00Z
---

# Update command and improvements design

Design behind the [Update command and improvements](../0041-update-command.md)
change and its [functional spec](spec.md).

## Context

[0032](../archive/0032-cli-managed-upgrades.md) built the self-update surface as
layers in `internal/cli/upgrade.go`: install-context detection, an injectable
latest-version provider, and a guarded `--apply`. This change reshapes that
command to an apply-by-default `update` verb and adds a cached ambient notice. The
hard constraint, inherited from qualitymd's agent-first posture, is that ordinary
commands must stay deterministic, offline, and clean on stdout/`--json` — so the
notice has to be cache-backed and tightly gated.

The paired `/quality upgrade` → `/quality update` skill-mode rename is a
mechanical rename of the skill files and the durable skill spec — no design
discussion; it carries no behavior change beyond the new name and the
`qualitymd update` commands the mode drives.

Relevant existing shape:

- `fetchLatestVersion latestVersionProvider` returns a bare `string`, swapped in
  tests to run offline.
- `applySupported(method)` is a fixed allowlist (`npm`, `homebrew`);
  `upgradeCommand(method)` returns a `(name, args)` pair run via
  `exec.CommandContext`. `recommendedCommand(installManagedStandalone)` already
  emits the GitHub-hosted installer one-liner per OS.
- The installer is idempotent, serves update, resolves `latest`, and honors
  `QUALITYMD_NO_INPUT=1`. The GitHub `releases/latest` response already carries
  `assets[]` and `html_url`.

## Approach

### Rename and reshape the command

Rename the command to `update`. The default `RunE` applies; a `--check` flag
short-circuits to the advisory path; `--json` emits the result struct. `--apply`
is removed — applying is the default. Register a second cobra command `upgrade`
marked `Deprecated` (cobra prints the deprecation line automatically) that shares
`update`'s `RunE`, so `qualitymd upgrade [--check]` keeps working for one cycle.
The file moves `internal/cli/upgrade.go` → `update.go`; symbols rename to match.

Apply gating keeps 0032's refusal for unowned channels, just relocated to the
default path: unknown/archive/Go-source still refuse and print guidance.

### Widen the latest-version provider

The provider returns a struct so readiness and the notes URL ride the one
injectable, offline-testable network call:

```go
type latestRelease struct {
    Version         string // tag/version for the owning channel ("" when none)
    Ready           bool   // target is actually retrievable for this channel
    ReleaseNotesURL string // omitted when not known
}
type latestVersionProvider func(context.Context, installMethod) (latestRelease, error)
```

`checkUpdate` computes `newer := updateAvailable(current, latest.Version)`
(unchanged SemVer logic) and reports `UpdateAvailable = newer && latest.Ready`.
The default (apply) path refuses with a clear diagnostic when `newer &&
!latest.Ready`, before any mutation.

`Ready` is resolved from the response already fetched — no extra round trip:

- GitHub-backed channels: parse `assets[]`; `Ready` requires both
  `qualitymd_<os>_<arch>.tar.gz` (reusing the installer's GOOS/GOARCH → token
  mapping in one helper) and `checksums.txt`. `ReleaseNotesURL` = the release
  `html_url`.
- npm: the registry `/latest` *is* the installable latest, so `Ready` is true
  whenever a version returns; `ReleaseNotesURL` left empty (no GitHub URL in the
  registry response — see Alternatives).

### Managed standalone self-apply

Add `installManagedStandalone` to `applySupported`, and have `upgradeCommand`
return the installer invocation for it — the same one-liner `recommendedCommand`
emits, run non-interactively through the shell:

- Unix: `sh -c "curl -fsSL <install.sh raw URL> | QUALITYMD_NO_INPUT=1 sh"`
- Windows: the existing PowerShell one-liner with `$env:QUALITYMD_NO_INPUT=1`.

The installer already stages under `releases/<version>`, flips the
`bin/qualitymd` symlink atomically, and verifies `--version`, so the binary adds
no replacement logic. npm/Homebrew apply paths are unchanged.

### Ambient notice: cache + non-blocking refresh

This is the load-bearing design decision, because a CLI process is short-lived —
the "refresh in a background task, surface next run" pattern can't be a goroutine
that outlives the process.

**Cache.** A small JSON record under `$QUALITYMD_HOME` (default `~/.qualitymd`),
e.g. `.qualitymd-update-cache`: `{ latestVersion, releaseNotesURL, ready,
checkedAt }`.

**Emit.** A cobra root `PersistentPostRunE` reads the cache and, when all gates
pass, writes one stderr line naming current → latest, the `qualitymd update`
command, and the notes URL when known. Gates: cache says a newer, ready version
exists; stderr is a terminal (`term.IsTerminal`); not `--json` (inspect the
invoked command's flag); no CI env (`CI`); opt-out unset
(`QUALITYMD_NO_UPDATE_CHECK`); not a dev build. Any failure to read the cache is
swallowed — the notice never affects exit code or output.

**Refresh.** When checks are enabled and interactive and the cache is older than
the TTL (e.g. 24h), the foreground process spawns a *detached* `qualitymd`
subprocess that performs the check and rewrites the cache, then returns
immediately without waiting. A `QUALITYMD_UPDATE_REFRESH=1` env marker on the
child makes it do only the refresh (no notice, no recursive spawn). The
triggering command uses the existing cache and exits; the fresh value appears on
a later invocation — exactly the non-blocking, surface-next-run semantics the
spec requires. The explicit `update`/`update --check` paths also rewrite the
cache as a side effect, so an interactive check keeps it warm.

**Opt-out.** `QUALITYMD_NO_UPDATE_CHECK=1` disables both the refresh and the
notice. (A config-file key can follow if/when qualitymd grows a config surface;
the env var is the documented contract now.)

## Alternatives

- **Keep the provider returning `string`; separate readiness/notes calls.**
  Rejected — readiness and notes come from the same `releases/latest` response;
  one struct keeps the network call and test seam singular.
- **HEAD each asset URL to probe readiness.** Rejected — `assets[]` is already in
  the JSON; a HEAD per asset adds latency for no new information.
- **Synchronous bounded fetch (short timeout) for the notice.** Rejected — even a
  1s fetch violates "ordinary commands MUST NOT block on the network" and would
  tax every interactive command.
- **Refresh only during explicit `update`/`update --check`.** Rejected as the
  sole mechanism — it never warms the cache for users who don't run checks, which
  is exactly the audience the notice serves. Kept as a secondary warming path.
- **Detached subprocess vs. a goroutine.** A goroutine dies on process exit; the
  detached subprocess is the only reliable way to refresh without blocking. Cost:
  cross-platform process detachment and an extra short-lived process (see risks).
- **Construct an npm release-notes URL from the version.** Rejected — assumes the
  `v`-prefixed tag format and can 404; honest omission beats a broken link.
- **Hard rename with no alias.** Rejected — the version-pinned `/quality` skill
  and existing scripts call `upgrade`; a deprecated alias avoids breaking the
  paired flow mid-transition.

## Trade-offs and risks

- The ambient notice reverses 0032's offline-by-default rule. The risk is output
  contamination for agents/CI; the mitigation is strict gating (stderr-only,
  `--json`/TTY/CI/opt-out/dev-build), so agent and pipeline invocations never see
  it. This is the crux to get right in review and tests.
- Detaching a refresh subprocess differs by platform (`Setpgid`/`Setsid` on Unix,
  `CREATE_NEW_PROCESS_GROUP`/`DETACHED_PROCESS` on Windows) and spawns a brief
  network-touching process the user didn't directly invoke — acceptable given the
  opt-out and TTL, but it must be invisible (no output, no error propagation).
- Widening the provider return type touches every call site and its tests, but
  all are confined to `internal/cli`.
- The deprecated `upgrade` alias is one extra command to carry and later remove;
  cheap, and it keeps the skill working across the version bump.

## Open questions

- **TTL and notice cadence.** A 24h refresh, and show the notice on every
  qualifying run until updated? Or rate-limit the notice itself to avoid nagging?
- **Refresh trigger placement.** Spawn the detached refresh from
  `PersistentPreRun` (so it overlaps the command) or `PersistentPostRun` (after
  output)? Pre-run overlaps better but must not delay the command.
- **Homebrew latest source.** Homebrew installs currently resolve "latest" via
  GitHub (a pre-existing 0032 quirk), so readiness checks GitHub assets for brew.
  Worth fixing the Homebrew provider, but out of this change's scope — flagged,
  not silently widened.
- **`releaseReady` in `--json`.** Expose an explicit readiness field, or let
  agents infer from `updateAvailable: false` + a populated `latestVersion`?
  Leaning to the latter to avoid schema growth.
