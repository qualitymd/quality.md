---
type: Change Case
title: String model-identity fields in evaluation data
description: Collapse the structured `areaId`, `factorId`, and `requirementId` fields in persisted Evaluation routine and report JSON into single canonical qualified model-reference strings, keep the `*Id` names, and reserve `root` as an Area name so the string form is lossless.
status: Done
tags: [evaluation, records, references, schema]
timestamp: 2026-06-26T00:00:00Z
---

# String model-identity fields in evaluation data

This parent concept captures the why and status; the detail lives in its
children:

- [Functional spec](0120-string-model-identity-fields/spec.md) - what the change
  must do.
- [Design doc](0120-string-model-identity-fields/design.md) - how the data
  contract, parser/encoder, schema, and lint changes land together.

## Motivation

Persisted Evaluation routine and report JSON carries the same model-identity
concept in three different physical shapes:

- `ratingLevelId` is a plain string;
- `areaId` is a string array (the Area path; `[]` for the root);
- `factorId` and `requirementId` are objects
  (`{declaringAreaId: [...], factorPath: [...]}` and
  `{declaringAreaId: [...], requirementName: "..."}`).

Anyone reading or consuming the report data meets three encodings for one idea —
"which model element is this?" — and has to index, dedupe, and compare a mix of
scalars, arrays, and nested objects. The structured shapes are not richer than a
string: [`SPECIFICATION.md`](../SPECIFICATION.md) already defines a canonical
qualified model reference for each kind (`area:<path>`,
`factor:<area>::<path>`, `requirement:<area>::<name>`, `rating:<id>`) and a
strict name grammar (no `/`, `:`, `.`, or spaces in names) whose entire purpose
is to make that string parse back to the composite identity unambiguously. The
object is just the *parsed* form of the string, persisted in full.

So the identity fields can collapse to single canonical-reference strings with no
loss of information, giving every payload one uniform identity shape. This is the
direction [`SPECIFICATION.md`](../SPECIFICATION.md) already points: it requires
durable machine-readable artifacts to carry *qualified* references and forbids
only the *unqualified* (prefixless) forms there.

Two things make this a deliberate reversal rather than a free refactor, and the
spec must account for both:

- It reverses an explicit prior decision. Cases
  [0058](archive/0058-model-reference-identifiers.md) and
  [0059](archive/0059-unqualified-model-references.md) made it a stated non-goal
  to replace the structured `areaPath`/`factorPath` arrays with string
  references in durable machine artifacts, and
  [`json-conventions.md`](../specs/evaluation/records/json-conventions.md)
  currently *mandates* the structured shapes and *forbids* the string forms in
  persisted routine JSON. This case rewrites that rule.
- The string form is lossless only if the rendering is unambiguous. The one gap
  is the root-Area token: `area:root` (and the `root` segment in
  `factor:root::…` / `requirement:root::…`) renders the empty root path, but the
  name grammar currently also admits a child Area literally named `root`, which
  would render to the same string. `root` is not reserved today. Collapsing to
  strings requires reserving `root` as a forbidden Area name (the only
  reservation needed — `/` and `:` are already excluded from names).

## Scope

Covered:

- Persist `areaId`, `factorId`, and `requirementId` — and their repeated forms
  (`areaIds`, `factorIds`, `localRequirementIds`, `rootFactorIds`,
  `childAreaIds`) and the secondary Factor lists (`factors`, subject
  `factorIds`) — as single canonical qualified model-reference strings in all
  Evaluation routine and report JSON, dropping the nested `declaringAreaId`,
  `factorPath`, and `requirementName` fields.
- Keep the `*Id` field names; do not rename to `*Ref`.
- Reserve `root` as a forbidden Area name across lint, the model schema, and
  [`SPECIFICATION.md`](../SPECIFICATION.md), so the string form is provably
  lossless.
- Bump the Evaluation data `SchemaVersion` for the payload-shape change; reject
  the old structured shapes rather than migrating them.
- Reconcile the durable Evaluation record specs (the JSON-conventions identity
  rule, payload-kinds, report-tree ordering) and regenerate the committed
  `evaluation-data.schema.json`.

Deferred / non-goals:

- No change to the *human-facing* rendered reports (Markdown tables, titles,
  ratings, trails, links) beyond what re-deriving an identity string from data
  requires; report values shown to people are titles and labels, not these IDs.
