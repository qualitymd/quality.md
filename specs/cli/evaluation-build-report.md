---
type: Functional Specification
title: qualitymd evaluation build-report
description: Render report.md and report.json from evaluation records.
tags: [cli, command, evaluation]
timestamp: 2026-06-18T00:00:00Z
---

# qualitymd evaluation build-report

`qualitymd evaluation build-report <run>` derives `report.md` and `report.json`
from the run's assessment, analysis, and recommendation records. It renders
recorded judgment; it **MUST NOT** infer or recompute ratings.

The command **MUST** fail without writing a partial report when the run is not
renderable. It **MUST** be deterministic and idempotent: unchanged records produce
byte-identical report files.

`--fail-at-or-below <level>` turns the command into a CI gate. The command still
writes both report files on a successful render. It exits `1` when the root
aggregate rating is equal to or worse than `<level>`, exits `0` when better, and
exits `2` when `<level>` is not in the run's rating scale. A root *not assessed*
result fails the gate.
