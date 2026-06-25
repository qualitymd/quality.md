---
type: Runtime Guide
title: Authoring Agent Harnessability
description: Guidance for modeling Agent Harnessability as a model-wide factor.
tags: [quality, authoring, guide]
---

# Authoring Agent Harnessability

Read this when:

- creating, revising, reviewing, or evaluating the Agent Harnessability factor, its sub-factors, or its projection boundary.

Depends on:

- `../authoring.md`

---

## Agent Harnessability for agent-collaborated composite roots

For an agent-collaborated composite root, propose **Agent Harnessability** as a
model-wide umbrella factor, using `agent-harnessability` as the recommended
stable key:

```yaml
agent-harnessability:
  title: Agent Harnessability
  description: >
    The degree to which the project's checked-in materials, tools, workflows,
    feedback signals, standards, and action limits equip an AI agent to
    understand the project, take scoped work, operate the environment, preserve
    and resume state, verify its output, and stay safely bounded while preserving
    clear human direction, review, and accountability.
```

Agent Harnessability is a deliberate umbrella: carry no requirements on the
parent — its assessment lives in the sub-factors, and its rating comes from
rolling those up, the same way a grouping area has no local rating. Decompose it
into independently assessable sub-factors, each with its own non-overlapping
share. Name the quality the project exhibits, not the harness
artifact itself: guides, sensors, sandboxes, scripts, and gates are evidence for
the factor, not the factor.

- **Do** include Agent Harnessability by default for an agent-collaborated
  composite root. *A thin or absent harness is not a reason to drop the factor; it
  is a low rating and a finding, the same way no tests is not a reason to omit
  `testability`.*
- **Do** use these sub-factors as the Agent Harnessability decomposition,
  adapting the example requirements to the entity's domain:
  - **agent-accessibility** — the degree to which decision-relevant knowledge,
    structure, intent, and observable behavior are present and intelligible in
    materials an agent can reach in-context at the moment of work. Example
    requirements: a stable minimal agent entry point points to deeper material
    without blowing the context budget; deeper context is progressively
    disclosed only when relevant; decision-relevant knowledge is durably recorded
    and discoverable; needed behavior signals are reachable and machine-parseable.
    Boundary: this is in-context machine reachability and context selectivity,
    not general human understandability; cede in-flight handoff/progress records
    to `continuity`, durable decision-record existence to traceability, and raw
    instrumentation to observability.
  - **task-specifiability** — the degree to which a unit of work can be handed to
    an agent as a self-contained, scoped assignment with explicit success criteria
    and known boundaries. Example requirements: goals, non-goals, success
    criteria, done criteria, decomposition, and starting points are articulable
    before work begins; completion checks compare the result with the original
    task before declaring success. Boundary: this frames the task;
    self-verifiability checks whether the criteria can be confirmed.
  - **agent-operability** — the degree to which an agent, including a fresh
    session, can establish and operate the working environment from recorded
    materials: the tools, data, accounts, systems, and known starting state the
    work runs on. Example requirements: a fresh session reaches ready-to-work
    state without human-led setup; needed access and inputs are recorded; tools
    expose agent-useful affordances and output; the act-then-observe loop is
    bounded and inspectable. Boundary: this turns operability toward the agent's
    working environment; `continuity` covers resuming with prior state and
    decisions; containment-of-action confines the access this equips.
  - **continuity** — the degree to which an agent can preserve state and resume
    useful work across long-running tasks, compaction, interruption, handoff, and
    fresh sessions. Example requirements: progress or handoff artifacts capture
    current state, decisions made, remaining work, verification status, blockers,
    and next steps; resumptions do not depend on unrecoverable chat history;
    progress records reduce false completion and context-anxiety failure modes.
    Boundary: agent-operability covers a fresh session reaching a ready-to-work
    environment, while continuity covers resuming with prior state, decisions,
    and progress; agent-accessibility covers durable reachable knowledge, while
    continuity covers the progress and handoff record of an in-flight task.
  - **self-verifiability** — the degree to which the project gives an agent
    machine-readable verification signals it can run or inspect on demand to
    confirm whether its own work is correct, with remediation-bearing output.
    Example requirements: a
    deterministic signal returns objective pass/fail without human setup;
    inferential evals cover behavioral or non-deterministic outcomes without a
    deterministic oracle; failures are fast, grounded in concrete evidence,
    context-aware, and point toward the fix with the rule's rationale; traces,
    run logs, or evaluation records expose what happened and why. Boundary: this
    owns the actionability and evidentiary quality of feedback, including visible
    reviewable exceptions rather than a binary pass/fail-only signal;
    enforcement-of-standards constrains whether and how those exceptions may
    bypass a gate.
  - **enforcement-of-standards** — the degree to which stated quality standards
    hold regardless of agent behavior because non-compliant output is prevented by
    gates or equivalent domain-neutral controls rather than advisory prose.
    Example requirements: quality invariants are blocked or routed through
    reviewable input, output, or tool guardrails; controls are mutually consistent
    and high-signal; suppression escapes are constrained. Boundary: this prevents
    non-compliant output or action; self-verifiability makes the feedback
    actionable and observable; containment-of-action prevents out-of-scope action.
  - **containment-of-action** — the degree to which an agent's permitted actions
    are confined by enforced limits such as sandboxes, permission allowlists, and
    approval gates. Example requirements: an erroneous or unattended run cannot
    escalate scope, reach sensitive resources, or take consequential real-world
    action such as filing to a court, moving money, emailing a client, or deleting
    records without approval. Boundary: cede external threat posture to security
    where an area carries it; this confines the agent's own actions while working.
    Audit trails make consequential actions reviewable after the fact.
- **Do** keep Agent Harnessability separate from the agent harness constituent.
  *Agent Harnessability rates how each constituent equips an agent. The
  agent-harness area rates the checked-in, project-owned governing artifacts'
  own quality (see
  [Carry the recurring use-context constituents](model-structure.md#carry-the-recurring-use-context-constituents)).*
- **Do** recognize an existing `harnessability` factor with the legacy
  six-sub-factor shape as semantic coverage of the same model-wide concern.
  *When you are already authoring or revising the model, recommend renaming it to
  `agent-harnessability` / Agent Harnessability and adding `continuity` unless
  the project has an explicit reason to preserve the old key.*
- **Avoid** adding "improve the harness over time" as an eighth sub-factor. *That
  is the model-wide learn loop improving this equipping, not a sibling quality
  beside the seven.*
