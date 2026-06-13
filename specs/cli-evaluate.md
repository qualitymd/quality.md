# CLI: `evaluate` / `evaluate-model`

> Detail doc for the deep, agentic **semantic tier**. See [`cli.md`](./cli.md)
> for the full command surface and shared conventions, and
> [`cli-lint.md`](./cli-lint.md) for the fast structural tier.

| Command | Purpose | Cost / determinism | Output |
| --- | --- | --- | --- |
| `qualitymd evaluate [factor]` | Deep ISO 25040 evaluation of the **subject** against the model's requirements. | Expensive, non-reproducible. | Evaluation bundle on disk. |
| `qualitymd evaluate-model [factor]` | Deep evaluation of the **quality model itself** — are these the right requirements, well-specified, and complete against the real subject? | Expensive, non-reproducible. | Evaluation bundle on disk. |

`evaluate` and `evaluate-model` are **one engine with the target
swapped** (see below) — two verbs over the same pipeline, differing only in which
`(target, criteria, context)` triple they run.

## Factor

The optional positional `factor` names the factor or sub-factor subtree to
evaluate, by dotted path — `security`, `maintainability.testability`. Omitted =
the whole model. It scopes *which requirements* are evaluated, not which files.
(As an axis it is the **requirement selector**; the argument is named for its
value — a factor path.)

## Target source & scope

An evaluation is parameterized on **three independent axes**, not one:

| Axis | What it varies | Controlled by |
| --- | --- | --- |
| **What's evaluated** | the subject vs. the quality model | `evaluate` / `evaluate-model` |
| **Requirement selector** | *which requirements* (vertical slice of the model) | the positional `factor` |
| **Target source & scope** | *which instance of the subject, and which files within it* (horizontal slice) | the flags below |

The selector and the target scope are orthogonal: the selector slices the
*model*, the target scope slices the *code*. Both can be applied at once
(evaluate `security` over the current diff).

### Target source

Where the code under evaluation comes from. The **criteria source stays fixed**
(the resolved `QUALITY.md`); only the target moves.

| Flag | Target | Use case |
| --- | --- | --- |
| *(none)* | **Working tree**, as-is — including uncommitted changes. The default. | Point-in-time audit of what's on disk. |
| `--ref <git-ref>` | A **clean checkout** of a committed ref (branch, tag, SHA). Ignores the dirty working tree. | "Evaluate `main` as committed," reproducible-as-possible runs. |

These resolve what "the codebase" means — a plain implicit `.` conflates
"the working tree" with "the committed state." They are kept distinct.

### Target scope

Which files within the target the finders attend to. Independent of source.

| Flag | Scope | Use case |
| --- | --- | --- |
| *(none)* | Whole tree. | Default. |
| `--diff [<base>..<head>]` / `--since <ref>` | **Changed-files scope** — findings are prioritized and attributed to the changed region. | "What's the quality impact of this change?" — local edit, PR. |

`--diff` does **not** turn whole-system requirements ("no secrets committed,"
"coverage ≥ X") into diff-local ones — those are still evaluated against the
whole target. The report states which requirements were evaluated whole vs.
scoped to the change, so a clean diff is never mistaken for a clean system.

