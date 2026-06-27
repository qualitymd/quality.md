---
type: Functional Specification
title: Evaluation Advice
description: Requirements for required evaluation Advice, finding rankings, recommendations, recommendation rankings, and recommendation reports.
tags: [evaluation, advice, recommendations, reports]
timestamp: 2026-06-27T00:00:00Z
---

# Evaluation Advice

Companion to the [Evaluation Advice](../0150-evaluation-advice.md) Change Case.
This spec states *what* the change must do. It defers to the
[Evaluation](../../../specs/evaluation/evaluation.md) contract,
[`SPECIFICATION.md`](../../../SPECIFICATION.md), and
[Modeling quality across domains](../../../docs/guides/model-quality-across-domains.md)
as normative context for evaluation semantics and domain-agnostic language.

The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

Evaluation currently produces evidence, ratings, roll-up analysis, and reports,
but its durable contracts either forbid generated recommendations or treat
Advice as optional. That leaves users without a persisted answer to "what should
we do with this quality signal?" and pushes recommendation judgment into
ad-hoc follow-up. Advice should be the required late phase that turns completed
assessment and analysis into traceable, domain-agnostic quality-management
recommendations.

The change keeps the evaluation layers distinct. Requirement Findings remain the
evidence observations. Analysis explains rating impact. A Finding ranking orders
existing Findings for advice synthesis and report navigation without creating a
new finding layer. Recommendations state next quality-management moves, which
can include review, preservation, monitoring, evidence improvement, Model
improvement, or intentional deferment. A Recommendation ranking orders those
moves by expected quality impact and confidence, not by implementation effort.

## Vocabulary

- **Advice phase** - the required evaluation phase after roll-up analysis and
  before report build that produces Finding ranking, Recommendations, coverage
  accounting, and Recommendation ranking.
- **Finding ranking** - a derived ordering of existing Requirement Findings
  within the evaluation scope. It does not change Finding severity, confidence,
  type, evidence, or rating effect.
- **Recommendation** - quality-management advice about the next warranted move
  for the evaluated entity, the Model, the evidence base, or future review. It is
  not synonymous with a repair task.
- **Recommendation ranking** - a derived ordering of Recommendations by expected
  quality-management impact and confidence.
- **Coverage accounting** - the Advice-phase ledger that records whether every
  persisted Finding is addressed by one or more Recommendations or is explicitly
  not advice-driving.

## Scope

Covered: required Advice in complete evaluations; the
`FindingRankingResult`, `RecommendationResult`, and
`RecommendationRankingResult` payload kinds; coverage accounting for every
Finding; reportability gates; generated `report.md` Advice summaries; generated
`recommendations.md` and `recommendations/<NNN>-<slug>.md`; and skill guidance
for producing domain-agnostic Advice.

Deferred / non-goals: no CLI judgment generation command; no recommendation
application or issue creation changes; no effort, ROI, quick-win, backlog
priority, or numeric score fields; no machine-readable frontmatter or YAML
appendix in generated recommendation Markdown; no compatibility mode for
reportable runs that omit Advice.

## Requirements

### Advice phase and reportability

R1. Evaluation **MUST** include a required Advice phase after completed
Requirement rating and Factor/Area roll-up analysis and before
`qualitymd evaluation report build`. The Advice phase sequence **MUST** be:
Finding ranking, Recommendation generation, coverage accounting,
Recommendation ranking, Advice closure.

> Rationale: Advice depends on completed evidence and roll-up judgment, and
> reports should project persisted Advice rather than synthesize it during
> rendering.
>
> Durable spec: modify `SPECIFICATION.md`, `specs/evaluation/evaluation.md`,
> `specs/evaluation/protocol.md`, `specs/evaluation/orchestration.md`,
> `specs/skills/quality-skill/evaluation.md`, and
> `specs/skills/quality-skill/workflows/evaluate.md` - replace optional or
> forbidden recommendation wording with required late-phase Advice and its
> ordered substeps.

