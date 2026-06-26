---
type: Design Doc
title: Model query commands
description: How the read-only `qualitymd model` group projects a model's elements and canonical reference IDs from a neutral home, and how the evaluate workflow consumes them.
status: Draft
tags: [cli, model, query]
timestamp: 2026-06-26T00:00:00Z
---

# Model query commands

## Context

Implements the [0125 functional spec](spec.md): a read-only `model` group
(`tree`, `list`, `get`) that projects a model's elements, their canonical
reference IDs, and their containment from a model file, plus wiring the
`/quality` evaluate workflow to consume those IDs instead of hand-deriving them.

The spec leaves two questions to design: where the shared projection and the
reference encoding should live so `model` depends on neither `status` nor
`evaluation`, and how the verbs are sequenced. The skill-integration
requirements (SK1–SK4) also need a concrete data flow.

## Approach

### A neutral home in `internal/model`

`internal/model` already owns the typed `Spec`/`Area`/`Factor`/`Requirement`
and is imported by both `status` and `evaluation`; nothing internal sits below
it except `document`. It is the natural home for both the projection and the
canonical-reference grammar.

1. **Move the reference grammar into `internal/model`.** The path types
   (`AreaPath`, `FactorPath`), the canonical encoders
   (`AreaPath.Reference()` / `FactorReference()` / `RequirementReference()` and
   their unqualified forms), and the parsers
   (`ParseAreaReference` / `ParseFactorReference` / `ParseRequirementReference`,
   which already take a `*model.Spec`) move from `internal/evaluation/types.go`
   and `internal/evaluation/model_reference.go` into a new
   `internal/model/reference.go`. `evaluation` updates its call sites to the
   `model.*` names. Rating references (`RatingReference`, rating-scale parsing)
   stay in `evaluation`: they address the rating scale, not model structure, and
   `model` never emits them.

   This is what lets the `model` command both *emit* canonical IDs and *resolve*
   `--area`/`<id>` inputs without importing `evaluation` — the parser it needs
   comes home with the encoder.

2. **Add a projection walk.** A single `internal/model/projection.go` walk
   produces the element set both verbs render from:

   ```go
   type Element struct {
       ID        string      // canonical reference, e.g. "factor:client-app::performance"
       Kind      Kind        // area | factor | requirement
       Label     string      // Title, falling back to the map key
       ParentID  string      // "" for the root area
       Children  []*Element  // populated for the tree projection
   }
   ```

   One recursive walk visits the rooted area, its factors (sub-factors nested),
   its requirements, then child areas — the documented order in **SB4**. `list`
   flattens it; `tree` keeps it nested; `get <id>` resolves the id with the
   moved parser and returns that element plus its immediate relations. Ordering
   is lexicographic by map key, reusing the same `sortedKeys` discipline
   `status` already applies, so output is byte-stable with no sort flag.

3. **Retire the duplicate walk in `status`.** `modelAccumulator.walkAreas` /
   `walkFactors` accumulate counts and coverage, not IDs; once the shared
   projection exists, `status` builds its shape/coverage from a single pass over
   the same projection rather than its own recursion. This removes the second
   model-tree walk the change case flagged.

### CLI surface

`internal/cli/model.go` adds the `model` parent and its three subcommands,
registered in `root.go` in the common group beside `lint`/`status`/`spec`/
`schema`. Each verb shares one preflight: resolve `[path]` (default
`./QUALITY.md`), `document.Parse` → `model.Decode`, then project. The shared
helpers already in `root.go` carry the exit-code contract (**SB5**):

- a read/parse failure returns the internal coded error → exit `70`, with a
  message that points at `lint`;
- an unknown flag, an invalid `--type` value, or an `--area`/`<id>` that does
  not resolve returns `usageError` → exit `2`;
- the group never returns `ExitProblems` (`1`) — it reports no negative result.

`--json` selects the structured encoder; the human form is the default. Styling
is applied only to the human form and never reorders or changes facts (**SB2**).

