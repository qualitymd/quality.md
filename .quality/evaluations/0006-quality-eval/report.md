# Evaluation Report

## Verdict

- **Root area:** QUALITY.md Project
- **Evaluation level:** not recorded
- **Rigor:** Standard
- **Evaluation verdict:** 🟡 Minimum
- **Rationale:** The root local rating is target, but aggregate quality is held to minimum by minimum-rated child areas for docs, distribution, and evaluation history.

## Scope

Standard.

- **Narrowing:** whole recorded run
- **In scope:** QUALITY.md Project; Agent Harness; qualitymd CLI; Installation and Distribution; Documentation and Examples; Evaluation History; Format Specification; QUALITY.md Project QUALITY.md; /quality Skill; Tooling Specs Bundle
- **Out of scope:** Applying recommendations, editing evaluated source, editing `QUALITY.md`, or creating external issues.; Migrating or rewriting stale historical evaluation records.; Exhaustive release matrix verification, hosted registry validation, or live install-channel tests beyond repository-visible smoke-check definitions.

## Selected findings and limitations

- `assessments/001-root-agents-can-reach-the-minimum-project-context-and-deeper-routed-guidance.json` at `AGENTS.md:5` [Low]: AGENTS.md names README.md and CONTRIBUTING.md as required context and links task-specific guides.
- `assessments/002-root-project-guidance-makes-task-boundaries-and-done-criteria-explicit.json` at `AGENTS.md:33` [Low]: Routine changes are scoped directly, while guide routing and change-case rules define when durable design history is needed.
- `assessments/003-root-agents-can-discover-the-tools-and-commands-needed-to-work-the-project.json` at `mise.toml:38` [Low]: The check task aggregates format, tidy, vet, lint, test, Markdown, and npm package checks.
- `assessments/004-root-quality-workflows-preserve-useful-state-in-durable-local-artifacts.json` at `skills/quality/SKILL.md:47` [Low]: The skill requires evaluate to write numbered records through the CLI and feedback logs under .quality/logs/.
- `assessments/005-root-verification-signals-are-discoverable-and-remediation-bearing.json` at `command:mise run check` [Low]: The full local check passed across gofmt, tidy, dprint, vet, tests, npm package verification, and golangci-lint.
- `assessments/006-root-core-project-standards-are-backed-by-enforceable-checks-or-reviewable-gates.json` at `.github/workflows/ci.yml:24` [Low]: CI intentionally stays in lockstep with mise run check.
- `assessments/007-root-agent-action-limits-and-mutation-gates-are-visible-before-consequential-changes.json` at `skills/quality/SKILL.md:161` [Low]: The user interaction contract requires decision briefs before confirmation-sensitive mutations.
- `assessments/008-format-spec-the-format-specification-admits-a-single-interpretation.json` at `SPECIFICATION.md:236` [Low]: The schema section defines required model properties and the model-wide ratingScale.
- Additional selected findings or limitations are available in `report.json`.

## Evidence basis

No command or source evidence basis was recorded in findings.

## Next action

- [001-refresh-pinned-installer-examples](recommendations/001-refresh-pinned-installer-examples.md) - The install guide no longer contains stale v0.5.1 pins, pinned examples cannot age silently, and the docs and distribution assessments reach target on re-evaluation.

## Area breakdown

