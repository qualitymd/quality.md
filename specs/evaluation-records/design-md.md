---
type: Functional Specification
title: design.md
description: Runtime design.md contract for evaluation runs.
tags: [evaluation, records, cli, skill]
timestamp: 2026-06-22T00:00:00Z
---

# design.md

This spec is part of the [Evaluation records](../evaluation-records.md) contract.
It owns the independently reviewable runtime contract below; shared
responsibility, runtime-not-OKF status, schema-version, and historical-record
compatibility rules live in the parent spec.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", and "MAY" are to be interpreted as described
in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Evaluation Design

`design.md` is a Markdown artifact recording the resolved evaluation frame. The
CLI seeds the file when it creates a run, and the skill authors the initial
content before assessment evidence collection or record writes begin.

The initial `design.md` **MUST** record:

- the resolved model file and its relationship to the run's `model.md` snapshot;
- the evaluation mode and chosen rigor;
- the evaluation scope or narrowing;
- the in-scope areas;
- out-of-scope or deferred areas;
- methodological constraints or rating limitations known before assessment.

> Rationale: a design artifact can make the run reproducible only when it records
> the inputs and known limits before the evaluator starts producing judgment.
> Retrospective findings belong in formal records and reports. — 0056

`design.md` **MUST NOT** carry assessment findings, rating rationale, or
recommendation reasoning as formal judgment. Those belong in assessment,
analysis, recommendation, and report artifacts.

`design.md` **MAY** be updated during a run when the resolved evaluation frame
changes, but updates must preserve enough original context for a reader to
distinguish the initial frame from later corrections.