R2. `qualitymd evaluation status <run>` **MUST** report a run as not reportable
unless its effective data contains one valid `FindingRankingResult`, one or more
valid `RecommendationResult` payloads, and one valid
`RecommendationRankingResult`.

> Rationale: a complete evaluation should always produce a next quality-management
> signal, even when that signal is review, preservation, monitoring, or
> intentional deferment rather than remediation.
>
> Durable spec: modify `specs/cli/evaluation-status.md`,
> `specs/evaluation/protocol.md`, and
> `specs/evaluation/records/payload-kinds.md` - define Advice payloads as
> reportability requirements.

R3. Reports **MUST** render Advice only from persisted Advice payloads and
**MUST NOT** rank Findings, create Recommendations, rank Recommendations, or add
coverage rationales during rendering.

> Rationale: the agent owns judgment and the CLI owns deterministic projection;
> report build must not become a hidden evaluator.
>
> Durable spec: modify `specs/evaluation/evaluation.md`,
> `specs/cli/evaluation-report.md`, `specs/evaluation/reports/report-tree.md`,
> and `specs/skills/quality-skill/reporting.md` - preserve the renderer boundary
> while allowing persisted Advice to render.

### Finding ranking

R4. `qualitymd evaluation data set` **MUST** accept a
`FindingRankingResult` payload kind that ranks every persisted Requirement
Finding in the evaluation scope exactly once.

> Durable spec: modify `specs/evaluation/records/payload-kinds.md` and
> `specs/cli/evaluation-data.md` - add the new agent-written payload kind and
> validation rule.

R5. Each `FindingRankingResult` entry **MUST** reference an existing
Requirement Finding, include a 1-based `rank`, include a priority tier
`P1`, `P2`, `P3`, or `P4`, and include a short rationale for the tier or rank.

> Rationale: Finding ranking requires evaluative judgment, but the judgment must
> be visible and constrained enough for review.
>
> Durable spec: modify `specs/evaluation/records/payload-kinds.md` - define
> ranked Finding entry fields and tier enum.

R6. `FindingRankingResult` **MUST NOT** create Findings, change Finding fields,
introduce new evidence, change ratings, or carry recommendation impact,
effort, ROI, or quick-win metadata.

> Rationale: Finding ranking is an advice input and report view over existing
> evidence, not a second Finding layer or task-planning surface.
>
> Durable spec: modify `SPECIFICATION.md`,
> `specs/evaluation/records/payload-kinds.md`, and
> `specs/evaluation/reports/report-tree.md` - distinguish Finding ranking from
> Findings and Recommendations.

R7. The `/quality` evaluation workflow **MUST** rank Findings using constrained
judgment over severity, Finding type, confidence, rating influence, scope
breadth, blocked evidence, and advice relevance, then sort mechanically inside
equivalent priority tiers by severity, rating influence, confidence, and model
order unless a stated rationale overrides the mechanical order.

> Rationale: pure field sorting misses advice relevance, while unconstrained
> ranking is hard to review; tiered judgment plus deterministic tie-breakers
> gives both usefulness and auditability.
>
> Durable spec: modify `specs/skills/quality-skill/evaluation.md` and
> `specs/skills/quality-skill/workflows/evaluate.md` - add ranking judgment
> guidance.

### Recommendations

R8. `qualitymd evaluation data set` **MUST** accept a `RecommendationResult`
payload kind for one Recommendation produced by the Advice phase.

> Durable spec: modify `specs/evaluation/records/payload-kinds.md` and
> `specs/cli/evaluation-data.md` - add the new agent-written payload kind.

R9. Each `RecommendationResult` **MUST** include a run-local ID, title, why it
matters, recommended next move, expected benefit, how to know it worked, impact,
confidence, and trace references.

> Rationale: these are the smallest user-facing fields that make a
> recommendation intelligible and handoff-ready without exposing internal
> taxonomy.
>
> Durable spec: modify `specs/evaluation/records/payload-kinds.md`,
> `specs/evaluation/reports/report-tree.md`, and
> `specs/skills/quality-skill/evaluation.md` - define required Recommendation
> content and report fields.

