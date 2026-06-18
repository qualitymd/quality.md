---
type: Functional Specification
title: qualitymd evaluation add-record
description: Write schema-conformant evaluation records from judgment payloads.
tags: [cli, command, evaluation]
timestamp: 2026-06-18T00:00:00Z
---

# qualitymd evaluation add-record

`qualitymd evaluation add-record` writes one evaluation record into an existing
run folder. The record contract is [Evaluation records](../evaluation-records.md).

## Commands

```text
qualitymd evaluation add-record assessment <run>
qualitymd evaluation add-record analysis <run>
qualitymd evaluation add-record recommendation <run>
```

Each invocation writes exactly one record.

## Input

The command reads one JSON judgment payload from `--file <path>`, from
`--file -`, or from stdin when stdin is not a terminal. The payload **MUST NOT**
include CLI-owned fields such as `schemaVersion` or local record numbers.

## Requirements

The command **MUST** verify that `<run>` is an existing evaluation run folder,
strictly decode one JSON document, validate required fields, stamp
`schemaVersion: 1`, derive deterministic filenames, and write atomically.

Assessment and recommendation records use local `NNN` numbering. Analysis
records are keyed by target slug and may replace an existing analysis record for
that target.

On success, human output **MUST** report the written path on stderr. Under
`--json`, stdout **MUST** contain a receipt with `schemaVersion`, `path`, `kind`,
and, for analysis records, `created`.
