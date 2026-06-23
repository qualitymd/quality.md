---
type: Functional Specification
title: Setup context checkpoint - functional spec
description: Requirements for replacing /quality setup's final open-ended discovery questions with a compact human context checkpoint.
tags: [quality-skill, setup, ux]
timestamp: 2026-06-23T00:00:00Z
---

# Setup context checkpoint - functional spec

Companion to [Setup context checkpoint](../0072-setup-context-checkpoint.md).
This spec states what the change must do; the [design doc](design.md) covers how
the runtime skill and durable spec carry it.

The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear
in all capitals.

## Background / Motivation

Setup's human context questions ground the model in who depends on the work,
what outcomes matter, who maintains or operates it, which invisible stakeholders
exist, and which material facts are unknown. Presented as several open-ended
questions, they create poor response ergonomics: the user can easily answer only
the final broad prompt and leave earlier context uncorrected. The workflow should
preserve the same modeling dimensions while making the user review a draft,
correct it tersely, and see the highest-value missing facts.

## Scope

This change governs the discovery prompt shape for setup's human context inputs:
primary users/outcomes, maintainers/collaborators, other stakeholders, and
missing/not-agent-accessible context. It does not change the first five setup
discovery questions, the final review gate, setup's mutation boundary, the
QUALITY.md format, or CLI behavior.

## Requirements

After the first five setup discovery questions, `setup` **MUST** present a
human context checkpoint instead of asking separate open-ended questions for
primary users/outcomes, maintainers/collaborators, other stakeholders, and
missing context.

The human context checkpoint **MUST** present the repository-inferred context as
a draft for confirmation or correction, with confidence labels and evidence
notes where useful.

The checkpoint **MUST** cover primary users/outcomes,
maintainers/collaborators, other stakeholders, and missing or
not-agent-accessible context.

The checkpoint **MUST** make the user's response task correction-oriented: the
user can confirm, correct, fill in terse fragments, or point to
agent-accessible evidence the setup pass missed.

The checkpoint **MUST** state that unanswered low-confidence or not-visible
items will be recorded as Unknown, open questions, or low-confidence inference
as appropriate, not treated as confirmed facts.

The checkpoint **MUST** prioritize the highest-value missing facts in a short
list visible at the bottom of the prompt: who the evaluated thing is for, what
outcome matters most, and whether data, compliance, availability, or
business-criticality constraints exist.

`setup` **MUST NOT** end human-context discovery with a broad catch-all question
that can obscure the primary users/outcomes, maintainer/collaborator, and other
stakeholder dimensions.

`setup` **MUST** preserve source/provenance distinctions when it authors
`QUALITY.md`: user-confirmed context, repository-inferred context, and unknown
or not-agent-accessible context must remain distinguishable.

## Acceptance criteria

- The runtime setup workflow presents the human context dimensions as one draft
  checkpoint instead of four separate open-ended questions.
- The checkpoint covers primary users/outcomes, maintainers/collaborators, other
  stakeholders, and missing/not-agent-accessible context.
- The checkpoint tells the user unanswered material gaps will be recorded as
  Unknown, open questions, or low-confidence inference rather than confirmed
  fact.
- The prompt's bottom section asks for the highest-value corrections: who the
  product/entity is for, what outcome matters most, and whether data,
  compliance, availability, or business-criticality constraints exist.
- The durable setup workflow spec mirrors the runtime behavior.
- Setup model authoring still distinguishes user-provided context,
  repository-inferred context, and unknown/not-agent-accessible context.
- Markdown formatting/checks pass.

## Durable spec changes

### To add

None

### To modify

- `specs/skills/quality-skill/workflows/setup.md` - replace the four separate
  open-ended human-context discovery questions with a single context checkpoint
  that preserves the same dimensions, prioritizes high-value corrections, and
  records omitted low-confidence context honestly.

### To rename

None

### To delete

None
