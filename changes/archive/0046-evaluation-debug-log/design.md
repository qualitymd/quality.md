---
type: Design Doc
title: Evaluation debug log - design doc
description: Design for adding a process-only debug log to evaluation runs.
tags: [evaluation, records, skill]
timestamp: 2026-06-21T00:00:00Z
---

# Evaluation debug log - design doc

Design behind the [Evaluation debug log](../0046-evaluation-debug-log.md) and
its [functional spec](spec.md).

## Context

The change adds a diagnostic prose artifact to every evaluation run. The CLI
already owns deterministic run-folder creation; the skill already owns
judgment-oriented prose in `design.md` and `plan.md`. The new log sits between
those: seeded mechanically by the CLI, then hand-authored by the skill only for
notable evaluation-process events.

## Approach

`qualitymd evaluation create` seeds `debug-log.md` beside `model.md`,
`design.md`, and `plan.md`. The seeded body is intentionally small:

- a title;
- a purpose statement that names the process-only boundary;
- an `## Events` heading.

The CLI does not parse or validate the log beyond creating it as a normal run
artifact. Existing run loading and report assembly continue to ignore
`debug-log.md`; this keeps reports derived only from formal assessment,
analysis, recommendation, plan, and model data.

The `/quality` skill guidance is the main behavior change. Evaluate mode should
add concise entries when the evaluation process has diagnostic value: ambiguity,
coverage adjustment, interruption/resume, record correction, tooling failure,
redaction, prompt-injection handling, or routing a subject command into formal
assessment evidence. The skill guidance should explicitly forbid using the log
as a second evidence store.

## Alternatives

Dedicated `qualitymd evaluation debug append` command:

- Rejected for this change. It would create a cleaner mutation boundary, but the
  current need is diagnostic prose rather than a machine contract. A command can
  be added later if repeated hand-authored logs drift in shape.

Store process notes in `plan.md`:

- Rejected. `plan.md` records intended method and coverage; process events are
  chronological and may happen after records or reports are written. Mixing the
  two makes resumed-run diagnostics harder to scan.

Add debug notes to generated reports:

- Rejected. Reports are quality outputs for readers. Process diagnostics are
  useful when debugging an evaluation, but they should not change report
  authority or normal triage flow.

Do not seed the file unless needed:

- Rejected. Seeding makes the boundary visible on every run and gives agents a
  stable location without requiring a separate creation decision. Empty
  event-less logs remain lightweight.

## Trade-offs & risks

The main risk is misuse: evaluators may put subject-quality findings into the
debug log because it is easy to edit. The mitigation is explicit spec and skill
language: the log can reference formal records, but those records remain the only
source for evidence, ratings, recommendations, and report output.

The other risk is artifact sprawl. A short seeded template and no structured
parser keep the cost low. If later examples show inconsistent logs, a CLI append
command or validator can be introduced with observed data rather than guessed
schema.

## Open questions

None for this change.
