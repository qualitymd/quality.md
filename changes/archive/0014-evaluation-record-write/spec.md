---
type: Functional Specification
title: Evaluation record write — functional spec
description: The qualitymd evaluation add-record command that writes one schema-conformant evaluation record from a judgment payload.
tags: [evaluation, cli, skill]
timestamp: 2026-06-17T00:00:00Z
---

# Evaluation record write — functional spec

Companion to [Evaluation record write](../0014-evaluation-record-write.md). This
spec states *what* the record-writing command must do.

The on-disk contract every written record must satisfy — the run-folder layout,
record schemas, required fields, `schemaVersion`, and the CLI-writes /
skill-judges division of responsibility — is the source of truth in the
[evaluation record format spec](../0012-evaluation-record-format/spec.md). This
spec governs only the command that persists records under that contract; it does
not restate the field set.

`evaluation add-record` inherits invocation-wide behavior — non-interactive operation,
stdout/stderr separation, determinism, plain-output rules, exit-code categories,
and the `--json` convention — from the [CLI spec](../../../specs/cli.md). This
sub-spec states only what is particular to writing a record.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Scope

Covered: the command surface, the judgment-payload input channel, the local
record numbering and filename derivation, the placement and `schemaVersion`
stamping the CLI performs, the validation and rejection behavior, and the
command's outputs and exit codes — for assessment, analysis, and recommendation
records.

Deferred:

- Scaffolding a run folder and its `model.md` / `design.md` / `plan.md` (change
  0013).
- Rendering `report.md` / `report.json` and gating on the outcome (change 0015).

## Command surface

- The command **MUST** be spelled `qualitymd evaluation add-record` and **MUST**
  take the target run folder as a positional argument, so the writer always acts
  on an explicit run.
- The record kind **MUST** be selected explicitly, as a subcommand:
  `qualitymd evaluation add-record assessment <run>` writes one assessment record
  and `qualitymd evaluation add-record analysis <run>` writes one analysis
  record, and `qualitymd evaluation add-record recommendation <run>` writes one
  recommendation record. Each invocation writes exactly **one** record.

  *Note:* the kind is a subcommand rather than a flag because the two payloads
  have different required fields and land in different subdirectories; making the
  kind structural keeps each subcommand's contract distinct. Placing the writer
  under `evaluation` keeps the full run lifecycle discoverable in one command
  namespace.
