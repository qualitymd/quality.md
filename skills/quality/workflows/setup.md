# Setup Workflow

Run this workflow to create or update a useful first `QUALITY.md`. Setup writes
the selected `QUALITY.md` and writes a current-run workflow feedback log under
`.quality/logs/`; evaluation, quality-log entries under `.quality/log/`, external
issues, recommendation handoff, and recurring-review automation are follow-on
workflows.

## Workflow

```text
Preflight
- verify CLI support
- resolve model file
- emit run frame
- create the current-run workflow feedback log

Read context
- inspect setup signals
- build setup brief
- identify missing context
- update feedback log for material workflow-experience events

Ask discovery questions
- teach and ask questions 1-5, then present a human context checkpoint covering
  questions 6-9
- iterate one at a time when there is no structured question affordance
- page through a structured question tool when one is available
- update feedback log for material workflow-experience events

Review and confirm
- recap every question with its final answer
- wait for an explicit review-gate response before authoring
- update feedback log for material workflow-experience events

Write QUALITY.md
- scaffold missing file with qualitymd init
- gate existing-file edits with a decision brief
- synthesize body first, then frontmatter
- update feedback log for material workflow-experience events

Verify and close
- run qualitymd lint
- classify maturity with the Top 10 checklist
- report status and next-step choices
- finalize the workflow feedback log
```

## Preflight

1. Verify the CLI prerequisite from `SKILL.md`.
2. Resolve the model file: explicit path when supplied, otherwise `QUALITY.md`
   in the current working directory. Do not walk parent directories.
3. Emit the run frame:

   ```text
   /quality run
   - Mode: setup
   - Model file: <resolved path>
   - Scope: contextual QUALITY.md setup
   - Mutation: QUALITY.md + workflow feedback log under .quality/logs/
   - Artifacts: QUALITY.md, .quality/logs/<timestamp>-setup-feedback-log.md
   - Next gate: create feedback log, context analysis, discovery, lint, maturity inspection
   ```

