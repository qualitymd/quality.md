---
type: Functional Specification
title: Modularize quality skill modes
description: Requirements for splitting /quality mode procedures into separate reference files.
tags: [skill, authoring]
timestamp: 2026-06-18T00:00:00Z
---

# Modularize quality skill modes

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Background / Motivation

The `/quality` skill's root prompt has to be read every time the skill is used,
so it should carry the routing and always-on rules, not every mode's full
procedure. Mode procedures are easier to author and review as separate
reference files.

## Requirements

The installable skill root `SKILL.md` **MUST** remain the router. It **MUST**
retain argument parsing, shared CLI prerequisites, global safety rules, config,
and artifact-contract guidance.

Mode-specific procedures **MUST** live in separate files under
`skills/quality/modes/`:

- `setup.md`
- `wizard.md`
- `evaluate.md`
- `improve.md`

`SKILL.md` **MUST** tell the agent to read the matching mode file before
executing a mode.

Supporting docs **MUST** live under `skills/quality/resources/`.

The refactor **MUST NOT** change mode behavior.
