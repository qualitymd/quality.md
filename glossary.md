# Glossary

This glossary defines shared QUALITY.md terms, concepts, and fixed
vocabularies.

For fixed vocabulary tables, `Value` is the canonical persisted value. `Label`
is the human display label. Labels, markers, aliases, and case variants are not
accepted as structured data values.

## Analysis status

Whether an Area or Factor analysis has been completed, skipped as empty, not
analyzed, or blocked.

| Label           | Value          | Description                                |
| --------------- | -------------- | ------------------------------------------ |
| ✅ Analyzed     | `analyzed`     | Analysis completed with a recorded result. |
| ⬜ Empty        | `empty`        | Scope had nothing applicable to analyze.   |
| ⚪ Not Analyzed | `not_analyzed` | Analysis has not been performed.           |
| ⛔ Blocked      | `blocked`      | Analysis could not be completed.           |

## Area

An entity or set of entities with quality requirements subject to evaluation.

## Assessment status

Whether a Requirement assessment has been completed, partially completed, not
assessed, or blocked.

| Label                 | Value                | Description                                            |
| --------------------- | -------------------- | ------------------------------------------------------ |
| ✅ Assessed           | `assessed`           | Requirement was assessed with usable findings.         |
| 🟡 Partially Assessed | `partially_assessed` | Assessment is incomplete but contains usable judgment. |
| ⚪ Not Assessed       | `not_assessed`       | No assessment was completed.                           |
| ⛔ Blocked            | `blocked`            | Assessment could not be completed.                     |

## Confidence

Confidence in the recorded judgment based on available evidence.

| Label     | Value    | Description                                                    |
| --------- | -------- | -------------------------------------------------------------- |
| 🟢 High   | `high`   | Strong evidence supports the judgment.                         |
| 🔵 Medium | `medium` | Adequate evidence supports the judgment with some uncertainty. |
| 🟡 Low    | `low`    | Limited evidence supports the judgment.                        |
| ⚪ None   | `none`   | No meaningful confidence is available.                         |

## Data kind

Kind of structured Evaluation payload stored for a run.

| Label                           | Value                         | Description                                       |
| ------------------------------- | ----------------------------- | ------------------------------------------------- |
| 📋 Evaluation Manifest          | `EvaluationManifest`          | Evaluation metadata, scope, and run context.      |
| 🧭 Evaluation Frame             | `EvaluationFrame`             | Top-level evaluation planning frame.              |
| 🗺️ Area Evaluation Frame         | `AreaEvaluationFrame`         | Planned evaluation frame for an Area.             |
| 📋 Requirement Evaluation Frame | `RequirementEvaluationFrame`  | Planned assessment frame for a Requirement.       |
| 🔎 Requirement Assessment       | `RequirementAssessmentResult` | Judgment evidence and findings for a Requirement. |
| 🎚️ Requirement Rating            | `RequirementRatingResult`     | Rating assigned to a Requirement.                 |
| 🧩 Factor Analysis Frame        | `FactorAnalysisFrame`         | Planned analysis frame for a Factor.              |
| 📊 Factor Analysis              | `FactorAnalysisResult`        | Synthesized judgment for a Factor.                |
| 🏗️ Area Analysis Frame           | `AreaAnalysisFrame`           | Planned analysis frame for an Area.               |
| 📈 Area Analysis                | `AreaAnalysisResult`          | Synthesized judgment for an Area.                 |
| 🔝 Finding Ranking              | `FindingRankingResult`        | Ordered finding priority set.                     |
| 💡 Recommendation               | `RecommendationResult`        | Proposed improvement action.                      |
| 🏁 Recommendation Ranking       | `RecommendationRankingResult` | Ordered recommendation priority set.              |
| 📦 Evaluation Output            | `EvaluationOutputResult`      | Generated report-output index.                    |

## Factor

A quality characteristic or attribute through which an Area's quality is
described. A Factor groups connected Requirements and can be decomposed into
sub-factors.

## Finding

A single observation produced by a Requirement Assessment. A Finding records
what was observed, the criteria applied, its evidence basis, and its quality or
rating effect.

## Finding basis

Evidence support state for a finding's basis.

| Label             | Value            | Description                                   |
| ----------------- | ---------------- | --------------------------------------------- |
| ✅ Verified       | `verified`       | Basis is directly supported by evidence.      |
| 🟡 Plausible      | `plausible`      | Basis is reasonable but not fully verified.   |
| ⚪ Not Assessed   | `not_assessed`   | Basis support was not assessed.               |
| ⬜ Not Applicable | `not_applicable` | Basis support does not apply to this finding. |

## Finding coverage

Whether a finding is addressed by recommendations or intentionally not
advice-driving.

| Label                          | Value                         | Description                                               |
| ------------------------------ | ----------------------------- | --------------------------------------------------------- |
| ✅ Addressed by Recommendation | `addressed_by_recommendation` | Finding is covered by one or more recommendations.        |
| ⬜ Not Advice Driving          | `not_advice_driving`          | Finding is intentionally not driving recommendation work. |

## Finding rank

