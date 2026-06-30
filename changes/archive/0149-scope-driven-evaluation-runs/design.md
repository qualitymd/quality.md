---
type: Design Doc
title: Scope-driven evaluation runs — design
description: How create captures requested/planned scope, how the run report becomes the scoped area report, and how the headline concept is removed.
status: Done
tags: [evaluation, cli, skill, reports]
timestamp: 2026-06-27T00:00:00Z
---

# Scope-driven evaluation runs — design

## Context

Answers the [functional spec](spec.md). The headline bug
([parent motivation](../0149-scope-driven-evaluation-runs.md#motivation)) is a
symptom: scope lives only as a lossy `--narrowing` slug plus an agent-authored
frame enumeration, and the headline is selected positionally from that
enumeration. The fix makes scope a CLI-owned run parameter and lets the report
_be_ the scoped Area report.

## Approach

### Manifest record

`create` writes one new CLI-owned payload, `data/run-manifest.json`
(kind `RunManifest`), holding two same-shaped scope objects:

```jsonc
{
  "schemaVersion": <n>,
  "kind": "RunManifest",
  "number": 7,
  "model": "QUALITY.md",
  "requestedScope": { "areaId": "area:backend", "factorFilter": ["factor:backend::security"] },
  "plannedScope":   { "areaId": "area:backend", "factorFilter": ["factor:backend::security"] }
}
```

- `requestedScope` is faithful: fields present only when requested; absent when
  not. A bare run has `requestedScope: {}`.
- `plannedScope` is the normalization — `areaId` defaults to the root Area,
  `factorFilter` defaults to `[]`. It is the one always-well-formed field every
  consumer reads.
- The two diverge only when nothing was requested (planned defaults to the root
  Area). The CLI computes `plannedScope` once, against the snapshot it just froze.

Registration mirrors `EvaluationOutputResult`: a `DataKind` in `data.go`, a
contract in `data_contract.go`, a label in `display.go`, and the
agents-MUST-NOT-write rule enforced on the `data set` path.

### `create` flow

1. Snapshot the model (unchanged).
2. Resolve `--area`/`--factor` against the snapshot; reject unknown IDs (R6) and
   any Factor outside the Area (R7) before creating the numbered folder.
3. Normalize to `plannedScope` (root default).
4. Derive the slug from `plannedScope` (R8) and create `NNNN-<scope>-eval`.
5. Write `run-manifest.json`.

The skill resolves the user's natural-language ask to canonical IDs against the
live model _before_ calling `create` (R9, R11); `create` re-validates against the
frozen snapshot, so a live-vs-snapshot drift surfaces as a create-time error
rather than a silent mismatch.

### Planned expansion (the one new derivation)

`plannedScope` stores no enumeration; the planned set is computed on demand from
`plannedScope` + snapshot, used for the report's resolved-scope counts (R12) and
the coverage check (R19):

- **Area component** — `areaId` plus its descendant Areas.
- **Factor component** — when `factorFilter` is empty, every Factor of the
  in-scope Areas; otherwise the listed Factors plus their descendant Factors.
- **Requirement component** — the Requirements the in-scope Factors trace to, plus
  the in-scope Areas' Requirements.

The exact parent/trace relationships come from the model schema; the rule is
deterministic and judgment-free, so the CLI can compute it.

### Report reshape

`report build` drops headline resolution entirely (`resolveScopedFactorHeadline`
/ `resolveScopedAreaHeadline` / `resolveRootHeadline` and the headline plan
fields). Instead it reads `plannedScope` and renders `report.md` as the report
for that Area, narrowed by `factorFilter`:

- **Title** from `plannedScope` — `Backend`, `Backend — Security, Reliability`,
  or the root Area's label for a full run.
- **Lead rating + summary** from that Area's local-and-descendant analysis,
  computed over the filtered Factors; marked partial via limits when
  `factorFilter` is non-empty (R16).
- **Body** keeps the existing run index (subject reports, coverage, limits).

A single-Factor run is just the Area filtered to one Factor: the Area's
filtered roll-up equals that Factor's rating, and the Factor's own subject report
stays linked in the index. No special case.

`EvaluationOutputResult` loses `headlineResultRef`/`headlineReportRef`. It keeps
`runReportRef` (→ `report.md`) and gains a single deterministic
`scopedAreaAnalysisRef` (→ the Area analysis `report.md` renders), so consumers
still have a structured pointer to the run's lead result.

For a full run, `report.md` and the root Area's report are the same content, so
the scoped Area's report _is_ `report.md` — no duplicate file is emitted for the
scoped Area. Child Areas still get their per-Area subject reports in the index.

## Spec response

- **Scope data (R1–R4)** — `RunManifest` with faithful `requestedScope` and
  normalized `plannedScope`; expansion derived, not stored.
- **Flags/validation (R5–R8)** — `--area`/`--factor`, snapshot validation,
  coherence, derived slug, in the `create` flow above.
- **Skill (R9–R11)** — resolve-before-create; CLI owns defaulting/expansion.
- **Report (R12–R16, R19)** — report-as-Area-report from `plannedScope`; headline
  removed; requested line; partial marking; coverage as planned-vs-produced.
- **Frame/list (R17, R18)** — frame sheds scope; `list` reads the manifest.

Determinism (R14, AC8) follows structurally: the report's subject is a function
of `plannedScope` alone, with no agent-authored ordering in the path.

## Alternatives

- **Fold scope into `EvaluationFrame.inputs`.** Rejected: the frame is
  agent-authored and post-create, so `create` could not persist or validate scope,
  and the report would still derive its subject from an agent artifact — the exact
  fragility being removed.
- **Store the planned expansion in the manifest.** Rejected: the expansion is a
  pure function of `plannedScope` + snapshot, so storing it adds a field that can
  drift or go stale on an interrupted run. Derive it.
- **Keep a headline subject, fix only its selection rule.** Rejected: once
  `plannedScope` names the subject, the report _is_ that Area's report; a separate
  selection step is needless surface that can still be got wrong.
- **Keep scope (resolved set) in the frame for the coverage check.** Rejected:
  the CLI-derived planned expansion is a stronger expected-set baseline than an
  agent-declared list, so the frame need not carry scope at all.

## Trade-offs & risks

- **New payload kind.** `RunManifest` adds contract/registration/label/test
  surface — accepted as the keystone that makes scope authoritative.
- **Clean break.** Dropping `--narrowing` and pre-manifest run support means old
  runs are unreadable by the new `list`/`report` paths. Acceptable under
  early-alpha rules; no shim.
- **Expansion rule is now a CLI contract.** Counts and coverage depend on the
  Area/Factor/Requirement relationships being modeled correctly; mis-modeling
  shows up in coverage. Mitigated by tests over a fixture model.

## Open questions

- Whether to lift `RunManifest` into its own `run-manifest-json.md` artifact spec
  now or leave the contract in `payload-kinds.md` (spec lists it as _suggested_).
- Exact field name for the output-result pointer to the scoped Area analysis
  (`scopedAreaAnalysisRef` proposed).