| Area | Path | Area Rating | Area + Sub-Areas Rating | Factors |
| --- | --- | --- | --- | --- |
| QUALITY.md Project | `/` | 🔵 Target | 🟡 Minimum | Agent Harnessability: 🔵 Target; Agent Harnessability / Agent Accessibility: 🔵 Target; Agent Harnessability / Task Specifiability: 🔵 Target; Agent Harnessability / Agent Operability: 🔵 Target; Agent Harnessability / Continuity: 🔵 Target; Agent Harnessability / Self-Verifiability: 🔵 Target; Agent Harnessability / Enforcement of Standards: 🔵 Target; Agent Harnessability / Containment of Action: 🔵 Target |
| Agent Harness | `agent-harness` | 🔵 Target | 🔵 Target | Completeness: 🔵 Target; Coherence: 🔵 Target; Currentness: 🔵 Target; Assessability: 🔵 Target |
| qualitymd CLI | `cli` | 🔵 Target | 🔵 Target | Specification Conformance: 🔵 Target; Automation Compatibility: 🔵 Target; Usability: 🔵 Target; Maintainability: 🔵 Target |
| Installation and Distribution | `distribution` | 🟡 Minimum | 🟡 Minimum | Installability: 🟡 Minimum; Release Readiness: 🔵 Target |
| Documentation and Examples | `docs` | 🟡 Minimum | 🟡 Minimum | Approachability: 🔵 Target; Domain Range: 🔵 Target; Currentness: 🟡 Minimum |
| Evaluation History | `evaluation-history` | 🟡 Minimum | 🟡 Minimum | Reportability: 🟡 Minimum; Traceability: 🔵 Target |
| Format Specification | `format-spec` | 🔵 Target | 🔵 Target | Clarity: 🔵 Target; Consistency: 🔵 Target; Completeness: 🔵 Target; Verifiability: 🔵 Target; Domain Agnosticism: 🔵 Target |
| QUALITY.md Project QUALITY.md | `quality-md` | 🔵 Target | 🔵 Target | Context Grounding: 🔵 Target; Model Structure: 🔵 Target; Evaluability: 🔵 Target; Lifecycle Maintenance: 🔵 Target |
| /quality Skill | `quality-skill` | 🔵 Target | 🔵 Target | Judgment Grounding: 🔵 Target; Mutation Safety: 🔵 Target; Agent-Mediated UX: 🔵 Target; Workflow Completeness: 🔵 Target |
| Tooling Specs Bundle | `specs-bundle` | 🔵 Target | 🔵 Target | Traceability: 🔵 Target; Consistency: 🔵 Target |

## Area details

### QUALITY.md Project

- **Path:** /
- **Area rating:** 🔵 Target
  - All seven local Agent Harnessability requirements meet target.
- **+ Sub-Areas rating:** 🟡 Minimum
  - The root local rating is target, but aggregate quality is held to minimum by minimum-rated child areas for docs, distribution, and evaluation history.
- **Factor Agent Harnessability:** 🔵 Target
  - The project equips agents with context, operability, continuity, verification, standards, and containment.
- **Factor Agent Harnessability / Agent Accessibility:** 🔵 Target
  - Project context and routed guidance are reachable.
- **Factor Agent Harnessability / Task Specifiability:** 🔵 Target
  - Boundaries and done criteria are explicit.
- **Factor Agent Harnessability / Agent Operability:** 🔵 Target
  - Tools and commands are discoverable.
- **Factor Agent Harnessability / Continuity:** 🔵 Target
  - Workflow state is preserved in durable local artifacts.
- **Factor Agent Harnessability / Self-Verifiability:** 🔵 Target
  - Verification signals are runnable and actionable.
- **Factor Agent Harnessability / Enforcement of Standards:** 🔵 Target
  - Core standards are checked or review-gated.
- **Factor Agent Harnessability / Containment of Action:** 🔵 Target
  - Mutation gates and action limits are visible.
- **Analysis record:** `analysis/root.json`

### Agent Harness

- **Path:** agent-harness
- **Area rating:** 🔵 Target
  - The checked-in harness orients agents and routes them to deeper guidance with assessable controls.
- **+ Sub-Areas rating:** 🔵 Target
  - Leaf area aggregate equals the target local rating.
- **Factor Completeness:** 🔵 Target
  - Harness covers setup, scoped work, verification, mutation, and handoff guidance.
- **Factor Coherence:** 🔵 Target
  - Instructions preserve the split between skill judgment and CLI mechanics.
- **Factor Currentness:** 🔵 Target
  - Harness names current repo layout, workflows, commands, and compatibility expectations.
