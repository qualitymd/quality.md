---
type: Functional Specification
title: Scope-driven evaluation runs — functional spec
description: Capture requested/planned scope as CLI-owned run data at create time and render the run report as the scoped area report.
status: Done
tags: [evaluation, cli, skill, reports]
timestamp: 2026-06-27T00:00:00Z
---

# Scope-driven evaluation runs — functional spec

Companion to the [Scope-driven evaluation runs](../0149-scope-driven-evaluation-runs.md)
change case. This spec states *what* the change must do; the [design doc](design.md)
covers *how*. It defers to the [Evaluation](../../../specs/evaluation/evaluation.md)
contract and the format spec [`SPECIFICATION.md`](../../../SPECIFICATION.md) as
normative sources.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

The run report's "headline" subject is selected positionally: `report build`
reads `EvaluationFrame.inputs.factorIds` and features the first listed Factor
with an analysis. A full evaluation lists every Factor in model order, so the
headline collapses to whichever Factor appears first in the model and reports its
local rating instead of the run's overall result — masking, in the observed case,
a 🔴 Unacceptable root Area behind a 🟡 Minimum Security Factor.

The deeper cause is that scope is never captured as authoritative structured
data. `create` records scope only as a lossy `--narrowing` folder slug, and the
structured scope that drives the headline is reconstructed by the agent into the
frame after the fact. This change captures scope as a CLI-owned run parameter at
`create` and makes the run report *be* the scoped Area's report, so there is no
subject to select and no selection to get wrong.

## Vocabulary

- **Requested scope** — the faithful record of what a run was created for:
  `{ areaId?, factorFilter? }`. A field is absent when nothing of that kind was
  requested.
- **Planned scope** — requested scope **normalized**: same shape, with `areaId`
  always present (defaulting to the root Area) and `factorFilter` always an array
  (possibly empty).
- **`factorFilter`** — a list of canonical Factor IDs that narrows a run *within*
  its `areaId`. An empty filter means *no narrowing* — the whole Area, not zero
  Factors.
- **Planned expansion** — the concrete set of Areas, Factors, and Requirements a
  planned scope covers. It is **derived** from `plannedScope` and the model
  snapshot; it is never stored.

## Scope

Covered: the `RunManifest` record and its two scope fields; the `create`
`--area`/`--factor` flags and validation; the run report rendered as the scoped
Area report; removal of the headline concept and of scope from
`EvaluationFrame`; manifest-sourced `list` and coverage; and the skill's
resolve-before-create flow.

Deferred / non-goals: no `--narrowing` or pre-`RunManifest` back-compat
(early-alpha clean break); no change to judgment routines or rating semantics;
the natural-language ask stays in the evaluate feedback log, not the manifest.

## Requirements

### Run manifest and scope data

R1. `qualitymd evaluation create` **MUST** write a CLI-owned `RunManifest` record
into the run containing at least `schemaVersion`, the run number, the model path,
`requestedScope`, and `plannedScope`.

> Rationale: scope must be authoritative from the first mutation; the deterministic
> CLI owns the run boundary, so it — not a later agent artifact — records scope.

>> Durable spec: modify `specs/evaluation/records/payload-kinds.md` — add the
>> `RunManifest` kind; modify `specs/cli/evaluation-create.md` and
>> `specs/evaluation/records/data-layout.md` — `create` writes it into the run
>> layout.

R2. Agents **MUST NOT** write `RunManifest` through `qualitymd evaluation data set`.

>> Durable spec: modify `specs/evaluation/records/payload-kinds.md` — list
>> `RunManifest` as CLI-owned, mirroring `EvaluationOutputResult`.

R3. `requestedScope` **MUST** record exactly what was supplied to `create`,
without applying defaults or expansion; `areaId` and `factorFilter` are each
present only when an Area or Factor, respectively, was requested.

>> Durable spec: modify `specs/cli/evaluation-create.md`,
>> `specs/evaluation/records/payload-kinds.md` — define `requestedScope` as the
>> faithful record.

R4. `create` **MUST** record `plannedScope` as the normalization of
`requestedScope`: `areaId` set to the requested Area or, when none was requested,
the root Area; `factorFilter` set to the requested Factor IDs or an empty array.
`plannedScope` **MUST NOT** store a structural expansion.

> Rationale: one always-well-formed field localizes the root default to `create`,
> so no downstream consumer re-derives it inconsistently; the expansion is a pure
> function of `plannedScope` and the snapshot and is therefore derived, not stored.

>> Durable spec: modify `specs/cli/evaluation-create.md`,
>> `specs/evaluation/records/payload-kinds.md` — define `plannedScope` and its
>> normalization.

### `create` flags and validation

