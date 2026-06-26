# QUALITY.md

> 🚧 **Early alpha — under active construction.** The format, `/quality` skill, and
> `qualitymd` CLI are still evolving and **breaking changes are expected.** Run
> `/quality update` to stay current with the latest skill and CLI versions.

**QUALITY.md** is an [open format](./SPECIFICATION.md) for modeling quality:
what matters most, why, and how, for software, documentation, data, services,
operations, or whatever else your team tends. Use it with the `/quality` agent
skill to align people and AI agents on what good means in that context.

A QUALITY.md file is a Markdown file with a quality model and supporting
context. The `/quality` skill helps set up the file, evaluate quality, and
evolve the model as you learn. That agentic workflow is the primary experience;
the `qualitymd` CLI provides support tooling for validating QUALITY.md files,
managing quality evaluations, and maintaining a QUALITY.md workspace.

## Install

1. Install the agent skill:

```sh
npx skills add qualitymd/quality.md
```

2. Install the CLI:

```sh
npm install -g quality.md
```

## Usage

Invoke the `/quality` skill to manage quality for your project:

```text
/quality setup                                  Get started working with QUALITY.md
/quality                                        Get read-only guidance on the next workflow
/quality evaluate                               Evaluate the quality of your project
/quality evaluate Security                      Evaluate a specific quality factor or characteristic
/quality evaluate Payments                      Evaluate a specific area or project component
/quality evaluate Payments Reliability
                                                Evaluate an area's specific quality
```

For exact or ambiguous scoped evaluations, use qualified model references such
as `/quality evaluate area:payments` or
`/quality evaluate factor:payments::reliability`.

Most users should work with `QUALITY.md` through their coding agent, the
`/quality` skill, or direct edits. The CLI is primarily support tooling for
validation, status, and evaluation records.

To keep the model visible to agents, add a short note to `AGENTS.md` or
`CLAUDE.md`:

```text
See [QUALITY.md](./QUALITY.md) for how this project models and evaluates quality.
```

## Why QUALITY.md

### Manage Quality Debt

As work changes faster, quality erodes quietly unless expectations stay visible
and current. One useful software-specific framing names three accumulating
debts:

- **Technical debt** — code drifting from where it should be.
- **Cognitive debt** — the mounting burden of understanding complex or
  under-documented systems.
- **Intent debt** — software diverging from what stakeholders actually need.

QUALITY.md makes a team's quality expectations explicit and checkable, so those
gaps stay visible and addressable instead of compounding.

