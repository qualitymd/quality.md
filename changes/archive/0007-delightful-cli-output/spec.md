---
type: Functional Specification
title: delightful CLI output — functional spec
description: What qualitymd's human output surface must do — a cohesive styled terminal experience layered over the unchanged agent-facing plain and JSON paths.
tags: [cli, dx, output]
timestamp: 2026-06-17T00:00:00Z
---

# delightful CLI output — functional spec

Companion to the [Delightful human CLI output](../0007-delightful-cli-output.md)
change. This spec states *what* the human surface must do; the
[design doc](design.md) covers *how*, on the existing stack.

It inherits the cross-cutting contract — the
[baseline](../../../specs/cli.md#baseline), exit codes, and the `--json` and
output conventions — from the [CLI spec](../../../specs/cli.md). The durable
result of this change is folded into that spec's **Human output styling** and
**Binary version** conventions and the [`spec` sub-spec](../../../specs/cli/spec.md)'s
paging clause.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Scope

Covered: the brand palette, the styled `lint` and `init` human output, `spec`
paging, `--help` examples, and `--version` resolution.

Deferred: the global quiet/verbosity control, freezing the styled layout as
stable API, and any interactive TUI.

## The gate

The single rule the whole surface obeys:

- Styling and paging **MUST** apply only when the destination stream is a
  terminal and `NO_COLOR` is unset. This is the same gate the baseline already
  states for color and the `spec` renderer.
- When the gate is closed — a pipe, a redirect, a non-terminal writer, or
  `NO_COLOR` — output **MUST** be the canonical plain form, byte-for-byte.
  Existing plain output **MUST NOT** change.
- `--json` documents **MUST NOT** be affected in any way: no styling, no glyphs,
  no terminal control sequences, identical bytes.

## Brand palette

- The harness (help, errors, version) and the commands' own human output **MUST**
  draw their colors from one shared palette, so the tool reads as a single
  program rather than a styled harness wrapping unstyled commands.
- The palette **SHOULD** be applied through the stack's idiomatic colorscheme
  hook rather than per-command ad-hoc colors.

## Styled lint output

When the gate is open, `lint` human output:

- **MUST** carry the same facts as the plain form — for each finding the
  severity, rule id, message, and location label per the
  [`lint` sub-spec](../../../specs/cli/lint.md#human-output).
- **SHOULD** lead each finding with a severity glyph and color, and render the
  count summary with color.
- **SHOULD** append a clickable `file:line` to a finding's location label when
  the finding carries a source position, without dropping the label.
- **MUST** still report a valid file and the applied-repair count, matching the
  plain form's information.

## Styled init output

When the gate is open, the `init` success confirmation **SHOULD** lead with a
success glyph and render its next-action command in the palette's accent. It
**MUST** keep writing to standard error and **MUST** carry the created path and
the next command, matching the plain confirmation the
[`init` sub-spec](../../../specs/cli/init.md#reporting) requires.

## spec paging

`spec` **MAY** display its rendered specification through the user's pager when
stdout is a terminal, per the [paging convention](../../../specs/cli.md#conventions).
Paging **MUST** be skipped when the gate is closed, so a redirect reproduces the
artifact byte-for-byte; and it **MUST NOT** be load-bearing — if no pager is
available or it fails to start, `spec` writes the rendered output directly.

## Help examples

`init`, `lint`, and `spec` **SHOULD** each carry a short block of runnable
example invocations in their help, so `--help` teaches the command's common uses.

## Version

`--version` **MUST** report the release-stamped version when present, and
otherwise **MUST** fall back to the Go toolchain's embedded module build
information — the module version for an installed release, otherwise a
development label carrying the VCS revision — so the reported version is never a
bare placeholder. This is the durable rule in the CLI spec's
[Binary version](../../../specs/cli.md#conventions) convention.
