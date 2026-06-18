---
type: Functional Specification
title: qualitymd evaluation set-planned-coverage
description: Write planned coverage metadata for an evaluation run.
tags: [cli, command, evaluation]
timestamp: 2026-06-18T00:00:00Z
---

# qualitymd evaluation set-planned-coverage

`qualitymd evaluation set-planned-coverage <run>` writes optional planned
coverage metadata for an existing evaluation run. The artifact contract is
defined by [Evaluation records](../evaluation-records.md#planned-coverage).

## Input

The command reads one JSON payload from `--file <path>`, from `--file -`, or
from stdin when stdin is not a terminal.

The payload **MUST** contain `schemaVersion: 1`, `assessments`, and `analyses`.
Assessment entries **MUST** contain ordered `targetPath` and `requirement`.
Analysis entries **MUST** contain ordered `targetPath`.

## Requirements

The command **MUST** verify that `<run>` is an existing evaluation run folder,
strictly decode one JSON document, reject unknown fields, reject duplicate
planned assessment keys, reject duplicate planned analysis keys, and write a
canonical `planned-coverage.json` artifact at the run root.

Planned assessment keys are ordered `targetPath` plus `requirement`. Planned
analysis keys are ordered `targetPath`.

The command **MUST** replace any existing `planned-coverage.json` for the run.

On success, human output **MUST** report the written path on stderr. Under
`--json`, stdout **MUST** contain a receipt with `schemaVersion`, `path`, and
`nextActions`.
