---
type: Design Doc
title: Create-run subject validation design
description: How create-run validates the subject before creating durable run files.
tags: [evaluation, cli, usability]
timestamp: 2026-06-18T00:00:00Z
---

# Create-run subject validation design

## Context

The [Create-run subject validation spec](spec.md) answers the E14/E29 UX gap:
invalid `--subject` values can fail after the run folder and subdirectories are
already created.

## Approach

Move model snapshot preparation ahead of evaluation directory creation. A helper
resolves the subject path, rejects obvious directory values such as `.`, checks
that the resolved path is a file, and returns the exact `model.md` bytes to
write later.

`CreateRun` then creates the evaluation directory, computes the run number, and
writes the already prepared model snapshot. If subject validation fails, no
evaluation directory or numbered run folder has been created yet.

## Alternatives

**Rollback partial folders on failure.** Rejected for this change. Prevalidation
is simpler and avoids deleting anything that might have been concurrently
created.

**Accept a directory subject and search for `QUALITY.md`.** Rejected. The flag
is documented as a path to the model file to snapshot; accepting directories
would add implicit lookup behavior.

## Trade-offs and risks

The model-altitude path becomes slightly stricter because the source path is
validated as a file before embedding it in the meta-model. That matches the
documented meaning of `--subject`.
