---
type: Design Doc
title: Domain constituent kinds and stewardship concerns — design doc
description: The two-axis constituent-kind taxonomy (stewardship concerns × audience/purpose), the secure/safeguard direction-of-harm split, the three-projections rule, and where each edit lands.
tags: [skill, guide, authoring, areas, constituents]
timestamp: 2026-06-24T00:00:00Z
---

# Domain constituent kinds and stewardship concerns — design doc

Design behind the
[Domain constituent kinds and stewardship concerns](../0076-domain-constituent-kinds.md)
change case and its [functional spec](spec.md).

## Context

[0074](../0074-composite-root-areas.md) named the composite root shape
and the two recurring use-context constituents. It deliberately left domain
constituents to "vary with what is modeled" and pointed at "Ground high-leverage
concerns in normative artifacts" for the rest. In practice that is not enough: a
setup-authored model enumerates constituents by walking the repository's folders,
so it models the services it can see and silently omits the constituent _kinds_
that have no folder — tests that are thin, specs that do not exist yet, docs
beyond a README, a threat model that was never written. The
[acquire-roi model](../0076-domain-constituent-kinds.md#motivation) that prompted
this case modeled `server`, `db`, harness, and self-check, deferred the other
services in Scope, and never mentioned tests or specs at all — they were not
deferred, not unknowns, just invisible. See the change case
[motivation](../0076-domain-constituent-kinds.md#motivation) for the full problem
statement; the spec states the required content.

The asymmetry to fix: factor coverage is prescribed sharply, use-context
constituents are prescribed sharply, but domain-constituent coverage is left to
improvisation. The guide needs a _generator_ for constituent kinds that is
domain-agnostic in form yet inferable once a domain is named.

## Approach

### Two axes generate the kinds

A constituent kind is found at the intersection of two questions, neither of
which names a domain:

1. **Which stewardship concern leaves an artifact here?** Caring for any entity
   carries a recurring set of concerns; each tends to leave an authored,
   inspectable artifact that is a candidate constituent.
2. **Whom does that artifact serve, and for what job?** The same concern splits
   into several constituents when they serve different audiences or purposes.

Naming the axes by _function_ keeps them domain-neutral. "Enable its audiences"
rather than "docs" — because "docs" already presumes a written-document domain; a
data set's enabling artifact is a data dictionary, a service's is a runbook.

### Axis A — stewardship concerns, in two bands

The lifecycle band is roughly sequential; the protective band is cross-cutting
and bidirectional. The split matters: a linear "phase" framing cannot hold a
concern that is woven through every phase.

| Band       | Concern       | What it leaves (software instance)                        | Conventional factors                         |
| ---------- | ------------- | --------------------------------------------------------- | -------------------------------------------- |
| Lifecycle  | discover      | research notes, design explorations, decision records     | traceability, credibility                    |
| Lifecycle  | define        | requirements, specs, interface contracts, schema          | completeness, consistency, traceability      |
| Lifecycle  | realize       | source, services                                          | reliability, performance, maintainability    |
| Lifecycle  | verify        | tests, validation suites                                  | coverage, determinism, maintainability       |
| Lifecycle  | enable        | tutorial / how-to / reference / explanation (Diátaxis)    | currentness, completeness, understandability |
| Lifecycle  | operate       | runbooks, deploy config, monitoring-as-code, SLOs         | reliability, operability, observability      |
| Lifecycle  | maintain      | migrations, upgrade & deprecation guides, changelog       | maintainability, modifiability, currentness  |
| Protective | **secure**    | threat model, security policy, secrets/auth config        | security                                     |
| Protective | **safeguard** | safety case, privacy/impact assessment, output guardrails | safety, privacy, compliance                  |

The table is the part that makes the kinds _inferable_ — it is the answer to the
original ask ("we should be able to infer what these composite parts are from the
domain"). It ships as one illustrative, domain-scoped column, explicitly a prompt
rather than a checklist.

### secure vs. safeguard — the direction of harm

The protective band is two concerns, not one, split by who is being protected
from whom:

- **secure** — guard the entity _from_ the world: breach, tampering, injection,
  theft. Inward-facing.
- **safeguard** — guard stakeholders and the environment, internal and external,
  _from_ the entity: a destructive action, an unsafe output, privacy harm to a
  data subject. Outward-facing.

They are orthogonal — an entity can be hardened against attackers yet routinely
harm its own users, and vice versa — and they generalize cleanly: a data set
(secure: breach; safeguard: privacy/bias harm to subjects); an agent harness
(secure: injection/tampering; safeguard: harmful actions or outputs); a service
(secure: attack; safeguard: failure that injures the public). The pair earns its
keep in QUALITY.md's own agentic context of use, where the harness constituent
has a clean secure-vs-safeguard reading the guidance currently has no vocabulary
for. Privacy and compliance stop being their own band entries and _distribute_
across the pair (data-breach-privacy is secure; subject-harm-privacy is
safeguard).

### Axis B — audience × purpose

The same concern yields several constituents when audience or purpose diverges.
[Diátaxis](https://diataxis.fr/) is exactly this axis applied to the _enable_
concern: tutorial / how-to / reference / explanation are four constituents with
different factor families (a tutorial's quality is "does a novice succeed";
reference's is "accurate, complete, fast to look up"), not one "documentation"
area. The same split recurs elsewhere — a maintainer-facing CI gate and a
user-facing acceptance suite are two _verify_ constituents.

Axis B is wired to the body's **Needs**, which already enumerate stakeholders, so
it introduces no new roster: each audience a Need names should have an enabling
and verifying constituent that is modeled or accounted for. Axis B is also _how
you discover a constituent's internal shape_ — a constituent that fans out by
audience or purpose is itself composite or a collection, which nests directly
into 0074's decomposition-shape material rather than competing with it.

### The three-projections rule resolves the factor collision

`operate`, `maintain`, and `secure` are also factor names. Left unaddressed, Axis
A would appear to rename qualities as activities — exactly the "name the quality,
not the practice" trap one level down. The resolution is that a fundamental
concern **projects** into the model up to three ways:

- as a **factor** — the quality lens applied across areas (`security` of the
  server);
- as a **constituent** — the artifact that pursues it (a security policy);
- as an **audience** — who it serves or protects (an auditor).

They share a name because they share a root concern, not because they duplicate.
The author names the projection being modeled and models it once; the security
_of_ the server is a factor on the server area, while the security policy is its
own area. That is the rule that keeps the table from double-counting against the
existing factor guidance.

### Where it lands

- `authoring.md`, "Area" section: new subsection **"Cover the domain's
  constituent kinds"**, placed between "Ground high-leverage concerns in
  normative artifacts" and "Carry the recurring use-context constituents" — so
  the reading order is normative artifacts → domain constituents → use-context
  constituents. Carries the two axes, the secure/safeguard pair, the
  three-projections rule, the earn-it guardrail, and the illustrative table.
- `authoring.md`, "Carry the recurring use-context constituents": a one-line
  cross-reference where it already names domain constituents.
- `skills/quality/workflows/setup.md`: a model-building bullet that enumerates
  constituent kinds for a composite root and accounts for each.
- `top-10-quality-md-checks.md`, Check 8: a missing-domain-constituent routing
  finding plus a condensed-checklist line, both earned-not-roster.
- Guide specs `authoring-md.md` and `top-10-quality-md-checks-md.md` plus their
  `log.md`, and a `CHANGELOG.md` note.

## Alternatives

- **A flat list of expected areas (tests, docs, specs, security).** Rejected —
  reproduces the "universal roster of areas every model must carry" anti-pattern
  the guide explicitly warns against, and is not domain-agnostic. The two-axis
  generator yields the same instances for software while staying a prompt
  inferable from any domain.
- **A single "explanation" / "docs" constituent.** Rejected — it collapses the
  four Diátaxis modes (different audiences, different factor families) into one
  area and leaks a written-document domain. Axis B replaces it.
- **`secure` as a lifecycle stage.** Rejected — security is cross-cutting, woven
  through every stage, not a position in a sequence. Modeled as a protective
  cross-cutting concern instead.
- **One protective concern named `security`.** Rejected — it hides the
  direction-of-harm distinction (protect the entity vs. protect others from it),
  which has different stakeholders, failure modes, and governing artifacts, and
  is especially load-bearing for agentic projects.
- **Rename Axis A entries to avoid factor-name overlap (e.g. "running" instead of
  "operate").** Rejected — the overlap is real and informative; the
  three-projections rule explains it rather than hiding it behind invented
  vocabulary.
- **Promote the taxonomy to `SPECIFICATION.md`.** Deferred — these are authoring
  heuristics, not format semantics; keeping them guide-only avoids
  over-constraining the format, consistent with 0074.
- **Full treatment in getting-started.md.** Rejected — that guide is deliberately
  thin for starter models ("pick two to five factors"); importing the table
  there fights its purpose. At most an optional one-line pointer.

## Trade-offs & risks

- **Added vocabulary.** Nine concerns plus two axes plus the three-projections
  rule is a lot. Mitigated by one additive subsection and a single table; a
  primary-subject author can skip it, and the table is a prompt, not required
  reading.
- **Over-proliferation.** Nine concerns × audience fan-out invites a sprawling
  area tree. Mitigation: the earn-it guardrail is restated verbatim (owned
  inspectable artifact, divergent factor family, traced to a Need or Risk; prompt
  not quota), and the audience×purpose split routes into the existing
  composite/collection shapes so "documentation" stays one node until its modes
  genuinely diverge.
- **Factor/constituent double-counting.** The very overlap that makes the table
  inferable risks counting `security` twice. Mitigation: the three-projections
  rule is stated explicitly with the security example.
- **Domain leak.** A table with a software column risks reading as the default
  domain. Mitigation: the column is marked illustrative and domain-scoped per
  `AGENTS.md`, and the axis terms are functions, not artifacts.

## Open questions

- Does `operate` plus `maintain` over-fragment for most models, or is the
  distinction (run it today vs. keep it changeable over time) worth two entries?
  Kept distinct here because they map to two stakeholders the Needs already name.
- Should the protective pair eventually graduate into named factors in the
  stable-stakes factor guidance, or stay constituent-side only? Left as factors
  by reference for now via the three-projections rule.
- Is `discover` reliably an inspectable artifact across domains, or is it often
  only an unknown to record? Treated as a candidate kind that is frequently
  carried as a finding rather than a populated area.
