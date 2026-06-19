# CLI Quick Reference

Use these commands and routing patterns as the starting point for `qualitymd`
work.

## Output modes

| Goal                                      | Command form                          |
| ----------------------------------------- | ------------------------------------- |
| Human-readable command help               | `qualitymd <command> --help`          |
| Machine-readable state or receipts        | `qualitymd <command> --json`          |
| Active format specification text          | `qualitymd spec`                      |
| JSON payload from stdin                   | `--file -` on commands that accept it |
| Version and compatibility facts for skill | `qualitymd version --json`            |

Use `--json` when a command offers it and the agent must inspect, route from,
or carry the result forward. Use human output for display or diagnostics only.

## CLI introspection

Use command help before guessing a command shape:

```sh
qualitymd --help
qualitymd lint --help
qualitymd status --help
qualitymd evaluation create-run --help
qualitymd evaluation add-record --help
qualitymd evaluation set-planned-coverage --help
qualitymd evaluation show-status --help
qualitymd evaluation build-report --help
```

Use `qualitymd version --json` to inspect the visible CLI version,
development-build state, commit when known, bundled specification version, and
whether the install is in the skill's supported range.

## Quick reference

| Task                      | Command                                                                   |
| ------------------------- | ------------------------------------------------------------------------- |
| Check CLI version         | `qualitymd --version`                                                     |
| Read version metadata     | `qualitymd version --json`                                                |
| Check for CLI upgrades    | `qualitymd upgrade --check`                                               |
| Apply CLI upgrade         | `qualitymd upgrade --apply`                                               |
| Read format rules         | `qualitymd spec`                                                          |
| Create a starter model    | `qualitymd init [path]`                                                   |
| Validate a model          | `qualitymd lint [path]`                                                   |
| Fix simple lint issues    | `qualitymd lint --fix [path]`                                             |
| Inspect project status    | `qualitymd status [path] --json`                                          |
| Create evaluation run     | `qualitymd evaluation create-run [--subject <path>] [--narrowing <slug>]` |
| Add assessment record     | `qualitymd evaluation add-record assessment <run> --file <path-or-->`     |
| Add analysis record       | `qualitymd evaluation add-record analysis <run> --file <path-or-->`       |
| Add recommendation record | `qualitymd evaluation add-record recommendation <run> --file <path-or-->` |
| Set planned coverage      | `qualitymd evaluation set-planned-coverage <run> --file <path-or-->`      |
| Check reportability       | `qualitymd evaluation show-status <run>`                                  |
| Build report              | `qualitymd evaluation build-report <run>`                                 |

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
- Create run -> qualitymd evaluation create-run [--subject <path>] [--narrowing <slug>]
- Add judgment records -> qualitymd evaluation add-record assessment|analysis|recommendation <run> --file <path-or-->
- Ready to report? -> qualitymd evaluation show-status <run>
- Build report -> qualitymd evaluation build-report <run>
```

### Resuming or diagnosing a run

```text
Run incomplete or stale?
- Inspect project state -> qualitymd status [path] --json
- Inspect run readiness -> qualitymd evaluation show-status <run>
- Missing planned coverage? -> qualitymd evaluation set-planned-coverage <run> --file <path-or-->
- Missing records? -> qualitymd evaluation add-record assessment|analysis|recommendation <run> --file <path-or-->
- Reportable? -> qualitymd evaluation build-report <run>
```

### Upgrading

```text
CLI missing, stale, or incompatible?
- Inspect visible CLI -> qualitymd version --json
- Check install-aware action -> qualitymd upgrade --check
- Apply only when confirmed -> qualitymd upgrade --apply
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
qualitymd evaluation create-run [--subject <path>] [--narrowing <slug>]
qualitymd evaluation add-record assessment <run> --file <path-or-->
qualitymd evaluation add-record analysis <run> --file <path-or-->
qualitymd evaluation add-record recommendation <run> --file <path-or-->
qualitymd evaluation show-status <run>
qualitymd evaluation build-report <run>
```

### Gate on report results

```sh
qualitymd evaluation build-report <run> --fail-at-or-below <level>
```

## Command rules

- Use `--json` when a command offers it and the agent must consume the result.
- Prefer `qualitymd status [path] --json` for readiness, model shape, evaluation
  history, stale-run signals, and active recommendation counts.
- Use `--file -` when passing JSON through stdin.
- Do not continue past missing evaluation commands by manually creating files.
- Keep generated run paths exactly as the CLI reports them.
