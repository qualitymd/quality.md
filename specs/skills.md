# Skills specification

> **Status:** *skeleton.* This document fixes the **skill layer** — the agentic
> orchestration that sits on top of the deterministic `qualitymd` CLI. It records
> the decisions settled so far and marks open detail with `TODO`. It is
> downstream of the resource-based CLI surface (see [`cli.md`](./cli.md)); field
> names, payload shapes, and skill bodies are expected to firm up as that surface
> and the meta-model settle.

This is the umbrella document for the skill layer. Individual skills are *not* yet
written; this spec defines the contract they orchestrate against.

## Purpose

`qualitymd` splits responsibility along one boundary:

- **The CLI is deterministic.** It parses and inspects the model, resolves
  targets, runs `bash` assessments, classifies results, persists evaluations,
  rolls up factors, renders reports, and diffs runs. It **never calls a model.**
- **Skills carry judgment and orchestration.** A skill drives the evaluation
  loop, judges `prompt` assessments, gathers evidence, decides scope and
  saturation, and writes verdicts back through the CLI.

This is the inversion from earlier drafts, where the agentic engine lived *inside*
the CLI (`evaluate` / `evaluate-model` as monolithic agentic commands). Those
commands are now **skills**, not CLI commands; the CLI's job shrinks to
deterministic management (`init`, `lint`, and the `model` / `evaluation` /
`result` resource trees). Authoring — scaffolding a first model and fixing an
existing one — is also skill work (`setup-quality-md`, `improve-quality-md`).

### The boundary is the Assessment dichotomy

The split is not arbitrary — it falls exactly on the distinction the format
already draws between the two kinds of assessment:

| Assessment | Kind | Who runs it | How it is recorded |
| --- | --- | --- | --- |
| `bash` | computational | **CLI** executes it, applies the rating condition | `result run` |
| `prompt` | inferential | **Skill** judges it | `result set` |

A consequence worth stating: a model made **entirely of `bash` requirements is
checkable with no skill and no model calls at all** — pure CI. Skills are required
only where `prompt` requirements exist.

## The skill set

| Skill | What it does | Model used |
| --- | --- | --- |
| `setup-quality-md` | **Onboards** a project: scaffolds and drafts a first `QUALITY.md` grounded in the project's real needs and risks, and sets up `.quality/`. | — (authors the model) |
| `evaluate-quality` | **Evaluates** the **subject** (the system/component) against its model. | the project's `QUALITY.md` |
| `improve-quality-md` | **Diagnoses then improves** the **`QUALITY.md` artifact itself** — critiques its well-formedness, then proposes/applies fixes to factors, requirements, and the Markdown body. | the built-in quality meta-model |

The three skills cover the model's life cycle: `setup-quality-md` brings a model
into being, `evaluate-quality` uses it to assess the subject, and
`improve-quality-md` keeps the model itself well-formed.

Mechanically, `improve-quality-md`'s **diagnose** phase **is** `evaluate-quality`
with the built-in meta-model wired in as the model and the user's `QUALITY.md`
supplied as the target — so it reuses one underlying loop and contract, then adds
an **improve** phase that turns the diagnosis into edits. `TODO`: specify the
exact wiring (how the bundled `QUALITY-META-MODEL.md` is selected as the model and
the user file passed as the target).

## Orchestration contract

The CLI is **strictly resource-based**: it exposes `model`, `evaluation`, and
`result` resources and never emits a prompt or assumes an iteration order. Prompt
composition and ordering are the skill's job. There is intentionally **no `next`
cursor** in the CLI.

The canonical loop a skill runs:

1. `evaluation create` — create/reset the living per-target run.
2. `result list --status pending,stale --json` — find the work.
3. For each result:
   - `result show <req> --json` — fetch the **resolved data** (prompt text,
     target manifest, rating scale). The skill composes this into a prompt.
   - Judge it, then `result set <req> --rating <level> --evidence …` — *or*, for
     a `bash` assessment, `result run <req>` (the CLI executes and classifies).
