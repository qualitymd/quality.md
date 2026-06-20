# Quality Evaluation Summary

| Field          | Value                       |
| -------------- | --------------------------- |
| Subject        | QUALITY.md                  |
| Run            | `0005-subject-quality-eval` |
| Scope          | Full evaluation             |
| Rigor          | Deep                        |
| Overall rating | 🔵 Target                   |
| Full report    | [report.md](report.md)      |
| Machine report | [report.json](report.json)  |

## Summary

The root is a grouping target. All three child deliverables now aggregate to target after the spec, README, CLI, and evaluation-history fixes, so the whole project meets the release quality bar.

| Target               | Local rating     | Overall rating | Driver                                                                                                                                                                                             |
| -------------------- | ---------------- | -------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| QUALITY.md           | n/a (structural) | 🔵 Target      | The root is a grouping target. All three child deliverables now aggregate to target after the spec, README, CLI, and evaluation-history fixes, so the whole project meets the release quality bar. |
| qualitymd CLI        | 🔵 Target        | 🔵 Target      | Leaf target: aggregate equals local rating. The release-relevant CLI conformance gaps found in this run are fixed and verified.                                                                    |
| Format specification | 🔵 Target        | 🔵 Target      | Leaf target: aggregate equals local rating. The format spec satisfies every in-scope requirement at target or better.                                                                              |
| README               | 🔵 Target        | 🔵 Target      | Leaf target: aggregate equals local rating. The front door now explains and demonstrates the tool adequately for release.                                                                          |

## Top Issues

1. **low**\
   Scalar placeholders such as <string> and <level-name> intentionally do not define detailed character-set or length bounds, leaving edge cases to conforming tools.
   `SPECIFICATION.md:137`
   Assessment: `assessment-results/001-format-spec-the-format-specification-is-complete.json`
2. **low**\
   The not-assessed boundary still depends on evaluator judgment about evidence sufficiency, but the report distinction is explicit.
   `SPECIFICATION.md:344`
   Assessment: `assessment-results/002-format-spec-the-format-specification-admits-a-single-interpretation.json`
3. **low**\
   Roll-up intentionally has no numeric aggregation formula, so exact rating inference remains evaluator judgment constrained by required distinctions.
   `SPECIFICATION.md:356`
   Assessment: `assessment-results/006-format-spec-each-rule-is-observable-or-testable.json`

## Recommendations

No recommendation records exist for this run.

## Scope & Limitations

Scope: **Full evaluation**

In scope: QUALITY.md; qualitymd CLI; Format specification; README

- The spec now includes valid minimal/suggested examples and explicit invalid counter-examples for missing required title, missing direct-requirement factors, and list-valued assessment
