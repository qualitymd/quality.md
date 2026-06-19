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

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Commands

```text
qualitymd evaluation add-record assessment <run>
qualitymd evaluation add-record analysis <run>
qualitymd evaluation add-record recommendation <run>
```

Each invocation writes exactly one record. The record kind is a subcommand rather
than a `--kind` flag because each kind has a different required-field set and
lands in a different subdirectory, so the structural split keeps each contract
distinct and discoverable.

## Input

The command reads one JSON judgment payload from `--file <path>`, from
`--file -`, or from stdin when stdin is not a terminal. When `--file` is absent
and stdin is a terminal, the command **MUST** fail with a usage error that tells
the caller to pass `--file <path>` or pipe JSON on stdin.

The command **MUST** strictly decode a single well-formed JSON document. Empty
input, trailing garbage, or multiple documents **MUST** be a usage error. The
payload **MUST NOT** include CLI-owned fields such as `schemaVersion` or a local
record number; a payload that supplies one **MUST** be rejected as a usage error,
not silently honored or stripped, so the CLI-writes / skill-judges division
cannot be bypassed through the payload.

## Requirements

The command **MUST** verify that `<run>` is an existing evaluation run folder,
strictly decode one JSON document, validate required fields, stamp
`schemaVersion: 1`, derive deterministic filenames, and write atomically. The
command never scaffolds a run; a missing or malformed `<run>` is an internal
error.

Validation is contract-specific. The command **MUST** reject, at minimum: an
assessment with `notAssessed` true and a non-null `rating`; a `rating` not in
the run's rating scale; a finding missing `locator`, `observation`, or
`category`; and a recommendation missing its required structured fields. An
assessment `rating` **MUST** be validated against the rating scale in the run's
`model.md` snapshot, not the live working-directory `QUALITY.md`, so a model edit
mid-run cannot silently change which ratings are valid and the run stays
internally consistent; an unparseable `model.md` is an internal error.

Decode and payload-validation failures **MUST** map to a usage error (exit `2`);
a missing or malformed run target and I/O failures **MUST** map to an internal
error (exit `70`). Rejection and failed writes are atomic: they **MUST** consume
no `NNN` and leave no partial file.

Assessment and recommendation records use local `NNN` numbering, computed by a
stateless directory scan (one past the highest present, no counter file). On a
detected `NNN` collision from a concurrent writer, the command **MUST** recompute
the next number once and retry, then fail as an internal error naming the
contended directory. Analysis records are keyed by target slug and may replace an
existing analysis record for that target.

Assessment payloads **MAY** include `supersedes`, a list of assessment IDs or
paths that the new record corrects for the same ordered `targetPath` and
`requirement`. When present, the command **MUST** write the list into the
assessment record.

Recommendation payloads **MAY** include `supersedes`, a list of recommendation
IDs or paths that the new record corrects or replaces as active advice. When
present, the command **MUST** write the list into recommendation metadata and
the rendered recommendation body.

The CLI renders the recommendation body and YAML frontmatter deterministically
from the structured payload fields in a fixed key order — the body is generated,
not passed through from the payload — so identical judgment input produces
byte-identical records and the human layout stays a CLI concern.

On success, human output **MUST** report the written path on stderr. Under
`--json`, stdout **MUST** contain a receipt with `schemaVersion`, `path`, `kind`,
and, for analysis records, `created`.