- **Factor Assessability:** 🔵 Target
  - Harness quality is inspectable through guides, hooks, workflows, and feedback logs.
- **Analysis record:** `analysis/agent-harness.json`

### qualitymd CLI

- **Path:** cli
- **Area rating:** 🔵 Target
  - All CLI requirements meet target with passing local checks and structured support surfaces.
- **+ Sub-Areas rating:** 🔵 Target
  - Leaf area aggregate equals the target local rating.
- **Factor Specification Conformance:** 🔵 Target
  - Implemented behavior and tests pass the current functional surface.
- **Factor Automation Compatibility:** 🔵 Target
  - Agent-consumed commands expose JSON or stable non-interactive behavior.
- **Factor Usability:** 🔵 Target
  - Help and status output support human recovery and next actions.
- **Factor Maintainability:** 🔵 Target
  - Go formatting, tests, vetting, linting, and package checks pass.
- **Analysis record:** `analysis/cli.json`

### Installation and distribution

- **Path:** distribution
- **Area rating:** 🟡 Minimum
  - Distribution gates are strong, but installability is held down by stale pinned installer examples.
- **+ Sub-Areas rating:** 🟡 Minimum
  - Leaf area aggregate equals the minimum local rating.
- **Factor Installability:** 🟡 Minimum
  - Stale v0.5.1 pins can steer agents and CI to old tooling.
- **Factor Release Readiness:** 🔵 Target
  - Release and smoke workflows verify GitHub, npm, and Homebrew channels.
- **Analysis record:** `analysis/distribution.json`

### Documentation and examples

- **Path:** docs
- **Area rating:** 🟡 Minimum
  - Docs are strong on approachability and domain range, but currentness is held down by stale v0.5.1 pinned installer examples.
- **+ Sub-Areas rating:** 🟡 Minimum
  - Leaf area aggregate equals the minimum local rating.
- **Factor Approachability:** 🔵 Target
  - Introductory docs foreground the /quality workflow and position the CLI as support tooling.
- **Factor Domain Range:** 🔵 Target
  - Docs represent non-software domains while preserving agentic context of use.
- **Factor Currentness:** 🟡 Minimum
  - Install examples pin stale v0.5.1 values.
- **Analysis record:** `analysis/docs.json`

### Evaluation history

- **Path:** evaluation-history
- **Area rating:** 🟡 Minimum
  - History remains inspectable and distinct, but the prior run is stale and non-reportable under the current record contract.
- **+ Sub-Areas rating:** 🟡 Minimum
  - Leaf area aggregate equals the minimum local rating.
- **Factor Reportability:** 🟡 Minimum
  - The prior run lacks required areaPath fields and is not reportable.
- **Factor Traceability:** 🔵 Target
  - The repository separates evaluation records, quality-log entries, and workflow feedback logs.
- **Analysis record:** `analysis/evaluation-history.json`

### Format specification

- **Path:** format-spec
- **Area rating:** 🔵 Target
  - All five format-spec requirements meet target across clarity, consistency, completeness, verifiability, and domain-agnosticism.
- **+ Sub-Areas rating:** 🔵 Target
  - Leaf area aggregate equals the target local rating.
- **Factor Clarity:** 🔵 Target
  - Rules and schema terms are sufficiently unambiguous.
- **Factor Consistency:** 🔵 Target
  - Terminology and roll-up semantics stay aligned.
- **Factor Completeness:** 🔵 Target
  - Required model, evaluation, and report behavior is specified.
- **Factor Verifiability:** 🔵 Target
  - Normative rules map to observable document or tool behavior.
- **Factor Domain Agnosticism:** 🔵 Target
  - The spec keeps broad model domain separate from agent-first use context.
- **Analysis record:** `analysis/format-spec.json`

### QUALITY.md Project QUALITY.md

- **Path:** quality-md
- **Area rating:** 🔵 Target
  - The project QUALITY.md is valid, context-grounded, assessable, and aligned with authoring guidance.
