# Evaluation v2 Update Log

## 2026-06-26

- **Revision**: Updated the payload-kind and report-tree contracts for
  [0132 - Remove info finding severity](../../changes/archive/0132-remove-info-finding-severity.md).
  Evaluation finding severity no longer includes `info`; report severity sorting
  now orders only `critical`, `high`, `medium`, and `low`.

- **Revision**: Updated the routine, payload-kind, report-tree, and CLI data
  contracts for
  [0131 - Area findings in evaluation reports](../../changes/archive/0131-area-findings.md).
  `AreaAnalysisResult.findings` now carries traceable Area Findings with closed
  type/severity/confidence and Factor relationship enums; reports render Area
  Findings in Area pages and matching Findings in Factor pages without adding
  recommendations or global rankings.

- **Revision**: Updated the routine and report-tree contracts for
  [0122 - Finding-level candidate actions](../../changes/archive/0122-finding-candidate-actions.md).
  Requirement assessment MAY record non-binding, finding-local candidate actions
  and MUST NOT synthesize, aggregate, or prioritize them; finding detail sections
  MUST NOT render candidate actions in Evaluation v0.

- **Revision**: Updated the records data-layout and report-tree contracts for
  [0108 - Short evaluation report filenames](../../changes/archive/0108-short-evaluation-report-filenames.md).
  The root Area report remains `report.md`, while descendant Area, Factor, and
  Requirement reports use short subject-aware filenames derived from structural
  model IDs.

## 2026-06-25

- **Creation**: Added the Evaluation v2 durable spec bundle for
  [0094 - Replace evaluation workflow](../../changes/0094-replace-evaluation-workflow.md),
  including shared workflow invariants, protocol, orchestration, routine,
  record, and report contracts.
