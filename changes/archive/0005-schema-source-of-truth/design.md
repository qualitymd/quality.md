---
type: Design Doc
title: Single source of truth for the structural schema — design doc
description: How the linter derives structural checks from one typed schema declaration.
tags: [lint, spec, schema, design]
timestamp: 2026-06-17T00:00:00Z
---

# Single source of truth for the structural schema — design doc

Design behind the
[Single source of truth for the structural schema](../0005-schema-source-of-truth.md)
change and its [functional spec](spec.md).

## Context

The linter used to encode structural format facts directly in rule code: each
node checker had its own valid-key map and shape switch. That duplicated the
public format spec and made drift easy: a documented key could be forgotten by
lint, or lint could keep rejecting a key the spec allows.

## Approach

Add `internal/schema` as the owner of the structural schema concept. It defines
node kinds, valid properties, YAML shapes, presence level, collection element
shape/kind, the model-content required-any group, and the rating-scale minimum.

`internal/lint` now imports that declaration and derives structural checks from
it:

- unknown-key checks call `schema.Node.Property`;
- scalar/map/sequence shape checks read `schema.Property.Shape`;
- secondary-factor list items and rating override values read
  `ElementShape`/`ValueShape`;
- rating-scale minimum reads `MinItems`; and
- the empty-model check reads the model-content `RequiredAny` group.

Semantic lint remains in `internal/lint`: duplicate rating levels, unknown
secondary factor references, unknown rating override keys, empty factor/target
warnings, and rule-specific messages still belong to lint. The schema says what
the structure is; lint decides which finding to emit and how to phrase it.

Spec consistency is enforced by `internal/schema` tests that parse the YAML
schema snippets in `SPECIFICATION.md` and compare their keys and presence labels
to the typed declaration. The test also pins the documented rating-scale minimum
and model-content group. This keeps the public format spec and implementation
source aligned without generating docs.

## Alternatives

- **Embedded data file.** Rejected for now. A YAML/JSON schema file would add
  parsing and validation machinery with only one implementation consumer.
- **Generate `SPECIFICATION.md` from the schema.** Rejected as too much
  machinery for this change. The public prose remains authored prose; tests catch
  drift.
- **Use Go model structs and YAML tags as the schema.** Rejected. Struct tags
  capture key names but not presence, recommendations, collection item shapes,
  the model-content required-any rule, or rating-scale minimum.
- **Keep the schema inside `internal/lint`.** Rejected. The schema models the
  format structure, not a lint implementation detail; a neutral package makes
  that ownership clear.

## Trade-offs & risks

- The schema declaration is still an implementation artifact, not a public
  generated schema format. A second consumer may justify a data-file format or
  generated docs later.
- Rule-specific messages still contain human wording and key names. The structural
  facts come from the schema; wording remains lint's job.
- The spec-consistency test checks the YAML schema snippets and selected prose
  anchors, not every explanatory paragraph. It is a targeted drift guard, not a
  documentation generator.