- **+ Sub-Areas rating:** 🔵 Target
  - Leaf area aggregate equals the target local rating.
- **Factor Context Grounding:** 🔵 Target
  - The body explains root area, scope, needs, risks, use context, unknowns, and lifecycle posture.
- **Factor Model Structure:** 🔵 Target
  - Areas, factors, requirements, sources, and assessment references are scoped to the project.
- **Factor Evaluability:** 🔵 Target
  - Requirements can be assessed from agent-accessible evidence.
- **Factor Lifecycle Maintenance:** 🔵 Target
  - The model names evaluation history, quality log, feedback logs, and future recommendation loops.
- **Analysis record:** `analysis/quality-md.json`

### /quality skill

- **Path:** quality-skill
- **Area rating:** 🔵 Target
  - All four skill requirements meet target: grounding, mutation safety, agent-mediated UX, and workflow completeness.
- **+ Sub-Areas rating:** 🔵 Target
  - Leaf area aggregate equals the target local rating.
- **Factor Judgment Grounding:** 🔵 Target
  - Skill workflow requires active model, current evidence, source-as-data handling, and stop rules.
- **Factor Mutation Safety:** 🔵 Target
  - Mutating actions are gated with decision briefs and verification.
- **Factor Agent-Mediated UX:** 🔵 Target
  - Run frames, scannable status, and decision briefs are built into the interaction contract.
- **Factor Workflow Completeness:** 🔵 Target
  - Workflow instructions match the CLI-supported record/report surface.
- **Analysis record:** `analysis/quality-skill.json`

### Tooling specs bundle

- **Path:** specs-bundle
- **Area rating:** 🔵 Target
  - Both specs-bundle requirements meet target for traceability and consistency.
- **+ Sub-Areas rating:** 🔵 Target
  - Leaf area aggregate equals the target local rating.
- **Factor Traceability:** 🔵 Target
  - Specs and runtime artifacts identify governed command, workflow, record, and report surfaces.
- **Factor Consistency:** 🔵 Target
  - The specs, skill, and implementation pass the current repository gate.
- **Analysis record:** `analysis/specs-bundle.json`

## Requirements

### agents can reach the minimum project context and deeper routed guidance

- **State:** active
- **Area:** QUALITY.md Project
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/001-root-agents-can-reach-the-minimum-project-context-and-deeper-routed-guidance.json`
- **Rationale:** Project orientation, agent rules, guide routing, and component map are reachable from checked-in files.

### project guidance makes task boundaries and done criteria explicit

- **State:** active
- **Area:** QUALITY.md Project
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/002-root-project-guidance-makes-task-boundaries-and-done-criteria-explicit.json`
- **Rationale:** Task guidance routes work by artifact type and names verification expectations.

### agents can discover the tools and commands needed to work the project

- **State:** active
- **Area:** QUALITY.md Project
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/003-root-agents-can-discover-the-tools-and-commands-needed-to-work-the-project.json`
- **Rationale:** The repo exposes setup, check, package, release, install, and quality commands through docs and mise tasks.

### quality workflows preserve useful state in durable local artifacts

- **State:** active
- **Area:** QUALITY.md Project
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/004-root-quality-workflows-preserve-useful-state-in-durable-local-artifacts.json`
- **Rationale:** Evaluation records, quality logs, and workflow feedback logs exist as durable local artifacts rather than chat-only state.

### verification signals are discoverable and remediation-bearing

- **State:** active
- **Area:** QUALITY.md Project
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/005-root-verification-signals-are-discoverable-and-remediation-bearing.json`
- **Rationale:** Local status, lint, evaluation status, CI, hooks, and mise checks give actionable feedback.

### core project standards are backed by enforceable checks or reviewable gates

- **State:** active
- **Area:** QUALITY.md Project
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/006-root-core-project-standards-are-backed-by-enforceable-checks-or-reviewable-gates.json`
- **Rationale:** Core formatting, Go quality, schema, packaging, and release expectations have runnable gates or explicit review routes.

