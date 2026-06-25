---
schemaVersion: 1
title: Reconcile stale evaluation history
gap: The previous evaluation run remains inspectable but is not reportable under the current record contract because historical records and plan coverage lack areaPath fields.
evidenceLocators:
  - qualitymd evaluation status .quality/evaluations/0005-subject-quality-eval
  - assessments/030-evaluation-history-evaluation-history-and-model-change-history-remain-inspectable-and-distinct.json
assessmentResultRecords:
  - assessments/030-evaluation-history-evaluation-history-and-model-change-history-remain-inspectable-and-distinct.json
remediationOptions:
  - Leave the stale run as historical context and document that pre-0.11 runs may be non-reportable
  - Add an explicit archive or migration path for stale pre-areaPath evaluation records
recommendedOption: Add an explicit archive or migration path for stale pre-areaPath evaluation records
doneCriterion: qualitymd status QUALITY.md --json no longer reports a stale non-reportable prior run, or the repository has an explicit documented archive state that keeps stale runs out of current-readiness gaps.
---

# Reconcile stale evaluation history

## Gap

The previous evaluation run remains inspectable but is not reportable under the current record contract because historical records and plan coverage lack areaPath fields.

## Evidence locators

- `qualitymd evaluation status .quality/evaluations/0005-subject-quality-eval`
- `assessments/030-evaluation-history-evaluation-history-and-model-change-history-remain-inspectable-and-distinct.json`

## Remediation options

- Leave the stale run as historical context and document that pre-0.11 runs may be non-reportable
- Add an explicit archive or migration path for stale pre-areaPath evaluation records

## Recommended option

Add an explicit archive or migration path for stale pre-areaPath evaluation records

## Done criterion

qualitymd status QUALITY.md --json no longer reports a stale non-reportable prior run, or the repository has an explicit documented archive state that keeps stale runs out of current-readiness gaps.
