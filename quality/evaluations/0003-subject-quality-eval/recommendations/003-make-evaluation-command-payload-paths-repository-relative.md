---
schemaVersion: 1
title: Make evaluation command payload paths repository-relative
gap: '''qualitymd evaluation show-status --json'' emits an absolute filesystem path in the ''path'' field while ''create-run --json'' emits a repository-relative path. The absolute path is machine-varying, breaking determinism (specs/cli.md:99-100) and cross-command consistency, and contradicting the repository-relative requirement in specs/cli/evaluation-create-run.md:33. add-record and set-planned-coverage receipts share the behavior.'
evidenceLocators:
  - specs/cli/evaluation-create-run.md:33
  - specs/cli.md:99-100
  - qualitymd evaluation show-status <run> --json
assessmentRecords:
  - quality/evaluations/0003-subject-quality-eval/assessments/014-cli-the-cli-follows-its-functional-specifications.json
remediationOptions:
  - Normalize the 'path' field (and the human 'Wrote ...' footer's machine-relevant references) to a repository-relative path in show-status, add-record, set-planned-coverage, and build-report receipts.
  - Keep the absolute path only on human stderr footers but make every --json payload path repository-relative.
  - Document absolute paths as intended and relax the create-run relative-path MUST (not recommended; loses cross-machine diffability).
recommendedOption: 'Option 2: emit repository-relative paths in all --json payloads (matching create-run) while leaving human stderr free to show a fuller path; this restores determinism and cross-command consistency without losing human ergonomics.'
doneCriterion: show-status --json (and the other evaluation receipts) emit a repository-relative 'path' identical in form to create-run, with no machine-varying value in any --json payload. Re-evaluate in a new numbered run.
---

# Make evaluation command payload paths repository-relative

## Gap

'qualitymd evaluation show-status --json' emits an absolute filesystem path in the 'path' field while 'create-run --json' emits a repository-relative path. The absolute path is machine-varying, breaking determinism (specs/cli.md:99-100) and cross-command consistency, and contradicting the repository-relative requirement in specs/cli/evaluation-create-run.md:33. add-record and set-planned-coverage receipts share the behavior.

## Evidence locators

- `specs/cli/evaluation-create-run.md:33`
- `specs/cli.md:99-100`
- `qualitymd evaluation show-status <run> --json`

## Remediation options

- Normalize the 'path' field (and the human 'Wrote ...' footer's machine-relevant references) to a repository-relative path in show-status, add-record, set-planned-coverage, and build-report receipts.
- Keep the absolute path only on human stderr footers but make every --json payload path repository-relative.
- Document absolute paths as intended and relax the create-run relative-path MUST (not recommended; loses cross-machine diffability).

## Recommended option

Option 2: emit repository-relative paths in all --json payloads (matching create-run) while leaving human stderr free to show a fuller path; this restores determinism and cross-command consistency without losing human ergonomics.

## Done criterion

show-status --json (and the other evaluation receipts) emit a repository-relative 'path' identical in form to create-run, with no machine-varying value in any --json payload. Re-evaluate in a new numbered run.
