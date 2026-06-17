---
type: Functional Specification
title: Implement the /quality skill — functional spec
description: The delta for building the /quality skill — conform to the durable skill spec, package an invocable skill, and settle the open items and gaps that spec leaves.
tags: [skill, quality, evaluation]
timestamp: 2026-06-17T00:00:00Z
---

# Implement the /quality skill — functional spec

> 🚧 **Draft.** Companion to
> [Implement the /quality skill](../0010-implement-quality-skill.md). This spec
> states *what* the change must do; a design doc will cover *how*.

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
[open items and gaps](#open-items-and-gaps-settle-before-design) below, syncing the
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
- It **MUST** be an **invocable skill** — discoverable and runnable as `/quality`
  with the mode, altitude, target-file, scope, and effort arguments the spec
  defines.
- It **MUST** drive the deterministic CLI for every mechanical step and treat its
  output as the source of truth, per
  [Driving the CLI](../../specs/skills/quality-skill/quality-skill.md#driving-the-cli);
  it **MUST NOT** reimplement scaffolding, structural validation, or the format
  rules.
- Its runtime artifacts **MUST** satisfy the
  [Reporting](../../specs/skills/quality-skill/quality-skill.md#reporting) contract
  and **SHOULD** match the shape of the
  [worked example bundle](../../specs/skills/quality-skill/examples/index.md),
  which serves as the acceptance reference.
- Before the change reaches **Done**, every
  [open item](#open-items-and-gaps-settle-before-design) **MUST** be resolved and
  its resolution reflected back into the durable skill spec (and the example where
  affected), so the enduring spec matches what was built.

## Open items and gaps (settle before Design)

A review of the durable skill spec surfaced the items below. Each is a *what*
decision this change MUST settle while **Draft** — recorded here so it stays
visible rather than discovered mid-implementation — with a **recommended
resolution** to confirm or revise. On resolution, each is carried into
[`quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md) (and the
example where affected); see the parent's
[Affected specs & docs](../0010-implement-quality-skill.md#affected-specs--docs).

### Blocking — resolve before Design

1. **Packaging and location.** Where the skill lives and how it is packaged is
   undefined — the repo has no skill artifact, `skills/` tree, or `.claude/`
   directory today.
   - **Resolution:** add a top-level `skills/quality/` source home (sibling to
     `cmd/`, `internal/`, `specs/`) containing the versioned `/quality` Agent
     Skill artifact, including `SKILL.md`. Distribution and installation mechanics
     are Design sub-decisions; the *source home* is settled now.

2. **The `model` altitude's criteria source.** Evaluating *the model itself* needs
   criteria for what makes a `QUALITY.md` *good*, but the spec never says where
   they come from; grounding (`qualitymd spec`) supplies only the **format** rules.
   The repo already ships a built-in
   [quality meta-model](../../internal/diagnostics/quality-model/QUALITY-META-MODEL.md)
   — itself a `QUALITY.md` whose subject is a `QUALITY.md` — for exactly this, yet
   the skill spec neither references it nor is there a CLI command that emits it
   (only `spec`, which emits the format specification).
   - **Recommendation:** treat meta-evaluation as ordinary evaluation with the
     meta-model as the active model and the user's `QUALITY.md` as its subject. Add
     a deterministic CLI surface that emits the bundled meta-model — mirroring
     `qualitymd spec` (e.g. `qualitymd spec --meta-model`, or a `qualitymd
     meta-model` command) — and have the `model` altitude ground in it the same way
     the `subject` altitude grounds the format rules, so it rates against declared
     criteria, never an invented standard. Re-point the meta-model's `source` at the
     user's model at runtime. This unifies the two altitudes under one workflow and
     reuses the already-maintained (change 0009) meta-model.

3. **What `setup` does.** `setup` is a first-class mode in the arguments and
   examples and the wizard routes to it, but — unlike `evaluate`/`improve`/`wizard`
   — it has no behavioral section, is absent from the operating-model loop, and is
   not in the workflow.
   - **Recommendation:** specify `setup` as two steps — (1) drive
     [`init`](../../specs/cli/init.md) for the deterministic, `lint`-valid skeleton
     (never reimplement scaffolding), then (2) apply judgment to populate it from
     the subject: propose a `title`, factors, at least one requirement per factor,
     and the body sections (Overview/Scope/Needs/Risks), each inferred from
     inspecting the source. This preserves the determinism/judgment split and makes
     `setup` more than an `init` alias (which would add nothing the CLI lacks). It
     MUST end `lint`-valid.

### Spec-affecting — resolvable during implementation

4. **Default target-file resolution.** The spec says the skill resolves a default
   file "the way the CLI does," but the CLI's file-argument convention is itself
   listed under *To be specified* in [`cli.md`](../../specs/cli.md).
   - **Recommendation:** adopt the one concrete CLI precedent now — default to
     `QUALITY.md` in the current working directory (as
     [`init`](../../specs/cli/init.md)'s output target does), accept an explicit
     path override, and error clearly when none is found. No directory-tree walk or
     multi-file discovery (YAGNI). When the CLI's convention lands, the skill defers
     to it; until then this is the rule, recorded so it is not silently divergent.

5. **`improve`'s post-apply re-evaluation.** The workflow re-evaluates within one
   run, but Reporting says a re-assessment "produces a new evaluation folder rather
   than editing an existing record."
   - **Recommendation:** the post-apply re-evaluation writes a **new** numbered
     folder (run N applies; run N+1 re-rates and links back to N), and the workflow
     diagram is corrected to show this. Reusing the folder would either mutate
     write-once records (forbidden) or mix two subject revisions under one model
     snapshot (breaking the "one folder = one evaluation at one revision"
     guarantee). The done-criterion is confirmed against the new folder's rating.

6. **Machine-readable report output.** Reporting **SHOULD** render a
   machine-readable form "for a gate or tool," but no argument selects output
   format.
   - **Recommendation:** since each run already writes a durable folder, always emit
     a structured `report.json` sibling to `report.md` (the
     [`--json` convention](../../specs/cli.md#conventions) governs its shape and
     field stability), rather than threading an output-format argument through
     `/quality`. Both renderings of the same result exist; a consumer reads
     whichever it needs. No new invocation argument.

7. **The wizard's model-outline source.** The wizard inspects "the targets and
   factors the model declares," but [`lint`](../../specs/cli/lint.md) emits
   findings, not a model outline, and no command emits one (`lint.md` only
   *reserves* a future `model-summary` info rule) — in tension with "every
   mechanical step driven through the CLI."
   - **Recommendation:** the wizard's orientation comes from the same single-file
     read the workflow's *Read the resolved target file* step already performs —
     parsing the model's own declared frontmatter outline is a judgment-free
     structural read of one declared file, not source resolution, so it sits inside
     the existing read rather than needing a new CLI surface. Prefer a `qualitymd`
     outline / `model-summary` surface once it lands; until then the direct read is
     the source.

### Durable-spec corrections to apply during the sync

Consistency fixes, not new behavior — applied to
[`quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md) (and the
example) when the durable spec is brought into sync.

8. **Folder `<scope>` naming.** The example's `0001-payments-quality-eval` is a
   whole-model **subject** run, yet `payments` is the subject's name, not the
   "altitude and narrowing" the spec says `<scope>` encodes — the default
   whole-model case has no derivable slug.
   - **Recommendation:** make the slug deterministic and altitude-first —
     `NNNN-<altitude>[-<narrowing>]-quality-eval`, where `<altitude>` is `subject`
     or `model` and `<narrowing>` is the scoped target/factor name (omitted when
     whole-model). Re-slug the example to `0001-subject-quality-eval`. Trade-off:
     less evocative than the subject name, but predictable and sortable; if a
     subject-name slug is preferred for whole-model runs, state that as the explicit
     rule instead — either way the rule must determine the example's name.

9. **`Limitations` vs *not assessed*.** The example report separates a
   **Limitations** section (effort ceilings, point-in-time scans, single-test
   confidence) from per-requirement *not assessed* outcomes; the spec folds both
   under one "what was not assessed" phrase.
   - **Recommendation:** name two distinct report elements in Reporting — (a) *not
     assessed* outcomes (a Rating Result where evidence was absent, shown per
     requirement and roll-up) and (b) a **Limitations** statement that bounds how
     far a *rated* outcome should be trusted and reconciles coverage against the
     plan. The example already implements both; the spec just needs the vocabulary.

10. **Done-criterion for a *not assessed* gap.** The spec defines a done-criterion
    as "the target rating level the requirement should reach," but a *not assessed*
    gap's fix is to *become assessable, then* reach a level — the wording example
    recommendation 002 improvises.
    - **Recommendation:** broaden the definition to "the outcome the in-scope
      requirement should reach against its `criterion` — for a rated gap, a target
      rating level; for a *not assessed* gap, becoming assessable and reaching at
      least the acceptable floor."

11. **Say it once.** The conformance-vs-deference point is restated in four-plus
    places (draft banner, *Frontmatter and metadata*, *Driving the CLI*, the
    *Conformance to the format spec* section, the *Workflow* intro).
    - **Recommendation:** keep the *Conformance to the format spec* section as the
      single home; cut the draft banner to a one-line pointer and replace the other
      repeats with links. Leave *Frontmatter and metadata* its one distinct rule —
      format rules and vocabulary are grounded at runtime while the evaluation
      process is owned — stated once, linking to Conformance for the rest. Per
      [the spec-writing guide](../../docs/guides/write-functional-specs.md).
