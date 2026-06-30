---
type: Functional Specification
title: Single source of truth for the structural schema — functional spec
description: Requirements for a single authoritative QUALITY.md structural schema that the linter derives from.
tags: [lint, spec, schema]
timestamp: 2026-06-17T00:00:00Z
---

# Single source of truth for the structural schema — functional spec

Companion to the
[Single source of truth for the structural schema](../0005-schema-source-of-truth.md)
change. This spec states _what_ the change must do; a design doc covers _how_. The
format itself is defined in [`SPECIFICATION.md`](../../../SPECIFICATION.md); this spec
does not restate it.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be interpreted
as described in IETF RFC 2119.

## Scope

Covers a single authoritative definition of the QUALITY.md _structural schema_ and
the linter deriving its structural validation from it. **Deferred:** generating
prose specs or user docs from the source; any runtime configuration surface;
unknown-key typo suggestions; and changes to the rule catalog's severities or
membership. See the [change's scope](../0005-schema-source-of-truth.md#scope).

## Requirements

- There **MUST** be a single authoritative definition of the structural schema:
  the valid frontmatter keys for each node kind (root, target, factor,
  requirement), which keys are required versus optional, and the rating-scale
  shape — enough to drive the linter's structural rules.
- The linter **MUST** derive its structural validation — unknown-key, node-shape,
  and required-property checks — from this definition. A second, independently
  hand-maintained copy of the valid-key set **MUST NOT** exist.
- Adding, removing, or renaming a valid key **MUST** require editing only the
  schema definition, not the rule logic.
- The durable format spec ([`SPECIFICATION.md`](../../../SPECIFICATION.md) and the
  relevant [`specs/`](../../../specs/index.md) concepts) **MUST** be reconciled to
  agree with the definition, and the project **SHOULD** make that consistency
  mechanically checkable rather than maintained by hand.
- The change **MUST** be behavior-preserving: for any document, the linter's
  findings **MUST NOT** change — except where the documented format and the
  current rule code disagree today, which **MUST** be reconciled in favor of the
  documented format.
