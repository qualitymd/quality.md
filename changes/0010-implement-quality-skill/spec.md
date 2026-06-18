---
type: Functional Specification
title: Implement the /quality skill — functional spec
description: The delta for building the /quality skill — conform to the durable skill spec, package an invocable skill, and settle the open items and gaps that spec leaves.
tags: [skill, quality, evaluation]
timestamp: 2026-06-17T00:00:00Z
---

# Implement the /quality skill — functional spec

> Companion to
> [Implement the /quality skill](../0010-implement-quality-skill.md). This spec
> states *what* the change must do; the design doc covers *how*.

This change **implements an existing durable spec** rather than authoring a new
behavioral contract. The *what* — operating model, invocation, evaluation
workflow, grounding, effort levels, and the reporting and artifact contract —
already lives in the
[`/quality` skill spec](../../specs/skills/quality-skill/quality-skill.md) and its
[worked example bundle](../../specs/skills/quality-skill/examples/index.md). Per
[one source of truth](../../docs/guides/write-functional-specs.md), this spec does
**not** restate that contract; it states only the **delta** this change owns and
the **open items and gaps** the durable spec leaves for the change to settle.

The key words **MUST**, **SHOULD**, and **MAY** are to be interpreted as described
in IETF RFC 2119.

## Scope

