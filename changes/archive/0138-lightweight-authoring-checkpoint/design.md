---
type: Design Doc
title: Lightweight Authoring Checkpoint - design
description: How the /quality skill adds a conversational checkpoint for direct QUALITY.md edits.
tags: [skill, authoring, ux]
timestamp: 2026-06-27T00:00:00Z
---

# Lightweight Authoring Checkpoint - design

## Context

Answers the [functional spec](spec.md) for change case
[0138](../0138-lightweight-authoring-checkpoint.md). The change is prompt and
documentation behavior: no CLI or code path edits are needed. The runtime skill
already has setup discovery, recommendation follow-up confirmation, decision
briefs, and quality-log rules; direct `QUALITY.md` edits need a lighter bridge
between ordinary user requests and those existing mutation safeguards.

## Approach

### Add direct authoring as a routed path, not a public workflow

Update the root skill prompt and durable skill spec so free-form requests to edit
an existing `QUALITY.md` route to direct model authoring unless they match setup,
evaluate, tooling `update`, or recommendation follow-up. Keep public workflows
unchanged: `setup`, `evaluate`, and `update` remain the public workflow names,
and recommendation follow-up remains a non-public follow-up path.

This avoids overloading `/quality update`, which already means skill/CLI
maintenance, and avoids introducing a command-shaped `/quality author` surface.

### Reuse existing authoring guides

Direct authoring reads `guides/authoring.md` first and then only the routed
sub-guides matching the likely mutation surface. The new entry-guide guidance
lists the compact considerations the agent should preserve:

- intent;
- target;
- rationale;
- judgment effect;
- unresolved unknowns; and
- quality-log routing.

Those considerations are not a questionnaire. They are the agent's checklist for
deciding whether it can infer enough or needs a short follow-up.

### Insert a lightweight checkpoint before mutation

Add a reusable checkpoint pattern to the root prompt:

```text
I’m reading this as: <plain-language intent>.

I’ll update <QUALITY.md target>, keep <boundary> unchanged, and <write/not write>
a quality-log entry because <reason>.

Anything you want adjusted first? You can say `looks good`, or tell me any
concerns, goals, needs, worries, or constraints I should account for.
```

The checkpoint is a confirmation surface for ordinary direct edits. It names the
mutation clearly enough that `looks good` is explicit approval. The wording also
solicits the user's concerns and goals without requiring a form.

### Escalate only when model risk is high

The existing decision-brief contract remains the heavier gate. Direct authoring
uses it when the edit changes rating semantics, removes coverage, shifts scope or
apex, or otherwise carries high judgment risk. In those cases, the checkpoint's
teaching and the decision brief can be combined into one block so the user is not
asked twice.

### Keep quality-log judgment unchanged

Direct authoring does not change what counts as a meaningful model change. It
only adds another confirmed source of such changes after setup. The quality-log
entry uses the checkpoint or decision-brief rationale; wording-only, formatting,
typo, and non-judgment-changing body edits still do not write the log.

## Spec response

- **Direct authoring dispatch** - satisfied by adding routing language to the
  root prompt and durable skill spec.
- **Intent inference and follow-up** - satisfied by explicit infer-first and
  material-follow-up guidance in the root prompt, plus authoring-guide
  considerations.
- **Intent checkpoint and confirmation** - satisfied by the reusable checkpoint
  text and `looks good` confirmation rule.
- **Quality-log routing** - satisfied by aligning the runtime and durable
  quality-log guidance with confirmed direct authoring.

## Alternatives

- **A full setup-like questionnaire for every edit.** Rejected. It captures
  context, but it makes routine edits too slow and conflicts with the requested
  conversational feel.
- **A formal direct-authoring workflow with a run frame.** Rejected for this
  slice. Direct authoring is a common path inside the skill, not a new public
  invocation surface.
- **Always use decision briefs.** Rejected. Decision briefs are right for risky
  model changes, but too heavy for obvious wording, body-context, or narrow model
  edits.
- **Skip confirmation for obvious edits.** Rejected. The lightweight checkpoint
  preserves explicit approval while making the approval path casual and short.

## Trade-offs & risks

- The phrase `looks good` is intentionally broad. The runtime guidance limits it
  to checkpoints that clearly name the mutation, which keeps casual approval from
  applying to vague or unresolved edits.
- This change relies on agent judgment to decide when follow-up is material. The
  follow-up triggers and high-risk escalation list keep that judgment bounded
  without turning it into a rigid form.

## Open questions

None.
