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
The guide **MUST** teach the strict name grammar for Area names, Factor names,
and Rating Level IDs, and that Requirement statements remain natural-language
keys outside that grammar. The guide **MUST** distinguish human display titles
from structured IDs and canonical model references.
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

The guide **MUST** include at least one worked Markdown body section example
showing the common section shape, concise judgment context, agent-accessible
support, section-scoped unknowns, section-scoped open questions, and the review
state line.

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
(naming the agent surface and model used, for example `Codex (GPT-5.5)`), and
that the human review advances only when a person reads and endorses the
section.

> Rationale: the body is largely agent-authored, so the only freshness signal
> worth trusting is when a person last stood behind the section. — 0044

The guide **MUST** teach that the rating scale should be reviewed after the body
and before writing requirements, so the shared rating vocabulary fits the
root area's decision context.

The guide **MUST** teach that the recommended four-level Rating Scale keeps
stable Rating Level IDs as `outstanding`, `target`, `minimum`, and
`unacceptable`, and that its default human display titles are
`🟢 Outstanding`, `🔵 Target`, `🟡 Minimum`, and `🔴 Unacceptable`. It **MUST**
frame the emoji markers as a human scanning aid, not as rating identity,
ordering, semantics, or a conformance requirement, and it **SHOULD** discourage
emoji-only titles.

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

The guide **MUST** teach the three decomposition shapes an area — the root
included — can take: **primary-subject** (one entity, one factor family),
**collection** (many entities of the same kind, judged by whole-set concerns),
and **composite** (many entities of different kinds, each with its own
largely-disjoint factor family, the node's own concern being cross-part
coherence). It **MUST** teach that the shapes are recursive and composable — a
node of one shape **MAY** hold children of another, to any depth — that a
near-disjoint factor family is a first-class signal to split a part into its own
area, and that the root-factor coverage aim applies per primary-subject node
while a composite root carries only the factors that recur across its
constituents.

The guide **MUST** teach that QUALITY.md's agentic context of use — not the
modeled domain — makes two constituents recur in a composite root regardless of
domain: the **agent harness** and the **QUALITY.md self-check**. It **MUST**
present them as expected-but-earned constituents rather than a required roster,
**MUST** teach that the QUALITY.md self-check stays on the learn loop and out of
the entity's roll-up, and **MUST** teach that the agent harness is partly
normative and plays the dual area/assessment-reference role.

> Rationale: a root factored as one primary subject silently equates the entity
> with a single constituent and drops the other high-leverage artifacts; naming
> the composite shape and the recurring use-context constituents keeps them
> visible while preserving domain agnosticism. — 0074

The guide **MUST** teach that, for a composite root, the author enumerates
**domain constituents** by **constituent kind** inferred from the entity's
quality domain — not only the components the repository already has folders for —
using two generators: a **stewardship-concern** axis and an **audience ×
purpose** axis. The stewardship-concern axis **MUST** comprise a **lifecycle**
band (discover, define, realize, verify, enable, operate, maintain) and a
cross-cutting **protective** pair, **secure** (guard the entity from harm by the
world) and **safeguard** (guard stakeholders and the environment from harm by the
entity); the guide **MUST** name each concern by its function rather than a
domain-specific artifact and **MUST NOT** present `safeguard` as a synonym for
`secure`. The audience × purpose axis **MUST** cite Diátaxis once as that lens
applied to the *enable* concern and **MUST** be derivable from the body's Needs.

The guide **MUST** teach the **three-projections rule**: a stewardship concern
projects as a **factor**, a **constituent**, and an **audience**, so shared names
reflect a shared concern rather than duplication, and the author models the
projection meant rather than double-counting (the security *of* an area is a
factor; a security policy is a constituent). The guide **MUST** instruct the
author to account for each implied constituent kind — model it, defer it in
Scope, mark it out of Scope, or record it as an unknown — treating silence as a
coverage gap, and **MUST** teach carrying a germane, high-leverage kind as an area
even when its artifact is thin or missing, recording the gap as a finding. It
**MUST** keep the inclusion test earned (owned, inspectable artifact; divergent
factor family; traced to a Need or Risk) and **MUST NOT** present the kinds as a
roster every model must carry.

> Rationale: 0074 named the composite shape but left domain constituents to
> "vary with what is modeled", so a setup-authored model enumerated constituents
> by folder and silently dropped the kinds without a folder — tests, specs, docs,
> a threat model. A domain-agnostic generator (stewardship concerns ×
> audience/purpose) makes those kinds inferable once the domain is named, while
> the three-projections rule keeps factors, constituents, and audiences from
> double-counting. — 0076
