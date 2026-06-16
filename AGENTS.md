# quality.md

This is the public repository for the **QUALITY.md** project: the `qualitymd`
CLI, the format spec, and the docs shipped to users.

`QUALITY.md` is a plain-text format for a *quality model* — a file that declares
the quality requirements for a software system or component and scores them. Each
file pairs **YAML frontmatter** (the recursive target-tree quality model) with a
**Markdown body** (its documentation). The specified CLI is deterministic: it
scaffolds, lints, resolves, records, rolls up, and reports; skills perform the
judgment and record rating results through the CLI.

See [`README.md`](README.md) for the overview, [`CONTRIBUTING.md`](CONTRIBUTING.md)
for dev setup and the release process, and [`SPECIFICATION.md`](SPECIFICATION.md) for
the draft format spec.

## Conventions

### Referencing ISO standards

The conceptual model behind `QUALITY.md` is informed by ISO/IEC 25000 (SQuaRE)
and ISO/IEC/IEEE 29148, acknowledged in the README. That lineage is deliberately
kept in the background.

**Do not cite or refer to specific ISO standards in public source code or
artifacts** (READMEs, user-facing docs, CLI output, code comments, templates,
schemas) unless the user explicitly requests it, or surfacing the standard is the
specific purpose of that file. Prefer QUALITY.md's own vocabulary (Targets,
Factors, Requirements) over ISO terminology.

ISO standards **may** be referenced in [`SPECIFICATION.md`](SPECIFICATION.md),
where the provenance of a design decision is relevant.

### Agent guidance files

- `CLAUDE.md` and `GEMINI.md` are symlinks to this file — edit `AGENTS.md` only.
  (Both symlinks are gitignored.)
