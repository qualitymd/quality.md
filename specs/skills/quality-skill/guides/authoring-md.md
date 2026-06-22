---
type: Functional Specification
title: QUALITY.md authoring guide
description: Contract for the skill's authoring guide — the canonical reference and best-practices guide for understanding and working with QUALITY.md files.
tags: [skill, quality, guide]
timestamp: 2026-06-18T00:00:00Z
---

# QUALITY.md authoring guide

This spec governs the **authoring guide** the [`/quality` skill](../quality-skill.md) ships at
[`skills/quality/guides/authoring.md`](../../../../skills/quality/guides/authoring.md)
— the document the skill reads when creating, populating, reviewing, or
improving a QUALITY.md file. It has 1:1 coverage with that document: this spec
is its contract, and the guide is its implementation. The format the guide
teaches is defined by [`SPECIFICATION.md`](../../../../SPECIFICATION.md), the source
of truth the guide **conforms to**.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", and
"MAY" are to be interpreted as described in
[RFC 2119](../../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../../docs/reference/rfc8174.md) when, and only when, they appear
in all capitals.

## Purpose

The guide exists to be the **canonical reference and best-practices guide for
understanding and working with QUALITY.md files**: one document a reader can
stay inside to learn what each concept of the format *is* and how to author a
useful model. It serves both the skill (its reader at runtime) and a human
author. The guide **MUST** state this purpose in its own preamble, so the
document declares the job it is built to do.

## Scope

The defining scope decision is **self-containedness**. The guide is *all-inclusive*:
it restates the format's concepts, properties, and rating vocabulary in full
rather than deferring to [`SPECIFICATION.md`](../../../../SPECIFICATION.md) for the
definitions. A reader **MUST** be able to understand and author a QUALITY.md
from the guide alone, without opening the format spec.

This is a deliberate exception to the skill's general "don't embed a drifting
copy of the format" rule, which binds the skill's *metadata and prompt* — those
are grounded at runtime from `qualitymd spec` (see
[Invocation](../quality-skill.md#frontmatter-and-metadata)). The guide is a bundled
**reference resource**, not the prompt, and a self-contained reference is the
whole point of this one: the bundled [`SPECIFICATION.md`](../../../../SPECIFICATION.md)
copy is the same kind of artifact. The duplication is paid for by the conformance
duty below, which makes the format spec the authority whenever the two disagree.

> Rationale: a guide a reader must cross-reference against the spec to use is not
> a single comprehensive guide. Self-containedness is the purpose; the
> conformance duty is the price that keeps it honest.

**In scope:** authoring a useful QUALITY.md — its file shape, the model
concepts (rating scale, areas, factors, requirements), the Markdown body, and
durable authoring practices such as body-first model development, rating-scale
fit, deriving factors from Needs/Risks, writing assessable requirements, and
recording each section's unknowns and open questions.

**Non-goals:** the guide does **not** document the *evaluation* process (Define →
Assess and Rate → Analyze → Advise → Report) beyond what bears on authoring; that
is the skill's own workflow (see
[Evaluation workflow](../evaluation.md#evaluation-workflow)). It does **not** restate the CLI
surface, and it is **not** a normative spec — `SPECIFICATION.md` and the durable
specs remain the contracts; the guide is instructional. It also does **not**
carry first-run workflow sequencing after `qualitymd init`; that procedural
playbook belongs to the getting-started guide, which depends on this guide.

## Requirements

### Conformance

The guide **MUST conform to** [`SPECIFICATION.md`](../../../../SPECIFICATION.md):
every concept definition, property, presence level (Required / Recommended /
Optional), and rating-vocabulary term it states **MUST** match the format spec.
The guide **MUST** teach that `title` is required on the Model, every Area,
every Factor, and every Rating Level, and that requirements do not have a
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
  QUALITY.md concept (the file, the model, rating scale, area, factor,
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

- Concept chapters should be ordered as an author builds a useful first model:
  the file and model frame first, then the Markdown body, then the rating scale,
  then the area tree (area → factor → requirement), then maintenance.

These conventions exist so the guide reads consistently and a reader can either
read a chapter straight through to understand a concept or jump to "Working with…"
when mid-task. They are the guide's editorial contract, not constraints on the
QUALITY.md format itself.

### Best-practice coverage

The guide **MUST** teach body-first model development: fill the Markdown body
before expanding the model tree, because the body supplies the evaluable
judgment context for choosing the rating scale, factors, requirements, and
scope, and for later judging whether the model still fits the root area.

The guide **MUST** state desired outcomes for the recommended body sections:
Overview, Scope, Needs, and Risks. The guide **MUST** teach that the body is
context for building the model, understanding the model's purpose, using the
model in evaluation, evaluating the model's quality, and deciding whether the
model still fits the root area.

The guide **MUST** teach a common section shape under which each section records
its own unknowns (broad areas of uncertainty about the section's topic that may
not resolve to a single answer) and open questions (specific questions with one
particular answer, still unresolved), scoped to the section and distinct from a
`not assessed` result, and **MUST** teach stating "none known" rather than
omitting them.

> Rationale: a single catch-all Known gaps section sat far from the content it
> qualified and was skimmed past; uncertainty belongs with the section it
> concerns. — 0044

The guide **MUST** teach that body sections are themselves evaluable: a later
human or agent should be able to judge their completeness, thoroughness,
recency, root area specificity, grounding, agent-accessibility, and open
questions. The guide **MUST** teach sections to be concise, rigorous, and
self-explanatory, with supporting detail cited rather than copied when citation
keeps the body readable.

The guide **MUST** use and define **agent-accessible** support: support
available to the evaluating agent through the repository, cited local paths,
configured tools, linked public sources, or explicitly provided context. The
guide **MUST** instruct authors to cite material support when it is
agent-accessible and to record material support that is not agent-accessible as
a first-class limitation in the relevant section's unknowns or open questions.
The guide **MUST NOT** require a separate `Access gaps` line in every section.

The guide **MUST** teach a per-section review state line that records the last
human review (citing a named person) distinctly from the last agent review
(naming the agent), and that the human review advances only when a person reads
and endorses the section.

> Rationale: the body is largely agent-authored, so the only freshness signal
> worth trusting is when a person last stood behind the section. — 0044

The guide **MUST** teach that the rating scale should be reviewed after the body
and before writing requirements, so the shared rating vocabulary fits the
root area's decision context.

The guide **MUST** teach that initial factors should be derived from the body's
Needs and Risks, and that requirements should make the body context assessable.

The guide **MUST** teach authors to name factors as quality characteristics the
area can exhibit to a degree, not as practices, workflows, lifecycle phases,
authoring techniques, or evaluation tactics. It **SHOULD** show how to choose
narrower, better-established attributes when a broad label hides the real
concern, and it **SHOULD** distinguish product/tooling qualities from
data/document/model qualities.

The guide **MUST** teach that when one guide, spec, or checklist defines a
coherent assessment that bears on several factors, authors should write one
requirement, connect it to the affected factors through `factors`, and reference
the governing entity once. It **MUST** teach authors to split such requirements
only when the referenced entity defines claims whose results could legitimately
diverge.
