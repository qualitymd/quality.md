# CLI specification

Functional specification for the `qualitymd` CLI surface.

> **Status:** high-level *functional* spec. It fixes the command surface, the
> deterministic-CLI / skill-layer split, and the output contracts. Exact JSON
> field names, lint-rule identifiers, and resource payload shapes are
> illustrative and expected to be tuned during implementation. Items still
> genuinely open are collected under [Open questions](#open-questions).
>
> This document was rewritten around a **resource-based** surface. The deep
> agentic commands (`evaluate`, `evaluate-model`) are no longer CLI
> commands — they are **skills** that orchestrate the deterministic surface below
> (see [`skills.md`](./skills.md)).

This is the umbrella document. Per-command and per-layer detail lives in:

- [`cli-init.md`](./cli-init.md) — scaffolding a first `QUALITY.md` from a static
  template (deterministic, offline).
- [`cli-lint.md`](./cli-lint.md) — the fast, deterministic structural validator
  (frontmatter and Markdown body rules).
- [`skills.md`](./skills.md) — the **skill layer**: how the skills
  (`setup-quality-md`, `evaluate-quality`, `improve-quality-md`) orchestrate the
  `model` / `evaluation` / `result` resources, and the CLI ↔ skill interface.
- [`cli-evaluate.md`](./cli-evaluate.md) — the deterministic **evaluation
  lifecycle** (the `evaluation` and `result` resources): data model, on-disk
  layout, the `bash` execution path, the report rollup, and the **authoritative
  CLI ↔ skill interface payloads**.
- [`cli-federation.md`](./cli-federation.md) — how multiple `QUALITY.md` models in
  one repository compose: discovery, ownership and scope, one run per model, the
  per-model-gated tree report. Cross-cutting rather than per-command.

## The split: deterministic CLI, judgment in skills

`qualitymd` draws one hard line:

- **The CLI is deterministic and never calls a model.** It parses and inspects
  the model, resolves targets, runs `bash` assessments, classifies results,
  persists evaluations, rolls up factors, and renders reports.
- **Skills carry judgment and orchestration** (see [`skills.md`](./skills.md)).
  A skill drives the evaluation loop, judges `prompt` assessments, gathers
  evidence, and writes verdicts back through the CLI.

The line falls exactly on the format's two assessment kinds:

| Assessment | Kind | Who runs it | Recorded via |
| --- | --- | --- | --- |
| `bash` | computational | **CLI** executes and classifies it | `result run` |
| `prompt` | inferential | **Skill** judges it | `result set` |

A consequence worth stating: a model made **entirely of `bash` requirements is
evaluable with no skill and no model calls at all** — pure, reproducible CI.

## Command surface

Two deterministic top-level commands, plus three resources — `model` (read-only),
`evaluation` (a run), and `result` (a requirement-result within a run). Resource
nouns are **singular**; the verb carries the cardinality.

### Top-level

| Command | Purpose | Output |
| --- | --- | --- |
| `qualitymd init` | Scaffold a starter `QUALITY.md` from a static template, for the author to fill in. | `./QUALITY.md` on disk. |
| `qualitymd lint [file]` | Structural validation of the `QUALITY.md` against the format spec — parses, conforms to the schema, references resolve. | JSON findings + summary on stdout. |

### `qualitymd model` — read-only inspection

| Command | Purpose | Output |
| --- | --- | --- |
| `model show [--requirement <path>] [--json]` | The parsed model: factor tree, requirements by full path, **resolved** targets (globs expanded), loaded `prompt`/`bash` text, the rating scale. `--requirement` narrows to one fully-resolved requirement. | JSON |

### `qualitymd evaluation` — the run (alias `eval`)

| Command | Purpose | Transition |
| --- | --- | --- |
| `evaluation create [--model <path>] [--target <path>] [--from <id>]` | Create or reset the living per-target run; enumerate requirements as `pending`. `--from` carries forward still-valid results. | → **open** |
| `evaluation list` | List runs (active + archived). | — |
| `evaluation show [<id>]` | Run manifest: status, rollup, verdict. | — |
| `evaluation report [<id>] [--fail-on <level>]` | Render the report; non-zero exit when `--fail-on` is tripped. | — |
| `evaluation archive [<id>] --as <name>` | Snapshot current state to `archive/<name>/`. | → **archived** |
| `evaluation delete <id>` | Discard a run. | → abandoned |

### `qualitymd result` — a requirement-result within a run

| Command | Purpose | Transition |
| --- | --- | --- |
| `result list [--status pending,stale,…] [--json]` | Query results by state. There is **no `next` cursor** — the skill orders the work. | — |
| `result show <req> [--json]` | The resolved data for one requirement (prompt text, target manifest, scale). The skill composes the prompt; the CLI emits none. | — |
| `result run <req\|--all>` | Execute `bash` assessment(s), classify, record. `bash` only. | pending → **recorded** / **errored** |
| `result set <req> --rating <level> --evidence …` | Record a `prompt` verdict (the skill's judgment). The diffable artifact. | pending → **recorded** |
| `result skip <req> --reason …` | Deliberately not assessed. | pending → **skipped** |
| `result reset <req>` | Return to pending, to re-judge. | → **pending** |

### Structural tier vs. management tier

The deterministic surface answers two different questions:

- **Structural — `lint`.** *Is this a well-formed `QUALITY.md`?* Frontmatter
  parses, every factor has `requirements`, `factors`, or both, every requirement
  declares exactly one assessment, `prompt`/`target` paths resolve, the `ratings`
  scale is well-shaped. Cheap enough for every save and a pre-commit hook.
  Modeled on Google's `design.md lint` (see [`cli-lint.md`](./cli-lint.md)).
- **Management — `model` / `evaluation` / `result`.** *What is the recorded state
  of evaluating this subject?* Deterministic inspection of the model and CRUD over
  the evaluation lifecycle, including running `bash` assessments. The *judgment*
  over `prompt` requirements is not here — it is the skill layer
  ([`skills.md`](./skills.md)).

The tiers are complementary. `lint` proves the *file* is correct without an
opinion on whether the requirements are any good; the `improve-quality-md`
skill forms that opinion but cannot run meaningfully on a file that does not
parse. **Lint first, evaluate second.**

`init` is step zero — a deterministic scaffold that produces the first
`QUALITY.md` for the author to fill in (see [`cli-init.md`](./cli-init.md)). The
intended order is **init → lint → evaluate** (the latter via skills). This also
settles the previously-deferred `check` verb: the structural CI gate is `lint`;
the evaluation gate is `evaluation report --fail-on`. There is no separate
`check`.

## The evaluation lifecycle

The `evaluation` and `result` commands manage a single, durable model: a **living
per-target run**. Full detail lives in
[`cli-evaluate.md`](./cli-evaluate.md) and is summarized in
[`skills.md`](./skills.md#the-evaluation-lifecycle-the-skills-drive). In brief:

- **One run per (model, target)**, re-run in place, stored under
  `.quality/evaluations/<slug>/`. Git history is the timeline.
- **Always mutable; no finalize/seal.** Git commits are the audit layer.
- **Commit everything.** Serialization is deterministic and **volatile metadata
  (timestamps, durations) is segregated** from verdicts, so a PR diff shows only
  rating/evidence changes — *the evaluation is a reviewable PR artifact* is a
  primary design goal.
- **Manual archive** — `evaluation archive --as <name>` snapshots to
  `.quality/evaluations/archive/<name>/`.
- Run states: `open → complete` (derived) `→ archived`. Result states:
  `pending / recorded / skipped / errored / stale`. On re-run, only `stale`
  results return to `pending`.

The canonical skill loop and the bash-only, skill-free CI path:

```sh
# Skill loop (judgment over prompt requirements)
qualitymd evaluation create
qualitymd result list --status pending,stale --json     # skill orders + composes prompts
#   per result: qualitymd result show <req> --json → judge → qualitymd result set <req> …
qualitymd evaluation report --fail-on fail

# Bash-only, no skill, no model calls
qualitymd evaluation create && qualitymd result run --all && qualitymd evaluation report --fail-on fail
```

## Conventions

These apply across all commands; per-command flags live in the detail docs.

- **Singular commands, plural paths.** Resource commands are singular
  (`evaluation create`); the on-disk collection is plural
  (`.quality/evaluations/`).
- **Implicit current run.** `evaluation` and `result` commands default to the
  living run for the current model + target; an `<id>` is only needed for
  historical or archive operations.
- **Strictly resource-based.** The CLI exposes state; it never emits a prompt and
  never assumes an iteration order. Prompt composition and ordering are the
  skill's job — there is no `next` cursor.
- **Model file resolution.** `-f, --file <path>` selects a single `QUALITY.md`;
  default `./QUALITY.md`. `lint` also accepts the file as a positional argument
  and `-` for stdin. With `-f` omitted against a directory holding more than one
  `QUALITY.md`, a command operates on the whole set as a *federation* (discovery,
  composition, a tree-shaped result); see [`cli-federation.md`](./cli-federation.md).
- **Paths are model-relative.** `target` and `prompt` paths inside a `QUALITY.md`
  resolve relative to that file's directory, not the working directory.
- **`--json` for automation.** Available on every command and emits a
  schema-stable, versioned object (a `schemaVersion` field) so a skill or CI
  harness can parse results without screen-scraping. `lint` and `model show`
  default to JSON.
- **Interactivity.** A command runs **non-interactive** — never blocks on input —
  if stdin/stdout is not a TTY, `--non-interactive` is passed, or `--json` is
  passed. Otherwise it **may** prompt; today only [`init`](./cli-init.md) does.
  Any prompt is written to **stderr**, never contaminating stdout.

### Exit codes

A non-zero exit from a *gate* is distinct from a non-zero exit caused by the tool
failing to run, so a caller can tell "the quality is bad" from "the command
broke." One shared three-code convention:

| Code | Meaning |
| --- | --- |
| `0` | Success — the command ran and the gate (if any) passed. Inspection and report-only runs that complete also exit `0`. |
| `1` | **Gate verdict failure** — the command ran fine, but the bar was not met: `lint` found an `error`, or `evaluation report --fail-on` tripped. "The quality is bad." |
| `2` | **Tool failure** — bad flags, unreadable/absent file, internal error. "The command broke." |

`result run` rating a requirement poorly is **not** a gate trip — it ran fine and
records a low rating, exiting `0`. The gate is `evaluation report --fail-on`,
separately. Detail docs restate the codes a command emits but never reassign
these meanings.

### Advisory output (optional)

A command's `--json` output **may** carry an advisory `nextActions` array —
ordered, machine-readable suggestions (a runnable command, a reason, a priority).
These are *advisory only*: they never change exit status, and ordering the actual
work is the skill's job, not the CLI's. They exist so a human or harness reading
raw output has a hint, not so the CLI drives a workflow.

## The `.quality/` home

`./.quality/` is a model's quality home: the committed evaluation runs, plus
project configuration. A project with a single `QUALITY.md` has exactly one,
beside that file — the layout below.

```
.quality/
  config.yaml
  evaluations/
    <slug>/                # the living run for one (model, target)
    archive/
      <name>/              # manual snapshots (eval archive --as <name>)
```

```yaml
# ./.quality/config.yaml  (illustrative)
defaults:
  failOn: fail
```

**In a federation** this layout repeats per model — each `QUALITY.md` keeps its
own `.quality/` holding *its* runs, resolved relative to that model's directory,
so a component's evaluation travels with the component. Configuration does **not**
repeat: `config.yaml` and any shared rating scale are federation-level, resolved
once at the **scan root**'s `.quality/`. That scan-root home does double duty —
federation config plus, when a root model sits there, that root model's own runs —
exactly as a lone model's `.quality/` already carries both. *Runs decentralize;
configuration centralizes* (see [`cli-federation.md`](./cli-federation.md#output)).

The **meta-model** — the built-in `QUALITY.md` whose subject is a project's *own*
`QUALITY.md` — ships bundled with the CLI as a normal `QUALITY.md`-schema file. It
is what `improve-quality-md` uses as its model in its **diagnose** phase (see
[`skills.md`](./skills.md#improve-quality-md)). For now it is **fixed and not
configurable**: there is no project-local override or extension. (A
replace/extend mechanism is a possible later addition — see
[Open questions](#open-questions).) `lint` is governed by the spec, not this
config, so it needs nothing here.

Everything under `.quality/` is committed (see
[The evaluation lifecycle](#the-evaluation-lifecycle)), so git — not a retention
policy — manages history.

## Open questions

- **The CLI ↔ skill interface payloads.** The shapes of `result show` (what a
  skill needs to judge a `prompt`), `result set` (the verdict / diffable
  artifact), and the staleness hash are the real contract. They are now defined —
  with illustrative, provisionally-named schemas — authoritatively in
  [`cli-evaluate.md`](./cli-evaluate.md#the-interface-payloads-cli--skill-contract);
  what remains open is field-name tuning and the staleness-hash serialization.
- **Postable CI summaries.** Whether the CLI emits a human-facing PR comment /
  check annotation directly (a dedicated flag) or leaves that to the CI harness
  consuming `--json`. Relevant to `evaluation report --fail-on`.
- **`spec` verb (design.md prior art).** Printing the format spec for injection
  into agent/skill prompts is a plausible ergonomic addition. Out of scope for v1.
- **Configurable meta-model.** The diagnostic meta-model is fixed and bundled for
  now. Whether to later allow a project-local model to replace or extend it (and
  how — append vs. deep-merge) is deferred.
