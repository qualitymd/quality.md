# CLI specification

Functional specification for the `qualitymd` CLI surface.

> **Status:** high-level *functional* spec. It fixes the command surface, the
> two evaluation tiers, and the output contracts. Exact agent topology, JSON
> field names, lint-rule identifiers, and rigor parameters are illustrative and
> expected to be tuned during implementation. Items still genuinely open are
> collected under [Open questions](#open-questions).

This is the umbrella document. Per-command detail lives in:

- [`cli-init.md`](./cli-init.md) — scaffolding a first `QUALITY.md` from a static
  template (deterministic, offline).
- [`cli-lint.md`](./cli-lint.md) — the fast, deterministic structural validator
  (frontmatter and Markdown body rules).
- [`cli-evaluate.md`](./cli-evaluate.md) — the deep agentic evaluation engine
  (`evaluate` and `evaluate-model`).
- [`cli-compare.md`](./cli-compare.md) — comparison mode: the same engine run
  over multiple targets against one shared model (`compare`).

## Command surface

| Command | Purpose | Cost / determinism | Output |
| --- | --- | --- | --- |
| `qualitymd init` | Scaffold a first `QUALITY.md` from a static starter template, for the author to fill in. | Cheap, fully deterministic, offline. | `./QUALITY.md` on disk. |
| `qualitymd lint [file]` | Structural validation of the `QUALITY.md` file itself against the format spec — does it parse, conform to the schema, and resolve its references? | Cheap, fully deterministic. | JSON findings + summary on stdout. |
| `qualitymd evaluate [factor]` | Deep ISO 25040 evaluation of the **subject** against the model's requirements. | Expensive, non-reproducible. | Evaluation bundle on disk. |
| `qualitymd evaluate-model [factor]` | Deep evaluation of the **quality model itself** — are these the right requirements, well-specified, and complete against the real subject? | Expensive, non-reproducible. | Evaluation bundle on disk. |
| `qualitymd compare <target> <target> [...] [factor]` | Evaluate **multiple targets against one shared model** and rank/diff them per requirement. | Expensive, non-reproducible (× N targets). | Comparison bundle on disk. |

### Two tiers

The surface splits cleanly into two tiers that answer different questions:

- **Structural tier — `lint`.** *Is this a well-formed `QUALITY.md`?* A
  deterministic check that the file conforms to the format spec: frontmatter
  parses, every factor has `requirements`, `factors`, or both, every
  requirement declares exactly one assessment, `prompt`/`target` paths resolve,
  the `ratings` scale is well-shaped. Cheap enough to run on every save and in a
  pre-commit hook. Modeled on Google's `design.md lint` (see
  [`cli-lint.md`](./cli-lint.md)).
- **Semantic tier — `evaluate` / `evaluate-model` / `compare`.** *Is the
  system actually good, are these the right requirements, and how do candidates
  rank?* An agentic, judgment-based audit that reads the codebase. Expensive and
  non-reproducible. `evaluate` gives an **absolute** verdict on one target (see
  [`cli-evaluate.md`](./cli-evaluate.md)); `compare` gives a **relative** verdict
  across multiple targets sharing one model (see
  [`cli-compare.md`](./cli-compare.md)).

The tiers are complementary, not redundant. `lint` proves the *file* is correct
without an opinion on whether the requirements are any good; `evaluate-model`
forms exactly that opinion but cannot run meaningfully on a file that does not
parse. Lint first, evaluate second.

