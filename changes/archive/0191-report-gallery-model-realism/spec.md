---
type: Functional Specification
title: Report gallery model realism functional spec
description: Delta contract for deepening the LedgerLite gallery model to realistic scale, guide conformance, sensor-grounded assessments, and workflow-faithful report prose.
tags: [examples, report-gallery, authoring]
timestamp: 2026-07-02T00:00:00Z
---

# Report gallery model realism functional spec

This spec governs the content contract of the
[report gallery](../../../examples/report-gallery/software-service/README.md) as
deepened by this case, building on the exemplar contract of
[0190](../0190-report-gallery-exemplar/spec.md), whose requirements remain in
force except where a requirement here explicitly extends one. Normative
references: the authoring guide family —
[`skills/quality/guides/authoring.md`](../../../skills/quality/guides/authoring.md)
and its routed sub-guides ([`body.md`](../../../skills/quality/guides/authoring/body.md),
[`model-structure.md`](../../../skills/quality/guides/authoring/model-structure.md),
[`factors.md`](../../../skills/quality/guides/authoring/factors.md),
[`requirements.md`](../../../skills/quality/guides/authoring/requirements.md),
[`rating-scale.md`](../../../skills/quality/guides/authoring/rating-scale.md),
[`agent-harnessability.md`](../../../skills/quality/guides/authoring/agent-harnessability.md),
[`agent-harness.md`](../../../skills/quality/guides/authoring/agent-harness.md),
[`quality-changelog.md`](../../../skills/quality/guides/authoring/quality-changelog.md))
— the setup workflow
([`skills/quality/workflows/setup.md`](../../../skills/quality/workflows/setup.md)),
the evaluate workflow
([`skills/quality/workflows/evaluate.md`](../../../skills/quality/workflows/evaluate.md)),
the quality changelog format contract in
[`skills/quality/SKILL.md`](../../../skills/quality/SKILL.md), and the evaluation
data contract enforced by `internal/evaluation`. Informational: the
harness-engineering framing of feedforward guides and feedback sensors
(computational vs. inferential controls) that the assessment-grounding
requirements demonstrate.

The key words "MUST", "SHOULD", and "MAY" are to be interpreted as described in
BCP 14 when, and only when, they appear in all capitals.

## Background / motivation

The 0190 exemplar demonstrates model structure but stops short of the
authoring guides' own coverage bar: most factors carry one requirement, the
payments-domain lenses a knowledgeable reader expects first (security,
auditability) are absent, no constituent covers the implementation itself, two
`quality-md` factor names are practice names the factors guide explicitly
warns against, and nearly every assessment reads as a bespoke inferential
investigation. The deepened gallery must read as a model that earned its
factors from its body, meets the stable-stakes coverage bar per constituent,
and makes the economics of assessment visible: a small named sensor catalog
that many requirements reuse, a deliberate spread of assessment kinds, and a
visible maturation path from judgment to computation.

## Scope

Covered: the LedgerLite `QUALITY.md` frontmatter and body, the generated
`0001-full-eval` run content, the synthetic quality changelog, the gallery
README, and the generator content and payload tables. Deferred: a second
non-software example; a fictional source tree; a multi-run loop demonstration;
a model-wide root `currentness` factor (the composite-root recurring-factor
move). Non-goals: changes to report format, record schema, CLI, or skill
behavior; edits to the authoring guides — where the exemplar and a guide
conflict, the exemplar bends and any guide amendment is a separate scoped
change.

## Assumptions & dependencies

- Change Case 0190 has landed (`Done`, archived): this spec builds on the
  embedded content files, generalized payload tables, and gallery content it
  introduced, and its spec's requirements stay in force except where extended
  here.
- `evaluation.SetData` validates payloads against the evaluation data
  contract; whatever it accepts is contract-conformant.
- The factors guide's coverage aim (roughly ten factors for a primary-subject
  node, applied per constituent at a composite root, with deliberately narrow
  constituents exempt) is stable; requirement 2's per-area families are
  calibrated against it.
- `prettier` formats all checked-in Markdown, including generated gallery
  files.

## Requirements

### Model structure and factors

1. The model MUST add a `codebase` constituent area covering the
   implementation itself, carrying at minimum a `maintainability` factor
   decomposed into sub-factors (such as `analyzability`, `modifiability`,
   `testability`), a `consistency` factor for architecture conformance, and a
   locally refined `security` factor.

   > Rationale: harness engineering's maintainability and architecture-fitness
   > regulation families have no home in the 0190 model because no area covers
   > the implementation; the body's "maintainers and coding agents change it
   > weekly" already earns the constituent. — 0191

   > Durable spec: none.

