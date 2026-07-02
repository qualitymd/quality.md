---
type: Runtime Guide
title: Authoring factors
description: Factor naming, coverage, descriptions, stable-stakes factors, and sub-factor guidance for QUALITY.md models.
tags: [quality, authoring, guide]
---

# Authoring factors

Read this when:

- creating, revising, reviewing, or evaluating factors or sub-factors other than agent harnessability.

Depends on:

- `../authoring.md`

---

## Factor

A **factor** is a quality characteristic — a _lens_ through which an area's
quality is described. It groups the requirements assessed through it.
Conventional factor families differ by quality domain: a software product might
be viewed through `reliability`, `security`, or `maintainability`, a document or
data set through `completeness`, `credibility`, or `currentness`. These examples
are illustrative, non-exhaustive, and may overlap. A factor may decompose into
**sub-factors**: finer characteristics that together make up the parent (e.g.
`reliability` into `availability`, `fault tolerance`, `recoverability`). A
sub-factor is itself a factor of the same shape, nested to any depth.

Factor identity is local to its area: factors of the same name on two
different areas are distinct. A factor's stable ID is the declaring area ID plus
its path of factor names from that area's `factors` map.

Within one area, avoid reusing the same factor name anywhere in that area's
factor tree. The full factor IDs stay distinct, but requirement `factors`
entries use scalar names, so repeated names in one area are hard for authors,
agents, and reports to read unambiguously.

### Properties

| Property       | Presence    | What it is                                                        |
| -------------- | ----------- | ----------------------------------------------------------------- |
| `title`        | Required    | Human-readable label for reports and status output.               |
| `description`  | Recommended | The characteristic, defined operationally for this entity.        |
| `factors`      | Optional    | Sub-factors — recursively a factor.                               |
| `requirements` | Optional    | [Requirements](requirements.md) uniquely relevant to this factor. |

### Working with factors

#### Judge factor quality at two layers

Factor authoring has two related quality questions: whether each individual
factor is a good lens, and whether the set of factors is a good coverage model
for the area.

For a **factor set**, aim for comprehensive, proportionate, sustainable
coverage:

- **Comprehensive** — the set covers what matters. Consequential needs, risks,
  stakeholder concerns, stable stakes, and known unknowns are represented
  somewhere in the model or explicitly marked out of scope.
- **Proportionate** — the set gives most attention to what matters most. Depth,
  rigor, and granularity increase with consequence, risk, volatility, evidence
  cost, and decision value.
- **Sustainable** — the set remains usable over time. It is navigable, not
  needlessly redundant, affordable to maintain, and supportive of workflow
  rather than a parallel bureaucracy.

For an **individual factor**, first ask whether it is consequential, then improve
its shape:

- **Consequential** — the factor names a quality concern that matters for this
  entity. This is the admission rule: a concern that does not matter here does
  not earn a factor.
- **Bounded** — the factor has a clear boundary and grain. It is singular enough
  to reason about, distinct from neighboring factors, and placed at the right
  area or sub-factor level.
- **Operational** — the factor can be turned into assessable requirements,
  evidence, findings, and ratings. It is not merely an aspiration or theme.
- **Traceable** — the factor can be traced backward to needs, risks,
  stakeholders, or stable stakes, and forward to requirements that actually test
  it.
- **Neutral** — the factor names the quality to be judged, not a preferred
  implementation, workflow, tactic, or metric.

A consequential factor that is weak on the other qualities should usually be
refined, not dropped. Improve its name, boundary, placement, requirements,
evidence strategy, or assessment cadence until the model can judge it well
enough for the decisions it supports.

#### Choose factors that name the concerns that matter here

- **Do** pick the focused set of quality characteristics that genuinely drive this
  entity's quality. **Avoid** importing a standard checklist of characteristics
  wholesale.
- **Do** derive initial factors from the body's Needs and Risks. _Needs point at
  the outcomes quality should preserve; Risks point at the failure modes worth
  assessing._
- **Do** reconcile major Needs and Risks back to factors after drafting them.
  _If an important concern has no factor, either add the factor, mark the concern
  out of scope, or note the unresolved concern as an unknown in the relevant
  body section._
- **Do** justify each factor by something concrete about _this_ root area — who
  depends on it, what it is for, where it runs — not by its presence on a general
  list. _A characteristic with no user, no failure mode, and no decision riding on
  it here does not earn a factor, however standard it is elsewhere._
- **Do** prefer general-purpose, conventional factor names for the quality
  domain once a concern earns a factor. _For software product quality, examples
  include `reliability`, `security`, `usability`, `maintainability`,
  `performance`, `compatibility`, and `portability`; other domains have their
  own conventional factor families. Requirements and assessments map those
  lenses to the root area's unique quality expectations._
- **Avoid** inventing bespoke factor names for the subject's domain when a
  conventional quality attribute covers the concern. _The factor names the
  quality lens; the requirement says what that quality means for this entity._
