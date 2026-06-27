---
type: Design Doc
title: Evaluation Report CTA
description: Design for value-oriented human report CTAs in /quality evaluate closeouts.
tags: [quality-skill, evaluation, reports, agent-mediated-ux]
timestamp: 2026-06-27T00:00:00Z
---

# Evaluation Report CTA

Companion to the [Evaluation Report CTA](../0151-evaluation-report-cta.md)
Change Case and its [functional spec](spec.md).

## Context

The change is an agent-mediated UX correction. The generated report tree already
has the right human artifacts: `report.md` is the run entrypoint and
`recommendations.md` is the recommendation index/detail gateway. The gap is the
agent closeout after `/quality evaluate`: it currently names those reports like
artifacts instead of directing the user to the next reading action and explaining
why each report matters.

## Approach

Keep the generated reports, report paths, CLI receipt, and machine-readable
`data/evaluation-output-result.json` contract unchanged. Update only the
closeout contract and runtime guidance.

The closeout keeps the existing status-first fields, but replaces the generic
`Reports` line with a primary `Open next` field:

```text
**Open next:** `<run>/report.md` - the decision-ready evaluation result: rating, evidence basis, limits, top findings, and top recommendations.
```

`recommendations.md` stays on the `Recommendations` line because that line
already carries advice context. It now points at the full run-relative path and
states the report's value:

```text
**Recommendations:** `<run>/recommendations.md` - the action-planning report: ranked recommendations, why they matter, expected benefit, and how to know each worked.
```

The closeout does not mention `data/evaluation-output-result.json` in the report
CTA. The machine index remains available through the CLI receipt and generated
report links, but it is not a useful primary action for the user after a
successful evaluation.

## Spec response

- R1 and R2 are satisfied by making `Open next` point to `<run>/report.md` and
  giving it the decision-ready value proposition.
- R3 is satisfied by keeping the advice artifact on `Recommendations` with the
  full `<run>/recommendations.md` path and action-planning value proposition.
- R4 is satisfied by using `<run>/...` paths in the runtime template and durable
  closeout contract.
- R5 is satisfied by excluding generated data/index files from the CTA.
- R6 and R7 are satisfied by preserving the existing rating, scope, evidence,
  limitations, changed, not-done, and next-action fields.

## Alternatives

**Keep a `Reports` field and improve its wording.** Rejected because it still
makes the user parse a category instead of seeing the next action. It also puts
both reports at equal weight even though `report.md` is the first read.

**Put both reports under `Open next`.** Rejected because the field becomes a
mini artifact list. Keeping `recommendations.md` on the `Recommendations` line
preserves the distinction between the first read and the advice-planning path.

**Mention the report index as supporting data.** Rejected for the closeout CTA.
The index is still a real artifact, but surfacing it here competes with the
human reading path.

**Change generated report headings or content.** Rejected as out of scope. The
current reports already contain the content the CTA promises.

## Trade-offs & risks

`Open next` is slightly imperative, but that is appropriate at workflow close:
the user needs the next reading target. The phrase also avoids implying that the
agent will open the file automatically.

The value propositions must stay true as report contents evolve. If future
report-tree changes remove top findings, top recommendations, limits, or
recommendation detail fields, the closeout contract should be revisited in the
same change.

## Open questions

None. The design chooses `Open next` for the primary report CTA and keeps
`recommendations.md` on the `Recommendations` line.
