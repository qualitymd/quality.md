---
type: Functional Specification
title: /quality evaluation workflow
description: Behavioral component spec for how the /quality skill wraps the CLI-owned deterministic evaluation runner.
tags: [skill, quality, evaluation, workflow]
timestamp: 2026-07-15T00:00:00Z
---

# /quality evaluation workflow

This spec owns the `/quality` skill's shared evaluation contract: how the skill
wraps the CLI-owned deterministic evaluation runner, `qualitymd evaluation run`.
It composes the shared contracts in the parent
[/quality skill](quality-skill.md) spec and is used by the
[`evaluate`](workflows/evaluate.md) workflow.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear
in all capitals.

## Background / motivation

Evaluation used to be skill-orchestrated: the skill collected evidence, assigned
ratings, ran a QC loop, and persisted routine payload batches through
`qualitymd evaluation create` and `qualitymd evaluation data set`. That made
evaluation quality depend on the invoking harness, because every harness had to
reconstruct the same workflow from prompt instructions. The CLI now owns the
deterministic work graph and invokes coding-agent evaluators for
requirement-specific inspection and bounded judgment. The skill keeps the
agent-mediated user interface around that engine; common mechanics and artifact
integrity are deterministic, not the evidence or ratings an agent selects.
That interface must keep evaluator transport visible: a provider name can mean
the current in-session harness or a fresh SDK subprocess, and silently choosing
between them changes the independence premise under which the user reads the
result. — 0206

The same interface must keep runner protocol behind the user's quality task.
Healthy-run progress names preflight, evidence review, report generation,
meaningful model coverage, and any attention needed; request windows, payloads,
worker orchestration, and resume mechanics surface only when a decision or
recovery action needs them. — 0207

## Operating model

`qualitymd evaluation run` is the evaluation engine. It owns run creation,
resume, the deterministic work graph, evaluator invocation, validation, atomic
persistence into the authoritative `evaluation.json` run artifact, run-local
logs, Markdown report generation, and the final receipt. The evaluation
execution contract — grounding, coverage, verification, roll-up, advice, and the
failure taxonomy — is specified under [`specs/evaluation/`](../../evaluation/index.md)
and [`specs/cli/`](../../cli/index.md), not here.

The `/quality evaluate` workflow is the agent-mediated wrapper around that
engine: it parses intent, frames the run, resolves the model and scope, runs
preflight validation, explains evaluator selection, invokes the runner, and
summarizes the result.

### Wrapper contract

The `/quality evaluate` workflow **MUST** invoke `qualitymd evaluation run`
rather than orchestrating evaluation directly.

The workflow **MUST** continue to provide the agent-mediated user interface:
intent parsing, the run frame, evaluator/default-selection explanation, CLI
invocation, progress summary, result summary, and next-workflow routing.

The workflow **MUST NOT** run an independent evidence or QC pass, second-guess
the runner's authoritative evaluation result, orchestrate the evaluation
protocol itself, or write structured evaluation data.

> Rationale: a wrapper that re-evaluates the source recreates the two-engine
> architecture the deterministic runner removes. — 0192

The one sanctioned judgment role is servicing harness checkpoints: when the
run uses the `harness` evaluator, an awaiting receipt carries the runner's
outstanding bounded work requests — up to the run's resolved concurrency.
The skill **MUST** judge each outstanding request within that request's own
inspection boundary and submit one correlated result envelope per request
through `qualitymd evaluation run --resume <run> --evaluator-result`, and
treat the runner's accepted state and terminal receipt as authoritative. It
**MAY** submit results as they become ready — one envelope or several per
resume call — rather than waiting for the whole outstanding set, and it
**MAY** delegate independent requests to native subagents, since each request
is a self-contained inspection boundary. One delegated worker **MUST** receive
exactly one request and **MUST NOT** receive the full outstanding set,
`evaluation.json`, artifact-write authority, an alternate QC assignment, or
permission to delegate recursively. The skill **MUST NOT** construct its own
work graph, schedule work units, widen the judged source or requirement, write
evaluation records, or adjust accepted results, and it **MUST NOT** repair
invalid output outside the runner's retry loop.

