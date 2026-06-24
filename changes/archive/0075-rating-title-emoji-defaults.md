---
type: Change Case
title: Rating title emoji defaults
description: Make emoji-prefixed Rating Level titles the default starter and setup display convention while keeping the formal format neutral.
status: Done
tags: [scaffold, skill, rating-scale, docs]
timestamp: 2026-06-24T00:00:00Z
---

# Rating title emoji defaults

A **Change Case** to make lightweight emoji markers the default for Rating Level
titles in starter `QUALITY.md` files and `/quality setup` guidance while keeping
the formal QUALITY.md format neutral. The stable Rating Level IDs stay plain
machine identifiers (`outstanding`, `target`, `minimum`, `unacceptable`);
human-facing titles carry a visible marker by default.

Detail lives in:

- [Functional spec](0075-rating-title-emoji-defaults/spec.md) - what the change
  must do.
- [Design doc](0075-rating-title-emoji-defaults/design.md) - how the default
  display convention lands.

## Motivation

Emoji-prefixed rating titles are genuinely useful in the agent-first workflow.
They make repeated ratings easier to scan in generated reports, review tables,
and model frontmatter without changing the model semantics. The project already
uses them in its own `QUALITY.md`, and evaluation reporting already preserves
Rating Level titles for human output. The remaining gap is the default authored
surface: CLI scaffolds, setup guidance, and examples still seed plain titles,
which means new models miss the visual affordance unless an author adds it
manually.

The default should be opinionated where it helps humans, but not formalized as a
format rule. QUALITY.md remains domain agnostic and accessible to authors who
prefer plain text or a custom scale.

## Scope

Covered:

- Seed emoji-prefixed titles for the standard four-level Rating Scale in CLI
  scaffolds and `/quality setup` default guidance.
- Preserve plain stable `level` IDs and all machine-facing references.
- Update user-facing examples and guide contracts where they describe the
  recommended default scale.
- Keep `SPECIFICATION.md` neutral: titles are human-readable labels and may
  contain visual markers, but the format does not require emoji.

Deferred / non-goals:

- No new schema field for icon/color/display metadata.
- No lint requirement or warning that requires emoji.
- No migration command for existing `QUALITY.md` files.
- No change to rating semantics, roll-up, or machine report identifiers.

## Affected artifacts

### Code

- [x] `internal/scaffold/skeleton.md` - seed emoji-prefixed Rating Level titles.
- [x] `internal/scaffold/skeleton-minimal.md` - seed emoji-prefixed Rating Level
      titles.
- [x] `internal/scaffold/scaffold_test.go` - assert scaffolded emoji-prefixed
      titles while stable IDs stay plain.

### Durable specs

- [x] `specs/skills/quality-skill/workflows/setup.md` - specify that setup's
      recommended standard scale uses emoji-prefixed human titles while stable
      IDs remain plain.
- [x] `specs/skills/quality-skill/guides/authoring-md.md` - specify that the
      authoring guide presents emoji-prefixed titles as the recommended default
      display convention.
- [x] `specs/skills/quality-skill/workflows/log.md` and
      `specs/skills/quality-skill/guides/log.md` - record spec updates.

### Format spec

- [x] `SPECIFICATION.md` - no change; the suggested-scale appendix remains plain
      to keep the formal spec neutral.

### Durable docs (README and bundled skill)

- [x] `README.md` - update the illustrative example where the public starter
      default should be visible there.
- [x] `skills/quality/guides/authoring.md` - present emoji-prefixed titles as the
      default display convention for the recommended standard scale.
- [x] `skills/quality/workflows/setup.md` - tell setup to author the recommended
      scale with emoji-prefixed titles.
- [x] `skills/quality/guides/getting-started.md` - align the rating-scale
      checklist with the default display titles.

### Release

- [x] `CHANGELOG.md` - add Unreleased CLI and `/quality Skill` notes.

## Children

- [Functional spec](0075-rating-title-emoji-defaults/spec.md) - required
  scaffold, setup, and documentation behavior.
- [Design doc](0075-rating-title-emoji-defaults/design.md) - landing strategy and
  trade-offs.

## Status

`Done`. Implemented the scaffold defaults, setup and authoring guidance, durable
skill specs/logs, README example, changelog, and scaffold tests. Verified with
`go test ./internal/scaffold` and `mise run check`. Archived.
