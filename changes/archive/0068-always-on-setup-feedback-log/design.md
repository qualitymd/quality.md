---
type: Design Doc
title: Always-on setup feedback log - design doc
description: Design for creating, updating, and finalizing a setup feedback log throughout every /quality setup run.
tags: [skill, setup, logging, feedback]
timestamp: 2026-06-23T00:00:00Z
---

# Always-on setup feedback log - design doc

Design behind the
[Always-on setup feedback log](../0068-always-on-setup-feedback-log.md) and its
[functional spec](spec.md).

## Context

The 0066 feedback log made setup feedback possible but optional and close-only.
That design kept the first version small, but it left a missing file ambiguous:
maybe the run was clean, maybe the run was interrupted, maybe the agent missed
the close step, or maybe the installed skill did not yet know about feedback
logs.

This change keeps the same artifact home and safety boundary while changing the
timing model. The log becomes the setup run's local experience record: created
near the start, updated only when something useful happens, and finalized at
close.

## Approach

Skill-only. No CLI command, flag, config key, parser, or Go package is added.
The skill writes and updates one Markdown file for the current setup run:

```text
.quality/logs/<started-at>-setup-feedback-log.md
```

The workflow order is:

1. Run the normal CLI/version preflight and resolve the model file.
2. Emit the `/quality run` frame, including the feedback-log mutation and
   planned artifact path.
3. Create `.quality/logs/` on demand and write the initial feedback log with
   `status: in-progress`.
4. Update the same file at phase boundaries or when material workflow-experience
   events occur.
5. Finalize the same file at normal close, or mark it failed/interrupted when a
   stop can be recorded without obscuring the user-facing stop reason.

The run frame stays before the first filesystem mutation. That preserves the
existing interaction contract while still making the feedback log always present
for the body of the workflow.

### File updates

The "never overwrite an existing feedback log" rule is interpreted at the run
identity boundary. The skill must not replace another run's file, but it may
rewrite the current run's file in place as its status, timestamps, timeline, and
section summaries change.

The current-run file is small and hand-authored, so whole-file replacement is
acceptable. There is no need for a patch log, lock file, append-only event file,
or CLI seeder. The skill should preserve the section order and keep entries
concise.

### Metadata shape

Frontmatter is for triage: which workflow, which agent/model/tooling, which
platform, what state the run reached, and what redaction posture was used. The
body is for maintainer-readable notes. This split avoids turning prose sections
into ad hoc metadata while still keeping the file useful when copied out of
context.

The field names use explicit lifecycle clocks:

- `started-at`, `updated-at`, `completed-at` rather than a single overloaded
  `timestamp`.
- `status` for run state.
- `outcome` for the setup maturity result.

Blank scalar values are acceptable while the run is in progress. At close, the
skill fills terminal fields it knows and leaves genuinely unavailable fields
blank or summarized.

### Body shape

The added `Timeline` section is the working area for progressive updates. It
keeps phase changes, retries, stops, and finalization notes out of the thematic
sections.

The thematic sections remain summaries:

- `Friction and errors`
- `UX/AX observations`
- `Efficiency and speed`
- `What worked well`
- `Suggested improvements`
- `Redaction note`

At close, empty sections get a short explicit note such as `None observed.` That
makes a clean run visible without making it noisy.

### Stop handling

If setup reaches a stop after the feedback log exists, the skill should update
the file before returning the stop response when that can be done cheaply and
without changing the stop reason. Examples: lint failure, user declines an
existing-file edit, or authoring stops on a missing input.

If setup stops before the feedback log can be created, the existing preflight
stop behavior wins. In particular, missing or unsupported CLI support may prevent
feedback-log creation because setup's hard rule is to stop rather than
hand-author setup artifacts with stale or missing CLI support.

## Alternatives

- **Keep optional close-only logging.** Rejected. It keeps low-value runs out of
  the tree, but a missing log remains ambiguous and interrupted runs lose the
  most useful process feedback.
- **Create only at close, but always.** Rejected. It fixes clean-run ambiguity
  but still loses interrupted-run feedback.
- **Append-only log updates.** Rejected. It would preserve every intermediate
  note, but the artifact is an experience summary, not an audit trail. Whole-file
  updates keep the final file readable.
- **Separate `events.jsonl` plus Markdown summary.** Rejected. That adds a
  machine-readable artifact and coordination burden before there is a parser,
  report, or second consumer.
- **CLI helper for create/update/finalize.** Rejected for now. The CLI would add
  determinism, but the content is still judgmental prose and no other tool reads
  it. A helper can be revisited if feedback logs need validation or multiple
  workflows adopt the same lifecycle.
- **Use evaluation `debug-log.md`.** Rejected again. The feedback log remains a
  central workflow-improvement artifact, not an evaluation run audit.

## Trade-offs & risks

- **More files.** Every setup run creates a log, including clean runs. Mitigation:
  keep clean-run logs terse and local; nothing transmits them automatically.
- **More write points.** Progressive updates widen where setup touches
  `.quality/logs/`. Mitigation: update only the current run's single file and
  only for material workflow-experience events.
- **Transcript drift.** A timeline can become a chat transcript. Mitigation:
  record phase changes and notable workflow experience, not every prompt,
  answer, or model-authoring decision.
- **Preflight edge cases.** Some stops can happen before the log exists.
  Mitigation: create the log immediately after run-frame emission when CLI
  support and model-path metadata are available, and document that earlier
  preflight failures may still stop without a log.
- **Redaction burden.** Always writing a file increases the chance of shareable
  artifacts containing too much context. Mitigation: keep the existing absolute
  prohibitions on secrets and raw prompt-injection text, sanitize sensitive
  project context, and record the redaction posture in frontmatter.

## Open questions

None. The design keeps the 0066 directory and filename, keeps the implementation
skill-only, and settles the current-run update rule: updating the current file is
allowed; overwriting another run's file is not.