- No migration of existing completed evaluation runs under `.quality/`; they are
  regenerated on the next run, not rewritten.
- No new query language or selector surface; this only changes how identities are
  persisted, reusing the existing reference parsers.
- No change to the QUALITY.md frontmatter format or to how authors write Factor
  references inside a `QUALITY.md` file.

## Affected artifacts

Durable **spec** changes are itemized in the functional spec's
[Durable spec changes](0120-string-model-identity-fields/spec.md#durable-spec-changes)
section. The index below is the full skimmable list to reconcile before
In-Review.

### Code

- [x] `internal/evaluation/data_contract.go` - the source of truth for the
      shapes: `schemaForField` emits the structured array/object schemas for
      `dataAreaID`/`dataFactorID`/`dataRequirementID` (and `walkModelReferences`
      / the contract validator special-case the structured keys). Emit and
      validate canonical-string shapes instead; collapse the repeated and
      secondary-Factor list fields to string arrays.
- [x] `internal/evaluation/data.go` - the `areaIDFrom` / `factorIDFrom` /
      `requirementIDFrom` parsers and the `example*` builders decode/emit the
      structured shapes; switch them to render/parse the canonical reference
      strings (reusing the existing `parseAreaRef` / `parseFactorRef` /
      `parseRequirementRef` and the `types.go` reference encoders). `SchemaVersion`
      lives in `internal/evaluation/types.go` and bumps here.
- [x] `internal/evaluation/report_tree.go` - `factorIDJSON` / `requirementIDJSON`
      and the report-ref emission write the structured shapes into report
      outputs; emit canonical strings. Confirm no report consumer reads the split
      `declaringAreaId` / `factorPath` components directly (re-derive by parsing
      the string where one does).
- [x] `internal/evaluation/types.go` - `SchemaVersion` bump.
- [x] `internal/lint/` (`model.go`, `result.go`, `rules.go`) - add a diagnostic
      reserving `root` as a forbidden Area name.
- [x] `internal/schema/` - express the `root` Area-name reservation where the
      structural schema can (mirrors the existing strict-name pattern).
- [x] Tests - `internal/evaluation/evaluation_test.go`,
      `internal/cli/evaluation_test.go`, `internal/lint/rules_test.go`: update
      the structured-shape payload assertions to the new string shapes and add
      the `root` reservation case.

### Generated schema

- [x] `internal/evaluation/evaluation-data.schema.json` - committed generated
      artifact; regenerate after `data_contract.go` changes (the staleness guard
      in `evaluation_test.go` enforces this).
- [x] `quality.schema.json` - reviewed; add the `root` Area-name reservation if
      the model schema can express it (it carries the strict name grammar). No
      identity-field shapes live here.

### Format spec

- [x] [`SPECIFICATION.md`](../SPECIFICATION.md) - reserve `root` as a forbidden
      Area name in the name-grammar section; no other change (it already defines
      the qualified references this change persists).

### Durable specs

- [x] `specs/evaluation/records/json-conventions.md` - rewrite "Identity And
      References": persisted identity fields are canonical qualified
      model-reference strings, not structured objects/arrays; the prior
      "structured identity, never rendered refs" rule is replaced.
- [x] `specs/evaluation/records/payload-kinds.md` - align the per-kind identity
      field descriptions with the string shape.
- [x] `specs/evaluation/reports/report-tree.md` - align the ordering rules that
      name "declaring Area ID and structural Factor path / Requirement name" with
      the string shape.
- [x] `specs/cli/lint-rules.md` and `specs/cli/lint.md` - add the reserved-`root`
      Area-name diagnostic.

### Durable docs / bundled skill

- [x] `skills/quality/SKILL.md` and `specs/skills/quality-skill/quality-skill.md`
      - update Evaluation artifact identity wording to canonical qualified
      model-reference strings and align `/quality` compatibility metadata with
      the `0.16` CLI line.

### Suggested new durable specs

- The Evaluation data contract (`data_contract.go` + the generated
  `evaluation-data.schema.json`) is the source of truth for these shapes but has
  no 1:1 durable artifact spec. Worth considering an
  `evaluation-data-schema-json.md` artifact spec so the persisted identity
  contract has a durable home rather than living only in `json-conventions.md`
  prose and generated JSON. Suggesting only; not a precondition to land.

## Status

`Done`. Implemented, verified, and archived with the Evaluation data contract, report output, lint/schema validation, generated schemas, focused tests, and durable specs in sync.
