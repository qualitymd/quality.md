---
type: Functional Specification
title: Bulk data set — functional spec
description: Requirements for `qualitymd evaluation data set`'s array-only batch stdin contract, written all-or-nothing in one invocation.
tags: [cli, evaluation, data, agent]
timestamp: 2026-06-26T00:00:00Z
---

# Bulk data set — functional spec

Companion to the [Bulk data set](../0126-bulk-data-set.md) change case. This
spec states _what_ the new `data set` stdin contract must do; a later design doc
covers _how_. The payload kinds, their fields, and the canonical reference
grammar are defined by [`SPECIFICATION.md`](../../../SPECIFICATION.md) (normative)
and the Evaluation data records under
[`specs/evaluation/`](../../../specs/evaluation/index.md); the invocation-wide CLI
contract — output posture, exit codes, agent accessibility — is defined by
[`specs/cli.md`](../../../specs/cli.md) (normative). The current `data set` behavior
this spec amends lives in
[`specs/cli/evaluation-data.md`](../../../specs/cli/evaluation-data.md).

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

`data set` reads exactly one JSON object per invocation, so the evaluate
workflow persists a model's tens of routine payloads one invocation at a time —
~115 agent Bash round-trips for the cited acquire-roi-next run. A batch contract
turns that loop into a single write. See the change case
[Motivation](../0126-bulk-data-set.md#motivation) for the originating evidence.

The contract is **array-only**: the array _replaces_ the single-object form
rather than adding a second path, so the tool keeps one transport contract to
validate, document, and test. A one-payload write is a one-element array. This
is a deliberate clean break from the v0 single-object rule, appropriate while
the format is early alpha.

The batch is **atomic**: one bad element rejects the whole batch and nothing is
written. An agent that submitted ~115 payloads gets every failure at once and
fixes them in one pass, and a run's `data/` directory is never left
half-populated by a partial batch.

## Scope

Covered: the stdin contract, validation aggregation, atomicity, duplicate-path
rejection, the write receipt, and dry-run for `qualitymd evaluation data set`.

Out of scope (unchanged by this case):

- Per-payload validation: unknown/misspelled fields, wrong types, out-of-range
  enums, missing required fields, and model-snapshot reference resolution all
  keep their current behavior, now applied per element.
- Kind routing and canonical `data/**` path derivation per payload.
- The persisted payload schema and `evaluation-data.schema.json` — the array is
  a transport envelope, not a stored shape.
- The other `data` verbs (`list`, `get`, `kinds`, `example`, `schema`,
  `verify`), `evaluation create`, and stdin source (still stdin only; no
  `--file`).

Deferred / non-goals (each absence is deliberate):

- No `--partial` / continue-on-error mode — v0 is strictly all-or-nothing.
- No NDJSON or streaming input form — a single JSON array only.
- No new exit-code vocabulary — failure reuse the invocation-wide codes in
  [`specs/cli.md`](../../../specs/cli.md).

## Requirements

### Input format

- **IN1 — Array only.** `data set` **MUST** read its stdin as a single JSON array
  whose elements are Evaluation payload objects, and **MUST** reject input that is
  not a JSON array (including a bare JSON object) with a usage error that names
  the expected array form.

  > > Rationale: the array replaces the single-object form rather than adding a
  > > second path, so there is one transport contract. A one-payload write is a
  > > one-element array. — 0126

- **IN2 — Non-empty.** `data set` **MUST** reject an empty array (`[]`) as a usage
  error; a set with no payloads is a mistake, not a successful no-op.

- **IN3 — Heterogeneous kinds.** A batch **MAY** mix payload kinds (e.g.
  `*EvaluationFrame`, `*AnalysisResult`, `*RatingResult` together); `data set`
  **MUST** route each element by its own `kind`, exactly as for a single payload.

### Validation and atomicity

- **AT1 — Validate the whole batch first.** `data set` **MUST** apply the full
  per-element contract — structural, typed, required-field, enum, and
  model-snapshot reference resolution — to **every** element before writing any
  element.

- **AT2 — All-or-nothing.** If **any** element fails validation, `data set`
  **MUST** write **no** element, **MUST** exit non-zero, and **MUST NOT** leave
  the run's `data/` directory partially modified by the batch.

  > > Rationale: a partial batch leaves a run in a state neither the agent nor a
  > > later reader can reason about; atomicity makes a failed `data set` a no-op to
  > > retry after fixing. — 0126

- **AT3 — Aggregated, indexed diagnostics.** When a batch fails validation,
  `data set` **MUST** report each invalid element identified by its position in
  the array (and the offending field where one is known), and **SHOULD** report
  every invalid element, not only the first, so the whole batch can be corrected
  in one pass.

  > > Rationale: stopping at the first error would re-impose the round-trip cost
  > > the batch contract removes — the agent would resubmit once per error. — 0126

- **AT4 — Snapshot fail-closed.** As today, if the run's `model-snapshot.md` is
  missing or unparseable, `data set` **MUST** fail closed for the whole batch
  rather than accept unbound model references. The snapshot **SHOULD** be loaded
  once per invocation, not once per element.

### Idempotency and collisions

- **ID1 — Per-path idempotent overwrite.** `data set` **MUST** overwrite each
  element's derived `data/**` path by default, so re-submitting a batch of the
  same routine outputs is idempotent.

- **ID2 — Reject intra-batch path collisions.** If two elements in one batch
  derive the **same** canonical `data/**` path, `data set` **MUST** reject the
  batch as a usage error naming the colliding elements, rather than silently
  letting a later element overwrite an earlier one.

  > > Rationale: two payloads writing one path in a single batch is almost always a
  > > duplicate-authoring mistake; surfacing it beats a silent last-writer-wins. —
  > > 0126

### Receipt and dry-run

- **RC1 — Batch receipt.** Under `--json`, `data set` **MUST** emit a write
  receipt for the batch — a count of payloads written and the canonical path of
  each, in input order — not the stored JSON artifacts. The receipt **MUST**
  describe a successful batch as a whole.

- **RC2 — Dry-run the whole batch.** Under `--dry-run`, `data set` **MUST**
  perform the same whole-batch validation (AT1–AT4) and collision check (ID2) and
  **MUST NOT** persist any element. It **MUST** report what the batch would write
  (the same set RC1 would confirm) without modifying the run.

### Boundary (unchanged)

- **BD1 — `EvaluationOutputResult` rejected anywhere.** `data set` **MUST** reject
  a batch containing an `EvaluationOutputResult` element, at any position; that
  payload stays CLI-owned and generated by `evaluation report build`.

- **BD2 — No `--file`.** `data set` **MUST NOT** accept a `--file` flag; the batch
  array arrives on stdin.

### Skill integration (evaluate workflow)

The round-trip cost this case removes only materializes when the evaluate
workflow stops persisting one payload per invocation (see
[Motivation](#background--motivation)). Shipping the array contract without
wiring it into that workflow would leave the originating pain unsolved, so
adoption is a required outcome of this case, not an optional follow-up. Rollout
MAY land the CLI change and the skill edit as separate commits; the contract
below is what the change must deliver.

- **SK1 — Batch the writes.** The evaluate workflow **MUST** persist a scope's
  routine payloads with a single batched `data set` invocation per scope —
  authoring the array from the payloads it has assembled — rather than one
  invocation per Requirement, Factor, or Area.

  > > Rationale: the per-element loop is the ~115-invocation cost the array
  > > contract exists to remove; the win is realized in the workflow, not the CLI.
  > > — 0126

- **SK2 — One whole-batch dry-run.** Where the workflow validates before
  persisting, it **MUST** dry-run the assembled batch once (one `data set
--dry-run`) rather than dry-running each payload, relying on the aggregated,
  indexed diagnostics (AT3) to correct the whole batch in one pass.

  > > Rationale: per-payload dry-runs reintroduce the round-trips the batch
  > > contract removes; whole-batch validation matches whole-batch writes. — 0126

## Acceptance criteria

- [x] A one-element array writes exactly the file the old single-object form
      wrote, to the same canonical path; a bare JSON object on stdin is now a
      usage error naming the expected array form.
- [x] An empty array (`[]`) exits non-zero as a usage error.
- [x] A batch mixing kinds (e.g. a frame, an analysis result, and a rating result
      for one scope) writes each to its own canonical path in one invocation.
- [x] A batch with one invalid element writes **nothing**, exits non-zero, and the
      run's `data/` directory is byte-identical to before the invocation.
- [x] The failure report identifies each invalid element by array index and names
      the offending field where known, reporting all invalid elements rather than
      stopping at the first.
- [x] Two elements deriving the same `data/**` path reject the batch as a usage
      error naming the collision; nothing is written.
- [x] `--json` on a successful batch emits a receipt with the written count and
      each canonical path in input order, not the stored artifacts.
- [x] `--dry-run` on a valid batch persists nothing and reports the would-write
      set; `--dry-run` on an invalid batch reports the same aggregated diagnostics
      as a real write.
- [x] A batch containing an `EvaluationOutputResult` at any position is rejected.
- [x] The evaluate workflow persists a scope's routine payloads with a single
      batched `data set` per scope rather than one invocation per element (SK1),
      and validates the assembled batch with one `--dry-run` before the write
      rather than per payload (SK2).

## Durable spec changes

Durable **specs** this change rewrites — the
[`specs/`](../../../specs/index.md) bundle and
[`SPECIFICATION.md`](../../../SPECIFICATION.md). Each subsection is required.

### To add

None.

### To modify

- `specs/cli/evaluation-data.md` — replace the single-object stdin rules with the
  array contract: rewrite "`data set` **MUST** read exactly one structured JSON
  payload from stdin …" to read a JSON array (IN1–IN3); replace "`data set`
  **MUST NOT** accept batch payloads in v0." with the atomicity, aggregated-error,
  intra-batch-collision, and batch-receipt rules (AT1–AT4, ID1–ID2, RC1–RC2);
  update the command listing to the `payloads.json` array form. Carry the _why_
  — round-trip cost and the atomic-no-op retry property — into its Background /
  Motivation.
- `specs/cli.md` — update the `data set` signature in the command listing and
  review the "Structured input" note for any single-payload phrasing.
- `specs/cli/evaluation-create.md` — update the `evaluation create` next-action
  command to show the array input filename (`payloads.json`) rather than the old
  single-payload example.
- `specs/evaluation/orchestration.md` — update orchestrator persistence guidance
  to assemble accepted payloads into one JSON array and persist the batch through
  `data set`.
- `specs/skills/quality-skill/evaluation.md` — update the high-level evaluation
  workflow overview to describe batched routine payload persistence.
- `specs/skills/quality-skill/workflows/evaluate.md` — add the normative
  requirement that routine payloads are persisted in batches via a single
  `data set`, not one invocation per element.

### To rename

None.

### To delete

None.

## Open questions

None. The [design doc](design.md) settles write staging/rollback and the batch
receipt schema.