> The three-debt framing draws on Margaret-Anne Storey, *The Triple Debt of
> Software Development* ([arXiv:2603.22106](https://arxiv.org/abs/2603.22106)).

### Evaluate Agent Harnessability — and your agent harness

**Agent Harnessability** names how accessible and operable a project is to an
agent — whether its context is visible, its tasks specifiable, its checks
runnable, and its agent's actions safely bounded.

QUALITY.md turns that agent-facing quality into an explicit signal, and you can
evaluate it from two angles:

```text
/quality evaluate Agent Harnessability     How well the whole project equips an agent
/quality evaluate Agent Harness            The quality of the harness your project owns
```

The first rates the project as a whole — where an agent is blocked by missing
context, unclear tasks, weak feedback, advisory-only standards, or unsafe action
scope. The second rates the harness itself — the agent guidance, skills,
environment setup, checks, and guardrails your project owns and maintains —
as artifacts judged on their own quality. Either way, the `/quality` skill turns
the gaps into concrete improvement work.

> The term *harnessability* and the guides-and-sensors framing draw on
> Birgitta Böckeler, *[Harness Engineering](https://martinfowler.com/articles/harness-engineering.html)*
> (martinfowler.com), and OpenAI,
> *[Harness engineering](https://openai.com/index/harness-engineering/)*.

### Quality Beyond Software

QUALITY.md is useful whenever people need to be deliberate about the quality of
something they maintain: software systems, documentation sets, data products,
research reports, service operations, clinical handoffs, legal contracts,
engineering design reviews, classroom plans, household budgets, and other
things people care for.

Those cases do not share one checklist. They share the need to make "good"
visible enough for people to evaluate, learn, and improve, and for coding agents
and AI assistants to follow.

[Install QUALITY.md](#install), then run `/quality setup` with your coding agent
to create a `QUALITY.md` for your project.

## Working with QUALITY.md

A `QUALITY.md` file is your project's **reward signal for quality** — the
explicit, shared definition of *good* that aligns your team, AI assistants, and
coding agents with what matters in *this* context. You capture it once, then run a
**quality loop** that keeps both the work and the signal sharp.

**The quality loop:**

1. **Align** — agree on what *good* means for this context, and capture it as the
   signal.
2. **Evaluate** — grade the work against the signal; every gap is a gap in the
   *work* or a gap in the *bar*.
3. **Learn** — feed what you learned back into the signal as the product and its
   risks evolve.
4. **Improve** — close the gaps in the work.

### Getting started

1. **Run `/quality setup` to make your quality bar visible.** The skill runs a
   guided setup workflow that inspects available context, asks concrete setup
   questions with recommended defaults, and writes a `QUALITY.md` with likely
   quality factors, assessable requirements, needs, risks, unknowns, open
   questions, and agent-accessibility gaps.
   → *A shared starting point for what good means here: the expectations already
   visible in the project, the judgments that still need human input, and the
   context gaps agents and contributors need to close before they can know what
   matters most.*

2. **Run `/quality evaluate` to see where you stand.** The skill analyzes each
   area of your project against the quality requirements defined in your model,
   then provides ratings, findings, evidence limits, and opportunities to improve
   both the work and `QUALITY.md` itself.
   → *A concrete read on where the work meets the bar, where it falls short, and
   where the model needs to become clearer, more complete, or better grounded.*

3. **Refine the model as you learn.** Revise `QUALITY.md` when the evaluation
   reveals missing context, unclear requirements, changed risks, or a quality bar
   that no longer matches the project.
   → *A quality bar that gets clearer and more useful every time you use it.*

   **Tip:** Discuss and apply changes to `QUALITY.md` with your agent. The
   `/quality` skill applies authoring best practices so model changes stay
   well-formed, grounded, and current, with meaningful changes recorded in the
   quality log.

4. **Review the evaluation.** Weigh the quality ratings for each area and
   factor, use the findings to understand the evidence behind the gaps, and turn
   the most important gaps into follow-up work.
   → *A clear read on what to fix now, what to track for later, and what evidence
   or model context is still missing.*

### Keeping the skill and CLI current

**Run `/quality update` to update the `/quality` agent skill and `qualitymd`
CLI.** Keep them up to date to take advantage of the latest improvements to the
efficiency and efficacy of working with QUALITY.md files.

### Keeping the loop running

Once you have a model, the loop keeps going at whatever cadence fits your team:

- **On demand.** Run `/quality evaluate` whenever you need a read — before a
  release, during review, when you inherit an unfamiliar project, or when
  something just feels off. Scope it down when you don't need the whole model:
  `/quality evaluate Payments` or `/quality evaluate Payments Reliability`.
  Qualified references such as `/quality evaluate area:payments` remain
  available for exact addressing.

- **On a cadence.** Make the model and its latest evaluation a recurring team
  review — per sprint, monthly, or whatever rhythm maintainers already use.
  Close gaps in the work,
  and sharpen the model where the bar proved wrong, unclear, or out of date — so
  the shared definition stays current as the product evolves.

- **Recurring.** Use Codex automations, Claude Code routines, or another
  maintainer-owned workflow when you want that cadence to run without someone
  remembering. Keep the loop tied to review habits, not CI or release gates by
  default.

## Example QUALITY.md

The example below is an illustrative software product model. It is not a
default domain or factor set for QUALITY.md — a model can just as well describe
documentation, a data set, a research report, or a service or operation, each
with the factors that matter for that kind of work. See
[Modeling quality across domains](docs/guides/model-quality-across-domains.md)
for a worked non-software documentation example and guidance on keeping model
content domain agnostic.

```markdown
---
title: ACME Payments API
ratingScale:
  - level: outstanding
    title: 🟢 Outstanding
    description: The work clearly exceeds the shared quality bar.
    criterion: "Consistently exceeds the requirement with clear margin."
  - level: target
    title: 🔵 Target
    description: The work meets the shared quality bar.
    criterion: "Meets the expected quality bar."
  - level: minimum
    title: 🟡 Minimum
    description: The work is acceptable, but has gaps worth improving.
    criterion: "Meets the lowest acceptable bar, with visible gaps."
  - level: unacceptable
    title: 🔴 Unacceptable
    description: The work is below the shared quality bar.
    criterion: "Falls below the minimum acceptable bar."
areas:
  payments:
    title: Payments
    source: ./services/payments
    factors:
      reliability:
        title: Reliability
        description: Payment outcomes remain correct under ordinary failures.
        requirements:
          confirmed-payments-are-durable:
            title: confirmed payments are durable
            assessment: >
              A payment is reported as confirmed only after the transaction is
              durably recorded and can be reconciled after a process restart.
      security:
        title: Security
        description: Payment data and privileged operations are protected.
        requirements:
          secrets-stay-out-of-source:
            title: secrets stay out of source
            assessment: >
              Credentials, API keys, and payment-provider tokens are loaded from
              approved runtime configuration and do not appear in source,
              fixtures, logs, or checked-in documentation.
---

# Quality model: ACME Payments API

## Overview

This model describes the quality bar for the ACME Payments API. Good payment
behavior means confirmed transactions are durable, failures are visible, and
payment-provider access is handled without exposing secrets.

## Scope

This model covers the service code, configuration, tests, and operational docs
under `./services/payments`. It does not cover the external payment provider or
the accounting system that consumes payment events.
```

## Format

### Specification

The full format is specified in [`SPECIFICATION.md`](SPECIFICATION.md).

### File Structure

A QUALITY.md file has two layers:

1. **YAML frontmatter** — the structured quality model.
2. **Markdown body** — the judgment context, rationale, scope, needs, risks,
   unknowns, and open questions that help people and agents build, interpret,
   and evaluate the model.

The document begins with the YAML frontmatter. The Markdown body can be empty,
but it is where the model explains its purpose and context.

### Model Schema

The root model is an area plus a model-wide `ratingScale`.

```yaml
title: <string>                 # Required
description: <string>           # Optional
ratingScale:                    # Required, ordered best to worst
  - level: <rating-level-id>    # Required, unique within the scale
    title: <string>             # Required
    description: <string>       # Recommended
    criterion: <string>         # Required
source: <string>                # Optional
factors:                        # Optional*
  <factor-name>:
    title: <string>             # Required
    description: <string>       # Recommended
    factors:                    # Optional
      <sub-factor-name>: <Factor>
    requirements:               # Optional
      <requirement-name>: <Requirement>
requirements:                   # Optional*
  <requirement-name>:
    title: <string>             # Required human-facing statement
    assessment: <string>        # Required, exactly one
    factors: [<factor-name>]    # Required for direct area requirements
    ratings:                    # Optional per-level criteria
      <rating-level-id>: <criterion>
areas:                          # Optional*
  <area-name>: <Area>
```

At least one of `factors`, `requirements`, or `areas` must be supplied.
Areas can nest recursively. `ratingScale` exists only on the root model.

### Core Concepts

| Concept      | Meaning                                                        |
| ------------ | -------------------------------------------------------------- |
| Model        | The root quality model in a QUALITY.md file.                   |
| Area         | The thing being evaluated.                                     |
| Source       | The material assessed for an area, such as a path or selector. |
| Factor       | A quality dimension that matters for an area.                  |
| Requirement  | A specific quality expectation.                                |
| Assessment   | The means of checking a requirement against an area source.    |
| Finding      | An observation produced by an assessment.                      |
| Rating Scale | The ordered model-wide scale used to rate results.             |

## CLI Quick Reference

| Task                    | Command                          |
| ----------------------- | -------------------------------- |
| Show format spec        | `qualitymd spec`                 |
| Show frontmatter schema | `qualitymd schema`               |
| Create a starter file   | `qualitymd init [path]`          |
| Validate a file         | `qualitymd lint [path]`          |
| Fix lint issues         | `qualitymd lint --fix [path]`    |
| Show project status     | `qualitymd status [path] --json` |
| Show version info       | `qualitymd version --json`       |
| Check for updates       | `qualitymd update --check`       |
| Show command help       | `qualitymd <command> --help`     |

## Status

The QUALITY.md format, `qualitymd` CLI, and `/quality` skill are early and
under active development. Expect the format and tooling to change as they
mature.

## Contributing

Contributor setup and local tasks live in [`CONTRIBUTING.md`](CONTRIBUTING.md).
