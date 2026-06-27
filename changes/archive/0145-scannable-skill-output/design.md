---
type: Design Doc
title: Scannable Skill Output - design
description: How /quality runtime guidance adopts labeled, five-second-scan output templates.
tags: [docs, skill, ux, workflows]
timestamp: 2026-06-27T00:00:00Z
---

# Scannable Skill Output - design

## Context

Answers the [functional spec](spec.md) for change case
[0145](../0145-scannable-skill-output.md). The shared UX guide now defines
scannable output, and the `/quality` runtime skill needs to absorb that guidance
in the places users see the skill's output: review gates, workflow openings,
closeouts, status summaries, and next-workflow prompts.

## Approach

### Treat labels as templates, not prose decoration

Add concrete text templates where the runtime currently lists required content
without giving a shape. The templates use bold labels for fields users scan:
`Changed`, `Verification`, `Not changed`, `Next`, `Answer`, `Boundary`, and
workflow-specific variants. This follows the updated shared UX guide while
preserving existing workflow semantics.

### Keep established interaction types

Do not invent new gates. Existing run frames, decision briefs, review gates, and
closed-choice prompts stay the same interaction types. The change only makes
their text fallback and closeouts more scannable.

### Update runtime and durable specs together

Each runtime template change has a corresponding durable spec update:

- root skill contract and direct authoring -> `quality-skill.md`;
- workflow templates -> matching `workflows/*.md` specs;
- recommendation result -> recommendation follow-up guide and behavior specs;
- top-10 and getting-started summaries -> matching guide specs;
- evaluation closeout -> workflow and reporting specs.

### Preserve boundaries

The templates name what did not happen where that boundary matters, especially
for read-only review, setup, evaluate, update, and recommendation follow-up.
They do not change mutation behavior, artifact creation, evaluation data,
ratings, or CLI behavior.

## Spec response

- **Shared UX guidance** - satisfied by keeping the existing scannable-output
  guide update in scope.
- **Runtime interaction contract** - satisfied by updating `SKILL.md`.
- **Workflow output templates** - satisfied by updating setup, evaluate, review,
  improve, and update workflow guidance plus durable specs.
- **Runtime guide output templates** - satisfied by updating recommendation
  follow-up, top-10 checks, and getting-started guidance plus durable specs.
- **Verification** - satisfied by source inspection, Markdown formatting, and
  full repository checks.

## Alternatives

- **Only update the shared UX guide.** Rejected. The skill runtime guidance is
  what agents execute; leaving it prose-shaped would make adoption uneven.
- **Add a universal closeout macro.** Rejected. The project uses Markdown
  guidance, not a template engine, and each workflow has distinct boundaries.
- **Rewrite every output example in one pass.** Rejected. The scoped surfaces are
  the ones identified by audit as user-facing templates or closeouts; explanatory
  authoring prose can stay prose.

## Trade-offs & risks

- More templates can make runtime guidance longer. The added length is localized
  to user-facing output surfaces and replaces ambiguity with concrete shapes.
- Over-labeling can feel stiff. The templates use labels where multiple facts
  compete for attention; single-idea prose remains acceptable.

## Open questions

None.
