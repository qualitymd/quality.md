---
type: Functional Specification
title: Top 10 QUALITY.md checks
description: Contract for the skill's quick QUALITY.md inspection checklist used to produce routing findings for read-only orientation and model review.
tags: [skill, quality, guide, checklist]
timestamp: 2026-06-23T00:00:00Z
---

# Top 10 QUALITY.md checks

This spec governs the **Top 10 QUALITY.md checks** guide the [`/quality` skill](../quality-skill.md)
ships at
[`skills/quality/guides/top-10-quality-md-checks.md`](../../../../skills/quality/guides/top-10-quality-md-checks.md).
The guide is a bounded, fast inspection checklist for the `QUALITY.md` file
itself. It produces findings about model state, model usefulness, and lifecycle
routing; it does not evaluate or rate the root area.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../../docs/reference/rfc8174.md) when, and only when, they
appear in all capitals.

## Motivation

Read-only orientation needs more than raw status counts to recommend the right
lifecycle step, but it should not become a full model audit or quality
evaluation. A short, shared checklist gives orientation and model-review
workflows a consistent way to inspect the current `QUALITY.md`, surface
actionable findings, and route to setup, getting-started, authoring,
evaluation, recommendation follow-up, or history work. The checklist also keeps
durable setup assumptions visible over time: project posture, stakeholder
needs, agent/collaboration fit, missing context, and quality-loop expectations.

## Purpose

The guide exists to quickly assess the current state, quality, and lifecycle of
a QUALITY.md file. Its output is a small set of routing findings that explain
why the next workflow should be setup, getting-started, authoring/model review,
evaluation, recommendation follow-up, history/reconciliation, or update.

The guide does not re-run setup. It checks whether the current `QUALITY.md`
still preserves the setup assumptions and model qualities needed for useful
evaluation, authoring, and maintenance.

## Scope

In scope: read-only inspection of `qualitymd status --json`, the area
`QUALITY.md` file, and evaluation-history signals summarized by status JSON.
The checklist covers model lifecycle state, project posture, stakeholder and
needs coverage, agent and collaboration fit, Markdown body context and missing
context, root area/scope alignment, rating-scale fit, area/factor shape,
requirement and assessment quality, and quality-loop maintenance signals.

Non-goals: the checklist does not inspect root area source files, produce
evaluation artifacts, rate the root area, fully audit every requirement, or
replace the authoring and getting-started guides. It produces routing findings,
not an Evaluation Report.

## Requirements

### Runtime Use

The skill root prompt **MUST** tell agents to read the checklist when quickly
inspecting a QUALITY.md file's current state, quality, or lifecycle.

Read-only orientation **MUST** use the checklist after status probing when a
`QUALITY.md` exists and is structurally valid, unless the user asked only for
raw status. Orientation may skip checklist inspection when the model is missing,
the model is invalid, or CLI support is missing/stale enough that routing is
already decided.

Other modes **MAY** use checklist findings as context when they need to explain
why model authoring, evaluation, recommendation follow-up, or
history/reconciliation is the next workflow.

### Inspection Boundary

The checklist **MUST** stay read-only. It **MUST NOT** edit `QUALITY.md`, inspect
root area source files, read evaluation report bodies, create evaluation records,
or rate the root area.

The checklist should use status JSON for mechanical signals and read the
`QUALITY.md` file only for bounded model-usefulness inspection. It should not
perform an exhaustive audit of every model node.

The checklist **MUST NOT** require lifecycle, risk tolerance, modeling rigor,
collaboration context, stakeholder needs, or quality-loop posture to appear in
fixed sections. It should treat them as present when they are explicit, current,
and usable anywhere in the Markdown body or model context.

### Finding shape

The checklist **MUST** define a concise finding shape containing:

- check id;
- finding;
- evidence;
- impact; and
- route.

Findings **MUST** be routing-oriented. Evidence should cite status fields,
section names, property paths, counts, or short locators rather than long
quotations.

Routes should use skill workflow language such as setup, getting-started,
authoring, evaluate, recommendation follow-up, history, or update.

The checklist **MUST** treat the authoring guide as the quality reference for
what good authoring looks like. It should route starter or placeholder models to
getting-started for first-run process, and route populated models with
best-practice gaps to authoring/model review.

### Required Checks

The checklist **MUST** contain ten checks:

