---
title: "<the system, component, or artifact this model is about>"
ratingScale:
  # Each level carries a description and a criterion. The description fixes what
  # the level *means* across the whole model — its standing and intent — and is
  # never overridden. The criterion is the default rule for rating a
  # requirement's findings at that level; an individual requirement may replace
  # it under its own `ratings` (e.g. a measured threshold) without changing what
  # the level means. Outstanding, Target, and Minimum are all acceptable; only
  # Unacceptable falls below the floor. The title is the human label used in
  # reports; the default titles include lightweight emoji markers for scanning,
  # but the stable `level` values carry the machine meaning.
  - level: outstanding
    title: 🟢 Outstanding
    description: "The stretch band — reached only with significant extra effort."
    criterion: "Exceeds the requirement; satisfies it with margin to spare."
  - level: target
    title: 🔵 Target
    description: "The level to aim for — achievable at reasonable cost and effort."
    criterion: "Satisfies the requirement."
  - level: minimum
    title: 🟡 Minimum
    description: "The acceptable floor — less than you'd aim for, but consciously agreed as good enough to rely on."
    criterion: "Falls short of the target but remains acceptable."
  - level: unacceptable
    title: 🔴 Unacceptable
    description: "Below the floor — not good enough to rely on."
    criterion: "Does not meet the requirement to an acceptable degree."
factors:
  # factors here hang off the model root, so they describe the whole artifact —
  # the qualities that matter most across all of it. This is the usual starting
  # point. A factor names a quality characteristic. Illustrative examples from
  # different domains include reliability for a service, credibility for a
  # report, freshness for a data set, or fairness for a process. Say what this
  # factor means here, why it matters and to whom, and how it differs from the
  # other factors.
  quality-name:
    title: "<short display label for this quality>"
    # description — what this quality means for the artifact, and to whom.
    # e.g. "Reliability: the service keeps serving correct responses under load
    #       and recovers cleanly from failure; the ops team and downstream
    #       callers depend on it."
    # e.g. "Credibility: the report's claims are traceable to current evidence
    #       so reviewers can trust its conclusions."
    description: "<what this quality means for the artifact, and to whom it matters>"
    requirements:
      # A requirement is one assessable expectation, stated as a claim an
      # evaluator can judge — e.g. "error messages give the user a clear path to
      # recovery" or "reported figures cite current source data".
      "<state one expectation you can assess>":
        # assessment — how an evaluator should inspect or measure the claim. It
        # can be stated inline, or reference an entity that already defines how to
        # check it — a spec, style guide, runbook, or test plan. An entity you own
        # and reference can be an area in its own right, which keeps the
        # dependency traceable.
        # Reference existing documentation (often the simplest):
        # e.g. "Conform to the error-handling rules in docs/errors.md."
        # e.g. "Run the load test in perf/loadtest.md and compare to its SLOs."
        # e.g. "Check the cited data source dates against the reporting window."
        # Or state it inline:
        # e.g. "Trigger the documented error cases and check each message names
        #       the cause and a next step."
        # e.g. "Hand the onboarding guide to someone unfamiliar with the process
        #       and see whether they can complete the first task unaided."
        assessment: "<how an evaluator should inspect or measure it>"

# areas are optional and narrower. Reach for one when a distinct part of the
# artifact — for example a service, module, document, data table, operating
# procedure, or curriculum unit — deserves its own factors or requirements that
# wouldn't fit cleanly at the top level. The whole-artifact factors above still
# apply; an area just adds focus where a part needs it. An area takes the same
# shape as the root: its own `factors` (and their `requirements`), direct
# `requirements` with factor references under `factors`, or further nested
# `areas`.
# areas:
#   "<name a thing to evaluate>":
#     title: "<short display label for this area>"
#     factors:
#       "<a quality specific to this part>":
#         title: "<short display label for this quality>"
#         description: "<what it means for this area, and to whom>"
#         requirements:
#           "<an expectation you can assess on this area>":
#             assessment: "<how an evaluator should inspect or measure it>"
---

# <the system, component, or artifact this model is about>

The frontmatter above fixes *what* is assessed and *how* it is rated. This body
records the judgment context a reader needs to trust the model, evaluate the
model's quality, and weigh future findings. Fill in each section and delete
these hints as you go.

Close each section with its unknowns, its open questions, and a state line. An
unknown is a broad area of uncertainty about that section's topic that may not
resolve to a single answer; an open question is a specific question with one
particular answer, still unresolved. Write "none known" rather than leaving
either out. If important supporting context is not agent-accessible — not in the
repo, cited paths, configured tools, linked public sources, or context provided
to the agent — record that limitation where it bears on the section. The state
line records who last stood behind the section: `Reviewed` names the person who
endorsed it, `agent-reviewed` the last agent pass — advance `Reviewed` only when
a person actually reviews it.

## Overview

What is this, who relies on it, and what does "good" look like for them? Quality
is value to the people who depend on the thing, so name them and the value they
expect.

*Unknowns* — <broad uncertainties about the root area, or "none known">
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
