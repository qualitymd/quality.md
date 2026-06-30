---
type: Functional Specification
title: Requirement Findings Only - functional spec
description: What Evaluation must do to make Requirement Findings the only finding layer and require finding-backed ratings.
tags: [evaluation, findings, reports, skill, cli]
timestamp: 2026-06-27T00:00:00Z
---

# Requirement Findings Only - functional spec

Companion to the
[Requirement Findings Only](../0142-requirement-findings-only.md) change case.
This spec states _what_ the change must do.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

The Evaluation evidence trail should have one canonical layer. Requirement
Findings bind observations to Requirement criteria and evidence; Requirement
ratings interpret those findings against the configured Rating Scale; Factor and
Area analysis roll those lower-level results up. Allowing a rated Requirement
with no findings breaks that chain, while Area Findings duplicate it by adding a
second synthesis finding layer above Requirements.

Keeping findings at Requirement level makes ratings auditable without assuming a
fixed Rating Scale such as target/sub-target. A configured Rating Level is valid
only when the paired Requirement Findings support that selected level. Factor and
Area outputs explain roll-up judgment through `ratingDrivers`, rationale,
confidence, limits, and incomplete inputs, not through new findings.

## Scope

Covered: Evaluation data schema version 3, CLI `data set` validation, report
generation, durable Evaluation specs, `/quality` runtime and durable skill specs,
tests, generated schema, logs, and release notes.

Deferred:

- recommendation generation;
- semantic proof that a finding is sufficient for a configured Rating Level;
- migrations, fallback readers, compatibility aliases, or legacy report rendering
  for schema version 2 Evaluation runs; and
- new machine-readable capability versions beyond the existing CLI/skill
  compatibility boundary.

## Requirements

### Requirement Findings and ratings

- Requirement Findings **MUST** be the only Evaluation findings.

  > Rationale: Requirements already bind evidence to criteria and may connect to
  > multiple Factors, so they are the right granularity for evidence.
  >
  > Durable spec: modify `SPECIFICATION.md`,
  > `specs/evaluation/records/payload-kinds.md`, and
  > `specs/evaluation/routines/routine-contracts.md`.

- A `RequirementRatingResult` with `status: rated` **MUST** be backed by a paired
  `RequirementAssessmentResult` for the same Requirement whose status is
  `assessed` or `partially_assessed` and whose `findings` array contains at
  least one Requirement Finding.

  > Rationale: a selected Rating Level without findings has no checkable
  > evidence trail.
  >
  > Durable spec: modify `SPECIFICATION.md`,
  > `specs/evaluation/records/payload-kinds.md`, and
  > `specs/cli/evaluation-data.md`.

- Requirement rating **MUST** remain scale-agnostic: it must justify the selected
  configured Rating Level against the model's applied criteria and **MUST NOT**
  assume fixed meanings such as target, sub-target, pass, or fail.

  > Durable spec: modify `SPECIFICATION.md` and
  > `specs/evaluation/routines/routine-contracts.md`.

- If the paired Requirement Findings are absent, too weak, or insufficient to
  distinguish the selected configured Rating Level from adjacent levels, the
  Requirement Rating Result **MUST NOT** use `status: rated`.

  > Durable spec: modify `SPECIFICATION.md`,
  > `specs/evaluation/routines/routine-contracts.md`, and
  > `specs/skills/quality-skill/evaluation.md`.

- A rated Requirement **MUST** carry non-empty `ratingDrivers` that cite the
  paired Requirement Assessment through `inputRefs`; the driver references
  **SHOULD** select the specific Requirement Findings that drove the rating.

  > Rationale: the CLI can validate that drivers exist and point to persisted
  > inputs, while the skill remains responsible for judgment sufficiency.
  >
  > Durable spec: modify `specs/evaluation/records/payload-kinds.md` and
  > `specs/cli/evaluation-data.md`.

### Roll-up analysis

- Factor and Area analysis **MUST NOT** produce findings.

  > Durable spec: modify `SPECIFICATION.md`,
  > `specs/evaluation/routines/routine-contracts.md`, and
  > `specs/evaluation/records/payload-kinds.md`.

- A Factor or Area analysis scope with `status: analyzed` and a `ratingLevelId`
  **MUST** carry non-empty `ratingDrivers`.

  > Rationale: after Area Findings are removed, rating drivers are the durable
  > roll-up explanation layer.
  >
  > Durable spec: modify `specs/evaluation/records/payload-kinds.md` and
  > `specs/cli/evaluation-data.md`.