Both tiers assume a model already exists. `init` (see
[`cli-init.md`](./cli-init.md)) is step zero — a deterministic template scaffold
that produces the first `QUALITY.md` for the author to fill in. Tailoring that
scaffold to a real codebase is not a special `init` mode; it falls out of
`evaluate-model`, whose structured suggestions propose the requirements
the code warrants (see [Agent-friendly CI patterns](#agent-friendly-ci-patterns)).
The intended order is init → lint → evaluate.

This split also settles the previously-deferred `check` verb: the CI gate is
`lint` (deterministic, fast). The deep `evaluate` commands are report-only in v1
(see [`cli-evaluate.md`](./cli-evaluate.md#flags-exit-codes)); a semantic gate is
a later candidate. There is no separate `check`.

## Shared conventions

These apply across all commands; per-command flags are documented in the detail
docs.

- **Model file resolution.** `-f, --file <path>` selects the `QUALITY.md` model
  file; default `./QUALITY.md`. `lint` also accepts the file as a positional
  argument and `-` for stdin.
- **Paths are model-relative.** `target` and `prompt` paths inside a
  `QUALITY.md` resolve relative to that file's directory, not the working
  directory — consistent across `lint` (existence checks) and `evaluate`
  (assessment).
- **Output format.** `lint` emits JSON by default (agent-consumable). The deep
  commands write a Markdown-first bundle to disk; see
  [`cli-evaluate.md`](./cli-evaluate.md#output-the-evaluation-bundle) and
  [`cli-compare.md`](./cli-compare.md#output-the-comparison-bundle).
- **Agent-friendly output.** Every command speaks a stable, machine-readable
  contract and tells the caller what to do next, so an agent can drive the CLI in
  CI without hard-coding the workflow. See
  [Agent-friendly CI patterns](#agent-friendly-ci-patterns).
- **Interactivity and the non-interactive path.** A command runs **non-interactive**
  — no prompts, never blocks on input — if **any** of these hold:
  - stdin/stdout is **not a TTY** (piped, redirected, or running in CI), detected
    automatically; **or**
  - **`--non-interactive`** is passed — a global flag, valid on every command, that
    forces this path explicitly even on a terminal; **or**
  - **`--json`** is passed — requesting machine-readable output signals automation,
    so the command must not block on a prompt.

  Otherwise (an interactive terminal with none of the above) a command **may**
  prompt. Today only [`init`](./cli-init.md) does; the rest are non-interactive
  unconditionally. Any prompt is written to **stderr**, so it never contaminates
  stdout (JSON or otherwise). This one rule lets CI and agents guarantee "never
  block on input" without knowing which subcommand might prompt.

## Agent-friendly CI patterns

The CLI is designed to be driven by a **coding agent in CI**, not only read by a
human in a terminal. A QUALITY.md gate typically runs inside an automated PR
workflow where the agent's job is not just to learn the verdict but to *act on
it* — fix the file, fill in a requirement, address a finding, and re-run. These
patterns are shared across commands so the agent can chain them without
hard-coding the workflow.

### Structured next-action suggestions

Every command's output carries a `nextActions` array: an ordered, machine-readable
list of suggested follow-up steps. Each entry is a runnable command plus the
*reason* it is suggested and a priority, so an agent can decide what to do next
from the output alone — and a human reads the same suggestions as guidance.

```json
"nextActions": [
  {
    "command": "qualitymd lint",
    "reason": "Validate the scaffold before evaluating.",
    "priority": "recommended"
  }
]
```

`priority` is a three-level ordinal — an agent acts on `required` before
`recommended` before `optional`:

| Priority | Meaning |
| --- | --- |
| `required` | Needed before the output or the workflow can be trusted or continued — the top remediation. Still *advisory* (the exit code, not this field, carries the gate verdict); it means "do this first to make progress," not "you are failed." |
| `recommended` | The normal happy-path next step, or a finding worth addressing. The default. |
| `optional` | Exploratory or nice-to-have. |

When a next-action remediates a specific finding, its `priority` tracks that
finding's **severity** (`error → required`, `warning → recommended`,
`info → optional`); when it is plain workflow progression (init → lint →
evaluate), it is `recommended`. This keeps next-action priority and finding
severity as one ladder rather than two unrelated vocabularies.

The suggestions are what make the **authoring/evaluation loop self-describing**
rather than tribal knowledge — each command points at the next, and a failure
points at its own remedy:

- `init` → run `lint` once the placeholders are filled in.
- `lint` **with errors** → fix the specific keys the findings name (the finding's
  `path` + `message` is the actionable unit); re-run `lint`.
- `lint` **clean** → run `evaluate-model` to pressure-test whether these
  are the right requirements.
- `evaluate-model` → apply the `proposedChange` snippets on its findings to
  `QUALITY.md` (this is also how a thin scaffold gets fleshed out into a
  codebase-tailored model — see [`cli-init.md`](./cli-init.md#purpose)), then
  `lint`, then `evaluate`.
- `evaluate` **with unsatisfied requirements** → address the highest-priority
  finding, with its `path:line` evidence as the starting point.

Suggestions are *advisory*: they never change exit status, and an agent is free
to ignore them. They encode the recommended path, not a mandate.

### Machine-readable result contract

- **Stable JSON on stdout.** `--json` is available on every command and
  emits a schema-stable, versioned object (a `schemaVersion` field) so an agent
  can parse results without screen-scraping. `lint` already defaults to JSON; the
  deep commands write their Markdown-first bundle to disk *and* can emit a JSON
  result (findings + `nextActions` + a pointer to the bundle) to stdout for the
  agent. The on-disk bundle is for human audit; stdout JSON is the agent's sink.
- **Exit codes separate verdict from failure.** A non-zero exit from a *gate*
  (findings that fail the gate) is distinct from a non-zero exit caused by the
  tool itself failing to run (bad flags, unreadable file, internal error), so an
  agent can tell "the quality is bad" from "the command broke" and react
  differently. The two cases never share a code. Every command uses one shared
  three-code convention:

  | Code | Meaning |
  | --- | --- |
  | `0` | Success — the command ran and the gate (if any) passed. Report-only runs that complete also exit `0`. |
  | `1` | **Gate verdict failure** — the command ran fine, but the quality bar was not met: `lint` found an `error`, `evaluate --fail-on` tripped, `compare --fail-on-regression` tripped. "The quality is bad." |
  | `2` | **Tool failure** — the command could not produce a trustworthy verdict: bad flags, unreadable/absent file, parse-of-CLI-input error, internal error. "The command broke." |

  An agent keys off this split: `1` means *act on the findings*; `2` means *fix
  the invocation*. Detail docs restate the codes a command actually emits, but
  never reassign these meanings.
- **Non-blocking and idempotent in automation.** Whenever a command is on the
  non-interactive path (see [Shared conventions](#shared-conventions) — not a TTY,
  `--non-interactive`, or `--json`), it never prompts, so every scripted or agent
  invocation is re-runnable. Deterministic commands (`init`, `lint`) produce
  identical output on identical input; the deep commands are non-reproducible by
  nature but still never block on input. Interactive prompts (only `init`, only on
  a terminal) are a convenience layer that pre-fills the scaffold; the
  non-interactive path writes the same file with placeholders instead.

### Postable summaries (open)

Beyond stdout JSON, an automated-PR run often wants a summary posted as a PR
comment or a check annotation. Whether the CLI emits that directly (a dedicated
flag) or leaves it to the CI harness consuming the JSON is still open (see
[Open questions](#open-questions)).

## Configuration — `./.quality/`

`./.quality/` is the project's quality home: configuration plus the evaluation
bundles under `evaluations/`.

```yaml
# ./.quality/config.yaml  (illustrative)
defaults:
  rigor: high
requirementsDiagnosticModel:
  mode: extend          # extend | replace the built-in QUALITY-META-MODEL.md
  file: ./.quality/QUALITY-META-MODEL.md
```

The diagnostic model for `evaluate-model` ships built-in as `QUALITY-META-MODEL.md`, a
normal QUALITY.md-schema file bundled with the CLI. This config can replace it
or extend it with a project-local meta model. In `extend` mode, v1 appends the
configured model's factors/requirements to the built-in model; it does not try
to deep-merge same-named factors. The config may also carry defaults (rigor,
output dir, ignore globs, model selection) as those settle. `lint` is governed
by the spec, not this config, so it needs nothing here in v1.

## Open questions

- **Dedicated `check` verb.** Resolved: dropped in favor of the two-tier split
  (`lint` for structure, `evaluate` for substance). Reopen only if a single
  convenience verb proves worth the redundancy.
- **Postable CI summaries.** The machine-readable sink is settled — stdout JSON
  with `nextActions`, see [Agent-friendly CI patterns](#agent-friendly-ci-patterns).
  What remains open is the *human*-facing PR artifact: whether the CLI emits a PR
  comment / check annotation directly (a dedicated flag) or
  leaves that to the CI harness. Relevant to `evaluate` (with
  [`--fail-on`](./cli-evaluate.md#flags-exit-codes)) and `compare` (with
  [`--fail-on-regression`](./cli-compare.md#flags-exit-codes)).
- **`spec` and `diff` verbs (design.md prior art).** `design.md` ships `spec`
  (print the format spec for injection into agent prompts) and `diff` (compare
  two files). `spec` is a plausible addition for agent ergonomics. The `diff`
  niche — comparing two *targets* against a model, and base-vs-head regression —
  is now filled by [`compare`](./cli-compare.md); a separate file-level `diff`
  for tracking model drift across commits remains a candidate. Out of scope for
  v1.
- **Factor-selector grammar.** Dotted factor path only, or also glob/target
  scoping (e.g. evaluate `security` only over `./src/api`)? See
  [`cli-evaluate.md`](./cli-evaluate.md#factor).
- **Bundle retention.** ISO 5.6.3.4 "manage disposition of evaluation data" — do
  we prune old bundles, or leave that to the user / `.gitignore`?