- **Do** anchor each factor to a stakeholder whose concern it carries — the user
  who needs it to work, the maintainer who needs to change it, the operator who
  runs it. _Where stakeholders disagree on what "good enough" means, surface the
  conflict rather than averaging it into one criterion._
- **Consider** whether a factor names something a stakeholder _experiences_ (the
  system is available, decisions are correct) or something _internal_ that matters
  only because it produces that experience (low coupling, clear structure). _Keep
  internal factors tied to the outcome they serve, so an evaluator can tell a real
  weakness from a stylistic preference._

#### Cover the domain's stable stakes before specializing

Every quality domain has stable-stakes characteristics: concerns that predictably
affect trust, cost, change, use, operation, or stewardship for that kind of
entity. They matter even when the current root area is immature, performs poorly
on them, or lacks evidence about them.

- **Do** identify the quality domain before finalizing root factors. _A software
  product, document, data set, research report, model, service operation, and
  human process each has a different conventional factor family._
- **Do** include the domain's common stable-stakes factors for the root area, or
  explicitly justify why each omitted one is out of scope, delegated to a child
  area, or still unresolved as an unknown. _A sparse root model should be a
  conscious decision, not the result of only modeling the first risks that came
  to mind._
- **Do** treat roughly ten factors as a reasonable aim for a **primary-subject
  node** (see [Choose the decomposition shape](model-structure.md#choose-the-decomposition-shape-primary-subject-collection-or-composite)).
  _Fewer than eight should trigger a coverage review; four to six is usually too
  thin unless that node is deliberately narrow, temporary, or mostly delegated to
  child areas._ At a **composite** root the aim applies **per constituent**, not at
  the root: each primary-subject constituent earns its own ~ten-factor family,
  while the composite root itself carries only the factors that recur across
  constituents — typically those tracing to stewardship concerns (currentness,
  traceability, consistency, maintainability), each refined per child.
- **Avoid** dropping a conventional factor because the current artifact lacks
  evidence or performs poorly on it. _No tests is not a reason to omit
  `testability`; it is a reason to write requirements that make testability, test
  evidence, or the absence of evidence assessable._
- **Do** keep stable-stakes factors conventional and put subject-specific
  interpretation in requirements. _The factor names the durable quality lens; the
  requirements say what that quality means for this root area._
- **Avoid** padding the model with factors no stakeholder would notice and no
  decision would use. _The aim is coverage adequacy, not ceremony._

#### Name the quality, not the practice

- **Do** name a factor as a quality characteristic the area can exhibit to a
  degree: `completeness`, `consistency`, `credibility`, `currentness`,
  `understandability`, `traceability`, `assessability`, `maintainability`,
  `modifiability`, `testability`. _The name should read as the thing being
  judged, not the work someone does to improve it._
- **Avoid** factor names that describe a workflow, lifecycle phase, authoring
  technique, or evaluation tactic. _For example, `lifecycle-stewardship` is a
  practice; `maintainability`, `modifiability`, or `currentness` names the
  quality it is meant to protect. `grounding` is a tactic or metaphor;
  `credibility` or `traceability` names the observable quality._
- **Do** choose the narrower attribute when a broad label hides the real concern.
  _If the concern is missing expected context, use `completeness`; if claims
  contradict, use `consistency`; if claims lack believable support, use
  `credibility`; if context is stale, use `currentness`; if readers cannot
  interpret it, use `understandability`; if origin, rationale, or dependency
  paths are unclear, use `traceability`._
- **Do** draw from established attribute families as prompts, not quotas.

#### Write the description as an operational definition

- **Do** define the characteristic as the _degree or capability to achieve some
  end under the conditions that matter_, and say why it matters and to whom. A
  useful shape: _"\<factor\> is the degree to which \<entity\> \<achieves some
  end\> under \<conditions\>; it matters because \<stakeholder concern\>."_
- **Avoid** an adjective or a synonym for the factor name ("Reliability: how
  reliable it is"). _That tells a reader nothing and doesn't distinguish it from
  its siblings._

- **Do** write each factor's description so the factors on an area read as a
  distinct, non-overlapping set. _Overlapping factors make it ambiguous where a
  requirement belongs and double-count concerns in roll-up._

#### Decompose into sub-factors only when it aids understanding

- **Consider** sub-factors when a factor carries more than one distinct concern
  that's clearer assessed apart than together.
- **Avoid** decomposing a factor whose requirements already speak for themselves.
  _Decompose only as far as it helps._
- **Consider** decomposing a factor when its concern resists any direct
  assessment — a sub-factor you _can_ observe or proxy is more useful than a parent
  you can only assert. _Decompose for measurability, not only for readability._

- **Do** treat a child area's factor that shares a name with an ancestor's as a
  local refinement tailored to the child. _They're distinct factors with
  distinct IDs; a requirement in the child area references the child factor, not
  the ancestor._
- **Avoid** reusing the same factor name in multiple places inside one area's
  own factor tree. _Scalar `factors` references name local factors by name, so
  duplicate names inside one area make the intended connection ambiguous._

---
