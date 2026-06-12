# quality.md

A CLI for evaluating **`QUALITY.md`** specifications — a plain-text representation
of a quality model that declares the quality requirements for a software system
or component, and scores them.

A `QUALITY.md` file has two parts: **YAML frontmatter** holding the structured
quality model, and a **Markdown body** documenting it.

## Install

```sh
# npm / npx (no toolchain required)
npx quality.md check

# Homebrew
brew install qualitymd/tap/qualitymd

# Go
go install github.com/qualitymd/quality.md/cmd/qualitymd@latest
```

## Usage

```sh
qualitymd check            # evaluate ./QUALITY.md
qualitymd check path.md    # evaluate a specific file
qualitymd --help
```

`check` evaluates every requirement and prints a grouped pass/fail report. It
exits non-zero if any requirement fails, so it drops straight into CI.

## Spec format

> 🚧 The `QUALITY.md` format is still in flux. See [`docs/spec.md`](docs/spec.md)
> for the current draft schema.

## Contributing

Development setup, tasks, and the release process live in
[`CONTRIBUTING.md`](CONTRIBUTING.md).
