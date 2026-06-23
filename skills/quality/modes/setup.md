# Setup Workflow

Run this workflow to create or update a useful first `QUALITY.md`. Setup writes
only the selected `QUALITY.md`; evaluation, quality-log entries, external
issues, recommendation handoff, and recurring-review automation are follow-on
workflows.

## Workflow

```text
Preflight
- verify CLI support
- resolve model file
- emit run frame

Read context
- inspect setup signals
- build setup brief
- identify missing context

Ask discovery questions
- confirm or correct inferred defaults
- use compact prompt by default
- use short sequence when needed

Write QUALITY.md
- scaffold missing file with qualitymd init
- gate existing-file edits with a decision brief
- synthesize body first, then frontmatter

Verify and close
- run qualitymd lint
- inspect readiness with Top 10 checks
- report status and next-step choices
```

## Preflight

1. Verify the CLI prerequisite from `SKILL.md`.
2. Resolve the model file: explicit path when supplied, otherwise `QUALITY.md`
   in the current working directory. Do not walk parent directories.
3. Emit the run frame:

   ```text
   /quality run
   - Mode: setup
   - Model file: <resolved path>
   - Scope: contextual QUALITY.md setup
   - Mutation: QUALITY.md only
   - Artifacts: QUALITY.md
   - Next gate: context analysis, discovery, lint, readiness inspection
   ```

4. Read the authoring sections a first model needs, not the whole guide. For a
   first model read [`../guides/authoring.md`](../guides/authoring.md) sections
   The QUALITY.md file, Quality Model, The Markdown body, and Rating Scale, plus
   the Area and Factor sections enough to shape candidates. Defer deep
   authoring detail until you actually author requirements or ratings, then read
   the relevant sections. Read
   [`../guides/getting-started.md`](../guides/getting-started.md) when setup is
   continuing from a starter/immature model or needs first-run iteration
   guidance.

## Read Context

Inspect available repository context before asking setup questions. Keep this
bounded to setup signals, not source-quality judgment.

Useful signals include README/docs, package metadata, repository structure,
tests/build scripts, contributor docs, existing agent instructions, and visible
work-management or recurring-review hints.

Treat the current directory as the default root area convention unless the user
supplied an explicit model path or context strongly indicates a narrower root.

When the root spans multiple workspaces, packages, or services, delegate a
bounded `Explore` fan-out only when proportional. Use one bounded pass per
component to capture purpose, entry points, external systems, risk surfaces, and
test/CI coverage. Feed the result into candidate Area shape. Do not turn this
into source-quality evaluation.

Build this setup brief as working context:

```text
Root area: <default> (<strongly inferred | weakly inferred | assumed>, <evidence>)
Domain: <default> (<confidence>, <evidence>)
Lifecycle: <default> (<confidence>, <evidence>)
Risk tolerance: <default> (<confidence>, <evidence>)
Modeling rigor: <default> (<confidence>, <evidence>)
Collaboration: <default> (<confidence>, <evidence>)
Primary users/outcomes: <default> (<confidence>, <evidence>)
Maintainers/collaborators: <default> (<confidence>, <evidence>)
Other stakeholders: <default> (<confidence>, <evidence>)
Missing context: <specific gaps>
Review posture: <default when visible> (<confidence>, <evidence>)
Candidate model shape: <factors and child areas>
```

Use only these confidence labels:

- `strongly inferred`
- `weakly inferred`
- `assumed`

The setup brief is not a new artifact. It guides the discovery prompt and then
gets distilled into `QUALITY.md` where the assumptions shape the model.

## Ask Discovery Questions

Ask the user to confirm or correct the setup assumptions before writing
`QUALITY.md`, unless they explicitly ask to accept all inferred defaults. Each
question must include a recommended answer and confidence label.

Default to one compact prompt:

```text
I inspected the repo for setup signals. Please confirm or correct these setup
assumptions before I write QUALITY.md:

1. Root area: Should this QUALITY.md model the whole current project, or a
   narrower area?
   Recommended: <default> (<confidence>)

2. Domain: What kind of thing is this model evaluating?
   Recommended: <default> (<confidence>)

3. Lifecycle: Which stage best fits?
   Options: exploratory, pre-release, active production, maintenance, sunset
   Recommended: <default> (<confidence>)

4. Risk tolerance: How costly is poor quality here?
   Options: high tolerance, moderate tolerance, low tolerance
   Recommended: <default> (<confidence>)

5. Modeling rigor: How detailed should the first quality model be?
   Options: lightweight, standard, high-assurance
   Recommended: <default> (<confidence>)

6. Primary users and outcomes: Who needs the evaluated thing to work, and what
   outcomes matter most?
   Recommended: <default> (<confidence>)

7. Maintainers and collaborators: Who has to change, operate, review, or rely on
   this work?
   Recommended: <default> (<confidence>)

8. Other stakeholders: Are there customers, operators, compliance, support,
   data, security, business, or other stakeholders not visible in the repo?
   Recommended: <default> (<confidence>)

9. Missing context: I think these important inputs are not visible:
   <specific gaps>. What else should the model record as unknown or not
   agent-accessible?

10. Review posture: Should the model record a recurring review expectation?
    Options: none for now, per sprint/iteration, monthly, before major
    releases/planning, custom
    Recommended: <default> (<confidence>)
```

