---
title: QUALITY.md Project
description: >
  The open QUALITY.md format and the agent-first quality-management toolchain
  around it: the /quality skill, qualitymd CLI, docs, specs, packaging,
  distribution, repository-owned agent harness, and this repository's own
  QUALITY.md model.
ratingScale:
  - level: outstanding
    title: 🟢 Outstanding
    description: "The stretch band: the artifact exceeds the quality requirement with meaningful margin."
    criterion: "Exceeds the requirement; satisfies it with margin to spare."
  - level: target
    title: 🔵 Target
    description: "The expected good state: the artifact satisfies the quality requirement."
    criterion: "Satisfies the requirement."
  - level: minimum
    title: 🟡 Minimum
    description: "The acceptable floor: the artifact falls short of the goal but remains good enough to rely on."
    criterion: "Falls short of the target but remains acceptable for pre-release use."
  - level: unacceptable
    title: 🔴 Unacceptable
    description: "Below the floor: the artifact is not good enough to rely on."
    criterion: "Does not meet the requirement to an acceptable degree."
factors:
  # Agent harnessability rates the broader equipping capability across the
  # project. The agent-harness area below rates the checked-in steering and
  # owned-control artifacts themselves.
  agent-harnessability:
    title: Agent Harnessability
    description: >
      The degree to which the project's checked-in materials, tools, workflows,
      feedback signals, standards, and action limits equip an AI assistant or
      coding agent to understand the project, take scoped work, operate the
      environment, preserve and resume state, verify its output, and stay safely
      bounded while preserving clear human direction, review, and accountability.
    factors:
      agent-accessibility:
        title: Agent Accessibility
        description: >
          Decision-relevant project knowledge is reachable, selective, and
          intelligible to an agent working in context.
        requirements:
          agents-can-reach-the-minimum-project-context-and-deeper-routed-guidance:
            title: agents can reach the minimum project context and deeper routed guidance
            assessment: >
              Inspect AGENTS.md, README.md, CONTRIBUTING.md, docs/guides/,
              skills/quality/, and specs/ to confirm a fresh agent can find the
              project purpose, operating rules, authoring/evaluation guidance,
              and task-specific references without relying on private memory.
      task-specifiability:
        title: Task Specifiability
        description: >
          Work can be handed to an agent as a scoped assignment with visible
          success criteria, boundaries, and verification expectations.
        requirements:
          project-guidance-makes-task-boundaries-and-done-criteria-explicit:
            title: project guidance makes task boundaries and done criteria explicit
            assessment: >
              Inspect AGENTS.md, CONTRIBUTING.md, docs/guides/, specs/, and
              changes/ to confirm agent-facing work can be scoped by artifact,
              expected outcome, relevant guide, and verification command.
      agent-operability:
        title: Agent Operability
        description: >
          A fresh agent session can establish and operate the working
          environment from recorded project materials.
        requirements:
          agents-can-discover-the-tools-and-commands-needed-to-work-the-project:
            title: agents can discover the tools and commands needed to work the project
            assessment: >
              Inspect CONTRIBUTING.md, mise.toml, install.md, README.md, and
              workflow specs to confirm setup, build, test, lint, release, and
              qualitymd commands are recorded with enough context for agent use.
      continuity:
        title: Continuity
        description: >
          Agent work can preserve state and resume across interruption,
          compaction, handoff, and fresh sessions.
        requirements:
          quality-workflows-preserve-useful-state-in-durable-local-artifacts:
            title: quality workflows preserve useful state in durable local artifacts
            assessment: >
              Inspect .quality/evaluations/, .quality/changelog/, .quality/logs/,
              specs/skills/quality-skill/, and skills/quality/workflows/ to
              confirm evaluation records, quality changelog entries, and workflow
              feedback logs capture decisions, progress, verification, and
              remaining gaps without depending on chat history.
      self-verifiability:
        title: Self-Verifiability
        description: >
          The project gives agents runnable or inspectable feedback signals they
          can use to check their own work.
        requirements:
          verification-signals-are-discoverable-and-remediation-bearing:
            title: verification signals are discoverable and remediation-bearing
            assessment: >
              Inspect mise.toml, CONTRIBUTING.md, .github/workflows/, .githooks/,
              qualitymd lint/status behavior, and relevant specs to confirm an
              agent can run or inspect checks and understand failures well enough
              to act on them.
      enforcement-of-standards:
        title: Enforcement of Standards
        description: >
          Stated quality standards hold through gates or equivalent controls,
          not advisory prose alone.
        requirements:
          core-project-standards-are-backed-by-enforceable-checks-or-reviewable-gates:
            title: core project standards are backed by enforceable checks or reviewable gates
            assessment: >
              Inspect mise.toml, .github/workflows/, .githooks/, .prettierrc.json,
              .golangci.yml, specs, and qualitymd lint behavior to confirm
              formatting, code quality, schema validity, packaging, and release
              expectations are enforced or explicitly routed through review.
      containment-of-action:
        title: Containment of Action
        description: >
          Agent-permitted actions are bounded by project rules, local tooling,
          and approval expectations so work stays inside intended scope.
        requirements:
          agent-action-limits-and-mutation-gates-are-visible-before-consequential-changes:
            title: agent action limits and mutation gates are visible before consequential changes
            assessment: >
              Inspect AGENTS.md, skills/quality/SKILL.md, skills/quality/workflows/,
              specs/skills/quality-skill/, and repository hooks/workflows to
              confirm file mutation, tooling updates, recommendation follow-up,
              release actions, and external handoff require explicit gates where
              consequential.
