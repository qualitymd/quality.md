---
type: Change Case
title: Delightful human CLI output
description: Give qualitymd's human surface a cohesive brand palette, styled lint and init output, runnable help examples, spec paging, and an informative --version, without touching the agent-facing plain and JSON paths.
status: Done
tags: [cli, dx, output]
timestamp: 2026-06-17T00:00:00Z
---

# Delightful human CLI output

`qualitymd` already runs on the Charm stack — [Cobra](https://github.com/spf13/cobra)
for structure, [Fang](https://github.com/charmbracelet/fang) for the harness,
[Glamour](https://github.com/charmbracelet/glamour) for `spec` — so help, errors,
and the rendered specification are already styled. The _commands' own_ human
output was not: `lint` and `init` printed flat monochrome text from
`fmt.Fprintf`, help carried no examples, a long `spec` scrolled off-screen, and
`--version` printed a bare `dev`. The styled harness and the unstyled command
output read as two different programs. This change closes that gap on the human
path only. The detail lives in its children:

- [Functional spec](0007-delightful-cli-output/spec.md) — what the human surface
  must do.
- [Design doc](0007-delightful-cli-output/design.md) — how it is built on the
  existing stack without disturbing the agent contract.

## Motivation

The "moment of delight" the Charm tools are known for comes from _cohesion_ —
every surface sharing one palette and a consistent visual grammar. `qualitymd`
had the harness half of that and not the command half. The fix is cheap because
the stack is already present; the work is wiring a single brand palette through
Fang and the command renderers, adding glyphs and a colored summary to `lint`,
dressing up the `init` confirmation, teaching `spec` to page, giving each command
runnable `--help` examples, and making `--version` informative for
`go install` and local builds.

The hard constraint is the CLI's
[agent-accessibility baseline](../../specs/cli.md#agent-accessibility): the
surface is contractually deterministic and machine-drivable. So every bit of
delight is a **terminal-only** layer over a canonical plain form, gated on a TTY
destination and `NO_COLOR`, and the `--json` documents and redirected/piped
bytes are untouched. No interactive prompts, no spinners, no nondeterminism enter
the core paths.

## Scope

Covered:

1. **Brand palette.** One Fang colorscheme (`WithColorSchemeFunc`) the harness
   and the command renderers both draw from, so help, errors, version, and
   command output share a look.
2. **Styled `lint`.** A severity glyph and color per finding, a clickable
   `file:line` appended to the location label where a source position exists, and
   a colored count summary — on the terminal path only, with the existing plain
   text preserved byte-for-byte off-terminal.
3. **Styled `init`.** A success glyph and an accented next-command on the
   terminal path; the existing plain confirmation off-terminal.
4. **`spec` paging.** Render through the user's pager (`$PAGER`, else `less`)
   when stdout is a terminal; verbatim direct write otherwise, with paging never
   load-bearing.
5. **Help examples.** A `cobra` `Example` block on `init`, `lint`, and `spec`.
6. **Informative `--version`.** Fall back to the Go toolchain's embedded module
   build info (module version, else a `dev` label with the VCS revision) when the
   binary was not release-stamped.

Deferred:

- The global quiet/verbosity control that would suppress the human footers (still
  owned by the [CLI spec](../../specs/cli.md)).
- Freezing the exact styled layout as stable API; only the plain and JSON forms
  are contracts.
- Any interactive TUI (forms, prompts, spinners) — excluded by the baseline.

## Affected specs & docs

Updated before this change reached **Done**:

- [x] [`specs/cli.md`](../../specs/cli.md) — added the **Human output styling**
      and **Binary version** conventions, and narrowed the matching
      _To be specified_ entries.
- [x] [`specs/cli/spec.md`](../../specs/cli/spec.md) — added the optional paging
      behavior, riding the same terminal-detection rule as rendering.

No change to [`specs/cli/lint.md`](../../specs/cli/lint.md) or
[`specs/cli/init.md`](../../specs/cli/init.md): both already delegate styling and
exact layout to the CLI spec's output conventions, which this change fills in.
The implementation in `internal/cli` (the shared `style.go`, the styled
renderers, the paging helper, the examples, and the version resolution) is this
change's implementation, tracked by the functional spec and design doc.

## Status

`Done`. Implemented and archived after adding the brand palette and shared
styling helper, styling the `lint` and `init` human output behind a TTY/`NO_COLOR`
gate while preserving the plain and JSON paths byte-for-byte, teaching `spec` to
page, adding `--help` examples, recovering an informative `--version` from build
info, updating the two durable specs, and adding focused tests.
