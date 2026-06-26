---
type: How-to Guide
title: Designing agent-mediated UX
description: How to design workflows that users experience through an AI assistant or coding agent.
tags: [agents, ux, workflows, contributing]
timestamp: 2026-06-24T00:00:00Z
---

# Designing agent-mediated UX

Use this guide when designing a workflow that a user experiences through an AI
assistant or coding agent.

**Agent-mediated UX** is the user's experience of a product, workflow, or task
as rendered by an agent: prompts, progress updates, questions, decisions,
confirmations, summaries, tool output, and generated artifacts. The agent is the
interface, so the workflow's words and structure deserve the same care as a
screen-based UI.

Prefer this term over "user-agent experience." `User agent` already has a
browser and HTTP meaning, while "agent-mediated" names the important design
fact: the user's path is mediated by an assistant.

## Core principle

Make the current state and next action obvious.

Every user-facing step should let a busy reader answer four questions quickly:

- Where are we in the workflow?
- What needs my attention?
- What is the recommended action?
- What will happen if I confirm?

If the user has to infer the main question from surrounding explanation, the
interaction is doing too much work in prose.

## Output hierarchy

Start with the status, then the primary action, then supporting context.

Good agent output has a visible left edge: headings and labels tell the user
what kind of information each line carries. Use short blocks instead of long
paragraphs when a user is choosing, confirming, or reviewing.

```text
**Discovery** - Step 2 of 6

**Should this `QUALITY.md` model the whole project?**

**Why it matters:** Sets the boundary for what this model evaluates.
**Recommended:** Whole current project
**Confidence:** High, based on `README.md` and repository layout.

**Answer:** Say "yes" to accept, or name a narrower area/path.
```

Avoid burying the primary call to action after rationale. Rationale matters, but
it supports the choice.

## Emphasis

Use Markdown emphasis as interaction structure.

- **Bold the primary question or call to action.** In each interaction block,
  the user's main task should be the strongest visual element.
- **Bold labels** such as `Recommended`, `Why it matters`, `Confidence`,
  `Changed`, `Validation`, `Important gaps`, and `Next`.
- Use *italics* for soft notes or caveats, not for required actions.
- Use `code` for exact files, commands, fields, model references, IDs, and
  literal values.
- Do not bold whole paragraphs, repeated prose, or every heading-like phrase.

Bold should make the output skimmable if the user reads only the labels and the
primary call to action. If everything is emphasized, nothing is.

## Progress

Show progress in workflows with multiple phases.

```text
**/quality setup**
Step 2 of 6: **Discovery**

Preflight: done
Context scan: done
Discovery: in progress
Review: next
Write: pending
Verify: pending
```

Keep progress factual. Do not turn it into a transcript of internal reasoning.
Update progress when the user's mental model would otherwise drift, especially
before a long context scan, after a tool-dependent phase, and before a mutation
gate.

## Discovery questions

A good agent-mediated question includes:

- the primary question;
- why the answer matters;
- the recommended answer, with the recommended option first for closed choices;
- confidence and evidence when available;
- the shortest acceptable response, preferably the option number for closed
  choices.

The question itself should be visually primary. The supporting fields make the
choice fast and informed.

```text
**Which lifecycle stage best fits this project?**

**Why it matters:** The stage calibrates which risks and requirements are
realistic to assess now.

**Recommended:** Active production
**Confidence:** Med, based on release notes and support docs.

**Answer:** Accept the recommendation, or choose exploratory, pre-release,
active production, maintenance, or sunset.
```

When a workflow uses a structured question tool, keep the teaching and rationale
in the surrounding message. Tool option labels are too small to carry the full
interaction design.

### Closed-choice questions

For small closed-choice sets, use numbered options and put the recommended
option first. The user's shortest accept path should be `1`.

For a true binary confirmation, especially a mutation gate, use `y`/`n` when the
question naturally means yes or no. Accept obvious aliases such as `yes`, `no`,
`1`, `apply`, or `skip` when they match the displayed options, but make `y` and
`n` the visible shortest responses.

