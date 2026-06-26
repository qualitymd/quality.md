---
type: Functional Specification
title: Evaluation orchestration overhaul — functional spec
description: Requirements for a single best-quality evaluate workflow — rigor removed, exhaustive coverage, parallel-by-default collection, and an always-on two-pronged QC phase.
status: Done
tags: [skill, evaluation, orchestration, qc]
timestamp: 2026-06-26T00:00:00Z
---

# Evaluation orchestration overhaul — functional spec

> The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
> described in BCP 14 when, and only when, they appear in all capitals.

Companion note: this is the delta contract for the
[0129 change case](../0129-evaluation-orchestration-overhaul.md). It governs the
behavior of the `/quality` **evaluate** workflow and its shared evaluation
semantics. The binding sources of truth it defers to are the durable component
spec [`specs/skills/quality-skill/evaluation.md`](../../../specs/skills/quality-skill/evaluation.md)
(normative) and the format spec
[`SPECIFICATION.md`](../../../SPECIFICATION.md) (normative, for *Assess and Rate* /
*Analyze* semantics). The durable contracts these requirements land in are mapped
in [Durable spec changes](#durable-spec-changes).

## Background / Motivation

Evaluation rigor (`quick`/`standard`/`deep`) traded coverage for cost. Parallel
subagent fan-out removes the cost reason to ever under-cover, so the dial becomes
a footgun that produces shallow passes reading as whole coverage. The durable
knob is **scope**, already expressed by model references and `--narrowing`.

Verification was conditional and one-sided: it re-checked findings the skill
*had* (guarding false positives) but never hunted for findings it *missed*
(false negatives). The two failure modes need different machinery, so this change
promotes verification into a first-class QC phase with both prongs. See the
[change case motivation](../0129-evaluation-orchestration-overhaul.md#motivation)
for the full rationale; the durable *why* is promoted into `evaluation.md` on
landing.

## Scope

Covered: the evaluate workflow's coverage contract, execution strategy
(parallel/serial), and the QC phase; the run frame, feedback-log, and invocation
surface where they reference rigor.

Non-goals:

- The CLI never gains a rigor concept; this change touches none of it.
- **Modeling rigor** (setup discovery) and per-requirement **assessment rigor**
  (model-authoring guidance) are different concepts and are out of scope.
- Recommendation/advise behavior is unchanged.

## Assumptions & dependencies

- The harness *may or may not* expose a subagent capability (e.g. a Task/Agent
  tool). Requirements that mention fan-out are conditioned on that capability;
  the coverage and QC obligations do not depend on it.
- Subagents, where used, return text/structured data to the orchestrator and
  cannot themselves persist the run's authoritative ratings.
- `qualitymd evaluation` records, schemas, and report build are unchanged and
  remain the only writers of authoritative evaluation data.

## Requirements

### Single workflow, no rigor dial

- **R1.** The evaluate workflow **MUST NOT** expose or accept an evaluation rigor
  selector. The `quick`/`standard`/`deep` levels, the `--rigor` argument, and the
  `/quality evaluate deep` invocation form **MUST** be removed.

  > Rationale: a coverage dial whose only justification was cost is obviated by
  > parallel fan-out, and a `quick` pass silently under-covers. Scope, not rigor,
  > is the durable knob. — 0129

- **R2.** The evaluate run frame **MUST NOT** render a `Rigor:` field, and the
  evaluate feedback log **MUST NOT** carry a `rigor:` frontmatter field.

- **R3.** Scope **MUST** remain the only mechanism by which a user bounds an
  evaluation's breadth: full evaluation by default, narrowed by an Area or Factor
  reference resolved to `--narrowing`. Removing rigor **MUST NOT** remove or
  weaken any scope-resolution behavior.

### Exhaustive coverage contract

- **R4.** Every evaluate run **MUST** assess every in-scope Requirement against a
  full read of the in-scope Area `source`. There is no shallow mode.

  > Rationale: with fan-out, exhaustive is no longer slow, so partial coverage
  > has no remaining justification. — 0129

- **R5.** Each in-scope Requirement **MUST** end the run in one of two terminal
  evidentiary states: rated against the rating scale on cited verified evidence,
  or recorded as *not assessed* with a stated reason. The run **MUST NOT** leave
  an in-scope Requirement silently unexamined, and **MUST NOT** assign a level to
  fill an evidence gap (per `SPECIFICATION.md` *Assess and Rate*).

- **R6.** The report **MUST** state what was not assessed, so no run reads as
  whole coverage when it is not. (Unchanged obligation; restated because the
  rigor table that previously carried it is removed.)

### Execution strategy: parallel by default, serial fallback

- **R7.** Where the harness exposes a subagent capability, the workflow **SHOULD**
  fan out independent collection and QC work (per-Area, or per-Requirement within
  an Area) to subagents running concurrently, to hold exhaustive coverage at low
  wall-clock cost.

- **R8.** Where no subagent capability is present, the workflow **MUST** perform
  the same exhaustive coverage and the same two-pronged QC phase serially. The
  mandatory contract is coverage + QC; parallelism is only how it goes fast.

  > Rationale: the skill runs in many harnesses; making fan-out a hard
  > requirement would break the ones without subagents. The quality contract must
  > not depend on the execution strategy. — 0129

- **R9.** A subagent **MUST** be given the resolved scope, the relevant
  Requirements, the secret-handling rule, the evaluated-source-as-data rule, and
  an instruction to return structured findings only — never files, ratings, or
  final judgment. Subagent-collected evidence **MUST** meet the same locator and
  verification rules as orchestrator-collected evidence (every claim about code,
  CLI, or tool behavior verified by an executed command or search; every locator
  a `file:line` or exact searchable string).

### Orchestrator-owned judgment (preserved)

- **R10.** Roll-up judgment and all authoritative ratings (Requirement, Factor,
  Area, and headline) **MUST** be produced by the orchestrating skill, after the
  QC phase converges. Subagents **MUST NOT** produce the run's authoritative
  ratings.

### QC phase — always on, two-pronged

- **R11.** Every evaluate run **MUST** run a QC phase after initial collection and
  before final roll-up. The QC phase has two prongs — **verify** (R12) and
  **completeness sweep** (R13) — and both **MUST** run on every run, regardless of
  scope size.

- **R12. Verify prong (false-positive guard).** The QC phase **MUST**
  re-establish, by re-running the verifying command or search (not by re-reading
  the earlier observation), every finding that binds any roll-up rating — not only
  the headline — and every finding recorded at low confidence. If a binding
  finding fails re-check, the affected rating **MUST** be re-derived before it is
  reported, and the stale rating **MUST NOT** be asserted.

  > Rationale: re-reading cannot catch a stale or hallucinated first read; only a
  > re-run can. Generalizing from headline-only to all roll-up-binding findings
  > follows from removing the rigor levels that scoped it narrowly. — 0129

- **R13. Completeness sweep (false-negative guard).** The QC phase **MUST**
  perform a coverage check that:
  - confirms every in-scope Requirement reached a terminal evidentiary state
    (R5), failing the sweep if any was silently skipped or marked *not assessed*
    without a stated reason;
  - re-examines, with an adversarial gap/risk lens, every Area or Requirement
    whose first pass produced only `strength` findings or no findings;
  - escalates for an independent second look any Requirement rated on a single
    weak observation.

  > Rationale: the highest-risk place for a missed finding is a zone the first
  > pass reported as clean; "found nothing" usually means "didn't look hard
  > enough." This is the prong that addresses false negatives. — 0129

- **R14.** Findings surfaced by the completeness sweep **MUST** re-enter
  collection and then be verified under R12 before they can bind a rating. The
  verify and completeness-sweep prongs **MAY** run concurrently where the harness
  supports it.

### QC convergence and stop condition

- **R15.** The collection→QC loop **MUST** terminate. It **MUST** stop when a
  completeness-sweep round surfaces no new in-scope findings and every in-scope
  Requirement is in a terminal evidentiary state (R5) — the converged state.

- **R16.** The loop **MUST** be bounded by a fixed maximum number of
  re-collection rounds so it cannot spin. If the bound is reached before
  convergence, the workflow **MUST** proceed to roll-up and **MUST** report every
  zone left unexamined or unresolved as an explicit limitation; it **MUST NOT**
  silently drop the remaining coverage.

  > Rationale: an unbounded "re-collect until dry" loop can fail to converge on a
  > large or adversarial source; a bound plus honest disclosure of the residue
  > beats both an infinite loop and a silent truncation. — 0129

## Durable spec changes

### To add

None. The QC phase folds into the existing shared evaluation contract in
`evaluation.md`. (It is a behavioral component, but it is tightly entangled with
collection and roll-up; a future split into a `evaluation-qc.md` child spec is
worth revisiting if `evaluation.md` grows, but is not earned now.)

### To modify

- `specs/skills/quality-skill/evaluation.md` — replace the **Rigor levels**
  section with the single-workflow coverage contract (R4–R6), the
  parallel-by-default/serial-fallback execution rule (R7–R9), the preserved
  orchestrator-owned invariant (R10), and the two-pronged QC phase with its
  convergence bound (R11–R16); remove `rigor` from the spec's companion note and
  intro framing.
- `specs/skills/quality-skill/quality-skill.md` — remove the `Rigor:` run-frame
  field and the `--rigor`/`deep` invocation examples; drop the **Rigor** frame
  bullet and the cross-references to *Rigor levels* (per R1, R2).
- `specs/skills/quality-skill/index.md` — update the `evaluation.md` one-line
  description to drop "rigor" and name the QC phase (per R11).
- `specs/skills/quality-skill/workflows/evaluate/feedback-log.md` — remove the
  `rigor` field from the feedback-log frontmatter contract (per R2).

### To rename

None.

### To delete

None.

## Open questions

- **Re-collection bound value (R16).** The spec fixes that a bound MUST exist and
  be honestly disclosed when hit, but leaves the concrete number to the design
  doc / skill prompt (current lean: two re-collection rounds). It is a prompt
  heuristic, not a conformance surface, so it is intentionally not pinned here.
