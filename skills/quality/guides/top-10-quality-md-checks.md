---
type: Runtime Guide
title: Top 10 QUALITY.md Checks
description: Quick model-readiness and routing checks for QUALITY.md files.
---

# Top 10 QUALITY.md Checks

Use this checklist for a fast, read-only inspection of a QUALITY.md file's
lifecycle state and model maturity. The result is not an Evaluation Report and
does not rate the evaluated source. It produces routing findings that read-only
orientation and model-review workflows can use to choose the next public
workflow.

Two distinct axes come out of this inspection, and they must not be conflated:

- *Lifecycle state* is where the model sits in the evaluation lifecycle and is
  owned by `qualitymd status` (`readiness`): missing, invalid, ready to evaluate
  (valid, no runs yet), has evaluation history, or needs reconciliation.
- *Maturity* is how developed the model is, and is this checklist's judgment:
  `starter`, `immature`, or `evaluation-ready`.

A model can be lifecycle `ready-to-evaluate` (valid, no runs) while still being a
`starter` in maturity. Report both; do not read the CLI's lifecycle signal as a
maturity verdict.

This checklist does not re-run setup. It checks whether the current
`QUALITY.md` still preserves the setup assumptions and model qualities needed
for useful evaluation, authoring, and maintenance.

## Inputs

- `qualitymd status [path] --json`
- the area `QUALITY.md`
- evaluation history summarized by status JSON

Do not inspect evaluated source files for this checklist. Do not read evaluation
report bodies. Keep the inspection bounded to the model file and status signals.

Do not require lifecycle, risk tolerance, modeling rigor, collaboration context,
stakeholder needs, or quality-loop posture to appear in fixed sections. Treat
them as present when they are explicit, current, and usable anywhere in the
Markdown body or model context.

## Finding shape

Report only findings that affect routing or model usefulness. Use this shape:

```text
QUALITY.md inspection findings
- <check id>: <finding>
  Evidence: <status field, section, or property>
  Impact: <why this affects lifecycle/model usefulness>
  Route: <setup | getting-started | authoring | evaluate | recommendation follow-up | history | update>
```

Keep evidence short. Cite section names, property paths, counts, or status JSON
fields rather than quoting long passages.

## Condensed close checklist

Setup's close step uses this condensed pass to classify maturity. Read the full
checks below only when the maturity call is borderline or a read-only
orientation needs the complete inspection.

- Lint passes and `qualitymd status` reports a valid model. For a freshly
  authored model, lifecycle is `ready-to-evaluate` (valid, no runs). (check 1)
- Overview, Scope, Needs, and Risks are substantive, each closing with its own
  unknowns and open questions, and important missing context is named. (check 5)
- Project posture — lifecycle, risk tolerance, modeling rigor — and the primary
  stakeholder needs are explicit and specific to this project, not generic.
  (checks 2-4)
- Factors map to the body's needs and risks, and requirements are concrete with a
  usable assessment each. (checks 8-9)
- For a composite root, the domain's constituent kinds (e.g. specs, tests, docs,
  operations, security/safety) are each modeled or consciously accounted for, not
  silently missing. (check 8)
- For an agent-collaborated composite root, Agent Harnessability or its
  sub-factors are present as model-wide factors unless clearly not germane. (check
  8)

Classify maturity from how many of these hold:

- `starter` — scaffold or thin body; posture or needs mostly generic; few
  concrete requirements.
- `immature` — substantive in parts, but with material gaps in posture, needs,
  factors, or assessable requirements.
- `evaluation-ready` — posture, needs, factors, and assessable requirements are
  present and specific enough to evaluate.

## The Checks

### 1. Model lifecycle state

Use `qualitymd status --json` to identify whether the model is missing, invalid,
valid with no history, valid with history, or needs evaluation reconciliation.

- Finding when missing or invalid: route to setup or lint repair before any model
  quality judgment.
- Finding when history needs reconciliation: route to history/reconciliation
  before a fresh evaluation unless the user explicitly wants a new run.

### 2. Project Posture

Inspect whether the model captures the project context that calibrates the
quality bar: lifecycle, risk tolerance, and intended modeling rigor.

- Finding when lifecycle is absent, stale, or contradicted by the body: route to
  authoring.
- Finding when risk tolerance is unclear enough that requirements cannot tell
  acceptable gaps from unacceptable gaps: route to authoring.
- Finding when the model is too thin or too heavy for its stated modeling rigor:
  route to authoring.