areas:
  format-spec:
    title: Format Specification
    description: >
      The normative QUALITY.md format and evaluation-semantics specification.
    source: ./SPECIFICATION.md
    factors:
      clarity:
        title: Clarity
        description: >
          The specification's rules have one settled meaning for independent
          authors, implementers, evaluators, and report renderers.
        requirements:
          the-format-specification-admits-a-single-interpretation:
            title: the format specification admits a single interpretation
            assessment: >
              Assess SPECIFICATION.md for defined terms, clear obligation
              strength, unambiguous field semantics, explicit malformed-content
              handling, and examples that do not carry hidden rules.
      consistency:
        title: Consistency
        description: >
          The specification uses one concept vocabulary and keeps examples,
          schema descriptions, and evaluation semantics aligned.
        requirements:
          the-format-specification-is-internally-consistent:
            title: the format specification is internally consistent
            assessment: >
              Compare terminology, schema rules, examples, model-reference
              grammar, and evaluation/report semantics within SPECIFICATION.md
              and against quality.schema.json where structural schema is named.
      completeness:
        title: Completeness
        description: >
          The specification covers enough document shape, frontmatter schema,
          body semantics, and evaluation/report behavior for conforming
          implementations and authors.
        requirements:
          the-format-specification-is-complete-enough-to-implement-and-author-from:
            title: the format specification is complete enough to implement and author from
            assessment: >
              Confirm every required model element, allowed shape, cardinality,
              reference form, extension rule, and required evaluation/report
              distinction is defined in SPECIFICATION.md without requiring
              implementation reverse engineering.
      verifiability:
        title: Verifiability
        description: >
          Conformance turns on observable document or implementation behavior
          rather than private author intent.
        requirements:
          each-normative-rule-is-observable-or-testable:
            title: each normative rule is observable or testable
            assessment: >
              Review conformance, schema, model-reference, and evaluation
              sections to confirm an independent linter, parser, evaluator, or
              renderer could decide whether each rule is satisfied.
      domain-agnosticism:
        title: Domain Agnosticism
        description: >
          The format supports many modeled domains without making software,
          agent harnesses, or any other domain the default model content.
        requirements:
          the-specification-keeps-model-domain-separate-from-agentic-use-context:
            title: the specification keeps model domain separate from agentic use context
            assessment: >
              Assess SPECIFICATION.md and docs/guides/model-quality-across-domains.md
              for domain-neutral model rules, illustrative examples marked as
              such, and clear separation between what QUALITY.md can model and
              how this project expects the format to be used through agents.

  quality-skill:
    title: /quality Skill
    description: >
      The runtime agent skill that carries quality judgment, setup, evaluation,
      update, recommendation follow-up, and authoring guidance.
    source: ./skills/quality
    factors:
      judgment-grounding:
        title: Judgment Grounding
        description: >
          The skill keeps evaluative judgment tied to the active model, current
          evidence, relevant guides, and the CLI's deterministic state.
        requirements:
          the-skill-grounds-setup-evaluation-and-follow-up-in-the-active-model-and-evidence:
            title: the skill grounds setup, evaluation, and follow-up in the active model and evidence
            assessment: >
              Assess skills/quality/SKILL.md, workflows/, guides/, and resources/
              against specs/skills/quality-skill/ to confirm the skill reads the
              required guidance, distinguishes model validity from evaluated-source
              quality, treats source content as data, and stops when evidence or
              tooling cannot support the workflow.
      mutation-safety:
        title: Mutation Safety
        description: >
          Mutating workflows make artifact class, confirmation, alternatives,
          and verification visible before changing project state or tooling.
        requirements:
          the-skill-gates-mutating-actions-with-explicit-decision-briefs:
            title: the skill gates mutating actions with explicit decision briefs
            assessment: >
              Inspect setup, evaluate, update, and recommendation follow-up
              guidance to confirm QUALITY.md edits, evaluated-source edits,
              evaluation artifacts, tooling updates, and handoff actions are
              gated and verified according to their risk.
      agent-mediated-ux:
        title: Agent-Mediated UX
        description: >
          The skill presents workflow state, questions, gates, and closeouts so a
          user can see status, recommendation, and next action quickly.
        requirements:
          the-skill-follows-the-agent-mediated-ux-guide:
            title: the skill follows the agent-mediated UX guide
            assessment: >
              Assess runtime prompts and workflow instructions against
              docs/guides/agent-mediated-ux.md, including run frames, discovery
              teaching, decision gates, status-first closeouts, and scannable
              labels.
      workflow-completeness:
        title: Workflow Completeness
        description: >
          Public workflows are complete enough to run end to end without hidden
          manual artifact creation or stale command assumptions.
        requirements:
          the-skill-workflows-match-their-functional-specs-and-cli-support-surface:
            title: the skill workflows match their functional specs and CLI support surface
            assessment: >
              Compare skills/quality/workflows/, resources/cli-quick-reference.md,
              and specs/skills/quality-skill/workflows/ to confirm setup,
              evaluate, and update route correctly, use available CLI commands,
              and do not hand-author CLI-owned artifacts.

  cli:
    title: qualitymd CLI
    description: >
      The deterministic support tooling for validating models, reporting status,
      managing evaluation records, exposing the active spec, and maintaining the
      skill/CLI compatibility pair.
    source: ./cmd/qualitymd
    factors:
      specification-conformance:
        title: Specification Conformance
        description: >
          CLI behavior matches the functional command specs and the active
          QUALITY.md format semantics.
        requirements:
          the-cli-follows-its-functional-specifications:
            title: the CLI follows its functional specifications
            assessment: >
              Assess cmd/qualitymd/, internal/, and tests against specs/cli.md,
              specs/cli/, specs/evaluation-records.md, and SPECIFICATION.md for
              command behavior, flags, exit codes, output contracts, workspace
              resolution, lint rules, and evaluation-record mechanics.
      automation-compatibility:
        title: Automation Compatibility
        description: >
          Agents, CI, scripts, and the /quality skill can drive commands without
          polluted streams, hidden prompts, or unstable structured data.
        requirements:
          cli-commands-expose-stable-machine-readable-behavior-where-agents-need-it:
            title: CLI commands expose stable machine-readable behavior where agents need it
            assessment: >
              Inspect CLI implementations and tests for JSON modes, stdout/stderr
              separation, deterministic ordering, exit status behavior, and
              non-interactive command paths used by the /quality skill.
      usability:
        title: Usability
        description: >
          Human users can discover commands, understand results, and recover
          from mistakes when they need the support tooling directly.
        requirements:
          the-cli-follows-the-project-cli-design-guide:
            title: the CLI follows the project CLI design guide
            assessment: >
              Assess command arguments, flags, help, errors, examples, ambient
              notices, and next actions against docs/guides/cli-design.md.
      maintainability:
        title: Maintainability
        description: >
          The Go implementation can be changed safely by maintainers and agents
          beyond what deterministic checks alone enforce.
        requirements:
          the-go-implementation-follows-the-project-go-style-guide:
            title: the Go implementation follows the project Go style guide
            assessment: >
              Assess cmd/qualitymd/ and internal/ against docs/guides/go-style.md,
              focusing on naming, package boundaries, error handling, interfaces,
              state, tests, and comments while leaving formatting and mechanical
              gates to mise run check.

  docs:
    title: Documentation and Examples
    description: >
      User-facing, contributor-facing, and normative-adjacent documentation that
      explains QUALITY.md, teaches workflows, and demonstrates domain breadth.
    source: ./docs
    factors:
      approachability:
        title: Approachability
        description: >
          Newcomers can understand what QUALITY.md is, why the /quality workflow
          is primary, and how to reach a first useful result.
        requirements:
          introductory-docs-foreground-the-agent-first-workflow-without-making-the-cli-the-main-surface:
            title: introductory docs foreground the agent-first workflow without making the CLI the main surface
            assessment: >
              Assess README.md, install.md, and docs/ to confirm users are guided
              first to the /quality skill and QUALITY.md file, with the CLI
              positioned as deterministic support tooling.
      domain-range:
        title: Domain Range
        description: >
          Documentation makes the format's broad modeled-domain range visible
          without defaulting to software product quality.
        requirements:
          docs-keep-modeled-domains-broad-while-preserving-the-agentic-use-context:
            title: docs keep modeled domains broad while preserving the agentic use context
            assessment: >
              Assess README.md, SPECIFICATION.md examples, docs/guides/model-quality-across-domains.md,
              and skill guide examples to confirm software examples are framed
              as illustrative, non-software examples are represented where needed,
              and AI-assistant/coding-agent language is used for context of use
              rather than as a universal modeled domain.
      currentness:
        title: Currentness
        description: >
          Documentation matches the shipped skill, CLI, format, install channels,
          and repository conventions.
        requirements:
          docs-reflect-the-current-command-install-release-and-workflow-surfaces:
            title: docs reflect the current command, install, release, and workflow surfaces
            assessment: >
              Compare README.md, install.md, CONTRIBUTING.md, docs/, specs/, and
              package/install metadata against implemented commands, workflows,
              release scripts, and skill metadata.

  specs-bundle:
    title: Tooling Specs Bundle
    description: >
      The OKF specifications for the deterministic qualitymd surface, evaluation
      records, reports, JSON Schema, and the /quality skill.
    source: ./specs
    factors:
      traceability:
        title: Traceability
        description: >
          Runtime behavior, tests, reports, and skill guidance can be traced back
          to durable functional specifications.
        requirements:
          tooling-specs-identify-the-behavior-they-govern-and-stay-linked-to-runtime-artifacts:
            title: tooling specs identify the behavior they govern and stay linked to runtime artifacts
            assessment: >
              Inspect specs/, internal/ tests, skills/quality/, and report
              examples to confirm specs name the command, workflow, record, or
              guide behavior they govern and runtime artifacts cite the right
              specs or guides.
      consistency:
        title: Consistency
        description: >
          The specs agree with each other and with the active runtime artifacts.
        requirements:
          the-specs-bundle-remains-internally-consistent-and-synchronized-with-implementation:
            title: the specs bundle remains internally consistent and synchronized with implementation
            assessment: >
              Compare specs/cli/, specs/evaluation-records/, specs/reports/,
              specs/skills/, and implemented behavior for command names, record
              fields, report contracts, workflow modes, and terminology.

  distribution:
    title: Installation and Distribution
    description: >
      The install scripts, npm packaging, release automation, Homebrew/GitHub
      distribution support, and smoke checks that get the CLI and skill to users.
    source: ./install
    factors:
      installability:
        title: Installability
        description: >
          Users and agents can install or update the supported tooling through
          documented channels on supported platforms.
        requirements:
          supported-install-and-update-paths-are-documented-tested-and-package-the-expected-binary:
            title: supported install and update paths are documented, tested, and package the expected binary
            assessment: >
              Inspect install/, npm/quality.md/, scripts/build-npm.mjs,
              scripts/check-npm-package.mjs, .goreleaser.yaml, .github/workflows/,
              install.md, and README.md for cross-platform install behavior,
              managed-install markers, package contents, and update guidance.
      release-readiness:
        title: Release Readiness
        description: >
          Release preparation and smoke checks reduce the risk of shipping a
          broken or incompatible skill/CLI pair.
        requirements:
          release-gates-verify-the-artifacts-and-channels-that-users-depend-on:
            title: release gates verify the artifacts and channels that users depend on
            assessment: >
              Inspect scripts/check-release.mjs, scripts/extract-release-notes.mjs,
              CHANGELOG.md, docs/guides/cut-a-release.md, .github/workflows/,
              and mise release tasks for checks that cover versioning,
              packaging, changelog, smoke installs, and compatibility metadata.

  # This area rates the checked-in steering and owned-control artifacts
  # themselves. The agent harnessability factor above rates the broader
  # equipping capability those artifacts help create.
  agent-harness:
    title: Agent Harness
    description: >
      The repository-owned steering and control artifacts that guide, verify, and
      bound AI assistant or coding-agent work, distinct from the broader agent
      harnessability capability they support.
    source: ./AGENTS.md
    factors:
      completeness:
        title: Completeness
        description: >
          The harness covers the feedforward and feedback controls agents need
          for setup, scoped work, verification, mutation, and handoff.
      coherence:
        title: Coherence
        description: >
          Harness instructions, skills, guides, hooks, workflows, and specs do
          not contradict each other or blur responsibility between skill judgment
          and CLI mechanics.
      currentness:
        title: Currentness
        description: >
          Harness guidance matches the current repository layout, workflow
          names, command surfaces, and compatibility expectations.
      assessability:
        title: Assessability
        description: >
          Harness quality can be checked through inspectable artifacts,
          representative workflow logs, and runnable verification signals.
    requirements:
      the-agent-harness-orients-agents-and-routes-them-to-deeper-guidance:
        title: the agent harness orients agents and routes them to deeper guidance
        factors:
          - completeness
          - coherence
          - currentness
          - assessability
        assessment: >
          Assess AGENTS.md, skills/quality/, docs/guides/, specs/skills/,
          .githooks/, .github/workflows/, mise.toml, and .quality/logs/ to
          confirm the checked-in harness defines project rules, routes to
          relevant guides, exposes verification, records workflow feedback, and
          distinguishes agent-harness artifacts from agent harnessability as a
          project-wide factor.

  # This area evaluates the concrete QUALITY.md artifact. Its `source` is the
  # model file itself; the requirement's assessment names the guide family used
  # to judge it.
  quality-md:
    title: QUALITY.md Project QUALITY.md
    description: >
      This repository's own QUALITY.md model, including the structured model,
      Markdown judgment context, setup assumptions, and maintenance posture.
    source: ./QUALITY.md
    factors:
      context-grounding:
        title: Context Grounding
        description: >
          The Markdown body explains root area, scope, domain breadth, use
          context, lifecycle, needs, risks, unknowns, and review state well
          enough for later humans and agents.
      model-structure:
        title: Model Structure
        description: >
          Areas, factors, requirements, sources, and assessment references are
          scoped, traceable, and shaped by the project's actual constituents.
      evaluability:
        title: Evaluability
        description: >
          Requirements are assessable from agent-accessible evidence and give
          future evaluations enough information to distinguish rating levels.
      lifecycle-maintenance:
        title: Lifecycle Maintenance
        description: >
          The model can evolve with the project, evaluation history, quality
          log, feedback logs, and future recommendations.
    requirements:
      the-quality-md-model-follows-the-active-authoring-guide-family:
        title: the QUALITY.md model follows the active authoring guide family
        factors:
          - context-grounding
          - model-structure
          - evaluability
          - lifecycle-maintenance
        assessment: >
          Assess ./QUALITY.md against skills/quality/guides/authoring.md and
          its routed sub-guides under skills/quality/guides/authoring/, especially
          whether the body credibly supports the model, factors come from visible
          needs and risks, requirements are assessable, sources are inspectable,
          agent harnessability is distinct from the agent-harness area, and
          unknowns or open questions are explicit.

  evaluation-history:
    title: Evaluation History
    description: >
      Local evaluation runs, reports, recommendations, and quality changelog entries
      that preserve what the project has learned from prior quality work.
    source: ./.quality
    factors:
      reportability:
        title: Reportability
        description: >
          Evaluation records remain complete enough for tools and humans to
          inspect, report, and distinguish history status from source quality.
      traceability:
        title: Traceability
        description: >
          Records, reports, recommendations, quality changelog entries, and feedback
          logs preserve the links between findings, model changes, workflow
          experience, and follow-up decisions.
    requirements:
      evaluation-history-and-model-change-history-remain-inspectable-and-distinct:
        title: evaluation history and model-change history remain inspectable and distinct
        factors:
          - reportability
          - traceability
        assessment: >
          Inspect .quality/evaluations/, .quality/changelog/, .quality/logs/, and
          specs/skills/quality-skill/ to confirm evaluation records, quality-changelog
          entries, and workflow feedback logs are stored separately, preserve
          their intended purpose, avoid secret values, and can be interpreted
          without treating stale or malformed history as evaluated-source quality.