Match the option labels to the question's framing. If the workflow stores an
internal value such as `lowTolerance`, but the user-facing question asks about
cost, present cost options and map the answer internally. Do not make the user
translate between axes while answering.

Good:

```text
**Question 4: How costly is poor quality here?**

**Why it matters:** This sets modeling rigor, risk strictness, and which Factors
need explicit Requirements.

1. High cost - poor quality can cause serious business, operational, financial,
   or trust impact. **Recommended**
2. Moderate cost - poor quality has meaningful cost but can often be contained.
3. Low cost - poor quality is usually recoverable or low impact.

**Confidence:** Medium, based on visible production workflows and external
integrations.

**Answer:** Reply `1` to accept the recommendation, or choose `2` or `3`.
```

Avoid:

```text
Options: high tolerance, moderate tolerance, low tolerance
Recommended: low tolerance
Answer: Reply accept to use low tolerance, or choose another option.
```

## Checkpoints

Use checkpoints when the agent has inferred context and needs correction rather
than open-ended invention.

Tables work well when the user is reviewing several related inferred values:

```text
**Human context checkpoint**

Please correct this draft with short fragments.

| Item | Draft | Confidence | What to do |
| --- | --- | --- | --- |
| Primary users | Maintainers and agent collaborators | Med | Confirm or correct |
| Other stakeholders | Not visible | Low | Name any, or leave Unknown |
| Missing context | Support expectations not found | Low | Point to docs or record as a gap |
```

The checkpoint's primary call to action should still be explicit. Do not end a
structured checkpoint with a broad catch-all question that makes the specific
dimensions feel optional.

## Decision gates

Before a workflow mutates files, creates external artifacts, sends messages, or
changes tooling, show the decision plainly.

```text
**Update existing `QUALITY.md`?**

**Changes:** Replace the starter model with a project-specific model.
**Evidence/reason:** The repository has enough context to draft Areas, Factors,
and initial Requirements.
**Recommended option:** Update `QUALITY.md` now.
**Alternatives:** Stop here, or only scaffold the file.
**Done criterion:** `qualitymd lint QUALITY.md` passes.

**Answer:** Reply `y` to update, or `n` to stop.
```

Decision gates should name both what will happen and what will not happen when
that boundary matters. For example, a setup workflow can write `QUALITY.md`
without running an evaluation, creating issues, or configuring automation.

## Closeout

Close with the outcome, validation, remaining gaps, and recommended next action.

```text
**Setup complete** ✅

**Changed:** `QUALITY.md`
**Validation:** lint passed
**Important gaps:** stakeholder context and support expectations are still thin.
**Not done:** no evaluation, no issues, no automations.
**Next:** continue iterating on the model before the first full evaluation.
```

Do not make the user reconstruct success from logs or command output. If a
verification step failed or could not run, say that directly and name the next
useful action.

## Emoji

Use emoji as semantic markers, not decoration.

Good uses:

- `✅` for a completed validation or workflow result.
- `⚠️` for an important gap or caution.
- `❓` for an unknown or unresolved question.
- `🟢`, `🔵`, `🟡`, and `🔴` for Rating Level display titles when the standard
  QUALITY.md scale fits.

Avoid emoji in every heading or label. Repeated decorative emoji reduces scan
quality and can make a serious workflow feel less trustworthy.

## Tone

Agent-mediated UX should be direct, calm, and operational.

- Say what is happening.
- Say why the user is being asked.
- Recommend a default when evidence supports one.
- Make terse answers acceptable.
- Avoid marketing language, cheerleading, and vague reassurance.

The user should feel guided, not managed.

## Checklist

Before shipping an agent-mediated workflow, check:

- The primary question or call to action is bolded in each interaction block.
- The recommendation and evidence are adjacent to the choice.
- Closed-choice questions put the recommended option first and accept `1` as the
  shortest confirmation, except binary confirmations that use `y`/`n`.
- The shortest acceptable user response is clear.
- Progress is visible for multi-step workflows.
- Mutation gates state the change, reason, alternatives, and done criterion.
- Closeout reports changed artifacts, validation, remaining gaps, and next
  action.
- Emoji, if present, carries status or semantics.