R5. `qualitymd evaluation create` **MUST NOT** accept `--narrowing`. It **MUST**
accept `--area <area-id>` (at most once) and `--factor <factor-id>` (repeatable).

>> Durable spec: modify `specs/cli/evaluation-create.md` — replace the
>> `--narrowing` flag with `--area`/`--factor`.

R6. If any `--area` or `--factor` value does not resolve against the run's model
snapshot, then `create` **MUST** fail without creating a numbered run folder.

>> Durable spec: modify `specs/cli/evaluation-create.md` — require snapshot
>> validation of scope arguments.

R7. If any `--factor` value does not belong to the resolved `--area`, then
`create` **MUST** fail without creating a numbered run folder.

> Rationale: a manifest whose Factors span Areas is incoherent; catching it at
> create keeps `plannedScope` self-consistent for every downstream reader.

>> Durable spec: modify `specs/cli/evaluation-create.md` — require single-Area
>> coherence between `--factor` and `--area`.

R8. `create` **MUST** derive the run folder slug from `plannedScope`: the root
Area with an empty `factorFilter` produces `NNNN-full-eval`; otherwise the slug
is built from the Area and filtered-Factor structural paths. The slug **MUST NOT**
be supplied by the caller.

>> Durable spec: modify `specs/cli/evaluation-create.md` — derive `<scope>` in
>> the `NNNN-<scope>-eval` grammar from `plannedScope`, not from `--narrowing`.

### Skill scope resolution

R9. When the user requests an Area and/or Factor, the `/quality` skill **MUST**
resolve it to canonical Area/Factor IDs, validate it against the model, and pass
it to `create` via `--area`/`--factor`.

>> Durable spec: modify `specs/skills/quality-skill/evaluation.md`,
>> `specs/skills/quality-skill/workflows/evaluate.md` — resolve requested scope to
>> IDs and pass it through the flags.

R10. The skill **MUST NOT** apply a root-Area default and **MUST NOT** enumerate
planned scope or its expansion; those are `create` responsibilities.

> Rationale: defaulting and expansion are deterministic CLI work; duplicating them
> in the skill reintroduces the divergence this case removes.

>> Durable spec: modify `specs/skills/quality-skill/evaluation.md` — forbid skill
>> defaulting/enumeration of scope.

R11. The skill **MUST** resolve requested scope before invoking `create`.

>> Durable spec: modify `specs/skills/quality-skill/workflows/evaluate.md` — move
>> scope resolution ahead of run creation.

### Run report as scoped area report

R12. `qualitymd evaluation report build` **MUST** render `report.md` as the
report for `plannedScope.areaId` narrowed by `factorFilter`: its title and lead
rating **MUST** reflect that Area and any filtered Factors.

> Rationale: once `plannedScope` names what the run is about, the report *is* that
> Area's report; there is no separate subject to feature.

>> Durable spec: modify `specs/cli/evaluation-report.md`,
>> `specs/evaluation/reports/report-tree.md` — `report.md` is the scoped Area
>> report titled from `plannedScope`.

R13. The run output **MUST NOT** carry a headline subject: `EvaluationOutputResult`
**MUST NOT** include `headlineResultRef` or `headlineReportRef`, and `report.md`
**MUST NOT** present a headline-subject selection.

>> Durable spec: modify `specs/evaluation/records/payload-kinds.md` — remove the
>> headline refs from `EvaluationOutputResult`; modify
>> `specs/evaluation/protocol.md` — remove the "Headline Result" section; modify
>> `specs/evaluation/orchestration.md` and `SPECIFICATION.md` — drop headline
>> wording.

R14. `report build` **MUST** determine the run report's Area and filter solely
from `plannedScope`; it **MUST NOT** depend on the ordering of any
agent-authored payload.

> Rationale: this is the fix — removing positional dependence on frame ordering
> makes the report deterministic for the same `plannedScope`.

>> Durable spec: modify `specs/cli/evaluation-report.md` — source the report's
>> subject from `plannedScope`, not frame inputs.

R15. `report.md` **MUST** render a requested-scope line from `requestedScope`,
showing "full evaluation" when `requestedScope` requested neither an Area nor a
Factor.

>> Durable spec: modify `specs/cli/evaluation-report.md` — requested-scope line
>> from `requestedScope`.

R16. If `factorFilter` is non-empty, then the run report's Area rating **MUST** be
marked a partial roll-up through the existing evaluation-limits mechanism.

> Rationale: a Factor-filtered Area is not a complete Area assessment; an unmarked
> rating would overclaim coverage.

>> Durable spec: modify `specs/cli/evaluation-report.md` — mark a filtered Area
>> rating as partial via limits.

### Frame, list, and coverage

