---
type: Functional Specification
title: qualitymd lint
description: Validate a QUALITY.md file's structure against the format spec.
tags: [cli, command, lint]
timestamp: 2026-06-17T00:00:00Z
---

# qualitymd lint

`lint` inherits the cross-cutting CLI contract — invocation, global flags, output
formats, exit codes, and agent accessibility — from the [CLI spec](../cli.md).
This file specifies only what is particular to `lint`.

`qualitymd lint` validates a `QUALITY.md` file's structure against the
[format specification](../../SPECIFICATION.md), fast and deterministically,
exiting non-zero on errors so it drops into CI.

**Boundary.** `lint` checks *format conformance* — whether the file is a valid
`QUALITY.md` — only. It does not assess whether the model is a *good* quality
model; that judgment lives in the evaluation skills, not the deterministic CLI.

The key words "MUST", "MUST NOT", "SHOULD", "RECOMMENDED", and "MAY" are to be
interpreted as described in IETF RFC 2119.

## Scope

Covered: the rule boundary for mechanical format validation, rule metadata,
finding messages, finding locations, human and JSON output, deterministic
ordering, in-place repair of fixable findings, and the initial rule set.

Deferred: suppression directives, rule selection or severity overrides, emitting
a rule catalog from `lint` itself, and repair output modes other than in-place
writes (for example patch output or emitting a full rewritten file to stdout).

## Flags

`lint` inherits the cross-cutting flags from the [CLI spec](../cli.md), including
`--json`, and adds one command-specific flag:

- `--fix` — apply every fixable finding that can be repaired deterministically,
  writing the repaired `QUALITY.md` back to the same path.

`--fix` **MAY** be combined with `--json`; the JSON output then reports the
post-repair findings and the repairs that were applied.

## Rule scope

A `QUALITY.md` is validated by a set of named **rules**, each checking one
property and emitting a **finding** when violated. A rule belongs in `lint` only
when **all four** of the following hold:

1. **Conformance-grounded.** It enforces a constraint that already exists in the
   [format specification](../../SPECIFICATION.md). `lint` introduces no
   requirements of its own; each rule traces to the spec clause it checks.
2. **Deterministic and mechanical.** Its verdict is computable from the file's
   structure alone — no judgment, no external data, the same result every run.
3. **Format, not goodness.** It decides whether the file *is a valid
   `QUALITY.md`*, never whether it is a *good* quality model. Anything requiring
   evaluative judgment is out of scope — that is the evaluation skills' work.
4. **Self-contained.** It is decidable from the single file plus its own declared
   cross-references (e.g. resolving a secondary factor name, a `ratings`
   override key). A rule MUST NOT resolve a target's `source` or read the
   entities a target evaluates; that is beyond format conformance.

Structural validation — valid keys, required/recommended/optional properties,
YAML value shapes, the model-content group, and the rating-scale minimum — MUST
derive from the single structural schema declaration that the linter consumes.
Rule logic MUST NOT maintain a second valid-key list independent of that schema.

### Severity

Every rule carries a fixed **severity** that governs how its finding surfaces:

- **error** — a spec **MUST**/**MUST NOT** is violated; the file is not a valid
  `QUALITY.md`. Errors are what make `lint` exit non-zero.
- **warning** — a spec **SHOULD**/**RECOMMENDED** that is still mechanically
  determinable is unmet (e.g. a recommended `description` is absent, or a target
  whose subtree reaches no requirement). Warnings do not affect the exit code.
- **info** — a non-judgmental structural observation (e.g. a model summary), not
  a defect.

## Rule naming

A rule's id is a stable public identifier: it appears in `--json` findings and in
any future suppression directive, so a rename is a breaking change. Rule ids
MUST follow these guidelines:

- **kebab-case**, lowercase, naming the **condition checked** as a noun phrase —
  not an imperative. Prefer `missing-criterion` over `criterion-required` or
  `no-missing-criterion`.
- Use one of two shapes: **`<defect>-<subject>`** for the defect a rule rejects
  (`missing-rating-scale`, `duplicate-level`, `unknown-factor`) or
  **`<subject>-<aspect>`** for a neutral observation about a subject (e.g. a
  future `model-summary` info rule). Every rule in the current set is
  `<defect>-<subject>`; the second shape is reserved for the `info` observations
  defined later.
