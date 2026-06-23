# Evaluation Records Update Log

## 2026-06-23

- **Revision**: Updated the [run folder](run-folder.md),
  [debug-log.md legacy compatibility](debug-log-md.md), and parent
  [Evaluation records](../evaluation-records.md) specs for
  [0073 - Evaluation feedback log](../../changes/archive/0073-evaluation-feedback-log.md).
  New evaluation runs no longer seed `debug-log.md`; historical runs that contain
  it remain tolerated as non-authoritative legacy process artifacts.

## 2026-06-22

- **Revision**: Updated [report outputs](report-outputs.md) for
  [0060 - Friendly path display](../../changes/archive/0060-friendly-path-display.md) so
  human Markdown reports may render display values such as `/` for the root Area
  while `report.json` continues to preserve structured path arrays and
  reference-compatible identifiers.

- **Revision**: Updated [report outputs](report-outputs.md) for
  [0059 - Unqualified model references](../../changes/archive/0059-unqualified-model-references.md)
  so fixed-type human report contexts may use unqualified references while
  `report.json` continues to preserve structured `areaPath` and `factorPath`
  arrays.

- **Revision**: Updated [report outputs](report-outputs.md) for
  [0058 - Model reference identifiers](../../changes/archive/0058-model-reference-identifiers.md)
  so human reports use canonical model references where a stable handle is
  displayed while `report.json` preserves structured `areaPath` and `factorPath`
  arrays.

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
