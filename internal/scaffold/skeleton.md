---
title: "<the system, component, or artifact this model is about>"
ratingScale:
  # Each level carries a description and a criterion. The description fixes what
  # the level *means* across the whole model — its standing and intent — and is
  # never overridden. The criterion is the default rule for rating a
  # requirement's findings at that level; an individual requirement may replace
  # it under its own `ratings` (e.g. a measured threshold) without changing what
  # the level means. Outstanding, Target, and Minimum are all acceptable; only
  # Unacceptable falls below the floor.
  - level: outstanding
    title: Outstanding
    description: "The stretch band — reached only with significant extra effort."
    criterion: "Exceeds the requirement; satisfies it with margin to spare."
  - level: target
    title: Target
    description: "The level to aim for — achievable at reasonable cost and effort."
    criterion: "Satisfies the requirement."
  - level: minimum
    title: Minimum
    description: "The acceptable floor — less than you'd aim for, but consciously agreed as good enough to ship."
    criterion: "Falls short of the target but remains acceptable."
  - level: unacceptable
    title: Unacceptable
    description: "Below the floor — not good enough to ship."
    criterion: "Does not meet the requirement to an acceptable degree."
factors:
  # Factors here hang off the model root, so they describe the whole artifact —
  # the qualities that matter most across all of it. This is the usual starting
  # point. A factor names a quality characteristic — e.g. reliability, security,
  # usability, maintainability. Say what it means here, why it matters and to
  # whom, and how it differs from the other factors.
  "<name a quality that matters>":
    title: "<short display label for this quality>"
    # description — what this quality means for the artifact, and to whom.
    # e.g. "Reliability: the system keeps serving correct responses under load
    #       and recovers cleanly from failure — the ops team and every
    #       downstream caller depend on it."
    # e.g. "Usability: a first-time developer can integrate with it without
    #       asking the team for help."
    description: "<what this quality means for the artifact, and to whom it matters>"
    requirements:
      # A requirement is one assessable expectation, stated as a claim an
      # evaluator can judge — e.g. "error messages give the user a clear path to
      # recovery" or "the public API can be adopted without reading the source".
      "<state one expectation you can assess>":
        # assessment — how an evaluator should inspect or measure the claim. It
        # can be stated inline, or reference an entity that already defines how to
        # check it — a spec, style guide, runbook, or test plan. An entity you own
        # and reference can be a target in its own right, which keeps the
        # dependency traceable.
        # Reference existing documentation (often the simplest):
        # e.g. "Conform to the error-handling rules in docs/errors.md."
        # e.g. "Run the load test in perf/loadtest.md and compare to its SLOs."
        # Or state it inline:
        # e.g. "Trigger the documented error cases and check each message names
        #       the cause and a next step."
        # e.g. "Hand the public API docs to someone unfamiliar with it and see
        #       whether they can complete the quickstart unaided."
        assessment: "<how an evaluator should inspect or measure it>"

# Targets are optional and narrower. Reach for one when a distinct part of the
# artifact — a service, a module, a document (e.g. "checkout-api",
# "auth-service", "design-system") — deserves its own factors or requirements
# that wouldn't fit cleanly at the top level. The whole-artifact factors above
# still apply; a target just adds focus where a part needs it. A target takes
# the same shape as the root: its own `factors` (and their `requirements`),
# direct `requirements` with factor references under `factors`, or further
# nested `targets`.
# targets:
#   "<name a thing to evaluate>":
#     title: "<short display label for this target>"
#     factors:
#       "<a quality specific to this part>":
#         title: "<short display label for this quality>"
#         description: "<what it means for this target, and to whom>"
#         requirements:
#           "<an expectation you can assess on this target>":
#             assessment: "<how an evaluator should inspect or measure it>"
---

# <the system, component, or artifact this model is about>

The frontmatter above fixes *what* is assessed and *how* it is rated. This body
explains *why* — the context a reader needs to trust the model, and an evaluator
needs to weigh it. Fill in each section and delete these hints as you go.

Close each section with its unknowns, its open questions, and a state line. An
unknown is a broad area of uncertainty about that section's topic that may not
resolve to a single answer; an open question is a specific question with one
particular answer, still unresolved. Write "none known" rather than leaving
either out. The state line records who last stood behind the section: `Reviewed`
names the person who endorsed it, `agent-reviewed` the last agent pass — advance
`Reviewed` only when a person actually reviews it.

## Overview

What is this, who relies on it, and what does "good" look like for them? Quality
is value to the people who depend on the thing, so name them and the value they
expect.

*Unknowns* — <broad uncertainties about the subject, or "none known">
*Open questions* — <specific unresolved questions, or "none">

*Reviewed — <name>, <date>; agent-reviewed — <agent>, <date>*

## Scope

Draw the boundary: what this model covers, and what it deliberately leaves out.
What is left out is an exclusion by design, not a failing.

*Unknowns* — <broad uncertainties about the boundary, or "none known">
*Open questions* — <specific unresolved questions, or "none">

*Reviewed — <name>, <date>; agent-reviewed — <agent>, <date>*

## Needs

The outcomes your stakeholders are counting on. Each requirement above should
answer to a need here — this is what tells an evaluator how much each one
matters.

*Unknowns* — <broad uncertainties about the needs, or "none known">
*Open questions* — <specific unresolved questions, or "none">

*Reviewed — <name>, <date>; agent-reviewed — <agent>, <date>*

## Risks

What goes wrong, and for whom, when a need goes unmet? Naming the worst outcomes
is what keeps a rating meaningful rather than mechanical.

*Unknowns* — <broad uncertainties about the risks, or "none known">
*Open questions* — <specific unresolved questions, or "none">

*Reviewed — <name>, <date>; agent-reviewed — <agent>, <date>*
