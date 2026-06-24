---
type: Design Doc
title: Rating title emoji defaults — design doc
description: Landing strategy for emoji-prefixed Rating Level title defaults across scaffold, setup guidance, and docs without changing rating semantics.
tags: [scaffold, skill, rating-scale, docs]
timestamp: 2026-06-24T00:00:00Z
---

# Rating title emoji defaults — design doc

Design behind the
[Rating title emoji defaults](../0075-rating-title-emoji-defaults.md) change case
and its [functional spec](spec.md).

## Context

Human reports already resolve Rating Level titles from the model, and the
project's own `QUALITY.md` uses emoji-prefixed titles. The default authoring
surfaces still seed plain titles, so new files miss the same scanning affordance
unless the author knows to add it. The design needs to make the default
friendlier while preserving the split between display labels and stable Rating
Level IDs.

## Approach

Land the change as a display-title default, not a semantic feature:

- Update both embedded scaffold templates to seed the same four title strings:
  `🟢 Outstanding`, `🔵 Target`, `🟡 Minimum`, `🔴 Unacceptable`.
- Keep every `level` value, criterion, and ordering unchanged.
- Update `/quality setup` runtime guidance and its durable setup workflow spec so
  the recommended scale is described as plain stable IDs with emoji-prefixed
  human titles.
- Update the authoring guide and its durable guide spec so authors understand the
  markers are useful defaults, not conformance rules.
- Update the README starter example because it is a user-facing example of the
  default scale; leave evaluation example fixtures unchanged unless tests require
  regeneration, because those are historical/reporting examples rather than
  starter defaults.
- Leave `SPECIFICATION.md` normative text unchanged. If touched at all, keep it
  to the non-normative suggested scale appendix.

The implementation relies on existing behavior: frontmatter parsing already
treats `title` as a scalar string, JSON Schema does not restrict Unicode, and
report rendering already uses model titles for human output while machine JSON
and gates use stable IDs.

## Alternatives

- **Add explicit icon/color fields.** Rejected for now. The need is a default
  display affordance, not a richer presentation model.
- **Require emoji in the format.** Rejected. It would over-specify a human style
  choice and make plain-text house styles invalid for no semantic gain.
- **Keep all public examples plain while changing only the scaffold.** Rejected.
  That hides the default from readers and makes setup/scaffold behavior look
  surprising.
- **Regenerate all evaluation examples with emoji titles.** Rejected unless a
  test forces it. Historical examples demonstrate report contracts and stable
  IDs; changing them adds churn without improving the starter default.

## Trade-offs & risks

- **Unicode in scaffolds.** The repository already uses emoji in `QUALITY.md` and
  tests, so this is not a new encoding dependency.
- **Accessibility.** Emoji-only labels would be weak; the default keeps the word
  label after the marker, and guidance should say titles remain customizable.
- **Spec drift.** The formal format stays neutral, so durable skill specs and
  guides carry the display convention while `SPECIFICATION.md` continues to
  define only the required human-readable `title`.

## Open questions

- Should the non-normative suggested scale in `SPECIFICATION.md` stay plain to
  emphasize neutrality, or show the preferred default to keep all starter
  examples aligned?
