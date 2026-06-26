---
type: Change Case
title: Introspection-first CLI workflow conventions
description: Rename and refocus the `/quality` skill's former `cli-quick-reference.md` as `cli-workflow-conventions.md`, retaining workflow conventions the CLI's own help cannot carry, routing all command/flag/payload discovery to CLI introspection, and resolving the skill spec's embedded-reference contradiction.
status: Done
tags: [skill, cli, introspection, docs]
timestamp: 2026-06-26T00:00:00Z
---

# Introspection-first CLI workflow conventions

This parent concept captures the why and status; the detail lives in its
children:

- [Functional spec](0127-introspection-first-cli-reference/spec.md) - what the
  change must do.

No design doc is planned. The resource rename question is settled in
implementation: the refocused resource is renamed to
`cli-workflow-conventions.md`.

## Motivation

The skill spec already mandates introspection over embedding:
`specs/skills/quality-skill/quality-skill.md:589-591` — the skill "should
discover the CLI's available commands and flags from the CLI itself rather than
embedding a list that drifts." Yet `skills/quality/resources/cli-quick-reference.md`
is exactly such an embedded list — roughly 40% of its 183 lines duplicate command
shapes and flags that the CLI's per-command `--help` and discovery commands
(`evaluation data kinds` / `data example` / `data schema`) already provide — and
the *same* skill spec prescribes reading that resource (`quality-skill.md:146-147`).
The spec contradicts itself, and the embedded listing is a live drift/typo
source: the [0126](0126-bulk-data-set.md) breaking change otherwise has to
hand-edit a `data set` invocation listing that `--help` already carries. Landed
[0125](0125-model-query-commands.md) also added `model tree`/`list`/`get`
rows to the same embedded listing; those command-shape rows should now disappear
with the rest of the duplicated CLI surface rather than becoming another manual
maintenance obligation.

The other ~50–60% of the resource is content `--help` *cannot* carry — skill
conventions with no CLI home: the `.quality/` workspace-artifact layout, the
feedback-log sequencing, the narrowing-slug rule, the do/don't command rules, and
the cross-command orchestration sequences. So the fix is not "delete the resource
and point at help" — that would orphan real content — but to **split**: strip the
duplicative listing, keep and refocus the resource on what is uniquely the
skill's, and make the CLI's introspection channels the single source for command,
flag, and payload knowledge.

## Scope

Covered: renaming `cli-quick-reference.md` to `cli-workflow-conventions.md`,
trimming it to the non-introspectable workflow conventions, routing the skill's
command/flag/payload discovery to the CLI's introspection channels (per-command
`--help`, `--json`, `evaluation data kinds`/`example`/`schema`), and reconciling
the skill spec so it does not both prescribe an embedded command listing and
forbid one.

Settled premises (not open questions):

- **Refocus, do not delete.** The resource survives, scoped to the conventions
  the CLI cannot describe; only the duplicative command/flag listing leaves.
- **Structured channels are the source of truth.** Command, flag, and payload
  facts the skill routes on come from the CLI's stable structured forms
  (`--json`, `data schema`, `data example`), not from a hand-maintained copy and
  not by parsing human-formatted help tables whose wording is not guaranteed
  stable.

Deferred / non-goals:

- No new CLI introspection command — the surface (`--help`, `data
  kinds`/`example`/`schema`) already exists; this case consumes it, it does not
  extend it.
- No change to the CLI's help output or to any command's behavior.
- No change to the *other* skill resources or workflows beyond the references
  that point at this resource.

## Affected artifacts

Durable **spec** changes are itemized in the functional spec's
[Durable spec changes](0127-introspection-first-cli-reference/spec.md#durable-spec-changes)
section. The index below is the full skimmable list to reconcile before
Done.

### Code

- [x] None. Reviewed: the introspection surface this case relies on — per-command
      `--help`, `evaluation data kinds`/`example`/`schema` — already exists in
      `internal/cli/`. No code change.

### Format spec

- [x] [`SPECIFICATION.md`](../../SPECIFICATION.md) - reviewed; no change. The format
      and rating grounding the skill reads from `qualitymd spec` is unaffected.

### Durable specs

- [x] `specs/skills/quality-skill/quality-skill.md` - reconcile the contradiction:
      `:146-147` prescribes reading `cli-quick-reference.md`; `:589-594` mandates
      introspection over an embedded, drifting command list. Scope the resource's
      prescribed role to the conventions it uniquely holds, and make the
      introspection channels the source of truth for command/flag/payload
      knowledge. (Itemized in the spec's Durable spec changes.)

### Durable docs / bundled skill

- [x] `skills/quality/resources/cli-quick-reference.md` →
      `skills/quality/resources/cli-workflow-conventions.md` - rename the
      resource and strip the duplicative command/flag listing (the `--help` /
      `data kinds`/`example`/`schema` overlap, including the `model tree` /
      `model list` / `model get` rows landed by 0125); retain and refocus on the
      non-introspectable conventions (workspace-artifact layout, feedback-log
      sequencing, narrowing-slug rule, do/don't rules, cross-command
      orchestration).
- [x] `skills/quality/SKILL.md` - update the reference and trigger (`:38-39`) to
      reflect the renamed/refocused resource and the introspection-first posture
      (reach for `--help`/`--json`/discovery commands for command shapes).
- [x] `skills/quality/resources/index.md` - update the resource's description
      (`:8`) and name.

### Suggested new durable specs

- None. The introspection-first posture belongs in the existing skill spec; no
  new durable spec is earned.

## Relationship to 0126

Both cases touched the former `cli-quick-reference.md`. This case absorbs that
shared-file edit by renaming/refocusing the resource and removing the
single-payload `data set` listing wholesale, so [0126](0126-bulk-data-set.md)
does not need to edit the renamed resource.

## Status

`Done`. Renamed/refocused the skill resource to
`cli-workflow-conventions.md`, removed duplicated command/flag listings, updated
the skill root, resource index/log, and durable skill spec, and verified Markdown
formatting plus bundled-skill links. No CLI code changed. Archived after checks
passed.