- Finding when production, maintenance, or sunset posture is stated but not
  reflected in factors or requirements: route to authoring.

### 3. Stakeholder and Needs Coverage

Inspect whether the model makes the relevant stakeholder needs visible enough to
justify the factors and requirements. Consider primary users,
collaborators/maintainers, and other affected stakeholders.

- Finding when primary users or user outcomes are unclear: route to
  getting-started for starter models or authoring for populated models.
- Finding when collaborator or maintainer needs are absent despite being central
  to the project's quality: route to authoring.
- Finding when other stakeholders are implied but their needs are not stated:
  route to authoring.
- Finding when needs are generic enough that the same text could fit almost any
  project: route to authoring.

### 4. Agent and Collaboration Fit

Inspect whether the model supports the assumed agent-heavy workflow plus the
named human collaboration context.

- Finding when future agents would need private memory, unavailable tools, or
  uncited context to understand or apply the model: route to authoring.
- Finding when the collaboration context is unclear enough to leave review,
  onboarding, governance, or handoff expectations implicit: route to authoring.
- Finding when open source, cross-functional, customer-facing, or external
  contributor collaboration is implied but not reflected in factors,
  requirements, or body context: route to authoring.

### 5. Body Context and Missing Context

Inspect whether the Markdown body gives enough evaluable judgment context to
build, use, and evaluate the model: Overview, Scope, Needs, and Risks should be
present and substantive, each closing with its own unknowns and open questions.
Important missing or non-agent-accessible context should be explicit rather than
invisible.

- Finding when body sections are empty, placeholder-like, or generic in a starter
  model: route to getting-started. Use the authoring guide as the quality
  reference for what the body should accomplish.
- Finding when a section omits its unknowns or open questions while the rest of
  the body leaves unresolved questions: route to getting-started for first-run
  process or authoring for best-practice guidance.
- Finding when material support is referenced or implied but is not
  agent-accessible, and the gap prevents a reader or agent from judging whether
  the body is complete, current, grounded, or sufficient: route to authoring.
- Finding when important missing context is not named, even though the model
  depends on it for scope, needs, risks, or assessment evidence: route to
  authoring.

### 6. Root area and scope alignment

Check whether the root title, body scope, file location, and root or child
`source` values describe the same evaluated root area. The current directory is
the default root area convention unless the model clearly narrows or relocates
scope.

- Finding when the title names the repository but the body/sources are narrower:
  route to getting-started or authoring.
- Finding when source coverage includes unrelated/generated/supporting artifacts:
  route to authoring.
- Finding when the model overrides the current-directory convention without
  explaining the evaluated boundary: route to authoring.
- Finding when exclusions or boundary decisions are important but implicit:
  route to authoring.

### 7. Rating scale fit

Check whether the rating scale is understandable and fits the body's decision
context, including lifecycle, risk tolerance, and modeling rigor.

- Finding when level descriptions or criteria are generic enough that findings
  cannot distinguish `target` from `minimum`: route to authoring.
- Finding when a custom scale exists but the body does not explain why: route to
  authoring.
- Finding when the scale implies a stricter or looser bar than the stated project
  posture: route to authoring.

### 8. Area and factor shape

Check whether the area tree is small enough to understand, specific enough to
represent distinct evaluated entities, and shaped by the body's needs and risks.
Factors should be meaningful quality lenses, not vague labels alone.

- Finding when all concerns are flattened into the root despite clear sub-entities
  in the body: route to authoring.
- Finding when the body describes distinct constituent artifacts of different
  kinds (e.g. the running artifact, its requirements, its docs) but all factors
  are held at the root as one family — a composite entity flattened into a single
  primary-subject root: route to authoring.
- Finding when an agent-collaborated project's body shows an owned, high-leverage
  agent harness or QUALITY.md self-check that is not modeled as a constituent:
  route to authoring. These are modeled by default given the context of use; do
  not flag a harness-less or throwaway project, which hits the not-germane
  disqualifier.
- Finding when an agent-collaborated composite root does not carry
  Agent Harnessability or its sub-factors among its model-wide factors: route to
  authoring. Do not flag a non-agent-collaborated, harness-less, throwaway, or
  narrowly scoped entity where the factor is not germane. Treat an existing
  `harnessability` factor with the legacy six-sub-factor shape as semantic
  coverage, but route model-authoring work to rename it to
  `agent-harnessability` / Agent Harnessability and add `continuity` unless the
  project has an explicit reason to preserve the old key. A thin or absent
  harness is rating evidence, not a reason to omit the factor.
