---
type: Functional Specification
title: QUALITY.md getting-started guide
description: Contract for the skill's first-run guide for turning an initialized QUALITY.md skeleton into a first useful model.
tags: [skill, quality, guide]
timestamp: 2026-06-19T00:00:00Z
---

# QUALITY.md getting-started guide

This spec governs the **getting-started guide** the [`/quality` skill](../quality-skill.md)
ships at
[`skills/quality/guides/getting-started.md`](../../../../skills/quality/guides/getting-started.md).
The guide is the setup follow-on: the document a human or agent reads after
`qualitymd init` has created a valid skeleton and before the first real
evaluation model exists. It assumes the reader has already read the
[authoring guide](authoring.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../../docs/reference/rfc8174.md) when, and only when, they
appear in all capitals.

## Motivation

`qualitymd init` can create a valid file, but validity is only the first gate.
The initial authoring job needs process more than reference: the reader should
already have the authoring guide's best practices, then use a focused sequence
to turn the skeleton into a first useful model. Keeping that job in its own guide
lets setup hand off to a playbook without making the authoring reference carry
first-run workflow.

## Purpose

The guide exists to help a reader turn an initialized QUALITY.md skeleton into
a first useful model through an ordered process. After following it, the file
should still lint cleanly and each step should have reached a stated outcome, so
`/quality wizard` can route a meaningful next workflow.

## Scope

In scope: first-run population after `qualitymd init`, including ordered steps,
desired outcomes for each step, lightweight completion checks, links to the
authoring guide for best-practice detail, validation, and routing to the next
skill workflow.

Non-goals: the guide does not restate the full format reference, carry durable
authoring best-practice rationale, teach the full evaluation workflow, prescribe
domain-specific model templates, or promise that the first model is complete.
It should link to the authoring guide for concept reference and best-practice
detail.

## Requirements

### Runtime Use

The skill root prompt **MUST** tell agents to read the getting-started guide
after setup creates an initial `QUALITY.md`, or when a user asks how to make the
first useful model from a skeleton.

The getting-started guide **MUST** require the authoring guide as a prerequisite.
The skill should already have read the authoring guide before using
getting-started for first-run model population.

Setup mode **MUST** route successful initialization to the getting-started guide
before sending the user to wizard or evaluation.

Wizard mode **SHOULD** prefer the getting-started guide over the broad authoring
guide when the user has just initialized a skeleton or asks how to start from
one.

### Guide Content

The guide **MUST** identify its starting point as a valid scaffold produced by
`qualitymd init`, and **MUST** tell the reader to resolve lint errors before
building the model.

The guide **MUST** cover these first-pass jobs:

- filling the recommended Markdown body stubs;
- confirming or lightly adapting the seeded rating scale;
- naming the root subject from the body context;
- deciding whether the root `source` can stay implicit;
- selecting a small first set of Factors;
- replacing placeholder Requirements with assessable Requirements;
- running validation/status commands before the next skill workflow.

Each first-pass job **MUST** include a desired outcome and should link to the
relevant authoring-guide section for best-practice detail. The guide may include
brief completion checks, but it should not duplicate the full rationale or
concept guidance from the authoring guide.

The guide **MUST** preserve this process order: Markdown body first, rating
scale second, then subject/source alignment, Factors, Requirements, validation,
and wizard routing.

The guide **MUST** distinguish structural validity from model usefulness. A
valid skeleton is not enough for evaluation; the guide should ask whether the
body context is specific, current, grounded in agent-accessible support where
support is material, and sufficient to evaluate model quality; whether
Requirements are assessable; and whether evidence can distinguish important
rating levels.

The guide **MUST** end by routing to `/quality wizard` as the normal next
workflow once the first model is valid and useful enough.

### Relationship to Authoring Guide

The getting-started guide **MUST** depend on and link to the authoring guide for
format reference, concept guidance, and authoring best practices. It should stay
procedural and first-run focused rather than duplicating the comprehensive
reference.
