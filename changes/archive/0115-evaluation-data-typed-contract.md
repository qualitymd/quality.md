---
type: Change Case
title: Type-safe, model-bound Evaluation v2 data
description: Make `qualitymd evaluation data set` validate payloads against typed per-kind definitions and the run's model snapshot, generate the schema and populated examples from one source, and add a `data schema` discovery command.
status: Done
tags: [cli, evaluation, data, validation, agent-mediated]
timestamp: 2026-06-26T00:00:00Z
---

# Type-safe, model-bound Evaluation v2 data

This parent concept captures the why and status; the detail lives in its
children:

- [Functional spec](0115-evaluation-data-typed-contract/spec.md) — what the
  change must do.
- [Design doc](0115-evaluation-data-typed-contract/design.md) — how it is built,
  and why.

## Motivation

A field evaluation of an external project surfaced two failures in the
Evaluation v2 data surface, both confirmed against the code.

1. **Loose validation lets bad data through.** `qualitymd evaluation data set`
   decodes every payload into `map[string]any` and checks only a handful of
   named fields per kind. Unknown and misspelled fields pass silently: a
   `RequirementAssessmentResult` whose findings used `title`/`summary` instead of
   `description`/`evidence` was accepted and persisted, and only surfaced when the
   built report rendered blank. The evaluator's own feedback log records the
   wasted loop ("First report pass rendered blank Description/Evidence because
   findings used `title/summary`").

2. **No model binding lets invented IDs through.** `data set` validates
   `areaId` / `factorId` / `requirementId` / `ratingLevelId` only for
   filesystem-safe characters — never that the referenced Area, Factor,
   Requirement, or Rating Level exists in the model. Any safe free-form string is
   accepted and written to a path, so the persisted graph can reference model
   nodes that do not exist. Existence-checkers already live in
   [`internal/evaluation/model_reference.go`](../../internal/evaluation/model_reference.go)
   (`areaExists` / `factorExists` / `requirementExists`) but are wired only into
   tests, not into `data set`.

Both failures trace to one root: there is **no typed source of truth** for the
v2 data kinds. The "schema" is spread across three hand-maintained, drifting
copies — a per-kind validator, the `data example` builders, and the on-disk
shape — and an agent cannot discover the real sub-shapes. The same feedback log
shows the discovery cost: finding fields, the analysis status enum, and the
`ratingDrivers` / `unknowns` object shapes "are not in `data example` output
(arrays shown empty)… had to reverse-engineer via `--dry-run` probes and the
platform binary."

The repository already has the idiomatic answer to this in
[`internal/schema`](../../internal/schema): a typed Go model of the frontmatter,
a generated JSON Schema (`go generate`), and a drift/determinism test. The v2
data surface is the only structured surface that does **not** follow that
pattern. This case brings it into line: one typed source of truth feeds
validation, a generated schema, and populated examples — and adds the semantic
model-binding pass as the v2-data analogue of what `qualitymd lint` already does
for the frontmatter.

## Scope

Covered:

- define each accepted Evaluation v2 data kind (and its shared sub-shapes and ID
  value types) as typed Go definitions in one place — the single source of truth;
- strict structural validation in `data set`: reject unknown / misspelled fields,
  wrong field types, and out-of-range enum values, naming the offending field;
- required-field validation per kind;
- **model-binding** validation: resolve every Area, Factor, Requirement, and
  Rating Level reference in a payload against the run's `model-snapshot.md` and
  reject references absent from it;
- run the full validation under `--dry-run` so the agent's probe loop catches the
  same errors a real write would;
- generate the JSON Schema and the populated `data example` payloads from the
  typed definitions, guarded by a drift/determinism test;
- a new `qualitymd evaluation data schema [<kind>]` discovery command;
- accept the kebab-case form of a `<kind>` argument alongside the canonical name;
- a `qualitymd evaluation data verify <run>` command that re-validates every
  persisted payload in a run against the full contract (the migration / self-check
  path for the tightened acceptance rules);
- improve the `/quality` skill's payload-discovery procedure to use
  `data schema` (and the now-populated `data example`) as the authoritative source
  of payload shape, and reframe `data set --dry-run` as validation of an authored
  payload rather than the loop used to _sniff_ shape — updating both the durable
  skill spec's discovery requirement and the runtime skill resources.

Deferred / non-goals:

- no change to the v2 wire format: `schemaVersion` stays `1` and a currently
  _valid_ payload's on-disk JSON is byte-unchanged — only acceptance is tightened;
- no change to the set of accepted kinds, the protocol order, or the routine
  prompt contracts;
- no retroactive rewrite of already-persisted data: existing files are only
  re-validated on demand via `data verify` or on re-`set`;
- no runtime JSON-Schema engine — validation is typed decode plus a per-kind
  semantic pass; a schema library, if used, is a generate-time tool only (settled
  in the design doc);
- no change to the QUALITY.md format itself or `SPECIFICATION.md`.

## Affected artifacts

Derived by sweeping the repo for the data-set/validate/example code path, the v2
data kinds, and the `data` command surface across `internal/`, `cmd/`, `specs/`,
`skills/`, and `docs/`.

### Code

- [ ] [`internal/evaluation/data.go`](../../internal/evaluation/data.go) —
      replace `map[string]any` decode + the per-kind `validateDataPayload` switch + the empty-array example builders with typed decode and registry dispatch.
- [ ] New package (e.g. `internal/evaluation/datakind/`) — typed per-kind
      definitions, shared sub-shapes, ID value types, enums, the kind registry,
      and `Validate(spec)`; plus a `gen/` generator behind `//go:generate`,
      mirroring [`internal/schema`](../../internal/schema).
- [ ] Generated, checked-in JSON Schema artifact for the v2 data kinds, guarded
      by a drift/determinism test (mirrors `internal/schema/jsonschema_test.go`).
- [ ] [`internal/cli/evaluation.go`](../../internal/cli/evaluation.go) — add
      `data schema` and `data verify`; populate `data example`; accept kebab-case
      `<kind>` arguments.
- [ ] [`internal/evaluation/load.go`](../../internal/evaluation/load.go) —
      extract the snapshot→`*model.Spec` load into a shared helper for the
      model-binding pass.
- [ ] [`internal/evaluation/model_reference.go`](../../internal/evaluation/model_reference.go)
      — reuse the existing `areaExists` / `factorExists` / `requirementExists`
      helpers from the `data set` path (no longer test-only).
- [ ] [`internal/evaluation/report_v2.go`](../../internal/evaluation/report_v2.go)
      — optionally consume the typed structs instead of re-parsing maps (in scope
      only if it falls out cleanly; otherwise deferred).
- [ ] Tests for the above (strict-decode rejections, model-binding rejections,
      generated-artifact drift, example round-trip).

### Format spec

- [ ] None — `SPECIFICATION.md` and the QUALITY.md format are unaffected; this is
      the Evaluation v2 tooling contract, not the model format.

### Durable specs

- [ ] [`specs/cli/evaluation-data.md`](../../specs/cli/evaluation-data.md) —
      strengthen the `data set` validation requirements (reject unknown fields,
      model-binding, dry-run parity); document `data schema` and `data verify`;
      tie `data example` to the schema.
- [ ] [`specs/evaluation-v2/records/payload-kinds.md`](../../specs/evaluation-v2/records/payload-kinds.md)
      — "the CLI **MUST** validate each accepted kind" becomes structural + typed + model-bound, from one source of truth.
- [ ] [`specs/evaluation-v2/records/json-conventions.md`](../../specs/evaluation-v2/records/json-conventions.md)
      — note that persisted identity fields are resolved against the run's model
      snapshot, and that unknown fields are rejected.
- [ ] [`specs/skills/quality-skill/quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md)
      — the structured-payload discovery requirement (currently `data kinds` →
      `data example` → `data set --dry-run`) gains `data schema` as the
      authoritative shape source and reframes `--dry-run` as payload validation,
      not shape discovery.

### Durable docs / bundled skill

- [ ] [`skills/quality/resources/cli-quick-reference.md`](../../skills/quality/resources/cli-quick-reference.md)
      — make `data schema` the "discover payload shapes" entry; reframe the
      `--dry-run` entry as validation; drop the reverse-engineering workaround.
- [ ] [`skills/quality/workflows/evaluate.md`](../../skills/quality/workflows/evaluate.md)
      — point payload-shape discovery at `data schema`; `--dry-run` validates an
      authored payload.
- [ ] [`skills/quality/SKILL.md`](../../skills/quality/SKILL.md) — update the
      data-discovery command surface to include `data schema`.

### Suggested new durable specs

- A 1:1 artifact spec for the generated JSON Schema file (artifact-spec naming,
  e.g. `…-schema-json.md`), once the artifact path is settled in the design.
  Suggesting is enough; creating it is not a precondition for landing.

## Status

`Done`. Implemented, verified, and archived. The data surface now uses a single
typed contract for validation, schema, examples, model-binding, and verification;
the durable specs and `/quality` skill resources are aligned.
