---
type: Change Case
title: Model constituents by default
description: Flip the skill's constituent-coverage guidance from earn-it/defer-freely to model-by-default with a short "don't model" list, so a first-pass model is as full and sufficient as the evidence supports.
status: Done
tags: [skill, authoring, setup, constituents, coverage]
timestamp: 2026-06-24T00:00:00Z
---

# Model constituents by default

A **Change Case** to reset how the `/quality` skill decides which constituents of
a composite root become their own areas. Today an area is something a constituent
must **earn** — a three-part test (owned, inspectable artifact; divergent factor
family; traced to a Need or Risk) wrapped in "a prompt, not a quota" /
"not a roster" language. Every pressure points toward fewer areas, and
"defer it in Scope" sits as a cost-free peer of "model it." A real setup run on a
multi-service monorepo produced a flat root with five model-wide factors and
deferred every per-constituent area (`server`, `client`, the data pipeline,
tests) to "the next iteration."

This case flips the default: **enumerate the constituents the domain implies and
model each as its own area by default.** A constituent skips its own area only
when one of two simple disqualifiers holds, and a germane concern is *never*
silently omitted — its absence must surface as a ratable element of the model.

Detail lives in:

- [Functional spec](0080-model-constituents-by-default/spec.md) — what the
  guidance must say.
- [Design doc](0080-model-constituents-by-default/design.md) — the model-by-default
  reframe, the disqualifier set, the no-silent-omission rule, and the (a)
  minimal-area vs (b) requirement-elsewhere routing.

## Motivation

A first-pass model should be as full and sufficient as the evidence supports.
Setup runs roughly once per project, the user is most engaged then, and an
under-built model that defers its real constituents teaches the wrong lesson:
that a flat root is the normal starting shape and coverage is something to add
later. It rarely gets added later.

The current guidance is not buggy — it does exactly what it says. The problem is
the encoded philosophy: it guards hard against *over*-modeling (quota, roster,
earn-it) with no counterweight against *under*-modeling, and it makes deferral a
neutral accounting option. So when the agent is uncertain, every explicit
pressure removes areas, and a deferral note silently passes the maturity bar.

The deeper correction is that a germane concern whose artifact is thin or missing
is one of the highest-value early signals a model can carry. "Nothing to evaluate
yet" must therefore never be a reason to omit — the absence has to land *in the
model* as a poor rating (a minimal area with a missing-anchor finding, or a
requirement on an existing area that rates poorly because the artifact is
missing), not as a sentence in Scope.

## Scope

Covered:

- Reframe the constituent-coverage guidance to **model-by-default** in the
  authoring guide, the setup workflow, and the Top 10 checks (and their spec
  mirrors).
- Replace the three-part "earn an area" inclusion test with a short, explicit
  **"don't give a constituent its own area unless…"** list of two disqualifiers:
  *no distinct concerns* (fold into an existing area) and *not germane / outside
  the boundary* (out of scope).
- Add the **no-silent-omission** rule: a germane concern is always surfaced as a
  ratable element — a minimal area with a missing-anchor finding, or a requirement
  on an existing area — never dropped to prose; with the (a)-vs-(b) routing
  criterion.
- Add a **first-pass completeness bar**: a composite model that leaves a
  non-disqualified constituent unmodeled (or merely deferred) is not
  `evaluation-ready`; the Top 10 check flags it and a bare deferral note does not
  satisfy it.
- Demote **deferral** from a routine accounting option to a narrow,
  blocker-recorded exception ("next iteration" / "first model is thin" are not
  reasons).
- Align the use-context constituents (agent harness, QUALITY.md self-check) and the
  getting-started first-model wording to the model-by-default default.

Deferred / non-goals:

- No change to the QUALITY.md format or schema, and **no change to
  `SPECIFICATION.md`** — this is authoring *judgment* (how to build a good model),
  not format *semantics* (what a model is).
- No change to the CLI.
- Re-checking this repo's own dogfooded `QUALITY.md` against the new rule is a
  follow-up, not part of this case.
- The stewardship-concern generator, the audience×purpose axis, the
  three-projections rule, and the motivation/taxonomy vocabulary discipline are
  unchanged; this case touches only the inclusion default and its enforcement.

## Affected artifacts

Derived by sweeping the skill and its spec mirror for the earn-it / quota /
roster / defer language and the constituent-coverage checks. Empty kinds are
deliberate.

### Code

- No impact. This case changes skill guidance and its functional specs only.

### Durable specs

- [x] `specs/skills/quality-skill/guides/authoring-md.md` — restate the
      constituent-coverage contract as model-by-default with the two disqualifiers
      and the no-silent-omission rule; replace the "earn an area" / "roster" MUSTs.
- [x] `specs/skills/quality-skill/workflows/setup.md` — require the
      model-by-default constituent walk and the first-pass completeness bar in the
      maturity classification / completion criteria; demote deferral.
- [x] `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md` — flip the
      check from "earned defaults, do not flag" to flagging a non-disqualified
      constituent left unmodeled or merely deferred.
- [x] `specs/skills/quality-skill/guides/log.md` — record the contract revision.

### Format spec

- [x] `SPECIFICATION.md` — no change; constituent-coverage is authoring judgment,
      not format semantics.

### Durable docs (bundled skill, guides, README)

- [x] `skills/quality/guides/authoring.md` — canonical statement: model-by-default,
      the two disqualifiers, the no-silent-omission rule, the (a)-vs-(b) routing;
      align the "Carry the recurring use-context constituents" subsection.
- [x] `skills/quality/workflows/setup.md` — model-by-default constituent walk in
      Synthesize; first-pass completeness bar in Verify and Close.
- [x] `skills/quality/guides/top-10-quality-md-checks.md` — flip the
      constituent-coverage findings and the closing "earned, not a roster" note.
- [x] `skills/quality/guides/getting-started.md` — align the "small first model"
      wording so it guards assessability, not under-coverage.
- [x] `README.md` — no change.
- [x] `docs/` guides — no change (repo-development guides, not quality authoring).

### Release

- [x] `CHANGELOG.md` — add an Unreleased `/quality Skill` note for the
      model-by-default reframe.

## Children

- [Functional spec](0080-model-constituents-by-default/spec.md) — required guidance
  behavior for this case.
- [Design doc](0080-model-constituents-by-default/design.md) — the reframe, the
  disqualifier set, and the routing.

## Status

`Done`. Implementation complete and the Affected artifacts list reconciled.
Landed: the authoring guide's model-by-default constituent-coverage rewrite (two
disqualifiers, the no-silent-omission rule with (a)/(b) routing, deferral demoted,
the use-context subsection aligned); the setup constituent walk and the
first-pass completeness bar in Verify and Close; the Top 10 constituent-coverage
findings and closing note; the getting-started assessability rewording; the durable
authoring-guide, setup-workflow, and Top-10 specs (with 0080 rationales); the
guide/workflow spec logs; and the CHANGELOG note. `SPECIFICATION.md`, `README.md`,
and `docs/` assessed no-change; no code change. Verified with `mise run check`
(dprint markdown format, bundle link resolution, Go vet/lint/test all pass) and
`qualitymd lint QUALITY.md`. Archived to `archive/` on landing. Follow-up: re-check
this repo's own dogfooded `QUALITY.md` against the new rule.
