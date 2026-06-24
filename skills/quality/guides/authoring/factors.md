---
type: Runtime Guide
title: Authoring Factors
description: Factor naming, coverage, descriptions, stable-stakes factors, and sub-factor guidance for QUALITY.md models.
tags: [quality, authoring, guide]
---

# Authoring Factors

Read this when:

- creating, revising, reviewing, or evaluating Factors or sub-factors other than Agent Harnessability.

Depends on:

- `../authoring.md`

---

## Factor

A **factor** is a quality characteristic — a *lens* through which an area's
quality is described. It groups the requirements assessed through it.
Conventional factor families differ by quality domain: a software product might
be viewed through `reliability`, `security`, or `maintainability`, a document or
data set through `completeness`, `credibility`, or `currentness`. These examples
are illustrative, non-exhaustive, and may overlap. A factor may decompose into
**sub-factors**: finer characteristics that together make up the parent (e.g.
`reliability` into `availability`, `fault tolerance`, `recoverability`). A
sub-factor is itself a factor of the same shape, nested to any depth.

Factor identity is local to its area: factors of the same name on two
different areas are distinct. A Factor's stable ID is the declaring Area ID plus
its path of Factor names from that Area's `factors` map.

### Properties

| Property       | Presence    | What it is                                                        |
| -------------- | ----------- | ----------------------------------------------------------------- |
| `title`        | Required    | Human-readable label for reports and status output.               |
| `description`  | Recommended | The characteristic, defined operationally for this entity.        |
| `factors`      | Optional    | Sub-factors — recursively a factor.                               |
| `requirements` | Optional    | [Requirements](requirements.md) uniquely relevant to this factor. |

### Working with factors

#### Choose factors that name the concerns that matter here

- **Do** pick the focused set of quality characteristics that genuinely drive this
  entity's quality. **Avoid** importing a standard checklist of characteristics
  wholesale.
- **Do** derive initial factors from the body's Needs and Risks. *Needs point at
  the outcomes quality should preserve; Risks point at the failure modes worth
  assessing.*
- **Do** reconcile major Needs and Risks back to factors after drafting them.
  *If an important concern has no factor, either add the factor, mark the concern
  out of scope, or note the unresolved concern as an unknown in the relevant
  body section.*
- **Do** justify each factor by something concrete about *this* root area — who
  depends on it, what it is for, where it runs — not by its presence on a general
  list. *A characteristic with no user, no failure mode, and no decision riding on
  it here does not earn a factor, however standard it is elsewhere.*
- **Do** prefer general-purpose, conventional factor names for the quality
  domain once a concern earns a factor. *For software product quality, examples
  include `reliability`, `security`, `usability`, `maintainability`,
  `performance`, `compatibility`, and `portability`; other domains have their
  own conventional factor families. Requirements and assessments map those
  lenses to the root area's unique quality expectations.*
- **Avoid** inventing bespoke factor names for the subject's domain when a
  conventional quality attribute covers the concern. *The factor names the
  quality lens; the requirement says what that quality means for this entity.*
- **Do** anchor each factor to a stakeholder whose concern it carries — the user
  who needs it to work, the maintainer who needs to change it, the operator who
  runs it. *Where stakeholders disagree on what "good enough" means, surface the
  conflict rather than averaging it into one criterion.*
- **Consider** whether a factor names something a stakeholder *experiences* (the
  system is available, decisions are correct) or something *internal* that matters
  only because it produces that experience (low coupling, clear structure). *Keep
  internal factors tied to the outcome they serve, so an evaluator can tell a real
  weakness from a stylistic preference.*

#### Cover the domain's stable stakes before specializing

Every quality domain has stable-stakes characteristics: concerns that predictably
affect trust, cost, change, use, operation, or stewardship for that kind of
entity. They matter even when the current root area is immature, performs poorly
on them, or lacks evidence about them.

- **Do** identify the quality domain before finalizing root factors. *A software
  product, document, data set, research report, model, service operation, and
  human process each has a different conventional factor family.*
