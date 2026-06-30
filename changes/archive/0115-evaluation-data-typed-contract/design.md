---
type: Design Doc
title: Type-safe, model-bound Evaluation v2 data — design doc
description: Technical approach for typed per-kind definitions, strict decode, model-binding validation, and generated schema/examples in the Evaluation v2 data surface.
tags: [cli, evaluation, data, validation, go]
timestamp: 2026-06-26T00:00:00Z
---

# Type-safe, model-bound Evaluation v2 data — design doc

> 🚧 **Provisional sketch.** The change case is in **Draft**; this design is
> captured in parallel and settles once the [functional spec](spec.md) is
> reviewed. No code is written until the case reaches **In-Progress**.

Design behind the
[Type-safe, model-bound Evaluation v2 data](../0115-evaluation-data-typed-contract.md)
change case and its [functional spec](spec.md).

## Context

`internal/evaluation/data.go` decodes every payload into `map[string]any` and
validates it through three hand-maintained, drifting representations: a per-kind
`validateDataPayload` switch, the `data example` builders (which hard-code empty
arrays), and the implicit on-disk shape. Identity fields are checked only for
path-safe characters, never against the model. The spec turns this into four
obligations: strict typed validation (R1–R3), model-binding (R4, R7), one source
of truth feeding schema and examples (R6, R8, R9), and on-demand re-validation
(R11) — without changing the wire format (R12).

The repository already solves the structural half once, for the frontmatter:
`internal/schema` is a typed Go model with a `//go:generate`'d JSON Schema and a
drift/determinism test (`jsonschema_test.go`). This design applies that same
pattern to the v2 data kinds and adds the semantic (model-binding) pass that the
frontmatter delegates to `qualitymd lint`.

## Approach

**1. One typed source of truth — a `datakind` package.** A new package (working
name `internal/evaluation/datakind`) defines one Go struct per accepted kind, the
shared sub-shapes (`Finding`, `RatingDriver`, `Unknown`, `ScopedAnalysis`,
`CriteriaResult`, `EvidenceTarget`, …), the ID value types (`AreaID []string`;
`RequirementID{DeclaringAreaID, RequirementName}`;
`FactorID{DeclaringAreaID, FactorPath}`), and the enums as typed string constants
with a `Valid()` method (the status, confidence, finding-`type`, and `severity`
sets currently inlined in `validateDataPayload`). A registry maps each kind name →
`{ new() any; pathFor(v) (string, error); validate(v, *model.Spec) error }`. This
table is the single dispatch point every verb shares (answers R6).

**2. Strict structural decode (R1–R3).** `SetData` peeks `kind`, looks it up in
the registry, then decodes with `json.NewDecoder` + `DisallowUnknownFields()` into
the kind's struct. Unknown/misspelled fields and wrong types fail at the decoder;
enum and required-field checks run in the struct's `Validate`. Required fields use
explicit presence checks (or pointers for the few required scalars), since Go zero
values don't distinguish "missing." Error messages name the offending field; a
Levenshtein "did you mean `description`?" hint is a cheap optional polish given the
writer is an agent. Canonical-JSON re-marshal is preserved for stable on-disk
output, so a valid payload is byte-identical to today (R12).

**3. Model-binding pass (R4, R7).** Lift the snapshot→`*model.Spec` load from
`load.go` (`document.Parse(model-snapshot.md)` + `model.Decode`) into a shared
`loadRunSpec(runAbs)`. After structural decode, each kind's `Validate(spec)`
resolves its IDs with the existing `areaExists` / `factorExists` /
`requirementExists` helpers and `spec.RatingScale` (for `ratingLevelId`,
`criteriaResults[].ratingLevelId`, and the frame `ratingLevelIds`). A missing or
unparseable snapshot is a hard error naming the snapshot (R7) — fail closed. The
same `Validate` runs under `--dry-run` (R5), so the probe loop is faithful.

**4. Generated schema and populated examples (R6, R8, R9).** Mirror
`internal/schema`: a `datakind/gen` generator behind `//go:generate` emits a
checked-in JSON Schema (one document, one `$def` per kind), and the `data example`
output becomes a marshaled struct instance with hand-authored representative
sample values (one real `Finding`, `RatingDriver`, `Unknown`, etc.). Because the
example is an instance of the same struct the schema describes and the decoder
accepts, it cannot drift. `datakind/jsonschema_test.go` asserts determinism, that
every accepted kind has a `$def`, and that every example both validates against
its schema and round-trips through strict decode — this is the drift guard R6
requires.

