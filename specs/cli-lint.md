# CLI: `lint`

> Detail doc for the fast, deterministic **structural tier**. See
> [`cli.md`](./cli.md) for the full command surface and shared conventions, and
> [`cli-evaluate.md`](./cli-evaluate.md) for the deep semantic tier.
>
> This command is modeled on Google's `design.md lint` — the closest prior art
> for QUALITY.md. The shape is deliberately the
> same: parse the file, run a fixed set of rules, emit structured JSON findings
> an agent can act on, exit non-zero on errors.

```bash
qualitymd lint QUALITY.md
qualitymd lint --json QUALITY.md
cat QUALITY.md | qualitymd lint -
```

## Purpose

`lint` answers one question, deterministically: **is this a well-formed
`QUALITY.md`?** It validates the file against the format spec
(`../SPECIFICATION.md`) without any judgment about whether the
requirements are *good* — that is `evaluate-model`'s job (see
[`cli-evaluate.md`](./cli-evaluate.md#one-engine-two-targets)). Because it is
pure parsing and static checks, it is cheap enough to run on every save, in a
pre-commit hook, and as the structural CI gate.

It checks four things, in order:

1. **Parse** — the frontmatter exists, is fenced by `---`, and is valid YAML.
2. **Schema** — the parsed model conforms to the spec: `factors` present; every
   factor carries at least one of `requirements`/`factors` (either, or both);
   every requirement declares exactly one assessment; `ratings` (if present) is
   well-shaped.
3. **References** — `prompt` paths and `target` paths/globs resolve on disk
   (model-relative; see [`cli.md`](./cli.md#shared-conventions)).
4. **Body** — the Markdown body carries the spine sections and each factor has
   matching prose (see [Body rules](#body-rules)).

When the working tree holds more than one `QUALITY.md`, `lint` discovers and
validates the whole [federation](./cli-federation.md) by default: the rules below
run per model, and a small set of [cross-file rules](#federation-rules) run over
the set.

## Linting rules

The linter runs a fixed set of rules against a parsed `QUALITY.md` — over the
frontmatter and the Markdown body. Each rule produces findings at a fixed severity.
Rule identifiers are illustrative and
expected to be tuned during implementation. All rules share one shape and one
severity philosophy, stated next.

### Rule idiom

Every lint rule follows one shape, borrowed from `design.md lint`'s rule
descriptor. Stating it once keeps the rule set uniform: a new rule is read
against this contract, not invented fresh. The rule tables below are this
descriptor in tabular form — the **Rule** column is `name`, **Severity** is the
default `severity`, and **What it checks** is `description`.

- **One rule, one defect, pure.** A rule is a named unit `{ name, severity,
  description, run }`. `run` is a pure function of the *parsed model* — parsed
  frontmatter and extracted body headings in, findings out. Its only I/O is the
  on-disk existence checks for `prompt`/`target` references; it never mutates the
  model, rewrites the file, or judges prose.
- **Names name the fault, kebab-case.** A rule reads as the thing it flags —
  either the defect directly (`missing-factors`, `broken-target`,
  `duplicate-section`) or the element-and-aspect it validates
  (`assessment-count`, `factor-shape`, `section-order`). Never the verb of the
  check. The `-summary` suffix is reserved for the `info` rule that *reports*
  rather than flags (`model-summary`).
- **Three severities, gate-safe by default.** `error` — an unambiguous spec
  violation; fails the gate (exit `1`). `warning` — lower-confidence, or
  legitimately fine sometimes (a forward-declared glob, an empty collection);
  surfaced, never blocking. `info` — a summary, never a defect. The ladder is set
  so a **minimal, well-formed file passes CI**: when in doubt, a rule is a
  warning, not an error.
- **Findings are `{ severity, path, message }`.** `path` is the dotted locator
  into the model (quoting map keys that contain spaces; see [Output](#output));
  body-level rules may omit it. `message` is one complete sentence stating the
  violation and, where the rule checks against a closed set, naming the valid
  alternatives. A finding may carry its own severity, overriding the rule's
  default — the descriptor allows one rule to span severities (as design.md's
  `broken-ref` does), though no current QUALITY.md rule needs it. The related
  discretion every rule does exercise is *whether to emit at all*: `unknown-key`
  warns on a typo-shaped key but stays silent on a genuinely custom one.
- **A rule may be fixable.** A rule whose remedy is mechanical and
  behavior-preserving — reordering levels or sections that carry no semantic
  weight — may carry an optional `fix` alongside `run`: also a *pure* function
  (parsed model → corrected model), so the rule stays side-effect-free and only
  the command layer writes. `lint --fix` applies these in place, and the finding
  carries the patch so an agent can apply it straight from the JSON. Only
  **canonical-order** rules qualify — `section-order` and `rating-level-order` —
  because their fix cannot change any evaluation result. A rule that flags a
  *semantic* violation (a duplicate, an unreachable level, an unknown level) is
  never auto-fixed: choosing the remedy needs author intent, so it stays a
  `nextActions` suggestion.
- **The format grows through users.** Unrecognized keys and body sections are
  preserved silently; only typo-shaped ones warn. A rule never errors on
  extension it does not understand.

### Frontmatter rules

| Rule | Severity | What it checks |
| --- | --- | --- |
| `parse-error` | error | Frontmatter missing, unterminated, or not valid YAML (no opening/closing `---`, malformed YAML). Aborts the run — no later rules can run. |
| `missing-factors` | error | The required top-level `factors` key is absent or empty. |
| `factor-shape` | error | A factor declares **neither** `requirements` nor `factors` — the spec requires at least one. |
| `assessment-count` | error | A requirement does not declare exactly one assessment: zero, or both `prompt` and `bash`. |
| `assessment-shape` | error | A `prompt` or `bash` value is not a single scalar string — e.g. a YAML list/sequence or a map. A requirement carries one prompt or one command, never a list of them. |
| `broken-prompt-ref` | error | A `prompt` given as a path (rather than inline text) points to a file that does not exist. |
| `broken-ratings-ref` | error | A `ratings` value given as a path (a shared scale file) points to a file that does not exist or does not parse as a rating scale. |
| `broken-target` | warning | A `target` path or glob resolves to no files on disk. Warning, not error — a glob may legitimately match nothing yet. |
| `target-escapes-scope` | warning | A `target` resolves outside the model's own directory subtree (e.g. climbs above it with `../`), re-governing code another model owns. Federation keeps a requirement's reach local; an outward target is a smell, not always a mistake. |
| `empty-collection` | warning | A `requirements` or `factors` map is present but empty. |
| `ratings-shape` | warning | `ratings` is present but malformed: fewer than two levels defined (a scale needs at least two). |
| `unknown-rating-level` | error | A per-requirement `ratings` override names a level not defined in the scale. The spec treats this as a configuration error: an override may only re-state conditions for levels the scale already declares. |
| `duplicate-rating-level` | error | Two entries in the `ratings` sequence — inline or shared file — declare the same `level` name. As an ordered sequence the scale no longer gets duplicate-key rejection from YAML, so the linter enforces level-name uniqueness (mirrors `duplicate-section`). |
| `unreachable-rating-level` | warning | A `ratings` level can never be selected because an earlier, better-ranked level's `bashCondition` subsumes it under first-match-wins (e.g. a looser numeric threshold sitting above a stricter one). Best-effort over recognizable numeric comparisons; silent when conditions are not statically comparable, and never fired on a condition-less intermediate band a `prompt` scale legitimately uses. |
| `rating-level-order` | warning | A per-requirement `ratings` override lists its levels in a different order than the scale (the configured frontmatter scale, or the default `outstanding`/`target`/`minimum`/`unacceptable`, best to worst). The override is a by-name patch, so order is cosmetic — matching scale order keeps it scannable. **Fixable** (`lint --fix` reorders to scale order); mirrors `section-order`. |
| `unknown-key` | warning | A key looks like a typo of a known schema key (`factor:` → `factors:`, `requirement:` → `requirements:`, `prompts:` → `prompt:`, `rating:` → `ratings:`). Genuinely custom extension keys stay silent — the format grows through users, like design.md. |
| `model-summary` | info | Summary counts: factors, leaf requirements, and the split of assessment types (`prompt` vs `bash`). |

> **`prompt` vs path.** A `prompt` value is treated as a file reference when it
> resolves to an existing path and as inline criteria text otherwise — matching
> the spec's "text \| path" union. `broken-prompt-ref` only fires for values that
> *look* like a path (e.g. start with `./`, `../`, or `/`, or end in `.md`) but
> do not exist, so inline prose is never mistaken for a broken reference.

> **`ratings` inline vs path.** Likewise, a `ratings` value is treated as a
> shared-scale file reference when it is a path string and as an inline scale sequence
> otherwise; `broken-ratings-ref` fires only for a path-shaped value that does not
> resolve or parse as a scale.

### Body rules

The linter parses the body only far enough to extract its `##`/`###` headings — it
never judges prose. The recommended body sections and what each captures are
defined in the format spec (`../SPECIFICATION.md#markdown-body`); these rules
check their *shape*. All of them mirror a `design.md lint` shape except
`factor-without-prose` (see [The one rule past precedent](#the-one-rule-past-precedent)).

| Rule | Severity | What it checks |
| --- | --- | --- |
| `missing-overview` | warning | The body has no **Overview** section (or leading prose). |
| `missing-factors-section` | warning | The body has no **Factors** section. |
| `section-order` | warning | Recognized `##` sections appear out of canonical order. **Fixable** (`lint --fix` reorders to canonical order). |
| `duplicate-section` | error | A `##` heading appears more than once. |
| `factor-without-prose` | warning | A frontmatter factor has no matching `###` subsection under **Factors**. |

Only **Overview** and **Factors** are presence-checked, and only as warnings; every
other section is optional and never warned for absence. The sole body-level error is
`duplicate-section`, an unambiguous mistake.

#### Recognized sections

Canonical section names and their aliases, in canonical order, resolved before
`section-order` and the presence rules run — the same mechanism as design.md's
`spec-config.yaml` + `resolveAlias()`. Unknown sections are preserved silently; the
format grows through users.

| Canonical | Aliases |
| --- | --- |
| Overview | Summary |
| Scope | Boundary, In scope, Out of scope |
| Needs | Quality needs |
| Risks | Risk, Stakes, What's at stake |
| Factors | The quality model, The model, Quality factors |
| Known gaps | Accepted risks, Limitations |

#### The one rule past precedent

`factor-without-prose` is the single rule that steps beyond `design.md lint`. It
matches frontmatter factor *names* against body *headings* — a frontmatter↔prose
cross-reference, which design.md deliberately never does (it cross-references only
YAML against YAML). We cross that line on purpose: this rule is what enforces the
"prose version of the model," and it stays low-risk by comparing only heading text
to factor keys (case-insensitive), never parsing prose. It is a **warning**, so
model/prose drift nudges without blocking the gate.

### Federation rules

When `lint` runs over a [federation](./cli-federation.md) — more than one
`QUALITY.md` discovered — most rules above run per model. A few defects are only
visible across the *set*, so they run once over the discovered models; unlike the
single-model rules, these take the whole set as input rather than one parsed
model.

| Rule | Severity | What it checks |
| --- | --- | --- |
| `mixed-rating-scales` | warning | Models in one federation define or reference **different** rating scales. A tree report is only commensurable on a shared scale; share one by reference (see [`cli-federation.md`](./cli-federation.md#shared-rating-scale)). Sometimes a subtree deliberately rates on its own scale, hence a warning. |
| `target-overlap` | warning | An ancestor and a descendant requirement's `target` globs resolve to overlapping files, so the same code is governed two ways. Often intentional for a cross-cutting ancestor requirement, hence a warning, not an error. |

Coverage ("does every significant component have a model?") and the cross-file
*consistency of meaning* between requirements are judgments, not structural
facts, so they belong to `evaluate-model`'s set-level pass (see
[`cli-evaluate.md`](./cli-evaluate.md#set-level-evaluation)), not `lint`.

## Output

JSON by default — agent-consumable, the same contract as `design.md lint`, and
carrying the two fields [`cli.md`](./cli.md#machine-readable-result-contract)
requires of *every* command: a `schemaVersion` (so an agent can parse results
without screen-scraping) and a `nextActions` array (so the lint→fix→re-run loop
is self-describing rather than tribal knowledge).

```json
{
  "schemaVersion": "1",
  "findings": [
    {
      "severity": "error",
      "rule": "assessment-count",
      "path": "factors.security.requirements.\"no secrets committed to the repository\"",
      "message": "Requirement declares both `prompt` and `bash`; exactly one assessment is required."
    },
    {
      "severity": "warning",
      "rule": "broken-target",
      "path": "factors.maintainability.factors.reusability.requirements.\"shared domain types come from the common package\".target",
      "message": "target glob \"./src/**/*.ts\" matched no files."
    },
    {
      "severity": "info",
      "rule": "model-summary",
      "message": "3 factors, 6 requirements (4 prompt, 2 bash)."
    }
  ],
  "summary": { "errors": 1, "warnings": 1, "info": 1 },
  "nextActions": [
    {
      "command": "qualitymd lint",
      "reason": "Fix `factors.security.requirements.\"no secrets committed to the repository\"`: declare exactly one of `prompt`/`bash`, then re-run.",
      "priority": "required"
    }
  ]
}
```

Each finding carries `severity`, the `rule` that produced it, a `path` (a dotted
locator into the model, quoting map keys that contain spaces), and a `message`.
The `summary` tallies findings by severity.

`schemaVersion` is the stable top-level version string from the shared contract.
`nextActions` follows the shared
[next-action shape](./cli.md#structured-next-action-suggestions): each error
finding emits a `required` action naming the offending `path`, each warning a
`recommended` one; a clean file emits a single `recommended` action — run
`evaluate-model` to pressure-test the requirements. The `priority` of a
finding-derived action tracks the finding's severity (`error → required`,
`warning → recommended`, `info → optional`), as in
[`cli.md`](./cli.md#structured-next-action-suggestions).

## Flags, exit codes

Flags (shared flags are in [`cli.md`](./cli.md#shared-conventions)):

- `file` — positional path to the `QUALITY.md` file, or `-` for stdin. Defaults
  to `./QUALITY.md` / `-f`. With no `file`/`-f` and more than one `QUALITY.md`
  present, `lint` validates the whole [federation](./cli-federation.md) (see
  [Federation rules](#federation-rules)).
- `--json` — emit JSON output. JSON only in v1 (and the `lint` default), so the
  flag is a no-op for now; a human-readable text format is a possible later
  addition.
- `--fix` — apply the auto-fixable findings (the canonical-order rules
  `section-order` and `rating-level-order`) in place, then re-run the checks and
  report what remains. Behavior-preserving by construction — it only reorders
  levels or sections that carry no evaluation weight — so it never changes a
  verdict and never touches a semantic finding. Without `--fix`, `lint` only
  reports.

Exit codes follow the shared three-code convention (see
[`cli.md`](./cli.md#machine-readable-result-contract)):

- **`0`** — no `error` findings. Warnings and info do not fail the gate.
- **`1`** — **gate verdict failure:** at least one finding is an `error` (a real
  spec violation, including `parse-error`). The file *was* read and checked; it
  is simply not well-formed.
- **`2`** — **tool failure:** `lint` could not run the checks at all — the file
  is unreadable or absent, a bad flag was passed, or an internal error occurred.
  This is distinct from `parse-error` (a malformed-but-readable file), which is a
  finding and exits `1`.

This makes `lint` the deterministic structural gate for CI and pre-commit: it
fails the build (`1`) only on real spec violations, while warnings surface
lower-confidence issues without blocking — and an agent can still tell a bad file
(`1`) from a broken invocation (`2`).

## Open questions

- **`broken-target` severity.** Warning (current) vs. error. A glob matching
  nothing is sometimes intentional (forward-declared paths) and sometimes a
  typo; warning is the safer default but loses gate strength.
- **Text output format.** JSON is the v1 contract for agents; whether to add a
  pretty/grouped text renderer (like `qualitymd check`'s original report) for
  human terminal use is open.
- **`bash` command sanity.** Should `lint` do any static check on `bash`
  assessment strings (e.g. non-empty), or leave all command validity to
  `evaluate`? Currently out of scope — `lint` validates structure, not commands.
- **Subfactor heading depth.** `factor-without-prose` assumes factors map to `###`
  subsections. How deep nested subfactors are expected to go in prose (`####` and
  beyond) before the coverage check stops insisting is unsettled.
- **`owner:` frontmatter field.** Research favored recording *who is accountable*
  for the model (routing CI failures, codeowners-style automation) — the one
  structured field clean enough to consider. Deferred: it overlaps the Overview/Needs
  prose and risks staleness. Other proposed enums (`priority`, `criticality`,
  `severity`, `source`) are deliberately *not* adopted — they invite box-ticking and
  false precision, and `--fail-on` already carries that weight.
