---
type: Command Specification
title: qualitymd lint
description: Validate a QUALITY.md file's structure against the format spec.
tags: [cli, command, lint]
timestamp: 2026-06-16T00:00:00Z
---

# qualitymd lint

**Version 0.1 — Draft · Placeholder**

> 🚧 **Placeholder.** This sub-spec is a stub. The cross-cutting CLI contract
> that `lint` inherits — invocation, global flags, output formats, exit codes,
> and agent accessibility — is specified in the [CLI spec](../cli.md). This file
> will specify only what is particular to `lint`.

`qualitymd lint` validates a `QUALITY.md` file's structure against the
[format specification](../../SPECIFICATION.md), fast and deterministically,
exiting non-zero on errors so it drops into CI.

**Boundary.** `lint` checks *format conformance* — whether the file is a valid
`QUALITY.md` — only. It does not assess whether the model is a *good* quality
model; that judgment lives in the evaluation skills, not the deterministic CLI.

## To be specified

- The full rule set: frontmatter parses; `ratingScale` present with ≥2 ordered,
  uniquely named levels each carrying a `criterion`; at least one of
  `factors`/`requirements`/`targets`; every requirement has exactly one non-empty
  scalar `assessment`; secondary factor names resolve in scope; `ratings`
  override keys name real scale levels; `title`/`ratingScale` only on the root.
- The finding schema (rule id, severity, message, location) for `--format json`.
- Severity levels and which map to a non-zero [exit code](../cli.md#exit-codes).
- `lint`-specific flags beyond the [global flags](../cli.md#global-flags).
