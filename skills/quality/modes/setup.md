# Setup Mode

Use setup to create the minimal valid `QUALITY.md` skeleton and then route the
user toward guided population.

## Decision Tree

```text
Resolve target file
- missing? run qualitymd init [path]
- present? do not overwrite; validate it

Run qualitymd lint
- errors? stop and report lint findings
- valid? read getting-started + authoring and begin guided population, then wizard for next action
```

## Procedure

1. Verify the CLI prerequisite from `SKILL.md`.
2. Resolve the target file.
3. Emit the run frame:

   ```text
   /quality run
   - Mode: setup
   - Target file: <resolved path>
   - Scope: model skeleton/readiness
   - Mutation: `QUALITY.md` when missing or when an existing model change is confirmed
   - Artifacts: `QUALITY.md`
   - Next gate: lint result, guided population, then wizard next step
   ```

4. If no target file exists, run `qualitymd init [path]`.
5. If the target file exists and setup would change it, use a decision brief
   before editing:

   ```text
   Decision: update existing `QUALITY.md`?
   - Changes:
   - Evidence/reason:
   - Recommended option:
   - Alternatives:
   - Done criterion / verification:
   ```

6. Run `qualitymd lint [path]`; stop on errors and report the CLI findings.
7. Read [`../guides/authoring.md`](../guides/authoring.md) and
   [`../guides/getting-started.md`](../guides/getting-started.md), then begin
   guided population in the same run: draft the body's Overview, Scope, Needs,
   Risks, each with its unknowns, open questions, and any material support that
   is not agent-accessible, and propose project-specific Factors and Requirements
   to replace the placeholders with the user. Do not stop at naming the next
   step. Route to `wizard` after population for the next evaluation workflow.

`setup` creates a valid skeleton; it does not invent a complete quality model
without user/project context. For authoring judgment, read
[`../guides/authoring.md`](../guides/authoring.md).
