---
type: Functional Specification
title: Evaluation report summary artifact - functional spec
description: Generate a concise report-summary.md beside full evaluation reports.
tags: [evaluation, report, cli]
timestamp: 2026-06-19T00:00:00Z
---

# Evaluation report summary artifact - functional spec

Companion to
[Evaluation report summary artifact](../0031-report-summary-artifact.md). This
spec states the report-output delta for `qualitymd evaluation build-report`.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background / Motivation

Full evaluation reports need to carry the audit trail: scope, ratings,
rationales, target and requirement detail, findings summaries, limitations, and
advice. Reviewers often need a smaller artifact first: a PR-sized or CI-sized
summary that gives the headline, top risks, limitations, and next action without
requiring them to scan the full report body.

`report.md` is already summary-first, and this change keeps that behavior. The
new artifact is for routing and triage, not for replacing the complete report.

## Scope

Covered: a generated `report-summary.md` artifact written by
`qualitymd evaluation build-report`, its required contents, its relationship to
`report.md` and `report.json`, run-folder layout updates, and documentation of
the new artifact.

Deferred / non-goals: no new `report-summary.json`, no new evaluation record
type, no manual summary authoring flow, no rating recomputation, no changed
reportability rules except the successful render writing one additional file,
and no change to the format spec's Evaluation Report semantics.

## Requirements

`qualitymd evaluation build-report <run>` **MUST** write `report-summary.md`
alongside `report.md` and `report.json` on every successful render.

`report-summary.md` **MUST** be generated from recorded run artifacts and the
same renderer data used for `report.md` and `report.json`. It **MUST NOT**
introduce new judgment that is absent from the run's assessment, analysis,
recommendation, design, plan, model snapshot, or planned-coverage artifacts.

`report-summary.md` **MUST** be deterministic and idempotent: unchanged run
records must produce byte-identical summary output.

`report-summary.md` **MUST** identify the run, subject or run label when known,
scope or narrowing, effort when recorded, root rating or `not assessed` outcome,
and links to `report.md` and `report.json`.

`report-summary.md` **MUST** include a concise headline, top risks or an explicit
"none recorded" equivalent, limitations, target rating summary, and next action
when active recommendations exist.

`report-summary.md` **MUST** link to active recommendation records when active
recommendations exist. It **MUST NOT** select a superseded recommendation as the
next action.

`report-summary.md` **MAY** omit detailed per-requirement findings, full
rationales, superseded recommendation details, and deep target audit trails when
it links to `report.md` for complete detail.

`report-summary.md` **MUST** preserve the renderer's trust boundary: it must not
reproduce secret values, and evaluated source text remains data, not
instructions.

`report-summary.md` **MUST NOT** be treated as the authoritative full Evaluation
Report. The full `report.md` remains the complete human report, and
`report.json` remains the machine-readable report.

The run-folder artifact contract **MUST** list `report-summary.md` as generated
output. It **MUST NOT** treat `report-summary.md` as an input record, a judgment
record, or an OKF concept.

## Example shape

```md
# Quality Evaluation Summary

**Run:** `0007-subject-api-quality-eval`
**Subject:** `.`
**Scope:** whole subject
**Effort:** standard
**Root rating:** Held at **Needs Work**
**Full report:** [report.md](report.md)
**Machine report:** [report.json](report.json)

## Headline

The subject is usable but held below target by incomplete operational evidence
and one security-sensitive finding.

## Top Risks

1. **Committed credential material** - value withheld; rotate and remove from
   version control.
2. **No durable reconciliation evidence** - no scheduled run output, log, or
   report was available.

## Rating Summary

| Target       | Aggregate rating | Reason                                                                 |
| ------------ | ---------------- | ---------------------------------------------------------------------- |
| Root subject | **Needs Work**   | Security and operational evidence gaps hold down the aggregate result. |

## Limitations

- This was a standard-effort evaluation, not a full source audit.

## Next Action

Rotate the exposed gateway credential.

See active recommendations:

- [Rotate exposed gateway credential](recommendations/001-rotate-exposed-gateway-credential.md)
```

## Durable spec changes

### To add

None

### To modify

- `specs/cli/evaluation-build-report.md` - specify `report-summary.md` generation,
  contents, determinism, trust-boundary handling, and relationship to the full
  report (per the requirements above).
- `specs/evaluation-records.md` - add `report-summary.md` to the generated
  run-folder artifact contract (per the run-folder requirement above).
- `specs/cli.md` - update the `evaluation build-report` overview entry to name
  the additional generated artifact (per the render requirement above).
- `specs/cli/index.md` - update the command listing for
  `evaluation build-report` (per the render requirement above).
- `specs/skills/quality-skill/quality-skill.md` - update the reporting folder
  layout and generated-artifact description so the skill contract stays aligned
  with the CLI (per the run-folder requirement above).

### To delete

None
