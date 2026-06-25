---
schemaVersion: 1
title: Refresh pinned installer examples
gap: The install guide's agent/CI examples still pin v0.5.1 even though the active CLI and skill pair are 0.11.0, so automated users can be steered to stale tooling.
evidenceLocators:
  - install.md:83
  - install.md:87
  - install.md:99
  - assessments/023-docs-docs-reflect-the-current-command-install-release-and-workflow-surfaces.json
  - assessments/026-distribution-supported-install-and-update-paths-are-documented-tested-and-package-the-expected-binary.json
assessmentResultRecords:
  - assessments/023-docs-docs-reflect-the-current-command-install-release-and-workflow-surfaces.json
  - assessments/026-distribution-supported-install-and-update-paths-are-documented-tested-and-package-the-expected-binary.json
remediationOptions:
  - Replace the literal stale pins with the current release tag during each release prep
  - Use a documented placeholder such as vX.Y.Z and explain that agents/CI should substitute the intended release
recommendedOption: Use a documented placeholder such as vX.Y.Z and explain that agents/CI should substitute the intended release
doneCriterion: The install guide no longer contains stale v0.5.1 pins, pinned examples cannot age silently, and the docs and distribution assessments reach target on re-evaluation.
---

# Refresh pinned installer examples

## Gap

The install guide's agent/CI examples still pin v0.5.1 even though the active CLI and skill pair are 0.11.0, so automated users can be steered to stale tooling.

## Evidence locators

- `install.md:83`
- `install.md:87`
- `install.md:99`
- `assessments/023-docs-docs-reflect-the-current-command-install-release-and-workflow-surfaces.json`
- `assessments/026-distribution-supported-install-and-update-paths-are-documented-tested-and-package-the-expected-binary.json`

## Remediation options

- Replace the literal stale pins with the current release tag during each release prep
- Use a documented placeholder such as vX.Y.Z and explain that agents/CI should substitute the intended release

## Recommended option

Use a documented placeholder such as vX.Y.Z and explain that agents/CI should substitute the intended release

## Done criterion

The install guide no longer contains stale v0.5.1 pins, pinned examples cannot age silently, and the docs and distribution assessments reach target on re-evaluation.
