---
type: Change Case
title: Evaluation record write
description: Specify the CLI command that takes skill-supplied judgment and writes one schema-conformant, correctly-numbered evaluation record.
status: Done
tags: [evaluation, cli, skill]
timestamp: 2026-06-17T00:00:00Z
---

# Evaluation record write

A unit of work that specifies the deterministic command which **writes** an
evaluation record from skill-supplied judgment. Detail lives in the child:

- [Functional spec](0014-evaluation-record-write/spec.md) — the `evaluation add-record`
  command surface that persists one record.
- [Design doc](0014-evaluation-record-write/design.md) — how the subcommands map
  skill judgment to validated, atomically-numbered records.

## Motivation

Per the [evaluation record format](0012-evaluation-record-format.md), the
deterministic `qualitymd` CLI **owns writing every record** — file creation,
serialization, run-folder numbering, and `schemaVersion` stamping — while the
skill supplies only judgment. That contract has a writer-shaped hole: nothing yet
specifies the command that takes a structured judgment payload and persists a
correctly-shaped record into the right run. This change fills it.

Because the CLI is the only writer, validation is **inherent**, not a separate
step: a payload that violates the record contract is rejected at write time, so
there is no path by which a malformed record reaches disk and no `validate`
command to add.

## Scope

Covered: a spec-only change defining `qualitymd evaluation add-record` — the command that
reads a judgment payload, assigns the local record number, derives the filename,
stamps `schemaVersion`, validates against the
[record contract](0012-evaluation-record-format/spec.md), and writes one record
into the given run. Assessment, analysis, and recommendation records are
writable. Payloads can be piped on standard input or supplied with
`--file <path>`.

Deferred: run scaffolding (change 0013) and report rendering / gating (change
0015). No code here.

## Affected specs & docs

Created or updated during `In-Progress`, before this change reaches `In-Review`
(not now — `Design` authorizes only this change's own folder plus the change
parent, bundle index, and bundle log):

- [x] `specs/cli/evaluation-add-record.md` (new sub-spec) — the
      `evaluation add-record` command; list it in [`specs/cli.md`](../../specs/cli.md) and
      [`specs/cli/index.md`](../../specs/cli/index.md).
- [x] [`specs/skills/quality-skill/quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md)
      — update the evaluation workflow/reporting contract to say assessment,
      analysis, and recommendation records are written through
      `qualitymd evaluation add-record`.
- [x] [`specs/skills/quality-skill/examples/`](../../specs/skills/quality-skill/examples/)
      — refresh reference artifacts as needed for CLI-written record numbering and
      recommendation records.
- [x] [`README.md`](../../README.md) — move `evaluation add-record` from *planned*
      to built when it lands.
- [x] [`skills/quality/SKILL.md`](../../skills/quality/SKILL.md) — drive the
      evaluation loop through `evaluation add-record` rather than writing records
      directly.

## Status

`Done`. Implemented and archived after implementing `qualitymd evaluation add-record assessment|analysis|recommendation` with schema validation and atomic numbered writes.
