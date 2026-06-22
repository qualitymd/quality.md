---
type: Change Case
title: Self-describing evaluation record input
description: Make qualitymd evaluation record writing discoverable and validatable from inside the tool — payload-documenting help, a no-persist dry-run, aggregated key-named validation errors — and repair the skill surfaces that drifted from the binary.
status: Done
tags: [cli, evaluation, skill, ergonomics]
timestamp: 2026-06-22T00:00:00Z
---

# Self-describing evaluation record input

A **Change Case** making the `qualitymd evaluation` record-writing surface
self-describing, and repairing the `/quality` skill surfaces that drifted from
it. The detail lives in its
[functional spec](0055-evaluation-input-ergonomics/spec.md).

## Motivation

A field evaluation run against `qualitymd` 0.8.0 (run `0001-quality-eval`,
captured in an external project's `observations.md`) found that the evaluation
*judgment* went smoothly — model valid, evidence clean, ratings bound to
evidence — but the **record-writing layer cost ~25–30 failed or probe CLI
invocations**, every one of them avoidable. The root cause is that an
author-supplied structured payload (the assessment / analysis / recommendation
JSON read from `--file`/stdin) is the most agent-hostile surface in the binary,
and the tool does nothing to make it reachable:

- **The shape is invisible.** `evaluation <kind> add --help` lists flags only —
  no fields, no required-ness, no enums, no example. The recommendation schema in
  particular took ~12 probes to reverse-engineer (`remediationOptions` is an
  array of *strings*, not objects; `doneCriterion` and `evidenceLocators` are
  required but each surfaced only after everything else was already correct).
- **There is no way to check a payload without writing one.** `qualitymd schema`
  emits only the QUALITY.md frontmatter schema; there is no record schema and no
  dry-run. So the only schema oracle is "submit and read the rejection," and
  **every probe persists a real numbered record** — which then had to be `rm`'d
  by hand, brushing against the skill's own "never hand-create record files"
  rule.
- **Validation forces serial round-trips.** Each bad payload reveals exactly one
  problem; the caller fixes it, resubmits, and discovers the next. (Current HEAD
  already echoes camelCase keys like `evidenceLocators is required`, improving on
  0.8.0's TitleCased struct names — but the one-at-a-time behavior remains.)
- **The skill's own pointers are broken or stale.** `SKILL.md` cites
  `../../specs/evaluation-records.md` as the record source of truth, but that
  path does not resolve in the *published* skill (only `SKILL.md`, `modes/`,
  `guides/`, `resources/` ship). The one in-bundle payload reference,
  `resources/cli-quick-reference.md`, had all three example payloads wrong for
  0.8.0.

The durable record contract itself is *not* the problem — it is well specified
in [`specs/evaluation-records/`](../../specs/evaluation-records/index.md). The
problem is purely that the CLI does not surface that contract, and the skill
cannot reach it. This case closes both gaps, and it is the operational
counterpart to the new
[Structured input](../../docs/guides/cli-design.md#structured-input) section of
the CLI design guide — the principles land here as binding requirements and a
working surface.

## Scope

Covered:

- **Payload-documenting help.** Each evaluation write command
  (`assessment add`, `analysis set`, `recommendation add`, and the `plan.md`
  coverage read by `status`) documents its payload contract — every field, type,
  required-ness, and enum values — and shows at least one complete valid example,
  discoverable from `--help` without reading source.
- **Validate without persisting.** A `-n/--dry-run` path on the write commands
  validates a payload fully and reports what would happen without creating,
  numbering, or serializing any record.
- **Aggregated, caller-vocabulary errors.** Validation reports the whole set of
  problems in one pass and names each field by the JSON/YAML key the author
  wrote, with expected type and allowed values — never an internal symbol.
- **Drift-proofing.** The payload examples surfaced in help (and in the skill's
  quick reference) are generated from, or golden-tested against, the record
  structs / spec, so they cannot silently drift again.
- **Skill surfacing.** `SKILL.md` stops citing the unshipped spec and points at
  the in-tool help / dry-run as the authority; `cli-quick-reference.md` payloads
  are corrected and pinned to fixtures; the evaluate mode's record-writing
  procedure uses help/dry-run for discovery.
- **Packaging guard.** A check that every relative link in the *published* skill
  bundle resolves, so a dead in-skill citation fails CI rather than a field run.
- **Document the implicit `analysis set` contract** (#8): that it writes the
  whole area set atomically and how `childAnalysisRecords` links resolve.

Deferred / non-goals:

- **Report-rendering defects** (the `report-summary.md` Scope field sourced from
  `## Rigor`, and recommendation link text lossily Title-Cased from the slug
  instead of using the record's own `title`). These are real but separable
  report-build cosmetics, not input ergonomics; track them in a small follow-up
  case.
- **`evaluation <kind> remove`.** The dry-run removes the need to write probe
  records, so in-tool cleanup is not pursued here; add it only if a real
  consumer needs it.
- A standalone `evaluation <kind> schema` *command*. If machine-readable JSON
  Schema is wanted, prefer a flag on the existing command over a new noun; the
  functional spec decides.
- No change to evaluation semantics, the rating vocabulary, the on-disk record
  layout, or the judgment content of any record. This is a discovery,
  validation, and surfacing change only.

## Affected artifacts

Durable **spec** changes are itemized in the functional spec's
[Durable spec changes](0055-evaluation-input-ergonomics/spec.md#durable-spec-changes)
section. The index below is the full skimmable list, reconciled before
In-Review.

Code (edited only once the case is `In-Progress`):

- [x] [`internal/cli/evaluation.go`](../../internal/cli/evaluation.go) — add the
      `-n/--dry-run` flag to the write commands; populate `Long`/`Example` for
      each with the payload contract and a complete valid example.
- [x] [`internal/evaluation/write.go`](../../internal/evaluation/write.go) —
      aggregate validation (collect all problems, not the first), echo JSON keys,
      and add the no-persist dry-run path.
- [x] [`internal/evaluation/input_contract.go`](../../internal/evaluation/input_contract.go)
      and [`internal/evaluation/types.go`](../../internal/evaluation/types.go) —
      the write options, field metadata, payload help, and canonical examples.
- [x] [`internal/evaluation/create.go`](../../internal/evaluation/create.go) and
      [`internal/evaluation/planned_coverage.go`](../../internal/evaluation/planned_coverage.go)
      — seed the `plan.md` coverage-shape example and aggregate planned-coverage
      validation problems.
- [x] [`internal/evaluation/atomic.go`](../../internal/evaluation/atomic.go) —
      remove the obsolete write helper after numbering moved into the write
      planning stage.
- [x] [`internal/cli/evaluation_test.go`](../../internal/cli/evaluation_test.go)
      and `internal/evaluation/*_test.go` — cover dry-run, aggregated errors, and
      help/example-vs-struct golden assertions.

Specs:

- [x] [`specs/cli.md`](../../specs/cli.md) — add the cross-cutting structured-input
      contract (help documents the payload, `-n/--dry-run` validate-without-persist,
      aggregated key-named validation errors) mirroring the design guide's new
      section.
- [x] [`specs/cli/evaluation-assessment.md`](../../specs/cli/evaluation-assessment.md),
      [`specs/cli/evaluation-create.md`](../../specs/cli/evaluation-create.md),
      [`specs/cli/evaluation-analysis.md`](../../specs/cli/evaluation-analysis.md),
      [`specs/cli/evaluation-recommendation.md`](../../specs/cli/evaluation-recommendation.md)
      — per-command dry-run, help-contract, and `plan.md` coverage-stub
      requirements.
- [x] [`specs/cli/evaluation-status.md`](../../specs/cli/evaluation-status.md) —
      apply the aggregated/key-named rule to the `plan.md` coverage validation.
- [x] [`specs/evaluation-records/analysis-record.md`](../../specs/evaluation-records/analysis-record.md)
      and [`specs/evaluation-records.md`](../../specs/evaluation-records.md) —
      state the `analysis set` whole-set atomicity and `childAnalysisRecords`
      link-resolution contract; designate the record specs as the generation
      source for help/quick-ref examples.
- [x] [`specs/skills/quality-skill/quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md)
      — update the record-writing procedure to discover the contract via
      help/dry-run rather than the unshipped spec.

Docs:

- [x] [`docs/guides/cli-design.md`](../../docs/guides/cli-design.md) — added the
      **Structured input** section plus the Help, Documentation, and Errors
      additions (landed alongside this case as the reasoning behind it).
- [x] [`skills/quality/SKILL.md`](../../skills/quality/SKILL.md) — drop the dead
      `../../specs/evaluation-records.md` citation; point at the in-tool
      help/dry-run authority.
- [x] [`skills/quality/resources/cli-quick-reference.md`](../../skills/quality/resources/cli-quick-reference.md)
      — correct the stale assessment/analysis/recommendation payloads and pin
      them to fixtures.
- [x] [`skills/quality/modes/evaluate.md`](../../skills/quality/modes/evaluate.md)
      — discover record shapes via help/dry-run in the procedure.
- [x] [`scripts/check-npm-package.mjs`](../../scripts/check-npm-package.mjs) — add
      the published-bundle relative-link resolution guard.
- [x] [`CHANGELOG.md`](../../CHANGELOG.md) — add the 0055 entry.

`SPECIFICATION.md` is **not** affected: this case changes how the CLI surfaces and
validates record input and how the skill reaches it, not the QUALITY.md format or
evaluation semantics.

## Children

- [Functional spec](0055-evaluation-input-ergonomics/spec.md) — what the
  self-describing record surface must require.
- [Design doc](0055-evaluation-input-ergonomics/design.md) — how the Go code
  delivers it: the decode/validate/plan/commit split, the golden-tested canonical
  example, and the skill/packaging guards.

## Status

`Done`. Landed and archived after review. Verified with `mise run check`.
