---
type: Functional Specification
title: /quality recommendation follow-up
description: Post-evaluation follow-up workflow for applying or handing off /quality evaluation recommendations.
tags: [skill, quality, recommendation]
timestamp: 2026-06-22T00:00:00Z
---

# /quality recommendation follow-up

This spec owns the `/quality` skill's post-evaluation recommendation follow-up
route: how the skill helps users act on compatible evaluation recommendation
artifacts after `improve` resolves to recommendation focus. It composes the
shared contracts in the parent [/quality skill](quality-skill.md) spec and the
recommendation artifact contract in [/quality reporting](reporting.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear
in all capitals.

## Purpose and routing

Recommendation follow-up is selected by `/quality improve` when the user asks to
apply, act on, improve from, or hand off a compatible existing evaluation
recommendation artifact. It is not a separate public workflow and does not
replace `evaluate`: recommendations remain evaluation-related artifacts.

The workflow's purpose is to turn a selected recommendation into one of two
productive outcomes: a confirmed local apply or an issue-tracker handoff.

Before recommendation inspection, history inspection, outcome selection, local
apply, issue creation, quality changelog writes, or any other tool-dependent follow-up
work, the workflow **MUST** emit a concise follow-up frame. The frame **MUST**
name the recommendation or `resolving…`, the outcome if already requested or
`resolving…`, mutation surfaces, expected artifacts, and next gate. It
**MUST NOT** render a command-style header that implies recommendation follow-up
is a new public `/quality` invocation.

## Outcomes

Recommendation follow-up **MUST** offer only two explicit outcomes:

1. apply a confirmed recommendation option now;
2. hand off the recommendation to an issue tracker.

When the user has not already chosen one of those outcomes, the choice is a
single-select closed-choice intent rendered per the shared
[progressive-enhancement contract](quality-skill.md#user-interaction-contract):
through a fit-for-purpose native option picker when present, otherwise the text
fallback with numbered options and an `Answer` line. When recommendation evidence
supports one path, that path **SHOULD** be the recommended option (option `1` in
the text fallback); otherwise the workflow **MUST NOT** invent a recommendation
solely to satisfy the prompt shape.

The workflow **MUST NOT** present `defer`, `skip`, or `keep open` as formal
follow-up options. If the user does not choose apply or handoff, the workflow
**MUST** stop without mutating evaluated source, `QUALITY.md`, the quality changelog,
or external systems.

## Apply now

Before applying a recommendation, the skill **MUST** present a decision brief
that names the recommendation, selected option, artifact class being changed,
evidence or reason, risk when relevant, done criterion, and verification path.
The primary apply question or call to action **MUST** be visually emphasized,
and the decision brief **MUST** keep the changed artifacts, evidence/reason,
recommended option, alternatives, and done criterion in a consistent, scannable
shape. Because local apply is a true binary mutation gate after the
recommendation option and mutation surface are selected, the brief **MUST** show
`y`/`n` as the visible shortest answer path.

The skill **MUST NOT** edit evaluated source files, edit `QUALITY.md`, or write
the quality changelog until the user explicitly confirms the recommendation option and
mutation surface. It **MUST NOT** treat an obvious or recommended fix as consent.

After applying a confirmed option, the skill **SHOULD** verify the done criterion
with the narrowest useful evidence. When the done criterion is rating-bound or
depends on the QUALITY.md model, the skill **SHOULD** create a new numbered
evaluation run for the affected scope and compare the new evidence to the
recommendation's done criterion.

The result report **MUST** state the recommendation, outcome, applied option,
changed artifacts, verification performed, rating movement when known, remaining
gaps or limits, and next action. It **SHOULD** include `Not done` when the
mutation or verification boundary matters, such as no evaluation rerun, no issue
creation, no model change, or no quality changelog entry. It **MUST** use the shared
agent-mediated UX contract: status first, with a primary outcome line and a clear
next action when work remains. The block **MUST** use labeled fields for
recommendation, applied option, changed artifacts, verification, rating movement,
remaining gaps, not-done boundary, and next action. If verification is
incomplete, the result **MUST** be labeled limited rather than fully confirmed.

When a confirmed apply changes the QUALITY.md model, the skill **MUST** write one
quality changelog entry for the coherent model change, cross-linking the source
evaluation run and recommendation when applicable. Evaluated-source fixes that do
not change the model **MUST NOT** write the quality changelog.

Before applying a model-changing recommendation, the skill **MUST** read the
authoring entry guide, the routed authoring sub-guide for the model element being
changed, and the quality changelog authoring sub-guide.

## Issue-tracker handoff

Issue-tracker handoff **MUST** produce issue-ready content that includes the
recommendation number, recommendation ID, and title; source evaluation run;
affected area/factor or requirement; current rating when known; target or done
criterion; evidence summary with locators; suggested implementation option;
verification path; and links or paths to the generated report and recommendation
artifact.

Recommendation follow-up **MUST** treat numeric user input such as
`recommendation 1`, `rec #1`, and `1` as a recommendation-number selection from
the ranked recommendations report. It **MUST** treat `qrec_...` input as exact
recommendation ID selection.

Creating an external issue **MUST** require explicit user confirmation,
available issue-tracker tooling, and a decision brief that names the external
artifact to create, the local artifacts that will not change, the evidence or
reason, the recommended option, alternatives, and verification. Because external
issue creation is a true binary mutation gate, the text-fallback brief **MUST**
show `y`/`n` as the visible shortest answer path. Where the issue-tracker tooling
will itself prompt to authorize the creation, the skill **SHOULD** render the
confirmation through that native gate and keep the brief's teaching in the
preceding message rather than stacking a second text gate, per the shared
[no-double-gate rule](quality-skill.md#decision-briefs); this **MUST NOT** weaken
the confirmation requirement. If tooling is unavailable or the user has not
confirmed external creation, the skill **MUST** stop after producing issue-ready
text.

Issue-tracker handoff **MUST NOT** mutate evaluated source, `QUALITY.md`, or the
quality changelog.