---

# Quality model - QUALITY.md Project

## Overview

This model governs the QUALITY.md project: the open format
([`SPECIFICATION.md`](./SPECIFICATION.md)), the `/quality` agent skill that makes
the format usable through an AI-assistant/coding-agent workflow, the
deterministic `qualitymd` CLI, the docs and specs that explain and constrain the
system, the install and release paths that deliver it, the repository-owned agent
harness, this repository's own `QUALITY.md`, and the local evaluation history.

Good quality here means people and teams using AI assistants or coding agents can
make quality expectations explicit for many kinds of maintained entities, evaluate
evidence against those expectations, learn from the results, and improve the work
without losing the intent behind the model. The modeled domains QUALITY.md can
serve are broad: software, documentation, data sets, research or analytical
reports, services, operations, processes, AI assistants, agent harnesses, and
other entities. The project is not broad in its context of use: the primary
experience is agent- and skill-first, with the CLI as deterministic support
tooling.

The governing sense of good is fitness for purpose first, backed by conformance
where this repository defines a normative contract. The format specification,
skill functional specs, CLI specs, design guides, and authoring guides are
sources of truth for the artifacts they govern.

_Unknowns_ — human ownership/review expectations, real adopter feedback across
domains, support burden, and private roadmap or issue-tracker priorities are not
agent-accessible.
_Open questions_ — what human review cadence should endorse this model after
agent-authored changes?