1. model lifecycle state;
2. project posture;
3. stakeholder and needs coverage;
4. agent and collaboration fit;
5. body context and missing context;
6. root area and scope alignment;
7. rating scale fit;
8. area and factor shape;
9. requirement and assessment quality; and
10. quality-loop maintenance signals.

The model-lifecycle-state check **MUST** use `qualitymd status --json` to
identify whether the model is missing, invalid, valid with no history, valid
with history, or needs evaluation reconciliation.

The project-posture check **MUST** inspect whether the model captures the
project context that calibrates the quality bar: lifecycle, risk tolerance, and
intended modeling rigor.

The stakeholder-and-needs-coverage check **MUST** inspect whether primary user
needs, collaborator/maintainer needs, and other affected stakeholder needs are
visible enough to justify the model's factors and requirements.

The agent-and-collaboration-fit check **MUST** inspect whether the model supports
the assumed agent-heavy workflow plus the named human collaboration context.

The body-context-and-missing-context check **MUST** inspect the recommended
Markdown body sections: Overview, Scope, Needs, and Risks, along with each
section's unknowns and open questions. Those unknowns and open questions **MUST**
be treated as author-declared context, distinct from a `not assessed` evaluation
result. The check **SHOULD** flag missing or non-agent-accessible support as a
model-usefulness finding when it prevents a reader or agent from evaluating
whether the body context is complete, current, grounded, or sufficient.

The root-area-and-scope-alignment check **MUST** inspect whether the root title,
body scope, file location, and root or child `source` values describe the same
evaluated root area. It should treat the current directory as the default root
area convention unless the model clearly narrows or relocates scope.

The rating-scale-fit check **MUST** inspect whether the rating scale is
understandable and fits the body's decision context, including lifecycle, risk
tolerance, and modeling rigor.

The area-and-factor-shape check **MUST** inspect whether the area tree is small
enough to understand, specific enough to represent distinct evaluated entities,
and shaped by the body's needs and risks. It **MUST** flag a composite entity
flattened into a single primary-subject root — distinct constituent artifacts of
different kinds described in the body, but all factors held at the root as one
family — and **SHOULD** flag a missing expected use-context constituent (an
owned, high-leverage agent harness or QUALITY.md self-check not modeled as a
constituent). It **MUST** treat such constituents as earned expected defaults of
the context of use, not a required roster, and **MUST NOT** flag a harness-less
or throwaway project for their absence.

The requirement-and-assessment-quality check **MUST** inspect whether
requirements are concrete enough to produce findings and ratings, and whether
each `assessment` gives the evaluator a usable means of assessment.

The quality-loop-maintenance-signals check **MUST** inspect evaluation history,
active recommendations, and visible model context to decide whether the next
workflow is maintenance rather than new authoring or evaluation. It **MUST NOT**
require or recommend CI or release gating by default.

Every check that would block evaluation **MUST** distinguish model
usefulness from root area quality. A valid but vague model is a model-authoring
finding, not evidence that the root area is low quality.

### Maturity vs lifecycle state

The checklist **MUST** keep two axes distinct and **MUST NOT** present them as one
blended classification:

- *Lifecycle state* — where the model sits in the evaluation lifecycle — is owned
  by `qualitymd status` (`readiness`): missing, invalid, ready to evaluate (valid
  with no runs), has evaluation history, or needs reconciliation.
- *Maturity* — how developed the model is — is the checklist's own judgment:
  `starter`, `immature`, or `evaluation-ready`.

The checklist **MUST NOT** treat the CLI's lifecycle `ready-to-evaluate` signal
as a maturity verdict; a valid model with no runs can still be a `starter`.

### Condensed close checklist

The guide **MUST** include a condensed checklist that the setup workflow uses to
classify maturity at close without reading every check. The condensed checklist
**MUST** cover model validity, body context, project posture and stakeholder
needs, and factor/requirement quality, and **MUST** map its result to the
`starter`/`immature`/`evaluation-ready` maturity levels. The full checks remain
available for borderline maturity calls and complete read-only orientation.

### Summary Judgment

The checklist **MUST** end by reporting the two axes separately: a lifecycle
classification (missing, invalid, ready to evaluate, has evaluation history, or
needs reconciliation) and, for a valid model, a maturity classification
(`starter`, `immature`, or `evaluation-ready`).
