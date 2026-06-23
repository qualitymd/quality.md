---
type: Functional Specification
title: Setup missing-context provenance — functional spec
description: Requirements for /quality setup missing-context discovery provenance.
tags: [quality-skill, setup, agent-accessibility]
timestamp: 2026-06-23T00:00:00Z
---

# Setup missing-context provenance — functional spec

Companion to
[Setup missing-context provenance](../0070-setup-missing-context-provenance.md).
This spec states what the change must do; the [design doc](design.md) covers how
the runtime skill and durable spec carry it.

The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear
in all capitals.

## Background / Motivation

Setup turns repository-visible evidence and user corrections into a first useful
`QUALITY.md`. When material context is absent from the repo, setup must keep that
absence visible. A prompt option that says low-evidence product purpose or
operational context is "sufficiently understood" converts tacit knowledge into
model evidence without provenance, which is exactly what the missing-context
question is meant to prevent.

## Scope

This change governs `/quality setup` missing-context discovery, including
recommended choices and any generated structured-question options. It covers
material project-specific context that setup has identified as low-confidence,
not visible, or not agent-accessible.

It does not require setup to discover every possible external fact, add CLI
prompt generation, or change the QUALITY.md format.

## Requirements

`setup` **MUST** treat material context as agent-accessible only when it is
available through repository content, cited local paths, configured tools,
linked public sources, or explicit user-provided setup context.

`setup` missing-context discovery **MUST** distinguish context that is visible
from agent-accessible evidence, context that should be recorded as unknown or
not agent-accessible, and context the user provides during setup.

Generated missing-context choices **MUST NOT** invite the user to assume that
product purpose, operational context, stakeholder needs, telemetry,
security/compliance posture, incident history, SLAs, production metrics, or
similar material project-specific facts are understood when the setup brief marks
them `Low` confidence or not visible from evidence.

> Rationale: the missing-context question exists to prevent guessing. If setup
> already marked a project-specific fact as low/no-evidence, an "assume it is
> understood" option silently turns tacit maintainer knowledge into evidence the
> agent cannot inspect. — 0070

When a generated option excludes an identified material gap from unknowns, that
option **MUST** make the provenance explicit: either the user is providing the
missing context during setup, or the user is pointing setup to agent-accessible
evidence it missed.

For material gaps with low or no evidence, the recommended option **SHOULD**
record the gaps as unknowns or open questions.

When a user provides missing context during setup, the generated `QUALITY.md`
body **MUST** preserve that provenance clearly enough that a later reader can
tell the context came from explicit setup input rather than repository
inspection.

The setup workflow **MUST NOT** treat tacit maintainer, operator, contributor,
or stakeholder knowledge as available evidence unless it has been explicitly
provided or cited.

## Acceptance criteria

- Given setup identifies product purpose as `Low` because the README only names a
  product or repository without explaining what it does, the missing-context
  prompt does not offer an option that says product purpose is sufficiently
  understood unless that option asks the user to provide or cite the missing
  context.
- Given setup identifies runtime telemetry, security/compliance posture, incident
  history, SLAs, and production metrics as not visible, the recommended choice
  records them as unknown/not-agent-accessible unless the user provides context
  or evidence.
- Given the user provides product or operational context during setup, the
  generated `QUALITY.md` includes that context and makes its setup-provided
  provenance visible while preserving any remaining unavailable telemetry,
  security/compliance, incident, SLA, or production-metric inputs as unknowns or
  open questions.
- Given the user wants to exclude an identified gap from unknowns, setup records
  the agent-accessible evidence or explicit provided context that justifies
  excluding it.
- Transcript review of `/quality setup` can show that each material
  missing-context gap was recorded, explicitly answered by the user, or tied to
  evidence; no material gap is silently dropped.

## Durable spec changes

### To add

None

### To modify

- `specs/skills/quality-skill/workflows/setup.md` — add the missing-context
  provenance requirements above to the setup workflow contract.

### To rename

None

### To delete

None
