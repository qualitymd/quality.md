---
schemaVersion: 1
title: Add worked invalid counter-examples to the specification
gap: SPECIFICATION.md shows valid worked examples and describes invalidity in prose, but provides no worked invalid counter-example blocks. The requirement 'the format's constructs are shown with valid and invalid examples' asks for both.
evidenceLocators:
  - SPECIFICATION.md:446-498
  - SPECIFICATION.md:250-278
assessmentRecords:
  - quality/evaluations/0004-subject-quality-eval/assessments/007-format-spec-the-format-s-constructs-are-shown-with-valid-and-invalid-examples.json
remediationOptions:
  - Add a small set of invalid counter-example blocks next to the key constructs (e.g. a requirement with a list-valued assessment, a direct-target requirement missing factors, a rating scale with one level), each paired with why it is non-conforming.
  - Add a single combined 'invalid examples' appendix covering the most common mistakes.
  - Reference an external conformance test suite of invalid files instead of inlining counter-examples.
recommendedOption: 'Option 1: inline a few targeted invalid counter-examples beside the constructs they illustrate, so readers see valid and invalid side by side at the point of definition.'
doneCriterion: The spec shows worked invalid counter-examples for its key constructs; the requirement 'the format's constructs are shown with valid and invalid examples' reaches at least target. Re-evaluate in a new numbered run.
---

# Add worked invalid counter-examples to the specification

## Gap

SPECIFICATION.md shows valid worked examples and describes invalidity in prose, but provides no worked invalid counter-example blocks. The requirement 'the format's constructs are shown with valid and invalid examples' asks for both.

## Evidence locators

- `SPECIFICATION.md:446-498`
- `SPECIFICATION.md:250-278`

## Remediation options

- Add a small set of invalid counter-example blocks next to the key constructs (e.g. a requirement with a list-valued assessment, a direct-target requirement missing factors, a rating scale with one level), each paired with why it is non-conforming.
- Add a single combined 'invalid examples' appendix covering the most common mistakes.
- Reference an external conformance test suite of invalid files instead of inlining counter-examples.

## Recommended option

Option 1: inline a few targeted invalid counter-examples beside the constructs they illustrate, so readers see valid and invalid side by side at the point of definition.

## Done criterion

The spec shows worked invalid counter-examples for its key constructs; the requirement 'the format's constructs are shown with valid and invalid examples' reaches at least target. Re-evaluate in a new numbered run.
