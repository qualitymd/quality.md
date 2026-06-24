---
type: Functional Specification
title: Quality-domain agnosticism guide and secondary illustrations - functional spec
description: Requirements for the contributor guide and aligned examples that keep QUALITY.md repo content quality-domain agnostic while preserving the agent-first use context.
tags: [docs, doctrine, domain-agnostic, examples]
timestamp: 2026-06-24T00:00:00Z
---

# Quality-domain agnosticism guide and secondary illustrations - functional spec

Companion to
[Quality-domain agnosticism guide and secondary illustrations](../0083-quality-domain-agnosticism.md).
This spec states what the contributor guide and related wording alignment must
do. No design doc is required: this change adds doctrine and examples, not a new
implementation surface.

The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as described
in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

QUALITY.md is meant to model quality across domains, but software examples are
the easiest examples for this repo to reach for. That creates a slow drift toward
software as the implied default, executable checks as the implied assessment
oracle, and familiar software factor families as implied default Factors. The
guide must make the agnosticism claim operational: contributors need to know
which secondary domains stress the format, how to frame examples as
illustrative, and how to preserve the project's agent-first use context without
turning AI or harness quality into the default modeled domain.

## Scope

Covered: contributor guidance, public/domain enumeration wording, the README
example note, the bundled authoring/setup guidance lists, and logs for those
durable changes.

Deferred / non-goals: no change to the QUALITY.md schema, rating semantics,
roll-up semantics, CLI behavior, scaffold output, or historical archived Change
Cases.

## Requirements

### Add the domain-agnostic contributor guide

The docs bundle **MUST** include a guide that explains how contributors keep
example quality-model content domain agnostic.

The guide **MUST** define why secondary domains matter, including stress axes
such as source materiality, assessment oracle, constituency, and stakes.

The guide **MUST** name a canonical secondary-domain set for worked or
substantial examples: documentation or written corpora, data sets or data
products, research or analytical reports, and services or operations.

The guide **MUST** include one complete worked non-software example that uses the
same Model shape as the software example while changing only the domain-carried
Factors and Assessments.

> Rationale: the repo needs a reusable proof point that the format carries more
> than software quality, and contributors need a small set of high-value domains
> to reach for instead of inventing examples ad hoc. - 0083

### Keep example content explicitly illustrative

Repo content that includes concrete model examples **MUST** mark them as
illustrative unless they define this project's own Model or a normative format
rule.

Substantial examples **SHOULD** pair the familiar software/product anchor with a
cite-worthy secondary domain, unless the example is brief enough that a balanced
pair would be heavier than the point it serves.

Lists of domains, contexts, Factors, Requirements, Findings, recommendations, or
quality families **MUST** be framed as illustrative, non-exhaustive, and
potentially overlapping when they could otherwise read as default taxonomy.

The guidance **MUST NOT** import an external standard's characteristic list as a
default QUALITY.md factor family.

### Preserve the agent-first use context

The guide and `AGENTS.md` **MUST** distinguish the modeled domain from the
project's agent- and skill-first use context.

Contributor guidance **MUST NOT** remove AI assistant, coding agent,
agent-accessible evidence, harness, or skill-workflow wording merely because it
is AI-specific when that wording describes how QUALITY.md is used or this
project's own operating context.

Contributor guidance **MUST** flag AI or harness wording only when it makes that
domain sound inherent to every QUALITY.md file, normal for all models, or the
default modeled domain.

### Align durable wording and routing

`AGENTS.md` and the docs guide index **MUST** route contributors to the new guide
before adding or reviewing example quality-model content.

The README example note, the non-normative `SPECIFICATION.md` lineage language,
the bundled authoring guide, and the setup workflow **SHOULD** use the same
illustrative secondary-domain vocabulary where they enumerate domains.

The relevant docs, change, and skill-spec logs **MUST** record the guide and
guidance alignment.

## Durable spec changes

### To add

None

### To modify

- `SPECIFICATION.md` - align the non-normative lineage domain list with the
  secondary-domain vocabulary (per the alignment requirement above). No
  normative format rule changes.

### To rename

None

### To delete

None
