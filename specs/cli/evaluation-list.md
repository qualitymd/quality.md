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
as [`evaluation create`](evaluation-create.md), anchored by `--model <model>`
when supplied and otherwise by `QUALITY.md` in the current working directory. It
**MUST** list runs with a valid `EvaluationManifest` in deterministic
`EvaluationManifest.run.number` order and **MUST NOT** write or modify any run.

`--model <model>` **MUST** select the `QUALITY.md` file whose model-relative
workspace supplies the evaluation history. Listed run paths **MUST** be relative
to that selected model's workspace root.

Under `--json`, stdout **MUST** contain `schemaVersion` and `runs`. Each entry
**MUST** identify the run path, root area, evaluation data artifact count,
reportability, gap count, `requestedScope`, and `plannedScope`. Entries for
artifact-backed runs **MUST** carry the runner `lifecycle` status, so an
`awaiting_evaluator` run stays distinguishable from a generic incomplete
entry.

The command **MUST** use the current evaluation run inspection path. If a
recognized run has a broken evaluation skeleton, the command **MUST** fail
with the same diagnostic as `qualitymd evaluation status <run>`.

The command **MAY** accept `--state all|complete|reportable|incomplete|awaiting`
to filter the listed runs. `complete` and `reportable` are equivalent;
`awaiting` selects runs whose lifecycle is `awaiting_evaluator`. Unknown state
filters **MUST** be usage errors.
