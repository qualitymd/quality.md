# QUALITY.md

**QUALITY.md** is an agent-friendly file format and companion agent skill and
CLI for continuously improving the quality of coding agent and AI assistant projects/harnesses.

## Why QUALITY.md

As software — and the agents that write it — moves faster, quality erodes
quietly through three accumulating debts:

- **Technical debt** — code drifting from where it should be.
- **Cognitive debt** — the mounting burden of understanding complex or
  under-documented systems.
- **Intent debt** — software diverging from what stakeholders actually need.

QUALITY.md makes a team's quality expectations explicit and checkable, so those
gaps stay visible and addressable instead of compounding.

> The three-debt framing draws on Margaret-Anne Storey, *The Triple Debt of
> Software Development* ([arXiv:2603.22106](https://arxiv.org/abs/2603.22106)).

## Install

1. Install the agent skill:

```sh
npx skills add qualitymd/quality.md
```

2. Install the CLI:

```sh
npm install quality.md -g
```

## Usage

Create and check a first model with the CLI:

```sh
qualitymd init
qualitymd lint
```

Expected first result:

```text
Created QUALITY.md

Next: qualitymd lint QUALITY.md
QUALITY.md is valid.
```

Invoke the `/quality` skill to manage quality for your project:

```text
/quality setup                                  Get started working with QUALITY.md
/quality wizard                                 Have your AI assistant/agent help you manage quality
/quality evaluate                               Evaluate the quality of your project
/quality evaluate security                      Evaluate a specific quality factor or characteristic
/quality evaluate payments-api                  Evaluate specific target or project component
/quality evaluate payments-api maintainability  Evaluate a target's specific quality
```

## The Format

A `QUALITY.md` file combines a structured quality model with plain-language
rationale. The structured model names what is being evaluated, which quality
factors matter, what requirements define those factors, and how each requirement
should be assessed. The Markdown body explains the context: what the work is,
what "good" means, where the boundaries are, and why these standards matter.

### Specification

The format is specified in [`SPECIFICATION.md`](SPECIFICATION.md).

### Example QUALITY.md

```markdown
---
title: Support Inbox
ratingScale:
  - level: outstanding
    title: Outstanding
    description: The work clearly exceeds the shared quality bar.
    criterion: "Consistently exceeds the requirement with clear margin."
  - level: target
    title: Target
    description: The work meets the shared quality bar.
    criterion: "Meets the expected quality bar."
  - level: minimum
    title: Minimum
    description: The work is acceptable, but has gaps worth improving.
    criterion: "Meets the lowest acceptable bar, with visible gaps."
  - level: unacceptable
    title: Unacceptable
    description: The work is below the shared quality bar.
    criterion: "Falls below the minimum acceptable bar."
targets:
  triage:
    title: Triage
    source: ./support
    factors:
      responsiveness:
        title: Responsiveness
        description: Customers receive timely, useful attention.
        requirements:
          "urgent messages are visible":
            assessment: >
              New messages are classified so urgent customer-impacting issues
              are separated from routine requests.
      accuracy:
        title: Accuracy
        description: Replies are correct, complete, and grounded in policy.
        requirements:
          "answers cite the current policy":
            assessment: >
              Customer-facing replies use the active support policy and do not
              rely on outdated guidance or unsupported assumptions.
---

# Quality model: Support Inbox

## Overview

This model describes the quality bar for daily support triage. Good support
means urgent issues are easy to see, routine requests still move, and customers
receive answers grounded in the current policy.

## Scope

This model covers message triage and written replies in the support workspace.
It does not cover billing system behavior or product incident response.
```

An agent that reads this file can evaluate support work against the stated
requirements, produce findings, and rate the results against the model's scale.

For a completed evaluation run, the CLI renders a concise summary and the full
report from records supplied by the `/quality` skill:

```text
Wrote quality/evaluations/0001-subject-quality-eval/report-summary.md, quality/evaluations/0001-subject-quality-eval/report.md, and quality/evaluations/0001-subject-quality-eval/report.json
```

A summary excerpt looks like this:

```md
# Quality Evaluation Summary

| Field          | Value  |
| -------------- | ------ |
| Overall rating | Target |

## Top Issues

None recorded.
```

## The Specification

The full format specification lives at [`SPECIFICATION.md`](SPECIFICATION.md).
What follows is a condensed reference.

### File Structure

A `QUALITY.md` file has two layers:

1. **YAML frontmatter** - the structured quality model.
2. **Markdown body** - the context, rationale, scope, needs, risks, and known
   gaps that help people and agents interpret the model.

The document begins with the YAML frontmatter. The Markdown body can be empty,
but it is where the model explains itself.

### Model Schema

The root model is a target plus a model-wide `ratingScale`.

```yaml
title: <string>                 # Required
description: <string>           # Optional
ratingScale:                    # Required, ordered best to worst
  - level: <level-name>         # Required, unique within the scale
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
      <requirement-statement>: <Requirement>
requirements:                   # Optional*
  <requirement-statement>:
    assessment: <string>        # Required, exactly one
    factors: [<factor-name>]    # Required for direct target requirements
    ratings:                    # Optional per-level criteria
      <level-name>: <criterion>
targets:                        # Optional*
  <target-name>: <Target>
```

At least one of `factors`, `requirements`, or `targets` must be supplied.
Targets can nest recursively. `ratingScale` exists only on the root model.

### Core Concepts

| Concept      | Meaning                                                         |
| ------------ | --------------------------------------------------------------- |
| Model        | The root quality model in a `QUALITY.md` file.                  |
| Target       | The thing being evaluated.                                      |
| Source       | The material assessed for a target, such as a path or selector. |
| Factor       | A quality dimension that matters for a target.                  |
| Requirement  | A specific quality expectation.                                 |
| Assessment   | The means of checking a requirement against a target source.    |
| Finding      | An observation produced by an assessment.                       |
| Rating Scale | The ordered model-wide scale used to rate results.              |

### Evaluation Semantics

Each requirement has exactly one `assessment`. An evaluator performs that
assessment against the target source, records findings, and rates the
requirement's findings against the model's `ratingScale`. A finding records what
was observed; the rating result records how that observation compares with the
model's scale.

```text
assessment -> findings -> rating result
```

## The CLI

> **The CLI is an early work in progress.** Today the binary ships
> `qualitymd init`, `qualitymd lint`, `qualitymd spec`, `qualitymd status`, and the
> `qualitymd evaluation` run-record surface.

`qualitymd` draws one hard line: the **CLI never asks an AI model to judge your
project.** It creates and checks `QUALITY.md` files, shows what the file covers,
writes evaluation records for the `/quality` skill, renders reports, and can
fail CI when ratings fall below your chosen bar. The judgment work happens in
the skill, not in the CLI.

### Common commands

| Goal                   | Command                          |
| ---------------------- | -------------------------------- |
| Show the format rules  | `qualitymd spec`                 |
| Create a starter file  | `qualitymd init [path]`          |
| Check a file           | `qualitymd lint [path]`          |
| Fix simple lint issues | `qualitymd lint --fix [path]`    |
| Show project status    | `qualitymd status [path] --json` |
| Show version info      | `qualitymd version --json`       |
| Check for updates      | `qualitymd upgrade --check`      |
| Show command help      | `qualitymd <command> --help`     |

Typical local loop:

```sh
qualitymd spec
qualitymd init
qualitymd lint
qualitymd status --json
```

The `/quality` skill uses additional evaluation commands behind the scenes.
The detailed command guide lives in the bundled
[`CLI Quick Reference`](skills/quality/resources/cli-quick-reference.md).

## Conceptual model

The way `QUALITY.md` frames quality is informed by the **ISO/IEC 25000 (SQuaRE)**
family of software-quality standards — particularly ISO/IEC 25010 — and, for the
shape of a well-formed requirement, **ISO/IEC/IEEE 29148**. We acknowledge these
as the conceptual lineage, not a conformance target: `QUALITY.md` borrows their
ideas and vocabulary where they help and diverges where they don't (it uses
*Factors* where ISO says *characteristics*), optimizing first for a practical,
readable format.

## Status

The `QUALITY.md` format, `qualitymd` CLI, and `/quality` skill are early and under active development. The format is specified in [`SPECIFICATION.md`](SPECIFICATION.md). Expect the format and tooling to change as they mature.

## Contributing

Contributor setup and local tasks live in [`CONTRIBUTING.md`](CONTRIBUTING.md).
The release runbook lives in
[`docs/guides/cut-a-release.md`](docs/guides/cut-a-release.md).
