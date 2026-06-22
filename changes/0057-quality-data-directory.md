---
type: Change Case
title: Quality data directory
description: Move qualitymd and /quality support artifacts under .quality/ and add a root config pointer for workspace resolution.
status: Design
tags: [workspace, config, evaluation, lint, skill]
timestamp: 2026-06-22T00:00:00Z
---

# Quality data directory

A **Change Case** defining the QUALITY.md workspace envelope and the `.quality/`
quality data directory convention. The detail lives in its
[functional spec](0057-quality-data-directory/spec.md).

## Motivation

QUALITY.md currently splits tool support files across two namespaces:
`.quality/config.yaml` for configuration and `quality/evaluations/` plus
`quality/log/` for generated evaluation and model-history artifacts. That split
is harder to explain, consumes a plain `quality/` directory that projects may
want for authored quality material, and makes the tool-owned surface look less
cohesive than common project-local tool directories.

The change should establish a single project-local quality data directory while
preserving a larger workspace concept: the workspace is the resolved operating
context for one `QUALITY.md`; `.quality/` is only the directory where
qualitymd and the `/quality` skill store support data for that workspace.

## Scope

Covered:

- Define the QUALITY.md workspace as the resolved operating context for a
  selected `QUALITY.md` file.
- Define `.quality/` as the default quality data directory for that workspace.
- Move default evaluation records from `quality/evaluations/` to
  `.quality/evaluations/`.
- Move the default quality log from `quality/log/` to `.quality/log/`.
- Keep `.quality/config.yaml` as the default config file.
- Add optional root `config` frontmatter as a qualitymd tooling convention that
  points to the config file and is not part of the normative Model.
- Keep unknown keys as lint errors by default, with the unknown-key rule
  internally configurable so `config` is allowed at the root.

Deferred / non-goals:

- No backward compatibility routing for existing `quality/` artifact
  directories.
- No public `.quality/config.yaml` lint-rule configuration surface yet.
- No broad extension system for arbitrary frontmatter keys.
- No change to the normative QUALITY.md Model or rating semantics.
- No migration command.

## Affected artifacts

Durable **spec** changes are itemized in the functional spec's
[Durable spec changes](0057-quality-data-directory/spec.md#durable-spec-changes)
section. The index below is the full skimmable list to reconcile before
In-Review.

Code:

- [ ] New or existing internal workspace/path-resolution package - resolve the
      selected model file, repository root, config file, quality data directory,
      evaluation directory, and quality log directory from one shared path model.
- [ ] `internal/evaluation/` - default evaluation runs to
      `.quality/evaluations/`, consume the shared workspace resolver, and keep
      `--evaluation-dir` precedence.
- [ ] `internal/status/` - report and inspect evaluation history from the
      resolved workspace evaluation directory.
- [ ] `internal/lint/` - allow root `config`, validate its shape/path, and keep
      other unknown keys as errors through rule options rather than ad hoc
      schema exceptions.
- [ ] `internal/cli/` tests and command expectations - update path receipts,
      status JSON, lint findings, and help-adjacent examples as needed.
- [ ] `dprint.json` - update generated report ignore globs from
      `quality/evaluations/**` to `.quality/evaluations/**`.

Specs:

- [ ] [`specs/cli/evaluation-create.md`](../specs/cli/evaluation-create.md) -
      default and precedence for the workspace evaluation directory.
- [ ] [`specs/cli/status.md`](../specs/cli/status.md) - default status history
      path and workspace config resolution.
- [ ] [`specs/cli/lint.md`](../specs/cli/lint.md) - clarify that qualitymd lint
      applies the default qualitymd lint profile: normative model checks plus
      documented tooling-key checks.
- [ ] [`specs/cli/lint-rules.md`](../specs/cli/lint-rules.md) - add the root
      `config` validation rule and the internally configurable unknown-key rule
      behavior.
- [ ] [`specs/skills/quality-skill/reporting.md`](../specs/skills/quality-skill/reporting.md)
      - align evaluation artifact paths and workspace config resolution.
- [ ] [`specs/skills/quality-skill/quality-log.md`](../specs/skills/quality-skill/quality-log.md)
      - move the default quality log to `.quality/log/`.
- [ ] [`specs/skills/quality-skill/modes/setup.md`](../specs/skills/quality-skill/modes/setup.md),
      [`modes/evaluate.md`](../specs/skills/quality-skill/modes/evaluate.md),
      and [`modes/wizard.md`](../specs/skills/quality-skill/modes/wizard.md) -
      align mode contracts with the quality data directory.
- [ ] [`specs/skills/quality-skill/quality-skill.md`](../specs/skills/quality-skill/quality-skill.md)
      - align shared skill config and artifact terminology.
- [ ] [`specs/quality-schema-json.md`](../specs/quality-schema-json.md) -
      confirm the companion JSON Schema remains a normative-format schema and
      does not describe the qualitymd-only `config` convention.

Runtime skill and docs:

- [ ] [`skills/quality/SKILL.md`](../skills/quality/SKILL.md) - use the quality
      data directory terminology, root `config` pointer, `.quality/evaluations/`,
      and `.quality/log/`.
- [ ] [`skills/quality/modes/setup.md`](../skills/quality/modes/setup.md),
      [`modes/evaluate.md`](../skills/quality/modes/evaluate.md), and
      [`modes/wizard.md`](../skills/quality/modes/wizard.md) - update mode
      paths and artifact descriptions.
- [ ] [`skills/quality/guides/authoring.md`](../skills/quality/guides/authoring.md)
      and
      [`guides/recommendation-follow-up.md`](../skills/quality/guides/recommendation-follow-up.md)
      - update quality log paths.
- [ ] [`skills/quality/resources/cli-quick-reference.md`](../skills/quality/resources/cli-quick-reference.md)
      - update evaluation path examples and config notes.
- [ ] [`docs/guides/use-quality-skill.md`](../docs/guides/use-quality-skill.md)
      - describe the quality data directory and updated defaults.
- [ ] [`README.md`](../README.md) - update any user-facing CLI/skill artifact
      wording if the existing overview needs the new term.
- [ ] [`CHANGELOG.md`](../CHANGELOG.md) - add the unreleased entry when the
      implementation lands.

`SPECIFICATION.md` is **not** affected in this case: `config` is a qualitymd
tooling convention, not a normative Model property.

## Children

- [Functional spec](0057-quality-data-directory/spec.md) - what the workspace,
  quality data directory, config pointer, and lint behavior must require.
- [Design doc](0057-quality-data-directory/design.md) - how workspace
  resolution, `.quality/` defaults, root `config`, and strict lint rule options
  should be implemented.

## Status

`Design`. The functional spec is settled and the design doc has been created.
Implementation has not started.