- Roll-up `ratingDrivers` **MUST** cite lower-level routine outputs through
  `inputRefs` and **MUST NOT** introduce evidence or claims that are absent from
  those referenced outputs.

  > Durable spec: modify `SPECIFICATION.md`,
  > `specs/evaluation/routines/routine-contracts.md`, and
  > `specs/skills/quality-skill/evaluation.md`.

- Area Analysis Results **MUST NOT** contain a `findings` field,
  `factorRelationships`, or any Area Finding object shape.

  > Durable spec: modify `specs/evaluation/records/payload-kinds.md` and
  > `specs/cli/evaluation-data.md`.

### CLI validation and schema

- Evaluation data schema version **MUST** be bumped from `2` to `3`.

  > Durable spec: modify `specs/evaluation/evaluation.md` and
  > `specs/cli/evaluation-data.md`.

- `qualitymd evaluation data set` **MUST** reject payloads whose
  `schemaVersion` is not `3`.

  > Durable spec: modify `specs/cli/evaluation-data.md`.

- `qualitymd evaluation data set` **MUST** validate cross-payload rated-result
  invariants against the effective run data: existing persisted payloads overlaid
  with every candidate payload in the batch.

  > Rationale: batch order should not matter, and incremental correction should
  > be possible when the required paired payload already exists.
  >
  > Durable spec: modify `specs/cli/evaluation-data.md`.

- `qualitymd evaluation data set` **MUST** reject a rated Requirement result when
  the effective run data lacks a paired assessed or partially assessed
  Requirement Assessment with at least one finding.

  > Durable spec: modify `specs/cli/evaluation-data.md`.

- `qualitymd evaluation data set` **MUST** reject rated Requirement, Factor, or
  Area analysis results with empty or absent `ratingDrivers`.

  > Durable spec: modify `specs/cli/evaluation-data.md`.

- `qualitymd evaluation data set` **MUST** reject `ratingDrivers[].inputRefs`
  that do not resolve to existing routine outputs in the effective run data.

  > Durable spec: modify `specs/cli/evaluation-data.md`.

- `qualitymd evaluation data verify` **MUST** enforce the same schema version,
  structural, model-bound, and cross-payload invariants as `data set`.

  > Durable spec: modify `specs/cli/evaluation-data.md`.

- The CLI **MUST NOT** migrate, transform, accept, or render schema version 2
  Evaluation payloads as current data.

  > Durable spec: modify `specs/evaluation/evaluation.md` and
  > `specs/cli/evaluation-data.md`.

### Reports and skill behavior

- Area and Factor reports **MUST NOT** render `Findings` sections.

  > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
  > `specs/skills/quality-skill/reporting.md`.

- Requirement reports **MUST** continue to render Requirement Findings from
  `RequirementAssessmentResult.findings`.

  > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

- `/quality evaluate` **MUST** treat a rated Requirement with no Requirement
  Findings as invalid workflow output to correct before persistence or reporting.

  > Durable spec: modify `specs/skills/quality-skill/evaluation.md` and
  > `specs/skills/quality-skill/workflows/evaluate.md`.

- `/quality evaluate` **MUST** use Requirement Findings for evidence, and use
  rating drivers, rationale, limits, incomplete inputs, and confidence for Factor
  and Area roll-up explanation.

  > Durable spec: modify `specs/skills/quality-skill/evaluation.md`,
  > `specs/skills/quality-skill/workflows/evaluate.md`, and
  > `specs/skills/quality-skill/reporting.md`.

## Durable spec changes

### To add

None.

### To modify

- `SPECIFICATION.md`
- `specs/evaluation/evaluation.md`
- `specs/evaluation/protocol.md`
- `specs/evaluation/routines/routine-contracts.md`
- `specs/evaluation/records/payload-kinds.md`
- `specs/evaluation/reports/report-tree.md`
- `specs/cli/evaluation-data.md`
- `specs/skills/quality-skill/evaluation.md`
- `specs/skills/quality-skill/workflows/evaluate.md`
- `specs/skills/quality-skill/reporting.md`

### To rename

None.

### To delete

None.

## Verification

- `go test ./internal/evaluation`
- `go test ./...`
- `mise run fmt-md-check`
- `mise run check`
- Search for active Area Finding and `AreaAnalysisResult.findings` references,
  excluding archived Change Cases and historical changelog entries.
