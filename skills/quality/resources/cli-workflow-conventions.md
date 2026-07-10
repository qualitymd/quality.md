---
type: Runtime Resource
title: CLI workflow conventions
description: Workflow conventions for how the /quality skill drives qualitymd CLI introspection and artifacts.
---

# CLI workflow conventions

Use this resource for workflow conventions the CLI cannot teach by itself:
artifact layout, feedback-log sequencing, scoped-evaluation naming, and
cross-command orchestration. Do not treat it as a command reference.

## Introspection first

Discover command shapes, flags, and payload contracts from the CLI at runtime:

- Use command help before guessing an invocation shape.
- Use `--json` when a command offers it and the skill must inspect, route from,
  or carry the result forward.
- Use human output for display or diagnostics only.
- Use `qualitymd version --json` to inspect the visible CLI version,
  development-build state, commit when known, bundled specification version, and
  whether the install is in the skill's supported range.
- Use `qualitymd spec` for active format rules and rating vocabulary.
- Use `qualitymd evaluation run --dry-run --json` to preview a resolved
  evaluation — model, scope, evaluator, concurrency, and work-unit counts —
  without invoking an evaluator or writing evaluation data.
- Use `qualitymd evaluation data kinds`, `qualitymd evaluation data schema`, and
  `qualitymd evaluation data example` only when inspecting a historical
  multi-file run's payloads; new runs keep structured data in `evaluation.json`
  written by the runner.

Prefer stable structured channels (`--json`, schemas, examples) over parsing
human-formatted help tables when the result will drive routing or authored
artifacts.

## Workspace artifacts

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

## Starting or repairing a model

Sequence the work this way:

1. If no `QUALITY.md` exists, create one through the CLI starter workflow.
2. If a file exists, validate it through the CLI lint workflow.
3. If fixes are available and appropriate, use the CLI fix path.
4. When routing state is needed, inspect status through the structured status
   output.
5. When format rules are needed, read the active CLI-bundled spec.

## Evaluating

Sequence an evaluation this way:

1. Validate the model before invoking the runner.
2. Inspect current workspace status from structured status output.
3. Create the evaluate feedback log after the run frame.
4. Resolve any natural-language scope to canonical references through the model
   introspection command's structured list output.
5. Optionally preview the resolved run with the runner's dry-run JSON output,
   and confirm or explain evaluator selection.
6. Run the evaluation through `qualitymd evaluation run` with explicit flags
   and record the reported run path in the feedback log.
7. For a harness-backed run, service each `awaiting_evaluator` receipt: judge
   only the supplied bounded work request and submit the result envelope with
   `--resume <run> --evaluator-result -`, repeating until a terminal receipt.
8. Maintain workflow feedback for material process events.
9. Summarize the receipt and the generated reports.

The runner owns run creation, evaluator invocation, structured data, and report
generation; do not author or persist evaluation payloads for new runs.

## Resuming or diagnosing a run

Sequence recovery work this way:

1. Inspect workspace status from structured status output.
2. List evaluation runs through structured run-list output.
3. Inspect run readiness through the CLI evaluation-status path.
4. Record process ambiguity or recovery notes in the current evaluate feedback
   log; do not duplicate assessment evidence.
5. Resume a failed or cancelled runner run with
   `qualitymd evaluation run --resume <run>`, keeping the run's recorded
   evaluator; a different evaluator means a new run. A run awaiting harness
   judgment is resumable the same way: resuming without a result re-emits the
   pending work request.
6. For a historical multi-file run with missing data, inspect data kinds,
   schemas, and examples; validate with dry-run; then persist through the CLI
   data write path.
7. Rebuild reports through the CLI report-build path only when the run is
   reportable.

## Updating

When the CLI is missing, stale, incompatible, or uncertain:

1. Inspect the visible CLI through structured version output.
2. Use the install-aware update check to identify remediation.
3. Apply an update only after confirmation.

## Evaluation scope

Use `--area <area-ref>` and repeatable `--factor <factor-ref>` for scoped
evaluations. Pass canonical `area:` and `factor:` references resolved from the
model. Let `qualitymd evaluation run` record the run manifest, apply the root
default, and derive the run-folder slug.

## Command rules

- Use `--json` when a command offers it and the agent must consume the result.
- Prefer structured status output for readiness, model shape, evaluation
  history, and stale-run signals.
- Use model introspection output for canonical area, factor, and requirement IDs.
- Create and execute evaluations only through `qualitymd evaluation run`; the
  data write path exists for historical multi-file runs only.
- Do not continue past missing evaluation commands by manually creating files.
- Never manually create evaluation run folders or structured data files.
- Keep generated run paths exactly as the CLI reports them.