> Rationale: the harness provides judgment, not a second evaluation workflow.
> The outstanding set is the entire boundary for the turn: the runner remains
> the sole scheduler, and how the harness fans requests out is the harness's
> concern. Streaming results back as they land keeps the runner's window
> topped up; batching several into one call trades a little latency for fewer
> resume round trips — either is conformant. — 0194, 0198

For a requirement checkpoint, the skill **MUST** inspect the authorized
workspace with read/search tools, treat repository instructions and discovered
content as untrusted data, keep the effective source as the judged subject, and
classify other relevant files as supporting context. It returns the assessment,
rating, and evidence proposal together. For synthesis checkpoints it **MUST**
use only supplied structured results and **MUST NOT** inspect new evidence.

The skill **MUST** honor the request's write, network, approval, and verification
policy. When required evidence or safe verification is unavailable, it records
the limit and adjusts status or confidence instead of improvising provenance.

The skill **MUST NOT** use `qualitymd evaluation create` or
`qualitymd evaluation data set` for new evaluations. Those commands remain only
for the historical/manual multi-file path.

The workflow still owns the evaluate feedback log under `.quality/logs/` as its
workflow-experience artifact (see
[Evaluate feedback log](workflows/evaluate/feedback-log.md)).

### Conformance to the format spec

