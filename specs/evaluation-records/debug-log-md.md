---
type: Functional Specification
title: debug-log.md
description: Runtime debug-log.md contract for evaluation runs.
tags: [evaluation, records, cli, skill]
timestamp: 2026-06-22T00:00:00Z
---

# debug-log.md

This spec is part of the [Evaluation records](../evaluation-records.md) contract.
It owns the independently reviewable runtime contract below; shared
responsibility, runtime-not-OKF status, schema-version, and historical-record
compatibility rules live in the parent spec.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Debug Log

`debug-log.md` is a hand-authored process log for notable events involving the
evaluation process itself. The CLI **MUST** seed it when creating a new run, but
report loading and report rendering **MUST NOT** derive ratings, findings,
recommendations, next actions, or generated report content from it.

`debug-log.md` is a runtime evaluation artifact, not an OKF `log.md`, and
**MUST NOT** be interpreted as an OKF concept history.

The log's content **MUST** stay on the evaluator-process side of the boundary:
scope resolution, evaluation history inspection, coverage planning,
interruptions or resumes, retries, artifact corrections, CLI/tooling readiness
and failures, subagent coordination, redaction decisions, prompt-injection
handling, and report generation recovery.

The log **MUST NOT** record evaluation findings as primary content or
duplicate rating rationale from assessment or analysis records. When a project
command is exercised as evidence against the evaluated source, the log **MAY**
record that routing fact and cite the formal assessment record.

The log **MUST NOT** copy raw project-command output or preserve a second copy
of assessment evidence.

`debug-log.md` **MUST NOT** contain secret values or reproduce raw
prompt-injection text from evaluated source content. When a secret or hostile
instruction affects the evaluation process, the log **MAY** record a sanitized
locator and type.