4. `evaluation report [--fail-on …]` — render the report and gate.

The bash-only, skill-free path collapses to:

```sh
qualitymd evaluation create && qualitymd result run --all && qualitymd evaluation report --fail-on …
```

### The CLI ↔ skill interface

The three payloads that *are* the contract are **specified authoritatively in
[`cli-evaluate.md`](./cli-evaluate.md#the-interface-payloads-cli--skill-contract)**
(under the `result` commands), with **illustrative JSON** for each; their field
names are provisional and expected to be tuned. This section summarizes what each
carries and links to the authoritative schema rather than re-specifying it.

- **[`result show` output](./cli-evaluate.md#result-show-output)** — what the skill
  needs to judge a `prompt` requirement: full requirement path + factor path,
  assessment kind, resolved prompt text, the resolved target file manifest (paths;
  contents optional/by reference), the rating scale with its `promptCondition`s
  (and any per-requirement overrides), the result's current state, and the
  sufficiency ("done") guidance for judging.
- **[`result set` input](./cli-evaluate.md#result-set-input)** — the verdict shape:
  requirement ref, rating level, structured evidence (summary + evidence items),
  and optional rationale. This becomes the diffable artifact, so its serialization
  has a stable field order and no inline volatile metadata.
- **[Staleness hash](./cli-evaluate.md#staleness-hash)** — what gets hashed (the
  requirement's resolved definition and the resolved target contents); a recorded
  result becomes `stale` on re-run if either changes.

## The evaluation lifecycle the skills drive

Settled; **specified in full in
[`cli-evaluate.md`](./cli-evaluate.md)** (the `evaluation` / `result` resources).
Summary, with links to the authoritative sections:

- **Living per-target run**, re-run in place. One evaluation per (model, target),
  stored under `.quality/evaluations/<slug>/`. Git history is the timeline.
- **Always mutable; no finalize/seal.** Git commits are the audit layer — see
  [Why no finalize: git is the audit layer](./cli-evaluate.md#why-no-finalize-git-is-the-audit-layer).
- **Commit everything.** Serialization is deterministic and **volatile metadata
  (timestamps, durations) is segregated** from verdicts so PR diffs show only
  rating/evidence changes — the *evaluation is a reviewable PR artifact* is a
  primary design goal (see
  [Segregating volatile metadata](./cli-evaluate.md#segregating-volatile-metadata-for-clean-diffs)).
- **Manual archive** — `evaluation archive --as <name>` snapshots to
  `.quality/evaluations/archive/<name>/`.
- [**Run states**](./cli-evaluate.md#run-states): `open → complete` (derived)
  `→ archived`. [**Result states**](./cli-evaluate.md#result-states):
  `pending / recorded / skipped / errored / stale` — on re-run, only `stale`
  results return to `pending`.

## Per-skill detail

### `setup-quality-md`

- **Orchestrates:** onboarding — scaffold via `qualitymd init`, then draft a
  starter model: elicit the project's real needs and risks, turn them into a small
  set of factors and requirements, write the Markdown body sections, and set up
  `.quality/`. Leans on `lint` to land a structurally clean file.
- **Owns (judgment):** what the project's first quality model should *say* — which
  needs and risks matter, which become factors and requirements, what the starter
  body documents. Authoring judgment grounded in the actual project.
- **Excludes:** evaluating anything (no run is created); the deep well-formedness
  critique (that is `improve-quality-md`, the natural next step once a model
  exists).
- **Outcome:** a lint-clean starter `QUALITY.md` grounded in the project, a
  set-up `.quality/`, and a user oriented to the evaluate / improve loop.

`TODO`: how much the starter model should infer from the codebase vs. ask the
user; the handoff to `improve-quality-md`.

### `evaluate-quality`

- **Orchestrates:** the loop above against the project's `QUALITY.md`.
- **Owns (judgment):** how to judge a `prompt` requirement — gather evidence
  against the resolved target, rate against the scale's `promptCondition`, decide
  when evidence is sufficient (saturation), write structured evidence back.
- **Excludes:** authoring/fixing the model; judging the model's well-formedness.
- **Outcome:** a complete, committed evaluation under `.quality/evaluations/<slug>/`
  with per-requirement ratings + evidence, a rendered report, and a CI pass/fail.

`TODO`: judging methodology (evidence standards, saturation, rigor levels).

### `improve-quality-md`

Two phases: **diagnose**, then **improve**.

- **Orchestrates (diagnose):** the `evaluate-quality` loop, with the built-in
  meta-model as the model and the user's `QUALITY.md` as the target; plus
  `model show` to inspect the structure under review. Produces a meta-evaluation
  of the model's well-formedness.
- **Orchestrates (improve):** turns that diagnosis into edits — proposing and (on
  the user's confirmation) applying fixes to factors, requirements, and the
  Markdown body, re-linting to keep the file clean, and re-running the diagnose
  loop to confirm the issues are resolved.
- **Owns (judgment):** the well-formedness critique — is each factor a genuine
  quality attribute (not a component or activity), distinct, operationalized; is
  the factor *set* complete and non-overlapping; is each requirement necessary,
  unambiguous, singular, verifiable; is the requirement *set* coherent. Expressed
  in `QUALITY.md` vocabulary (Factors / requirements). And the authoring judgment
  to fix each issue without distorting the model.
- **Excludes:** evaluating the *subject* (that is `evaluate-quality`); authoring a
  model from nothing (that is `setup-quality-md`).
- **Outcome:** a meta-evaluation naming specific mis-cast / ambiguous /
  unverifiable / overlapping elements, plus an improved, lint-clean `QUALITY.md`
  with those issues addressed. It owns the fix-it loop that diagnosis alone leaves
  open.

`TODO`: the meta-model is normative in
`internal/diagnostics/quality-model/QUALITY-META-MODEL.md`; specify how the skill
references it and how its criteria map onto the critique output. `TODO`: the edit
protocol for the improve phase (propose-then-apply, confirmation, re-lint /
re-diagnose).

## Distribution

Skills ship from the **public repo** and are installed into a user's agent. As
public artifacts they follow the project's convention: prefer `QUALITY.md`'s own
vocabulary and keep standards references in the background (this spec, like other
docs under `specs/`, may cite provenance where it is relevant).

`TODO`: decide the skills' home in the repo (e.g. `skills/`), packaging, and
whether any are marked internal.

## Relationship to the other specs

- [`cli.md`](./cli.md) — the resource-based command surface the skills orchestrate
  against.
- [`cli-evaluate.md`](./cli-evaluate.md) — the deterministic evaluation lifecycle
  and the **authoritative CLI ↔ skill interface payloads** this contract leans on.
- [`cli-compare.md`](./cli-compare.md) — `evaluation compare` as a deterministic
  diff of stored runs (A/B and base-vs-head regression gating).
- [`cli-lint.md`](./cli-lint.md), [`cli-init.md`](./cli-init.md) — the
  deterministic commands; largely unaffected.
- [`../SPECIFICATION.md`](../SPECIFICATION.md) — the format itself, including the
  Requirement / Assessment / Rating separation this contract leans on.

## Open questions

- The three interface payloads (`result show`, `result set`, staleness hashing)
  are defined illustratively in
  [`cli-evaluate.md`](./cli-evaluate.md#the-interface-payloads-cli--skill-contract);
  their field names are provisional and the staleness-hash serialization is still open.
- Judging methodology shared by the evaluate loop (evidence, saturation, rigor) —
  reused by both `evaluate-quality` and `improve-quality-md`'s diagnose phase.
- The exact meta-model wiring for `improve-quality-md`'s diagnose phase.
- The improve-phase edit protocol (propose-then-apply, confirmation, re-lint /
  re-diagnose).
- How much `setup-quality-md` infers from the codebase vs. asks the user, and its
  handoff to `improve-quality-md`.
- Skill packaging and home in the public repo.
