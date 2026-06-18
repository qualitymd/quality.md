---
type: Design Doc
title: Skill rigor and efficiency — design doc
description: How the /quality skill prompt operationalizes effort, requires verified evidence and pinned locators, re-checks rating-binding findings, batches artifact writes, and allows deep subagent fan-out — all as prompt structure, independent of the CLI.
tags: [skill, evaluation, rigor, design]
timestamp: 2026-06-17T00:00:00Z
---

# Skill rigor and efficiency — design doc

How the [Skill rigor and efficiency](../0017-skill-rigor-efficiency.md) change is
built — the technical approach behind its [functional spec](spec.md). The spec
says *what* the evaluation behavior must hold; this doc says *how* the prompt
makes it so, and why that way.

## Context

The only artifact this change edits is the prompt:
[`skills/quality/SKILL.md`](../../skills/quality/SKILL.md). There is no code,
schema, or CLI surface to design — the work is entirely *prompt engineering* that
turns five loosely stated process habits into observable, decidable rules an
evaluation run can be audited against. The change is **independent of the CLI
work** (0012–0016): every rule below holds against the current CLI and the
existing [Artifact Contract](../../skills/quality/SKILL.md#artifact-contract).
Where the CLI later supplies a record-writing surface, a rule here can be
superseded, but none of them depends on it.

Two facts shape every decision:

- **Correctness is prompt-bound.** Unlike the deterministic CLI, the skill's
  behavior cannot be pinned by a unit test (carried forward from
  [0010's risks](../../changes/archive/0010-implement-quality-skill/design.md)).
  The lever this change pulls is to make each rule produce a **durable artifact
  the run can be judged on** — the applied breadth lands in `plan.md`, every
  claim carries a cited command/locator in the assessment record, and a failed
  re-check blocks the headline. Rigor is enforced by what the artifacts must
  contain, not by trust in the model.
- **The artifact contract is fixed.** This change sharpens *how* the existing
  artifacts are produced; it adds no new files and changes no required layout or
  record fields. The evidence and locator requirements ride on fields the
  [assessment record](../../skills/quality/SKILL.md#artifact-contract) already
  has (`findings[].locator`, `findings[].evidence`, `criterionSource`).

## Approach

The change edits four places in the skill body — the **Arguments/Effort** note,
the **Evaluation** workflow steps, and a small new rigor subsection — plus the
**Artifact Contract** prose. Each spec requirement maps to one concrete prompt
construct below.

### 1. Operationalized effort

Today the prompt describes effort in one line of intent ("`quick` covers
high-risk hotspots only; `deep` covers the full in-scope source"). This change
replaces that with a **decision table** the skill applies at plan time, written
so an auditor can read `plan.md` and decide whether the run obeyed it:

| Effort     | Breadth (which requirements)            | Verification depth                                |
| ---------- | --------------------------------------- | ------------------------------------------------- |
| `quick`    | apex + high-risk in-scope only          | evidence MAY be sampled                           |
| `standard` | every in-scope requirement              | targeted evidence sufficient to bind each rating  |
| `deep`     | every in-scope requirement, full source | adversarially verify every rating-binding finding |

The operative design decision is **recording, not just selecting**: the skill
writes the chosen effort *and the concrete requirement set it produced* into
`plan.md` during step 6 (the existing `Write … plan.md` step). Breadth becomes
decidable from the artifact rather than asserted, which is what makes `quick`'s
"hotspots only" auditable — the omitted requirements are visible by their absence
from the recorded set. No CLI flag or schema change is needed; the run folder
already carries `plan.md`.

### 2. Evidence rigor and pinned locators

Two rules attach to the assess step (current step 7) and the record fields:

- **Verify-then-claim.** Any claim about code/CLI/tool behavior must be backed by
  an executed command or search, and the record must cite it. This rides on the
  existing finding `evidence` array: the prompt now requires each behavioral
  claim to carry an `evidence` entry whose `ref` is the command run or the search
  performed (e.g. `{ "kind": "command", "ref": "qualitymd lint QUALITY.md" }`,
  `{ "kind": "search", "ref": "rg 'apiKey' -n" }`), alongside the existing
  `kind: source` file references. The prompt states the prohibition directly:
  *do not assert code/CLI/tool behavior from memory.*
- **Pinned locators.** Every `finding.locator` must be a verifiable position — a
  `file:line` or an exact searchable string — not a recalled or approximate one.
  This is a tightening of the field's existing meaning, not a new field. "Exact
  searchable string" is the escape hatch for artifacts without stable line
  numbers (rendered output, generated files), keeping the rule satisfiable
  everywhere.

The design choice here is to express both rules as **constraints on already-required
record fields** rather than new structure. That keeps the [JSON shape a public contract](../../skills/quality/SKILL.md#artifact-contract) stable (no
`schemaVersion` bump) and means a gate reading `assessments/*.json` can mechanically
check "every finding has a pinned locator and a verifying evidence ref" without
knowing this change happened.

### 3. Rating-binding re-check

The one or two findings that determine the headline rating get an independent
second pass before they drive the report. The design:

- The skill identifies the **rating-binding findings** — the findings whose
  rating, by the model's roll-up rules, sets the top-line rating (typically the
  apex or the single worst contributing finding). These are already implicit in
  the analysis roll-up; the prompt names them explicitly as the re-check set.
- The re-check **re-runs the verifying command or search** rather than rereading
  the first observation. Re-running is the point: a stale or hallucinated first
  read is exactly what the re-check defends against, so reusing it would defeat
  the rule.
- A re-check is recorded so it is auditable. The natural home is a second
  `evidence` entry on the binding finding (the re-run command/search), making the
  double-verification visible in the same record that carries the original.
- **Gate on failure.** If a binding finding fails re-check, the report must not
  assert that headline rating. The prompt routes this back into the workflow:
  treat the failed re-check as a corrected finding, re-derive the affected rating,
  and let the corrected result drive `report.md`/`report.json`. The report is
  always the render over verified records, never over a finding that did not
  survive its second look.

This sits between the assess step and the report step in the workflow, after
analysis roll-up has identified what binds the headline and before `report.*` is
written.

### 4. Execution efficiency — batched writes

The current workflow writes "source-of-record JSON assessment and analysis
records as each is completed" (step 8) — one serial write per record. The spec
asks for compute-then-emit with batched independent writes. The design:

- **Phase the workflow into judge then emit.** The skill computes all
  per-requirement judgments first (holding the structured results in working
  context), then emits the assessment records together. Analysis roll-up and the
  reports follow, since they depend on the assessments.
- **Batch the independent writes.** Per-requirement assessment records are
  mutually independent, so they are written in parallel / as a single batch
  rather than one round trip each. Concretely, when the running agent can issue
  parallel tool calls, the skill emits the assessment files in one parallel
  block; absent that, it still batches rather than interleaving a write between
  each judgment.
- **Content-invariant.** Batching changes only *when and how many round trips*,
  never record content, field names, file names, or the required layout. The
  acceptance check is that a batched run and a serial run produce byte-identical
  artifact trees.

This is a `SHOULD` for the compute-first ordering and a `MUST` for "not one
serial write per round trip," matching the spec. No artifact contract change — it
constrains the *emission* of the same files.

### 5. Deep subagent fan-out

At `deep` effort only, and only when the in-scope work justifies it, the skill
may fan assessment out to subagents. The architecture is **orchestrator keeps
judgment; workers gather evidence**:

- **What fans out.** Per-requirement or per-target *assessment* — the
  evidence-gathering and per-requirement rating that is naturally parallel.
- **What stays central.** All roll-up judgment — aggregate, factor, and headline
  ratings — stays with the orchestrating skill. A subagent never decides the
  top-line rating; it returns structured findings, and the orchestrator
  aggregates. This keeps the single run folder, single model snapshot, single
  plan, and deterministic record names intact (the boundary
  [0010 already anticipated](../../changes/archive/0010-implement-quality-skill/design.md)
  for deep fan-out).
- **Workers return structured findings, not files.** Subagents return the same
  finding structure the assessment record expects; the orchestrator writes the
  records (composing cleanly with the batched-write phase in §4 — fan-out is
  another way to produce the judgments that the single emit phase then writes).
- **Same evidence bar.** Subagent-returned evidence must meet the §2 rigor rules
  — pinned locators and a cited verifying command/search. The orchestrator is
  responsible for the rating-binding re-check (§3) on the central side, since it
  owns the headline; it re-runs the verification itself rather than trusting a
  worker's word for the finding that binds the rating.

Fan-out is a `MAY` and an optional tactic, not a contract: a `deep` run on a
single small target need not spawn anything.

## Alternatives

- **Effort as prose intent vs. a recorded decision table (§1).** Keeping effort
  as a one-line description (status quo) was rejected: "hotspots only" is not
  auditable if the omitted requirements are never recorded. Writing the applied
  requirement set into `plan.md` makes breadth decidable from the artifact, which
  is the spec's actual requirement, at the cost of a slightly longer `plan.md`.
- **New evidence/locator schema fields vs. constraining existing ones (§2).**
  Adding a dedicated `verifiedBy` field or a `locatorKind` enum was rejected: it
  bumps the public JSON contract and forces gates to learn new fields, when the
  existing `evidence` array and `locator` field already carry the information.
  Tightening their required meaning keeps `schemaVersion` stable.
- **Re-read vs. re-run for the re-check (§3).** Re-reading the first observation
  was rejected outright — it cannot catch a stale or hallucinated first read,
  which is the whole failure mode. Re-running the verifying command/search is the
  only re-check that defends against it. The cost is a second execution for the
  one or two binding findings, which is bounded and falls on exactly the findings
  that matter most.
- **Re-check everything vs. only the binding findings (§3).** Re-checking every
  finding was rejected as disproportionate; `deep` already adversarially verifies
  every binding finding, and `standard` cannot afford a universal second pass. The
  headline is what consumers act on, so it is where the guaranteed re-check lands.
- **Streaming vs. batched writes (§4).** Serial write-as-you-go (status quo) is
  simpler to reason about but pays a round trip per record and interleaves I/O
  with judgment. Compute-then-batch was chosen because the records are
  independent and the per-round-trip cost was the concrete inefficiency the run
  exposed; the content-invariant rule keeps it from being observable in the
  output.
- **Subagents own roll-up vs. orchestrator-only roll-up (§5).** Letting a
  subagent compute a factor or headline rating was rejected: it fragments the
  single-judgment guarantee and risks divergent roll-up rules across workers.
  Confining workers to evidence + per-requirement findings, with all aggregation
  central, preserves one coherent verdict and the deterministic single-run-folder
  shape.

## Trade-offs & risks

- **Prompt length vs. enforceability.** Each rule adds prompt text. The mitigation
  is that every rule lands as an artifact constraint (recorded plan breadth, cited
  evidence, recorded re-check, byte-identical batched output) rather than pure
  exhortation, so the rules are checkable rather than merely stated.
- **Re-check cost on `standard`.** The guaranteed re-check adds a second execution
  on `standard` runs, not just `deep`. This is deliberate — the headline is the
  highest-stakes output — and bounded to one or two findings.
- **Batching is agent-capability-dependent.** Parallel tool calls are not
  available in every host. The rule degrades to "batched, not one-per-round-trip"
  so it stays satisfiable without parallelism, and the content-invariant rule
  guarantees the same artifacts either way.
- **Fan-out evidence trust.** Subagent-returned findings are only as good as their
  cited evidence. The §2 bar applies to workers, and the orchestrator independently
  re-runs verification for the rating-binding findings, so the headline never rests
  on an unverified worker claim.
- **Supersession by the CLI.** When 0012–0016 add a record-writing surface, the
  batched-write rule (§4) and parts of the evidence-recording mechanism may be
  absorbed by the CLI. That is expected; the rules are written to stand alone until
  then and to compose with, not block, that work.

## Open questions

None blocking. Deferred follow-ups:

- The exact `evidence.kind` vocabulary for verifying refs (`command`, `search`
  alongside `source`) is finalized during implementation against the current
  example bundle, so the JSON stays consistent with what the worked example
  already shows.
- If a future CLI add-record surface lands, revisit whether the re-check evidence
  is better modeled as a CLI-recorded verification step than as a second
  `evidence` entry. Out of scope for 0017.
