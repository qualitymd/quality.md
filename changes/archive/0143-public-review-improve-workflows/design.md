---
type: Design Doc
title: Public Review and Improve Workflows - design
description: How to expose review and improve as public /quality workflows while keeping the initial implementation as focus-routed stubs.
tags: [skill, workflows, docs, ux]
timestamp: 2026-06-27T00:00:00Z
---

# Public Review and Improve Workflows - design

Companion to the
[Public Review and Improve Workflows](../0143-public-review-improve-workflows.md)
change case and its [functional spec](spec.md).

## Context

The change makes `/quality review` and `/quality improve` public immediately,
without pretending their deep workflow designs are complete. The implementation
is a skill/docs change: update the root skill router, add runtime and durable
workflow stubs, and simplify the README around evaluate, review, and improve.

The central design constraint is that `evaluate` already has a bounded contract.
`review` and `improve` should route users into useful next steps quickly while
preserving the existing mutation boundaries for direct model authoring and
recommendation follow-up.

## Approach

Add `review` and `improve` beside `setup`, `evaluate`, and `update` in the
public workflow surface.

Use **focus** as the routing field both workflows expose to users:

```text
review focus:  evaluation | model | concern
improve focus: evaluation | model | concern | recommendation
```

The runtime workflow files are intentionally short operator playbooks:

- `skills/quality/workflows/review.md` opens with a run frame, infers or asks for
  focus, stays read-only, and either performs a shallow review using existing
  guidance or stops with a clear preview of the deeper review path.
- `skills/quality/workflows/improve.md` opens with a run frame, infers focus and
  mutation surface, confirms both before mutation, then delegates to existing
  direct model-authoring or recommendation-follow-up guidance when applicable.

The durable workflow specs mirror those runtime files:

- `specs/skills/quality-skill/workflows/review.md`
- `specs/skills/quality-skill/workflows/improve.md`

The root durable skill spec and runtime `SKILL.md` own shared dispatch:

- public invocation list includes `review` and `improve`;
- workflow dispatch reads the matching workflow file;
- bare `/quality` orientation may recommend review/improve when they fit;
- `setup` and `update` are framed as support workflows;
- `evaluate` remains the only Evaluation-record/report writer.

`README.md` is reorganized around the user-facing promise:

```text
Evaluate -> Review -> Improve
```

Setup becomes the way to create or update the model before the loop. Update
stays maintenance for skill/CLI compatibility. Specialized sections such as
quality debt and Agent Harnessability can remain, but they move behind the
primary workflow explanation.

## Spec response

The public framing requirements are satisfied by changing the README hierarchy:
lead with the three verbs, show the five command examples, and explain setup and
update as support.

The invocation and routing requirements are satisfied by adding `review` and
`improve` to `SKILL.md`, the durable root skill spec, and both workflow indexes.
The root prompt remains the always-loaded router; workflow-specific behavior
stays in dispatched files.

The review-stub requirements are satisfied by making review read-only by
default, using a first-output run frame with `Focus`, and limiting initial
behavior to shallow inspection plus next-action recommendation.

The improve-stub requirements are satisfied by treating improve as read-only
until focus and mutation surface are confirmed. Existing gates remain
authoritative: model changes use direct authoring review gates, and older
recommendation artifacts use recommendation follow-up.

The boundary requirements are satisfied by not adding CLI behavior, not changing
Evaluation data, and not adding compatibility aliases or shims.

## Alternatives

### Keep review/improve private until fully designed

Rejected. The public value proposition is already clearer than the current
workflow list, and the skill can safely expose shallow workflows if they confirm
focus and stop at explicit boundaries.

### Add many specific subcommands

Examples: `/quality review model`, `/quality review evaluation`,
`/quality improve recommendation`, `/quality improve model`.

Rejected as the primary interface. Users should be able to invoke
`/quality review` or `/quality improve` and be guided quickly. Specific phrases
can still be accepted as focus hints.

### Rebrand recommendation follow-up as the improve workflow wholesale

Rejected. Recommendation follow-up is one improve route, not the whole workflow.
`improve` also needs to cover model improvement and ad hoc quality concerns.

### Make review write a persisted review artifact

Deferred. A persisted review artifact may be useful later, but the first public
workflow should stay read-only and avoid new artifact contracts.

## Trade-offs and risks

The main trade-off is public surface ahead of deep behavior. That is acceptable
only if the stubs are honest: they must confirm focus, name what they will and
will not do, and stop rather than inventing unsupported analysis or mutation.

There is a terminology risk around historical guidance that says "review model"
or "recommendation follow-up" is not public. Implementation needs a careful
search so those statements are rewritten without losing the older safety
boundaries they carried.

There is also a recommendation risk: current `evaluate` does not generate new
recommendations in the active protocol. `improve` must not imply that every
Evaluation result has recommendations to apply. It should use findings,
candidate actions, compatible historical recommendation artifacts, or ask for a
different focus.

## Open questions

- Should `improve focus: evaluation` prefer findings/candidate actions from the
  latest run, or should it ask the user to select a specific finding first?
- Should `/quality review` create a workflow feedback log later, or remain
  entirely read-only with no local artifact?
- Once deeper workflow design begins, should review and improve share one focus
  resolver module in prose/spec, or keep separate workflow-specific routing
  tables?
