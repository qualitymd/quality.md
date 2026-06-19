# CLI Quick Reference

Use these commands as the starting point for `qualitymd` work.

| Task                      | Command                                                                   |
| ------------------------- | ------------------------------------------------------------------------- |
| Check CLI version         | `qualitymd --version`                                                     |
| Read version metadata     | `qualitymd version --json`                                                |
| Check for CLI upgrades    | `qualitymd upgrade --check`                                               |
| Read format rules         | `qualitymd spec`                                                          |
| Create a starter model    | `qualitymd init [path]`                                                   |
| Validate a model          | `qualitymd lint [path]`                                                   |
| Inspect project status    | `qualitymd status [path] --json`                                          |
| Create evaluation run     | `qualitymd evaluation create-run [--subject <path>] [--narrowing <slug>]` |
| Add assessment record     | `qualitymd evaluation add-record assessment <run> --file <path-or-->`     |
| Add analysis record       | `qualitymd evaluation add-record analysis <run> --file <path-or-->`       |
| Add recommendation record | `qualitymd evaluation add-record recommendation <run> --file <path-or-->` |
| Set planned coverage      | `qualitymd evaluation set-planned-coverage <run> --file <path-or-->`      |
| Check reportability       | `qualitymd evaluation show-status <run>`                                  |
| Build report              | `qualitymd evaluation build-report <run>`                                 |

## Command rules

- Use `--json` when a command offers it and the agent must consume the result.
- Prefer `qualitymd status [path] --json` for readiness, model shape, evaluation
  history, stale-run signals, and active recommendation counts.
- Use `--file -` when passing JSON through stdin.
- Do not continue past missing evaluation commands by manually creating files.
- Keep generated run paths exactly as the CLI reports them.
