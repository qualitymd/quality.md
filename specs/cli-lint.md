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
(`../quality.md/docs/spec.md`) without any judgment about whether the
requirements are *good* — that is `evaluate-model`' job (see
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

## Linting rules

The linter runs a fixed set of rules against a parsed `QUALITY.md` — over the
frontmatter and the Markdown body. Each rule produces findings at a fixed severity.
Rule identifiers are illustrative and
expected to be tuned during implementation. Severities are chosen so a **minimal,
well-formed file still passes CI**: only unambiguous violations are errors.

### Frontmatter rules

| Rule | Severity | What it checks |
| --- | --- | --- |
| `parse-error` | error | Frontmatter missing, unterminated, or not valid YAML (no opening/closing `---`, malformed YAML). Aborts the run — no later rules can run. |
| `missing-factors` | error | The required top-level `factors` key is absent or empty. |
| `factor-shape` | error | A factor declares **neither** `requirements` nor `factors` — the spec requires at least one. |
| `assessment-count` | error | A requirement does not declare exactly one assessment: zero, or both `prompt` and `bash`. |
| `assessment-shape` | error | A `prompt` or `bash` value is not a single scalar string — e.g. a YAML list/sequence or a map. A requirement carries one prompt or one command, never a list of them. |
| `broken-prompt-ref` | error | A `prompt` given as a path (rather than inline text) points to a file that does not exist. |
| `broken-target` | warning | A `target` path or glob resolves to no files on disk. Warning, not error — a glob may legitimately match nothing yet. |
| `empty-collection` | warning | A `requirements` or `factors` map is present but empty. |
| `ratings-shape` | warning | `ratings` is present but malformed: an entry missing `displayName`, or fewer than two levels defined (a scale needs at least two). |
| `unknown-key` | warning | A key looks like a typo of a known schema key (`factor:` → `factors:`, `requirement:` → `requirements:`, `prompts:` → `prompt:`, `rating:` → `ratings:`). Genuinely custom extension keys stay silent — the format grows through users, like design.md. |
| `model-summary` | info | Summary counts: factors, leaf requirements, and the split of assessment types (`prompt` vs `bash`). |

> **`prompt` vs path.** A `prompt` value is treated as a file reference when it
> resolves to an existing path and as inline criteria text otherwise — matching
> the spec's "text \| path" union. `broken-prompt-ref` only fires for values that
> *look* like a path (e.g. start with `./`, `../`, or `/`, or end in `.md`) but
> do not exist, so inline prose is never mistaken for a broken reference.

### Body rules

The linter parses the body only far enough to extract its `##`/`###` headings — it
never judges prose. The recommended body sections and what each captures are
defined in the format spec (`../quality.md/docs/spec.md#markdown-body`); these rules
check their *shape*. All of them mirror a `design.md lint` shape except
`factor-without-prose` (see [The one rule past precedent](#the-one-rule-past-precedent)).

| Rule | Severity | What it checks |
| --- | --- | --- |
| `missing-overview` | warning | The body has no **Overview** section (or leading prose). |
| `missing-factors-section` | warning | The body has no **Factors** section. |
| `section-order` | warning | Recognized `##` sections appear out of canonical order. |
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

## Output

JSON by default — agent-consumable, the same contract as `design.md lint`.

```json
{
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
  "summary": { "errors": 1, "warnings": 1, "info": 1 }
}
```

Each finding carries `severity`, the `rule` that produced it, a `path` (a dotted
locator into the model, quoting map keys that contain spaces), and a `message`.
The `summary` tallies findings by severity.

## Flags, exit codes

Flags (shared flags are in [`cli.md`](./cli.md#shared-conventions)):

- `file` — positional path to the `QUALITY.md` file, or `-` for stdin. Defaults
  to `./QUALITY.md` / `-f`.
- `--json` — emit JSON output. JSON only in v1 (and the `lint` default), so the
  flag is a no-op for now; a human-readable text format is a possible later
  addition.

Exit codes:

- **`1`** if any finding is an `error`.
- **`0`** otherwise — warnings and info do not fail the gate.

This makes `lint` the deterministic structural gate for CI and pre-commit: it
fails the build only on real spec violations, while warnings surface
lower-confidence issues without blocking.

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
