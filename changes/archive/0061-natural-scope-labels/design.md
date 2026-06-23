---
type: Design Doc
title: Natural scope labels - design doc
description: How the /quality skill resolves natural scoped-evaluation labels while preserving stable model IDs.
tags: [skill, evaluation, ux, references]
timestamp: 2026-06-23T00:00:00Z
---

# Natural scope labels - design doc

## Context

This design answers the
[Natural scope labels functional spec](spec.md). The change makes natural Area
and Factor labels the documented happy path for scoped `/quality evaluate`
requests, while stable model references and record identifiers stay unchanged.

The design goal is to move technical reference syntax out of the user's normal
path without weakening the model-reference boundary established by 0058, 0059,
and 0060. Natural labels are an agent input convenience; they are not a new
durable identifier form.

## Approach

Keep scope resolution inside the `/quality` skill. The CLI and format
specification continue to operate on stable model IDs and qualified references;
the skill performs the human-edge interpretation before invoking CLI-backed
evaluation workflows.

Update the runtime skill instructions in two places:

- `SKILL.md` owns the argument contract and examples. It should describe
  one-label and two-label scoped evaluation as normal usage, then mention
  qualified references as exact-addressing syntax.
- `modes/evaluate.md` owns the operational decision tree. It should resolve
  qualified references first, then natural labels, and stop to clarify only when
  the grounded model has multiple plausible candidates.

Resolution should be deterministic once the model is inspected:

```text
/quality evaluate <label>
  -> match the label against Area titles/names and Factor titles/names
  -> if exactly one Area matches, evaluate that Area
  -> if exactly one Factor matches, evaluate that Factor in its declaring Area
  -> if repeated Factors match, ask for the Area
  -> if mixed Area/Factor matches remain, ask which kind or which candidate

/quality evaluate <area-label> <factor-label>
  -> resolve the Area label first
  -> resolve the Factor label within that Area
  -> evaluate that Factor for the resolved Area
```

Use exact matching only for this change. Matching should consider the required
`title` values because those are the primary human-facing labels, and also the
stable YAML names because users may type what they see in model keys or prior
technical output. Case-folding and punctuation normalization can be allowed by
the skill if useful, but the durable contract should not depend on fuzzy or
semantic matching.

Clarification prompts should be shaped around the missing decision. For the
common repeated-Factor case, ask:

```text
What area do you want to evaluate Reliability for?
```

Then list Area titles or names first:

```text
1. API
2. Webhooks
3. CLI
```

Qualified references can be secondary context when the labels alone are still
hard to distinguish:

```text
1. API (`area:api`)
2. Webhooks (`area:webhooks`)
3. CLI (`area:cli`)
```

For mixed Area/Factor ambiguity, use a broader clarification first, then fall
back to the repeated-Factor Area question if the user chooses the Factor path.

Update user-facing docs to lead with natural examples:

```text
/quality evaluate Security
/quality evaluate Payments Maintainability
/quality evaluate API Reliability
```

Move `area:...` and `factor:...` examples into exact-addressing prose. They
remain important for unambiguous reproduction, support, and advanced workflows,
but should no longer define the first-read interaction model.

## Alternatives

**Keep qualified references as the documented primary form.** Rejected because
it exposes implementation vocabulary at the user edge. Qualified references are
still needed, but they are not the natural way users name quality concerns.

**Add CLI support for natural scope labels now.** Rejected because the skill
already owns judgment, model inspection, and clarification. Moving label
resolution into the CLI would require a deterministic command UX for ambiguity
that is not needed for this skill-first workflow.

**Persist the user's natural label in evaluation records for readability.**
Rejected because records and `report.json` need stable identifiers. Human
reports can resolve titles and names from the model snapshot without making
natural labels part of the durable machine contract.

**Use fuzzy or semantic matching.** Rejected for this change because it makes
surprising matches harder to reason about. Exact title/name matching plus clear
clarification handles the primary ergonomics problem without adding hidden
guesswork.

**Always ask whether a single label means Area or Factor.** Rejected because it
would make the common unique case noisier. Ask only when the grounded model
contains real ambiguity.

## Trade-offs & Risks

Natural labels can collide. The skill should treat collisions as normal and ask
for the missing Area or kind, not as an error. The wording matters: asking "What
area do you want to evaluate Reliability for?" keeps the user in project
vocabulary instead of forcing them into model-reference syntax.

Title and name matching creates two possible labels for the same model element.
That is acceptable because both labels resolve to the same stable ID before any
artifact is written. When both title and name collide across different elements,
the skill should show options rather than pick one.

Documentation must keep a clear boundary between natural labels and durable
references. If examples only show natural labels, advanced users may not know how
to reproduce an exact scope. Keep a short exact-addressing section after the
natural examples.

Because this is skill-first behavior, there may be little or no Go code to
change. The implementation review should still search CLI examples and generated
help for stale primary qualified-reference examples, but should avoid expanding
the scope into CLI parsing unless a hard blocker appears.

## Open Questions

- Should matching be case-sensitive in the written skill contract, or should
  the skill explicitly normalize case for user convenience?
- Should clarification options always include qualified references in
  parentheses, or only when duplicate human labels make the list hard to
  distinguish?
