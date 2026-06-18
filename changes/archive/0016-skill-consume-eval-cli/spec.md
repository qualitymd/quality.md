---
type: Functional Specification
title: Skill consumes evaluation CLI - functional spec
description: Requirements that the /quality skill drive the evaluation CLI for run scaffolding, record writes, and report rendering instead of hand-authoring them.
tags: [skill, evaluation, cli]
timestamp: 2026-06-17T00:00:00Z
---

# Skill consumes evaluation CLI - functional spec

This spec states the delta for
[Skill consumes evaluation CLI](../0016-skill-consume-eval-cli.md). It constrains
the [`/quality` skill](../../../skills/quality/SKILL.md) prompt's evaluation flow so it
**consumes** the deterministic evaluation CLI rather than reproducing its work. It
specifies no CLI behavior; the commands themselves are governed by changes
[0013](../0013-evaluation-run-scaffold.md), [0014](../0014-evaluation-record-write.md),
and [0015](../0015-evaluation-report-build.md).

The enduring division of labor — the CLI owns serialization, numbering,
`schemaVersion` stamping, and report rendering; the skill owns judgment — is the
[evaluation-record contract](../0012-evaluation-record-format/spec.md). This spec
does not restate that contract; it requires the skill to defer to it through the
CLI.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be interpreted
as described in IETF RFC 2119.

## Scope

Covered: the skill's `evaluate` flow and the evaluation half of `improve` — how
they create the run folder, persist records, and produce the report.

Intentionally **deferred**: the CLI command behaviors (changes 0013/0014/0015); and
the skill's evaluation **rigor and efficiency** — effort levels, evidence rigor,
rating-binding re-check, batched writes, deep fan-out — owned by sibling change
[0017](../0017-skill-rigor-efficiency.md). This spec governs *which mechanical
surface produces the artifacts*, not how much judgment goes into them.

## Requirements

### Delegate run-folder creation

The skill **MUST** create the evaluation run folder by invoking
`qualitymd evaluation create-run`. It **MUST NOT** compute the run number or lay
out the run folder by hand.

### Delegate record writes

The skill **MUST** write every assessment, analysis, and recommendation record by
invoking `qualitymd evaluation add-record`, supplying its judgment content as input. It
**MUST NOT** hand-author record JSON or Markdown, assign record numbers, or stamp
`schemaVersion` itself.

### Delegate report rendering

The skill **MUST** inspect report renderability by invoking
`qualitymd evaluation show-status` before invoking
`qualitymd evaluation build-report`. If the status says the run is not
reportable, the skill **MUST** either add the missing judgment records through
`qualitymd evaluation add-record` or stop with the CLI-reported status; it
**MUST NOT** hand-repair the run folder.

The skill **MUST** produce `report.md` and `report.json` by invoking
`qualitymd evaluation build-report`. It **MUST NOT** hand-author either report file.

### Retain judgment

The skill **MUST** still produce all judgment content — findings, ratings,
rationales, roll-up inference, and recommendation prose — and supply it to the CLI
as the input the above commands persist. Delegating serialization does not delegate
judgment.

### Reference the record contract

The skill prompt's prose **Artifact Contract** section, which restates the record
schema and run-folder layout, **MUST** be replaced by a reference to the enduring
[evaluation-record contract](../0012-evaluation-record-format/spec.md). The skill
**MUST NOT** restate the record schema, field set, or folder layout in its own
prose, so the contract has one home.

### Fallback when the CLI lacks the commands

The skill's CLI prerequisite check **MUST** detect whether
`qualitymd evaluation create-run`, `qualitymd evaluation add-record`,
`qualitymd evaluation show-status`, and `qualitymd evaluation build-report` are
present. When any is missing, the skill **MUST**
stop and report which command is unavailable rather than fall back to
hand-authoring run folders, records, or reports — consistent with the skill's
existing prerequisite-failure behavior of stopping and helping the user install or
upgrade the CLI.
