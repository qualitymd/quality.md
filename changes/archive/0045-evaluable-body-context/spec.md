---
type: Functional Specification
title: Evaluable body context - functional spec
description: Requirements for treating the QUALITY.md Markdown body as concise, self-explanatory, agent-accessible judgment context for evaluating model quality.
tags: [authoring, body, skill, accessibility]
timestamp: 2026-06-21T00:00:00Z
---

# Evaluable body context - functional spec

Companion to [Evaluable body context](../0045-evaluable-body-context.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background / Motivation

The Markdown body is the durable judgment context behind a QUALITY.md model. It
does more than help a reader interpret frontmatter: it records what the subject
is, why quality matters, what decisions the model supports, what needs and risks
shaped the model, and which support is missing or inaccessible. That context is
also evidence for judging the model itself. A later reviewer should be able to
ask whether the body is complete, current, specific, grounded, and accessible
enough to justify and evaluate the model.

Because QUALITY.md work is agent-first, "accessible" must be explicit:
supporting context needs to be agent-accessible, or its absence needs to be
visible. A section can be concise and still rigorous, but it is too terse if a
later human or agent cannot evaluate the quality of the context it provides.

## Scope

Covered: the authoring guide's Markdown-body purpose statement, body-section
shape, support-citation guidance, progressive-disclosure guidance, and treatment
of non-agent-accessible support as an unknown or open question. Also covered are
the downstream guide/check surfaces that route authors and agents through that
body-context guidance.

Deferred / non-goals: no new required body sections; no required `Access gaps`
line in every body section; no frontmatter schema change; no lint rule; no CLI
parser behavior for support citations, recency, or agent-accessibility language.

## Requirements

### Body as evaluable judgment context

Authoring guidance **MUST** describe the Markdown body as evaluable judgment
context: context for building the model, understanding the model's purpose,
using the model in evaluation, evaluating the model's quality, and deciding
whether the model still fits the subject.

The guidance **MUST NOT** frame the body only as context for interpreting the
model.

> Rationale: a body that only explains the model can look sufficient while still
> failing to show whether the model is complete, current, grounded, or fit for
> purpose.

### Body quality dimensions

Authoring guidance **MUST** teach that body context should be written so a later
human or agent can evaluate the body's own quality, including completeness,
thoroughness, recency, subject specificity, grounding, agent-accessibility, and
open questions.

The guide **SHOULD** describe this as an evaluable standard for body sections,
not merely as a style preference.

### Concise, self-explanatory sections

Authoring guidance **MUST** teach that each body section should be concise,
rigorous, and self-explanatory. A section should state its conclusion clearly
enough to be reviewed on its own, but should cite supporting detail rather than
copying or summarizing large supporting sources.

The guidance **MUST** preserve progressive disclosure: the body carries the
judgment context and citations; supporting evidence, examples, raw data, and
long rationale live in cited sources when those sources are available.

### Agent-accessible support

Authoring guidance **MUST** use the term **agent-accessible** and define it.
Agent-accessible support is support available to the evaluating agent through
the repository, cited local paths, configured tools, linked public sources, or
explicitly provided context.

The guide **MUST** instruct authors to cite supporting sources that materially
ground a body section when those sources are agent-accessible.

### Inaccessible support as first-class context

Authoring guidance **MUST** instruct authors to record material supporting
context that is not agent-accessible as a first-class limitation in the relevant
body section. The limitation should be recorded under that section's unknowns or
open questions unless the section uses a clearer local form.

Examples of non-agent-accessible support include private dashboards,
permission-limited documents, uncited stakeholder memory, stale sources, or
sources known to exist but unavailable to the evaluating agent.

The guide **MUST NOT** require a separate `Access gaps` line in every section.

> Rationale: support accessibility bears on the section it qualifies. Folding it
> into scoped unknowns/open questions keeps the body rigorous without turning
> every section into a procedural checklist.

### Model-quality inspection

The quick model-quality checks and getting-started guidance **SHOULD** treat
missing or non-agent-accessible support as a model-usefulness finding when it
prevents a reader or agent from evaluating whether the body context is complete,
current, grounded, or sufficient.

The guidance **MUST** keep that finding distinct from subject quality. Missing
agent-accessible support weakens the model or body context; it is not, by
itself, evidence that the subject is low quality.

## Durable spec changes

### To add

None

### To modify

- `SPECIFICATION.md` — broaden body semantics from context for interpreting the
  Model to context for building, using, and evaluating the Model (per the
  body-as-evaluable-judgment-context requirement above).
- `specs/skills/quality-skill/guides/authoring-md.md` — require the runtime
  authoring guide to teach evaluable body context, body-quality dimensions,
  concise self-explanatory sections, agent-accessible support, and inaccessible
  support limitations (per all authoring-guidance requirements above).
- `specs/skills/quality-skill/guides/getting-started-md.md` — align first-pass body
  checks with the body-quality and agent-accessible-support requirements above.
- `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md` — add model
  inspection treatment for missing or non-agent-accessible support (per the
  model-quality-inspection requirement above).
- `specs/skills/quality-skill/quality-skill.md` — update only if setup or wizard
  routing needs explicit agent-accessibility language (per the model-quality
  inspection requirement above).
- `specs/cli/init.md` — update only if the starter scaffold text changes to
  mention support accessibility (per the inaccessible-support requirement above).

### To delete

None
