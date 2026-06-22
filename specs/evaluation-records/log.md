# Evaluation Records Update Log

## 2026-06-22

- **Revision**: Added the [design.md](design-md.md) artifact contract for
  [0056 - Prospective evaluation plan artifacts](../../changes/archive/0056-prospective-evaluation-plan-artifacts.md)
  and clarified that the initial [plan.md](plan-md.md) is prospective execution
  planning. `plan.md` updates now preserve the original plan through explicit
  amendments, while actual findings, rating rationale, and recommendation
  reasoning stay in formal records and reports.

- **Creation**: Added child specs for independently reviewable evaluation-record
  runtime contracts: [run folder](run-folder.md),
  [assessment result record](assessment-result-record.md),
  [analysis record](analysis-record.md), [plan.md](plan-md.md),
  [debug-log.md](debug-log-md.md),
  [recommendation record](recommendation-record-md.md), and
  [report outputs](report-outputs.md). The parent
  [Evaluation records](../evaluation-records.md) spec keeps shared
  responsibility, runtime-not-OKF, schema-version, and historical compatibility
  rules.