R10. `RecommendationResult.impact` **MUST** use one of `very_high`, `high`,
`medium`, or `low`, and `RecommendationResult.confidence` **MUST** use the
Evaluation confidence vocabulary.

> Rationale: impact needs enough range for ranking without colliding with Finding
> severity or implying numeric precision; confidence should reuse the existing
> evaluation vocabulary.
>
> Durable spec: modify `specs/evaluation/records/payload-kinds.md` - define
> Recommendation impact and confidence enums.

R11. `RecommendationResult` **MUST NOT** include effort, ROI, quick-win,
backlog-priority, or numeric score fields.

> Rationale: evaluation advice should be impact-led quality-management guidance,
> not backlog planning.
>
> Durable spec: modify `specs/evaluation/records/payload-kinds.md`,
> `specs/evaluation/evaluation.md`, and `specs/skills/quality-skill/evaluation.md`
> to forbid planning metadata in the v1 Recommendation contract.

R12. Each `RecommendationResult` **MUST** cite at least one trace reference to a
Finding, Rating Driver, Evaluation limit or incomplete input, or strength/high
rating evidence that supports review, preservation, monitoring, or raising the
bar.

> Rationale: required recommendations must stay grounded even when they are not
> repair work.
>
> Durable spec: modify `specs/evaluation/records/payload-kinds.md` and
> `SPECIFICATION.md` - require traceable recommendation support.

R13. The `/quality` evaluation workflow **MUST** generate Recommendations by
working from the ranked Findings and completed analysis toward the smallest
coherent quality-management move that materially addresses the evidence or
judgment limit.

> Rationale: this avoids both one-recommendation-per-finding chore lists and
> broad unactionable programs.
>
> Durable spec: modify `specs/skills/quality-skill/evaluation.md` and
> `specs/skills/quality-skill/workflows/evaluate.md` - add recommendation
> synthesis guidance.

R14. The `/quality` evaluation workflow **MUST** treat review, preservation,
monitoring, evidence improvement, Model improvement, and intentional deferment
with a review trigger as valid Recommendations when supported by the evaluation.

> Rationale: a complete evaluation should produce quality-management advice, not
> invent work when the best next move is review, preservation, monitoring, or
> no near-term change.
>
> Durable spec: modify `SPECIFICATION.md`,
> `specs/skills/quality-skill/evaluation.md`,
> `specs/skills/quality-skill/workflows/evaluate.md`, and
> `specs/skills/quality-skill/recommendation-follow-up.md` - broaden
> Recommendation semantics beyond repair tasks.

R15. Recommendation language in durable specs, runtime guidance, and examples
**MUST** remain quality-domain agnostic and **MUST NOT** assume software, code,
tests, CI, releases, bugs, issues, or product delivery as the default evaluated
domain.

> Rationale: QUALITY.md is quality-domain agnostic; Advice should work for
> documentation, data, research, services, household budgets, and other
> evaluated entities.
>
> Durable spec: modify `SPECIFICATION.md`,
> `specs/skills/quality-skill/evaluation.md`, and relevant examples - align
> Recommendation wording with
> `docs/guides/model-quality-across-domains.md`.

### Recommendation ranking and coverage accounting

R16. `qualitymd evaluation data set` **MUST** accept a
`RecommendationRankingResult` payload kind that ranks every
`RecommendationResult` in the run exactly once and records the completed Finding
coverage accounting.

> Durable spec: modify `specs/evaluation/records/payload-kinds.md` and
> `specs/cli/evaluation-data.md` - add the new agent-written payload kind and
> validation rule, including its role as the final Advice closure payload.

R17. Each `RecommendationRankingResult` ranked entry **MUST** reference an
existing `RecommendationResult`, include a 1-based `rank`, include impact,
include confidence, and include a ranking rationale.

