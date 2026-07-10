# /quality runtime resources update log

## 2026-07-10

- **Revision**: Updated
  [CLI workflow conventions](cli-workflow-conventions.md) for 0195 - Evaluation
  concurrency. Dry-run preview guidance now names resolved concurrency rather
  than execution strategy.

## 2026-07-09

- **Revision**: Updated
  [CLI workflow conventions](cli-workflow-conventions.md) and
  [Output policy](output-policy.md) for 0192 - Deterministic evaluation runner.
  Evaluation sequencing now routes new runs through `qualitymd evaluation run`
  (with dry-run preview and `--resume` recovery), and the data
  discovery/write path is scoped to historical multi-file runs.

## 2026-06-29

- **Revision**: Updated runtime resource headings for 0189 - Heading Sentence
  Case. Active resource headings now use sentence case while preserving formal
  QUALITY.md concept names.

- **Release prep**: Bumped the bundled
  [QUALITY.md specification](SPECIFICATION.md) resource to `0.9 (Draft)` for the
  `v0.26.0` release.

## 2026-06-27

- **Revision**: Updated the bundled
  [QUALITY.md specification](SPECIFICATION.md) for 0142 - Requirement Findings
  Only.
  The runtime resource now treats Requirement Findings as the only Evaluation
  findings, removes Area-analysis findings from the current contract, and
  requires rated results to carry rating drivers.

## 2026-06-26

- **Revision**: Updated the bundled
  [QUALITY.md specification](SPECIFICATION.md) for 0131 - Area findings in
  evaluation reports.
  The runtime resource now recognizes Area Findings as analysis-level findings
  that summarize material observations for one Area and may relate to Factors in
  that Area, without carrying advice or ranking fields.

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
