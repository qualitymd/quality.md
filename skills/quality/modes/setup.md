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
- valid? route to wizard for population or next action
```

## Procedure

1. Verify the CLI prerequisite from `SKILL.md`.
2. Resolve the target file.
3. If no target file exists, run `qualitymd init [path]`.
4. Run `qualitymd lint [path]`; stop on errors and report the CLI findings.
5. Hand off to `wizard` to guide model population and the next evaluation.

`setup` creates a valid skeleton; it does not invent a complete quality model
without user/project context. For authoring judgment, read
[`../resources/quality-md-guide.md`](../resources/quality-md-guide.md).
