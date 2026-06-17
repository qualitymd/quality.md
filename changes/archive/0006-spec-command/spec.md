---
type: Functional Specification
title: spec command — functional spec
description: What qualitymd spec must do — emit the bundled QUALITY.md format specification to stdout.
tags: [cli, command, spec]
timestamp: 2026-06-17T00:00:00Z
---

# spec command — functional spec

Companion to the [Specify and implement the spec command](../0006-spec-command.md)
change. This spec states *what* `qualitymd spec` must do; the
[design doc](design.md) covers *how* the specification is embedded and emitted.

`spec` inherits the cross-cutting contract — invocation, exit codes, the
baseline, and the `--json` convention — from the [CLI spec](../../../specs/cli.md);
this document states only what is particular to `spec`, and the settled sub-spec
it produces replaces the placeholder at
[`specs/cli/spec.md`](../../../specs/cli/spec.md).

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Scope

Covered: the output (its verbatim form and a human-readable rendering), the
argument and flag surface, and how `spec` meets the CLI
[baseline](../../../specs/cli.md#baseline). Intentionally **deferred** — see
[Deferred](#deferred) — are a structured form of the specification, sub-views,
and any `spec`-specific flags.

## Requirements

- `spec` **MUST** emit the bundled `QUALITY.md`
  [format specification](../../../SPECIFICATION.md) to stdout, sourced from the
  specification the binary was built with, so the command is self-contained and
  needs no file on disk or network access.
- When its output must be plain — stdout is not a terminal, or `NO_COLOR` is set,
  per the [baseline](../../../specs/cli.md#baseline) — `spec` **MUST** write the
  specification as **verbatim Markdown** and nothing else, so
  `qualitymd spec > SPECIFICATION.md` reproduces the artifact byte-for-byte. This
  is the payload, and it is deterministic: the same binary emits the same bytes
  every invocation. Any diagnostic goes to stderr.
- When stdout is a terminal, `spec` **SHOULD** render the Markdown formatted for
  readability, using the stack's terminal renderer. This rendering is a human
  convenience, not the payload: it is the same content as the verbatim form and
  **MUST NOT** change the bytes written when output must be plain. No flag selects
  it — the rendered-on-TTY / verbatim-when-plain split rides the baseline's
  existing terminal-detection rule, exactly as color does.
- `spec` **MUST NOT** require any argument, and **MUST** treat an unexpected
  argument or flag as a usage error (exit `2`).
- `spec` **MUST NOT** offer `--json`. It is the verbatim-artifact carve-out named
  in the [`--json` convention](../../../specs/cli.md#conventions): its output *is*
  the payload, so wrapping it adds nothing.
- On success `spec` **MUST** exit `0`; if it cannot emit the specification (an
  I/O failure writing stdout) it **MUST** exit `70`, per the
  [exit-code categories](../../../specs/cli.md#exit-codes).

## Deferred

Recorded so their absence reads as deliberate, not forgotten:

- A structured/JSON form of the specification. `spec` is the carve-out; should a
  machine-readable view of the format ever be needed, it is its own change.
- Sub-views — e.g. emitting only the `lint` rule set — and the flags that would
  select them.
- `spec`-specific flags beyond the cross-cutting CLI flags.
- `nextActions`. `spec` emits an artifact meant for redirection and has no
  obvious deterministic follow-up; it offers none until one is warranted.
