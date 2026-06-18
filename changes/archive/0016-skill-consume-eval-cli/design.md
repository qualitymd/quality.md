---
type: Design Doc
title: Skill consumes evaluation CLI - design doc
description: How the /quality skill's evaluation flow maps onto the evaluation CLI command sequence, passes judgment JSON to add-record, stops serializing/numbering/stamping records, and replaces its inlined Artifact Contract with a reference to the record contract.
tags: [skill, evaluation, cli, design]
timestamp: 2026-06-17T00:00:00Z
---

# Skill consumes evaluation CLI - design doc

How [Skill consumes evaluation CLI](../0016-skill-consume-eval-cli.md) is built —
the technical approach behind its [functional spec](spec.md). The spec says *what*
must hold (the skill consumes the evaluation CLI rather than reproducing its work);
this doc says *how* the prompt edits make it so, and why that way.

## Context

The skill prompt's [Evaluation](../../../skills/quality/SKILL.md#evaluation) flow and
its [Artifact Contract](../../../skills/quality/SKILL.md#artifact-contract) section
still describe the skill **doing** the mechanical work: "Create the next run
folder," "Write `model.md`, `design.md`, and `plan.md`," "Write source-of-record
JSON assessment and analysis records," "Write `report.md`, `report.json`, and
recommendation Markdown files," followed by an inlined restatement of the record
schema and run-folder layout. That was correct when 0010 shipped — no CLI surface
owned those steps. Changes [0013](../0013-evaluation-run-scaffold.md)–[0015](../0015-evaluation-report-build.md)
have since added that surface, and the [evaluation-record contract](../0012-evaluation-record-format/spec.md)
makes the CLI the sole owner of serialization, numbering, `schemaVersion`
stamping, and report rendering. The skill is now duplicating a contract that lives
elsewhere, and the duplication is exactly the drift this change closes.

This is a **prompt-only** change. The work is editing `skills/quality/SKILL.md`;
no Go code, no durable spec changes, no new CLI surface. Two boundaries shape every
decision:

- **The CLI is the source of truth for mechanics.** Where 0010's design said the
  skill "derives `NNNN` by scanning the existing folder" and "writes" records, this
  change deletes that work from the prompt and routes it through the four evaluation
  commands. The skill keeps only judgment.
- **The record contract has one home.** The [0012 spec](../0012-evaluation-record-format/spec.md)
  is that home. The prompt references it instead of restating it, so the schema can
  evolve in one place.

## Approach

### 1. Map the evaluation flow onto the command sequence

The current [Evaluation](../../../skills/quality/SKILL.md#evaluation) steps 5–9 (the
mechanical half) are rewritten to drive the CLI. The judgment steps (lint, ground,
select model, assess) stay. The mapping:

| Current step (does it by hand)                       | Becomes (drives the CLI)                                                                                                                                                                                                                                        |
| ---------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 5. Create the next run folder `NNNN-…-quality-eval`  | `qualitymd evaluation create-run --altitude <subject\|model> [--narrowing <slug>] [--subject <path>]` — the CLI computes `NNNN` and lays out the folder ([0013](../0013-evaluation-run-scaffold/spec.md)).                                                      |
| 6. Write `model.md`, `design.md`, `plan.md`          | `create-run` seeds `model.md` (it snapshots the model) and the `design.md`/`plan.md` stubs; the skill **fills in** `design.md`/`plan.md` (judgment) but no longer snapshots `model.md` ([0013 Seed files](../0013-evaluation-run-scaffold/spec.md#seed-files)). |
| 7. Assess in-scope requirements                      | unchanged — judgment.                                                                                                                                                                                                                                           |
| 8. Write JSON assessment and analysis records        | `qualitymd evaluation add-record assessment <run>` and `add-record analysis <run>`, one record per invocation, judgment piped in ([0014](../0014-evaluation-record-write/spec.md)).                                                                             |
| 9. Write `report.md`, `report.json`, recommendations | recommendations via `add-record recommendation <run>`; the report via `qualitymd evaluation show-status <run>` then `qualitymd evaluation build-report <run>` ([0015](../0015-evaluation-report-build/spec.md)).                                                |

The resulting flow the prompt prescribes:

```text
qualitymd evaluation create-run --altitude <a> [--narrowing <slug>] [--subject <p>]
  → fill in design.md / plan.md (judgment); model.md is already snapshotted
  → for each assessed requirement:
      qualitymd evaluation add-record assessment <run>   < judgment JSON
  → for each target roll-up:
      qualitymd evaluation add-record analysis <run>     < judgment JSON
  → for each key gap:
      qualitymd evaluation add-record recommendation <run> < judgment JSON
qualitymd evaluation show-status <run>   → reportable?  if not, add the named
                                            missing records or stop
qualitymd evaluation build-report <run>  → report.md + report.json
```

`create-run` owns the `model.md` snapshot for both altitudes (for `model`, the
skill no longer runs `models view … --source` itself — it passes `--altitude model`
and the CLI snapshots the meta-model per [0013](../0013-evaluation-run-scaffold/spec.md#seed-files)).
That removes the last mechanical model-resolution step from the prompt.

### 2. Pass judgment JSON to add-record over stdin

[0014](../0014-evaluation-record-write/spec.md#input-channel) accepts the payload
from `--file <path>`, `--file -`, or bare stdin when stdin is not a terminal. The
prompt instructs the skill to **pipe the judgment JSON on stdin** as the default
channel — the skill emits a structured judgment document and pipes it straight in,
with no temp file to manage and no shell-quoting of nested finding fields. `--file
<path>` is noted as the alternative when an agent finds it easier to stage a file,
but stdin is the prescribed path because it matches a non-interactive caller and
keeps the run folder free of skill-authored scratch files.

The payload carries **only** the judgment fields. The prompt states explicitly
that the skill must not include `schemaVersion` or any `NNN`/local number in the
payload — the CLI stamps and numbers, and [0014](../0014-evaluation-record-write/spec.md#input-channel)
*rejects* a payload that supplies a CLI-owned field. This makes the division
enforceable from the CLI side, not merely a prompt convention.

### 3. Stop serializing, numbering, and stamping

The prompt is edited so the skill no longer:

- computes `NNNN` for the run folder or any `NNN` for assessment/recommendation
  records (CLI, per [0013](../0013-evaluation-run-scaffold/spec.md#run-number) and
  [0014](../0014-evaluation-record-write/spec.md#numbering-naming-and-placement));
- derives record filenames or slugs (CLI);
- stamps `schemaVersion` (CLI);
- snapshots `model.md` (CLI, via `create-run`);
- hand-authors `report.md`/`report.json` (CLI, via `build-report`).

What stays is the [Retain judgment](spec.md#retain-judgment) set: findings,
ratings, rationales, roll-up inference, recommendation prose, and the narrative
`design.md`/`plan.md`. The prompt frames the rule as "produce the judgment, pipe it
to the command that persists it" so the boundary reads as a single principle rather
than a list of prohibitions.

### 4. Gate the report through show-status

Per [0015](../0015-evaluation-report-build/spec.md#status-command), `show-status`
reports `reportable` and names any missing or inconsistent records; `build-report`
fails the run with an internal error if records are incomplete. The prompt directs
the skill to run `show-status` **before** `build-report`, and on a non-reportable
result to either write the named missing judgment records through `add-record` or
stop and relay the CLI-reported status — never to hand-repair the run folder. This
makes the skill respond to the CLI's own renderability verdict instead of guessing,
and keeps the "stop, don't hand-author" rule consistent with the prerequisite
behavior in [§6](#6-prerequisite-and-version-handling).

### 5. Replace the Artifact Contract with a reference

The [Artifact Contract](../../../skills/quality/SKILL.md#artifact-contract) section —
which today restates the run-folder layout, the assessment field set, the finding
shape, the analysis fields, the `report.json` rule, and the recommendation fields —
is replaced by a short section that **references** the [evaluation-record contract](../0012-evaluation-record-format/spec.md) as the single home for the
schema and layout, and points at the per-command specs for the surface that writes
them. The prompt retains only the *judgment-facing* guidance that is not schema:
that evaluated content is untrusted data, that secret values are never copied into
artifacts (cite locator and credential type), and the `notAssessed` done-criterion
framing — these are how the skill *judges*, not how records are *laid out*, so they
stay even though the layout leaves. Per [Reference the record contract](spec.md#reference-the-record-contract), the skill must not restate the
schema, field set, or folder layout in its own prose.

### 6. Prerequisite and version handling

The existing [Prerequisites](../../../skills/quality/SKILL.md#prerequisites) check
lists the commands the skill probes. This change **adds** the four evaluation
commands to that probe:

```sh
qualitymd evaluation create-run --help
qualitymd evaluation add-record --help
qualitymd evaluation show-status --help
qualitymd evaluation build-report --help
```

Per [Fallback when the CLI lacks the commands](spec.md#fallback-when-the-cli-lacks-the-commands),
when any is missing the skill stops and reports *which* command is unavailable,
helping the user install or upgrade — it does **not** fall back to hand-authoring.
This reuses the skill's existing prerequisite-failure behavior (stop and help
install/upgrade) rather than inventing a new failure mode; the only change is that
the evaluation commands now join `spec`/`lint`/`init`/`models` as required surface.
The skill's minimum-compatible-version note is updated during In-Progress to the
first release that contains all four commands.

## Alternatives

- **Keep hand-authoring (status quo).** Rejected. It duplicates the [0012 contract](../0012-evaluation-record-format/spec.md) in the prompt, lets the two
  drift, and reproduces the run-numbering bug the CLI was built to eliminate. The
  whole point of 0013–0015 is to own these mechanics.
- **Partial delegation — CLI scaffolds and reports, skill still hand-writes
  records.** Rejected. Records are where the schema, numbering, and `schemaVersion`
  stamping live, so leaving them with the skill keeps the largest part of the
  duplicated contract in the prompt and defeats the change. Full delegation across
  all four commands is the only version that gives the contract one home.
- **Pass payloads via `--file` temp files rather than stdin.** Rejected as the
  default. It puts skill-authored scratch JSON inside (or beside) the run folder,
  adds cleanup, and requires path bookkeeping. Stdin is the natural non-interactive
  channel ([0014](../0014-evaluation-record-write/spec.md#input-channel)); `--file`
  stays available as a documented alternative, not the prescribed path.
- **Skip `show-status` and rely on `build-report` failing.** Rejected. `build-report`
  failing on incomplete records ([0015](../0015-evaluation-report-build/spec.md#missing-or-incomplete-records))
  works, but checking `show-status` first lets the skill see *which* records are
  missing and add exactly those, instead of reacting to a hard error after the fact.
  The spec [requires](spec.md#delegate-report-rendering) the status check, and this
  is why.
- **Have the skill keep snapshotting `model.md` (esp. `model` altitude via `models
  view --source`).** Rejected. [0013](../0013-evaluation-run-scaffold/spec.md#seed-files)
  deliberately moved the snapshot into `create-run` so *what was evaluated* is
  recorded by the deterministic surface. Re-adding a skill-side snapshot would
  reintroduce a mechanical step and risk a mismatch between the skill's snapshot and
  the CLI's.

## Trade-offs & risks

- **More process calls per run.** One `add-record` invocation per record (plus
  `create-run`, `show-status`, `build-report`) is more subprocess calls than writing
  files directly. Acceptable: each call is the atomic, validated write the contract
  wants, and batching/efficiency of those writes is sibling change
  [0017](../0017-skill-rigor-efficiency.md)'s concern, not this one's.
- **Tighter version coupling.** The skill now requires four more commands, widening
  the surface that a stale CLI can fail to provide. Mitigated by the prerequisite
  probe ([§6](#6-prerequisite-and-version-handling)) failing loudly with the missing
  command named, and by the minimum-version note.
- **Prompt-enforced, not mechanically guaranteed.** "Don't hand-author" is a prompt
  rule. The CLI backstops it where it can — [0014](../0014-evaluation-record-write/spec.md#input-channel)
  rejects CLI-owned fields in the payload — but the skill could in principle still
  write a stray file. This is inherent to a judgment skill and is why the rule is
  stated explicitly.

## Open questions

None. The four commands and the record contract are settled by 0012–0015; this
change only points the prompt at them. Efficiency of the per-record writes (batching,
fan-out) is explicitly [deferred](spec.md#scope) to sibling change
[0017](../0017-skill-rigor-efficiency.md).