> Durable spec: modify `specs/evaluation/records/payload-kinds.md` - define
> ranked Recommendation entry fields.

R18. The `/quality` evaluation workflow **MUST** rank Recommendations by impact
first, then by evidenced unblocking value, urgency when evidenced, confidence,
and stable Recommendation ID; it **MUST** require rationale for any order that
does not follow those tie-breakers. Recommendation ranking **MUST** happen after
Recommendation generation and coverage accounting.

> Rationale: Recommendation ranking should answer which quality-management move
> most improves the evaluated quality signal, not which move is easiest.
>
> Durable spec: modify `specs/skills/quality-skill/evaluation.md` and
> `specs/skills/quality-skill/workflows/evaluate.md` - add recommendation
> ranking judgment guidance and sequencing.

R19. `RecommendationRankingResult` **MUST** include coverage accounting for every
persisted Requirement Finding in the evaluation scope, with each Finding marked
exactly once as either `addressed_by_recommendation` or `not_advice_driving`.
Coverage accounting **MUST** be determined after `RecommendationResult` payloads
are drafted and before `RecommendationRankingResult` is finalized.

> Rationale: full coverage prevents findings from silently disappearing while
> still allowing the evaluator to decide that a Finding should not shape Advice.
>
> Durable spec: modify `specs/evaluation/records/payload-kinds.md`,
> `specs/cli/evaluation-status.md`, and `specs/evaluation/protocol.md` - define
> coverage accounting, sequence it after Recommendation generation, and require
> it for reportability.

R20. A coverage entry with `addressed_by_recommendation` **MUST** reference one
or more valid `RecommendationResult` IDs, and a coverage entry with
`not_advice_driving` **MUST** include a non-empty rationale.

> Durable spec: modify `specs/evaluation/records/payload-kinds.md` and
> `specs/cli/evaluation-status.md` - enforce coverage completeness and rationale
> requirements.

R21. A `RecommendationResult` **MAY** address multiple Findings when the
recommended move clearly covers them, but the coverage accounting **MUST NOT**
require one Recommendation per Finding.

> Rationale: coverage should ensure accountability without forcing duplicate or
> artificially small recommendations.
>
> Durable spec: modify `SPECIFICATION.md` and
> `specs/evaluation/records/payload-kinds.md` - distinguish coverage from
> one-to-one recommendation generation.

### Generated reports and artifacts

R22. `qualitymd evaluation report build <run>` **MUST** generate `report.md`
with a `Top Findings` table and a `Top Recommendations` table, each capped at 10
rows.

> Durable spec: modify `specs/evaluation/reports/report-tree.md`,
> `specs/cli/evaluation-report.md`, and `specs/skills/quality-skill/reporting.md`
> to add run-report Advice summary sections.

R23. The `Top Findings` table **MUST** render rows from `FindingRankingResult`
and link each row to the existing generated report location where the Finding is
defined or detailed.

> Rationale: Top Findings is a navigational view over existing Findings, not a
> new report-only finding list.
>
> Durable spec: modify `specs/evaluation/reports/report-tree.md` - define Top
> Findings rendering and links.

R24. The `Top Recommendations` table **MUST** render `Rank`,
`Recommendation`, `Impact`, and `Confidence`, and each Recommendation link
**MUST** point to the corresponding generated recommendation detail file.

> Durable spec: modify `specs/evaluation/reports/report-tree.md` - define the
> compact user-facing Recommendation table.

R25. `report.md` **MUST** always link to `recommendations.md`, including when
the run has 10 or fewer Recommendations.

> Rationale: the Recommendation index is the stable handoff entrypoint even when
> the run summary contains every Recommendation.
>
> Durable spec: modify `specs/evaluation/reports/report-tree.md` - require the
> index link from the run report.

R26. `qualitymd evaluation report build <run>` **MUST** generate
`recommendations.md` at the run root as the full Recommendation index.

