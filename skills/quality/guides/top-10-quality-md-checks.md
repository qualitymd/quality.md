# Top 10 QUALITY.md Checks

Use this checklist for a fast, read-only inspection of a QUALITY.md file's
current state, model quality, and lifecycle readiness. The result is not an
Evaluation report and does not rate the subject. It produces routing findings
that wizard and other modes can use to choose the next workflow.

## Inputs

- `qualitymd status [path] --json`
- the target `QUALITY.md`
- evaluation history summarized by status JSON

Do not inspect subject source files for this checklist. Do not read evaluation
report bodies. Keep the inspection bounded to the model file and status signals.

## Finding Shape

Report only findings that affect routing or model usefulness. Use this shape:

```text
QUALITY.md inspection findings
- <check id>: <finding>
  Evidence: <status field, section, or property>
  Impact: <why this affects lifecycle/model usefulness>
  Route: <setup | getting-started | authoring | evaluate | improve | history | update>
```

Keep evidence short. Cite section names, property paths, counts, or status JSON
fields rather than quoting long passages.

## The Checks

### 1. Lifecycle State

Use `qualitymd status --json` to identify whether the model is missing, invalid,
valid with no history, valid with history, or needs evaluation reconciliation.

- Finding when missing or invalid: route to setup or lint repair before any model
  quality judgment.
- Finding when history needs reconciliation: route to history/reconciliation
  before a fresh evaluation unless the user explicitly wants a new run.

### 2. Body Context

Inspect whether the Markdown body gives enough context to build and evaluate the
model: Overview, Scope, Needs, Risks, and Known gaps should be present and
substantive.

- Finding when body sections are empty, placeholder-like, or generic in a starter
  model: route to getting-started. Use the authoring guide as the quality
  reference for what the body should accomplish.
- Finding when Known gaps omits known unknowns while the rest of the body leaves
  unresolved questions: route to getting-started for first-run process or
  authoring for best-practice guidance.

### 3. Rating Scale Fit

Check whether the rating scale is understandable and fits the body's decision
context.

- Finding when level descriptions or criteria are generic enough that findings
  cannot distinguish `target` from `minimum`: route to authoring.
- Finding when a custom scale exists but the body does not explain why: route to
  authoring.

### 4. Subject and Source Alignment

Check whether the root title, body scope, file location, and root or child
`source` values describe the same evaluated subject.

- Finding when the title names the repository but the body/sources are narrower:
  route to getting-started or authoring.
- Finding when source coverage includes unrelated/generated/supporting artifacts:
  route to authoring.

### 5. Target Shape

Check whether the target tree is small enough to understand and specific enough
to represent distinct evaluated entities.

- Finding when all concerns are flattened into the root despite clear sub-entities
  in the body: route to authoring.
- Finding when child targets merely mirror the parent without distinct Factors or
  Requirements: route to authoring.

### 6. Factor Coverage

Check whether Factors reflect the body context: important needs and risks should
have a quality lens, and Factors should not be vague labels alone.

- Finding when major body needs/risks have no Factor: route to authoring.
- Finding when Factors are generic, overlapping, or unexplained: route to
  authoring.

### 7. Requirement Assessability

Check whether Requirements are concrete enough to produce findings and ratings.

- Finding when Requirements are aspirations rather than assessable expectations:
  route to authoring.
- Finding when a Requirement lacks observable evidence or criteria: route to
  authoring before evaluation.

### 8. Assessment Evidence

Check whether each Requirement's `assessment` gives the evaluator a usable means
of assessment, either inline or by referencing an entity that defines it.

- Finding when assessments are placeholders, circular, or vague: route to
  authoring.
- Finding when referenced assessment sources are not traceable from the model:
  route to authoring.

### 9. Evaluation Readiness

Decide whether the model is useful enough to evaluate without confusing model
weakness with subject quality.

- Finding when the model is valid but too vague to bind evidence to ratings:
  route starter/placeholder models to getting-started; route populated models to
  authoring.
- Finding when the model is valid, scoped, assessable, and unreconciled history
  is absent: route to evaluate.

### 10. Maintenance Signals

Use evaluation history and active recommendations to decide whether the next
workflow is maintenance rather than new authoring or evaluation.

- Finding when active recommendations exist: route to improve or recommendation
  review.
- Finding when the latest run is stale, incomplete, malformed, or unreportable:
  route to history/reconciliation.

## Summary Judgment

After the checks, classify the QUALITY.md lifecycle in one phrase:

- `missing`
- `invalid`
- `starter`
- `immature`
- `ready to evaluate`
- `has evaluation history`
- `needs reconciliation`

Use the finding routes to recommend one next workflow and list a few concrete
alternatives.
