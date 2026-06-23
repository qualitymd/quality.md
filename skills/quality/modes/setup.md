# Setup Mode

Use setup to create or update a useful first `QUALITY.md` through a short,
context-informed discovery flow. Setup writes only the selected `QUALITY.md`;
evaluation, quality-log entries, external issues, and recurring-review
automation are follow-on workflows.

## Decision Tree

```text
Verify CLI support and resolve model file
- CLI missing/stale? stop and route to update/install
- model missing? analyze context, ask discovery questions, run qualitymd init, author QUALITY.md
- model present? analyze context, ask discovery questions, confirm edits, author QUALITY.md

Validate and inspect readiness
- lint errors? report lint findings and route to continued QUALITY.md iteration
- valid? inspect with Top 10 checks, classify readiness, offer next-step choices
```

## Procedure

1. Verify the CLI prerequisite from `SKILL.md`.
2. Resolve the model file.
3. Emit the run frame:

   ```text
   /quality run
   - Mode: setup
   - Model file: <resolved path>
   - Scope: contextual QUALITY.md setup
   - Mutation: QUALITY.md only
   - Artifacts: QUALITY.md
   - Next gate: context analysis, discovery, lint, readiness inspection
   ```

4. Read [`../guides/authoring.md`](../guides/authoring.md). Read
   [`../guides/getting-started.md`](../guides/getting-started.md) when setup is
   continuing from a starter/immature model or needs first-run iteration
   guidance.
5. Inspect available repository context before asking setup questions. Keep this
   bounded to setup signals, not source-quality judgment. Useful signals include
   README/docs, package metadata, repository structure, tests/build scripts,
   contributor docs, existing agent instructions, and visible work-management or
   recurring-review hints. Build a short working fact sheet:

   ```text
   Root area: <path/title/default>
   Lifecycle: <default> (<strongly inferred | weakly inferred | assumed>, evidence)
   Risk tolerance: <default> (<confidence>, evidence)
   Modeling rigor: <default> (<confidence>, evidence)
   Collaboration: <default> (<confidence>, evidence)
   Needs: <candidate user/maintainer/other needs>
   Missing context: <specific gaps>
   Workflow posture: <work management / recurring review hints>
   Candidate model shape: <factors and child areas>
   ```

   Treat the current directory as the default root area convention unless the
   user supplied an explicit model path or context strongly indicates a narrower
   root.

6. Ask a compact discovery prompt, or a short sequence when the interaction
   surface needs it. Each question must include a recommended default and
   confidence signal:

   - lifecycle: exploratory, pre-release, active production, maintenance, or
     sunset;
   - risk tolerance: high, moderate, or low;
   - modeling rigor: lightweight, standard, or high-assurance;
   - collaboration context, assuming agent-heavy development and asking who else
     must align with the quality bar; and
   - needs map: primary user needs, collaborator/maintainer needs, other
     stakeholder needs, plus missing context.

   Seed the missing-context question from analysis, for example: `I think we're
   missing customer context and support workflow. What else important is not
   visible here?`

   Optional work-management and recurring-review questions may capture future
   posture, but they must not imply setup will create issues or automations. Ad
   hoc `/quality evaluate` is always available, not a selectable automation
   option. Recurring review should be cadence-oriented and maintainer-chosen; do
   not recommend CI or release gating by default.

7. If no model file exists, run `qualitymd init [path]` after discovery and
   before authoring content. If the model file exists and setup would change it,
   use a decision brief before editing:

   ```text
   Decision: update existing `QUALITY.md`?
   - Changes:
   - Evidence/reason:
   - Recommended option:
   - Alternatives:
   - Done criterion / verification:
   ```

8. Synthesize directly into `QUALITY.md`. Author the body first, then the
   frontmatter model:

   - Overview and Scope establish the root area, default current-directory
     scope, lifecycle, risk tolerance, modeling rigor, and key boundaries.
   - Needs and Risks capture user needs, maintainer/collaborator needs, other
     stakeholder needs, and the failure modes that matter.
   - Unknowns and open questions capture missing or non-agent-accessible context.
   - The rating scale stays standard unless the body shows a real mismatch.
   - Factors and child areas derive from the needs/risks and obvious component
     boundaries; do not ask users to design factors and areas cold.
   - Requirements are small, concrete, and assessable from agent-accessible
     evidence, or explicitly name missing evidence.
   - Include the `quality-md` self-check area when appropriate unless the user
     declines or the model file is not in the root area it governs. Use
     `quality-md` as the area key, `<Root Title> QUALITY.md` as the area title,
     an area `description`, and an explicit path-based `source` for the model
     file such as `./QUALITY.md`; do not use prose aliases such as `(this file)`
     for `source`. Add concise YAML comments around that area explaining that
     `source` is the `QUALITY.md` artifact being evaluated, while the
     requirement's `assessment` references the guide used to judge it. Prefer one
     area-level requirement that cites the active authoring guide once and lists
     each affected factor under `factors` when that guide defines one coherent
     judgment across the factors.

9. Run `qualitymd lint [path]`. Stop on lint errors, report the CLI findings,
   and route to continued `QUALITY.md` iteration. Do not recommend evaluation
   while the model is invalid.
10. When lint passes, inspect the resulting model with
    [`../guides/top-10-quality-md-checks.md`](../guides/top-10-quality-md-checks.md).
    This is a bounded model-readiness inspection, not a project evaluation.
    Classify readiness as `starter`, `immature`, or `ready to evaluate`.
11. Report setup completion status-first:

    ```text
    Setup complete
    - Changed: QUALITY.md
    - Validation: lint passed | lint failed
    - Readiness: starter | immature | ready to evaluate
    - Not done: no evaluation, no quality log, no issues, no automations
    - Next: continue iterating | run evaluation | set up recurring review | set up recommendation handoff | stop here
    ```

    If readiness is not `ready to evaluate`, list the most important model gaps
    and make continued iteration the recommended next step. Do not automatically
    take any next-step action.

`setup` creates or updates a useful first model; it does not invent a complete
quality model without user/project context, run an evaluation, write the quality
log, create external issues, or configure automation.
