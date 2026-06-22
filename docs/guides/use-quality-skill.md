---
type: How-to Guide
title: Use the /quality skill
description: Install the /quality skill, maintain the skill/CLI pair, run setup/wizard/evaluation modes, and follow up on recommendations.
tags: [skill, quality, evaluation]
timestamp: 2026-06-21T00:00:00Z
---

# Use the /quality skill

Use the `/quality` skill when you want an agent to set up, evaluate, maintain,
or follow up on quality recommendations for a project/entity or one of its
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

## Update an Existing Install

For an existing setup, use the skill-orchestrated update flow:

```text
/quality update
```

The update mode checks the installed `/quality` skill metadata, verifies the
visible `qualitymd` CLI, plans any skill and CLI updates, asks before applying
changes, and reports whether the agent session must be restarted or reloaded.

If `/quality update` is unavailable, reinstall the skill and check the CLI
manually:

```sh
npx skills add qualitymd/quality.md
qualitymd update --check
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
qualitymd update --check
qualitymd spec
qualitymd lint --help
qualitymd init --help
qualitymd evaluation create --help
qualitymd evaluation list --help
qualitymd evaluation status --help
qualitymd evaluation assessment --help
qualitymd evaluation analysis --help
qualitymd evaluation recommendation --help
qualitymd evaluation report --help
```

If the CLI is missing or stale, use the recommended action from
`qualitymd update --check`, or install through the GitHub-hosted managed
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
/quality update
```

`setup` creates and lints a skeleton `QUALITY.md`, then uses the bundled
authoring guide and getting-started guide to help populate the first useful
model. `wizard` inspects the model state and suggests concrete next actions.
`evaluate` writes a numbered evaluation run under `quality/evaluations/` by
default, including formal assessment records, generated reports, and a
process-only `debug-log.md` for notable evaluation orchestration events.
`update` plans and orchestrates paired skill/CLI maintenance without running a
quality evaluation.

Recommendations are produced by `evaluate`. After an evaluation, ask the agent
to follow up on a recommendation when you want to apply a confirmed option now
or hand it off to an issue tracker. The agent prepares issue-ready text, and it
creates an external issue only when suitable issue-tracker tooling is available
and you explicitly confirm creation.

`setup` and confirmed model-authoring or recommendation-apply workflows maintain
a **quality log** — dated entries under `quality/log/` that record meaningful,
evidence-linked changes to the model. `setup` seeds an inaugural entry; later
confirmed model changes add one entry per coherent model change. The log
preserves *why* the model changed where `git log` only shows the diff; it is
curated, not exhaustive. Issue-tracker handoff alone does not write the quality
log.

## Configure the evaluation directory

Create `.quality/config.yaml` to choose a different parent directory for
numbered evaluation runs:

```yaml
evaluationDir: quality/evaluations
```

The path must be repository-relative and must not escape the repository.
