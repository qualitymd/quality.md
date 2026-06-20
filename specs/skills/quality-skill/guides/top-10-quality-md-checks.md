---
type: Functional Specification
title: Top 10 QUALITY.md checks
description: Contract for the skill's quick QUALITY.md inspection checklist used to produce routing findings for wizard and related modes.
tags: [skill, quality, guide, checklist]
timestamp: 2026-06-19T00:00:00Z
---

# Top 10 QUALITY.md checks

This spec governs the **Top 10 QUALITY.md checks** guide the [`/quality` skill](../quality-skill.md)
ships at
[`skills/quality/guides/top-10-quality-md-checks.md`](../../../../skills/quality/guides/top-10-quality-md-checks.md).
The guide is a bounded, fast inspection checklist for the `QUALITY.md` file
itself. It produces findings about model state, model usefulness, and lifecycle
routing; it does not evaluate or rate the subject.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../../docs/reference/rfc8174.md) when, and only when, they
appear in all capitals.

## Motivation

Wizard needs more than raw status counts to recommend the right lifecycle step,
but it should not become a full model audit or subject evaluation. A short,
shared checklist gives wizard and related modes a consistent way to inspect the
current `QUALITY.md`, surface actionable findings, and route to setup,
getting-started, authoring, evaluation, improvement, or history work.

## Purpose

The guide exists to quickly assess the current state, quality, and lifecycle of
a QUALITY.md file. Its output is a small set of routing findings that explain
why the next workflow should be setup, getting-started, authoring/model review,
evaluation, improvement, history/reconciliation, or upgrade.

## Scope

In scope: read-only inspection of `qualitymd status --json`, the target
`QUALITY.md` file, and evaluation-history signals summarized by status JSON.
The checklist covers lifecycle state, Markdown body context, rating-scale fit,
subject/source alignment, target shape, Factor coverage, Requirement
assessability, assessment evidence, evaluation readiness, and maintenance
signals.

Non-goals: the checklist does not inspect subject source files, produce
evaluation artifacts, rate the subject, fully audit every Requirement, or replace
the authoring and getting-started guides. It produces routing findings, not an
Evaluation report.

## Requirements

### Runtime Use

The skill root prompt **MUST** tell agents to read the checklist when quickly
inspecting a QUALITY.md file's current state, quality, or lifecycle.

Wizard mode **MUST** use the checklist after status probing when a `QUALITY.md`
exists and is structurally valid, unless the user asked only for raw status.
Wizard may skip checklist inspection when the model is missing, the model is
invalid, or CLI support is missing/stale enough that routing is already decided.

Other modes **MAY** use checklist findings as context when they need to explain
why model authoring, evaluation, improvement, or history/reconciliation is the
next workflow.

### Inspection Boundary

The checklist **MUST** stay read-only. It **MUST NOT** edit `QUALITY.md`, inspect
subject source files, read evaluation report bodies, create evaluation records,
or rate the subject.

The checklist should use status JSON for mechanical signals and read the
`QUALITY.md` file only for bounded model-usefulness inspection. It should not
perform an exhaustive audit of every model node.

### Finding Shape

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
authoring, evaluate, improve, history, or upgrade.

The checklist **MUST** treat the authoring guide as the quality reference for
what good authoring looks like. It should route starter or placeholder models to
getting-started for first-run process, and route populated models with
best-practice gaps to authoring/model review.

### Required Checks

The checklist **MUST** contain ten checks:

1. lifecycle state;
2. body context;
3. rating scale fit;
4. subject and source alignment;
5. target shape;
6. Factor coverage;
7. Requirement assessability;
8. assessment evidence;
9. evaluation readiness; and
10. maintenance signals.

The body-context check **MUST** inspect the recommended Markdown body sections:
Overview, Scope, Needs, Risks, and Known gaps. The Known gaps guidance **MUST**
include known unknowns: missing context, unresolved questions, and evidence gaps
that prevented the rest of the body from being fully identified or defined.

The evaluation-readiness check **MUST** distinguish model usefulness from subject
quality. A valid but vague model is a model-authoring finding, not evidence that
the subject is low quality.

### Summary Judgment

The checklist **MUST** end by mapping findings to a lifecycle classification.
The classification vocabulary should include missing, invalid, starter,
immature, ready to evaluate, has evaluation history, and needs reconciliation.