> **"Did this change make quality *worse*?" is not a scope — it's a
> comparison.** Regression ("don't fail the PR for pre-existing debt, only for
> what this change introduces") is a base-vs-head run of the **same model over
> two targets**, which is [`compare`](./cli-compare.md), not a `--diff` flag.
> `--diff` answers "what did I touch"; `compare base head` answers "what got
> worse." See [`cli-compare.md`](./cli-compare.md).

## Conceptual model

### ISO 25040 process, one-shot

The engine follows the ISO/IEC 25040:2024 quality-evaluation process reference
model — **Define → Design → Plan → Execute → Conclude** — adapted for an agentic
coding context: no human "staffing & scheduling" (5.4.3.2) and no
"gain consensus" / human results review (5.4.3.4, 5.6.3.1). The run is one-shot.

Removing the human review steps has a consequence the design must honor: **the
adversarial machinery is the load-bearing replacement for human review.** It
carries the rigor ISO normally gets from people.

### One engine, two targets

`evaluate` and `evaluate-model` are the same pipeline with the
`(target, criteria)` pair swapped. ISO names both modes:

| | `evaluate` | `evaluate-model` |
| --- | --- | --- |
| ISO type | **T3** — checking requirements satisfaction | **T2** — qualification to a quality standard (diagnostic model) |
| Target entity | the subject | the `QUALITY.md` + its referenced prompts/standards |
| Criteria | the model's requirements | a diagnostic model of *what a good quality requirement is* |
| Context | the requirements | **the subject** |
| Findings | ways the subject falls short | requirement **defects** + coverage **gaps** |

The **T1–T4** labels are ISO's "four types of quality evaluation," defined in
ISO/IEC 25040:2024 §5.2.3.1 ("Establish the purpose"): T1 — suitability to a
specific context of use; T2 — qualification to a quality standard; T3 — checking
requirements satisfaction; T4 — suitability to the market. Each type fixes the
*source* of evaluation criteria — Table 2 maps T3 to a requirements
specification and T2 to a diagnostic model, which is precisely the criteria swap
in the table above.

ISO 25040 (§5.2.3.3 NOTE 2, and Table 3) explicitly notes that
requirements-conformity evaluation alone cannot *score* quality or tell you
whether your requirements are any *good* — for that you need a diagnostic model.
That is the reason `evaluate-model` exists as a first-class sibling.

The CLI **ships this diagnostic model built-in** as `QUALITY-META-MODEL.md` — a normal
QUALITY.md-schema file whose subject is another `QUALITY.md`. The bundled meta
model is a framework of *what a good quality model is* — its factors are
product-quality attributes of the model-as-artifact (the **ISO/IEC 25000
(SQuaRE)** vocabulary, ISO/IEC 25010 style), while its requirement-quality
criteria follow **ISO/IEC/IEEE 29148** (§5.2.5–5.2.6, the characteristics of a
well-formed requirement and requirement set) — plus the QUALITY.md spec. It is
overridable (`replace`) or extendable (`extend`) per project via the
`requirementsDiagnosticModel` config block (see
[`cli.md`](./cli.md#configuration--quality)).

`QUALITY-META-MODEL.md` is an implementation asset, not a second public root-file
convention. Users normally author `QUALITY.md`; the CLI loads its bundled
`QUALITY-META-MODEL.md` when running `evaluate-model`.

`evaluate-model` yields two finding types:

- **Defects** — an existing requirement is vague, unmeasurable, mis-targeted, or
  lacks gradation.
- **Gaps** — a requirement that *should* exist given the real codebase (e.g.
  "there is an auth module but nothing covers session fixation"), grounded in
  code evidence rather than generic best practice.

It is also QUALITY.md dogfooding its own format on itself.

### `QUALITY-META-MODEL.md` target binding

Although `QUALITY-META-MODEL.md` is parsed with the same schema as any other quality
model, `evaluate-model` binds its target specially:

```text
target   = the resolved project QUALITY.md + referenced prompts/standards
criteria = the bundled or configured QUALITY-META-MODEL.md requirements
context  = the subject code selected by --ref / --diff / --since
```

The bundled meta model should therefore avoid literal `target` paths that point
beside the bundled file. Its requirements describe properties of the project
model, while the command invocation supplies the concrete project model and code
context being evaluated. This is what lets the same engine evaluate a subject
against `QUALITY.md` and evaluate `QUALITY.md` against `QUALITY-META-MODEL.md`.

### T2 is also a requirement flavor, not just a command mode

The `evaluate` ⇒ T3 / `evaluate-model` ⇒ T2 mapping above is each
*command's* dominant evaluation type — not an exhaustive one. In ISO, the type
is a property of a *requirement* (it determines that requirement's criteria
source, per Table 2), so a single `QUALITY.md` can legitimately mix flavors:

- **T3-flavored requirement** — a concrete, agreed criterion the code must
  *satisfy*, with the criteria stated inline (e.g. "no secrets committed to the
  repository"). The requirement *is* its own criteria source.
- **T2-flavored requirement** — a requirement that says *qualify against this
  standard* (e.g. "conform to OWASP ASVS L2", "meet our org's API-design
  checklist") rather than spelling out every criterion. Its criteria source is
  the referenced standard / diagnostic model.

So `evaluate` over the codebase is *predominantly* T3 but may carry T2
requirements; ISO explicitly anticipates the combination (§5.2.3.1 EXAMPLE: T2
and T3 are combined in one evaluation when agreed requirements alone don't yield
sufficient criteria). When a requirement is T2-flavored, the engine resolves its
criteria from the named standard using the **same diagnostic-model machinery**
`evaluate-model` uses on the whole model — applied here to a single
requirement. This is the symmetry worth keeping: the diagnostic model is a
first-class input to the engine, consumed both as the *subject's criteria*
(`evaluate-model`) and as a *requirement's criteria source* (a T2
requirement inside `evaluate`).

> **Relationship to `lint`.** `lint` (see [`cli-lint.md`](./cli-lint.md)) and
> `evaluate-model` both scrutinize the `QUALITY.md` file, but at different
> depths. `lint` checks *form* deterministically — does it parse and conform to
> the spec. `evaluate-model` checks *substance* by judgment — are these
> the right requirements, well-stated, complete against the code. A file should
> pass `lint` before `evaluate-model` is worth running.

## The evaluation engine (shared)

Both deep commands run the same five steps and emit the same four-artifact
bundle. Only the `(target, criteria, context)` triple differs.

| Step | ISO | Artifact | Contents |
| --- | --- | --- | --- |
| Define | 5.2 | *(report header)* | purpose, mode, factor, criteria source, rigor — all fixed by the invocation |
| Design | 5.3 | `design.md` | decomposition of the selected subtree into requirements → *information needs*; the rating module per requirement (how it's assessed, evidence sources, how findings map to ratings); analysis/rollup method; planned outputs |
| Plan | 5.4 | `plan.md` | the concrete task graph: partitions, finder lenses, saturation strategy + stop criteria, verification strategy. No staffing/consensus. |
| Execute | 5.5 | `results.json` | per-requirement rating + findings + evidence/locations + verification verdicts + saturation stats |
| Conclude | 5.6 | `report.md` | overall rating, per-factor rollup, strengths/weaknesses, unsatisfied requirements, prioritized recommendations, limitations |

### Assessment paths

- **`bash` requirements** bypass the agentic harness entirely. The engine runs the
  command once and **classifies the result against the `ratings` scale's
  `bashCondition`s**, exactly as the format spec defines it
  (`../SPECIFICATION.md#computational-rating`): the levels are tested best
  to worst over a single `result` (`success`, `exit`, `stdout`, `stderr`) and the
  first matching level wins, with the worst level as the default fallback. So the
  verdict is *not* simply "exit zero" — a scale may band a numeric value the
  command prints on stdout (`double(result.stdout.trim()) >= 90`), and a
  requirement may carry a [per-requirement override](../SPECIFICATION.md#per-requirement-rating-overrides)
  for its own bands. Only the default `pass`/`fail` scale reduces to "best on a
  zero exit." Either way the path is **deterministic** — no model judgment, no
  adversarial loop — and a `bashCondition` that fails to evaluate (e.g. `.json()`
  on non-JSON output) is a model configuration error, surfaced as a tool failure
  (exit `2`), not silently scored.
- **`prompt` requirements** go through the adversarial harness below, judged
  against the scale's `promptCondition`s.

### Saturation: adversarial loop-until-dry

Findings have two independent failure axes, handled by different mechanisms:

- **Recall (saturation)** — *did we find everything?* Driven by: partitioning the
  search space, an independent ensemble of finder "lenses" per partition, a
  loop that runs until findings dry up, and a completeness critic.
- **Precision (verification)** — *are the findings real?* Driven by independent
  refuters that vote on each candidate. Contested findings are **surfaced with a
  confidence level, not silently dropped** — for a quality gate, a missed real
  problem is worse than a flagged borderline one.

Per requirement (or per partition), Execute runs:

```text
seen = {}            # semantically deduped findings
confirmed = []
dryRounds = 0
while dryRounds < K:
    # ensemble of independent finders, each a different lens (by file/dir,
    # by category, by sub-requirement), over target ∩ partition
    fresh = finders() - seen
    if fresh is empty:
        dryRounds += 1
        continue
    dryRounds = 0
    seen += fresh
    for f in fresh:
        # independent refuters each try to DISPROVE f; majority vote
        f.confidence = vote(refuters(f))   # confirmed | contested | refuted
    confirmed += [f for f in fresh if f.confidence != refuted]
# completeness critic: name any category/area/modality never examined;
# if any, enqueue as new partitions and resume the loop
```

Each requirement's **rating** is then derived from its confirmed findings
against the model's `ratings` scale: the best level when there are no confirmed
shortfalls, lower levels as severity/count grow. Ratings roll up to factor and
overall via the analysis method recorded in `design.md`.

### Rigor ladder

`--rigor` scales the harness; it maps to ISO's notion of evaluation rigor
(coverage of information needs, objectivity, transparency).

| Level | Behavior |
| --- | --- |
| `low` | Single finder pass per requirement; no adversarial loop. Lightest CI-gate setting. |
| `medium` | Bounded finder ensemble + one refutation pass + one completeness-critic round. |
| `high` *(default)* | Full loop-until-dry, `K=2` dry rounds, multi-vote refutation. |
| `max` | Loop-until-dry with higher `K`, more finders/voters, finer partitioning. |

Rigor scales four knobs: finders per round, dry-round threshold `K`, refutation
voters, and partition granularity. Default is configurable in `./.quality/`
(see [`cli.md`](./cli.md#configuration--quality)).

## Output: the evaluation bundle

```text
./.quality/evaluations/<YYYYMMDD-HHMMSS>-<mode>-<factor>/
  design.md      # Design step — how the evaluation will be conducted
  plan.md        # Plan step  — the concrete task graph that was run
  results.json   # Execute step — structured ratings, findings, verdicts
  report.md      # Conclude step — the primary human-facing artifact
```

`report.md` is the usable output; the other three exist for inspection,
auditing, and improving the process. The mode in the directory name
(`subject` / `model`) keeps the two evaluation modes from colliding. `<factor>`
is the dotted factor path that was selected, or `all` when the whole model is
evaluated (the positional `factor` omitted).

### `report.md`

Follows ISO 25040 Table 3 outputs:

- **Definition** — what was evaluated (mode, factor, criteria source,
  rigor, timestamp).
- **Verdict** — overall pass/fail + rating.
- **Per-factor rollup** — rating per factor/sub-factor.
- **Strengths & weaknesses.**
- **Unsatisfied requirements** — with evidence and `path:line` locations.
- **Recommendations** — prioritized corrective matters. For
  `evaluate-model`, these include concrete, ready-to-apply YAML edits and
  additions to `QUALITY.md`.
- **Limitations** — rigor used, saturation achieved, contested findings, and the
  non-reproducibility caveat.

### `results.json` (illustrative shape)

```json
{
  "definition": {
    "mode": "subject",
    "factor": "security",
    "rigor": "high",
    "criteriaSource": "./QUALITY.md",
    "target": "."
  },
  "ratings": { "pass": {}, "fail": {} },
  "factors": {
    "security": {
      "requirements": {
        "no secrets committed to the repository": {
          "assessment": "prompt",
          "rating": "fail",
          "findings": [
            {
              "id": "sec-001",
              "summary": "AWS key hard-coded in test fixture",
              "locations": ["test/fixtures/config.json:12"],
              "severity": "high",
              "confidence": "confirmed",
              "verification": { "voters": 3, "refuted": 0 }
            }
          ],
          "saturation": { "rounds": 4, "dryRounds": 2, "findersPerRound": 3 }
        }
      }
    }
  },
  "summary": { "overall": "fail", "byFactor": { "security": "fail" } }
}
```

For `evaluate-model`, each finding additionally carries a `kind`
(`defect` | `gap`) and an optional `proposedChange` holding the YAML
snippet/diff to apply to `QUALITY.md`.

### Finding severity & next-action priority

A deep-command finding carries a `severity` on a three-level scale —
**`high` / `medium` / `low`** — reflecting how badly the shortfall (or, for
`evaluate-model`, the defect/gap) undermines the requirement. This is the
deep-command analogue of `lint`'s `error` / `warning` / `info`.

Each finding emits a corresponding entry in the run's `nextActions` (see
[`cli.md`](./cli.md#structured-next-action-suggestions)), and the finding's
severity sets that action's `priority` on the same ladder:

| Finding `severity` | Next-action `priority` |
| --- | --- |
| `high` | `required` |
| `medium` | `recommended` |
| `low` | `optional` |

Two refinements specific to the deep commands:

- **Confidence caps priority.** A finding whose `confidence` is `contested` (the
  refuters did not agree it is real — see
  [Saturation](#saturation-adversarial-loop-until-dry)) is capped at
  `recommended`, never `required`: an agent should not be told a *disputed*
  finding blocks progress. `confirmed` findings map by severity as above.
- **The action depends on the verb.** For `evaluate`, the next-action is "address
  this requirement," carrying the finding's `path:line` evidence. For
  `evaluate-model`, it is "apply this `proposedChange` to `QUALITY.md`,"
  carrying the YAML snippet — then re-`lint` and re-evaluate.

As in `cli.md`, these suggestions are advisory: severity drives `priority`, but
the *gate* verdict is carried by the exit code (`--fail-on`), never by
`nextActions`.

## Flags, exit codes

Flags specific to the deep commands (shared flags are in
[`cli.md`](./cli.md#shared-conventions)):

- `--rigor low|medium|high|max` — coverage/cost dial (default `high`).
- `--ref <git-ref>` — evaluate a clean checkout of a committed ref instead of the
  working tree (see [Target source](#target-source)).
- `--diff [<base>..<head>]` / `--since <ref>` — scope findings to the changed
  region (see [Target scope](#target-scope)).
- `--fail-on <rating>` — exit non-zero when the overall rating is at or below
  `<rating>`. **Off by default** (report-only); opt-in for a CI gate. See exit
  codes below.
- `--name <name>` — override the derived bundle name.

Exit codes follow the shared three-code convention (see
[`cli.md`](./cli.md#machine-readable-result-contract)):

- **`0`** — the run completed. **This is the default**, report-only outcome
  *regardless of rating*: without `--fail-on`, even an all-`fail` evaluation exits
  `0`, because a non-deterministic deep audit is a poor *implicit* hard gate.
  `lint` remains the deterministic structural gate (see
  [`cli-lint.md`](./cli-lint.md)).
- **`1`** — **gate verdict failure:** `--fail-on <rating>` was passed and the
  overall rating landed at or below `<rating>`. The run succeeded; the subject (or,
  for `evaluate-model`, the model) just did not clear the bar. This is the opt-in
  semantic gate — the automated-PR use case (post a deep evaluation on a pull
  request and block on a bad rating) is its concrete motivation. It is opt-in
  precisely because the run is non-reproducible: a team chooses to accept that
  trade-off, the tool never imposes it. Pair with a machine-readable output sink
  (see [`cli.md`](./cli.md#shared-conventions)) when wiring into CI.
- **`2`** — **tool failure:** the run could not produce a trustworthy verdict — a
  bad flag, an unresolvable `--ref`/target, a `QUALITY.md` that does not parse
  (run `lint` first), or an internal error. A `2` never means "quality is bad,"
  so `--fail-on` and a broken run stay distinguishable.

> **`--fail-on` and `evaluate-model`.** For `evaluate`, the gated value is the
> subject's overall rating. For `evaluate-model`, it is the model's overall rating
> against the meta-model — defects and gaps roll up to a rating on the same scale,
> so `--fail-on` reads identically (e.g. "block if the model rates `fail`").

## Open questions

- **Semantic CI gate.** *Resolved:* default stays report-only; `--fail-on
  <rating>` is the opt-in semantic gate, motivated by the automated-PR use case
  (see [Flags, exit codes](#flags-exit-codes)). Still open: the **output sink** a
  CI gate writes to — stdout JSON and/or a postable PR-comment/check-annotation
  summary, vs. only the on-disk bundle. The bundle is the wrong primary output in
  CI; the sink format is a shared concern (see
  [`cli.md`](./cli.md#open-questions)).
- **Define step surfacing.** Folded into the `report.md`/`design.md` header
  (current plan) vs. a fifth `define.md` artifact.
- **Factor-selector grammar.** *Partly resolved:* file/target scoping is now its own axis
  (`--diff`/`--since`/`--ref`, see [Target source &
  scope](#target-source--scope)), separate from the requirement selector. Still
  open: whether the selector also gains glob/target scoping for the
  *non-diff* case (e.g. evaluate `security` only over `./src/api`).
- **Expressing a T2 requirement in the format.** A requirement that says
  "qualify against standard X" instead of stating its criteria inline is a
  distinct flavor (see [T2 is also a requirement
  flavor](#t2-is-also-a-requirement-flavor-not-just-a-command-mode)). How a
  `QUALITY.md` requirement references an external standard / diagnostic model as
  its criteria source — and how that interacts with the `prompt`/`bash`
  assessment split — is a spec question, not yet settled.
- **Cross-run diffing.** Should `report.md` reference the previous bundle for the
  same `(mode, factor)` to show what changed? Out of scope for v1. (See also
  the `diff` verb candidate in [`cli.md`](./cli.md#open-questions).)