R17. `EvaluationFrame` **MUST NOT** carry scope fields (`requestedScope`,
`areaIds`, `factorIds`); scope is read from `RunManifest`.

>> Durable spec: modify `specs/evaluation/records/payload-kinds.md` — remove
>> scope from the `EvaluationFrame` contract.

R18. `qualitymd evaluation list` **MUST** read each run's number and scope from
its `RunManifest` rather than by parsing the run folder name.

> Rationale: the folder slug is a derived mnemonic; reading the manifest makes the
> structured scope, not the lossy name, the source of truth.

>> Durable spec: modify `specs/cli/evaluation-list.md` — source run number and
>> scope from `RunManifest`.

R19. The coverage check **MUST** compare the planned expansion derived from
`plannedScope` and the model snapshot against the analysis artifacts actually
produced.

>> Durable spec: modify `specs/cli/evaluation-report.md`,
>> `specs/evaluation/protocol.md` — define coverage as planned-expansion vs
>> produced.

## Acceptance criteria

- AC1. Bare `create` (no scope flags) writes a `RunManifest` with `requestedScope`
  empty and `plannedScope = { areaId: <root>, factorFilter: [] }`, and the folder
  is `NNNN-full-eval`. *(R1, R3, R4, R8)*
- AC2. `create --factor <root Security>` (no `--area`) records `requestedScope`
  with that Factor and root `areaId`, an identical `plannedScope`, and `report.md`
  is titled for the root Area narrowed to Security. *(R3, R4, R12)*
- AC3. `create --area <backend>` records `plannedScope.factorFilter = []` and
  `report.md` is the Backend Area report (full roll-up). *(R4, R12)*
- AC4. `create --area <backend> --factor <backend Security> --factor <backend
  Reliability>` titles `report.md` for the Backend Area with both Factors, and
  the Area rating is marked partial. *(R12, R16)*
- AC5. `create --area <backend> --factor <mobile Security>` fails and creates no
  run folder. *(R7)*
- AC6. `create --factor <nonexistent>` fails against the snapshot and creates no
  run folder. *(R6)*
- AC7. A full run over a model whose root Area is Unacceptable produces a
  `report.md` titled for the root Area and rated Unacceptable — not a Minimum
  Factor. *(R12, R14)*
- AC8. Building the report twice from the same run yields an identical title and
  rating regardless of payload ordering. *(R14)*
- AC9. `evaluation list` reports a run's scope correctly after its folder name is
  manually altered. *(R18)*
- AC10. An `EvaluationFrame` payload containing scope fields is rejected by the
  data contract. *(R17)*
- AC11. `create --narrowing <slug>` is rejected as an unknown flag. *(R5)*

## Durable spec changes

### To add

- *(Suggested)* `specs/evaluation/records/run-manifest-json.md` — a 1:1 artifact
  spec for the new `run-manifest.json` record, if its contract is lifted out of
  `specs/evaluation/records/payload-kinds.md` (drives R1–R4). Not a precondition
  for landing; the contract may live in `payload-kinds.md`.

### To modify

- `specs/cli/evaluation-create.md` — `--area`/`--factor` flags, snapshot
  validation, single-Area coherence, `RunManifest` write, slug derived from
  `plannedScope` (per R1, R3–R8).
- `specs/cli/evaluation-list.md` — read number and scope from `RunManifest` (per
  R18).
- `specs/cli/evaluation-report.md` — `report.md` as the scoped Area report,
  requested-scope line, partial-roll-up marking, subject sourced from
  `plannedScope`, coverage as planned-vs-produced (per R12, R14–R16, R19).
- `specs/evaluation/protocol.md` — remove the "Headline Result" section; define
  coverage as planned-expansion vs produced (per R13, R19).
- `specs/evaluation/records/payload-kinds.md` — add the CLI-owned `RunManifest`
  kind; remove scope from `EvaluationFrame`; remove headline refs from
  `EvaluationOutputResult` (per R1–R4, R13, R17).
- `specs/evaluation/reports/report-tree.md` — `report.md` shape and title from
  `plannedScope` (per R12).
- `specs/evaluation/orchestration.md` — scope/headline flow (per R13).
- `specs/evaluation/records/data-layout.md` — run layout includes
  `run-manifest.json` (per R1).
- `specs/skills/quality-skill/evaluation.md` — resolve-before-create, pass scope
  flags, forbid skill defaulting/enumeration, frame sheds scope (per R9–R11, R17).
- `specs/skills/quality-skill/workflows/evaluate.md` — scope-resolution ordering
  and `create` invocation (per R9, R11).
- `specs/skills/quality-skill/quality-skill.md` — replace narrowing references
  (per R5, R9).
- `SPECIFICATION.md` — reframe headline/scope wording (per R13).

### To rename

None.

### To delete

None.
