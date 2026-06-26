---
type: Design Doc
title: Closed-choice setup UX — design doc
description: How /quality setup prompt guidance is updated for numbered, recommended-first closed choices.
tags: [skill, quality, setup, ux]
timestamp: 2026-06-26T00:00:00Z
---

# Closed-choice setup UX — design doc

Design behind the [Closed-choice setup UX](../0099-closed-choice-setup-ux.md)
change case and its [functional spec](spec.md).

## Context

The defect is in the agent-facing prompt contract, not in a parser or model
schema. Setup already infers defaults and confidence labels, but the runtime
prompt template lets a closed-choice question render as separate `Options`,
`Recommended`, and `Answer` lines. In practice, that made the recommended answer
harder than it should be and exposed internal vocabulary (`low tolerance`) in a
question framed around user-visible cost.

## Approach

Update the guidance at both durable and runtime layers:

1. The shared agent-mediated UX guide carries the general design rule and an
   example: numbered closed-choice options, recommended first, `1` as the
   shortest accept path, and option labels that match the question's axis.
2. The durable `/quality` skill parent spec and runtime `SKILL.md` repeat the
   general interaction contract so all workflows inherit the rule even when a
   contributor reads only the skill prompt.
3. The durable setup workflow spec and runtime setup workflow apply the rule to
   the concrete discovery questions. Lifecycle and Rating Scale keep their
   existing option vocabularies but are rendered in recommended-first order.
   Risk tolerance keeps the internal field and setup-brief meaning but presents
   cost labels to the user, mapping high cost to low tolerance, moderate cost to
   moderate tolerance, and low cost to high tolerance.

This is a documentation/spec/skill change. There is no Go code path to update
because setup discovery is agent-mediated, not CLI-driven.

## Alternatives

- **Only update the general UX guide.** Rejected: the actual setup runtime prompt
  would still contain the old `Options`/`Recommended` pattern, and future skill
  runs would not reliably inherit the guide's newer rule.
- **Rename the internal setup value from risk tolerance to cost.** Rejected: the
  model concept is still about tolerance and requirement strictness. The problem
  is the user-facing answer axis, not the stored concept.
- **Accept both numbered and word confirmations without requiring `1` as the
  default.** Rejected as too weak: accepting numbers is useful, but the important
  design constraint is that the recommended path is consistently the first and
  easiest path.

## Trade-offs & risks

- Reordering options means the visible option order can vary by inferred
  recommendation. That is deliberate; these are setup questions, not stable CLI
  enum displays.
- Mapping cost labels to tolerance values adds a small translation step for the
  agent. The runtime guidance names the mapping explicitly to avoid ambiguity.

## Open questions

None outstanding.
