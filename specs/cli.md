---
type: Functional Specification
title: qualitymd CLI
description: High-level requirements for the deterministic qualitymd command-line surface.
tags: [cli, specification]
timestamp: 2026-06-16T00:00:00Z
---

# qualitymd CLI

> 🚧 **Placeholder.** This spec is a stub. It will specify the cross-cutting
> contract every `qualitymd` command shares — invocation, flags, output, exit
> codes, and agent accessibility — independent of any one command's behavior.
> Per-command behavior lives in the command sub-specs under [Commands](#commands).

This document will specify the high-level requirements for the `qualitymd`
command-line interface. It is a companion to the
[`QUALITY.md` format specification](../SPECIFICATION.md): where this document
constrains the *tool*, the format spec constrains the *file*.

## Scope

This phase covers three commands — the **format-tooling layer**: commands that
operate on a single `QUALITY.md` file, hold no evaluation state, and never call a
model.

| Command                 | Purpose                                              |
| ----------------------- | ---------------------------------------------------- |
| [`init`](./cli/init.md) | Scaffold a starter `QUALITY.md` to fill in.          |
| [`lint`](./cli/lint.md) | Validate a file's structure against the format spec. |
| [`spec`](./cli/spec.md) | Emit the `QUALITY.md` format specification.          |

**Deferred.** The deeper surface that records per-target verdicts, rolls them up
the target tree, and gates CI on the outcome is out of scope for this phase. The
judgment-based evaluation defined in the format spec's
[Evaluation](../SPECIFICATION.md#evaluation) is carried by skills that orchestrate
the model, not by these commands.

## Commands

- [`init`](./cli/init.md) — scaffold a starter `QUALITY.md`.
- [`lint`](./cli/lint.md) — validate a file's structure.
- [`spec`](./cli/spec.md) — emit the format specification.

## To be specified

- The shared invocation form and the file / stdin argument convention.
- Global flags common to every command.
- Output formats (human and machine-readable) and their stability.
- Exit-code semantics.
- Agent-accessibility and CI requirements.
- Versioning of the binary and the format spec it targets.
