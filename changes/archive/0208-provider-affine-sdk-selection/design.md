---
type: Design Doc
title: Provider-affine SDK evaluator selection — design
description: Use the existing automatic discovery preview as a readiness oracle, apply harness-known provider affinity in the skill, pin the determined evaluator explicitly, and retain harness as the no-SDK fallback.
tags: [evaluation, evaluator, selection, skill, sdk, harness]
timestamp: 2026-07-15T00:00:00Z
---

# Provider-affine SDK evaluator selection — design

## Context

Answers the [functional spec](spec.md) for change case
[0208](../0208-provider-affine-sdk-selection.md). The standalone CLI already
probes both built-in SDK candidates and returns structured readiness evidence,
but its `auto` winner uses a fixed Codex-then-Claude order and cannot portably
know which agent harness invoked it. The `/quality` skill knows that harness
identity and can therefore add provider affinity without a new CLI flag,
environment convention, or receipt shape.

## Approach

### Resolve overrides before automatic discovery

The workflow retains its two strongest inputs:

1. an explicit evaluator request selects `harness`, `codex`, `claude`, or a
   configured profile directly; and
2. a non-`auto` `evaluation.evaluator` workspace value selects that configured
   evaluator.

A bare provider name maps to its built-in SDK evaluator. It no longer enters a
same-provider harness-versus-SDK interaction. Explicit and configured failures
remain terminal selection failures with the runner's concrete remedy.

### Use the CLI preview as the readiness oracle

With no override, the skill runs the existing read-only preview with an
explicit automatic evaluator and the already resolved model and scope:

```text
qualitymd evaluation run --dry-run --evaluator auto --json
  [--model <model>] [--area <area-ref>] [--factor <factor-ref>...]
```

The receipt already contains the CLI's selected evaluator plus both candidates'
`usable` state and authentication basis. The skill consumes that structured
data rather than reproducing executable or authentication probes.

Only a `missing_evaluator` discovery result means that no SDK candidate is
usable and permits the harness fallback. A source, model, scope, capability, or
other preview failure remains a stop; it is not reinterpreted as evidence that
the SDKs are unavailable.

### Apply provider affinity in the skill

The skill maps the known invoking harness to a built-in provider name:

```text
Codex harness  → codex candidate
Claude harness → claude candidate
other/unknown  → no affinity
```

If the matching candidate is usable, it becomes the selected evaluator even
when the receipt's fixed standalone order chose the other candidate. If the
matching candidate is not usable or there is no mapping, the receipt's selected
evaluator wins. If discovery reports `missing_evaluator` and the invoking
harness can service checkpoints, `harness` wins. Otherwise selection stops.

This creates a deterministic table:

| State                               | Selected evaluator                     |
| ----------------------------------- | -------------------------------------- |
| Explicit request                    | Requested evaluator                    |
| Non-`auto` config                   | Configured evaluator                   |
| Matching SDK candidate usable       | Matching `codex` or `claude` evaluator |
| No usable match, another SDK usable | CLI `auto` winner                      |
| No SDK usable, harness capable      | `harness`                              |
| No SDK usable, harness incapable    | Stop                                   |

### Pin and explain the result

The actual run always receives the determined evaluator explicitly. A
provider-affine Codex selection therefore invokes `--evaluator codex`, a CLI
fallback winner invokes its concrete built-in name, and the no-SDK fallback
invokes `--evaluator harness`. The recorded run is never left to re-resolve
`auto` after the preview.

The pre-mutation progress beat names:

- the evaluator and transport;
- whether explicit intent, configuration, provider affinity, CLI automatic
  discovery, or no-SDK fallback determined it; and
- that SDK selection starts fresh independent sessions, while harness selection
  uses the current session.

It does not ask whether the user prefers another transport. A later failure
does not trigger another evaluator; switching requires a new run.

### Documentation-only runtime implementation

The behavior lives in the durable and bundled `/quality` skill contracts. No
TypeScript or CLI grammar changes are needed: candidate probing, authentication
evidence, deterministic standalone ordering, dry-run behavior, and explicit
evaluator pinning already exist and are covered by current tests. Verification
therefore combines contract/runtime traceability, a stale-wording sweep, the
existing evaluator-selection tests, and the full repository and release gates.

## Spec response

- R1 uses the read-only automatic preview before considering harness and limits
  fallback to `missing_evaluator`.
- R2 applies provider affinity over the candidate array while preserving the
  receipt winner as the no-match SDK fallback.
- R3 maps bare Codex and Claude requests directly to the built-in SDK names.
- R4 gives the pre-mutation beat a selected transport and determination reason,
  with no choice invitation.
- R5 pins the concrete evaluator on the real run and keeps explicit,
  configured, and already-created runs outside the automatic fallback path.

## Alternatives

### Keep harness first and only improve the explanation

Rejected. Visibility would improve, but the ordinary evidence basis would
still be the current conversation when fresh isolated SDK sessions are ready.

### Let standalone CLI `auto` infer its parent harness

Rejected. Codex and Claude do not expose one portable documented parent-runtime
signal, and ambient inference would couple the CLI to private environment
details. The skill already has the identity needed for affinity.

### Add an invoker or preferred-provider CLI flag

Rejected. The agent-mediated workflow can resolve and pass a concrete
`--evaluator` using existing candidate data. A new flag would enlarge the CLI
surface without enabling a use case the current skill cannot serve.

### Invoke the matching evaluator without discovery

Rejected. Harness identity strongly suggests runtime availability but does not
prove that the local SDK executable and authentication path are usable. The
existing preview provides that evidence and exposes a ready alternate SDK.

### Ask the user on every unresolved run

Rejected. Evaluator identity should be observable, but deterministic precedence
and provider affinity supply a defensible default without recurring setup
friction.

## Trade-offs & risks

- Automatic selection adds one read-only runner preview before a real run. The
  extra process probes are small relative to an evaluation and prevent choosing
  from assumptions.
- A Claude harness may now prefer Claude even though standalone CLI `auto`
  remains Codex-first. The explanation must name provider affinity so the
  difference is not mistaken for inconsistent discovery.
- Explicitly pinning the preview result creates a small interval in which
  runtime readiness could change. The real run still owns failure reporting and
  must stop rather than silently select another evaluator.
- Skill behavior is specified in prose rather than executable selection code.
  Drift control relies on mirrored durable/runtime contracts, logs, targeted
  searches, and release review.

## Open questions

None.
