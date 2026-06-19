---
schemaVersion: 1
title: Make show-status --json emit repo-relative paths
gap: 'qualitymd evaluation show-status <run> --json returns a machine-varying absolute path in the ''path'' field and inside nextActions[].command, even when given a repo-relative run path. This violates payload determinism: the same input and file state produce checkout-location-dependent output that agents and CI cannot diff or cache.'
evidenceLocators:
  - internal/cli/evaluation.go:129-149
  - docs/guides/cli-design.md:224-232
assessmentRecords:
  - quality/evaluations/0004-subject-quality-eval/assessments/014-cli-the-cli-follows-its-functional-specifications.json
remediationOptions:
  - Normalize the run path to a repository-relative path before emitting it in the show-status payload and in derived nextActions commands, matching create-run's relative-path behavior.
  - Echo back exactly the path form the caller supplied (relative in, relative out) without absolutizing.
  - Document the absolute path as intentional and exclude it from determinism guarantees (not recommended; weakens the automation contract).
recommendedOption: 'Option 1: normalize to a repo-relative path so show-status matches create-run and the determinism principle holds across sibling commands.'
doneCriterion: show-status --json emits a repo-relative path for both 'path' and nextActions[].command given a relative run argument, and the requirement 'the CLI follows its functional specifications' reaches at least target. Re-evaluate in a new numbered run.
---

# Make show-status --json emit repo-relative paths

## Gap

qualitymd evaluation show-status <run> --json returns a machine-varying absolute path in the 'path' field and inside nextActions[].command, even when given a repo-relative run path. This violates payload determinism: the same input and file state produce checkout-location-dependent output that agents and CI cannot diff or cache.

## Evidence locators

- `internal/cli/evaluation.go:129-149`
- `docs/guides/cli-design.md:224-232`

## Remediation options

- Normalize the run path to a repository-relative path before emitting it in the show-status payload and in derived nextActions commands, matching create-run's relative-path behavior.
- Echo back exactly the path form the caller supplied (relative in, relative out) without absolutizing.
- Document the absolute path as intentional and exclude it from determinism guarantees (not recommended; weakens the automation contract).

## Recommended option

Option 1: normalize to a repo-relative path so show-status matches create-run and the determinism principle holds across sibling commands.

## Done criterion

show-status --json emits a repo-relative path for both 'path' and nextActions[].command given a relative run argument, and the requirement 'the CLI follows its functional specifications' reaches at least target. Re-evaluate in a new numbered run.
