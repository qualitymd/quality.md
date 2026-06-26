# /quality Runtime Skill Update Log

## 2026-06-26

- **Revision**: Updated [`SKILL.md`](SKILL.md) for 0131 - Area findings in
  evaluation reports.
  Runtime evaluation guidance now tells Area analysis to synthesize
  `AreaAnalysisResult.findings` for material Area/Factor report observations and
  keeps recommendations, impact, priority, effort, benefit, ROI, actions, and
  global top-finding rankings out of Area Findings.

- **Revision**: Updated [`SKILL.md`](SKILL.md) for 0130 - Self-contained
  per-kind data schema.
  Artifact Contract guidance now treats `qualitymd evaluation data schema <kind>`
  as the source for required fields and allowed enum values, with
  `data example <kind>` as one concrete instance and `data set --dry-run` as
  authored-payload validation only.

- **Revision**: Updated [`SKILL.md`](SKILL.md) for 0129 - Evaluation
  orchestration overhaul.
  Runtime guidance now removes the evaluation rigor argument and run-frame field,
  makes exhaustive in-scope coverage the only evaluate workflow, defaults
  collection/QC to subagent fan-out when available, and requires the two-pronged
  QC phase before roll-up.

- **Revision**: Updated [`SKILL.md`](SKILL.md) for 0128 - Agent-mediated skill
  alignment.
  Runtime guidance now gives read-only orientation a status-first shape and
  requires recommendation follow-up to emit a frame before inspection, outcome
  selection, or mutation.

- **Revision**: Reshaped the [`SKILL.md`](SKILL.md) User Interaction Contract for
  0121 - Scannable interaction hierarchy.
  The decision-brief template now leads with the question, renders choices as a
  separated `[y]`/`[n]` block with the alternative folded in, and demotes a capped
  set of `label:` rationale lines; the contract prose now requires hierarchy by
  position rather than bold alone, surviving bold-stripping.

- **Revision**: Retitled the run-frame header label from `Quality` to
  `QUALITY.md`.
  The run-frame template header is now `**QUALITY.md · <workflow>**` across the
  dispatcher and the setup/evaluate/update workflows.

- **Revision**: Updated the Evaluation data contract guidance for 0115 -
  Type-safe, model-bound Evaluation v2 data.
  The dispatcher now points agents to `qualitymd evaluation data schema [<kind>]`
  as the authoritative payload-shape source, examples as populated samples, and
  `data set --dry-run` as authored-payload validation.

- **Revision**: Renamed live Evaluation wording for 0116 - Drop the "Evaluation
  v2" naming.
  Runtime skill prose now uses plain "Evaluation" for the active workflow and
  data surface.

- **Revision**: Tightened the run-frame instruction for 0114 - Run frame as
  first output.
  The dispatcher now requires the run frame as the workflow's first output before
  any tool call, forbids gating emission on a tool result, and allows a
  best-known or `resolving…` value for a field (such as a many-Area scope) that
  needs a tool to resolve.

- **Revision**: Retitled the run frame and unified workflow vocabulary for 0110 -
  Run frame title and workflow vocabulary.
  The run-frame template header is now `**Quality · <workflow>**` instead of
  `**/quality run**`, the `Mode:` field is dropped (the workflow name moves into
  the header), and Arguments/Workflow Dispatch wording now says "workflow" rather
  than "mode".

- **Revision**: Updated the root runtime interaction contract for 0106 - Binary
  confirmation UX.
  Runtime guidance now distinguishes non-binary closed choices, which keep `1` as
  the shortest accept path, from true binary mutation confirmations, which show
  `y`/`n`.

- **Revision**: Updated the root runtime interaction contract for 0101 - Quality
  skill UX action clarity.
  Runtime guidance now requires explicit shortest-answer paths for user
  interactions, code spans for concrete operational examples, and numbered
  ambiguity choices for scoped evaluation prompts.

## 2026-06-24

- **Restructure**: Started the runtime skill content as an OKF-shaped bundle with
  root `index.md`, `schema.md`, and `log.md`; added guide indexes/logs; and split
  authoring guidance into routed sub-guides.
