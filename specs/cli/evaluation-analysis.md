---
type: Functional Specification
title: qualitymd evaluation analysis
description: Set and list analysis records.
tags: [cli, command, evaluation]
timestamp: 2026-06-19T00:00:00Z
---

# qualitymd evaluation analysis

`qualitymd evaluation analysis` is the analysis-record resource. The record
contract is [Evaluation records](../evaluation-records.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Commands

```text
qualitymd evaluation analysis set <run>
qualitymd evaluation analysis list <run>
```

Both verbs **MUST** accept either a positional run path or `--latest`, and
**MUST** error when both or neither are supplied.

`set` reads analysis JSON from `--file <path>`, from `--file -`, or from
non-terminal stdin. It **MUST** accept either one analysis object or an array of
analysis objects, validate each payload against the run's `model.md`, stamp
`schemaVersion: 1`, and upsert each analysis by area slug in order. It
**MUST** report all written paths. The payload **MUST NOT** include CLI-owned
fields such as `schemaVersion`.

`list` **MUST NOT** accept `--file`, **MUST NOT** write files, and **MUST** list
the run's analysis record paths. Under `--json`, stdout **MUST** include
`schemaVersion`, `path`, `kind`, and `records`.