- The named run **MUST** already exist as a run folder per the
  [record format spec](../0012-evaluation-record-format/spec.md#run-folder).
  `evaluation add-record` **MUST NOT** create the run; a missing or non-run
  target is a failure (see [Outputs and exit codes](#outputs-and-exit-codes)).

## Input channel

- The judgment payload **MUST** be read as a single JSON document from
  `--file <path>` when that flag is present. `--file -` **MUST** read from
  standard input.
- When `--file` is absent and standard input is not a terminal, the command
  **MUST** read the payload from standard input. When `--file` is absent and
  standard input is a terminal, the command **MUST** fail with a usage error that
  tells the caller to pass `--file <path>` or pipe JSON on standard input.
  JSON on stdin remains the natural channel for a non-interactive caller: the
  skill emits structured judgment and pipes it in, with no shell-quoting of
  nested fields.
- The payload carries **only the judgment content** — the fields the skill
  decides. It **MUST NOT** carry the fields the CLI owns: local `NNN` numbers are
  not part of the payload, and `schemaVersion` **MUST** be stamped by the CLI
  (top-level JSON for JSON records, YAML frontmatter for Markdown records), not
  trusted from input. A payload that supplies a CLI-owned field **MUST** be
  rejected rather than honored, so the division of responsibility cannot be
  bypassed through the payload.
- Reading anything other than a single well-formed JSON document — empty input,
  trailing garbage, or multiple documents — **MUST** be a usage error.

## Numbering, naming, and placement

The CLI performs every mechanical step the
[record format spec](../0012-evaluation-record-format/spec.md) assigns it:

- For an **assessment** record, `evaluation add-record` **MUST** assign the local
  `NNN` — a zero-padded sequence local to the run's `assessments/` directory —
  as one past the highest `NNN` already present, and **MUST** write the file as
  `assessments/NNN-<target>-<requirement>.json`.
- `<target>` and `<requirement>` in the filename **MUST** be slugs derived from
  the payload's `target` and `requirement`: lowercased, with each run of
  characters outside `[a-z0-9]` collapsed to a single hyphen and leading and
  trailing hyphens trimmed. The same rule produces both slug segments, so the
  derivation is one rule, deterministic, and filesystem-safe.
- For an **analysis** record, `evaluation add-record` **MUST** write the file as
  `analysis/<target>.json`, where `<target>` is the same slug of the payload's
  `target`. Analysis records are keyed by target, not numbered; writing an
  analysis record for a target that already has one **MUST** overwrite it (a
  target's roll-up is rewritten as its inputs change), and the command **MUST**
  report whether it created or replaced the file.
- For a **recommendation** record, `evaluation add-record` **MUST** assign the
  local `NNN` — a zero-padded sequence local to the run's `recommendations/`
  directory — as one past the highest `NNN` already present, and **MUST** write
  the file as `recommendations/NNN-<slug>.md`. The `<slug>` **MUST** be derived
  from the payload's `title` using the same slug rule as assessment filenames.
- `evaluation add-record` **MUST** stamp the `schemaVersion` integer (`1` for the
  current contract) on the written record: as a top-level JSON field for
  assessment and analysis records, and in YAML frontmatter for recommendation
  records.
- The record **MUST** be placed in the correct subdirectory (`assessments/` or
  `analysis/` or `recommendations/`) of the named run and nowhere else.

## Validation and rejection

`evaluation add-record` is the only writer, so it **MUST** reject any payload
that would produce a record violating the
[record contract](../0012-evaluation-record-format/spec.md) — there is no
separate validate step. At minimum it **MUST** reject:

- a payload missing any field the contract marks required for that record kind;
- an assessment payload with `notAssessed` true while `rating` is non-null (the
  contract requires `rating` be `null` when `notAssessed`);
- a `rating` that is not a level defined by the run's rating scale, or otherwise
  not a recognized value;
- a finding lacking any of `locator`, `observation`, or `category`;
- a recommendation payload missing `title`, `gap`, `evidenceLocators`,
  `remediationOptions`, `recommendedOption`, or `doneCriterion`, which are the
  structured fields the CLI uses to render the human-readable Markdown body
  required by the record contract;
- a payload supplying a CLI-owned field (see [Input channel](#input-channel)).

Rejection **MUST** be atomic: when a payload is rejected, **nothing is written**
— no partial record, no consumed `NNN`, no placeholder file. The numbering
sequence advances only when a record is actually written.

Successful writes **MUST** also be atomic from the caller's perspective: write a
complete temporary file in the destination directory, then move it into place
with an atomic rename. For numbered records, if the final rename detects an
`NNN` collision caused by a concurrent writer, the command **MUST** recompute the
next `NNN` once and retry; a second collision **MUST** fail with the internal
error category and a diagnostic naming the directory whose numbering is
contended.

## Outputs and exit codes

- On success, `evaluation add-record` **MUST** report the path of the written
  record. The human confirmation is written to standard error, keeping standard
  output clean.
- Under `--json`, `evaluation add-record` **MUST** emit a result receipt on
  stdout instead of the human confirmation, carrying at least `schemaVersion`,
  the written `path`, the record `kind`, and — for analysis — whether the file
  was `created` or replaced. Per the CLI spec it **MUST NOT** emit a success
  receipt when it refuses to write.
- A payload that violates the contract, or input that is not a single well-formed
  JSON document, **MUST** exit with the CLI spec's **usage error** category and a
  clear diagnostic naming the offending field or condition.
- A missing or non-run target, or an I/O failure that prevents writing, **MUST**
  exit with the CLI spec's **internal error** category.
- On successfully writing one record, `evaluation add-record` **MUST** exit `0`.
