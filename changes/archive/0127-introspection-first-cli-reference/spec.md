---
type: Functional Specification
title: Introspection-first CLI workflow conventions — functional spec
description: Requirements for renaming and refocusing the skill's former cli-quick-reference resource onto non-introspectable workflow conventions and routing command/flag/payload discovery to CLI introspection.
tags: [skill, cli, introspection, docs]
timestamp: 2026-06-26T00:00:00Z
---

# Introspection-first CLI workflow conventions — functional spec

Companion to the
[Introspection-first CLI workflow conventions](../0127-introspection-first-cli-reference.md)
change case. This spec states _what_ the refocus must achieve; no design doc is
planned. The skill's contract is defined by
[`specs/skills/quality-skill/quality-skill.md`](../../../specs/skills/quality-skill/quality-skill.md)
(normative); the CLI's introspection contract is defined by
[`specs/cli.md`](../../../specs/cli.md) (normative) and the self-teaching posture in
[Designing CLI interfaces](../../../docs/guides/cli-design.md) (informational).

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

The skill spec already requires the skill to "discover the CLI's available
commands and flags from the CLI itself rather than embedding a list that drifts"
(`quality-skill.md:589-591`), yet it also prescribes reading
`cli-quick-reference.md` (`:146-147`), ~40% of which is exactly such an embedded
list. This change brings the resource into compliance with the spec's own
introspection mandate and resolves the contradiction, while preserving the
~50–60% of the resource that is skill convention the CLI's help cannot carry. See
the change case [Motivation](../0127-introspection-first-cli-reference.md#motivation)
for the content breakdown and originating evidence. Landed 0125 added the
`model tree`/`model list`/`model get` rows to that embedded listing, so this case
also removes those rows as part of stripping duplicated command-shape content.

## Scope

Covered: renaming `cli-quick-reference.md` to `cli-workflow-conventions.md`, the
renamed resource's content boundary, the skill's routing of command/flag/payload
discovery to introspection, and the skill-spec reconciliation.

Out of scope (unchanged):

- The CLI's introspection surface — per-command `--help`, `evaluation data
kinds`/`example`/`schema`, `--json` — already exists; this case consumes it.
- The CLI's help output, command behavior, and exit codes.
- The skill's other resources and workflows, except references that point at the
  refocused resource.
- Format/rating grounding via `qualitymd spec`, which is already
  introspection-sourced.

Rename decision: the refocused resource is renamed from `cli-quick-reference.md`
to `cli-workflow-conventions.md`, because the command listing is removed and the
remaining content is workflow convention.

## Requirements

- **RF1 — Remove duplicative listings.** `cli-workflow-conventions.md` **MUST NOT**
  enumerate command shapes, flags, or payload field listings that the CLI's
  per-command `--help` and discovery commands (`evaluation data kinds`,
  `data example`, `data schema`) already provide, including the `model`
  command-shape rows landed by 0125.

  > > Rationale: an embedded copy of the CLI's own command/flag surface is the
  > > "list that drifts" the skill spec forbids (`quality-skill.md:589-591`); a
  > > breaking CLI change otherwise has to hand-edit it. — 0127

- **RF2 — Retain non-introspectable conventions.** The resource **MUST** retain
  the workflow conventions the CLI's help cannot express: the `.quality/`
  workspace-artifact layout, feedback-log sequencing, the narrowing-slug rule,
  the do/don't command rules, and cross-command orchestration sequences.

  > > Rationale: these have no CLI home; removing the resource wholesale would
  > > orphan them. The split keeps what is uniquely the skill's. — 0127

- **RF3 — Introspection-first routing.** The skill **MUST** direct its
  command, flag, and payload discovery to the CLI's introspection channels
  (per-command `--help`, `--json`, `evaluation data kinds`/`example`/`schema`)
  rather than to a hand-maintained listing.

- **RF4 — Structured channels for routed facts.** Where the skill consumes
  introspection for a fact it will route on, it **MUST** prefer the stable
  structured forms (`--json`, `data schema`, `data example`) over human-formatted
  help tables, whose wording is not a guaranteed-stable contract.

  > > Rationale: `data kinds`' human output is tab-delimited free text; the
  > > stability guarantees in `specs/cli.md` attach to `--json`/schema/example. —
  > > 0127

- **RF5 — Reconcile the spec.** `quality-skill.md` **MUST NOT** both prescribe an
  embedded command/flag listing and mandate introspection over embedding. The
  resource's prescribed role **MUST** be scoped to the conventions it uniquely
  holds (RF2), with introspection (RF3) the source of truth for command, flag,
  and payload knowledge.

## Acceptance criteria

- [x] `cli-workflow-conventions.md` contains no command/flag/payload listing that
      duplicates `--help` or the `data kinds`/`example`/`schema` discovery
      commands, including no `model tree`/`model list`/`model get` command-shape
      rows.
- [x] The resource still carries the workspace-artifact layout, feedback-log
      sequencing, narrowing-slug rule, do/don't rules, and orchestration
      sequences.
- [x] `SKILL.md` and `resources/index.md` reference
      `cli-workflow-conventions.md`, describe it by its narrowed role, and point
      command-shape questions at the CLI's introspection channels.
- [x] `quality-skill.md` no longer prescribes an embedded command listing; its
      reference to `cli-workflow-conventions.md` is scoped to conventions, and the
      introspection-over-embedding mandate reads as the source of truth for
      command/flag/payload knowledge with no contradiction.
- [x] No CLI command, help text, or exit code changed.

## Durable spec changes

Durable **specs** this change rewrites — the
[`specs/`](../../../specs/index.md) bundle and
[`SPECIFICATION.md`](../../../SPECIFICATION.md). Each subsection is required.

### To add

None.

### To modify

- `specs/skills/quality-skill/quality-skill.md` — update the prescribed read from
  `cli-quick-reference.md` (`:146-147`) to `cli-workflow-conventions.md`, scope
  that resource to the conventions it uniquely holds, and make the
  introspection-over-embedding rule (`:589-594`) the source of truth for
  command/flag/payload knowledge, so the two no longer contradict.

### To rename

None.

### To delete

None.

## Open questions

None.
