---
type: Change Case
title: Authoring guide replaces meta-model workflow
description: Replace the bundled quality meta-model workflow with a practical QUALITY.md authoring guide and remove public CLI/model-altitude compatibility surfaces.
status: In-Review
tags: [skill, cli, docs, evaluation]
timestamp: 2026-06-18T00:00:00Z
---

# Authoring guide replaces meta-model workflow

## Motivation

The bundled quality meta-model made model review look like another mechanical
evaluation surface, but authors need accessible guidance more than an internal
diagnostic model. Keeping the diagnostic model public also forced the skill,
CLI, and specs to preserve a `model` altitude that is now the wrong abstraction.

## Scope

Replace the skill-facing meta-model reference with a `quality-md-guide.md`
authoring guide, remove the bundled `qualitymd models` surface, and simplify
evaluation run creation so it no longer creates model-altitude runs.

## Affected specs & docs

- [README.md](../README.md)
- [specs/cli.md](../specs/cli.md)
- [specs/cli/evaluation-create-run.md](../specs/cli/evaluation-create-run.md)
- [specs/cli/index.md](../specs/cli/index.md)
- [specs/evaluation-records.md](../specs/evaluation-records.md)
- [specs/skills/index.md](../specs/skills/index.md)
- [specs/skills/quality-skill/quality-skill.md](../specs/skills/quality-skill/quality-skill.md)
- [docs/guides/use-quality-skill.md](../docs/guides/use-quality-skill.md)
- [skills/quality/SKILL.md](../skills/quality/SKILL.md)

## Children

- [Functional spec](0026-authoring-guide-remove-meta-model/spec.md)
- [Design doc](0026-authoring-guide-remove-meta-model/design.md)
