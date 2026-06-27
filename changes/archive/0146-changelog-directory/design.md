---
type: Design Doc
title: Changelog Directory - design
description: How the workspace, skill guidance, durable specs, and dogfood data adopt .quality/changelog/ and flat .quality/logs/.
tags: [skill, workspace, logging, changelog]
timestamp: 2026-06-27T00:00:00Z
---

# Changelog Directory - design

## Context

Answers the [functional spec](spec.md) for change case
[0146](../0146-changelog-directory.md). The change is a naming cleanup with one
small code default and a broad documentation/spec footprint.

## Approach

### Rename the model-change directory at the workspace boundary

Update the single code default, `DefaultQualityLogDir`, from `.quality/log` to
`.quality/changelog`. The `Workspace.Log` field can remain in code for this
change because it is an internal path slot and no public command or JSON field
currently exposes a `qualityLogDir` name. Public-facing text should use
"quality changelog" and `.quality/changelog/`.

### Keep feedback logs as flat workflow logs

Leave `DefaultFeedbackLogDir` as `.quality/logs`. The durable feedback-log
contract already requires filenames like
`<timestamp>-<workflow>-feedback-log.md`; the change reframes that directory as
the flat workflow-log home and clarifies that future log kinds should use
descriptive filenames rather than nested subdirectories.

### Update active guidance, not historical records

Update active specs, runtime skill files, current docs, code comments, examples,
and dogfood data. Leave archived Change Cases, historical bundle log entries,
and older root changelog entries intact unless they are being extended with a
new current entry. Historical records should preserve the path that was true
when they were written.

### Rename dogfood data with a timestamp

Move `.quality/log/2026-06-22-add-quality-md-area.md` to
`.quality/changelog/2026-06-22T000000Z-add-quality-md-area.md`. The exact time
is not recoverable from the old date-only filename, so midnight UTC is the least
surprising deterministic preservation value for an existing historical entry.
New entries will use run/write timestamps.

### No compatibility layer

Do not add fallback readers or dual writes. This is an early-alpha contract
rename; tests and docs should assert the new current path rather than preserving
the old path as live behavior.

## Spec response

- **Model-change changelog** - satisfied by changing the workspace default,
  renaming dogfood data, and updating the quality changelog spec/runtime
  guidance.
- **Workflow logs** - satisfied by leaving `.quality/logs/` flat and clarifying
  filename-based log kinds in the shared feedback-log contract.
- **Clean break** - satisfied by avoiding compatibility code and limiting old
  path references to historical records.
- **Verification** - satisfied by source inspection, targeted searches, Go tests,
  and Markdown formatting.

## Alternatives

- **Use `.quality/changelog.md`.** Rejected. A single file is easy to browse but
  creates avoidable append conflicts and makes agent writes harder to keep
  scoped.
- **Use `.quality/model-changes/`.** Rejected. It is precise, but "changelog" is
  more recognizable to users and matches the artifact's purpose.
- **Put feedback logs under `.quality/logs/feedback/`.** Rejected. Feedback is
  one current log kind; subfolders add hierarchy before the category needs it
  and make chronological scanning harder.
- **Keep `.quality/log/` and improve docs.** Rejected. The confusion is in the
  path itself, not only in the explanation.

## Trade-offs & risks

- Existing local workspaces with `.quality/log/` will need to rename the folder
  manually if they want current layout. That is acceptable during early alpha and
  avoids making agents reason about two live homes.
- The code keeps the internal `Workspace.Log` field name for now. That avoids a
  wider internal rename with no public payoff, but future CLI surfaces should
  choose `Changelog` naming before exposing the value.

## Open questions

None.
