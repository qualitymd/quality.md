# QUALITY.md Specification

**Version 0.1 — Draft**

This is the specification for the `QUALITY.md` standard: a file format and set of conventions to help model, evaluate, and improve quality. This specification uses terminology from the software development, but the `QUALITY.md` standard can work in any context.

## Conformance

Conforming use and applications of the standard must fulfill all normative requirements. Conformance requirements are described in this document via both descriptive assertions and key words with clearly defined meanings.

The key words “MUST”, “MUST NOT”, “REQUIRED”, “SHALL”, “SHALL NOT”, “SHOULD”, “SHOULD NOT”, “RECOMMENDED”, “MAY”, and “OPTIONAL” in the normative portions of this document are to be interpreted as described in IETF RFC 2119. These key words may appear in lowercase and still retain their meaning unless explicitly declared as non-normative.

A conforming use or application of QUALITY.md may provide additional functionality, but must not do so where explicitly disallowed or where doing so would result in non-conformance.

## Key Terms

- **Quality Model**: a structured, declarative description of what quality means for a given entity.
- **Entity**: a thing that is evaluated for quality.
- **Subject**: the entity a quality model as a whole describes — the system, component, or artifact a `QUALITY.md` file is about, named by the model's `title` when set.
- **Factor**: a quality (sub)characteristic or attribute — a lens such as reliability or security — through which an entity's quality is described; it groups the requirements assessed through it and may be decomposed into sub-factors.
- **Requirement**: a quality requirement for assessing and rating the quality of an entity.
- **Finding**: a single observation produced by assessing the source entities against a requirement — a unit of evidence such as a measured value, an inspection note, or a diagnostic result. A finding records *what was observed* and is not itself rated; the **findings** of a requirement are rated together.
- **Assessment**: the means for assessing an entity — measurement, specifications, inspection, checklists, diagnostics, etc.
- **Rating Scale**: a defined set of rating levels for a quality model.
- **Rating Level**: a single level on a rating scale, providing the default criterion for rating a requirement's findings.
- **Rating Result**: the outcome of rating a requirement's findings against the rating scale — a single rating level (or a *not assessed* outcome) assigned to the requirement, considering all of its findings together.
- **Target**: an entity or set of entities with quality requirements subject to evaluation.
- **Source**: the scope of entities defined by a target.

## QUALITY.md File

A `QUALITY.md` file is a markdown file with YAML frontmatter with a structured quality model and a markdown body.

The presence of a `QUALITY.md` file in a directory implies that the directory and all its sub-directories and their contents are the implicit source of any quality evaluation for the contained quality model. Exceptions to this are when a different entity (or set of entities) is selected by setting a custom `source` property on the model or by special tooling.

### YAML Frontmatter

`QUALITY.md` files MUST begin with a valid YAML frontmatter block containing the required **Model** properties specified below.

When authoring `QUALITY.md` frontmatter, null or empty optional properties SHOULD be omitted.

#### Model

The model represents a quality model: what things (**Targets**) are evaluated for quality, their important quality characteristics (**Factors**), measurable quality **Requirements**, assessment criteria, and **Rating** criteria for determining the level of quality.

```yaml
title: <string>                 # Optional; title of the entity whose quality is modeled
ratingScale:                    # Required; the rating scale
  - level: <level-name>         #   Required; unique within the scale
    title: <string>             #   Optional; human-readable label
    criterion: <string>         #   Required; used to rate a requirement's findings
factors:                        # Optional*
  <factor-name>: <Factor>
requirements:                   # Optional*
  <requirement-statement>: <Requirement>
targets:                        # Optional*
  <target-name>: <Target>
source: <string>                # Optional
```

*An entry on either factors, requirements, or targets MUST be supplied.

**Title**: An optional human-readable name for the entity whose quality is modeled. For software projects, this is typically the name of the product, system, or library. A `title` is RECOMMENDED for readable reports but MAY be omitted.

**Rating Scale**: This is the rating scale that provides the default criterion for how requirement assessments should be judged to arrive at a rating level result. Each **Rating Level** MUST declare a `level` name, unique within the scale, and a `criterion` used to rate a requirement's findings; a **title** for improved readability is OPTIONAL. Rating levels MUST be ordered from best (first) to worst (last). At least two rating levels MUST be supplied.

