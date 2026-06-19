---
type: Change Case
title: Require factor references
description: Require every requirement to be connected to at least one factor and align requirement-to-factor terminology.
status: Done
tags: [specification, lint, terminology]
timestamp: 2026-06-19T00:00:00Z
---

# Require factor references

## Motivation

Requirements that inform no factor make analysis and reporting carry a special
case: the requirement affects a target's local rating but no factor rating.
Requiring every requirement to be connected to at least one factor removes that
special case, makes each assessment answer an explicit quality characteristic,
and keeps factor roll-ups complete.

The public terminology should reserve "characterize" for factors describing
targets. Requirements reference factors; they are the assessable expectations
through which those factors are checked.

## Scope

This change makes requirements with no factor reference a lint error, updates
the QUALITY.md format language around factor references, and clarifies that a
requirement's `factors` references are only "secondary" when the requirement is
already nested under another factor.

Implementation begins only after the change advances to `In-Progress`.

## Affected specs & docs

- [SPECIFICATION.md](../../SPECIFICATION.md) — require every requirement to be
  connected to at least one factor, retain "lens" only as helpful shorthand, and
  clarify direct target requirements versus additional factor references.
- [specs/cli/lint.md](../../specs/cli/lint.md) — add the
  `missing-factor-reference` error rule and sync rule descriptions/messages.
- [README.md](../../README.md) — sync public overview terminology where it explains
  factors and requirement references.
- [internal/scaffold/skeleton.md](../../internal/scaffold/skeleton.md) — keep
  scaffold guidance aligned with required factor references.

## Children

- [Functional spec](0028-require-characterized-requirements/spec.md)
- [Design doc](0028-require-characterized-requirements/design.md)
