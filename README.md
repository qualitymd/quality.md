# quality.md

**`QUALITY.md`** is a plain-text *quality model*: one checked-in file that
declares what *good* means for a software system — its quality requirements —
and how each is checked. Where a test suite pins down behavior, a `QUALITY.md`
pins down quality — reliability, security, maintainability, and the rest —
stated explicitly instead of living in reviewers' heads.

The file has two parts: **YAML frontmatter** holding the structured model —
*factors* and the *requirements* under them — and a **Markdown body** documenting
what the system is, what *good* means for it, and why those are the right
requirements. It is written for whoever decides quality: **authors** declare the
model, **coding agents** read it to build and evaluate against, and **CI** gates
on the result.

> 🚧 **Alpha.** The `QUALITY.md` format and the `qualitymd` CLI are early and
> under active development. The format is specified in
> [`SPECIFICATION.md`](SPECIFICATION.md); the CLI is still a work in progress
> (see [The CLI](#the-cli)). Expect the format to change as it matures.

## What a QUALITY.md looks like

```markdown
---
factors:
  reliability:
    requirements:
      "acknowledged writes are durable":
        target: "./internal/storage"
        prompt: >
          A write is acknowledged to the client only after it is committed
          to durable storage. Failures surface as errors, never as false
          successes.
      "the write path is covered end-to-end":
        prompt: >
          The write path is exercised by automated tests that would fail if a
          write were lost or acknowledged before it was durable.
  security:
    requirements:
      "no secrets are committed":
        prompt: >
          No credentials, API keys, or tokens appear in source, config, or
          fixtures; secrets are loaded from the environment at runtime.
---

# Quality model — Orders API

## Overview

The Orders API is the public HTTP interface customers integrate against,
maintained by the platform team. "Good" here means it never silently loses or
corrupts an order. This model covers the service and its data layer; the
third-party payment provider is a dependency, not part of it.

## Factors

### Reliability

Customers build on our acknowledgements, so a confirmed write must be durable.
When durability and latency conflict, durability wins.

### Security

The API handles customer data, so access is authenticated and least-privilege.
```

The **frontmatter** is the structured model: **factors** (quality attributes
like *reliability* and *security*) and the **requirements** under them. Each
requirement carries exactly one assessment — a **`prompt`**, judged against the
requirement's intent — and an optional **`target`** narrowing what it applies to
(here, the storage layer). The **body** holds the reasoning the frontmatter
cannot: what the system is, what *good* means for it, and why these are the right
requirements — the same context a `prompt` is judged against.

A coding agent reads this file to evaluate the Orders API against it — judging
each `prompt` against the body, then reporting where the subject falls short. The
`qualitymd` CLI records and rolls up those verdicts deterministically; the
judging is the agent's part.

## The CLI

> **The CLI is an early work in progress.** Today the binary ships a single
> placeholder `check` command that predates the current spec. The surface below
> is specified under [`specs/`](specs/) but not yet built.

`qualitymd` draws one hard line: the **CLI is deterministic and never calls a
model** — it scaffolds and validates a `QUALITY.md`, resolves targets, records
and rolls up results, and gates CI — while **skills carry the judgment**, driving
the evaluation loop and judging each `prompt` against the model.

The deterministic CLI:

- **`qualitymd init`** *(planned)* — scaffold a starter `QUALITY.md` to fill in.
- **`qualitymd lint`** *(planned)* — validate a file's structure, fast and
  deterministic, exiting non-zero on errors so it drops into CI.
- **`qualitymd model` / `evaluation` / `result`** *(planned)* — inspect the
  model, manage a per-target evaluation run, record verdicts, and gate on the
  outcome with `evaluation report --fail-on`.

The deep, judgment-based evaluation of a subject against its model is carried by
**skills** that orchestrate those resources — not by a CLI command.

The only command in the shipping binary today is the placeholder
**`qualitymd check`**; the surface above fails with "unknown command" until it
lands. The full CLI and skill surface is specified under [`specs/`](specs/) — see
[`specs/cli.md`](specs/cli.md) and [`specs/skills.md`](specs/skills.md).

## Install

> **Status.** The format spec and CLI design are settled — see
> [`SPECIFICATION.md`](SPECIFICATION.md) and [`specs/`](specs/) — but
> implementation is in progress. Of the documented surface, only the placeholder
> **`check`** is currently built; **`init`**, **`lint`**, and the
> **`model`/`evaluation`/`result`** resources are planned. Don't expect the
> planned commands to run yet.

`qualitymd` has no tagged release yet; build the current binary from source with
Go 1.26+:

```sh
go install github.com/qualitymd/quality.md/cmd/qualitymd@latest
```

Pre-built binaries via npm (`npx quality.md`) and Homebrew
(`brew install qualitymd/tap/qualitymd`) arrive with the first tagged release.

## Specification

The `QUALITY.md` format is specified in [`SPECIFICATION.md`](SPECIFICATION.md),
and the `qualitymd` CLI under [`specs/`](specs/). These are the source of truth
for the format and the tool.

## Conceptual model

The way `QUALITY.md` frames quality is informed by the **ISO/IEC 25000 (SQuaRE)**
family of software-quality standards — particularly ISO/IEC 25010 — and, for the
shape of a well-formed requirement, **ISO/IEC/IEEE 29148**. We acknowledge these
as the conceptual lineage, not a conformance target: `QUALITY.md` borrows their
ideas and vocabulary where they help and diverges where they don't (it uses
*Factors* and *Subfactors* where ISO says *characteristics*), optimizing first
for a practical, readable format.

## Contributing

Development setup, tasks, and the release process live in
[`CONTRIBUTING.md`](CONTRIBUTING.md).
