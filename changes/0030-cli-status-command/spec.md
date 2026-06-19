---
type: Functional Specification
title: CLI status snapshot command — functional spec
description: Placeholder spec for the qualitymd status command; to be written when change 0030 is picked up in Draft.
status: Placeholder
tags: [cli, wizard]
timestamp: 2026-06-19T00:00:00Z
---

# CLI status snapshot command — functional spec

> **Placeholder.** This functional spec is not yet written. Change
> [0030](../0030-cli-status-command.md) is queued in `Draft` to capture
> motivation and scope. Before writing this spec, read
> [Writing functional specs](../../docs/guides/write-functional-specs.md) and
> [Designing CLI interfaces](../../docs/guides/cli-design.md).

## To be specified

- Command shape: `qualitymd status [path] [--json]`, read-only.
- Snapshot contents: model validity (lint result), model shape (target / factor
  / requirement counts and per-target source coverage), evaluation run history
  (run count, latest run, incomplete/stale runs), and open recommendation counts.
- Human output form (concise summary) vs. `--json` form (stable, agent-consumable
  schema); stdout/stderr contract per the output policy.
- Determinism: stable ordering and stable output for the same file state.
- Exit codes and behavior when `QUALITY.md` is absent or invalid.

## Durable spec changes

To be enumerated when the spec is written (expected: add `specs/cli/status.md`,
register in `specs/cli/index.md`; no `SPECIFICATION.md` format change).
