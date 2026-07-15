---
type: Design Doc
title: Intent-faithful evaluator selection — design
description: Probe-all auto discovery with a verified claude auth probe, structured authentication basis in candidate records, receipt-level candidate reporting, and skill-side transport disambiguation.
tags: [evaluation, evaluator, selection, discovery, skill]
timestamp: 2026-07-15T00:00:00Z
---

# Intent-faithful evaluator selection — design

## Context

Answers the [functional spec](spec.md) for change case
[0206](../0206-intent-faithful-evaluator-selection.md): make `auto` evaluator
discovery probe and report every built-in candidate, carry each candidate's
authentication basis as structured data, verify claude authentication with a
documented probe, and have the skill disambiguate harness-versus-SDK transport.
Discovery ordering (`codex` then `claude`) and the skill's four-tier precedence
do not change.

The spec's open question is settled: the claude runtime documents a
non-interactive probe. `claude auth status --json` (with `--json` listed as the
explicit, default output mode) prints a JSON object whose `loggedIn` field
reports authentication state, verified against the local runtime during this
design.

## Approach

### Probe-all discovery in `selectEvaluator`

The `auto` branch in `src/application/evaluation-run.ts` currently probes
codex, returns immediately when it is usable (reporting only that candidate),
and probes claude only on the fallback path. It becomes build-then-select:

1. Build the codex candidate (`which("codex")`, then `codexAuthenticated()`
   when the executable exists) and the claude candidate (`which("claude")`,
   then the new claude probe when the executable exists) — unconditionally.
2. Select the first usable candidate in the fixed order.
3. Attach the full two-entry `candidates` array to every `auto` result.
4. Compose the selection reason from the usable set: when more than one
   candidate is usable, the reason names the winner, states that deterministic
   discovery order decided the selection, and names each usable candidate not
   selected (R2).

The no-usable-candidate path still throws the `missing_evaluator` remedy list,
unchanged.

### Candidate record shape

Each candidate entry gains one structured field alongside the existing
`authenticated` boolean:

```text
authenticationBasis: "verified" | "assumed" | "unchecked"
```

- `verified` — a documented non-interactive probe ran and its result is what
  `authenticated` reports;
- `assumed` — the runtime's probe is unavailable (older CLI, unparseable
  output), so `authenticated` reflects the documented-assumption rule and the
  evidence prose says so;
- `unchecked` — no executable was found, so no probe ran.

The prose `evidence` array keeps its human-readable strings; consumers no
longer need to parse them (R3). The dry-run receipt's `schemaVersion: 3` is
kept: the change is additive (a new candidate field; candidates on the run
receipt), and version bumps here signal breaking shape changes, as when 0204
replaced the capability booleans.

### Claude authentication probe

`HostRuntimeService` (`src/services/host-runtime.ts`) gains a claude
counterpart to `codexAuthenticated`, and the `EvaluatorDiscovery` interface
passes it through. Its contract is three-valued so discovery can distinguish a
verified answer from an unavailable probe:

```text
claudeAuthenticated: () => boolean | null
```

Implementation: spawn `claude auth status --json` non-interactively, parse
stdout as JSON, and return `loggedIn === true`. Any spawn failure, non-JSON
output, or missing `loggedIn` field returns `null`, which discovery maps to
the `assumed` basis with today's usability rule (executable present ⇒ usable)
— so an older claude CLI without `auth status` degrades to current behavior
rather than turning falsely unusable. A parsed `loggedIn: false` yields
`authenticated: false`, basis `verified`, candidate not usable (R4).

Exit code is deliberately not the signal: `codex login status` uses exit
codes, but `claude auth status` is a status reporter and its logged-out exit
semantics are undocumented; the JSON field is the explicit contract.

### Receipt reporting

`evaluatorCandidates` is emitted wherever an `auto` selection is reported:
the dry-run receipt (already wired) and the returned run receipt (new). The
application decorates the runner's JSON receipt with the selection reason and
candidate array after an `auto` run. The run's `evaluation.json` is unchanged:
it keeps the selected evaluator's identity and capabilities, while candidates
and the decision reason remain receipt-level observability, so no schema-10
bump is spent on them.

### Skill-side transport disambiguation

