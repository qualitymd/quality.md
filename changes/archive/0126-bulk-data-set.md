---
type: Change Case
title: Bulk data set
description: Replace `qualitymd evaluation data set`'s single-object stdin contract with an array-only batch contract — one invocation writes many payloads, validated all-or-nothing — collapsing the evaluate workflow's per-element write loop from tens of invocations to one.
status: Done
tags: [cli, evaluation, data, agent]
timestamp: 2026-06-26T00:00:00Z
---

# Bulk data set

This parent concept captures the why and status; the detail lives in its
children:

- [Functional spec](0126-bulk-data-set/spec.md) - what the change must do.
- [Design doc](0126-bulk-data-set/design.md) - how the batch validation,
  staged writes, rollback, and receipt shape work.

## Motivation

`evaluation data set` reads exactly one JSON object per invocation
([`specs/cli/evaluation-data.md:34`](../../specs/cli/evaluation-data.md), enforced
by the `decodeDataPayload` "payload must contain one JSON object" guard in
`internal/evaluation/data.go`), and v0 explicitly forbids batch input
([`evaluation-data.md:59`](../../specs/cli/evaluation-data.md): "`data set`
**MUST NOT** accept batch payloads in v0").

So the `/quality` evaluate workflow loops: "**for each** Requirement / Factor /
Area … persist with `data set`" (`skills/quality/workflows/evaluate.md`). The
acquire-roi-next evaluation cited in landed
[0125](0125-model-query-commands.md)
produced **~115 separate `data set` invocations** — i.e. ~115 agent Bash
round-trips, each a model turn with its own latency and tokens.

The win here is **agent round-trips, not CLI internals.** A bulk contract
collapses the per-element loop into "assemble all payloads, flush once," and
turns the workflow's per-payload `--dry-run` validation into a single
whole-batch validate-then-write. The CLI's own parsing gets marginally more
complex (iterate + batch-failure semantics); the orchestration it drives gets
dramatically simpler.

This is the write-side companion to landed
[0125](0125-model-query-commands.md), which removes the _read_-side
friction of the same loop (query canonical IDs once instead of hand-deriving ~115
of them). Together they reshape the evaluate workflow's element loop: query the
in-scope IDs in one call, author all payloads, persist them in one call. 0125 is
now the baseline this case builds on; 0126 starts from its `model list --json`
snapshot-query workflow text rather than treating that edit as in flight.

## Scope

Covered: changing `evaluation data set`'s stdin contract from a single JSON
object to a **JSON array of payloads**, written **all-or-nothing** in one
invocation, and wiring the `/quality` evaluate workflow to author and persist
payloads in batches instead of one at a time.

Settled in design discussion (these are the spec's premises, not open
questions):

- **Exclusive bulk, not a second path.** The array contract _replaces_ the
  single-object contract; a one-payload write is a one-element array. The CLI is
  agent/skill-driven, so the ergonomic tax of `[{…}]` for a human one-off is
  negligible, and maintaining a single-object contract alongside an array
  contract is the "two paths" cost the case exists to avoid. Early-alpha
  clean-break rules (AGENTS.md) favor replacing over keeping both.
- **JSON array, not NDJSON.** Agents generate the whole structure at once; an
  array validates as a unit and matches how `model list --json` (landed in 0125)
  feeds IDs in.
- **All-or-nothing atomicity.** Validate the whole array; if any element fails,
  write nothing and report every failure by index. The agent fixes all errors in
  one pass and the run folder never ends up half-populated.
- **Intra-batch duplicate derived paths are an error,** not last-wins — almost
  always an agent mistake.

Deferred / non-goals:

- No `--partial` / continue-on-error mode — deferred until a real need; v0 is
  strictly atomic.
- No new input source — stdin only; the existing "`MUST NOT` accept a `--file`
  flag" rule stands.
- No change to per-payload validation, kind routing, canonical-path derivation,
  or the model-snapshot reference check — only the transport (one object → array)
  and the failure/receipt semantics change.
- No change to the persisted payload schema or to `evaluation-data.schema.json` —
  the array is a CLI transport envelope, not a stored shape.
- The CLI-owned `EvaluationOutputResult` stays rejected, anywhere in a batch.

## Affected artifacts

Durable **spec** changes are itemized in the functional spec's
[Durable spec changes](0126-bulk-data-set/spec.md#durable-spec-changes)
section. The index below is the full skimmable list to reconcile before
Done.

### Code

- [x] `internal/evaluation/data.go` - `decodeDataPayload` currently decodes one
      object and asserts EOF on a second; rework it to decode a JSON array of
      payload objects. `SetData` currently validates-and-writes one payload;
      split it into a validate-all pass (over every element, against one loaded
      `model-snapshot.md`) and a write-all pass that runs only if validation
      wholly succeeds, rejecting any two elements whose derived `data/**` paths
      collide.
- [x] `internal/cli/evaluation.go` - `newEvaluationDataSetCmd` and `readPayload`;
      the `--json` write receipt becomes a batch summary (count + the paths
      written), and `--dry-run` reports the whole batch's would-write set without
      persisting.
- [x] Tests - `internal/evaluation/` data tests (single-element array round-trip,
      multi-kind batch, one-bad-element rejects the whole batch with per-index
      reasons, duplicate-derived-path rejection, empty-array usage error,
      `EvaluationOutputResult` anywhere in a batch rejected) and any
      `internal/cli` test that exercises `data set` stdin and the receipt shape.

### Format spec

- [x] [`SPECIFICATION.md`](../../SPECIFICATION.md) - reviewed; no change. Bulk is a
      CLI transport envelope; the payload model and canonical reference grammar
      are unchanged.

### Durable specs

- [x] `specs/cli/evaluation-data.md` - invert the "single object" / "MUST NOT
      accept batch payloads in v0" rules to the array contract; state atomicity,
      duplicate-derived-path rejection, the batch receipt, and empty-array
      handling. (Itemized in the spec's Durable spec changes.)
- [x] `specs/cli.md` - review the "Structured input" note for single-payload
      phrasing and update it to cover batch envelopes.
- [x] `specs/cli/evaluation-create.md` - update the run-creation next action to
      show the array input filename.
- [x] `specs/evaluation/orchestration.md` - update orchestrator persistence
      guidance to assemble and persist a JSON-array batch.
- [x] `specs/evaluation/records/json-conventions.md` - reviewed for any
      single-payload-write phrasing; update only if it describes the `data set`
      transport.
- [x] `specs/skills/quality-skill/evaluation.md` - update the overview workflow
      to describe batched routine payload persistence.

### Durable docs / bundled skill

- [x] `skills/quality/workflows/evaluate.md` - collapse the per-Requirement /
      per-Factor / per-Area `data set` loops into batched writes (author all
      payloads for the scope, persist in one `data set`). Builds on the landed
      0125 snapshot-ID query text in the same file.
- [x] `skills/quality/resources/cli-workflow-conventions.md` - reviewed after
      [0127](0127-introspection-first-cli-reference.md) absorbs the shared
      resource edit. No 0126 edit is expected here because the stale
      single-payload `set <run> < payload.json` listing is removed wholesale with
      the former command listing.
- [x] `specs/skills/quality-skill/workflows/evaluate.md` - normative requirement
      that routine payloads are persisted in batches via a single `data set`, not
      one invocation per element.
- [x] `README.md` - reviewed; update only if it shows a single-payload `data set`
      example.

### Suggested new durable specs

- None. The batch envelope and failure semantics belong in the existing
  `specs/cli/evaluation-data.md`; no new durable spec is earned.

## Status

`Done`. Implemented array-only `evaluation data set`, batch validation and
receipt output, staged writes with rollback, durable spec/workflow updates, and
focused tests. Archived after `mise run check` passed.
