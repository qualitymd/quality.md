---
type: Functional Specification
title: Evaluation Report CTA
description: Requirements for value-oriented human report CTAs in /quality evaluate closeouts.
tags: [quality-skill, evaluation, reports, agent-mediated-ux]
timestamp: 2026-06-27T00:00:00Z
---

# Evaluation Report CTA

Companion to the
[Evaluation Report CTA](../0151-evaluation-report-cta.md) Change Case. This
spec states _what_ the change must do. It defers to the
[/quality reporting](../../../specs/skills/quality-skill/reporting.md),
[/quality evaluate](../../../specs/skills/quality-skill/workflows/evaluate.md),
[Evaluation report tree](../../../specs/evaluation/reports/report-tree.md), and
[Designing agent-mediated UX](../../../docs/guides/agent-mediated-ux.md)
contracts as normative context for report artifacts and user-facing closeout
shape.

The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

The evaluation closeout is the user's handoff from an agent-run assessment to
the durable report artifacts. A generic artifact list is not enough at that
moment: the user needs a clear reading path and a reason to follow it. The two
human reports serve different jobs. `report.md` is the decision-ready evaluation
result, while `recommendations.md` is the action-planning report. The closeout
should make that distinction visible in a five-second scan and avoid elevating
machine-oriented indexes as if they were user reading targets.

## Scope

Covered: the user-facing `/quality evaluate` closeout, durable skill specs that
govern reporting and evaluate completion, and runtime skill guidance for the
closeout template.

Deferred / non-goals: generated Markdown report content changes, report file
renames, CLI report-build receipt changes, recommendation follow-up behavior,
and machine-readable report index contracts.

## Requirements

R1. The `/quality evaluate` closeout **MUST** make the first report-reading call
to action point to the completed run's `report.md` path.

> Rationale: after evaluation, the most important next user action is reading
> the human evaluation result, not scanning a generic artifact list or inferring
> a run-relative filename.
>
> Durable spec: modify `specs/skills/quality-skill/reporting.md` and
> `specs/skills/quality-skill/workflows/evaluate.md` - require a primary
> report-reading CTA for the run-level report.

R2. The `/quality evaluate` closeout **MUST** describe `report.md` in
value-oriented terms as the decision-ready evaluation result containing the
rating, evidence basis, limits, top findings, and top recommendations.

> Rationale: "run report" names the artifact but does not tell the user why to
> open it or what decision support it contains.
>
> Durable spec: modify `specs/skills/quality-skill/reporting.md` - require a
> value description for the primary human report.

R3. The `/quality evaluate` closeout **MUST** name the completed run's
`recommendations.md` path and describe it in value-oriented terms as the
action-planning report containing ranked recommendations, why they matter,
expected benefit, and how to know each worked.

> Rationale: recommendation follow-up remains a separate workflow, but the user
> still needs a clear path to the advice artifact that supports planning.
>
> Durable spec: modify `specs/skills/quality-skill/reporting.md` and
> `specs/skills/quality-skill/workflows/evaluate.md` - require a user-facing
> recommendation report CTA.

R4. The `/quality evaluate` closeout **MUST** present human report paths as full
run-relative paths, not only bare filenames.

> Rationale: bare filenames force the user to join the artifact name with the
> run folder from another line; full paths make the CTA directly usable.
>
> Durable spec: modify `specs/skills/quality-skill/reporting.md` and
> `specs/skills/quality-skill/workflows/evaluate.md` - require full human report
> paths in the closeout.

R5. The `/quality evaluate` closeout **MUST NOT** include
`data/evaluation-output-result.json` or other machine-oriented report indexes in
the primary report-reading CTA.

> Rationale: the generated output index is useful to tooling and debugging, but
> it dilutes the user's next action after a successful evaluation.
>
> Durable spec: modify `specs/skills/quality-skill/reporting.md` - distinguish
> human report CTAs from machine report indexes in closeouts.

R6. The redesigned closeout **MUST** preserve the existing closeout facts:
rating, scope, evidence basis, known limitations, changed artifacts, not-done
boundary, and recommended next action.

> Rationale: strengthening the report CTA should not remove the evidence,
> boundary, or next-action signals that make the evaluation closeout trustworthy.
>
> Durable spec: modify `specs/skills/quality-skill/workflows/evaluate.md` -
> keep completion criteria aligned with the existing closeout contract while
> replacing the generic reports field.

R7. The redesigned closeout **MUST NOT** imply that evaluate applied
recommendations, edited evaluated source, edited `QUALITY.md`, wrote the quality
changelog, or created external issues.

> Rationale: a strong recommendation CTA can read like an action handoff unless
> the closeout preserves the evaluate workflow's mutation boundary.
>
> Durable spec: none; existing reporting and evaluate workflow specs already
> require the not-done boundary. Keep this as an implementation acceptance check.

## Durable spec changes

### To add

None

### To modify

- `specs/skills/quality-skill/reporting.md` - add the human report CTA contract
  for `report.md` and `recommendations.md`, and exclude machine report indexes
  from the primary closeout CTA (R1, R2, R3, R4, R5).
- `specs/skills/quality-skill/workflows/evaluate.md` - align evaluate
  completion criteria with full human report paths, report-reading CTA, and
  preserved closeout facts (R1, R3, R4, R6).

### To rename

None

### To delete

None

## Design decisions

- Use `Open next` as the primary report-reading CTA label.
- Keep `recommendations.md` on the `Recommendations` line with its value
  proposition.
