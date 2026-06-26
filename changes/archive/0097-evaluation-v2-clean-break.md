---
type: Change Case
title: Evaluation v2 clean break
description: Remove legacy evaluation compatibility and close the Evaluation v2 report, status, data-set, skill, and spec gaps found after the first real generated run.
status: Done
tags: [evaluation, cli, reports, skill, cleanup]
timestamp: 2026-06-26T00:00:00Z
---

# Evaluation v2 clean break

A **Change Case** to make Evaluation v2 the only runtime evaluation contract.
The previous replacement case introduced the new protocol, but the first real
generated run showed that legacy record/report paths, transition wording, and
underspecified v2 payload/report behavior still shape the product.

Detail lives in:

- [Functional spec](0097-evaluation-v2-clean-break/spec.md) - the removal
  punchlist and v2-only behavior this case must deliver.
- [Design doc](0097-evaluation-v2-clean-break/design.md) - how the clean break
  is implemented across CLI, evaluation graph validation, report rendering, and
  legacy removal.

## Motivation

Evaluation v2 is meant to turn evaluation into a structured agent judgment
protocol: routine outputs under `data/`, deterministic report trees, and clear
status over the v2 graph. In practice, the generated report from the first real
run did not match `evaluation-v2-sketch.md`: ratings rendered as status text,
finding details disappeared, requirement titles were lost, status still reported
legacy counts, and report build treated incomplete or mismatched v2 data as more
renderable than it was.

The repo still carries previous-workflow assumptions in code, specs, skill
instructions, examples, and tests. Keeping those paths as compatibility surfaces
makes the v2 workflow harder to reason about and easier to regress. This case is
therefore a clean break: no historical compatibility, no mixed runtime model,
and no `--file` data-set ergonomics carried forward from the transition.

## Scope

Covered:

- Remove previous evaluation record/runtime compatibility from active code,
  command behavior, durable specs, bundled skill guidance, and examples.
- Make Evaluation v2 status, data persistence, report build, and report
  rendering operate against one v2 data contract.
- Replace `qualitymd evaluation data set --file <path|->` with stdin-only input.
- Close the generated-report gaps found by comparing the acquire run with
  `evaluation-v2-sketch.md`.
- Reconcile still-current `evaluation-v2-sketch.md` protocol, JSON, orchestration,
  synthesis, report-tree, and routine-prompt details into durable v2 specs and
  skill guidance.
- Add tests and search checks that keep legacy evaluation names, artifacts, and
  compatibility paths out of live runtime surfaces.

Deferred / non-goals:

- Do not migrate, load, transform, or preserve old assessment/analysis/
  recommendation runs.
- Do not add batch/bulk data-set import yet; record it as a feedback-log
  opportunity after stdin-only single-payload input is settled.
- Do not reintroduce recommendation generation or report gates unless they are
  designed as explicit v2 features in a later case.
- Do not add `qualitymd evaluation data schema <kind>` or judgment QC routines
  in this case; keep both as explicit deferrals.
- Do not edit archived Change Cases or append-only history logs except to add
  this case's normal log entry.

## Affected artifacts

Derived by analysis: searched live code, specs, bundled skill files, examples,
README/docs, and tests for `--file`, `assessment add`, `analysis set`,
`report-summary.md`, `report.json`, `assessments/`, `analysis/`,
`RecordCounts`, `AddRecord`, `WriteRecords`, `GateReport`, `assessmentResults`,
`recommendations`, `legacy`, `evaluation data set`, and `evaluation report`.
Historical archive entries remain frozen except for this case once it lands.

### Code

- [x] `internal/evaluation/create.go` - remove legacy run directories, legacy
      plan coverage scaffolding, and `--file` next-action text.
- [x] `internal/evaluation/load.go`, `list.go`, and related status helpers -
      remove old record loading/counting and make run status v2-only.
- [x] `internal/evaluation/write.go`, `types.go`, and
      `input_contract.go` - delete or replace old record writer/input contracts
      such as `AddRecord`, `WriteRecords`, assessment results, analyses, and
      recommendations.
- [x] `internal/evaluation/report.go` and `report_v2.go` - remove the legacy
      `report-summary.md` / `report.json` path, remove or redesign
      `GateReport`, require complete v2 data validation, and fix report
      rendering gaps.
- [x] `internal/evaluation/data.go`, `planned_coverage.go`,
      `model_reference.go`, and `path.go` - align persistence, model identity,
      and ordering helpers with the v2-only graph.
