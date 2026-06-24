---
type: Runtime Guide
title: Authoring Requirements
description: Requirement statements, assessments, factor connections, claim splitting, and rating override guidance.
tags: [quality, authoring, guide]
---

# Authoring Requirements

Read this when:

- creating, revising, reviewing, or evaluating Requirements, Assessments, factor references, or rating overrides.

Depends on:

- `../authoring.md`

---

## Requirement

A **requirement** is an assessable quality expectation — the single unit the
model is built to judge. Its **statement** (the map key) is its identity in
reports. Its `assessment` produces **findings** — observations about the source —
and those findings are rated together to yield the requirement's **rating
result**: one level on the rating scale, or *not assessed*.

Every requirement must be connected to at least one factor:

- **By placement** — declared under a factor or sub-factor. That factor is its
  **primary** factor and the requirement joins that factor's roll-up.
- **By reference** — naming factors under `factors`. On a nested requirement
  these are **secondary** factors (it appears in additional roll-ups). On a
  requirement placed directly under an area, `factors` is **required** and names
  the factors directly.

However it is connected, a requirement is assessed once, against the source of
the area it sits on, and counts once in that area's local rating.

### Properties

| Property     | Presence                     | What it is                                                       |
| ------------ | ---------------------------- | ---------------------------------------------------------------- |
| `assessment` | Required                     | The means of assessing the source; produces the findings.        |
| `factors`    | Required for area-level reqs | factor references; secondary factors when nested under a factor. |
| `ratings`    | Optional                     | Per-requirement criterion overrides, keyed by Rating Level ID.   |

### Working with requirements

#### Operationalize a factor into an assessable property

An area does not directly manifest a factor like "reliability"; it exhibits
concrete properties that imply it. Bridging from the abstract characteristic to
something inspectable is the core authoring move.

- **Do** name a concrete, observable property of the entity that the factor
  depends on, then phrase the requirement as the expectation about that property.
  *Reliability (abstract) → "recovers from a dependency outage without data loss"
  (property) → the requirement and its assessment. The move is domain-general:
  for a reference document, credibility → "every load-bearing claim cites a
  verifiable source."*
- **Do** preflight each requirement by naming its **scale** (the dimension and
  unit a finding lands on — *time for a novice to complete a one-item order*), its
  **meter** (the agreed procedure that produces the finding — *median over 100 test
  users using only online help*), and the rating boundary it tests. *If you cannot
  name a scale and a meter, the statement is still a slogan — "easy to use" has
  neither.*
- **Avoid** jumping straight from a factor to a convenient metric with no property
  in between. *That is how you get measurements no one can tie back to why they
  matter.*

#### Write ratable requirement statements

