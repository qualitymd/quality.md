---
type: Design Doc
title: Setup missing-context provenance — design
description: Design for tightening /quality setup missing-context options without CLI changes.
tags: [quality-skill, setup, agent-accessibility]
timestamp: 2026-06-23T00:00:00Z
---

# Setup missing-context provenance — design

Design behind
[Setup missing-context provenance](../0070-setup-missing-context-provenance.md)
and its [functional spec](spec.md).

## Context

The setup workflow is prompt-driven by the bundled `/quality` skill. The CLI
does not construct discovery questions, so the implementation lives in the
runtime setup playbook and the durable setup workflow spec.

## Approach

Keep the nine-question setup shape unchanged. Add a targeted rule to the
missing-context discovery guidance: when setup presents fixed choices or a
recommended default for identified gaps, each choice must preserve provenance.
The acceptable outcomes are recording the gap, receiving explicit context from
the user, or being pointed to missed accessible evidence.

In the durable spec, place the requirement directly after the seeded
missing-context-question rule, where future changes to discovery choices are
most likely to look. In the runtime workflow, mirror it in the discovery
presentation section and in the write step so the authored body records
setup-provided provenance.

## Alternatives

One alternative was to change the authoring guide's definition of
agent-accessible support. That definition is already correct: explicitly
provided context is accessible, while private, stale, memory-only, or unavailable
support is not. The failure was setup option wording, not the underlying
authoring contract.

Another alternative was to require every missing-context question to be free
text. That would avoid bad fixed choices but lose the useful recommended-default
shape for setup runs where agents have a structured question surface. The
problem is not fixed choices themselves; it is choices that erase provenance.

## Trade-offs & risks

The rule adds another guardrail to an already detailed setup workflow, but it is
localized to the missing-context question. Agents still retain room to tailor
options to the project as long as they do not convert low/no-evidence facts into
assumptions.

The change relies on transcript review and prompt conformance rather than a
deterministic CLI test because the behavior is currently skill-only.

## Open questions

None.
