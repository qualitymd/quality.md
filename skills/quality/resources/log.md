# /quality Runtime Resources Update Log

## 2026-06-26

- **Revision**: Updated
  [CLI Workflow Conventions](cli-workflow-conventions.md) for 0130 -
  Self-contained per-kind data schema.
  Evaluation payload discovery now states that per-kind schemas are
  self-contained and authoritative for required fields and enum values, while
  examples are concrete valid instances.

- **Revision**: Updated
  [CLI Workflow Conventions](cli-workflow-conventions.md) for 0128 -
  Agent-mediated skill alignment.
  Setup feedback-log timing now matches the setup workflow: create the log after
  the setup preview when the run continues into discovery.

- **Refocus**: Renamed `cli-quick-reference.md` to
  `cli-workflow-conventions.md`, removing the embedded command/flag listing and
  retaining the workflow conventions the CLI's own introspection cannot carry.

## 2026-06-24

- **Restructure**: Added this resource index/log as part of the runtime skill
  OKF shape.
