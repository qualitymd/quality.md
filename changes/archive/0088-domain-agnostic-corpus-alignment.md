---
type: Change Case
title: Domain-agnostic corpus alignment
description: Close the residual gaps a content audit found against the 0083 domain-agnosticism guide — add a complete non-software worked example to the /quality example set, mark and cross-link the software example corpus, add the factors-earned rule to AGENTS.md, and re-scope the README modeled domain — all while preserving the agent-first use context.
status: Done
tags: [docs, doctrine, domain-agnostic, examples, skill]
timestamp: 2026-06-24T00:00:00Z
---

# Domain-agnostic corpus alignment

A **Change Case** to bring the repo's example corpus and front-door guidance into
alignment with the doctrine guide added by
[0083](0083-quality-domain-agnosticism.md). 0083 made the agnosticism claim
_operational_ — it added
[Modeling quality across domains](../../docs/guides/model-quality-across-domains.md)
and aligned the first-pass domain enumerations. A multi-agent content audit
against that now-authoritative guide then found the repo carries **no normative
(P0) violations** and strong register and care-not-conformance discipline — but a
recurring residual pattern the guide exists to prevent: every _worked_ example in
the repo is software, several are the _same_ one, and a few high-traffic guidance
surfaces still under-state the doctrine.

Each individual example is fine — illustrative-marked, with factors visibly
earned. The gap is at the level of the _corpus_: the guide's own documentation-set
worked example is currently the only non-software counterweight to the README
example, the `SPECIFICATION.md` minimal example, and the `0001` skill fixture
(rendered again by the report specs). This case closes that gap so the example set
_demonstrates_ invariance rather than asserting it, and brings the
highest-leverage guidance (the `AGENTS.md` summary, the README front door) into
full agreement with the guide.

Detail lives in:

- [Functional spec](0088-domain-agnostic-corpus-alignment/spec.md) - what the
  alignment must do.
- [Design doc](0088-domain-agnostic-corpus-alignment/design.md) - how the new
  non-software example is built (domain choice, bundle scope, placement) and how
  the README and positioning edits are shaped without losing the agent-first use
  context.

## Motivation

