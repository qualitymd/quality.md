---
type: Design Doc
title: Run frame title and workflow vocabulary — design
description: Implementation approach for retitling the /quality run frame and retiring "mode" in favor of "workflow".
tags: [skill, quality, ux, agent-mediated]
timestamp: 2026-06-26T00:00:00Z
---

# Run frame title and workflow vocabulary — design

Companion to the
[Run frame title and workflow vocabulary](../0110-run-frame-and-workflow-vocabulary.md)
change case and its [functional spec](spec.md).

## Context

This is a documentation-and-runtime-skill alignment change. There is no Go or
CLI behavior under design — the artifact under change is the agent's rendered
interaction surface (the run-frame header) and the vocabulary the durable specs
and runtime skill use to name the public workflows.

Two coupled edits: a localized retitle of the run frame, and a repository-wide
term swap from "mode" to "workflow" wherever "mode" means a `/quality` workflow.

## Approach

### Run frame header

Replace the rendered header `**/quality run**` with a title that names the
workflow and cannot be mistaken for a command:

```text
**Quality · <workflow>**
- **Model file:**
- **Scope:**
- **Rigor:**        (when applicable)
- **Mutation:**      (read-only, evaluation artifacts, evaluated source, QUALITY.md, quality log, feedback log, tooling)
- **Artifacts:**
- **Next gate:**
```

`<workflow>` is the resolved workflow name — `setup`, `evaluate`, or `update`.
The concrete renders become `**Quality · setup**`, `**Quality · evaluate**`, and
`**Quality · update**`. The `- **Mode:** <x>` line is removed from the template
and from all three concrete frames; the workflow name it carried now lives in the
header.

Everything else in the frame is unchanged. The read-only-vs-mutating distinction
stays in the `Mutation:` line, which already names the class of thing that may
change — so the spec's "MUST distinguish read-only from mutating" requirement is
still satisfied without a redundant title tag.

The separator is the middot (`·`): it visually fuses "Quality" and the workflow
into one title rather than a `label: value` pair, while staying clearly not a
slash command.

### Workflow vocabulary

Swap "mode" → "workflow" everywhere "mode" names a `/quality` workflow:

- runtime skill: `SKILL.md` Arguments and Workflow Dispatch wording, the
  `# Evaluate Mode` / `# Update Mode` H1s, and the recommendation-follow-up
  guide's "not a `/quality` mode".
- durable specs under `specs/skills/quality-skill/`: the Run frames section, the
  "public mode" / "cross-mode" / "non-mode" / "evaluation mode" phrasings, and
  the `mode` frontmatter tags on the workflow specs.
- durable docs: `docs/reference/versioning.md` ("The mode" → "The workflow") and
  `docs/guides/cut-a-release.md` ("skill modes" → "skill workflows").

Unrelated senses of "mode" are left untouched: failure modes, CLI output/JSON
modes, Diátaxis documentation modes, and the Go OS file mode. The sweep in the
change case enumerated these so the term swap can be applied by reading each hit
in context rather than blind find-and-replace.

The internal construct name "run frame" is kept, and the `### Run frames`
heading (anchor `#run-frames`) is left as-is so the append-only
`quality-log.md` reference does not break.

### Durable spec requirement

Add to the Run frames section of `quality-skill.md` the new constraints: the
header MUST name the resolved workflow, MUST NOT render `/quality run` or any
command-style header for a non-invocation, and the frame MUST NOT use a `Mode:`
field label. This sits alongside 0062's existing "MUST NOT emit `Mode: wizard`"
constraint and 0038's run-frame rationale, carrying the 0110 lesson forward in
the durable spec rather than only in the archived case.

### Sequence

1. Patch the durable run-frame contract in `quality-skill.md` (header rules +
   workflow vocabulary).
2. Patch the runtime `SKILL.md` template, Arguments, and Workflow Dispatch.
3. Patch the three workflow files' rendered frames and H1s.
4. Apply the term swap across the remaining spec and doc files.
5. Update append-only skill/spec logs for the durable and runtime files touched.
6. Search for residual `/quality run`, `**Mode:**`, and workflow-sense "mode" in
   the live surface; run the Markdown formatting check.

## Spec response

- [Run frame header](spec.md#run-frame-header) is satisfied by the retitled
  template in `SKILL.md`, the three concrete frames, and the new header
  constraints in `quality-skill.md`.
- [Workflow vocabulary](spec.md#workflow-vocabulary) is satisfied by the term
  swap across the runtime skill, durable specs, and the two docs, scoped to the
  workflow sense of "mode".

## Alternatives

- **Keep a header, just drop the slash (`Quality run`).** Rejected: "run" still
  implies a `run` workflow that does not exist. Naming the actual workflow is
  clearer.
- **Title tag for read-only (`Quality · evaluate — read-only`).** Rejected as
  redundant with the `Mutation:` line, which already carries the read-only vs.
  mutating class. Keeping the title to `Quality · <workflow>` avoids two sources
  of truth for the same fact.
- **Colon separator (`Quality: evaluate`).** Rejected: `Quality:` reads as a
  field label with a value, the same `label: value` shape the change is trying to
  move away from. The middot reads as a title.
- **Rename the "run frame" construct too.** Rejected: it is instructional, never
  user-rendered, and renaming breaks the `#run-frames` anchor referenced by
  append-only history. The leak is the rendered header, not the internal term.
- **Leave "mode" as an internal-only synonym.** Rejected per the user's
  direction to unify vocabulary; one concept, one name, removes the standing
  "Mode/workflow" ambiguity in `SKILL.md`.

## Trade-offs & risks

- The term swap touches many files, so the main risk is catching an unrelated
  "mode". The change case sweep already bucketed every hit; each edit is applied
  by reading context, not global replace.
- Renaming the workflow H1s (`# Evaluate Mode` → `# Evaluate Workflow`) changes
  in-file anchors, but those headings are not cross-referenced by anchor from
  other files (unlike `#run-frames`), so no inbound links break.

## Open questions

None.