- Finding when an agent-collaborated composite root has a germane agent-harness
  area carried with only one or two thin factors: route to authoring. Treat this
  as a coverage gap in the steering-materials area, not as evidence the harness is
  unimportant.
- Finding when an agent-harness area is defined or scoped as instructions only,
  omitting the feedback half of the harness or project-owned runtime controls:
  route to authoring. Define the harness as the whole engineered system around
  the model, then scope the area to checked-in steering and owned-control
  artifacts.
- Finding when project-owned runtime harness machinery is present in the repo
  (hooks, tool/MCP definitions, sandbox or permission policy, orchestration
  config, subagent config) but is neither modeled in the agent-harness area, given
  its own area, nor explicitly out of scope: route to authoring.
- Finding when agent-harness area requirements assume a software toolchain
  (lint/type-check/test/CI/deploy) but the project's served domain is not
  software: route to authoring. Rephrase toward how this project's steering
  materials point to verification, enforce standards, or bound action.
- Finding when a model carries two same-rooted projections of one concern (e.g. an
  Agent Harnessability factor and an agent-harness area) with no boundary note —
  neither a YAML comment nor a disambiguating description clause distinguishing
  them: route to authoring. A reader cannot tell the projections apart, and the
  double-count risk is invisible.
- Finding when the domain implies a germane constituent kind the body evidences or
  implies — for software product quality: tests, documentation modes
  (tutorial/how-to/reference/explanation), specs/requirements, operations, or a
  security/safety artifact; for a data product: schema, provenance, freshness, or
  lineage metadata — that the model neither models as an area nor surfaces
  as a ratable gap (a missing-anchor area, or a requirement on an existing area):
  route to authoring. For a germane kind, a bare deferral or Scope note does not
  satisfy coverage. Do not flag a kind that legitimately hits a disqualifier —
  folded into a sibling for lack of distinct concerns, or genuinely not germane
  (and so marked out of Scope) for a throwaway or narrowly scoped entity.
- Finding when child areas merely mirror the parent without distinct factors or
  requirements: route to authoring.
- Finding when major body needs/risks have no factor: route to authoring.
- Finding when factors are generic, overlapping, or unexplained: route to
  authoring.

Model every germane constituent as its own area by default; a constituent is left
unmodeled only when it hits a disqualifier (no distinct concerns → fold; not
germane → out of scope), and a germane concern with a thin or missing artifact is
surfaced as a ratable gap, never recorded only in prose.

### 9. Requirement and assessment quality

Check whether requirements are concrete enough to produce findings and ratings,
and whether each `assessment` gives the evaluator a usable means of assessment,
either inline or by referencing a traceable entity that defines it.

- Finding when requirements are aspirations rather than assessable expectations:
  route to authoring.
- Finding when a requirement lacks observable evidence or criteria: route to
  authoring before evaluation.
- Finding when assessments are placeholders, circular, or vague: route to
  authoring.
- Finding when referenced assessment sources are not traceable from the model:
  route to authoring.
- Finding when evidence or criteria cannot distinguish adjacent rating levels:
  route to authoring before evaluation.

### 10. Quality Loop Maintenance Signals

Use evaluation history, active recommendations, and visible model context to
decide whether the next workflow is maintenance rather than new authoring or
evaluation. When a model depends on recurring quality review or recommendation
handoff, that posture should be visible, but the checklist does not require or
recommend CI or release gating by default.

- Finding when active recommendations exist: route to recommendation review,
  apply, or issue-tracker handoff.
- Finding when the latest run is stale, incomplete, malformed, or unreportable:
  route to history/reconciliation. Treat malformed or incompatible records as
  history status, not evaluated-source quality evidence; do not suggest manual
  migration.
- Finding when the model implies a recurring quality review cadence but does not
  name how maintainers expect to revisit it: route to authoring.
- Finding when recommendation handoff is central to the workflow but the handoff
  destination is unknown or stale: route to recommendation follow-up or
  authoring, depending on whether active recommendations exist.

## Summary Judgment

After the checks, report the two axes separately.

Lifecycle state (from `qualitymd status`):

- `missing`
- `invalid`
- `ready to evaluate`
- `has evaluation history`
- `needs reconciliation`

Maturity (this checklist's judgment, for a valid model):

- `starter`
- `immature`
- `evaluation-ready`

Use the finding routes to recommend one next workflow and list a few concrete
alternatives.
