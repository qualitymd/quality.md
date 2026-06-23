---
type: Functional Specification
title: Setup discovery pedagogy
description: Authored per-question background and how-to-change-later copy, ask-every-question, Low/Med/High confidence, and a final review recap for /quality setup discovery.
tags: [skill, setup, ux, pedagogy]
timestamp: 2026-06-23T00:00:00Z
---

# Setup discovery pedagogy

This Change Case spec defines four deltas to the `/quality setup` discovery step:
authored per-question teaching copy in the runtime skill, asking every discovery
question regardless of inferred-default confidence, a `Low`/`Med`/`High`
confidence vocabulary, and a final review recap before authoring. It does not
change which questions are asked, their option sets, their recommended defaults,
the agent-agnostic presentation tiers from 0065, the QUALITY.md format, or any
`qualitymd` CLI behavior.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background / Motivation

The ten discovery questions do double duty: they capture the context the model
needs and they teach the user the dimensions a quality model spans. The current
workflow asserts this purpose but leaves the teaching to per-run agent
improvisation, so it is inconsistent and not human-tunable. Two further frictions
were logged in a real field run (see the
[parent change case](../0067-setup-discovery-pedagogy.md)): the
`strongly inferred`/`weakly inferred`/`assumed` confidence labels read awkwardly,
and once inferences are strong, per-question paging reads as overhead rather than
instruction.

This case resolves that tension in favor of teaching. Setup runs roughly once per
project, so spending interaction to make each dimension legible — and to leave
the user knowing why each answer shapes the model and how to change it later — is
worth more than minimizing round-trips. The teaching value must be authored into
the skill so it is reproducible and tunable, and applied to every question rather
than skipped when confidence is high.

## Scope

Covered:

- Authored per-question teaching copy for all ten discovery questions in the
  runtime setup workflow.
- A workflow-framing statement that setup optimizes for teaching over round-trip
  count.
- Asking every discovery question regardless of inferred-default confidence, and
  revising or removing guidance that contradicts it.
- A `Low`/`Med`/`High` confidence vocabulary, retaining the per-item evidence
  note.
- A final review recap of the full question/answer set with a last-chance comment
  before authoring `QUALITY.md`.

Deferred / non-goals:

- No QUALITY.md format change.
- No `qualitymd` CLI or Go code change.
- No change to which discovery questions are asked, their option sets, or their
  recommended defaults.
- No change to the 0065 agent-agnostic presentation tiers (structured paging vs
  one-at-a-time iteration) beyond the additions above.
- No change to setup's mutation boundary, the `quality-md` self-check Area, or the
  close/maturity step.

## Requirements

### Per-question pedagogy

The setup workflow **MUST** carry authored teaching copy for each of the ten
discovery questions in the runtime skill. For each question, that copy **MUST**
state the purpose of the question — why the dimension matters and what it shapes
in `QUALITY.md` — and **MUST** state how the user can change that answer later.

The teaching copy **MUST** be authored in the workflow itself, not left to per-run
agent improvisation. It **MUST** be written as copy the agent presents to the
user (prose surrounding the question), not as text constrained to a structured
tool's option or description fields.

The workflow **MUST** present a question's purpose and how-to-change-later
context to the user before or together with that question, on whatever
presentation surface the agent uses.

The workflow framing **MUST** state that setup optimizes for teaching the user
the quality-model dimensions over minimizing interaction round-trips, so the
per-question pedagogy is preserved rather than treated as removable overhead.

### Ask every question

The setup workflow **MUST** ask every one of the ten discovery questions on every
run, including questions whose inferred default is high-confidence. High
confidence in an inferred default **MUST NOT** be a reason to skip a question.

The workflow **MUST NOT** offer an escape that accepts all inferred defaults and
skips the remaining questions. Any prior guidance permitting "accept all defaults
and skip the remaining questions" **MUST** be removed or revised so it does not
contradict asking every question. A per-question fast confirm — the user accepts
the recommended default for a single question and advances without writing prose
— **MAY** remain, because it still presents that question and its teaching copy.

The workflow **MUST** continue to honor an explicit user request to see all
questions at once instead of one at a time, and **MUST NOT** lead with that
escape.

### Confidence vocabulary

The confidence vocabulary **MUST** be `Low`, `Med`, and `High`. The workflow
**MUST NOT** use the prior `strongly inferred`, `weakly inferred`, or `assumed`
labels.

Each inferred setup-brief item and each recommended discovery default **MUST**
carry one of `Low`, `Med`, or `High`, plus the short evidence note when evidence
exists. A default with no supporting evidence **MUST** be labeled `Low` and
**SHOULD** name the absence of evidence in its note, preserving the
"no-evidence, pure default" meaning the prior `assumed` label carried.

### Final review recap

After all ten discovery questions are answered and before writing `QUALITY.md`,
the setup workflow **MUST** present a final review recap that lists every
discovery question with its final answer.

The recap **MUST** invite the user to add a last free-text comment or correct any
answer before authoring proceeds. The workflow **MUST** incorporate corrections
the user makes at this step before authoring, and **MUST NOT** require the user to
add a comment to proceed.

The recap **MUST NOT** be the only place a question is surfaced; it supplements,
and does not replace, asking each question during discovery.

## Acceptance Criteria

- Each of the ten discovery questions has authored purpose and
  how-to-change-later copy in the runtime setup workflow.
- The teaching copy is authored in the workflow and presented as prose, not
  confined to structured-tool option fields.
- The workflow framing states that setup favors teaching over minimizing
  round-trips.
- Setup asks every discovery question on every run, including high-confidence
  ones; there is no accept-all-defaults-and-skip escape.
- A per-question fast confirm and the show-all-at-once escape remain available;
  neither escape is led with.
- Confidence is reported as `Low`/`Med`/`High` with the evidence note retained,
  and the prior inferred/assumed vocabulary appears nowhere in the live skill or
  setup spec.
- A no-evidence default is labeled `Low` and names the absence of evidence.
- After discovery and before authoring, setup presents a full question/answer
  recap and invites a last comment or correction, without requiring one.

## Durable spec changes

### To add

None.

### To modify

- `specs/skills/quality-skill/workflows/setup.md` - require authored per-question
  pedagogy (purpose + how-to-change-later) and the teaching-first framing; change
  the discovery prompt-form contract to require asking every question and to
  remove/revise the accept-all-defaults-and-skip escape (retaining show-all-at-once
  and a per-question fast confirm); change the confidence-vocabulary `MUST` from
  `strongly inferred`/`weakly inferred`/`assumed` to `Low`/`Med`/`High` while
  retaining the evidence note; and add the final-review-recap requirement.
- `specs/skills/quality-skill/quality-skill.md` - review the "confidence-labeled
  defaults" framing for consistency; update only if it pins the old vocabulary.
- Relevant OKF logs and indexes under `specs/` - record durable spec updates when
  they land; append-only `log.md` files keep historical confidence-vocabulary
  references frozen.

### To rename

None.

### To delete

None.
