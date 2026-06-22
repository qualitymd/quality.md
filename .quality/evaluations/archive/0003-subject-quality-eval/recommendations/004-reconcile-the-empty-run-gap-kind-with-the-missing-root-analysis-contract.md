---
schemaVersion: 1
title: Reconcile the empty-run gap kind with the missing-root-analysis contract
gap: A run with zero analysis records reports gap kind 'missing-analysis', but specs/cli/evaluation-show-status.md:29-32 mandates a 'missing-root-analysis' gap when no analysis record has an empty targetPath. The empty run has no root analysis record, so a caller branching on the documented gap kind misses the most common (fresh-run) case.
evidenceLocators:
  - specs/cli/evaluation-show-status.md:29-32
  - qualitymd evaluation show-status <empty-run> --json
assessmentRecords:
  - quality/evaluations/0003-subject-quality-eval/assessments/014-cli-the-cli-follows-its-functional-specifications.json
remediationOptions:
  - Emit 'missing-root-analysis' for the zero-analysis case so the implementation matches the spec's single documented gap kind.
  - Keep 'missing-analysis' for the zero case but amend the spec to document both kinds and when each applies.
  - Emit both gaps for the zero case (superset) so callers branching on either kind succeed.
recommendedOption: 'Option 1: align the implementation to the spec''s ''missing-root-analysis'' so a single documented gap kind covers the empty-run case; it is the smallest change that removes the divergence. Option 2 is acceptable if the team prefers to distinguish zero-records from records-without-root.'
doneCriterion: An empty run's show-status gap kind matches the spec exactly (either the CLI emits 'missing-root-analysis' or the spec documents 'missing-analysis' for the zero case). Re-evaluate in a new numbered run.
---

# Reconcile the empty-run gap kind with the missing-root-analysis contract

## Gap

A run with zero analysis records reports gap kind 'missing-analysis', but specs/cli/evaluation-show-status.md:29-32 mandates a 'missing-root-analysis' gap when no analysis record has an empty targetPath. The empty run has no root analysis record, so a caller branching on the documented gap kind misses the most common (fresh-run) case.

## Evidence locators

- `specs/cli/evaluation-show-status.md:29-32`
- `qualitymd evaluation show-status <empty-run> --json`

## Remediation options

- Emit 'missing-root-analysis' for the zero-analysis case so the implementation matches the spec's single documented gap kind.
- Keep 'missing-analysis' for the zero case but amend the spec to document both kinds and when each applies.
- Emit both gaps for the zero case (superset) so callers branching on either kind succeed.

## Recommended option

Option 1: align the implementation to the spec's 'missing-root-analysis' so a single documented gap kind covers the empty-run case; it is the smallest change that removes the divergence. Option 2 is acceptable if the team prefers to distinguish zero-records from records-without-root.

## Done criterion

An empty run's show-status gap kind matches the spec exactly (either the CLI emits 'missing-root-analysis' or the spec documents 'missing-analysis' for the zero case). Re-evaluate in a new numbered run.
