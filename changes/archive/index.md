# Archived changes

Completed changes, moved here from the bundle root when they reach **Done**.

# Changes

- [0001 — Example change](0001-example-change.md) - placeholder retired as a
  reference template for the Change concept shape (`Done`).
- [0002 — Specify the init command](0002-init-command.md) - settled and shipped
  `qualitymd init` (`Done`).
- [0003 — Implement the lint command](0003-implement-lint-command.md) - built
  `qualitymd lint` from the completed durable lint sub-spec (`Done`).
- [0004 — Specify and enforce agent accessibility](0004-specify-agent-accessibility.md) - added the CLI agent-accessibility contract, broadened `--json`, and enforced categorized exit codes plus `init --json` (`Done`).
- [0005 — Single source of truth for the structural schema](0005-schema-source-of-truth.md) - extracted the structural schema into one typed declaration consumed by `lint` (`Done`).
- [0006 — Specify and implement the spec command](0006-spec-command.md) -
  settled and shipped `qualitymd spec`, emitting the bundled format specification
  from the binary (`Done`).
- [0007 — Delightful human CLI output](0007-delightful-cli-output.md) - gave the
  human surface a shared brand palette, styled `lint` and `init` output, `--help`
  examples, `spec` paging, and an informative `--version`, all behind the
  TTY/`NO_COLOR` gate so the plain and JSON paths are untouched (`Done`).
- [0008 — Describe targets with title and description](0008-target-display-fields.md) -
  lets every target carry a recommended `title` and optional `description`, and
  reframes the root as a Model (`ratingScale` + Target properties) so
  `ratingScale` is the one Model-only key (`Done`).
- [0009 — Diagnose rating-scale soundness in the meta-model](0009-rating-scale-diagnostic.md) -
  adds a meta-model Functionality requirement that judges a model's rating scale
  and per-requirement criterion overrides for meaning, not only structure
  (`Done`).
