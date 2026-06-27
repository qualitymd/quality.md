# Evaluation v2 Update Log

## 2026-06-27

- **Revision**: Updated the Evaluation records data-layout and report-tree
  contracts for [0154 - Ranked Findings Reports](../../changes/archive/0154-ranked-findings-reports.md).
  Evaluation report build now writes `findings.md` as the full ranked findings
  index, keeps `report.md` Top Findings capped, links ranked finding statements
  to stable Requirement detail anchors, and renders per-finding Advice rank
  context on Requirement reports.

- **Revision**: Updated the Evaluation data, protocol, orchestration, routine,
  data-layout, and report-tree contracts for
  [0150 - Evaluation Advice](../../changes/archive/0150-evaluation-advice.md).
  Evaluation now requires finding ranking, recommendations, finding coverage
  accounting, and recommendation ranking before report build. Reports render Top
  Findings, Top Recommendations, `recommendations.md`, and recommendation detail
  pages from persisted Advice data.

- **Revision**: Updated the Evaluation create, data, protocol, orchestration,
  and report-tree contracts for
  [0149 - Scope-driven evaluation runs](../../changes/archive/0149-scope-driven-evaluation-runs.md).
  Runs now carry CLI-owned `RunManifest` scope data, `EvaluationFrame` no longer
  carries run scope, and `report.md` renders as the scoped Area report without
  headline selection from agent-authored payload ordering.

- **Revision**: Updated the Evaluation data, routine, and report-tree contracts
  for [0148 - Finding Basis](../../changes/archive/0148-finding-basis.md).
  Requirement Findings now use `basis` instead of `cause` for the finding-local
  explanation or support posture, reports render `Basis` labels, and the
  existing support status values are preserved.

- **Revision**: Updated the Evaluation data, routine, protocol, and report-tree
  contracts for
  [0142 - Requirement Findings Only](../../changes/archive/0142-requirement-findings-only.md).
  Requirement Findings are now the only Evaluation findings, rated Requirements
  must be backed by paired Requirement Assessment findings and rating drivers,
  Factor/Area roll-ups use drivers and rationale instead of new findings, and
  the active Evaluation data schema is version 3.

- **Revision**: Updated the Evaluation report entrypoint contracts for
  [0137 - Run Report Entrypoint](../../changes/archive/0137-run-report-entrypoint.md).
  `report.md` is now the run-level report, the root Area detail report is
  `root-area.md` when present, scoped Area/Factor runs use recorded
  Evaluation Frame scope for headline reportability, and
  `EvaluationOutputResult` carries explicit run and headline refs.

- **Revision**: Updated the Evaluation candidate action contracts for
  [0136 - Candidate Actions Payload](../../changes/archive/0136-candidate-actions-payload.md).
  Requirement Findings now use `candidateActions` instead of `actions`, each
  candidate action has an ID local to its containing Finding, and reports still
  omit candidate actions in Evaluation v0.

- **Revision**: Updated the Evaluation finding contracts for
  [0135 - Structured Finding Core](../../changes/archive/0135-structured-finding-core.md).
  Requirement and Area Findings now use one shared Finding Core with statement,
  condition, criteria, cause, effect, and evidence; reports render Requirement,
  Area, and Factor Findings through one table/detail shape; and finding IDs are
  documented as payload-local selectors rather than durable cross-run IDs.

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
