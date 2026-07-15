---
type: Design Doc
title: User-facing evaluation progress — design
description: Add a shared implementation boundary, translate evaluate protocol into user-facing phases, make evaluator alternatives informational or gated, and put selection before the first write.
tags: [evaluation, skill, agent-mediated-ux, progress]
timestamp: 2026-07-15T00:00:00Z
---

# User-facing evaluation progress — design

## Context

Answers the [functional spec](spec.md) for change case
[0207](../0207-user-facing-evaluation-progress.md). The runner must retain its
bounded request window, result envelopes, retry behavior, concurrency policy,
and resume loop; those are sound implementation mechanics. The defect is that
runtime guidance currently promotes them into ordinary user-facing progress and
combines an apparent evaluator choice with immediate continuation.

## Approach

### One implementation boundary, specialized by evaluate

The shared agent-mediated UX guide and `/quality` parent contract gain one
boundary: present the user's task state, not the implementation protocol. The
rule names five stable user-facing dimensions — state, scope or coverage,
attention, result or artifact, and next action — and permits protocol detail
only when it changes a decision or recovery action.

The evaluate workflow specializes that boundary with a small translation table
in runtime guidance:

| Runner or agent mechanic                                | Ordinary user-facing state                              |
| ------------------------------------------------------- | ------------------------------------------------------- |
| preflight commands, lint, and source resolution         | preflight complete or stopped                           |
| requirement checkpoint requests                         | evidence review in progress                             |
| work-unit and request-window totals                     | in-scope area or requirement coverage, only when useful |
| worker assignment, payload validation, and resume calls | omitted                                                 |
| terminal report projection                              | report ready / evaluation complete                      |
| classified failure with a resumable run                 | stopped, reason, and recovery action                    |

The protocol prose remains in the runtime workflow because the agent needs it
to operate correctly. A neighboring presentation rule prevents that procedure
from becoming a narrated transcript. This avoids weakening the runner contract
or duplicating it in a second, simplified workflow.

### Phase beats instead of loop beats

Ordinary progress becomes event-driven at user-meaningful transitions:

1. the run frame establishes scope, mutation, artifacts, and closeout;
2. a pre-mutation beat says preflight is complete, names the selected evaluator
   concisely, and states that evaluation artifacts plus the local feedback log
   are about to be written;
3. a long harness-backed run may say evidence review is in progress and whether
   attention is required, but does not narrate each request window; and
4. closeout reports the completed rating, evidence basis, reports, limits, and
   next action.

Quantitative context uses model coverage already available from status or
dry-run — for example, ten areas or 78 in-scope requirements. It does not turn
heterogeneous work-unit totals into a percentage or equate outstanding requests
with active workers.

### Evaluator selection is information or a gate

The provider-named ambiguity behavior from 0206 stays unchanged: when the user
said “Claude” or “Codex” and transport is ambiguous, the workflow presents a
single-select choice and waits.

When no evaluator was requested and default precedence selects `harness`, the
selection is informational:

```text
Evaluator: this session (default). It uses the context and authentication
already available here; no separate evaluator process will start. For an
independent future run, request `claude` or `codex`, or set
`evaluation.evaluator`.
```

It does not end with “unless you prefer otherwise” or another current-run
invitation. If an implementation chooses to offer a current-run change, the
existing closed-choice contract applies and the workflow waits before mutation.
This preserves discoverability of the independent path without turning every
default into a blocking picker.

### Selection before mutation

The evaluate flow changes from:

```text
lint → feedback-log write → evaluator selection → “first mutation” beat → run
```

to:

```text
lint → evaluator selection / optional dry-run → pre-mutation beat →
feedback-log write → run
```

The beat names both mutation surfaces. Feedback-log creation stays after model,
scope, and CLI preflight, so early structural stops still do not create it.

## Spec response

- R1 lands in the shared UX guide, durable parent skill contract, and runtime
  skill interaction contract; evaluate receives the concrete translation.
- R2 and R3 replace request-window narration with phase beats and optional
  model-coverage counts in the durable and runtime evaluate workflows.
- R4 narrows default-harness copy to informational future-path guidance while
  retaining a real gate for ambiguous or explicitly offered current-run
  choices.
- R5 reorders the durable flowchart, ordered procedure, and runtime steps so
  selection precedes the first write and the beat names both artifacts.

Verification is a requirement-by-requirement source ledger, targeted searches
for the superseded progress and “request it now” wording, Markdown/package
checks, and the complete repository and release gates.

## Alternatives

### Remove evaluator transport explanation entirely

Rejected. Evaluator identity affects the evidence basis, and 0206 correctly
made in-session versus independent evaluation observable. The fix is to present
that distinction without a false current-run invitation.

### Ask every default-harness run to choose a transport

Rejected. It would make the normal default a mandatory gate even when the user
expressed no evaluator preference. A gate remains correct only when ambiguity
or an offered current-run change creates a real choice.

### Change runner terminology or CLI stderr

Rejected. The runner and direct CLI need exact protocol diagnostics for
operators and recovery. The agent-mediated surface should summarize those
diagnostics rather than changing their underlying contract.

### Report completed work units as a percentage

Rejected. Work units are heterogeneous and include deterministic,
requirement-judgment, and synthesis work. A precise percentage would look more
meaningful than it is and would couple the UX to graph structure.

## Trade-offs and risks

- Less protocol detail can make a long run feel quieter. The phase guidance
  therefore permits a factual evidence-review beat and useful scope counts,
  plus a statement that no action is needed.
- “Future invocation” is less immediately discoverable than a current-run
  picker, but it avoids blocking every default run and prevents a false choice.
- The contract is guidance rather than executable UI code. Drift is controlled
  through duplicated durable/runtime review evidence and release checks, not a
  brittle prose snapshot test.

## Open questions

None.