Accept terse corrections. The user does not need to answer every item in prose
when the defaults are right.

Use a short sequence instead of the compact prompt when the interaction surface
needs it. Preserve all ten questions, defaults, confidence labels, and seeded
missing-context content. Prefer these groups:

1. Root area and domain.
2. Lifecycle, risk tolerance, and modeling rigor.
3. Primary users, maintainers/collaborators, and other stakeholders.
4. Missing context and review posture.

The collaboration question assumes agent-heavy development and asks which human
collaborators, reviewers, maintainers, or stakeholders also need to align with
the quality bar.

Review posture is context capture only. Do not ask for permission to create
issues, automations, CI gates, release gates, calendar events, Codex
automations, or Claude Code routines. Ad hoc `/quality evaluate` is always
available, not a selectable automation option.

You may ask an additional work-handoff question when repo context suggests an
issue tracker or handoff process:

```text
Work handoff: If evaluations produce recommendations later, where should
follow-up usually go?
Options: leave in evaluation report, GitHub Issues, Linear/Jira, maintainer
decides each time
Recommended: <default> (<confidence>)

Setup will not create issues or configure integrations.
```

Do not ask users to design factors, child Areas, Requirements, or rating levels
cold. Derive model shape from the setup brief, discovery answers, authoring
guide, and repository context.

## Write QUALITY.md

If no model file exists, run `qualitymd init [path]` after discovery and before
authoring content.

If the model file exists and setup would change it, use a decision brief before
editing:

```text
Decision: update existing `QUALITY.md`?
- Changes:
- Evidence/reason:
- Recommended option:
- Alternatives:
- Done criterion / verification:
```

Synthesize directly into `QUALITY.md`. Author the body first, then the
frontmatter model:

- Overview and Scope establish the root area, domain, lifecycle, risk
  tolerance, modeling rigor, and key boundaries.
- Needs and Risks capture primary user needs, maintainer/collaborator needs,
  other stakeholder needs, and the failure modes that matter.
- Unknowns and open questions capture missing or non-agent-accessible context.
- Review posture is recorded only when it affects how the model should be used.
- The rating scale stays standard unless the body shows a real mismatch.
- Factors and child Areas derive from project needs, risks, stakeholder
  concerns, component boundaries, and available evidence.
- Requirements are small, concrete, and assessable from agent-accessible
  evidence, or explicitly name missing evidence.
- Include the `quality-md` self-check Area when appropriate unless the user
  declines or the model file is not in the root area it governs. Use
  `quality-md` as the Area key, `<Root Title> QUALITY.md` as the Area title, an
  Area `description`, and an explicit path-based `source` for the model file
  such as `./QUALITY.md`; do not use prose aliases such as `(this file)` for
  `source`. Add concise YAML comments around that Area explaining that `source`
  is the `QUALITY.md` artifact being evaluated, while the Requirement's
  `assessment` references the guide used to judge it. Prefer one Area-level
  Requirement that cites the active authoring guide once and lists each affected
  Factor under `factors` when that guide defines one coherent judgment across
  the Factors.

## Verify and Close

Run `qualitymd lint [path]`. Stop on lint errors, report the CLI findings, and
route to continued `QUALITY.md` iteration. Do not recommend evaluation while the
model is invalid.

When lint passes, inspect the resulting model with
[`../guides/top-10-quality-md-checks.md`](../guides/top-10-quality-md-checks.md).
This is a bounded model-readiness inspection, not a project evaluation. Classify
readiness as `starter`, `immature`, or `ready to evaluate`.

Report setup completion status-first:

```text
Setup complete
- Changed: QUALITY.md
- Validation: lint passed | lint failed
- Readiness: starter | immature | ready to evaluate
- Important gaps: <none | concise model gaps>
- Not done: no evaluation, no quality log, no issues, no automations
- Next: continue iterating | run evaluation | set up recurring review | set up recommendation handoff | stop here
```

If readiness is not `ready to evaluate`, list the most important model gaps and
make continued iteration the recommended next step. Do not automatically take
any next-step action.

Setup creates or updates a useful first model; it does not invent a complete
quality model without user/project context, run an evaluation, write the quality
log, create external issues, or configure automation.