- Name the subject in `QUALITY.md` vocabulary — Target, Factor, Requirement,
  Rating Scale, Rating Level, Assessment — never ISO terms or implementation
  names.
- Keep it short; prefer two words. Do not prefix ids with `quality-` or similar;
  the context is implied.
- One concept per rule. When a rule would check two unrelated properties, split
  it into two rules.

## Rule authoring

Each rule carries a static **description** (a catalog entry for the rule itself)
and emits a **finding message** per violation (the concrete instance). The two
serve different readers and MUST NOT duplicate each other: the description is the
generic rule, the message is the specific occurrence. Reading one MUST NOT feel
like reading the other twice.

### Descriptions

A rule's `description` describes the check, for someone browsing the rule
catalog. It MUST:

- be **one present-tense sentence**, generic — describing the check, never a
  specific occurrence, with no interpolated values;
- **state the triggering condition**, leading with the subject in `QUALITY.md`
  vocabulary — e.g. "Flags a requirement whose `assessment` is missing, empty,
  or a list.";
- **name the exact frontmatter key it inspects** in backticks (`ratingScale`,
  `criterion`, `assessment`) — unless the rule targets the frontmatter block as a
  whole or any optional property, with no single key to name;
- **signal severity through its verb** — error rules *flag* or *reject*, warning
  rules *warn*, info rules *report* or *summarize* — so the catalog reads
  consistently; and
- stay **descriptive, not prescriptive** (the fix belongs in the message) and add
  information beyond the rule id rather than re-spelling it.

### Messages

A finding **message** describes one violation, for the author or agent who hit
it. It MUST:

- **name the specific element and where it is** — the requirement statement,
  factor name, or rating level — quoted by its key or name, not by reproducing
  its full content. (The structured `location` carries this for `--json`; the
  prose still names it for human and plain output.)
- **contrast found against required** — what is actually there versus what the
  format spec requires;
- **be self-contained and deterministic** — understandable without the rule
  description, the same for the same input, in `QUALITY.md` vocabulary with no
  implementation terms; and
- **carry tone to match severity** — errors state the violation plainly,
  warnings frame the recommendation they enforce, info is neutral observation; no
  blame, no exclamation.

A message SHOULD **be actionable when the fix is determinate** — pointing to the
expected shape or the valid set — and SHOULD stay concise, leading with the
problem. A rule that finds several instances SHOULD emit **one message per
instance**, each with its own location, rather than one bundled list, so each is
independently addressable.

### Fixability

A rule is **fixable** only when `qualitymd` can repair each finding with a
deterministic edit that does not presume the author's intent. Fixability is
independent of severity: an error MAY be fixable, and a warning MAY be
non-fixable.

A rule **MUST** be marked fixable only when all of the following hold:

- **Single mechanical edit.** The finding has one correct structural repair,
  with no choice among valid outcomes.
- **No presumed intent.** The repair does not decide what the quality model
  should say, choose domain content, infer stakeholder intent, decide that
  misplaced content was accidental, or rank quality levels.
- **Content-preserving.** The repair does not discard authored content, except
  when the content is an empty optional property whose absence is the required
  repair.
- **Strictly improving.** Applying the repair removes that finding and MUST NOT
  introduce another finding.
- **Idempotent.** Applying the repair twice has the same result as applying it
  once.
- **Local and explainable.** The repair is attributable to the finding's
  location and can be described without referring to implementation details.
- **Stable formatting.** Any surrounding YAML or Markdown rewrite is
  deterministic and avoids unrelated churn.

A rule **MUST NOT** be marked fixable merely because a placeholder could be
inserted. Placeholders are scaffold content owned by [`init`](./init.md), not a
lint repair for missing model semantics.

### Repair behavior

When `--fix` is passed, `lint` **MUST** apply every emitted finding whose
`fixable` value is `true`, then lint the repaired file again and report the
post-repair findings. Non-fixable findings remain findings; they are never
silently changed or suppressed.

Repairs are **transactional per file**:

- `lint` **MUST** compute repairs from the original parsed document.
- If two repairs cannot be applied together, `lint` **MUST** leave the file
  unchanged and report a repair failure rather than applying a partial set.
- If a write fails, `lint` **MUST** report the failure and **MUST NOT** report the
  file as repaired.
- Applying `--fix` twice **MUST** have the same effect as applying it once.

When repair fails before a post-repair lint result exists, `lint` **MUST** exit
non-zero and report the failure through the CLI's error-reporting path. It
**MUST NOT** emit a successful lint result for a file it did not repair.

In-place writes **MUST** preserve authored content outside the repaired
frontmatter nodes:

- Markdown body content **MUST** be preserved byte-for-byte.
- YAML map order, comments, scalar style, and whitespace outside repaired nodes
  **SHOULD** be preserved where the parser and emitter make that possible.
- Unrelated YAML keys **MUST NOT** be reordered or rewritten unless the rewrite is
  necessary to apply the repair deterministically.
- The write **SHOULD** be atomic from the caller's perspective: write a complete
  replacement and then replace the target path, rather than truncating the
  original before the replacement is ready.
- To avoid ambiguous replacement behavior, `lint --fix` **SHOULD** refuse to
  repair a linted path that is a symbolic link until symlink write semantics are
  specified.

This phase defines only in-place repair through `--fix`. `lint` **MUST NOT** emit
patches or full rewritten files as alternate repair output modes in this phase.

## Findings and output

`lint` emits zero or more **findings**. Human-readable output and JSON output
MUST report the same findings for the same input; `--json` changes only the
format.

### Finding schema

Under `--json`, `lint` MUST emit one JSON document on stdout with this shape:

```json
{
  "schemaVersion": 1,
  "path": "QUALITY.md",
  "valid": false,
  "summary": {
    "errors": 1,
    "warnings": 0,
    "info": 0,
    "fixable": 0,
    "fixed": 0
  },
  "findings": [
    {
      "ruleId": "missing-rating-scale",
      "severity": "error",
      "message": "The model root declares no `ratingScale`; a QUALITY.md model requires one rating scale.",
      "location": {
        "path": "QUALITY.md",
        "modelPath": ["ratingScale"],
        "label": "ratingScale"
      },
      "fixable": false
    }
  ],
  "repairs": [],
  "nextActions": []
}
```

Fields are stable public API:

- `schemaVersion` **MUST** be `1` until the JSON shape changes incompatibly.
- `path` **MUST** identify the linted file using the caller-facing path.
- `valid` **MUST** be `true` when there are no `error` findings and `false`
  otherwise.
- `summary.errors`, `summary.warnings`, and `summary.info` **MUST** equal the
  number of findings at each severity.
- `summary.fixable` **MUST** equal the number of findings whose `fixable` value
  is `true`.
- `summary.fixed` **MUST** equal the number of repairs applied during this run;
  it is `0` when `--fix` is not passed.
- `findings` **MUST** contain every emitted finding in deterministic order.
- `repairs` **MUST** contain every applied repair in deterministic order. It is
  an empty array when `--fix` is not passed or no repairs were applied.