**A suggested scale (non-normative).** When a graded scale fits and an author has no strong preference, the following four-level scale is a reasonable starting point, and a scaffolding tool MAY seed it. Its vocabulary names four best-to-worst bands — **outstanding** exceeds the goal, **target** meets it, **minimum** holds the acceptable floor, and **unacceptable** falls below it:

```yaml
ratingScale:
  - { level: outstanding,  title: Outstanding,  criterion: "Exceeds the requirement; satisfies it with margin to spare." }
  - { level: target,       title: Target,       criterion: "Satisfies the requirement." }
  - { level: minimum,      title: Minimum,      criterion: "Falls short of the goal but holds the acceptable floor." }
  - { level: unacceptable, title: Unacceptable, criterion: "Falls below the acceptable floor." }
```

**Factors**: quality characteristics or attributes that matter most for evaluating the overall quality of the entity.

**Requirements**: quality requirements that will be used to assess the quality of the entity. These are typically nested under a factor or target, but may be defined at the model root level when it is simpler to define a single requirement at the root and cross reference to multiple quality attributes.

**Targets**: more focused quality modeling for possible target entities. Not required but useful when a distinct set of factors or requirements would be more cohesively defined around a narrower target of evaluation than the scope implied by the source of the entire quality model.

**Source**: the location of the target entity subject to this quality model. This SHOULD be omitted on the top-level model to assume the default convention.

#### Target

A target is an entity, or set of entities, with quality requirements subject to evaluation. It is the recursive node of the quality model: the model root is itself the apex target, and every entry under `targets` is another target of the same shape, nested to any depth.

```yaml
factors:                        # Optional*
  <factor-name>: <Factor>
requirements:                   # Optional*
  <requirement-statement>: <Requirement>
targets:                        # Optional*
  <target-name>: <Target>
source: <string>                # Optional
```

*A target MAY declare no factors or requirements of its own and serve purely as a grouping node (holding only child targets), but each target SHOULD lead to at least one requirement somewhere in its subtree — its own, one carried by a factor, or one contributed by a descendant. A target whose subtree contains no requirements evaluates nothing.

