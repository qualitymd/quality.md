---
type: Runtime Resource
title: CLI Quick Reference
description: Quick reference for qualitymd CLI commands used by the /quality skill.
---

# CLI Quick Reference

Use these commands and routing patterns as the starting point for `qualitymd`
work.

## Output modes

| Goal                                      | Command form                                         |
| ----------------------------------------- | ---------------------------------------------------- |
| Human-readable command help               | `qualitymd <command> --help`                         |
| Machine-readable state or receipts        | `qualitymd <command> --json`                         |
| Active format specification text          | `qualitymd spec`                                     |
| Persist Evaluation v2 JSON                | `qualitymd evaluation data set <run> < payload.json` |
| Version and compatibility facts for skill | `qualitymd version --json`                           |

Use `--json` when a command offers it and the agent must inspect, route from,
or carry the result forward. Use human output for display or diagnostics only.

## CLI introspection

Use command help before guessing a command shape:

```sh
qualitymd --help
qualitymd lint --help
qualitymd status --help
qualitymd evaluation create --help
qualitymd evaluation data --help
qualitymd evaluation data kinds --help
qualitymd evaluation data example --help
qualitymd evaluation list --help
qualitymd evaluation status --help
qualitymd evaluation report --help
```

Use `qualitymd version --json` to inspect the visible CLI version,
development-build state, commit when known, bundled specification version, and
whether the install is in the skill's supported range.

## Quick reference

| Task                   | Command                                                    |
| ---------------------- | ---------------------------------------------------------- |
| Check CLI version      | `qualitymd --version`                                      |
| Read version metadata  | `qualitymd version --json`                                 |
| Check for CLI updates  | `qualitymd update --check`                                 |
| Apply CLI update       | `qualitymd update`                                         |
| Read format rules      | `qualitymd spec`                                           |
| Create a starter model | `qualitymd init [path]`                                    |
| Validate a model       | `qualitymd lint [path]`                                    |
| Fix simple lint issues | `qualitymd lint --fix [path]`                              |
| Inspect project status | `qualitymd status [path] --json`                           |
| Create evaluation run  | `qualitymd evaluation create [model] [--narrowing <slug>]` |
| List evaluation runs   | `qualitymd evaluation list [--json]`                       |
| Discover data kinds    | `qualitymd evaluation data kinds [--json]`                 |
| Print payload example  | `qualitymd evaluation data example <kind>`                 |
| Persist routine output | `qualitymd evaluation data set [-n] <run> < payload.json`  |
| List routine outputs   | `qualitymd evaluation data list <run> [--json]`            |
| Read routine output    | `qualitymd evaluation data get <run> --kind <kind> ...`    |
| Check reportability    | `qualitymd evaluation status <run>`                        |
| Build report           | `qualitymd evaluation report build <run>`                  |

Evaluation runs default under `.quality/evaluations/`. A repository can set
`evaluationDir` in the resolved workspace config file; the selected `QUALITY.md`
can point to that file with root `config` frontmatter, otherwise
`.quality/config.yaml` is used.

Other `.quality/` workspace artifacts: the quality log under `.quality/log/`
(singular) and workflow feedback logs under `.quality/logs/` (plural,
`<timestamp>-<workflow>-feedback-log.md`). Both are skill-written conventions
created on demand, with no `qualitymd` command yet; setup creates and updates its
current-run feedback log after CLI support is verified, and evaluate creates and
updates its current-run feedback log after the run frame.

## Decision trees

### Starting or repairing a model

```text
Need a QUALITY.md?
- No file yet? -> qualitymd init [path]
- File exists? -> qualitymd lint [path]
- Lint says fixes are available? -> qualitymd lint --fix [path]
- Need routing state? -> qualitymd status [path] --json
- Need format rules? -> qualitymd spec
```

### Evaluating

```text
Need to evaluate?
- Check model first -> qualitymd lint [path]
- Inspect current state -> qualitymd status [path] --json
- Create feedback log -> edit .quality/logs/<timestamp>-evaluate-feedback-log.md after the run frame
- Create run -> qualitymd evaluation create [model] [--narrowing <slug>]
- Discover payload shapes -> qualitymd evaluation data kinds; qualitymd evaluation data example <kind>
- Maintain workflow feedback -> update .quality/logs/<timestamp>-evaluate-feedback-log.md for material workflow-experience events only
- Validate routine output -> qualitymd evaluation data set <run> --dry-run < payload.json
- Persist routine output -> qualitymd evaluation data set <run> < payload.json
- Ready to report? -> qualitymd evaluation status <run>
- Build report -> qualitymd evaluation report build <run>
```

### Resuming or diagnosing a run

```text
Run incomplete or stale?
- Inspect project state -> qualitymd status [path] --json
- List runs -> qualitymd evaluation list --json
- Inspect run readiness -> qualitymd evaluation status <run>
- Unsupported old run shape? -> create a fresh Evaluation v2 run; do not hand-migrate records
- Process ambiguity or recovery? -> record concise notes in .quality/logs/<timestamp>-evaluate-feedback-log.md; do not duplicate assessment evidence
- Missing data? -> inspect data kinds/examples, validate with -n/--dry-run, then persist with qualitymd evaluation data set <run> < payload.json
- Reportable? -> qualitymd evaluation report build <run>
```

### Updating

```text
CLI missing, stale, or incompatible?
- Inspect visible CLI -> qualitymd version --json
- Check install-aware action -> qualitymd update --check
- Apply only when confirmed -> qualitymd update
```

## Common workflows

### Setup a new model

```sh
qualitymd init
qualitymd lint
qualitymd status --json
```

### Inspect readiness

```sh
qualitymd version --json
qualitymd lint [path]
qualitymd status [path] --json
```

### Create and complete an evaluation run

```sh
qualitymd evaluation create [model] [--narrowing <slug>]
qualitymd evaluation data kinds --json
qualitymd evaluation data example RequirementAssessmentResult
qualitymd evaluation data set <run> --dry-run < requirement-assessment-result.json
qualitymd evaluation data set <run> < requirement-assessment-result.json
qualitymd evaluation status <run>
qualitymd evaluation report build <run>
```

Use `--narrowing` for scoped evaluations. The slug should be the scope's full
structural path: Area path segments from the root, plus Factor path segments
when scoping to a Factor, joined with single hyphens.

## Command rules

- Use `--json` when a command offers it and the agent must consume the result.
- Prefer `qualitymd status [path] --json` for readiness, model shape, evaluation
  history, and stale-run signals.
- Persist routine JSON through `qualitymd evaluation data set`; use
  `qualitymd evaluation data example <kind>` to inspect payload shape.
- Use `-n/--dry-run` to validate new or materially revised payloads before
  writing.
- Do not continue past missing evaluation commands by manually creating files.
- Keep generated run paths exactly as the CLI reports them.
