---
type: Change Case
title: Area-local Factor References
description: Make Requirement factor references resolve only within their declaring Area and warn on ambiguous same-Area Factor names.
status: Done
tags: [format, lint, skill, authoring]
timestamp: 2026-06-27T00:00:00Z
---

# Area-local Factor References

A **Change Case** capturing the _why_ and _status_; the detail lives in its
children:

- [Functional spec](0141-area-local-factor-references/spec.md) - what the case
  must do.
- [Design doc](0141-area-local-factor-references/design.md) - how it's built,
  and why.

## Motivation

Requirement `factors` currently resolve against the declaring Area and its
ancestors. That lets a descendant Requirement contribute directly to an ancestor
Factor rating, which makes Factor ownership non-local: a root Factor can be held
up or down by Requirements scattered through descendant Areas. This complicates
authoring, linting, reporting, and agent planning, and it encourages authors to
use ancestor Factor references as a wildcard roll-up mechanism.

Factor references should stay local to the Area that owns the Requirement. A root
Area can still judge cross-cutting concerns with its own Requirements, and child
Areas can still declare same-named local Factors when the concern needs local
assessment. The model remains tree-shaped while preserving qualified Factor IDs
for exact addressing.

## Scope

Covered:

- Change `Requirement.factors` resolution from Area-or-ancestor to declaring
  Area only.
- Keep same-named Factors in different Areas valid and distinct.
- Add a lint warning when the same Factor name appears more than once inside one
  Area's recursive Factor tree, because scalar `factors` entries become
  ambiguous to authors.
- Align format semantics, lint specs, lint implementation, tests, runtime skill
  authoring guidance, durable skill specs, scaffold comments, shared UX example
  guidance, logs, and release notes.

Deferred:

- Introducing path-qualified `Requirement.factors` entries.
- Migration tooling for existing models.
- New evaluation record payload fields.
- Numeric roll-up or weighting changes.

## Affected artifacts

Derived by sweeping for Factor references, ancestor resolution, secondary
Factors, same-name Factors, descendant roll-up wording, and `unknown-factor`.

**Code**

- [x] `internal/lint/model.go` - make factor-reference resolution Area-local and
      support same-Area duplicate-name detection.
- [x] `internal/lint/rules.go` - update `unknown-factor`, `empty-factor`, and
      add the duplicate-name warning check.
- [x] `internal/lint/result.go` - add the new warning rule.
- [x] `internal/lint/rules_test.go` - cover local, ancestor, sibling,
      cross-Area duplicate, and same-Area duplicate cases.

**Durable specs** (substance in the [functional spec](0141-area-local-factor-references/spec.md))

- [x] `SPECIFICATION.md` - make explicit Factor references Area-local, remove
      descendant-Area contribution semantics, and add same-Area Factor-name
      warning guidance.
- [x] `specs/cli/lint-rules.md` - update `unknown-factor` and add the new
      duplicate Factor-name warning rule.
- [x] `specs/skills/quality-skill/quality-skill.md` - keep scoped Factor
      evaluation tied to local Factor connections.
- [x] `specs/skills/quality-skill/guides/authoring/factors.md` - require
      same-Area duplicate-name warning guidance.
- [x] `specs/skills/quality-skill/guides/authoring/requirements.md` - require
      Area-local Factor-reference guidance.

**Durable docs / bundled skill runtime**

- [x] `README.md` - clarify direct Area Requirement factor references.
- [x] `docs/guides/agent-mediated-ux.md` - remove the ancestor roll-up example.
- [x] `docs/log.md` - record the shared guide update.
- [x] `internal/scaffold/skeleton.md` - clarify local Area factor references in
      scaffold comments.
- [x] `skills/quality/guides/authoring/factors.md` - guide same-named child
      Factors and duplicate-name warning behavior.
- [x] `skills/quality/guides/authoring/requirements.md` - guide Area-local
      `factors` references.
- [x] `skills/quality/guides/log.md` - append runtime-guide history entry.
- [x] `skills/quality/log.md` - append runtime skill history entry when needed.
- [x] `specs/log.md` - append durable-spec history entry.
- [x] `specs/skills/quality-skill/guides/log.md` - append durable-guide history
      entry.
- [x] `CHANGELOG.md` - note the format/lint and skill guidance change.

No planned impact: evaluation JSON schema, `qualitymd model` canonical
references, install docs, or release tooling.

## Status

`Done`. Implemented and archived after lint, format semantics, durable specs,
runtime skill guidance, docs, scaffold comments, logs, and release notes were
updated. `go test ./...`, `go run ./cmd/qualitymd lint QUALITY.md --json`,
`mise run fmt-md-check`, and `mise run check` pass.
