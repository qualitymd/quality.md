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
present them as modeled by default unless a disqualifier fires (a harness-less or
throwaway project hits not-germane), rather than as a required roster. It
**MUST** teach that the QUALITY.md self-check follows the ordinary area pattern:
the `quality-md` key, a title of the form `<Root Title> QUALITY.md`, an explicit
path-based `source` such as `./QUALITY.md`, Factors that describe the model
artifact's qualities, and a Requirement that assesses the model against the
active authoring guide. It **MUST** teach that, when `quality-md` is in evaluation
scope, the skill assesses, analyzes, reports, and rolls it up like any other Area;
it **MUST NOT** teach that the QUALITY.md self-check is excluded from aggregate
rating, withheld from the root Area aggregation, reported only on a separate axis, or
given different evaluation semantics solely because its source is `QUALITY.md`.
It **MUST** teach that the agent harness is partly normative and plays the dual
area/assessment-reference role.

> Rationale: a root factored as one primary subject silently equates the entity
> with a single constituent and drops the other high-leverage artifacts; naming
> the composite shape and the recurring use-context constituents keeps them
> visible while preserving domain agnosticism. — 0074
>
> Rationale: `quality-md` already has ordinary Area structure and the CLI/report
> surfaces treat Areas generically. Keeping a separate aggregation exception in
> guidance made "full evaluation" ambiguous without a mechanical contract to
> enforce it. — 0082

The guide **MUST** teach that, for a composite root, the author enumerates
**domain constituents** by **constituent kind** inferred from the entity's
quality domain — not only the components the repository already has folders for —
using two generators: a **stewardship-concern** axis and an **audience ×
purpose** axis. The stewardship-concern axis **MUST** comprise a **lifecycle**
band (discover, define, realize, verify, enable, operate, maintain) and a
cross-cutting **protective** pair, **secure** (guard the entity from harm by the
world) and **safeguard** (guard stakeholders and the environment from harm by the
entity); the guide **MUST** name each concern by its function rather than a
domain-specific artifact, **MUST** present the protective pair as cross-cutting
stewardship under vulnerability — tracking who is exposed rather than which
lifecycle phase applies — and **MUST NOT** present `safeguard` as a synonym for
`secure`. The audience × purpose axis **MUST** cite Diátaxis once as that lens
applied to the *enable* concern and **MUST** be derivable from the body's Needs.

The guide **MUST** teach the **three-projections rule**: a stewardship concern
projects as a **factor**, a **constituent**, and an **audience**, so shared names
reflect a shared concern rather than duplication, and the author models the
projection meant rather than double-counting (the security *of* an area is a
factor; a security policy is a constituent). The guide **MUST** require the author
to encode the projection boundary in the emitted model when a model carries two or
more projections of one concern, rather than only reasoning about it during
modeling: on each modeled projection's node, a YAML comment that names the sibling
projection and the one-line distinction, and — when both projections are rated
nodes that surface in an evaluation report — a short disambiguating clause in each
node's `description` in addition to the comment. It **SHOULD** keep that clause to
the distinction from the sibling projection, consistent with the
"distinguishes, not enumerates" description rule. The guide **MUST** keep the
motivation-layer stewardship/care vocabulary (stewardship, care, tending,
vulnerability, concern) from modifying or replacing a taxonomy noun (factor,
area, requirement, constituent, audience): a concern is the source a factor
projects from, not a kind of factor. It **MUST** name the root's recurring factors
as **model-wide** (or cross-cutting) factors — which **MAY** be noted as tracing
to stewardship concerns — and **MUST NOT** render that link as "stewardship
factors" or "stewardship lenses"; the singular gloss "a factor is a quality lens"
is preserved. The guide **MUST** instruct the author to enumerate the implied
constituent kinds and **model each as its own area by default** — framing modeling
as the default outcome rather than something a constituent must earn — and
**MUST NOT** give the thinness of a first pass as a reason to defer or omit one. It
**MUST** state that a germane constituent is given its own area unless one of
exactly two disqualifiers holds — *no distinct concerns* (fold into a parent or
sibling) or *not germane / outside the boundary* (mark out of Scope) — and
**MUST NOT** list the absence of a constituent's artifact as a reason to omit a
germane constituent. The guide **MUST** teach that a germane concern is never
omitted by being recorded only in prose: when its artifact is absent or thin, its
absence **MUST** be surfaced as a ratable element of the model — a minimal area
carrying a missing-anchor finding, or a requirement on an existing area — and a
bare Scope or deferral note does **not** satisfy coverage. It **MUST** give the
routing criterion between the two (a minimal area when the kind would carry its own
factor family once it exists; a requirement on an existing area when the concern
folds in, the gap is partial, or a standalone area would be a single-finding stub).
The guide **MUST** present in-scope deferral as a narrow exception reserved for a
genuinely blocked constituent with the blocker recorded, **MUST NOT** accept "next
iteration" or "the first model is thin" as deferral reasons, and **MUST** scale
coverage to this entity — the kinds are a prompt traced to a Need or Risk, not a
roster every model must carry. It **MUST** distinguish the constituent question
(whether a tending leaves an owned, inspectable artifact, modeled as an area) from
the stewardship-quality question (whether it is done well, carried by factors),
reading a present artifact as evidence rather than proof.

