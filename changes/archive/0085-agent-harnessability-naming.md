---
type: Change Case
title: Agent Harnessability naming
description: Rename the harnessability factor guidance to Agent Harnessability, with a clearer accountability-preserving definition and current-model authoring guidance.
status: Done
tags: [skill, authoring, setup, factors, harnessability, agentic]
timestamp: 2026-06-24T00:00:00Z
---

# Agent Harnessability naming

A **Change Case** to rename the model-wide harnessability factor introduced in
0081 to **Agent Harnessability** for human-facing guidance and generated models,
using `agent-harnessability` as the recommended stable factor key.

The rename keeps the 0081 modeling structure: the factor remains the
agent-collaboration concern's model-wide quality projection, the **agent
harness** remains the constituent projection, and the **agent** remains the
audience. What changes is the label and definition. The new wording should teach
that the factor is about how the project equips an AI agent to work from the
project's own materials and tooling while preserving human direction, review, and
accountability.

Detail lives in:

- [Functional spec](0085-agent-harnessability-naming/spec.md) - what the guidance
  must say.
- [Design doc](0085-agent-harnessability-naming/design.md) - why the display name,
  key, and accountability wording are the right shape.

## Motivation

The current factor name, **harnessability**, is compact and tied to harness
engineering, but it is not self-explanatory in a model. It can read as a property
of the harness artifact itself rather than the agent-facing quality of the whole
project. That creates friction precisely where the factor is supposed to teach:
new model readers have to infer that harnessability means "equipped for agent
work," not "quality of the agent harness area."

The current short definition also leans on "limited human attention" and related
unsupervised-work language. That names a real scarcity, but it can imply that the
goal is reducing human responsibility. The factor should instead make the
responsibility boundary explicit: good agent harnessability reduces avoidable
re-specification and review toil while keeping human direction, review, and
accountability visible.

## Scope

Covered: rename the human-facing factor title to **Agent Harnessability**; use
`agent-harnessability` as the recommended key for newly authored or revised
models; update the factor definition to preserve human direction, review, and
accountability; keep the 0081 six-sub-factor decomposition; update setup,
authoring, and Top 10 guidance plus their spec mirrors; and update README and
CHANGELOG wording.

Deferred / non-goals: no change to the QUALITY.md format or schema, no CLI or Go
behavior change, no mass rewrite of historical Change Cases or append-only logs,
and no automatic migration of existing `harnessability` factors. Existing models
that already carry the old key remain structurally valid; model-authoring work can
rename them opportunistically when it is already editing the model.

## Affected artifacts

Derived by sweeping for `harnessability`, `Harnessability`, `limited human
attention`, `synchronous supervision`, and agent-collaborated composite guidance
across live skill guidance, spec mirrors, README/package docs, CHANGELOG,
`SPECIFICATION.md`, and CLI code. Grouped by kind; empty kinds are deliberate.

### Code

None - factor naming is skill authoring guidance and model content, not CLI or Go
behavior.

### Format spec (`SPECIFICATION.md`)

None - the format permits arbitrary factor keys and titles; this change governs
the `/quality` skill's authoring judgment, not the model schema.

### Durable specs (`specs/`)

The functional spec's [Durable spec changes](0085-agent-harnessability-naming/spec.md)
section is the authoritative breakdown. In summary:

- `specs/skills/quality-skill/guides/authoring-md.md` - rename the factor guidance
  to Agent Harnessability, require the accountability-preserving definition, and
  record the recommended `agent-harnessability` key.
- `specs/skills/quality-skill/workflows/setup.md` - require setup to propose
  `agent-harnessability` / Agent Harnessability for new agent-collaborated
  composite roots.
- `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md` - align the
  missing-factor check with Agent Harnessability while recognizing legacy
  `harnessability` coverage as stale naming rather than absence.
- `specs/skills/quality-skill/guides/log.md` and
  `specs/skills/quality-skill/workflows/log.md` - record the durable spec and
  workflow guidance revision when implemented.

### Durable docs

- `README.md` - update the Why QUALITY.md use case so it names Agent
  Harnessability and uses the new accountability-preserving definition.
- `npm/quality.md/README.md` - generated from the root README during package
  build; verify or regenerate if the implementation phase keeps the checked-in
  package README in sync.

### Bundled skill (`skills/quality/`)

- `skills/quality/guides/authoring.md` - runtime counterpart of the authoring-md
  spec changes.
- `skills/quality/workflows/setup.md` - runtime counterpart of the setup spec
  changes.
- `skills/quality/guides/top-10-quality-md-checks.md` - runtime counterpart of the
  Top 10 check.

### Install / scaffold

None.

### Changelog

- `CHANGELOG.md` - revise the unreleased harnessability entry to name Agent
  Harnessability and the accountability-preserving definition.

## Children

- [Functional spec](0085-agent-harnessability-naming/spec.md) - what the guidance
  must say.
- [Design doc](0085-agent-harnessability-naming/design.md) - display name, key,
  legacy handling, and accountability wording rationale.

## Status

`Done`. Landed and archived after implementation across the live skill guidance,
durable spec mirrors, README surfaces, and CHANGELOG. No CLI, Go, or format-spec
change was needed.
