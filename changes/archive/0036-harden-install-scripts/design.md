---
type: Design Doc
title: Harden install scripts and upgrade idiomatics — design doc
description: How the installer and upgrade-check fixes are built, and why the Homebrew cask stays.
tags: [cli, install, upgrade, homebrew]
timestamp: 2026-06-19T00:00:00Z
---

# Harden install scripts and upgrade idiomatics — design doc

Design behind the
[Harden install scripts and upgrade idiomatics](../0036-harden-install-scripts.md)
change case and its [functional spec](spec.md). Each fix is local to one surface;
this doc records the decisions that aren't obvious from the requirement.

## Context

Five narrow portability/convention fixes across `install/install.sh`,
`install/install.ps1`, and `internal/cli/upgrade.go`, plus a durable note that
the Homebrew cask is intentional. No new dependencies on the install surface; one
small Go dependency for SemVer.

## Approach

### Checksum tool fallback (`install.sh`)

Replace the single `command -v shasum` gate with a small helper that returns the
archive's SHA-256 from the first available of `sha256sum`, `shasum -a 256`,
`openssl dgst -sha256` (each prints the digest as the first whitespace field;
`openssl` output is normalized). If `checksums.txt` was fetched but no tool is
found, or the archive is absent from it, print a `warning: …` to stderr and
proceed; a positive mismatch still aborts. Verification stays best-effort by
design (the script is fetched over TLS from the same origin), but the *skip*
becomes visible instead of silent.

### TLS 1.2 (`install.ps1`)

Add, before the first request:

```powershell
try { [Net.ServicePointManager]::SecurityProtocol =
  [Net.ServicePointManager]::SecurityProtocol -bor [Net.SecurityProtocolType]::Tls12 } catch {}
```

`-bor` preserves any already-enabled protocols; the `try/catch` keeps PowerShell
7 (where the property is effectively a no-op) from erroring. This is the standard
PS 5.1 compatibility shim.

### PATH integration — deliberately asymmetric

- **Windows (`install.ps1`)**: write the per-user PATH via
  `[Environment]::SetEnvironmentVariable('Path', …, 'User')` when `bin` is absent,
  prepend it to `$env:Path` for the current session, and print a "restart your
  shell" note. This is what Scoop and rustup-init do; Windows has a clean,
  reversible per-user PATH registry API.
- **Unix (`install.sh`)**: print the exact `export PATH="<bin>:$PATH"` line and
  stop. No dotfile editing.

The asymmetry is intentional and matches platform norms: Windows offers a safe
structured API, while Unix profile editing is fragile (which file? login vs
interactive? piped `curl | sh` has no reliable TTY for consent). Printing the
line is the Deno-style floor; auto-editing dotfiles is out of scope.

### Non-interactive semantics (`install.sh` / `install.ps1`)

`--non-interactive` / `--yes` / `QUALITYMD_NO_INPUT=1` now gate the *human
guidance*, not a phantom prompt: interactive runs print the PATH guidance and a
friendly summary; non-interactive runs print a single completion line and skip
the guidance. The flags stay accepted (the install-smoke workflow and documented
CI forms pass them) and the install/PATH actions are identical either way — only
verbosity changes. This makes the flag observable without inventing a TTY-driven
prompt flow that piped installs can't honor.

### SemVer update detection (`upgrade.go`)

Add `golang.org/x/mod/semver` (already in the module graph) and rewrite
`updateAvailable`: normalize both versions, re-add the `v` prefix
`semver.Compare`/`semver.IsValid` expect, and report an update only when
`semver.Compare(latest, current) > 0`. When either side isn't valid SemVer, fall
back to the prior string-inequality behavior so unusual version strings still
surface *a* difference. Development builds (`dev…`) never report an update, as
today.

### Homebrew cask — keep, document

No config change. Add a comment in `.goreleaser.yaml` that `homebrew_casks` is the
GoReleaser-recommended path (the `brews` formula path was deprecated in v2.10 and
is removed in v3) and the quarantine hook is the documented unsigned-binary
pattern. Mirror the rationale in `CONTRIBUTING.md` and `cut-a-release.md`.

## Alternatives

- **Convert the cask to a Homebrew formula** (the original review item).
  Rejected: GoReleaser deprecated the `brews` formula path in v2.10 (removed in
  v3), Homebrew steers pre-built binaries to casks, and `brew audit --new`
  rejects binary-only formulae. It would move onto a deprecated path to satisfy a
  pre-2024 rule of thumb. The genuinely-more-idiomatic option, a `homebrew-core`
  *source* formula (`brew install qualitymd`), is notability-gated and
  externally maintained — recorded as a non-goal, not done here.
- **Auto-edit Unix shell profiles** (Bun-style). Rejected for now: requires TTY
  consent or a profile-detection heuristic that misfires across shells; printing
  the export line is safe and sufficient.
- **Hand-rolled SemVer comparison.** Rejected: `x/mod/semver` is small, already
  in the graph, and correct on prereleases/build metadata where a hand-rolled
  split would drift.
- **Drop `--non-interactive` entirely.** Rejected: the install-smoke workflow and
  documented CI forms pass it; removing it would break them with "unknown
  argument."

## Trade-offs & risks

- Adding `golang.org/x/mod` as a direct dependency is a (tiny, std-adjacent)
  surface increase; `go mod tidy` records it.
- Windows PATH mutation is a persistent side effect. It's per-user and
  idempotent, and matches install-tool expectations, but it is the one fix that
  writes outside the install root. Guarded by an "already present" check.
- Checksum verification remains best-effort (skipped, loudly, when no tool
  exists) rather than mandatory; making it mandatory would break minimal images
  that lack all three tools.

## Open questions

None blocking. Signing/notarizing the macOS binary (which would let the cask
quarantine hook be dropped) is deferred to a future change.
