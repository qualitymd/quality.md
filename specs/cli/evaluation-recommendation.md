---
type: Functional Specification
title: qualitymd evaluation recommendation
description: Add and list recommendation records.
tags: [cli, command, evaluation]
timestamp: 2026-06-19T00:00:00Z
---

# qualitymd evaluation recommendation

`qualitymd evaluation recommendation` is the recommendation-record resource. The
record contract is
[Recommendation record](../evaluation-records/recommendation-record-md.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Commands

```text
qualitymd evaluation recommendation add <run>
qualitymd evaluation recommendation list <run>
```

Both verbs **MUST** accept either a positional run path or `--latest`, and
**MUST** error when both or neither are supplied.

`add` reads recommendation JSON from `--file <path>`, from `--file -`, or from
non-terminal stdin. It **MUST** accept either one recommendation object or an
array of recommendation objects, validate each payload, render deterministic
runtime Markdown with `schemaVersion: 1`, append numbered recommendation records
in order, and report all written paths. The payload **MUST NOT** include
CLI-owned fields such as `schemaVersion`.

`list` **MUST NOT** accept `--file`, **MUST NOT** write files, and **MUST** list
the run's recommendation record paths. Under `--json`, stdout **MUST** include
`schemaVersion`, `path`, `kind`, and `records`.
