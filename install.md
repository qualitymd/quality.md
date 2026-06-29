# Install QUALITY.md

Install the `/quality` agent skill first; it is the primary agentic experience
for working with QUALITY.md. The `qualitymd` CLI is its deterministic
prerequisite and support tooling for setup, linting, format grounding, and
evaluation artifact mechanics.

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

## 2. Update an existing install

For an existing `/quality` setup, prefer the skill-orchestrated update flow:

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

## 3. Verify or install the CLI

Check whether `qualitymd` is available:

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

Released skill installs declare the `qualitymd` CLI SemVer range they support;
see [Versioning](docs/reference/versioning.md) for the compatibility policy.

If the CLI is missing or stale, prefer the GitHub-hosted managed installer.

macOS/Linux:

```sh
curl -fsSL https://raw.githubusercontent.com/qualitymd/quality.md/main/install/install.sh | sh
```

Windows PowerShell:

```powershell
iwr https://raw.githubusercontent.com/qualitymd/quality.md/main/install/install.ps1 -UseB | iex
```

For agents and CI, use non-interactive and pinned forms:

```sh
curl -fsSL https://raw.githubusercontent.com/qualitymd/quality.md/main/install/install.sh | QUALITYMD_NO_INPUT=1 QUALITYMD_VERSION=v0.5.1 sh
```

```powershell
.\install.ps1 -NonInteractive -Version v0.5.1
```

When the install directory is not already on your `PATH`, the shell installer
prints the exact `export PATH=...` line to add (it never edits your shell
profiles), and the PowerShell installer updates your per-user `PATH` and asks you
to open a new terminal. The `--non-interactive` / `-NonInteractive` flags (and
`QUALITYMD_NO_INPUT=1`) suppress that human-oriented guidance for CI and agent
runs; they do not change what is installed. To pass flags through a piped shell
install, use `sh -s --`:

```sh
curl -fsSL https://raw.githubusercontent.com/qualitymd/quality.md/main/install/install.sh | sh -s -- --version v0.5.1
```

Alternative channels remain supported. Install via npm:

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

After installation, `qualitymd update --check` reports the detected install
method, latest known version, release readiness, and the recommended update
action. `qualitymd update` applies by default through supported owner channels
such as managed standalone, npm, and Homebrew; unknown, source, and archive
installs receive manual guidance. Set `QUALITYMD_NO_UPDATE_CHECK=1` to disable
explicit update checks, the ambient cached update notice, and background cache
refresh.

## 4. Bootstrap a project

In the repository to evaluate, ask the installed skill to set up or guide you:

```text
/quality setup
/quality
```

`setup` inspects available context, asks a few setup questions with recommended
defaults, writes only `QUALITY.md`, and validates it through the CLI. Bare
`/quality` gives read-only guidance on the next public workflow, such as
`/quality evaluate` or scoped evaluations.

If you use `qualitymd init` directly, it creates a starter `QUALITY.md` and adds
a concise pointer to local agent instruction files by default. Pass
`--no-agent-instructions` to skip that pointer, then invoke `/quality setup` to
tailor the scaffold to the project.

## 5. Optional config

Create `.quality/config.yaml` in the selected `QUALITY.md` workspace to move
evaluation run folders away from the default `.quality/evaluations/` parent. If
your config file lives elsewhere, add root `config: <path>` frontmatter to the
selected `QUALITY.md` to point qualitymd to it.

```yaml
evaluationDir: tmp/evals
```

The path is relative to the selected `QUALITY.md` workspace and must not escape
the repository.
