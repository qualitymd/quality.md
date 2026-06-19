# QUALITY.md

**QUALITY.md** is an agent-friendly file format and companion agent skill and
CLI for continuously improving the quality of coding agent and AI assistant projects/harnesses.

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
    source: ./support
    factors:
      responsiveness:
        description: Customers receive timely, useful attention.
        requirements:
          "urgent messages are visible":
            assessment: >
              New messages are classified so urgent customer-impacting issues
              are separated from routine requests.
      accuracy:
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
title: <string>                 # Recommended
description: <string>           # Optional
ratingScale:                    # Required, ordered best to worst
  - level: <level-name>         # Required, unique within the scale
    title: <string>             # Optional
    description: <string>       # Recommended
    criterion: <string>         # Required
source: <string>                # Optional
factors:                        # Optional*
  <factor-name>:
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

`qualitymd` draws one hard line: the **CLI is deterministic and never calls a
model.** It scaffolds and validates a `QUALITY.md`, resolves target nodes and
their `source` manifests, records evaluation artifacts, renders reports, and
gates CI. The deep, judgment-based evaluation of a subject against its model is
carried by **skills**, not by any CLI command.

The deterministic surface:

- **`qualitymd init`** — scaffold a starter `QUALITY.md` to fill in.
- **`qualitymd lint`** — validate a file's structure, fast and deterministic,
  exiting non-zero on errors so it drops into CI.
- **`qualitymd spec`** — emit the bundled `QUALITY.md` format specification.
- **`qualitymd status`** — emit a deterministic project-state snapshot for
  routing, automation, and agent use.
- **`qualitymd evaluation create-run`** — create and number an evaluation run
  folder.
- **`qualitymd evaluation add-record`** — write assessment, analysis, and
  recommendation records from judgment payloads.
- **`qualitymd evaluation set-planned-coverage`** — write optional planned
  assessment and analysis coverage for resume diagnostics.
- **`qualitymd evaluation show-status`** — inspect whether a run is ready to
  render.
- **`qualitymd evaluation build-report`** — derive `report.md` / `report.json`
  and optionally gate with `--fail-at-or-below`.

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