2. The constituent factor families MUST expand to, at minimum: the API area
   carrying at least seven conventional factors including `security`,
   `reliability`, `compatibility`, and `testability` alongside the existing
   `correctness`, `operability`, and `performance`; the persistence area
   adding `auditability` and `security`; the operations area adding `security`
   and a capacity/performance-family factor; and the service-contract area
   adding `currentness` and `understandability`.

   > Rationale: security and auditability are the omissions a payments-domain
   > reader notices first; contract `currentness` carries the drift risk the
   > body already names; capacity gives the holiday-peak unknown a home. — 0191

   > Durable spec: none.

3. The `quality-md` area's factors MUST use quality names rather than practice
   names: `context-grounding` and `lifecycle-maintenance` are renamed (for
   example to `credibility` and `currentness`), and `evaluability` is
   reconciled with the sibling agent-harness area's `assessability`.

   > Rationale: the factors guide names both offenders — "grounding is a
   > tactic or metaphor," and lifecycle wording names a practice rather than
   > the quality it protects; sibling areas naming the same lens differently
   > reads as accident, not refinement. — 0191

   > Durable spec: none.

4. Every factor in the revised model MUST satisfy the factors guide's
   admission tests (consequential, bounded, operational, traceable, neutral),
   and each area's factor set MUST meet the guide's stable-stakes coverage bar
   or the body's model-shape section MUST justify each omission (at minimum
   `usability` and `portability`) as out of scope, delegated, or unresolved.

   > Durable spec: none.

5. Every factor MUST carry at least one requirement; `correctness`,
   `integrity`, and the API `security` factor MUST each carry two or more; and
   the model total SHOULD land in the 35–45 requirement band.

   > Rationale: one requirement per factor is the tell that made the 0190
   > model read as scaffold; the band keeps the gallery realistic while
   > staying browsable and within generated-artifact budget. — 0191

   > Durable spec: none.

6. The deepening MUST preserve 0190's exemplar moves unchanged in kind: the
   `agent-harnessability` umbrella with its seven sub-factors and
   projection-boundary encoding, the `balance-invariants` veto, the measured
   `ratings` override, the multi-factor referenced assessment, and the
   composite-root shape with its normative service-contract area.

   > Durable spec: none.

### Assessment grounding

7. The model MUST name a small stable sensor catalog (for example: contract
   tests, invariant suite, reconciliation job, latency rollup, lint,
   complexity check, dependency audit, drift detector, structural
   import-boundary tests), with each catalog sensor referenced by name from at
   least two requirement assessments and the naming consistent with the
   agent-harness area's sensor-catalog source.

   > Rationale: mature harnesses have few sensors reused by many checks, not
   > one sensor per check; consistent naming is what lets a reader see the
   > reuse. — 0191

   > Durable spec: none.

8. Requirement assessments MUST span three kinds, each appearing at least
   three times across the model: pure computational (a named sensor's
   deterministic result is the assessment), sensor-plus-guide (a sensor result
   read against the contract, a runbook, or the authoring guides), and
   inferential (judgment over cited materials where no sensor exists).

   > Durable spec: none.

9. At least one assessment MUST be a drift-detection check comparing a
   normative artifact against shipped behavior (contract-vs-handler or
   architecture conformance).

   > Durable spec: none.

10. At least one inferential requirement MUST be paired, in the generated run,
    with a finding and a recommendation whose done criteria establish a
    computational sensor for it.

    > Rationale: the loop converting judgment into computation is the
    > harness-engineering maturation story the gallery exists to show. — 0191

    > Durable spec: none.

11. The body MUST state the assessment posture in one place: assessments
    prefer recorded computational sensors where they exist, fall back to
    inferential review where they do not, and the quality loop's job is to
    shrink the inferential set — worded within the repository's register rules
    (no motivation-layer words modifying taxonomy nouns).

    > Durable spec: none.

### Body revision

12. The Markdown body MUST be re-derived as a whole so every factor family in
    the revised model finds its rationale in the fictional scenario before the
    frontmatter uses it — at minimum: security (unauthorized money movement,
    tenant trust), auditability (the "remains explainable" half of the
    governing sense of good), the codebase area and maintainability family
    (weekly agent-first change), compatibility (the v1 error-envelope
    deprecation question), capacity (holiday-peak load), and contract
    currentness (the drift risk). A reader MUST be able to reconstruct the
    factor tree's major branches from the body alone.

    > Rationale: the authoring order is body-first; a body patched to
    > footnote an already-written frontmatter teaches derivation backwards.
    > — 0191

    > Durable spec: none.