Verbs land in one slice — they share the preflight, the projection, and the
test fixture, so splitting them buys nothing.

### Skill integration (SK1–SK4)

No new CLI mechanism is needed; the snapshot is just a model file by path
(**SB1**/**SB2**). After step 8's `evaluation create` writes
`model-snapshot.md`, the evaluate workflow composes that path from the run
folder it already holds and runs, once:

```
qualitymd model list --json <run>/model-snapshot.md --area <resolved-scope>
```

The returned `id`/`kind`/`parentId` set becomes the source of truth for every
payload reference authored in steps 9–15 (`EvaluationFrame`,
`AreaEvaluationFrame`, `RequirementEvaluationFrame`, `FactorAnalysisFrame`,
`AreaAnalysisFrame`). The workflow stops deriving IDs from `QUALITY.md` text.
`model get`/`list` labels also serve the scope-resolution decision tree's
label → canonical-reference step (SK4). Because `model` has no run awareness,
the workflow — not the CLI — knows the run path and builds the snapshot path.

## Spec response

- **Boundary / SB1–SB7.** `model` reads a model *file* and depends only on
  `internal/model`; it emits no coverage, counts, provenance, or results because
  it has no access to the `status`/`evaluation` types that hold them. `--area`
  parses through the moved `ParseAreaReference`, so it accepts only the
  canonical `area:<path>` form and round-trips with emitted IDs (**SB7**).
- **T1–T3 / L1–L3 / G1–G2.** All render from the one projection: `tree` nested
  with `--area` rooting and `--depth` truncation, `list` flat with `--type`/
  `--area` filters, `get` resolving an id to its element and immediate relations.
  An unresolved `get` id returns `usageError` and suggests near matches by
  lexical distance over the projected ids (**G2**).
- **SB4 determinism** rests on lexicographic map-key ordering with no ordering
  flag — the verification-sensitive choice.
- **SK1–SK4** are satisfied by the snapshot-path query above; no CLI run
  resolution is introduced.

## Alternatives

- **New `internal/modelquery` package depending on `evaluation`** — rejected:
  the spec forbids `model` depending on `evaluation`, and the reference parser it
  needs lives there today.
- **Duplicate the tiny encoders in `model`, leave them in `evaluation` too** —
  rejected: two homes for one grammar is exactly the silent-divergence risk 0124
  and 0125 exist to remove.
- **Type aliases (`type AreaPath = model.AreaPath`) to avoid touching
  `evaluation` call sites** — rejected as a compatibility shim; the early-alpha
  rule prefers updating call sites to the moved names in one clean break.
- **Add a `--run` flag so `model` resolves the snapshot itself** — rejected:
  run-folder resolution is the `evaluation` surface this case keeps separate
  (**SB1**); the workflow composing the path costs nothing.
- **Slice the verbs across commits** — rejected: shared preflight and projection
  make a single slice simpler than coordinating three.

## Trade-offs & risks

- **Moving the reference grammar touches `evaluation` broadly.** The encoders
  and `AreaPath`/`FactorPath` are referenced across `evaluation`; the move is a
  wide but mechanical rename, covered by the existing evaluation tests plus the
  canonical-ID round-trip tests this case adds.
- **Two consumers, one walk.** Folding `status` onto the shared projection risks
  a subtle shift in its counts/coverage output; the `status` golden tests pin
  that and must stay green through the refactor.
- **Skill/CLI land separately.** If the skill edit lags the CLI verbs, the
  originating pain persists until it lands; the change case marks adoption a
  required outcome, not optional, to keep that from being dropped.

## Open questions

- **Suggestion ranking for `get`.** Whether near-match suggestions rank by
  edit distance, shared id prefix, or both — a quality detail, not a contract;
  resolved during implementation against the `get` acceptance tests.
- **Lifting a projection data-contract spec.** If the `--json` shape grows, the
  projection may warrant its own 1:1 artifact spec (the change case's "Suggested
  new durable specs"); deferred until the shape settles.
