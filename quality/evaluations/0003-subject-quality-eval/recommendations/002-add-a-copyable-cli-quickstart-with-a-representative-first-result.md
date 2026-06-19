---
schemaVersion: 1
title: Add a copyable CLI quickstart with a representative first result
gap: The README's Usage section routes a newcomer to the /quality skill rather than a short install-then-one-command CLI sequence that yields a visible result. The deterministic init then lint flow works end-to-end but is never shown as a quickstart, and no representative output is shown.
evidenceLocators:
  - README.md:20-31
  - README.md:194
  - README.md:204-205
assessmentRecords:
  - quality/evaluations/0003-subject-quality-eval/assessments/012-readme-the-readme-gets-a-newcomer-to-a-first-result-quickly.json
remediationOptions:
  - 'Add a ''Quickstart'' block: install, then ''qualitymd init'' and ''qualitymd lint QUALITY.md'', with the actual printed output (the init ''Created ... Next:'' line and ''QUALITY.md is valid.'') and the CI exit-code note.'
  - Point to the same flow but only describe it in prose without showing output.
  - Defer until the npm/npx packages are published and show the published-binary path.
recommendedOption: 'Option 1: a real install-then-one-command quickstart with the printed output is exactly what the requirement asks for and the underlying flow already works, so it is low-risk and high-value.'
doneCriterion: 'The requirement ''the README gets a newcomer to a first result quickly'' reaches at least target: a copyable install-then-one-command sequence with representative output appears early in the README. Re-evaluate in a new numbered run.'
---

# Add a copyable CLI quickstart with a representative first result

## Gap

The README's Usage section routes a newcomer to the /quality skill rather than a short install-then-one-command CLI sequence that yields a visible result. The deterministic init then lint flow works end-to-end but is never shown as a quickstart, and no representative output is shown.

## Evidence locators

- `README.md:20-31`
- `README.md:194`
- `README.md:204-205`

## Remediation options

- Add a 'Quickstart' block: install, then 'qualitymd init' and 'qualitymd lint QUALITY.md', with the actual printed output (the init 'Created ... Next:' line and 'QUALITY.md is valid.') and the CI exit-code note.
- Point to the same flow but only describe it in prose without showing output.
- Defer until the npm/npx packages are published and show the published-binary path.

## Recommended option

Option 1: a real install-then-one-command quickstart with the printed output is exactly what the requirement asks for and the underlying flow already works, so it is low-risk and high-value.

## Done criterion

The requirement 'the README gets a newcomer to a first result quickly' reaches at least target: a copyable install-then-one-command sequence with representative output appears early in the README. Re-evaluate in a new numbered run.
