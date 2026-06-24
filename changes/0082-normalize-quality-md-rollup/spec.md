---
type: Functional Specification
title: Normalize QUALITY.md self-check roll-up - functional spec
description: Requirements for treating the `quality-md` self-check area as an ordinary in-scope area for evaluation and roll-up while preserving quality-log behavior for model changes.
tags: [skill, authoring, evaluation, roll-up, quality-md]
timestamp: 2026-06-24T00:00:00Z
---

# Normalize QUALITY.md self-check roll-up - functional spec

Companion to the
[Normalize QUALITY.md self-check roll-up](../0082-normalize-quality-md-rollup.md)
change case. This spec states the delta for the `/quality` skill guidance and
its durable spec mirror. It governs the bundled skill
([`skills/quality/`](../../../skills/quality/)) and its functional-spec mirror
([`specs/skills/quality-skill/`](../../../specs/skills/quality-skill/)), and
defers the QUALITY.md format itself to
[`SPECIFICATION.md`](../../../SPECIFICATION.md).

The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

The current authoring guide makes the QUALITY.md self-check a recurring
use-context constituent, but then says it stays on a separate learn-loop axis and
out of the entity roll-up. That exception is not encoded in the format, CLI,
record schema, or report renderer; it exists only as evaluator guidance. The
result is avoidable ambiguity: a full evaluation can mean every modeled area, or
every modeled area except `quality-md`.

The simpler rule is that a modeled area behaves like an area. A `QUALITY.md`
artifact can be vague, stale, ungrounded, or unassessable; those are real quality
findings for a project whose primary experience depends on agent-readable quality
expectations. Model-change follow-up still needs special logging, but that is a
mutation-history concern, not a reason to exclude the area from evaluation.

## Scope

Covered: removing the special roll-up exclusion for the QUALITY.md self-check;
teaching that `quality-md` is evaluated, analyzed, reported, and rolled up like
any other in-scope area; preserving setup's `quality-md` area pattern; and
preserving the existing quality-log behavior for confirmed model changes.

Deferred / non-goals: no QUALITY.md schema or format change, no `SPECIFICATION.md`
change, no CLI/report schema change, no Go implementation change, no new
roll-up-exclusion flag, and no migration command. Re-evaluating this repository's
own current model after implementation is a follow-up.

## Requirements

### The self-check is an ordinary area

The authoring guide and its spec mirror **MUST** teach the QUALITY.md self-check
as a recurring use-context constituent that is modeled by default when germane,
using the normal area pattern: the `quality-md` key, a title of the form
`<Root Title> QUALITY.md`, an explicit path-based `source` such as
`./QUALITY.md`, factors that describe the model artifact's qualities, and a
requirement that assesses the model against the active authoring guide.

The guidance **MUST NOT** say that the QUALITY.md self-check is kept out of the
root area's roll-up, excluded from aggregate rating, reported only on a separate
axis, or handled by different evaluation semantics.

> Rationale: the area already has ordinary model structure and the CLI already
> treats it generically. The exception creates process ambiguity without a
> mechanical contract to enforce it. - 0082

### Full evaluation covers all in-scope areas uniformly

The evaluation workflow guidance and its spec mirror **MUST** state or preserve
the rule that an unnarrowed `/quality evaluate` covers every in-scope modeled area
with assessable requirements, including `quality-md` when present. Missing
assessment or analysis coverage for `quality-md` in a full run **MUST** be treated
as the same kind of incomplete evaluation coverage as missing coverage for any
other modeled area.

The guidance **MUST NOT** make `quality-md` opt-in, out-of-band, or excluded from
full evaluation by default.

> Rationale: "full evaluation" should mean the resolved model scope. Excluding a
> named area forces every evaluator to rediscover a convention that the model and
> records cannot express. - 0082

### Roll-up semantics are uniform

The authoring and evaluation guidance **MUST** teach that `quality-md` contributes
to factor, local, and aggregate ratings according to the same inferred,
importance-weighted roll-up judgment as other child areas when it is in scope.
Reports **MAY** label the area clearly as model quality, but **MUST NOT** remove
it from the aggregate result solely because its source is `QUALITY.md`.

> Rationale: weak model quality can make future evaluation and improvement work
> worse. If the area is modeled, its rating should be consequential like the
> rating of any other project artifact. - 0082

### Model-change logging remains distinct

Recommendation follow-up guidance **MUST** continue to require a quality-log
entry when a confirmed apply changes the QUALITY.md model in a meaningful way.
The guidance **SHOULD** frame this as mutation-history behavior: the quality log
records why the model changed; it does not create separate evaluation semantics
for the `quality-md` area.

Evaluated-source fixes outside `QUALITY.md` **MUST NOT** write the quality log
unless they also change the model.

### Learn-loop wording no longer implies exclusion

Live guidance **MUST NOT** use "learn loop" or similar wording to imply that
QUALITY.md model quality is excluded from ordinary evaluation or roll-up. If
"learn loop" terminology remains, it **MUST** refer only to how model changes are
recognized, justified, and logged after evaluation reveals that the model itself
should change.

## Acceptance criteria

- Searching live guidance and spec mirrors for `out of roll-up`,
  `kept out of roll-up`, `out of the entity's roll-up`, `learn-loop axis`, and
  `never averaged` finds no normative instruction that excludes `quality-md` from
  ordinary roll-up.
- `skills/quality/guides/authoring.md` says the QUALITY.md self-check follows
  ordinary area evaluation and roll-up semantics.
- `specs/skills/quality-skill/guides/authoring-md.md` mirrors the same rule.
- Evaluation guidance either remains clearly generic or explicitly says a full
  evaluation covers all in-scope modeled areas, including `quality-md`.
- Setup guidance still includes the `quality-md` area pattern with
  `source: ./QUALITY.md` and an authoring-guide assessment.
- Recommendation follow-up still requires a quality-log entry for meaningful
  confirmed model changes and does not require one for ordinary evaluated-source
  fixes.
- This repository's `QUALITY.md` no longer describes the `quality-md` area as an
  out-of-band learn-loop axis.
- Verification includes `qualitymd lint QUALITY.md --json`,
  `qualitymd status QUALITY.md --json`, and the repository's markdown formatting
  check.

## Durable spec changes

### To add

None

### To modify

- `specs/skills/quality-skill/guides/authoring-md.md` - remove the self-check
  roll-up exclusion and require ordinary area semantics for `quality-md` (per
  [The self-check is an ordinary area](#the-self-check-is-an-ordinary-area),
  [Roll-up semantics are uniform](#roll-up-semantics-are-uniform), and
  [Learn-loop wording no longer implies exclusion](#learn-loop-wording-no-longer-implies-exclusion)).
- `specs/skills/quality-skill/evaluation.md` - clarify that full evaluation
  covers all in-scope modeled areas uniformly, including `quality-md` when present
  (per [Full evaluation covers all in-scope areas uniformly](#full-evaluation-covers-all-in-scope-areas-uniformly)).

### To rename

None

### To delete

None
