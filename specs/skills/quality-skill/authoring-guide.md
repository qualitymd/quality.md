---
type: Functional Specification
title: QUALITY.md authoring guide
description: Contract for the skill's authoring guide resource — a single comprehensive guide to understanding and working with QUALITY.md files.
tags: [skill, quality, guide]
timestamp: 2026-06-18T00:00:00Z
---

# QUALITY.md authoring guide

This spec governs the **authoring guide** resource the [`/quality` skill](quality-skill.md) ships at
[`skills/quality/resources/quality-md-guide.md`](../../../skills/quality/resources/quality-md-guide.md)
— the document the skill reads when creating, populating, reviewing, or
improving a `QUALITY.md` file. It has 1:1 coverage with that document: this spec
is its contract, and the guide is its implementation. The format the guide
teaches is defined by [`SPECIFICATION.md`](../../../SPECIFICATION.md), the source
of truth the guide **conforms to**.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", and
"MAY" are to be interpreted as described in
[RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear
in all capitals.

## Purpose

The guide exists to be a **single, comprehensive guide to understanding and
working with `QUALITY.md` files**: one document a reader can stay inside to learn
what each concept of the format *is* and to do the jobs of authoring a model. It
serves both the skill (its reader at runtime) and a human author. The guide
**MUST** state this purpose in its own preamble, so the document declares the job
it is built to do.

## Scope

The defining scope decision is **self-containedness**. The guide is *all-inclusive*:
it restates the format's concepts, properties, and rating vocabulary in full
rather than deferring to [`SPECIFICATION.md`](../../../SPECIFICATION.md) for the
definitions. A reader **MUST** be able to understand and author a `QUALITY.md`
from the guide alone, without opening the format spec.

This is a deliberate exception to the skill's general "don't embed a drifting
copy of the format" rule, which binds the skill's *metadata and prompt* — those
are grounded at runtime from `qualitymd spec` (see
[Invocation](quality-skill.md#frontmatter-and-metadata)). The guide is a bundled
**reference resource**, not the prompt, and a self-contained reference is the
whole point of this one: the bundled [`SPECIFICATION.md`](../../../SPECIFICATION.md)
copy is the same kind of artifact. The duplication is paid for by the conformance
duty below, which makes the format spec the authority whenever the two disagree.

> Rationale: a guide a reader must cross-reference against the spec to use is not
> a single comprehensive guide. Self-containedness is the purpose; the
> conformance duty is the price that keeps it honest.

**In scope:** authoring a `QUALITY.md` — its file shape, the model concepts
(rating scale, targets, factors, requirements), and the Markdown body.

**Non-goals:** the guide does **not** document the *evaluation* process (Define →
Assess and Rate → Analyze → Advise → Report) beyond what bears on authoring; that
is the skill's own workflow (see [Evaluation workflow](quality-skill.md#evaluation-workflow)). It does **not** restate the CLI
surface, and it is **not** a normative spec — `SPECIFICATION.md` and the durable
specs remain the contracts; the guide is instructional.

## Requirements

### Conformance

The guide **MUST conform to** [`SPECIFICATION.md`](../../../SPECIFICATION.md):
every concept definition, property, presence level (Required / Recommended /
Optional), and rating-vocabulary term it states **MUST** match the format spec.
The guide **MUST** teach that `title` is required on the Model, every Target,
every Factor, and every Rating Level, and that Requirements do not have a
separate `title` because the requirement statement is their display text.
Conformance is the binding relationship, not deference — the guide phrases the
format in its own instructional voice rather than quoting it — but where the
guide and the format spec diverge, **the format spec governs and the guide MUST
be corrected to conform**. When `SPECIFICATION.md` changes, the guide **MUST** be
reconciled to it.

The guide **MUST** carry a one-line conformance note in its preamble stating that
it restates the format defined in `SPECIFICATION.md` and conforms to it, with the
spec governing on any conflict — so a reader and a future editor both know the
authority.

### Structure

The guide is organized for a reader who is authoring, mixing **reference** (what
a concept is) with **how-to** (the jobs of working with it):

- The guide **MUST** use a **single level** of concept hierarchy: each top-level
  `QUALITY.md` concept (the file, the model, rating scale, target, factor,
  requirement, the body) is its own chapter. A concept's sub-concepts and
  properties **MUST** be documented *within* their parent concept's chapter, not
  promoted to their own chapters.

  > Rationale: nesting concepts under concepts produced a heading swamp and a
  > confusing flat list of mixed altitudes; one level keeps the table of contents
  > short and every reader oriented to which chapter they are in.

- Each concept chapter **MUST** carry both layers: a **reference** part defining
  the concept and its properties, and a **how-to** part (a "Working with…"
  section) holding the authoring guidance.

- How-to guidance **MUST** be written as **directives** — a short imperative
  (**Do** / **Consider** / **Avoid** / **Don't**) each paired with a brief *why*.
  A directive should be promoted to a named **job** (a verb-phrased task
  heading grouping several directives) only when the task needs judgment — a
  choice among options, a tradeoff, or several steps; a single rule should
  stay a floating directive rather than acquire a job wrapper.

- Concept chapters should be ordered as an author builds a model: the file
  and model frame first, then the rating scale, then the target tree (target →
  factor → requirement), then the Markdown body, then maintenance.

These conventions exist so the guide reads consistently and a reader can either
read a chapter straight through to understand a concept or jump to "Working with…"
when mid-task. They are the guide's editorial contract, not constraints on the
`QUALITY.md` format itself.
