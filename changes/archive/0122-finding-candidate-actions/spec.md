---
type: Functional Specification
title: Finding-level candidate actions — functional spec
description: Type the finding actions field as non-binding candidate actions, instruct the assessor to populate them on shortcoming findings, and keep them out of the Evaluation v0 report and closeout.
tags: [evaluation, skill, schema, advise]
timestamp: 2026-06-26T00:00:00Z
---

# Finding-level candidate actions — functional spec

Companion to the
[Finding-level candidate actions](../0122-finding-candidate-actions.md) change
case. This spec states *what* the change must do.

**Normative references:**

- [Evaluation](../../../../specs/evaluation/evaluation.md) and its
  [routine contracts](../../../../specs/evaluation/routines/routine-contracts.md) —
  the Evaluation protocol, the agent/CLI division of labor, and the v0
  "no recommendation generation" boundary this change must not cross.
- [Evaluation payload kinds](../../../../specs/evaluation/records/payload-kinds.md) —
  the requirement that every persisted payload field is validated from one typed
  source of truth.
- [SPECIFICATION.md](../../../../SPECIFICATION.md) — the format-level Assess and
  Advice phases.

**Informational:** the option-B
[advise sketch](advise-sketch.md) explains the downstream consumer that shapes
why candidate actions are captured the way they are.

The key words "MUST", "MUST NOT", "SHOULD", "SHOULD NOT", and "MAY" are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background

The evaluation data contract already declares `findings[].actions`, the generated
schema carries it, and the report tree renders an `Actions` row — but the field is
an untyped `arrayOfAny()`, has no authoring guidance, and is never populated, so
every finding renders `Actions | —`. Findings are therefore purely diagnostic,
and the remediation context that is richest at assessment time (the assessor has
the source open) is discarded.

A future **Advise** phase will synthesize a final, prioritized recommendation set
but will only see finding *descriptions*, not the source. Capturing finding-local
remediation leads as typed raw material preserves that context for it. Evaluation
v0 forbids recommendation generation, so this raw material MUST stay finding-local
and non-binding and MUST NOT appear as recommendations in the v0 report or
closeout. The distinction this spec draws — **candidate action** (a finding-local,
non-binding lead produced during assessment) versus **recommendation** (a
synthesized, prioritized, selected remediation produced by the Advise phase) — is
load-bearing for staying inside the v0 boundary.

## Scope

Covered: the typed candidate-action shape; its CLI validation; the assessment
authoring guidance; the example payload; and the v0 report/closeout exclusion.

Deferred:

- The Advise phase (option B): synthesis, dedup, prioritization, option modeling,
  recommendation records, and any new payload kind. See [advise-sketch.md](advise-sketch.md).
- A stable per-action identifier for downstream reference (the Advise phase can
  address an action by finding id plus position until it needs more).

Non-goals:

- Presenting candidate actions to the user, in the report tree or the closeout —
  that would be recommendation presentation, which v0 prohibits.
- Aggregating, deduplicating, or ranking candidate actions across findings.
- Subjecting candidate actions to the headline-binding-finding evidence re-check;
  they are speculative leads, not evidence-backed findings.

## Assumptions & dependencies

- The finding shape is defined by one typed source of truth
  (`internal/evaluation/data_contract.go`) from which the discovery schema,
  validation, schema output, and example are derived, per
  [Evaluation payload kinds](../../../../specs/evaluation/records/payload-kinds.md).
  If finding structure ever moves out of that single source, the validation and
  example requirements below must follow it.
- Reports are deterministic projections over `data/`; removing a rendered row
  does not remove the underlying persisted data.

## Requirements

### Candidate-action shape

- The finding `actions` field **MUST** be an array of **candidate-action
  objects**. Each candidate-action object **MUST** carry a `description` string
  (the non-binding remediation lead) and **MAY** carry a `rationale` string (why
  it would help or what shortcoming it addresses). It **MUST NOT** carry other
  fields.

  > Rationale: `description` + optional `rationale` mirrors the finding's own
  > description/rationale pairing and gives the Advise phase enough to cluster on
  > without inventing speculative structure (cost/effort/impact) before there is a
  > consumer. — 0122

- The `actions` field **MUST** remain optional: a finding with no `actions`, or an
  empty `actions` array, **MUST** be a valid finding.

### Validation

- When persisting a `RequirementAssessmentResult`,
  `qualitymd evaluation data set` **MUST** validate each candidate-action object
  from the typed source of truth: it **MUST** reject a candidate action that omits
  `description`, that carries an unknown field, or whose `description`/`rationale`
  is not a string, before writing the payload.

  > Rationale: payload kinds are validated from one typed source of truth;
  > an untyped `arrayOfAny()` accepts malformed actions silently and defeats the
  > point of capturing structured raw material. — 0122

### Authoring guidance

- When the skill authors a `RequirementAssessmentResult`, it **MUST** record at
  least one candidate action on each `gap` and `risk` finding, and **MUST NOT**
  record candidate actions on `strength` findings.

  > Rationale: shortcoming findings are where remediation context exists; actions
  > on strengths are noise. Tying authoring to finding `type` keeps the raw
  > material targeted. — 0122

- The skill **SHOULD** ground the candidate-action shape from the
  `RequirementAssessmentResult` example
  (`qualitymd evaluation data example`) rather than from recall, and the example
  **MUST** include at least one finding carrying a candidate action.

### v0 boundary

- Candidate actions **MUST** be stored only within their finding in the
  `RequirementAssessmentResult`. The evaluation **MUST NOT** aggregate,
  deduplicate, prioritize, or otherwise synthesize candidate actions across
  findings.

  > Rationale: cross-finding synthesis is recommendation generation, which
  > Evaluation v0 forbids; keeping actions finding-local is what makes the harvest
  > layer v0-safe. — 0122

- The Evaluation v0 Markdown report **MUST NOT** render candidate actions, and the
  user-facing closeout **MUST NOT** present them. In particular, the finding
  detail table **MUST NOT** include the `Actions` row.

  > Rationale: rendering candidate actions would present them as the evaluation's
  > advice; v0 must not present generated recommendations. The data remains in
  > `data/` for the Advise phase. — 0122

## Durable spec changes

### To add

None.

### To modify

- `specs/evaluation/routines/routine-contracts.md` — state that requirement
  assessment MAY record finding-local candidate actions and MUST NOT synthesize,
  aggregate, or prioritize them (per the *Candidate-action shape* and *v0
  boundary* requirements above).
- `specs/evaluation/reports/report-tree.md` — state that finding detail does not
  render candidate actions in v0 (per the *v0 boundary* report requirement above).
- `specs/skills/quality-skill/evaluation.md` — state that the skill records
  non-binding candidate actions on `gap`/`risk` findings as raw material for a
  later Advise phase (per the *Authoring guidance* requirements above).
- `specs/skills/quality-skill/reporting.md` — state that candidate actions are not
  recommendations and that the v0 report and closeout exclude them (per the *v0
  boundary* requirements above).
- `SPECIFICATION.md` — in *Assess Requirements*, note that a Finding MAY carry
  non-binding candidate actions, distinct from the Advice phase's recommendations
  (per the *Candidate-action shape* and *v0 boundary* requirements above).

### To rename

None.

### To delete

None.
