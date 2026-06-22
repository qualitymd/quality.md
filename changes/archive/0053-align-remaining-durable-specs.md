---
type: Change Case
title: Align remaining durable specs
description: Split remaining large durable specs into parent and component/artifact contracts under the revised guidance.
status: Done
tags: [specs, cli]
timestamp: 2026-06-22T00:00:00Z
---

# Align remaining durable specs

This case applies the durable-spec granularity guidance from
[0052](0052-durable-spec-alignment.md) to the remaining internal
functional specs surfaced by audit: evaluation records, lint, and update notice
behavior.

- [Functional spec](0053-align-remaining-durable-specs/spec.md) - what the
  alignment must do.

## Motivation

After 0052 split the `/quality` skill spec, the same inventory test still found
older durable specs that mixed shared parent contracts with independently
reviewable component or artifact contracts. Leaving those mixed would make the
new guidance uneven: future contributors would have examples of both the new
shape and the old shape.

`SPECIFICATION.md` is intentionally out of scope: it is the primary public
format deliverable, not a normal internal functional spec for the system, and it
has a product reason to remain a single canonical artifact.

## Scope

In scope:

- split `specs/evaluation-records.md` into a shared parent plus child specs for
  run-folder and runtime record/artifact contracts;
- split `specs/cli/lint.md` so the command contract no longer owns the full
  lint rule system and output contract inline;
- split the cross-command ambient update notice out of `specs/cli/update.md`;
- update indexes, links, and logs needed to keep the OKF bundle navigable.

Deferred:

- changing CLI behavior, runtime record schemas, report formats, lint rules, or
  update behavior;
- splitting `SPECIFICATION.md`;
- creating new taxonomy beyond the parent/component/artifact distinction.

## Affected artifacts

### Code

None.

### Format spec

None.

### Durable specs

- `specs/evaluation-records.md` and `specs/evaluation-records/` - split shared
  evaluation-record invariants from child run-folder, record, artifact, and
  report-output contracts.
- `specs/cli/lint.md`, `specs/cli/lint-rules.md`, and
  `specs/cli/lint-output.md` - split the lint command contract from the rule
  system and output schema.
- `specs/cli/update.md` and `specs/cli/update-notice.md` - split the ambient
  update notice from the explicit update command.
- `specs/index.md`, `specs/cli/index.md`, and `specs/log.md` - update bundle
  navigation and history.

### Durable docs

None expected.

### Bundled skill

None.

### Install, scaffold, and packaging

None.

## Status

`Done`: landed and archived; no design doc required because the work was a
mechanical durable-spec split with no behavior change.
