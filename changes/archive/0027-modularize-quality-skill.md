---
type: Change Case
title: Modularize quality skill modes
description: Split mode-specific /quality skill procedures into separate mode files while keeping global routing and safety rules in SKILL.md.
status: Done
tags: [skill, authoring]
timestamp: 2026-06-18T00:00:00Z
---

# Modularize quality skill modes

## Motivation

The installable `/quality` skill was becoming hard to author because routing,
global safety rules, and each mode's procedure lived in one file. Splitting mode
procedures into focused reference files makes setup, wizard, evaluation, and
improve instructions easier to maintain without changing behavior.

## Scope

Move mode-specific procedure text from `skills/quality/SKILL.md` into
`skills/quality/modes/`, keep supporting docs under `skills/quality/resources/`,
and keep `SKILL.md` as the root router with global rules, argument parsing,
shared CLI prerequisites, config, and artifact contracts.

## Affected specs & docs

- [specs/skills/quality-skill/quality-skill.md](../../specs/skills/quality-skill/quality-skill.md)
- [skills/quality/SKILL.md](../../skills/quality/SKILL.md)

## Children

- [Functional spec](0027-modularize-quality-skill/spec.md)