Covered: packaging and wiring an invocable `/quality` skill that **conforms to**
the [durable skill spec](../../specs/skills/quality-skill/quality-skill.md); the
runtime evaluation artifacts it writes; and resolving the
[open items and gaps](#open-items-and-gaps) below, syncing the
durable spec to the resolution.

Deferred: everything the skill spec's own
[Deferred](../../specs/skills/quality-skill/quality-skill.md#deferred) section
defers — recording verdicts *through the CLI*, `improve` apply-staging mechanics,
and bundled `references/` assets — none of which this change introduces.

## Requirements

- The implementation **MUST** conform to the
  [`/quality` skill spec](../../specs/skills/quality-skill/quality-skill.md) in
  full: its
  [operating model](../../specs/skills/quality-skill/quality-skill.md#operating-model),
  [boundaries and hard rules](../../specs/skills/quality-skill/quality-skill.md#boundaries-and-hard-rules),
  [invocation](../../specs/skills/quality-skill/quality-skill.md#invocation),
  [evaluation workflow](../../specs/skills/quality-skill/quality-skill.md#evaluation-workflow),
  and [reporting](../../specs/skills/quality-skill/quality-skill.md#reporting)
  contract. Where this change and that spec would diverge, the spec governs and is
  corrected (by resolving an open item) rather than silently departed from.
- It **MUST** be an **invocable skill** installable from this repository through
  Agent Skills tooling (for example, `npx skills add qualitymd/quality.md`) and
  runnable as `/quality` with the mode, altitude, target-file, scope, and effort
  arguments the spec defines.
- Its `SKILL.md` description **MUST** be trigger-oriented, not merely documentary:
  it **MUST** cover quality management, quality evaluation, and quality improvement
  requests even when the user does not mention `QUALITY.md`; it **MUST** include
  trigger vocabulary for the supported modes (`setup`, `wizard`, `evaluate`,
  `improve`) and for generic quality terms such as factors, characteristics,
  attributes, criteria, Targets, Factors, Requirements, subject evaluation, and
  model evaluation/improvement; it **SHOULD** frame the subject broadly as a
  project/entity or one of its components/targets rather than enumerating a closed
  set of subject types; it **SHOULD NOT** include CLI implementation details that
  belong in the skill body; it **SHOULD** avoid triggering for generic copyediting
  or one-off "make this higher quality" requests that do not ask for systematic
  quality criteria, assessment, or management.
- It **MUST** drive the deterministic CLI for every mechanical step and treat its
  output as the source of truth, per
  [Driving the CLI](../../specs/skills/quality-skill/quality-skill.md#driving-the-cli);
  it **MUST NOT** reimplement scaffolding, structural validation, or the format
  rules.
- Its `setup` and `wizard` flows **MUST** detect whether the `qualitymd` CLI is
  missing or below the skill's required compatible version and facilitate install
  or upgrade before running CLI-dependent work.
- It **MUST** support the repository-local qualitymd system configuration file
  `.quality/config.yaml`, initially with an `evaluationDir` setting that controls
  where evaluation run folders are written. The skill reads it in this change; the
  CLI is expected to share the same config surface as future CLI evaluation
  commands land. When the file or setting is absent, the default **MUST** remain
  `quality/evaluations/`.
- Its runtime artifacts **MUST** satisfy the
  [Reporting](../../specs/skills/quality-skill/quality-skill.md#reporting) contract
  and **SHOULD** match the shape of the
  [worked example bundle](../../specs/skills/quality-skill/examples/index.md),
  which serves as the acceptance reference.
- Before the change reaches **Done**, every
  [open item](#open-items-and-gaps) **MUST** be resolved and
  its resolution reflected back into the durable skill spec (and the example where
  affected), so the enduring spec matches what was built.

## Open items and gaps

A review of the durable skill spec surfaced the items below. Each is a *what*
decision this change owns, **surfaced now while Draft** — recorded here so it
stays visible rather than discovered mid-implementation — with a **recommended
resolution** to confirm or revise. They resolve on the schedule their grouping
names: the **blocking** items before **Design**, the rest during **In-Progress**,
and all of them before **Done**. On resolution, each is carried into
[`quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md) (and the
example where affected); see the parent's
[Affected specs & docs](../0010-implement-quality-skill.md#affected-specs--docs).

### Blocking — resolve before Design

1. **Packaging and location.** Where the skill lives and how it is packaged is
   undefined — the repo has no skill artifact or `skills/` tree today.
   - **Resolution:** add a top-level `skills/quality/` source home (sibling to
     `cmd/`, `internal/`, `specs/`) containing the `/quality` Agent Skill
     artifact, including `SKILL.md`. The repository **MUST** be installable as an
     Agent Skills source with `npx skills add qualitymd/quality.md`, following the
     Basecamp-style skill repository pattern. Claude plugin or marketplace
     packaging **MAY** be added later as a secondary channel, but it is not the
     primary onboarding path for this change.

2. **The `model` altitude's criteria source.** Evaluating *the model itself* needs
   criteria for what makes a `QUALITY.md` *good*, but the spec never says where
   they come from; grounding (`qualitymd spec`) supplies only the **format** rules.
   The repo already ships a built-in
   [quality meta-model](../../internal/diagnostics/quality-model/QUALITY-META-MODEL.md)
   — itself a `QUALITY.md` whose subject is a `QUALITY.md` — for exactly this, yet
   the skill spec neither references it nor is there a CLI command that emits it
   (only `spec`, which emits the format specification).
   - **Resolution:** treat meta-evaluation as ordinary evaluation with the bundled
     quality meta-model as the active model and the user's `QUALITY.md` as its
     subject. Add a deterministic `qualitymd models` CLI surface for bundled
     models; its catalog **MUST** include the quality meta-model so the `model`
     altitude grounds in declared criteria, never an invented standard. The surface
     **MUST** expose both the verbatim Markdown model and a JSON representation of
     the same model for agents/tools that need structured access. Future sample
     models **MAY** be exposed through the same catalog. Re-point the meta-model's
     `source` at the user's model at runtime. This unifies the two altitudes under
     one workflow and reuses the already-maintained (change 0009) meta-model.

3. **What `setup` does.** `setup` is a first-class mode in the arguments and
   examples and the wizard routes to it, but — unlike `evaluate`/`improve`/`wizard`
   — it has no behavioral section, is absent from the operating-model loop, and is
   not in the workflow.
   - **Resolution:** specify `setup` as the minimal bootstrap path after the skill
     is installed: check whether the `qualitymd` CLI is available, verify its
     version is compatible with the skill, facilitate install or upgrade when it is
     missing or too old, then drive [`init`](../../specs/cli/init.md) to create a
     deterministic, `lint`-valid `QUALITY.md` skeleton, validate it with
     [`lint`](../../specs/cli/lint.md), and defer to `wizard` for guided population
     and refinement. `setup` **MUST NOT** reimplement scaffolding, validation, CLI
     installation tooling, or source-driven authoring judgment. `wizard` owns the
     cursory repository assessment: infer and confirm suggested targets/factors,
     ask clarifying questions when important model inputs are missing (purpose,
     entity type, stakeholder needs, risks, and similar context), and guide the user
     toward evaluating or improving the model once enough context exists.

### Spec-affecting — resolvable during implementation

4. **Default target-file resolution.** The spec says the skill resolves a default
   file "the way the CLI does," but the CLI's file-argument convention is itself
   listed under *To be specified* in [`cli.md`](../../specs/cli.md).
   - **Resolution:** adopt the one concrete CLI precedent now — default to
     `QUALITY.md` in the current working directory (as
     [`init`](../../specs/cli/init.md)'s output target does), accept an explicit
     path override, and error clearly when none is found. No directory-tree walk or
     multi-file discovery. When the CLI's convention lands, the skill defers to it;
     until then this is the skill's rule.

5. **`improve`'s post-apply re-evaluation.** The workflow re-evaluates within one
   run, but Reporting says a re-assessment "produces a new evaluation folder rather
   than editing an existing record."
   - **Resolution:** the post-apply re-evaluation writes a **new** numbered folder
     (run N applies; run N+1 re-rates and links back to N), and the workflow diagram
     is corrected to show this. Reusing the folder would either mutate write-once
     records or mix two subject revisions under one model snapshot, breaking the
     "one folder = one evaluation at one revision" guarantee. The done-criterion is
     confirmed against the new folder's rating.

6. **Artifact form — machine-readable, and not OKF concepts.** Reporting
   **SHOULD** render a machine-readable form "for a gate or tool," but no argument
   selects output format, and the assessment and analysis records are specified as
   write-once artifacts without fixing whether their source-of-record form is prose
   or structured data. Separately, the worked example stores each artifact as an
   **OKF concept** — a `type:`-bearing Markdown file registered in
   [`specs/schema.md`](../../specs/schema.md) — conflating a runtime output with an
   entry in this knowledge bundle.
   - **Resolution:** evaluation artifacts are **raw runtime outputs written into
     the evaluated repository, not OKF concepts** — none carries OKF frontmatter or
     a registered type. Their form follows their consumer: assessment and analysis
     records are source-of-record **data** and **MUST** be JSON
     (`assessments/*.json`, `analysis/*.json`), and each run also emits a structured
     `report.json` beside the human `report.md` — two renderings of one result, so a
     consumer reads whichever it needs and no output-format argument is threaded
     through `/quality`. `report.md`, the recommendations, `design`, and `plan` stay
     Markdown; `model.md` stays a verbatim snapshot. The worked example is
     re-captured in this raw form and its now-unused concept types are retired from
     [`specs/schema.md`](../../specs/schema.md) (see the parent's
     [Affected specs & docs](../0010-implement-quality-skill.md#affected-specs--docs)).

7. **The wizard's model-outline source.** The wizard inspects "the targets and
   factors the model declares," but [`lint`](../../specs/cli/lint.md) emits
   findings, not a model outline, and no command emits one (`lint.md` only
   *reserves* a future `model-summary` info rule) — in tension with "every
   mechanical step driven through the CLI."
   - **Resolution:** the wizard's model orientation comes from the same single-file
     read the workflow's *Read the resolved target file* step already performs.
     Parsing the model's declared target/factor outline is a judgment-free
     structural read of one declared file, not source resolution or validation, so
     it does not require a new CLI surface. When no model exists or the model is
     skeletal, `wizard` may also perform a cursory repository scan and ask
     clarifying questions per the `setup` resolution above. Prefer a future
     top-level `qualitymd outline QUALITY.md --json` surface once it lands; until
     then the direct read is the source.

8. **Evaluation output directory configuration.** The durable spec fixes the
   evaluation folder shape but not whether repositories can choose where those
   folders live.
   - **Resolution:** add a repository-local qualitymd system config at
     `.quality/config.yaml`, read by the skill in this change and reserved for
     future CLI use. It **MUST** support `evaluationDir`, defaulting to
     `quality/evaluations/` when unset. The path **MUST** be repository-relative,
     normalized, and forbidden from escaping the repository. It points to the parent
     directory that contains numbered run folders; the deterministic run-folder
     naming rule still applies inside it. Unknown keys **SHOULD** be surfaced as
     warnings and ignored. Other configuration keys are deferred until a real need
     appears.

### Durable-spec corrections to apply during the sync

Consistency fixes, not new behavior — applied to
[`quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md) (and the
example) when the durable spec is brought into sync.

9. **Folder `<scope>` naming.** The example's `0001-payments-quality-eval` is a
   whole-model **subject** run, yet `payments` is the subject's name, not the
   "altitude and narrowing" the spec says `<scope>` encodes — the default
   whole-model case has no derivable slug.
   - **Resolution:** make the slug deterministic and altitude-first —
     `NNNN-<altitude>[-<narrowing>]-quality-eval`, where `<altitude>` is `subject`
     or `model` and `<narrowing>` is the scoped target/factor name (omitted when
     whole-model). Re-slug the example to `0001-subject-quality-eval`. This is less
     evocative than the subject name, but predictable and sortable.

10. **`Limitations` vs *not assessed*.** The example report separates a
    **Limitations** section (effort ceilings, point-in-time scans, single-test
    confidence) from per-requirement *not assessed* outcomes; the spec folds both
    under one "what was not assessed" phrase.

- **Resolution:** name two distinct report elements in Reporting — (a) *not
  assessed* outcomes (a Rating Result where evidence was absent, shown per
  requirement and roll-up) and (b) a **Limitations** statement that bounds how
  far a *rated* outcome should be trusted and reconciles coverage against the
  plan. The example already implements both; the spec just needs the vocabulary.

11. **Done-criterion for a *not assessed* gap.** The spec defines a done-criterion
    as "the target rating level the requirement should reach," but a *not assessed*
    gap's fix is to *become assessable, then* reach a level — the wording example
    recommendation 002 improvises.
    - **Resolution:** broaden the definition to "the outcome the in-scope
      requirement should reach against its `criterion` — for a rated gap, a target
      rating level; for a *not assessed* gap, becoming assessable and reaching at
      least the acceptable floor."

12. **Say it once.** The conformance-vs-deference point is restated in four-plus
    places (draft banner, *Frontmatter and metadata*, *Driving the CLI*, the
    *Conformance to the format spec* section, the *Workflow* intro).
    - **Resolution:** keep the *Conformance to the format spec* section as the
      single home; cut the draft banner to a one-line pointer and replace the other
      repeats with links. Leave *Frontmatter and metadata* its one distinct rule —
      format rules and vocabulary are grounded at runtime while the evaluation
      process is owned — stated once, linking to Conformance for the rest. Per
      [the spec-writing guide](../../docs/guides/write-functional-specs.md).
