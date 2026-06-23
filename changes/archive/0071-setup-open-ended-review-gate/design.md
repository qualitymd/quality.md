---
type: Design Doc
title: Setup open-ended review gate - design
description: Design for making /quality setup's final review prompt friendly and open-ended.
tags: [quality-skill, setup, ux]
timestamp: 2026-06-23T00:00:00Z
---

# Setup open-ended review gate - design

Design behind
[Setup open-ended review gate](../0071-setup-open-ended-review-gate.md) and its
[functional spec](spec.md).

## Context

The setup workflow is prompt-driven by the bundled `/quality` skill. The CLI
does not construct the final recap or confirmation prompt, so the implementation
lives in the runtime setup playbook and the durable setup workflow spec.

## Approach

Keep the hard review gate from 0069 unchanged: discovery completion is not
permission to author, and setup still waits for the user to respond to the final
recap. Replace the narrow "looks good or corrections" framing with a direct
example prompt in the runtime workflow.

The wording keeps two paths visible:

- a short confirmation path for users who are ready to proceed;
- an open-ended path for late context that may not fit a single discovery
  answer.

In the durable spec, promote the wording's intent rather than only the exact
sentence: the contract is that setup invites broad, useful last-call context and
does not frame all non-confirmation input as error correction. The exact prompt
is still included as recommended wording so runtime agents converge on a shared
tone.

## Alternatives

Only changing `"corrections"` to `"comments"` was rejected because it still reads
like a narrow review action. Setup benefits from making examples concrete enough
that a user knows it is acceptable to add priorities, worries, wording, edge
cases, or missing repo context.

Making the prompt fully free-form without naming `"looks good"` was rejected
because the fast path matters. Users who are done should not have to infer how
to proceed.

Moving this to closeout was rejected because closeout happens after authoring.
The useful context needs to arrive before `QUALITY.md` is written.

## Trade-offs & risks

The prompt is slightly longer than the previous approval line. That is acceptable
because it appears once, at the final gate, and it can materially improve the
first model.

The listed examples are illustrative, not exhaustive. The durable spec keeps the
broader rule so future wording can stay friendly without treating the examples
as a closed set.

## Open questions

None.