### agent action limits and mutation gates are visible before consequential changes

- **State:** active
- **Area:** QUALITY.md Project
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/007-root-agent-action-limits-and-mutation-gates-are-visible-before-consequential-changes.json`
- **Rationale:** Agent rules and skill workflows make mutation class, confirmation, and action limits visible before consequential changes.

### the format specification admits a single interpretation

- **State:** active
- **Area:** Format Specification
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/008-format-spec-the-format-specification-admits-a-single-interpretation.json`
- **Rationale:** The specification defines vocabulary, schema shape, references, and evaluation semantics with settled obligation language.

### the format specification is internally consistent

- **State:** active
- **Area:** Format Specification
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/009-format-spec-the-format-specification-is-internally-consistent.json`
- **Rationale:** Terminology and rating/report semantics are consistent across schema and evaluation sections.

### the format specification is complete enough to implement and author from

- **State:** active
- **Area:** Format Specification
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/010-format-spec-the-format-specification-is-complete-enough-to-implement-and-author-from.json`
- **Rationale:** Model elements, references, extension points, evaluation records, and report expectations are specified without implementation reverse engineering.

### each normative rule is observable or testable

- **State:** active
- **Area:** Format Specification
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/011-format-spec-each-normative-rule-is-observable-or-testable.json`
- **Rationale:** Conformance turns on observable document, record, and report properties; judgment boundaries are explicit.

### the specification keeps model domain separate from agentic use context

- **State:** active
- **Area:** Format Specification
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/012-format-spec-the-specification-keeps-model-domain-separate-from-agentic-use-context.json`
- **Rationale:** The spec and domain guide describe broad modeled domains while preserving agent-first usage as context of use.

### the skill grounds setup, evaluation, and follow-up in the active model and evidence

- **State:** active
- **Area:** /quality Skill
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/013-quality-skill-the-skill-grounds-setup-evaluation-and-follow-up-in-the-active-model-and-evidence.json`
- **Rationale:** The skill requires version/status grounding, current evidence, source-as-data handling, and stop rules for weak support.

### the skill gates mutating actions with explicit decision briefs

- **State:** active
- **Area:** /quality Skill
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/014-quality-skill-the-skill-gates-mutating-actions-with-explicit-decision-briefs.json`
- **Rationale:** Setup, update, and recommendation follow-up mutation surfaces are gated, while evaluation limits itself to evaluation artifacts.

### the skill follows the agent-mediated UX guide

- **State:** active
- **Area:** /quality Skill
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/015-quality-skill-the-skill-follows-the-agent-mediated-ux-guide.json`
- **Rationale:** The skill contract requires status-first, evidence-led interaction, run frames, decision briefs, and scannable closeouts.

### the skill workflows match their functional specs and CLI support surface

- **State:** active
- **Area:** /quality Skill
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/016-quality-skill-the-skill-workflows-match-their-functional-specs-and-cli-support-surface.json`
- **Rationale:** Runtime workflows use available CLI commands, require dry-runs for record payloads, and prohibit hand-authored CLI-owned records.

### the CLI follows its functional specifications

- **State:** active
- **Area:** qualitymd CLI
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/017-cli-the-cli-follows-its-functional-specifications.json`
- **Rationale:** Implemented commands and tests satisfy the functional surface exercised by lint, status, evaluation record writes, and reporting.

### CLI commands expose stable machine-readable behavior where agents need it

- **State:** active
- **Area:** qualitymd CLI
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/018-cli-cli-commands-expose-stable-machine-readable-behavior-where-agents-need-it.json`
- **Rationale:** Version, lint, status, record writes, and package checks expose non-interactive JSON or stable command behavior.

### the CLI follows the project CLI design guide

