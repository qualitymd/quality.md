---
type: Functional Specification
title: qualitymd models
description: Requirements for listing and viewing bundled QUALITY.md models.
tags: [cli, models, specification]
timestamp: 2026-06-17T00:00:00Z
---

# qualitymd models

`qualitymd models` emits bundled `QUALITY.md` models that agents and tools use as
deterministic inputs. It does not inspect a user's local model, record evaluation
state, or perform judgment.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Scope

The command exposes a bundled-model catalog:

- `qualitymd models list`
- `qualitymd models view <name>`

The initial catalog **MUST** include `quality-meta-model`, the model used to
evaluate a `QUALITY.md` itself. Future bundled sample or diagnostic models
**MAY** be added to the same catalog.

Commands that inspect a user's local `QUALITY.md` are outside this namespace and
SHOULD be top-level commands, such as a future `qualitymd outline`.

## `models list`

`qualitymd models list` **MUST** emit the bundled catalog in deterministic order.

Default output is a human-readable table with each model's `name`, `title`, and
`description`. When stdout is a terminal and color is enabled, the command
**SHOULD** style the table with the shared CLI palette while preserving the same
facts. When stdout is redirected, piped, or `NO_COLOR` is set, stdout **MUST**
remain the plain table with no terminal control sequences.

With `--json`, stdout **MUST** be a stable JSON array:

```json
[
  {
    "name": "quality-meta-model",
    "title": "Quality meta-model",
    "description": "Criteria for evaluating a QUALITY.md model."
  }
]
```

## `models view <name>`

`qualitymd models view <name>` **MUST** emit the named bundled model.

Default output is Markdown. When stdout is a terminal and color is enabled, the
command **MAY** render the Markdown through the terminal renderer and pager using
the same rules as [`qualitymd spec`](spec.md). When stdout is redirected, piped,
or `NO_COLOR` is set, stdout **MUST** be plain Markdown. Without `--source`, that
plain Markdown **MUST** be the authored bundled model bytes.

With `--json`, stdout **MUST** be a stable JSON document:

```json
{
  "schemaVersion": 1,
  "name": "quality-meta-model",
  "title": "Quality meta-model",
  "description": "Criteria for evaluating a QUALITY.md model.",
  "model": {
    "title": "Quality meta-model",
    "source": "QUALITY.md",
    "ratingScale": []
  },
  "bodyMarkdown": "\n# Quality meta model\n..."
}
```

`model` is the parsed frontmatter using the same `QUALITY.md` model shape as the
file format. `bodyMarkdown` is the model's Markdown body after the closing
frontmatter fence, preserving its leading newline when present.

## Source rewrite

`models view <name>` **MUST** accept `--source <path>`. The flag rewrites the
emitted model's apex `source` to the supplied path before either Markdown or JSON
is emitted. The rewrite **MUST** affect only the root `source` node.

The flag exists so model-altitude evaluation can use the bundled
`quality-meta-model` as the active model while pointing its apex source at the
user's `QUALITY.md`.

## Errors

An unknown bundled model name **MUST** fail with a usage error. The command
**MUST** otherwise follow the parent CLI baseline for stdout/stderr separation,
non-interactivity, deterministic output, and exit-code categories.
