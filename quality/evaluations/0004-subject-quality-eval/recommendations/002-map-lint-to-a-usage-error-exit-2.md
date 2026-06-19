---
schemaVersion: 1
title: Map 'lint -' to a usage error (exit 2)
gap: qualitymd lint - returns a plain error and exits 70 (internal error), while the sibling qualitymd status - returns a usage error and exits 2. A bad stdin argument is a usage error per the exit-code table, and the divergence makes the CLI inconsistent across sibling commands.
evidenceLocators:
  - internal/cli/lint.go:31-33
  - internal/cli/status.go:25-27
  - docs/guides/cli-design.md:189-199
assessmentRecords:
  - quality/evaluations/0004-subject-quality-eval/assessments/014-cli-the-cli-follows-its-functional-specifications.json
remediationOptions:
  - Wrap the 'lint does not read from stdin yet' error in usageError so it maps to exit 2, matching status.
  - Implement stdin support for lint via '-' so the sentinel becomes valid input.
  - Leave as-is (not recommended; perpetuates the cross-command inconsistency).
recommendedOption: 'Option 1: return a usage error for now (cheap, removes the inconsistency); implement ''-'' support later as a separate enhancement.'
doneCriterion: qualitymd lint - exits 2 with a usage diagnostic, matching status; a test asserts the exit code. Re-evaluate in a new numbered run.
---

# Map 'lint -' to a usage error (exit 2)

## Gap

qualitymd lint - returns a plain error and exits 70 (internal error), while the sibling qualitymd status - returns a usage error and exits 2. A bad stdin argument is a usage error per the exit-code table, and the divergence makes the CLI inconsistent across sibling commands.

## Evidence locators

- `internal/cli/lint.go:31-33`
- `internal/cli/status.go:25-27`
- `docs/guides/cli-design.md:189-199`

## Remediation options

- Wrap the 'lint does not read from stdin yet' error in usageError so it maps to exit 2, matching status.
- Implement stdin support for lint via '-' so the sentinel becomes valid input.
- Leave as-is (not recommended; perpetuates the cross-command inconsistency).

## Recommended option

Option 1: return a usage error for now (cheap, removes the inconsistency); implement '-' support later as a separate enhancement.

## Done criterion

qualitymd lint - exits 2 with a usage diagnostic, matching status; a test asserts the exit code. Re-evaluate in a new numbered run.