13. The body's unknowns and open questions MUST be reconciled with the revised
    model: concerns the expansion now assesses (such as holiday-peak load
    under a capacity factor) are restated or resolved rather than left
    dangling, new unknowns are added where the scenario implies them, and
    review provenance lines are refreshed to postdate the revision.

    > Durable spec: none.

14. The model-shape section MUST be updated for the revised structure: where
    the codebase and other new areas sit in the worst-of roll-up and required
    margin, and the deliberate-omission justifications of requirement 4.

    > Durable spec: none.

### Guide and workflow fidelity

15. The revised gallery MUST conform to the authoring guide entry point and
    every routed sub-guide named in this spec's normative references, with no
    modeled exceptions: a reviewer walking each sub-guide's do/avoid items
    against the gallery finds no violation, and any exemplar-vs-guide conflict
    is resolved by bending the exemplar (guide amendments are out of scope per
    Scope).

    > Durable spec: none.

16. Names, references, and prose MUST follow repository conventions: the
    strict name grammar for model names, canonical model references where a
    stable text handle is needed, sentence-case headings, and the
    vocabulary-capitalization rules.

    > Durable spec: none.

17. The revised `QUALITY.md` MUST read as a plausible product of the
    `/quality` setup workflow subsequently matured through the quality loop:
    no structure, section, or register the setup workflow would not produce,
    with growth beyond a plausible first draft visible in the quality
    changelog rather than implied by an oversized initial model.

    > Durable spec: none.

18. All generated run prose — report, findings, recommendations, per-area,
    per-factor, and per-requirement pages, and data payloads — MUST be
    consistent with what the evaluate workflow would emit: finding copy
    follows the finding-core field jobs, evidence is named per finding the way
    the workflow names it (sensor invoked, result observed, guide compared
    against), confidence levels are justified in the workflow's terms,
    recommendations carry the workflow's background, expected-value, and
    done-criteria shape, and no finding field describes itself as synthetic
    (extending 0190's rule to all new content).

    > Durable spec: none.

19. Where a requirement's assessment is computational, its finding prose MUST
    read as a sensor result (named sensor, deterministic outcome, remediation
    pointer); where inferential, it MUST read as judgment with cited
    materials — so the assessment-kind spread of requirement 8 is visible in
    the generated report, not only in the model.

    > Durable spec: none.

### Evaluation run content

20. The regenerated run MUST cover the expanded requirement set while
    preserving 0190's run invariants: the stable
    `.quality/evaluations/0001-full-eval/` path with pinned identity, all four
    finding types, a confidence spread including at least one `low` or `none`,
    at least one requirement with multiple findings, a `not_assessed` /
    `not_rated` requirement handled correctly in roll-up, at least four ranked
    recommendations including a multi-finding one, and full finding
    disposition coverage.

    > Durable spec: none.

21. Findings MUST NOT cluster only in pre-existing areas: the codebase area
    and at least one `security` factor each receive at least one finding.

    > Rationale: an expansion that evaluates uniformly green produces no
    > report content and cannot demonstrate the new lenses working. — 0191

    > Durable spec: none.

22. Area and root roll-ups MUST remain consistent with the body's stated
    margins and veto semantics (money-touching areas at target-or-better
    posture, veto intact).

    > Durable spec: none.

### History and framing

23. The synthetic quality changelog MUST gain entries, dated before the run
    and conforming to the SKILL.md format contract, that account for the
    model's revised shape: at minimum an `add` covering the codebase area and
    security factors, and one sensor-maturation entry (an assessment moving
    from manual review to a named sensor) consistent with requirement 10's
    story.

    > Durable spec: none.

24. The gallery README MUST keep its synthetic-evidence caveat, regeneration
    and do-not-hand-edit instructions, and changelog link, and MUST state that
    assessments demonstrate the guide-and-sensor posture.

    > Durable spec: none.

### Gates

25. `mise run report-gallery` MUST regenerate the gallery byte-stable, the
    emitted `QUALITY.md` MUST pass `qualitymd lint`, and `mise run check`
    (including `report-gallery-check` and `fmt-md-check`) MUST pass with no
    diff, with the `0001-full-eval` path unchanged so Mintlify and docs links
    need no edit.

    > Durable spec: none.

## Durable spec changes

### To add

None

### To modify

None

### To rename

None

### To delete

None

## Open questions

- Whether the 35–45 requirement band (requirement 5) survives the
  setup-workflow plausibility check of requirement 17 at its upper end; if
  not, land nearer 35 and let the changelog carry more of the growth story.