4. Create the current run's workflow feedback log under
   `.quality/logs/<timestamp>-setup-feedback-log.md`, creating `.quality/logs/`
   on demand. Use a sortable UTC, filesystem-safe `<timestamp>` such as
   `2026-06-23T154233Z`, and never overwrite a feedback log from another run.
   The initial log must include the frontmatter and body sections in
   [Workflow feedback log](#workflow-feedback-log), with `status: in-progress`.
5. Read the authoring sections a first model needs, not the whole guide. For a
   first model read [`../guides/authoring.md`](../guides/authoring.md) sections
   The QUALITY.md file, Quality Model, The Markdown body, and Rating Scale, plus
   the Area and Factor sections enough to shape candidates. Defer deep
   authoring detail until you actually author requirements or ratings, then read
   the relevant sections. Read
   [`../guides/getting-started.md`](../guides/getting-started.md) when setup is
   continuing from a starter/immature model or needs first-run iteration
   guidance.

## Read Context

Inspect available repository context before asking setup questions. Keep this
bounded to setup signals, not source-quality judgment.

Useful signals include README/docs, package metadata, repository structure,
tests/build scripts, contributor docs, existing agent instructions, and visible
work-management or recurring-review hints.

Treat the current directory as the default root area convention unless the user
supplied an explicit model path or context strongly indicates a narrower root.

When the root spans multiple workspaces, packages, or services, delegate a
bounded `Explore` fan-out only when proportional. Use one bounded pass per
component to capture purpose, entry points, external systems, risk surfaces, and
test/CI coverage. Feed the result into candidate Area shape. Do not turn this
into source-quality evaluation.

Build this setup brief as working context:

```text
Root area: <default> (<Low | Med | High>, <evidence>)
Domain: <default> (<confidence>, <evidence>)
Lifecycle: <default> (<confidence>, <evidence>)
Risk tolerance: <default> (<confidence>, <evidence>)
Modeling rigor: <default> (<confidence>, <evidence>)
Rating scale: <default> (<confidence>, <evidence>)
Collaboration: <default> (<confidence>, <evidence>)
Primary users/outcomes: <default> (<confidence>, <evidence>)
Maintainers/collaborators: <default> (<confidence>, <evidence>)
Other stakeholders: <default> (<confidence>, <evidence>)
Missing context: <specific gaps>
Candidate model shape: <factors and child areas>
```

Use only these confidence labels:

- `High` — strong, specific repository evidence.
- `Med` — partial or indirect evidence.
- `Low` — weak or no evidence. When a default has no supporting signal, label it
  `Low` and name the absence in the note (e.g. `Low (no signal in repo)`), so the
  "pure default, no evidence" case stays legible.

Always keep the short evidence note next to the label; the label plus note is
what makes the three-level scale carry full meaning.

The setup brief is not a new artifact. It guides the discovery prompt and then
gets distilled into `QUALITY.md` where the assumptions shape the model.

## Ask Discovery Questions

Ask the user to confirm or correct the setup assumptions before writing
`QUALITY.md`. Questions 1-5 carry a recommended answer and confidence label.
The human context checkpoint then presents the repository-inferred context for
questions 6-9 as a draft the user can confirm, correct, fill in tersely, or
point to missed agent-accessible evidence for. These discovery inputs do double
duty: they capture context the model needs, and they teach the user the
dimensions a quality model spans.

Setup optimizes for *teaching* the user those dimensions over minimizing
interaction round-trips. Setup runs roughly once per project, so spending the
interaction to make each dimension legible — and to leave the user knowing why
each answer shapes the model — is worth more than saving round-trips. Ask every
one of questions 1-5 on every run, including ones whose inferred default is
high-confidence: high confidence is a reason to *recommend* a default firmly,
never to hide the question. Present the human context dimensions from questions
6-9 every run as a checkpoint — never drop or silently default them away to fit
an interaction surface — and do not optimize the teaching back out to save turns.

Each question and checkpoint item below carries authored teaching copy — **Why
it matters** (what the dimension shapes in `QUALITY.md`) — that you present to
the user as prose around the question or checkpoint. Surface this copy on
whatever interaction surface you use; it is the input's instruction value, not
optional flavor. You may state once, around the discovery flow or final recap,
that `QUALITY.md` is a living document and that setup answers can be revised
later.

The discovery inputs:

```text
1. Root area: Should this QUALITY.md model the whole current project, or a
   narrower area?
   Recommended: <default> (<confidence>)
   Why it matters: Sets the model's boundary — what this QUALITY.md evaluates and
   what falls outside it. Shapes the root Area, the Scope body, and the
   `quality-md` self-check Area's `source`.

2. Domain: What kind of thing is this model evaluating?
   Recommended: <default> (<confidence>)
   Why it matters: Names the kind of thing under evaluation (software, document,
   dataset, service, process, agent, and so on) so Factors and evidence use the
   right vocabulary. Shapes the Overview and the candidate Factor set.

3. Lifecycle: Which stage best fits?
   Options: exploratory, pre-release, active production, maintenance, sunset
   Recommended: <default> (<confidence>)
   Why it matters: The stage calibrates how much rigor and which risks matter
   now. Shapes Scope, Risks, and which Requirements are realistic to assess yet.

4. Risk tolerance: How costly is poor quality here?
   Options: high tolerance, moderate tolerance, low tolerance
   Recommended: <default> (<confidence>)
   Why it matters: How costly poor quality is drives modeling rigor and which
   Factors earn explicit Requirements rather than stay descriptive. Shapes the
   Risks section and Requirement strictness.

5. Rating scale: Should this model use the recommended four-level Rating Scale?
   Options: recommended four-level scale, pass/fail gate, custom scale needed
   Recommended: recommended four-level scale (<confidence>)
   Why it matters: Rating Levels are configurable in QUALITY.md; they are not
   baked into the format. The recommended scale keeps stable IDs as
   `outstanding`, `target`, `minimum`, `unacceptable`, and uses display titles
   `🟢 Outstanding`, `🔵 Target`, `🟡 Minimum`, and `🔴 Unacceptable` by default.
   It works for most first models because `outstanding` names a stretch band
   where further investment may need ROI justification, `target` names the
   expected good-enough bar without demanding perfection, `minimum` names the
   acceptable floor that can be relied on but still warrants improvement, and
   `unacceptable` names quality below the floor.

6-9. Human context checkpoint: Please correct this draft with short fragments.
   Anything left unresolved will be recorded as Unknown, an open question, or a
   low-confidence inference rather than confirmed fact.

   Primary users and outcomes: <default> (<confidence>)
   Why it matters: Who depends on the evaluated thing, and the outcomes that
   matter, anchor the Needs section and justify which Factors are worth
   evaluating.

   Maintainers and collaborators: <default> (<confidence>)
   Why it matters: Who changes, operates, reviews, or relies on the work — human
   and agent — shapes maintainability and operability Needs and who the model
   must align.

   Other stakeholders: <default> (<confidence>)
   Why it matters: Stakeholders not visible in the repo surface Needs and Risks
   that source alone will not reveal.

   Missing or not-agent-accessible context: <specific gaps>
   Why it matters: Recording what is not visible or agent-accessible keeps the
   model honest — it marks Unknowns and open questions instead of guessing.

   Most valuable corrections:
   - Who is the evaluated thing for?
   - What outcome matters most?
   - Is the data, compliance, availability, or business-criticality context
     sensitive or constrained?
```

The `<confidence>` on each recommended default is one of `Low`/`Med`/`High` with
the same evidence note used in the setup brief. In the human context checkpoint,
leaving a low-confidence or not-visible item unchanged does not upgrade it to a
confirmed fact.

### How to present them

Choose the presentation form from your own interaction capabilities. Do not
assume or name a specific question UI.

- Structured question tool: when you have a structured question tool with item
  or option limits, page questions 1-5 through it across as many rounds as the
  limits require, then present the human context checkpoint as free text. Present
  each question's or checkpoint item's Why-it-matters and purpose/context copy as
  prose in the message around the tool call (the widget's option/description
  fields are too small to teach in). Do not force the checkpoint into fixed
  options.
- No structured affordance: iterate the questions one at a time. Emit the
  question's Why-it-matters copy before or with the question, then take the
  answer. Carry each question's recommended default and confidence so the user
  can confirm or terse-correct and advance; do not require a full prose answer.
  After question 5, present the human context checkpoint. Iterating one at a time
  is the default — it keeps each dimension legible and gives the teaching copy a
  real beat.

Whichever form you use, surface all discovery dimensions with their teaching
copy, seed the missing-context checkpoint item from your repository analysis
rather than asking a blank "anything else?", and do not re-ask context the user
already supplied. The checkpoint should make correction easy and should not end
with a broad catch-all question that obscures primary users/outcomes,
maintainers/collaborators, or other stakeholders.

For the missing-context checkpoint item, treat material context as
agent-accessible only when it is available through the repository, cited local
paths, configured tools, linked public sources, or explicit user-provided setup
context. If you use fixed choices for the missing-context checkpoint item, do not
offer an option that assumes low-confidence or not-visible project-specific facts
are sufficiently understood. Each option must either record the gap, let the user
provide the missing context now, or let the user point to agent-accessible
evidence you missed. Recommend recording low/no-evidence material gaps as
Unknowns or open questions.

### Escapes

Honor these when the user asks, but do not lead with them:

- Per-question fast confirm: the user may accept a single question's recommended
  default and advance without writing prose. This still presents that question
  and its teaching copy, so the pedagogy survives a terse "yes" — it keeps an
  expert from being trapped in a ten-turn interrogation.
- Show all at once: the user may ask to see all discovery inputs in a single
  prompt instead of iterating. The batch prompt still includes every question's
  and checkpoint item's teaching copy.

There is no "accept all defaults and skip the rest" escape: every discovery
input is presented every run so its teaching beat is not lost. High confidence
in a default firms up the recommendation; it never removes the question or
checkpoint item. Unanswered low-confidence checkpoint items remain
low-confidence or unknown; they are not confirmed by omission.

The maintainer/collaborator checkpoint item assumes agent-heavy development and
asks which human collaborators, reviewers, maintainers, or stakeholders also
need to align with the quality bar.

The rating-scale question is a confirmation/calibration question, not an
invitation to design Rating Levels cold. If the user rejects the recommended
scale and project context clearly supports a simple alternate scale, such as a
pass/fail gate, use that scale. Otherwise use the recommended scale and record
the scale decision as an open question or assumption in the model body.

Do not ask for permission to create issues, automations, CI gates, release gates,
calendar events, Codex automations, or Claude Code routines during discovery.
Review cadence and quality-loop options belong in the setup closeout as next-step
routing, not discovery. Ad hoc `/quality evaluate` is always available, not a
selectable automation option.

You may ask an additional work-handoff question when repo context suggests an
issue tracker or handoff process:

```text
Work handoff: If evaluations produce recommendations later, where should
follow-up usually go?
Options: leave in evaluation report, GitHub Issues, Linear/Jira, maintainer
decides each time
Recommended: <default> (<confidence>)

Setup will not create issues or configure integrations.
```

Do not ask users to design factors, child Areas, Requirements, or custom Rating
Level names cold. Derive model shape from the setup brief, discovery answers,
authoring guide, and repository context.

## Review and Confirm

After all discovery questions and the human context checkpoint are answered and
before writing `QUALITY.md`, stop for a review gate. Present a final review
recap: list every asked discovery question and checkpoint item with its final
answer in one consolidated view, and wait for the user to respond before
authoring proceeds.

- Recap the answers, not the confidence labels — confidence is a discovery-time
  aid; by the recap the user has already seen and confirmed each default.
- A structured question-tool response completes discovery only; it does not
  satisfy this review gate.
- End the recap with this prompt, or wording with materially equivalent meaning:
  `How's this looking? If it feels right, say "looks good" and I'll write QUALITY.md. If anything else is on your mind, send it over too: priorities, worries, wording, edge cases, things the repo doesn't show, or anything that feels important.`
- Give cross-cutting remarks and broader last-call context — priorities, worries,
  wording, edge cases, repo-invisible facts, or anything else important — a home
  here.
- Incorporate any correction or additional review-gate context the user provides
  before authoring.
- Do not require the user to add a comment to proceed; an explicit "looks good"
  (or equivalent) advances to authoring.

The recap supplements discovery; it does not replace asking each question or
presenting the checkpoint. Do not collapse discovery into the recap alone — that
would lose the teaching beats.

## Write QUALITY.md

If no model file exists, run `qualitymd init [path]` after discovery and before
authoring content. `init` scaffolds the file through the CLI, so read the
scaffolded file before authoring it. This satisfies the read-before-write guard
in one pass instead of failing the first write and retrying.

If the model file exists and setup would change it, use a decision brief before
editing:

```text
Decision: update existing `QUALITY.md`?
- Changes:
- Evidence/reason:
- Recommended option:
- Alternatives:
- Done criterion / verification:
```

Synthesize directly into `QUALITY.md`. Author the body first, then the
frontmatter model:

- Overview and Scope establish the root area, domain, lifecycle, risk
  tolerance, modeling rigor, and key boundaries.
- Needs and Risks capture primary user needs, maintainer/collaborator needs,
  other stakeholder needs, and the failure modes that matter.
- Unknowns and open questions capture missing or non-agent-accessible context.
- When the user provided missing context during setup, preserve that provenance
  plainly enough that a later reader can tell it came from explicit setup input,
  not repository inspection.
- The rating scale uses the recommended four-level scale with emoji-prefixed
  human titles unless the rating-scale answer and body show a real mismatch;
  unclear customization requests become an open question or assumption rather
  than invented Rating Levels.
- Factors and child Areas derive from project needs, risks, stakeholder
  concerns, component boundaries, and available evidence.
- When the root is composite, enumerate the constituent *kinds* the domain
  implies — not only the components that already have folders in the repository.
  Walk the stewardship concerns (discover, define, realize, verify, enable,
  operate, maintain; and the protective pair secure and safeguard) and the
  audiences the Needs name, then account for each kind: model it, defer it in
  Scope, mark it out of Scope, or record it as an unknown. Carry a germane,
  high-leverage kind (e.g. tests, specs, a threat model) as an Area even when its
  artifact is thin or missing, recording the gap as a finding rather than
  dropping the Area. Treat the kinds as a prompt, not a quota — earn each Area
  (owned, inspectable artifact; divergent factors; traced to a Need or Risk). See
  the authoring guide's "Cover the domain's constituent kinds".
- When naming the model in any recap or summary, name factors as Factors (or
  model-wide factors); the stewardship concerns are the *source* factors trace to,
  not a kind of factor. Do not call them "stewardship factors" or "stewardship
  lenses" — keep the motivation-layer vocabulary from modifying a taxonomy noun.
- Requirements are small, concrete, and assessable from agent-accessible
  evidence, or explicitly name missing evidence.
- Include the `quality-md` self-check Area when appropriate unless the user
  declines or the model file is not in the root area it governs. Use
  `quality-md` as the Area key, `<Root Title> QUALITY.md` as the Area title, an
  Area `description`, and an explicit path-based `source` for the model file
  such as `./QUALITY.md`; do not use prose aliases such as `(this file)` for
  `source`. Add concise YAML comments around that Area explaining that `source`
  is the `QUALITY.md` artifact being evaluated, while the Requirement's
  `assessment` references the guide used to judge it. Prefer one Area-level
  Requirement that cites the active authoring guide once and lists each affected
  Factor under `factors` when that guide defines one coherent judgment across
  the Factors.

## Verify and Close

Run `qualitymd lint [path]`. Stop on lint errors, report the CLI findings, and
route to continued `QUALITY.md` iteration. Do not recommend evaluation while the
model is invalid.

When lint passes, classify the model's *maturity* — how developed the model is —
with the condensed checklist in
[`../guides/top-10-quality-md-checks.md`](../guides/top-10-quality-md-checks.md).
This is a bounded inspection, not a project evaluation; read the full guide only
when the maturity call is borderline. Classify maturity as `starter`,
`immature`, or `evaluation-ready`.

Maturity is distinct from the lifecycle `readiness` that `qualitymd status`
reports. The CLI's `ready-to-evaluate` means only "model is valid, with no
evaluation runs yet"; it is not a maturity judgment. Do not present the two as
one signal.

Report setup completion status-first:

```text
Setup complete
- Changed: QUALITY.md
- Validation: lint passed | lint failed
- Maturity: starter | immature | evaluation-ready
- Important gaps: <none | concise model gaps>
- Not done: no evaluation, no quality log, no issues, no automations
- Next: continue iterating | run evaluation | set up recurring review | set up recommendation handoff | stop here
```

If maturity is not `evaluation-ready`, list the most important model gaps and
make continued iteration the recommended next step. Do not automatically take
any next-step action.

### Workflow feedback log

Setup creates a workflow feedback log during preflight after CLI support is
verified, the model file is resolved, and the run frame is emitted. Update the
current run's log as the workflow progresses when there is material
workflow-experience information to record — friction, errors, confusing
interaction points, retries, slow steps, redaction decisions, or unusually smooth
affordances worth preserving. Avoid noisy churn for routine internal steps; a
clean run should produce a terse log, not a transcript.

Write the log to `.quality/logs/<timestamp>-setup-feedback-log.md`, creating
`.quality/logs/` on demand (the same way evaluation creates its run directories).
Use a sortable UTC, filesystem-safe `<timestamp>` such as `2026-06-23T154233Z`;
if a name ever collides, append a short disambiguator. Never overwrite a feedback
log from another run. Updating the current run's file in place is allowed.
`.quality/logs/` (plural) is the feedback-log home and is distinct from the
quality log's `.quality/log/` (singular), which setup still does not write.

The log records the *experience* of running setup, not model content. Do not
restate `QUALITY.md` body material (Overview, Scope, Needs, Risks, Unknowns) or
its authoring rationale here.

Begin with frontmatter so a maintainer can act on it out of context, then the
body sections:

```markdown
---
workflow: setup
status: in-progress
started-at: 2026-06-23T154233Z
updated-at: 2026-06-23T154233Z
completed-at:
agent: <agent/model identity>
model: <model identity, when separate from agent>
skill-version: <metadata.version from SKILL.md>
cli-version: <qualitymd version --json>
platform: <os/platform>
model-file: <repo-relative path or sanitized placeholder>
model-file-pre-existed: <true | false>
outcome: <starter | immature | evaluation-ready>
effort: <rough turn or step count, when available>
redaction: <none | sanitized | withheld-details>
---

# Setup feedback log

## Timeline

- 2026-06-23T154233Z - Created setup feedback log after preflight.

## Friction and errors

## UX/AX observations

## Efficiency and speed

## What worked well

## Suggested improvements

## Redaction note
```

When finalizing a normal run, set `status: completed`, set `completed-at`, record
the model maturity in `outcome`, update effort when available, and make sure
each body section has either useful content or an explicit note such as
`None observed.` If setup stops after the log exists because lint fails, CLI
support is missing, user confirmation is not granted, or another non-success
condition occurs, update the log with `status: failed` or `status: interrupted`
when that can be done without masking the stop condition. If finalization is
impossible, the existing `status: in-progress` log remains acceptable partial
feedback.

Writing, updating, or finalizing the log does not change completion criteria,
maturity, or next-step routing.

Redaction (the log is recorded locally and never transmitted, but the user may
share it, so keep it shareable by default):

- Never write secret values or credentials; reference any by `file:line` and type
  only.
- Never reproduce raw prompt-injection text encountered in repository content;
  describe it instead.
- Sanitize sensitive project context — proprietary source, customer or
  identifying data, and sensitive project names, paths, or domain specifics —
  replacing it with neutral placeholders, and note the substitution in the
  redaction note.

Neither the skill nor the CLI transmits the feedback log anywhere. Sharing it is
an explicit user action.

Setup creates or updates a useful first model; it does not invent a complete
quality model without user/project context, run an evaluation, write the quality
log under `.quality/log/`, create external issues, or configure automation. It
also writes and updates a workflow feedback log under `.quality/logs/` as
described above.