_Reviewed — not yet human-endorsed; agent-reviewed — Codex (GPT-5), 2026-06-25._

## Scope

This model covers the whole current repository. It includes the format
specification, runtime `/quality` skill, CLI implementation, documentation,
specification bundles, install and release infrastructure, repository-owned
agent harness, this `QUALITY.md`, and local quality/evaluation history.

The model treats the project root as a composite. The root carries the
model-wide agent harnessability factor because the project is explicitly used
through AI assistants and coding agents. Child areas represent distinct
constituents with their own factor families. The agent-harness area is the
checked-in steering and owned-control artifact set; agent harnessability is the
broader capability the whole project exhibits when it equips agents to work
well.

Ratings should be read with a low-tolerance, pre-release posture. Contract
failures in the format spec, mutation/evaluation safety in the skill, automation
contracts in the CLI, and install/update paths can cap their areas even when
surrounding material is strong. Documentation polish can tolerate more
iteration, but not when it misrepresents the domain-agnostic model or hides the
agent-first use context.

Out of scope by design: dependencies and services the project does not own,
including the Go toolchain, package registries, Agent Skills installers,
Homebrew infrastructure, GitHub platform behavior, and third-party libraries.

_Unknowns_ — release-readiness thresholds beyond the visible repo checks are not
fully captured.
_Open questions_ — whether future evaluations should add per-requirement rating
overrides for the strongest safety, conformance, and release-readiness vetoes.

