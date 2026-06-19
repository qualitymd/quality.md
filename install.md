# Install QUALITY.md

Install the `/quality` skill first; the `qualitymd` CLI is its deterministic
prerequisite for setup, linting, format grounding, and evaluation artifact
mechanics.

## 1. Install the skill

```sh
npx skills add qualitymd/quality.md
```

For local development in this repository, install from the working tree if your
Agent Skills installer supports local paths:

```sh
npx skills add .
```

Restart the target agent session if that agent only discovers skills at startup.

## 2. Verify or install the CLI

Check whether `qualitymd` is available:

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

Released skill installs declare the `qualitymd` CLI SemVer range they support;
see [Versioning](docs/reference/versioning.md) for the compatibility policy.

If the CLI is missing or stale, install a pre-built binary — via npm:

```sh
npm install -g quality.md
```

or Homebrew:

```sh
brew install qualitymd/tap/qualitymd
```

Or build the current CLI from source:

```sh
go install github.com/qualitymd/quality.md/cmd/qualitymd@latest
```

## 3. Bootstrap a project

In the repository to evaluate, ask the installed skill to set up or guide you:

```text
/quality setup
/quality wizard
```

`setup` creates and lints a skeleton `QUALITY.md` through the CLI. `wizard`
checks the model, identifies available targets/factors, and suggests concrete
next actions such as `/quality evaluate` or scoped evaluations.

## 4. Optional config

Create `.quality/config.yaml` to move evaluation run folders away from the
default `quality/evaluations/` parent:

```yaml
evaluationDir: quality/evaluations
```

The path must be repository-relative and must not escape the repository.
