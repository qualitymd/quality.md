---
type: Functional Specification
title: Harnessability factor — functional spec
description: What the /quality skill must teach about harnessability — a model-wide umbrella factor (with six sub-factors) proposed by default for agent-collaborated composite roots, the factor projection of the agent-collaboration concern.
tags: [skill, authoring, setup, factors, harnessability, agentic]
timestamp: 2026-06-24T00:00:00Z
---

# Harnessability factor — functional spec

Companion to the [Harnessability factor](../0081-harnessability-factor.md) change
case. This spec states *what* the skill's factor guidance must say; the
[design doc](design.md) covers *how* the six sub-factors were derived and verified
and *why* the concern leads with its factor projection. It governs the bundled
skill ([`skills/quality/`](../../../skills/quality/)) and its functional-spec
mirror ([`specs/skills/quality-skill/`](../../../specs/skills/quality-skill/)), and
defers the QUALITY.md format itself to
[`SPECIFICATION.md`](../../../SPECIFICATION.md).

The key words "MUST", "MUST NOT", "SHOULD", "SHOULD NOT", and "MAY" are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

The skill carries the agent-collaboration concern only as a **constituent** — the
agent harness — and 0080 made that constituent model-by-default. But the concern's
most useful projection across a composite root is a **factor**: *how well does
each part of this project equip an agent to work on it well?* That question recurs
across constituents (the server, the schema, the tests, the docs are each more or
less so), which is the guide's own signature of a model-wide factor. Leading with
the factor projection also makes the concern robust to thinness the way
`testability` already is — a thin or absent harness rates the factor low and
surfaces a finding everywhere, instead of being deferrable as an immature area.

The quality is the one the harness-engineering literature names: agent legibility,
steerability, and verifiability of the whole project — *guides and sensors* around
the model (Böckeler), *"what the agent can't see doesn't exist"* (OpenAI). This
case names that quality **harnessability** and decomposes it into six sub-factors
that trace the agent's working loop, so an evaluator assesses a project's
agent-equipping directly rather than inferring it from one steering artifact's
substance.

## Scope

Covered: harnessability as a model-wide umbrella factor and its six sub-factors;
the rule that setup and the authoring guide propose it by default for an
agent-collaborated composite root; the three-projections relationship between the
factor, the agent-harness constituent, and the agent audience; the boundary
discipline that keeps the sub-factors non-overlapping with each other, with
existing common factors, and with the harness constituent; the Top 10 check for
its coverage; and the README use case.

Deferred / non-goals: no change to the QUALITY.md format or schema, to
`SPECIFICATION.md`, or to the CLI — recommending a factor is authoring judgment,
not format semantics. The continuous improvement of this equipping over time stays
the existing **model-wide learn-loop concern**, not a seventh sub-factor. Fleet
orchestration and the cross-project human-attention economy sit above a project
Area and are out of scope. Re-checking this repo's own `QUALITY.md` against the new
factor is a follow-up.

## Requirements

### Harnessability is a model-wide umbrella factor

The authoring guide (and its spec mirror) **MUST** teach **harnessability** as a
model-wide (cross-cutting) factor for an agent-collaborated composite root: the
degree to which a project equips an agent to do good work on it largely
unsupervised — to perceive the project, take a scoped unit of work, act within it,
confirm its own output, and stay safely bounded — from the project's own
checked-in materials and tooling, rather than from out-of-band human knowledge or
synchronous supervision. The guide **MUST** state that it matters because the
scarce resource in agent-collaborated work is human attention.

The guide **MUST** present harnessability as a **deliberate umbrella**: it does not
roll up directly but decomposes into the sub-factors below, each independently
assessable and contributing one non-overlapping share. It **MUST** name the
*quality the project exhibits* and **MUST NOT** let the name denote the harness
artifact — guides, sensors, sandboxes, and scripts are *evidence*, not the factor
(per [Name the quality, not the practice](../../../skills/quality/guides/authoring.md)).

