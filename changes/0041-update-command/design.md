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
short-circuits to the advisory path; `--json` emits the result struct — an output
modifier that never changes whether the command applies, so the struct carries an
`applied` boolean that lets a consumer tell an apply from a `--check`. `--apply`
is removed — applying is the default. `upgrade` is renamed outright to `update`
with no alias; the paired `/quality` skill mode is renamed and version-pinned in
lockstep, so nothing keeps calling `qualitymd upgrade`.
The file moves `internal/cli/upgrade.go` → `update.go` and its test
`internal/cli/version_upgrade_test.go` → `version_update_test.go`; symbols rename
to match.

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
(strict SemVer precedence; `updateAvailable` returns false when either version is
not valid SemVer or is a prerelease — the old "report any difference" fallback is
dropped, so such builds never report or apply an update, per the spec) and reports
`UpdateAvailable = newer && latest.Ready && !isDevBuild()`.
A development build is a first-class no-update condition, not a side effect of
SemVer parsing: a single `isDevBuild()` predicate (read from the build-stamped
version metadata) short-circuits `checkUpdate` to `UpdateAvailable=false`, refuses
apply, and is the exact dev-build gate the ambient notice and refresh reuse — one
source of truth, so a dev build carrying a SemVer-looking stamp can't slip through.
The default (apply) path refuses with a clear diagnostic when `newer &&
!latest.Ready`, before any mutation. The result struct also carries `latest.Ready`
through to an explicit `latestVersionReady` field in `--json`, so a consumer can
tell "newer release exists but isn't downloadable yet" apart from "up to date" —
the legibility the spec's unconfirmed-state SHOULD calls for; inferring it from
`updateAvailable: false` + a populated `latestVersion` is ambiguous (that combo
also describes being current).

`Ready` is resolved from the response already fetched — no extra round trip:

