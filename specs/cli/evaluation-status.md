---
type: Functional Specification
title: qualitymd evaluation status
description: Inspect whether an evaluation run can be rendered into reports.
tags: [cli, command, evaluation]
timestamp: 2026-06-19T00:00:00Z
---

# qualitymd evaluation status

`qualitymd evaluation status <run>` reads an evaluation run and reports whether
it is complete enough for `qualitymd evaluation report build`.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Requirements

The command **MUST** accept either a positional run path or `--latest`, and
**MUST** error when both or neither are supplied. `--latest` resolves to the most
recent recognized run in the resolved evaluation directory for the selected
model.

`--model <model>` **MUST** select the `QUALITY.md` file whose model-relative
workspace supplies `--latest` history. When `--model` is supplied with a
relative positional `<run>` path, the command **MUST** resolve that path relative
to the selected model's workspace root. When `--model` is absent and a
positional `<run>` path is supplied, the command **MAY** preserve ordinary
filesystem-path behavior.

The command **MUST NOT** write files. It exits `0` when a current Evaluation
run can be inspected, even when it is not yet reportable. Missing, malformed,
unreadable, schema-incompatible, or structurally incomplete Evaluation payloads under
`data/` **MUST** produce typed gaps and make the run non-reportable.

A missing run folder or broken evaluation run skeleton fails as a command error.

Human output **MUST** include the run path, evaluation data artifact count,
reportability, and any gaps. Under `--json`, stdout **MUST** include
`schemaVersion`, `path`, `reportable`, `data`, gaps, and `nextActions`.

Every gap `kind` **MUST** be one of the typed evaluation-run gap kinds defined by
the implementation and the active evaluation contract. For Evaluation runs,
typed gaps come from the [Evaluation](../evaluation/evaluation.md)
record and reportability rules. Status routing **MUST** use those typed gap kinds
rather than interpreting free-form text in `detail`.

When a run contains Evaluation data, `status` **MUST** inspect the required
structured payload graph under `data/` and report missing, malformed, unreadable,
schema-incompatible, or structurally incomplete payloads as typed gaps.

Status **MUST NOT** expose active recommendation counts, planned coverage gaps,
or compatibility transforms.
