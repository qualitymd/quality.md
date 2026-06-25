---
coverage:
  assessmentResults:
    - areaPath: []
      requirement: "agents can reach the minimum project context and deeper routed guidance"
    - areaPath: []
      requirement: "project guidance makes task boundaries and done criteria explicit"
    - areaPath: []
      requirement: "agents can discover the tools and commands needed to work the project"
    - areaPath: []
      requirement: "quality workflows preserve useful state in durable local artifacts"
    - areaPath: []
      requirement: "verification signals are discoverable and remediation-bearing"
    - areaPath: []
      requirement: "core project standards are backed by enforceable checks or reviewable gates"
    - areaPath: []
      requirement: "agent action limits and mutation gates are visible before consequential changes"
    - areaPath: [format-spec]
      requirement: "the format specification admits a single interpretation"
    - areaPath: [format-spec]
      requirement: "the format specification is internally consistent"
    - areaPath: [format-spec]
      requirement: "the format specification is complete enough to implement and author from"
    - areaPath: [format-spec]
      requirement: "each normative rule is observable or testable"
    - areaPath: [format-spec]
      requirement: "the specification keeps model domain separate from agentic use context"
    - areaPath: [quality-skill]
      requirement: "the skill grounds setup, evaluation, and follow-up in the active model and evidence"
    - areaPath: [quality-skill]
      requirement: "the skill gates mutating actions with explicit decision briefs"
    - areaPath: [quality-skill]
      requirement: "the skill follows the agent-mediated UX guide"
    - areaPath: [quality-skill]
      requirement: "the skill workflows match their functional specs and CLI support surface"
    - areaPath: [cli]
      requirement: "the CLI follows its functional specifications"
    - areaPath: [cli]
      requirement: "CLI commands expose stable machine-readable behavior where agents need it"
    - areaPath: [cli]
      requirement: "the CLI follows the project CLI design guide"
    - areaPath: [cli]
      requirement: "the Go implementation follows the project Go style guide"
    - areaPath: [docs]
      requirement: "introductory docs foreground the agent-first workflow without making the CLI the main surface"
    - areaPath: [docs]
      requirement: "docs keep modeled domains broad while preserving the agentic use context"
    - areaPath: [docs]
      requirement: "docs reflect the current command, install, release, and workflow surfaces"
    - areaPath: [specs-bundle]
      requirement: "tooling specs identify the behavior they govern and stay linked to runtime artifacts"
    - areaPath: [specs-bundle]
      requirement: "the specs bundle remains internally consistent and synchronized with implementation"
    - areaPath: [distribution]
      requirement: "supported install and update paths are documented, tested, and package the expected binary"
    - areaPath: [distribution]
      requirement: "release gates verify the artifacts and channels that users depend on"
    - areaPath: [agent-harness]
      requirement: "the agent harness orients agents and routes them to deeper guidance"
    - areaPath: [quality-md]
      requirement: "the QUALITY.md model follows the active authoring guide family"
    - areaPath: [evaluation-history]
      requirement: "evaluation history and model-change history remain inspectable and distinct"
  analyses:
    - areaPath: []
    - areaPath: [format-spec]
    - areaPath: [quality-skill]
    - areaPath: [cli]
    - areaPath: [docs]
    - areaPath: [specs-bundle]
    - areaPath: [distribution]
    - areaPath: [agent-harness]
    - areaPath: [quality-md]
    - areaPath: [evaluation-history]
---

# Evaluation plan

## Rigor

Standard.

## Requirement Set

Assess all 30 in-scope Requirements from the current `model.md` snapshot:

- Root Agent Harnessability: 7 Requirements across agent accessibility, task
  specifiability, agent operability, continuity, self-verifiability, enforcement
  of standards, and containment of action.
- Format Specification: 5 Requirements.
- `/quality` Skill: 4 Requirements.
- qualitymd CLI: 4 Requirements.
- Documentation and Examples: 3 Requirements.
- Tooling Specs Bundle: 2 Requirements.
- Installation and Distribution: 2 Requirements.
- Agent Harness: 1 Requirement.
- QUALITY.md Project QUALITY.md: 1 Requirement.
- Evaluation History: 1 Requirement.

## Evidence Strategy

Use targeted repository evidence for each Area:

- Read the declared source path for each Area and the referenced assessment
  sources named in the Requirement.
- Use `rg` searches for terminology, command names, record fields, workflow
  names, and documented contracts.
- Run local commands where behavior is rating-relevant: `qualitymd lint`,
  `qualitymd status --json`, `qualitymd evaluation status`, command `--help`,
  `mise run check`, and focused packaging/schema checks as needed.
- Treat external distribution channels and private roadmap/support context as
  limitations unless represented by repository-visible scripts, workflows, or
  docs.

## Planned Limitations

This run will not publish releases, install from remote registries, inspect
private issue trackers, or run full cross-platform smoke tests. Findings about
release and install quality are based on checked-in scripts, workflows, package
metadata, and local command behavior.