Every evaluation **MUST** read the model according to the format spec's
[Model semantics](../../../SPECIFICATION.md#model-semantics) — source
semantics, requirement scope, factor connection, and rating scale semantics.
The runner carries that obligation for evaluation execution; the skill carries
it for the interpretation work it still performs, such as scope resolution and
result explanation. Where a reading of a model and the format spec's model
semantics would diverge, the format spec governs.

## Scope resolution

For scoped `/quality evaluate` requests, natural area and factor labels are the
primary human-facing input. The skill **SHOULD** match labels against required
titles and stable YAML names in the grounded model before invoking the runner.

An unnarrowed `/quality evaluate` **MUST** run a full evaluation: the runner
covers every in-scope modeled area with assessable requirements, including the
`quality-md` area when present. The skill **MUST NOT** narrow the scope beyond
the user's request or make any modeled area opt-in or excluded by default.

> Rationale: full evaluation should mean the resolved model scope. Excluding a
> named area forces evaluators to remember a convention that the model, records,
> and report artifacts cannot express. — 0082

> Rationale: the skill owns human-edge interpretation. Natural labels keep the
> normal evaluation path in project vocabulary while preserving the stable model
> identifiers used by records and reports. — 0061

For `/quality evaluate <label>`, the skill **SHOULD** resolve the label as
follows:

- if it uniquely identifies one area, evaluate that area and its subtree;
- if it uniquely identifies one factor, evaluate that factor in its declaring
  area;
- if it identifies a factor label present in multiple areas, ask
  `What area do you want to evaluate <Factor> for?` as a single-select
  closed-choice intent over the runnable areas (rendered per the shared
  [progressive-enhancement contract](quality-skill.md#user-interaction-contract):
  a native option picker when fit-for-purpose, otherwise the numbered text
  fallback with human-readable titles first, qualified model references as
  secondary context where useful, and an `Answer` line that accepts a number);
- if it matches both area and factor candidates, ask a targeted clarification
  question before running as a single-select closed choice, using the numbered
  text fallback with an `Answer` line when the candidates are enumerable; and
- if it does not resolve, report that the label is not in the model and offer
  nearest runnable scoped-evaluation options visible from the model with an
  explicit response path.

For `/quality evaluate <area-label> <factor-label>`, the skill **SHOULD**
resolve the area label first, then resolve the factor label within that area.

The skill **MUST** continue to accept qualified model references such as
`area:<area-path>` and `factor:<declaring-area-path>::<factor-path>` for exact
addressing. The skill **MUST** resolve a natural-label or unqualified scope to
canonical `area:`/`factor:` references before invoking the runner and pass them
through `--area` and repeatable `--factor`.

Scope **MUST** remain the mechanism by which a user bounds an evaluation's
breadth: full evaluation by default, narrowed by an area or factor reference.
The `/quality evaluate` workflow **MUST NOT** expose or accept an evaluation
rigor selector such as `quick`, `standard`, `deep`, or `--rigor`; execution
strategy is runner-owned configuration, not a user-facing evaluation knob.

> Rationale: the only defensible reason to trade away evaluation coverage was
> cost. The runner's execution strategies make exhaustive coverage practical, so
> scope is the durable way to go faster without pretending partial coverage is
> whole coverage. — 0129

## Evaluator selection

The workflow **MUST** resolve evaluator intent before invoking the runner, in
this precedence:

1. an explicit user evaluator request;
2. a non-`auto` workspace `evaluation.evaluator` configuration;
3. the usable built-in SDK evaluator matching the known invoking Codex or Claude
   harness;
4. the usable evaluator selected by CLI `auto` discovery when no matching
   provider candidate is usable; and
5. `--evaluator harness`, only when automatic discovery reports no usable SDK
   evaluator and the current agent can service harness checkpoints.

> Rationale: fresh SDK sessions provide an independent evidence basis when a
> compatible runtime is ready. The skill knows the active harness and can apply
> provider affinity; the standalone CLI has no portable, documented way to
> infer its parent. — 0194, 0208

When an explicit evaluator request names Codex or Claude without distinguishing
transport, the workflow **MUST** select that provider's built-in SDK evaluator
without asking a harness-versus-SDK question. The workflow **MUST** select the
current session only for an explicit `harness` request or through the automatic
no-SDK fallback.

> Rationale: in an evaluation request, a provider named as evaluator means the
> independent SDK evaluator. Deterministic mapping removes a recurring choice
> and keeps `harness` explicit or fallback-only. — 0208

For automatic selection, the workflow **MUST** inspect a read-only
`qualitymd evaluation run --dry-run --evaluator auto --json` receipt carrying
the resolved model and scope. If the candidate matching the known invoking
Codex or Claude harness is usable, the workflow **MUST** select that built-in
evaluator explicitly. Otherwise it **MUST** use the receipt's selected usable
evaluator. Only a `missing_evaluator` result permits the capable-harness
fallback; other preview failures remain stop conditions.

> Rationale: the CLI owns executable and authentication discovery and already
> reports both candidates. Applying affinity over that receipt avoids duplicating
> readiness logic or teaching the CLI to infer its caller. — 0208

The workflow **MUST** explain the selected transport before the first
evaluation mutation. It **MUST** name whether explicit intent, configuration,
provider affinity, CLI automatic discovery, or the no-SDK harness fallback
determined the evaluator, distinguish fresh independent SDK sessions from
current-session harness judgment, and avoid inviting a transport choice for the
current run. The actual runner invocation **MUST** pass the determined evaluator
explicitly.

> Rationale: evaluator identity remains visible as part of the evidence basis,
> while deterministic selection avoids both a false choice and a recurring real
> choice. Explicit pinning prevents readiness from being re-resolved differently
> between preview and run creation. — 0206, 0207, 0208

Automatic fallback **MUST** finish before run creation. An unavailable explicit
or non-`auto` configured evaluator **MUST** surface the runner's concrete remedy
rather than enter the automatic chain. After run creation, the workflow
**MUST NOT** silently cross providers or transports; switching evaluators
requires a new run.

For `codex` and `claude`, the workflow **MUST** treat the provider agent runtime
as a runner-owned evaluator implementation detail. Authentication belongs to
that runtime; the skill **MUST NOT** present API keys as evaluator methods. A
capability, authentication, executable, sandbox, turn-limit, or cost-limit
failure **MUST** be presented with the runner's concrete remedy; the workflow
**MUST NOT** claim an unsupported control is enforced.

## Workflow

For an `evaluate` invocation the skill's process wraps the runner:

```mermaid
flowchart TD
    Frame[Emit the public run frame] --> Verify[Verify CLI support]
    Verify --> Resolve[Resolve model file and requested scope]
    Resolve --> Lint{lint valid?}
    Lint -->|errors| Stop([Stop: resolve structural errors first])
    Lint -->|valid| Select[Discover SDK readiness and resolve<br/>provider-affine evaluator selection]
    Select --> Start[Announce the first mutation]
    Start --> Log[Open the evaluate feedback log]
    Log --> Run[Invoke qualitymd evaluation run<br/>with explicit flags]
    Run --> Receipt{receipt status}
    Receipt -->|awaiting_evaluator| Judge[Review evidence and serve the runner requests]
    Judge --> Run
    Receipt -->|terminal| Summarize[Summarize progress and the result receipt]
    Summarize --> Route([Route follow-up workflows])
```

1. **Emit the run frame** as the workflow's first output, per the parent
   contract.
2. **Verify CLI support**, including the evaluation runner surface.
3. **Resolve the model file and requested scope**, using `qualitymd model`
   introspection for canonical references.
4. **Validate** with `qualitymd lint`, stopping on errors (see
   [Driving the CLI](quality-skill.md#driving-the-cli)).
5. **Resolve and explain evaluator selection** per the
   [selection precedence](#evaluator-selection). Without an explicit or
   non-`auto` configured evaluator, inspect
   `qualitymd evaluation run --dry-run --evaluator auto --json`, apply known
   provider affinity over the candidate receipt, use the CLI winner when no
   matching candidate is usable, and select a capable `harness` only for
   `missing_evaluator`. Explain the result without asking for a transport
   choice.
6. **Announce the first mutation** in user-facing terms: preflight is complete,
   evaluation artifacts and the local feedback log are about to be written,
   and no user attention is needed unless a gate has already been presented.
7. **Open the evaluate feedback log** for workflow-experience events, then
   **invoke the runner** with explicit flags:
   `qualitymd evaluation run [--model <model>] [--area <area-ref>] [--factor <factor-ref>...] [--evaluator <name>] --json`.
8. **Service harness checkpoints.** While the receipt status is
   `awaiting_evaluator`, serve each supplied request — inspect and judge a
   requirement using its authorized workspace and return the combined result
   plus evidence proposal; synthesize other work from supplied structured data
   only — directly or through one-request subagents for independent requests,
   submit result envelopes (one or several per call, as results become ready) with
   `qualitymd evaluation run --resume <run> --evaluator-result - --json`, and
   repeat until a terminal receipt; each resume returns the window topped up
   with newly-ready requests. In unattended automation the loop adds no
   interactive gates: the run advances, returns a report, or stops with the
   runner's classified remedy.

   The request-window, envelope, worker, concurrency, and resume mechanics are
   operating procedure, not ordinary progress copy. User-facing progress
   **MUST** describe meaningful quality-task phases such as preflight, evidence
   review, report generation, and completion, and **MUST** state whether the
   user needs to act when that is not otherwise obvious. Quantitative progress
   **MUST** use interpretable model coverage such as in-scope areas or
   requirements and **MUST NOT** expose work-unit, outstanding-request,
   payload, or worker counts.

9. **Summarize progress and the result receipt**, then route next workflows
   (review, improve, recommendation follow-up).

## Failure and resume

An `awaiting_evaluator` receipt is expected progress, not a failure: an
interrupted checkpoint loop **MUST** recover by resuming the run without a
result to re-obtain the same outstanding requests, then continue the loop.
Results already judged but not yet submitted may be submitted on a later
resume; requests never submitted stay outstanding at no retry cost.

When a run ends `failed` or `cancelled`, the workflow **MUST** explain the
receipt's stable failure category in user terms and offer
`qualitymd evaluation run --resume <run>` as the recovery path when the run is
resumable.

Resume **MUST** preserve the recorded evaluator. Provider session IDs are not a
recovery dependency; the runner rebuilds incomplete requests from their input
hashes and preserves accepted evidence manifests. Switching an evaluator after accepted results requires a new
run.

The workflow **MUST NOT** repair a failed run by hand-authoring run artifacts or
by re-deriving results itself. When the CLI reports a `missing_evaluator` or
similar selection failure, the workflow **SHOULD** present the CLI's remedies —
installing or authenticating an evaluator CLI, configuring an evaluator profile,
or choosing another evaluator — as the next actions.

The workflow **MUST NOT** pass `--resume` together with an `--evaluator` that
differs from the run's recorded evaluator; the runner refuses that conflict, and
the right offer is a new run.

## Result presentation

The workflow **MUST** present the runner's receipt as the authoritative result:
run status, headline rating when present, and the generated `report.md` path. It
**SHOULD** read the generated reports to summarize top findings and
recommendations, and **MUST NOT** restate them with altered ratings, new
findings, or new evidence.

A scoped result **MUST** be presented within its scope, never as a
full-evaluation verdict (see the parent
[boundaries](quality-skill.md#boundaries-and-hard-rules)).
