---
type: Design Doc
title: Domain-agnostic corpus alignment - design doc
description: How the non-software worked example is built (its secondary domain, minimal-but-complete bundle scope, and placement), how the existing corpus is marked and cross-linked, and how the README modeled-domain re-scoping is shaped without removing agentic use-context wording. Records the domain-choice decision and the alternatives weighed.
tags: [docs, doctrine, domain-agnostic, examples, skill]
timestamp: 2026-06-24T00:00:00Z
---

# Domain-agnostic corpus alignment - design doc

Design behind the
[Domain-agnostic corpus alignment](../0088-domain-agnostic-corpus-alignment.md)
change case and its [functional spec](spec.md). The spec settles *what* must hold;
this doc settles the choices the spec deliberately leaves open — chiefly the new
example's domain and bundle scope, and the shape of the README edits.

## Context

[0083](../0083-quality-domain-agnosticism.md) added the doctrine guide and one
worked non-software example *inside it*. A multi-agent audit then confirmed the
residual pattern the guide is meant to prevent: the repo's other worked examples
are uniformly software. The spec's substantive new artifact is a second,
non-software reference example for the `/quality` skill, paired with marking and
cross-link edits and two front-door guidance fixes. Most of the case is editorial
and needs no design; three decisions do, and are settled here.

## Approach

### Decision 1 — the secondary domain: a data product

The new example models a **data set / data product**. The cite-worthy set in the
guide is documentation, data product, research report, and service. The repo
already spends two of those: the guide's own worked example is a *documentation
set*, and `0001` is a software *service* (a payments API). A data product is the
remaining choice that is maximally distinct from both, and it stresses the format's
abstractions where software does not:

- **Source materiality** — the evaluated thing is data, not a readable code path,
  so `source` points at a dataset and its schema/lineage rather than a file an agent
  reads top to bottom.
- **Assessment oracle** — quality rests on profiling and reviewer judgment against
  the data and its documentation, with no runnable correctness check, exercising the
  spec's no-executable-oracle requirement.

Concretely (final subject fixed in implementation): a small fictional reference or
analytics dataset, with two Areas such as **Schema & structure** and **Provenance &
freshness**, and factors *earned* from the entity's own risks — e.g. accuracy,
completeness, provenance, timeliness — never imported as a default set. Assessments
read as a reviewer's judgment ("a reviewer profiles column X against the documented
range and treats out-of-range values as defects"), and recommendations follow
naturally (backfill a gap, add provenance metadata, tighten a freshness SLA).

