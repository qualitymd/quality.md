---
type: Functional Specification
title: Model query commands — functional spec
description: Requirements for the read-only `qualitymd model` command group (tree, list, get) that projects a model's structure and canonical reference IDs.
tags: [cli, model, query]
timestamp: 2026-06-26T00:00:00Z
---

# Model query commands — functional spec

Companion to the [Model query commands](../0125-model-query-commands.md) change
case. This spec states _what_ the `model` command group must do; a later design
doc covers _how_. The model structure and the canonical reference grammar are
defined by [`SPECIFICATION.md`](../../../SPECIFICATION.md) (normative); the
invocation-wide CLI contract — output posture, exit codes, agent accessibility —
is defined by [`specs/cli.md`](../../../specs/cli.md) (normative). Design conventions
are in [Designing CLI interfaces](../../../docs/guides/cli-design.md)
(informational).

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

Authoring evaluation payloads requires the canonical reference IDs of a model's
elements, and no command emits them — agents hand-derive them from `QUALITY.md`.
`model` makes the model's elements, their canonical IDs, and their containment a
read-only, queryable projection for both an agent (`--json`) and a person (tree).
See the change case
[Motivation](../0125-model-query-commands.md#motivation) for the originating
evidence.

The group is deliberately bounded to _structure and identity_ so it does not
overlap two existing surfaces: `status` owns project state/readiness and source
coverage, and `evaluation` owns runs, payloads, reports, and snapshots. Keeping
those lines clean is a primary design constraint of this case, not an
afterthought.

## Scope

Covered: the `model` noun and its `tree`, `list`, and `get` verbs.

Deferred / non-goals (each absence is deliberate):

- No `model query` selector/expression language.
- No mutation — `model` is read-only.
- No model validation — `lint` owns diagnostics; `model` assumes a parseable
  model.
- No source/provenance, coverage, or aggregate count summaries — `status` owns
  those.
- No evaluation-run awareness — no `--run` flag and no run-folder resolution.
- No ratings or evaluation results.

## Requirements

### Boundary

- The `model` command group **MUST NOT** emit source declaration/provenance,
  source coverage, aggregate element-count summaries, evaluation results, or
  readiness state. Those remain the responsibility of `status` and `evaluation`.

  > > Rationale: `model` is introduced as a _non-overlapping_ surface; duplicating
  > > `status`'s coverage/counts or `evaluation`'s results would create two homes
  > > for one fact. — 0125

### Shared behavior (every verb)

- **SB1 — Target.** Each `model` verb **MUST** read a model file — the default
  `./QUALITY.md`, or a `[path]` argument when given — and **MUST NOT** be aware of
  evaluation runs, run IDs, or run-folder layout.

  > > Rationale: a run's `model-snapshot.md` is itself a model file, reachable by
  > > path; resolving runs would couple `model` to the `evaluation` surface this
  > > case keeps separate. — 0125

- **SB2 — Output posture.** Each verb **MUST** default to a human-readable form on
  `stdout` and **MUST** offer `--json` for the structured form. Terminal styling
  **MUST NOT** change the words, order, or facts of the output. No `--format`
  flag.

- **SB3 — Canonical IDs.** Every element a verb emits **MUST** carry its canonical
  qualified reference string, identical to the form persisted in evaluation
  payloads (the root area is `area:root`).

- **SB4 — Determinism.** For the same model file and content, each verb's output
  **MUST** be byte-stable under a single documented element ordering. No sort or
  ordering flag is offered.

- **SB5 — Exit codes.** A verb **MUST** exit `0` on success, `2` on a malformed
  invocation (unknown flag, invalid `--type` value, an `--area` or `<id>` that
  does not resolve), and `70` when the model file cannot be read or parsed. The
  group **MUST NOT** return `1`; it reports no negative result of its own.

  > > Rationale: an unresolved `--area`/`<id>` is a _usage_ error — the model is
  > > valid, the request named a non-element — not a found-problem (`1`). — 0125

- **SB6 — Flags.** Beyond the shared `[path]`, `--json`, `-h/--help`, and
  `--no-color`, each verb **MUST** accept only the flags named for it below and
  **MUST NOT** offer write, output-redirection, dry-run, quiet, or verbose flags.

- **SB7 — `--area` vocabulary.** Where a verb accepts `--area`, its value **MUST**
  be a canonical area reference (`area:<path>`), the same form `model` emits, so
  the value round-trips with `list`/`get` output. A value that is not a canonical
  area reference, or that does not resolve to an area in the model, is a usage
  error.

  > > Rationale: one addressing vocabulary. Accepting a bare path too would create
  > > a second way to name an area that the tool must then keep consistent. — 0125

### `model tree`

- **T1.** `model tree` **MUST** render the model as a hierarchy — the rooted area,
  its factors with sub-factors nested, its requirements, then child areas
  recursively. The human form is an indented tree; under `--json` it is a nested
  structure where every node carries its canonical `id`, `label`, and element
  `kind`.
- **T2.** `model tree` **MUST** accept `--area <area-id>` to root the output at
  that area's subtree instead of the whole model.
- **T3.** `model tree` **MUST** accept `--depth <n>` to limit nesting, where `0`
  emits only the rooted node.

### `model list`

- **L1.** `model list` **MUST** emit a flat enumeration of model elements, each
  with its canonical `id`, `kind` (`area` | `factor` | `requirement`), `label`,
  and parent `id` (absent for the root area).
- **L2.** `model list` **MUST** accept `--type` restricting output to one or more
  element kinds (`area`, `factor`, `requirement`), and `--area <area-id>`
  restricting output to one area's subtree; the two **MAY** be combined.
- **L3.** With no filter, `model list` **MUST** enumerate every element of every
  kind in the model.

### `model get`

- **G1.** `model get <id>` **MUST** accept a canonical reference id as a required
  positional argument and emit that element's `id`, `kind`, `label`, and its
  structural relationships — for an area, its immediate factor IDs, requirement
  IDs, and child-area IDs; for a factor, its sub-factor IDs — as an object under
  `--json` and a human-readable form otherwise. It **MUST NOT** require a `--type`
  hint, since the id prefix carries the kind.
- **G2.** If `<id>` does not resolve to an element in the model, `model get`
  **MUST** exit `2`, name the unresolved id, and **SHOULD** suggest the closest
  matching element ids.

### Skill integration (evaluate workflow)

The pain this case answers is the evaluate workflow hand-deriving canonical IDs
(see [Motivation](#background--motivation)). Shipping `model` without wiring it
into that workflow would leave the originating pain unsolved, so adoption is a
required outcome of this case, not an optional follow-up. Rollout MAY land the
CLI verbs and the skill edit as separate commits; the contract below is what the
change must deliver.

- **SK1 — Source of truth.** The evaluate workflow **MUST** obtain the canonical
  IDs of in-scope Areas, Factors, and Requirements from `model` (a single
  `model list --json` query for the scope), and **MUST** author every payload
  reference (`EvaluationFrame`, `AreaEvaluationFrame`,
  `RequirementEvaluationFrame`, `FactorAnalysisFrame`, `AreaAnalysisFrame`) from
  that result rather than deriving IDs from `QUALITY.md` text.

  > > Rationale: the orchestrator hand-deriving tens of IDs is the cited pain; a
  > > queried, authoritative ID set removes the silent-typo class at authoring
  > > time rather than catching it after the fact via identity resolution and the
  > > 0124 reference-`kind` enum. — 0125

- **SK2 — Query the snapshot.** Once the run exists, the workflow **MUST** query
  the run's `model-snapshot.md` by path, not the live `./QUALITY.md`, so the IDs
  it authors against are exactly the frozen model being evaluated. This relies on
  **SB1**: a snapshot is a model file reachable by path, with no run awareness in
  `model`.

- **SK3 — Scope alignment.** The workflow **MUST** request the scope it is
  evaluating with `--area <area-id>` (and `--type` when narrowing to one kind),
  so the projected ID set is exactly the set of elements the run authors payloads
  for. The `model list` filter vocabulary is anchored to evaluate's scope model
  (full vs. Area/Factor narrowing); the two surfaces **MUST** stay aligned.

- **SK4 — Label resolution (SHOULD).** The workflow **SHOULD** use `model`
  (`list` labels, `get`) to resolve a natural-label scope to its canonical
  `area:`/`factor:` reference, replacing ad-hoc matching against Area/Factor
  titles in the scope-resolution decision tree.

## Acceptance criteria

Shared:

- [ ] `model tree`, `model list`, and `model get area:root` each exit `0` against
      a valid `QUALITY.md`, with no `[path]` argument.
- [ ] Every verb accepts `--json`; piped or `--json` output carries no styling
      bytes and is byte-identical across repeated runs of the same file.
- [ ] Every id emitted by any verb round-trips: passing it to `model get`
      resolves the same element.
- [ ] Running a verb against an explicit `[path]` reads that file with no run
      awareness; pointing it at a `model-snapshot.md` yields that file's IDs.
- [ ] An unreadable or unparseable model file exits `70` with a human error that
      points at `lint`.
- [ ] No verb emits source/provenance, coverage, count summaries, or evaluation
      results.

`tree`:

- [ ] Human output nests root → factors (sub-factors nested) → requirements →
      child areas recursively; `--json` mirrors it with `id`/`label`/`kind` on
      every node.
- [ ] `--area area:agent-harness` roots the tree at that subtree; `--area` with a
      bare path or an unknown area exits `2`.
- [ ] `--depth 0` emits only the rooted node.

`list`:

- [ ] `model list --json` with no filter enumerates all areas, factors (with full
      nested factor paths), and requirements, each with `id`, `kind`, `label`,
      `parentId`.
- [ ] `--type factor` returns only factors; `--type bogus` exits `2` naming the
      allowed values.
- [ ] `--area area:agent-harness --type requirement` returns only that area's
      requirements.

`get`:

- [ ] `model get factor:client-app::performance` returns that factor's detail
      including its sub-factor ids.
- [ ] `model get requirement:nope::missing` exits `2`, names the id, and suggests
      near matches.

Skill integration:

- [ ] The evaluate workflow obtains in-scope element IDs from a
      `model list --json` query and authors every payload reference from that
      result — no ID is derived from `QUALITY.md` text.
- [ ] The query targets the run's `model-snapshot.md` by path, so authored IDs
      match the frozen evaluated model rather than a since-edited `QUALITY.md`.
- [ ] The query is scoped with `--area`/`--type` to the run's resolved scope, and
      the post-hoc identity-resolution check (with the 0124 `kind` enum) is a
      backstop, not the primary guard against ID typos.

## Durable spec changes

Durable **specs** this change rewrites — the
[`specs/`](../../../specs/index.md) bundle and
[`SPECIFICATION.md`](../../../SPECIFICATION.md). Each subsection is required.

### To add

- `specs/cli/model.md` — the durable command spec for the `model` group (its
  `tree`/`list`/`get` behavior, flags, output shapes, and the non-overlap
  boundary), per the requirements above.

### To modify

- `specs/cli.md` — add `model` to the command table and the `Commands` list (per
  the group's introduction above).
- `specs/cli/index.md` — register `model.md`.
- `specs/skills/quality-skill/workflows/evaluate.md` — add the normative
  requirement that payload references come from `model` (SK1–SK4), not from
  hand-derived `QUALITY.md` IDs.

### To rename

None.

### To delete

None.

## Open questions

- **Projection home and reference encoding.** Resolved by the
  [design doc](design.md#a-neutral-home-in-internalmodel): the shared projection
  and the reference grammar move into `internal/model`, so `model` depends on
  neither `status` nor `evaluation`.
- **Verb rollout sequencing.** Resolved by the design doc: the three verbs share
  one preflight and projection and land in a single slice.
