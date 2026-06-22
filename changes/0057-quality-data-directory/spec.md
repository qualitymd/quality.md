---
type: Functional Specification
title: Quality data directory
description: Define the QUALITY.md workspace, .quality/ data directory, config pointer, and strict lint behavior for tooling keys.
tags: [workspace, config, evaluation, lint, skill]
timestamp: 2026-06-22T00:00:00Z
---

# Quality data directory

This Change Case spec defines the delta for qualitymd and the `/quality` skill:
resolve one QUALITY.md workspace for a selected model file, use `.quality/` as
that workspace's default quality data directory, and allow a root `config`
frontmatter pointer without making it part of the normative QUALITY.md Model.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background / Motivation

The existing defaults split a single tooling concern across `.quality/` and
`quality/`: configuration is under `.quality/config.yaml`, while generated
evaluation records and quality-log entries live under `quality/evaluations/` and
`quality/log/`. A single quality data directory makes the tool-owned surface
easier to recognize, keeps the plain `quality/` namespace available for project
authored material, and matches common project-local tool directory conventions.

The directory name alone is not the whole concept. Tools still need a broader
workspace envelope: the selected `QUALITY.md`, the repository root used for
safe path resolution, the config file, and the data directories the CLI and
skill read and write.

## Scope

Covered:

- QUALITY.md workspace resolution for the selected model file.
- The default `.quality/` quality data directory.
- The default `.quality/config.yaml` config file.
- The optional root `config` frontmatter key that points to the config file.
- Evaluation directory defaults and precedence.
- Quality log default location.
- Strict unknown-key lint behavior with internal per-rule configuration.
- Runtime skill and durable spec terminology updates.

Deferred / non-goals:

- No migration command and no automatic fallback to `quality/`.
- No public lint-rule configuration surface in `.quality/config.yaml`.
- No broad extension registry for frontmatter keys.
- No normative QUALITY.md Model change.
- No rating, report, or evaluation-record schema change beyond path defaults.

## Terminology

A **QUALITY.md workspace** is the resolved operating context for one selected
`QUALITY.md` file. It includes the selected model path, repository root, config
path, quality data directory, evaluation directory, and quality log directory.

The **quality data directory** is the project-local directory where qualitymd
and the `/quality` skill keep support data for a QUALITY.md workspace. It
defaults to `.quality/`.

The **config file** is the YAML file that configures qualitymd tooling behavior
for a workspace. It defaults to `.quality/config.yaml`.

## Workspace Resolution

Tools that need workspace paths **MUST** resolve a QUALITY.md workspace from the
selected model file. When no model file is supplied, the selected model file is
`QUALITY.md` in the current working directory.

Workspace path resolution **MUST** find the repository root from the selected
model file before resolving config-owned artifact paths.

The root frontmatter key `config` **MAY** be present on the selected
`QUALITY.md` file. When present, it points to the workspace config file.

The `config` value **MUST** be a non-empty scalar repository-relative path.
It **MUST NOT** be absolute and **MUST NOT** escape the repository after path
normalization.

If `config` is absent, the config path defaults to `.quality/config.yaml`.

If the resolved config file is absent, tools **MUST** use built-in defaults.

The default quality data directory is `.quality/`.

Configured artifact paths **MUST** be repository-relative normalized paths.
They **MUST NOT** be absolute and **MUST NOT** escape the repository.

## Evaluation Directory

The default evaluation directory is `.quality/evaluations/`.

Evaluation creation and readers **MUST** resolve the evaluation directory using
this precedence:

1. explicit command override, such as `--evaluation-dir`;
2. `evaluationDir` in the resolved config file;
3. `.quality/evaluations/`.

`qualitymd evaluation create` **MUST** create numbered run folders under the
resolved evaluation directory.

`qualitymd status` **MUST** inspect evaluation history from the same resolved
evaluation directory.

Status output **MUST** report `.quality/evaluations` as the default evaluation
history path when no override or config file changes it.

