---
type: Change Case
title: Require characterized requirements
description: Require every requirement to be characterized by at least one factor and align factor-association terminology.
status: Design
tags: [specification, lint, terminology]
timestamp: 2026-06-18T00:00:00Z
---

# Require characterized requirements

## Motivation

Requirements that inform no factor make analysis and reporting carry a special
case: the requirement affects a target's local rating but no factor rating.
Requiring every requirement to be characterized by at least one factor removes
that special case, makes each assessment answer an explicit quality
characteristic, and keeps factor roll-ups complete.

The public terminology should describe the mechanics with "characterize" and
"characterized by" because that better matches quality-characteristic language.
"Lens" remains useful shorthand where it improves readability, but it should not
be the primary term for the requirement-to-factor relationship.

## Scope

This change makes uncharacterized requirements a lint error, updates the
QUALITY.md format language around requirement characterization, and clarifies
that a requirement's `factors` list is only "secondary" when the requirement is
already nested under another factor.

Implementation begins only after the change advances to `In-Progress`.

## Affected specs & docs

- [SPECIFICATION.md](../SPECIFICATION.md) — require every requirement to be
  characterized by at least one factor, retain "lens" only as helpful shorthand,
  and clarify direct target requirements versus additional factors.
- [specs/cli/lint.md](../specs/cli/lint.md) — add the
  `uncharacterized-requirement` error rule and sync rule descriptions/messages.
- [README.md](../README.md) — sync public overview terminology where it explains
  factors and requirement association.
- [internal/scaffold/skeleton.md](../internal/scaffold/skeleton.md) — keep
  scaffold guidance aligned with required characterization.

## Children

- [Functional spec](0028-require-characterized-requirements/spec.md)
- [Design doc](0028-require-characterized-requirements/design.md)
