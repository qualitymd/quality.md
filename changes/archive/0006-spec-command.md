---
type: Change Case
title: Specify and implement the spec command
description: Settle the durable spec sub-spec for `qualitymd spec` and implement the command, emitting the bundled QUALITY.md format specification to stdout.
status: Done
tags: [cli, command, spec]
timestamp: 2026-06-17T00:00:00Z
---

# Specify and implement the spec command

The [`qualitymd spec`](../../specs/cli/spec.md) sub-spec was still a placeholder
stub, and the command was unimplemented — invoking it failed with "unknown
command". [Change Case 0004](0004-specify-agent-accessibility.md) settled the
cross-cutting CLI contract and named `spec` as the deliberate verbatim-artifact
carve-out, deferring the command itself to "a separate change that inherits this
baseline." This change settles the durable sub-spec and lands the command. The
detail lives in its children:

- [Functional spec](0006-spec-command/spec.md) - what `spec` must do.
- [Design doc](0006-spec-command/design.md) - how the bundled specification is
  embedded and emitted.

## Motivation

`qualitymd spec` exists so an author or — more often — a coding agent can inject
the current `QUALITY.md` format rules into its context without reaching for an
external copy. It is the last of the three format-tooling commands in the CLI
spec's [Scope](../../specs/cli.md#scope) still unbuilt, and the only one whose
durable sub-spec is still a stub. The stub's open questions also predate the
agent-accessibility work and now conflict with it — it floats a `--format json`
form, which the settled [`--json` convention](../../specs/cli.md#conventions)
forbids both in spelling and in spirit (`spec` is the carve-out that emits a
verbatim artifact and offers no `--json`). Settling the sub-spec against the
current contract and implementing the command closes the format-tooling layer.

## Scope

Covered:

1. **Specify.** Replace the placeholder [`specs/cli/spec.md`](../../specs/cli/spec.md)
   with a settled sub-spec: the output is the bundled format specification —
   verbatim Markdown when output must be plain (so a redirect reproduces the
   artifact byte-for-byte) and rendered for readability when stdout is a terminal,
   the split riding the baseline's terminal detection as color does; `spec` takes
   no arguments; it is the verbatim-artifact carve-out and offers no `--json`; and
   it owes the full baseline (categorized exit codes, stderr diagnostics,
   determinism). This empties the sub-spec's **To be specified** list.
2. **Implement.** Register `qualitymd spec` in `internal/cli`, emitting the
   format specification — embedded at build time so the binary is
   self-contained — to stdout, rendering it for a terminal and writing verbatim
   Markdown when output must be plain, with tests covering the verbatim payload,
   the clean-stdout/redirect invariant, and the exit-code behavior.
3. **Document.** Mark `spec` as built in the [`README.md`](../../README.md) status.

Deferred (recorded as **deferred** in the sub-spec, not specified here):

- A structured/JSON form of the specification — `spec` stays the carve-out.
- Sub-views such as emitting only the `lint` rule set, and the flags that would
  select them.
- Any `spec`-specific flags beyond the cross-cutting CLI flags.

## Affected specs & docs

Updated before this change reaches **Done**:

- [x] [`specs/cli/spec.md`](../../specs/cli/spec.md) — replace the placeholder stub
      with the settled sub-spec (verbatim-Markdown output, no arguments, no
      `--json` carve-out, baseline conformance); empty its **To be specified**
      list, moving the structured-form and sub-view items to **deferred**.
- [x] [`README.md`](../../README.md) — move `spec` from planned to built in the
      command status.

The implementation in `internal/cli` (registering `spec`, embedding the
specification, and tests) is this change's implementation, tracked by the
functional spec and design doc rather than listed here. No change to
[`specs/cli.md`](../../specs/cli.md): it already lists `spec` and names it the
carve-out.

## Status

`Done`. Implemented and archived after settling the durable `spec` sub-spec,
adding the design doc, registering `qualitymd spec`, embedding the format
specification, adding focused tests, and updating the README status.
