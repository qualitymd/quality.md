---
type: Functional Specification
title: Run frame title and workflow vocabulary — functional spec
description: Requirements for retitling the /quality run frame away from a fake command and Mode label, and unifying the public surface on "workflow".
tags: [skill, quality, ux, agent-mediated]
timestamp: 2026-06-26T00:00:00Z
---

# Run frame title and workflow vocabulary — functional spec

Companion to the
[Run frame title and workflow vocabulary](../0110-run-frame-and-workflow-vocabulary.md)
change case. This spec states what the change must do; the
[design doc](design.md) covers how.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119 and RFC 8174 when, and only when, they
appear in all capitals.

## Background / Motivation

The run frame is a status-first preamble that lets the user catch a wrong
inference before the skill spends effort or mutates anything (the 0038
rationale). It earns its place. But its header leaks two pieces of vocabulary the
public surface does not have:

- The header renders literally as `**/quality run**`. There is no `run`
  invocation; the real ones are `/quality`, `/quality setup`, `/quality
  evaluate`, and `/quality update`. A command-style header advertises an
  invocation that does not exist — the same leak the skill already forbids for
  `status`, `next`, and `review`.
- The frame's first field is `Mode:`. "Mode" is an internal name for what the
  public surface elsewhere calls a workflow; the skill and spec use both names
  inconsistently. Surfacing the internal name in the rendered frame, as a
  field label that reads like an argument, compounds the confusion.

This change retitles the frame so its header names the workflow without looking
like a command and without a `Mode:` field, and it settles the vocabulary on
"workflow" so the public surface is described one way.

## Scope

Covered: the `/quality` run-frame header and field set, and the
public-surface concept name across the durable skill specs, bundled runtime
skill, and the durable docs that name it.

Not covered: which workflows exist, how invocations resolve, the run frame's
other required fields, the internal term "run frame" itself, CLI behavior, Go
implementation, the QUALITY.md format, evaluation records, reports, or rating
semantics. Unrelated senses of "mode" (failure modes, CLI output modes,
Diátaxis documentation modes, OS file modes) are out of scope.

## Assumptions & dependencies

- The run-frame requirement and rationale established by 0038 (and the
  wizard-removal constraint from 0062, which already forbids emitting `Mode:
  wizard`) remain in force; this case extends that line of constraints.
- The durable `/quality` skill specs under `specs/skills/quality-skill/` mirror
  runtime skill behavior closely enough that a contract change must update both
  durable specs and bundled runtime guidance.

## Requirements

### Run frame header

- The run frame header **MUST** identify the resolved workflow.
  > Rationale: the header is the strongest visual element of the frame; naming
  > the workflow there lets the field set carry only the variable run details. —
  > 0110

- The run frame header **MUST NOT** render a string that reads as an invokable
  `/quality` command for a token that is not a real invocation. In particular it
  **MUST NOT** render `/quality run`.
  > Rationale: a command-style header invites the user to type a command that
  > does not exist, the same leak the skill forbids for `status`, `next`, and
  > `review`. — 0110

- The run frame **MUST NOT** use a `Mode:` field label.
  > Rationale: "Mode" is internal vocabulary the public surface does not use, and
  > a field label reads as a settable argument. The workflow name belongs in the
  > header, not in a `Mode:` field. This extends 0062's constraint that the run
  > frame must not emit `Mode: wizard`. — 0110

- The run frame **MUST** continue to name the model file, scope, rigor level when
  applicable, mutation policy, expected artifacts, and next user-visible gate,
  and **MUST** continue to distinguish read-only work from mutating work and name
  the class of thing that may be changed.
  > Rationale: removing the `Mode:` label is the only field change; the frame's
  > informational contract from 0038 is otherwise unchanged. The read-only vs.
  > mutating distinction stays carried by the mutation field. — 0110

### Workflow vocabulary

- The durable `/quality` skill specs and bundled runtime skill **MUST** name the
  public-surface concept "workflow", not "mode", when referring to `setup`,
  `evaluate`, `update`, or the set of them. Recommendation follow-up **MUST** be
  described as "not a workflow" rather than "not a mode".
  > Rationale: one concept should carry one name. The public surface already
  > routes through `workflows/` and lists workflows; "mode" is the inconsistent
  > second name. — 0110

- This rename **MUST NOT** alter unrelated uses of "mode" (failure modes, CLI
  output/JSON modes, Diátaxis documentation modes, OS file modes).
  > Rationale: "mode" is a correct word in those senses; only its use as a name
  > for a `/quality` workflow is being retired. — 0110

- The internal term "run frame" **MUST** be retained as the name of the construct
  in specs and runtime guidance.
  > Rationale: "run frame" is instructional and never rendered to the user, and
  > the `#run-frames` spec anchor is referenced by append-only history; renaming
  > it would break that anchor for no user-facing benefit. — 0110

## Durable spec changes

### To add

None.

### To modify

- [`specs/skills/quality-skill/quality-skill.md`](../../../../specs/skills/quality-skill/quality-skill.md)
  - update the Run frames section to name the resolved *workflow*, forbid a
    command-style header and a `Mode:` field label, and retire "mode" as the
    public-surface concept across the shared contract. Driven by
    [Run frame header](#run-frame-header) and
    [Workflow vocabulary](#workflow-vocabulary).
- [`specs/skills/quality-skill/index.md`](../../../../specs/skills/quality-skill/index.md),
  [`specs/skills/quality-skill/evaluation.md`](../../../../specs/skills/quality-skill/evaluation.md),
  [`specs/skills/quality-skill/recommendation-follow-up.md`](../../../../specs/skills/quality-skill/recommendation-follow-up.md),
  [`specs/skills/quality-skill/workflows/evaluate.md`](../../../../specs/skills/quality-skill/workflows/evaluate.md),
  [`specs/skills/quality-skill/workflows/setup.md`](../../../../specs/skills/quality-skill/workflows/setup.md),
  [`specs/skills/quality-skill/workflows/update.md`](../../../../specs/skills/quality-skill/workflows/update.md),
  [`specs/skills/quality-skill/workflows/index.md`](../../../../specs/skills/quality-skill/workflows/index.md),
  [`specs/skills/quality-skill/guides/recommendation-follow-up-md.md`](../../../../specs/skills/quality-skill/guides/recommendation-follow-up-md.md)
  - retire "mode" as a name for a `/quality` workflow. Driven by
    [Workflow vocabulary](#workflow-vocabulary).

### To rename

None. The "run frame" construct name and the `#run-frames` anchor are retained.

### To delete

None.

## Validation check

If every requirement above is satisfied, the run frame will identify its workflow
in the header without rendering a fake command or a `Mode:` field, the frame's
other fields and its read-only/mutating distinction will be unchanged, and the
public surface will be described consistently as a set of workflows. That
achieves the motivation without changing which workflows exist, how invocations
resolve, the CLI, or the QUALITY.md format.
