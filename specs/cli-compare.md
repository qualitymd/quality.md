# CLI: `evaluation compare`

> **Status:** rewritten for the **deterministic compare** model. This doc replaces
> the superseded "CLI compiles a prompt and runs a coding agent over N targets"
> design (the inversion — see [`cli.md`](./cli.md) and [`skills.md`](./skills.md)).
> `evaluation compare` is now a **pure diff of already-stored evaluation runs**: it
> reads recorded results and reports what differs. It never re-judges, never runs
> `bash`, and never calls a model. Field names and output shapes are illustrative
> and expected to firm up in implementation; genuinely open items are under
> [Open questions](#open-questions).

This is the detail doc for `evaluation compare`. It is a sibling of
[`cli.md`](./cli.md) (the umbrella command surface),
[`cli-evaluate.md`](./cli-evaluate.md) (the evaluation lifecycle that produces the
runs being compared — its [Compare](./cli-evaluate.md#compare) section introduces
the verb and **defers here** for the fuller story), and
[`skills.md`](./skills.md) (the judgment layer, which compare does **not** touch).

## What compare is — and is not

`evaluation compare <a> <b>` answers a **relative** question — *what changed
between these two recorded runs?* — where
[`evaluation report`](./cli-evaluate.md#evaluation-report-id---fail-on-level---json)
answers an **absolute** one — *how does this run rate?* Compare produces a
requirements × runs **diff**, not a fresh verdict.

The load-bearing property is **determinism**. Compare is a function of two
*already-recorded* runs:

```text
compare(a, b) = diff(read(a), read(b))     # no evaluation, no model call, no bash
```

Both operands are stored evaluations — each a set of per-requirement `result`s
with ratings and evidence, produced earlier by the lifecycle in
[`cli-evaluate.md`](./cli-evaluate.md). Compare opens them, lines requirements up,
and reports per-requirement rating **deltas**. It is as deterministic and as cheap
as a text diff, and it is reproducible: the same two runs always diff the same way.

This is the whole inversion applied to comparison. In the superseded model, the
CLI compiled one orchestration prompt and an agent re-evaluated every target. Now
the *evaluation* of each target happens first, separately, through the normal
lifecycle and skill loop; **compare only diffs the records they leave behind.**

## What comparison answers

Two usage modes, the same machinery — the difference is only which two runs are
named and how the diff is framed:

- **A/B across targets sharing one model.** Two instances of "the same" subject
  measured against one `QUALITY.md`: old implementation vs. new, ours vs. a
  competitor, candidate A vs. candidate B. Each is a separate run (one per
  `(model, target)` — see [`cli-evaluate.md`](./cli-evaluate.md#the-run)); compare
  diffs the two records.
- **Base-vs-head regression gating.** The PR question — *did this change make any
  requirement worse than the baseline?* This is `compare <baseline> <head>` with
  the [`--fail-on-regression`](#regression-gating) gate turned on. It is not a
  separate feature: it is the A/B diff with a gate framing.

> **Subject vs. target.** A **target** is a concrete source a run was evaluated
> against (a path or a git ref — see
> [target](./cli-evaluate.md#evaluation-create---model-path---target-path---from-id)).
> A **subject** is the system-under-study a target is an *instance* of. They are
> 1:1 when comparing two distinct subjects (A/B selection), and 1:many in
> regression gating (one subject, base and head are two targets of it). Compare
> always operates on **runs**; "subject" is reserved for the system those runs
> evaluated.

### Validity rests on a shared model

A diff is only meaningful if both runs were evaluated against the **same resolved
`QUALITY.md`** — the same factors, the same requirements, the same `ratings` scale.
The model file is **subject-agnostic** (it declares requirements, never which code
they apply to — see
[`cli-evaluate.md`](./cli-evaluate.md#evaluation-create---model-path---target-path---from-id)),
which is exactly what lets one model rate many targets comparably.

Compare therefore checks model agreement as part of the diff rather than assuming
it: requirements present in one run but not the other are reported as
[newly-present / dropped](#per-requirement-deltas) (model or target drift), not
silently aligned. A regression gate that straddled two different models would be
meaningless, so a model mismatch beyond a tolerable set of added/removed
requirements is surfaced, not scored around.

> Compare does **not** re-impose a shared *evaluation regime* across the runs — it
> cannot, because the judging already happened. Two runs are comparable to the
> extent the skill judged them consistently (same rigor, same evidence standard);
> compare reports the recorded ratings as found. Keeping the judging consistent
> across the two runs is the skill's job ([`skills.md`](./skills.md)); compare's
> job is to diff faithfully and to flag the structural mismatches it *can* detect.

## Operands: what you can compare

Each operand `<a>` / `<b>` resolves to a stored set of recorded results. Because
**everything under `.quality/` is committed and git is the timeline** (see
[`cli-evaluate.md`](./cli-evaluate.md#why-no-finalize-git-is-the-audit-layer)),
compare leans on existing storage rather than a separate "comparison run" concept.
An operand may be:

- the **living run** for the current `(model, target)` — the working state under
  `.quality/evaluations/<slug>/`;
- a **named archive** — `archive/<name>`, a frozen snapshot taken with
  [`evaluation archive --as <name>`](./cli-evaluate.md#evaluation-archive-id---as-name);
- a **git revision** of the same run's files — e.g. `HEAD~1`, a tag, or a branch,
  resolving the committed `evaluations/<slug>/` tree at that rev.

The two canonical shapes both fall out of this:

- **Working vs. archived snapshot** — `compare archive/baseline .` diffs the
  current working run against a deliberately-named baseline. Use when you want a
  stable, labeled reference independent of git history (a release baseline, a
  "known good").
- **Working vs. a git rev of the same run** — `compare HEAD~1 .` or
  `compare main HEAD` diffs the run as committed at one revision against the run
  now. This is the regression-gate shape on a PR: the base branch's committed
  evaluation vs. the head's. It needs **no archive** — the commit history already
  holds the baseline, which is the point of committing everything.

Operands may be labeled for the report with `name=ref` syntax
(`compare base=main head=. ...`); otherwise a label is derived from the ref/path.
Compare resolves a single model (`-f`, default `./QUALITY.md`); comparing whole
*federations* of models is out of scope for v1 (see [Open questions](#open-questions)).

## The requirements × runs matrix

Compare lines the two runs up by requirement — keyed on the requirement's full
dotted path (the same key the result files derive their `<req-id>` from, see
[`cli-evaluate.md`](./cli-evaluate.md#on-disk-layout)) — and emits a two-column
matrix: requirements as rows, the two runs (`a`, `b`) as columns, each cell the
recorded rating (or the absence of one).

```text
requirement (row)                          a            b           delta
-------------------------------------------------------------------------------
secrets.no-committed-secrets               fail         pass        improved
secrets.tokens-short-lived                 pass         pass        unchanged
secrets.vault-loaded                       pass         fail        regressed
secrets.rotation-audited                   —            pass        newly-present
secrets.legacy-key-check                   pass         —           dropped
```

The matrix generalizes to more than two runs in principle, but the v1 verb is
**two operands** (`<a> <b>`): the regression gate is inherently base-vs-head, and a
two-column diff is what reads cleanly as a PR artifact. Wider N-way matrices are an
[open question](#open-questions).

### Per-requirement deltas

For each requirement, compare classifies the change against the model's **ordered**
`ratings` scale (worst → best is defined by the model, the same ordering
[`evaluation report`](./cli-evaluate.md#report-rollup) rolls up against):

| Delta | Meaning |
| --- | --- |
| `improved` | The requirement rates **better** in `b` than in `a` on the scale's ordering. |
| `unchanged` | Same rating in both. |
| `regressed` | The requirement rates **worse** in `b` than in `a`. The gate-relevant case. |
| `newly-present` | Recorded in `b`, absent from `a` (added requirement, or `a` left it `pending`/`skipped`). |
| `dropped` | Recorded in `a`, absent from `b` (removed requirement, or `b` left it `pending`/`skipped`). |

Two refinements keep the diff honest:

- **State, not just rating.** A result that is `pending`, `skipped`, or `errored`
  on one side has **no comparable rating** (see
  [result states](./cli-evaluate.md#result-states)). Compare reports the
  state-vs-rating asymmetry explicitly — e.g. `skipped → pass` is
  `newly-present`, `recorded → pending` is `dropped` — rather than treating a
  missing rating as a low one. It never invents a rating to diff against.
- **Direction comes from the scale, not from string comparison.** `improved` /
  `regressed` are defined purely by the requirement's position on the in-scope
  `ratings` scale. A per-requirement scale override
  ([`cli-evaluate.md`](./cli-evaluate.md#the-bash-path-execution--classification))
  is honored per requirement.

A roll-up summary accompanies the per-requirement deltas: counts of each delta
class, and — for the gate — **whether any `regressed` rows exist**.

## Output: a reviewable PR artifact

Compare's output is built to read like a **diff in a pull request**: a reviewer
scanning it should see, at a glance, what got better and what got worse. The
primary human-facing artifact frames regressions as **diff-like lines**, worst
news first:

```text
Comparison: base (main) → head (working)
Model: ./QUALITY.md   |   2 regressed, 1 improved, 1 unchanged, 1 added, 1 dropped

  REGRESSED
  - secrets.vault-loaded                 pass → fail
  - auth.tokens-short-lived              pass → fail

  IMPROVED
  + secrets.no-committed-secrets         fail → pass

  CHANGED COVERAGE
  ~ secrets.rotation-audited             (added)    — → pass
  ~ secrets.legacy-key-check             (dropped)  pass → —

  UNCHANGED  secrets.tokens-short-lived  pass
```

The `-`/`+`/`~` gutter and "base → head" framing are deliberate: a regression
reads exactly like a removed line in a code diff, so a reviewer's diff-reading
instincts transfer. The same content is available structured under `--json` (a
schema-stable object carrying the matrix, the per-requirement deltas with from/to
ratings and states, and the roll-up counts) for a CI harness to post as a PR
comment or check annotation. Whether the CLI posts that comment itself or leaves
it to the harness is shared with the
[`cli.md` open questions](./cli.md#open-questions).

Because both operands are committed records and the diff is deterministic, the
comparison itself is reproducible from the two refs — there is no bundle to seal
and nothing volatile to segregate. Compare reads; it does not write a run.

## Regression gating

Regression gating is `compare <baseline> <head>` with one flag:

- **`--fail-on-regression`** — exit **non-zero** if **any** requirement is
  `regressed` (rates worse in `b`/head than in `a`/baseline). This is the
  comparison analogue of
  [`evaluation report --fail-on`](./cli-evaluate.md#evaluation-report-id---fail-on-level---json):
  a PR gate that blocks on **new** shortfalls without failing on pre-existing debt
  (a requirement already `fail` in the baseline and still `fail` is `unchanged`,
  not a trip). **Off by default** — without it, compare is report-only and a diff
  is never a gate. Opt-in for the same reason `--fail-on` is: the underlying
  ratings were produced by a non-reproducible judging process, so gating is a
  team's explicit choice.

By default `--fail-on-regression` trips on **any** `regressed` row. Whether
`newly-present`/`dropped` (coverage drift) can also be made to trip, and whether
the gate can be scoped to a factor subtree, are
[open questions](#open-questions). The direction is fixed by operand order:
`<a>` is the baseline, `<b>` is the candidate, so put base first, head second.

## Exit codes

Compare reuses the shared three-code convention verbatim (see
[`cli.md`](./cli.md#exit-codes)):

| Code | Meaning | In compare |
| --- | --- | --- |
| `0` | Success — ran, and the gate (if any) passed. | The default, report-only outcome: the diff was produced. Exits `0` even when runs diverge sharply, *unless* `--fail-on-regression` is set and a regression was found. A clean diff and a divergent-but-ungated diff both exit `0`. |
| `1` | **Gate verdict failure** — ran fine, the bar was not met. | `--fail-on-regression` was passed and at least one requirement `regressed`. The comparison ran successfully; a regression was found. The *only* `1` in compare. |
| `2` | **Tool failure** — the command broke. | An operand that does not resolve (unknown ref/archive, absent run), a `QUALITY.md` that does not parse, `--fail-on-regression` with operands evaluated against **incompatible models**, a malformed stored result, or an internal error. Never means "a run is worse." |

The load-bearing distinction is the same one the lifecycle draws: **a found
regression is `1`, not `2`.** Finding that head is worse than base is a successful
comparison reporting bad news (gate trip); only a *broken invocation* — an
unresolvable operand, a model mismatch the gate can't reason over — is `2`. A
report-only run that surfaces regressions still exits `0`.

## Relationship to `evaluation compare` in the lifecycle doc

[`cli-evaluate.md`](./cli-evaluate.md#compare) introduces `evaluation compare` as
the minimal two-run diff the `evaluation` resource exposes and defers here for the
fuller story. The two are the **same verb**, not two: this doc is the authoritative
specification of its semantics (deltas, the matrix, gating, output, exit codes);
that section specifies only the minimal diff and points here. They must stay
consistent — in particular, both describe a **deterministic diff of recorded
results that never re-judges or calls a model**, and both allow either operand to
be the living run, a named archive, or a git revision.

## Open questions

- **N-way compare.** Whether the verb accepts more than two operands (a true
  requirements × N-runs matrix for selection/benchmarking), or stays strictly
  two-operand (base-vs-head) for v1 with N-way left to the harness.
- **Coverage-drift gating.** Whether `--fail-on-regression` can be configured to
  also trip on `dropped` (a requirement that lost its rating) or
  `newly-present`-as-fail, vs. trip on rating regressions only.
- **Factor-scoped compare/gate.** Whether compare accepts a trailing `factor` path
  to diff/gate a single factor subtree, mirroring the model's factor scoping, or
  is always whole-model.
- **Model-mismatch tolerance.** How much model drift between the two runs is
  tolerable before the diff is `2` rather than reported as added/dropped
  requirements — the exact boundary between "comparable with coverage drift" and
  "incomparable."
- **Federated operands.** Comparing two runs that are each a whole *federation* of
  models (see [`cli-federation.md`](./cli-federation.md)) — a tree-shaped diff —
  is out of scope for v1; compare resolves a single model for now.
- **Three-way / merge-base compare.** Whether a `base...head` (merge-base) form is
  useful for PR gating the way `git diff` three-dot range is, vs. plain two-operand.
