---
type: Runtime Resource
title: CLI Quick Reference
description: Quick reference for qualitymd CLI commands used by the /quality skill.
---

# CLI Quick Reference

Use these commands and routing patterns as the starting point for `qualitymd`
work.

## Output modes

| Goal                                      | Command form                     |
| ----------------------------------------- | -------------------------------- |
| Human-readable command help               | `qualitymd <command> --help`     |
| Machine-readable state or receipts        | `qualitymd <command> --json`     |
| Active format specification text          | `qualitymd spec`                 |
| Record JSON into a write command          | pipe on stdin (heredoc); no file |
| Version and compatibility facts for skill | `qualitymd version --json`       |

Use `--json` when a command offers it and the agent must inspect, route from,
or carry the result forward. Use human output for display or diagnostics only.

## CLI introspection

Use command help before guessing a command shape:

```sh
qualitymd --help
qualitymd lint --help
qualitymd status --help
qualitymd evaluation create --help
qualitymd evaluation list --help
qualitymd evaluation status --help
qualitymd evaluation assessment --help
qualitymd evaluation analysis --help
qualitymd evaluation recommendation --help
qualitymd evaluation report --help
```

Use `qualitymd version --json` to inspect the visible CLI version,
development-build state, commit when known, bundled specification version, and
whether the install is in the skill's supported range.

## Quick reference

| Task                          | Command                                                             |
| ----------------------------- | ------------------------------------------------------------------- |
| Check CLI version             | `qualitymd --version`                                               |
| Read version metadata         | `qualitymd version --json`                                          |
| Check for CLI updates         | `qualitymd update --check`                                          |
| Apply CLI update              | `qualitymd update`                                                  |
| Read format rules             | `qualitymd spec`                                                    |
| Create a starter model        | `qualitymd init [path]`                                             |
| Validate a model              | `qualitymd lint [path]`                                             |
| Fix simple lint issues        | `qualitymd lint --fix [path]`                                       |
| Inspect project status        | `qualitymd status [path] --json`                                    |
| Create evaluation run         | `qualitymd evaluation create [--model <path>] [--narrowing <slug>]` |
| List evaluation runs          | `qualitymd evaluation list [--json]`                                |
| Add assessment result records | pipe JSON \| `qualitymd evaluation assessment add [-n] <run>`       |
| Set analysis records          | pipe JSON \| `qualitymd evaluation analysis set [-n] <run>`         |
| Add recommendation records    | pipe JSON \| `qualitymd evaluation recommendation add [-n] <run>`   |
| List records                  | `qualitymd evaluation <kind> list <run>`                            |
| Check reportability           | `qualitymd evaluation status <run>`                                 |
| Build report                  | `qualitymd evaluation report build <run>`                           |
| Gate report                   | `qualitymd evaluation report gate <run> --at-or-below <level>`      |

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
- Create run -> qualitymd evaluation create [--model <path>] [--narrowing <slug>]
- Author initial design and plan -> edit design.md and plan.md before assessment evidence collection or record writes
- Planned coverage needed? -> edit plan.md coverage frontmatter before dependent record writes
- Maintain workflow feedback -> update .quality/logs/<timestamp>-evaluate-feedback-log.md for material workflow-experience events only
- Inspect payload shape -> qualitymd evaluation assessment add | analysis set | recommendation add --help
- Validate judgment records -> pipe JSON on stdin with -n/--dry-run
- Add judgment records -> pipe JSON on stdin to qualitymd evaluation assessment add | analysis set | recommendation add <run>
- Ready to report? -> qualitymd evaluation status <run>
- Build report -> qualitymd evaluation report build <run>
```

### Resuming or diagnosing a run

```text
Run incomplete or stale?
- Inspect project state -> qualitymd status [path] --json
- List runs -> qualitymd evaluation list --json
- Inspect run readiness -> qualitymd evaluation status <run>
- Incompatible historical record? -> treat as run status; inspect or create a fresh run, do not hand-migrate records
- Process ambiguity or recovery? -> record concise notes in .quality/logs/<timestamp>-evaluate-feedback-log.md; do not duplicate assessment evidence
- Missing or changed planned coverage? -> edit plan.md coverage frontmatter and add a plan amendment when the planned scope changed
- Missing records? -> inspect --help, validate with -n/--dry-run, then pipe JSON on stdin to qualitymd evaluation assessment add | analysis set | recommendation add <run>
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
qualitymd evaluation create [--model <path>] [--narrowing <slug>]

# Write records by piping JSON on stdin. Use --dry-run first, then remove it to commit.
qualitymd evaluation assessment add --dry-run <run> <<'JSON'
[
  {
    "areaPath": [],
    "requirement": "Has tests",
    "factorPaths": [],
    "ratingResult": {
      "kind": "rated",
      "level": "target",
      "rationale": "Evidence supports the target level."
    },
    "criterionSource": "rating-scale",
    "findings": [
      {
        "locator": "tests/example_test.go:1",
        "observation": "The requirement is covered by a focused test.",
        "category": "coverage",
        "severity": "low",
        "evidence": [
          {
            "kind": "source",
            "ref": "tests/example_test.go:1"
          }
        ]
      }
    ],
    "recommendations": []
  }
]
JSON

qualitymd evaluation analysis set --dry-run <run> <<'JSON'
[
  {
    "areaPath": [],
    "localRatingResult": {
      "kind": "rated",
      "level": "target",
      "rationale": "The local assessment result reaches target."
    },
    "factorRatingResults": [],
    "aggregateRatingResult": {
      "kind": "rated",
      "level": "target",
      "rationale": "The root local rating binds the aggregate rating."
    },
    "assessmentResultRecords": [
      "assessments/001-root-has-tests.json"
    ],
    "childAnalysisRecords": []
  }
]
JSON

qualitymd evaluation recommendation add --dry-run <run> <<'JSON'
[
  {
    "title": "Improve test coverage",
    "gap": "The evaluation found a requirement with thin test evidence.",
    "evidenceLocators": [
      "assessments/001-root-has-tests.json"
    ],
    "assessmentResultRecords": [
      "assessments/001-root-has-tests.json"
    ],
    "remediationOptions": [
      "Add focused tests for the requirement"
    ],
    "recommendedOption": "Add focused tests for the requirement",
    "doneCriterion": "The affected requirement reaches target with current test evidence."
  }
]
JSON

qualitymd evaluation status <run>
qualitymd evaluation report build <run>
```

### Gate on report results

```sh
qualitymd evaluation report build <run>
qualitymd evaluation report gate <run> --at-or-below <level>
```

## Command rules

- Use `--json` when a command offers it and the agent must consume the result.
- Prefer `qualitymd status [path] --json` for readiness, model shape, evaluation
  history, stale-run signals, and active recommendation counts.
- Pipe record JSON on stdin (a `<<'JSON'` heredoc works well); never write the
  payload to a scratch file. `--file <path>` exists for human replay and
  debugging only — agents should not use it.
- Use `qualitymd evaluation <kind> add|set --help` to inspect payload fields and
  `-n/--dry-run` to validate new or materially revised payloads before writing.
- Do not continue past missing evaluation commands by manually creating files.
- Keep generated run paths exactly as the CLI reports them.