**5. CLI surface.** `internal/cli/evaluation.go` gains `data schema [<kind>]`
(polymorphic, same shape as `data example <kind>`) and `data verify <run>`
(walks `data/**`, runs the full contract per file, exits non-zero on any failure —
R11). Kind-argument resolution normalizes kebab-case to the canonical name (R10).
The existing `set` / `list` / `get` / `kinds` / `example` verbs keep their
signatures; the polymorphism just moves from scattered `switch kind` blocks onto
the registry.

**6. Skill discovery (R13).** The durable skill spec's structured-payload
discovery requirement (`quality-skill.md`, currently `data kinds` → `data example`
→ `data set --dry-run`) and the runtime resources (`cli-quick-reference.md`,
`workflows/evaluate.md`, `SKILL.md`) are updated to lead with `data schema` as the
authoritative shape source and to cast `--dry-run` as validation of an authored
payload. This is what converts the new CLI capability into a behavior change: the
agent reads the contract instead of probing for it. No new mechanism — a
prose/spec edit that depends on R8/R9 landing first.

**7. Read side (optional).** `report_v2.go` can consume the typed structs instead
of re-parsing maps; in scope only if it falls out cleanly, otherwise deferred.

## Alternatives

- **Keep `map[string]any`, extend the hand validator.** Add per-kind key
  allowlists to catch wrong finding fields. Cheapest, but it adds a _third_
  hand-maintained copy of the shape, doesn't give a discoverable schema, and
  drifts from the examples — it treats the symptom, not the root (no source of
  truth). Rejected as the long-term answer; viable only as an interim patch.
- **Hand-model the schema like `internal/schema`.** Maximal consistency and zero
  deps, but `internal/schema` models a bespoke YAML shape system; v2 data is plain
  JSON, so reflecting structs is far less code than re-describing every field by
  hand. Rejected in favor of struct-derived generation.
- **Runtime JSON-Schema validation** (e.g. validate each payload against the
  emitted schema at `set` time). Doubles the validation path and still needs Go
  types for model-binding. Rejected: typed decode + `Validate(spec)` is simpler
  and keeps a single validation path; the schema is for _discovery_, not runtime
  enforcement.

## Trade-offs & risks

- **Dependency choice (the main open decision).** Runtime stays zero-dep
  (typed decode + `Validate`). For _generate-time_ schema emission, either
  hand-build the JSON Schema from the registry (zero deps, more code) or use a
  reflection library — `github.com/google/jsonschema-go` (now the canonical,
  LLM-oriented Go choice: inference via `jsonschema.For[T]()`, plus validation and
  serialization) or `invopop/jsonschema`. Recommendation: take the generate-time
  dep for far less code, keeping runtime dep-free — consistent with the
  frontmatter side's zero _runtime_ deps. Final call deferred to a prototype.
- **Required-vs-optional in Go.** Zero values hide "missing." Mitigated with
  explicit presence checks or pointers per kind; the example round-trip test
  guards against an under-specified `Validate`.
- **Tightening is a behavior change.** Payloads accepted before may now fail.
  Bounded by R12 (valid payloads unchanged) and surfaced for old runs by
  `data verify` (R11). No silent migration.
- **Snapshot drift.** Binding to the run's snapshot, not live `QUALITY.md`, keeps
  runs reproducible after the model changes — intended, and it removes any need to
  locate the live model.

## Open questions

- Package name and location: `internal/evaluation/datakind` vs a peer to
  `internal/schema`.
- Schema artifact path and whether it ships one combined document or one file per
  kind; and the matching `data schema` default (whole set vs require `<kind>`).
- Whether `data verify` is part of this case or a fast follow (spec'd here as
  in-scope; could split if the case grows).
- Whether to fold `report_v2.go` onto the typed structs now or defer.
- Prototype `RequirementAssessmentResult` end-to-end to confirm the hand-built vs
  reflected schema-generation decision before fanning out to all nine kinds.
