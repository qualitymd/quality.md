# QUALITY.md

**QUALITY.md** is an [open format](SPECIFICATION.md), agent skill, and CLI for
managing the quality of your AI assistant and coding-agent projects. Use it to [improve project quality](#getting-started), [engineer quality loops](#loop-engineering-the-quality-loop-stack), [compound learning](#the-outer-loop-dailyweekly), and achieve alignment around _what quality means for what matters most_.

QUALITY.md helps teams **move quality judgment up the loop
stack**, enabling a continuous and comprehensive approach to improving and maintaining project quality instead of ad hoc prompts, passive skills, reactive reviews, and low-level inspections.

> 🚧 **Early alpha — under active construction.** The format, `/quality` skill, and
> `qualitymd` CLI are still evolving. Breaking changes should be manageable with assistance from the latest skill updates. Run
> `/quality update` to stay current with the latest skill and CLI versions.

## Install

1. Install the agent skill:

```sh
npx skills add qualitymd/quality.md
```

2. Install the CLI:

```sh
npm install -g quality.md
```

The npm launcher installs a standalone platform executable. Managed installers,
Homebrew, and direct release archives are also supported; the executable does
not require Node.js, Bun, or Go at runtime.

**Stay up to date**: Invoke `/quality update` to update both the skill and CLI to the latest compatible versions to take advantage of the latest improvements.

## Getting started

### Setup

Invoke `/quality setup` for guided creation of your initial `QUALITY.md` tailored for your project.

### Loop engineering: the Quality Loop Stack

#### The Middle/Meta Loop (hourly/daily)

1. **Evaluate** with `/quality evaluate` to
   create a quality evaluation report. Under the hood the skill drives
   `qualitymd evaluation run`, the CLI's deterministic evaluation runner, so
   evaluations behave the same across agents and harnesses.
2. **Review** evaluation report quality ratings, assessment
   findings, and improvement recommendations by project area or quality factor.

#### The Inner Loop (continuously)

3. **Act** by implementing or delegating handoff-ready recommendations to people, agents, or other agent loops.

#### The Outer Loop (daily/weekly)

4. **Improve** by invoking `/quality improve` to refine
   `QUALITY.md` and compound learning.

> The three-loop framing and the term _middle loop_ draw on Annie Vella,
> _[The Middle Loop](https://annievella.com/posts/the-middle-loop/)_. The
> loop-engineering framing also draws on Latent Space,
> _[Loopcraft: The Art of Stacking Loops](https://www.latent.space/p/ainews-loopcraft-the-art-of-stacking)_.

## Why QUALITY.md

### Manage quality debt

As work changes faster, quality erodes quietly unless expectations stay visible
and current. One useful software-specific framing names three accumulating
debts:

- **Technical debt** — code drifting from where it should be.
- **Cognitive debt** — the mounting burden of understanding complex or
  under-documented systems.
- **Intent debt** — software diverging from what stakeholders actually need.

QUALITY.md makes a team's quality expectations explicit and checkable, so those
gaps stay visible and addressable instead of compounding.

> The three-debt framing draws on Margaret-Anne Storey,
> _[The Triple Debt of Software Development](https://arxiv.org/abs/2603.22106)_
> (arXiv:2603.22106).

### Evaluate agent harnessability and your agent harness

**Agent harnessability** names how accessible and operable a project is to an
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

> The term _harnessability_ and the guides-and-sensors framing draw on
> Birgitta Böckeler, _[Harness Engineering](https://martinfowler.com/articles/harness-engineering.html)_
> (martinfowler.com), and OpenAI,
> _[Harness engineering](https://openai.com/index/harness-engineering/)_.

### Quality beyond software

QUALITY.md is useful whenever people need to be deliberate about the quality of
something they maintain: software systems, documentation sets, data products,
research reports, service operations, clinical handoffs, legal contracts,
engineering design reviews, classroom plans, household budgets, and other
things people care for.

Those cases do not share one checklist. They share the need to make "good"
visible enough for people to evaluate, learn, and improve, and for coding agents
and AI assistants to follow.

### Choose quality factors deliberately

QUALITY.md models do not import a default checklist of factors. A good factor set
is **comprehensive**, **proportionate**, and **sustainable**: it covers what
matters, gives most attention to what matters most, and remains usable over time.

Each individual factor should be **consequential**, then refined until it is
bounded, operational, traceable, and neutral: a clear quality lens, grounded in
the entity's needs and risks, assessable through requirements, and not confused
with an implementation, workflow, tactic, or metric. The `/quality` skill applies
this guidance when setting up or improving a model; see
[Authoring factors](skills/quality/guides/authoring/factors.md) for the detailed
guidance.

[Install QUALITY.md](install.md), then run `/quality setup` with your coding agent
to create a `QUALITY.md` for your project.

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
Shared report terms and fixed vocabularies are defined in
[`glossary.md`](glossary.md).

### File structure

A QUALITY.md file has two layers:

1. **YAML frontmatter** — the structured quality model.
2. **Markdown body** — the judgment context, rationale, scope, needs, risks,
   unknowns, and open questions that help people and agents build, interpret,
   and evaluate the model.

The document begins with the YAML frontmatter. The Markdown body can be empty,
but it is where the model explains its purpose and context.

### Model schema

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
    factors: [<factor-name>]    # Required for direct area requirements; same-area factors only
    ratings:                    # Optional per-level criteria
      <rating-level-id>: <criterion>
areas:                          # Optional*
  <area-name>: <Area>
```

At least one of `factors`, `requirements`, or `areas` must be supplied.
Areas can nest recursively. `ratingScale` exists only on the root model.

### Core concepts

| Concept      | Meaning                                                        |
| ------------ | -------------------------------------------------------------- |
| Model        | The root quality model in a QUALITY.md file.                   |
| Area         | The thing being evaluated.                                     |
| Source       | The material assessed for an area, such as a path or selector. |
| Factor       | A quality dimension that matters for an area.                  |
| Requirement  | A specific quality expectation.                                |
| Assessment   | The means of checking a requirement against an area source.    |
| Finding      | An observation produced by an assessment.                      |
| Rating scale | The ordered model-wide scale used to rate results.             |

## CLI quick reference

| Task                    | Command                          |
| ----------------------- | -------------------------------- |
| Show format spec        | `qualitymd spec`                 |
| Show frontmatter schema | `qualitymd schema`               |
| Create a starter file   | `qualitymd init [path]`          |
| Validate a file         | `qualitymd lint [path]`          |
| Fix lint issues         | `qualitymd lint --fix [path]`    |
| Query model structure   | `qualitymd model tree [path]`    |
| Run an evaluation       | `qualitymd evaluation run`       |
| Show workspace status   | `qualitymd status [path] --json` |
| Show version info       | `qualitymd version --json`       |
| Check for updates       | `qualitymd update --check`       |
| Show command help       | `qualitymd <command> --help`     |

`qualitymd init` also adds a concise `QUALITY.md` pointer to local agent
instruction files by default. Use `qualitymd init --no-agent-instructions` to
skip that side effect.

## Status

The QUALITY.md format, `qualitymd` CLI, and `/quality` skill are early and
under active development. Expect the format and tooling to change as they
mature.

## Contributing

Contributor setup and local tasks live in [`CONTRIBUTING.md`](CONTRIBUTING.md).
