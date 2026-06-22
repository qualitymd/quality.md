---
type: Functional Specification
title: Durable spec alignment - functional spec
description: Requirements for aligning durable specs with artifact-spec and behavioral-component spec guidance.
tags: [specs, quality-skill]
timestamp: 2026-06-22T00:00:00Z
---

# Durable spec alignment - functional spec

Companion to the
[Durable spec alignment](../0052-durable-spec-alignment.md) change case. This
spec states what the alignment must do; no design doc is required unless the
implementation discovers a non-mechanical restructuring decision.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear
in all capitals.

## Scope

This change aligns durable specifications, not runtime behavior. It applies the
granularity guidance from
[Writing functional specs](../../../docs/guides/write-functional-specs.md#conventions)
to the cumulative `specs/` bundle, with the `/quality` skill specs as the
primary target.

In scope:

- classifying relevant durable specs as parent specs, 1:1 artifact specs, or
  behavioral component specs;
- adding `/quality` mode behavioral component specs;
- adding behavioral component specs for the cross-mode evaluation workflow,
  reporting contract, and convention-first quality log when those sections are
  independently reviewable;
- moving component-specific normative content out of the parent `/quality` skill
  spec when a child spec becomes the better home;
- keeping cross-links, indexes, and logs accurate.

Non-goals:

- changing `skills/quality/` runtime behavior;
- changing `qualitymd` CLI behavior or generated evaluation artifacts;
- renaming existing 1:1 artifact specs that already follow the artifact naming
  convention;
- making OKF itself carry purpose-specific concept-boundary guidance.

## Requirements

### Spec inventory

The implementation **MUST** inspect the durable `specs/` bundle and identify
which current or proposed concepts are parent specs, 1:1 artifact specs, or
behavioral component specs.

The implementation **SHOULD** leave already-aligned artifact specs in place,
including report artifact specs and `/quality` runtime guide artifact specs,
unless the audit finds a concrete naming or ownership mismatch.

### Component inventory

The implementation **MUST** inventory the parent `/quality` skill spec's major
headings and classify each one as shared invariant, behavioral component,
artifact contract, example, or deferred work.

The parent **MUST** retain shared invariants such as skill identity and metadata,
argument model, safety rules, CLI ownership, runtime resource routing, common
interaction posture, cross-component invariants, examples, and deferred items.

Independently reviewable behavioral components or artifact-adjacent contracts
**MUST** move to child specs unless the implementation records a concrete reason
they remain in the parent.

### /quality mode specs

The durable `/quality` skill spec set **MUST** include behavioral component specs
for the installed mode behaviors that are independently reviewable:

- `setup`;
- `wizard`;
- `evaluate`;
- `improve`; and
- `update`.

Each mode spec **MUST** be named for the mode capability rather than the literal
runtime Markdown filename. For example, the setup mode spec is a behavioral
component spec named for `setup`, not an artifact spec named for `setup.md`.

Each mode spec **MUST** state the mode's purpose, routing conditions, mutation
surface, required artifacts, stop conditions, and verification or completion
criteria when those differ from the shared parent contract.

### Cross-mode component specs

The durable `/quality` skill spec set **MUST** include behavioral component specs
for independently reviewable cross-mode contracts:

- the evaluation workflow;
- reporting and evaluation run artifacts; and
- the convention-first quality log.

Each cross-mode component spec **MUST** state the component's purpose, scope,
required artifacts or records, mutation surface, stop or correction conditions,
and verification or completion criteria when applicable.

### Parent skill spec

The parent `/quality` skill spec **MUST** keep shared contracts that apply across
modes: skill identity and metadata, argument model, shared safety rules, CLI
ownership, runtime resource routing, common interaction posture, evaluation
semantics, report/log relationships, and cross-mode invariants.

The parent spec **MUST NOT** retain a second full copy of component-specific
requirements that have a child component spec as their source of truth. It may
keep short summaries and links when they help a reader navigate the skill as a
whole.

### Bundle navigation

The alignment **MUST** update the relevant OKF indexes so a reader can discover
the parent skill spec, mode behavior specs, guide artifact specs, and example
artifacts without opening unrelated files first.

The alignment **MUST** add a `specs/log.md` entry summarizing the durable spec
changes.

## Durable spec changes

### To add

- `specs/skills/quality-skill/modes/setup.md` - behavioral component spec for
  `/quality setup`.
- `specs/skills/quality-skill/modes/wizard.md` - behavioral component spec for
  `/quality wizard`.
- `specs/skills/quality-skill/modes/evaluate.md` - behavioral component spec for
  `/quality evaluate`.
- `specs/skills/quality-skill/modes/improve.md` - behavioral component spec for
  `/quality improve`.
- `specs/skills/quality-skill/modes/update.md` - behavioral component spec for
  `/quality update`.
- `specs/skills/quality-skill/evaluation.md` - behavioral component spec for
  the cross-mode evaluation workflow.
- `specs/skills/quality-skill/reporting.md` - component spec for evaluation
  reporting and run artifacts.
- `specs/skills/quality-skill/quality-log.md` - component spec for the
  convention-first quality log.

### To modify

- `specs/skills/quality-skill/quality-skill.md` - keep shared skill contracts in
  the parent and replace full component-specific contracts with summaries and
  links where child specs become the source of truth.
- `specs/skills/quality-skill/index.md` - list the new mode specs and clarify
  the folder's parent/spec/component/artifact/example structure.

### To rename

None.

### To delete

None.