> Durable spec: modify `specs/evaluation/records/data-layout.md`,
> `specs/evaluation/reports/report-tree.md`, and `specs/cli/evaluation-report.md`
> to add the Recommendation index artifact.

R27. `qualitymd evaluation report build <run>` **MUST** generate one
human-readable recommendation detail file per `RecommendationResult` under
`recommendations/<NNN>-<slug>.md`, where `NNN` is the zero-padded
Recommendation rank and `<slug>` is derived from the Recommendation title.

> Durable spec: modify `specs/evaluation/records/data-layout.md` and
> `specs/evaluation/reports/report-tree.md` - add Recommendation detail artifact
> paths.

R28. Each generated Recommendation detail file **MUST** render the
Recommendation title, rank, impact, confidence, why it matters, recommended next
move, expected benefit, how to know it worked, and trace links or references.

> Durable spec: modify `specs/evaluation/reports/report-tree.md` - define
> handoff-ready Recommendation detail content.

R29. Generated Recommendation Markdown **MUST NOT** include YAML frontmatter or a
machine-readable YAML appendix in this change.

> Rationale: the Markdown artifacts should stay human-first; structured
> inspection remains available through persisted `data/` payloads and future
> CLI support.
>
> Durable spec: modify `specs/evaluation/reports/report-tree.md` - exclude
> machine-oriented Recommendation Markdown metadata for this slice.

### Examples and quality control

R30. Advice QC in the `/quality` evaluation workflow **MUST** use progressive
validation: each Advice artifact is checked when it is created, and final Advice
closure checks only cross-artifact consistency.

> Rationale: shifting quality checks left catches missing fields, invalid
> references, unsupported recommendations, and coverage gaps before later Advice
> steps build on them.
>
> Durable spec: modify `specs/skills/quality-skill/evaluation.md` and
> `specs/skills/quality-skill/workflows/evaluate.md` - add progressive Advice QC.

R31. Progressive Advice QC **MUST** check the Finding ranking before
Recommendation generation, each Recommendation before coverage accounting,
coverage accounting before Recommendation ranking, and Recommendation ranking
before final Advice closure.

> Durable spec: modify `specs/skills/quality-skill/evaluation.md` and
> `specs/skills/quality-skill/workflows/evaluate.md` - define Advice QC
> checkpoints.

R32. The final Advice closure check **MUST** verify that the Advice payload set is
internally consistent, every Recommendation is ranked exactly once, every
Finding has a coverage disposition, every addressed Finding points to a valid
Recommendation, every not-advice-driving Finding has a rationale, and report
build has all required Advice inputs. It **MUST NOT** be the first check for
required Recommendation fields, trace validity, no-new-evidence discipline, or
coverage completeness.

> Rationale: final closure should catch integration errors, not serve as the
> first quality gate for each artifact.
>
> Durable spec: modify `specs/skills/quality-skill/evaluation.md` and
> `specs/skills/quality-skill/workflows/evaluate.md` - define final Advice
> closure.

R33. The durable spec or example updates for this change **MUST** include at
least one non-software example, preferably a household budget or similarly
non-software quality context, covering a Top Finding, a Recommendation,
coverage accounting, and a review-or-raise-bar Recommendation.

> Rationale: the example should test the domain-agnostic contract far away from
> software development and task-backlog assumptions.
>
> Durable spec: modify relevant evaluation, reporting, or skill example specs -
> add a non-software Advice example.

## Acceptance criteria

- AC1. `qualitymd evaluation status <run>` reports a run with complete
  assessment and roll-up data but no Advice payloads as not reportable. *(R2)*
- AC2. A run with one or more Requirement Findings, a valid
  `FindingRankingResult`, at least one valid `RecommendationResult`, a valid
  `RecommendationRankingResult`, and complete coverage accounting is reportable
  when all non-Advice reportability requirements are also met. *(R2, R19, R20)*
- AC3. `data set` rejects a `FindingRankingResult` that omits a persisted
  Finding, references a non-existent Finding, duplicates a Finding, or carries
  disallowed effort/ROI/quick-win fields. *(R4, R6)*
