---
type: Command Specification
title: qualitymd spec
description: Emit the QUALITY.md format specification.
tags: [cli, command, spec]
timestamp: 2026-06-16T00:00:00Z
---

# qualitymd spec

**Version 0.1 — Draft · Placeholder**

> 🚧 **Placeholder.** This sub-spec is a stub. The cross-cutting CLI contract
> that `spec` inherits — invocation, global flags, output formats, exit codes,
> and agent accessibility — is specified in the [CLI spec](../cli.md). This file
> will specify only what is particular to `spec`.

`qualitymd spec` emits the bundled `QUALITY.md`
[format specification](../../SPECIFICATION.md) to stdout, so an author or — more
often — a coding agent can inject the current format rules into its context
without reaching for an external copy.

## To be specified

- The default output: the format spec as Markdown.
- Whether and how to emit a structured form under `--format json`.
- Any sub-views (e.g. emitting the `lint` rule set), and the flags that select
  them.
- `spec`-specific flags beyond the [global flags](../cli.md#global-flags).
