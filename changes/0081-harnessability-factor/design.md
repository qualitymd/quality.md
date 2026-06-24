---
type: Design Doc
title: Harnessability factor — design
description: Why the agent-collaboration concern leads with its factor projection, how the six sub-factors were derived and adversarially verified, and the boundary discipline that keeps them non-overlapping.
tags: [skill, authoring, setup, factors, harnessability, agentic]
timestamp: 2026-06-24T00:00:00Z
---

# Harnessability factor — design

## Context

Answers the [functional spec](spec.md) for the
[0081 change case](../0081-harnessability-factor.md). The work is a guidance
addition across the bundled skill and its spec mirror, plus a README use case — no
code. The design questions are *which projection* the agent-collaboration concern
should lead with, *what* the sub-factors are, and *how* to keep them from
double-counting against each other, the existing common factors, and the agent
harness constituent.

## Approach

### 1. Lead with the factor projection, keep the constituent

The guide's [three-projections rule](../../../skills/quality/guides/authoring.md)
already says a stewardship concern projects as a **factor**, a **constituent**, and
an **audience**. The agent-collaboration concern was modeled only as a constituent
(the agent harness). The fix is not to move it but to *add the projection that was
missing*: harnessability as the **factor**, the agent harness as the
**constituent** (unchanged, still model-by-default per 0080), the agent as the
**audience**.

The factor is the better *lead* because the quality recurs across constituents —
every part of a repo is more or less legible, steerable, and verifiable to an agent
— which is the guide's own definition of a model-wide factor. And a factor is never
dropped for thinness, where an area can be deferred as immature. That property is
the whole point: the field failure (a "thin" harness deferred to "next iteration")
is structurally impossible once the concern is a factor, exactly as `testability`
survives a project with no tests.

### 2. Derive the sub-factors from the harness-engineering discourse, then verify adversarially

The six sub-factors are not invented; they are the convergent dimensions of the
2025–2026 harness-engineering literature (Böckeler/Thoughtworks' *guides and
sensors*; OpenAI's *agent legibility* and *system-of-record*; Augment's
*constraint harnesses / feedback loops / quality gates*), clustered and then run
through an adversarial overlap critique against QUALITY.md's existing factor
families. They trace the agent's working loop:

| Loop step                     | Sub-factor               |
| ----------------------------- | ------------------------ |
| perceive the project          | agent-accessibility      |
| receive a scoped unit of work | task-specifiability      |
| operate the environment       | agent-operability        |
| confirm its own output        | self-verifiability       |
| stay bounded (output)         | enforcement-of-standards |
| stay bounded (action)         | containment-of-action    |

The critique forced four corrections that shaped the final set: it dropped a
*runtime-transparency* candidate (a renamed `observability`), folding only its
distinct slice — in-context reachability of runtime data — into
`agent-accessibility`; it dropped a *self-correctiveness* candidate as a meta-loop
that improves the others (routed to the model-wide learn loop, not a sibling); it
added `task-specifiability` as the genuinely-missing front-of-loop dimension; and
it pinned the remediation-bearing-feedback requirement to `self-verifiability`
alone, leaving `agent-accessibility` only static guidance prose.

### 3. Split boundedness into output vs action

The critique's fused *agent-boundedness* candidate bundled two qualities that hold
independently: enforcement of standards on *output* (a project can block
non-compliant output) and containment of *action* (yet still permit a destructive
action). They are modeled as two sub-factors so one rating cannot absorb two
non-correlated signals. `enforcement-of-standards` is a quality gate on the agent's
output; `containment-of-action` is blast-radius control and explicitly cedes
external threat posture to the `security` common factor where an area carries it.

### 4. Police the boundaries

Each sub-factor carries a boundary note in the spec because the failure mode of an
umbrella factor is silent double-counting. The two that matter most: the
factor-vs-constituent boundary (harnessability rates *each constituent's*
agent-equipping; the agent-harness *area* rates the *steering artifact's own*
quality; harnessability is not assessed *on* that area, which would recurse the
same evidence), and the sub-factor-vs-common-factor boundaries
(`agent-accessibility` vs `understandability`/`observability`,
`agent-operability` vs `operability`, `self-verifiability` vs `testability`,
`containment-of-action` vs `security`). The discipline is the same the guide
already applies to `secure` projecting three ways: name the projection meant, model
it once.