- `nextActions` follows the [CLI spec's convention](../cli.md#conventions).

Each finding object **MUST** contain:

- `ruleId` — the stable rule id from the [Rules](#rules) table.
- `severity` — one of `error`, `warning`, or `info`.
- `message` — the per-instance message described in [Messages](#messages).
- `location` — the location object described in [Locations](#locations).
- `fixable` — whether this finding has a deterministic repair under
  [Fixability](#fixability).

Each repair object **MUST** contain:

- `ruleId` — the rule id whose finding was repaired.
- `message` — a deterministic description of the edit that was applied.
- `location` — the location object for the repaired finding.

The JSON object **MUST NOT** include human styling, terminal control sequences,
or implementation-only fields. Additional documented fields MAY be added later
when they do not change the meaning of the required fields.

### Locations

A finding location names the smallest stable place the finding can be attached
to. It is not required to be a byte-perfect source span.

Under `--json`, `location` **MUST** contain:

- `path` — the linted file path, matching the top-level `path`.
- `modelPath` — an array path into the frontmatter model, using strings for map
  keys and numbers for list indexes. For example, the first rating level's
  missing `criterion` is `["ratingScale", 0, "criterion"]`; a requirement under
  a factor is `["factors", "<factor-name>", "requirements",
  "<requirement-statement>", "assessment"]`.
- `label` — a concise human-readable rendering of the same location, suitable
  for human output.

When a finding attaches to the frontmatter block as a whole,
`modelPath` **MUST** be an empty array and `label` **MUST** be `frontmatter`.
When a finding attaches to a missing key, `modelPath` **MUST** include the
missing key at the place it would appear.

If source position is available, `location` **SHOULD** also include 1-based
`line` and `column` fields for the start of the relevant YAML node. Source
positions are advisory: callers MUST treat `modelPath` as the stable machine
location.

### Human output

Human-readable output **MUST** include, for each finding, the severity, rule id,
message, and location label. It **SHOULD** summarize the total errors and
warnings. When `--fix` applies repairs, human-readable output **MUST** also
report how many repairs were applied. Styling and exact layout are governed by
the CLI's output conventions and are not part of this sub-spec.

When there are no findings, human-readable output **SHOULD** report that the file
is valid. Under `--json`, a valid file is represented by `"valid": true`, zero
counts in `summary.errors`, `summary.warnings`, and `summary.info`, and an empty
`findings` array.

### Exit status

Without `--fix`, `lint` exits non-zero when it emits one or more `error`
findings. Warnings and info findings do not affect the exit code.

With `--fix`, `lint` exits non-zero when repair fails or when the post-repair
lint result still contains one or more `error` findings. A run that fixes all
errors and leaves only warnings exits zero.

### Ordering and blocking

Findings **MUST** be emitted in deterministic order:

1. Earlier source position first, when both findings have source positions.
2. Otherwise, shallower `modelPath` first, then lexicographic comparison of path
   segments with numeric indexes ordered numerically.
3. For the same location, `error` before `warning` before `info`.
4. For the same location and severity, lexicographic `ruleId`.

A malformed or structurally invalid parent **MAY** block rules that depend on
that parent's parsed shape. For example, absent frontmatter or invalid YAML emits
`invalid-frontmatter` and prevents model-level rules from running; a malformed
`ratingScale` may prevent per-level checks from running. A blocked downstream
rule **MUST NOT** emit a speculative finding.

## Rules

The rules below are one per mechanically checkable constraint in the
[format specification](../../SPECIFICATION.md), derived under
[Rule scope](#rule-scope) and named under [Rule naming](#rule-naming). Each cites
the constraint it enforces; the *description* column follows
[Descriptions](#descriptions), and the fixability columns follow
[Fixability](#fixability). This initial set grows as the format spec does.

**Empty is absent.** A required property whose value is null or empty is treated
identically to a missing one and raises the same error (an empty `criterion`
raises `missing-criterion`); the format spec draws no distinction (see
[YAML Frontmatter](../../SPECIFICATION.md)). An optional property that is present
but null or empty is never reported by a `missing-*` warning — those fire only on
true absence — but by `empty-property`.

### Errors

Each enforces a **MUST** — its finding means the file is not a valid
`QUALITY.md`, and `lint` exits non-zero.

| Rule                   | Enforces (format spec)                                                                    | Description                                                                                       | Fixable | Fixable rationale                                                                  |
| ---------------------- | ----------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- | ------- | ---------------------------------------------------------------------------------- |
| `invalid-frontmatter`  | *YAML Frontmatter* — a file MUST begin with a valid YAML frontmatter block                | Rejects frontmatter that is absent, not valid YAML, or not the model's shape.                     | No      | Parse and shape failures do not have one safe structural repair.                   |
| `missing-rating-scale` | *Model* — `ratingScale` is required                                                       | Flags a model that declares no `ratingScale`.                                                     | No      | The scale defines model semantics; the format has no single required scale.        |
| `too-few-levels`       | *Rating Scale* — at least two rating levels MUST be supplied                              | Flags a `ratingScale` with fewer than two rating levels.                                          | No      | Adding levels requires choosing scale vocabulary and criteria.                     |
| `missing-level-name`   | *Rating Scale* — each level MUST declare a `level` name                                   | Flags a rating level that declares no `level` name.                                               | No      | A generated name would define rating vocabulary for the author.                    |
| `duplicate-level`      | *Rating Scale* — a `level` name MUST be unique within the scale                           | Flags two rating levels that share a `level` name.                                                | No      | Repair requires choosing which level to rename, merge, or remove.                  |
| `missing-criterion`    | *Rating Scale* — each level MUST declare a `criterion`                                    | Flags a rating level that declares no `criterion`.                                                | No      | Criterion text defines rating semantics and cannot be inferred mechanically.       |
| `empty-model`          | *Model* — an entry on `factors`, `requirements`, or `targets` MUST be supplied            | Flags a model root that supplies no entry under `factors`, `requirements`, or `targets`.          | No      | Repair requires choosing what kind of model content to add.                        |
| `misplaced-root-key`   | *Target* — a non-root target MUST NOT declare `title` or `ratingScale`                    | Flags a non-root target that declares `title` or `ratingScale`.                                   | No      | Repair requires deciding whether to remove, move, or reinterpret authored content. |
| `invalid-assessment`   | *Requirement* — a requirement MUST declare exactly one `assessment` as a non-empty scalar | Flags a requirement whose `assessment` is missing, empty, or a list rather than a single scalar.  | No      | Repair requires choosing the requirement's assessment text.                        |
| `unknown-factor`       | *Requirement* — each secondary `factors` name MUST resolve to a factor in scope           | Flags a requirement whose secondary `factors` entry names no factor on its target or an ancestor. | No      | Repair requires choosing the intended in-scope factor or adding a new one.         |
| `unknown-rating-key`   | *Requirement* — each `ratings` override key MUST name a level of the rating scale         | Flags a `ratings` override key that names no level in the model's `ratingScale`.                  | No      | Repair requires choosing the intended rating level or changing the scale.          |

### Warnings

Each enforces a mechanically determinable **SHOULD**/**RECOMMENDED** — its
finding is advisory and does not affect the exit code.

| Rule                         | Enforces (format spec)                                                       | Description                                                                                                  | Fixable | Fixable rationale                                                          |
| ---------------------------- | ---------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------ | ------- | -------------------------------------------------------------------------- |
| `missing-title`              | *Model* — a `title` is RECOMMENDED                                           | Warns when the model root declares no `title`.                                                               | No      | A title names the subject and cannot be inferred mechanically.             |
| `missing-level-description`  | *Rating Scale* — a level `description` is RECOMMENDED                        | Warns when a rating level declares no `description`.                                                         | No      | A description states the level's meaning and requires authored content.    |
| `missing-factor-description` | *Factor* — a factor SHOULD declare a `description`                           | Warns when a factor declares no `description`.                                                               | No      | A description states what the factor means and requires authored content.  |
| `empty-factor`               | *Factor* — a factor SHOULD lead to at least one requirement                  | Warns when a factor leads to no `requirements` nested under it, tagging it, or reached through a sub-factor. | No      | Repair requires adding, moving, or tagging requirements with model intent. |
| `empty-target`               | *Target* — each target SHOULD lead to a requirement somewhere in its subtree | Warns when a target's subtree reaches no `requirements`.                                                     | No      | Repair requires choosing target content or restructuring the target tree.  |
| `empty-property`             | *YAML Frontmatter* — null or empty optional properties SHOULD be omitted     | Warns when an optional property is present but null or empty instead of omitted.                             | Yes     | Removing the empty optional property is the required structural repair.    |

### Not checked

- **Rating-level order.** The spec requires levels ordered best-to-worst, but
  that ordering is semantic and cannot be verified mechanically, so no rule
  enforces it (fails [Rule scope](#rule-scope) criterion 2).
- **Body heading.** The body's top-level heading SHOULD name the model's subject
  (matching `title` when set), but whether a heading *names* a subject is a
  semantic judgment rather than a string match, so no rule enforces it (fails
  [Rule scope](#rule-scope) criterion 2).
- **`info` rules.** The initial set defines none; the severity is reserved for
  future non-judgmental observations.
