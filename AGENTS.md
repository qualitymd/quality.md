# QUALITY.md

Public repo for **QUALITY.md**: `qualitymd` CLI, format spec, user docs.

`QUALITY.md` = YAML frontmatter target tree + Markdown docs. CLI is
deterministic/mechanical; skills judge and record through it.

## Key Docs

| Path                                                                                             | Purpose                                                                                                                                                                                                                                 |
| ------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [`README.md`](README.md)                                                                         | Project overview, current implementation status, and install notes.                                                                                                                                                                     |
| [`CONTRIBUTING.md`](CONTRIBUTING.md)                                                             | Development setup, `mise run` tasks, project layout, and release process.                                                                                                                                                               |
| [`SPECIFICATION.md`](SPECIFICATION.md)                                                           | Source of truth for the `QUALITY.md` format and evaluation semantics: Model / Target / Factor / Requirement / RatingScale schema, assessment -> finding -> rating-result chain, and Markdown body guidance.                             |
| [`specs/skills/quality-skill/authoring-guide.md`](specs/skills/quality-skill/authoring-guide.md) | Governing functional spec for the `/quality` skill's authoring guide; defines the guide's purpose, scope, conformance to `SPECIFICATION.md`, structure, and editorial contract.                                                         |
| [`skills/quality/resources/quality-md-guide.md`](skills/quality/resources/quality-md-guide.md)   | Runtime guide bundled with the `/quality` skill; gives agents and humans a self-contained practical guide to understanding and authoring `QUALITY.md` files, governed by the authoring-guide spec and conforming to `SPECIFICATION.md`. |

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
| `changes/` | Change Cases — formal work records with spec/design history. |

A Change Case records significant work: motivation, status, affected durable
specs/docs, a functional spec, and optional design doc.

### Smoke testing

- Do not add smoke-test scripts, utilities, fixtures, or code to repo.
- Temporary helpers only in `tmp/` or throwaway dirs; remove when done.

### Guides

Before work, read relevant [`docs/guides/`](docs/guides/index.md):

| When you are…                                   | Read                                                               |
| ----------------------------------------------- | ------------------------------------------------------------------ |
| Creating or advancing a Change Case             | [Working with change cases](docs/guides/work-with-change-cases.md) |
| Writing a functional spec (the `specs/` bundle) | [Writing functional specs](docs/guides/write-functional-specs.md)  |
| Writing a design doc                            | [Writing design docs](docs/guides/write-design-docs.md)            |
| Reading or editing any OKF bundle               | [Working with OKF](docs/guides/work-with-okf.md)                   |
| Designing or reshaping a CLI command            | [Designing CLI interfaces](docs/guides/cli-design.md)              |
| Adding a type or package to the Go code         | [Designing Go packages](docs/guides/design-go-packages.md)         |

Routine prompted edits do not require a Change Case. Use `changes/` only when
the user asks for a Change Case, when continuing an existing `changes/NNNN-*`
item, or when the work needs durable spec/design/review history. Other
changes/modifications follow the normal routine change guide: make the scoped
edit, update directly relevant docs, tests, and specs, and verify.

### Agent guidance files

- `CLAUDE.md`, `GEMINI.md` symlink here. Edit `AGENTS.md` only.
- Both symlinks gitignored.
