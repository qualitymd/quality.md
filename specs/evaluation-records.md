---
type: Functional Specification
title: Evaluation records
description: The deterministic on-disk contract for QUALITY.md evaluation run records.
tags: [evaluation, records, cli, skill]
timestamp: 2026-06-21T00:00:00Z
---

# Evaluation records

This spec defines the runtime record contract for a QUALITY.md evaluation run:
folder names, record names, record schemas, `schemaVersion`, and the division of
responsibility between the deterministic `qualitymd` CLI and the judging skill.
The evaluation semantics are defined by
[`SPECIFICATION.md` → Evaluation](../SPECIFICATION.md#evaluation).

This contract is a standalone spec — not prompt prose, and not nested under
`cli/` — because it has two consumers: the CLI that writes records and the skill
that supplies judgment. A single cited source of truth keeps the two surfaces
from drifting in a way duplicated prose could not.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../docs/reference/rfc2119.md) and
[RFC 8174](../docs/reference/rfc8174.md) when, and only when, they appear in all
capitals.

## Responsibility

The CLI **MUST** own file creation, serialization, run-folder numbering,
record numbering, `model.md` snapshotting, `schemaVersion` stamping, and
`report-summary.md` / `report.md` / `report.json` rendering. The skill supplies judgment content:
findings, ratings, rationales, roll-up judgment, and recommendations. The CLI
owns the `model.md` snapshot — rather than the skill — because the snapshot is
mechanically resolvable content, and keeping the record of *what was evaluated*
off the judging skill preserves the CLI-writes / skill-judges division.

The skill **MUST NOT** hand-author, number, serialize, or stamp evaluation
records when the corresponding CLI command exists. A separately tracked
numbering counter previously drifted and produced a real run-number collision;
numbering is therefore CLI-owned and derived from a single directory scan
(one past the highest present), so the on-disk folders are the single source of
truth and two writers cannot claim the same number from a stale counter.

## Runtime, Not OKF

Evaluation records are raw runtime outputs, not OKF concepts. A run folder
**MUST NOT** be treated as an OKF bundle and **MUST NOT** contain OKF
`index.md`, `log.md`, or `schema.md` semantics. Runtime Markdown frontmatter in
recommendation records is machine metadata, not OKF concept frontmatter.

## Runtime Contracts

Independently reviewable run-folder, record, artifact, and report-output
contracts live in child specs:

- [Run folder](evaluation-records/run-folder.md) - runtime folder naming and layout.
- [Assessment result record](evaluation-records/assessment-result-record.md) - JSON assessment result records.
- [Analysis record](evaluation-records/analysis-record.md) - JSON analysis records.
- [plan.md](evaluation-records/plan-md.md) - plan artifact and planned coverage metadata.
- [debug-log.md](evaluation-records/debug-log-md.md) - process-only debug log artifact.
- [Recommendation record](evaluation-records/recommendation-record-md.md) - Markdown recommendation records.
- [Report outputs](evaluation-records/report-outputs.md) - shared report-model and generated report-output invariants.

## Schema Version

Every JSON record (`assessments/*.json`, `analysis/*.json`,
`report.json`)
**MUST** carry top-level `schemaVersion: 1`.

Every CLI-written recommendation Markdown record **MUST** carry runtime YAML
frontmatter with `schemaVersion: 1`.

## Historical and Non-CLI Records

The current CLI writer is strict: new records **MUST** satisfy the active
contract and carry the active `schemaVersion`. Readers that inspect evaluation
history are tolerant: historical, partial, hand-edited, copied, or non-CLI
records can be present in a run folder without making ordinary history
inspection fail.

An individual record that cannot be trusted under the current contract makes
only that run non-reportable. Status/list readers **SHOULD** surface it as a
run gap that names the record path and reason, preserving any run metadata and
record-file counts that can be determined without trusting the malformed
payload. Tools **MUST NOT** migrate, rewrite, silently skip, or reinterpret old
record shapes as a compatibility mechanism. A fresh evaluation or explicit
correction through the current CLI is the forward path.

At minimum, incompatible-record gaps distinguish malformed JSON or runtime
frontmatter, unreadable records, missing `schemaVersion`, unsupported
`schemaVersion`, and structurally incomplete current-schema records.
