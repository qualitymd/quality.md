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

At a workflow's opening, one question comes before those four: did the agent
understand what I asked, and what is the plan? The [Opening](#opening) section
covers how to answer it immediately.

If the user has to infer the main question from surrounding explanation, the
interaction is doing too much work in prose.

A feedback invitation is a gate, not decoration. If an agent asks what the user
wants adjusted, asks for concerns, or names a confirmation phrase such as
`looks good`, the agent must wait for the response before mutating. Do not offer
agency in one sentence and remove it in the next by saying the agent is
proceeding anyway.

## Channels and progressive enhancement

Design the interaction, not the rendering.

The unit of design is the _intent_ — for example, "a single-select closed choice
with a recommended default and rationale" — not the Markdown that happens to
render it. A given runtime may render that intent through a native interaction
affordance (a selectable option list, a confirm gate, a multi-select) or through
plain Markdown. Treat the affordance as **progressive enhancement** over a text
rendering that always works on its own.

Two rules keep this both useful and agent-agnostic:

- **Always author a working text rendering.** A harness with no question tool
  must still get a complete, well-structured interaction. The numbered list and
  the `y`/`n` gate described later are the _fallback renderings_ of these
  intents, not a lesser path.
- **Keep the semantics in the message, not the widget.** Native option labels
  are small, vary in how much they display, and are sometimes truncated or
  stripped. The widget carries the _selection mechanic_; the surrounding message
  carries the _teaching_ — the question, why it matters, the recommendation, and
  the evidence. Never compress design-critical rationale into an option label.

Express the choice as intent plus affordance category, never as a named tool:
"present a single-select closed choice with the recommended option first," and
let whatever the runtime offers fulfill it. When in doubt, detect the
capability rather than assume it: _if the runtime exposes a structured
single-select affordance, render through it; otherwise emit the numbered-list
fallback._ Both branches are part of the design.

### Native interaction affordances to watch for

Look for these affordance categories and map each interaction to the richest one
that fits. Describe them by capability, never by a specific tool name.

- **Single-select closed choice** — the numbered-list interaction; maps to an
  option/radio picker.
- **Multi-select** — "which of these apply"; maps to a checklist picker.
- **Binary confirmation / approval gate** — the mutation gate; maps to a native
  confirm or approve affordance.
- **Plan or diff review-and-approve** — a richer gate where the artifact under
  review is itself rendered for inspection before approval.
- **Permission or tool-authorization prompt** — harness-level. Do not
  reimplement it in prose: if the harness will already prompt to authorize the
  mutation, a second hand-rolled gate is redundant friction.
- **Free-text input** — when the answer is open-ended; the right choice when a
  picker would be wrong.
- **Progress or task-list indicator** — native status UI that can replace or
  supplement the textual progress block.

### When a native affordance is _not_ fit-for-purpose

Reaching for a widget is not always right. Prefer the text rendering when:

- **Cardinality is open or unknown** — do not force a picker over a list you
  cannot enumerate; use free text.
- **The rationale will not survive in labels** — you may still use the widget
  for the pick, but the teaching and evidence stay in the message.
- **The harness already gates it** — do not stack a prose gate on top of a
  permission prompt for the same mutation.
- **The widget cannot place the recommendation and evidence next to the choice**
  — keep them adjacent in the message so the core principle still holds.

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

## Scannable output

Design output for a five-second scan. A user should be able to tell what
happened, what matters, and what to do next without reading every sentence.

Prefer labeled blocks over paragraphs when output contains status, choices,
verification, risks, boundaries, or next actions. Good labels include:

- **Planned edit**
- **Why**
- **Approach**
- **Boundary**
- **Recommendation**
- **Risk**
- **Changed**
- **Verification**
- **Not changed**
- **Next**
- **Answer**

Use paragraphs for one idea. Use bullets, numbered lists, or labels when there
are multiple independent facts.

Rules:

- Put the result, recommendation, or ask first.
- Keep paragraphs to one short idea.
- Use numbered lists for ranked options or ordered steps.
- Use bullets for unordered facts.
- Keep list items parallel.
- Put the shortest response path on its own line or label.
- Do not bury the required user action after rationale.

Good:

```text
**Changed**

- Added archived Change Case `0144 - Pointed Review Gates`.
- Updated `docs/guides/agent-mediated-ux.md`.
- Updated `/quality` runtime and durable skill guidance.

**Verification**

- `mise run fmt-md-check` passed.
- `mise run check` passed.
- Worktree is clean.
```

Avoid:

```text
I updated the guide and the skill guidance to reflect the new review-gate
approach, and I also added the archived Change Case with its spec and design,
then I ran the formatting and repository checks, both of which passed, and the
worktree is clean.
```

Even when the content is correct, long single-flow paragraphs make the user parse
too much before finding the outcome and next action.

## Emphasis

Use Markdown emphasis as interaction structure.

- **Bold the primary question or call to action.** In each interaction block,
  the user's main task should be the strongest visual element.
- **Bold labels** such as `Recommended`, `Why it matters`, `Confidence`,
  `Changed`, `Validation`, `Important gaps`, and `Next`.
- Use _italics_ for soft notes or caveats, not for required actions.
- Use `code` for exact files, commands, fields, model references, IDs, and
  literal values.
- Do not bold whole paragraphs, repeated prose, or every heading-like phrase.

Bold should make the output skimmable if the user reads only the labels and the
primary call to action. If everything is emphasized, nothing is.

Do not rely on bold alone for hierarchy. Some surfaces strip or flatten
emphasis, and the moment bold disappears, a layout held together only by bold
collapses into an undifferentiated wall. Carry the hierarchy in **position,
blank-line separation, and indentation** first, and treat bold as reinforcement
of a layout that already reads when emphasis is removed. A practical test: strip
every `**` from a block and check that the question and the response path are
still obvious.

## Opening

Open every workflow by confirming intent and previewing the path, before doing
the work.

The opening does two jobs:

- **Intent reflection.** Say back what you understood the user to want — the
  resolved workflow, the target, the scope — so a wrong inference is caught now,
  while it is cheap to correct, rather than after a long silent run.
- **Path preview.** State what will happen: what is read-only versus mutating,
  which artifacts will change, and where the workflow ends.

The carrier for both is a concise **run frame** emitted as the **first output**,
before discovery or any tool call. The frame's value is that the user can catch a
wrong inference before the agent acts; that value is lost if the agent reads
files and runs commands first and only frames the work afterward.

```text
**QUALITY.md · evaluate**

- **Model file:** `QUALITY.md`
- **Scope:** full evaluation
- **Mutation:** evaluation artifacts and feedback log under `.quality/`
- **Artifacts:** numbered evaluation run, structured data, Markdown report tree
- **Next gate:** report findings, ratings, and limits before any follow-up action
```

When a field genuinely needs a tool to resolve — a scope that spans many
modeled areas, for instance — emit the frame anyway with a best-known or
`resolving…` value, then confirm the resolved value in a later message. A
provisional frame still beats silence.

Avoid the **silent runway**: a long sequence of reads and commands before the
user sees any frame. From the user's side it is indistinguishable from a stall,
and it removes the early checkpoint the opening exists to provide.

For non-workflow tasks that still need meaningful reading or inference before a
mutation — for example, direct `QUALITY.md` authoring that requires the agent to
read the model and authoring guidance — send a short acknowledgement before that
work begins. The acknowledgement should name the understood task, say what the
agent will inspect, and state whether a feedback gate will happen before files
change.

```text
I’ll treat this as a `QUALITY.md` model change. I’m going to inspect the current
model and relevant authoring guidance, then I’ll show the intended edit for
feedback before changing files.
```

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

When a workflow renders a question through a native single-select or multi-select
affordance, keep the teaching and rationale in the surrounding message. Tool
option labels are too small to carry the full interaction design. See
[Channels and progressive enhancement](#channels-and-progressive-enhancement)
for when to reach for the affordance and when to stay in text.

### Closed-choice questions

A closed-choice question is an _intent_: a single-select pick with a recommended
default. Render it through a native option picker when one is present and
fit-for-purpose; otherwise render the numbered-list fallback below. Either way
the recommended option comes first and the rationale stays in the message.

For the numbered-list fallback, use numbered options and put the recommended
option first. The user's shortest accept path should be `1`.

For a true binary confirmation, especially a mutation gate, prefer a native
confirm affordance when present. In the text fallback, use `y`/`n` when the
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

**Why it matters:** This sets modeling rigor, risk strictness, and which factors
need explicit requirements.

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

## Review gates

Use review gates when the agent has inferred an edit or plan and needs the user
to correct, approve, or add context before mutation. A review gate can be lighter
than a binary decision gate: it may invite feedback in conversational terms and
accept a phrase such as `looks good` as approval. It is still a real gate.

Distinguish three interaction shapes:

- **Informational preview** - tells the user what the agent is about to inspect
  or do next, without asking for a response.
- **Review gate** - asks for corrections, concerns, goals, needs, constraints,
  or `looks good`; the agent waits for the response before mutating.
- **Decision gate** - asks the user to choose between proceeding and a
  non-mutating alternative; the agent waits for explicit approval before
  mutating.

When a review gate authorizes a content or model edit, it should state the
inferred intent and why the change appears needed, then describe the planned
change in simple conversational prose with a visible value prop. Prefer a
natural `so that` clause when it fits: the user can correct the goal as well as
the mechanics. The gate should not end with only a generic "anything to adjust?"
prompt. When the agent can infer a likely goal, risk boundary, or scope
assumption, name the most consequential one and ask the user to react to that
specific interpretation. Use labeled blocks when the gate carries more than two
distinct facts, such as intent, rationale, approach, boundary, logging, and
answer path. Keep important boundaries, logging decisions, and approval paths
visible without turning the gate into a formal questionnaire. For example:

```text
**Planned edit:** Add root-level security coverage.

**Why:** Whole-repo evaluation should catch security-sensitive concerns even
when the details live in child areas.

**Approach:** Add a root `security` factor with one requirement that checks
descendant areas for security-related factors, requirements, and evidence.

**Boundary:** This is broad cross-area security coverage, not only a roll-up of
`api-service/security`.

**Also:** I’ll add one `.quality/changelog/` entry because this changes model judgment.

**Answer:** Say `go`, or say if you want narrower appsec-only scope.
```

After this block, stop. If the agent should proceed without waiting, do not ask
for feedback; render an informational preview instead. A request for feedback
followed by "I’m proceeding" is a false affordance:

```text
Anything you want adjusted first? I’m proceeding with that interpretation.
```

For edits that reshape future judgment — such as `QUALITY.md` model changes,
rating criteria changes, scope shifts, or removal of coverage — prefer a review
gate even when the intent seems clear. Clear intent reduces follow-up questions;
it does not remove the user's chance to catch a wrong inference before the
mutation.

Use numbered planned-action lists only when the edit has multiple independent
parts that would be hard to scan in prose. The default review gate is a short
plain-language plan, not a task list.

A specific steering axis is better than an exhaustive menu. If the most
consequential assumption is broad security versus narrower appsec, ask that. Do
not dilute it with every possible concern, naming preference, edge case, and
constraint unless those dimensions are equally likely to change the edit.

Avoid long single-flow prose for review gates:

```text
I’m reading this as: add a model-wide security factor at the root so whole-repo
evaluation can judge security-sensitive concerns even when the detailed factor
lives in a child area. I’d implement that as a root security factor with one
requirement that inspects descendant areas for security-related factors,
requirements, and evidence...
```

The content may be right, but the user has to parse too much before finding the
decision.

## Decision gates

Before a workflow mutates files, creates external artifacts, sends messages, or
changes tooling, show the decision plainly.

A gate is a binary-confirmation intent. Render it through a native confirm or
approve affordance when one is present and fit-for-purpose; the text block below
is its fallback rendering. When the mutation is something the harness will
already prompt to authorize — a tool-permission or approval prompt — do not
stack a second prose gate on top of it; that is redundant friction. For a
file-mutating plan, prefer a native plan-or-diff review affordance when present
so the user inspects the artifact before approving.

A gate has one job: make the user see _what is being asked_ and _how to respond_
at a glance. Lead with the question, render the choices as a visually separated
block, and demote the rationale below them. Do not stack the question and its
supporting labels at the same weight — that is the flat-wall failure, where the
call to action is just one more bolded line at the bottom and disappears
entirely if bold is not rendered.

```text
**Update existing `QUALITY.md`?**
Replace the starter model with a project-specific model.

  [y] Update `QUALITY.md` now  — recommended
  [n] Stop, or only scaffold the file

Reason: repo has enough context to draft areas, factors, and requirements.
Done when: `qualitymd lint QUALITY.md` passes.
Not changed: no evaluation, no issues, no automations.
```

The question is the only bolded line; the indented `[y]`/`[n]` block is the next
thing the eye lands on; the supporting context sits below in plain `label:`
lines so the gate still reads when emphasis is stripped.

Rules for a gate:

- **Question first and visually primary.** It outranks every supporting label.
- **Choices as a separated block**, one per line, with the recommended option
  marked inline. Fold "Alternatives" into the `[n]` line — for a binary gate the
  two choices already carry it.
- **Cap supporting fields at about three.** Beyond that the rationale buries the
  choice. Keep what the user needs to decide (reason, done criterion) and the
  boundary line when it matters; drop the rest.
- **Name what will not happen** when that boundary matters. For example, a setup
  workflow can write `QUALITY.md` without running an evaluation, creating issues,
  or configuring automation.

Avoid the flat stack, where the question, five supporting labels, and the call
to action all carry equal weight and the choice is a prose afterthought:

```text
**Apply update plan?**

**Changes:** CLI only (0.15.0 → 0.16.0)
**Evidence/reason:** The loaded skill requires CLI >=0.16.0.
**Recommended option:** Run the owner command.
**Alternatives:** Skip and leave the CLI at 0.15.0.
**Done criterion / verification:** Re-run `qualitymd version --json`.
**Answer:** Reply `y` to apply, or `n` to skip.
```

Six labels compete for attention, the `y`/`n` choice is the last line, and with
bold removed nothing distinguishes the ask from its rationale.

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
- `🟢`, `🔵`, `🟡`, and `🔴` for rating level display titles when the standard
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

- The workflow opens with a run frame as its first output, confirming intent and
  previewing the path before any tool call.
- Interactions render through a fit-for-purpose native affordance when one is
  present, with a complete text fallback and the rationale carried in the
  message rather than the widget labels.
- No prose gate is stacked on a mutation the harness already prompts to
  authorize.
- Output passes the five-second scan: result, importance, and next action are
  visible without reading every sentence.
- Multi-fact responses use labels, bullets, or numbered lists instead of dense
  paragraphs.
- The primary question or call to action is the strongest element in each
  interaction block, by position and structure — not bold alone.
- Reading only the first line and the choice block, with emphasis stripped,
  makes the ask and the shortest response obvious.
- The recommendation and evidence are adjacent to the choice.
- Closed-choice questions put the recommended option first and accept `1` as the
  shortest confirmation, except binary confirmations that use `y`/`n`.
- The shortest acceptable user response is clear.
- Progress is visible for multi-step workflows.
- Long non-workflow reads or inference begin with a quick acknowledgement rather
  than a silent runway.
- Any feedback invitation is a real gate: the agent waits after asking for
  adjustments, concerns, or `looks good`.
- Content and model review gates state the inferred intent, why the change
  appears needed, the planned change, the value prop, and the most consequential
  assumption the user should react to before mutation.
- Review gates use labels when they need to carry intent, approach, boundary,
  logging, and answer path.
- Mutation gates lead with the question, render choices as a separated block,
  cap supporting fields at about three, and state the change, reason, and done
  criterion.
- Closeout reports changed artifacts, validation, remaining gaps, and next
  action.
- Emoji, if present, carries status or semantics.