> Rationale: 0074 named the composite shape but left domain constituents to
> "vary with what is modeled", so a setup-authored model enumerated constituents
> by folder and silently dropped the kinds without a folder — tests, specs, docs,
> a threat model. A domain-agnostic generator (stewardship concerns ×
> audience/purpose) makes those kinds inferable once the domain is named, while
> the three-projections rule keeps factors, constituents, and audiences from
> double-counting. — 0076
>
> Rationale: the stewardship concerns are care — an activity of tending whose
> artifact is its *trace* — so the claim that earns a constituent comes from a
> Need or Risk the entity presents, not from the generator list, and a present
> artifact is evidence (an *area*) distinct from whether the tending is done well
> (a *factor*). The protective pair is stewardship under vulnerability, which is
> why it cross-cuts the lifecycle rather than sitting in it. Framing only: the
> nine concerns, the two axes, and the earn-it test are unchanged. — 0077
>
> Rationale: the stewardship/care grounding leaked across the projection boundary
> it defines — the guide called recurring factors "stewardship lenses" and a live
> setup run reported "stewardship factors," demoting a term of art to a
> subcategory of the philosophical word. Keeping the motivation-layer vocabulary
> from modifying a taxonomy noun preserves consistent terms without retracting the
> grounding. — 0079
>
> Rationale: the earn-it inclusion test (0076–0077) supplied only
> anti-over-modeling pressure — "a prompt, not a quota," "earn each area" — with no
> counterweight, and "defer it in Scope" sat as a cost-free peer of "model it." A
> live setup run on a multi-service monorepo produced a flat root and deferred
> every per-constituent area to "the next iteration," passing maturity because a
> deferral note counted as "accounted for." Inverting to model-by-default with two
> disqualifiers, a no-silent-omission rule (a germane concern's absent artifact is
> surfaced as a ratable gap, never prose), and a completeness bar makes a first-pass
> model as full as the evidence supports. The stewardship/audience generators, the
> three-projections rule, and the vocabulary discipline are unchanged; only the
> inclusion default and its enforcement change. — 0080
>
> Rationale: the three-projections rule prevented double-counting for the *author*
> but left the boundary invisible to a *reader* of the emitted model — a field
> setup run produced a correct model carrying both an Agent Harnessability factor
> and an agent-harness area, yet a reader still asked whether the factor was meant
> to replace the area. Requiring the boundary to be encoded in the model (a YAML
> comment per projection, plus a description clause where both projections are rated
> nodes that surface in a report) promotes the reasoning into the artifact, the same
> way durable specs must carry their rationale forward. — 0087

The guide **MUST** teach **Agent Harnessability** as a model-wide umbrella factor
for an agent-collaborated composite root, using `agent-harnessability` as the
recommended stable key for new or revised models. It **MUST** define Agent
Harnessability as the degree to which the project's checked-in materials, tools,
workflows, feedback signals, standards, and action limits equip an AI agent to
understand the project, take scoped work, operate the environment, preserve and
resume state, verify its output, and stay safely bounded while preserving clear
human direction, review, and accountability. The guide **MUST** present Agent
Harnessability as the factor projection of the agent-collaboration concern, with
the agent harness remaining the constituent projection and the agent remaining the
audience projection. It **MUST** keep the factor/constituent boundary explicit:
Agent Harnessability rates how each constituent equips an agent, while the
agent-harness area rates the steering artifact's own quality, and Agent
Harnessability is not assessed on the agent-harness area as a recursion of the same
evidence. As the canonical instance of the projection-boundary rule above, the guide
**MUST** require this boundary to be encoded in the model — a YAML comment on both
the `agent-harnessability` factor and the agent-harness area, plus a disambiguating
clause in each description.