> Rationale: the concern recurs across constituents, which is the guide's own test
> for a model-wide factor; naming the quality (not the harness) keeps it from
> reading as a renamed constituent. — 0081

### Proposed by default, never dropped for thinness

The setup workflow and the authoring guide (and their spec mirrors) **MUST** have
the author propose harnessability by default when the root is an
agent-collaborated composite. The thinness or absence of a project's harness
**MUST NOT** be given as a reason to omit the factor; a thin or absent harness is
expressed as a low rating with a finding, not a dropped factor.

> Rationale: this is the factor-side counterpart to 0080's no-silent-omission rule,
> and the same robustness the guide already gives `testability` ("no tests is not a
> reason to omit `testability`"). It is the direct fix for the field run that
> deferred a "thin" harness. — 0081

### The six sub-factors

The guide (and its spec mirror) **MUST** present harnessability's sub-factors as
the following six, each named as a quality, each with an operational definition and
assessable example requirements, and each carrying a boundary note that keeps its
roll-up share distinct. The example requirements are illustrative, not an
exhaustive or mandatory checklist.

1. **agent-accessibility** — the degree to which a project's decision-relevant
   knowledge, structure, intent, and observable behavior are present and
   intelligible in materials an agent can reach *in-context at the moment of work*
   (durably recorded docs, a navigable map of where things live, scoped
   instruction files, consistently structured and searchable materials, decision
   records, machine-parseable signals of observable behavior), so the agent can
   build an accurate model of the project from the project itself. *Matters
   because what an agent cannot reach or parse in-context effectively does not
   exist to it.*
   - Example requirements (illustrative; the materials differ by domain — a
     version-controlled repo for software, a shared matter file and document store
     for a legal case, recorded accounts and rules for a budgeting project): a
     stable, minimal agent entry point exists, points to deeper material rather
     than inlining it, and stays within a context budget; decision-relevant
     knowledge is durably recorded and discoverable in-context rather than held
     out-of-band or in a person's head; runtime behavior the agent needs is
     reachable and machine-parseable in-context.
   - Boundary: vs `understandability` (general human comprehensibility) — this is
     in-context *machine* reachability and intelligibility; route human-prose
     clarity to `understandability`. vs `traceability` — cede "is a decision
     durably recorded" to `traceability`; scope this to in-context reachability of
     that record. vs `observability` — if the area carries `observability`, assess
     only the in-context reachability of runtime data here, not raw instrumentation.

2. **task-specifiability** — the degree to which a project lets a unit of work be
   handed to an agent as a self-contained, scoped assignment with explicit success
   criteria and known boundaries, so the agent knows what "done" means and where
   the work stops before it starts. *Matters because an ambiguous or unbounded task
   forces the human to re-specify work mid-flight — the supervision harnessability
   exists to remove.*
   - Example requirements: units of work can be expressed with explicit success
     criteria the agent can target and check against; scope and non-goals are
     articulable; task entry points (where to start, what to read first) are
     discoverable in-context.
   - Boundary: vs `agent-accessibility` (can the agent understand the *project*) —
     this is whether a *task* can be framed with goal, criteria, and bounds. vs
     `self-verifiability` — criteria are *stated* here, *mechanically checked*
     there.

3. **agent-operability** — the degree to which an agent, including a fresh session
   with no prior memory, can establish and operate the project's working
   environment — reaching a known, ready-to-work state and connecting to the tools,
   data, accounts, and systems the work runs on — from recorded materials rather
   than tacit knowledge or human-led setup. *Matters because an agent that cannot
   reliably stand up and operate the working environment can neither act, observe,
   nor verify.*
   - Example requirements (illustrative; the working environment differs by domain
     — a runnable dev environment for software, connected accounts and a ledger for
     a budgeting project, a loaded matter with its documents and research tools for
     a legal case): a fresh session reaches a ready-to-work state from recorded
     setup, not human-led onboarding; the access, credentials, and inputs the work
     needs are recorded, not tacit; the agent works from a consistent, known
     starting state and can act without disrupting live state others depend on
     (isolated, reproducible execution — e.g. per-session worktrees — where the
     domain supports it); the act-then-observe cycle (e.g. a software inner loop)
     stays within a bounded, stated turnaround.
   - Boundary: vs `operability` (the evaluated entity's general fitness to be run or
     operated in its live setting, often by humans) — agent-operability is that same
     lens turned on the *agent's working environment*: whether an agent doing the
     project's work can stand it up and operate it, fast and repeatably, just as
     `agent-accessibility` turns `understandability` toward the agent. vs
     `agent-accessibility` — accessibility is whether the agent can *know* the
     project; this is whether it can *operate* its workspace. vs
     `containment-of-action` — this *equips* working access; that *confines* it. The
     bring-up confirmation is owned here ("can the working environment stand up");
     confirming a *change* is correct is `self-verifiability`.

4. **self-verifiability** — the degree to which a project gives an agent objective,
   machine-readable signals it can run on demand to confirm whether its own change
   is correct, whose output carries the remediation the agent needs to self-correct
   (an on-demand, objective check the agent runs on its own work and reads the
   result of). *Matters because without a check it can run, "looks done" is the
   only signal the agent has.*
   - Example requirements (illustrative; the check differs by domain — pass/fail
     tests, type checks, and build exit codes for software; a citation and
     rules-of-court check for a legal filing; a reconciliation that proves the
     ledger balances for a budgeting project): a single command or action returns an
     objective pass/fail for a change without human setup; failure output is
     remediation-bearing (names what is wrong, points toward the fix);
     non-deterministic or behavioral outcomes have runnable evals that gate the
     work, and prior real-world failures are captured back as coverage.
   - Boundary: vs `testability` (general ease of testing the entity) — this is
     whether the *agent itself* can obtain a fast, objective, remediation-bearing
     pass/fail in-loop. It **owns** the actionability-of-feedback requirement (a
     check's failure message carrying the fix); `agent-accessibility` keeps only
     static guidance prose. vs `enforcement-of-standards` — a check here is a signal
     the agent *chooses to consume*; a gate there *binds regardless of the agent*.

5. **enforcement-of-standards** — the degree to which a project's stated quality
   standards hold *regardless of agent behavior* — non-compliant output is
   prevented by deterministic gates rather than depending on the agent honoring
   advisory prose. *Matters because LLM compliance with instructions is
   probabilistic; a written rule an agent may ignore is not a guarantee.*
   - Example requirements (illustrative; the gate differs by domain — structural
     tests, linters, and CI invariants that block a merge for software; a template
     that refuses a filing missing a mandatory disclosure for a legal case; a close
     that will not complete until every transaction is categorized and the books
     reconcile for a budgeting project): stated quality invariants are enforced by
     tooling that blocks the non-compliant output, not by docs alone; enforced
     controls are mutually consistent and high-signal, not contradictory rules that
     push the agent into over-engineering spirals or feedback overload; suppression
     escape hatches (inline disables, forced-pass) are constrained so the agent
     fixes rather than silences violations.
   - Boundary: named as a quality (non-compliant output is *prevented*), not a tally
     of gates. vs `self-verifiability` — bind-regardless-of-agent vs
     agent-consumed-signal. vs `containment-of-action` — this prevents
     *non-compliant output*; that confines *out-of-scope action*.

6. **containment-of-action** — the degree to which an agent's permitted actions are
   confined so out-of-scope or unsafe action is contained by enforced limits
   (sandboxing, permission allowlists, action-blocking hooks) rather than trusted
   to advisory prose. *Matters because an unbounded agent's mistaken action has
   unbounded blast radius.*
   - Example requirements: the agent's permitted actions are confined so an
     erroneous or unattended run cannot escalate scope, reach sensitive resources,
     or take a consequential real-world action without an approval gate (e.g. file
     to a court, move money, email a client, or delete records); containment is
     enforced (sandboxed filesystem/network, scoped permissions), not merely
     documented.
   - Boundary: named as a quality (out-of-scope action is *prevented*), not a tally
     of sandboxes. vs `security` (the entity's posture against external threats) —
     if the area carries `security`, assess threat posture there and confine this to
     containment of the *agent's own* actions while doing the project's work, to
     avoid double-counting sandbox/permission evidence. vs `enforcement-of-standards`
     — confines *action*, not *output*.

> Rationale: the six trace the agent's loop — perceive (agent-accessibility),
> receive direction (task-specifiability), operate (agent-operability), verify
> (self-verifiability), stay bounded (enforcement-of-standards +
> containment-of-action). The two boundedness sub-factors were split from a single
> candidate because enforcement-of-output and containment-of-action can hold
> independently (a project can gate its output yet still permit a destructive
> action). — 0081

### Three projections of the agent-collaboration concern

The guide (and its spec mirror) **MUST** present harnessability, the agent harness,
and the agent as the three projections of one concern under the existing
[three-projections rule](../../../skills/quality/guides/authoring.md): the
**factor** (harnessability, a quality of every constituent), the **constituent**
(the agent harness — the owned steering artifact, per 0080), and the **audience**
(the agent the project serves). It **MUST** keep the agent harness as its
constituent rather than replacing it with the factor.

The guide **MUST** give the double-count boundary between the factor and the
constituent: harnessability-the-factor rates how well *each constituent equips an
agent*; the agent-harness *area* rates the *steering artifact's own quality* (is
its map accurate, current, a map not a manual). The guide **MUST** state that
harnessability is not assessed *on* the agent-harness area as a recursion of the
same evidence.

> Rationale: same discipline the guide already applies to `secure` projecting as a
> factor, a threat-model constituent, and an auditor audience — name the projection
> meant and model it once. — 0081

### Continuous improvement stays model-wide

The guide **MUST NOT** add a sub-factor for improving the harness over time. The
guide **MUST** route that concern to the model's existing model-wide
continuous-improvement / learn-loop concern, not to a harnessability sub-factor.

> Rationale: improvement-of-the-equipping is a meta-loop over the other six; making
> it a sibling would let improvement-of-X be rated beside X in the same roll-up. —
> 0081

### Top 10 check

The Top 10 checks guide (and its spec mirror) **MUST** flag an agent-collaborated
composite root that does not carry harnessability (or its sub-factors) among its
model-wide factors, routing the author to the authoring guide. The check **MUST
NOT** flag a non-agent-collaborated or throwaway entity that would not carry the
factor (the same not-germane carve-out the constituent checks use).

### README use case

`README.md` **MUST** gain a `### Evaluate and Improve Agent Harnessability`
subsection under `## Why QUALITY.md`, framing a `QUALITY.md` as a way to model and
raise harnessability — the explicit signal that steers an agent and the basis for
evaluating its work. It **MUST** cite the harness-engineering sources (Böckeler /
Thoughtworks, OpenAI, Augment Code). Secondary URLs surfaced in research **MUST**
be verified before citing.

## Durable spec changes

### To add

None.

### To modify

- `specs/skills/quality-skill/guides/authoring-md.md` — require the harnessability
  model-wide factor, its umbrella character, the six sub-factors, the factor
  projection of the use-context agent-collaboration concern, and the
  factor/constituent double-count boundary (per Harnessability is a model-wide
  umbrella factor, The six sub-factors, Three projections, and Continuous
  improvement stays model-wide above).
- `specs/skills/quality-skill/workflows/setup.md` — require setup to propose
  harnessability by default for an agent-collaborated composite root, never
  dropping it for thinness (per Proposed by default above).
- `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md` — add the
  harnessability-coverage check with the not-germane carve-out (per Top 10 check
  above).
- `specs/skills/quality-skill/guides/log.md` and
  `specs/skills/quality-skill/workflows/log.md` — record the contract revision.

### To rename

None.

### To delete

None.
