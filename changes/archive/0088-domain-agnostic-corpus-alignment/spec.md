---
type: Functional Specification
title: Domain-agnostic corpus alignment - functional spec
description: Requirements for adding a non-software worked example, marking and cross-linking the software example corpus, completing the AGENTS.md domain-agnostic summary, and re-scoping the README modeled domain, so the repo demonstrates QUALITY.md's domain agnosticism while preserving the agent-first use context.
tags: [docs, doctrine, domain-agnostic, examples, skill]
timestamp: 2026-06-24T00:00:00Z
---

# Domain-agnostic corpus alignment - functional spec

Companion to
[Domain-agnostic corpus alignment](../0088-domain-agnostic-corpus-alignment.md).
This spec states what the example and guidance alignment must do. The doctrine it
serves is settled in
[Modeling quality across domains](../../../docs/guides/model-quality-across-domains.md),
which is the source of truth; this spec does not restate the guide's rules, it
brings repo artifacts into agreement with them. The
[design doc](design.md) settles _how_ the new example is built.

The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

QUALITY.md is meant to model quality across domains, but the easiest examples for
this repo to reach for are software. 0083 supplied the doctrine guide and one
worked non-software example _within the guide_. A multi-agent audit against that
guide found no normative violations, but confirmed that the repo's _worked_
examples are otherwise uniformly software — the README example, the
`SPECIFICATION.md` minimal example, and the `0001` skill fixture (which the report
specs render against) — so software still reads as the default domain. Two
guidance surfaces also under-state the doctrine: the `AGENTS.md` summary omits the
earned-factors rule, and the README opening pins the modeled domain to AI/agent
projects. This change closes those gaps so the example _corpus_ demonstrates
invariance and the front-door guidance states the doctrine in full, without
disturbing the project's correct agent- and skill-first use context.

## Scope

Covered: a new non-software worked example in the `/quality` reference example set;
domain-illustrative marking and guide cross-links for the existing software
examples; the `AGENTS.md` earned-factors rule; README modeled-domain re-scoping;
optional reinforcement edits; and logs for the durable changes.

Deferred / non-goals: no QUALITY.md schema, rating, roll-up, or evaluation-semantic
change; no CLI or Go change; no re-doing 0083's enumeration alignment; no change to
the `0001` example's modeled content; no rewrite of archived Change Cases or
append-only history beyond new log entries.

## Requirements

### Add a non-software worked example

The `/quality` skill reference example set **MUST** include a complete, reportable
worked example whose modeled entity is non-software, in a cite-worthy secondary
domain (documentation or written corpus, data set or data product, research or
analytical report, or service or operation) distinct from the existing
software-service example and from the guide's documentation-set worked example.

> Rationale: the repo's worked examples were uniformly software, so agnosticism was
> asserted but never demonstrated within the skill's own example set. A second,
> reportable fixture in a distinct domain makes the example set itself the proof,
> and gives the report specs a non-software artifact to render against. - 0088

The new example **MUST** use the same Model structure and runtime-record shape as
the existing example, changing only domain-carried content (the Factors,
Requirements, and Assessments), so the comparison demonstrates that the format
needs no domain-specific dialect.

The new example's Factors **MUST** be earned from the modeled entity's own needs
and risks, and **MUST NOT** be adopted wholesale from an external standard's
characteristic list.

The new example's Assessments **MUST** describe checks appropriate to a domain with
no runnable oracle — human judgment against the artifact — rather than leaning on
an executable check.

> Rationale: software's runnable check is a _proxy_ for judgment; a domain that
> removes the proxy stresses the assessment oracle in a way software does not, which
> is the point of a secondary example. - 0088

The new example **MUST** be marked illustrative and domain-scoped, consistent with
the marking requirement below.

### Mark and cross-link the existing example corpus

The reference example set's index **MUST** carry an explicit _domain-illustrative_
marking: that its examples model particular domains (a software service; the new
secondary domain) and are not the default modeled domain. This is distinct from,
and in addition to, the existing _fiction_ disclaimer about invented subjects,
revisions, and `file:line` locators.