Priority tier assigned when ranking findings for attention.

| Label         | Value | Description                         |
| ------------- | ----- | ----------------------------------- |
| 🔴 P1 Highest | `P1`  | Top-priority finding for attention. |
| 🟠 P2 High    | `P2`  | High-priority finding.              |
| 🟡 P3 Medium  | `P3`  | Medium-priority finding.            |
| ⚪ P4 Low     | `P4`  | Lower-priority finding.             |

## Finding severity

Severity of the finding's quality concern or evaluation significance.

| Label       | Value      | Description                                        |
| ----------- | ---------- | -------------------------------------------------- |
| 🔴 Critical | `critical` | Severe concern requiring urgent attention.         |
| 🔴 High     | `high`     | Important concern with substantial quality impact. |
| 🟡 Medium   | `medium`   | Meaningful concern worth addressing.               |
| 🔵 Low      | `low`      | Minor concern or low-impact observation.           |

## Finding type

Classification of what a finding contributes to the Evaluation judgment.

| Label       | Value      | Description                                                     |
| ----------- | ---------- | --------------------------------------------------------------- |
| 🚩 Gap      | `gap`      | Current shortfall against the quality bar.                      |
| ⚠️ Risk      | `risk`     | Plausible future or conditional quality concern.                |
| ✅ Strength | `strength` | Evidence of quality meeting or exceeding expectations.          |
| ℹ️ Note      | `note`     | Useful observation that is not itself a gap, risk, or strength. |

## Quality rating

A quality rating is the Rating Level assigned to evaluated work. Rating Levels
are configured by the quality model, not by a fixed Evaluation enum.

These labels and values come from this project's `QUALITY.md` Rating Scale.

| Label           | Value          | Description                                                                                    |
| --------------- | -------------- | ---------------------------------------------------------------------------------------------- |
| 🟢 Outstanding  | `outstanding`  | The stretch band: the artifact exceeds the quality requirement with meaningful margin.         |
| 🔵 Target       | `target`       | The expected good state: the artifact satisfies the quality requirement.                       |
| 🟡 Minimum      | `minimum`      | The acceptable floor: the artifact falls short of the goal but remains good enough to rely on. |
| 🔴 Unacceptable | `unacceptable` | Below the floor: the artifact is not good enough to rely on.                                   |

## Rating result

Whether a Rating Result contains an assigned rating or records that the subject
was not assessed.

| Label           | Value          | Description                             |
| --------------- | -------------- | --------------------------------------- |
| ✅ Rated        | `rated`        | The result contains an assigned rating. |
| ⚪ Not Assessed | `not_assessed` | The subject was not assessed.           |

## Rating status

Whether a rating result has been assigned, not assigned, or blocked.

| Label        | Value       | Description                   |
| ------------ | ----------- | ----------------------------- |
| ✅ Rated     | `rated`     | A rating level was assigned.  |
| ⚪ Not Rated | `not_rated` | No rating level was assigned. |
| ⛔ Blocked   | `blocked`   | Rating could not be assigned. |

## Recommendation

A proposed improvement action produced from Evaluation findings and judgment.

## Recommendation impact

Expected quality improvement from completing a recommendation.

| Label        | Value       | Description                                                |
| ------------ | ----------- | ---------------------------------------------------------- |
| ⬥⬥ Very high | `very_high` | Expected to materially improve important quality outcomes. |
| ⬥ High       | `high`      | Expected to meaningfully improve quality.                  |
| ● Medium     | `medium`    | Expected to provide useful but bounded improvement.        |
| ○ Low        | `low`       | Expected to provide small or localized improvement.        |

## Report kind

Kind of generated Markdown report artifact.

| Label              | Value             | Description                   |
| ------------------ | ----------------- | ----------------------------- |
| 📄 Run             | `run`             | Run entrypoint report.        |
| 🗺️ Area             | `area`            | Area report.                  |
| 🧩 Factor          | `factor`          | Factor report.                |
| 📋 Requirement     | `requirement`     | Requirement report.           |
| 🔝 Findings        | `findings`        | Findings index report.        |
| 📚 Recommendations | `recommendations` | Recommendations index report. |
| 💡 Recommendation  | `recommendation`  | Recommendation detail report. |

## Requirement

An assessable quality expectation. A Requirement has a stable Requirement name,
a title, an Assessment, zero or more explicit Factor references, and optional
per-level criterion overrides.

## Run gap kind

Kind of missing, unreadable, malformed, or incomplete run data blocking
reportability.

| Label                         | Value                        | Description                                            |
| ----------------------------- | ---------------------------- | ------------------------------------------------------ |
| 📭 Missing Evaluation Data    | `missing-evaluation-data`    | Required payload is absent.                            |
| ⚠️ Malformed Evaluation Data   | `malformed-evaluation-data`  | Payload cannot be parsed or has the wrong structure.   |
| 🚫 Unreadable Evaluation Data | `unreadable-evaluation-data` | Payload exists but cannot be read.                     |
| 🧩 Incomplete Evaluation Data | `incomplete-evaluation-data` | Payload is readable but lacks required usable content. |
