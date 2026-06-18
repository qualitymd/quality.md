---
type: Change
title: Create-run subject validation
description: Validate create-run subject paths before creating run folders.
status: In-Review
tags: [evaluation, cli, usability]
timestamp: 2026-06-18T00:00:00Z
---

# Create-run subject validation

A unit of work that turns the experiment program's `create-run --subject .`
finding into safer CLI behavior. Detail lives in the child:

- [Functional spec](0022-create-run-subject-validation/spec.md) - what subject
  validation must do.
- [Design doc](0022-create-run-subject-validation/design.md) - how creation is
  ordered so invalid subjects leave no partial run.

## Motivation

The E14 improve/re-evaluate experiment found that `qualitymd evaluation
create-run --subject .` failed after creating an empty run skeleton. The command
should validate the subject path before making durable run artifacts, so a bad
flag does not consume a run number or leave cleanup work.

## Scope

Covered: validate `--subject` as a repository-relative file path before creating
the evaluation directory or run folder, and keep model-altitude source
validation consistent with subject-altitude snapshots.

Deferred: accepting directory subjects, changing evaluation directory path
rules, and adding rollback for failures that happen after creation starts.

## Affected specs & docs

Created or updated during `In-Progress`, before this change reaches
`In-Review`:

- [x] [`specs/cli/evaluation-create-run.md`](../specs/cli/evaluation-create-run.md)
      - document subject path validation and no-partial-run behavior for invalid
      subjects.

## Status

`In-Review`. See the [status lifecycle](index.md#status-lifecycle).
