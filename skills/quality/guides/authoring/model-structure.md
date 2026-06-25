---
type: Runtime Guide
title: Authoring Model Structure
description: Area, source, decomposition, traceability, constituent, and use-context guidance for QUALITY.md models.
tags: [quality, authoring, guide]
---

# Authoring Model Structure

Read this when:

- creating or revising Areas, Sources, child Areas, decomposition shape, traceability edges, normative artifacts, constituent kinds, or recurring use-context constituents.

Depends on:

- `../authoring.md`

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
| `factors`      | Optional\* | [Factors](factors.md) scoped to this area's subtree.                   |
| `requirements` | Optional\* | [Requirements](requirements.md) assessed against this area's source.   |
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

#### Make the traceability graph visible

- **Do** treat the model as a graph: a tree of areas plus the assessment
  references that link a requirement to the entity supplying its criteria. Which
  entity is the criteria for which — and how quality depends along those edges —
  is often the most valuable thing the model records.
- **Do** make each assessment either inline or a reference to another entity, and
  reference that entity by the same selector used as its area's `source`, so a
  reader can follow the edge from one area to the next.
- **Do** stop the chain where ownership or value stops: model referenced entities
  as areas while you own them; let a guide that governs its own kind be its own
  assessment rather than adding another layer.

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
  refine it. The factor-coverage aim (see [Cover the domain's stable stakes](factors.md#cover-the-domains-stable-stakes-before-specializing)) applies at this
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
├── quality-md     (primary-subject constituent)
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
  the area and record the missing anchor as a finding within it.*
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

#### Cover the domain's constituent kinds

Deciding the root is composite (above) raises the next question: *which*
constituents. Enumerating them by walking the repository's folders misses thin,
scattered, or absent constituents. Enumerate by **constituent kind** instead,
inferred from the entity's quality domain and traced to a Need or Risk. Two
questions surface them:

**1. Which stewardship concern is at work here?** To *steward* is to care for
something held in trust, answerable to its stakeholders and its future. That care
carries a recurring set of concerns; each is an ongoing activity of tending that
usually leaves an authored, inspectable artifact as its trace — the candidate
constituent. Read the artifact as evidence of tending, not as the tending itself:
its presence is not the care, and its absence is not the lack of it. They fall in
two bands:

- *Lifecycle* (recurring, not a once-through pipeline): **discover** the problem
  and solution space → **define** what it should be → **realize** it → **verify**
  it → **enable** its audiences to use it → **operate** it → **maintain** it. The
  arrows show a typical first pass; the concerns are standing — revisited for as
  long as the entity is tended, with *operate* and *maintain* never concluding.
- *Protective* (cross-cutting and bidirectional): stewardship under
  **vulnerability** — a party that can be harmed and cannot fully protect itself.
  **secure** guards the entity from harm by the world (breach, tampering,
  injection); **safeguard** guards stakeholders and the environment, internal and
  external, from harm by the entity (a destructive action, an unsafe output,
  privacy harm to a data subject). The two are orthogonal: an entity can be
  hardened against attackers yet routinely harm its own users. This band tracks
  who is exposed, not which phase the work is in, so it cross-cuts the lifecycle
  rather than sitting inside it.

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
- **Do** encode the projection boundary in the model when you carry two or more
  projections of one concern. *Reasoning the boundary is not enough — a later
  reader sees only the emitted model, not your reasoning, and re-litigates whether
  the projections are redundant. On each modeled projection's node, add a YAML
  comment naming its sibling projection and the one-line distinction; and when both
  projections are rated nodes that surface in an evaluation report, add a short
  disambiguating clause to each `description` too, because comments do not survive
  rendering. Keep the clause to the distinction — that is the same "distinguishes,
  not enumerates" rule descriptions already follow. The Agent Harnessability factor
  vs. the agent-harness area is the canonical instance.*
- **Avoid** letting a stewardship word modify or replace the taxonomy noun for
  its projection. *A concern is the *source* a factor projects from, not a kind of
  factor. Write "model-wide factors" (and, if useful, that they trace to
  stewardship concerns) — not "stewardship factors" or "stewardship lenses," which
  demote the term of art to a subcategory of the philosophical word. The singular
  gloss "a factor is a quality lens" is fine; it names what a factor is.*
- **Do** enumerate the constituent kinds the entity's domain implies, then
  **model each as its own area by default**. *Modeling is the default outcome, not
  something a constituent must earn; the thinness of a first pass is never a reason
  to defer or omit one. Build the fullest, most sufficient model the evidence
  supports in a single pass — setup runs roughly once per project.*
- **Do** give a germane constituent its own area unless one of exactly two
  disqualifiers holds — and then fold it or scope it out, never silently drop it:
  - **No distinct concerns** — its quality is already fully judged by its parent's
    or a sibling's factors, so a separate area would assess nothing new. Fold the
    concern into that area.
  - **Not germane / outside the boundary** — the domain does not imply it for this
    entity, or it belongs to another system or owner. Mark it out of Scope.
- **Do** surface every germane concern as a ratable element of the model — an area
  or a requirement that produces a rating — and **never** omit one by recording it
  only in prose. *A Scope or "deferred" note is not coverage. "No artifact yet" is
  not a skip reason: for a germane concern an absent or thin artifact is the
  highest-value early finding, so its absence must land where the evaluation can
  see it. Route it one of two ways, by the same distinct-concerns test:*
  - **(a) a minimal area with a missing-anchor finding** when the kind would carry
    its own factor family once it exists — a high-leverage, first-class kind (a
    test suite, a threat model, specs) or a whole-constituent gap. *The empty area
    rates poorly on its own axis and reserves the shape; folder-walking drops it.
    This is the [normative-artifact rule](#ground-high-leverage-concerns-in-normative-artifacts)
    applied to whole constituents.*
  - **(b) a requirement on an existing area** when the concern folds into that
    area's factors, the gap is partial or a matter of degree, or a standalone area
    would be a single-finding stub. *The gap rates under the parent.*
- **Do** treat in-scope deferral ("modeled later") as a narrow exception, not a
  routine option peer to modeling. *Reserve it for a constituent genuinely blocked
  — e.g. on an undecided boundary — and record the specific blocker. "Next
  iteration" and "the first model is thin" are not deferral reasons, and a blocked
  constituent still surfaces its gap per the rule above.*
- **Do** scale coverage to *this* entity. *Model what is germane here and traces
  to a Need or Risk; a throwaway script carries almost none.*
- **Do** derive the audience side from the body's Needs. *The Needs already name
  the stakeholders; each audience a Need names should have an enabling — and
  verifying — constituent that is modeled or consciously accounted for.*
- **Do** keep the constituent question separate from the stewardship-quality
  question. *Whether a tending leaves an owned, inspectable artifact is the area
  question; whether it is done well — attuned to this entity's situation, not
  merely present or complete — is the factor question. A populated area can still
  rate poorly; an empty one that carries a missing-anchor finding is itself a true
  signal.*
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

Many cells collapse into one area, several merge, and some are legitimately
absent.

#### Carry the recurring use-context constituents

QUALITY.md is domain-agnostic in *what* it models, but its assumed context of use
is an agent/AI-assistant-collaborated project. That context of use — not the
modeled domain — makes two constituents recur in a composite root when germane:

- the **agent harness** — the project-owned harness system around the model,
  scoped in the Area to checked-in governing artifacts such as agent entry points,
  guidance files, skills, prompts, owned hooks, tool/MCP definitions,
  sandbox/permission policy, and orchestration config; read
  [Agent harness authoring](agent-harness.md) for its area guidance; and
- the **QUALITY.md self-check** — the model's own quality.

Distinguish these from **domain constituents**, which vary with what is modeled
(a data set's schema and collection methodology; a document's terminology
standard and sources) — enumerate those with [Cover the domain's constituent kinds](#cover-the-domains-constituent-kinds).

- **Do** model the harness and the self-check by default in a composite root,
  justified by the context of use. *They follow the same rule as any constituent:
  modeled unless a disqualifier fires.*
- **Do** model the QUALITY.md self-check as an ordinary area when it is germane.
  *Use the `quality-md` key, a title of the form `<Root Title> QUALITY.md`, an
  explicit path-based `source` such as `./QUALITY.md`, factors that describe the
  model artifact's qualities, and a requirement that assesses the model against
  the active authoring guide family. When `quality-md` is in evaluation scope, assess,
  analyze, report, and roll it up like any other area; its source being
  `QUALITY.md` does not remove it from factor, local, or aggregate ratings.*

#### Write a description that distinguishes, not enumerates

- **Do** state what the area *is* and how it differs from its siblings/parent.
  **Avoid** restating its factors or requirements in the description. *The
  description identifies the evaluated entity; factors and requirements define
  what is assessed about it.*
