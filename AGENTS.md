# QUALITY.md

Public repo for **QUALITY.md**: `qualitymd` CLI, format spec, user docs.

`QUALITY.md` = YAML frontmatter target tree + Markdown docs. CLI is
deterministic/mechanical; skills judge and record through it.

See [`README.md`](README.md), [`CONTRIBUTING.md`](CONTRIBUTING.md),
[`SPECIFICATION.md`](SPECIFICATION.md).

## Conventions

### Instruction style

Keep this file extremely concise. Brevity over grammar.

### Referencing ISO standards

- Keep ISO lineage background.
- Do not cite specific ISO standards in public code/artifacts unless requested
  or the file's purpose.
- Use QUALITY.md terms: Targets, Factors, Requirements.
- [`SPECIFICATION.md`](SPECIFICATION.md) may cite ISO for provenance.

### Open Knowledge Format (OKF) bundles

OKF bundles register concept types in root `schema.md`:

| Folder     | What it holds                                                |
| ---------- | ------------------------------------------------------------ |
| `specs/`   | Specifications for the deterministic `qualitymd` surface.    |
| `docs/`    | Project documentation, organized by the four Diátaxis modes. |
| `changes/` | Incremental work items — a spec and design doc per change.   |

### Smoke testing

- Do not add smoke-test scripts, utilities, fixtures, or code to repo.
- Temporary helpers only in `tmp/` or throwaway dirs; remove when done.

### Guides

Before work, read relevant [`docs/guides/`](docs/guides/index.md):

| When you are…                                   | Read                                                              |
| ----------------------------------------------- | ----------------------------------------------------------------- |
| Proposing or tracking a unit of work            | [Proposing a change](docs/guides/propose-a-change.md)             |
| Writing a functional spec (the `specs/` bundle) | [Writing functional specs](docs/guides/write-functional-specs.md) |
| Writing a design doc                            | [Writing design docs](docs/guides/write-design-docs.md)           |
| Reading or editing any OKF bundle               | [Working with OKF](docs/guides/work-with-okf.md)                  |

### Agent guidance files

- `CLAUDE.md`, `GEMINI.md` symlink here. Edit `AGENTS.md` only.
- Both symlinks gitignored.
