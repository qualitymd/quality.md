---
type: Runtime Resource
title: Output policy
description: Policy for consuming and summarizing command output in /quality workflows.
---

# Output policy

Use CLI output in the form that matches the task.

## Prefer machine-readable output for agent consumption

When a command offers `--json`, use it for data the agent must inspect or carry
forward. Do not parse styled or tabular human output when a JSON form exists.

Use JSON for:

- run creation receipts;
- data write receipts;
- status or report data when the command offers JSON; and
- any future command that emits structured findings or summaries.

## Use human output for the user-facing closeout

Human output is appropriate for short summaries to the user. Keep it concise and
name the important generated paths, ratings, failures, or next actions.

## Treat verbatim artifacts as artifacts

`qualitymd spec`, `qualitymd evaluation data get`, and
`qualitymd evaluation data example <kind>` emit verbatim artifacts. Treat those
outputs as artifacts, not as JSON receipts. Do not add `--json` to commands whose
stdout is already the JSON artifact.

## Keep stdout and stderr separate

Preserve the CLI contract:

- stdout is payload;
- stderr is diagnostics, progress, and human next-action footers; and
- redirected or piped stdout must remain safe to consume.

## Do not infer success from prose

Use exit status and structured receipts/status where available. If a command
fails, report the command and the diagnostic; do not continue by guessing the
missing state.
