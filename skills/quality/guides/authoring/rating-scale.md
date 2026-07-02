---
type: Runtime Guide
title: Authoring rating scales
description: Rating scale, rating level, criterion, roll-up, veto, and not-assessed guidance for QUALITY.md models.
tags: [quality, authoring, guide]
---

# Authoring rating scales

Read this when:

- creating, revising, reviewing, or evaluating rating scales, rating level criteria, roll-up expectations, veto requirements, or `not assessed` handling.

Depends on:

- `../authoring.md`

---

## Rating scale

The **rating scale** is the fixed set of levels every requirement result is
rated against — the model's shared vocabulary for "how good." It is a list of
**rating levels**, ordered best (first) to worst (last), with at least two
levels.

Each level does two distinct jobs through two properties:

- **`description`** — what the level _means_: its standing in the scale and its
  intent. Fixed for the whole model; never overridden.
- **`criterion`** — the default rule for deciding whether a requirement's
  findings _land at_ that level. A requirement may override its own criterion
  (see [Override criteria only when the shared scale can't
  express the gradient](requirements.md#override-criteria-only-when-the-shared-scale-cant-express-the-gradient));
  the description never changes.

### Properties (per level)

| Property      | Presence    | What it is                                                      |
| ------------- | ----------- | --------------------------------------------------------------- |
| `level`       | Required    | Rating level ID; unique within the scale.                       |
| `title`       | Required    | Human-readable label for reports.                               |
| `description` | Recommended | What the level means across the model (fixed).                  |
| `criterion`   | Required    | Default rule for rating a requirement's findings at this level. |

### Working with the rating scale

#### Reuse the suggested four-level scale unless you have a reason not to

- **Consider** the **outstanding > target > minimum > unacceptable** scale as a
  default — a stretch level, the level to aim for, the floor you've agreed to
  rely on, and below the floor. _A shared four-band vocabulary is enough for
  most models and keeps reports comparable._
- **Do** use the default display titles `🟢 Outstanding`, `🔵 Target`,
  `🟡 Minimum`, and `🔴 Unacceptable` when the standard scale fits. _The emoji is
  only a scanning aid for human reports and frontmatter; the plain `level` IDs
  still carry identity, ordering, and references. Use plain or custom titles when
  a project style demands it, but avoid emoji-only labels._
- **Do** choose a different scale only when the root area demands it (e.g. a
  pass/fail gate wants two levels). _The scale should fit how decisions are
  actually made about this entity._
- **Do** review the scale after writing the body. _The body should reveal whether
  "good enough" and "excellent" need sharper meaning for this root area._
- **Do** test adjacent levels against the same plausible finding. _If one
  finding could reasonably satisfy both `target` and `minimum`, sharpen the
  criteria until the boundary is usable._
- **Do** let the model's job pick the number of levels: a model that gates against
  an agreed bar can use two levels (pass / fail); a model that judges _how good_,
  compares options, or surfaces strengths and weaknesses earns a graded scale. _A
  pass/fail model deliberately gives up the ability to say "how good" — sometimes
  exactly what you want, sometimes a loss._
- **Do** fix the _required margin_ in the body: how far above `minimum` this
  root area must actually land to be good enough, and why (safety, reputation,
  regulatory exposure, temporary use). _The same four levels can serve very
  different domains and stakes; what differs is the band the root area must
  reach. State it so an evaluator knows whether an all-`minimum` result is fine
  or alarming._
- **Do** calibrate the levels against concrete exemplars where you can — a real
  artifact you would call `outstanding` and one you would call `unacceptable`.
  _Thresholds set against known cases are far more defensible than ones invented
  in the abstract; if no real artifact would ever earn `outstanding`, the level is
  decorative._
- **Avoid** inventing a custom scale before the body reveals a real need for
  one.

#### Keep `description` about meaning and `criterion` about rating

- **Do** write `description` as what the level _is_ ("the floor you've agreed to
  rely on") and `criterion` as the test a result must pass ("falls short of
  target but remains acceptable"). _Conflating them makes per-requirement
  criterion overrides impossible to write cleanly._
- **Avoid** putting thresholds or measurements in `description`. _Those belong in
  `criterion`, where a requirement can override them._
- **Do** prefer a _measurable_ boundary in `criterion` whenever the concern admits
  one — a threshold, a count, a proportion, a band. _Unclear, unmeasurable criteria
  are a leading source of requirements debt; a measurable boundary is what makes a
  finding land at one level rather than hover between two._
- **Avoid** subjective or comparative wording in a `criterion` ("user-friendly",
  "better than", "as appropriate", "minimal"). _Such terms cannot reliably separate
  two adjacent levels; if you cannot operationalize a criterion, treat it as a
  smell pointing back at the requirement statement._

---

#### Decide how ratings roll up

_Read this subsection when you are reasoning about roll-up or evaluating; a first
model can defer it._

An area's local rating and a factor's rating are inferred from the requirement
results beneath them. The format fixes no aggregation formula, so how those
results combine is a modeling decision you should make and communicate.

- **Do** decide and state how an area's requirements combine into its rating when
  it is not obvious. _Two defensible defaults: **worst-of** (the weakest finding
  sets the rating — right when any unacceptable requirement makes the whole
  untrustworthy) and **most-common / median** (right when requirements genuinely
  compensate for one another)._
- **Avoid** implying that rating levels _average_. _The scale is ordinal
  (outstanding > target > minimum > unacceptable); the arithmetic mean of ordinal
  levels has no meaning — three `target`s and one `unacceptable` is not "slightly
  above minimum."_
- **Do** identify the requirements that can **veto** a rating — a single
  `unacceptable` finding that makes the root area not good enough no matter how
  strong everything else is (secrets stored in plaintext, data loss on failover).
  _Most requirements trade off against each other; a few cap the rating. Sharpen a
  veto requirement's `unacceptable` criterion and name its role in the body, since
  a compensating roll-up can otherwise hide a critical problem behind strong
  siblings._
- **Do** say when requirements are not equally important, so a reader does not read
  the roll-up as one-vote-each.
- **Do** keep a `not assessed` result distinct from a low rating in roll-up.
  _Missing evidence is not a failure; it must stay visible rather than count as a
  weak pass._

---
