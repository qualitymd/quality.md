# Authoring QUALITY.md

The canonical reference and best-practices guide for understanding and working
with QUALITY.md files: the format concepts, authoring practices, and jobs
attached to each.

This guide conforms to [`SPECIFICATION.md`](../resources/SPECIFICATION.md). The
specification governs on any conflict.

## Contents

- [The QUALITY.md file](#the-qualitymd-file)
- [Quality Model](#quality-model)
- [The Markdown body](#the-markdown-body)
- [Rating Scale](#rating-scale)
- [Area](#area)
- [Factor](#factor)
- [Requirement](#requirement)
- [When to update QUALITY.md](#when-to-update-qualitymd)

---

## The QUALITY.md file

A QUALITY.md file is a Markdown file with two parts:

- **YAML frontmatter** — the **quality model**: a structured, declarative
  description of what quality means for the entity being evaluated.
- **Markdown body** — the evaluable judgment context for the model: what the
  root area is, why quality matters, what decisions the model supports, and what
  context is missing or inaccessible.

The whole file represents a single root **area** — the top entity whose
quality is modeled. Everything else nests beneath it.

The file's location carries meaning: a `QUALITY.md` in a directory makes that
directory and everything under it (`**/*`) the default scope of evaluation,
unless an area narrows it with a `source`.

### Working with the file

- **Do** put `QUALITY.md` in the directory whose contents are the root area. The
  default scope is that directory and all subdirectories — no `source` needed in
  the common case.
- **Do** keep the root `title`, body Overview and Scope, file location, and
  root `source` aligned on the same evaluated root area. *If those disagree, the
  model may be valid YAML but misleading about what is actually being judged.*

---

## Quality Model

The **quality model** is the frontmatter: the root node of the file. It carries
one model-wide property, the `ratingScale`, plus all the properties of an
[Area](#area) — because the model root *is* the apex area and tops the
area tree.

### Properties

| Property       | Presence   | What it is                                                               |
| -------------- | ---------- | ------------------------------------------------------------------------ |
| `title`        | Required   | Human-readable name of the entity whose quality is modeled.              |
| `description`  | Optional   | A concise statement of what the model's area is.                         |
| `ratingScale`  | Required   | The [rating scale](#rating-scale) — the levels every result is rated on. |
| `factors`      | Optional\* | The [factors](#factor) — lenses through which quality is described.      |
| `requirements` | Optional\* | The [requirements](#requirement) assessed against the root source.       |
| `areas`        | Optional\* | Child [areas](#area) — more focused models nested to any depth.          |
| `source`       | Optional   | Scope override; omit at the root to take the directory default.          |

\* At least one of `factors`, `requirements`, or `areas` must be present.

`ratingScale` is the only property unique to the model; everything else it
shares with Area.

Area names, Factor names, and Rating Level IDs use the same strict name grammar:
letters or digits at both ends, with letters, digits, `_`, or `-` inside.
Requirement statements stay natural language and are not constrained by that
grammar. When a tool needs a stable text handle, it uses canonical model
references such as `area:root`, `area:api`, `factor:api::reliability`, or
`rating:target`.

### Working with the model

- **Do** name the entity (`title`) and think through the body's Overview first —
  what the thing is and what "good" means for it — before listing factors.

#### Build from context, then scale, then model tree

- **Do** fill the Markdown body before expanding the frontmatter. *The body is
  where the model gets its judgment context: what the root area is, what decisions
  the model supports, what quality means here, and which risks matter enough to
  assess.*
- **Do** confirm the rating scale after the body and before writing requirements.
  *The scale is the shared vocabulary for turning future findings into ratings;
  requirements are easier to write once "unacceptable", "minimum", "target", and
  "outstanding" mean something for this root area.*
- **Do** derive factors and requirements from the body context. *Those factors
  and requirements should express the root area's needs, risks, scope, and unknowns,
  not lead them.*
- **Do** trace at least one important concern from body to model before expanding
  the tree. *A useful trace reads like: a need names the outcome, a risk names
  the failure mode, a factor names the quality lens, and a requirement names the
  inspectable expectation.*

#### Keep the root lean when child areas carry the detail

- **Consider** declaring only model-wide factors at the root and pushing
  narrower factors/requirements down to child areas. *Model-wide means the factors
  that recur across the root's constituents (often stewardship lenses like
  currentness or traceability), not an arbitrary subset; a flat root with
  everything at one level is harder to read and maintain.*
- **Avoid** modeling every property at the root "to be safe." *An entry on
  factors, requirements, **or** areas is enough.*
- **Do** make sure the model reaches requirements somewhere in the area tree
  before treating it as evaluable. *Factors can describe concerns, but without
  requirements they do not give an evaluator anything to rate.*

#### Make the traceability graph visible

- **Do** treat the model as a graph: a tree of areas plus the assessment
  references that link a requirement to the entity supplying its criteria. Which
  entity is the criteria for which — and how quality depends along those edges —
  is often the most valuable thing the model records.
- **Do** make each assessment either inline or a reference to another entity, and
  reference that entity by the same selector used as its area's `source`,
  so a reader can follow the edge from one area to the next.
- **Do** stop the chain where ownership or value stops: model referenced entities
  as areas while you own them; let a guide that governs its own kind be its own
  assessment rather than adding another layer.

---

## The Markdown body

The body is evaluable judgment context: what the root area is, why its quality
matters, what decisions the model supports, which needs and risks shaped it, and
what context is missing or inaccessible. It should provide enough concise,
self-explanatory context for a later human or agent to justify the model,
evaluate the model's quality, and decide whether the model still fits the
root area.

A strong body makes its completeness, thoroughness, recency, grounding,
agent-accessibility, and open questions visible instead of implicit.

**Agent-accessible** support is available to the evaluating agent through the
repository, cited local paths, configured tools, linked public sources, or
explicitly provided context. If important support exists but is private,
permission-limited, stale, only known from memory, or unavailable to the agent,
record that limitation in the relevant section's unknowns or open questions.

The body is optional and fixes no required sections; you may rename, reorder, or
replace these. They're recommended starting points:

| Section      | Desired outcome                                                                                                                                                         |
| ------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Overview** | A reader can say what the root area is, who depends on it, and why its quality matters. This names the real entity, not just the repo or file where `QUALITY.md` lives. |
| **Scope**    | A reader can tell what is included, what is excluded for now, and where the model boundary sits. Out-of-scope is not a deficiency.                                      |
| **Needs**    | A reader can see the outcomes the root area must support and the users, operators, maintainers, or downstream systems those outcomes serve.                             |
| **Risks**    | A reader can see the failures that would make the root area untrustworthy, unusable, unsafe, expensive, or hard to change. These are raw material for initial factors.  |

### Shape of a body section

Write each section — including ones you add — to a common shape, so the body
reads consistently as it grows:

1. **Purpose** — open with one line on why this section matters for *this*
   root area, not in the abstract. *If the line would read the same for any
   project, it isn't earning its place.*
2. **Contents** — concise, self-explanatory judgment context for this section.
   State the section's conclusion clearly enough that it can be reviewed on its
   own; cite supporting detail instead of copying it; include enough specificity
   to evaluate completeness, thoroughness, recency, and grounding in
   agent-accessible support.
3. **Unknowns & open questions** — captured for every section, scoped to what
   that section covers. An **unknown** is a broad area of uncertainty within the
   section's topic that may not resolve to a single answer; an **open question**
   is sharper — a specific question about that section with a particular answer,
   still unresolved. Both are context that feeds the model, not commentary on it.
   *Write "none known" when there are none, so the absence reads as considered,
   not skipped.*
4. **State** — close with the review-provenance line.

### Example body section

This example shows the section shape in use. It is illustrative; adapt the
domain, cited support, unknowns, and open questions to the actual root area.

```markdown
## Needs

Daily support triage quality matters because support leads use this model to
decide whether the inbox is safe to hand off between shifts.

Support leads need urgent customer-impacting messages surfaced before routine
account questions. Agents need enough current policy context to answer without
guessing. Maintainers need triage rules that are inspectable in
`support/policies/triage.md` and reflected in saved views under `support/views/`.

*Unknowns* — holiday launch escalation load is based on last year's notes, which
are not agent-accessible.
*Open questions* — what response-time target should apply to enterprise-contract
escalations?

*Reviewed — Ada Lovelace, 2026-05; agent-reviewed — Codex (GPT-5.5), 2026-06.*
```

### Mark the state of a section

Close each section with its unknowns, its open questions, and a state line. Both
are scoped to what the section covers — context that feeds the model, not
commentary on it. An **unknown** is a broad area of uncertainty within the
section's topic that may not resolve to a single answer; an **open question** is
sharper — a specific question about it that has a particular answer, still
unresolved (see [Keep Scope, unknowns, open questions, and "not assessed"
distinct](#keep-scope-unknowns-open-questions-and-not-assessed-distinct)).

Because the body is largely agent-authored, the freshness signal worth trusting
is not when a section last changed but when a person last stood behind it. The
state line carries two reviews — the last human review (cite the person) and the
last agent review (name the agent surface and model used):

```markdown
## Risks

A regional outage is the failure that would most erode trust: orders silently
drop instead of failing over. Cost overrun is a distant second.

*Unknowns* — failover under a full regional outage is untested.
*Open questions* — should orders fail over to another region, or degrade in place?

*Reviewed — Margaret Hamilton, 2026-05; agent-reviewed — Codex (GPT-5.5), 2026-06.*
```

A section with nothing outstanding still says so:

```markdown
*Unknowns* — none known.
*Open questions* — none.

*Reviewed — Margaret Hamilton, 2026-05; agent-reviewed — Codex (GPT-5.5), 2026-06.*
```

- **Do** capture `Unknowns` and `Open questions` for every section, writing "none
  known" rather than omitting them. *On a high-leverage file, an explicit "none"
  reads as considered; a blank reads as skipped.*
- **Do** cite a named person in `reviewed`. *An anonymous review carries no
  accountability; the name is what makes it a trust signal.*
- **Do** cite the agent surface and model in `agent-reviewed`, for example
  `Codex (GPT-5.5)`. *The agent name alone is too ambiguous once model behavior
  changes across versions.*
- **Do** advance `reviewed` only when that person actually read and endorsed the
  section — never for an agent or mechanical edit.
- **Do** read `agent-reviewed` newer than `reviewed` as the warning state: the
  section has agent changes not yet human-endorsed.
- **Do** treat a missing `reviewed` as **unreviewed** — agent-touched, not yet
  vetted. *Absence is honest; never backfill a name and date a person didn't
  earn.*

### Working with the body

- **Do** write the body before expanding the model tree. *It is the fastest way
  to discover what factors and requirements the frontmatter should express.*
- **Do** write the body so it can be evaluated for quality in its own right.
  *A later reviewer should be able to judge whether the context is complete
  enough, current enough, specific enough, grounded enough, and accessible enough
  to support the model.*
- **Do** cite supporting detail when it materially grounds a section, and flag
  important support that is not agent-accessible. *The body should not become an
  evidence dump, but a later evaluator must be able to tell what the judgment
  rests on and what context could not be inspected.*
- **Do** capture, in Needs and Risks, why some requirements matter more than
  others. *Importance and gaps both depend on this context.*
- **Do** use the body to explain any rating-scale change. *A custom scale should
  answer a real decision need visible in the body.*
- **Do** treat Needs as *benefits to realize*, not only outcomes to protect, and
  Risks as the problems that erode them. *A model that lists only failure modes
  can rate every requirement at target while the root area is still not worth
  relying on — because nothing weighs whether the benefits, on the whole,
  outweigh the residual problems.*
- **Consider** separating two reasons a concern matters: how much its failure
  hurts (stakeholder, safety, or business stakes) and how far its failure spreads
  (how much else it forces to change). *A concern can be low-stakes yet
  high-blast-radius, or the reverse; saying which helps the next evaluator weigh
  roll-ups.*
- **Consider** noting where two factors or requirements pull against each other
  (tighter access control vs. faster onboarding, latency vs. cost) and which way
  you have chosen to lean. *A model that hides its trade-offs invites an evaluator
  to "fix" a deliberate compromise.*

#### Say which sense of "good" this model uses

Quality is not one thing. A root area can be judged by **conformance** (does it
match its specification?), by **fitness for purpose** (does it serve the user's
real need?), or by **value** (is it worth its cost?). These can disagree — a
root area can meet its spec yet fail the need, or serve the need while departing
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

#### Keep Scope, unknowns, open questions, and "not assessed" distinct

- **Do** use **Scope** for concerns outside the model's remit, and a section's
  **Unknowns** for in-scope concerns you've deliberately deferred or cannot yet
  define.
- **Do** record known unknowns under the **Unknowns** of the section they bear
  on, when missing context or weak evidence prevented that part of the body from
  being fully identified. *This keeps uncertainty visible without pretending it
  has already been evaluated.*
- **Consider** keeping an **open question** distinct from an unknown: an unknown
  is a broad area of uncertainty within the section's topic that may not resolve
  to a single answer; an open question is a specific question about that section
  with one particular answer, still unresolved. Both are input to the model, not
  statements about how the model is built. *Note each in the section it bears on
  — an open question is a standing prompt for the next review.*
- **Don't** confuse a declared unknown (your standing declaration) with a **not
  assessed** result (an evaluator's per-run finding that evidence was missing).
- **Do** record low *confidence* in an assessment, not only its absence. *A
  requirement rated `target` on one stale benchmark or a single reviewer differs
  from one rated on sustained evidence; note that fragility alongside the concern
  it qualifies. "Rated but barely trusted" is neither "not assessed" nor "no
  gap," and it is often where the next evaluation should look first.*

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
  (see [Override criteria only when the shared scale can't
  express the gradient](#override-criteria-only-when-the-shared-scale-cant-express-the-gradient));
  the description never changes.

### Properties (per level)

| Property      | Presence    | What it is                                                      |
| ------------- | ----------- | --------------------------------------------------------------- |
| `level`       | Required    | Rating Level ID; unique within the scale.                       |
| `title`       | Required    | Human-readable label for reports.                               |
| `description` | Recommended | What the level means across the model (fixed).                  |
| `criterion`   | Required    | Default rule for rating a requirement's findings at this level. |

### Working with the rating scale

#### Reuse the suggested four-level scale unless you have a reason not to

- **Consider** the **outstanding > target > minimum > unacceptable** scale as a
  default — a stretch level, the level to aim for, the floor you've agreed to
  rely on, and below the floor. *A shared four-band vocabulary is enough for
  most models and keeps reports comparable.*
- **Do** use the default display titles `🟢 Outstanding`, `🔵 Target`,
  `🟡 Minimum`, and `🔴 Unacceptable` when the standard scale fits. *The emoji is
  only a scanning aid for human reports and frontmatter; the plain `level` IDs
  still carry identity, ordering, and references. Use plain or custom titles when
  a project style demands it, but avoid emoji-only labels.*
- **Do** choose a different scale only when the root area demands it (e.g. a
  pass/fail gate wants two levels). *The scale should fit how decisions are
  actually made about this entity.*
- **Do** review the scale after writing the body. *The body should reveal whether
  "good enough" and "excellent" need sharper meaning for this root area.*
- **Do** test adjacent levels against the same plausible finding. *If one
  finding could reasonably satisfy both `target` and `minimum`, sharpen the
  criteria until the boundary is usable.*
- **Do** let the model's job pick the number of levels: a model that gates against
  an agreed bar can use two levels (pass / fail); a model that judges *how good*,
  compares options, or surfaces strengths and weaknesses earns a graded scale. *A
  pass/fail model deliberately gives up the ability to say "how good" — sometimes
  exactly what you want, sometimes a loss.*
- **Do** fix the *required margin* in the body: how far above `minimum` this
  root area must actually land to be good enough, and why (safety, reputation,
  regulatory exposure, temporary use). *The same four levels can serve very
  different domains and stakes; what differs is the band the root area must
  reach. State it so an evaluator knows whether an all-`minimum` result is fine
  or alarming.*
- **Do** calibrate the levels against concrete exemplars where you can — a real
  artifact you would call `outstanding` and one you would call `unacceptable`.
  *Thresholds set against known cases are far more defensible than ones invented
  in the abstract; if no real artifact would ever earn `outstanding`, the level is
  decorative.*
- **Avoid** inventing a custom scale before the body reveals a real need for
  one.

#### Keep `description` about meaning and `criterion` about rating

- **Do** write `description` as what the level *is* ("the floor you've agreed to
  rely on") and `criterion` as the test a result must pass ("falls short of
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

## Area

An **area** is an entity, or set of entities, with quality requirements subject
to evaluation. It is the recursive node of the model: the root model is the root
area, and every entry under `areas` is another area, nested to any depth.

An **entity** is anything evaluated for quality. Most often it is a concrete
artifact — a document, a source tree, a config file — but it can also be a
service, a behavior, or a process. The same entity can be an area in one
relationship and the reference supplying another requirement's criteria in
another.

An area's **source** defines *what* it evaluates; its **factors** and
**requirements** define *what is assessed* about it. An area may also be a pure
**grouping area** — only child areas, no requirements of its own.

### Properties

| Property       | Presence   | What it is                                                             |
| -------------- | ---------- | ---------------------------------------------------------------------- |
| `title`        | Required   | Display name in reports; the map key is the Area name.                 |
| `description`  | Optional   | What the area is — the entity or scope it covers.                      |
| `factors`      | Optional\* | [Factors](#factor) scoped to this area's subtree.                      |
| `requirements` | Optional\* | [Requirements](#requirement) assessed against this area's source.      |
| `areas`        | Optional\* | Child areas, nested to any depth.                                      |
| `source`       | Optional   | Where the entities live; inherits the nearest ancestor's when omitted. |

\* At least one of `factors`, `requirements`, or `areas` must be present, same
as the model. An area with only child `areas` is a grouping node; even then,
each area should lead to at least one requirement somewhere in its subtree — a
subtree with no requirements evaluates nothing.

### Working with areas

An Area's stable ID is its path of Area names from the root. The root Area ID is
empty and renders as `area:root` when a canonical model reference is needed.

- **Do** point an area at the thing whose quality you actually care about.
  **Avoid** modeling fixtures, generated output, or build scaffolding as areas.
  *Those usually support evaluation; they are not the root area being evaluated.*

#### Choose areas that are authored, inspectable artifacts

- **Do** prefer areas that are authored and inspectable — source trees, specs,
  tests, docs, configs, checklists — over runtime behavior or process you cannot
  open and examine. *The model is judged by inspecting the entity; point it at
  something an evaluator can actually read.*
- **Consider** whether the entity is **constitutive** (it *is* the product — the
  source, the document) or **normative** (it *governs* the product — a spec, style
  guide, or checklist others are judged against). *Normative artifacts are
  high-leverage: quality invested once propagates to everything they govern, often
  the fastest early win for a low-maturity root area.*
- **Do** match the area's grain to the concern: a single test and a test *suite*
  are different areas. *A collection is judged by whole-set concerns no member
  has — coverage, balance, non-redundancy, navigability. Point `source` at the
  whole set and write requirements at that grain.*

#### Choose the decomposition shape: primary-subject, collection, or composite

Every area — the root included — decomposes in one of three shapes. The shapes
are a question you ask *at each node*, not a one-time root classification, and
they compose: a node of one shape can hold children of another, to any depth.

- **Primary-subject** — one entity with one factor family; children, if any, only
  refine it. The factor-coverage aim (see [Cover the domain's stable stakes](#cover-the-domains-stable-stakes-before-specializing)) applies at this
  kind of node.
- **Collection** — many entities of the *same* kind (a test suite, several like
  services). The node carries whole-set concerns no member has — coverage,
  balance, non-redundancy, navigability — and each member carries its own family.
- **Composite** — many entities of *different* kinds, each with its own,
  largely-disjoint factor family. The node's own concern is the **coherence
  between** its parts: whether they stay aligned, conformant, and current with
  each other. That cross-part coherence is exactly the assessment-edge graph from
  [Make the traceability graph visible](#make-the-traceability-graph-visible).

- **Do** treat a near-disjoint factor family as a first-class signal to split a
  part into its own area. *When part of the entity is best judged through
  characteristics that would not apply to its siblings, it has earned an area —
  the sharpest form of the split test below.*
- **Do** ask the shape question again at every child, and let the shapes nest. *A
  composite node commonly holds a collection child, and a collection member or a
  constituent can itself be composite.*
- **Avoid** flattening a composite into one factored node. *Holding every part's
  factors at one level either collapses the model to a single part or jams
  incompatible factor families into one list.*

A non-trivial entity is usually **composite** at the root: a running artifact,
the requirements that define it, its docs, and other parts each carry their own
factor family, so no single domain factor list belongs at the root. The shape
recurses — here a composite root holds primary-subject constituents alongside a
collection child (illustrative; adapt the parts to the actual entity, which need
not be software):

```text
root  (composite)
├── harness        (primary-subject constituent)
├── quality-md     (primary-subject constituent; learn loop, kept out of roll-up)
└── apps           (collection)
    ├── apps/product-a   (primary-subject — or itself composite)
    └── apps/product-b   (primary-subject — or itself composite)
```

#### Split off a child area only when it has distinct factors or requirements

- **Consider** a child area when a part of the root area is best evaluated
  through its own factors or requirements that would not apply to its siblings.
  *For example, in a software product model, a payment integration might have
  its own security concerns; in a support-operations model, urgent triage might
  have distinct responsiveness concerns.*
- **Avoid** child areas that merely re-state the parent's requirements. *Child
  areas don't inherit the parent's requirements, but duplicating them by hand
  is a maintenance trap — let the ancestor's source cover the shared scope.*
- **Do** place a requirement on the most specific area whose source actually
  owns the concern, not on a broad ancestor that merely contains it. *Assessing a
  payment-integration concern at the root re-judges it at every level and blurs
  who owns the finding; a second evaluator should be able to place the same
  finding at the same area.*
- **Avoid** splitting off a child whose rating is meaningless on its own. *If it is
  only interpretable alongside its siblings, keep the concern in the parent; split
  when the child could be evaluated, and trusted, in isolation.*

- **Do** leave a grouping area's `source` implicit and let its children narrow
  it. *A grouping area has no local rating; its aggregate reflects only its
  children.*

#### Define `source` only to narrow or relocate

- **Do** omit `source` on the root to take the directory default, and on a child
  to inherit the nearest ancestor's scope.
- **Do** set `source` (a path or glob, relative to the file) when an area
  evaluates a specific subtree. *An ancestor's source may overlap a child's; even
  so, each requirement evaluates against its own area's source — the overlap
  just means the child's entities also fall within the ancestor's scope.*
- **Do** compare each `source` to the area's title and description. *The source
  should select the entity being modeled, not a convenient proxy, generated
  output, or a broader repository area that happens to contain it.*

#### An entity can be both an area and an assessment reference

- **Do** decide by *role in this relationship*, not by identity: an entity is an
  area where its quality is the thing being managed, and an assessment
  reference where it supplies the criteria for judging another area. A spec is
  both — an area in its own right and the criteria for judging the code that
  implements it.
- **Do** point `source` at the whole entity, not a thin entry point or proxy.
- **Avoid** treating "area or reference?" as exclusive. *An owned entity used
  as criteria usually deserves an area of its own as well.*
- **Do** spend early effort on the normative artifacts that play both roles — a
  spec, standard, or style guide. *Improving one raises every area judged against
  it, so it is often the highest-leverage place to invest quality early.*

#### Ground high-leverage concerns in normative artifacts

Some concerns are *high-leverage* — improving them propagates to everything they
touch — and the strongest of these are usually governed by a **normative
artifact** an evaluator can point at. Treat both the concern and its anchor as
modeling obligations, not optional polish.

- **Do** model the concerns that are high-leverage and germane to the entity's
  quality domain as areas, even when no governing artifact exists yet. *A germane
  high-leverage concern left unmodeled is a coverage gap, not neutrality; carry
  the area and record the missing anchor as a finding within it. For software
  product quality these illustratively include requirements/intent definition,
  data quality, and interface contracts; other domains have their own — a
  terminology standard for a document, a schema or collection methodology for a
  dataset, operating thresholds for a service. Treat the list as a prompt, never a
  quota.*
- **Do** anchor a high-leverage area in a normative artifact — a spec, standard,
  requirements doc, contract, or style guide — and assess the area against it
  rather than embedding all the criteria inline. *An external anchor lets two
  evaluators reach the same finding and gives the requirement something to trace
  to; quality invested in the anchor propagates to everything it governs.*
- **Consider** the absence of a normative anchor a recorded finding, not silence.
  *Where the domain implies a governing artifact and none exists — or an area is
  assessed only by criteria embedded in its `assessment` — name that as a gap.
  Assess that the expectation is **governed and current**, referenced by review,
  automation, or the model itself; an artifact that exists but is stale or
  load-bearing on no one is not an anchor.*
- **Do** locate this concern in **maintainability** (an ungoverned high-leverage
  concern is cognitive and intent debt) and in the model's **evaluability and
  traceability** (with no anchor, criteria are forced inline and trace only to the
  author's judgment). **Avoid** linking it to the governed domain factor itself.
  *Having a data-quality standard is not the same as having good data; binding the
  two conflates governance with the governed property and double-counts.*
- **Do** keep the requirement **standing**, instantiated per project, and
  re-evaluated. *The rating — not the requirement's presence — carries whether the
  artifacts exist today; a requirement that appears because a doc is missing and
  vanishes once it is written is the defect-backlog anti-pattern. Expect new
  high-leverage concerns to imply new expected areas as the model matures.*
- **Do** calibrate leverage to *this* entity and trace it to a recorded Need or
  Risk. **Avoid** importing a universal roster of areas every model must carry.
  *"High-leverage and germane here" is the inclusion test; an unbounded checklist
  of expected areas collapses the model into generic best practice — the same trap
  as importing a standard factor list.*

This is the model's own meta-principle applied one level down: the self-check area
already anchors the model in this guide because normative artifacts are
high-leverage. The same reasoning applies to the concerns the evaluated entity is
made of.

#### Cover the domain's constituent kinds

Deciding the root is composite (above) raises the next question: *which*
constituents. Enumerating them by walking the repository's folders finds only the
constituents that already have a home; the thin, scattered, or missing ones —
often the highest-leverage early findings — stay invisible. Enumerate by
**constituent kind** instead, inferred from the entity's quality domain, then
account for each kind. Two questions generate the kinds.

**1. Which stewardship concern leaves an artifact here?** Caring for any entity
carries a recurring set of concerns, and each tends to leave an authored,
inspectable artifact that is a candidate constituent. They fall in two bands:

- *Lifecycle* (roughly sequential): **discover** the problem and solution space →
  **define** what it should be → **realize** it → **verify** it → **enable** its
  audiences to use it → **operate** it → **maintain** it.
- *Protective* (cross-cutting and bidirectional): **secure** — guard the entity
  from harm by the world (breach, tampering, injection); and **safeguard** —
  guard stakeholders and the environment, internal and external, from harm by the
  entity (a destructive action, an unsafe output, privacy harm to a data
  subject). The two are orthogonal: an entity can be hardened against attackers
  yet routinely harm its own users.

**2. Whom does each artifact serve, and for what job?** A concern splits into
several constituents when they serve different audiences or purposes — by
audience (learner, practitioner, operator, maintainer, downstream, auditor, the
public) and by job (acquire vs. apply, act vs. understand).
[Diátaxis](https://diataxis.fr/) is this lens applied to the *enable* concern:
tutorial, how-to, reference, and explanation are four constituents with different
factor families, not one "documentation" area. The same audience×purpose split
recurs elsewhere — a maintainer-facing CI gate and a user-facing acceptance suite
are two *verify* constituents.

A concern projects into the model in up to three ways — as a **factor** (the
quality lens), a **constituent** (the artifact that pursues it), and an
**audience** (who it serves or protects). They share a name because they share a
root concern, not because they duplicate one another.

- **Do** name the projection you are modeling and model it once. *`secure` is a
  security factor on each area, a threat-model constituent, and an auditor
  audience; the security of the server is a factor on the server area, while the
  security policy is its own area. Modeling the same projection twice
  double-counts.*
- **Do** enumerate the constituent kinds the entity's domain implies, then
  account for each: model it, defer it in Scope, mark it out of Scope, or record
  it as an unknown. *Silence is a coverage gap, not neutrality — the same failure
  as a sparse factor list (see [Cover the domain's stable stakes](#cover-the-domains-stable-stakes-before-specializing)).*
- **Do** derive the audience side from the body's Needs. *The Needs already name
  the stakeholders; each audience a Need names should have an enabling — and
  verifying — constituent that is modeled or consciously accounted for.*
- **Do** carry a germane, high-leverage kind as an area even when its artifact is
  thin or missing, recording the gap as a finding within it. *A missing test
  suite or absent threat model is a high-value early signal; an empty area with a
  missing-anchor finding surfaces it, where folder-walking drops it. This is the
  [normative-artifact rule](#ground-high-leverage-concerns-in-normative-artifacts)
  applied to whole constituents.*
- **Avoid** importing the list as a roster every model must carry. *The kinds are
  a prompt, not a quota; a kind earns an area only when it leaves an owned,
  inspectable artifact, its factor family diverges from its siblings, and it
  traces to a Need or Risk. A throwaway script earns almost none.*
- **Consider** the audience×purpose split as the way to find a constituent's
  internal shape: a constituent that fans out by audience or job is itself
  composite or a collection (see [Choose the decomposition shape](#choose-the-decomposition-shape-primary-subject-collection-or-composite)).

The kinds and their conventional factor families are inferable once the domain is
named (illustrative — software product quality as one domain; adapt the column to
the actual entity, which need not be software):

| Stewardship concern | Software-product instance                                 | Conventional factors                         |
| ------------------- | --------------------------------------------------------- | -------------------------------------------- |
| discover            | research notes, design explorations, decision records     | traceability, credibility                    |
| define              | requirements, specs, interface contracts, schema          | completeness, consistency, traceability      |
| realize             | source, services                                          | reliability, performance, maintainability    |
| verify              | tests, validation suites                                  | coverage, determinism, maintainability       |
| enable              | tutorial / how-to / reference / explanation (Diátaxis)    | currentness, completeness, understandability |
| operate             | runbooks, deploy config, monitoring-as-code, SLOs         | reliability, operability, observability      |
| maintain            | migrations, upgrade & deprecation guides, changelog       | maintainability, modifiability, currentness  |
| secure              | threat model, security policy, secrets/auth config        | security                                     |
| safeguard           | safety case, privacy/impact assessment, output guardrails | safety, privacy, compliance                  |

Treat the table as a prompt for *this* entity's domain, not a checklist: many
cells collapse into one area, several merge, and some are legitimately absent.

#### Carry the recurring use-context constituents

QUALITY.md is domain-agnostic in *what* it models, but its assumed context of use
is an agent/AI-assistant-collaborated project. That context of use — not the
modeled domain — makes two constituents recur in a composite root regardless of
what the root entity is:

- the **agent harness** — the instructions that steer the agent working the
  project (its agent guidance files, skills, and prompts); and
- the **QUALITY.md self-check** — the model's own quality.

Distinguish these from **domain constituents**, which vary with what is modeled
(a data set's schema and collection methodology; a document's terminology
standard and sources) — enumerate those with [Cover the domain's constituent kinds](#cover-the-domains-constituent-kinds). Domain constituents change from
project to project; these two use-context constituents recur across them.

- **Consider** the harness and the self-check expected constituents of a composite
  root, justified by the context of use — but **earn** them, do not assume them.
  *The inclusion test is unchanged: high-leverage, germane here, an owned and
  inspectable artifact, traced to a Need or Risk. A harness-less or throwaway
  project carries neither; this is not a roster every model must hold.*
- **Do** keep the QUALITY.md self-check on the **learn loop**, out of the entity's
  roll-up. *The model's own quality is never averaged into the root area's rating
  (see [When to update QUALITY.md](#when-to-update-qualitymd)); model it as a
  constituent whose quality is tracked and reported on its own axis.*
- **Do** treat the agent harness as partly **normative** — it governs agent
  behavior, so it plays the dual area/assessment-reference role (see [An entity
  can be both an area and an assessment
  reference](#an-entity-can-be-both-an-area-and-an-assessment-reference)). *Watch
  for double-counting if its influence is also assessed inside a domain
  constituent.*

#### Write a description that distinguishes, not enumerates

- **Do** state what the area *is* and how it differs from its siblings/parent.
  **Avoid** restating its factors or requirements in the description. *The
  description identifies the evaluated entity; factors and requirements define
  what is assessed about it.*

#### Decide how ratings roll up

*Read this subsection when you are reasoning about roll-up or evaluating; a first
model can defer it.*

An area's local rating and a factor's rating are inferred from the requirement
results beneath them. The format fixes no aggregation formula, so how those
results combine is a modeling decision you should make and communicate.

- **Do** decide and state how an area's requirements combine into its rating when
  it is not obvious. *Two defensible defaults: **worst-of** (the weakest finding
  sets the rating — right when any unacceptable requirement makes the whole
  untrustworthy) and **most-common / median** (right when requirements genuinely
  compensate for one another).*
- **Avoid** implying that rating levels *average*. *The scale is ordinal
  (outstanding > target > minimum > unacceptable); the arithmetic mean of ordinal
  levels has no meaning — three `target`s and one `unacceptable` is not "slightly
  above minimum."*
- **Do** identify the requirements that can **veto** a rating — a single
  `unacceptable` finding that makes the root area not good enough no matter how
  strong everything else is (secrets stored in plaintext, data loss on failover).
  *Most requirements trade off against each other; a few cap the rating. Sharpen a
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

| Property       | Presence    | What it is                                                     |
| -------------- | ----------- | -------------------------------------------------------------- |
| `title`        | Required    | Human-readable label for reports and status output.            |
| `description`  | Recommended | The characteristic, defined operationally for this entity.     |
| `factors`      | Optional    | Sub-factors — recursively a factor.                            |
| `requirements` | Optional    | [Requirements](#requirement) uniquely relevant to this factor. |

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
  it here does not earn a factor, however standard it is elsewhere; pull from a
  catalog as a prompt, never as a quota.*
- **Do** prefer general-purpose, conventional factor names for the quality
  domain once a concern earns a factor. *For software product quality, examples
  include `reliability`, `security`, `usability`, `maintainability`,
  `performance`, `compatibility`, and `portability`; other domains have their
  own conventional factor families. These examples are illustrative,
  non-exhaustive, and sometimes overlapping. Requirements and assessments are
  where the model maps those lenses to the root area's unique quality
  expectations.*
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
  product, document, data set, model, service operation, and human process each
  has a different conventional factor family; these domains are illustrative,
  non-exhaustive, and may overlap in a real model.*
- **Do** include the domain's common stable-stakes factors for the root area, or
  explicitly justify why each omitted one is out of scope, delegated to a child
  area, or still unresolved as an unknown. *A sparse root model should be a
  conscious decision, not the result of only modeling the first risks that came
  to mind.*
- **Do** treat roughly ten factors as a reasonable aim for a **primary-subject
  node** (see [Choose the decomposition shape](#choose-the-decomposition-shape-primary-subject-collection-or-composite)).
  *Fewer than eight should trigger a coverage review; four to six is usually too
  thin unless that node is deliberately narrow, temporary, or mostly delegated to
  child areas.* At a **composite** root the aim applies **per constituent**, not at
  the root: each primary-subject constituent earns its own ~ten-factor family,
  while the composite root itself carries only the factors that recur across
  constituents — typically stewardship lenses (currentness, traceability,
  consistency, maintainability), each refined per child.
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
- **Do** draw from established attribute families as prompts, not quotas. *In
  software and tooling domains, qualities such as maintainability and testability
  can fit; in data, document, and model domains, qualities such as completeness,
  credibility, currentness, traceability, and understandability often fit. These
  examples are illustrative, non-exhaustive, and may overlap.*

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
  governing entity once. *That keeps the criterion edge visible without repeating
  the same assessment under every factor.*
- **Do** split the requirement only when the referenced entity defines claims
  whose results could legitimately diverge. *A `QUALITY.md` can follow its
  authoring guide for assessability while still lacking credibility; those are
  separate requirements. A single conformance judgment that feeds assessability,
  traceability, and maintainability is one requirement with several factor
  references.*

```yaml
requirements:
  "the quality model follows its authoring guide":
    factors:
      - fitness-for-purpose
      - credibility
      - assessability
      - traceability
      - maintainability
    assessment: >
      Assess QUALITY.md against ./skills/quality/guides/authoring.md,
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

## When to update QUALITY.md

A `QUALITY.md` is expected to evolve through two loops: the **improve loop** fixes
the root area against the model, and the **learn loop** reviews the model itself
against reality. The two never fold together — the model's own quality is never
averaged into the root area's rating.

- **Do** revise when a discovery changes the context or content of the
  evaluation — a new factor that matters, a requirement whose assessment changed,
  a scope that shifted.
- **Do** update the model when an evaluation finding shows the model no longer
  reflects the root area's real scope, risks, or decision needs. *That is model
  drift, not merely a weak root area rating.*
- **Do** keep the body current with the frontmatter. *A model whose body no
  longer explains its factors misleads the next evaluator.*
- **Avoid** using `QUALITY.md` as a defect backlog. *Evaluated-source defects belong in
  the root area's normal planning system unless they also change what quality
  means or how it should be assessed.*
- **Do** distinguish *recalibration* (a deliberate decision to reset a criterion
  because you have learned what is achievable) from *drift* (the model silently
  falling out of step with the root area). *Recalibration is healthy: after a
  breakthrough, raise `minimum` so the new floor sticks; after hitting a real
  constraint, lower a `target` consciously and say why in the body.*
- **Avoid** sharpening criteria only to keep ratings green. *The review's job is to
  keep the rubric valid, not passing; locked baselines and an honest "not assessed"
  guard against gaming it.*
- **Do** treat a finding that no existing requirement anticipated as a signal to
  add a requirement or factor. *A real weakness your model could not express is the
  strongest evidence the model is incomplete — the model improves by being used,
  not only by being authored.*
- **Do** periodically check that satisfying the requirement set would actually
  deliver the body's Needs. *If the model can be fully green while the root area
  still fails its purpose, the requirement set is incomplete — that is model drift,
  not a strong root area.*

### Logging a model change

When the learn loop actually changes the model, record it in the **quality log**
— a curated, evidence-linked timeline of meaningful model changes under
`.quality/log/`, one dated entry per change (`YYYY-MM-DD-<slug>.md`). The log
preserves the *why* a commit message scrolls away; its format contract lives in
[`SKILL.md`](../SKILL.md). The judgment is *what* counts as meaningful:

- **Do** log a change that alters what the model *is* or *how it judges*: adding,
  removing, or renaming an Area, Factor, or Requirement; changing the rating
  scale, a criterion, or a relative weight; shifting scope; changing the apex or
  required margin; or applying an evaluation recommendation. *That is the
  judgment layer git cannot reconstruct from a diff.*
- **Do** state whether a criterion move is deliberate *recalibration* or a *drift
  correction*, and cross-link the evaluation run and recommendation behind it when
  the change came from one. *The evidence link is the entry's whole value over
  `git log`.*
- **Do** write **one entry per coherent change** — a confirmed recommendation
  apply, a model-authoring change, or the initial population — not one per field
  touched. *The unit of record is the decision, not the edit.*
- **Avoid** logging Markdown-body wording, typo, or formatting changes, or
  evaluated-source fixes that leave the model unchanged. *Those are not model
  changes; git already records them, and logging them turns a curated timeline
  into noise.*
- **Avoid** treating the log as a second evaluation record or a defect backlog.
  *It references evaluation runs; it never copies them.*
