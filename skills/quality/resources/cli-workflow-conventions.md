---
type: Runtime Resource
title: CLI Workflow Conventions
description: Workflow conventions for how the /quality skill drives qualitymd CLI introspection and artifacts.
---

# CLI Workflow Conventions

Use this resource for workflow conventions the CLI cannot teach by itself:
artifact layout, feedback-log sequencing, scoped-evaluation naming, and
cross-command orchestration. Do not treat it as a command reference.

## Introspection First

Discover command shapes, flags, and payload contracts from the CLI at runtime:

- Use command help before guessing an invocation shape.
- Use `--json` when a command offers it and the skill must inspect, route from,
  or carry the result forward.
- Use human output for display or diagnostics only.
- Use `qualitymd version --json` to inspect the visible CLI version,
  development-build state, commit when known, bundled specification version, and
  whether the install is in the skill's supported range.
- Use `qualitymd spec` for active format rules and rating vocabulary.
- Use `qualitymd evaluation data kinds`, `qualitymd evaluation data schema`, and
  `qualitymd evaluation data example` for Evaluation payload discovery. For one
  kind, `data schema <kind>` is self-contained and is the source for required
  fields and allowed enum values; `data example <kind>` is one concrete valid
  instance.

Prefer stable structured channels (`--json`, schemas, examples) over parsing
human-formatted help tables when the result will drive routing or authored
artifacts.

## Workspace Artifacts

Evaluation runs default under `.quality/evaluations/` relative to the selected
`QUALITY.md`. A repository can set `evaluationDir` in the resolved workspace
config file; the selected `QUALITY.md` can point to that file with root `config`
frontmatter, otherwise `.quality/config.yaml` beside the selected model is used.
Relative tooling paths are model-relative and must remain inside the Git
repository root found from that model.

Other `.quality/` workspace artifacts:

- `.quality/changelog/` - the quality changelog. Written only for confirmed,
  meaningful model changes.
- `.quality/logs/` - flat workflow and process logs. Setup and evaluate create
  current-run feedback logs here; future log kinds should use clear filenames,
  not subfolders.

Feedback logs are skill-written conventions created on demand. Setup creates and
updates its current-run feedback log after the setup preview when the run
continues into discovery. Evaluate creates and updates its current-run feedback
log after the run frame. Record material workflow-experience events only; do not
duplicate assessment evidence.

## Starting Or Repairing A Model

Sequence the work this way:

1. If no `QUALITY.md` exists, create one through the CLI starter workflow.
2. If a file exists, validate it through the CLI lint workflow.
3. If fixes are available and appropriate, use the CLI fix path.
4. When routing state is needed, inspect status through the structured status
   output.
5. When format rules are needed, read the active CLI-bundled spec.

## Evaluating

Sequence an evaluation this way:

1. Validate the model before judgment.
2. Inspect current workspace status from structured status output.
3. Create the evaluate feedback log after the run frame.
4. Create the run through the CLI and record the run path in the feedback log.
5. Query in-scope canonical IDs from the run's `model-snapshot.md`, never from
   the live `QUALITY.md`; use the model introspection command's structured list
   output and scope it to the resolved evaluation scope.
6. Discover payload kinds, schemas, and examples through the CLI data discovery
   commands.
7. Maintain workflow feedback for material process events.
8. Validate authored routine output with the CLI data dry-run path before
   writing.
9. Persist routine output through the CLI data write path.
10. Check reportability through the CLI evaluation-status path.
11. Build the report through the CLI report-build path.

Use the queried model IDs as the source of truth for every authored payload
reference. The post-hoc identity-resolution check is a backstop, not the primary
guard.

## Resuming Or Diagnosing A Run

Sequence recovery work this way:

1. Inspect workspace status from structured status output.
2. List evaluation runs through structured run-list output.
3. Inspect run readiness through the CLI evaluation-status path.
4. Record process ambiguity or recovery notes in the current evaluate feedback
   log; do not duplicate assessment evidence.
5. For missing data, inspect data kinds, schemas, and examples; validate with
   dry-run; then persist through the CLI data write path.
6. Build the report only after the run is reportable.

## Updating

When the CLI is missing, stale, incompatible, or uncertain:

1. Inspect the visible CLI through structured version output.
2. Use the install-aware update check to identify remediation.
3. Apply an update only after confirmation.

## Evaluation Scope

Use `--area <area-id>` and repeatable `--factor <factor-id>` for scoped
evaluations. Pass canonical `area:` and `factor:` references resolved from the
model. Let `qualitymd evaluation create` write `EvaluationManifest`, apply the root
default, and derive the run-folder slug.

## Command Rules

- Use `--json` when a command offers it and the agent must consume the result.
- Prefer structured status output for readiness, model shape, evaluation
  history, and stale-run signals.
- Use model introspection output for canonical Area, Factor, and Requirement IDs.
- Persist routine JSON only through the CLI data write path.
- Use the CLI payload schema and populated example commands to inspect payload
  shape.
- Use dry-run to validate new or materially revised payloads before writing.
- Do not continue past missing evaluation commands by manually creating files.
- Never manually create evaluation run folders or structured data files.
- Keep generated run paths exactly as the CLI reports them.