- GitHub-backed channels (managed standalone, and the archive/source guidance
  paths): parse `assets[]`; `Ready` requires both `qualitymd_<os>_<arch>.tar.gz`
  (reusing the installer's GOOS/GOARCH → token mapping in one helper) and
  `checksums.txt`. `ReleaseNotesURL` = the release `html_url`.
- npm: the registry `/latest` *is* the installable latest, so `Ready` is true
  whenever a version returns; `ReleaseNotesURL` left empty (no GitHub URL in the
  registry response — see Alternatives).
- Homebrew: route to a new `latestHomebrewVersion` that reads the tap cask the
  project publishes — `GET https://raw.githubusercontent.com/qualitymd/homebrew-tap/main/Casks/qualitymd.rb`
  — and parses its `version "X.Y.Z"` line (a small regex; the cask is generated,
  so the line is stable). The cask *is* what `brew upgrade qualitymd/tap/qualitymd`
  installs, so like npm `Ready` is true whenever a version parses — never gated on
  the GitHub release tag. Unlike npm, brew gets a real `ReleaseNotesURL`: the
  cask's `url` pins `releases/download/v<version>/…`, so the notes URL is
  `releases/tag/v<version>` for the same parsed version — known-good, not
  constructed-and-hoped (contrast the npm rejection in Alternatives). This drops
  brew out of `defaultLatestVersion`'s `default:` GitHub branch, which was the
  bug: brew was resolving latest from the GitHub tag instead of the tap.

Brew apply is unchanged (`brew upgrade qualitymd/tap/qualitymd`), but a user's
*local* Homebrew can still lag the published tap until it refreshes — which a
version check can't see. So the post-apply `qualitymd --version` verify that the
managed standalone path already performs is now required on every apply path (the
spec's Apply verification requirement), brew most notably: if the running version
didn't advance after `brew upgrade`, `update` reports that honestly (suggesting a
tap refresh) instead of a false "applied." Tap-source readiness fixes the advisory
signal; the verify keeps apply honest against local staleness.

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
command, and the notes URL when known. Gates: the cache says a newer, ready
version exists, `ambientAllowed()` holds (opt-out unset / not a dev build / no CI
env / stderr is a terminal — see Refresh), and the command is not `--json`
(inspect the invoked command's flag). Any failure to read the cache is
swallowed — the notice never affects exit code or output. The notice shows on
every qualifying run (no rate-limit, no persisted dismissal) — the gating already
confines it to interactive non-`--json` humans, and a single quiet stderr line
per run matches npm/gh/brew.

**Refresh.** Notice and refresh share one `ambientAllowed()` predicate — opt-out
unset, not a dev build, no CI env (`CI`), stderr interactive (`term.IsTerminal`) —
and the notice additionally checks `--json`. The refresh spawns only when
`ambientAllowed()` holds and the cache is older than the TTL (20h). It gates on CI
explicitly rather than leaning on "TTY implies not-CI": some CI runners allocate a
pseudo-TTY, and a CI run must never spawn a network subprocess the user didn't
invoke. The refresh needs no `--json` gate of its own — `--json` only suppresses
the visible line, not the background warming — and skipping dev builds avoids a
network subprocess for a build that can never report an update. When those gates pass, the foreground process spawns a *detached*
`qualitymd` subprocess from the root `PersistentPreRunE` — so the child has
maximum lead time, overlaps the real command, and still spawns even if the command
later errors before post-run. The child performs the check and rewrites the cache, then the parent
returns immediately without waiting. A `QUALITYMD_UPDATE_REFRESH=1` env marker on
the child makes it do only the refresh (no notice, no recursive spawn). The
triggering command uses the existing cache and exits; the fresh value appears on
a later invocation — exactly the non-blocking, surface-next-run semantics the
spec requires. The explicit `update`/`update --check` paths also rewrite the
cache as a side effect, so an interactive check keeps it warm.

**Opt-out.** `QUALITYMD_NO_UPDATE_CHECK=1` disables both the refresh and the
notice — the documented contract for now. A config-file key is the intended
*primary* opt-out for centrally managed and air-gapped fleets once qualitymd grows
a config surface (an env var doesn't persist across a managed fleet the way a
checked-in config setting does, which is exactly the audience the opt-out serves),
so the env var should not ossify as the only knob. This mirrors how Codex exposes
its update check solely as a `check_for_update_on_startup` config flag documented
"set to false only if your updates are centrally managed," with no env-var gate.

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
  `v`-prefixed tag format and can 404; honest omission beats a broken link. (Brew
  is different: the cask *pins* the release it installs, so the tag is known-good,
  not assumed — hence brew gets a notes URL and npm doesn't.)
- **Resolve brew latest via `brew info` / a local brew invocation.** Rejected —
  reads the user's *local* tap state, which is exactly the staleness we don't want
  to treat as authoritative, and shells out to brew on a plain version check.
  Reading the published cask over HTTP gives the canonical "what brew would install
  after a tap refresh." Local staleness is instead caught by the post-apply verify.
- **Keep brew on the GitHub release tag (defer the fix).** Rejected — 0041's
  readiness gate would then confirm GitHub assets, a signal `brew upgrade` doesn't
  use, advertising upgrades the tap can't yet deliver. The gate is the reason this
  had to be fixed here, not deferred.
- **Keep a deprecated `upgrade` alias.** Rejected — the only callers are the
  version-pinned `/quality` skill (renamed in lockstep) and scripts; carrying an
  alias whose no-flag behavior silently flips from advisory to apply is more
  surprising than a clean rename, and there is no external compatibility contract
  to preserve.

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
- Reading the brew tap cask couples the CLI to a published-file shape (the cask's
  `version` line and tap path) the project itself generates via GoReleaser, so it
  is stable but not contractual. Mitigation: a parse miss is treated like any
  failed latest check — no update reported, never an error — and the apply path's
  `--version` verify backstops a wrong/stale read. The tap path is a single
  constant if the repo or cask name ever changes.

## Resolved decisions

- **TTL and notice cadence.** 20h refresh TTL (under a day, so a user who runs at
  the same clock time daily isn't perpetually just short of the window — the same
  reason Codex uses a sub-24h, 20h TTL); show the notice on every qualifying run
  until updated, with no rate-limit and no persisted dismissal. The existing gating
  confines it to interactive non-`--json` humans, and one quiet stderr line per run
  matches npm/gh/brew. Revisit only if it proves noisy.
- **Refresh trigger placement.** Spawn the detached refresh from
  `PersistentPreRunE` (max lead time, overlaps the command, fires even on later
  error) and emit the notice from `PersistentPostRunE` (after stdout). Detached
  and non-blocking, so neither hook delays the command.
- **`latestVersionReady` in `--json`.** Expose an explicit readiness field rather
  than inferring from `updateAvailable: false` + a populated `latestVersion`,
  which is ambiguous with "up to date." Cheap, and it satisfies the spec's
  legibility SHOULD for the unconfirmed state.
- **Homebrew latest source.** Pulled into scope (not deferred): resolve brew latest
  and readiness from the published tap cask, not the GitHub release tag, and extend
  the post-apply `--version` verify to the brew path. This fixes the pre-existing
  0032 bug — brew was comparing against GitHub — and implements 0032's unfulfilled
  "compare against the owning channel" SHOULD. Deferring it would leave 0041's
  readiness gate confirming the wrong signal for brew. Mechanics in *Widen the
  latest-version provider*.