- **State:** active
- **Area:** qualitymd CLI
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/019-cli-the-cli-follows-the-project-cli-design-guide.json`
- **Rationale:** Command help documents payload fields, examples, dry-run behavior, JSON receipts, and next-action diagnostics.

### the Go implementation follows the project Go style guide

- **State:** active
- **Area:** qualitymd CLI
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/020-cli-the-go-implementation-follows-the-project-go-style-guide.json`
- **Rationale:** Formatting, vetting, linting, tests, and package checks pass on the implementation.

### introductory docs foreground the agent-first workflow without making the CLI the main surface

- **State:** active
- **Area:** Documentation and Examples
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/021-docs-introductory-docs-foreground-the-agent-first-workflow-without-making-the-cli-the-main-surface.json`
- **Rationale:** The README frames /quality as the primary experience and the CLI as support tooling.

### docs keep modeled domains broad while preserving the agentic use context

- **State:** active
- **Area:** Documentation and Examples
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/022-docs-docs-keep-modeled-domains-broad-while-preserving-the-agentic-use-context.json`
- **Rationale:** Docs state broad modeled-domain range while retaining AI assistant/coding-agent use context.

### docs reflect the current command, install, release, and workflow surfaces

- **State:** active
- **Area:** Documentation and Examples
- **Rating:** 🟡 Minimum
- **Assessment result record:** `assessments/023-docs-docs-reflect-the-current-command-install-release-and-workflow-surfaces.json`
- **Rationale:** Most docs match current workflows, but install examples still pin an obsolete v0.5.1 release despite the active 0.11.0 CLI/skill pair.

### tooling specs identify the behavior they govern and stay linked to runtime artifacts

- **State:** active
- **Area:** Tooling Specs Bundle
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/024-specs-bundle-tooling-specs-identify-the-behavior-they-govern-and-stay-linked-to-runtime-artifacts.json`
- **Rationale:** Specs are organized by CLI, reports, evaluation records, and skill workflows, and runtime guidance points back to those governed surfaces.

### the specs bundle remains internally consistent and synchronized with implementation

- **State:** active
- **Area:** Tooling Specs Bundle
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/025-specs-bundle-the-specs-bundle-remains-internally-consistent-and-synchronized-with-implementation.json`
- **Rationale:** The checked implementation, skill, and package tests pass against the current specs and documented surfaces.

### supported install and update paths are documented, tested, and package the expected binary

- **State:** active
- **Area:** Installation and Distribution
- **Rating:** 🟡 Minimum
- **Assessment result record:** `assessments/026-distribution-supported-install-and-update-paths-are-documented-tested-and-package-the-expected-binary.json`
- **Rationale:** Install scripts and smoke checks cover the expected channels, but pinned install examples are stale and can steer agents/CI to an old version.

### release gates verify the artifacts and channels that users depend on

- **State:** active
- **Area:** Installation and Distribution
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/027-distribution-release-gates-verify-the-artifacts-and-channels-that-users-depend-on.json`
- **Rationale:** Release and smoke workflows cover GitHub archives, npm, Homebrew, checksums, notes, and package assembly.

### the agent harness orients agents and routes them to deeper guidance

- **State:** active
- **Area:** Agent Harness
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/028-agent-harness-the-agent-harness-orients-agents-and-routes-them-to-deeper-guidance.json`
- **Rationale:** AGENTS.md provides project context, component routing, guide routing, agent action rules, vocabulary conventions, and quality-domain guidance.

### the QUALITY.md model follows the active authoring guide family

- **State:** active
- **Area:** QUALITY.md Project QUALITY.md
- **Rating:** 🔵 Target
- **Assessment result record:** `assessments/029-quality-md-the-quality-md-model-follows-the-active-authoring-guide-family.json`
- **Rationale:** The model is valid, broad in domain, evidence-oriented, and explicit about use context, risks, scope, sources, and lifecycle unknowns.

### evaluation history and model-change history remain inspectable and distinct

