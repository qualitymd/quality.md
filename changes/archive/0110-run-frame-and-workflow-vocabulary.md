---
type: Change Case
title: Run frame title and workflow vocabulary
description: Retitle the /quality run frame so it stops rendering a fake "/quality run" command and a "Mode:" label, and unify the public-surface concept on "workflow" instead of "mode".
status: Done
tags: [skill, quality, ux, agent-mediated]
timestamp: 2026-06-26T00:00:00Z
---

# Run frame title and workflow vocabulary

This parent concept captures the why and status; the detail lives in its
children:

- [Functional spec](0110-run-frame-and-workflow-vocabulary/spec.md) - what the
  change must do.
- [Design doc](0110-run-frame-and-workflow-vocabulary/design.md) - how it is
  implemented, and why.

## Motivation

Every public `/quality` workflow opens with a run frame whose header is rendered
literally as `**/quality run**` and whose first field is `**Mode:**`. Both
tokens surface vocabulary the public surface does not actually have:

- `/quality run` reads as an invokable command, but there is no `run`
  invocation. The real invocations are `/quality`, `/quality setup`,
  `/quality evaluate`, and `/quality update`. A user who copies what the frame
  shows and types `/quality run` is typing something the skill never defines —
  the same class of leak that the skill is explicitly told to avoid for
  `status`, `next`, and `review` pseudo-commands.
- `Mode:` is internal vocabulary the skill itself waffles on. `SKILL.md` writes
  "Mode/workflow" and "not as a mode run" in the same breath as routing every
  request through `workflows/`, and the durable spec calls the same thing a
  "public mode." One concept carries two names, and the frame surfaces the
  internal one to the user.

The run frame is worth keeping — a status-first preamble lets the user catch a
wrong inference before the skill mutates anything (the 0038 rationale). The
problem is purely the leaked labels in its header. This case retitles the frame
and settles the vocabulary so the public surface is described consistently as a
set of **workflows**.

## Scope

Covered:

- retitle the run frame header to identify the workflow without rendering a
  command-style string, and remove the `Mode:` field (the workflow name moves
  into the title);
- add a durable run-frame requirement forbidding a header that reads as an
  invokable command or a `Mode:` field label, so the lesson survives;
- unify the public-surface concept on "workflow", retiring "mode" as a name for
  a `/quality` workflow across the durable skill specs, bundled runtime skill,
  and the two durable docs that use it; and
- update append-only skill/spec/changes logs and indexes for the touched durable
  artifacts.

Deferred / non-goals:

- no change to the run frame's required *fields* beyond removing `Mode:` (model
  file, scope, rigor, mutation, artifacts, next gate are unchanged);
- no change to which workflows exist or how invocations resolve;
- the internal term "run frame" is retained — it is instructional, never
  rendered to users, and renaming it would break the `#run-frames` anchor that
  append-only history references;
- no CLI, Go, format-schema, rating, roll-up, evaluation-record, or report
  change; and
- no change to unrelated uses of "mode" (failure modes, CLI output modes,
  Diátaxis doc modes, OS file modes).

## Affected artifacts

### Code

- [ ] None - no Go, CLI, or generated report implementation change.

### Format spec

- [ ] None - `SPECIFICATION.md` and the QUALITY.md format are unaffected.

### Durable specs

- [x] [`specs/skills/quality-skill/quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md)
      - update the Run frames section: name the resolved *workflow*, forbid a
      command-style header and a `Mode:` field label; retire "mode" as the
      public-surface concept across the shared contract.
- [x] [`specs/skills/quality-skill/index.md`](../../specs/skills/quality-skill/index.md)
      - "cross-mode" → "cross-workflow"; "non-mode" follow-up wording.
- [x] [`specs/skills/quality-skill/evaluation.md`](../../specs/skills/quality-skill/evaluation.md)
      - "cross-mode"/"evaluate mode"/"evaluation mode" → workflow wording.
- [x] [`specs/skills/quality-skill/recommendation-follow-up.md`](../../specs/skills/quality-skill/recommendation-follow-up.md)
      - "non-mode"/"runtime mode" → workflow wording.
- [x] [`specs/skills/quality-skill/workflows/evaluate.md`](../../specs/skills/quality-skill/workflows/evaluate.md)
      - frontmatter tag and prose "mode" → "workflow".
- [x] [`specs/skills/quality-skill/workflows/setup.md`](../../specs/skills/quality-skill/workflows/setup.md)
      - frontmatter tag and "dispatched as a mode" → "workflow".
- [x] [`specs/skills/quality-skill/workflows/update.md`](../../specs/skills/quality-skill/workflows/update.md)
      - frontmatter tag and prose "mode" → "workflow"; "public `/quality` run
      frame" wording.
- [x] [`specs/skills/quality-skill/workflows/index.md`](../../specs/skills/quality-skill/workflows/index.md)
      - "dispatched as a mode" → "workflow".
- [x] [`specs/skills/quality-skill/workflows/evaluate/feedback-log.md`](../../specs/skills/quality-skill/workflows/evaluate/feedback-log.md)
      - run-frame reference wording (term retained, no change required unless it
      says "mode").
- [x] [`specs/skills/quality-skill/guides/recommendation-follow-up-md.md`](../../specs/skills/quality-skill/guides/recommendation-follow-up-md.md)
      - "not a `/quality` mode" → "workflow".
- [x] [`specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md`](../../specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md)
      - "Other modes" → "Other workflows".

### Durable docs / bundled skill

- [x] [`skills/quality/SKILL.md`](../../skills/quality/SKILL.md)
      - retitle the run-frame template, drop the `Mode:` line, retire "mode" in
      Arguments and Workflow Dispatch wording.
- [x] [`skills/quality/workflows/evaluate.md`](../../skills/quality/workflows/evaluate.md)
      - retitle the rendered frame, drop `Mode:`; H1 "Evaluate Mode" → "Evaluate
      Workflow".
- [x] [`skills/quality/workflows/setup.md`](../../skills/quality/workflows/setup.md)
      - retitle the rendered frame, drop `Mode:`.
- [x] [`skills/quality/workflows/update.md`](../../skills/quality/workflows/update.md)
      - retitle the rendered frame, drop `Mode:`; H1 "Update Mode" → "Update
      Workflow".
- [x] [`skills/quality/guides/recommendation-follow-up.md`](../../skills/quality/guides/recommendation-follow-up.md)
      - "not a `/quality` mode" → "workflow".
- [x] [`skills/quality/resources/cli-quick-reference.md`](../../skills/quality/resources/cli-quick-reference.md)
      - run-frame references (term retained; confirm no "mode" wording).
- [x] [`docs/reference/versioning.md`](../../docs/reference/versioning.md)
      - "The mode" describing `/quality update` → "The workflow".
- [x] [`docs/guides/cut-a-release.md`](../../docs/guides/cut-a-release.md)
      - "skill modes" → "skill workflows".

### Suggested new durable specs

- None. The existing `/quality` skill specs already own the run-frame contract
  and workflow vocabulary.

## Status

`Done`. Implemented across the durable `/quality` skill specs, the bundled
runtime skill, and the two durable docs. The run-frame header is now
`**Quality · <workflow>**` with no `Mode:` field, the spec's Run frames section
forbids a command-style header or `Mode:` label, and "mode" is retired in favor
of "workflow" for the public-surface concept (recommendation follow-up is
described as a post-evaluation follow-up that is not a public workflow). Append-only
skill, spec, doc, and changes logs are updated. Verified with
`mise run fmt-md-check` and a residual-leak sweep of the live surface. Archived
on landing. See the [status lifecycle](../index.md#status-lifecycle).
