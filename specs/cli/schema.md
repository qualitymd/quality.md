---
type: Functional Specification
title: qualitymd schema
description: Emit the companion JSON Schema for QUALITY.md frontmatter.
tags: [cli, command, schema]
timestamp: 2026-06-22T00:00:00Z
---

# qualitymd schema

`qualitymd schema` emits the bundled
[companion JSON Schema](../quality-schema-json.md) for QUALITY.md frontmatter to
stdout, so an author or tool can vendor the schema — for editor validation and
autocomplete, or third-party validation — without reaching for an external copy.

`schema` inherits invocation-wide behavior, stdout/stderr separation,
determinism, plain-output rules, and exit-code categories from the
[CLI spec](../cli.md). This sub-spec states only the behavior particular to
`schema`. The artifact's content and guarantees are specified in
[`quality.schema.json`](../quality-schema-json.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Scope

Covered: the emitted artifact, its terminal rendering, and the command's argument
and flag surface.

Deferred:

- A JSON form of the specification _prose_ — a separate deferral in
  [`spec`](spec.md), distinct from this structural schema.
- `schema`-specific flags beyond cross-cutting CLI flags.

## Requirements

- `schema` **MUST** emit the bundled companion JSON Schema to stdout, sourced
  from the schema the binary was built with, so the command is self-contained and
  needs no file on disk or network access.
- When output must be plain — stdout is not a terminal, or `NO_COLOR` is set —
  `schema` **MUST** write the schema as verbatim JSON and nothing else.
  `qualitymd schema > quality.schema.json` reproduces the artifact
  byte-for-byte.
- When stdout is a terminal, `schema` **MAY** syntax-highlight the JSON and page
  it for readability per the [paging convention](../cli.md#conventions). This
  rendering is a human convenience; it **MUST NOT** change the bytes written when
  output must be plain.
- `schema` **MUST NOT** require any argument, and **MUST** treat an unexpected
  argument or flag as a usage error.
- `schema` **MUST NOT** offer `--json`. It is the verbatim-artifact carve-out
  named in the [`--json` convention](../cli.md#conventions); its output is
  already the JSON payload.
- On success `schema` **MUST** exit `0`. If it cannot emit the schema, it
  **MUST** exit `70`.
