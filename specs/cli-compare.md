# CLI: `compare`

> Detail doc for the **comparison** mode of the semantic tier. See
> [`cli.md`](./cli.md) for the full command surface and shared conventions, and
> [`cli-evaluate.md`](./cli-evaluate.md) for the single-target evaluation engine
> this verb is built on.

| Command | Purpose | Cost / determinism | Output |
| --- | --- | --- | --- |
| `qualitymd compare <target> <target> [<target>...] [factor]` | Evaluate **multiple targets against one shared quality model** and rank/diff them per requirement. | Expensive, non-reproducible (× N targets). | Comparison bundle on disk. |

`compare` answers a **relative** question — *how do these targets rank against
the model, and what differs between them?* — where [`evaluate`](./cli-evaluate.md)
answers an **absolute** one — *how good is this subject?* That difference in
output contract (a requirements × targets **matrix**, not a single verdict) is
why it is a distinct verb rather than a flag on `evaluate`.

## Why this is a distinct mode of use

QUALITY.md projects tend to fall into one of three usage modes, and a single
`QUALITY.md` is unlikely to serve more than one of them at a time:

- **Single codebase** — the ordinary case. One subject, an absolute verdict →
  [`evaluate`](./cli-evaluate.md).
- **A/B comparison** — two instances of "the same" subject: `base` vs. `head`
  of a PR, old vs. new implementation, ours vs. a competitor. → `compare` with
  N = 2.
- **N-way analysis** — a benchmarking/selection effort: several candidate
  implementations, vendor packages, or forks measured against one rubric. →
  `compare` with N ≥ 2.

