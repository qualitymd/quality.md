---
type: Design Doc
title: Implement the /quality skill — design doc
description: How the /quality skill is packaged for Agent Skills installation, how setup verifies the qualitymd CLI prerequisite, how the qualitymd models CLI surface and raw JSON evaluation artifacts are built, and how the open items resolve into the durable spec.
tags: [skill, quality, evaluation, design]
timestamp: 2026-06-17T00:00:00Z
---

# Implement the /quality skill — design doc

How the [Implement the /quality skill](../0010-implement-quality-skill.md) change
is built — the technical approach behind its
[functional spec](spec.md). Where the spec says _what_ must hold (and defers the
behavioral contract to the durable
[`/quality` skill spec](../../../specs/skills/quality-skill/quality-skill.md)), this
doc says _how_ the implementation makes it so, and why that way.

## Context

The skill is fully specified but unbuilt: the
[durable skill spec](../../../specs/skills/quality-skill/quality-skill.md) owns the
operating model, invocation, evaluation workflow, and reporting contract, and a
[worked example bundle](../../../specs/skills/quality-skill/examples/index.md) shows
its output. This change packages an invocable skill that **conforms to** that
contract, adds the one CLI surface the contract leans on (`qualitymd models`), and
settles the [open items and gaps](spec.md#open-items-and-gaps) the durable spec
left — syncing the durable spec to the resolutions before **Done**.

Three things shape every decision below:

- **Judgment in the skill; determinism in the CLI.** The skill never reimplements
  scaffolding, structural validation, or the format rules — it drives the
  [`qualitymd` CLI](../../../specs/cli.md) for each and treats its output as the
  source of truth. Anything mechanical that the skill needs and the CLI lacks is
  CLI work this change must also deliver (that is `qualitymd models`).
- **The durable spec is the conformance target, not a script.** The skill is one
  _implementation_ of an evaluator; this doc is free to choose concrete artifacts,
  schemas, and ordering so long as the result satisfies the contract.
- **The open items already carry recommended resolutions.** This doc confirms the
  three [blocking](spec.md#blocking--resolve-before-design) ones (their resolution
  is a precondition for **Design**) and works out the engineering behind them; the
  rest are noted where they affect the build and land in the durable spec during
  **In-Progress**.

## Approach

### 1. Inputs this design builds on

The blocking items resolve as the spec recommends, confirmed here:

- **Item 1 — packaging.** Source home `skills/quality/`, settled. The deferred
  Design sub-decision — _distribution and installation_ — is decided in
  [§2](#2-packaging-for-agent-skills): the repo is an Agent Skills source installed
  with `npx skills add qualitymd/quality.md`; Claude plugin packaging is a
  possible secondary channel, not the main onboarding path.
- **Item 2 — the `model` altitude's criteria.** Meta-evaluation is ordinary
  evaluation with the bundled
  [quality meta-model](../../../internal/models/quality-meta-model.md)
  as the active model and the user's `QUALITY.md` as its subject, reached through a
  new `qualitymd models` surface ([§4](#4-the-qualitymd-models-command)).
- **Item 3 — `setup`.** A minimal bootstrap after skill installation — verify or
  repair the CLI prerequisite, drive `init`, validate with `lint`, then hand to
  `wizard` — owned by the skill prompt ([§3](#3-the-quality-skill-artifact),
  [§6](#6-driving-the-cli-and-resolving-inputs)).

The spec-affecting items (4, 5, 7) and the durable-spec corrections (8–11) are
handled in [§5](#5-evaluation-artifacts), [§6](#6-driving-the-cli-and-resolving-inputs),
and the [sync plan](#7-durable-spec-sync-in-progress).

### 2. Packaging for Agent Skills

The repository is the skill source users install first. The primary onboarding
path mirrors Basecamp-style skill repositories:

```sh
npx skills add qualitymd/quality.md
```

The CLI remains a prerequisite, not the installer for the skill. The installed
skill's `setup`/`wizard` flow detects whether that prerequisite is present and
compatible, then facilitates install or upgrade when needed (§6).

```
quality.md/                       # repo root = Agent Skills source
  skills/
    quality/
      SKILL.md                    # the authored skill artifact (item 1 source home)
  install.md                      # agent-readable full setup guide
  cmd/ internal/ specs/ …         # CLI, durable specs, and docs
```

- **Repository shape.** `skills/quality/SKILL.md` is the canonical distributable
  artifact. The top-level `skills/` view is intentionally simple so generic Agent
  Skills tooling can detect and copy it for Codex, Claude Code, Cursor, Gemini CLI,
  and other agents without a Claude-specific manifest.
- **Local dogfood install.** Contributors should be able to install the skill from
  the working tree during development (for example, `npx skills add .`, subject to
  the installer-supported local path syntax) so dogfooding does not require pushing
  to GitHub first. The published install command remains
  `npx skills add qualitymd/quality.md`.
- **Install guide.** A root `install.md` gives agents the full bootstrap objective:
  install or upgrade the `qualitymd` CLI, install the skill with
  `npx skills add qualitymd/quality.md`, verify both, and restart the agent session
  if the target agent requires it. README points to the same path but stays shorter.
- **Invocation.** The skill artifact is named `quality`, so the durable spec's
  `/quality` remains the canonical spelling. Some agents may expose installed
  skills differently; docs note that agent-specific invocation syntax is an adapter
  detail, while the skill contract stays `quality`.
- **Secondary channels.** Claude plugin or marketplace packaging can be added later
  if useful, but it is not required by this change and should not force a namespaced
  invocation or lock the skill release cadence to the CLI binary.

### 3. The `/quality` skill artifact

`skills/quality/SKILL.md` is a single prompt-only skill (bundled `references/`
assets stay deferred per the durable spec). Its frontmatter makes it invocable and
self-describing; its body carries the **evaluation process** the skill owns.

- **Frontmatter.** `name` and a trigger-oriented `description`. Argument hints,
  invocation parsing, and tool guidance live in the body rather than additional
  frontmatter fields, matching current Agent Skills authoring guidance. The
  description is deliberately broader than literal `QUALITY.md` mentions: it should
  trigger for quality management, setup, wizard, evaluation, and improvement
  requests such as "improve security quality" or "evaluate this component's
  reliability characteristic," then route through `wizard`/`setup` when no model
  exists. It uses broad quality vocabulary (factors, characteristics, attributes,
  criteria) plus `QUALITY.md` terms (Targets, Factors, Requirements — not ISO or
  implementation terms), frames the subject as a project/entity or
  component/target instead of a closed list of types, mentions subject evaluation
  and model evaluation/improvement, avoids generic copyediting triggers, and leaves
  CLI implementation details to the skill body. Final description text:
  ```yaml
  description: Use when a user wants setup, wizard guidance, evaluation, or improvement for quality management of a project/entity or one of its components/targets. Trigger for requests about quality factors, characteristics, attributes, criteria, Targets, Factors, Requirements, improving a quality factor such as security/reliability/usability, evaluating a subject against quality criteria, or evaluating/improving the QUALITY.md model itself.
  ```
  Model-invocation stays **enabled** (no
  `disable-model-invocation`) so "evaluate this against its `QUALITY.md`" can reach
  the skill, while the read-only [`wizard`](../../../specs/skills/quality-skill/quality-skill.md#wizard)
  default and the confirm-before-apply rule keep that safe.
- **Description criteria.** The text above follows criteria that must be carried
  into the durable skill spec's Frontmatter and metadata section: optimize for
  trigger matching rather than documentation; include the supported mode keywords
  (`setup`, `wizard`, `evaluate`, `improve`); include broad quality vocabulary
  users naturally ask with (`quality management`, `quality evaluation`, `quality
improvement`, factors, characteristics, attributes, criteria); include
  `QUALITY.md` vocabulary (Targets, Factors, Requirements) for users who know the
  format; frame the assessed thing broadly as a project/entity and scoped
  component/target rather than enumerating a closed list of subject types; include
  subject evaluation and model evaluation/improvement; exclude CLI implementation
  detail from the metadata; and avoid triggering for generic copyediting or one-off
  "make this higher quality" requests that lack systematic quality criteria or
  assessment.
- **Prerequisite contract.** The skill body includes a short Prerequisites section
  naming the minimum compatible `qualitymd` CLI version. During implementation that
  minimum is set to the first release/tag that contains every CLI surface this skill
  depends on, including `qualitymd models`.
- **Arguments.** The skill takes the raw `$ARGUMENTS` string and resolves the five
  parameters itself, each with a default, so a bare invocation is valid. The grammar
  is mixed (positional `evaluate`/`improve`/`setup`/`wizard` and `model`, a scope
  name, a `--effort` value, a path); a free-form parse fits it better than fixed
  positional `arguments:` slots. A bare name resolves against the already-grounded
  model to a target or factor, with `target`/`factor` keywords to disambiguate.
- **Tools.** The skill body names the read-and-drive surface the workflow always
  needs — `qualitymd` invocations, CLI-presence/version checks, file reads/searches,
  and writes under the resolved evaluation directory. It does not rely on
  frontmatter tool grants. Install or upgrade commands for a missing/stale CLI use
  normal agent permissioning and platform detection rather than being silently
  pre-granted. The skill deliberately does **not** pre-grant edits to the _subject_
  or the `QUALITY.md`; under `improve` those go through normal permissioning after
  the explicit confirmation the spec requires, so an apply is never silent.
- **Body.** The process the skill owns, carried in the prompt: the operating model
  and boundaries, the [workflow](../../../specs/skills/quality-skill/quality-skill.md#workflow),
  effort levels, and the artifact contract. It **MUST NOT** embed the format/schema
  rules or rating vocabulary — those are grounded at runtime from `qualitymd spec`
  (§6) — keeping the prompt from drifting out of sync with
  [`SPECIFICATION.md`](../../../SPECIFICATION.md).
- **Modes.** One skill branches on mode. `setup` (item 3) is the bootstrap path:
  confirm `qualitymd` is installed and compatible, facilitate install/upgrade when
  needed, drive [`init`](../../../specs/cli/init.md), validate with
  [`lint`](../../../specs/cli/lint.md), then hand to `wizard` for guided population —
  reimplementing none of scaffolding, validation, CLI installation tooling, or
  authoring judgment. `wizard` performs the same prerequisite check before any
  CLI-dependent action, then owns the cursory repository assessment: read the model
  outline, infer and confirm targets/factors, ask clarifying questions when key
  model inputs are missing, and route to a concrete next `/quality …` action.

### 4. The `qualitymd models` command

A new deterministic surface emits the **bundled** models the skill grounds the
`model` altitude in. It is the mechanical half of item 2: the skill drives it
rather than carrying meta-model criteria in its prompt.

- **Surface.** Two subcommands, mirroring existing conventions:
  - `qualitymd models list` — the catalog. Human table by default; under `--json`,
    a stable array of `{ name, title, description }` so the skill can discover what
    is available rather than embedding a list.
  - `qualitymd models view <name>` — emits one bundled model. By default it follows
    the same Markdown rendering split as [`qualitymd spec`](../../../specs/cli/spec.md):
    formatted Markdown through the terminal renderer and pager when stdout is a
    TTY, and **verbatim Markdown** byte-for-byte when output must be plain (redirect,
    pipe, or `NO_COLOR`). With `--json`, it emits a stable JSON document for
    agents/tools. A
    `--source <path>` flag re-points the emitted model's apex `source` to the
    user's file (item 2's "re-point at runtime"), a deterministic single-node
    rewrite reusing the frontmatter-edit primitive
    [`lint --fix`](../../../specs/cli/lint.md#repair-behavior) already has. Without it,
    the model is emitted as authored. The `--source` rewrite applies before either
    Markdown or JSON rendering, so the two forms describe the same active model.
    Representative JSON shape:
    ```json
    {
      "schemaVersion": 1,
      "name": "quality-meta-model",
      "title": "QUALITY.md Meta-Model",
      "description": "Criteria for evaluating a QUALITY.md model.",
      "model": { "source": "QUALITY.md", "ratingScale": [], "targets": {} },
      "bodyMarkdown": "\n# QUALITY.md Meta-Model\n..."
    }
    ```
- **Catalog contents.** One entry today — `quality-meta-model` — with future sample
  models exposed through the same catalog. `models` is the **bundled-model
  catalog**; future commands that inspect a user's local `QUALITY.md` should be
  top-level verbs (for example, `qualitymd outline`) rather than a nearby singular
  `qualitymd model` namespace.
- **Go shape.** A new `internal/models` package owns the bundled assets and catalog:
  `//go:embed` the meta-model Markdown (relocated under the package, since `go:embed`
  cannot reach a sibling directory and the file has no Go consumer today), exposing
  `Catalog()`, `Get(name)`, a source-rewrite helper, and a structured view that
  parses through `internal/document` and decodes through `internal/model`.
  `internal/cli/models.go` adds `newModelsCmd()` (with `list`/`view` children),
  registered in [`root.go`](../../../internal/cli/root.go) beside `init`/`lint`/`spec`.
  It inherits the CLI baseline — determinism, stdout-is-payload, exit categories,
  and `--json` where meaningful.
- **How the `model` altitude uses it.** The skill runs
  `qualitymd models view quality-meta-model --source <resolved target file>` and
  writes the result as the run's `model.md` snapshot; the meta-model is the active
  model, the user's `QUALITY.md` its subject. No invented standard, and the criteria
  stay the already-maintained meta-model.

### 5. Evaluation artifacts

The skill writes a numbered run folder into the **evaluated** repository's
configured evaluation directory. These are **raw runtime outputs, not OKF
concepts** (item 6): none carries OKF frontmatter or a registered type. Form
follows consumer.

```
quality/evaluations/                    # default parent; configurable via .quality/config.yaml
  0001-subject-quality-eval/        # NNNN-<altitude>[-<narrowing>]-quality-eval (item 9)
    model.md                        # verbatim model snapshot (Markdown)
    design.md  plan.md              # inputs, method (Markdown)
    assessments/001-<target>-<requirement>.json   # source-of-record data (JSON)
    analysis/<target>.json                         # source-of-record data (JSON)
    report.md   report.json         # one result, two renderings
    recommendations/001-<slug>.md   # triageable units (Markdown)
```

- **Evaluation directory config.** `.quality/config.yaml` is the repository-local
  qualitymd system config. The skill reads it before choosing `NNNN`; future CLI
  evaluation commands should use the same surface rather than inventing a parallel
  config file. Supported shape for this change:
  ```yaml
  evaluationDir: quality/evaluations
  ```
  `evaluationDir` is the parent directory that contains numbered run folders. It is
  repository-relative, normalized before use, and invalid if absolute or if it
  escapes the repository. Missing config or missing `evaluationDir` means
  `quality/evaluations/`. Unknown keys are surfaced as warnings and ignored, so a
  typo is visible without making future configuration impossible.
- **Naming (item 9).** Altitude-first and deterministic:
  `NNNN-<altitude>[-<narrowing>]-quality-eval`, `<altitude>` ∈ `subject`/`model`,
  `<narrowing>` the scoped target/factor slug (omitted whole-model). The example
  re-slugs to `0001-subject-quality-eval`. `NNNN` is the next integer the skill
  derives by scanning the existing folder — a mechanical, deterministic step.
- **Why JSON for assessments/analysis.** They are the **source of record** for
  Assess-and-Rate and Analyze; a gate or tool reads them directly, so structured
  data beats prose. `report.md` (human) and `report.json` (machine) are **two
  renderings of one result**, both emitted every run, so no output-format argument
  threads through the skill. `model.md`, `design.md`, `plan.md`, and recommendations
  stay Markdown — snapshot, narrative inputs, and triage units a person reads.
- **JSON shape principle.** JSON files use stable, generic top-level fields tied to
  the evaluation workflow, not fields invented for one factor or requirement.
  Domain-specific details live under `attributes` on the smallest relevant object.
  That keeps the schema useful for gates while allowing a security finding, docs
  gap, test failure, missing evidence, or prompt-injection observation to carry its
  own metadata without changing the public shape.
- **Record schemas (schemaVersion'd; final field names settle with the example
  re-capture).** Rating levels are stored as the model's own `level` slugs (grounded,
  not hard-coded). Representative shapes:

  _assessment_ — one per in-scope requirement, the assess→finding→rating chain:

  ```json
  {
    "schemaVersion": 1,
    "target": "Sparrow Payments API", "targetPath": ["root"],
    "requirement": "No credentials are committed to the repository",
    "factors": ["Security", "Secrets handling"],
    "rating": "unacceptable", "notAssessed": false,
    "criterionSource": "ratings-override",
    "findings": [
      { "locator": "internal/gateway/client.go:48",
        "observation": "Live gateway secret key committed in plaintext; matches active-key format.",
        "category": "secret", "severity": "critical",
        "evidence": [
          { "kind": "source", "ref": "internal/gateway/client.go:48" }
        ],
        "attributes": { "credentialType": "gateway secret key" } }
    ],
    "rationale": "…single live secret lands here regardless of clean findings elsewhere.",
    "recommendations": ["001-rotate-committed-gateway-key"]
  }
  ```

  A _not assessed_ requirement sets `rating: null`, `notAssessed: true`, and states
  the absent evidence in `rationale`. Each finding uses the same generic top-level
  fields: a primary `locator`, the observed fact, an open `category`, optional
  `severity`, supporting `evidence`, and optional `attributes` for category-specific
  metadata. A secret finding may set `category: "secret"` and
  `attributes.credentialType`, but **never** the secret value; a prompt-injection
  observation may set `category: "prompt-injection"` and is recorded, not followed
  (per
  [Boundaries](../../../specs/skills/quality-skill/quality-skill.md#boundaries-and-hard-rules)).

  _analysis_ — one per target node, the inferred weighted roll-up, citing what it
  derives from:

  ```json
  {
    "schemaVersion": 1,
    "target": "Webhooks", "targetPath": ["root", "webhooks"], "parent": "root",
    "localRating": "target", "aggregateRating": "minimum", "notAssessed": false,
    "factors": [ { "name": "Security", "rating": "target", "rationale": "…",
                   "subFactors": [], "contributing": [
                     { "requirement": "Every outbound webhook is signed…",
                       "assessmentRef": "006-webhooks-signing", "secondary": false } ] } ],
    "localRationale": "…", "aggregateRationale": "Delivery child (Minimum) pulls aggregate below local.",
    "derivedFrom": { "assessments": ["006-webhooks-signing"], "childAnalyses": ["delivery"] }
  }
  ```

  `factors` nests sub-factors recursively; `contributing` makes a secondary-factor
  requirement declared on a descendant identifiable under the factor it lenses.

  _report.json_ — the **render over** the records (not an independent copy): the
  top-level Rating and rationale, Scope, per-target ratings, the _not assessed_
  roll-up, Limitations, Advice, and a minimal finding summary for single-file gate
  and dashboard consumption. Each summary item carries only generic fields
  (`assessmentRef`, `category`, optional `severity`, `locator`, `observation`, and
  recommendation refs); full detail stays in `assessments/*.json`, so `report.json`
  remains a render/cache over the source records rather than a second source of
  truth. `report.md` is the same result for a person.

  ```json
  {
    "schemaVersion": 1,
    "rating": "unacceptable",
    "notAssessed": false,
    "findings": [
      {
        "assessmentRef": "assessments/001-root-no-committed-credentials.json",
        "category": "secret",
        "severity": "critical",
        "locator": "internal/gateway/client.go:48",
        "observation": "Live gateway secret key committed in plaintext.",
        "recommendations": ["001-rotate-committed-gateway-key"]
      }
    ]
  }
  ```

- **`improve` re-evaluation (item 5).** The post-apply re-assessment writes a
  **new** numbered folder: run _N_ applies the chosen option; run _N+1_ re-rates the
  affected scope and links back to _N_ (its `design.md` cites "re-evaluation of
  _N_ after applying recommendation 00X"). Reusing _N_ would mutate write-once records
  or mix two subject revisions under one model snapshot. The done-criterion is checked
  against _N+1_'s rating.

### 6. Driving the CLI and resolving inputs

- **Prerequisite check.** Before `setup`, `wizard`, `evaluate`, or `improve` uses
  CLI-dependent behavior, the skill checks `qualitymd` with `command -v` and
  `qualitymd --version`, compares the result with the skill's minimum compatible
  version, and stops on missing or stale versions with a concrete install/upgrade
  path. With user approval, it may run the appropriate platform install command;
  either way it verifies again with `qualitymd --version` before continuing. The
  skill facilitates this; it does not bundle or reimplement the CLI installer. For
  local dogfooding before a release tag exists, a development binary is acceptable
  when `qualitymd --version` identifies it as a dev build and the required commands
  (`spec`, `lint`, `init`, `models list`, `models view`) are present.
- **Mechanical steps.** `lint` gates judgment — the skill runs it first and refuses
  to rate a file with errors; `spec` grounds the format/schema rules and rating
  vocabulary; `init` scaffolds under `setup`; `models view` supplies the meta-model.
  The skill consumes `--json` where offered (`lint`, the `init` receipt,
  `models list`, and `models view` when it needs structured model data) and treats
  `spec` and default `models view` as Markdown artifacts, using their plain/verbatim
  output path for snapshots.
- **Introspection.** The skill discovers commands and flags from `qualitymd --help`
  / `qualitymd <cmd> --help` rather than embedding a list that drifts.
- **Default target file (item 4).** `QUALITY.md` in the current working directory,
  an explicit path overriding it, a clear error when neither resolves — matching
  [`init`](../../../specs/cli/init.md)'s output target. No tree walk or multi-file
  discovery; the skill defers to the CLI's file convention once
  [`cli.md`](../../../specs/cli.md#to-be-specified) lands one.
- **Model outline (item 7).** The `wizard`'s and the workflow's orientation comes
  from the **single-file read** the _Read the resolved target file_ step already
  performs — parsing the declared target/factor outline is a judgment-free
  structural read of one file, not source resolution or validation, so it needs no
  new CLI surface. A future top-level `qualitymd outline QUALITY.md --json` is
  preferred once it lands.
- **Source resolution.** Resolving each in-scope target's `source` to the entities
  to assess is the skill's judgment-adjacent read of the repository; no CLI
  record/resolve surface exists yet (it stays [deferred](../../../specs/skills/quality-skill/quality-skill.md#deferred)).

### 7. Durable-spec sync (In-Progress)

Implementation does not begin until the change is **In-Progress**; this doc only
plans the durable edits the spec's
[Affected specs & docs](../0010-implement-quality-skill.md#affected-specs--docs)
already commit to, so nothing outside this change folder is touched yet:

- [`quality-skill.md`](../../../specs/skills/quality-skill/quality-skill.md) — fold in
  the resolutions: model-altitude criteria via `models`, `setup`, the
  trigger-description criteria in Frontmatter and metadata, default-file rule,
  `.quality/config.yaml` with `evaluationDir`, the `improve` new-folder re-eval
  (correct the workflow diagram), the raw JSON artifact form, altitude-first folder
  naming, the **Limitations vs _not assessed_** distinction (item 10), the _not
  assessed_ done-criterion (item 11), and the _say-it-once_ consolidation of the
  conformance point (item 12).
- the [example bundle](../../../specs/skills/quality-skill/examples/index.md) and
  [`specs/schema.md`](../../../specs/schema.md) — re-capture assessments/analysis as
  JSON, add `report.json`, drop OKF frontmatter, re-slug to
  `0001-subject-quality-eval`, and retire the now-unused `Assessment Record` /
  `Analysis Record` / `Evaluation Report` / `Recommendation` types.
- new [`specs/cli/models.md`](../../../specs/cli/models.md) plus
  [`cli.md`](../../../specs/cli.md)/[`cli/index.md`](../../../specs/cli/index.md) — specify
  the `models` surface.
- [`README.md`](../../../README.md), root `install.md`, and `docs/` — introduce
  skill-first onboarding with `npx skills add qualitymd/quality.md`, the CLI
  prerequisite and setup verification, `/quality` modes/altitudes,
  `.quality/config.yaml` / `evaluationDir` as shared qualitymd system config, and
  the `specs/skills` indexes that point at the built skill. Root `install.md`
  structures the CLI installation step around verification: check
  `qualitymd --version`, install or upgrade through the documented release channel
  for the first public CLI version that contains `models`, verify again, then
  install the skill with `npx skills add qualitymd/quality.md` and verify it with
  `npx skills list`. The exact package-manager command is filled in when the
  release channel exists; the install guide's invariant is "documented install
  command plus verification," not a guessed package manager.

### 8. Dogfooding verification (In-Progress)

Before **Done**, verify the skill against this repository without treating the
rating as a release gate. The goal is artifact shape, invocation, and prerequisite
handling.

- Install the skill from the working tree with the local-path form the Agent Skills
  installer supports; published docs still show `npx skills add qualitymd/quality.md`.
- Use a local development `qualitymd` binary when it reports a dev version and
  exposes the required commands, so dogfooding can happen before the first release
  that contains `models`.
- Run a quick model-altitude evaluation against this repo's `QUALITY.md` (for
  example, `/quality evaluate model QUALITY.md --effort quick`) and confirm it
  writes the expected `NNNN-model-quality-eval/` artifact set under the resolved
  evaluation directory.
- Do not commit ad hoc dogfood run output unless intentionally re-captured as a
  durable example. During In-Progress, add the default `quality/evaluations/` path
  to `.gitignore`; configured non-default evaluation directories remain
  repository/user-owned and are not guessed.

### 9. Acceptance checklist

Before **Done**, the implementation satisfies this checklist:

- The repository contains `skills/quality/SKILL.md` with the final description text
  from [§3](#3-the-quality-skill-artifact), body-level argument guidance for
  mode/altitude/path/scope/effort, and the prompt body for the required modes.
- The skill installs from the published repo path (`npx skills add
qualitymd/quality.md`) and from the working tree for dogfooding when the installer
  supports a local path.
- `setup`, `wizard`, `evaluate`, and `improve` check for `qualitymd`, accept a
  compatible release, accept a dev build that exposes the required commands, and
  stop with a concrete install/upgrade path when the CLI is missing or stale.
- `qualitymd models list` and `qualitymd models list --json` expose the bundled
  `quality-meta-model`.
- `qualitymd models view quality-meta-model` renders like `qualitymd spec` for TTY
  users and emits byte-for-byte Markdown on plain output.
- `qualitymd models view quality-meta-model --json` emits stable JSON with catalog
  metadata, the decoded model, and `bodyMarkdown`.
- `qualitymd models view quality-meta-model --source QUALITY.md` applies the same
  source rewrite to both Markdown and JSON forms.
- `.quality/config.yaml` with `evaluationDir` changes the parent directory for
  numbered evaluation runs; missing config defaults to `quality/evaluations/`;
  absolute paths and repository escapes fail; unknown keys warn and are ignored.
- `quality/evaluations/` is ignored by default during In-Progress so dogfood output
  is not accidentally committed.
- `/quality evaluate model QUALITY.md --effort quick` can run against this repo
  using a local dev CLI, writes the expected artifact tree under the resolved
  evaluation directory, and does not treat the rating as a release gate.
- `assessments/*.json`, `analysis/*.json`, and `report.json` parse as JSON and use
  generic top-level shapes with domain-specific finding details under
  `attributes`.
- `report.json` includes minimal finding summaries by reference and does not
  duplicate full assessment detail.
- The worked example is re-captured in raw runtime form: no OKF frontmatter on
  runtime artifacts, JSON assessment/analysis records, `report.json`, and the
  altitude-first folder slug.
- Durable specs/docs listed in the parent change are updated before **Done**.

## Alternatives

- **Skill distribution (item 1).** _Claude plugin marketplace first_ was rejected:
  it makes the onboarding agent-specific, introduces namespaced invocation in Claude
  Code, and inverts the desired flow. _CLI-installs-the-skill_ was also rejected as
  the primary path because the user may not have the CLI yet; it remains a possible
  convenience later. `npx skills add qualitymd/quality.md` wins because it lets the
  skill be the entry point and leaves CLI setup to `setup`/`wizard` diagnostics.
- **`models view` JSON.** Considered keeping `view` as a pure verbatim-artifact
  carve-out like `spec`; rejected because a `QUALITY.md` model has a useful
  structured tree that agents and gates may need without reparsing Markdown/YAML.
  The default remains Markdown with the same human-rendered vs plain/verbatim split
  as `spec`, so terminal users get a readable view and `model.md` snapshots stay
  byte-for-byte; `--json` exposes the same rewritten model as data.
- **Re-pointing `source`.** Considered leaving `view` purely verbatim and having the
  skill note the real subject only in `design.md`; chose the `--source` flag so the
  snapshot _is_ the evaluated model and the substitution stays a mechanical CLI step
  rather than skill-side YAML editing.
- **`report.json` shape.** Considered a full inline twin duplicating every finding;
  chose a structured roll-up that **references** the assessment/analysis records,
  with only minimal finding summaries inlined for single-file gate/dashboard
  consumption. This preserves "the report is the render over these records, not an
  independent copy" while letting simple consumers avoid opening every assessment
  file for a top-level failure list.
- **Artifact format.** Keeping all records as Markdown (rejected — not a machine
  source of record) and making _everything_ JSON including the report (rejected — the
  human `report.md` and the narrative `design.md`/`plan.md` earn their prose).
- **Configuration scope.** Considered adding broader `.quality/config.yaml`
  settings for default target file, default effort, output formats, severity
  thresholds, retention, and install commands. Rejected for 0010: only
  `evaluationDir` solves a concrete repository-layout need. The config is still a
  shared qualitymd system config rather than skill-only config, because future CLI
  evaluation commands should honor it. Default target file and effort remain
  invocation rules, report formats are fixed by the artifact contract, thresholds
  belong to evaluation/gate work, retention is operational cleanup, and install
  commands belong in `install.md`.
- **Meta-model home.** Embedding from `internal/diagnostics/` in place (rejected —
  `go:embed` is package-local and a dedicated `internal/models` package is the honest
  home for a bundled-model catalog meant to grow).
- **Skill vs slash command.** A `commands/quality.md` slash command would work
  identically for invocation, but an Agent Skill (`skills/quality/SKILL.md`) is the
  right primitive: it supports broad Agent Skills installation and future bundled
  `references/` assets in the skill directory.

## Trade-offs & risks

- **Two version streams.** Skill installation and CLI installation are separate, so
  users can have a fresh skill with an old CLI. The mitigation is explicit: every
  entry mode checks the CLI version before relying on `models`, `spec`, `lint`, or
  `init`, and setup gives a concrete upgrade path.
- **Installer diversity.** `qualitymd` may be installed through Homebrew, `go
install`, npm, or another channel. The skill should prefer the documented install
  guide and verify the resulting binary rather than trying to infer every package
  manager state perfectly.
- **Bundled vs user-file surfaces.** `models` is intentionally plural because it is
  a bundled catalog. Future commands that inspect one user-supplied `QUALITY.md`
  should be top-level verbs (`outline`, `requirements`, `sources`) rather than a
  singular `model` namespace, avoiding `model`/`models` proximity.
- **JSON schemas are public API.** Once a gate reads `assessments/*.json` or
  `report.json`, their shape is a contract; `schemaVersion` is the lever, and field
  names are deliberately finalized during the example re-capture before anything
  consumes them. Top-level fields stay workflow-generic; bespoke assessment details
  belong under `attributes` so one factor's evidence type does not leak into the
  shared schema.
- **Skill correctness is prompt-bound.** Unlike the CLI, the skill's ratings can't
  be pinned by a unit test. The mitigation is structural: every mechanical step is
  the deterministic CLI's, and the [worked example](../../../specs/skills/quality-skill/examples/index.md)
  is the acceptance reference for artifact shape.
- **Skill cadence can outrun docs.** Because the skill can be updated independently
  from the CLI, a prompt-only fix can ship quickly. The cost is compatibility
  discipline: the skill's Prerequisites section and docs must be updated whenever a
  new CLI surface becomes required.
- **Boundaries stay prompt-enforced.** The secret-by-reference, prompt-injection-as-data,
  and scoped-not-whole rules are carried by the prompt, not mechanically guaranteed —
  inherent to a judgment skill, and the reason they are stated as hard rules.

## Open questions

None for this design. Deferred follow-ups:

- A future `qualitymd outline QUALITY.md --json` command may replace the
  skill's direct structural read of a user's model. It is not part of 0010.
- A future evaluation-management CLI surface may persist verdicts and gate CI
  through the CLI, likely under an `evaluations` command family. It is not part of 0010.
- Optional installer/UI metadata such as `agents/openai.yaml` is deferred until the
  Agent Skills installer or target agent documentation makes it necessary. For
  0010, `skills/quality/SKILL.md` is the canonical skill artifact.
- Deep effort may fan out assessment across targets/factors when the running agent
  supports subagents, but that is an optional implementation tactic, not a contract.
  The contract remains one run folder, one model snapshot, one plan, deterministic
  record names, and no duplicate writes.
