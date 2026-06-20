# Authoring QUALITY.md

The canonical reference and best-practices guide for understanding and working
with `QUALITY.md` files: the format concepts, authoring practices, and jobs
attached to each.

This guide conforms to [`SPECIFICATION.md`](../resources/SPECIFICATION.md). The
specification governs on any conflict.

## Contents

- [The `QUALITY.md` file](#the-qualitymd-file)
- [Quality Model](#quality-model)
- [The Markdown body](#the-markdown-body)
- [Rating Scale](#rating-scale)
- [Target](#target)
- [Factor](#factor)
- [Requirement](#requirement)
- [When to update QUALITY.md](#when-to-update-qualitymd)

---

## The `QUALITY.md` file

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

- **Do** put `QUALITY.md` in the directory whose contents are the subject target. The
  default scope is that directory and all subdirectories — no `source` needed in
  the common case.
- **Do** keep the root `title`, body Overview and Scope, file location, and
  root `source` aligned on the same evaluated subject. *If those disagree, the
  model may be valid YAML but misleading about what is actually being judged.*

---

## Quality Model

The **quality model** is the frontmatter: the root node of the file. It carries
one model-wide property, the `ratingScale`, plus all the properties of a
[Target](#target) — because the model root *is* the apex target and tops the
target tree.

### Properties

| Property       | Presence   | What it is                                                               |
| -------------- | ---------- | ------------------------------------------------------------------------ |
| `title`        | Required   | Human-readable name of the entity whose quality is modeled.              |
| `description`  | Optional   | A concise statement of what the model's target is.                       |
| `ratingScale`  | Required   | The [rating scale](#rating-scale) — the levels every result is rated on. |
| `factors`      | Optional\* | The [factors](#factor) — lenses through which quality is described.      |
| `requirements` | Optional\* | The [requirements](#requirement) assessed against the root source.       |
| `targets`      | Optional\* | Child [targets](#target) — more focused models nested to any depth.      |
| `source`       | Optional   | Scope override; omit at the root to take the directory default.          |

\* At least one of `factors`, `requirements`, or `targets` must be present.

`ratingScale` is the only property unique to the model; everything else it
shares with Target.

### Working with the model

- **Do** name the entity (`title`) and think through the body's Overview first —
  what the thing is and what "good" means for it — before listing factors.

#### Build from context, then scale, then model tree

- **Do** fill the Markdown body before expanding the frontmatter. *The body is
  where the model gets its judgment context: what the subject is, what decisions
  the model supports, what quality means here, and which risks matter enough to
  assess.*
- **Do** confirm the rating scale after the body and before writing Requirements.
  *The scale is the shared vocabulary for turning future findings into ratings;
  Requirements are easier to write once "unacceptable", "minimum", "target", and
  "outstanding" mean something for this subject.*
- **Do** derive Factors and Requirements from the body context. *Factors and
  Requirements should express the subject's needs, risks, scope, and known gaps,
  not lead them.*
- **Do** trace at least one important concern from body to model before expanding
  the tree. *A useful trace reads like: a Need names the outcome, a Risk names
  the failure mode, a Factor names the quality lens, and a Requirement names the
  inspectable expectation.*

#### Keep the root lean when child targets carry the detail

- **Consider** declaring only model-wide factors at the root and pushing
  narrower factors/requirements down to child targets. *A flat root with
  everything at one level is harder to read and maintain.*
- **Avoid** modeling every property at the root "to be safe." *An entry on
  factors, requirements, **or** targets is enough.*
- **Do** make sure the model reaches Requirements somewhere in the target tree
  before treating it as evaluable. *Factors without Requirements can describe
  concerns, but they do not give an evaluator anything to rate.*

#### Make the traceability graph visible

- **Do** treat the model as a graph: a tree of targets plus the assessment
  references that link a requirement to the entity supplying its criteria. Which
  entity is the criteria for which — and how quality depends along those edges —
  is often the most valuable thing the model records.
- **Do** make each assessment either inline or a reference to another entity, and
  reference that entity by the same canonical path used as its target `source`,
  so a reader can follow the edge from one target to the next.
- **Do** stop the chain where ownership or value stops: model referenced entities
  as targets while you own them; let a guide that governs its own kind be its own
  assessment rather than adding another layer.

---

## The Markdown body

The body gives the context an evaluator needs to interpret and weigh the model:
why these factors, why these requirements, and what matters most.

The body is optional and fixes no required sections; you may rename, reorder, or
replace these. They're recommended starting points:

| Section        | Desired outcome                                                                                                                                                       |
| -------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Overview**   | A reader can say what the subject is, who depends on it, and why its quality matters. This names the real entity, not just the repo or file where `QUALITY.md` lives. |
| **Scope**      | A reader can tell what is included, what is excluded for now, and where the model boundary sits. Out-of-scope is not a deficiency.                                    |
| **Needs**      | A reader can see the outcomes the subject must support and the users, operators, maintainers, or downstream systems those outcomes serve.                             |
| **Risks**      | A reader can see the failures that would make the subject untrustworthy, unusable, unsafe, expensive, or hard to change. These are raw material for initial Factors.  |
| **Known gaps** | A reader can see current weaknesses, weak evidence, and known unknowns: missing context, unresolved questions, and evidence gaps not fully identified elsewhere.      |

### Working with the body

- **Do** write the body before expanding the model tree. *It is the fastest way
  to discover what Factors and Requirements the frontmatter should express.*
- **Do** capture, in Needs and Risks, why some requirements matter more than
  others. *Importance and gaps both depend on this context.*
- **Do** use the body to explain any rating-scale change. *A custom scale should
  answer a real decision need visible in the body.*
- **Do** treat Needs as *benefits to realize*, not only outcomes to protect, and
  Risks as the problems that erode them. *A model that lists only failure modes
  can rate every requirement at target while the subject is still not worth
  shipping — because nothing weighs whether the benefits, on the whole, outweigh
  the residual problems.*
- **Consider** separating two reasons a concern matters: how much its failure
  hurts (stakeholder, safety, or business stakes) and how far its failure spreads
  (how much else it forces to change). *A concern can be low-stakes yet
  high-blast-radius, or the reverse; saying which helps the next evaluator weigh
  roll-ups.*
- **Consider** noting where two Factors or Requirements pull against each other
  (tighter access control vs. faster onboarding, latency vs. cost) and which way
  you have chosen to lean. *A model that hides its trade-offs invites an evaluator
  to "fix" a deliberate compromise.*

#### Say which sense of "good" this model uses

Quality is not one thing. A subject can be judged by **conformance** (does it
match its specification?), by **fitness for purpose** (does it serve the user's
real need?), or by **value** (is it worth its cost?). These can disagree — a
subject can meet its spec yet fail the need, or serve the need while departing
from spec.

- **Do** name the governing sense of "good" in the Overview, so a reader knows
  whether a passing model means "conforms" or "fits."
- **Do** make both visible where stakeholders would disagree — "meets the spec"
  vs. "meets the need" — rather than letting one silently win. *Different
  stakeholders rate the same finding differently; record the contested
  expectation instead of burying it in a single criterion.*
- **Do** attribute the model's judgments to the stakeholders they serve (users,
  operators, maintainers, downstream systems). *A quality no named stakeholder
  would miss is rarely worth modeling.*

#### Keep Scope, Known gaps, and "not assessed" distinct

- **Do** use **Scope** for concerns outside the model's remit, and **Known gaps**
  for in-scope concerns you've deliberately deferred or cannot yet define.
- **Do** record known unknowns as **Known gaps** when missing context, unresolved
  questions, or weak evidence prevented the rest of the body from being fully
  identified. *This keeps uncertainty visible without pretending it has already
  been evaluated.*
- **Do** record genuinely in-scope-but-deferred concerns as **Known gaps**. *A
  declared gap is clearer than a surprise not-assessed result.*
- **Don't** confuse **Known gaps** (your standing declaration) with a **not
  assessed** result (an evaluator's per-run finding that evidence was missing).
- **Do** record low *confidence* in an assessment, not only its absence. *A
  requirement rated `target` on one stale benchmark or a single reviewer differs
  from one rated on sustained evidence; note that fragility in Known gaps. "Rated
  but barely trusted" is neither "not assessed" nor "no gap," and it is often
  where the next evaluation should look first.*

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
| `title`       | Required    | Human-readable label for reports.                               |
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
- **Do** review the scale after writing the body. *The body should reveal whether
  "good enough" and "excellent" need sharper meaning for this subject.*
- **Do** test adjacent levels against the same plausible finding. *If one
  finding could reasonably satisfy both `target` and `minimum`, sharpen the
  criteria until the boundary is usable.*
- **Do** let the model's job pick the number of levels: a model that gates against
  an agreed bar can use two levels (pass / fail); a model that judges *how good*,
  compares options, or surfaces strengths and weaknesses earns a graded scale. *A
  pass/fail model deliberately gives up the ability to say "how good" — sometimes
  exactly what you want, sometimes a loss.*
- **Do** fix the *required margin* in the body: how far above `minimum` this
  subject must actually land to be good enough, and why (safety, reputation,
  regulatory exposure, throwaway tool). *The same four levels serve a life-critical
  system and a one-off script; what differs is the band the subject must reach.
  State it so an evaluator knows whether an all-`minimum` result is fine or
  alarming.*
- **Do** calibrate the levels against concrete exemplars where you can — a real
  artifact you would call `outstanding` and one you would call `unacceptable`.
  *Thresholds set against known cases are far more defensible than ones invented
  in the abstract; if no real artifact would ever earn `outstanding`, the level is
  decorative.*
- **Avoid** inventing a custom scale before the body reveals a real need for
  one.

#### Keep `description` about meaning and `criterion` about rating

- **Do** write `description` as what the level *is* ("the floor you've agreed to
  ship at") and `criterion` as the test a result must pass ("falls short of
  target but remains acceptable"). *Conflating them makes per-requirement
  criterion overrides impossible to write cleanly.*
- **Avoid** putting thresholds or measurements in `description`. *Those belong in
  `criterion`, where a requirement can override them.*
- **Do** prefer a *measurable* boundary in `criterion` whenever the concern admits
  one — a threshold, a count, a proportion, a band. *Unclear, unmeasurable criteria
  are a leading source of requirements debt; a measurable boundary is what makes a
  finding land at one level rather than hover between two.*
- **Avoid** subjective or comparative wording in a `criterion` ("user-friendly",
  "better than", "as appropriate", "minimal"). *Such terms cannot reliably separate
  two adjacent levels; if you cannot operationalize a criterion, treat it as a
  smell pointing back at the requirement statement.*

---

## Target

A **target** is an entity, or set of entities, with quality requirements subject
to evaluation. It is the recursive node of the model: the root model is the apex
target, and every entry under `targets` is another target, nested to any depth.

An **entity** is anything evaluated for quality. Most often it is a concrete
artifact — a document, a source tree, a config file — but it can also be a
service, a behavior, or a process. The same entity can be a target in one
relationship and the reference supplying another requirement's criteria in
another.

A target's **source** defines *what* it evaluates; its **factors** and
**requirements** define *what is assessed* about it. A target may also be a pure
**grouping target** — only child targets, no requirements of its own.

### Properties

| Property       | Presence   | What it is                                                             |
| -------------- | ---------- | ---------------------------------------------------------------------- |
| `title`        | Required   | Display name in reports; the map key stays the identifier.             |
| `description`  | Optional   | What the target is — the entity or scope it covers.                    |
| `factors`      | Optional\* | [Factors](#factor) scoped to this target's subtree.                    |
| `requirements` | Optional\* | [Requirements](#requirement) assessed against this target's source.    |
| `targets`      | Optional\* | Child targets, nested to any depth.                                    |
| `source`       | Optional   | Where the entities live; inherits the nearest ancestor's when omitted. |

\* At least one of `factors`, `requirements`, or `targets` must be present, same
as the model. A target with only child `targets` is a grouping node; even then,
each target should lead to at least one requirement somewhere in its subtree — a
target whose subtree has no requirements evaluates nothing.

### Working with targets

- **Do** point a target at the thing whose quality you actually care about.
  **Avoid** modeling fixtures, generated output, or build scaffolding as targets.
  *Those usually support evaluation; they are not the subject being evaluated.*

#### Choose targets that are authored, inspectable artifacts

- **Do** prefer targets that are authored and inspectable — source trees, specs,
  tests, docs, configs, checklists — over runtime behavior or process you cannot
  open and examine. *The model is judged by inspecting the entity; point it at
  something an evaluator can actually read.*
- **Consider** whether the entity is **constitutive** (it *is* the product — the
  source, the document) or **normative** (it *governs* the product — a spec, style
  guide, or checklist others are judged against). *Normative artifacts are
  high-leverage: quality invested once propagates to everything they govern, often
  the fastest early win for a low-maturity subject.*
- **Do** match the target's grain to the concern: a single test and a test *suite*
  are different targets. *A collection is judged by whole-set concerns no member
  has — coverage, balance, non-redundancy, navigability. Point `source` at the
  whole set and write requirements at that grain.*

#### Split off a child target only when it has distinct factors or requirements

- **Consider** a child target when a part of the subject is best evaluated
  through its own factors or requirements that would not apply to its siblings
  (e.g. a payment integration with its own security concerns).
- **Avoid** child targets that merely re-state the parent's requirements. *Child
  targets don't inherit the parent's requirements, but duplicating them by hand
  is a maintenance trap — let the ancestor's source cover the shared scope.*
- **Do** place a requirement on the most specific target whose source actually
  owns the concern, not on a broad ancestor that merely contains it. *Assessing a
  payment-integration concern at the root re-judges it at every level and blurs
  who owns the finding; a second evaluator should be able to place the same
  finding at the same target.*
- **Avoid** splitting off a child whose rating is meaningless on its own. *If it is
  only interpretable alongside its siblings, keep the concern in the parent; split
  when the child could be evaluated, and trusted, in isolation.*

- **Do** leave a grouping target's `source` implicit and let its children narrow
  it. *A grouping target has no local rating; its aggregate reflects only its
  children.*

#### Define `source` only to narrow or relocate

- **Do** omit `source` on the root to take the directory default, and on a child
  to inherit the nearest ancestor's scope.
- **Do** set `source` (a path or glob, relative to the file) when a target
  evaluates a specific subtree. *An ancestor's source may overlap a child's; even
  so, each requirement evaluates against its own target's source — the overlap
  just means the child's entities also fall within the ancestor's scope.*
- **Do** compare each `source` to the target's title and description. *The source
  should select the entity being modeled, not a convenient proxy, generated
  output, or a broader repository area that happens to contain it.*

#### An entity can be both a target and an assessment reference

- **Do** decide by *role in this relationship*, not by identity: an entity is a
  target where its quality is the thing being managed, and an assessment
  reference where it supplies the criteria for judging another target. A spec is
  both — a target in its own right and the criteria for judging the code that
  implements it.
- **Do** point `source` at the whole entity, not a thin entry point or proxy.
- **Avoid** treating "target or reference?" as exclusive. *An owned entity used
  as criteria usually deserves a target of its own as well.*
- **Do** spend early effort on the normative artifacts that play both roles — a
  spec, standard, or style guide. *Improving one raises every target judged against
  it, so it is often the highest-leverage place to invest quality early.*

#### Write a description that distinguishes, not enumerates

- **Do** state what the target *is* and how it differs from its siblings/parent.
  **Avoid** restating its factors or requirements in the description. *The
  description identifies the evaluated entity; factors and requirements define
  what is assessed about it.*

#### Decide how ratings roll up

A target's local rating and a factor's rating are inferred from the requirement
results beneath them. The format fixes no aggregation formula, so how those
results combine is a modeling decision you should make and communicate.

- **Do** decide and state how a target's requirements combine into its rating when
  it is not obvious. *Two defensible defaults: **worst-of** (the weakest finding
  sets the rating — right when any unacceptable requirement makes the whole
  untrustworthy) and **most-common / median** (right when requirements genuinely
  compensate for one another).*
- **Avoid** implying that rating levels *average*. *The scale is ordinal
  (outstanding > target > minimum > unacceptable); the arithmetic mean of ordinal
  levels has no meaning — three `target`s and one `unacceptable` is not "slightly
  above minimum."*
- **Do** identify the requirements that can **veto** a rating — a single
  `unacceptable` finding that makes the subject not good enough no matter how
  strong everything else is (secrets stored in plaintext, data loss on failover).
  *Most requirements trade off against each other; a few cap the target. Sharpen a
  veto requirement's `unacceptable` criterion and name its role in the body, since
  a compensating roll-up can otherwise hide a critical problem behind strong
  siblings.*
- **Do** say when requirements are not equally important, so a reader does not read
  the roll-up as one-vote-each.
- **Do** keep a `not assessed` result distinct from a low rating in roll-up.
  *Missing evidence is not a failure; it must stay visible rather than count as a
  weak pass.*

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
| `title`        | Required    | Human-readable label for reports and status output.            |
| `description`  | Recommended | The characteristic, defined operationally for this entity.     |
| `factors`      | Optional    | Sub-factors — recursively a factor.                            |
| `requirements` | Optional    | [Requirements](#requirement) uniquely relevant to this factor. |

### Working with factors

#### Choose factors that name the concerns that matter here

- **Do** pick the handful of quality characteristics that genuinely drive this
  entity's quality. **Avoid** importing a standard checklist of characteristics
  wholesale.
- **Do** derive initial factors from the body's Needs and Risks. *Needs point at
  the outcomes quality should preserve; Risks point at the failure modes worth
  assessing.*
- **Do** reconcile major Needs and Risks back to Factors after drafting them.
  *If an important concern has no Factor, either add the Factor, mark the concern
  out of scope, or record the unresolved concern in Known gaps.*
- **Do** justify each Factor by something concrete about *this* subject — who
  depends on it, what it is for, where it runs — not by its presence on a general
  list. *A characteristic with no user, no failure mode, and no decision riding on
  it here does not earn a Factor, however standard it is elsewhere; pull from a
  catalog as a prompt, never as a quota.*
- **Do** anchor each Factor to a stakeholder whose concern it carries — the user
  who needs it to work, the maintainer who needs to change it, the operator who
  runs it. *Where stakeholders disagree on what "good enough" means, surface the
  conflict rather than averaging it into one criterion.*
- **Consider** whether a Factor names something a stakeholder *experiences* (the
  system is available, decisions are correct) or something *internal* that matters
  only because it produces that experience (low coupling, clear structure). *Keep
  internal factors tied to the outcome they serve, so an evaluator can tell a real
  weakness from a stylistic preference.*

#### Write the description as an operational definition

- **Do** define the characteristic as the *degree or capability to achieve some
  end under the conditions that matter*, and say why it matters and to whom. A
  useful shape: *"\<Factor\> is the degree to which \<entity\> \<achieves some
  end\> under \<conditions\>; it matters because \<stakeholder concern\>."*
- **Avoid** an adjective or a synonym for the factor name ("Reliability: how
  reliable it is"). *That tells a reader nothing and doesn't distinguish it from
  its siblings.*

- **Do** write each factor's description so the factors on a target read as a
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

#### Operationalize a factor into an assessable property

A target does not directly manifest a Factor like "reliability"; it exhibits
concrete properties that imply it. Bridging from the abstract characteristic to
something inspectable is the core authoring move.

- **Do** name a concrete, observable property of the entity that the Factor
  depends on, then phrase the requirement as the expectation about that property.
  *Reliability (abstract) → "recovers from a dependency outage without data loss"
  (property) → the requirement and its assessment.*
- **Do** preflight each requirement by naming its **scale** (the dimension and
  unit a finding lands on — *time for a novice to complete a one-item order*), its
  **meter** (the agreed procedure that produces the finding — *median over 100 test
  users using only online help*), and the rating boundary it tests. *If you cannot
  name a scale and a meter, the statement is still a slogan — "easy to use" has
  neither.*
- **Avoid** jumping straight from a Factor to a convenient metric with no property
  in between. *That is how you get measurements no one can tie back to why they
  matter.*

#### Write ratable requirement statements

- **Do** phrase the map key as the thing you expect to be true ("Every public
  endpoint requires authentication", "p99 request latency stays within budget").
  *The statement is what shows up in reports; it should read as a claim that can
  be true or false to a degree.*

- **Do** write requirements specific enough that a single result stands on its
  own. *A vague requirement produces a vague rating.*
- **Do** write Requirements that make the body context assessable. *A Requirement
  should turn an important need, risk, or known gap into an expectation an
  evaluator can inspect.*
- **Do**, for behavioral qualities (reliability, recoverability, security under
  attack), phrase the statement around the *triggering condition and operating
  environment*, not just the steady state: "When a downstream dependency times out,
  requests fail over within 2 s and surface a degraded-mode response." *A bare "is
  reliable" hides the condition that makes the quality observable.*
- **Do** make every requirement *verifiable*: its assessment must name a concrete
  method by which a finding could be produced. *If you cannot state how the source
  would be examined, the requirement is not assessable yet, however well it reads.*
- **Do** name the evidence the assessment draws on and check it can actually speak
  to the claim. *A latency claim needs runtime telemetry; a structure claim needs
  the source. If the only available evidence cannot address the claim, the
  requirement returns "not assessed" no matter how well written — narrow it to what
  the evidence supports, or record the gap in Known gaps.*
- **Do** apply the discard test: if a requirement were deleted, would any decision
  about this subject change? *If not, it is ritual — drop it or move it to Known
  gaps. Requirements imported to "be thorough" inflate the model and dilute the
  ratings that drive choices.*

#### Let risk decide where requirements go deep

- **Do** spend requirement detail where risk exposure (likelihood × impact) is
  highest — the failure modes named in body Risks. *A high-risk concern deserves a
  sharply bounded requirement with measured criteria; a low-risk one can stay
  coarse or live in Known gaps.*
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
- **Do** make sure every assessment answers "which Factor's question does this
  finding help rate?" *An assessment that measures something only because it is
  easy to measure, tied to no Factor's concern, is noise that dilutes the report.*
- **Avoid** an assessment so narrowly metric-shaped that satisfying the letter
  abandons the intent ("coverage ≥ 80%" met by trivial tests). *Where a single
  number invites gaming, pair it with an inspection of intent or phrase the
  statement as the outcome you actually want.*

#### Reference an external assessment; don't copy it

- **Do** point at the spec, doc, or checklist that defines the assessment, naming
  it once.
- **Avoid** extracting, summarizing, or duplicating that content into the
  requirement. *Duplicated criteria drift out of sync with their origin.*
- **Do** reference that entity by the same canonical path used as its own
  target's `source`. *That shared path is the edge between the two targets; it is
  what makes the dependency traceable.*
- **Do** point at the specific applicable part of the referenced entity, and pin a
  version where it matters. *An unversioned, whole-document reference leaves the
  verification scope unbounded — name the section or rule that defines the
  assessment.*

#### Connect to factors deliberately

- **Do** rely on **placement** for the primary factor (nest the requirement under
  it), and add `factors` only to pull the result into additional (**secondary**)
  factor roll-ups.
- **Do** declare `factors` explicitly for a requirement placed directly under a
  target — it's required there, and `null`/`[]`/empty entries don't satisfy it.
- **Do** name in `factors` only factors that are **in scope** — declared on this
  target or on an ancestor. *A named factor that lives on a sibling or a
  descendant doesn't resolve; you can reach up to an ancestor's factor, but not
  across or down.*

#### Split by assessable claim, not by factor

- **Do** size a requirement to one claim you want to rate as a single judgment —
  not to one factor or one section of an assessment. *The requirement is the unit
  assessed and rated, once.*
- **Do** connect a claim that reads through several lenses to multiple factors
  (placement for the primary, `factors` for the rest) instead of copying it into
  a per-factor requirement. *One assessment, counted once, feeding several factor
  roll-ups.*
- **Avoid** a set of requirements that reference the same entity with the same
  assessment, sliced one per factor. *That fragments a single judgment and
  re-assesses the entity repeatedly.*
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
  legitimately diverge. *A subject can match its specification yet fail the user's
  real need, or serve the need while departing from spec; one result cannot carry
  both judgments.*

#### Override criteria only when the shared scale can't express the gradient

- **Consider** a `ratings` override when a requirement has a natural measured
  threshold or a distinct qualitative spectrum (e.g. latency bands).
- **Do** key overrides by existing rating levels and change *only* the
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
  a requirement, marked out of Scope, or recorded in Known gaps.
- **Do** check **consistency**: no two requirements make conflicting claims about
  the same source, and the same term or unit means the same thing across
  statements and criteria.
- **Avoid** **redundancy**: two statements that would always rate together against
  the same source are one requirement. *(Same test as for splitting — could their
  results legitimately diverge?)*
- **Do** move unresolved holes to Known gaps rather than shipping a vague,
  unratable requirement to stand in for them.

## When to update QUALITY.md

A `QUALITY.md` is expected to evolve through two loops: the **improve loop** fixes
the subject against the model, and the **learn loop** reviews the model itself
against reality. The two never fold together — the model's own quality is never
averaged into the subject's rating.

- **Do** revise when a discovery changes the context or content of the
  evaluation — a new factor that matters, a requirement whose assessment changed,
  a scope that shifted.
- **Do** update the model when an evaluation finding shows the model no longer
  reflects the subject's real scope, risks, or decision needs. *That is model
  drift, not merely a weak subject rating.*
- **Do** keep the body current with the frontmatter. *A model whose body no
  longer explains its factors misleads the next evaluator.*
- **Avoid** using `QUALITY.md` as a defect backlog. *Subject defects belong in
  the subject's normal planning system unless they also change what quality
  means or how it should be assessed.*
- **Do** distinguish *recalibration* (a deliberate decision to reset a criterion
  because you have learned what is achievable) from *drift* (the model silently
  falling out of step with the subject). *Recalibration is healthy: after a
  breakthrough, raise `minimum` so the new floor sticks; after hitting a real
  constraint, lower a `target` consciously and say why in the body.*
- **Avoid** sharpening criteria only to keep ratings green. *The review's job is to
  keep the rubric valid, not passing; locked baselines and an honest "not assessed"
  guard against gaming it.*
- **Do** treat a finding that no existing requirement anticipated as a signal to
  add a requirement or Factor. *A real weakness your model could not express is the
  strongest evidence the model is incomplete — the model improves by being used,
  not only by being authored.*
- **Do** periodically check that satisfying the requirement set would actually
  deliver the body's Needs. *If the model can be fully green while the subject
  still fails its purpose, the requirement set is incomplete — that is model drift,
  not a strong subject.*