### 5. Generalize the operate sub-factor and keep every definition domain-agnostic

QUALITY.md is quality-domain agnostic in what a model evaluates, so the sub-factor
definitions must not assume the evaluated project is software under development. The
first synthesis defined the operate sub-factor as *bringing up a development
environment* — a runnable build, a checked-in setup script, per-worktree isolation,
an inner loop — which silently made software-product quality the default. It is
renamed **agent-operability** and redefined as the agent's ability to establish and
operate its *working environment* — the tools, data, accounts, and systems the work
runs on — from recorded materials; a runnable dev environment is one instance, a
budgeting project's connected accounts and a legal case's loaded matter are others.

The rename also names the *quality* (operability), not a state (readiness), per
[Name the quality, not the practice](../../../skills/quality/guides/authoring.md),
and completes a pattern the loop's core already half-followed: **agent-accessibility**,
**agent-operability**, and **self-verifiability** are each a general common factor
scoped to the agent's loop and pinned against it in their boundary
(`understandability`, `operability`, `testability` for perceive / operate / verify).
The cold-start, reconstitute-from-recorded-materials emphasis the old name carried
moves into the description, where subject-specific interpretation belongs.

The other five keep domain-neutral definitions with software as one illustrative
instance among several (a legal filing's citation check, a budget's reconciliation),
not the definition. `containment-of-action`'s sandbox/permission language is agentic
*use context*, not a software *modeled domain*, so it stays.

## Alternatives

- **Make the agent harness a first-class area and stop there (status quo + 0080).**
  Rejected: it measures the concern by one steering artifact's substance, and an
  area is deferrable when thin — the exact field failure. The factor is the
  robust projection; the constituent stays, it just isn't the lead.
- **Name the factor `legibility` (or fold into `understandability`).** Rejected on
  the user's call: harnessability is broader than legibility (it adds steerability
  and verifiability), and a deliberate umbrella that maps to the
  "harness engineering" term of art is more useful than stretching
  `understandability` to an agent audience. `agent-accessibility` carries the
  legibility slice as a sub-factor.
- **Name the operate sub-factor plain `operability` (no prefix).** Rejected: the
  model already carries `operability` as the common factor for the *operate*
  stewardship concern (the entity's fitness to be run in its live setting), so an
  unprefixed name would give two different qualities one identifier and collapse the
  boundary that separates them. The `agent-` prefix buys the distinction, the same
  way `agent-accessibility` is distinguished from `understandability`.
- **Keep the four-factor "static equipping" spine** (drop `task-specifiability`).
  Rejected: the front of the loop — can the *work* be scoped and given success
  criteria — is distinct from perceiving the project, and was smeared across
  `agent-accessibility`/`agent-operability` in the first synthesis.
- **Keep boundedness fused.** Considered; rejected by the user in favor of the
  split, since enforcement-of-output and containment-of-action are independently
  variable.
- **Keep `runtime-transparency` / `self-correctiveness` as sub-factors.** Rejected
  by the overlap critique: the first is renamed `observability`; the second is a
  meta-loop belonging to the model-wide continuous-improvement concern.
- **Push harnessability into `SPECIFICATION.md`.** Rejected for the same layering
  reason 0080 used: which factors an entity carries is authoring judgment, not
  format semantics.

## Trade-offs & risks

- **Umbrella complexity.** Six sub-factors is a lot to carry; the mitigation is the
  loop mnemonic and the rule that the umbrella does not roll up directly — each
  sub-factor is assessed and rated on its own.
- **Boundary creep.** The sub-factors sit close to several common factors; without
  the boundary notes an evaluator double-counts. The notes are load-bearing, not
  decoration, and the factor-vs-constituent recursion is called out explicitly.
- **Scope confusion.** Harnessability is project-scoped; readers may expect it to
  cover fleet orchestration or the cross-project attention economy. The spec scopes
  those out so the umbrella definition and the sub-factor set stay matched.

## Open questions

- Whether this repo's own dogfooded `QUALITY.md` should gain the harnessability
  factor (and how its sub-factors rate here) is left to the follow-up named in the
  change case, so this case stays a pure guidance addition.
- Whether `enforcement-of-standards` and `containment-of-action` should later
  generalize beyond agent-collaborated roots (they resemble general CI-gate and
  least-privilege qualities) is left open; this case scopes them to harnessability.