- **State:** active
- **Area:** Evaluation History
- **Rating:** 🟡 Minimum
- **Assessment result record:** `assessments/030-evaluation-history-evaluation-history-and-model-change-history-remain-inspectable-and-distinct.json`
- **Rationale:** History is inspectable and distinct, but the prior run is stale and not reportable under the current record contract.

## Findings

- `assessments/001-root-agents-can-reach-the-minimum-project-context-and-deeper-routed-guidance.json` at `AGENTS.md:5`: AGENTS.md names README.md and CONTRIBUTING.md as required context and links task-specific guides.
- `assessments/002-root-project-guidance-makes-task-boundaries-and-done-criteria-explicit.json` at `AGENTS.md:33`: Routine changes are scoped directly, while guide routing and change-case rules define when durable design history is needed.
- `assessments/003-root-agents-can-discover-the-tools-and-commands-needed-to-work-the-project.json` at `mise.toml:38`: The check task aggregates format, tidy, vet, lint, test, Markdown, and npm package checks.
- `assessments/004-root-quality-workflows-preserve-useful-state-in-durable-local-artifacts.json` at `skills/quality/SKILL.md:47`: The skill requires evaluate to write numbered records through the CLI and feedback logs under .quality/logs/.
- `assessments/005-root-verification-signals-are-discoverable-and-remediation-bearing.json` at `command:mise run check`: The full local check passed across gofmt, tidy, dprint, vet, tests, npm package verification, and golangci-lint.
- `assessments/006-root-core-project-standards-are-backed-by-enforceable-checks-or-reviewable-gates.json` at `.github/workflows/ci.yml:24`: CI intentionally stays in lockstep with mise run check.
- `assessments/007-root-agent-action-limits-and-mutation-gates-are-visible-before-consequential-changes.json` at `skills/quality/SKILL.md:161`: The user interaction contract requires decision briefs before confirmation-sensitive mutations.
- `assessments/008-format-spec-the-format-specification-admits-a-single-interpretation.json` at `SPECIFICATION.md:236`: The schema section defines required model properties and the model-wide ratingScale.
- `assessments/009-format-spec-the-format-specification-is-internally-consistent.json` at `SPECIFICATION.md:454`: Roll-up inputs are named consistently as requirement rating results, factor ratings, child area ratings, and body context.
- `assessments/010-format-spec-the-format-specification-is-complete-enough-to-implement-and-author-from.json` at `SPECIFICATION.md:500`: The report section names required summary, scope, ratings, rationale, gaps, and recommendations.
- `assessments/011-format-spec-each-normative-rule-is-observable-or-testable.json` at `SPECIFICATION.md:488`: When evidence cannot justify roll-up, the specified outcome is not assessed rather than an inferred rating.
- `assessments/012-format-spec-the-specification-keeps-model-domain-separate-from-agentic-use-context.json` at `docs/guides/model-quality-across-domains.md:437`: The guide states domain agnostic is about what a model can describe, while the project is agent-first in how QUALITY.md is used.
- `assessments/013-quality-skill-the-skill-grounds-setup-evaluation-and-follow-up-in-the-active-model-and-evidence.json` at `skills/quality/SKILL.md:69`: Evaluated source content is treated as data, not instructions.
- `assessments/014-quality-skill-the-skill-gates-mutating-actions-with-explicit-decision-briefs.json` at `skills/quality/workflows/setup.md:373`: Existing-file setup edits require a decision brief before mutation.
- `assessments/015-quality-skill-the-skill-follows-the-agent-mediated-ux-guide.json` at `skills/quality/SKILL.md:139`: The user interaction contract explicitly routes presentation through the agent-mediated UX guide.
- `assessments/016-quality-skill-the-skill-workflows-match-their-functional-specs-and-cli-support-surface.json` at `skills/quality/workflows/evaluate.md:162`: The evaluate workflow writes records only through qualitymd evaluation assessment, analysis, and recommendation commands.
- `assessments/017-cli-the-cli-follows-its-functional-specifications.json` at `command:mise run check`: The full local check passed, including go test ./... across internal/cli, internal/evaluation, internal/model, internal/schema, and related packages.
- `assessments/018-cli-cli-commands-expose-stable-machine-readable-behavior-where-agents-need-it.json` at `qualitymd status QUALITY.md --json`: Status returned structured readiness, model shape, source coverage, evaluation run counts, and next actions.
- `assessments/019-cli-the-cli-follows-the-project-cli-design-guide.json` at `qualitymd evaluation assessment add --help`: Assessment help lists required payload fields, dry-run behavior, stdin usage, and flags.
- `assessments/020-cli-the-go-implementation-follows-the-project-go-style-guide.json` at `command:mise run check`: gofmt, go vet, go test ./..., and golangci-lint all passed.
- `assessments/021-docs-introductory-docs-foreground-the-agent-first-workflow-without-making-the-cli-the-main-surface.json` at `README.md:8`: The README says the /quality skill helps setup, evaluate, and evolve the model, and that the agentic workflow is primary.
- `assessments/022-docs-docs-keep-modeled-domains-broad-while-preserving-the-agentic-use-context.json` at `README.md:102`: The README's Quality Beyond Software section lists software, documentation, data, reports, services, operations, clinical handoffs, legal contracts, and other contexts.
- `assessments/023-docs-docs-reflect-the-current-command-install-release-and-workflow-surfaces.json` at `install.md:83`: The non-interactive pinned installer example uses QUALITYMD_VERSION=v0.5.1; PowerShell and sh -s examples repeat the stale v0.5.1 pin at nearby lines.
- `assessments/024-specs-bundle-tooling-specs-identify-the-behavior-they-govern-and-stay-linked-to-runtime-artifacts.json` at `QUALITY.md:21`: The repository model identifies SPECIFICATION.md, skills/quality, specs/skills, cmd/qualitymd, internal, and specs/cli as the major governed components.
- `assessments/025-specs-bundle-the-specs-bundle-remains-internally-consistent-and-synchronized-with-implementation.json` at `command:mise run check`: The full gate passed, including tests and the npm package/skill relative-link check.
- `assessments/026-distribution-supported-install-and-update-paths-are-documented-tested-and-package-the-expected-binary.json` at `install.md:83`: Agent/CI examples pin v0.5.1 while the visible CLI and skill metadata are 0.11.0.
- `assessments/027-distribution-release-gates-verify-the-artifacts-and-channels-that-users-depend-on.json` at `.github/workflows/install-smoke.yml:31`: The install-smoke workflow verifies managed shell/PowerShell installers, npm packages, and Homebrew channel behavior across supported OS matrices.
- `assessments/028-agent-harness-the-agent-harness-orients-agents-and-routes-them-to-deeper-guidance.json` at `AGENTS.md:46`: The guides table routes agents to release, change-case, spec, design, OKF, UX, domain-modeling, CLI, and Go guidance.
- `assessments/029-quality-md-the-quality-md-model-follows-the-active-authoring-guide-family.json` at `QUALITY.md:530`: The body explains the root area, broad quality domain, and agent-first use context; qualitymd lint reports zero findings.
- `assessments/030-evaluation-history-evaluation-history-and-model-change-history-remain-inspectable-and-distinct.json` at `qualitymd evaluation status .quality/evaluations/0005-subject-quality-eval`: The prior run reports 15 assessments and 4 analyses, but all old record shapes lack required areaPath fields and the run is not reportable.

## Advice

- [001-refresh-pinned-installer-examples](recommendations/001-refresh-pinned-installer-examples.md) [active] - The install guide no longer contains stale v0.5.1 pins, pinned examples cannot age silently, and the docs and distribution assessments reach target on re-evaluation.
- [002-reconcile-stale-evaluation-history](recommendations/002-reconcile-stale-evaluation-history.md) [active] - qualitymd status QUALITY.md --json no longer reports a stale non-reportable prior run, or the repository has an explicit documented archive state that keeps stale runs out of current-readiness gaps.
