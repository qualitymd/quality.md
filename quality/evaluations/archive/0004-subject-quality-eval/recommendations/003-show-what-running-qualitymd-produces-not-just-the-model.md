---
schemaVersion: 1
title: Show what running qualitymd produces, not just the model
gap: The README shows a realistic QUALITY.md excerpt but describes the payoff only in prose; it never shows what running qualitymd against that file produces (no lint output, init receipt, or report excerpt), so a reader cannot see the show-don't-tell payoff.
evidenceLocators:
  - README.md:61-118
  - README.md:120-121
assessmentRecords:
  - quality/evaluations/0004-subject-quality-eval/assessments/011-readme-the-readme-shows-the-format-and-its-payoff-by-example.json
remediationOptions:
  - 'Add a short rendered output block right after the example: a lint result plus an abbreviated report (root rating, one target, one finding) produced from the shown file.'
  - Link to a committed sample run under quality/evaluations/ and inline a trimmed report excerpt.
  - Add only a one-line sample lint/init result as a minimal payoff signal.
recommendedOption: 'Option 1: a trimmed rendered report excerpt plus a one-line lint result immediately after the example is the most direct proof of payoff and stays next to the example the reader just saw.'
doneCriterion: 'The requirement ''the README shows the format and its payoff by example'' reaches at least target: the README displays representative qualitymd output produced from the shown example. Re-evaluate in a new numbered run.'
---

# Show what running qualitymd produces, not just the model

## Gap

The README shows a realistic QUALITY.md excerpt but describes the payoff only in prose; it never shows what running qualitymd against that file produces (no lint output, init receipt, or report excerpt), so a reader cannot see the show-don't-tell payoff.

## Evidence locators

- `README.md:61-118`
- `README.md:120-121`

## Remediation options

- Add a short rendered output block right after the example: a lint result plus an abbreviated report (root rating, one target, one finding) produced from the shown file.
- Link to a committed sample run under quality/evaluations/ and inline a trimmed report excerpt.
- Add only a one-line sample lint/init result as a minimal payoff signal.

## Recommended option

Option 1: a trimmed rendered report excerpt plus a one-line lint result immediately after the example is the most direct proof of payoff and stays next to the example the reader just saw.

## Done criterion

The requirement 'the README shows the format and its payoff by example' reaches at least target: the README displays representative qualitymd output produced from the shown example. Re-evaluate in a new numbered run.
