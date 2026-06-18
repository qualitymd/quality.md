# 001 — README misstates which commands are built

**Target:** readme → approachability
**Related requirement:** *the README reflects what the CLI and spec actually provide*
**Severity:** medium — this is the binding constraint on the whole-model rating.

## Gap

`README.md:161-162` states: "The planned commands above other than `init` and
`lint` fail with 'unknown command' until they land." This implies `models` and
`spec` are not built and fail — but both are built and run successfully
(`./dist/qualitymd models list`, `./dist/qualitymd spec`), and the README's own
Install section (`README.md:169`) lists `init`, `lint`, `models`, and `spec` as
currently built. The sentence is both inaccurate and self-contradictory, and per
the model's Risks an overstating/misleading README is a first-order failure.

## Evidence

- `README.md:161-162` — the inaccurate claim.
- `README.md:148-156` — command list (models, spec present today).
- `README.md:169` — Install section lists models and spec as built.
- Verified: `./dist/qualitymd models list --json` and `./dist/qualitymd spec`
  both succeed.

## Options

1. Reword line 161 to scope "fail with unknown command" to the *planned*
   resources only (`evaluation`/`result`), leaving `init`, `lint`, `models`,
   `spec` as built.
2. Delete the sentence and rely on the existing `(planned)` markers in the
   command list.

## Recommended

**Option 1** — keep the helpful "unknown command until they land" signal but
correct its scope to the genuinely planned commands. Lowest-risk, preserves
intent, removes the contradiction.

## Done criterion

`the README reflects what the CLI and spec actually provide` reaches at least
**target**: every command's stated status matches the built CLI, with no
statement claiming a built command fails. Re-evaluate in a new numbered run.
