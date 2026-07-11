---
type: Change Case
title: Resolver-dispatched source selectors
description: Resolve every source selector through a per-kind resolver that feeds the bounded, hashed evidence bundle, so non-path selectors evaluate through the same audited contract instead of dead-ending as unavailable.
status: Done
---

# Resolver-dispatched source selectors

## Motivation

`SPECIFICATION.md` defines source as "a selector describing the material
evaluated by an area." A selector could be a path, a glob, a saved query, or a
prose description of a body of evidence ("all specs", "the deployed API",
"open tickets in the support queue").
[0196](0196-spec-faithful-model-reading.md) made the path/glob selector
kind spec-faithful and stopped silent empty evidence, and deferred the general
architecture to this case with the reasoning recorded in its
[considerations sketch](0196-spec-faithful-model-reading/considerations.md):
separate _resolution_ (gathering, per selector kind) from _judgment_ (rating
against a captured bundle), with the bounded, hashed evidence bundle as the
contract between them.

Today a non-path selector is stat-ed as a filesystem path, resolves to
nothing, and fails the affected work as `source_unavailable`. That failure is
loud — the 0196 fix — but it is also a dead end: a document whose selector the
format permits cannot be evaluated at all, and the failure message misdiagnoses
it as missing material rather than an unsupported selector kind.

The 2026-07-11 re-validation of 0196 confirmed the bundle contract is the
right seam. Every guarantee the runner gives — resumability, the input-hash
guard on pending harness work, evidence-bound re-judgment when material
changes — held through the bundle, independent of how the bundle was gathered.
A resolver that needs tools (a query, an agent) can therefore be
non-deterministic in _gathering_ while the run keeps reproducibility _of
record_: the captured bundle is still bounded, hashed, persisted, and
re-judgeable, and the source-as-data boundary is preserved.

The same re-validation surfaced a related design pressure, recorded here so
the design weighs it: a broad selector (the whole-workspace document default)
makes pending and completed work sensitive to any workspace write — one added
file invalidates a pending harness request and re-judges completed units on
resume. Correct, but expensive in agent-driven runs where the agent itself
writes files. Selector granularity and bundle scope are part of the resolution
design, not an afterthought.

## Scope

Covered:

- a per-kind resolution step between the effective source selector and the
  packaged bundle, with path/glob resolution unchanged;
- dispatching resolution of non-deterministic selector kinds to the invoking
  harness through the existing checkpoint transport, with the returned
  material validated and captured into the bundle before dependent judgment;
- a loud, correctly classified failure for a selector kind no resolver serves,
  distinct from material that is genuinely missing;
- run-artifact provenance: which resolver produced each bundle, from which
  selector, so audit and re-judgment read the same for deterministic and
  harness-gathered evidence;
- settling the format question this depends on — settled 2026-07-11, before
  Design: `SPECIFICATION.md` commits to non-filesystem selectors, and a
  selector stays a bare string with detected kind (glob metacharacters → glob;
  existing filesystem entry → path; otherwise prose); see the
  [functional spec](0197-resolver-dispatched-source-selectors/spec.md)'s
  settled questions.

Deferred / non-goals:

- **Out-of-tree source refs** (`../shared`) stay workspace-contained; 0196
  recorded that narrowing as intentional and this case does not revisit it.
- **Evaluator adapter fragility** (the `codex` 403 / `claude` schema-mismatch
  failures) is a separate concern from resolution; no case yet.
- **Default-selector granularity** — narrowing or partitioning the
  whole-workspace default to reduce re-judgment churn — is recorded as design
  pressure above but not required here.
- **Saved-query registries and live-system monitoring** as selector kinds of
  their own; this case establishes the seam, not a catalog of resolvers.

## Affected artifacts

Reconciled at implementation (In-Progress).

- **Code:** `internal/runner/source.go` (kind detection, resolver dispatch in
  `areaSourceBundle`, capture of resolver-returned material),
  `internal/runner/graph.go` (`resolveSource` work units gating judgment),
  `internal/runner/engine.go` (plan-time `selector_unsupported` check,
  per-unit guard, resolution acceptance), `internal/runner/requests.go` (the
  `resolveSource` request and schema), `internal/runner/harness.go` and
  `internal/runner/concurrent.go` (kind-aware reuse, harness runtime
  attribution, per-unit guard), `internal/runner/artifact.go` (the per-area
  `sources` provenance record; artifact schema version 6),
  `internal/runner/runner.go` (kind pinning at creation, pinned kinds on
  resume, receipt dispatch plan), `internal/runner/dryrun.go` (the preview
  dispatch plan), `internal/evaluator/evaluator.go` and
  `internal/evaluator/harness.go` (the `SourceResolution` capability, the
  `selector_unsupported` failure category). `internal/model/source.go` is
  unchanged — kind detection stays runner-owned per the
  [design](0197-resolver-dispatched-source-selectors/design.md). Tests
  alongside each (`internal/runner/source_test.go`,
  `internal/runner/runner_test.go`).
- **Durable specs:** modified `specs/evaluation/runner.md` (§Source packaging
  becomes detection + resolution + packaging; §Failure taxonomy),
  `specs/evaluation/protocol.md` and `specs/evaluation/orchestration.md`
  (the `resolveSource` move/unit and its dependency),
  `specs/evaluation/evaluator-contract.md` (the source-resolution capability;
  the harness evaluator serves resolution), `specs/evaluation/evaluation-json.md`
  (the `sources` provenance record; schema version 6),
  `specs/cli/evaluation-run.md` (dispatch plan in previews and receipts),
  `specs/skills/quality-skill/evaluation.md` (serving resolution requests),
  and `SPECIFICATION.md` (the format commits to non-filesystem selectors;
  `source` stays a single string — settled Q2/Q1). See the per-requirement
  annotations in the
  [functional spec](0197-resolver-dispatched-source-selectors/spec.md).
- **Bundled skill:** `skills/quality/workflows/evaluate.md` — the harness
  checkpoint loop serves `resolveSource` resolution requests alongside
  judgment.
- **Docs / generated artifacts:** `mintlify/specification.mdx` (regenerated
  from `SPECIFICATION.md`) and `CHANGELOG.md` release notes. `qualitymd` help
  text needed no change — no command help describes source selection.
  `quality.schema.json` and `specs/quality-schema-json.md` are deliberately
  unchanged — the settled Q1 keeps `source` a bare string.

## Children

- [Functional spec](0197-resolver-dispatched-source-selectors/spec.md)
- [Design doc](0197-resolver-dispatched-source-selectors/design.md)