- AC4. `data set` rejects a `RecommendationResult` without title, why it
  matters, recommended next move, expected benefit, how to know it worked,
  impact, confidence, or trace support. *(R9, R12)*
- AC5. `data set` rejects Recommendation impact outside
  `very_high | high | medium | low` and Recommendation confidence outside the
  Evaluation confidence vocabulary. *(R10)*
- AC6. `data set` rejects a `RecommendationResult` with effort, ROI, quick-win,
  backlog-priority, or numeric score fields. *(R11)*
- AC7. `data set` rejects a `RecommendationRankingResult` that omits,
  duplicates, or references a non-existent Recommendation. *(R16, R17)*
- AC8. `data set` or status validation rejects coverage accounting where an
  addressed Finding has no Recommendation reference, a not-advice-driving
  Finding has no rationale, or any persisted Finding is unaccounted for. *(R19,
  R20)*
- AC9. `evaluation report build` renders `report.md` with capped Top Findings
  and Top Recommendations tables and a link to `recommendations.md`. *(R22-R25)*
- AC10. `evaluation report build` writes `recommendations.md` and
  `recommendations/<NNN>-<slug>.md` files whose detail pages include the core
  user-facing Recommendation fields and trace links. *(R26-R28)*
- AC11. Generated Recommendation Markdown contains no YAML frontmatter or YAML
  appendix. *(R29)*
- AC12. Durable specs or examples include a non-software Advice example that
  demonstrates Top Findings, Recommendations, coverage accounting, and
  review-or-raise-bar advice. *(R15, R33)*
- AC13. Durable skill workflow guidance shows the Advice sequence as Finding
  ranking, Recommendation generation, coverage accounting, Recommendation
  ranking, and Advice closure, with progressive QC checkpoints at each artifact
  boundary. *(R1, R18, R19, R30-R32)*

## Durable spec changes

### To add

None

### To modify

- `SPECIFICATION.md` - make Advice required, define domain-agnostic
  Recommendations, Finding ranking, Recommendation ranking, and coverage
  accounting, and replace optional Advice wording.
- `specs/evaluation/evaluation.md` - remove the v0 recommendation-generation
  prohibition and add the Advice phase as a shared invariant.
- `specs/evaluation/protocol.md` - place Advice after roll-up and before report
  build; define reportability dependencies.
- `specs/evaluation/orchestration.md` - update orchestration order and
  responsibilities for agent-authored Advice and CLI-rendered reports.
- `specs/evaluation/records/payload-kinds.md` - add
  `FindingRankingResult`, `RecommendationResult`, and
  `RecommendationRankingResult` contracts, enums, references, and coverage
  accounting.
- `specs/evaluation/records/data-layout.md` - add `recommendations.md` and
  `recommendations/<NNN>-<slug>.md` generated report paths.
- `specs/evaluation/reports/report-tree.md` - render Top Findings, Top
  Recommendations, Recommendation index, and Recommendation detail reports while
  preserving renderer non-judgment.
- `specs/cli/evaluation-data.md` - accept and validate the new Advice payload
  kinds.
- `specs/cli/evaluation-status.md` - make Advice and coverage accounting part
  of reportability.
- `specs/cli/evaluation-report.md` - render persisted Advice artifacts without
  generating judgment.
- `specs/skills/quality-skill/evaluation.md` - update skill evaluation behavior
  to produce required Advice.
- `specs/skills/quality-skill/workflows/evaluate.md` - add Advice production,
  ranking, coverage, progressive QC, and closure steps.
- `specs/skills/quality-skill/reporting.md` - align reporting closeout with
  Advice artifacts.
- `specs/skills/quality-skill/recommendation-follow-up.md` and guide specs -
  keep follow-up compatible with Recommendations as quality-management advice,
  not only repair tasks.
- Relevant spec logs - record the contract change.

### To rename

None

### To delete

None
