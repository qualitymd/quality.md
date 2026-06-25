---
type: Functional Specification
title: Named Requirement identity — functional spec
description: Target behavior for stable Requirement names and qualified Requirement references.
tags: [format, requirements, model-references]
timestamp: 2026-06-25T00:00:00Z
---

# Named Requirement identity — functional spec

Companion to [Named Requirement identity](../0093-requirement-identity.md). This
spec states the target format and validation behavior for named Requirements.

The key words **MUST**, **SHOULD**, and **MAY** are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background

The current format identifies a Requirement by its map key, which is the
natural-language Requirement statement. Evaluation records need stable
Requirement references that survive wording edits and map cleanly into generated
data paths. Areas and Factors already separate stable names from display titles;
Requirements should follow the same pattern.

## Scope

This change covers the target Requirement frontmatter shape, Requirement naming
rules, qualified Requirement model references, and validation behavior.

Deferred:

- Automated migration tooling or authoring assistance for existing QUALITY.md
  files.
- Any evaluation v2 record schema that consumes Requirement references.

## Requirements

### Named Requirement shape

A Requirement entry **MUST** use a stable Requirement name as the map key.

The Requirement name **MUST** match the same strict name grammar used by Area
names, Factor names, and Rating Level IDs:

```regex
^[A-Za-z0-9](?:[A-Za-z0-9_-]*[A-Za-z0-9])?$
```

A Requirement object **MUST** include `title` as a non-empty scalar.

A Requirement object **MUST** include `assessment` as a non-empty scalar.

`assessment` continues to name or state the basis and means for assessing the
Requirement. It may be inline prose or a reference to a specification, guide,
checklist, standard, or other artifact that defines how the Requirement should
be assessed.

A Requirement object **MAY** include `description` as a non-empty scalar for a
human-facing gloss or explanation.

A Requirement object **MAY** include `ratings` with the same semantics as the
current format: a map keyed by Rating Level IDs from the Model's Rating Scale,
with non-empty scalar criteria.

Current statement-key Requirement entries are not a compatibility input in this
format version. A legacy Requirement whose map key is the natural-language
statement fails the Requirement name grammar, and a Requirement without `title`
fails the required-title rule.

### Factor connections

Every Requirement **MUST** remain connected to at least one Factor.

A Requirement declared under a Factor or sub-factor is connected by placement;
the containing Factor remains its primary Factor.

A Requirement declared under a Factor or sub-factor **MAY** include `factors` for
secondary Factor references.

A Requirement declared directly under an Area **MUST** include `factors` with at
least one non-empty scalar entry.

Each explicit Factor reference **MUST** resolve to a Factor in scope on the
declaring Area or one of its ancestors.

### Requirement uniqueness

Requirement names **MUST** be unique within the declaring Area.

For uniqueness, all Requirements declared directly under an Area and all
Requirements declared under that Area's Factors or sub-factors are considered to
belong to the declaring Area.

### Requirement model references

Delimiter note: this draft uses `::` to stay aligned with the current Factor
reference shape, `factor:<declaring-area-path>::<factor-path>`. The delimiter
can still be revised before this case moves to `In-Progress` if the broader
model-reference grammar changes consistently.

Qualified Requirement references **MUST** use:

```text
requirement:<declaring-area-path>::<requirement-name>
```

The root declaring Area is written as `root`, for example:

```text
requirement:root::release-notes-current
requirement:webhooks/delivery::retry-window
```

Tools that render qualified Requirement references **MUST** use the
`requirement:` prefix.

Tools that parse qualified Requirement references **MUST** reject references
whose Area path or Requirement name fails the strict name grammar or whose
referenced Requirement does not exist.

Tools **MUST NOT** persist unqualified Requirement references in durable
machine-readable artifacts.

## Durable spec changes

### To add

None

### To modify

- `SPECIFICATION.md` - replace statement-key Requirement identity with named
  Requirement entries, add `title` and optional `description`, keep `assessment`,
  preserve Factor connection semantics, and add qualified Requirement references
  per the requirements above.
- `specs/quality-schema-json.md` - update schema expectations for named
  Requirements.
- `specs/cli/lint-rules.md` - update validation rules for Requirement names,
  titles, uniqueness, Factor references, ratings, and compatibility behavior.
- `specs/cli/init.md` - update starter examples when they include Requirements.

### To rename

None

### To delete

None
