---
type: Functional Specification
title: qualitymd spec
description: Emit the QUALITY.md format specification.
tags: [cli, command, spec]
timestamp: 2026-06-17T00:00:00Z
---

# qualitymd spec

`qualitymd spec` emits the bundled
[QUALITY.md format specification](../../SPECIFICATION.md) to stdout, so an author or coding
agent can inject the current format rules into its context without reaching for
an external copy.

`spec` inherits invocation-wide behavior, stdout/stderr separation, determinism,
plain-output rules, and exit-code categories from the [CLI spec](../cli.md). This
sub-spec states only the behavior particular to `spec`.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Scope

Covered: the emitted artifact, its terminal rendering, and the command's argument
and flag surface.

Deferred:

- A structured or JSON form of the specification.
- Sub-views, such as emitting only one command's rule set.
- `spec`-specific flags beyond cross-cutting CLI flags.

## Requirements

- `spec` **MUST** emit the bundled format specification to stdout, sourced from
  the specification the binary was built with, so the command is self-contained
  and needs no file on disk or network access.
- When output must be plain — stdout is not a terminal, or `NO_COLOR` is set —
  `spec` **MUST** write the specification as verbatim Markdown and nothing else.
  `qualitymd spec > SPECIFICATION.md` reproduces the artifact byte-for-byte.
- When stdout is a terminal, `spec` should render the Markdown formatted for
  readability using the stack's terminal renderer. This rendering is a human
  convenience; it **MUST NOT** change the bytes written when output must be plain.
- When stdout is a terminal, `spec` can display the rendered specification
  through the user's pager per the [paging convention](../cli.md#conventions), so
  a long spec scrolls. Paging is skipped when stdout is not a terminal, so a
  redirect still reproduces the artifact byte-for-byte.
- `spec` **MUST NOT** require any argument, and **MUST** treat an unexpected
  argument or flag as a usage error.
- `spec` **MUST NOT** offer `--json`. It is the verbatim-artifact carve-out named
  in the [`--json` convention](../cli.md#conventions); its output is the payload.
- On success `spec` **MUST** exit `0`. If it cannot emit the specification, it
  **MUST** exit `70`.
