---
type: Functional Specification
title: Self-describing evaluation record input - functional spec
description: Requirements for payload-documenting help, a no-persist dry-run, and aggregated key-named validation on the qualitymd evaluation write commands, plus the skill surfaces that reach them.
tags: [cli, evaluation, skill, ergonomics]
timestamp: 2026-06-22T00:00:00Z
---

# Self-describing evaluation record input - functional spec

Companion to the
[Self-describing evaluation record input](../0055-evaluation-input-ergonomics.md)
change case. This spec states _what_ the `qualitymd evaluation` record-writing
surface must require so that an author-supplied payload is discoverable and
validatable from inside the tool. It defers the on-disk record contract to
[`specs/evaluation-records/`](../../../specs/evaluation-records/index.md), the
cross-cutting CLI contract to [`specs/cli.md`](../../../specs/cli.md), the design
principles to
[`docs/guides/cli-design.md` → Structured input](../../../docs/guides/cli-design.md#structured-input),
and evaluation semantics to [`SPECIFICATION.md`](../../../SPECIFICATION.md). It
changes how record input is surfaced and validated; it does not change the record
layout, the rating vocabulary, or the judgment content of any record.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

The `qualitymd` record-writing commands consume an author-supplied structured
payload (JSON read from `--file`/stdin) and surface nothing about its shape. A
field evaluation run (`0001-quality-eval`, against 0.8.0) spent ~25–30 failed or
probe invocations reverse-engineering record schemas from rejection messages —
and because the only schema oracle was "submit and read the error," every probe
persisted a real numbered record that then had to be deleted by hand. The
recommendation schema alone took ~12 probes. None of this touched the evaluation
_judgment_, which went cleanly; all of it was the input layer.

The durable record contract is already well specified under
[`specs/evaluation-records/`](../../../specs/evaluation-records/index.md). The
failure is that the CLI never exposes it and the skill cannot reach it: help
lists flags only, there is no validate-without-persist path, validation reveals
one problem at a time, and the skill's record source-of-truth citation does not
ship in the published bundle. This spec makes the record surface self-describing,
turning a write-and-fail discovery loop into help-and-dry-run.

## Scope

Covered: payload-documenting help, a no-persist dry-run, aggregated
key-named validation, drift-proofing of the surfaced examples, and the skill
surfaces (`SKILL.md`, quick reference, evaluate mode) plus a packaging guard that
reach or protect them.

Deferred:

- **Report-rendering defects** — the `report-summary.md` Scope field lifted from
  `## Rigor`, and recommendation link text Title-Cased from the slug instead of
  the record `title`. Real, but report-build cosmetics, not input ergonomics;
  a separate follow-up case.
- **`evaluation <kind> remove`** — the dry-run removes the need to write probe
  records, so in-tool cleanup is not built here.

Non-goals:

- No standalone `evaluation <kind> schema` command. If machine-readable JSON
  Schema is wanted, it is a flag on the existing command, not a new noun.
- No change to evaluation semantics, the rating vocabulary, the on-disk record
  layout, or the judgment content of any record.

## Requirements

### Payload-documenting help

- Each evaluation command that reads an author-supplied payload —
  `evaluation assessment add`, `evaluation analysis set`,
  `evaluation recommendation add` — **MUST** document its payload contract in
  `--help`: every field, its type, whether it is required, and the allowed values
  for any enum (e.g. finding `category`/`severity`, `criterionSource`).
  _Rationale: the 0.8.0 help listed flags only, forcing schema discovery by
  rejection; an agent cannot infer a payload's shape from context._
- That help **MUST** include at least one complete, valid example payload — every
  required field present, real enum values — not a fragment.
  _Rationale: the recommendation record took ~12 probes precisely because no
  worked example existed; required fields surfaced one at a time._
- The `plan.md` coverage contract read by `evaluation status` **SHOULD** be
  discoverable the same way (documented shape with an example, or a seeded
  commented stub), since it is the same write-and-fail trap on a different
  payload. _Rationale: the run hit two failed `status` round-trips guessing
  `assessmentResults`/`analyses` keys and object-vs-string entries._

### Validate without persisting

- Each evaluation write command **MUST** accept `-n/--dry-run`, which validates
  the payload fully and reports what would happen — what records would be written
  and where — without creating, numbering, serializing, or otherwise persisting
  any record or mutating the run folder.
  _Rationale: the only schema oracle was a real write, so every probe left a
  numbered junk record that had to be `rm`'d — uncomfortably close to the skill's
  "never hand-create record files" rule._
- Under `--dry-run`, exit status **MUST** reflect validity (zero when the payload
  would be accepted, the usage-error category when it would be rejected), and
  `--json` **MUST** emit a receipt describing the validation outcome, so the
  dry-run is usable unattended.

### Aggregated, caller-vocabulary validation

- Payload validation **MUST** report every problem it can detect in a single
  pass, rather than returning on the first failure.
  _Rationale: sequential required-field checks forced fix-one-resubmit-discover-
  the-next round-trips across assessment, analysis, and recommendation writes._
- Each validation error **MUST** name the offending field by the JSON key the
  author wrote, and **SHOULD** state the expected type and, for an enum, the
  allowed values. It **MUST NOT** echo an internal or language-level field name
  the caller never typed. _Rationale: 0.8.0 reported TitleCased struct names like
  `Donecriterion is required`, forcing the reader to re-camelCase and guess word
  boundaries. Current HEAD already emits the JSON key for the explicit checks;
  this requirement makes that the contract and extends it to the unknown-field
  and type-mismatch paths._

### Drift-proofing the surfaced contract

- The example payloads surfaced in command help and in the skill's quick
  reference **MUST** be generated from, or golden-tested against, the record
  structs or the
  [`specs/evaluation-records/`](../../../specs/evaluation-records/index.md)
  contract, such that a divergence between a surfaced example and the accepted
  schema fails a test. _Rationale: the shipped `cli-quick-reference.md` had all
  three example payloads wrong for 0.8.0 because they were hand-maintained;
  "keep in sync" is not a mechanism._

### Skill reach and packaging

- `SKILL.md` **MUST NOT** cite a record source-of-truth path that is absent from
  the published skill bundle; it **MUST** instead direct record writing to the
  in-tool authority (command help and `--dry-run`).
  _Rationale: `SKILL.md` named `../../specs/evaluation-records.md` as the source
  of truth, but only `SKILL.md`, `modes/`, `guides/`, `resources/` ship — the
  single most useful artifact was unreachable at the point of use._
- The published skill bundle **MUST** be checked so that every relative link in
  its shipped files resolves within the bundle, failing CI when one does not.
  _Rationale: a dead in-skill citation should fail a build, not a field run._
- The skill's record-writing procedure (`modes/evaluate.md`) and quick reference
  (`resources/cli-quick-reference.md`) **MUST** reflect the discover-via-help /
  validate-via-dry-run flow and carry correct payloads.

### Documenting the implicit analysis contract

- The analysis record specs **MUST** state that `evaluation analysis set` writes
  the whole area set atomically and how `childAnalysisRecords` links resolve
  (by path, to records the same `set` writes).
  _Rationale: the run inferred — correctly but unsupported by any doc — that a
  single `analysis set` writes all areas atomically, which is the only reason a
  `childAnalysisRecords` path to a not-yet-written record validated._

## Durable spec changes

### To add

- None. The record contract already lives in
  [`specs/evaluation-records/`](../../../specs/evaluation-records/index.md); this
  case surfaces and validates it rather than introducing new durable specs.

### To modify

- `specs/cli.md` — add the cross-cutting structured-input contract: help
  documents the payload, `-n/--dry-run` validates without persisting, validation
  is aggregated and names the caller's JSON keys (per the payload-help, dry-run,
  and validation requirements).
- `specs/cli/evaluation-assessment.md`, `specs/cli/evaluation-analysis.md`,
  `specs/cli/evaluation-recommendation.md` — per-command `--dry-run` and
  help-contract requirements (per the payload-help and dry-run requirements).
- `specs/cli/evaluation-create.md` — seeded `plan.md` coverage-shape discovery
  (per the coverage discoverability requirement).
- `specs/cli/evaluation-status.md` — apply the aggregated/key-named rule to the
  `plan.md` coverage validation and make the coverage shape discoverable (per the
  coverage and validation requirements).
- `specs/evaluation-records/analysis-record.md` and `specs/evaluation-records.md`
  — state `analysis set` whole-set atomicity and `childAnalysisRecords`
  link-resolution, and designate the record specs/structs as the generation
  source for surfaced examples (per the analysis-contract and drift-proofing
  requirements).
- `specs/skills/quality-skill/quality-skill.md` — update the record-writing
  procedure to discover the contract via help/dry-run rather than the unshipped
  spec (per the skill-reach requirements).

### To rename

- None.

### To delete

- None.