- [x] `internal/cli/evaluation.go` - remove `--file`, read one data payload
      from stdin, and remove legacy command/help/status wording.
- [x] `internal/status/status.go` - stop projecting legacy evaluation counts or
      active legacy recommendations into project status.
- [x] Tests under `internal/**` - replace old-run fixtures and writer/report
      tests with v2-only status, data, and report regression coverage.
- [x] Live-code legacy inventory - search runtime, CLI, tests, and fixtures for
      previous-evaluation leftovers such as old record writers, legacy report
      builders, compatibility checks, old next-action text, and old artifact
      names; delete or replace them rather than preserving historical
      compatibility.

### Format spec (`SPECIFICATION.md`)

- [x] `SPECIFICATION.md` - confirm evaluation semantics describe the v2-only
      runtime contract and do not imply legacy report/record compatibility.

### Reference sketch

- [x] `evaluation-v2-sketch.md` - reconcile the sketch with the durable v2 specs
      and generated report behavior, then either keep it clearly as a historical
      sketch or mark it superseded by the durable specs.

### Durable specs (`specs/`)

The functional spec's
[Durable spec changes](0097-evaluation-v2-clean-break/spec.md#durable-spec-changes)
section is authoritative. In summary:

- [x] `specs/evaluation-v2/**` - tighten payload schemas, status vocabulary,
      output assembly, report tree, validation, ordering, routine prompt
      contracts, JSON conventions, shared refs, orchestration, synthesis
      defaults, and sketch deferrals.
- [x] `specs/evaluation-records.md` and `specs/evaluation-records/**` - remove
      from the active runtime source of truth or archive/delete the previous
      record contracts with no compatibility promise.
- [x] `specs/reports/**` - remove or replace legacy `report-summary.md`,
      `report.md`, and `report.json` contracts as active runtime specs.
- [x] `specs/cli/evaluation-create.md`, `evaluation-data.md`,
      `evaluation-status.md`, `evaluation-list.md`, `evaluation-report.md`, and
      `specs/cli.md` - align command behavior and output with v2-only
      evaluation.
- [x] `specs/skills/quality-skill/**` - mirror the v2-only workflow, report,
      example, and no-compatibility contract.

### Durable docs

- [x] `README.md` - update only if evaluation prose still implies old records,
      recommendations, reports, or CLI-managed evaluation as the primary
      experience.
- [x] `docs/guides/**` - update any live guide that mentions previous
      evaluation artifacts, `--file`, or compatibility behavior.
- [x] `CHANGELOG.md` - add the user-facing clean-break note when implemented.

### Bundled skill (`skills/quality/`)

- [x] `skills/quality/SKILL.md` - remove `--file` and any compatibility or old
      artifact contract language.
- [x] `skills/quality/workflows/evaluate.md` - make stdin data persistence,
      v2-only status/report behavior, feedback-log learning, and no
      recommendations explicit.
- [x] `skills/quality/resources/cli-quick-reference.md` and
      `output-policy.md` - align command examples and report-consumption rules.
- [x] `skills/quality/guides/**` - update evaluation-related guidance that
      still reads old report bodies or old record examples.

### Evaluation feedback log

- [x] Review the evaluation feedback log for workflow improvements and add only
      v2-aligned changes to this case; record batch/bulk import, schema helper
      commands, judgment QC routines, or recommendation generation as explicit
      follow-up opportunities unless they are required to close the v2 clean
      break.

### Examples and fixtures

- [x] `specs/skills/quality-skill/examples/**` - remove, replace, or clearly
      re-scope legacy example runs so active examples do not teach the previous
      runtime contract.
- [x] Test fixtures under `internal/**` - add an acquire-shaped v2 regression
      fixture that exercises status vocabulary, rating drivers, findings,
      evidence refs, factor-declared Requirement titles, and report ordering.

### Install / scaffold

None expected unless the implementation sweep finds evaluation scaffolding in
install templates.

## Children

- [Functional spec](0097-evaluation-v2-clean-break/spec.md) - the removal
  punchlist and v2-only behavior this case must deliver.
- [Design doc](0097-evaluation-v2-clean-break/design.md) - how the clean break
  is implemented across CLI, evaluation graph validation, report rendering, and
  legacy removal.

## Status

`Done`. Evaluation v2 is now the only active runtime evaluation path; the
previous record/report compatibility surface has been removed from active code,
durable specs, bundled skill guidance, examples, and tests.
