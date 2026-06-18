---
type: Change
title: Evaluation record format
description: Make the evaluation artifact contract an enduring deterministic-surface spec the CLI writes and the skill consumes.
status: Done
tags: [evaluation, specs, cli, skill]
timestamp: 2026-06-17T00:00:00Z
---

# Evaluation record format

A unit of work that lifts the evaluation **artifact contract** out of the
`/quality` skill prompt and into the enduring [`specs/`](../../specs/index.md)
bundle. Detail lives in the children:

- [Functional spec](0012-evaluation-record-format/spec.md) — the deterministic
  evaluation-record contract.
- [Design doc](0012-evaluation-record-format/design.md) — where the contract
  lives in `specs/` and how the CLI and skill consume it.

## Motivation

The contract for a run's artifacts — folder layout, record schemas, field set,
`schemaVersion` — currently lives as prose in the
[`/quality` skill prompt](../../skills/quality/SKILL.md#artifact-contract). We are
moving to a model where the deterministic `qualitymd` CLI **writes** every
evaluation record (the skill supplies judgment; the CLI owns serialization, run
numbering, schema stamping, and report rendering). For that, the contract must
become an enduring deterministic-surface spec — one source of truth both the CLI
implementation and the skill consume — not prose duplicated inside a prompt.

## Scope

Covered: a spec-only change establishing the evaluation-record contract in
`specs/`. The spec defines the run-folder layout and naming, each record's schema
and required fields, the `schemaVersion` convention, the not-OKF statement, and
the CLI-writes / skill-judges division of responsibility.

Deferred: the CLI command surface that produces these records (separate changes
0013/0014/0015), and editing the skill prompt to reference the spec instead of
inlining it (done when this change reaches In-Progress). No code here.

## Affected specs & docs

Created or updated during `In-Progress`, before this change reaches `In-Review`
(not now — `Design` authorizes only this change's own folder plus the change
parent, bundle index, and bundle log):

- [x] `specs/evaluation-records.md` (new concept) — the enduring
      evaluation-record contract; register the concept type in
      [`specs/schema.md`](../../specs/schema.md) and list it in
      [`specs/index.md`](../../specs/index.md).
- [x] [`specs/skills/quality-skill/quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md)
      — replace the embedded runtime artifact contract with a reference to
      `specs/evaluation-records.md`.
- [x] [`specs/skills/quality-skill/examples/`](../../specs/skills/quality-skill/examples/)
      — refresh the reference evaluation artifacts so they conform to the durable
      record contract, including recommendation frontmatter.
- [x] [`skills/quality/SKILL.md`](../../skills/quality/SKILL.md) — replace the
      inlined **Artifact Contract** section with a reference to the new spec.

## Status

`Done`. Implemented and archived after lifting the evaluation-record contract into the enduring `specs/evaluation-records.md` spec the CLI writes and the skill consumes, and syncing the skill prompt and reference examples.
