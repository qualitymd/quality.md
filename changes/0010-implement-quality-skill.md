---
type: Change
title: Implement the /quality skill
description: Build the specified-but-unimplemented /quality evaluation skill, conforming to the format spec's Evaluation contract and driving the qualitymd CLI for every mechanical step.
status: Draft
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

A design doc follows once the spec is settled and the change moves to **Design**.

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
  altitudes, grounds the format and rating vocabulary from `qualitymd spec`, and
  drives `init`/`lint`/`spec` for every mechanical step.
- The evaluation artifacts the
  [Reporting](../specs/skills/quality-skill/quality-skill.md#reporting) contract
  requires — model snapshot, design, plan, write-once assessment and analysis
  records, report, and recommendations — matching the shape of the
  [example bundle](../specs/skills/quality-skill/examples/index.md).
- Settling the
  [open items and gaps](0010-implement-quality-skill/spec.md#open-items-and-gaps-settle-before-design)
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
[open item Q1](0010-implement-quality-skill/spec.md#open-items-and-gaps-settle-before-design));
the durable spec/doc deltas that accompany it are:

- [ ] [`specs/skills/quality-skill/quality-skill.md`](../specs/skills/quality-skill/quality-skill.md)
      — resolve the open items into the durable spec (model-altitude criteria,
      `setup`, default-file resolution, machine-readable report output, `improve`
      re-evaluation folder) and reconcile it with the worked example (folder
      `<scope>` naming, the *Limitations* vs *not assessed* distinction, the
      done-criterion for a *not assessed* gap), so spec and implementation agree.
- [ ] [`README.md`](../README.md) and the [`docs/`](../docs/index.md)
      tutorials/how-tos — introduce `/quality` to users: how to invoke it, the modes
      and altitudes, and how it pairs with the CLI. The exact files are settled during
      Design.
- [ ] [`specs/skills/index.md`](../specs/skills/index.md) and
      [`specs/skills/quality-skill/index.md`](../specs/skills/quality-skill/index.md)
      — point at the implementation once its location (Q1) is decided.

## Status

`Draft`. See the [status lifecycle](index.md#status-lifecycle). While Draft, work
stays inside this change's folder: the
[functional spec](0010-implement-quality-skill/spec.md) is settled and the open
questions resolved before the change advances to **Design**, and no durable
spec/doc or implementation file is touched until **In-Progress**.
