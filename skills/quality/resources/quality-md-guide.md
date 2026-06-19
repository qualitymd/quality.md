# Authoring QUALITY.md

A single, comprehensive guide to understanding and working with `QUALITY.md`
files: the concepts that make up the format, and the jobs you do when working
with each.

This guide conforms to [`SPECIFICATION.md`](SPECIFICATION.md). The specification governs on any conflict.

## Contents

- [The QUALITY.md file](#the-qualitymd-file)
- [Quality Model](#quality-model)
- [Rating Scale](#rating-scale)
- [Target](#target)
- [Factor](#factor)
- [Requirement](#requirement)
- [The Markdown body](#the-markdown-body)
- [When to update QUALITY.md](#when-to-update-qualitymd)

---

## The QUALITY.md file

A `QUALITY.md` file is a Markdown file with two parts:

- **YAML frontmatter** — the **quality model**: a structured, declarative
  description of what quality means for the entity being evaluated.
- **Markdown body** — the context a reader needs to interpret the model: why
  these factors, why these requirements, what matters most.

The whole file represents a single apex **target** — the top entity whose
quality is modeled. Everything else nests beneath it.

The file's location carries meaning: a `QUALITY.md` in a directory makes that
directory and everything under it (`**/*`) the default scope of evaluation,
unless a target narrows it with a `source`.

### Working with the file

#### Place the file at the root of what it evaluates

- **Do** put `QUALITY.md` in the directory whose contents are the subject target. The
  default scope is that directory and all subdirectories — no `source` needed in
  the common case. *Location is the simplest, most reliable scope declaration.*

---

## Quality Model

The **quality model** is the frontmatter: the root node of the file. It carries
one model-wide property, the `ratingScale`, plus all the properties of a
[Target](#target) — because the model root *is* the apex target and tops the
target tree.

### Properties

| Property       | Presence    | What it is                                                               |
| -------------- | ----------- | ------------------------------------------------------------------------ |
| `title`        | Recommended | Human-readable name of the entity whose quality is modeled.              |
| `description`  | Optional    | A concise statement of what the model's target is.                       |
| `ratingScale`  | Required    | The [rating scale](#rating-scale) — the levels every result is rated on. |
| `factors`      | Optional\*  | The [factors](#factor) — lenses through which quality is described.      |
| `requirements` | Optional\*  | The [requirements](#requirement) assessed against the root source.       |
| `targets`      | Optional\*  | Child [targets](#target) — more focused models nested to any depth.      |
| `source`       | Optional    | Scope override; omit at the root to take the directory default.          |

\* At least one of `factors`, `requirements`, or `targets` must be present.

`ratingScale` is the only property unique to the model; everything else it
shares with Target.

### Working with the model

#### Start from the subject, not the structure

- **Do** name the entity (`title`) and think through the body's Overview first —
  what the thing is and what "good" means for it — before listing factors.
  *Factors and requirements are answers; you need the question first.*

#### Keep the root lean when child targets carry the detail

- **Consider** declaring only model-wide factors at the root and pushing
  narrower factors/requirements down to child targets. *A flat root with
  everything at one level loses the structure that makes a model readable.*
- **Avoid** modeling every property at the root "to be safe." *An entry on
  factors, requirements, **or** targets is enough; empty scaffolding is noise.*

---

## Rating Scale

The **rating scale** is the fixed set of levels every requirement result is
rated against — the model's shared vocabulary for "how good." It is a list of
**rating levels**, ordered best (first) to worst (last), with at least two
levels.

Each level does two distinct jobs through two properties:

- **`description`** — what the level *means*: its standing in the scale and its
  intent. Fixed for the whole model; never overridden.
- **`criterion`** — the default rule for deciding whether a requirement's
  findings *land at* that level. A requirement may override its own criterion
  (see [Override the criteria for one
  requirement](#override-criteria-only-when-the-shared-scale-cant-express-the-gradient));
  the description never changes.

### Properties (per level)

| Property      | Presence    | What it is                                                      |
| ------------- | ----------- | --------------------------------------------------------------- |
| `level`       | Required    | The level's name; unique within the scale.                      |
| `title`       | Optional    | Human-readable label for reports.                               |
| `description` | Recommended | What the level means across the model (fixed).                  |
| `criterion`   | Required    | Default rule for rating a requirement's findings at this level. |

### Working with the rating scale

#### Reuse the suggested four-level scale unless you have a reason not to

- **Consider** the **outstanding > target > minimum > unacceptable** scale as a
  default — a stretch level, the level to aim for, the floor you've agreed to
  ship at, and below the floor. *A shared four-band vocabulary is enough for most
  models and keeps reports comparable.*
- **Do** choose a different scale only when the subject demands it (e.g. a
  pass/fail gate wants two levels). *The scale should fit how decisions are
  actually made about this entity.*

#### Keep `description` about meaning and `criterion` about rating

- **Do** write `description` as what the level *is* ("the floor you've agreed to
  ship at") and `criterion` as the test a result must pass ("falls short of
  target but remains acceptable"). *Conflating them makes per-requirement
  criterion overrides impossible to write cleanly.*
- **Avoid** putting thresholds or measurements in `description`. *Those belong in
  `criterion`, where a requirement can override them.*

---

## Target

A **target** is an entity, or set of entities, with quality requirements subject
to evaluation. It is the recursive node of the model: the root model is the apex
target, and every entry under `targets` is another target, nested to any depth.

A target's **source** defines *what* it evaluates; its **factors** and
**requirements** define *what is assessed* about it. A target may also be a pure
**grouping target** — only child targets, no requirements of its own.

### Properties

| Property       | Presence    | What it is                                                             |
| -------------- | ----------- | ---------------------------------------------------------------------- |
| `title`        | Recommended | Display name; overrides the map key in reports.                        |
| `description`  | Optional    | What the target is — the entity or scope it covers.                    |
| `factors`      | Optional\*  | [Factors](#factor) scoped to this target's subtree.                    |
| `requirements` | Optional\*  | [Requirements](#requirement) assessed against this target's source.    |
| `targets`      | Optional\*  | Child targets, nested to any depth.                                    |
| `source`       | Optional    | Where the entities live; inherits the nearest ancestor's when omitted. |

\* A target may declare none of its own and serve purely as a grouping node, but
each target should lead to at least one requirement somewhere in its subtree — a
target whose subtree has no requirements evaluates nothing.

### Working with targets

#### Target the primary artifact, not its supporting cast

- **Do** point a target at the thing whose quality you actually care about.
  **Avoid** modeling fixtures, generated output, or build scaffolding as targets.
  *A model of the test fixtures rates the wrong thing.*

#### Split off a child target only when it has distinct factors or requirements

- **Consider** a child target when a part of the subject is best evaluated
  through its own factors or requirements that would not apply to its siblings
  (e.g. a payment integration with its own security concerns).
- **Avoid** child targets that merely re-state the parent's requirements. *Child
  targets don't inherit the parent's requirements, but duplicating them by hand
  is a maintenance trap — let the ancestor's source cover the shared scope.*

#### Use a grouping target to organize, not to evaluate

- **Do** leave a grouping target's `source` implicit and let its children narrow
  it. *A grouping target has no local rating; its job is to hold child targets,
  and its aggregate reflects only them.*

#### Define `source` only to narrow or relocate

- **Do** omit `source` on the root to take the directory default, and on a child
  to inherit the nearest ancestor's scope.
- **Do** set `source` (a path or glob, relative to the file) when a target
  evaluates a specific subtree. *Be aware an ancestor's source may overlap a
  child's, so the ancestor's requirements can also land on the child's entities.*

#### Write a description that distinguishes, not enumerates

- **Do** state what the target *is* and how it differs from its siblings/parent.
  **Avoid** restating its factors or requirements in the description. *The
  description fixes what is evaluated; factors and requirements fix what is
  assessed about it.*

---

## Factor

A **factor** is a quality characteristic — a *lens* such as `reliability`,
`security`, or `maintainability` — through which a target's quality is described.
It groups the requirements assessed through it. A factor may decompose into
**sub-factors**: finer characteristics that together make up the parent (e.g.
`reliability` into `availability`, `fault tolerance`, `recoverability`). A
sub-factor is itself a factor of the same shape, nested to any depth.

Factor identity is local to its target: factors of the same name on two
different targets are distinct.

### Properties

| Property       | Presence    | What it is                                                     |
| -------------- | ----------- | -------------------------------------------------------------- |
| `description`  | Recommended | The characteristic, defined operationally for this entity.     |
| `factors`      | Optional    | Sub-factors — recursively a factor.                            |
| `requirements` | Optional    | [Requirements](#requirement) uniquely relevant to this factor. |

### Working with factors

#### Choose factors that name the concerns that matter here

- **Do** pick the handful of quality characteristics that genuinely drive this
  entity's quality. **Avoid** importing a standard checklist of characteristics
  wholesale. *A factor nothing contributes to is a lens over nothing.*

#### Write the description as an operational definition

- **Do** define the characteristic as the *degree or capability to achieve some
  end under the conditions that matter*, and say why it matters and to whom. A
  useful shape: *"\<Factor\> is the degree to which \<entity\> \<achieves some
  end\> under \<conditions\>; it matters because \<stakeholder concern\>."*
- **Avoid** an adjective or a synonym for the factor name ("Reliability: how
  reliable it is"). *That tells a reader nothing and doesn't distinguish it from
  its siblings.*

#### Make sibling factors non-overlapping

- **Do** write each factor's description so the factors on a target read as a
  distinct, non-overlapping set. *Overlapping factors make it ambiguous where a
  requirement belongs and double-count concerns in roll-up.*

#### Decompose into sub-factors only when it aids understanding

- **Consider** sub-factors when a factor carries more than one distinct concern
  that's clearer assessed apart than together.
- **Avoid** decomposing a factor whose requirements already speak for themselves.
  *Decompose only as far as it helps; the same description and
  distinct-from-siblings rules apply at every level.*

#### Refine, don't redefine, an ancestor's factor

- **Do** treat a child target's factor that shares a name with an ancestor's as a
  *refinement* tailored to the child. *They're technically distinct factors;
  write the child's description to say how it specializes the ancestor's
  concern.*

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
  requirement placed directly under a target, `factors` is **required** and names
  the factors directly.

However it is connected, a requirement is assessed once, against the source of
the target it sits on, and counts once in that target's local rating.

### Properties

| Property     | Presence                       | What it is                                                       |
| ------------ | ------------------------------ | ---------------------------------------------------------------- |
| `assessment` | Required                       | The means of assessing the source; produces the findings.        |
| `factors`    | Required for target-level reqs | Factor references; secondary factors when nested under a factor. |
| `ratings`    | Optional                       | Per-requirement criterion overrides, keyed by rating level.      |

### Working with requirements

#### Write the statement as the expectation itself

- **Do** phrase the map key as the thing you expect to be true ("Every public
  endpoint requires authentication", "p99 request latency stays within budget").
  *The statement is what shows up in reports; it should read as a claim that can
  be true or false to a degree.*

#### Make each requirement meaningful on its own

- **Do** write requirements specific enough that a single result stands on its
  own. *A vague requirement produces a vague rating and a weak report.*

#### Give each requirement exactly one assessment

- **Do** declare one `assessment` — inline criteria, a measurement procedure, an
  inspection checklist, a diagnostic, or a path to a doc describing one.
- **Avoid** stacking several independent assessments under one statement. *Split
  it into separate requirements instead — each result must be independently
  ratable.*

#### Reference assessment sources; don't copy them

- **Do** point at the spec, doc, or checklist that defines the assessment, naming
  it once.
- **Avoid** extracting, summarizing, or duplicating that content into the
  requirement. *Duplicated criteria drift out of sync with their source.*

#### Connect to factors deliberately

- **Do** rely on **placement** for the primary factor (nest the requirement under
  it), and add `factors` only to pull the result into additional (**secondary**)
  factor roll-ups.
- **Do** declare `factors` explicitly for a requirement placed directly under a
  target — it's required there, and `null`/`[]`/empty entries don't satisfy it.

#### Override criteria only when the shared scale can't express the gradient

- **Consider** a `ratings` override when a requirement has a natural measured
  threshold or a distinct qualitative spectrum (e.g. latency bands).
- **Do** key overrides by existing rating levels and change *only* the
  `criterion`. **Avoid** touching a level's `description`, order, or `title` —
  those stay fixed across the model.

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

---

## The Markdown body

The body documents *why* this is the right model — the context an evaluator
needs to weigh the requirements. Where the frontmatter fixes *what* is assessed
and *how it's rated*, the body explains *why these factors, why these
requirements, and what matters most*.

The body is optional and fixes no required sections; you may rename, reorder, or
replace these. They're recommended starting points:

| Section        | What it captures                                                                       |
| -------------- | -------------------------------------------------------------------------------------- |
| **Overview**   | What the subject is, who depends on it, and what "good" means here.                    |
| **Scope**      | What the model covers and deliberately leaves out. Out-of-scope ≠ deficiency.          |
| **Needs**      | Stakeholder outcomes the requirements answer to — the source of how much each matters. |
| **Risks**      | What goes wrong, and for whom, if a need isn't met.                                    |
| **Known gaps** | In-scope concerns deliberately deferred, each with a brief reason.                     |

### Working with the body

#### Write the body to justify the model's weighting

- **Do** capture, in Needs and Risks, why some requirements matter more than
  others. *An evaluation weighs requirements by importance and names key gaps —
  both lean on the body for that judgment.*

#### Keep Scope, Known gaps, and "not assessed" distinct

- **Do** use **Scope** for concerns outside the model's remit, and **Known gaps**
  for in-scope concerns you've deliberately deferred.
- **Do** record genuinely in-scope-but-deferred concerns as **Known gaps**. *A
  requirement with no available evidence will rate as not assessed — better to
  declare the gap than to surprise a reader with it.*
- **Don't** confuse **Known gaps** (your standing declaration) with a **not
  assessed** result (an evaluator's per-run finding that evidence was missing).
  *They look similar but mean different things in a report.*

---

## When to update QUALITY.md

#### Update when you learn something that changes the model

- **Do** revise when a discovery changes the context or content of the
  evaluation — a new factor that matters, a requirement whose assessment changed,
  a scope that shifted.
- **Do** keep the body current with the frontmatter. *A model whose body no
  longer explains its factors misleads the next evaluator.*
