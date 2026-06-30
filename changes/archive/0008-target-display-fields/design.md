---
type: Design Doc
title: target title and description — design doc
description: How the title/description target fields and the Model = Target + ratingScale framing land in the schema, the linter, and the spec.
tags: [specification, schema, lint, design]
timestamp: 2026-06-17T00:00:00Z
---

# target title and description — design doc

Design behind the
[Describe targets with title and description](../0008-target-display-fields.md)
change and its [functional spec](spec.md).

## Context

A target's only label is its map key, and the spec frames the root as an _apex
target_ with a "non-root target MUST NOT declare `title`/`ratingScale`"
prohibition. The change adds `title` and `description` to every target and
reframes the root as a **Model** = the model-wide `ratingScale` plus the
**Target** properties. The structural schema already models `Model` and `Target`
as two distinct `schema.Node`s, so the framing is mostly catching the prose up to
the code. The interesting question is how little code actually has to move.

## Approach

Three property additions in [`internal/schema`](../../../internal/schema/schema.go)
drive everything downstream:

- `Target` gains `title` (`ScalarShape`, `RecommendedPresence`) and `description`
  (`ScalarShape`, `OptionalPresence`).
- `Model` gains `description` (`ScalarShape`, `OptionalPresence`). Its `title`
  stays `RecommendedPresence`.

`SPECIFICATION.md` absorbs the matching prose: the Target YAML snippet lists
`title` (`# Recommended`) and `description` (`# Optional`); the Model snippet adds
`description` (`# Optional`); and the Target intro plus the "shares the
structure … but for two keys" sentence are rewritten as _Model = Target
properties + `ratingScale`_, with `ratingScale` named the one Model-only
property. The evaluation prose's "root target" wording (the top of the target
tree) stays — it remains true under the new framing.

Two consequences fall out without changing rule logic:

- **`misplaced-root-key` self-narrows.** The rule
  ([`rules.go`](../../../internal/lint/rules.go)) fires in `checkSchemaProperties`
  when a key is _not_ a property of the current node but _is_ a property of
  `Model` — i.e. a Model-only key on a target. Once `title`/`description` are
  Target properties, the only key in `Model` and not in `Target` is
  `ratingScale`, so the branch fires for `ratingScale` alone. The diagnostic
  already interpolates the key name, so it reads "…declares `ratingScale`;
  `ratingScale` is only valid on the model root." The rule catalog description
  (`result.go`) still changes so the static rule text names `ratingScale` rather
  than a generic root-only key.
- **The consistency test needs no edit.** `schema_test.TestSpecificationSchemaSnippetsMatchDeclaration`
  is data-driven: it parses each `#### Heading` YAML snippet into a key→presence
  map and compares it to the `Node` declaration in both directions. Adding the
  properties to the declarations and the snippets in lockstep keeps it green; the
  test code is untouched.

So the only test change is in
[`rules_test.go`](../../../internal/lint/rules_test.go): the "nested target
title" case currently _expects_ `misplaced-root-key` and must flip to a valid
model (no findings), joined by a case covering `description` on a nested target.
The "nested target rating scale" case stays as the surviving root-only assertion,
with a second `ratingScale` fixture kept for the catalog's fixture-count guard.
Lastly, [`specs/cli/lint.md`](../../../specs/cli/lint.md)'s `misplaced-root-key`
row and the rule catalog description change from root-only-key wording to
`ratingScale`.

## Alternatives

- **Compose `Model.Properties` from `Target.Properties` + `ratingScale`.**
  Tempting — it would encode "Model = Target + `ratingScale`" literally and kill
  the duplicated lists. With `title` and `description` now sharing presence across
  the two nodes, composition is closer than before. Rejected anyway: the Model
  carries the `model-content` `RequiredAny` group a Target lacks, and inserts
  `ratingScale` _between_ `description` and `factors` rather than appending it, so
  composition would still need a positional splice plus the extra `RequiredAny` —
  more code, and less obvious, than the two explicit lists it replaces. The drift
  composition guards against is already caught by the consistency test, so
  explicit lists win.
- **Add a `missing-target-title` warning** to mirror `missing-title` and
  `missing-factor-description`, enforcing the RECOMMENDED `title`. Rejected for
  this change: "missing" warnings are explicit rules, not derived from
  `RecommendedPresence`, so omitting the rule means a target without a `title` is
  simply not nagged (it still has its map key). That keeps the change to
  fields-and-framing. A follow-up may add the warning if the omission proves
  costly.
- **Keep the "apex target / MUST NOT" prose and only add the fields.** Rejected:
  once `title` and `description` are shared, `ratingScale` is the lone
  difference, and stating it as a prohibition rather than a type fact is exactly
  the awkwardness this change set out to remove.

## Trade-offs & risks

- `description` is `Optional` on both the Model and a Target, so the field reads
  uniformly across the tree: neither the root nor a nested target is nagged for
  omitting it. The two nodes still are not byte-identical minus `ratingScale` —
  the Model carries the `model-content` `RequiredAny` group a Target lacks — but
  that is a structural difference, not field-presence drift.
- `Target.title` is `Recommended` but unenforced (no warning rule). A reader
  comparing the schema's `RecommendedPresence` against lint behavior could expect
  a warning; the absence is intentional and documented here and in the spec.
- The reframing rewrites prose the [schema-source](../0005-schema-source-of-truth.md)
  consistency test only spot-checks (snippets and selected anchors), so the
  narrative paragraphs are reviewed by hand. The structural keys/presence remain
  test-guarded.

## Open questions

- None blocking. The `missing-target-description` warning is a possible follow-up,
  recorded above rather than carried here.
