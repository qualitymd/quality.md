---
type: Design Doc
title: Planned coverage status design
description: Design for optional planned coverage metadata and status comparison.
tags: [evaluation, status, cli]
timestamp: 2026-06-18T00:00:00Z
---

# Planned coverage status design

## Context

The [functional spec](spec.md) answers the interruption/resume failure found in
the E11 artifact-write experiments: after a partial run, `show-status` can tell
that analysis is missing, but it cannot know which intended assessments were
never written. The scratch planned-coverage prototype showed that an optional
manifest can close that gap with a set comparison.

This design keeps `design.md` and `plan.md` as the human-readable evaluation
scope record, and adds a small CLI-owned metadata artifact only when a run needs
machine-checkable resume diagnostics.

## Approach

Add an optional `planned-coverage.json` artifact at the evaluation run root.
When the file is absent, `show-status` uses the existing reportability checks
unchanged.

The artifact is written through a CLI command rather than hand-authored in the
run folder. The expected command shape is:

```sh
qualitymd evaluation set-planned-coverage <run> --file <path>
```

The command validates and rewrites the JSON into canonical form. It does not
infer coverage from the model; the skill supplies the intended coverage after it
has resolved scope, effort, and limitations in `design.md` / `plan.md`.

The stored shape should be intentionally small:

```json
{
  "schemaVersion": 1,
  "assessments": [
    {
      "targetPath": ["Target", "Factor"],
      "requirement": "Requirement text"
    }
  ],
  "analyses": [
    {
      "targetPath": ["Target"]
    }
  ]
}
```

`show-status` loads the artifact, normalizes each planned entry into a stable
key, and compares it with written records:

- assessment keys: ordered `targetPath` plus `requirement`
- analysis keys: ordered `targetPath`

Missing planned entries become renderability gaps. Unexpected assessment or
analysis records outside the plan are also reported as gaps for the initial
implementation because a planned run with unplanned records is ambiguous. The
gap detail should identify whether the record is missing or unexpected and
include the target path and requirement when present.

Gap output should be sorted by kind, target path, and requirement so repeated
runs and tests are deterministic. `build-report` needs no separate planned
coverage logic because it already refuses runs with renderability gaps.

The quality skill should write planned coverage only when it materially helps
resume diagnostics, especially for standard, deep, concurrent-write, or
interruption-prone runs. Quick scoped runs may keep relying on `plan.md` alone.

## Alternatives

- Parse `plan.md` prose for intended coverage. Rejected because section labels
  and prose wording are useful for humans but too brittle for reportability.
- Require planned coverage for every run. Rejected because quick scoped
  evaluations should stay lightweight and current runs must keep working.
- Ask the skill to compare planned coverage manually. Rejected because resume
  diagnostics belong in deterministic status output, not prompt discipline.
- Add planned fields to every assessment and analysis record. Rejected because
  it changes the record payload contract and makes a run-level concern repeat
  across records.

## Trade-offs and Risks

This adds a second machine-readable artifact to evaluation runs. The cost is
worth it only if the artifact remains optional and small.

The artifact can drift from `plan.md` if the skill changes scope after writing
it. The skill prompt should write planned coverage after the plan is settled and
rewrite it when the plan changes before records are written.

An incomplete planned-coverage file can create false confidence. The design
reduces that risk by treating the artifact as an explicit plan supplied by the
skill, not as proof of whole-model coverage.

Unexpected-record gaps may be too strict for some correction workflows. If that
blocks legitimate use, a later replacement/superseding change can distinguish
intentional corrections from unplanned records.

## Open Questions

- Should the writer command be named `set-planned-coverage`, or should it be a
  flag on `create-run` once the run plan is known?
- Should unexpected records be hard reportability gaps or warnings in the
  durable status spec?
- Should `report.json` include planned coverage summaries later, or should this
  remain status-only?
