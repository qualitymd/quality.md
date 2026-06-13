# quality.md

A CLI for evaluating **`QUALITY.md`** specifications — a plain-text representation
of a quality model that declares the quality requirements for a software system
or component, and scores them.

A `QUALITY.md` file has two parts: **YAML frontmatter** holding the structured
quality model, and a **Markdown body** documenting it.

## Install

```sh
# npm / npx (no toolchain required)
npx quality.md lint

# Homebrew
brew install qualitymd/tap/qualitymd

# Go
go install github.com/qualitymd/quality.md/cmd/qualitymd@latest
```

## Usage

```sh
qualitymd init             # scaffold a starter ./QUALITY.md to fill in
qualitymd lint             # validate the file's structure (fast, deterministic)
qualitymd evaluate         # deep evaluation of the subject against the model
qualitymd --help
```

`lint` checks that the file is a well-formed `QUALITY.md` — it parses, conforms
to the schema, and its references resolve — and exits non-zero on errors, so it
drops straight into CI. `evaluate` runs the deep, judgment-based audit of the
subject against the model's requirements.

## Spec format

> 🚧 The `QUALITY.md` format is still in flux. See [`SPECIFICATION.md`](SPECIFICATION.md)
> for the current draft schema.

## Conceptual model

The way `QUALITY.md` frames quality is informed by the **ISO/IEC 25000 (SQuaRE)**
family of software-quality standards — particularly ISO/IEC 25010 — and, for the
shape of a well-formed requirement, **ISO/IEC/IEEE 29148**. We acknowledge these
standards as the conceptual lineage for the format. `QUALITY.md` is not intended
to strictly conform to them, however: it borrows ideas and vocabulary where they
help, diverges where they don't (for example, it uses *Factors* and *Subfactors*
where ISO says *characteristics*), and optimizes first for being a practical,
readable format.

## Contributing

Development setup, tasks, and the release process live in
[`CONTRIBUTING.md`](CONTRIBUTING.md).
