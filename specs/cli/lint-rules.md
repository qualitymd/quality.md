---
type: Functional Specification
title: qualitymd lint rules
description: Rule-system, rule-authoring, and rule-catalog contract for qualitymd lint.
tags: [cli, command, lint, rules]
timestamp: 2026-06-22T00:00:00Z
---

# qualitymd lint rules

This spec owns the rule-system and rule-catalog contract for
[`qualitymd lint`](lint.md): which checks belong in lint, how rule identifiers,
descriptions, messages, severity, and fixability work, and the initial rule set.
The command invocation, flags, repair execution, exit status, and ordering rules
live in [`qualitymd lint`](lint.md). The finding/output schema lives in
[`qualitymd lint output`](lint-output.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", "RECOMMENDED", and "MAY" are to be
interpreted as described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Rule scope

A QUALITY.md is validated by a set of named **rules**, each checking one
property and emitting a **finding** when violated. A rule belongs in `lint` only
when **all four** of the following hold:

1. **Conformance-grounded.** It enforces a constraint that already exists in the
   [format specification](../../SPECIFICATION.md). `lint` introduces no
   requirements of its own; each rule traces to the spec clause it checks.
2. **Deterministic and mechanical.** Its verdict is computable from the file's
   structure alone — no judgment, no external data, the same result every run.
3. **Format, not goodness.** It decides whether the file *is a valid
   QUALITY.md file*, never whether it is a *good* quality model. Anything requiring
   evaluative judgment is out of scope — that is the evaluation skills' work.
4. **Self-contained.** It is decidable from the single file plus its own declared
   cross-references (e.g. resolving a factor reference, a `ratings`
   override key). A rule MUST NOT resolve a area's `source` or read the
   entities an area evaluates; that is beyond format conformance.

Structural validation — valid keys, required/recommended/optional properties,
YAML value shapes, the model-content group, and the rating-scale minimum — MUST
derive from the single structural schema declaration that the linter consumes.
Rule logic MUST NOT maintain a second valid-key list independent of that schema.

The unknown-key rule **MUST** be internally configurable by rule options so a
documented qualitymd tooling key can be allowed without hard-coding a one-off
exception into schema traversal. The default qualitymd lint profile **MUST**
allow root `config` and **MUST** keep all other unknown root keys as error
findings. Unknown nested keys inside Areas, Factors, Requirements, and Rating
Levels **MUST** remain error findings by default.

Root `config` is a qualitymd tooling convention, not a normative Model property.
`lint` **MUST** accept it only at the root and **MUST** validate its value with
the `invalid-config` rule.

### Severity

Every rule carries a fixed **severity** that governs how its finding surfaces:

- **error** — a spec **MUST**/**MUST NOT** is violated; the file is not a valid
  QUALITY.md file. Errors are what make `lint` exit non-zero.
- **warning** — a spec **SHOULD**/**RECOMMENDED** that is still mechanically
  determinable is unmet (e.g. a recommended `description` is absent, or an area
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
- Use one of two shapes: **`<defect>-<concept>`** for the defect a rule rejects
  (`missing-rating-scale`, `duplicate-level`, `unknown-factor`) or
  **`<concept>-<aspect>`** for a neutral observation about a concept (e.g. a
  future `model-summary` info rule). Every rule in the current set is
  `<defect>-<concept>`; the second shape is reserved for the `info` observations
  defined later.
- Name the concept in QUALITY.md vocabulary — Area, Factor, Requirement,
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
- **state the triggering condition**, leading with the concept in QUALITY.md
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
  description, the same for the same input, in QUALITY.md vocabulary with no
  implementation terms; and
- **carry tone to match severity** — errors state the violation plainly,
  warnings frame the recommendation they enforce, info is neutral observation; no
  blame, no exclamation.

A message should **be actionable when the fix is determinate** — pointing to the
expected shape or the valid set — and should stay concise, leading with the
problem. A rule that finds several instances should emit **one message per
instance**, each with its own location, rather than one bundled list, so each is
independently addressable.

### Fixability

A rule is **fixable** only when `qualitymd` can repair each finding with a
deterministic edit that does not presume the author's intent. Fixability is
independent of severity: an error can be fixable, and a warning can be
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

| Rule                       | Enforces (format spec)                                                                    | Description                                                                                      | Fixable | Fixable rationale                                                                  |
| -------------------------- | ----------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------ | ------- | ---------------------------------------------------------------------------------- |
| `invalid-frontmatter`      | *YAML Frontmatter* — a file MUST begin with a valid YAML frontmatter block                | Rejects frontmatter that is absent, not valid YAML, or not the model's shape.                    | No      | Parse and shape failures do not have one safe structural repair.                   |
| `missing-rating-scale`     | *Model* — `ratingScale` is required                                                       | Flags a model that declares no `ratingScale`.                                                    | No      | The scale defines model semantics; the format has no single required scale.        |
| `missing-title`            | *Model / Area / Factor / Rating Scale* — each displayable element requires a `title`      | Flags a Model, Area, Factor, or Rating Level with no non-empty scalar `title`.                   | No      | A title names the element for readers and cannot be inferred mechanically.         |
| `too-few-levels`           | *Rating Scale* — at least two rating levels MUST be supplied                              | Flags a `ratingScale` with fewer than two rating levels.                                         | No      | Adding levels requires choosing scale vocabulary and criteria.                     |
| `missing-level-name`       | *Rating Scale* — each level MUST declare a `level` name                                   | Flags a rating level that declares no `level` name.                                              | No      | A generated name would define rating vocabulary for the author.                    |
| `duplicate-level`          | *Rating Scale* — a `level` name MUST be unique within the scale                           | Flags two rating levels that share a `level` name.                                               | No      | Repair requires choosing which level to rename, merge, or remove.                  |
| `invalid-config`           | *qualitymd tooling* — root `config` MUST be a safe repository-relative scalar path        | Flags a root `config` value that is empty, non-scalar, absolute, or escapes the repository.      | No      | Repair requires choosing the intended workspace config file path.                  |
| `missing-criterion`        | *Rating Scale* — each level MUST declare a `criterion`                                    | Flags a rating level that declares no `criterion`.                                               | No      | Criterion text defines rating semantics and cannot be inferred mechanically.       |
| `empty-model`              | *Model* — an entry on `factors`, `requirements`, or `areas` MUST be supplied              | Flags a model root that supplies no entry under `factors`, `requirements`, or `areas`.           | No      | Repair requires choosing what kind of model content to add.                        |
| `misplaced-root-key`       | *Area* — a Area MUST NOT declare `ratingScale`                                            | Flags an area that declares `ratingScale`.                                                       | No      | Repair requires deciding whether to remove, move, or reinterpret authored content. |
| `invalid-assessment`       | *Requirement* — a requirement MUST declare exactly one `assessment` as a non-empty scalar | Flags a requirement whose `assessment` is missing, empty, or a list rather than a single scalar. | No      | Repair requires choosing the requirement's assessment text.                        |
| `unknown-factor`           | *Requirement* — each factor reference MUST resolve to a factor in scope                   | Flags a requirement whose `factors` entry references no factor on its area or an ancestor.       | No      | Repair requires choosing the intended in-scope factor or adding a new one.         |
| `missing-factor-reference` | *Requirement* — every requirement MUST be connected to at least one factor                | Flags a direct area-level requirement with no non-empty scalar factor reference.                 | No      | Repair requires choosing the intended factor or moving the requirement under one.  |
| `unknown-rating-key`       | *Requirement* — each `ratings` override key MUST name a level of the rating scale         | Flags a `ratings` override key that names no level in the model's `ratingScale`.                 | No      | Repair requires choosing the intended rating level or changing the scale.          |

### Warnings

Each enforces a mechanically determinable **SHOULD**/**RECOMMENDED** — its
finding is advisory and does not affect the exit code.

| Rule                         | Enforces (format spec)                                                   | Description                                                                                                                        | Fixable | Fixable rationale                                                                             |
| ---------------------------- | ------------------------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------- | ------- | --------------------------------------------------------------------------------------------- |
| `missing-level-description`  | *Rating Scale* — a level `description` is RECOMMENDED                    | Warns when a rating level declares no `description`.                                                                               | No      | A description states the level's meaning and requires authored content.                       |
| `missing-factor-description` | *Factor* — a factor SHOULD declare a `description`                       | Warns when a factor declares no `description`.                                                                                     | No      | A description states what the factor means and requires authored content.                     |
| `empty-factor`               | *Factor* — a factor SHOULD lead to at least one requirement              | Warns when a factor leads to no `requirements` declared under it, referencing it under `factors`, or reached through a sub-factor. | No      | Repair requires adding or moving requirements, or adding factor references with model intent. |
| `empty-area`                 | *Area* — each area SHOULD lead to a requirement somewhere in its subtree | Warns when an area's subtree reaches no `requirements`.                                                                            | No      | Repair requires choosing area content or restructuring the area tree.                         |
| `empty-property`             | *YAML Frontmatter* — null or empty optional properties SHOULD be omitted | Warns when an optional property is present but null or empty instead of omitted.                                                   | Yes     | Removing the empty optional property is the required structural repair.                       |

### Not checked

- **Rating-level order.** The spec requires levels ordered best-to-worst, but
  that ordering is semantic and cannot be verified mechanically, so no rule
  enforces it (fails [Rule scope](#rule-scope) criterion 2).
- **Body heading.** The body's top-level heading should name the root area
  (matching `title` when set), but whether a heading names the root area is a
  semantic judgment rather than a string match, so no rule enforces it (fails
  [Rule scope](#rule-scope) criterion 2).
- **`info` rules.** The initial set defines none; the severity is reserved for
  future non-judgmental observations.