The guide **MUST** present Agent Harnessability as a deliberate umbrella that
carries no requirements of its own and decomposes into seven independently
assessable sub-factors, its rating coming from rolling those up:
**agent-accessibility**, **task-specifiability**, **agent-operability**,
**continuity**, **self-verifiability**, **enforcement-of-standards**, and
**containment-of-action**. Each sub-factor **MUST** be named as a quality, carry an
operational definition, include illustrative requirements that remain
quality-domain agnostic, and include boundary guidance that prevents
double-counting with sibling sub-factors, common factors, and the agent-harness
constituent.

The guide **MUST** define **continuity** as the degree to which an agent can
preserve state and resume useful work across long-running tasks, compaction,
interruption, handoff, and fresh sessions. Its illustrative requirements **SHOULD**
cover progress or handoff artifacts that capture current state, decisions made,
remaining work, verification status, blockers, and next steps; resumptions that do
not depend on unrecoverable chat history; and progress records that reduce false
completion or context-anxiety failure modes. Its boundary guidance **MUST**
distinguish it from agent-operability, which covers a fresh session reaching a
ready-to-work environment, and from agent-accessibility, which covers durable,
reachable knowledge.

The guide **MUST** strengthen the existing sub-factor examples so
agent-accessibility covers progressive disclosure and context selectivity;
task-specifiability covers decomposition, explicit done criteria, and completion
discipline against the original task; agent-operability covers tool affordance
quality, fresh-session readiness, act/observe loops, and agent-useful output;
self-verifiability covers deterministic checks, inferential evals, trace/run
evidence, and actionable machine-readable feedback; enforcement-of-standards covers
input, output, and tool guardrails or equivalent domain-neutral controls; and
containment-of-action covers least privilege, approval gates, sandboxing,
sensitive-resource boundaries, and auditability. Self-verifiability **MUST** name
good verification signals as fast, actionable, grounded in concrete evidence,
context-aware, and suppressible through visible reviewable exceptions; distinguish
deterministic signals from inferential evals for behavioral or non-deterministic
outcomes; and keep the boundary explicit that self-verifiability surfaces
exceptions while enforcement-of-standards constrains escapes.

The guide **MUST** route improvement of the harness over time to the model-wide
continuous-improvement or learn-loop concern, not to an eighth Agent Harnessability
sub-factor. The guide **MUST** state that Agent Harnessability is proposed by
default for an agent-collaborated composite root and **MUST NOT** give the thinness
or absence of the harness as a reason to omit it; thinness or absence is a low
rating and finding. The guide **SHOULD** treat an existing `harnessability` factor
with the legacy six-sub-factor shape as semantic coverage of the same concern, and
**SHOULD** recommend renaming it to `agent-harnessability` / Agent Harnessability
and adding `continuity` during model-authoring work unless the project has an
explicit reason to preserve the old key.

> Rationale: the agent-collaboration concern was present as a modeled-by-default
> constituent after 0080, but its cross-cutting quality projection was still
> missing. Because the concern recurs across constituents, it belongs as a
> model-wide factor. Decomposing it into the agent's working loop lets authors
> assess legibility, task framing, operability, self-checking, output enforcement,
> and action containment without double-counting the steering artifact itself. —
> 0081
>
> Rationale: long-running-agent practice makes state preservation and handoff too
> central to leave implicit under accessibility or operability. A project that
> cannot resume with the right state, decisions, remaining work, and verification
> status is not fully harnessable. — 0089

The guide **MUST** give the agent-harness area a concrete modeling template at
parity with the QUALITY.md self-check template: it **MUST** identify the area as an
enable and partly normative constituent for checked-in steering materials — agent
entry points, agent guidance files, skills, prompts, and related instructions — and
**MUST** offer a richer illustrative factor family rather than one or two placeholder
factors. Candidate factors **SHOULD** include or cover `completeness`, `accuracy`,
`currentness`, `understandability`, `coherence`, `selectivity`, `discoverability`
or `triggerability`, `maintainability`, `trustworthiness`, and `assessability`.
The family **MUST** be presented as an illustrative prompt earned per entity, not
as a required or default roster.

The guide **MUST** direct that agent-harness area factors and requirements remain
agnostic to the served domain — the domain the project model is about — and
**MUST NOT** present software-specific mechanisms such as lint, type-check, test,
or CI gates as harness requirements except as one domain's instance of
domain-neutral expectations about verification, standards, or action boundaries.

> Rationale: the guide had a concrete template for the QUALITY.md self-check but
> left the sibling harness constituent thinly specified. Explicit harness guidance
> is licensed by the agentic use context, but the harness serves whatever domain the
> project models, so its requirements must not leak software assumptions. — 0089
