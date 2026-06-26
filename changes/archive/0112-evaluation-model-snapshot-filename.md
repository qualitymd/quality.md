---
type: Change Case
title: Evaluation model snapshot filename
description: Rename the Evaluation v2 run-directory model snapshot from model.md to model-snapshot.md so the filename signals it is a frozen point-in-time copy of the working-tree model, not the live model.
status: Done
tags: [cli, evaluation, reports, naming]
timestamp: 2026-06-26T00:00:00Z
---

# Evaluation model snapshot filename

This parent concept captures the why and status; the detail lives in its
children:

- [Functional spec](0112-evaluation-model-snapshot-filename/spec.md) - what the
  change must do.
- [Design doc](0112-evaluation-model-snapshot-filename/design.md) - how it is
  built, and why.

## Motivation

When the CLI creates an Evaluation v2 run, it writes a frozen copy of the
resolved working-tree model into the run folder as `model.md`. Unlike the other
run artifacts, this file has a living counterpart — the working-tree
`QUALITY.md` model it was copied from — and the bare name `model.md` does not
say which one a reader is looking at. In a run directory whose whole purpose is
to be a point-in-time evaluation receipt, the one file that mirrors a mutable
upstream should name itself as a snapshot.

The durable specs already call it "the run's `model.md` snapshot" and "the model
snapshot" in prose; the filename should stop fighting that. `model-snapshot.md`
carries the one bit a reader needs — that this is a frozen capture, not the
canonical editable model — and reads as run-scoped alongside the other receipts.

## Scope

Covered:

- the Evaluation v2 run-folder snapshot filename the CLI writes on
  `evaluation create`;
- every CLI read of that snapshot — run-folder validation, model load for report
  rating-label resolution, the v2 report builder, and `status` staleness
  comparison — and the error messages that name the file;
- the seed-layout test that asserts the snapshot file is created;
- durable CLI, Evaluation v2 report, and `/quality` skill contracts that name
  the snapshot file, and the bundled skill workflow guidance; and
- this repository's own two tracked active dogfood runs, whose snapshot files and
  run-local prose name `model.md`.

Deferred / non-goals:

- no compatibility path that reads or accepts the old `model.md` filename — this
  is a clean break, consistent with the Evaluation v2 clean-break stance;
- no runtime migration of evaluation runs in other repositories;
- no change to frozen `.quality/evaluations/archive/` runs, which are not
  enumerated by the CLI and record past state;
- no change to run-folder names, the `data/` tree, report filenames, report
  content, or any other run artifact; and
- no change to the meaning of the snapshot or how staleness is computed.

## Affected artifacts

### Code

- [x] `internal/evaluation/create.go` - write the snapshot to
      `model-snapshot.md` and update the write error message; introduce the
      shared snapshot-filename constant.
- [x] `internal/evaluation/path.go` - validate run folders by the new snapshot
      filename and update the "missing" error message.
- [x] `internal/evaluation/load.go` - parse the snapshot from the new filename.
- [x] `internal/evaluation/report_v2.go` - read the snapshot from the new
      filename for rating-label resolution.
- [x] `internal/status/status.go` - read the snapshot from the new filename for
      staleness comparison and update the read error message.
- [x] `internal/evaluation/evaluation_test.go` - assert the seeded snapshot file
      uses the new filename; review other affected CLI tests.

### Dogfood evaluation runs (this repo)

- [x] `.quality/evaluations/0005-subject-quality-eval/model.md` and
      `.quality/evaluations/0006-quality-eval/model.md` - rename the tracked
      active-run snapshot files so this repo's own runs keep validating.
- [x] `.quality/evaluations/0006-quality-eval/design.md` and `plan.md` - update
      run-local prose that names `model.md`.

### Format spec

- [x] None - `SPECIFICATION.md` does not govern the Evaluation v2 run-folder
      snapshot filename. (Deliberate.)

### Durable specs

- [x] `specs/cli/evaluation-create.md` - name the seeded snapshot
      `model-snapshot.md`.
- [x] `specs/cli/status.md` - name the staleness-comparison snapshot
      `model-snapshot.md`.
- [x] `specs/evaluation-v2/reports/report-tree.md` - name the rating-label
      source snapshot `model-snapshot.md`.
- [x] `specs/skills/quality-skill/evaluation.md` - name the snapshot the create
      step writes `model-snapshot.md`.
- [x] `specs/skills/quality-skill/reporting.md` - name the run-folder snapshot
      the report contract requires `model-snapshot.md`.

### Durable docs / bundled skill

- [x] `skills/quality/workflows/evaluate.md` - name the snapshot the create step
      writes `model-snapshot.md`.

### Suggested new durable specs

- None.

## Status

`Done`. Implemented, verified, and archived. The snapshot is read/written through
one `evaluation.ModelSnapshotFile` constant with no old-name compatibility path;
`mise run check` passes. See the
[status lifecycle](../index.md#status-lifecycle).
