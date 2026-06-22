---
type: Functional Specification
title: Evaluation run folder
description: Runtime folder naming and layout contract for QUALITY.md evaluation runs.
tags: [evaluation, records, cli, skill]
timestamp: 2026-06-22T00:00:00Z
---

# Evaluation run folder

This spec is part of the [Evaluation records](../evaluation-records.md) contract.
It owns the independently reviewable runtime contract below; shared
responsibility, runtime-not-OKF status, schema-version, and historical-record
compatibility rules live in the parent spec.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Run Folder

Each run folder is named:

```text
NNNN[-<narrowing>]-quality-eval
```

`NNNN` is one zero-padded sequence across the evaluation directory, computed by directory scan immediately before creation (one past the highest present). `<narrowing>` is omitted for an unnarrowed run. An existing folder at the computed name **MUST** fail rather
than be overwritten or silently retried at the next number: under a deterministic
shared sequence a collision signals concurrent or corrupt state that must
surface, not be papered over. The scan-then-create step is safe but not
serializable — concurrent invocations are resolved by the loser failing on an
existing-folder error — and a lock is deferred until concurrent use is a real
requirement.

A run folder contains:

```text
model.md
design.md
debug-log.md
plan.md
assessments/
  NNN-<area>-<requirement>.json
analysis/
  <area>.json
recommendations/
  NNN-<slug>.md
report.md
report-summary.md
report.json
```

`assessments/` and `recommendations/` each use their own local `NNN`
sequence.

`debug-log.md` is a process-only Markdown artifact. `plan.md` may be body-only,
or it may carry YAML frontmatter with optional planned coverage metadata.
