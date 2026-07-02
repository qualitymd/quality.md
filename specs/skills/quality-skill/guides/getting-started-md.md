---
type: Functional Specification
title: QUALITY.md getting-started guide
description: Contract for the skill's first-run guide for turning a QUALITY.md with important model gaps into a first useful model.
tags: [skill, quality, guide]
timestamp: 2026-06-23T00:00:00Z
---

# QUALITY.md getting-started guide

This spec governs the **getting-started guide** the [`/quality` skill](../quality-skill.md)
ships at
[`skills/quality/guides/getting-started.md`](../../../../skills/quality/guides/getting-started.md).
The guide is the setup follow-on: the document a human or agent reads after
setup leaves a valid `QUALITY.md` with important model gaps, or when a user asks
how to keep iterating on the first useful model. It assumes the reader has
already read the [authoring guide](authoring.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../../docs/reference/rfc8174.md) when, and only when, they
appear in all capitals.

## Motivation

Setup can create a valid file, but validity is only the first gate. The initial
authoring job needs process more than reference: the reader should already have
the authoring guide's best practices, then use a focused sequence to turn the
current file into a first useful model. Keeping that job in its own guide lets
setup hand off to a playbook without making the authoring reference carry
first-run workflow.

## Purpose

The guide exists to help a reader turn a QUALITY.md with important model gaps
into a first useful model through an ordered process. After following it, the
file should still lint cleanly and each step should have reached a stated
outcome, so the skill can recommend a meaningful next public workflow.

## Scope

In scope: first-run iteration after setup, including ordered steps, desired
outcomes for each step, lightweight completion checks, links to the authoring
guide for best-practice detail, validation, and routing to the next skill
workflow.

Non-goals: the guide does not restate the full format reference, carry durable
authoring best-practice rationale, teach the full evaluation workflow, prescribe
domain-specific model templates, or promise that the first model is complete.
It should link to the authoring guide for concept reference and best-practice
detail.

## Requirements

### Runtime use

The skill root prompt **MUST** tell agents to read the getting-started guide
after setup leaves a valid `QUALITY.md` with important model gaps, or when a user
asks how to keep iterating on the first useful model.

The getting-started guide **MUST** require the authoring guide as a prerequisite.
The skill should already have read the authoring guide before using
getting-started for first-run model population.

The setup workflow **MAY** route models with important gaps to the
getting-started guide after its own context-informed setup work. It **MUST NOT**
treat the guide as a separate required phase before setup writes `QUALITY.md`.

Read-only orientation **SHOULD** prefer the getting-started guide over the broad
authoring guide when the user has just initialized a skeleton or asks how to
start from one.

### Guide content

The guide **MUST** identify its starting point as a valid `QUALITY.md` with
important model gaps, and **MUST** tell the reader to resolve lint errors before
building the model.

The guide **MUST** cover these first-pass jobs:

- filling the recommended Markdown body stubs;
- reviewing setup assumptions such as root area, domain, lifecycle, risk
  tolerance, modeling rigor, collaboration context, stakeholder needs, and missing
  context;
- confirming or lightly adapting the seeded rating scale;
- naming the root area from the body context;
- deciding whether the root `source` can stay implicit;
- selecting a small first set of factors;
- replacing placeholder requirements with assessable requirements;
- running validation/status commands before the next skill workflow.

Each first-pass job **MUST** include a desired outcome and should link to the
relevant authoring-guide section for best-practice detail. The guide may include
brief completion checks, but it should not duplicate the full rationale or
concept guidance from the authoring guide.

The guide **MUST** preserve this process order: Markdown body first, rating
scale second, then root area/source alignment, factors, requirements,
validation, and next-workflow guidance.

The guide **MUST** distinguish structural validity from model usefulness. A
valid skeleton is not enough for evaluation; the guide should ask whether the
body context is specific, current, grounded in agent-accessible support where
support is material, and sufficient to evaluate model quality; whether
requirements are assessable; and whether evidence can distinguish important
rating levels.

The guide **MUST** end by offering next workflow choices once the first model is
valid and useful enough: continue iterating on `QUALITY.md`, run
`/quality evaluate`, or stop. It should direct users to continue authoring when
model usefulness is not yet sufficient and to `/quality update` when tooling is
stale.
The next-workflow guidance **MUST** use a labeled recommendation, numbered
options, and an explicit `Answer` line.

### Relationship to authoring guide

The getting-started guide **MUST** depend on and link to the authoring guide for
format reference, concept guidance, and authoring best practices. It should stay
procedural and first-run focused rather than duplicating the comprehensive
reference.
