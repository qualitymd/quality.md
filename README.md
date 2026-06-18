# QUALITY.md

**`QUALITY.md`** is a plain-text *quality model*: one checked-in file that
declares what *good* means for a software system — its quality requirements —
and how each is checked. Where a test suite pins down behavior, a `QUALITY.md`
pins down quality — reliability, security, maintainability, and the rest —
stated explicitly instead of living in reviewers' heads.

The file has two parts: **YAML frontmatter** holding the structured model — a
recursive tree of *targets*, scoped *factors*, and *requirements* — and a
**Markdown body** documenting what the system is, what *good* means for it, and
why those are the right requirements. It is written for whoever decides quality:
**authors** declare the model, **coding agents** read it to build and evaluate
against, and **CI** gates on the result.

> 🚧 **Alpha.** The `QUALITY.md` format, `qualitymd` CLI, and `/quality` skill
> are early and under active development. The format is specified in
> [`SPECIFICATION.md`](SPECIFICATION.md). Expect the format and tooling to change
> as they mature.

## What a QUALITY.md looks like

```markdown
---
title: Orders API
ratingScale:
  - level: outstanding
    title: Outstanding
    criterion: "Exceeds the requirement; satisfies it with margin to spare."
  - level: target
    title: Target
    criterion: "Satisfies the requirement."
  - level: minimum
    title: Minimum
    criterion: "Falls short of the goal but holds the acceptable floor."
  - level: unacceptable
    title: Unacceptable
    criterion: "Falls below the acceptable floor."
targets:
  api:
    source: ./internal/api
    requirements:
      "accepted orders are durable":
        assessment: >
          A write is acknowledged to the client only after it is committed to
          durable storage. Failures surface as errors, never as false successes.
    factors:
      reliability:
        description: The API behaves predictably under ordinary and failure conditions.
        requirements:
          "the write path is covered end-to-end":
            assessment: >
              The write path is exercised by automated tests that would fail if a
              write were lost or acknowledged before it was durable.
      security:
        description: Customer data and privileged operations are protected.
        requirements:
          "no secrets are committed":
            assessment: >
              No credentials, API keys, or tokens appear in source, config, or
              fixtures; secrets are loaded from the environment at runtime.
  docs:
    source: ./docs
    requirements:
      "integration docs describe the order lifecycle":
        assessment: >
          The documentation explains how an order moves from request to durable
          acknowledgement, including failure responses and retry guidance.
---

# Quality model — Orders API

## Overview

The Orders API is the public HTTP interface customers integrate against,
maintained by the platform team. "Good" here means it never silently loses or
corrupts an order. This model covers the service and its data layer; the
third-party payment provider is a dependency, not part of it.

## Targets and factors

### api

The API target covers the HTTP boundary and storage path.

#### Reliability

Customers build on our acknowledgements, so a confirmed write must be durable.
When durability and latency conflict, durability wins.

#### Security

The API handles customer data, so access is authenticated and least-privilege.
```

The **frontmatter** is the structured model: **targets** (things evaluated),
**factors** (quality lenses scoped to the target where they are declared), and
**requirements** under either a target or a factor. Each requirement carries one
**`assessment`**; the evaluator turns that assessment into a finding, then rates
the finding against the scale criteria. A target's **`source`** identifies the
material assessed, and child `targets:` decompose or narrow the subject. The
**body** holds the reasoning the frontmatter cannot: what the system is, what
*good* means for it, and why these are the right requirements.

A coding agent with the `/quality` skill reads this file to evaluate the Orders
API against it — performing each `assessment` against the target source, then
reporting where the subject falls short. The `qualitymd` CLI owns deterministic
mechanical steps such as scaffolding, linting, spec grounding, and bundled model
access; judging is the skill's part.

## Skill-first onboarding

Install the skill first, then make sure the CLI prerequisite is available:

```sh
npx skills add qualitymd/quality.md
qualitymd --version
```

Then use the skill in your agent:

```text
/quality setup
/quality wizard
/quality evaluate
/quality evaluate model
```

`setup` and `wizard` verify that the `qualitymd` CLI is installed and exposes the
commands the skill depends on. If it is missing or stale, they stop and help you
install or upgrade it before continuing. See [`install.md`](install.md) for the
full bootstrap flow.

## The CLI

> **The CLI is an early work in progress.** Today the binary ships
> `qualitymd init`, `qualitymd lint`, `qualitymd models`, and `qualitymd spec`.
> The evaluation record/gate surface is planned but not yet built.

`qualitymd` draws one hard line: the **CLI is deterministic and never calls a
model** — it scaffolds and validates a `QUALITY.md`, resolves target nodes and
their `source` manifests, records and rolls up results, and gates CI — while
**skills carry the judgment**, driving the evaluation loop and performing each
`assessment` against the model.

The deterministic CLI:

- **`qualitymd init`** — scaffold a starter `QUALITY.md` to fill in.
- **`qualitymd lint`** — validate a file's structure, fast and deterministic,
  exiting non-zero on errors so it drops into CI.
- **`qualitymd models`** — list and view bundled `QUALITY.md` models, including
  the quality meta-model used to evaluate a `QUALITY.md` itself.
- **`qualitymd spec`** — emit the bundled `QUALITY.md` format specification.
- **`qualitymd evaluation` / `result`** *(planned)* — manage a per-target
  evaluation run, record verdicts, and gate on the outcome with
  `evaluation report --fail-on`.

The deep, judgment-based evaluation of a subject against its model is carried by
**skills** that orchestrate those resources — not by a CLI command.

The planned commands above other than `init` and `lint` fail with "unknown
command" until they land.

## Install

> **Status.** The format spec is settled — see
> [`SPECIFICATION.md`](SPECIFICATION.md) — but implementation is in progress.
> Of the documented CLI surface, **`init`**, **`lint`**, **`models`**, and
> **`spec`** are currently built; the **`evaluation`/`result`** resources are
> planned.

Install the `/quality` skill with Agent Skills tooling:

```sh
npx skills add qualitymd/quality.md
```

`qualitymd` has no tagged release yet; build the current binary from source with
Go 1.26+:

```sh
go install github.com/qualitymd/quality.md/cmd/qualitymd@latest
```

Pre-built binaries via npm (`npx quality.md`) and Homebrew
(`brew install qualitymd/tap/qualitymd`) arrive with the first tagged release.

## Specification

The `QUALITY.md` format is specified in [`SPECIFICATION.md`](SPECIFICATION.md),
the source of truth for the format.

## Conceptual model

The way `QUALITY.md` frames quality is informed by the **ISO/IEC 25000 (SQuaRE)**
family of software-quality standards — particularly ISO/IEC 25010 — and, for the
shape of a well-formed requirement, **ISO/IEC/IEEE 29148**. We acknowledge these
as the conceptual lineage, not a conformance target: `QUALITY.md` borrows their
ideas and vocabulary where they help and diverges where they don't (it uses
*Factors* where ISO says *characteristics*), optimizing first for a practical,
readable format.

## Contributing

Development setup, tasks, and the release process live in
[`CONTRIBUTING.md`](CONTRIBUTING.md).