_Reviewed — not yet human-endorsed; agent-reviewed — Codex (GPT-5), 2026-06-25._

## Needs

Primary users need QUALITY.md to remain domain agnostic in what it can model,
while being effective in its intended context of use: AI assistants and coding
agents working with people. They need the model shape to work for software,
documents, data, research, services, operations, processes, and other maintained
entities without importing a default factor checklist.

People and teams using AI assistants or coding agents need `/quality setup`,
`/quality evaluate`, `/quality update`, recommendation follow-up, and direct
`QUALITY.md` edits to feel coherent, bounded, and useful. Agents need
agent-accessible context, scoped tasks, runnable or inspectable verification,
durable evaluation records, and action limits. Maintainers and contributors need
the specs, guides, CLI implementation, skill, and docs to stay aligned as the
project changes.

Format implementers need a portable specification. Skill users need the agent
workflow to foreground judgment and keep deterministic mechanics in the CLI. CLI
users, scripts, and CI need stable command behavior, structured output, and
clear failure modes. Package and install-channel users need release artifacts
that install, update, and report compatibility correctly.

Later evaluation recommendations should default to GitHub Issues for handoff
when external tracking is useful, but setup does not create issues or configure
integrations.

_Unknowns_ — real adopter feedback and support burden across non-software and
software domains are not visible in the repository.
_Open questions_ — which adopter domains should be sampled first when future
examples or evaluations test domain breadth?

