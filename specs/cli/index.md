# qualitymd CLI commands

The format-tooling commands specified for this phase. See the parent
[CLI spec](../cli.md) for the cross-cutting contract they all share.

# Commands

- [qualitymd init](init.md) - scaffold a starter `QUALITY.md` to fill in.
- [qualitymd lint](lint.md) - validate a file's structure against the format spec.
- [qualitymd spec](spec.md) - emit the QUALITY.md format specification.
- [qualitymd schema](schema.md) - emit the companion JSON Schema for QUALITY.md frontmatter.
- [qualitymd status](status.md) - emit a deterministic project-state snapshot.
- [qualitymd evaluation create](evaluation-create.md) - create a numbered
  evaluation run folder.
- [qualitymd evaluation data](evaluation-data.md) - persist and inspect
  Evaluation v2 structured data.
- [qualitymd evaluation list](evaluation-list.md) - list evaluation runs.
- [qualitymd evaluation status](evaluation-status.md) - inspect whether a run can
  be rendered.
- [qualitymd evaluation report](evaluation-report.md) - build evaluation reports.
- [qualitymd version](version.md) - show structured CLI and bundled spec version
  metadata.
- [qualitymd update](update.md) - apply or check for CLI updates through managed
  install channels.

# Components

- [qualitymd lint rules](lint-rules.md) - rule-system, rule-authoring, and rule
  catalog contract for `qualitymd lint`.
- [qualitymd lint output](lint-output.md) - finding, JSON, repair-result, and
  human-output contract for `qualitymd lint`.
- [qualitymd update notice](update-notice.md) - cross-command ambient update
  notice contract.
