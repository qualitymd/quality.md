---
type: Design Doc
title: Evaluation history compatibility - design
description: Implementation approach for tolerant evaluation-history inspection without migrations.
tags: [evaluation, records, compatibility, status]
timestamp: 2026-06-20T00:00:00Z
---

# Evaluation history compatibility - design

Companion to
[Evaluation history compatibility - functional spec](spec.md).

## Context

The current evaluation package has one main run loader, and that loader is used
both by report-building paths that need trusted records and by history/status
paths that only need to describe what exists. This couples ordinary history
inspection to the strictest interpretation of the current record contract.

The design keeps record writers strict and changes the read side by introducing
a tolerant inspection model. A run with incompatible records remains visible,
but report generation refuses to trust it.

## Approach

Introduce a small run-inspection layer in `internal/evaluation` that can read a
run folder, basic metadata, record filenames, and any records that decode cleanly
under the current contract. Record files that cannot be decoded or trusted are
captured as `EvaluationRunGap` values rather than returned as fatal loader
errors.

Use that layer for history-oriented commands:

- `qualitymd status`
- `qualitymd evaluation list`
- `qualitymd evaluation status`
- `qualitymd evaluation <record-kind> list`

Keep the existing strict write path for `assessment add`, `analysis set`, and
`recommendation add`. Report build and gate should call through the same
inspection/status path first; if compatibility gaps exist, they return a clear
non-reportable-run diagnostic instead of attempting report assembly.

The compatibility gap set should be deliberately small and record-oriented:
malformed JSON, malformed recommendation frontmatter, missing `schemaVersion`,
unsupported `schemaVersion`, structurally incomplete current-schema record, and
unreadable record. Existing semantic gaps such as duplicate assessment results,
missing analysis, and invalid severity still belong in normal reportability
validation after a record decodes.

Counts should be computed from discovered filenames where possible. Trusted
report assembly should still use decoded current-schema records only.

## Alternatives

Making `Load` fully permissive was rejected because it blurs whether a caller is
working with trusted records or merely discovered history. It would be too easy
for report assembly to accidentally continue with partial data.

Adding schema migrations was rejected because the requested posture is not to
maintain old record shapes. Historical records can remain as audit artifacts; a
fresh evaluation is the forward path.

Silently skipping bad historical runs was rejected because `--latest` and status
would become misleading. A malformed latest run should remain the latest run and
should explain why it is not reportable.

## Trade-offs & risks

The main risk is introducing two read paths that drift. The mitigation is to make
the tolerant path the shared entry point for run status, and reserve strict
report assembly for the point after status has already proven there are no
blocking gaps.

Filename-based counts are useful for history but can overstate usable record
data when files are malformed. JSON output should keep that distinction visible
through gaps/problems rather than imply every counted file was trusted.

Recommendation Markdown is the oddest record type because its machine metadata
lives in frontmatter and its body is human-readable. Compatibility handling
should treat frontmatter failures as record gaps and avoid interpreting the body
as a fallback schema.

## Resolved During Implementation

- `evaluation list --json` includes a `gaps` count per run so callers can
  distinguish an empty incomplete run from one with blocking diagnostics.
- Record-list commands include discovered record filenames, including malformed
  record files. Detailed trust/reportability diagnostics stay in
  `evaluation status <run>`.
