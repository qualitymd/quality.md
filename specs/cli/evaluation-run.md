---
type: Functional Specification
title: qualitymd evaluation run
description: Execute a complete evaluation run with the deterministic runner.
tags: [cli, command, evaluation, runner]
timestamp: 2026-07-09T00:00:00Z
---

# qualitymd evaluation run

`qualitymd evaluation run` executes a complete evaluation run for a resolved
`QUALITY.md` model. It inherits the cross-cutting CLI contract from
[qualitymd CLI](../cli.md) and is the command surface over the
[evaluation runner](../evaluation/runner.md), which schedules the work graph
per the [orchestration contract](../evaluation/orchestration.md), invokes
evaluators under the [evaluator contract](../evaluation/evaluator-contract.md),
and persists the [`evaluation.json`](../evaluation/evaluation-json.md) run
artifact.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background / motivation

Evaluation needs the same results whether it is run through Codex, Claude, a
direct OpenAI or Anthropic API key, or a future runtime. `evaluation run` gives
the workflow one deterministic owner: the CLI executes the work graph and the
selected evaluator provides bounded judgment, so orchestration quality no
longer depends on the invoking harness. — 0192

## Scope

This spec covers the command surface: flags, evaluator selection, dry run,
resume, cancellation reporting, output, and exit behavior. The work graph,
concurrency resolution, run-local logging, and failure taxonomy are the
[evaluation runner](../evaluation/runner.md)'s contract.

Deferred:

- positional natural-language scope input, until the explicit flag surface
  settles;
- selecting the reserved `shell` and `manual` evaluator names, which have no
  implementations yet.

## Arguments and flags

The command takes no positional arguments.

- `--model <path>` — selected `QUALITY.md` file to evaluate; defaults to
  `QUALITY.md`.
- `--evaluation-dir <path>` — override the model-relative evaluation directory.
- `--area <area-ref>` — canonical area reference for the run scope.
- `--factor <factor-ref>` — canonical factor reference for a scoped evaluation;
  repeatable.
- `--evaluator <name>` — evaluator to use: `auto`, a built-in name (including
  `harness`), or a configured profile.
- `--resume <run>` — resume an existing run from its `evaluation.json`.
- `--evaluator-result <path|->` — submit one harness result envelope for the
  awaiting work request, from a file or stdin; valid only with `--resume`.
- `-n/--dry-run` — preview the resolved run without invoking an evaluator or
  writing evaluation data.
- `--json` — emit a machine-readable receipt on stdout.

## Requirements

The command **MUST** execute a complete evaluation run for the resolved model.
Run creation, resume, work-graph execution, evaluator invocation, result
validation, persistence, report generation, and the final receipt are owned by
the [evaluation runner](../evaluation/runner.md), not by the invoking agent or
skill.

> Rationale: if a skill or agent still owns any of these phases, the workflow
> keeps two orchestrators and cannot be made deterministic across harnesses.
> — 0192

### Evaluator selection

The command **MUST** resolve the evaluator with this precedence: `--evaluator`,
then `evaluation.evaluator` in the resolved config file, then `auto`. An
omitted `--evaluator` **MUST** behave exactly as `--evaluator auto`.

> Rationale: default evaluation should work for subscription-based Codex or
> Claude users without requiring a config file on first use. — 0192

`auto` **MUST** use deterministic local discovery, in order: a ready `codex`
CLI, then a ready `claude` CLI, then configured API profiles in alphabetical
order whose API key environment variable is present.

`auto` **MUST** consider a CLI-backed candidate usable only after verifying
that its executable, authentication state, and required non-interactive
structured-output capabilities are available; where a CLI documents no
non-interactive authentication probe, readiness assumes authentication and the
evidence says so. An unusable candidate **MUST** be skipped, and dry-run JSON
**MUST** report each considered candidate's readiness evidence and the final
selection reason, never credential values.

> Rationale: command presence alone let a dry run describe an unauthenticated
> or incompatible evaluator as ready. — 0194

`auto` **MUST NOT** infer a parent agent harness from undocumented environment
variables; a harness-backed run is selected explicitly by the skill or caller
via `--evaluator harness`.