A domain-agnostic format proven only against software is not proven. 0083 supplied
the doctrine and one worked non-software example _inside the guide_; the rest of
the repo's worked examples remain software, so a contributor reading the README,
the spec's minimal example, or the skill's reference evaluation still sees
software as the de facto "real" example. The audit confirmed this is the dominant
residual pattern, alongside two cheap, high-leverage guidance gaps: the `AGENTS.md`
domain-agnostic summary omits the guide's most concrete and most-violated rule
(factors are earned per Model; never adopt a standard's characteristic list), and
the README opening pins the _modeled domain_ to AI-assistant and coding-agent
projects rather than keeping that agentic framing in the _use context_ where it
belongs. `AGENTS.md` is the single source for the gitignored `CLAUDE.md`/`GEMINI.md`
symlinks, so the summary gap propagates to every agent reading this repo's
instructions — the highest-leverage single fix in the set.

## Scope

Covered:

- Add one complete, reportable non-software worked example to the `/quality` skill
  reference example set, in a cite-worthy secondary domain distinct from both the
  existing software-service example (`0001`) and the guide's documentation-set
  example.
- Give the existing software example corpus an explicit _domain-illustrative_
  marking — distinct from the current _fiction_ disclaimer about invented subjects
  and locators — and cross-link the guide's worked non-software example.
- Add the factors-earned-per-Model / no-default-characteristic-list rule to the
  `AGENTS.md` domain-agnostic summary.
- Re-scope the README so the opening presents the modeled domain as agnostic
  (preserving the agent- and skill-first use context) and the Agent Harnessability
  material reads as one illustrative factor family a project _earns_ for an
  agent-collaborated entity, with a domain-neutral counterweight — not a built-in
  capability or default factor.
- Optional reinforcements: a non-software bracket in the Top 10 readiness check, a
  domain-illustrative note on the `SPECIFICATION.md` Appendix B example and a
  lineage clause, a software-neutral rendering fixture in the report-summary spec,
  and a care-vs-mechanism nod in `SKILL.md`.
- Record the new example and alignments in the relevant logs and the changelog.

Deferred / non-goals:

- No QUALITY.md schema, rating, roll-up, or evaluation-semantic change.
- No CLI or Go behavior change.
- No re-doing 0083's enumeration alignment; this builds on it.
- No adoption of any external standard's factor family as a default QUALITY.md
  taxonomy.
- No change to the existing `0001` example's modeled content — it is correct and
  stays; only its bundle-level marking and cross-links change.
- No rewrite of historical Change Cases or append-only history beyond new log
  entries for this work; the frozen archived "stewardship lens/factor" wording in
  `0074`/`0077` stays as history.

## Affected artifacts

Derived by analysis, not recall: a multi-agent audit swept the whole repo — README
and identity docs, `SPECIFICATION.md`, the bundled `skills/quality/`, the
`specs/skills/quality-skill/` bundle (including the `0001` example), the
`specs/cli`/`reports`/`evaluation-records` specs, and `docs/` + `changes/` — against
the two normative sections of the guide. Grouped by kind below; empty kinds are
deliberate.

### Code

None - documentation, doctrine, bundled-skill content, and spec-example content
only.

### Format spec (`SPECIFICATION.md`)

- [x] `SPECIFICATION.md` - add a one-sentence model-shape invariance note to the
      Appendix B minimal example, and (optional) a lineage clause that the project
      takes the named traditions' boundaries and vocabulary, not their
      characteristic lists as default factors. No normative format rule change.

### Durable specs (`specs/`)

The functional spec's
[Durable spec changes](0088-domain-agnostic-corpus-alignment/spec.md#durable-spec-changes)
section is authoritative. In summary:

- [x] `specs/skills/quality-skill/examples/0002-<slug>/` - add the new non-software
      worked example bundle (model plus its reportable runtime trail).
- [x] `specs/skills/quality-skill/examples/index.md` - add the domain-illustrative
      marking, register the new example, and generalize the `0001`-specific shared
      note.
- [x] `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md` - add a
      non-software constituent bracket beside the software illustration.
- [x] `specs/reports/report-summary-md.md` - re-cast (or pair) the software-only
      rendering fixture so the rendered factor names do not echo the standard
      characteristic list.
- [x] `specs/skills/quality-skill/reporting.md` - cross-reference the guide and the
      new example to reinforce invariance.
- [x] the relevant `specs/` log(s) recording the above.

### Durable docs

- [x] `AGENTS.md` - add the factors-earned-per-Model / no-default-characteristic-list
      rule to the domain-agnostic summary (the single source for the gitignored
      `CLAUDE.md`/`GEMINI.md` symlinks).
- [x] `README.md` - re-scope the opening modeled-domain sentence; frame the Agent
      Harnessability section as an illustrative, earned factor family with a
      domain-neutral counterweight; link the guide from the Example QUALITY.md.
- [x] `docs/log.md` - record the alignment (and, if touched, correct the stale 0083
      catalog description at `docs/log.md:18`).
- [x] `CHANGELOG.md` - a documentation/example release note.

### Bundled skill (`skills/quality/`)

- [x] `skills/quality/SKILL.md` - optional care-vs-mechanism half-sentence dropped
      to keep the runtime prompt concise; no skill contract change required.

### Install / scaffold

None - no scaffolded QUALITY.md content changes.

## Children

- [Functional spec](0088-domain-agnostic-corpus-alignment/spec.md) - what the
  alignment must do.
- [Design doc](0088-domain-agnostic-corpus-alignment/design.md) - how the new
  example and the positioning edits are built.

## Status

`Done`. Implemented and archived. Added the `0002-city-bike-stations-quality-eval`
non-software data-product fixture with the same reportable runtime artifact shape
as `0001`; marked the example corpus as domain-illustrative; added the earned
Factors rule to `AGENTS.md`; re-scoped the README modeled-domain framing while
preserving the agent-first use context; added the Appendix B invariance note and
lineage clause in `SPECIFICATION.md`; reinforced the Top 10, reporting, and
report-summary specs; and recorded the work in the relevant logs and
`CHANGELOG.md`. No CLI or Go behavior changed.