## Quality Log

The default quality log directory is `.quality/log/`.

The `/quality` skill **MUST** write setup-created inaugural entries and confirmed
model-change log entries under the resolved workspace's quality log directory.

This change does not introduce a public `logDir` config key. A configurable log
directory remains deferred.

## Lint Behavior

`config` is a qualitymd tooling convention, not a normative Model property.
`qualitymd lint` **MUST** accept `config` at the root of the selected
`QUALITY.md` when its value satisfies the config-path requirements above.

`qualitymd lint` **MUST** report an error when root `config` is empty, not a
scalar, absolute, or repository-escaping.

Other unknown root keys **MUST** remain error findings by default.

Unknown nested keys inside Areas, Factors, Requirements, and Rating Levels
**MUST** remain error findings by default.

The unknown-key lint rule **MUST** be internally configurable by rule options so
allowed keys can be declared without hard-coding one-off exceptions in the
schema traversal. The initial default allowed root key outside the normative
Model schema is `config`.

No user-facing lint-rule configuration file syntax is specified by this change.

The companion JSON Schema emitted by `qualitymd schema` remains a structural
schema for the normative QUALITY.md frontmatter. This change **MUST NOT** make
the schema command describe the qualitymd-only `config` convention as a
normative Model property.

## Acceptance Criteria

- With no config file, `qualitymd evaluation create --model QUALITY.md` creates
  `.quality/evaluations/0001-quality-eval`.
- With no config file, status reports the evaluation history path as
  `.quality/evaluations`.
- With `config: .quality/config.yaml` at the root of `QUALITY.md`, lint reports
  no unknown-key finding for `config`.
- With invalid root `config` values such as an empty value, a map, an absolute
  path, or `../outside.yaml`, lint reports an error.
- With any other unknown root key, lint reports an error.
- With an unknown nested key under an Area, Factor, Requirement, or Rating
  Level, lint reports an error.
- With `.quality/config.yaml` containing `evaluationDir: tmp/evals`, evaluation
  creation and status use `tmp/evals`.
- With both `--evaluation-dir custom/evals` and config
  `evaluationDir: tmp/evals`, evaluation creation uses `custom/evals`.
- The `/quality` skill setup and model-change workflows describe and write the
  quality log under `.quality/log/`.
- Generated report formatter ignores and docs no longer name
  `quality/evaluations/` or `quality/log/` as the default live artifact paths.

## Durable spec changes

### To add

None.

### To modify

- `specs/cli/evaluation-create.md` - change the default evaluation directory
  and config-resolution precedence to the workspace quality data directory
  requirements above.
- `specs/cli/status.md` - change status evaluation-history resolution and
  reported default path to the workspace evaluation directory requirements
  above.
- `specs/cli/lint.md` - clarify that qualitymd lint applies the default
  qualitymd lint profile, including documented tooling-key checks in addition
  to normative Model checks.
- `specs/cli/lint-rules.md` - specify root `config` validation and the
  internally configurable unknown-key rule default.
- `specs/skills/quality-skill/reporting.md` - update evaluation artifact
  defaults and workspace config resolution for the runtime skill contract.
- `specs/skills/quality-skill/quality-log.md` - move the default quality log
  location to `.quality/log/` and retain `logDir` as deferred.
- `specs/skills/quality-skill/modes/setup.md` - align setup artifacts with the
  `.quality/` quality data directory.
- `specs/skills/quality-skill/modes/evaluate.md` - align evaluate artifacts and
  run-frame wording with the workspace evaluation directory.
- `specs/skills/quality-skill/modes/wizard.md` - align history and quality-log
  inspection with the quality data directory.
- `specs/skills/quality-skill/quality-skill.md` - align shared skill config and
  artifact terminology with the workspace and quality data directory.
- `specs/quality-schema-json.md` - confirm the schema remains normative-format
  only and does not include the qualitymd-only `config` convention.

### To rename

None.

### To delete

None.