Two wording changes in the bundled skill, with normative versions in the skill
specs per the spec's durable-spec map:

- **Provider-named intent (R5).** `skills/quality/SKILL.md` (evaluator
  selection order) and `skills/quality/workflows/evaluate.md` (the
  resolve-and-explain step) gain a rule: when the user's evaluator request
  names a provider that matches both the current harness vendor and an SDK
  evaluator kind (for example "have Claude evaluate this" inside a Claude
  harness), treat the request as ambiguous — ask a single-select closed
  choice (per the agent-mediated UX contract) offering in-session `harness`
  judgment and the independent SDK evaluator, naming the shortest path to
  each: an explicit request now, `evaluation.evaluator` config for a durable
  default. The question fires only on that vendor collision; "evaluate this"
  or an explicit `harness`/`claude`/`codex` request asks nothing.
- **Default-harness explanation (R6).** The existing explain-the-transport
  duty is extended: when `harness` was selected by precedence tier 3 (not by
  explicit request or config), the explanation states that judgment runs
  in-session with the session's own context and authentication, and names the
  independent SDK-evaluator alternative and how to choose it.

## Spec response

- **R1, R2** — build-then-select discovery with the composed reason; the
  candidates array is attached on every `auto` result and now also emitted on
  the returned run receipt.
- **R3** — the `authenticationBasis` field; `unchecked` covers the
  no-executable case the spec's two-value distinction does not reach.
- **R4** — the three-valued `claudeAuthenticated` probe; the "where documented"
  condition is discharged by the verified `claude auth status --json`
  contract, and probe unavailability degrades to the assumed basis rather than
  unusability.
- **R5, R6** — the skill wording changes above; normative text lands in
  `specs/skills/quality-skill/evaluation.md` and its evaluate workflow spec.
- **R7** — `test/application/evaluator-selection.test.ts` drops the
  `candidates` length-1 short-circuit assertions and asserts: two candidates
  reported with codex usable; the reason naming claude as usable-but-not-
  selected when both are ready; `verified` versus `assumed` bases from a
  stubbed discovery (`claudeAuthenticated` returning `true`/`false`/`null`);
  claude selection gated on the probe; and the decorated run receipt. The CLI
  integration test uses deterministic fake runtimes to cover the dry-run
  receipt, while other integration discovery stubs gain `claudeAuthenticated`.

## Alternatives

- **Reorder or capability-score `auto`.** Rejected: an explicit non-goal.
  Determinism is the compatibility promise; `evaluation.evaluator` already
  pins a preference, and R2 makes the ordering's effect visible instead.
- **Interactive CLI evaluator picker.** Rejected: the CLI is non-interactive
  by contract; the skill owns interaction, which is where R5 puts it.
- **Persist candidates into `evaluation.json`.** Rejected: costs a schema
  bump (v10 → v11) for data whose value is at decision time; the selection
  reason and candidates are carried by the decision-time receipts instead.
- **Exit-code claude probe (mirroring `codex login status`).** Rejected: the
  logged-out exit code is undocumented, and a status reporter plausibly exits
  `0` either way; the JSON `loggedIn` field is the explicit contract.
- **Treat probe unavailability as unusable.** Rejected: an older claude CLI
  would silently lose evaluation capability it has today; degrading to the
  assumed basis is honest and non-breaking.

## Trade-offs & risks

- **Discovery latency.** Probe-all means every `auto` selection spawns both
  status probes even when codex wins — the claude CLI adds roughly a second of
  startup. Acceptable against evaluation-run duration; probes only run when
  the executable exists, and parallelizing the two spawns remains open if it
  ever matters.
- **Probe output stability.** `--json` is documented as the explicit default,
  but the `loggedIn` field's stability across claude releases is not a stated
  contract. The parse-failure → `assumed` fallback bounds the damage to
  reverting one candidate to today's behavior.
- **Verified-unauthenticated tightens selection.** Today a logged-out claude
  with the executable present is selected and fails at run time; after R4 it
  is skipped at discovery (or `missing_evaluator` fires with remedies). That
  is the intended honesty, but it surfaces earlier than users may expect.

## Open questions

None. Before implementation, `claude auth status --json` was verified in an
isolated API-key-only environment: it reports `loggedIn: true` and identifies
the environment credential source without returning the credential value.
