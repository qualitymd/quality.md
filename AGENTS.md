# QUALITY.md

## Project Context

QUALITY.md is an open format for modeling a project's quality for the purpose
of evaluation, team/agent alignment, and continuous improvement.

Read [`README.md`](README.md) and [`CONTRIBUTING.md`](CONTRIBUTING.md)
before you continue for important project context and development guidance.

The QUALITY.md experience is largely agent and skill-first. Users do not
typically use the CLI for most use cases. Instead, the CLI and edits to
`QUALITY.md` are managed by the agent skill. Users are still encouraged to edit
`QUALITY.md` manually or with thoughtful AI assistance, especially the Markdown
body. User-facing docs, guides, explainers, etc. should foreground the
`/quality` agent skill or the `QUALITY.md` file itself and only highlight the
CLI if necessary.

## Major Components

| Component         | Where to look                                                                                                                                                         |
| ----------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| QUALITY.md format | [`SPECIFICATION.md`](SPECIFICATION.md) is the source of truth for the model schema, Markdown body guidance, and evaluation semantics.                                 |
| `/quality` skill  | Runtime files live in [`skills/quality/`](skills/quality/); functional specs and guide outlines live in [`specs/skills/quality-skill/`](specs/skills/quality-skill/). |
| `qualitymd` CLI   | Source starts at [`cmd/qualitymd/`](cmd/qualitymd/) and [`internal/`](internal/); CLI specs live in [`specs/cli/`](specs/cli/) and [`specs/cli.md`](specs/cli.md).    |

## Working Rules

### Instruction style

Keep this file extremely concise. Brevity over grammar.

### Routine changes

Routine prompted edits do not require a Change Case. Use `changes/` only when
the user asks for a Change Case, when continuing an existing `changes/NNNN-*`
item, or when the work needs durable spec/design/review history. Other routine
changes follow the normal change guide: make the scoped edit, update directly
relevant docs, tests, and specs, and verify.

### Smoke testing

- Do not add smoke-test scripts, utilities, fixtures, or code to the repo.
- Temporary helpers only in `tmp/` or throwaway dirs; remove when done.

## Guides

Before work, read the relevant [`docs/guides/`](docs/guides/index.md):

| When you are…                                   | Read                                                               |
| ----------------------------------------------- | ------------------------------------------------------------------ |
| Cutting or verifying a release                  | [Cut a release](docs/guides/cut-a-release.md)                      |
| Creating or advancing a Change Case             | [Working with change cases](docs/guides/work-with-change-cases.md) |
| Writing a functional spec (the `specs/` bundle) | [Writing functional specs](docs/guides/write-functional-specs.md)  |
| Writing a design doc                            | [Writing design docs](docs/guides/write-design-docs.md)            |
| Reading or editing any OKF bundle               | [Working with OKF](docs/guides/work-with-okf.md)                   |
| Designing or reshaping a CLI command            | [Designing CLI interfaces](docs/guides/cli-design.md)              |
| Adding a type or package to the Go code         | [Designing Go packages](docs/guides/design-go-packages.md)         |

## Repository Conventions

### Naming QUALITY.md

- Use `QUALITY.md` in backticks for the literal file/path agents should read or
  edit.
- Use QUALITY.md plain for the project, format, standard, or product concept.
- Use **QUALITY.md** only for first-mention emphasis in user-facing intro prose.
- Prefer no bold in agent instructions, specs, and dense technical docs.

### Open Knowledge Format (OKF) bundles

OKF bundles register concept types in the root `schema.md`:

| Folder     | What it holds                                                |
| ---------- | ------------------------------------------------------------ |
| `specs/`   | Specifications for the deterministic `qualitymd` surface.    |
| `docs/`    | Project documentation, organized by the four Diátaxis modes. |
| `changes/` | Change Cases — formal work records with spec/design history. |

A Change Case records significant work: motivation, status, affected durable
artifacts, a functional spec, and optional design doc.

### Referencing ISO standards

- Keep ISO lineage background.
- Do not cite specific ISO standards in public code/artifacts unless requested
  or relevant to the file's purpose.
- Use QUALITY.md terms: Targets, Factors, Requirements.
- [`SPECIFICATION.md`](SPECIFICATION.md) may cite ISO for provenance.

### Agent guidance files

- `CLAUDE.md` and `GEMINI.md` symlink to this file. Edit `AGENTS.md` only.
- Both symlinks are gitignored.
