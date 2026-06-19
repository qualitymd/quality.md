# qualitymd CLI commands

The format-tooling commands specified for this phase. See the parent
[CLI spec](../cli.md) for the cross-cutting contract they all share.

# Commands

- [qualitymd init](init.md) - scaffold a starter `QUALITY.md` to fill in.
- [qualitymd lint](lint.md) - validate a file's structure against the format spec.
- [qualitymd spec](spec.md) - emit the `QUALITY.md` format specification.
- [qualitymd status](status.md) - emit a deterministic project-state snapshot.
- [qualitymd evaluation create-run](evaluation-create-run.md) - create a
  numbered evaluation run folder.
- [qualitymd evaluation add-record](evaluation-add-record.md) - write
  assessment, analysis, and recommendation records.
- [qualitymd evaluation set-planned-coverage](evaluation-set-planned-coverage.md)
  - write planned coverage metadata for a run.
- [qualitymd evaluation show-status](evaluation-show-status.md) - inspect whether
  a run can be rendered.
- [qualitymd evaluation build-report](evaluation-build-report.md) - derive
  `report.md` and `report.json`, with an optional CI gate.
