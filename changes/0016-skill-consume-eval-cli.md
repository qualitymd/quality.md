---
type: Change
title: Skill consumes evaluation CLI
description: Constrain the /quality skill to drive the evaluation CLI for run scaffolding, record writes, and report rendering instead of hand-authoring them.
status: In-Review
tags: [skill, evaluation, cli]
timestamp: 2026-06-17T00:00:00Z
---

# Skill consumes evaluation CLI

A unit of work that makes the [`/quality` skill prompt](../skills/quality/SKILL.md)
**consume** the deterministic evaluation CLI rather than reproduce its work. Detail
lives in the child:

- [Functional spec](0016-skill-consume-eval-cli/spec.md) — how the skill must
  drive the CLI during evaluation.
- [Design doc](0016-skill-consume-eval-cli/design.md) — how the skill's
  evaluation flow maps onto the evaluation CLI command sequence.

## Motivation

With the CLI owning serialization, run-folder numbering, `schemaVersion` stamping,
and report rendering (changes
[0012](0012-evaluation-record-format.md)–[0015](0015-evaluation-report-build.md)), the
skill still hand-creates run folders, hand-authors JSON/Markdown records, and
hand-writes report files. That duplicates the contract in two places, lets the two drift, and
reproduces the numbering bug the CLI was built to eliminate. This change closes the
gap from the skill side: the skill stops doing the mechanical work and delegates it
to the CLI, keeping only judgment.

## Scope

Covered: a spec-only delta on the skill's `evaluate` flow (and the evaluation half
of `improve`) requiring it to create the run folder via `qualitymd evaluation create-run`,
write every assessment, analysis, and recommendation record via
`qualitymd evaluation add-record`, inspect run renderability via
`qualitymd evaluation show-status`, and produce `report.md`/`report.json` via
`qualitymd evaluation build-report` — and to replace the prose **Artifact
Contract** section with a reference to the enduring
[evaluation-record contract](0012-evaluation-record-format/spec.md). It also states
the fallback when the CLI lacks these commands.

Deferred — owned elsewhere, not restated here:

- The CLI command behaviors themselves (changes 0013/0014/0015).
- The skill's evaluation **rigor and efficiency** — effort levels, evidence rigor,
  rating-binding re-check, batched writes, deep fan-out — owned by sibling change
  [0017](0017-skill-rigor-efficiency.md). This change governs *which mechanical
  surface produces the artifacts*; 0017 governs *how much judgment and verification
  goes into them*. They share the skill prompt but not requirements.

No code here — `Design` authorizes only this change's own folder plus the
change parent, bundle index, and bundle log.

## Dependencies

- **Depends on** [0013](0013-evaluation-run-scaffold.md) (`evaluation create-run`),
  [0014](0014-evaluation-record-write.md) (`evaluation add-record`), and
  [0015](0015-evaluation-report-build.md) (`evaluation show-status` and
  `evaluation build-report`): the skill can only delegate to commands that exist.
- **Sibling of** [0017](0017-skill-rigor-efficiency.md): both edit the skill
  prompt's evaluation flow. This change owns the CLI-consumption delta only.

## Affected specs & docs

Created or updated during `In-Progress`, before this change reaches `In-Review`
(not now — `Design` authorizes only this change's own folder plus the change
parent, bundle index, and bundle log):

- [x] [`skills/quality/SKILL.md`](../skills/quality/SKILL.md) — drive run-folder
      creation, record writes, and report rendering through the CLI in the
      [Evaluation](../skills/quality/SKILL.md#evaluation) flow; replace the
      [Artifact Contract](../skills/quality/SKILL.md#artifact-contract) section with a
      reference to the enduring record contract.
- [x] [`specs/skills/quality-skill/quality-skill.md`](../specs/skills/quality-skill/quality-skill.md)
      — sync the durable skill workflow with the CLI-driven evaluation flow and the
      new record contract reference.
- [x] [`docs/guides/use-quality-skill.md`](../docs/guides/use-quality-skill.md)
      — update prerequisite guidance so users can verify the evaluation CLI commands
      the skill now requires.

## Status

`In-Review`. See the [status lifecycle](index.md#status-lifecycle).
