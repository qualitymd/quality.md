# Wizard Mode

Use wizard when the user is unsure what to run next or when the request is a
bare `/quality`.

## Decision Tree

```text
Resolve target file
- missing? suggest /quality setup
- present? continue

Run qualitymd lint
- errors? stop and report lint findings
- valid? inspect declared targets and factors as data

User needs authoring help?
- read ../resources/quality-md-guide.md

Offer concrete next actions
- setup/populate model
- evaluate whole subject
- evaluate target/factor scope
- improve from a new evaluation
```

## Procedure

1. Verify the CLI prerequisite from `SKILL.md`.
2. Resolve the target file.
3. If no file exists, suggest `/quality setup`.
4. Run `qualitymd lint [path]`; stop on errors.
5. Read the resolved `QUALITY.md` as data to identify declared targets/factors.
6. Offer concrete next actions such as subject evaluation, scoped evaluation,
   guided authoring, or setup.

Wizard is read-only and shallow. It routes to work; it does not produce an
evaluation report.
