---
type: Design Doc
title: Setup context checkpoint - design
description: Design for replacing /quality setup's final open-ended discovery questions with a compact human context checkpoint.
tags: [quality-skill, setup, ux]
timestamp: 2026-06-23T00:00:00Z
---

# Setup context checkpoint - design

Design behind [Setup context checkpoint](../0072-setup-context-checkpoint.md)
and its [functional spec](spec.md).

## Context

Setup is prompt-driven by the bundled `/quality` skill. The CLI does not build
or display setup discovery prompts, so the implementation lives in the runtime
setup playbook and the durable setup workflow spec.

Questions 6-9 are not wrong model inputs. They are the right dimensions in a
form that asks too much composition from the user. The design therefore changes
the presentation form, not the model semantics: setup still needs users,
outcomes, maintainers, stakeholders, and unknowns before writing a useful first
`QUALITY.md`.

## Approach

Keep questions 1-5 as discrete discovery questions because each teaches a
specific modeling decision: boundary, domain, lifecycle, risk tolerance, and
Rating Scale. Replace questions 6-9 with one human context checkpoint.

The checkpoint is a compact draft built from the setup brief:

- primary users/outcomes;
- maintainers/collaborators;
- other stakeholders;
- missing or not-agent-accessible context.

The user interaction shifts from "answer these essay questions" to "correct this
draft." The prompt states that silence does not confirm low-confidence facts and
that omitted material gaps will be recorded as Unknown, open questions, or
low-confidence inference. A short bottom section asks for the highest-value
corrections so a user who scans only the end still sees the important asks.

The durable spec names the checkpoint as a discovery input rather than treating
it as a final review shortcut. The existing final review gate remains unchanged:
after discovery, setup still recaps answers and waits for explicit permission
before authoring.

## Alternatives

Keeping four separate questions but shortening the wording was rejected because
the core issue is prompt shape. Four open-ended prompts still make the broad
missing-context question the easiest one to answer while earlier stakeholder
context is skipped.

Moving all human context to the final review gate was rejected because discovery
answers shape the first authored model. The review gate should catch late
context, not become the first place users see these dimensions.

Turning the checkpoint into fixed choices was rejected because the facts are
project-specific and often low-confidence. Fixed choices can make invisible
stakeholder or compliance context look understood when it is not.

Dropping the human context dimensions was rejected because they anchor Needs,
Risks, Factors, and Unknowns. The UX problem is response burden, not relevance.

## Trade-offs & risks

The checkpoint has less per-item rhythm than one-at-a-time questions. The risk is
acceptable because these four dimensions are closely related and the checkpoint
still includes purpose copy and a prioritized correction list.

There is a risk that agents treat an unchanged draft as confirmation. The
runtime and durable spec both guard against that: silence on low-confidence or
not-visible context records Unknowns or low-confidence inference, not confirmed
fact.

## Open questions

None.
