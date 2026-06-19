---
type: Functional Specification
title: Evaluation command surface redesign - functional spec
description: Requirements for a noun/verb-consistent qualitymd evaluation surface with plan.md-folded coverage and a separated report gate.
tags: [cli, evaluation, surface]
timestamp: 2026-06-19T00:00:00Z
---

# Evaluation command surface redesign - functional spec

Companion to the
[Evaluation command surface redesign](../0039-evaluation-command-surface.md)
change case. This spec states *what* the reshaped `qualitymd evaluation` surface
must require. It defers the on-disk record contract to
[`specs/evaluation-records.md`](../../../specs/evaluation-records.md), the
cross-cutting CLI contract to [`specs/cli.md`](../../../specs/cli.md), and
evaluation semantics to
[`SPECIFICATION.md`](../../../SPECIFICATION.md). It changes the command surface
and the planned-coverage artifact location; it does not change evaluation
semantics, rating vocabulary, or the judgment content of any record.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

The evaluation surface accreted one command at a time, and the names now fight
the pipeline they model. Verbs repeat the noun the `evaluation` parent already
carries (`create-run`, `add-record`, `show-status`, `build-report`); the three
record kinds share one `add-record` verb despite behaving differently on disk;
`build-report` overloads rendering with a CI gate; there is no way to enumerate
runs or inspect a run's records; and a removed `model`-altitude leaves dead
residue. The fix is not a bigger surface but a single rule applied uniformly, so
the surface stops accreting exceptions. Because the CLI is pre-1.0 with one
in-repo consumer (the `/quality` skill), the names can be broken outright in the
same change that updates the skill.

## Scope

Covered: the structure, names, arguments, flags, batching, run resolution,
output streams, and `--file` handling of the `qualitymd evaluation` commands; the
relocation of planned coverage into `plan.md`; and the removal of the altitude
residue.

Out of scope: report rendering content (beyond separating the gate), evaluation
semantics, rating vocabulary, record judgment schemas (field meaning is
unchanged — only the command that writes them changes), `show`/`remove` record
verbs, and any non-evaluation command.

## Requirements

### Surface structure

- The surface **MUST** follow one rule: a concept with more than one operation is
  a **noun** carrying its operations as verb subcommands; a single-operation
  action on the run itself is a **bare verb** directly under `evaluation`.

  > Rationale: the surface grew exceptions because no rule decided noun-vs-verb.
  > One litmus — "does this concept have more than one operation?" — settles every
  > case and gives new operations an obvious home. — 0039

- The run is the implicit primary resource. Its single-operation lifecycle
  actions — `create`, `list`, `status` — **MUST** be bare verbs directly under
  `evaluation`, not nested under a noun.

- The record kinds and the report — each with more than one operation — **MUST**
  be nouns: `assessment`, `analysis`, `recommendation`, and `report`, each with
  its verbs as subcommands.

- No command name **MUST** repeat the `evaluation` parent noun (no `-run`,
  `-record`, `-report`, `-status`, `-coverage` suffixes on the verbs).

### Run resolution

- Every run-scoped command **MUST** accept the run folder as a positional
  argument exactly as the CLI reported it.

- Every run-scoped command **MUST** accept `--latest` as an alternative to the
  positional argument, resolving to the most recent run in the resolved
  evaluation directory. A command **MUST** error rather than guess when neither a
  positional run nor `--latest` is given, and **MUST** error when both are given.

  > Rationale: `--latest` removes path-threading for resume and human use, but a
  > silent "newest run" default could target the wrong run; resolution stays
  > explicit. — 0039

### `evaluation create`

- `qualitymd evaluation create` **MUST** create a numbered run folder, replacing
  the former `create-run`. It **MUST** accept `--subject <path>` (defaulting to
  `QUALITY.md`), `--narrowing <slug>`, an evaluation-directory override, and
  `--json`.

- The create path **MUST NOT** expose any altitude option, flag, or receipt
  field. New runs are always subject-altitude; the removed `model`-altitude
  option, its unreachable guard, and the `altitude` receipt field **MUST** be
  gone.

  > Rationale: `model`-altitude was removed but its residue lingered as an
  > always-`"subject"` field and an unreachable guard — dead surface that misleads
  > readers into thinking altitude is still a dimension. — 0039

### `evaluation list`

- `qualitymd evaluation list` **MUST** enumerate the runs in the resolved
  evaluation directory. Under `--json` it **MUST** emit a machine-readable list
  where each entry identifies the run path and enough state to route from
  (at least: subject, narrowing when present, record counts, and whether the run
  is reportable).

  > Rationale: callers had no command to find runs and were told to keep the exact
  > paths the CLI printed; `list` makes runs discoverable and unblocks resume. —
  > 0039

- `list` **SHOULD** support filtering by run state (such as incomplete, stale, or
  reportable) so resume and diagnosis can target runs directly.

- `list` is read-only and **MUST NOT** write or modify any run.

### Record resources

- `assessment`, `analysis`, and `recommendation` **MUST** each be a noun with its
  own verbs, replacing the single `add-record` verb. Each verb **MUST** name the
  storage semantics honestly:
  - `assessment add <run>` appends a numbered assessment record;
  - `analysis set <run>` upserts the analysis for its target (replacing any
    existing one);
  - `recommendation add <run>` appends a numbered recommendation record.

  > Rationale: routing all three through `add-record` called an upsert "add" and
  > hid that the kinds differ on disk; verb-per-resource puts the
  > append-vs-replace contract in the surface. — 0039

- Each record noun **MUST** provide a `list <run>` verb that enumerates the
  written records of that kind for the run, with a `--json` form for the agent
  consumer.

