---
name: quality
description: "Setup or work with QUALITY.md files or the qualitymd CLI; model, evaluate, or improve project or harness quality, get wizard quality advice, anything concerning quality factors/attributes/characteristics relevant to project context"
---

## Purpose

Drive quality management work for a project/entity through `QUALITY.md` and the
`qualitymd` CLI. Keep judgment in the skill and mechanical work in the CLI.

## Prerequisates

- Read [`resources/SPECIFICATION.md`](resources/SPECIFICATION.md)
- Read [`resources/quality-md-guide.md`](resources/quality-md-guide.md) when
  creating, populating, reviewing, or improving a `QUALITY.md` file.

## CLI Reference

```sh
qualitymd --version
qualitymd spec
qualitymd lint --help
qualitymd init --help
qualitymd evaluation create-run --help
qualitymd evaluation add-record --help
qualitymd evaluation set-planned-coverage --help
qualitymd evaluation show-status --help
qualitymd evaluation build-report --help
```

Accept a local development build when those commands are present. If the CLI is
missing or stale, stop and help the user install or upgrade it before continuing.

## Arguments

Parse the user's request from free-form arguments:

- Mode: `wizard` by default when direction is unclear; otherwise `evaluate`,
  `improve`, `setup`, or `wizard`.
- Target file: explicit path if supplied; otherwise `QUALITY.md` in the current
  working directory. Do not walk parent directories.
- Scope: whole model by default, or a named target/factor. Use explicit
  `target`/`factor` words to disambiguate.
- Effort: `standard` by default; `quick` covers apex and high-risk in-scope
  requirements only; `standard` covers every in-scope requirement with targeted
  evidence; `deep` covers every in-scope requirement against a full source read
  and adversarial verification of rating-binding findings.

When a bare request is ambiguous, run `wizard`: inspect state, summarize the
concrete runnable options, and ask the user which action to take.

## Mode Dispatch

After resolving the mode, read the matching mode file before acting:

- `setup` → [`modes/setup.md`](modes/setup.md)
- `wizard` → [`modes/wizard.md`](modes/wizard.md)
- `evaluate` → [`modes/evaluate.md`](modes/evaluate.md)
- `improve` → [`modes/improve.md`](modes/improve.md)

## Config

Read `.quality/config.yaml` from the repository root when present. Supported now:

```yaml
evaluationDir: quality/evaluations
```

Rules:

- Default to `quality/evaluations/` when the file or key is absent.
- Treat `evaluationDir` as the parent directory for numbered run folders.
- Require a repository-relative normalized path.
- Reject absolute paths and paths that escape the repository.
- Warn and ignore unknown keys.

## Artifact Contract

The evaluation record layout and field contract lives in
[`specs/evaluation-records.md`](../../specs/evaluation-records.md). Treat that
spec, plus the `qualitymd evaluation ...` command help, as the source of truth.
Do not restate the schema or folder layout in this prompt.
