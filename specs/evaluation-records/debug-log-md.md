---
type: Functional Specification
title: debug-log.md legacy compatibility
description: Legacy runtime debug-log.md compatibility for historical evaluation runs.
tags: [evaluation, records, cli, skill, legacy]
timestamp: 2026-06-23T00:00:00Z
---

# debug-log.md legacy compatibility

This spec is part of the [Evaluation records](../evaluation-records.md) contract.
It documents compatibility for historical evaluation runs that contain
`debug-log.md`. Current evaluation workflow feedback is defined by the
[/quality evaluate feedback log](../skills/quality-skill/workflows/evaluate/feedback-log.md)
spec and is written under `.quality/logs/`.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Legacy artifact

Historical evaluation runs may contain a run-root `debug-log.md` file. Readers
that inspect evaluation history **SHOULD** tolerate the file when present.

New evaluation runs **MUST NOT** require or seed `debug-log.md`.

`debug-log.md` **MUST NOT** be authoritative for ratings, findings,
recommendations, next actions, generated report content, or reportability.

When present, `debug-log.md` **MUST NOT** be treated as an assessment record,
rating rationale, report, or evidence store. Evaluation findings and rating
rationale belong in formal records.

Tools **MUST NOT** migrate, rewrite, or reinterpret historical `debug-log.md`
files as current feedback logs. A fresh evaluation writes current workflow
feedback to `.quality/logs/<timestamp>-evaluate-feedback-log.md`.
