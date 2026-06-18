---
type: Functional Specification
title: Skill rigor and efficiency - functional spec
description: Requirements on the /quality skill's evaluation behavior — operationalized effort, verified evidence, rating-binding re-checks, batched writes, and optional deep fan-out.
tags: [skill, evaluation, rigor]
timestamp: 2026-06-17T00:00:00Z
---

# Skill rigor and efficiency - functional spec

This spec states the delta for
[Skill rigor and efficiency](../0017-skill-rigor-efficiency.md). It constrains the
behavior of the [`/quality` skill](../../specs/skills/quality-skill/quality-skill.md)
prompt during evaluation; it specifies no CLI behavior. Where a requirement is
later superseded by CLI-written records, it stands on its own until then.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Scope

Covered: the skill's `evaluate` flow and the evaluation half of `improve`.
Intentionally **deferred**: any change to the `qualitymd` CLI; `wizard` and
`setup`, which are read-only or scaffolding and do not produce ratings.

## Requirements

### Operationalized effort

The skill **MUST** select assessment breadth and verification depth from the
`Effort` argument by these observable rules:

- `quick` — assess only apex and high-risk in-scope requirements; evidence MAY be
  sampled rather than exhaustive.
- `standard` — assess every in-scope requirement; gather targeted evidence
  sufficient to bind each rating.
- `deep` — assess every in-scope requirement against a full read of each target's
  source, and adversarially verify every rating-binding finding.

The selected effort and the resulting requirement set **MUST** be recorded in the
run's `plan.md` so the breadth actually applied is decidable from the artifacts.

### Evidence rigor

Any claim the skill makes about code, CLI, or tool behavior **MUST** be verified
by an actually executed command or search, and the assessment record **MUST**
cite the command or locator that verified it. The skill **MUST NOT** assert such
behavior from memory.

Every finding's locator **MUST** be pinned to a verifiable position — a
`file:line` or an exact searchable string — not a recalled or approximate one.

### Rating-binding re-check

The one or two findings that determine the headline (top-line) rating **MUST** be
independently re-checked before they drive `report.md` / `report.json`. The
re-check **MUST** re-run the verifying command or search rather than reuse the
first observation, and the report **MUST NOT** assert a headline rating whose
binding finding failed re-check.

### Recommendation actionability

Recommendation records **MUST** remain independently triageable without the
conversation in front of the reader. When the affected package, path, workflow,
maintainer surface, or verification route is inferable from the evidence, the
skill **SHOULD** include that route hint in the recommendation's existing text
fields rather than require a new schema field.

### Execution efficiency

The skill **SHOULD** compute all judgments first, then emit artifacts. Independent
artifact writes (for example, per-requirement assessment records) **MUST** be
emitted in parallel or batched rather than one serial write per round trip.
Batching **MUST NOT** change artifact content or the required layout.

### Deep subagent fan-out

At `deep` effort only, and only when the in-scope work justifies it, the skill
**MAY** fan out per-requirement or per-target assessment to subagents that return
structured findings. Roll-up judgment — aggregate, factor, and headline ratings —
**MUST** stay with the orchestrating skill. Subagent-returned evidence **MUST**
meet the same [evidence rigor](#evidence-rigor) requirements.
