---
type: Runtime Guide
title: Authoring the Markdown body
description: Markdown-body guidance for QUALITY.md judgment context, unknowns, open questions, and review state.
tags: [quality, authoring, guide]
---

# Authoring the Markdown body

Read this when:

- writing or revising the Markdown body;
- reviewing body grounding, unknowns, open questions, or review provenance.

Depends on:

- `../authoring.md`

---

## The Markdown body

The body is evaluable judgment context: what the root area is, why its quality
matters, what decisions the model supports, which needs and risks shaped it, and
what context is missing or inaccessible. It should provide enough concise,
self-explanatory context for a later human or agent to justify the model,
evaluate the model's quality, and decide whether the model still fits the
root area.

A strong body makes its completeness, thoroughness, recency, grounding,
agent-accessibility, and open questions visible instead of implicit.

**Agent-accessible** support is available to the evaluating agent through the
repository, cited local paths, configured tools, linked public sources, or
explicitly provided context. If important support exists but is private,
permission-limited, stale, only known from memory, or unavailable to the agent,
record that limitation in the relevant section's unknowns or open questions.

The body is optional and fixes no required sections; you may rename, reorder, or
replace these. They're recommended starting points:

| Section      | Desired outcome                                                                                                                                                         |
| ------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Overview** | A reader can say what the root area is, who depends on it, and why its quality matters. This names the real entity, not just the repo or file where `QUALITY.md` lives. |
| **Scope**    | A reader can tell what is included, what is excluded for now, and where the model boundary sits. Out-of-scope is not a deficiency.                                      |
| **Needs**    | A reader can see the outcomes the root area must support and the users, operators, maintainers, or downstream systems those outcomes serve.                             |
| **Risks**    | A reader can see the failures that would make the root area untrustworthy, unusable, unsafe, expensive, or hard to change. These are raw material for initial factors.  |

### Shape of a body section

Write each section — including ones you add — to a common shape, so the body
reads consistently as it grows:

1. **Purpose** — open with one line on why this section matters for _this_
   root area, not in the abstract. _If the line would read the same for any
   project, it isn't earning its place._
2. **Contents** — concise, self-explanatory judgment context for this section.
   State the section's conclusion clearly enough that it can be reviewed on its
   own; cite supporting detail instead of copying it; include enough specificity
   to evaluate completeness, thoroughness, recency, and grounding in
   agent-accessible support.
3. **Unknowns & open questions** — captured for every section, scoped to what
   that section covers. An **unknown** is a broad area of uncertainty within the
   section's topic that may not resolve to a single answer; an **open question**
   is sharper — a specific question about that section with a particular answer,
   still unresolved. Both are context that feeds the model, not commentary on it.
   _Write "none known" when there are none, so the absence reads as considered,
   not skipped._
4. **State** — close with the review-provenance line.

### Example body section

This example shows the section shape in use. It is illustrative; adapt the
domain, cited support, unknowns, and open questions to the actual root area.

```markdown
## Needs

Daily support triage quality matters because support leads use this model to
decide whether the inbox is safe to hand off between shifts.

Support leads need urgent customer-impacting messages surfaced before routine
account questions. Agents need enough current policy context to answer without
guessing. Maintainers need triage rules that are inspectable in
`support/policies/triage.md` and reflected in saved views under `support/views/`.

*Unknowns* — holiday launch escalation load is based on last year's notes, which
are not agent-accessible.
*Open questions* — what response-time target should apply to enterprise-contract
escalations?

*Reviewed — Ada Lovelace, 2026-05; agent-reviewed — Codex (GPT-5.5), 2026-06.*
```

### Mark the state of a section

Because the body is largely agent-authored, the freshness signal worth trusting
is not when a section last changed but when a person last stood behind it. The
state line carries two reviews — the last human review (cite the person) and the
last agent review (name the agent surface and model used):

```markdown
## Risks

A regional outage is the failure that would most erode trust: orders silently
drop instead of failing over. Cost overrun is a distant second.

*Unknowns* — failover under a full regional outage is untested.
*Open questions* — should orders fail over to another region, or degrade in place?

*Reviewed — Margaret Hamilton, 2026-05; agent-reviewed — Codex (GPT-5.5), 2026-06.*
```

A section with nothing outstanding still says so:

```markdown
*Unknowns* — none known.
*Open questions* — none.

*Reviewed — Margaret Hamilton, 2026-05; agent-reviewed — Codex (GPT-5.5), 2026-06.*
```

- **Do** capture `Unknowns` and `Open questions` for every section, writing "none
  known" rather than omitting them. _On a high-leverage file, an explicit "none"
  reads as considered; a blank reads as skipped._
- **Do** cite a named person in `reviewed`. _An anonymous review carries no
  accountability; the name is what makes it a trust signal._
- **Do** cite the agent surface and model in `agent-reviewed`, for example
  `Codex (GPT-5.5)`. _The agent name alone is too ambiguous once model behavior
  changes across versions._
- **Do** advance `reviewed` only when that person actually read and endorsed the
  section — never for an agent or mechanical edit.
- **Do** read `agent-reviewed` newer than `reviewed` as the warning state: the
  section has agent changes not yet human-endorsed.
- **Do** treat a missing `reviewed` as **unreviewed** — agent-touched, not yet
  vetted. _Absence is honest; never backfill a name and date a person didn't
  earn._

### Working with the body

- **Do** write the body before expanding the model tree. _It is the fastest way
  to discover what factors and requirements the frontmatter should express._
- **Do** write the body so it can be evaluated for quality in its own right.
  _A later reviewer should be able to judge whether the context is complete
  enough, current enough, specific enough, grounded enough, and accessible enough
  to support the model._
- **Do** cite supporting detail when it materially grounds a section, and flag
  important support that is not agent-accessible. _The body should not become an
  evidence dump, but a later evaluator must be able to tell what the judgment
  rests on and what context could not be inspected._
- **Do** capture, in Needs and Risks, why some requirements matter more than
  others. _Importance and gaps both depend on this context._
- **Do** use the body to explain any rating-scale change. _A custom scale should
  answer a real decision need visible in the body._
- **Do** treat Needs as _benefits to realize_, not only outcomes to protect, and
  Risks as the problems that erode them. _A model that lists only failure modes
  can rate every requirement at target while the root area is still not worth
  relying on — because nothing weighs whether the benefits, on the whole,
  outweigh the residual problems._
- **Consider** separating two reasons a concern matters: how much its failure
  hurts (stakeholder, safety, or business stakes) and how far its failure spreads
  (how much else it forces to change). _A concern can be low-stakes yet
  high-blast-radius, or the reverse; saying which helps the next evaluator weigh
  roll-ups._
- **Consider** noting where two factors or requirements pull against each other
  (tighter access control vs. faster onboarding, latency vs. cost) and which way
  you have chosen to lean. _A model that hides its trade-offs invites an evaluator
  to "fix" a deliberate compromise._

#### Say which sense of "good" this model uses

Quality is not one thing. A root area can be judged by **conformance** (does it
match its specification?), by **fitness for purpose** (does it serve the user's
real need?), or by **value** (is it worth its cost?). These can disagree — a
root area can meet its spec yet fail the need, or serve the need while departing
from spec.

- **Do** name the governing sense of "good" in the Overview, so a reader knows
  whether a passing model means "conforms" or "fits."
- **Do** make both visible where stakeholders would disagree — "meets the spec"
  vs. "meets the need" — rather than letting one silently win. _Different
  stakeholders rate the same finding differently; record the contested
  expectation instead of burying it in a single criterion._
- **Do** attribute the model's judgments to the stakeholders they serve (users,
  operators, maintainers, downstream systems). _A quality no named stakeholder
  would miss is rarely worth modeling._

#### Keep scope, unknowns, open questions, and "not assessed" distinct

- **Do** use **Scope** for concerns outside the model's remit, and a section's
  **Unknowns** for in-scope concerns you've deliberately deferred or cannot yet
  define.
- **Do** record unknowns and open questions under the section they bear on, so
  missing context stays visible without pretending it has already been evaluated.
- **Don't** confuse a declared unknown (your standing declaration) with a **not
  assessed** result (an evaluator's per-run finding that evidence was missing).
- **Do** record low _confidence_ in an assessment, not only its absence. _A
  requirement rated `target` on one stale benchmark or a single reviewer differs
  from one rated on sustained evidence; note that fragility alongside the concern
  it qualifies. "Rated but barely trusted" is neither "not assessed" nor "no
  gap," and it is often where the next evaluation should look first._

---
