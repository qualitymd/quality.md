---
type: Functional Specification
title: Report regression coverage
description: Focused automated coverage for high-risk generated report behavior.
tags: [evaluation, report, tests]
timestamp: 2026-06-18T00:00:00Z
---

# Report regression coverage

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Scope

This change covers automated test coverage for generated evaluation reports.

It does not add committed subject fixtures, visual/browser report testing, or a
full benchmark harness.

## Requirements

The test suite **MUST** cover seeded secret-style report behavior: the report
surfaces the safety finding and recommendation path while not copying a seeded
secret value into generated report artifacts.

The test suite **MUST** cover prompt-injection-style report behavior: hostile
source content remains finding data, and the report exposes the finding without
following or rendering prompt-disclosure sentinel text.

The test suite **MUST** cover not-assessed report behavior: root and
requirement ratings remain `notAssessed` with `rating: null`.

The test suite **MUST** cover dotted-path limitation extraction so paths such as
`docs/production-telemetry.md` are not split into corrupt fragments.

The test suite **MUST** cover structural-root and empty-recommendation rendering:
structural grouping targets render distinctly and empty recommendations render
as an empty array.
