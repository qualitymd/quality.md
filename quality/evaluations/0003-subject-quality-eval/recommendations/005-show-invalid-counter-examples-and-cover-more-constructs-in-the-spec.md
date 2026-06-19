---
schemaVersion: 1
title: Show invalid counter-examples and cover more constructs in the spec
gap: The requirement asks for both valid and invalid examples, but SPECIFICATION.md enumerates invalid cases only in prose (no worked invalid YAML), and several constructs (direct-target requirement, sub-factors, multi-target trees, source/glob, ratings overrides) have no worked example at all.
evidenceLocators:
  - SPECIFICATION.md:249-250
  - SPECIFICATION.md:275-277
  - SPECIFICATION.md:476-498
assessmentRecords:
  - quality/evaluations/0003-subject-quality-eval/assessments/007-format-spec-the-format-s-constructs-are-shown-with-valid-and-invalid-examples.json
remediationOptions:
  - Add a small set of valid/invalid example pairs for the high-traffic constructs (assessment cardinality, factors connection, direct-target requirement), each a short YAML snippet with a one-line why-invalid.
  - Add invalid counter-examples only for the constructs whose rules are most error-prone, deferring full construct coverage.
  - Add a single combined example file exercising sub-factors, multi-target trees, and source/glob to broaden valid coverage.
recommendedOption: 'Option 1: paired valid/invalid snippets for the most error-prone constructs directly satisfy the requirement and most help implementers; Option 3 can follow to broaden valid coverage.'
doneCriterion: 'The requirement ''the format''s constructs are shown with valid and invalid examples'' reaches at least target: worked invalid counter-examples accompany the principal constructs. Re-evaluate in a new numbered run.'
---

# Show invalid counter-examples and cover more constructs in the spec

## Gap

The requirement asks for both valid and invalid examples, but SPECIFICATION.md enumerates invalid cases only in prose (no worked invalid YAML), and several constructs (direct-target requirement, sub-factors, multi-target trees, source/glob, ratings overrides) have no worked example at all.

## Evidence locators

- `SPECIFICATION.md:249-250`
- `SPECIFICATION.md:275-277`
- `SPECIFICATION.md:476-498`

## Remediation options

- Add a small set of valid/invalid example pairs for the high-traffic constructs (assessment cardinality, factors connection, direct-target requirement), each a short YAML snippet with a one-line why-invalid.
- Add invalid counter-examples only for the constructs whose rules are most error-prone, deferring full construct coverage.
- Add a single combined example file exercising sub-factors, multi-target trees, and source/glob to broaden valid coverage.

## Recommended option

Option 1: paired valid/invalid snippets for the most error-prone constructs directly satisfy the requirement and most help implementers; Option 3 can follow to broaden valid coverage.

## Done criterion

The requirement 'the format's constructs are shown with valid and invalid examples' reaches at least target: worked invalid counter-examples accompany the principal constructs. Re-evaluate in a new numbered run.
