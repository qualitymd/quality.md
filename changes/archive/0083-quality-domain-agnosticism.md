---
type: Change Case
title: Quality-domain agnosticism guide and secondary illustrations
description: Add a contributor-doctrine guide on modeling quality across domains, with stress axes, a canonical secondary-domain set, and one worked non-software example; standardize domain enumerations across docs and skill guidance.
status: Done
tags: [docs, doctrine, domain-agnostic, examples]
timestamp: 2026-06-24T00:00:00Z
---

# Quality-domain agnosticism guide and secondary illustrations

A **Change Case** to make QUALITY.md's domain agnosticism demonstrated rather
than only asserted. The repo already says a Model can describe software,
documents, data sets, services, operations, processes, and other evaluated
entities, but most durable examples still lean on software/product quality. This
change adds a contributor guide that explains how to keep examples
domain-agnostic, names a stable set of secondary knowledge-work domains, and
carries one complete non-software example.

Detail lives in:

- [Functional spec](0083-quality-domain-agnosticism/spec.md) - what the guide and
  aligned examples must cover. No design doc: this is an editorial and doctrine
  change with no implementation design to settle.

## Motivation

A domain-agnostic format needs examples that test more than the software case.
Without a stronger contributor doctrine, examples can quietly imply that
software/product quality is the default domain, that executable checks are the
normal assessment oracle, or that familiar software factor families are the
default factor set. A durable guide gives contributors a shared way to choose
secondary examples, keep illustrative content clearly marked, and preserve the
agent- and skill-first use context without making AI or harness work sound like
the default modeled domain.

## Scope

Covered:

- Add a contributor guide for modeling quality across domains.
- Define the stress axes that make a non-software domain useful as a test of the
  format.
- Name a canonical secondary-domain set for substantial examples:
  documentation or written corpora, data sets or data products, research or
  analytical reports, and services or operations.
- Include one full worked non-software example.
- Route agents and contributors to the guide from `AGENTS.md` and the guides
  index.
- Align short domain enumerations in the README, SPECIFICATION.md lineage note,
  the bundled authoring guide, and setup workflow.
- Record the durable doc and bundled skill changes in the relevant logs.

Deferred / non-goals:

- No QUALITY.md schema, rating, roll-up, or evaluation semantic change.
- No CLI or Go behavior change.
- No adoption of any external standard's factor family as a default QUALITY.md
  taxonomy.
- No rewrite of historical Change Cases or append-only history beyond new log
  entries for this work.

## Affected artifacts

### Code

None - this is documentation, doctrine, and bundled skill guidance only.

### Format spec (`SPECIFICATION.md`)

- [x] `SPECIFICATION.md` - align the non-normative lineage domain list with the
      canonical secondary-domain set. No normative format rule changes.

### Durable specs (`specs/`)

- [x] `specs/skills/quality-skill/guides/log.md` - record the bundled authoring
      and setup guide alignment. No guide-contract requirement changed.

### Durable docs

- [x] `docs/guides/model-quality-across-domains.md` - add the contributor-doctrine
      guide.
- [x] `docs/guides/index.md` and `docs/log.md` - register and log the guide.
- [x] `AGENTS.md` - route contributors and agents to the guide before adding or
      reviewing example quality-model content.
- [x] `README.md` - make the Example QUALITY.md note explicitly non-default and
      name secondary domains.
- [x] `CHANGELOG.md` - add a documentation release note.

### Bundled skill (`skills/quality/`)

- [x] `skills/quality/guides/authoring.md` - align the illustrative domain list.
- [x] `skills/quality/workflows/setup.md` - align the setup domain prompt's
      illustrative list.

### Install / scaffold

None - no scaffolded QUALITY.md content changes.

## Children

- [Functional spec](0083-quality-domain-agnosticism/spec.md) - what the guide and
  aligned examples must cover.

## Status

`Done`. Implemented, verified, and archived. The new guide, routing updates,
domain-list alignment, logs, and changelog entry are in place. No code,
format-schema, or CLI behavior changed.
