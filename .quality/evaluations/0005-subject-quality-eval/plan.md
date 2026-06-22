---
coverage:
  assessmentResults:
    - targetPath: [format-spec]
      requirement: the format specification is complete
    - targetPath: [format-spec]
      requirement: the format specification admits a single interpretation
    - targetPath: [format-spec]
      requirement: the format specification separates rules from rationale
    - targetPath: [format-spec]
      requirement: the format specification defines its terms before use
    - targetPath: [format-spec]
      requirement: the format specification is internally consistent
    - targetPath: [format-spec]
      requirement: each rule is observable or testable
    - targetPath: [format-spec]
      requirement: the format's constructs are shown with valid and invalid examples
    - targetPath: [format-spec]
      requirement: the format specifies its core and how it extends and evolves
    - targetPath: [format-spec]
      requirement: the format specification is well-structured and readable
    - targetPath: [readme]
      requirement: the README says what QUALITY.md is and who it's for
    - targetPath: [readme]
      requirement: the README shows the format and its payoff by example
    - targetPath: [readme]
      requirement: the README gets a newcomer to a first result quickly
    - targetPath: [readme]
      requirement: the README reflects what the CLI and spec actually provide
    - targetPath: [cli]
      requirement: the CLI follows its functional specifications
    - targetPath: [cli]
      requirement: the CLI follows the project CLI design guide
  analyses:
    - targetPath: []
    - targetPath: [format-spec]
    - targetPath: [readme]
    - targetPath: [cli]
---

# Evaluation plan

## Rigor

Deep.

## Requirement set

Cover all 15 requirements in the current root model:

- Format specification: completeness; single interpretation; rules vs
  rationale; terms before use; internal consistency; observability/testability;
  valid and invalid examples; core/evolution; structure/readability.
- README: what/who; format plus payoff; first result; current CLI/spec accuracy.
- CLI: functional specification conformance; project CLI design-guide
  conformance.

## Procedure

1. Verify model validity and inspect evaluation history.
2. Read the full in-scope source documents and relevant CLI implementation.
3. Execute representative CLI behavior checks for rating-binding findings.
4. Fix confirmed release-relevant gaps.
5. Re-run focused tests and behavior checks.
6. Record assessment and analysis records through the CLI.
7. Build the report and use it as the final subject-quality evidence for the
   release gate.

## Deferred areas

- Hosted CI and publication checks are deferred to the release runbook.
- Exhaustive CLI matrix testing is deferred; focused behavior checks and package
  tests cover the release-relevant regressions found in this run.
