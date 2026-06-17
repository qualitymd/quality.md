---
type: Design Doc
title: delightful CLI output — design doc
description: How qualitymd styles its human output, pages spec, and resolves its version on the existing Charm stack without disturbing the agent-facing paths.
tags: [cli, dx, output, design]
timestamp: 2026-06-17T00:00:00Z
---

# delightful CLI output — design doc

Design behind the [Delightful human CLI output](../0007-delightful-cli-output.md)
change and its [functional spec](spec.md).

## Context

The stack is already in place: Fang renders styled help, errors, and version;
Glamour renders the `spec` Markdown. What was missing was the commands' own
output and a shared palette. The agent-accessibility baseline forbids changing
the plain and JSON paths, so the whole design is "a styled branch behind one
gate, plain branch unchanged."

## One gate, one palette

`internal/cli/style.go` is the single home for the human-presentation concept:

- `colorEnabled(w io.Writer) bool` generalizes the predicate `spec.go` already
  used — `NO_COLOR` unset **and** `w` exposes an `Fd()` that is a terminal. Every
  command routes through it, so the gate is defined once. The old
  `shouldRenderSpec` is replaced by it.
- `brandColorScheme(c lipgloss.LightDarkFunc) fang.ColorScheme` starts from
  `fang.DefaultColorScheme` and overrides the title and program name to the brand
  green, wired through `fang.WithColorSchemeFunc` in `root.go`. The command
  renderers use the same charmtone colors, so harness and commands match.
- A small set of `lipgloss` styles and glyph constants (`✓ ✗ ⚠ •`) the styled
  renderers share. Colors are charmtone keys chosen to stay legible on light and
  dark terminals, so no background detection is needed for the command output.

Keeping this in `internal/cli` rather than a new package follows the
[Go package guide](../../../docs/guides/design-go-packages.md): it is
CLI-presentation behavior with a single consumer, and the rule of three says not
to extract a shared package before a second one appears.

## Styled renderers as a branch over the plain form

`lint` and `init` each split into a plain renderer and a styled renderer behind
`colorEnabled`:

- `renderLintPlain` is the existing code verbatim — it is the canonical, stable
  form and the one the tests assert byte-for-byte.
- `renderLintStyled` carries the same facts with a glyph, color, and a colored
  `lintSummary`. `findingLocation` appends a clickable `path:line` to the label
  when the finding has a source position, so the label (which the `lint` sub-spec
  requires in human output) is never dropped.
- `renderInitHuman` does the same split inline for the one-line confirmation,
  still writing to stderr.

Because the gate is closed for a `bytes.Buffer`, the existing CLI tests exercise
the plain branch unchanged; new unit tests call the styled renderers directly and
assert on the glyphs and substrings (not on ANSI), so they are robust whether or
not lipgloss emits color under `go test`.

## Paging spec

`spec.writeSpec` keeps writing verbatim bytes when the gate is closed. When open,
it renders with Glamour as before and hands the bytes to `page`:

- `page` only engages a pager when the writer is the real `*os.File` terminal; a
  non-file writer falls back to a direct write.
- It prefers `$PAGER` (run via `sh -c`), else `less -R -F` (`-R` passes the
  rendered ANSI; `-F` quits if it fits one screen). If neither is available or
  the pager fails to start, it writes the bytes directly.
- Paging is therefore never load-bearing, and the byte-for-byte redirect
  invariant is untouched because redirects close the gate.

## Version resolution

`fang` only consults `debug.ReadBuildInfo` when handed an empty version, and the
binary always passed the `"dev"` placeholder, so `go install` and local builds
showed `dev`. `buildInfo()` in `root.go` now resolves it ourselves:

- Release builds keep the goreleaser-stamped `version`/`commit`.
- Otherwise it reads `debug.ReadBuildInfo`: an installed module release reports
  its module version; a local build keeps the friendly `dev` label and attaches
  the `vcs.revision`, which Fang renders as `dev (<short-sha>)`.

This is strictly more informative than the bare placeholder and stays accurate
for releases.

## Alternatives considered

- **Interactive `huh`/Bubble Tea touches** (prompts, spinners) — rejected: they
  violate the non-interactive baseline and inject nondeterminism into captured
  output.
- **Letting Fang derive the version from an empty string** — rejected: it yields
  `unknown (built from source)` for a local build, which is less friendly than
  `dev (<short-sha>)` for the common contributor case.
- **A separate `style`/`theme` package** — deferred under the rule of three; one
  consumer does not justify the extraction yet.