- **Do** phrase the map key as the thing you expect to be true ("Every record
  has a unique key", "p99 request latency stays within budget").
  *The statement is what shows up in reports; it should read as a claim that can
  be true or false to a degree.*

- **Do** write requirements specific enough that a single result stands on its
  own. *A vague requirement produces a vague rating.*
- **Do** write requirements that make the body context assessable. *A requirement
  should turn an important need, risk, or noted unknown into an expectation an
  evaluator can inspect.*
- **Do**, for behavioral qualities (reliability, recoverability, security under
  attack), phrase the statement around the *triggering condition and operating
  environment*, not just the steady state: "When a downstream dependency times out,
  requests fail over within 2 s and surface a degraded-mode response" — or, for a
  support process, "when an urgent ticket arrives after hours, it is acknowledged
  within 15 minutes." *A bare "is reliable" hides the condition that makes the
  quality observable.*
- **Do** make every requirement *verifiable*: its assessment must name a concrete
  method by which a finding could be produced. *If you cannot state how the source
  would be examined, the requirement is not assessable yet, however well it reads.*
- **Do** name the evidence the assessment draws on and check it can actually speak
  to the claim. *A latency claim needs runtime telemetry; a structure claim needs
  the source. If the only available evidence cannot address the claim, the
  requirement returns "not assessed" no matter how well written — narrow it to what
  the evidence supports, or note the gap as an unknown in the relevant section.*
- **Do** apply the discard test: if a requirement were deleted, would any decision
  about this root area change? *If not, it is ritual — drop it or note it as an
  unknown in the relevant section. Imported requirements meant to "be thorough" inflate
  the model and dilute the
  ratings that drive choices.*

#### Let risk decide where requirements go deep

- **Do** spend requirement detail where risk exposure (likelihood × impact) is
  highest — the failure modes named in body Risks. *A high-risk concern deserves a
  sharply bounded requirement with measured criteria; a low-risk one can stay
  coarse or be noted as an unknown.*
- **Avoid** spreading equal effort across all requirements. *Uniform depth spends
  scarce judgment on concerns that will not change a decision.*

#### Give each requirement exactly one assessment

- **Do** declare one `assessment`, stated **inline** (criteria, a measurement
  procedure, an inspection checklist, a diagnostic) or as a **reference** to an
  entity that defines one.
- **Avoid** stacking several independent assessments under one statement. *Split
  it into separate requirements instead — each result must be independently
  ratable.*
- **Do** write an inline assessment so a second person would gather the same
  evidence and reach the same rating — name the window, the sample, the tool, or
  the checklist. *"Review the write path" is not reproducible; "review production
  write-path telemetry over a representative window and the recovery-test results"
  is.*
- **Consider** matching assessment rigor to stakes: a requirement that gates a
  release deserves a repeatable, reproducible method; a low-stakes or exploratory
  one can name a lighter inspection. *Over-specifying rigor everywhere wastes
  effort; under-specifying it where it counts hides risk.*
- **Do** make sure every assessment answers "which factor's question does this
  finding help rate?" *An assessment that measures something only because it is
  easy to measure, tied to no factor's concern, is noise that dilutes the report.*
- **Avoid** an assessment so narrowly metric-shaped that satisfying the letter
  abandons the intent ("coverage ≥ 80%" met by trivial tests). *Where a single
  number invites gaming, pair it with an inspection of intent or phrase the
  statement as the outcome you actually want.*

#### Reference an external assessment; don't copy it

- **Do** point at the spec, doc, or checklist that defines the assessment, naming
  it once.
- **Avoid** extracting, summarizing, or duplicating that content into the
  requirement. *Duplicated criteria drift out of sync with their origin.*
- **Do** reference that entity by the same selector used as its own
  area's `source`. *That shared selector is the edge between the two areas; it is
  what makes the dependency traceable.*
- **Do** point at the specific applicable part of the referenced entity, and pin a
  version where it matters. *An unversioned, whole-document reference leaves the
  verification scope unbounded — name the section or rule that defines the
  assessment.*

#### Use one referenced assessment when one guide governs several factors

- **Do** write one requirement when a single guide, spec, or checklist defines a
  coherent quality judgment that bears on several factors. Put the requirement at
  the area level, list every affected factor in `factors`, and reference the
  governing entity once.
- **Do** split the requirement only when the referenced entity defines claims
  whose results could legitimately diverge.

```yaml
requirements:
  "the quality model follows its authoring guide family":
    factors:
      - fitness-for-purpose
      - credibility
      - assessability
      - traceability
      - maintainability
    assessment: >
      Assess QUALITY.md against ./skills/quality/guides/authoring.md and its
      routed sub-guides under ./skills/quality/guides/authoring/,
      especially whether the body credibly supports the model, factors come
      from visible needs and risks, requirements are assessable, sources are
      inspectable, and unknowns or open questions are explicit.
```

#### Connect to factors deliberately

- **Do** rely on **placement** for the primary factor (nest the requirement under
  it), and add `factors` only to pull the result into additional (**secondary**)
  factor roll-ups.
- **Do** declare `factors` explicitly for a requirement placed directly under an
  area — it's required there, and `null`/`[]`/empty entries don't satisfy it.
- **Do** name in `factors` only factors that are **in scope** — declared on this
  area or on an ancestor. *A named factor that lives on a sibling or a
  descendant doesn't resolve; you can reach up to an ancestor's factor, but not
  across or down.*

#### Split by assessable claim, not by factor

- **Do** size a requirement to one claim you want to rate as a single judgment —
  not to one factor or one section of an assessment. *The requirement is the unit
  assessed and rated, once.*
- **Do** connect a claim that reads through several lenses to multiple factors
  (placement for the primary, `factors` for the rest) instead of copying it into
  a per-factor requirement.
- **Avoid** a set of requirements that reference the same entity with the same
  assessment, sliced one per factor.
- **But** keep genuinely independent claims separate even when they share a
  reference — the test is whether their results could legitimately diverge (one
  strong, one weak). *Many requirements can draw on one rich entity; that is not
  duplication.*
- **Avoid** joining two assessable claims with "and" / "or" / "and/or" in one
  statement — that is usually two requirements. *A conjunction signals a
  non-singular requirement whose halves could rate differently.*
- **Do**, when a statement uses totality words ("every", "all", "never"), make the
  assessment say how exhaustiveness is checked — enumerate the population or sample
  it. *A totality claim is only verifiable if the assessment establishes the
  "every."*
- **Do** keep conformance and fitness as separate requirements when they can
  legitimately diverge. *A root area can match its specification yet fail the user's
  real need, or serve the need while departing from spec; one result cannot carry
  both judgments.*

#### Override criteria only when the shared scale can't express the gradient

*Read this subsection when you are actually authoring rating overrides; a first
model rarely needs it.*

- **Consider** a `ratings` override when a requirement has a natural measured
  threshold or a distinct qualitative spectrum (e.g. latency bands).
- **Do** key overrides by existing Rating Level IDs and change *only* the
  `criterion`. **Avoid** touching a level's `description`, order, or `title` —
  those stay fixed across the model.
- **Do** treat a measured override as a pair: the value you aim for (lands at
  `target`) and the range you will still accept (lands at `minimum`). *If a
  requirement has a natural number but you cannot name both, you have not yet
  decided what "good enough" means.*

```yaml
"p99 request latency stays within budget":
  factors: [reliability]
  assessment: >
    Measure p99 request latency over a representative production window.
  ratings:
    outstanding: "p99 at or under 150 ms."
    target: "p99 at or under 300 ms."
    minimum: "p99 at or under 500 ms."
    unacceptable: "p99 above 500 ms."
```

#### Validate the requirement set before treating the model as done

Individual requirements can each be sound while the set as a whole has holes,
conflicts, or duplicates. Run a closing pass over the set:

- **Do** check **completeness**: every Need and Risk the body raises is covered by
  a requirement, marked out of Scope, or noted as an unknown in the relevant
  section.
- **Do** check **consistency**: no two requirements make conflicting claims about
  the same source, and the same term or unit means the same thing across
  statements and criteria.
- **Avoid** **redundancy**: two statements that would always rate together against
  the same source are one requirement. *(Same test as for splitting — could their
  results legitimately diverge?)*
- **Do** note unresolved holes as unknowns rather than keeping a vague,
  unratable requirement to stand in for them.
