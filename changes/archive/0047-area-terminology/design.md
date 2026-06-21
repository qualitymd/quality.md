---
type: Design Doc
title: Area terminology changeover - design doc
description: Implementation approach for replacing Target and Subject-facing vocabulary with Area and root area.
tags: [terminology, schema, evaluation, cli, skill]
timestamp: 2026-06-21T00:00:00Z
---

# Area terminology changeover - design doc

Design behind the [Area terminology changeover](../0047-area-terminology.md)
and its [functional spec](spec.md).

## Context

This is a draft-format breaking change. The implementation should not carry
Target compatibility aliases, and it should also retire Subject where Subject was
only a user-facing name for the evaluated root. The one concept is Area; the root
node is the root area.

The change crosses almost every layer because Target is both a YAML property
(`targets:`) and a Go/API vocabulary (`TargetPath`, `targetPath`, target
summaries, target reports). The safest implementation is a mechanical big-bang
rename with explicit guardrails around the places where an empty root path can
otherwise mask a missing field.

## Approach

Start at the model/schema boundary, then let compiler and tests drive the rest.

Rename `internal/model.Target` to `Area`, `Spec.Targets` to `Spec.Areas`, and the
recursive YAML/JSON tags from `targets` to `areas`. In `internal/schema`, rename
`TargetKind` to `AreaKind`, `PropertyTargets` to `PropertyAreas`, and the node
definition from `Target` to `Area`. Removing `targets` from the allowed property
set is enough for lint to reject legacy files as non-conforming instead of
treating them as a second accepted shape.

Leave the existing `source` property name in place. The implementation should
only update comments, docs, diagnostics, and report prose so `source` is
described as the material selected for an Area. Avoid substituting broader words
such as material, content, evidence, or scope into the schema.

Use one pass to rename traversal concepts:

- `walkTargets` -> `walkAreas`;
- `TargetPath` -> `AreaPath`;
- `TargetSummary` / `TargetDetails` report concepts -> Area equivalents;
- `ChildTargets` -> `ChildAreas`;
- human labels such as `Target Summary`, `Target Details`, and `Target Ratings`
  -> Area headings.

Keep the path identity model unchanged. An Area path is still an ordered list of
stable model keys, and the root area remains the empty path. The rename changes
the concept and serialized field names, not the path algorithm.

Record decoding needs one non-mechanical guardrail: required path fields must be
checked for field presence before trusting zero values. A missing `areaPath` and
an explicit root `areaPath: []` both decode to an empty slice in Go; old
`targetPath` records must not be accidentally interpreted as root-area records.
The record loader should validate raw JSON object keys, or use a small decode
shim, so `targetPath`-only records become incompatible historical records.

Evaluation run creation should replace `--subject` with `--model`, because the
flag selects the `QUALITY.md` model file to snapshot, not a second evaluated
concept. New run folders should drop the altitude word:

```text
0001-quality-eval
0002-cli-quality-eval
```

The run-number scanner can continue to look for the leading `NNNN-` sequence so
old run folders remain visible as historical folders. New report building should
only trust current-schema records.

Reports should replace `Subject` key details with `Root area` when naming the
evaluated root, and use direct model/file wording when naming the snapped
`QUALITY.md` file. Summary rendering should always include a row for the root
area when a root-area analysis exists, even with no child areas. That test should
use a root-only model with root-level factors so it exercises the fresh-project
case that prompted this change.

The `/quality` skill should change its parsing grammar after the CLI/schema
rename is in place:

- bare names resolve against Areas or Factors;
- two bare names are an `<area> <factor>` pair;
- explicit disambiguators become `area` and `factor`;
- run frames say `Model file` or `QUALITY.md file` for the file path and `Scope`
  for the Area/Factor selection;
- record examples use `areaPath`.

The dogfood `QUALITY.md`, scaffold, README, npm README, and maintained example
evaluation bundle should be updated in the same implementation pass so no live
artifact teaches `targets:` or `targetPath`.

## Alternatives

Display-only rename:

- Rejected. It would improve reports but leave `targets:` and `targetPath` in
  the authoring and record contracts, preserving the cognitive split this change
  is meant to remove.

Formal `Entity`, user-facing `Area`:

- Rejected. It is more standards-aligned, but it creates two canonical nouns.
  The chosen product direction values one term that users and agents can repeat
  consistently.

Accept both `targets:` and `areas:` for one release:

- Rejected. The format is still draft, and compatibility aliases would spread
  throughout linting, records, examples, skill instructions, and docs. The cost
  is not worth preserving a short-lived draft term.

Keep `--subject` and only rename Target:

- Rejected. Once root area is the formal root descriptor, Subject becomes a
  second user-facing noun for the evaluated thing. `--model` names what the flag
  actually selects: the `QUALITY.md` model file.

Rename the default rating level `target`:

- Rejected. The problematic overload is resolved by renaming the model node.
  The rating scale vocabulary remains useful and is independent of Area.

Rename `source`:

- Rejected. `source` remains the clearest schema key for a path, glob, URL, or
  selector that identifies what an Area evaluates. Alternatives either collide
  with existing concepts (`scope`, `evidence`) or read poorly as YAML fields
  (`material`, `corpus`, `basis`).

## Trade-offs & Risks

The main cost is churn. This change will touch most tests, examples, generated
reports, and skill guidance. The mitigation is to land it as one intentional
breaking pass rather than splitting machine fields from prose labels.

Historical evaluation runs become non-reportable under current readers. That is
acceptable for this draft-format change, but status/history views should explain
the incompatibility rather than failing opaquely.

Area can blur with Factor if authors use quality-dimension names as areas. The
spec and authoring guidance need clear contrast: Area is what is evaluated;
Factor is the quality lens. Examples should show concrete areas such as CLI,
README, Documentation, API, and Tests.

The root path is the highest-risk data bug. Tests must cover old `targetPath`
records, missing `areaPath`, and explicit root `areaPath: []` so compatibility
rejection does not break valid root-area records.

## Open Questions

None for the design. Exact wording polish can happen while updating the durable
docs, but the term choices and compatibility posture are settled by the
functional spec.
