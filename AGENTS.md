# quality.md

This is the public repository for the **QUALITY.md** project: the `qualitymd`
CLI, the format spec, and the docs shipped to users.

A `QUALITY.md` file pairs **YAML frontmatter** (the recursive target-tree quality
model) with a **Markdown body** (its documentation). The `qualitymd` CLI is
deterministic and mechanical; skills perform the judgment and record results
through it.

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

### Open Knowledge Format (OKF) bundles

Several directories are authored as **OKF** bundles, each registering its concept
types in a root `schema.md`:

| Folder     | What it holds                                                |
| ---------- | ------------------------------------------------------------ |
| `specs/`   | Specifications for the deterministic `qualitymd` surface.    |
| `docs/`    | Project documentation, organized by the four Diátaxis modes. |
| `changes/` | Incremental work items — a spec and design doc per change.   |

### Guides

Before working in this repo, consult the relevant how-to guide in
[`docs/guides/`](docs/guides/index.md). Each is task-oriented:

| When you are…                                   | Read                                                              |
| ----------------------------------------------- | ----------------------------------------------------------------- |
| Proposing or tracking a unit of work            | [Proposing a change](docs/guides/propose-a-change.md)             |
| Writing a functional spec (the `specs/` bundle) | [Writing functional specs](docs/guides/write-functional-specs.md) |
| Writing a design doc                            | [Writing design docs](docs/guides/write-design-docs.md)           |
| Reading or editing any OKF bundle               | [Working with OKF](docs/guides/work-with-okf.md)                  |

### Agent guidance files

- `CLAUDE.md` and `GEMINI.md` are symlinks to this file — edit `AGENTS.md` only.
  (Both symlinks are gitignored.)