> **Subject vs. target.** A **target** is a concrete source the engine measures
> — a path or a git ref (see [target source](./cli-evaluate.md#target-source)). A
> **subject** is the system-under-study a target is an *instance* of. They are
> 1:1 in N-way analysis (each target is a different subject), but 1:many in A/B
> (one subject, two targets — `base` and `head`). `compare` always operates on
> **targets**; "subject" is reserved for the system those targets instantiate.

The A/B and N-way modes are the **same machinery** — the only difference is the
count of targets and how the result is framed. Crucially, the **regression
gate** ("did this PR make any requirement *worse* than base?") is just
`compare base head` with gate framing; it is not a separate feature and not a
`--diff` flag on `evaluate` (see
[`cli-evaluate.md`](./cli-evaluate.md#target-scope)).

## Conceptual model

Comparison is the [single-target engine](./cli-evaluate.md#the-evaluation-engine-shared)
run **N times with the target varied**, under **identical conditions**, plus a
final synthesis step. Formally: one criteria source, N `(target, criteria,
context)` triples sharing `criteria`, then a roll-up across targets.

```text
shared:  criteria  = resolved QUALITY.md (the ruler)
         factor    = same requirement subtree for every target
         rigor     = same level for every target
         design    = one design.md, applied to all targets

for each target T in targets:
    results[T] = engine(target=T, criteria, context)   # cli-evaluate.md pipeline

comparison = synthesize(results)   # matrix + ranking + per-requirement winner + deltas
```

### Validity rests on a shared ruler

A comparison is only meaningful if every target is measured the same way. Two
invariants are therefore load-bearing — violating either makes the ranking
meaningless rather than merely imprecise:

1. **One criteria source.** The same resolved `QUALITY.md`, the same `factor`,
   the same `ratings` scale across all targets. The model file is **subject-
   agnostic**: it declares requirements, never which code they apply to. That
   agnosticism is what lets one model rank many subjects.
2. **One evaluation regime.** The same `--rigor`, the same partitioning strategy,
   the same finder/refuter ensemble, the same saturation stop criteria for every
   target. `compare` builds **one `design.md`** and applies it to all targets,
   rather than designing each independently.

### Non-reproducibility is *amplified* — handle it explicitly

The single-target engine is already non-reproducible
([`cli-evaluate.md`](./cli-evaluate.md#iso-25040-process-one-shot)). Comparison
makes this worse: a rank difference between target A and target B may be
**finder noise, not a real quality difference.** A comparison that hides its own
noise floor is worse than no comparison. So `compare`:

- runs every target under the identical regime above;
- **surfaces per-requirement confidence** in the matrix, not just a winner;
- marks a per-requirement delta as **`tie` / `inconclusive`** when the rating
  gap is within the noise the confidence levels imply, rather than forcing a
  spurious winner;
- at `--rigor max`, *may* run each target more than once to report **variance**,
  so a reported difference can be distinguished from run-to-run jitter.

## Targets

Each positional `<target>` is a [target source](./cli-evaluate.md#target-source)
in its own right. A target may be:

- a **path** to a working tree (`./impl-a`, `../competitor`),
- a **git ref** of the current repo (`main`, `v1.2.0`, a SHA) — the A/B and
  regression cases, and
- *(candidate)* a remote repo reference, for cross-repo analysis.

The optional trailing `factor` (a dotted factor path, as in
[`evaluate`](./cli-evaluate.md#factor)) applies identically to every target.

Targets may be named for the report with `name=path` syntax
(`compare baseline=main candidate=feature-x`); otherwise a label is derived from
the path/ref.

### Declaring targets in config

A project whose *purpose* is comparison (an analysis/benchmarking effort) can
declare its default targets in `./.quality/config.yaml` so `compare` needs no
positional targets — matching the observation that comparison projects and
single-codebase projects are different projects with different ergonomics. The
`QUALITY.md` model file itself stays subject-agnostic regardless.

```yaml
# ./.quality/config.yaml  (illustrative)
compare:
  targets:
    - name: baseline
      ref: main
    - name: candidate
      path: ../candidate-impl
```

## Output: the comparison bundle

```text
./.quality/comparisons/<YYYYMMDD-HHMMSS>-<factor>/
  design.md         # the single shared Design step applied to every target
  plan.md           # the concrete task graph, per target
  targets/
    <label>/results.json   # full single-target Execute output per target
  comparison.json   # the matrix: ratings & confidence per (requirement × target)
  comparison.md     # the primary human-facing artifact — ranking, deltas, winner
```

Each `targets/<label>/results.json` is exactly a single-target
[`results.json`](./cli-evaluate.md#resultsjson-illustrative-shape), so a
comparison is fully decomposable into the individual evaluations that produced
it. `comparison.json`/`comparison.md` add only the cross-target synthesis.

### `comparison.md`

- **Definition** — the shared model, factor, rigor, the target set with labels
  and resolved sources, timestamp.
- **Ranking** — overall ordering of targets, with the overall rating each
  achieved and a confidence/variance note.
- **Matrix** — requirements (rows) × targets (columns), each cell a rating +
  confidence; per-row winner highlighted, ties marked.
- **Notable deltas** — requirements where targets diverge most, with evidence and
  `path:line` locations per target.
- **Per-target strengths & weaknesses** — a brief absolute read of each, so the
  comparison doesn't erase the standalone verdict.
- **Limitations** — rigor used, saturation achieved per target, contested/tied
  findings, variance if measured, and the amplified non-reproducibility caveat.

### `comparison.json` (illustrative shape)

```json
{
  "definition": {
    "factor": "security",
    "rigor": "high",
    "criteriaSource": "./QUALITY.md",
    "targets": [
      { "label": "baseline",  "source": { "ref": "main" } },
      { "label": "candidate", "source": { "path": "../candidate-impl" } }
    ]
  },
  "ranking": [
    { "label": "candidate", "overall": "pass", "confidence": "high" },
    { "label": "baseline",  "overall": "fail", "confidence": "high" }
  ],
  "matrix": {
    "no secrets committed to the repository": {
      "baseline":  { "rating": "fail", "confidence": "confirmed" },
      "candidate": { "rating": "pass", "confidence": "confirmed" },
      "winner": "candidate"
    },
    "secrets are loaded from a vault": {
      "baseline":  { "rating": "pass", "confidence": "contested" },
      "candidate": { "rating": "pass", "confidence": "contested" },
      "winner": "tie"
    }
  }
}
```

## Flags, exit codes

Shared flags are in [`cli.md`](./cli.md#shared-conventions); the deep-command
flags (`--rigor`, `--name`) carry over from
[`cli-evaluate.md`](./cli-evaluate.md#flags-exit-codes) and apply to **every**
target identically. Comparison-specific:

- `--baseline <label>` — treat one target as the reference and express the others
  as **deltas against it** (the A/B / regression framing). Without it, all targets
  are peers and the output is a symmetric ranking.
- `--fail-on-regression` — exit non-zero if any requirement rates **worse** in a
  non-baseline target than in the `--baseline` target. The comparison analogue of
  [`evaluate --fail-on`](./cli-evaluate.md#flags-exit-codes): a PR gate that
  blocks on *new* shortfalls without failing on pre-existing debt. Off by default
  (report-only), opt-in, and requires `--baseline`.

Exit codes:

- **Report-only by default** (exit `0` unless the run errors), as with
  [`evaluate`](./cli-evaluate.md#flags-exit-codes).
- **`--fail-on-regression`** opts into a regression gate, for the automated-PR
  use case. Opt-in for the same reason `evaluate --fail-on` is: the run is
  non-reproducible, so gating is a team's explicit choice, never imposed.

## Open questions

- **Remote/cross-repo targets.** Whether a target may be a remote repo reference
  (clone-and-evaluate) for cross-repo analysis, or whether v1 limits targets to
  local paths and refs of the current repo.
- **Variance budget at `max`.** Whether repeated-run variance reporting is
  in-scope for v1, and how many runs is enough to call a delta real without
  blowing up cost (N targets × R runs).
- **Tie threshold.** How the "within the noise floor → `tie`" decision is made —
  a fixed rating-distance rule, a confidence-overlap rule, or variance-derived.
- **Shared vs. per-target partitioning.** Comparable targets (two refs of one
  repo) can share a partitioning; structurally different targets (two unrelated
  codebases) may not. How much of the "one evaluation regime" invariant can hold
  when targets aren't structurally comparable.
- **Single-target `--diff` over a comparison.** Whether `compare` accepts a
  `--diff`-style scope per target, or comparison is always whole-target.
