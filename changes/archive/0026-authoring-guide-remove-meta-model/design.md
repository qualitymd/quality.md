---
type: Design Doc
title: Authoring guide replaces meta-model workflow design
description: Implementation approach for removing the bundled model surface and adding a QUALITY.md authoring guide.
tags: [skill, cli, evaluation, authoring-guide]
timestamp: 2026-06-18T00:00:00Z
---

# Authoring guide replaces meta-model workflow design

## Context

This design answers the [functional spec](spec.md). The change removes a public
diagnostic model workflow and replaces the skill-facing reference with direct
authoring guidance.

## Approach

Break the skill reference symlink and add `skills/quality/references/quality-md-guide.md`
as an ordinary Markdown guide. Keep `references/SPECIFICATION.md` as the format
source of truth.

Delete the `qualitymd models` command and the `internal/models` package. Remove
the command from the root command tree and delete its command tests.

Simplify evaluation run creation by removing the CLI `--altitude` flag and
rejecting non-subject altitude values in the internal API. Keep the existing
`NNNN-subject-quality-eval` folder shape for new runs so existing report and
path parsing code needs minimal churn and historical `NNNN-model-quality-eval`
folders can still be scanned for numbering.

Sync durable specs and docs to describe subject-only evaluation and authoring
guidance instead of bundled model criteria.

## Alternatives

- **Keep `qualitymd models` but empty the catalog.** Rejected because it keeps a
  public command whose only current purpose was the removed meta-model workflow.
- **Rename the bundled model into a bundled guide command.** Rejected because
  the skill can carry a reference guide directly; the CLI should stay mechanical.
- **Rename new run folders to `NNNN-quality-eval`.** Deferred because it would
  force broader runtime/report fixture churn without materially improving the
  removal of model-altitude evaluation.

## Risks

The main compatibility risk is stale references to `qualitymd models` in docs or
tests. The mitigation is a repository-wide search for `quality-meta-model`,
`qualitymd models`, and model-altitude wording, excluding archived historical
records.
