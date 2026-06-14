# CLI: the evaluation lifecycle (`evaluation` / `result`)

> **Status:** rewritten for the **deterministic evaluation lifecycle**. This doc
> replaces the superseded "CLI compiles a prompt and runs a coding agent" model
> (the inversion — see [`cli.md`](./cli.md) and [`skills.md`](./skills.md)). It is
> the detail doc for the `evaluation` (alias `eval`) and `result` resources: their
> data model, on-disk layout, per-command behavior, the `bash` execution path, and
> the report rollup. Field names and the staleness hash are illustrative and
> expected to firm up in implementation; genuinely open items are marked `TODO` or
> collected under [Open questions](#open-questions).

This document covers the two resources that manage a recorded evaluation:

| Resource | What it is | Where detail lives |
| --- | --- | --- |
| `evaluation` (`eval`) | A **living per-target run** — one per (model, target), re-run in place. | here |
| `result` | A **requirement-result within a run** — one per selected requirement. | here |

It is a sibling of [`cli.md`](./cli.md) (the umbrella command surface),
[`skills.md`](./skills.md) (the judgment/orchestration layer that drives these
resources), and [`cli-compare.md`](./cli-compare.md) (the fuller, pending-rewrite
treatment of `evaluation compare`).

## The boundary this doc lives on

`qualitymd` draws one hard line (see [`cli.md`](./cli.md#the-split-deterministic-cli-judgment-in-skills)):
**the CLI is deterministic and never calls a model.** Everything in this doc is
that deterministic surface — it parses the model, resolves targets, runs `bash`
assessments and classifies them, persists the run, rolls up factors, renders the
report, and diffs runs. The **judgment** — composing prompts, judging `prompt`
assessments, deciding when evidence is sufficient — belongs to the skill layer
([`skills.md`](./skills.md)) and is cross-linked, never duplicated, here.

The line falls on the format's two assessment kinds. A `bash` requirement is
*computational*: the CLI runs it and classifies the result (`result run`). A
`prompt` requirement is *inferential*: a skill judges it and records the verdict
(`result set`). A model made **entirely of `bash` requirements is evaluable with
no skill and no model calls** — pure, reproducible CI.

## Data model

### The run

A run is the durable record of evaluating **one model against one target**. There
is exactly one *living* run per (model, target); re-running re-enters it in place
rather than forking a new one. Git history — not a chain of immutable run objects
— is the timeline (see [Why no finalize](#why-no-finalize-git-is-the-audit-layer)).

A run carries:

- its **identity** — the (model, target) pair, and a derived `slug` (the on-disk
  directory name);
- the **selected requirement set** — every requirement enumerated at
  `evaluation create`, each with a `result`;
- a derived **status** and **rollup** (computed, not stored as truth — see
  [Run states](#run-states)).

### The result

A result is the record for **one requirement within a run**. It carries:

- the requirement's full **path** (factor → … → requirement);
- its **assessment kind** (`bash` / `prompt`), carried from the model;
- a **state** (see [Result states](#result-states));
- once assessed, a **rating** (a level from the model's `ratings` scale) and
  **evidence** — for `bash`, the captured `result` fields and which `bashCondition`
  matched; for `prompt`, the skill's structured evidence;
- **provenance** for staleness: the hash inputs captured when the rating was
  recorded (see [Staleness](#staleness)).

The result is the diffable artifact. Its rating + evidence are what a PR reviewer
reads; volatile metadata is segregated out (see [On-disk layout](#on-disk-layout)).

## Run states

A run's status is **derived**, never authored:

| Status | Meaning | How reached |
| --- | --- | --- |
| `open` | At least one result is still `pending`. | `evaluation create`; any result returned to `pending`. |
| `complete` | Nothing is `pending` — every result is `recorded`, `skipped`, or `errored`. | derived once the last `pending` clears. |
| `archived` | A frozen snapshot. | `evaluation archive --as <name>`. |

`complete` is a watermark, not a seal — a complete run stays fully mutable, and a
re-run that marks results `stale` (then `pending`) drops it back to `open`. There
is no `finalize`.

## Result states

| State | Meaning | Entered from | Via |
| --- | --- | --- | --- |
| `pending` | Enumerated, not yet assessed. | (initial); `stale`; `result reset` | `evaluation create`, `result reset` |
| `recorded` | Assessed, carries a rating + evidence. | `pending` | `result run` (bash), `result set` (prompt) |
| `skipped` | Deliberately not assessed, carries a reason. | `pending` | `result skip` |
| `errored` | A `bash` command **could not run** (or its `bashCondition` failed to evaluate). Not a low rating. | `pending` | `result run` |
| `stale` | Was `recorded`, but the model subtree or target it was judged against has changed. | `recorded` | detected on `evaluation create` / re-run |

Two transition rules carry weight:

- **On re-run, only `stale` results return to `pending`.** `recorded`, `skipped`,
  and `errored` results are left intact unless their hash inputs changed — that is
  what keeps re-running cheap and the diff small. A fresh `evaluation create`
  re-hashes every result; any whose inputs moved become `stale → pending`.
- **`errored` is not a rating.** It means the command itself failed — non-existent
  binary, a `bashCondition` that does not evaluate (e.g. `.json()` on non-JSON
  output, a CEL type error). It is a *model/environment* problem, distinct from a
  command that **ran and matched a low level** (that is a legitimate `recorded`
  result with a poor rating). See [The bash path](#the-bash-path-execution--classification).

```text
                         result run (ran, classified)
            ┌──────────────────────────────────────────────┐
            ▼                                               │
        recorded ──── result reset ───► pending ◄───────────┤
            │                              ▲                │
   inputs changed                          │ result run / set
       (re-hash)                           │ (skipped/errored too)
            ▼                              │
          stale ──── on re-run ────────────┘
```

## On-disk layout

A living run is stored under a slug directory; archives live beside it. Everything
under `.quality/` is committed (see [`cli.md`](./cli.md#the-quality-home)).

```text
.quality/
  config.yaml
  evaluations/
    <slug>/                  # the living run for one (model, target)
      evaluation.json        # manifest: identity, selected set, rollup, verdict
      results/
        <req-id>.json        # one per result — rating + evidence (the diffable artifact)
      report.md              # rendered human-facing report (deterministic)
      .run/                  # SEGREGATED volatile metadata — see below
        meta.json            # timestamps, durations, host, tool version
    archive/
      <name>/                # frozen snapshot (evaluation archive --as <name>)
```

The `<slug>` derives from the (model, target) pair so the same pair always
re-enters the same directory. `<req-id>` derives deterministically from the
requirement's full path (so renaming a requirement is a visible add/remove in the
diff, not an in-place edit).

### Segregating volatile metadata for clean diffs

The **primary design goal is that the evaluation is a reviewable PR artifact** —
a diff should show *only* what a reviewer cares about: which ratings and evidence
changed. To get that:

- `evaluation.json` and `results/<req-id>.json` hold **verdicts only** — rating,
  evidence, the matched condition, the staleness provenance. Their serialization
  is **deterministic** (stable key order, normalized whitespace, sorted
  collections) so re-running an unchanged result produces a byte-identical file
  and an empty diff.
- All **volatile metadata** — wall-clock timestamps, command durations, host,
  tool version — is written to `.run/meta.json`, *segregated* from the verdicts.
  A reviewer can ignore `.run/` in review; a re-run touches it every time without
  noising up the verdict diff.

`TODO`: whether `.run/` is committed alongside verdicts (full audit, noisier
history) or gitignored (cleanest diffs, weaker provenance) — the design leans
toward committed-but-ignored-in-review, but this is not finally settled.

## `evaluation` commands

All commands default to the **living run for the current model + target**; an
explicit `<id>` is only needed for historical, archive, or compare operations
(see [`cli.md`](./cli.md#conventions)).

### `evaluation create [--model <path>] [--target <path>] [--from <id>]`

Create or re-enter the living per-target run and enumerate its requirements.

- Parses the model, resolves targets, and **enumerates every selected requirement
  as a `result`**. A first run starts every result `pending`.
- On an **existing** run, it re-enters in place: it re-hashes each `recorded`
  result and marks any whose inputs changed `stale → pending`. Untouched results
  keep their state. This is the re-run path; it never forks a new directory.
- `--from <id>` **carries forward** still-valid results from another run (an
  archive, or a prior run for a different target): for each requirement present in
  both, if the source result's hash inputs still match, its rating + evidence are
  copied in as `recorded`; otherwise the requirement is left `pending`. Carry-forward
  is how a new target run reuses still-applicable verdicts instead of re-judging
  from scratch.

Transition: → **open** (or → **complete** if `--from` carried forward a full set).

### `evaluation list`

List runs — the active run plus archived snapshots — with status and verdict. No
transition.

### `evaluation show [<id>] [--json]`

The run manifest: status, the selected requirement set with each result's state
and rating, and the rolled-up verdict (see [Report rollup](#report-rollup)). Read-only.

### `evaluation report [<id>] [--fail-on <level>] [--json]`

Render the report from the recorded results and, when `--fail-on` is set, gate.

- Renders `report.md` (the deterministic human-facing artifact — see
  [The report](#the-report)) and, under `--json`, the machine-readable rollup.
- `--fail-on <level>` exits **non-zero** when the overall rolled-up rating lands
  **at or below** `<level>` on the model's scale. **Off by default** — without it,
  `report` is inspection and exits `0` regardless of rating. This is the **single
  gate** in the lifecycle (see [Exit codes](#exit-codes)).

### `evaluation archive [<id>] --as <name>`

Snapshot the run's current state to `.quality/evaluations/archive/<name>/`. The
living run is untouched and stays mutable; the archive is a frozen copy for
comparison or record. Transition: → **archived** (the snapshot; the live run stays
`open`/`complete`).

### `evaluation delete <id>`

Discard a run (living or archived). Transition: → abandoned.

### `evaluation compare <a> <b> [--json]`

A **deterministic diff of two stored runs** — see [Compare](#compare).

## `result` commands

A requirement is addressed by its full dotted path (`<req>`). Commands default to
the current living run.

### `result list [--status pending,stale,…] [--json]`

Query results by state. `--status` takes a comma-separated set
(`pending,stale,recorded,skipped,errored`). There is **no `next` cursor** — the
CLI never orders the work or composes a prompt; the skill does
(see [`skills.md`](./skills.md#orchestration-contract)). No transition.

### `result show <req> [--json]`

The **resolved data** for one requirement — everything a skill needs to judge a
`prompt` requirement, with the CLI emitting no prompt of its own. No transition.

This payload is the **authoritative CLI ↔ skill contract**; the schema below is
*illustrative — field names are provisional* and expected to be tuned in
implementation (see [Interface payloads](#the-interface-payloads-cli--skill-contract)).
`cli.md` and `skills.md` cross-reference this section rather than re-specifying it.

### `result run <req | --all>`

Execute and classify `bash` assessment(s). **`bash` only** — pointing it at a
`prompt` requirement is a usage error (exit `2`). For each target requirement, the
CLI runs the command and classifies the captured `result` against the scale's
`bashCondition`s (see [The bash path](#the-bash-path-execution--classification)).

Transition: `pending → recorded` (command ran, a level matched) **or**
`pending → errored` (command could not run, or a `bashCondition` failed to
evaluate). A poor-but-legitimate rating is `recorded`, **not** a gate trip — it
exits `0` (the gate is `evaluation report --fail-on`).

### `result set <req> --rating <level> --evidence …`

Record a **`prompt` verdict** — the skill's judgment. Writes the rating level (a
declared scale level) and structured evidence to the result file. This is the
**diffable artifact** whose schema *is* the PR-review experience; the input
schema is defined authoritatively below
([Interface payloads](#the-interface-payloads-cli--skill-contract)).
Transition: `pending → recorded`.

### `result skip <req> --reason …`

Deliberately mark a requirement not-assessed, with a recorded reason. Counts as
non-`pending` for run completeness. Transition: `pending → skipped`.

### `result reset <req>`

Return a result to `pending`, discarding its rating/evidence, to be re-judged or
re-run. Transition: → **pending**.

## The interface payloads (CLI ↔ skill contract)

Three payloads constitute the open contract between the deterministic CLI and the
judging skill: the **`result show` output** a skill consumes to judge a `prompt`
requirement, the **`result set` input** that records its verdict, and the
**staleness hash** that decides when a recorded verdict is re-opened. These are
specified **authoritatively here**; [`cli.md`](./cli.md#open-questions) and
[`skills.md`](./skills.md#the-cli--skill-interface) cross-reference this section
rather than restate it.

The schemas below are **illustrative — field names are provisional** and expected
to firm up in implementation, consistent with this spec's stance on field naming.
What is contractual is the *information* each payload carries, not the exact keys.

### `result show` output

Everything a skill needs to judge one `prompt` requirement, fully resolved by the
CLI so the skill performs no model parsing or glob expansion of its own:

- **`requirementPath`** — the full dotted requirement path, and **`factorPath`** —
  the factor → … chain it sits under (for grouping and rollup context).
- **`assessmentKind`** — `prompt` or `bash`. (`result show` resolves both; a skill
  judges only `prompt`, and routes `bash` to `result run`.)
- **`assessment`** — the resolved assessment text: the loaded **`prompt`** text
  for a `prompt` requirement, *or* the **`bash`** command for a `bash` requirement
  (whichever the kind selects).
- **`target`** — the **resolved target manifest**: the list of files the target
  glob expands to, each with its model-relative `path`. File **contents are
  optional / by reference** — included only when the CLI is asked to inline them;
  otherwise the skill reads the listed paths itself (keeps the payload small and
  the skill in control of what it loads).
- **`ratings`** — the in-scope rating **scale**: its levels ordered **best → worst**,
  each with the `promptCondition` (for `prompt`) or `bashCondition` (for `bash`)
  that defines it.
- **`ratingOverrides`** — any
  [per-requirement rating overrides](../SPECIFICATION.md#per-requirement-rating-overrides)
  that replace the default scale's bands for this requirement (absent when none).
- **`state`** — the result's current [state](#result-states) (`pending`, `stale`, …).
- **`sufficiency`** — **"done" guidance for `prompt` judging**: how to decide when
  evidence is sufficient to rate (saturation). Provisionally a short prose note
  plus any model-supplied hints; the judging methodology that fills this in is the
  skill's (see [`skills.md`](./skills.md#the-cli--skill-interface)).

Illustrative JSON:

```json
{
  "schemaVersion": 1,
  "requirementPath": "security.input-validation.rejects-malformed-input",
  "factorPath": ["security", "input-validation"],
  "assessmentKind": "prompt",
  "assessment": "Assess whether the request handler rejects malformed input before it reaches business logic, with a clear error and no partial side effects.",
  "target": {
    "glob": "./src/handlers/**/*.ts",
    "files": [
      { "path": "src/handlers/intake.ts" },
      { "path": "src/handlers/validate.ts" }
    ]
  },
  "ratings": {
    "levels": [
      { "level": "pass", "promptCondition": "All untrusted inputs are validated at the boundary and rejected with a clear error." },
      { "level": "weak", "promptCondition": "Most inputs are validated, but at least one path reaches logic unchecked." },
      { "level": "fail", "promptCondition": "Inputs are not systematically validated at the boundary." }
    ],
    "order": "bestToWorst"
  },
  "ratingOverrides": null,
  "state": "pending",
  "sufficiency": "Rate once you have traced every handler in the target manifest; do not rate from a single file if others remain unread."
}
```

### `result set` input

The verdict a skill records for a `prompt` requirement — **this is the diffable
artifact**, so its on-disk serialization has a **stable field order** and carries
**no volatile metadata inline** (timestamps, durations, and host live in `.run/`;
see [Segregating volatile metadata](#segregating-volatile-metadata-for-clean-diffs)).
It carries:

- **`requirement`** — the requirement ref (full dotted path) being recorded.
- **`rating`** — the chosen **level**, which MUST be one of the scale's declared
  levels (validated against the scale resolved in `result show`).
- **`evidence`** — structured support for the rating: a **`summary`** (one- or
  two-line justification) plus an **`items`** array, each item either a
  `file:line` **location** (a citation into the target) or a free-text **note**.
- **`rationale`** *(optional)* — longer reasoning beyond the summary, when the call
  is non-obvious.

Illustrative JSON (the recorded artifact; field order is stable):

```json
{
  "schemaVersion": 1,
  "requirement": "security.input-validation.rejects-malformed-input",
  "rating": "weak",
  "evidence": {
    "summary": "Boundary validation covers JSON bodies but query params reach logic unchecked.",
    "items": [
      { "location": "src/handlers/intake.ts:42", "note": "validates and rejects malformed JSON body" },
      { "location": "src/handlers/validate.ts:17", "note": "query params parsed but never validated before use" },
      { "note": "no shared validation middleware across handlers" }
    ]
  },
  "rationale": "Bodies are guarded but the unvalidated query-param path is a real gap, so this rates weak rather than pass."
}
```

On the command line the same payload is supplied via `--rating <level>` and
repeated `--evidence` flags (or `--evidence-file` for the structured form); the
CLI normalizes it into the stable on-disk shape above.

### Staleness hash

A recorded result becomes [`stale`](#result-states) when what it was judged
against changes. The hash is computed over **two inputs**, and a recorded result
goes `stale` if **either** changes:

1. **The requirement's resolved definition** — its name / full path, the
   **resolved target glob set** (the patterns, so a manifest-changing glob edit
   counts), the **assessment text** (`prompt` or `bash`), and the **applicable
   rating conditions** (the in-scope scale levels plus any per-requirement
   overrides). If the *requirement itself* moves, the prior verdict no longer
   describes the same question.
2. **The resolved target contents** — a hash over the **expanded file set** the
   target glob resolves to (the manifest from [`result show`](#result-show-output)
   and the bytes of those files). If the code under the requirement changes, the
   prior verdict no longer describes the current subject.

`evaluation create` re-computes both on re-entry; a mismatch on either marks the
result `stale` (and, on re-run, `stale → pending`). A non-`stale` `recorded`
result is left untouched — this is what keeps re-running cheap (see
[Staleness](#staleness)).

`TODO`: the **precise serialization / canonicalization** is unsettled — the exact
byte-level normalization of input (1) (key ordering, whitespace), and whether the
target contribution (2) hashes raw file *contents* (precise, but a whitespace edit
re-opens everything) or a coarser signal (manifest + revision/mtime). It is also
open whether a `prompt` result should hash anything the skill saw beyond the
resolved target. Tracked under [Open questions](#open-questions).

`result run` is the whole deterministic computational tier. For a `bash`
requirement it:

1. **Runs the command** in the model's directory (paths are model-relative — see
   [`cli.md`](./cli.md#conventions)), capturing the `result` fields the format
   defines: `result.success` (zero exit), `result.exit`, `result.stdout`,
   `result.stderr` (see `../SPECIFICATION.md#computational-rating`).
2. **Classifies** the `result` against the scale's `bashCondition`s — CEL booleans
   evaluated best-to-worst; the **first** matching level wins; if none match, the
   **worst** level is the fallback (the scale denies by default). A requirement may
   carry a [per-requirement override](../SPECIFICATION.md#per-requirement-rating-overrides)
   supplying its own bands. So a verdict is **not** merely "exit zero" — a scale may
   band a numeric value the command prints (`double(result.stdout.trim()) >= 90`);
   only the default `pass`/`fail` scale reduces to "best on a zero exit."
3. **Records** the matched level as the result's rating, with the captured
   `result` fields and the matched condition as evidence.

### `errored` vs. a legitimate low rating

This distinction is the crux of the bash path:

| Outcome | State | Meaning |
| --- | --- | --- |
| Command ran; a level matched (even the worst). | `recorded` | A real classification. A `fail` here is the subject genuinely failing the condition — a valid verdict, exit `0`. |
| Command **could not run** (binary not found, non-zero from the runner itself, timeout). | `errored` | No trustworthy verdict — an environment problem. |
| A `bashCondition` **could not evaluate** (`.json()` on non-JSON, `double()` on non-numeric, CEL type error). | `errored` | A **model configuration** error, surfaced as such — never silently scored. |

The first row is the common, intended case: `bash` assessments are *supposed* to
produce ratings, including bad ones. `errored` is reserved for "the measurement
itself broke." A `result run` that produces a low rating still exits `0`; the gate
is separate.

## Staleness

When `evaluation create` re-enters a run, each `recorded` result is re-hashed; if
its hash inputs changed, it becomes `stale` (and on re-run, `stale → pending`).
Staleness is what keeps a re-run from re-judging unchanged work while catching work
the change invalidated.

The hash inputs — the requirement's resolved definition and the resolved target
contents — and the open question of their precise serialization are specified
authoritatively under
[The interface payloads → Staleness hash](#staleness-hash). This section is the
*behavioral* view (when re-hashing happens and what it does to state); that one is
the *contract* view (what is hashed).

## `--from` carry-forward

`--from <id>` (on `evaluation create`) seeds a new or re-entered run from another
run's results. For every requirement common to both, the source's rating + evidence
are copied in as `recorded` **iff** the source result's staleness hash still holds
against the new run's model/target; otherwise the requirement starts `pending`.
This makes a new-target run reuse still-valid verdicts rather than re-judging from
scratch — the carry-forward analogue of in-place re-running, across runs.

## Report rollup

`evaluation report` rolls per-requirement ratings up to factors, sub-factors, and
an overall verdict, **deterministically** — no judgment. Each requirement
contributes its recorded rating (a level on the model's scale); a factor's rating
is the rollup of its requirements and sub-factors.

- The rollup method is **worst-wins by default** (a factor is no better than its
  weakest requirement), matching a quality *gate*'s intent.
- `skipped` results are excluded from the rollup but reported as gaps in coverage.
- `pending` / `stale` / `errored` results mean the rollup is **incomplete**; the
  report states this rather than scoring around it, so a partial run is never
  mistaken for a clean one.

`TODO`: whether the rollup is configurable (e.g. weighted, or "any-fail-fails" vs.
a banded average) — for now, worst-wins, deterministic, stated in the report.

### The report

`report.md` is the deterministic human-facing artifact, rendered from recorded
results (no model authoring it — the inversion). It carries:

- **Definition** — model, target, scope, requirement set, timestamp (the timestamp
  comes from `.run/`, not the verdict files).
- **Verdict** — overall rolled-up rating, and whether a `--fail-on` gate would trip.
- **Per-factor rollup** — rating per factor / sub-factor.
- **Per-requirement results** — rating, the matched `bashCondition` or the skill's
  evidence, and `path`/location where applicable.
- **Coverage** — `skipped` requirements (with reasons) and any `pending` / `stale`
  / `errored` results that leave the run incomplete.

It is rendered, deterministic prose; under `--json`, `evaluation report` emits the
same rollup as a schema-stable object.

## Compare

`evaluation compare <a> <b>` is a **deterministic diff of two stored runs** — no
re-evaluation, no model calls. It reads both runs' recorded results and reports,
per requirement, what differs:

- rating changes (`a → b`), and the **direction** (improved / regressed / same on
  the scale's ordering);
- requirements present in one run but not the other (model or target drift);
- state differences (e.g. `recorded` in one, `skipped`/`errored`/`pending` in the
  other).

Either side may be the living run, an archive (`archive/<name>`), or a git revision
of the run files — git history *is* the timeline, so "compare against last week" is
"compare against that revision."

This section is intentionally light: `evaluation compare` is the **same-model,
two-run** diff. The fuller comparison story — **multiple targets against one shared
model**, the requirements × targets matrix, A/B and regression gating — lives in
[`cli-compare.md`](./cli-compare.md). Treat that doc as authoritative for compare;
this section specifies only the minimal two-run diff the `evaluation` resource
exposes.

## Why no finalize: git is the audit layer

There is deliberately **no `finalize`/`seal`** transition. A run is **always
mutable**; correctness of the record is provided by **git commits**, not an
immutable in-tool state:

- **Commit everything** under `.quality/` — the manifest, the per-result verdicts,
  the rendered report. The commit history is the audit trail.
- Because verdict serialization is deterministic and volatile metadata is
  segregated, a commit's diff is exactly the rating/evidence changes a reviewer
  should see — the run **is** a reviewable PR artifact.
- `evaluation archive --as <name>` exists for the cases where a *named*, frozen
  snapshot is wanted independent of git (a labeled baseline to compare against),
  not as a finalize step.

## Exit codes

The lifecycle reuses the shared three-code convention verbatim in meaning (see
[`cli.md`](./cli.md#exit-codes)):

| Code | Meaning | In this lifecycle |
| --- | --- | --- |
| `0` | Success — ran, and the gate (if any) passed. | Inspection (`show`, `list`), `result run`/`set`/`skip` (even when the recorded rating is poor), and `evaluation report` *without* `--fail-on`. |
| `1` | **Gate verdict failure** — ran fine, the bar was not met. | `evaluation report --fail-on <level>` and the overall rating landed at or below `<level>`. The *only* `1` in this lifecycle. |
| `2` | **Tool failure** — the command broke. | Bad flags, an unresolvable model/target, a `QUALITY.md` that does not parse (run `lint` first), `result run` on a `prompt` requirement, an `errored` measurement surfaced as a tool problem, internal error. |

The load-bearing point: **a poorly-rated `result run` exits `0`.** Running a `bash`
assessment that classifies to a low level is a *successful* measurement — it ran
and recorded a verdict. The quality gate is `evaluation report --fail-on`,
separately, so "the command broke" (`2`), "the quality is bad" (`1`, opt-in), and
"recorded a low rating" (`0`) stay distinguishable.

## Open questions

- **Staleness hash serialization** — the precise canonicalization of the hash
  inputs, the contents-vs-coarse target signal, and what a `prompt` result hashes
  beyond the resolved target (see [Staleness hash](#staleness-hash)). The
  authoritative definition lives here; [`cli.md`](./cli.md#open-questions) and
  [`skills.md`](./skills.md#the-cli--skill-interface) cross-reference it.
- **`.run/` commit policy** — committed-but-ignored-in-review vs. gitignored (see
  [Segregating volatile metadata](#segregating-volatile-metadata-for-clean-diffs)).
- **Rollup method** — whether worst-wins is configurable
  (see [Report rollup](#report-rollup)).
- **Interface payload field names** — the `result show` / `result set` schemas are
  defined [here](#the-interface-payloads-cli--skill-contract) but their field names
  are provisional and expected to be tuned in implementation.
- **Compare depth** — how much of the multi-target / matrix story `evaluation
  compare` carries vs. defers to [`cli-compare.md`](./cli-compare.md) once that is
  rewritten.
- **Federation** — how `evaluation` / `result` address a federated tree of models
  (one run per model, a tree-shaped report); needs the federation rewrite (see
  [`cli-federation.md`](./cli-federation.md)).
