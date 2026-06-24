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
| Designing or reshaping an agent-run workflow    | [Designing agent-mediated UX](docs/guides/agent-mediated-ux.md)    |
| Designing or reshaping a CLI command            | [Designing CLI interfaces](docs/guides/cli-design.md)              |
| Adding a type or package to the Go code         | [Designing Go packages](docs/guides/design-go-packages.md)         |
| Writing Go code                                 | [Go style](docs/guides/go-style.md)                                |

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

### Quality-domain agnostic examples

- QUALITY.md is quality-domain agnostic. Concrete quality model content in this
  repo is illustrative unless it defines this project's model or a normative
  format rule.
- This includes example Areas, Factors, Requirements, Assessments, criteria,
  Rating Levels, Findings, recommendations, and quality-domain examples.
- When citing quality domains or factor families, make clear examples are brief,
  illustrative, overlapping, not exhaustive, and not mutually exclusive.
- Prefer domain-neutral principles first, then examples from several domains
  when examples help. Do not imply software product quality is the default use.

### Agentic use context

Domain agnostic does not mean context neutral. QUALITY.md is domain agnostic in
what a quality model can describe: software, documents, data sets, services,
operations, processes, AI assistants, agent harnesses, or other evaluated
entities. This project is not context neutral in how QUALITY.md is used.

The primary experience is agent- and skill-first. AI assistants and coding
agents read, author, evaluate, and improve `QUALITY.md`; the `/quality` skill
carries judgment; and the CLI provides deterministic support tooling. Preserve
that agentic use context in docs, specs, examples, and skill content.

Do not remove references to AI assistants, coding agents, agent-accessible
evidence, harnesses, skill workflows, or agent collaboration when they describe
how QUALITY.md is used. Do not treat that operating context as the default
modeled quality domain.

Do not flag a phrase solely because the evaluated project is an AI assistant,
coding agent, or harness. Those are valid project/use contexts. Flag only when
the wording makes that domain sound inherent to QUALITY.md, normal for all
QUALITY.md files, or the default model content.

Decision test:

- Use context: who uses QUALITY.md, through what workflow, with what tools.
  Agentic/AI language is appropriate and often preferred.
- Model domain: what a `QUALITY.md` evaluates. Keep this domain agnostic unless
  the example is explicitly scoped.
- Project/use context: the concrete project or harness being evaluated. AI
  assistant or harness language is acceptable when scoped to that project.
- Project self-description: this repository's own format, skill, CLI, docs, and
  examples. Agentic/AI language is appropriate when true for this project.

Examples:

- Good: "Use QUALITY.md with the `/quality` agent skill to align coding agents
  and teams."
- Good: "Record context that is agent-accessible."
- Good: "This example models an AI assistant harness."
- Good: "Evaluate and improve the quality of AI assistant projects and
  harnesses." Acceptable when describing that use case or this project's
  agentic tooling, not the universal scope of QUALITY.md.
- Good: "The CLI is support tooling for the agent-first workflow."
- Needs scoping: "QUALITY.md helps improve AI assistant quality." Better:
  "QUALITY.md can model AI assistant quality; it can also model other domains."
- Avoid: "A QUALITY.md normally evaluates a codebase or agent harness."
- Avoid: "QUALITY.md is for evaluating AI assistant projects and harnesses."
- Avoid: "Default factors include security, reliability, usability, and
  maintainability."
- Avoid: removing "agent-accessible" because it sounds AI-specific.

Satisfiable check before adding concrete quality model content:

- State or imply the example's quality domain.
- Make illustrative status clear unless it is a format rule or this project's
  own model.
- Reword if a reader could mistake it for a universal default.
- Frame brief lists as non-exhaustive and possibly overlapping.
- Use software product quality as the default only for explicitly software
  topics.

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