> Rationale: Claude and Codex expose different subprocess environments, and an
> internal variable is not a cross-harness compatibility contract. — 0194

If discovery finds no usable evaluator, then the command **MUST** fail
non-interactively with the `missing_evaluator` failure category and list the
available remedies.

Evaluator names, profiles, and the configuration surface are defined by the
[evaluator contract](../evaluation/evaluator-contract.md).

### Harness checkpoints

`--evaluator harness` **MUST** select the reserved harness evaluator: bounded
judgment is supplied by the invoking agent harness through checkpoints, per
the [evaluator contract](../evaluation/evaluator-contract.md#harness-evaluator).

When a harness-backed run reaches an evaluator work unit, the command **MUST**
persist the awaiting checkpoint, exit `0`, and emit a receipt with the stable
status `awaiting_evaluator` carrying the complete bounded work request — run
reference, request identity, work-unit identity and kind, subject,
instructions, context, bounded source package, expected result schema, input
hash, and correlation ID.

> Rationale: awaiting harness judgment is expected progress, not a failure
> that automation should retry from the beginning. — 0194

`--resume <run> --evaluator-result <path|->` **MUST** accept exactly one
harness result envelope from a file or stdin and **MUST** advance
deterministic work until the next evaluator checkpoint or the terminal run
receipt. Submitting a result for a run whose evaluator is not `harness`, or
without `--resume`, **MUST** fail as a usage error; submitting when no request
is pending **MUST** fail with `run_state_invalid`.

Resuming an awaiting run without `--evaluator-result` **MUST** re-emit the
same pending request when its rebuilt input hash matches the checkpoint, and
**MUST** fail with `run_state_invalid` and recommend a new run when it cannot
rebuild the same request. Result correlation, validation, retry, and identity
binding are the [runner](../evaluation/runner.md#harness-checkpoints) and
[orchestration](../evaluation/orchestration.md#harness-checkpoints) contracts.

### Dry run

`qualitymd evaluation run --dry-run --json` **MUST** emit a deterministic
machine-readable preview containing the resolved model, requested and planned
scope, selected evaluator with its kind and selection reason, resolved
concurrency, work-unit counts, `expectedRunPath`, and `nextActions`.

A dry run **MUST NOT** create the run folder and **MUST NOT** invoke an
evaluator for judgment.

### Resume

`--resume <run>` **MUST** resume the named run when its `evaluation.json` is
resume-compatible per the
[`evaluation.json` contract](../evaluation/evaluation-json.md#resume-compatibility).
If compatibility verification fails, then the command **MUST** fail with
`run_state_invalid` and report that starting a new run is the remedy.

If `--resume` is combined with `--evaluator` naming a different evaluator than
the run manifest records, then the command **MUST** refuse the run as a usage
error and report the conflict.

> Rationale: a run's judgments stay attributable to one evaluator profile.
> Re-evaluating with a different evaluator is a new run, not a resume. — 0192

### Cancellation

When the run is interrupted by SIGINT or SIGTERM, the command **MUST** report
the run as `cancelled` rather than failed and leave the run resumable, per the
[orchestration cancellation semantics](../evaluation/orchestration.md#cancellation).

### Output

Human progress diagnostics **MUST** go to stderr; stdout stays reserved for the
command payload. Under `--json`, stdout **MUST** carry the run receipt,
including the run path, status, generated report reference, rating when
available, classified failure when present, and `nextActions`.

If evaluator selection or another pre-run failure prevents the run, then under
`--json` the command **MUST** emit a failure receipt of the shape
`{"schemaVersion": …, "status": "failed", "failure": {"category": …, "detail": …}}`.

Failure categories in receipts **MUST** use the
[runner failure taxonomy](../evaluation/runner.md#failure-taxonomy).

### Exit codes

The command **MUST** exit `0` when the run completes or checkpoints at
`awaiting_evaluator`. When the run finishes `failed` or `cancelled`, the
command **MUST** emit its receipt and exit `1`. Usage errors, including the
resume evaluator conflict, **MUST** exit `2`. Internal errors **MUST** exit
`70`.
