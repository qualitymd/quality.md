---
type: Design Doc
title: Area-local Factor References - design
description: How lint resolution, warnings, specs, and skill guidance adopt Area-local Factor references.
tags: [format, lint, skill, authoring]
timestamp: 2026-06-27T00:00:00Z
---

# Area-local Factor References - design

## Context

Answers the [functional spec](spec.md) for change case
[0141](../0141-area-local-factor-references.md). The change simplifies
Requirement-to-Factor resolution by making explicit `factors` entries local to
the Requirement's declaring Area. It also adds an advisory warning for duplicate
Factor names inside one Area's recursive Factor tree, because scalar references
cannot distinguish those Factors.

## Approach

### Make the resolver Area-local

Replace the lint resolver's ancestor walk with a direct search of the declaring
Area's Factor tree:

```go
func (s *runState) resolveFactor(area *areaRef, name string) *factorRef {
    return findFactor(area.factors, name)
}
```

This single resolver is used by both `unknown-factor` validation and
`empty-factor` analysis. Keeping one path prevents a Requirement reference from
being invalid for one rule while still counting toward another rule's coverage.

### Warn on duplicate names inside one Area tree

Add a warning rule, `duplicate-factor-name`, to the lint rule catalog. The rule
walks each Area's Factor tree, groups Factors by scalar name, and emits a warning
on every Factor participating in a duplicate group.

The rule is scoped per Area:

- same name under different Areas: valid, no warning;
- same name under one Area at different paths: warning;
- duplicate YAML keys under one `factors` map: existing YAML traversal can expose
  both key nodes, so the same grouped check reports them when present.

The warning is not fixable. Choosing whether to rename, merge, or split Factors
is authored model judgment.

### Update lint copy and tests

`unknown-factor` wording changes from "declaring area or an ancestor" to
"declaring area." Tests invert the current ancestor-resolution fixture so a
child Area referencing a root-only Factor now errors. Add tests proving local
shadowing works when both root and child declare the same name.

Add focused duplicate-name tests:

- duplicate names under one Area's Factor tree warn;
- duplicate names across different Areas do not warn;
- duplicate names do not affect structural validity unless other errors exist.

### Update durable semantics and guidance

Update `SPECIFICATION.md` so it no longer describes descendant Area Requirements
as direct contributors to ancestor Factor ratings. Root-level cross-cutting
judgment should be modeled as root-level Requirements that assess the root
source and child coverage, not as descendant Requirements tagging an ancestor
Factor.

Update runtime and durable authoring guides so they teach:

- `factors` entries name same-Area Factors only;
- same-named child Area Factors are distinct local refinements;
- repeating a Factor name inside one Area's Factor tree is discouraged because it
  makes scalar references hard to read.

Update the shared UX example that previously planned to connect descendant
Security Requirements to a root Factor; the new example should add root Security
coverage without descendant opt-in references.

### Logs and release notes

Record the format/lint change in `specs/log.md`, the runtime guide updates in
`skills/quality/guides/log.md`, the skill guidance update in
`skills/quality/log.md`, the shared guide update in `docs/log.md`, and the
user-facing change in `CHANGELOG.md`.

## Spec response

- **Format semantics** - satisfied by editing `SPECIFICATION.md` and removing
  descendant explicit-reference roll-up language.
- **Lint behavior** - satisfied by the Area-local resolver, new warning rule,
  updated messages, and focused tests.
- **Skill and authoring guidance** - satisfied by runtime guide edits plus
  durable guide spec edits.
- **Verification** - satisfied by focused lint tests, broad Go tests, and
  Markdown formatting checks.

## Alternatives

- **Keep ancestor references and only discourage them in guidance.** Rejected.
  The complexity is in the formal contract and linter, so guidance alone would
  leave agents with two competing models.
- **Make duplicate Factor names inside one Area an error.** Rejected. Canonical
  Factor IDs include the full path, so these Factors remain structurally
  addressable; the problem is authoring clarity for scalar `factors` entries.
- **Introduce path-qualified `factors` entries now.** Rejected. That is a larger
  schema change and not required to remove ancestor roll-up complexity.
- **Warn whenever an Area shadows an ancestor Factor name.** Rejected. Same-name
  Factors on different Areas are a useful local refinement pattern and already
  have distinct qualified IDs.

## Trade-offs & risks

- Existing models that rely on ancestor references will become invalid. This is
  acceptable under the early-alpha clean-break policy, and the repair is explicit:
  add a same-named local Factor or move the Requirement to the Area that owns the
  Factor.
- Root cross-cutting Factors no longer get automatic descendant Requirement
  contributors. Authors must model root coverage with root Requirements, which is
  more explicit but sometimes more verbose.
- Duplicate-name warnings may surface in hand-authored models that are otherwise
  valid. The warning is advisory and does not block use.

## Open questions

None.
