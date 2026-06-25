---
type: Design Doc
title: Setup workflow scope trim - design
description: Design for trimming setup by deleting future-workflow preference capture, replacing maturity/evaluation-ready closeout with important gaps, and keeping feedback logging and pedagogical discovery.
tags: [skill, setup, ux]
timestamp: 2026-06-25T00:00:00Z
---

# Setup workflow scope trim - design

Answers the [functional spec](spec.md). This is a guidance and durable-spec
mirror change: no CLI, Go, format-schema, rating, roll-up, or evaluation-record
behavior changes.

## Context

Setup has three jobs that directly produce a useful initial `QUALITY.md`:

1. Gather enough grounded context to author the first model.
2. Write or update the model within setup's narrow mutation boundary.
3. Validate the file and name important model gaps.

The current workflow also asks or reports on future workflow concerns:
recommendation handling, handoff destination, recurring review, automation
posture, and model maturity/evaluation-readiness. Those are useful later, but
they do not improve first-model authoring enough to justify their setup surface.

## Approach

### Delete future-workflow preference capture

Remove the optional work-handoff question entirely from runtime setup and its
durable spec. Also remove review cadence, recurring review, and automation
posture from setup discovery and first-run continuation guidance. Keep only the
boundary sentence that setup does not create issues, configure integrations, or
configure automations.

Recommendation follow-up keeps owning concrete apply/handoff decisions once
recommendations exist. Evaluation remains the path that produces recommendations.

### Replace maturity with important gaps

Setup closeout should stop saying `Maturity:` and stop producing
`starter`/`immature`/`evaluation-ready`. Instead:

- lint reports validity;
- important gaps report model-usefulness concerns;
- the next action is one of `continue iterating on QUALITY.md`, `run /quality
  evaluate`, or `stop here`.

The "important gaps" check reuses the same underlying concerns that made the old
maturity check useful — thin body context, unmodeled germane constituents, missing
Agent Harnessability coverage, vague requirements — but reports the actual gaps
instead of compressing them into a readiness label.

### Keep Top 10 as orientation, not setup closeout

Top 10 QUALITY.md checks remain useful for read-only orientation and model review,
but setup no longer depends on their condensed close checklist. The guide is
reframed to inspect lifecycle state and model usefulness, not to classify
`starter`/`immature`/`evaluation-ready`.

### Keep feedback logs, revise outcome

The setup feedback log remains mandatory because it improves the setup process.
Its setup-specific `outcome` field changes from a maturity result to workflow
outcome values:

- `completed`
- `completed-with-important-gaps`
- `lint-failed`
- `failed`
- `interrupted`

These values describe the setup run without implying evaluation readiness.

## Alternatives

- **Keep the work-handoff question and wire it into recommendation follow-up.**
  Rejected for this case. That would make setup remember future workflow policy
  and still require later confirmation. The simpler boundary is to ask when a real
  recommendation exists.
- **Keep maturity but rename `evaluation-ready`.** Rejected. Renaming preserves
  the extra classification layer. The desired closeout is validity plus concrete
  gaps.
- **Remove or reduce the core discovery questions.** Rejected as premature. The
  retained questions teach the dimensions of a quality model and guard against
  invented context.
- **Remove workflow feedback logs.** Rejected. The feedback log is not model
  content, but it is the mechanism for improving setup itself.

## Trade-offs & risks

- Setup may give a less compact closeout without the single maturity label.
  Mitigation: list only important gaps and recommend one immediate next action.
- Read-only orientation still needs a way to route vague or immature models. The
  Top 10 guide keeps route-oriented findings and model-usefulness language.
- Existing historical logs and archived cases still mention maturity/readiness.
  They remain past-state records and are not rewritten.

## Open questions

None.
