# CLI: the evaluation lifecycle (`evaluation` / `result`)

> **Status:** rewritten for the **deterministic evaluation lifecycle**. This doc
> replaces the superseded "CLI compiles a prompt and runs a coding agent" model
> (the inversion — see [`cli.md`](./cli.md) and [`skills.md`](./skills.md)). It is
> the detail doc for the `evaluation` (alias `eval`) and `result` resources: their
> data model, on-disk layout, per-command behavior, and the report rollup. Field
> names and the staleness hash are illustrative and
> expected to firm up in implementation; genuinely open items are marked `TODO` or
> collected under [Open questions](#open-questions).

This document covers the two resources that manage a recorded evaluation:

| Resource | What it is | Where detail lives |
| --- | --- | --- |
| `evaluation` (`eval`) | A **living per-target run** — one per (model, target), re-run in place. | here |
| `result` | A **requirement-result within a run** — one per selected requirement. | here |

It is a sibling of [`cli.md`](./cli.md) (the umbrella command surface) and
[`skills.md`](./skills.md) (the judgment/orchestration layer that drives these
resources).

## The boundary this doc lives on

`qualitymd` draws one hard line (see [`cli.md`](./cli.md#the-split-deterministic-cli-judgment-in-skills)):
**the CLI is deterministic and never calls a model.** Everything in this doc is
that deterministic surface — it parses the model, resolves targets, persists the
run, records verdicts, rolls up factors, and renders the report. The
**judgment** — composing prompts, judging the requirements, deciding when
evidence is sufficient — belongs to the skill layer ([`skills.md`](./skills.md))
and is cross-linked, never duplicated, here.

The line is the mechanics / judgment boundary. Every assessment is an inferential
`prompt`: a skill judges it and records the verdict (`result set`), while the CLI
does the deterministic work around that — resolving the requirement, persisting
the verdict, and rolling it up.

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
- a **state** (see [Result states](#result-states));
- once assessed, a **rating** (a level from the model's `ratings` scale) and the
  skill's structured **evidence**;
- **provenance** for staleness: the hash inputs captured when the rating was
  recorded (see [Staleness](#staleness)).

The result is the diffable artifact. Its rating + evidence are what a PR reviewer
reads; volatile metadata is segregated out (see [On-disk layout](#on-disk-layout)).

## Run states

A run's status is **derived**, never authored:

| Status | Meaning | How reached |
| --- | --- | --- |
| `open` | At least one result is still `pending`. | `evaluation create`; any result returned to `pending`. |
| `complete` | Nothing is `pending` — every result is `recorded` or `skipped`. | derived once the last `pending` clears. |
| `archived` | A frozen snapshot. | `evaluation archive --as <name>`. |

`complete` is a watermark, not a seal — a complete run stays fully mutable, and a
re-run that marks results `stale` (then `pending`) drops it back to `open`. There
is no `finalize`.

## Result states

| State | Meaning | Entered from | Via |
| --- | --- | --- | --- |
| `pending` | Enumerated, not yet assessed. | (initial); `stale`; `result reset` | `evaluation create`, `result reset` |
| `recorded` | Assessed, carries a rating + evidence. | `pending` | `result set` |
| `skipped` | Deliberately not assessed, carries a reason. | `pending` | `result skip` |
| `stale` | Was `recorded`, but the model subtree or target it was judged against has changed. | `recorded` | detected on `evaluation create` / re-run |

Two transition rules carry weight:

- **On re-run, only `stale` results return to `pending`.** `recorded` and
  `skipped` results are left intact unless their hash inputs changed — that is
  what keeps re-running cheap and the diff small. A fresh `evaluation create`
  re-hashes every result; any whose inputs moved become `stale → pending`.
- **A low rating is still a `recorded` result.** A requirement the skill judges to
  land at the worst level is a legitimate verdict — `recorded` with a poor rating —
  not a separate failure state. The gate is `evaluation report --fail-on`, applied
  separately to the rollup.

```text
                          result set (judged, recorded)
            ┌──────────────────────────────────────────────┐
            ▼                                               │
        recorded ──── result reset ───► pending ◄───────────┤
            │                              ▲                │
   inputs changed                          │ result set
       (re-hash)                           │ (skipped too)
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
explicit `<id>` is only needed for historical or archive operations
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

Render the report from the recorded results and gate on the rolled-up verdict.

- Renders `report.md` (the deterministic human-facing artifact — see
  [The report](#the-report)) and, under `--json`, the machine-readable rollup.
- `--fail-on <level>` exits **non-zero** when the overall rolled-up rating lands
  **at or below** `<level>` on the model's scale. It **defaults to `unacceptable`**
  — so under the default scale `report` gates only on `unacceptable` (a rating
  below the acceptable floor), treating anything at or above `minimum` as releasable
  ("at or above Minimum is acceptable"). `--fail-on <level>` tightens the bar (e.g.
  `--fail-on minimum`). This is the **single gate** in the lifecycle (see
  [Exit codes](#exit-codes)).

### `evaluation archive [<id>] --as <name>`

Snapshot the run's current state to `.quality/evaluations/archive/<name>/`. The
living run is untouched and stays mutable; the archive is a frozen copy for
comparison or record. Transition: → **archived** (the snapshot; the live run stays
`open`/`complete`).

### `evaluation delete <id>`

Discard a run (living or archived). Transition: → abandoned.

## `result` commands

A requirement is addressed by its full dotted path (`<req>`). Commands default to
the current living run.

### `result list [--status pending,stale,…] [--json]`

Query results by state. `--status` takes a comma-separated set
(`pending,stale,recorded,skipped`). There is **no `next` cursor** — the
CLI never orders the work or composes a prompt; the skill does
(see [`skills.md`](./skills.md#orchestration-contract)). No transition.

### `result show <req> [--json]`

The **resolved data** for one requirement — everything a skill needs to judge a
`prompt` requirement, with the CLI emitting no prompt of its own. No transition.

This payload is the **authoritative CLI ↔ skill contract**; the schema below is
*illustrative — field names are provisional* and expected to be tuned in
implementation (see [Interface payloads](#the-interface-payloads-cli--skill-contract)).
`cli.md` and `skills.md` cross-reference this section rather than re-specifying it.

### `result set <req> --rating <level> --evidence …`

Record a **verdict** — the skill's judgment. Writes the rating level (a declared
scale level) and structured evidence to the result file. This is the **diffable
artifact** whose schema *is* the PR-review experience; the input schema is defined
authoritatively below
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
- **`assessment`** — the resolved **`prompt`** text the skill judges against.
- **`target`** — the **resolved target manifest**: the list of files the target
  glob expands to, each with its model-relative `path`. File **contents are
  optional / by reference** — included only when the CLI is asked to inline them;
  otherwise the skill reads the listed paths itself (keeps the payload small and
  the skill in control of what it loads).
- **`ratings`** — the in-scope rating **scale**: its levels ordered **best → worst**,
  each with the `promptCondition` that defines it.
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
      { "level": "outstanding", "promptCondition": "Exceeds the requirement; meets it with margin to spare" },
      { "level": "target", "promptCondition": "Meets the requirement" },
      { "level": "minimum", "promptCondition": "Falls short of the goal but stays at the acceptable floor" },
      { "level": "unacceptable", "promptCondition": "Falls below the acceptable floor" }
    ],
    "order": "bestToWorst"
  },
  "ratingOverrides": null,
  "state": "pending",
  "sufficiency": "Rate once you have traced every handler in the target manifest; do not rate from a single file if others remain unread."
}
```

### `result set` input

The verdict a skill records for a requirement — **this is the diffable
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
  "rating": "minimum",
  "evidence": {
    "summary": "Boundary validation covers JSON bodies but query params reach logic unchecked.",
    "items": [
      { "location": "src/handlers/intake.ts:42", "note": "validates and rejects malformed JSON body" },
      { "location": "src/handlers/validate.ts:17", "note": "query params parsed but never validated before use" },
      { "note": "no shared validation middleware across handlers" }
    ]
  },
  "rationale": "Bodies are guarded but the unvalidated query-param path is a real gap, so this rates minimum rather than target."
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
   counts), the **assessment text** (the `prompt`), and the **applicable rating
   conditions** (the in-scope scale levels plus any per-requirement overrides). If
   the *requirement itself* moves, the prior verdict no longer describes the same
   question.
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
open whether a result should hash anything the skill saw beyond the resolved
target. Tracked under [Open questions](#open-questions).

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
- `pending` / `stale` results mean the rollup is **incomplete**; the report states
  this rather than scoring around it, so a partial run is never mistaken for a
  clean one.

`TODO`: whether the rollup is configurable (e.g. weighted, or "any-fail-fails" vs.
a banded average) — for now, worst-wins, deterministic, stated in the report.

### The report

`report.md` is the deterministic human-facing artifact, rendered from recorded
results (no model authoring it — the inversion). It carries:

- **Definition** — model, target, scope, requirement set, timestamp (the timestamp
  comes from `.run/`, not the verdict files).
- **Verdict** — overall rolled-up rating, and whether a `--fail-on` gate would trip.
- **Per-factor rollup** — rating per factor / sub-factor.
- **Per-requirement results** — rating, the skill's evidence, and `path`/location
  where applicable.
- **Coverage** — `skipped` requirements (with reasons) and any `pending` / `stale`
  results that leave the run incomplete.

It is rendered, deterministic prose; under `--json`, `evaluation report` emits the
same rollup as a schema-stable object.

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
| `0` | Success — ran, and the gate passed. | Inspection (`show`, `list`), `result set`/`skip` (even when the recorded rating is poor), and `evaluation report` whose rolled-up rating cleared the `--fail-on` bar (by default, anything at or above `minimum`). |
| `1` | **Gate verdict failure** — ran fine, the bar was not met. | `evaluation report` and the overall rating landed at or below `--fail-on <level>` (by default `unacceptable`). The *only* `1` in this lifecycle. |
| `2` | **Tool failure** — the command broke. | Bad flags, an unresolvable model/target, a `QUALITY.md` that does not parse (run `lint` first), internal error. |

The load-bearing point: **recording a low rating exits `0`.** A requirement the
skill judges to land at a low level is a *successful* record — it ran and recorded
a verdict. The quality gate is `evaluation report` (which gates on the rolled-up
rating, by default `unacceptable`), separately, so "the command broke" (`2`), "the
quality is bad" (`1`, at `report`), and "recorded a low rating" (`0`, at
`result set`) stay distinguishable.

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
- **Federation** — how `evaluation` / `result` address a federated tree of models
  (one run per model, a tree-shaped report); needs the federation rewrite (see
  [`cli-federation.md`](./cli-federation.md)).
