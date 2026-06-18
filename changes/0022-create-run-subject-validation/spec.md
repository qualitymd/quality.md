---
type: Functional Specification
title: Create-run subject validation
description: Validate create-run subject paths before creating run folders.
tags: [evaluation, cli, usability]
timestamp: 2026-06-18T00:00:00Z
---

# Create-run subject validation

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Scope

This change covers validation of `qualitymd evaluation create-run --subject`
before run-folder creation.

It does not change evaluation directory validation, run numbering, run-folder
names, or support directory subjects.

## Requirements

`qualitymd evaluation create-run` **MUST** validate `--subject` before creating
the evaluation directory or run folder.

The subject path **MUST** be repository-relative, must not escape the repository,
and **MUST** resolve to a file, not a directory.

When the subject path is invalid, `create-run` **MUST** fail without creating a
numbered run folder.

For subject altitude, the validated subject file **MUST** be snapshotted into
`model.md`.

For model altitude, the validated subject file path **MUST** be used as the
source path in the bundled quality meta-model snapshot.
