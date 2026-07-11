---
type: Functional Specification
title: Resolver-dispatched source selectors — functional spec
description: Requirements for resolving source selectors per kind into the bounded evidence bundle, including harness-dispatched resolution and provenance of record.
---

# Resolver-dispatched source selectors — functional spec

Companion to the
[Resolver-dispatched source selectors](../0197-resolver-dispatched-source-selectors.md)
change case. This spec states _what_ the change must do; a design doc will
cover _how_ once the case reaches Design. It builds on the resolution model
recorded in 0196's
[considerations sketch](../archive/0196-spec-faithful-model-reading/considerations.md):

```
source selector ──[ resolver (per kind) ]──► bounded, hashed evidence bundle ──[ evaluator ]──► judgment
```

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Scope

Covered: per-kind selector resolution, harness-dispatched resolution for
non-deterministic kinds, classified failure for unsupported kinds, and bundle
provenance. Deferred: out-of-tree refs, evaluator adapter fragility,
default-selector granularity, and a catalog of resolver kinds beyond the seam
itself — see the change case
[scope](../0197-resolver-dispatched-source-selectors.md#scope).

The two questions that gated this requirement set are settled — the format
commits to non-filesystem selectors, and kind is detected from the bare
selector string — see [Settled questions](#settled-questions).

## Requirements

### Resolution

**R1 — Kind-dispatched resolution.** The runner **MUST** resolve each
effective source selector according to its selector kind, detected from the
bare selector string in this order: a selector containing glob metacharacters
is a glob; otherwise a selector that names an existing filesystem entry is a
path; otherwise the selector is prose. A selector that is absolute or escapes
the workspace **MUST** remain a filesystem selector under the existing
workspace-containment rules and **MUST NOT** fall back to prose. Path and glob
selectors **MUST** keep the deterministic resolution and packaging contract of
[`specs/evaluation/runner.md`](../../specs/evaluation/runner.md) §Source
packaging unchanged. A run's detected kinds **MUST** be recorded at run
creation and honored on resume, so a filesystem change mid-run cannot silently
re-dispatch a selector to a different resolver.

> _Why:_ the resolver seam must not alter the spec-faithful path/glob behavior
> 0196 restored; it generalizes around it. Detection order makes the
> filesystem interpretation always win, so prose is only ever the meaning of a
> selector that cannot be filesystem material — and `source` stays a
> hand-authorable scalar with an unchanged frontmatter shape (Q1). — 0196,
> Q1/Q2 resolutions
> _Durable spec:_ modify `specs/evaluation/runner.md` §Source packaging —
> introduce the resolution step, kind detection, and its dispatch rule; modify
> `SPECIFICATION.md` — commit the format to non-filesystem selectors with
> `source` remaining a single string (Q2).

**R2 — Judgment consumes only captured bundles.** Material gathered by any
resolver **MUST** be captured into the same bounded, hashed, persisted source
bundle before judgment, and a judgment work unit **MUST NOT** receive evidence
that is not part of a captured bundle.

> _Why:_ the runner's guarantees — resumability, input-hash guards,
> evidence-bound re-judgment, source-as-data — are properties of the resolved
> bundle, not of how it was gathered; the 2026-07-11 re-validation confirmed
> they hold through the bundle alone. An agent-gathered "all specs" is
> legitimate only through the bundle contract, never as free repo exploration
> by the judge.
> _Durable spec:_ modify `specs/evaluation/runner.md` §Source packaging —
> state the bundle contract as resolver-independent.

**R3 — Harness-dispatched resolution.** When an effective selector's kind has
no deterministic resolver and the run's evaluator dispatch supports harness
checkpoints, the runner **MUST** emit a resolution work request through the
same checkpoint transport as harness judgment, and **MUST** validate and
capture the returned material into the bundle (paths, content, hashes) before
any dependent judgment work is dispatched.

> _Why:_ a prose or live-system selector needs tools to gather; the harness
> checkpoint transport already carries correlated, validated, resumable work
> both ways. Resolution and judgment stay distinct requests so the gatherer
> and the judge are never the same uncontrolled step.
> _Durable spec:_ modify `specs/evaluation/protocol.md` (resolution request
> kind and envelope), `specs/evaluation/evaluator-contract.md` (the
> capability that advertises resolution support), and
> `specs/skills/quality-skill/evaluation.md` (the harness-side workflow
> serves resolution work requests alongside judgment requests).

**R4 — Unsupported selector kinds fail loudly and distinctly.** If no resolver
is available for an effective selector — its kind is unsupported, or the
selected evaluator dispatch cannot serve resolution requests — then the runner
**MUST** fail the affected work with a classified failure naming the selector
and why it cannot be resolved, distinct from `source_unavailable` (material
missing), and **MUST NOT** fall back to the document default or an empty
bundle.

> _Why:_ today a prose selector fails as `source_unavailable`, which
> misdiagnoses an unsupported kind as missing material and tells the author to
> "add the material it names". The remedy differs (change the selector or the
> evaluator, not the material), so the classification must too.
> _Durable spec:_ modify `specs/evaluation/runner.md` §Failure taxonomy — add
> the classification and its remedy wording.

### Record

**R5 — Bundle provenance of record.** For every packaged bundle, the run
artifact **MUST** record the selector, its resolved kind, and the resolver
that produced it (deterministic walk or harness resolution, with the harness
identity when applicable), alongside the existing hashes.

> _Why:_ reproducibility of record — audit, diffing, and re-judgment must read
> the same for deterministic and harness-gathered evidence, and a reviewer
> must be able to tell which evidence an agent gathered.
> _Durable spec:_ modify `specs/evaluation/evaluation-json.md` — provenance
> fields on the packaged-source record.

**R6 — Source-as-data across resolvers.** Resolver-returned material **MUST**
be presented to judging evaluators as data under the same standing safety
instructions as walked source, applying the source-as-data invariant in
[`specs/evaluation/evaluation.md`](../../specs/evaluation/evaluation.md#shared-invariants).

> _Why:_ the data-not-instructions boundary is the safety property that makes
> non-deterministic gathering acceptable; it must not weaken with the
> resolver.
> _Durable spec:_ none — the invariant already binds all packaged source;
> R2/R3 route resolver output through it.

## Settled questions

Both gating questions were settled 2026-07-11, before Design.

- **Q1 — Selector typing: bare string with detected kind.** `source` stays a
  single scalar; the runner detects kind in R1's order (glob metacharacters →
  glob; existing filesystem entry → path; otherwise prose). Chosen for
  hand-authoring friendliness and an unchanged frontmatter shape —
  `quality.schema.json` needs no change. The known hazard — a mistyped path
  that names nothing detects as prose and dispatches resolution instead of
  failing `source_unavailable` — is accepted; provenance (R5), loud per-kind
  surfacing in dry-run and receipts, and the resolution instructions'
  path-like-selector guidance mitigate it (see the [design](design.md)
  trade-offs).
- **Q2 — Format commitment: the format commits.** `SPECIFICATION.md` commits
  to non-filesystem selectors: a source selector is a string describing the
  material an area evaluates, resolved per kind by evaluating tools, with path
  and glob remaining the deterministic filesystem kinds. R3 and the resolution
  seam stand as drafted.

## Durable spec changes

Covers the [`specs/`](../../specs/index.md) bundle and
[`SPECIFICATION.md`](../../SPECIFICATION.md).

### To add

None.

### To modify

- `specs/evaluation/runner.md` — §Source packaging gains the per-kind
  resolution step with kind detection and path/glob unchanged (R1) and the
  resolver-independent bundle contract (R2); §Failure taxonomy gains the
  unsupported-selector classification (R4).
- `specs/evaluation/protocol.md` — the `resolveSource` move and its
  gathering-not-judgment flow (R3).
- `specs/evaluation/orchestration.md` — the `resolveSource` work unit in the
  graph vocabulary and its capture-before-judgment dependency (R2, R3;
  reconciled at implementation — the orchestration contract owns the graph
  shape the case changes).
- `specs/evaluation/evaluator-contract.md` — the capability advertising
  resolution support (R3).
- `specs/cli/evaluation-run.md` — the per-area source dispatch plan in
  dry-run previews and run receipts (R1's surfacing; reconciled at
  implementation).
- `specs/skills/quality-skill/evaluation.md` — the harness-side workflow
  serves resolution work requests (R3).
- `specs/evaluation/evaluation-json.md` — bundle provenance fields (R5).
- `SPECIFICATION.md` — per the Q2 resolution: commit the format to
  non-filesystem selectors, with `source` remaining a single string and kind
  detected rather than declared (Q1). `quality.schema.json` and
  `specs/quality-schema-json.md` are deliberately unchanged — the frontmatter
  shape does not move.

### To rename

None.

### To delete

None.
