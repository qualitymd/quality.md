---
type: How-to Guide
title: Use the /quality skill
description: Install the /quality skill, verify the qualitymd CLI prerequisite, and run setup, wizard, and evaluation modes.
tags: [skill, quality, evaluation]
timestamp: 2026-06-17T00:00:00Z
---

# Use the /quality skill

Use the `/quality` skill when you want an agent to set up, evaluate, or improve
the quality model for a project/entity or one of its components/targets.

## Install the skill

```sh
npx skills add qualitymd/quality.md
```

For local dogfooding from this repository, use the installer's local-path form
when supported:

```sh
npx skills add .
```

Restart the target agent if it discovers skills only at session startup.

## Verify the CLI prerequisite

The skill drives the deterministic `qualitymd` CLI for setup, linting, format
grounding, bundled model access, and evaluation record/report mechanics:

```sh
qualitymd --version
qualitymd models list --json
qualitymd evaluation create-run --help
qualitymd evaluation add-record --help
qualitymd evaluation show-status --help
qualitymd evaluation build-report --help
```

If the CLI is missing or stale, build the current binary from source until a
tagged release channel exists:

```sh
go install github.com/qualitymd/quality.md/cmd/qualitymd@latest
```

## Run the skill

In the repository you want to evaluate:

```text
/quality setup
/quality wizard
/quality evaluate
/quality evaluate model
```

`setup` creates and lints a skeleton `QUALITY.md`. `wizard` inspects the model
state and suggests concrete next actions. `evaluate` writes a numbered evaluation
run under `quality/evaluations/` by default.

## Configure the evaluation directory

Create `.quality/config.yaml` to choose a different parent directory for
numbered evaluation runs:

```yaml
evaluationDir: quality/evaluations
```

The path must be repository-relative and must not escape the repository.
