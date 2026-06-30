---
type: Functional Specification
title: Closed-choice setup UX — functional spec
description: Requirements for numbered, recommended-first closed-choice prompts in /quality setup.
tags: [skill, quality, setup, ux]
timestamp: 2026-06-26T00:00:00Z
---

# Closed-choice setup UX — functional spec

Companion to the [Closed-choice setup UX](../0099-closed-choice-setup-ux.md)
change case. This spec states _what_ the change must do; the
[design doc](design.md) covers _how_.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119 and RFC 8174.

## Background / Motivation

Setup discovery questions are part of the product surface. When a prompt asks on
one axis and offers choices on another, or when the recommended choice is not the
fastest valid answer, the user spends effort translating the interface instead
of calibrating the quality model. Small closed-choice sets should behave like
defaults in a form: the recommended option is first, numbered, and confirmed by
`1`.

## Scope

This change governs user-facing closed-choice prompts in `/quality setup` and
the shared agent-mediated UX guidance that the skill follows. It does not change
the discovery dimensions, the model fields written to `QUALITY.md`, the human
context checkpoint, review gate, authoring behavior, feedback logs, CLI
commands, or evaluation/report behavior.

## Requirements

- For small closed-choice discovery questions, setup **MUST** present numbered
  options, put the recommended option first, mark it as recommended, and make
  `1` the shortest confirmation.

  > Rationale: the recommendation is the default path. Keeping it at `1` makes
  > the easiest answer also the most likely correct answer, while still leaving
  > alternatives visible. — 0099

- Closed-choice option labels **MUST** match the question's visible axis. If the
  stored model or setup brief uses different internal vocabulary, setup **MUST**
  map the visible answer to that internal value without making the user translate
  while answering.

  > Rationale: "How costly is poor quality?" should be answered with cost
  > choices, not tolerance labels. Internal vocabulary can remain precise without
  > leaking into the interaction. — 0099

- The risk-tolerance discovery question **MUST** keep asking about the cost of
  poor quality, but its visible closed-choice options **MUST** be cost labels
  such as high cost, moderate cost, and low cost. Setup **MUST** map those
  visible answers to the existing risk-tolerance meaning used in the setup brief
  and authored model.

- The lifecycle and Rating Scale discovery questions **MUST** also follow the
  numbered, recommended-first closed-choice pattern when rendered as fixed
  choices.

- The runtime skill's general user interaction contract **SHOULD** state the
  closed-choice rule so other `/quality` workflow prompts do not reintroduce the
  older "Options + Recommended + accept" shape when a small fixed choice set is
  used.

- The agent-mediated UX guide **MUST** document the rule with a good example and
  checklist item so future workflow design work has a durable source of truth.

## Durable spec changes

### To add

None.

### To modify

- [`specs/skills/quality-skill/quality-skill.md`](../../../specs/skills/quality-skill/quality-skill.md)
  — add the general closed-choice rule to the user interaction contract (per the
  runtime interaction contract requirement above).
- [`specs/skills/quality-skill/workflows/setup.md`](../../../specs/skills/quality-skill/workflows/setup.md)
  — update setup discovery and prompt-form requirements for numbered,
  recommended-first closed choices and the risk-cost-to-tolerance mapping (per
  the setup discovery requirements above).

### To rename

None.

### To delete

None.
