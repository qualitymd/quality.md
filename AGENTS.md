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

### Early-alpha compatibility

QUALITY.md is early alpha. Breaking changes are acceptable when they keep the
current model, skill, CLI, specs, and docs simpler and clearer.

Do not author backward-compatibility shims, legacy aliases, fallback readers,
dual writers, migration commands, deprecated command paths, or legacy specs
unless an active spec or release task explicitly requires them. Prefer clean
breaks: update the current contract, tests, docs, examples, and release notes
together.

When legacy compatibility code, specs, or docs are found in active surfaces,
remove them as part of the scoped change when safe. Preserve historical records
in `changes/archive/`, changelogs, and append-only logs unless the task is
explicitly cleaning history.

### Smoke testing

- Do not add smoke-test scripts, utilities, fixtures, or code to the repo.
- Temporary helpers only in `tmp/` or throwaway dirs; remove when done.

## Guides

Before work, read the relevant [`docs/guides/`](docs/guides/index.md):

| When you are…                                     | Read                                                                           |
| ------------------------------------------------- | ------------------------------------------------------------------------------ |
| Cutting or verifying a release                    | [Cut a release](docs/guides/cut-a-release.md)                                  |
| Creating or advancing a Change Case               | [Working with change cases](docs/guides/work-with-change-cases.md)             |
| Writing a functional spec (the `specs/` bundle)   | [Writing functional specs](docs/guides/write-functional-specs.md)              |
| Writing a design doc                              | [Writing design docs](docs/guides/write-design-docs.md)                        |
| Reading or editing any OKF bundle                 | [Working with OKF](docs/guides/work-with-okf.md)                               |
| Designing or reshaping an agent-run workflow      | [Designing agent-mediated UX](docs/guides/agent-mediated-ux.md)                |
| Adding or reviewing example quality-model content | [Modeling quality across domains](docs/guides/model-quality-across-domains.md) |
| Designing or reshaping a CLI command              | [Designing CLI interfaces](docs/guides/cli-design.md)                          |
| Adding a type or package to the Go code           | [Designing Go packages](docs/guides/design-go-packages.md)                     |
| Writing Go code                                   | [Go style](docs/guides/go-style.md)                                            |

## Repository Conventions

### Naming QUALITY.md

- Use QUALITY.md plain by default, including when referring to QUALITY.md in
  the abstract or conceptual sense.
- Use `QUALITY.md` in backticks when describing a concrete instance of a
  QUALITY.md in an operational use case.
- Use bold/emphasized **QUALITY.md** only for first-mention emphasis in user-facing intro prose.
- Prefer no bold in agent instructions, specs, and dense technical docs.

### QUALITY.md vocabulary capitalization

- Capitalize formal model concepts when used as type names or terms of art:
  Model, Area, Factor, Requirement, Assessment, Finding, Rating Scale,
  Rating Level, Rating Result, Evaluation Report.
- Use lowercase in ordinary prose: area, factor, requirement, assessment,
  finding, rating, recommendation.
- Use backticks for concrete YAML fields, file names, commands, and literal
  values: `areas`, `factors`, `requirements`, `ratingScale`, `QUALITY.md`,
  `qualitymd`.
- Prefer lowercase in README, guides, tutorials, and user-facing prose unless
  capitalization improves precision.

### Keep the motivation and taxonomy registers distinct

- The stewardship/care core language — stewardship, care, tending, vulnerability,
  concern — is **motivation-layer**: it describes *why* a concern exists and what
  it means to tend an entity. The taxonomy — factor, area, requirement,
  constituent, audience — names the slots in the Model.
- Do not let a motivation-layer word modify or replace a taxonomy noun. A
  stewardship concern *projects into* a factor/constituent/audience; it is not one.
  Avoid "stewardship factor" / "stewardship lens" / "care requirement" — they
  demote a term of art to a subcategory of the philosophical word.
- Name the root's recurring factors **model-wide** or **cross-cutting factors**
  (the established terms). You may note they *trace to* stewardship concerns, but
  do not render that link by making "stewardship" an adjective on the taxonomy
  noun.
- The singular gloss "a factor is a quality *lens*" is fine — it defines what a
  factor is. Only a philosophical word substituting for the noun is the problem.

### Quality-domain agnostic examples and agentic use context

Authoritative rules:
[Modeling quality across domains](docs/guides/model-quality-across-domains.md)
(sections "Rules for domain-agnostic example content" and "Agentic use context").
Read it before adding or reviewing example quality-model content. Summary:

- QUALITY.md is quality-domain agnostic. Concrete model content here is
  illustrative unless it defines this project's own model or a normative format
  rule; mark it so and frame domain/factor lists as non-exhaustive and
  overlapping.
- Lead with domain-neutral principles; do not make software product quality the
  default. For worked examples, pair software/product quality with one cite-worthy
  secondary domain, balanced.
- Factors are earned per Model from the modeled entity's own risks and needs; do
  not adopt an external standard's characteristic list as a default factor family.
- Domain agnostic is not context neutral: the *modeled domain* stays agnostic, but
  the *use context* is agent- and skill-first. Preserve agentic/AI references that
  describe how QUALITY.md is used; do not treat that operating context as the
  default modeled domain. Flag AI/harness wording only when it sounds inherent to
  all QUALITY.md files, not when it correctly describes a use context or this
  project. The guide's decision test resolves which register a phrase is in.
- Use-context constituents (agent harness, QUALITY.md self-check) may have explicit
  guidance, but their factors and requirements stay agnostic to the served domain.

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
- Use QUALITY.md vocabulary instead of ISO terms.
- [`SPECIFICATION.md`](SPECIFICATION.md) may cite ISO for provenance.

### Agent guidance files

- `CLAUDE.md` and `GEMINI.md` symlink to this file. Edit `AGENTS.md` only.
- Both symlinks are gitignored.
