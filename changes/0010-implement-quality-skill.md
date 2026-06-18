---
type: Change
title: Implement the /quality skill
description: Build the specified-but-unimplemented /quality evaluation skill, conforming to the format spec's Evaluation contract and driving the qualitymd CLI for every mechanical step.
status: In-Progress
tags: [skill, quality, evaluation]
timestamp: 2026-06-17T00:00:00Z
---

# Implement the /quality skill

A **Change** that builds the
[`/quality` skill](../specs/skills/quality-skill/quality-skill.md) — the judgment
companion to the deterministic [`qualitymd` CLI](../specs/cli.md). The skill is
fully specified (operating model, invocation, evaluation workflow, reporting and
artifact contract) with a
[worked example bundle](../specs/skills/quality-skill/examples/index.md), but no
implementation exists yet. This change implements it. Detail lives in the child:

- [Functional spec](0010-implement-quality-skill/spec.md) — what the change must
  do. It **defers the behavioral contract** to the durable
  [skill spec](../specs/skills/quality-skill/quality-skill.md) and states only the
  delta plus the open items and gaps that spec leaves for this change to settle.
- [Design doc](0010-implement-quality-skill/design.md) — how the change is built,
  now that the spec is settled and the change is in **In-Progress**: the skill
  packaged for Agent Skills installation, the CLI prerequisite check in
  `setup`/`wizard`, the `qualitymd models` CLI surface, the raw JSON evaluation
  artifacts, and how the open items resolve into the durable spec.

## Motivation

The CLI's deterministic format-tooling layer is built —
[`init`](../specs/cli/init.md) scaffolds, [`lint`](../specs/cli/lint.md)
validates, [`spec`](../specs/cli/spec.md) emits the format rules — but the
judgment layer those commands were designed to pair with does not exist. The
`QUALITY.md` standard's value is the *evaluation* — assessing a subject against
its model, rating the evidence, rolling it up, and advising — and the
[format spec](../SPECIFICATION.md#evaluation) assigns that work to a skill, not
the CLI. Without the `/quality` skill the repository specifies an evaluation
contract and ships a worked example of its output, but nothing produces that
output. This change closes the gap by building the skill the spec already
describes.

## Scope

Covered:

- A packaged, invocable `/quality` skill that realizes the
  [skill spec](../specs/skills/quality-skill/quality-skill.md)'s
  `evaluate`/`improve`/`setup`/`wizard` modes at the **subject** and **model**
  altitudes, is installable with `npx skills add qualitymd/quality.md`, grounds the
  format and rating vocabulary from `qualitymd spec`, verifies the `qualitymd` CLI
  prerequisite before CLI-dependent work, and drives the deterministic CLI for
  every mechanical step.
- The deterministic **`qualitymd models`** CLI command that emits the bundled
  models — including the built-in
  [quality meta-model](../internal/models/quality-meta-model.md)
  the **model** altitude evaluates against. The skill drives the CLI for every
  mechanical step and must not reimplement it, so emitting a bundled model is CLI
  work this otherwise skill-focused change also delivers; its surface is
  [open item 2](0010-implement-quality-skill/spec.md#open-items-and-gaps).
- The evaluation artifacts the
  [Reporting](../specs/skills/quality-skill/quality-skill.md#reporting) contract
  requires — model snapshot, design, plan, write-once assessment and analysis
  records, report, and recommendations — matching the shape of the
  [example bundle](../specs/skills/quality-skill/examples/index.md).
- Settling the
  [open items and gaps](0010-implement-quality-skill/spec.md#open-items-and-gaps)
  the skill spec leaves, and syncing the durable skill spec to the resolution
  before this change reaches **Done**.

Deferred (inheriting the skill spec's own
[Deferred](../specs/skills/quality-skill/quality-skill.md#deferred) surface):

- Recording per-target verdicts *through the CLI* and gating CI on them — the CLI
  record/log/gate surface is itself deferred in [`cli.md`](../specs/cli.md).
- `improve` apply-staging/isolation mechanics and bundled `references/` assets,
  both already deferred in the skill spec.

## Affected specs & docs

The durable artifacts this change creates or updates, decided up front. The skill
*implementation* itself is the change's primary product (its location is
[open item 1](0010-implement-quality-skill/spec.md#open-items-and-gaps));
the durable spec/doc deltas that accompany it are:

- [ ] [`specs/skills/quality-skill/quality-skill.md`](../specs/skills/quality-skill/quality-skill.md)
      — resolve the open items into the durable spec: model-altitude criteria,
      `setup`, default-file resolution, the `SKILL.md` trigger-description criteria
      in Frontmatter and metadata, `.quality/config.yaml` with `evaluationDir`, the
      `improve` re-evaluation folder, the raw (non-OKF) machine-readable artifact
      form (JSON assessment/analysis records plus a `report.json`), the
      altitude-first folder-naming convention, the *Limitations* vs *not assessed*
      distinction, and the *not assessed* done-criterion — so the spec matches what
      is built.
- [ ] the [worked example bundle](../specs/skills/quality-skill/examples/index.md)
      and [`specs/schema.md`](../specs/schema.md) — re-capture the example to match
      those resolutions: convert the assessment/analysis records to JSON and add
      `report.json`, **drop the OKF frontmatter so the artifacts are raw runtime
      outputs**, re-slug the folder altitude-first (`0001-subject-quality-eval`),
      and retire the now-unused `Assessment Record` / `Analysis Record` /
      `Evaluation Report` / `Recommendation` types from `specs/schema.md`.
- [ ] [`specs/cli.md`](../specs/cli.md), [`specs/cli/index.md`](../specs/cli/index.md),
      and new `specs/cli/models.md` — specify the `qualitymd models` surface that
      exposes bundled models, including the quality meta-model used by
      model-altitude evaluation.
- [ ] [`README.md`](../README.md), root `install.md`, and the
      [`docs/`](../docs/index.md) tutorials/how-tos — introduce skill-first
      onboarding with `npx skills add qualitymd/quality.md`, the `qualitymd` CLI
      prerequisite and setup verification, `/quality` invocation, modes, altitudes,
      `.quality/config.yaml` / `evaluationDir`, and how the skill pairs with the
      CLI.
- [ ] [`specs/skills/index.md`](../specs/skills/index.md) and
      [`specs/skills/quality-skill/index.md`](../specs/skills/quality-skill/index.md)
      — point at the implementation once its location (open item 1) is decided.

## Status

`In-Progress`. See the [status lifecycle](index.md#status-lifecycle). The
[functional spec](0010-implement-quality-skill/spec.md) is settled and the
[design doc](0010-implement-quality-skill/design.md) resolves the technical
approach. Implementation is now authorized: add the skill artifact, add the
`qualitymd models` CLI surface, sync the durable specs/docs and example bundle,
and verify the result before advancing the change to **Done**.
