---
type: Functional Specification
title: Report gallery exemplar functional spec
description: Delta contract for expanding the LedgerLite gallery into a best-practice QUALITY.md exemplar with a realistic generated evaluation.
tags: [examples, report-gallery, authoring]
timestamp: 2026-07-02T00:00:00Z
---

# Report gallery exemplar functional spec

This spec governs the content and regeneration contract of the
[report gallery](../../examples/report-gallery/software-service/README.md).
Normative references: the authoring guide family under
[`skills/quality/guides/authoring.md`](../../skills/quality/guides/authoring.md)
(the practices the exemplar must demonstrate), the quality changelog format
contract in [`skills/quality/SKILL.md`](../../skills/quality/SKILL.md), and the
evaluation data contract enforced by `internal/evaluation` through
`qualitymd evaluation data set`. Informational: the evaluate workflow
([`skills/quality/workflows/evaluate.md`](../../skills/quality/workflows/evaluate.md))
whose finding and recommendation conventions the synthetic payloads imitate, and
[Engineering quality loops](../../mintlify/loops.mdx) for the maturation story
the changelog demonstrates.

The key words "MUST", "SHOULD", and "MAY" are to be interpreted as described in
BCP 14 when, and only when, they appear in all capitals.

## Background / motivation

The gallery is the linked-from-docs, end-to-end demonstration of QUALITY.md.
Its current model is structurally valid but practices almost none of the
authoring guidance, so the example teaches the format while silently
mis-teaching authoring. The exemplar closes that gap with one fictional service
modeled the way the guides say to model it, evaluated once through the real
report pipeline, with a curated model-change history behind it.

## Scope

Covered: the LedgerLite `QUALITY.md` (frontmatter and body), the synthetic
payloads behind the single `0001-full-eval` run, a synthetic quality changelog,
the gallery README, and the generator that emits all of it. Deferred: a second
non-software example; a checked-in fictional source tree; a multi-run loop
demonstration. Non-goals: any change to report format, record schema, CLI, or
skill behavior.

## Assumptions & dependencies

- `evaluation.CreateRun` numbers the run `0001` in an empty evaluations
  directory and accepts the pinned manifest identity; the Mintlify home page and
  `docs/guides/reporting-design.md` reference the `0001-full-eval` path.
- `evaluation.SetData` validates payloads against the evaluation data contract;
  whatever it accepts is by definition contract-conformant.
- `prettier` formats all checked-in Markdown, including generated gallery files.

## Requirements

### Model and body

1. The gallery `QUALITY.md` MUST model LedgerLite as a composite root whose
   child areas include, at minimum: the public API, a normative service-contract
   area, persistence, operations, an agent-harness area, and a `quality-md`
   self-check area.

   > Rationale: composite decomposition, a normative-artifact area, and the two
   > recurring use-context constituents are the model-structure guide's core
   > moves; an exemplar that lacks any of them cannot demonstrate the guide.

   > Durable spec: none.

2. The root MUST declare `agent-harnessability` as a model-wide umbrella factor
   with all seven sub-factors from the authoring guide
   (`agent-accessibility`, `task-specifiability`, `agent-operability`,
   `continuity`, `self-verifiability`, `enforcement-of-standards`,
   `containment-of-action`), carry no requirements on the umbrella itself, and
   encode the factor-vs-area projection boundary as YAML comments plus
   disambiguating `description` clauses on both the factor and the
   agent-harness area.

   > Durable spec: none.

3. The Markdown body MUST contain Overview, Scope, Needs, Risks, and a
   model-shape section, each following the body guide's section shape: a
   purpose-bearing opening, judgment content, explicit _Unknowns_ and _Open
   questions_ lines (using "none known"/"none" where empty), and a review
   provenance line naming a fictional human reviewer and an agent surface with
   model.

   > Durable spec: none.

4. The body MUST state the governing sense of "good", the required margin above
   `minimum` with its reason, the roll-up posture, and MUST name the veto
   requirement (requirement 6) and trace at least one concern explicitly from a
   Need and Risk through a factor to a requirement.

   > Durable spec: none.

5. The model MUST include at least one requirement with a measured `ratings`
   criterion override, at least one area-level requirement connected to three
   or more factors through one referenced assessment, and requirement
   assessments that (across the model) reference the fictional service
   contract, runbooks/docs, the authoring guide family, and computational
   sensors (deterministic checks such as contract tests, invariant tests, and
   telemetry queries an agent could run or inspect).

   > Rationale: the sensor-grounded assessments demonstrate the
   > harness-engineering posture (feedforward guides, feedback sensors) the
   > gallery is meant to model. — 0190

   > Durable spec: none.

6. The persistence area MUST carry a balance-integrity requirement written as a
   veto: its `unacceptable` boundary is sharpened (via override or criterion
   wording) and the body names its veto role.

   > Durable spec: none.

### Evaluation run content

7. The generated run MUST remain a single evaluation at the stable path
   `.quality/evaluations/0001-full-eval/` with pinned run identity, regenerated
   byte-stable by `mise run report-gallery`.

   > Rationale: `mintlify/index.mdx` and `docs/guides/reporting-design.md` link
   > to that path; `report-gallery-check` diffs the tree in CI.

   > Durable spec: none.

8. The synthetic payloads MUST include all four finding types, a confidence
   spread that includes at least one `low` or `none`, at least one requirement
   with more than one finding, at least one requirement whose assessment status
   is `not_assessed` with a `not_rated` rating that stays visible in roll-up
   limits rather than counting as a low rating, and finding copy that follows
   the evaluate workflow's finding-core field jobs without describing itself as
   synthetic inside finding fields.

   > Rationale: the "synthetic" caveat belongs to the README and evaluation
   > limits, not every sentence of copy; self-referential copy is the main
   > realism failure of the current gallery.

   > Durable spec: none.

9. The run MUST include at least four ranked recommendations with background,
   expected value, and done criteria; at least one MUST trace to two or more
   findings, and at least one MUST address the not-assessed requirement by
   restoring assessability. Finding coverage MUST give every finding a
   disposition.

   > Durable spec: none.

### History and framing

10. The generator MUST emit a synthetic quality changelog under the example's
    `.quality/changelog/` conforming to the SKILL.md format contract, with at
    least entries of kind `add`, `rename`, `recalibrate`, and
    `drift-correction`, dated before the run and consistent with the current
    model (the recalibration matching the measured override in requirement 5,
    the rename matching the guide's `agent-harnessability` naming rule).

    > Rationale: the changelog is what makes the gallery read as a model
    > matured through the quality loop rather than authored in one sitting.

    > Durable spec: none.

11. The gallery README MUST keep the synthetic-evidence caveat, explain that
    prior evaluation runs are not retained in the gallery, link the quality
    changelog, and keep the regeneration and do-not-hand-edit instructions.

    > Durable spec: none.

12. The whole gallery MUST pass the repository gates: `mise run report-gallery`
    followed by `mise run check` (including `report-gallery-check` and
    `fmt-md-check`) succeeds with no diff.

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

None — the deferred items (second domain example, fictional source tree,
multi-run loop) are recorded in Scope.
