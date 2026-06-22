---
type: Functional Specification
title: /quality reporting
description: Component spec for /quality evaluation reports, run artifacts, records, recommendation files, and report correction behavior.
tags: [skill, quality, evaluation, reporting]
timestamp: 2026-06-22T00:00:00Z
---

# /quality reporting

This spec owns the `/quality` skill's evaluation reporting and run-artifact
contract: the evaluation folder shape, generated reports, records,
recommendations, correction behavior, and reportability expectations. It
composes the shared contracts in the parent [/quality skill](quality-skill.md)
spec and the judgment workflow in
[/quality evaluation workflow](evaluation.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear
in all capitals.

## Reporting

The skill produces an **Evaluation Report** that conforms to
[Report](../../../SPECIFICATION.md#report) — the Rating and its rationale, the
Scope, the per-area requirement/factor/local/aggregate ratings with
rationales, and the Advice. *Not assessed* outcomes **MUST** appear wherever they
occur, distinct from rated outcomes.

Every evaluation that finds gaps **MUST** also emit its Advice as discrete,
triageable **recommendation** artifacts — recommendations are a product of
evaluation and the input to
[recommendation follow-up](recommendation-follow-up.md).

A rating level's name can collide with QUALITY.md structural vocabulary —
most often the suggested scale's **Area** level against a **Area** entity.
Wherever a level name could be read as a structural term, the report **MUST**
qualify it: name the level with a qualifier (the **Area** rating level;
*rated* **Area**; *meets* **Area**; *held at* **Unacceptable**) rather than
a bare noun, and keep structural areas introduced by their `Area:` heading
label. The same applies to any author-named level coinciding with *Area*,
*Factor*, or *Requirement*.

The CLI creates a numbered evaluation folder per run, so each run is a durable,
routable record. The default parent directory is `quality/evaluations/`, but a
repository may set `.quality/config.yaml`:

```yaml
evaluationDir: quality/evaluations
```

`evaluationDir` names the parent directory that contains numbered run folders.
The shared folder and record contract is defined by
[`Evaluation records`](../../evaluation-records.md), with artifact-specific
details in its [child specs](../../evaluation-records/index.md).
It **MUST** be repository-relative, normalized before use, and rejected when it is
absolute or escapes the repository. Missing config or missing `evaluationDir`
uses the default. Unknown config keys should be surfaced as warnings and
ignored.

Runtime evaluation artifacts are raw outputs in the evaluated repository, not
OKF concepts. They **MUST NOT** carry OKF frontmatter or require registration in
`specs/schema.md`. Alongside the report and its recommendations the folder
captures four further artifacts that make the run auditable and reproducible —
a snapshot of the model evaluated, the run's **design** (its inputs), its
execution **plan** (its method), and its process-only **debug log**:

```
quality/evaluations/
  0001[-<narrowing>]-quality-eval/
    model.md
    design.md
    debug-log.md
    plan.md
    assessments/
      001-<area>-<requirement>.json
      002-<area>-<requirement>.json
    analysis/
      <area>.json
      <child-area>.json
    report-summary.md
    report.md
    report.json
    recommendations/
      001-<slug>.md
      002-<slug>.md
```

The folder name **MUST** be deterministic:
`NNNN[-<narrowing>]-quality-eval`, where `<narrowing>` is the scoped
area/factor slug, omitted for a whole-model run. `NNNN` is the next integer in
the resolved evaluation directory.

Together these separate the three things an audit must tell apart — the *inputs*
(design), the *method* (plan), and the *result* (report) traced to a fixed model
(snapshot):

- The folder **MUST** include a **snapshot of the `QUALITY.md` as evaluated** —
  the model state the ratings were produced against. A rating is only meaningful
  against the model that defined its criteria, and that model may change after the
  run; the snapshot makes the report a self-contained, reproducible record whose
  findings trace to the exact requirements and `source` selectors in force at
  evaluation time. It is a verbatim capture, not a runtime judgment, and
  should record the revision (e.g. commit) of the evaluated source it was taken
  against.
- The folder **MUST** include a **design** artifact recording the evaluation's
  resolved parameters — mode, model file, scope, and rigor (see
  [Arguments](quality-skill.md#arguments)) — and a citation of the `model.md` snapshot the run is
  bound to. It is the authoritative record of *what* was evaluated and *under what
  inputs*, so a later reader or re-run can reproduce the setup. The report's
  **Scope** statement is the reader-facing summary of this; the full
  parameterization lives here, stated once.
- The folder **MUST** include a **plan** artifact recording the run's *method* —
  how the skill covers each in-scope area's `source` at the chosen rigor (per
  [Rigor levels](evaluation.md#rigor-levels)): the entities or hotspots to assess, their
  order, and any diagnostics to run. The report's statement of what was *not
  assessed* (see [Rigor levels](evaluation.md#rigor-levels)) **MUST** reconcile actual
  coverage against this plan, so divergence between intended and achieved coverage
  is visible rather than silent. The plan **MUST** name the rigor level and the
  concrete requirements selected by that rigor. The design and plan together
  **MUST** record enough concise report context for the CLI-rendered summary
  layer: rigor, scope or narrowing, in-scope requirement set, out-of-scope or
  deferred areas, headline evidence basis, and limitations that constrain the
  rating.
- The folder can include optional **planned coverage** metadata when the run
  needs machine-checkable resume diagnostics. The skill supplies the intended
  assessment requirements and analysis areas as `coverage:` frontmatter in
  `plan.md` after the plan is settled.
- The folder **MUST** include a **debug log** recording notable events involving
  the evaluation process itself: scope ambiguity, history inspection, coverage
  adjustment, interruptions or resumes, retries, record corrections, tooling
  failures, redaction decisions, prompt-injection handling, or report generation
  recovery. `debug-log.md` is not an assessment record, rating rationale, report,
  or evidence store. Evaluation findings and rating rationale belong in the
  formal records. When a project command is exercised as evidence against the
  evaluated source, the log may record only the routing fact and point to the
  assessment record; it must not copy raw command output or preserve a second
  copy of assessment evidence.
- The folder **MUST** capture the **assessment result records** the Evaluate phase
  produces as JSON — one artifact per in-scope requirement, holding its findings
  (each with its locator), the rating inferred against the requirement's
  `criterion`, and a brief rationale: the assess → finding → rating chain of
  [Grounding judgment](evaluation.md#grounding-judgment). A *not assessed* requirement gets a
  record too, with `ratingResult.kind: not-assessed`, and a rationale stating
  the absent evidence. Each record is **written atomically and never mutated** —
  a re-assessment produces a new evaluation folder rather than editing an
  existing record. The skill writes assessment result records through
  `qualitymd evaluation assessment add`; the CLI owns serialization,
  numbering, and `schemaVersion`.
- The folder **MUST** capture the **analysis records** the Analyze phase produces
  as JSON — one write-once artifact per area node — holding that node's inferred
  **local** and **aggregate** ratings and its **factor** ratings, each with a brief
  rationale naming the binding constraints (the inferred, weighted roll-up of
  [Grounding judgment](evaluation.md#grounding-judgment)). Each record **MUST cite the records
  it derives from**: the in-scope **assessment result records** behind its local rating,
  and its **children's analysis records** behind its aggregate — so the chain leaf
  → node → root is explicit and a *not assessed* outcome is visible wherever it
  propagates. The skill writes analysis records through
  `qualitymd evaluation analysis set`; the CLI owns serialization and
  `schemaVersion`.

Assessment, analysis, and report JSON files **MUST** use stable generic
top-level fields tied to the evaluation workflow, not fields invented for one
factor or requirement. Domain-specific details live under `attributes` on the
smallest relevant object.

An assessment result record's finding uses generic fields:

- `locator`
- `observation`
- `category`
- `severity`
- `evidence`
- optional `attributes`

For example, a secret finding may use `category: "secret"` and
`severity: "critical"` with `attributes.credentialType`; it must not include the
secret value. A prompt-injection observation may use
`category: "prompt-injection"` and is recorded, not followed. Severity values
come from the canonical evaluation-record vocabulary and reports render their
display titles.

The report is the **render over these records**, not an independent copy:
`report-summary.md` is the concise human triage artifact, `report.md` is the
full human rendering, and `report.json` is the machine-readable rendering of the
same result, produced by `qualitymd evaluation report build`. The assessment result records are the source of record for
Assess-and-Rate and the analysis records for Analyze, and the report's
per-requirement and per-area sections derive from them (the report adds the
Advise and Report layers and the reader-facing framing). `report.json` should
inline only minimal generic finding summaries by assessment-record reference for
single-file consumers; full finding detail remains in `assessments/*.json`. This
keeps the report from drifting and makes every rating in it traceable — leaf
finding → assessment result record → analysis record → report — to the immutable records
that produced it.

Human Markdown report labels are resolved from the run's `model.md` snapshot:
Model, Area, Factor, and Rating Level titles are primary display text, with
stable identifiers retained where the report needs traceability. `report.json`
preserves stable identifiers for machines.

The CLI-rendered report artifacts are specified by the durable report specs:
[`report-summary.md`](../../reports/report-summary-md.md),
[`report.md`](../../reports/report-md.md), and
[`report.json`](../../reports/report-json.md). The concise summary **MUST** read
as a decision brief for human readers: key details, Verdict, Area Breakdown,
Selected Findings, Recommended Actions, and Scope & Limitations. Its key details
use reader-facing labels including "Full evaluation" for an unnarrowed run and
"Evaluation verdict" for the in-scope root Area's Area-with-descendants verdict.
Its Recommended Actions section surfaces copyable Recommendation IDs for
follow-up prompts. The full `report.md` remains verdict-first before detailed
area and requirement sections. The JSON report **MUST** expose the same
summary-layer data with non-null scope, empty arrays for empty collections,
explicit rating objects for null or not-assessed ratings, typed lifecycle state
for assessment and recommendation digests, typed next-step state, typed
missing-metadata entries, and a structural Area-only rating state for area
groups. The skill must treat those typed report states as the routing source
rather than inferring state from `null`, absent fields, or `active` booleans
alone.

Like the report, the design, plan, assessment, and analysis records reference any
secret value by `file:line` and type only (see
[Boundaries](quality-skill.md#boundaries-and-hard-rules)).

A worked reference instance of this layout — model snapshot, design, plan,
assessment result records, analysis records, report, and recommendations — is in
[`examples/`](examples/index.md).

Each recommendation file **MUST** stand on its own as a unit a reader can triage
and route without the report or the session in front of them. It **MUST** state:
the gap it closes, with the evidence and `file:line` locators behind it; a small
set of remediation **options**; exactly one option marked **recommended**; and a
**done-criterion** expressed as the outcome the in-scope requirement should reach
against its `criterion`: for a rated gap, a target rating level; for a *not
assessed* gap, becoming assessable and reaching at least the acceptable floor.
That is what recommendation follow-up verifies, with a scoped re-evaluation when
the done criterion is rating-bound. When the evidence or source structure makes
ownership inferable, the recommendation should
name the route hint in existing text, such as the affected package, path,
workflow, maintainer surface, or verification command. Like the report, a
recommendation references any secret value by `file:line` and type only (see
[Boundaries](quality-skill.md#boundaries-and-hard-rules)). The skill writes recommendation
records through `qualitymd evaluation recommendation add`; the CLI owns
Markdown frontmatter, numbering, and stable rendering.

When correcting an already written recommendation, the skill should write a
new recommendation record with `supersedes` pointing at the stale
recommendation, rather than appending ambiguous advice with no active-state
signal. Appending a correction without `supersedes` leaves the run reportable and
renders both files, so the report's primary Next Action can still point at the
stale original — the ambiguity is silent. Superseding makes the active advice
unambiguous while preserving the audit trail.

When correcting an already written assessment, the skill should write a new
assessment result record with `supersedes` pointing at the stale assessment, then
replace the affected analysis record so it references the active assessment. This
analysis step is required for assessment results — and not for recommendations — because
analysis ratings bind to assessment references, so a corrected assessment left
unpaired with its analysis would let a roll-up silently rely on stale judgment.

- A report **MUST** state the **Scope** it was produced under, so a scoped result
  is never mistaken for a whole-model verdict.
- A report **MUST** distinguish *not assessed* outcomes from the report's
  **Limitations** statement. *Not assessed* is a Rating Result where evidence was
  absent, shown per requirement and roll-up. **Limitations** bounds how far a
  rated outcome should be trusted and reconciles actual coverage against the
  plan.
- The CLI **MUST** render all report forms: concise prose for triage in
  `report-summary.md`, full prose for a person in `report.md`, and a
  machine-readable form in `report.json`. The underlying result is the same;
  only the rendering differs.
