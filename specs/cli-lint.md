# CLI: `lint`

Detail doc for the fast, deterministic structural tier. See [`cli.md`](./cli.md)
for the full command surface and [`../SPECIFICATION.md`](../SPECIFICATION.md) for
the normative format.

```bash
qualitymd lint QUALITY.md
qualitymd lint --json QUALITY.md
cat QUALITY.md | qualitymd lint -
```

## Purpose

`lint` answers one question: **is this a well-formed `QUALITY.md`?** It parses
the target-node frontmatter, validates the recursive schema, resolves local
references, checks body shape, and exits non-zero only on structural errors. It
does not judge whether the requirements are good; skills do that.

It checks, in order:

1. **Parse** - fenced frontmatter exists and is valid YAML.
2. **Schema** - the root is a target node; every target, factor, requirement,
   and rating scale is well-shaped.
3. **References** - `assessment`, `source`, and shared `ratings` paths resolve
   when they are path-shaped.
4. **Inheritance** - factor refinement, secondary-factor references,
   and federation/baseline inheritance are structurally coherent.
5. **Body** - recommended sections and target-scoped factor prose are present in
   a parseable shape.

When more than one `QUALITY.md` is discovered, lint validates each model and then
runs federation rules over the discovered target tree.

## Rule Idiom

Rules are pure checks over the parsed model plus on-disk reference existence. A
finding has `severity`, `rule`, `path`, and `message`. Paths are dotted locators
into the target tree, quoting map keys that contain spaces, for example
`targets.api.factors.security.requirements."no secrets are committed"`.

Severities:

- `error` - definite spec violation; exits `1`.
- `warning` - likely issue or low-confidence issue; does not fail the gate.
- `info` - summary or non-blocking context.

Unknown extension keys are preserved silently unless they look like typos of
known keys.

## Frontmatter Rules

| Rule | Severity | What it checks |
| --- | --- | --- |
| `parse-error` | error | Missing, unterminated, or invalid YAML frontmatter. Aborts later rules. |
| `target-node-shape` | error | The root and every `targets.*` entry is a mapping target node or the scalar `source` shorthand. Known target-node fields have the right shapes. |
| `open-target-vocabulary` | info | Target names are user vocabulary. A catalog miss may be reported as context, never as a warning or error. |
| `source-shape` | error | `source` is a scalar path/glob/URL or a list of scalar entries. |
| `broken-source` | warning | A path/glob `source` resolves to no files. A forward-declared glob may be valid, so this is not an error. |
| `source-escapes-scope` | warning | A local path/glob `source` resolves outside the model's directory subtree. |
| `factor-shape` | error | A factor entry is not a mapping, or known factor fields have the wrong shape. |
| `factor-refinement` | error | A descendant factor with an inherited name redefines the inherited lens instead of only adding compatible detail or requirements. |
| `requirement-shape` | error | A requirement entry is not a mapping or declares no non-empty `assessment`. |
| `assessment-shape` | error | `assessment` is not a single scalar string. A requirement has one assessment, never a list or map. |
| `broken-assessment-ref` | error | A path-shaped `assessment` points to a file that does not exist. |
| `secondary-factor-scope` | error | A requirement's secondary `factors` entry names no factor visible on the current target or an ancestor. |
| `containment-resolution` | error | The inherited target/factor/requirement set cannot be resolved deterministically. |
| `baseline-resolution` | warning | A configured rolling baseline cannot be loaded. Baseline content must be visible when configured. |
| `broken-ratings-ref` | error | A path-shaped top-level `ratings` value does not exist or does not parse as a rating scale. |
| `ratings-shape` | error | A rating scale is not a sequence with at least two unique `level` entries ordered best to worst. |
| `unknown-rating-level` | error | A per-requirement `ratings` map names a level not in the active scale. |
| `rating-level-order` | warning | A per-requirement `ratings` map is written in an order different from the active scale. Fixable. |
| `empty-collection` | warning | `targets`, `requirements`, or `factors` is present but empty. |
| `unknown-key` | warning | A key looks like a typo of a known schema key: `target` -> `targets`, `sources` -> `source`, `prompt`/`prompts` -> `assessment`, `factor` -> `factors`, `requirement` -> `requirements`, `rating` -> `ratings`. |
| `model-summary` | info | Summary counts: targets, factors, direct requirements, lensed requirements, and secondary-factor links. |

The table names every condition in the five schema rules:

- Rule 1, open target vocabulary: `open-target-vocabulary`.
- Rule 2, scoped factor identity/refinement: `factor-refinement`.
- Rule 3, containment inheritance: `containment-resolution`.
- Rule 4, rolling baseline: `baseline-resolution`.
- Rule 5, nest vs. federate: federation rules below plus
  `target-node-shape`.

`assessment` and `ratings` values are treated as file references only when they
are path-shaped (`./`, `../`, `/`, URL, or a recognized file extension);
otherwise they are inline text.

## Body Rules

The linter extracts headings only. It never judges prose.

| Rule | Severity | What it checks |
| --- | --- | --- |
| `missing-overview` | warning | The body has no **Overview** section or leading prose. |
| `missing-targets-factors-section` | warning | The body has no **Targets and factors** or **Factors** section. |
| `section-order` | warning | Recognized `##` sections appear out of canonical order. Fixable. |
| `duplicate-section` | error | A recognized `##` heading appears more than once. |
| `target-without-prose` | warning | A frontmatter target has no matching subsection under **Targets and factors**. |
| `factor-without-prose` | warning | A factor has no matching prose under the target where it is declared. |

Canonical body sections, in order:

| Canonical | Aliases |
| --- | --- |
| Overview | Summary |
| Scope | Boundary, In scope, Out of scope |
| Needs | Quality needs |
| Risks | Risk, Stakes, What's at stake |
| Targets and factors | Factors, The quality model, The model, Quality factors |
| Known gaps | Accepted risks, Limitations |

`factor-without-prose` is target-scoped: `targets.api.factors.security` is matched
against prose for the `api` target, not against an unrelated `security` factor on
another target.

## Federation Rules

When a federation is discovered, single-model rules run per file and these rules
run over the set:

| Rule | Severity | What it checks |
| --- | --- | --- |
| `mixed-rating-scales` | warning | Models in one federation use different active rating scales. A tree report is commensurable only on a shared scale. |
| `federation-graft` | error | A child model cannot be grafted as a target subtree because its root target address conflicts with an ancestor target address. |
| `cross-file-factor-redefinition` | error | A federated child redefines an inherited factor meaning instead of refining it. |
| `source-overlap` | warning | Ancestor and descendant `source` globs overlap in a way that may double-govern the same files. |

Coverage and semantic consistency are judgment, not lint; they belong to the
`improve-quality-md` and `evaluate-quality` skills.

## Output

JSON is the v1 contract:

```json
{
  "schemaVersion": "1",
  "findings": [
    {
      "severity": "error",
      "rule": "requirement-shape",
      "path": "targets.api.requirements.\"accepted orders are durable\"",
      "message": "Requirement declares no `assessment`."
    },
    {
      "severity": "warning",
      "rule": "broken-source",
      "path": "targets.web.source",
      "message": "source glob \"./web/**/*.ts\" matched no files."
    },
    {
      "severity": "error",
      "rule": "secondary-factor-scope",
      "path": "targets.api.factors.security.requirements.\"no secrets are committed\".factors[0]",
      "message": "Secondary factor \"reliability\" is not declared on this target or an ancestor."
    },
    {
      "severity": "info",
      "rule": "model-summary",
      "message": "3 targets, 4 factors, 2 direct requirements, 5 lensed requirements."
    }
  ],
  "summary": { "errors": 2, "warnings": 1, "info": 1 },
  "nextActions": [
    {
      "command": "qualitymd lint",
      "reason": "Fix `targets.api.requirements.\"accepted orders are durable\"`: declare an `assessment`, then re-run.",
      "priority": "required"
    }
  ]
}
```

## Flags And Exit Codes

Flags:

- `file` - positional `QUALITY.md` path, or `-` for stdin. Defaults to
  `./QUALITY.md`. With no explicit file and multiple discovered models, validates
  the federation.
- `--json` - emit JSON. JSON is the default v1 output.
- `--fix` - apply canonical-order fixes (`section-order`, `rating-level-order`)
  and re-run. It never changes semantic findings.

Exit codes:

- `0` - no `error` findings.
- `1` - at least one `error` finding.
- `2` - the command could not run: bad flags, unreadable file, or internal error.

## Open Questions

- Whether `broken-source` should remain a warning or become an error in stricter
  modes.
- Whether to persist a pretty text output in addition to JSON.
- Whether target/factor prose matching should require nested heading depth or only
  normalized names.
- Whether an `owner:` extension becomes a first-class target-node field later.
