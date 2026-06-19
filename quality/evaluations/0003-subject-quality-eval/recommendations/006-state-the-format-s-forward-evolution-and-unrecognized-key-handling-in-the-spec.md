---
schemaVersion: 1
title: State the format's forward-evolution and unrecognized-key handling in the spec
gap: SPECIFICATION.md specifies the minimal core and structural extension but delegates the forward-evolution/versioning policy to an external docs/reference/versioning.md, and an unrecognized frontmatter key that is neither defined nor a sanctioned extension has no explicit reader rule (it leans toward making the file non-conforming).
evidenceLocators:
  - SPECIFICATION.md:3
  - SPECIFICATION.md:16-19
  - SPECIFICATION.md:113-117
  - SPECIFICATION.md:431-443
assessmentRecords:
  - quality/evaluations/0003-subject-quality-eval/assessments/008-format-spec-the-format-specifies-its-core-and-how-it-extends-and-evolves.json
remediationOptions:
  - Add a short 'Versioning and evolution' section to SPECIFICATION.md stating a forward-compatibility rule for unrecognized frontmatter keys (ignore-and-preserve vs reject) and a version-bump/compatibility policy, even if details stay in the external doc.
  - State only the unrecognized-frontmatter-key rule now and keep deferring the version policy to the external doc, with an explicit in-spec pointer summarizing its intent.
  - Record the missing evolution policy as an explicit Known gap in the spec while pre-1.0, naming the intended direction.
recommendedOption: 'Option 1: a brief in-spec section closes the ''how it evolves'' half of the requirement and removes the unrecognized-key ambiguity; Option 3 is an honest interim if the team wants to stay minimal pre-1.0.'
doneCriterion: 'The requirement ''the format specifies its core and how it extends and evolves'' reaches at least target: the spec states how the format versions forward and how readers treat unrecognized frontmatter content. Re-evaluate in a new numbered run.'
---

# State the format's forward-evolution and unrecognized-key handling in the spec

## Gap

SPECIFICATION.md specifies the minimal core and structural extension but delegates the forward-evolution/versioning policy to an external docs/reference/versioning.md, and an unrecognized frontmatter key that is neither defined nor a sanctioned extension has no explicit reader rule (it leans toward making the file non-conforming).

## Evidence locators

- `SPECIFICATION.md:3`
- `SPECIFICATION.md:16-19`
- `SPECIFICATION.md:113-117`
- `SPECIFICATION.md:431-443`

## Remediation options

- Add a short 'Versioning and evolution' section to SPECIFICATION.md stating a forward-compatibility rule for unrecognized frontmatter keys (ignore-and-preserve vs reject) and a version-bump/compatibility policy, even if details stay in the external doc.
- State only the unrecognized-frontmatter-key rule now and keep deferring the version policy to the external doc, with an explicit in-spec pointer summarizing its intent.
- Record the missing evolution policy as an explicit Known gap in the spec while pre-1.0, naming the intended direction.

## Recommended option

Option 1: a brief in-spec section closes the 'how it evolves' half of the requirement and removes the unrecognized-key ambiguity; Option 3 is an honest interim if the team wants to stay minimal pre-1.0.

## Done criterion

The requirement 'the format specifies its core and how it extends and evolves' reaches at least target: the spec states how the format versions forward and how readers treat unrecognized frontmatter content. Re-evaluate in a new numbered run.