> Rationale: telling a reader the company and locators are invented does not tell
> them the _domain_ is one of many; the audit found the fiction disclaimer standing
> in for a domain-illustrative marking it does not provide. - 0088

Repo locations that present a worked software example as a lead illustration — the
README Example QUALITY.md and the `SPECIFICATION.md` minimal example — **SHOULD**
note that the Model shape is invariant across domains and **SHOULD** link the
guide's worked non-software example, unless the surrounding text already makes the
point and a link would only clutter it.

### Complete the AGENTS.md domain-agnostic summary

The `AGENTS.md` "Quality-domain agnostic examples and agentic use context" summary
**MUST** state that Factors are earned per Model from the modeled entity's own risks
and needs, and **MUST NOT** be adopted wholesale from an external standard's
characteristic list.

> Rationale: this is the guide's most concrete and most-violated rule, and it is the
> one rule the summary omitted. `AGENTS.md` is the single source for the gitignored
> `CLAUDE.md`/`GEMINI.md` symlinks, so the omission reaches every agent that reads
> this repo's instructions. - 0088

### Re-scope the README modeled domain

The README opening **MUST** present what a QUALITY.md models as domain-agnostic, and
**MUST NOT** state or imply that AI-assistant and coding-agent projects are the
default modeled domain. It **MUST** preserve the agent- and skill-first _use
context_ (the `/quality` skill as the primary experience, the CLI as support
tooling).

> Rationale: the modeled domain and the use context are different registers (per the
> guide's decision test). Foregrounding the agentic _workflow_ is correct; pinning
> the modeled _domain_ to AI/agent projects is the flagged anti-pattern. - 0088

The README's Agent Harnessability material **MUST** read as one illustrative factor
family a project may earn for an agent-collaborated entity — framed alongside a
domain-neutral statement of why a QUALITY.md helps — not as a built-in QUALITY.md
capability or a default factor. All agentic and harness-engineering wording that
describes the use context **MUST** be preserved; this requirement re-scopes the
_modeled domain_, it does not remove agentic references.

### Optional reinforcements

The following **MAY** be included where they fit without bloating the file they
touch, and **MAY** be deferred:

- A non-software constituent bracket in the Top 10 readiness check, beside the
  existing (correctly scoped) software illustration.
- A software-neutral (or paired) rendering fixture in the report-summary spec, so
  the only factor-naming content in the CLI/reports specs does not echo the standard
  characteristic list.
- A lineage clause in `SPECIFICATION.md` stating the project takes the named
  traditions' boundaries and vocabulary, not their characteristic lists as default
  factors.
- A half-sentence in `SKILL.md` noting the skill's judgment role is the care the
  format serves.

### Record the work

The relevant `specs/`, `docs/`, and `changes/` logs and the `CHANGELOG.md` **MUST**
record the new example and the alignment edits before the case reaches `In-Review`.

## Durable spec changes

### To add

- `specs/skills/quality-skill/examples/0002-<slug>/` - the new non-software worked
  example bundle: the `model.md` and its reportable runtime trail in the same shape
  as `0001` (per _Add a non-software worked example_ above). The design doc fixes
  the bundle's exact file set and size.

### To modify

- `specs/skills/quality-skill/examples/index.md` - add the domain-illustrative
  marking, register the new example, and generalize the `0001`-specific shared note
  so it no longer reads as the bundle's only subject (per _Mark and cross-link the
  existing example corpus_ and _Add a non-software worked example_).
- `SPECIFICATION.md` - add the Appendix B model-shape invariance note and the
  optional lineage clause (per _Mark and cross-link the existing example corpus_ and
  _Optional reinforcements_). No normative format rule changes.
- `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md` - add the
  non-software constituent bracket (per _Optional reinforcements_).
- `specs/reports/report-summary-md.md` - re-cast or pair the rendering fixture (per
  _Optional reinforcements_).
- `specs/skills/quality-skill/reporting.md` - cross-reference the guide and the new
  example (per _Mark and cross-link the existing example corpus_).

### To rename

None

### To delete

None
