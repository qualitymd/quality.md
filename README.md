# QUALITY.md

**QUALITY.md** is an [open format](./SPECIFICATION.md) for modeling quality:
what matters most, why, and how. Use it with the `/quality` agent skill to
continuously improve AI assistant and coding agent projects.

A QUALITY.md file is a Markdown file with a quality model and supporting
context. The `/quality` skill helps set up the file, evaluate quality, and
evolve the model as you learn. The `qualitymd` CLI provides support tooling for
validating QUALITY.md files, managing quality evaluations, and maintaining a
QUALITY.md workspace.

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
/quality wizard                                 Have your AI assistant/agent help you manage quality
/quality evaluate                               Evaluate the quality of your project
/quality evaluate security                      Evaluate a specific quality factor or characteristic
/quality evaluate payments-api                  Evaluate a specific area or project component
/quality evaluate payments-api maintainability  Evaluate an area's specific quality
```

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

## Example QUALITY.md

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
areas:
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
    factors: [<factor-name>]    # Required for direct area requirements
    ratings:                    # Optional per-level criteria
      <level-name>: <criterion>
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

| Task                  | Command                          |
| --------------------- | -------------------------------- |
| Show format spec      | `qualitymd spec`                 |
| Create a starter file | `qualitymd init [path]`          |
| Validate a file       | `qualitymd lint [path]`          |
| Fix lint issues       | `qualitymd lint --fix [path]`    |
| Show project status   | `qualitymd status [path] --json` |
| Show version info     | `qualitymd version --json`       |
| Check for updates     | `qualitymd update --check`       |
| Show command help     | `qualitymd <command> --help`     |

## Status

The QUALITY.md format, `qualitymd` CLI, and `/quality` skill are early and
under active development. Expect the format and tooling to change as they
mature.

## Contributing

Contributor setup and local tasks live in [`CONTRIBUTING.md`](CONTRIBUTING.md).
