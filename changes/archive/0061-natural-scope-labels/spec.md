---
type: Functional Specification
title: Natural scope labels
description: Make natural Area and Factor labels the primary scoped-evaluation input for the /quality skill.
tags: [skill, evaluation, ux, references]
timestamp: 2026-06-23T00:00:00Z
---

# Natural scope labels

This Change Case spec defines the delta for `/quality` skill guidance and
user-facing docs: scoped evaluation should be documented around natural Area and
Factor labels, while the skill continues to resolve those labels to stable model
IDs before any evaluation artifacts are written.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background / Motivation

The current documentation foregrounds canonical references like
`factor:root::security`. Those references are the right durable and exact form,
but they make the normal skill interaction feel like a tool API instead of an
agent workflow. A user should be able to say the Area or Factor they recognize
from the model, and the skill should perform the grounded resolution work.

This preserves the stable identifier model while making the default interaction
more natural. Ambiguity is still handled explicitly: the skill asks for the
missing judgment instead of guessing.

## Scope

Covered:

- Natural one-label scoped evaluation input.
- Natural two-label scoped evaluation input.
- Area-first disambiguation for repeated Factor labels.
- Continued support for qualified model references.
- Clear separation between user-facing labels and durable record identifiers.
- README, durable skill spec, and runtime skill guidance updates.

Deferred / non-goals:

- No change to the QUALITY.md format specification or model-reference grammar.
- No new durable machine-readable fields.
- No requirement for fuzzy, partial, or semantic label matching.
- No requirement that the CLI parse natural labels directly.

## Terminology

A **natural label** is user-provided text meant to identify an Area or Factor by
the human vocabulary visible in the model. The skill may match it against Area
titles, Area names, Factor titles, and Factor names.

A **qualified model reference** is the canonical exact form such as
`area:payments` or `factor:payments::maintainability`.

## Requirements

The `/quality` skill **SHOULD** document natural labels as the primary
user-facing scoped-evaluation input.

> Rationale: the skill owns human-edge judgment. Users should not have to learn
> canonical model-reference syntax before asking for a familiar quality concern.
> — 0061

The `/quality` skill **MUST** continue to accept qualified model references for
exact addressing, disambiguation, and advanced workflows.

The skill **MUST** resolve every scoped evaluation to stable model IDs before
creating evaluation artifacts or writing record payloads.

Evaluation records and `report.json` **MUST NOT** persist natural labels in
place of structured `areaPath`, `factorPath`, or rating `level` identifiers.

For `/quality evaluate <label>`, the skill **SHOULD** resolve the label against
the current `QUALITY.md` model as follows:

- If it uniquely identifies one Area, evaluate that Area.
- If it uniquely identifies one Factor, evaluate that Factor in its declaring
  Area.
- If it identifies a Factor label present in multiple Areas, ask which Area the
  user wants to evaluate that Factor for.
- If it matches both Area and Factor candidates, ask a targeted clarification
  question before evaluating.
- If it does not resolve, report that the label is not in the model and offer
  the nearest runnable scoped-evaluation options visible from the model.

For `/quality evaluate <area-label> <factor-label>`, the skill **SHOULD**
resolve the Area label first, then resolve the Factor label within that Area.

When a repeated Factor label needs an Area, the clarification prompt **SHOULD**
use this wording:

```text
What area do you want to evaluate <Factor> for?
```

The clarification options **SHOULD** lead with human-readable Area titles or
names. Qualified references may appear as secondary context where traceability
matters.

User-facing examples **SHOULD** lead with natural-label forms, such as:

```text
/quality evaluate Security
/quality evaluate Payments Maintainability
/quality evaluate API Reliability
```

Qualified-reference examples **SHOULD** be presented as exact-addressing or
advanced syntax, not as the primary usage pattern.

## Acceptance Criteria

- `README.md` scoped-evaluation examples foreground natural labels.
- `README.md` still documents qualified references as exact or advanced syntax.
- `skills/quality/SKILL.md` lists one-label and two-label scoped evaluation as
  normal invocation variants.
- `skills/quality/SKILL.md` no longer describes bare names as legacy shorthand.
- `skills/quality/modes/evaluate.md` treats one-label and two-label resolution
  as primary scope-resolution behavior.
- `skills/quality/modes/evaluate.md` specifies the repeated-Factor clarification
  prompt: "What area do you want to evaluate <Factor> for?"
- `specs/skills/quality-skill/quality-skill.md` defines natural labels as the
  primary scoped-evaluation input and keeps qualified references as supported
  exact addressing.
- `specs/skills/quality-skill/evaluation.md` aligns evaluation-mode scope
  behavior with natural one-label and two-label inputs.
- Durable machine-readable artifacts continue to use stable identifiers, not
  natural labels.
- No change is made to `SPECIFICATION.md` unless implementation discovers a
  format-level contract gap.

## Durable spec changes

### To add

None.

### To modify

- `specs/skills/quality-skill/quality-skill.md` - make natural labels the
  primary scoped-evaluation input, define one-label and two-label resolution,
  keep qualified references as exact-addressing input, and preserve durable
  stable identifiers according to the requirements above.
- `specs/skills/quality-skill/evaluation.md` - align evaluation-mode scope
  behavior, ambiguity handling, and clarification prompts with the requirements
  above.

### To rename

None.

### To delete

None.