- **Do** include the domain's common stable-stakes factors for the root area, or
  explicitly justify why each omitted one is out of scope, delegated to a child
  area, or still unresolved as an unknown. *A sparse root model should be a
  conscious decision, not the result of only modeling the first risks that came
  to mind.*
- **Do** treat roughly ten factors as a reasonable aim for a **primary-subject
  node** (see [Choose the decomposition shape](model-structure.md#choose-the-decomposition-shape-primary-subject-collection-or-composite)).
  *Fewer than eight should trigger a coverage review; four to six is usually too
  thin unless that node is deliberately narrow, temporary, or mostly delegated to
  child areas.* At a **composite** root the aim applies **per constituent**, not at
  the root: each primary-subject constituent earns its own ~ten-factor family,
  while the composite root itself carries only the factors that recur across
  constituents — typically those tracing to stewardship concerns (currentness,
  traceability, consistency, maintainability), each refined per child.
- **Avoid** dropping a conventional factor because the current artifact lacks
  evidence or performs poorly on it. *No tests is not a reason to omit
  `testability`; it is a reason to write requirements that make testability, test
  evidence, or the absence of evidence assessable.*
- **Do** keep stable-stakes factors conventional and put subject-specific
  interpretation in requirements. *The factor names the durable quality lens; the
  requirements say what that quality means for this root area.*
- **Avoid** padding the model with factors no stakeholder would notice and no
  decision would use. *The aim is coverage adequacy, not ceremony.*

#### Name the quality, not the practice

- **Do** name a factor as a quality characteristic the area can exhibit to a
  degree: `completeness`, `consistency`, `credibility`, `currentness`,
  `understandability`, `traceability`, `assessability`, `maintainability`,
  `modifiability`, `testability`. *The name should read as the thing being
  judged, not the work someone does to improve it.*
- **Avoid** factor names that describe a workflow, lifecycle phase, authoring
  technique, or evaluation tactic. *For example, `lifecycle-stewardship` is a
  practice; `maintainability`, `modifiability`, or `currentness` names the
  quality it is meant to protect. `grounding` is a tactic or metaphor;
  `credibility` or `traceability` names the observable quality.*
- **Do** choose the narrower attribute when a broad label hides the real concern.
  *If the concern is missing expected context, use `completeness`; if claims
  contradict, use `consistency`; if claims lack believable support, use
  `credibility`; if context is stale, use `currentness`; if readers cannot
  interpret it, use `understandability`; if origin, rationale, or dependency
  paths are unclear, use `traceability`.*
- **Do** draw from established attribute families as prompts, not quotas.

#### Write the description as an operational definition

- **Do** define the characteristic as the *degree or capability to achieve some
  end under the conditions that matter*, and say why it matters and to whom. A
  useful shape: *"\<factor\> is the degree to which \<entity\> \<achieves some
  end\> under \<conditions\>; it matters because \<stakeholder concern\>."*
- **Avoid** an adjective or a synonym for the factor name ("Reliability: how
  reliable it is"). *That tells a reader nothing and doesn't distinguish it from
  its siblings.*

- **Do** write each factor's description so the factors on an area read as a
  distinct, non-overlapping set. *Overlapping factors make it ambiguous where a
  requirement belongs and double-count concerns in roll-up.*

#### Decompose into sub-factors only when it aids understanding

- **Consider** sub-factors when a factor carries more than one distinct concern
  that's clearer assessed apart than together.
- **Avoid** decomposing a factor whose requirements already speak for themselves.
  *Decompose only as far as it helps.*
- **Consider** decomposing a factor when its concern resists any direct
  assessment — a sub-factor you *can* observe or proxy is more useful than a parent
  you can only assert. *Decompose for measurability, not only for readability.*

- **Do** treat a child area's factor that shares a name with an ancestor's as a
  *refinement* tailored to the child. *They're technically distinct factors;
  write the child's description to say how it specializes the ancestor's
  concern.*

---
