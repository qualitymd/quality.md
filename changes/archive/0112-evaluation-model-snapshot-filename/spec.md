---
type: Functional Specification
title: Evaluation model snapshot filename - functional spec
description: Requirements for renaming the Evaluation v2 run-folder model snapshot from model.md to model-snapshot.md across the CLI write, every read, and the durable contracts that name it.
tags: [cli, evaluation, reports]
timestamp: 2026-06-26T00:00:00Z
---

# Evaluation model snapshot filename - functional spec

Companion to the
[Evaluation model snapshot filename](../0112-evaluation-model-snapshot-filename.md)
change case. This spec states what the rename must do.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119 and RFC 8174 when, and only when, they
appear in all capitals.

## Background / Motivation

On `evaluation create`, the CLI writes a frozen copy of the resolved
working-tree model into the run folder. That copy is the only run artifact with
a living counterpart — the mutable working-tree `QUALITY.md` model — and the
bare filename `model.md` does not distinguish the frozen snapshot from its
editable source. The directory's purpose is to be a point-in-time receipt, so
the file that mirrors a mutable upstream should name itself as a snapshot. The
durable specs already describe it as "the model snapshot"; `model-snapshot.md`
makes the filename match.

## Scope

Covered: the Evaluation v2 run-folder snapshot filename the CLI writes, every
CLI read of that snapshot and the error messages that name it, the seed-layout
test, the durable contracts that name the snapshot file, and this repository's
own tracked active dogfood runs.

Not covered: compatibility reading of the old filename, runtime migration of
runs in other repositories, frozen archived runs the CLI does not enumerate,
run-folder names, the `data/` tree, report filenames, report content, and the
meaning of the snapshot or the staleness computation.

## Requirements

### Snapshot filename contract

- On `evaluation create`, the CLI **MUST** write the run-folder model snapshot to
  `model-snapshot.md` at the run root.

- The CLI **MUST NOT** write the run-folder model snapshot to `model.md`.

- The CLI **MUST** read the run-folder model snapshot from `model-snapshot.md`
  everywhere it consumes the snapshot: run-folder validation, model load for
  report rating-label resolution, the Evaluation v2 report builder, and `status`
  staleness comparison.

  > Rationale: the snapshot is the only run artifact with a mutable upstream, so
  > its filename should state that it is a frozen capture rather than the live
  > model. The write and every read move together because a split name would make
  > a freshly created run unreadable. - 0112

- The CLI **MUST NOT** retain a compatibility path that reads, accepts, or
  validates the old `model.md` snapshot filename.

  > Rationale: Evaluation v2 took a clean break from previous run shapes; a
  > dual-name reader would carry the old name forward indefinitely and weaken the
  > deterministic run layout. - 0112

### Operator-facing messages

- CLI error messages that name the missing or unreadable run-folder snapshot
  **MUST** name `model-snapshot.md`, not `model.md`.

### Boundaries

- This change **MUST NOT** migrate, rewrite, or re-validate evaluation runs in
  other repositories.

- This change **MUST NOT** alter frozen archived runs under
  `.quality/evaluations/archive/`, which the CLI does not enumerate.

- This change **MUST NOT** change run-folder names, the structured `data/` tree,
  generated report filenames, report content, or the staleness computation
  beyond the filename the comparison reads.

- This repository's own tracked active (non-archived) dogfood runs **MUST**
  carry the renamed snapshot file and **MUST** keep validating under the new
  filename contract.

## Acceptance Criteria

- Creating an Evaluation v2 run produces a run-folder snapshot named
  `model-snapshot.md` and no run-folder file named `model.md`.
- Run-folder validation recognizes a run by `model-snapshot.md` and rejects a
  folder missing it, with an error message naming `model-snapshot.md`.
- The model load path, the Evaluation v2 report builder, and `status` staleness
  comparison all read `model-snapshot.md`.
- `status` reports a created run as non-stale immediately after creation and as
  stale once the selected model file bytes differ from the snapshot.
- No source file, durable spec, bundled skill, or tracked active dogfood run
  reads or writes the old `model.md` snapshot filename.
- The seed-layout test asserts `model-snapshot.md` is created and `model.md` is
  not.
- `mise run check` passes.

## Durable spec changes

### To add

None.

### To modify

- [`specs/cli/evaluation-create.md`](../../../specs/cli/evaluation-create.md) -
  name the seeded run-folder snapshot `model-snapshot.md`. Driven by
  [Snapshot filename contract](#snapshot-filename-contract).
- [`specs/cli/status.md`](../../../specs/cli/status.md) - name the
  staleness-comparison snapshot `model-snapshot.md`. Driven by
  [Snapshot filename contract](#snapshot-filename-contract).
- [`specs/evaluation-v2/reports/report-tree.md`](../../../specs/evaluation-v2/reports/report-tree.md)
  - name the rating-label source snapshot `model-snapshot.md`. Driven by
    [Snapshot filename contract](#snapshot-filename-contract).
- [`specs/skills/quality-skill/evaluation.md`](../../../specs/skills/quality-skill/evaluation.md)
  - name the snapshot the create step writes `model-snapshot.md`. Driven by
    [Snapshot filename contract](#snapshot-filename-contract).
- [`specs/skills/quality-skill/reporting.md`](../../../specs/skills/quality-skill/reporting.md)
  - name the run-folder snapshot the report contract requires
    `model-snapshot.md`. Driven by
    [Snapshot filename contract](#snapshot-filename-contract).

### To rename

None. (The renamed artifact is a generated run-folder file, not a durable spec
file.)

### To delete

None.

## Validation check

If every requirement is satisfied, a freshly created Evaluation v2 run names its
model snapshot `model-snapshot.md`, every CLI read and operator message uses that
name, no old-name compatibility path lingers, and this repository's own tracked
runs keep validating — so the snapshot filename states that it is a frozen
capture without widening into run migration, report-layout changes, or staleness
semantics.