- An `add`/`set` verb **MUST** accept its judgment payload from `--file <path>`,
  from `--file -`, or from non-terminal stdin, and **MUST** accept either a single
  record object or an array of records of that kind in one invocation, writing
  each in order and reporting all written paths.

  > Rationale: the skill is instructed to batch independent record writes, but one
  > process per record made that impossible; array input collapses N writes into
  > one call. — 0039

- The payload contract is otherwise unchanged from
  [`specs/evaluation-records.md`](../../../specs/evaluation-records.md): the CLI
  stamps `schemaVersion`, numbering, and filenames, and the payload **MUST NOT**
  carry them.

### `report`

- `report` **MUST** be a noun with two verbs: `report build <run>` renders
  `report-summary.md`, `report.md`, and `report.json` from the run's records;
  `report gate <run>` evaluates the rendered result against a threshold.

- `report gate <run>` **MUST** accept the threshold level, **MUST** read the
  already-rendered `report.json`, and **MUST** decide a pass/fail exit code
  without writing or modifying any run file. It **MUST** error when no rendered
  report exists for the run.

  > Rationale: `build-report --fail-at-or-below` overloaded rendering with a gate
  > and mutated the run folder even when used only as a check; a separate,
  > side-effect-free `gate` is the home a second report operation earned. — 0039

- `report build` **MUST NOT** carry a gate flag. Gating is `report gate`'s job;
  the two compose as `build` then `gate`.

### Planned coverage in `plan.md`

- Planned coverage **MUST NOT** be a standalone command or a standalone
  `planned-coverage.json` artifact. The `set-planned-coverage` command and that
  file are removed.

- `plan.md` **MUST** be a YAML-frontmatter + Markdown-body artifact. Optional
  frontmatter declares the run's intended coverage — the planned assessments
  (keyed by target plus requirement) and planned analyses (keyed by target) — and
  the Markdown body remains the run's prose plan. A run **MUST** remain valid and
  reportable when the coverage frontmatter is absent.

  > Rationale: planned coverage was a second structured file and a second write
  > path for what is one concept — the run's plan. Folding it into `plan.md`
  > frontmatter mirrors the `QUALITY.md` frontmatter+body shape and removes a
  > command and an artifact. — 0039

- The coverage frontmatter is hand-authored as part of `plan.md`. `status`
  **MUST** parse and validate it at read time, normalizing coverage keys to the
  same identities used for written records, and **MUST** surface malformed
  coverage as a gap rather than failing silently. Validation that previously ran
  at write time (rejecting duplicate planned keys, canonicalizing) **MUST** run at
  read time instead.

  > Rationale: moving coverage onto the hand-authored `plan.md` removes its
  > write-time CLI validation; validating where it is consumed preserves drift
  > detection and tolerates hand-authored keys by comparing on normalized
  > identity, not raw strings. — 0039

- `status` **MUST** continue to report the same coverage-drift gaps it does today
  — planned-but-unwritten records and written-but-unplanned records — now sourced
  from `plan.md` frontmatter instead of `planned-coverage.json`.

### Output streams and `--file`

- Across the whole evaluation surface, output **MUST** follow one rule: the
  requested data or result goes to stdout (including every `--json` payload and
  every read command's report), and side-effect confirmation lines go to stderr.
  The former split where `show-status` printed to stdout while the writers printed
  to stderr is replaced by this single rule.

- `--file <path>` and `--file -` (stdin) **MUST** be accepted by exactly the write
  verbs that take a payload (`assessment add`, `analysis set`,
  `recommendation add`) and **MUST NOT** be offered by the read commands
  (`list`, `status`, `report build`, `report gate`).

  > Rationale: `--file`/`-` meaning stdin for some evaluation subcommands and being
  > meaningless for others forced callers to memorize which; binding it to payload
  > writers makes one predictable rule. — 0039

## Durable spec changes

### To add

- `specs/cli/evaluation-create.md` — renamed from `evaluation-create-run.md`;
  drops the altitude residue (per the `evaluation create` requirements).
- `specs/cli/evaluation-list.md` — the new `list` command (per the
  `evaluation list` requirements).
- `specs/cli/evaluation-status.md` — renamed from `evaluation-show-status.md`;
  sources coverage gaps from `plan.md` (per the run-resolution, planned-coverage,
  and output-stream requirements).
- `specs/cli/evaluation-assessment.md`, `specs/cli/evaluation-analysis.md`,
  `specs/cli/evaluation-recommendation.md` — split from
  `evaluation-add-record.md`, one per record noun with honest verbs and batched
  input (per the record-resource requirements).
- `specs/cli/evaluation-report.md` — renamed from `evaluation-build-report.md`;
  splits `build` from `gate` (per the `report` requirements).

### To modify

- `specs/cli/index.md` — rewrite the evaluation command list to the new surface
  (per the surface-structure requirements).
- `specs/evaluation-records.md` — fold planned coverage into `plan.md`
  frontmatter, update the run-folder layout to drop `planned-coverage.json`, and
  update command-name references (per the planned-coverage requirements).
- `specs/cli.md` — update evaluation command-name references and confirm the
  cross-cutting stdout/stderr and `--file` contract matches the normalized surface
  (per the output-stream and `--file` requirements).

### To delete

- `specs/cli/evaluation-create-run.md`, `specs/cli/evaluation-add-record.md`,
  `specs/cli/evaluation-show-status.md`, `specs/cli/evaluation-build-report.md` —
  superseded by the renamed/split specs above.
- `specs/cli/evaluation-set-planned-coverage.md` — the command and artifact are
  removed; the contract moves into `specs/evaluation-records.md` and `plan.md`
  frontmatter.
