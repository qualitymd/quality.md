---
type: Functional Specification
title: qualitymd evaluation assessment-result
description: Add and list assessment result records.
tags: [cli, command, evaluation]
timestamp: 2026-06-19T00:00:00Z
---

# qualitymd evaluation assessment-result

`qualitymd evaluation assessment-result` is the assessment result record
resource. The record contract is [Evaluation records](../evaluation-records.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Commands

```text
qualitymd evaluation assessment-result add <run>
qualitymd evaluation assessment-result list <run>
```

Both verbs **MUST** accept either a positional run path or `--latest`, and
**MUST** error when both or neither are supplied.

`add` reads assessment result JSON from `--file <path>`, from `--file -`, or
from non-terminal stdin. It **MUST** accept either one assessment result object
or an array of assessment result objects, validate each payload against the
run's `model.md`, stamp `schemaVersion: 1`, append numbered assessment result
records in order, and report all written paths. The payload **MUST NOT** include
CLI-owned fields such as `schemaVersion`.

`list` **MUST NOT** accept `--file`, **MUST NOT** write files, and **MUST** list
the run's assessment result record paths. Under `--json`, stdout **MUST** include
`schemaVersion`, `path`, `kind`, and `records`.
