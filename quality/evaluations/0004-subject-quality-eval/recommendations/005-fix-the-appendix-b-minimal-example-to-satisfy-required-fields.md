---
schemaVersion: 1
title: Fix the Appendix B minimal example to satisfy required fields
gap: SPECIFICATION.md Appendix B's minimal example omits the required 'title' field on both rating levels and on the 'reliability' factor, contradicting the schema rules the example is meant to illustrate. The requirement 'the format specification is internally consistent' explicitly tests that every example agrees with the rule it illustrates.
evidenceLocators:
  - SPECIFICATION.md:481-493
  - SPECIFICATION.md:165-170
  - SPECIFICATION.md:220
assessmentRecords:
  - quality/evaluations/0004-subject-quality-eval/assessments/005-format-spec-the-format-specification-is-internally-consistent.json
remediationOptions:
  - Add the required 'title' fields to the two rating levels and the factor in the Appendix B example.
  - Relax the schema so 'title' is recommended rather than required (not recommended; broadens scope and weakens display guarantees).
  - Replace Appendix B with the already-conforming Appendix A scale plus a complete factor block.
recommendedOption: 'Option 1: add the missing title fields; it is the smallest change and keeps the minimal example genuinely conforming.'
doneCriterion: The Appendix B example includes all required title fields and would pass lint; the requirement 'the format specification is internally consistent' reaches at least target. Re-evaluate in a new numbered run.
---

# Fix the Appendix B minimal example to satisfy required fields

## Gap

SPECIFICATION.md Appendix B's minimal example omits the required 'title' field on both rating levels and on the 'reliability' factor, contradicting the schema rules the example is meant to illustrate. The requirement 'the format specification is internally consistent' explicitly tests that every example agrees with the rule it illustrates.

## Evidence locators

- `SPECIFICATION.md:481-493`
- `SPECIFICATION.md:165-170`
- `SPECIFICATION.md:220`

## Remediation options

- Add the required 'title' fields to the two rating levels and the factor in the Appendix B example.
- Relax the schema so 'title' is recommended rather than required (not recommended; broadens scope and weakens display guarantees).
- Replace Appendix B with the already-conforming Appendix A scale plus a complete factor block.

## Recommended option

Option 1: add the missing title fields; it is the smallest change and keeps the minimal example genuinely conforming.

## Done criterion

The Appendix B example includes all required title fields and would pass lint; the requirement 'the format specification is internally consistent' reaches at least target. Re-evaluate in a new numbered run.
