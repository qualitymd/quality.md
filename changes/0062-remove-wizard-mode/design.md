---
type: Design Doc
title: Remove wizard mode - design doc
description: How the /quality skill removes wizard from its public contract while preserving safe orientation.
tags: [skill, ux, contract]
timestamp: 2026-06-23T00:00:00Z
---

# Remove wizard mode - design doc

## Context

This design answers the
[Remove wizard mode functional spec](spec.md). The change removes `wizard` as a
public `/quality` mode without replacing it with a public `status`, `next`, or
review command. Bare and ambiguous `/quality` requests stay safe and read-only.

## Approach

Treat the change as a public-contract reduction, not a rename.

Keep the runtime skill root as the only place that describes ambiguous request
handling. `SKILL.md` should list `setup`, `evaluate`, and `update` as the public
modes, then separately state that bare or unclear `/quality` input produces
read-only orientation. That orientation may use `qualitymd status --json` and
the bounded top-10 model checklist, but it must not identify itself as a mode or
advertise a command name.

Remove the wizard mode files rather than renaming them:

```text
skills/quality/modes/wizard.md
specs/skills/quality-skill/modes/wizard.md
```

The durable parent skill spec should absorb the durable part of the old wizard
contract as shared routing/orientation language: read-only, shallow, status-first
inspection for bare or ambiguous input, with next actions limited to public
workflows. This avoids introducing an internal mode document whose name could
become a new accidental contract.

Update setup and getting-started handoffs to recommend concrete public next
workflows directly. For example, after guided setup, recommend evaluation when
the model is ready, model repair/improvement when it is not, or update when
tooling is stale. Do not tell users to run another read-only command.

Keep `qualitymd status` as CLI support tooling. The skill may still consume it
internally, and docs may mention it when describing CLI prerequisites or support
commands, but `/quality status` should not appear as a public skill invocation.

For `/quality wizard`, prefer a narrow compatibility note in runtime guidance:
if a user sends it, respond read-only that `wizard` is deprecated/removed from
the public surface and point to `setup`, `evaluate`, `update`, or recommendation
follow-up. Do not list it in public docs or invocation examples.

## Alternatives

**Rename wizard to status.** Rejected because the user explicitly does not want
`status` or `next` to become public contract. A rename would preserve the same
surface area under a different noun.

**Keep wizard as an internal mode file.** Rejected because internal mode files
are easy to expose in run frames, docs, and examples. The durable behavior is
small enough to live in the parent skill routing contract.

**Make bare `/quality` default to evaluate.** Rejected because evaluation writes
artifacts. The safe default is an important interaction invariant.

**Remove read-only orientation entirely.** Rejected because ambiguous requests
still need a safe path that helps the user choose a public workflow without
mutating the workspace.

## Trade-offs & Risks

Removing a documented invocation is a public-contract change. A compatibility
response for `/quality wizard` reduces surprise for users who learned the old
surface, but it should be framed as deprecated behavior and kept out of public
examples.

Moving orientation language into the parent skill spec makes that spec carry a
little more shared behavior. That is acceptable here because the behavior is a
cross-cutting safety rule for ambiguous input, not an independently public
workflow.

Search cleanup needs to distinguish live contract from history. Archived Change
Cases and append-only logs can retain old `wizard` references; runtime files,
durable specs, current docs, README, and changelog should not advertise it as
current public surface.

## Open Questions

- Should `/quality wizard` compatibility behavior remain for one release, or be
  removed from runtime guidance immediately?
