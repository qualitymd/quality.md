---
type: How-to Guide
title: Use the /quality skill
description: Install the /quality skill, maintain the skill/CLI pair, and run setup, wizard, evaluation, and improvement modes.
tags: [skill, quality, evaluation]
timestamp: 2026-06-17T00:00:00Z
---

# Use the /quality skill

Use the `/quality` skill when you want an agent to set up, evaluate, improve, or
maintain the quality workflow for a project/entity or one of its
components/targets.

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

## Upgrade an existing install

For an existing setup, use the skill-orchestrated upgrade flow:

```text
/quality upgrade
```

The upgrade mode checks the installed `/quality` skill metadata, verifies the
visible `qualitymd` CLI, plans any skill and CLI updates, asks before applying
changes, and reports whether the agent session must be restarted or reloaded.

If `/quality upgrade` is unavailable, reinstall the skill and check the CLI
manually:

```sh
npx skills add qualitymd/quality.md
qualitymd upgrade --check
```

## Verify the CLI prerequisite

The skill drives the deterministic `qualitymd` CLI for setup, linting, format
grounding, and evaluation record/report mechanics. Released skill installs use
the CLI SemVer range declared in `skills/quality/SKILL.md`
`metadata.requires-qualitymd-cli`; see
[Versioning](../reference/versioning.md) for the compatibility policy.

```sh
qualitymd --version
qualitymd version --json
qualitymd upgrade --check
qualitymd spec
qualitymd lint --help
qualitymd init --help
qualitymd evaluation create --help
qualitymd evaluation list --help
qualitymd evaluation status --help
qualitymd evaluation assessment-result --help
qualitymd evaluation analysis --help
qualitymd evaluation recommendation --help
qualitymd evaluation report --help
```

If the CLI is missing or stale, use the recommended action from
`qualitymd upgrade --check`, or install through the GitHub-hosted managed
installer:

```sh
curl -fsSL https://raw.githubusercontent.com/qualitymd/quality.md/main/install/install.sh | sh
```

## Run the skill

In the repository you want to evaluate:

```text
/quality setup
/quality wizard
/quality evaluate
/quality upgrade
```

`setup` creates and lints a skeleton `QUALITY.md`, then uses the bundled
authoring guide and getting-started guide to help populate the first useful
model. `wizard` inspects the model state and suggests concrete next actions.
`evaluate` writes a numbered evaluation run under `quality/evaluations/` by
default. `upgrade` plans and orchestrates paired skill/CLI maintenance without
running a quality evaluation.

## Configure the evaluation directory

Create `.quality/config.yaml` to choose a different parent directory for
numbered evaluation runs:

```yaml
evaluationDir: quality/evaluations
```

The path must be repository-relative and must not escape the repository.
