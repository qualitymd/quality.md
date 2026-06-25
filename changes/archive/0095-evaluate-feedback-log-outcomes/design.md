---
type: Design Doc
title: Evaluate feedback log outcomes - design
description: Design for keeping evaluate feedback logs outside evaluation runs and making outcome values describe workflow process state.
tags: [skill, evaluation, logging, feedback]
timestamp: 2026-06-25T00:00:00Z
---

# Evaluate feedback log outcomes - design

Answers the [functional spec](spec.md). This is a guidance and durable-spec
mirror change only.

## Context

Evaluation feedback logging already has the right artifact boundary:

- formal evaluation judgment lives in Evaluation v2 data and generated reports;
- historical `debug-log.md` is legacy run-folder compatibility only;
- workflow feedback lives under `.quality/logs/` and is local, hand-authored, and
  non-authoritative.

The remaining weakness is the `outcome` field. Values like `reported` and
phrases like "report outcome" are easy to read as evaluation semantics. The field
should instead classify the workflow's terminal process state.

## Approach

Keep the existing artifact location and lifecycle. Change only the
workflow-specific `outcome` vocabulary and surrounding wording:

- Runtime `evaluate.md` frontmatter example gets explicit allowed values.
- Runtime finalization says to record the workflow outcome.
- Durable evaluate feedback-log spec owns the terminal value list.
- Parent summaries and logs keep the location/boundary story intact.

The value set is intentionally small and operational:

```text
completed-reportable
stopped-lint
stopped-model
stopped-source
stopped-tooling
failed
interrupted
```

These values are enough to improve the workflow while avoiding a second report
channel.

## Alternatives

- **Use generic `reported | stopped | failed | interrupted`.** Rejected because
  `reported` sounds report-semantic and `stopped` hides the most useful stop
  category.
- **Put the log back in the evaluation run folder.** Rejected because workflow
  feedback should not sit beside formal evidence and reports where it can be
  mistaken for part of the evaluation record.
- **Make the CLI own feedback logs.** Rejected. These are hand-authored workflow
  observations; the CLI owns deterministic mechanics and structured evaluation
  artifacts.

## Trade-offs & risks

- The value list may need expansion after more real runs. Keep that as a future
  compatibility decision; do not invent categories now.
- The log remains non-authoritative, so it cannot be used to recover missing
  evaluation state. That is deliberate.

## Open questions

None.
