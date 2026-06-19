---
schemaVersion: 1
title: Get a newcomer to a visible first result with output and CI exit codes
gap: The README gives a copyable install-then-command sequence but shows no representative output for any command and mentions CI exit-code behavior only in passing, so a newcomer runs commands without knowing what success looks like.
evidenceLocators:
  - README.md:22-35
  - README.md:223-230
  - README.md:207-208
assessmentRecords:
  - quality/evaluations/0004-subject-quality-eval/assessments/012-readme-the-readme-gets-a-newcomer-to-a-first-result-quickly.json
remediationOptions:
  - Annotate the typical local loop with representative output for at least lint and status, and add a one-line CI snippet showing the gate exit code (build-report --fail-at-or-below).
  - Add a dedicated 'first result' walkthrough section with expected output blocks.
  - Add only the CI exit-code example, deferring per-command output.
recommendedOption: 'Option 1: annotate the existing local loop with output and add a short CI exit-code example; it reuses the sequence already present and closes both missing pieces at once.'
doneCriterion: 'The requirement ''the README gets a newcomer to a first result quickly'' reaches at least target: representative output and CI exit-code behavior are shown. Re-evaluate in a new numbered run.'
---

# Get a newcomer to a visible first result with output and CI exit codes

## Gap

The README gives a copyable install-then-command sequence but shows no representative output for any command and mentions CI exit-code behavior only in passing, so a newcomer runs commands without knowing what success looks like.

## Evidence locators

- `README.md:22-35`
- `README.md:223-230`
- `README.md:207-208`

## Remediation options

- Annotate the typical local loop with representative output for at least lint and status, and add a one-line CI snippet showing the gate exit code (build-report --fail-at-or-below).
- Add a dedicated 'first result' walkthrough section with expected output blocks.
- Add only the CI exit-code example, deferring per-command output.

## Recommended option

Option 1: annotate the existing local loop with output and add a short CI exit-code example; it reuses the sequence already present and closes both missing pieces at once.

## Done criterion

The requirement 'the README gets a newcomer to a first result quickly' reaches at least target: representative output and CI exit-code behavior are shown. Re-evaluate in a new numbered run.
