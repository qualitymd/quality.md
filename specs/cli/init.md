---
type: Functional Specification
title: qualitymd init
description: Scaffold a starter QUALITY.md to fill in.
tags: [cli, command, init]
timestamp: 2026-06-16T00:00:00Z
---

# qualitymd init

> 🚧 **Placeholder.** This sub-spec is a stub. The cross-cutting CLI contract
> that `init` inherits — invocation, global flags, output formats, exit codes,
> and agent accessibility — will be specified in the [CLI spec](../cli.md). This
> file will specify only what is particular to `init`.

`qualitymd init` scaffolds a starter `QUALITY.md` for an author to fill in.

## To be specified

- The scaffold contents: a seeded [rating scale](../../SPECIFICATION.md#model), a
  minimal target → factor → requirement skeleton, and the recommended Markdown
  body sections (Overview, Scope, Needs, Risks, Known gaps) as headed stubs.
- Target path and the rule for refusing to overwrite an existing file (and the
  explicit flag that permits it).
- Emitting to stdout for piping (`-`).
- `init`-specific flags beyond the cross-cutting flags in the
  [CLI spec](../cli.md).