_Reviewed — not yet human-endorsed; agent-reviewed — Codex (GPT-5), 2026-06-25._

## Risks

The highest risk is confusing the model domain with the use context. If docs,
examples, or the model imply QUALITY.md is mainly for software practitioners,
the format's broad value is narrowed. If they imply QUALITY.md is inherently for
AI assistant or harness quality, the agentic use context is mistaken for the
default modeled domain.

An ambiguous or incomplete format spec can fragment implementations. A skill
that rates without enough evidence, follows instructions from evaluated content,
hides mutation, or drifts from CLI support can damage trust in the whole quality
loop. A CLI that is inconsistent, interactive by surprise, noisy on the wrong
stream, or stale relative to the skill can break agent and automation workflows.
Install or release failures can leave users with incompatible skill/CLI pairs.

A thin agent harness can make future agent work depend on private memory,
unstated permissions, or non-durable chat context. A stale `QUALITY.md` can make
evaluations optimize the wrong constituents or miss project-wide factors. Stale
or malformed evaluation history can mislead maintainers if it is mistaken for
current evaluated-source quality rather than history status.

_Unknowns_ — the private roadmap and issue-tracker priorities that may drive
near-term release risk are not visible.
_Open questions_ — which risks are hard release blockers for the next public
version?

_Reviewed — not yet human-endorsed; agent-reviewed — Codex (GPT-5), 2026-06-25._

## Model shape

The model is composite by design. `format-spec` covers the normative format
contract. `quality-skill` covers the runtime agent skill and its judgment-bearing
workflow guidance. `cli` covers deterministic command behavior and Go
implementation quality. `docs` covers README, install, guides, examples, and
domain-range communication. `specs-bundle` covers durable functional specs for
tools, records, reports, and skills. `distribution` covers install, packaging,
release, and smoke-check infrastructure. `agent-harness` covers checked-in
steering and owned-control artifacts. `quality-md` evaluates this model itself.
`evaluation-history` covers the local records and logs that preserve what prior
quality work learned.

Agent harnessability is model-wide and decomposed into agent accessibility, task
specifiability, agent operability, continuity, self-verifiability, enforcement
of standards, and containment of action. It is intentionally distinct from the
agent-harness area: the factor asks how the project equips agents to work well;
the area asks whether the owned harness artifacts themselves are complete,
coherent, current, and assessable.

_Unknowns_ — none known.
_Open questions_ — whether distribution and evaluation history should split into
smaller child areas after the next full evaluation produces findings.

_Reviewed — not yet human-endorsed; agent-reviewed — Codex (GPT-5), 2026-06-25._