A target shares the structure of the model root but for two keys — `title` and `ratingScale` — which are declared only on the model root (see [Model](#model)).

**Factors**: quality characteristics scoped to this target's subtree. A factor declared on a target applies to that target and its descendants, not to unrelated targets.

**Requirements**: quality requirements assessed against this target's source. A requirement is assessed once, at the target that declares it, against that target's source.

**Targets**: more focused child targets, nested to any depth. Child targets do not inherit their parent's requirements. However, an ancestor target's source selector may select entities that overlap with its descendant child target selectors, resulting in the ancestor's requirements also being evaluated on the same entities targeted by the child's source selector.

**Source**: the location of the entities this target evaluates. Paths and globs resolve relative to the containing `QUALITY.md` file. When a target omits `source`, it inherits the source scope of its nearest ancestor that declares one; a grouping target MAY leave it implicit and let its child targets narrow it.

#### Factor

A factor is a quality characteristic — such as `reliability`, `security`, or `maintainability` — through which a target's quality is described. A factor MAY be decomposed into **sub-factors**: finer characteristics that together make up the parent — for example `reliability` into `availability`, `fault tolerance`, and `recoverability`. A sub-factor is itself a Factor of the same shape, nested to any depth.

```yaml
description: <string>           # Recommended
factors:                        # Optional; sub-factors, recursively Factor
  <factor-name>: <Factor>
requirements:                   # Optional
  <requirement-statement>: <Requirement>
```

**Description**: a concise statement of the quality characteristic as it applies to this entity. A factor SHOULD declare a `description`, and that description SHOULD:

- **define the characteristic operationally** — what it means here, phrased as the degree or capability to achieve some end under the conditions that matter, not merely an adjective or a synonym for the factor name;
- **convey why it matters and to whom** — the stakeholder concern the factor answers; and
- **distinguish it from its sibling factors**, so the factors on a target read as a non-overlapping set.

A description SHOULD NOT restate, enumerate, or stand in for the factor's requirements: the measurable expectations belong in `requirements` and are judged through the rating scale. The description fixes *what the lens is*; the requirements fix *what is assessed through it*.

A useful shape is one or two sentences — "*\<Factor\> is the degree to which \<entity\> \<achieves some end\> under \<relevant conditions\>; it matters here because \<stakeholder concern\>.*" For example: *Reliability is the degree to which the orders API continues to accept and durably record orders under load and partial failure; it matters because an acknowledged order that is later lost is unrecoverable for the customer.*

**Requirements**: the quality requirements assessed through this factor's lens, each evaluated against the source of the target on which it is declared. A factor SHOULD lead to at least one requirement — one nested directly beneath it, one added by a refinement on a descendant target, or one declared elsewhere that tags this factor as a secondary factor (see [Requirement](#requirement)). A factor that nothing contributes to is a lens over nothing.

**Sub-factors**: a map of finer characteristics that decompose this factor, each a Factor of the same shape. Decompose a factor when it carries more than one distinct concern that is clearer assessed apart than together — but only as far as it aids understanding; a factor whose requirements already speak for themselves needs no sub-factors. The guidance above applies at every level: a sub-factor SHOULD carry its own `description`, SHOULD be distinguishable from its siblings, and SHOULD lead to at least one requirement. A sub-factor's requirements are assessed through its lens, and an evaluation infers a rating for the sub-factor that rolls up into the parent factor's rating (see [Analyze](#analyze)). A factor MAY hold both its own direct requirements and sub-factors.

Factor identity is local to its target. Factors of the same name declared on two different targets are distinct factors. Child targets that declare factors the same as an ancestor target's should be refinements of the ancestor target's factor tailored to the child target.

#### Requirement

A requirement is an assessable quality expectation — the single unit the model is built to judge. It pairs a **statement** (its map key, and its identity in reports) with an `assessment` that produces the findings; those findings are rated together to yield the requirement's **Rating Result**.

A requirement's placement sets the factors it informs. Nested under a factor or sub-factor, that factor is its **primary** lens and the requirement joins the factor's roll-up. Placed directly under a target with no `factors`, it is **unlensed**, informing only that target's local rating. In any placement a requirement MAY name **secondary** factors under `factors` to inform additional lenses, so one result can feed several factor views at once — and a directly-placed requirement MAY use this to attach itself to one or more factors with no primary among them. However it is placed, a requirement always contributes to its target's local rating (counted once) and is assessed once, against the source of the target on which it is declared.

```yaml
assessment: <string>            # Required; the means of assessing the source, producing findings
factors:                        # Optional; additional factors in scope this result also informs
  - <factor-name>
ratings:                        # Optional; per-requirement criterion overrides
  <level-name>: <criterion>     #   keyed by a level of the model's rating scale
```

**Assessment**: the means of assessing the target's source for this requirement — inline criteria, a measurement procedure, an inspection checklist, a diagnostic, or a path to a document describing one. A requirement MUST declare exactly one `assessment` as a single non-empty scalar; a missing, empty, or list-valued `assessment` is invalid. The assessment produces the requirement's **findings** — one or more observations, each recording *what was observed* and not itself rated (see [Assess and Rate](#assess-and-rate)). When a single statement needs several independent assessments, split it into several requirements rather than listing assessments under one.

**Factors (secondary)**: an optional list of factor names this requirement's result should also inform, beyond the factor it is nested under (its **primary** factor). Each name MUST resolve to a factor in scope — one declared on the target where the requirement sits, or on an ancestor target. A secondary factor lets one result appear in additional factor roll-ups without duplicating the requirement; the result is still counted once in the target's local rating (see [Analyze](#analyze)). A requirement declared directly under a target, with no nesting factor, MAY use this list to attach itself to one or more factors.

**Ratings (criterion overrides)**: an optional map that overrides the rating-scale criteria for this requirement, for use when the scale's shared criteria cannot express the gradient that matters. Each key MUST name a level of the model's rating scale; its value replaces that level's criterion for this requirement only. Overrides change the criteria, not the levels, their order, or their titles. Use them when a requirement has a natural measured threshold or a distinct qualitative spectrum:

```yaml
requirements:
  "p99 request latency stays within budget":
    assessment: >
      Measure p99 request latency over a representative production window.
    ratings:
      outstanding: "p99 at or under 150 ms."
      target: "p99 at or under 300 ms."
      minimum: "p99 at or under 500 ms."
      unacceptable: "p99 above 500 ms."
```

### Markdown Body

The Markdown body documents *why* this is the right quality model — the context a reader needs to interpret the requirements and an evaluator needs to weigh them. Where the frontmatter fixes *what* is assessed and *how it is rated*, the body explains *why these factors, why these requirements, and what matters most*. The [Analyze](#analyze) and [Advise](#advise) phases draw on it: when an evaluator weighs requirements by importance or names the key gaps behind a held-down rating, the stakeholder context that justifies those judgments lives here.

The body is OPTIONAL, and the format fixes no required section set. An author MAY use, rename, reorder, or replace the sections below. A conforming tool MUST preserve body content it does not recognize rather than dropping or rejecting it.

The sections below are RECOMMENDED, non-normative starting points:

| Section        | What it captures                                                                                                                  |
| -------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| **Overview**   | What the subject is, who depends on it, and what "good" means here.                                                               |
| **Scope**      | The model boundary: what it covers and deliberately leaves out. Out-of-scope concerns are exclusions by design, not deficiencies. |
| **Needs**      | Stakeholder outcomes the requirements answer to — the source of how much each requirement matters.                                |
| **Risks**      | What goes wrong, and for whom, if a need is not met.                                                                              |
| **Known gaps** | In-scope quality concerns deliberately deferred, each with a brief reason.                                                        |

Two distinctions keep the body and the evaluation from blurring. **Scope** is for concerns outside the model's remit; **Known gaps** is for in-scope concerns deliberately deferred. And **Known gaps** is the author's standing declaration, distinct from a *not assessed* outcome, which an evaluator determines per run when evidence is missing (see [Assess and Rate](#assess-and-rate)).

The body MAY open with a top-level heading; when present it SHOULD name the model's subject — matching the `title`, when set.

#### Example

```markdown
# Acme Checkout API — Quality model

## Overview

The Acme Checkout API accepts and settles customer payments. Good means every charge
is authenticated, charged exactly once, and protected from unauthorized access to
cardholder data.

## Scope

This model covers the API service and its payment-processor integration
(`./payments`). The upstream card network and the bank settlement system are out of
scope — Acme does not own them.

## Needs

- Customers are charged exactly once for what they bought.
- Cardholder data is never exposed to an unauthorized party.
- On-call engineers can confirm whether a charge settled.

## Risks

A double charge or a leaked card number is the worst outcome: both are unrecoverable
for the customer and carry regulatory exposure. An unpatched dependency is the most
likely path to the latter.

## Known gaps

- Sustained peak-load behavior is in scope but not yet modeled.
- Failed-charge reconciliation evidence is not yet produced, so that requirement
  currently evaluates as *not assessed*.
```

## Evaluation

Evaluation assesses a model's targets against their requirements, rates the evidence, rolls the results up the target tree, and advises on what to improve. It proceeds in five phases: **Define**, **Assess and Rate**, **Analyze**, **Advise**, and **Report**.

### Define

Determine the scope of the evaluation. By default the scope is the whole model: every target, and within each target every requirement. The scope MAY be narrowed by a filter:

- by **target** — restrict evaluation to a given target and its subtree;
- by **factor** — restrict evaluation to the requirements tied to a given factor (including those that tag it as a secondary factor); or
- both.

For every target in scope, resolve the **source** entities to be evaluated from the target's `source`. A narrowed scope qualifies every result that follows: ratings are understood within the scope, and a scoped evaluation MUST NOT be presented as a whole-model verdict.

### Assess and Rate

For each requirement in scope, evaluated against its target node:

1. **Assess** the source entities using the requirement's **assessment**, producing the **findings** — the evidence for this requirement. Each finding records an observation (e.g., a measured value, inspection note, or diagnostic result) and is not itself rated.
2. **Rate** the findings *together* against the rating scale criterion (or the requirement's criterion overrides), producing the requirement's **Rating Result**: a single rating level. When there are no findings or the evidence is insufficient to rate against the scale, the requirement MUST be recorded as **not assessed** rather than assigned a rating level.

This produces one Rating Result per requirement in scope.

### Analyze

Roll the requirement results up the model tree (requirement → factor → target → root, with sub-factors rolling up into their parent factor) by inference. Roll-up is not computed; an evaluator infers it by judgment against the rating scale. For each target — child targets first — the evaluator infers:

- a **factor rating**, for each of the target's factors — its sub-factors first, deepest first: the level that best characterizes the factor considering, together, the rating results of every requirement tied to it (both those nested under it and those that tag it as a secondary factor) and the ratings of its sub-factors. A sub-factor is rated the same way, and its rating rolls up into its parent factor — as a child target's aggregate rolls up into its parent;
- a **local rating**: the level that best characterizes the target considering all of its own requirement results together (each requirement counted once, whatever factors it touches); and
- an **aggregate rating**: the level that best characterizes the target considering its local rating together with the aggregate ratings of its child targets.

Factor ratings and the local rating are two reads over a target's requirements, but not always of the same set: the local rating is the whole-set verdict over the target's *own* requirements (each counted once, whatever factors it touches), while a factor rating is a cross-cutting lens over every requirement tied to that factor (a requirement may tag several) — including any declared on a descendant target that tags it as a secondary factor. A factor rating therefore MAY range wider than the local rating. A target with no requirements of its own (a grouping target) has no local rating; its aggregate considers only its child targets. A leaf target's aggregate rating equals its local rating.

Roll-up considerations (guidance for inference, not computation):

- The rolled-up level SHOULD reflect the target's overall state, but a serious shortfall in an important requirement SHOULD NOT be masked by many satisfactory ones — weigh requirements by how much each matters to the target's quality.
- Requirements, factors, or targets recorded as *not assessed* are excluded from the rating but MUST be noted. When too little has been assessed to responsibly infer a level, the roll-up is itself recorded as *not assessed*.
- Each roll-up SHOULD record a brief **rationale** naming what most determined the level (its binding constraints).

Tools MAY extend roll-up with explicit weights, thresholds, or computed aggregation; this standard prescribes only inferential judgment.

### Advise

From the analysis, the evaluator advises on improvement:

- **Key gaps** — the shortcomings most responsible for held-down ratings (the binding constraints surfaced during Analyze), together with any *not assessed* areas material to the verdict.
- **Options** — for each key gap, the available options for remediation.
- **Recommendation** — for each key gap, a recommended option and the rationale for it.

Advice is inferential and advisory; it informs the report but does not change any rating.

### Report

Produce the **Evaluation Report**: the structured result of the evaluation, suitable for rendering to a person, a gate, or a tool. A report presents at least:

- the **Rating** — the in-scope root target's aggregate rating — and its **rationale**;
- the **Scope** the rating was produced under;
- for each target in scope (root first, recursively): each requirement's findings summary, rating, and rationale; each factor's rating and rationale — including each sub-factor's rating and rationale at every depth; and the target's local and aggregate ratings, each with rationale; and
- the **Advice** — key gaps, options, and recommendations.

*Not assessed* outcomes MUST be shown wherever they occur, distinct from rated outcomes, at every level of the report.

Note: a factor's rating MAY draw on requirements declared on descendant targets that tag it as a secondary factor (see [Analyze](#analyze)). When it does, the report SHOULD make those contributing requirements identifiable under the factor, even though each is listed in full under the target that declares it.

## Appendix A: Sample Evaluation Report

This appendix is **non-normative**. It illustrates one rendering — for a human reader — of the Evaluation Report defined in **Report**. A tool MAY render the same underlying result differently (e.g., as JSON for a gate).

### The model evaluated

A condensed view of the `QUALITY.md` under evaluation, for reference. Its rating scale is the suggested four-level scale — **Outstanding** > **Target** > **Minimum** > **Unacceptable** — ordered best to worst.

- **Acme Checkout API** (root target, source `./`)
  - Factors:
    - **Security**, decomposed into sub-factors **Access control** and **Dependency integrity**
    - **Reliability**
  - Requirements:
    - *Every public endpoint requires authentication* (Security → Access control)
    - *No dependencies with known critical or high vulnerabilities* (Security → Dependency integrity)
    - *p99 request latency ≤ 300 ms* (Reliability)
    - *Automated tests cover the checkout flow end to end* (Reliability)
  - **Payment Processor** (child target, source `./payments`)
    - Factors: **Security**, **Reliability**
    - Requirements:
      - *Cardholder data is encrypted at rest* (Security)
      - *Gateway calls are idempotent on retry* (Security; tags **Reliability** as a secondary factor)
      - *Failed charges reconcile within 24 hours* (Reliability)

### The report

---

**Rating: Minimum** *(Acme Checkout API — aggregate, whole model)*

**Rationale.** Held at **Minimum** by a single binding constraint: an unpatched high-severity dependency vulnerability at the root. Every other area meets **Target** or better, and the Payment Processor subtree meets **Target** — so lifting the one Security gap would raise the overall rating.

**Scope.** Whole model; no target or factor filter applied. Source resolved from `./` (root) and `./payments` (Payment Processor). One requirement *not assessed* (see below).

---

#### Target: Acme Checkout API *(root)*

**Aggregate: Minimum** — the root's own local rating binds; the Payment Processor subtree (**Target**) does not pull it down further.
**Local: Minimum** — four root requirements; the dependency-vulnerability shortfall is security-critical and is not offset by the three meeting Target or better.

Factors:

- **Security — Minimum.** Bound by the **Dependency integrity** sub-factor; **Access control** is solid. Sub-factors:
  - **Access control — Target.** Every public endpoint sits behind authentication, with no exceptions; no evidence of the defense-in-depth Outstanding would require.
  - **Dependency integrity — Minimum.** One unpatched high-severity dependency advisory holds this sub-factor — and so the parent Security factor — below Target.
- **Reliability — Target.** Latency is well within budget; end-to-end test coverage meets but does not exceed the bar.

Requirements:

- *Every public endpoint requires authentication* — **Target**
  - *Findings:* 42 of 42 public routes sit behind the auth middleware; 0 unauthenticated routes found.
  - *Rationale:* Full coverage with no exceptions meets the Target criterion; no evidence of the defense-in-depth that Outstanding would require.
- *No dependencies with known critical or high vulnerabilities* — **Minimum**
  - *Findings:* 0 critical, 1 high-severity advisory in a transitive dependency; a patched release is available.
  - *Rationale:* No critical vulnerabilities clears Unacceptable, but an open high-severity advisory keeps it at Minimum.
- *p99 request latency ≤ 300 ms* — **Outstanding**
  - *Findings:* p99 measured at 180 ms over a 7-day window under production load.
  - *Rationale:* Comfortably inside the threshold with sustained margin.
- *Automated tests cover the checkout flow end to end* — **Target**
  - *Findings:* End-to-end suite covers the happy path and three failure paths; the partial-refund path is uncovered.
  - *Rationale:* Core flow is covered (Target); a known coverage gap keeps it short of Outstanding.

#### Target: Payment Processor

**Aggregate: Target** — a leaf target, so its aggregate equals its local rating.
**Local: Target** — both assessed requirements meet Target; one requirement is *not assessed* and is noted but excluded from the rating.

Factors:

- **Security — Target.** Encryption at rest and idempotent gateway calls both meet the bar.
- **Reliability — Target.** Idempotent gateway retries guard against double-charging on replay; the failed-charge reconciliation requirement is *not assessed*, so this rating rests on the idempotency evidence alone and is noted as incomplete.

Requirements:

- *Cardholder data is encrypted at rest* — **Target**
  - *Findings:* Card fields encrypted with AES-256; keys held in the managed KMS, rotated quarterly.
  - *Rationale:* Meets the encryption-at-rest criterion; no field-level exceptions found.
- *Gateway calls are idempotent on retry* — **Target** *(Security; also lensed under Reliability)*
  - *Findings:* Idempotency keys present on all charge calls; a replay test confirmed no double-charge.
  - *Rationale:* Meets the criterion. Counts once in the local rating while informing both the Security and Reliability lenses.
- *Failed charges reconcile within 24 hours* — **Not assessed**
  - *Findings:* None — no reconciliation report or job output was available to assess against.
  - *Rationale:* Insufficient evidence to rate; recorded as *not assessed* rather than assigned a level.

---

#### Advice

- **Key gap — open high-severity dependency vulnerability (root → Security).** The single constraint holding the whole model at Minimum.
  - *Options:* (a) upgrade the dependency to the patched release; (b) replace it with an unaffected library; (c) accept the risk with a compensating control and a tracked exception.
  - *Recommended:* **(a) upgrade to the patched release** — lowest effort, removes the binding constraint, and is expected to lift root Security and the overall rating to Target.
- **Coverage gap — reconciliation requirement not assessed (Payment Processor → Reliability).** The Payment Processor rating is incomplete until this is evaluated.
  - *Options:* (a) stand up the reconciliation report so the requirement can be assessed; (b) narrow or retire the requirement if reconciliation is out of scope.
  - *Recommended:* **(a) produce the reconciliation evidence** so the requirement can be rated and the subtree's rating reflects full coverage.
- **Minor — partial-refund path uncovered (root → Reliability).** Not rating-binding today; addressing it would move the end-to-end test requirement toward Outstanding.
