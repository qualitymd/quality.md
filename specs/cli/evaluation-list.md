---
type: Functional Specification
title: qualitymd evaluation list
description: List evaluation runs.
tags: [cli, command, evaluation]
timestamp: 2026-06-19T00:00:00Z
---

# qualitymd evaluation list

`qualitymd evaluation list` enumerates recognized evaluation runs in the
resolved evaluation directory.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Requirements

The command **MUST** resolve the evaluation directory using the same precedence
as [`evaluation create`](evaluation-create.md). It **MUST** list recognized run
folders in deterministic run-number order and **MUST NOT** write or modify any
run.

Under `--json`, stdout **MUST** contain `schemaVersion` and `runs`. Each entry
**MUST** identify the run path, subject, record counts, and reportability, and
**MAY** include narrowing when present. It **MAY** include a gap count so callers
can distinguish an incomplete run with blocking diagnostics from an empty run.

The command **MUST NOT** abort solely because an individual record in a
recognized run is malformed, schema-incompatible, or structurally incomplete
under the current record contract. Such a run remains listed as incomplete or
problematic; detailed diagnostics live in `qualitymd evaluation status <run>`.

The command **MAY** accept `--state all|reportable|incomplete` to filter the
listed runs. Unknown state filters **MUST** be usage errors.
