---
type: Functional Specification
title: qualitymd evaluation run
description: Execute a complete evaluation run with the deterministic runner.
tags: [cli, command, evaluation, runner]
timestamp: 2026-07-11T00:00:00Z
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

Evaluation needs the same orchestration, validation, and artifact guarantees
whether it is run through Codex, Claude, or an invoking agent harness.
`evaluation run` gives the workflow one deterministic owner while the selected
coding-agent evaluator discovers requirement-specific context and provides
bounded judgment. Identical evidence or ratings across runtimes are not
promised. — 0192, 0201

## Scope

This spec covers the command surface: flags, evaluator selection, dry run,
resume, cancellation reporting, output, and exit behavior. The work graph,
concurrency resolution, run-local logging, and failure taxonomy are the
[evaluation runner](../evaluation/runner.md)'s contract.

Deferred: positional natural-language scope input, until the explicit flag
surface settles.

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
- `--evaluator-result <path|->` — submit one or more harness result envelopes
  (a single object or a JSON array) for outstanding work requests, from a
  file or stdin; valid only with `--resume`.
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
agent runtime, then a ready `claude` agent runtime. It **MUST NOT** select a
configured profile implicitly or discover `harness`.

`auto` **MUST** consider an SDK-backed candidate usable only after verifying
that its executable, authentication state, structured output, fresh context,
instruction isolation, read-only workspace inspection, disabled network,
non-interactive policy, and cancellation capabilities are available; where a CLI documents no
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

The built-in names are `harness`, `codex`, and `claude`. Configured profiles may
use only `codex` or `claude`. Authentication belongs to the selected runtime;
the command does not define API-key evaluator methods or manage credentials.

### Harness checkpoints

`--evaluator harness` **MUST** select the reserved harness evaluator: bounded
judgment is supplied by the invoking agent harness through checkpoints, per
the [evaluator contract](../evaluation/evaluator-contract.md#harness-evaluator).

When a harness-backed run reaches evaluator work, the command **MUST**
persist the awaiting checkpoint, exit `0`, and emit a receipt with the stable
status `awaiting_evaluator` carrying the outstanding bounded work requests in
emission order — up to the run's resolved concurrency, per the
[runner harness contract](../evaluation/runner.md#harness-checkpoints). Each
carried request is complete: run reference, request identity, work-unit
identity and kind, subject, instructions, context, applicable body guidance,
inspection policy, expected result schema, input hash, correlation ID, and — for a retrying
request — the classified failure of the rejected attempt.

> Rationale: awaiting harness judgment is expected progress, not a failure
> that automation should retry from the beginning. Each request is complete
> and self-contained so the harness can judge it directly or hand it to a
> subagent without any other run context. — 0194, 0198

`--resume <run> --evaluator-result <path|->` **MUST** accept one harness
result envelope or a JSON array of them — any subset of the outstanding
requests, one envelope per request — from a file or stdin, and **MUST**
advance deterministic work and window top-up until the next awaiting
checkpoint or the terminal run receipt. Submitting results for a run whose
evaluator is not `harness`, or without `--resume`, **MUST** fail as a usage
error; submitting when no request is pending **MUST** fail with
`run_state_invalid`.

Resuming an awaiting run without `--evaluator-result` **MUST** re-emit the
same outstanding request set when each rebuilt input hash matches its
checkpoint entry, and **MUST** fail with `run_state_invalid` and recommend a
new run when it cannot rebuild the same requests. Result correlation,
validation, partial submission, retry, and identity binding are the
[runner](../evaluation/runner.md#harness-checkpoints) and
[orchestration](../evaluation/orchestration.md#harness-checkpoints) contracts.

### Dry run

`qualitymd evaluation run --dry-run --json` **MUST** emit a deterministic
machine-readable preview containing the resolved model, requested and planned
scope, selected evaluator with its kind and selection reason, resolved
concurrency (for a harness-backed run, the outstanding-window cap),
work-unit counts, selected evaluator capabilities, the read-only/network/
approval/verification/instruction inspection policy, each in-scope area's
effective selector and detected path, glob, or prose form, `expectedRunPath`,
and `nextActions`. It **MUST NOT** claim resolver selection, bundle size, file
count, byte caps, or a static per-area evidence package. Run receipts carry the
same effective source metadata pinned at creation.

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