> Note for the reviewer: the case was approved with the example phrased as "a
> documentation-set or service evaluation." On reflection both duplicate an existing
> example (the guide's doc set; `0001`'s service), so a data product is the higher-value
> pick and is still squarely "a non-software fixture." This is flagged as an
> [open question](#open-questions); documentation-set and research-report are recorded
> fallbacks.

### Decision 2 — bundle scope: minimal, but complete

The example set's contract is that each fixture is "intentionally complete and
reportable: one assessment per in-scope requirement and one analysis per area node."
The new bundle honors that contract but is deliberately **small** — roughly two
Areas and four to six Requirements total — and reuses the suggested four-level
rating scale. It carries the same file set as `0001` (`model.md`, `design.md`,
`plan.md`, the `assessments/` and `analysis/` records, the generated
`report-summary.md`/`report.md`/`report.json`, and one or two recommendations),
scaled down. The generated reports are produced from the runtime records and
regenerated when inputs change, exactly as the `0001` bundle note already states.

This keeps the fixture demonstrative without doubling the maintenance surface: the
point is invariance of *shape*, which a small complete bundle shows as well as a
large one.

### Decision 3 — placement and corpus marking

The bundle lives at `specs/skills/quality-skill/examples/0002-<slug>/`, parallel to
`0001-quality-eval`, because the report specs already reference the example set
there. `examples/index.md` then gains three things: (1) a short domain-illustrative
paragraph stating the examples model particular domains and are not the default
modeled domain; (2) a new **Examples** entry for `0002`; and (3) a rewrite of the
"Shared across this bundle … the subject ('Sparrow Payments')" note so the
genuinely shared facts (the rating scale, the raw runtime shape, the fiction of
subjects and locators) stay shared while the Sparrow-specific detail moves into the
`0001` entry.

### Decision 4 — README re-scoping, by minimal edit

Two edits, both subtractive of *implication* rather than of agentic content:

- The opening sentence changes from "continuously improve AI assistant and coding
  agent projects" to a domain-agnostic object ("align your team and AI agents on
  what *good* means — for software, docs, data, services, or whatever you tend"),
  keeping the `/quality`-skill use context intact. The existing "That agentic
  workflow is the primary experience" sentence stays untouched — it is correct
  use-context framing.
- The "Evaluate and Improve Agent Harnessability" section keeps all of its harness
  references (the guide forbids removing them) but is reframed as *one example* of a
  factor family a project earns for an agent-collaborated entity, sitting beside a
  domain-neutral statement of why a QUALITY.md helps, so it and the
  software-specific "Manage Quality Debt" section read as two illustrations rather
  than the two canonical reasons.

## Alternatives

- **Fold this into 0083.** Rejected — 0083 is `In-Review` and editorial; a new
  reportable fixture plus a README rework is new scope with a real design decision.
  A separate case keeps 0083 landable and records the fixture's rationale on its own
  durable spec.
- **Replace the `0001` Sparrow example with a non-software one** (one example, not
  two). Rejected — `0001` exercises genuinely software-shaped runtime behavior
  (secret-by-reference, prompt-injection-as-data, a layered binding constraint) worth
  keeping. The goal is a *pair* that shows invariance, not swapping which single
  domain is privileged.
- **Full-parity bundle** (same size as `0001`). Rejected — more maintenance for no
  extra demonstrative value; invariance of shape does not need scale.
- **`model.md` only, no generated reports.** Rejected — the report specs render
  against the example set and the bundle contract requires reportability; a model
  with no trail would neither exercise the contract nor demonstrate end-to-end
  invariance.
- **Domain = research/analytical report instead of a data product.** Viable and kept
  as the first fallback; rejected as the primary only because it reads closer to the
  guide's documentation example, whereas a data product stresses source-materiality
  harder.
- **Put the example in the bundled runtime `skills/quality/`.** Rejected — `0001`
  lives in the `specs/` example bundle and the report specs reference it there; keep
  the new one parallel.

## Trade-offs & risks

- **Maintenance.** A second reportable bundle adds upkeep (reports regenerate from
  inputs). Mitigation: keep it small and document it as generated-from-inputs, like
  `0001`.
- **Divergence from the approved phrasing.** The case was approved naming a
  documentation-set or service fixture; this picks a data product. Mitigation:
  flagged as an open question with recorded fallbacks; the choice stays within "add a
  non-software fixture."
- **README is a near-marketing surface.** Re-scoping risks weakening the valid,
  useful agentic positioning. Mitigation: re-scope, do not remove; preserve every
  use-context and harness reference; prefer a lead-in sentence and a hedge over
  restructuring.
- **Two technical examples.** Software + data are both "technical" next to the
  guide's prose documentation example. Accepted: across the guide (documentation),
  `0001` (service), and `0002` (data product) the repo now spans three distinct
  cite-worthy contexts.

## Open questions

- **Secondary domain:** data product (recommended) vs. documentation-set or
  research-report (fallbacks) — confirm before implementation.
- **Fixture subject and size:** the concrete fictional dataset and the exact
  Area/Requirement count (target ~2 Areas, ~4-6 Requirements).
- **Optional reinforcements:** whether all four (Top 10 bracket, report-summary
  fixture, lineage clause, `SKILL.md` nod) land in this case or some defer. Current
  lean: include the cheap ones, defer none.
- **Stale `docs/log.md:18`:** correct the 0083 catalog description ("Anthropic
  knowledge-work plugins sample") here, or leave it to 0083 to fix before it lands.
