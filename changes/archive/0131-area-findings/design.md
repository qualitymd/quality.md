---
type: Design Doc
title: Area findings in evaluation reports — design
description: How Area Findings are added to Area analysis data and rendered in Area and Factor reports.
tags: [evaluation, reports, skill]
timestamp: 2026-06-26T00:00:00Z
---

# Area findings in evaluation reports — design

## Context

Answers the [functional spec](spec.md) for change case
[0131](../0131-area-findings.md). The spec adds Area Findings as analysis-phase
findings persisted on `AreaAnalysisResult.findings`, projected into Area and
Factor reports, and kept deliberately separate from recommendations, global
top-finding synthesis, and impact/priority ranking.

The current evaluation implementation is already built around map-backed
payload contracts in `internal/evaluation/data_contract.go` and deterministic
report projection in `internal/evaluation/report_tree.go`. This design extends
those seams instead of adding a new payload kind or a parallel typed model.

## Approach

### Extend `AreaAnalysisResult`, not the data-kind set

Area Findings live inside the existing `AreaAnalysisResult` payload:

```json
{
  "schemaVersion": 2,
  "kind": "AreaAnalysisResult",
  "areaId": "area:payments",
  "findings": [
    {
      "id": "retry-guarantees-unclear",
      "type": "gap",
      "severity": "high",
      "confidence": "medium",
      "summary": "Retry guarantees are not clear enough to rely on payment webhook delivery.",
      "rationale": "Synthesizes retry and observability gaps into one Payments-level concern.",
      "inputRefs": [
        {
          "kind": "RequirementAssessmentResult",
          "subject": {
            "requirementId": "requirement:payments::retry-window"
          },
          "selector": "findings[retry-window-unspecified]"
        }
      ],
      "factorRelationships": [
        {
          "factorId": "factor:payments::reliability",
          "relationship": "primary-driver",
          "rationale": "Directly limits the reliability rating."
        }
      ]
    }
  ],
  "localAnalysis": {
    "status": "analyzed"
  },
  "localAndDescendantAnalysis": {
    "status": "analyzed"
  }
}
```

`dataContracts[DataKindAreaAnalysis]` changes from the generic
`analysisResultContract` helper to an Area-specific object that includes the
same `localAnalysis` and `localAndDescendantAnalysis` fields plus optional
`findings`. Factor analysis keeps the generic helper unchanged.

The field is optional at validation time for clean compatibility with existing
runs and partial in-progress runs. Consumers normalize absence to an empty array.
New skill-authored Area analysis results should write `findings: []` when no Area
Findings exist, matching the repo's repeated-field convention.

### Reuse the existing object-contract machinery

Add two local contract helpers in `data_contract.go`:

- `areaFindingContract()`
- `areaFindingFactorRelationshipContract()`

The Area Finding contract uses existing closed vocabularies where they already
exist:

- `type`: the current finding type set, including `note`;
- `severity`: the current finding severity set, including `info`;
- `confidence`: the analysis confidence set;
- `inputRefs`: the existing `routineRefContract()`.

Those vocabularies are encoded with the existing `enum(...)` helper so
`qualitymd evaluation data schema AreaAnalysisResult` exposes the closed sets as
JSON Schema `enum` values and `data set` rejects out-of-set values before
writing. The Area Finding and Factor relationship contracts remain closed object
contracts, so forbidden advice/ranking fields are rejected through the same
unknown-field path as misspellings.

The contract adds `summary` rather than reusing Requirement Finding
`description`. Reports can fall back to `description` only for defensive display
if historical or hand-authored data ever appears, but the schema and examples use
`summary`.

### Add focused semantic validation after generic model binding

The existing recursive model-reference walker already validates that any
`factorId` resolves in the model. Area Findings need two validations the generic
walker cannot infer:

1. Area Finding IDs are unique within one `AreaAnalysisResult.findings` array.
2. `factorRelationships[].factorId` belongs to the same declaring Area as the
   containing `AreaAnalysisResult.areaId`.

Add an Area-analysis-specific post-validation branch after generic structural and
model-reference validation. That branch parses `areaId`, walks
`findings[*].factorRelationships[*].factorId`, parses each Factor reference, and
compares `factor.DeclaringArea` to the Area result ID.

Keep the check in the evaluation package's data-contract layer so `data set`,
`data set --dry-run`, `data verify`, schema examples, and reportability all
agree about what a valid payload is.

### Report collection keeps findings with the owning Area

`collectEvaluationAreaAnalysis` already stores the entire Area analysis payload
on `evaluationArtifacts.area(...).Analysis`. No additional collector state is
needed. Report rendering reads `objectSlice(area.Analysis["findings"])`.

Factor reports need the owning Area's analysis. The report tree already has the
Factor ID's declaring Area, so the Factor report can look up:

```text
out.area(areaKey(factor.ID.DeclaringArea)).Analysis.findings
```

Then it filters to findings whose `factorRelationships[].factorId` equals the
current Factor reference.

This keeps Area Findings owned by Area analysis while allowing Factor reports to
render their local view without duplicating data.

### Sorting is deterministic, not an importance score

Add small rank maps in `report_tree.go` for:

- finding type;
- finding severity;
- finding confidence;
- Factor relationship.

Area report sort key:

```text
type rank, severity rank, confidence rank, original index
```

Factor report sort key:

```text
type rank, severity rank, relationship rank, confidence rank, original index
```

Unknown values should sort after known values while still rendering as their raw
value or existing display fallback. The validator should prevent unknown values
for new data; the defensive fallback only keeps report generation robust for old
or hand-edited artifacts.

### Rendering shape

Area reports add a `Findings` section near the other Area analysis material:
after the summary/header and before rating drivers is the natural location. The
section should render an empty-state row when no findings exist.

The summary table columns:

```text
Finding | Type | Severity | Confidence | Related Factors | Summary
```

Factor reports add a `Findings` section, using the same heading as Area reports
because the report context already names the Factor. Its table should include the
local relationship instead of all related Factors:

```text
Finding | Type | Severity | Relationship | Confidence | Summary
```

Area Finding detail rendering can stay compact in the first implementation:
`rationale`, `inputRefs`, and Factor relationship rationales can be visible in a
detail table beneath the summary table if present. The key design point is that
the report does not synthesize new prose: it renders fields already in
`AreaAnalysisResult.findings`.

### Skill and durable spec updates

The runtime skill changes are prompt-contract changes, not new tooling:

- `analyzeArea` is instructed to produce `AreaAnalysisResult.findings`.
- The skill is told to synthesize only from its frame and persisted routine
  outputs, not from new source inspection.
- Existing QC language is extended so roll-up-binding and low-confidence Area
  Findings are included in verification.

The durable spec updates mirror the same phase boundary so the reason survives
after the Change Case is archived: Requirement assessment observes, Area analysis
synthesizes Area Findings, report generation projects persisted data, and advice
remains deferred.

## Spec response

- **Area Finding data contract** — satisfied by extending only
  `AreaAnalysisResult`, with optional `findings` normalized to empty and local ID
  uniqueness checked in Area-analysis validation. Closed vocabularies are emitted
  as schema enums and enforced by the existing validator.
- **Factor relationships** — satisfied by an embedded relationship object plus
  relationship enum validation and same-Area post-validation.
- **Skill workflow** — satisfied by prompt/spec edits to `analyzeArea` and QC;
  no CLI workflow command changes are needed.
- **Reports** — satisfied by rendering Area-owned findings directly in Area
  reports and filtered findings in Factor reports.
- **Report ordering** — satisfied by explicit rank maps in report rendering,
  without introducing importance or priority fields.

## Alternatives

- **New `AreaFindingResult` data kind.** Rejected because Area Findings are part
  of Area analysis, share its lifecycle, and need the containing `areaId` for
  scope. A separate data kind would add paths, identity, report collection, and
  status behavior without a separate routine phase.
- **Top-level `EvaluationFinding` objects.** Rejected for this slice because the
  spec intentionally avoids cross-Area findings and global top-finding synthesis.
  Containment under `AreaAnalysisResult` keeps the scope honest.
- **Put findings inside `localAnalysis` and `localAndDescendantAnalysis`.**
  Rejected because the same Area Finding could explain both local and descendant
  analysis, and duplicating it across both scopes would create identity and
  drift problems. Rating drivers remain the per-scope explanation surface.
- **Use `description` instead of `summary`.** Rejected for Area Findings because
  the report table needs a short statement distinct from optional rationale and
  from Requirement Finding detail prose. Requirement Findings keep their existing
  shape.
- **Compute report order from severity only.** Rejected because type and local
  Factor relationship are important report-local context. The chosen rank order
  is still deterministic and explicitly not an importance score.

## Trade-offs & risks

- **Two finding shapes.** Requirement Findings and Area Findings share type and
  severity but differ in fields (`description` versus `summary`,
  `factorRelationships`). This is intentional because they are produced by
  different phases, but docs and report labels need to keep the distinction clear.
- **Optional `findings`.** Allowing absence preserves old runs, but new skill
  output should write empty arrays for clarity. Tests should cover both absent and
  empty cases.
- **Same-Area validation is custom.** The generic reference walker cannot enforce
  relationship locality, so a targeted validation branch is required. Tests should
  cover a valid same-Area Factor and an invalid sibling-Area Factor.
- **Closed objects reject future-looking fields.** The first implementation
  intentionally rejects fields such as `impact`, `priority`, `weight`, `actions`,
  and `recommendations`. A later recommendation or global-ranking case will need
  an explicit schema change rather than relying on ignored extension fields.
- **Selectors remain stringly typed.** `inputRefs.selector` can point to
  `findings[<id>]`, but the CLI does not currently validate selector existence.
  This is acceptable for this case because routine references are already the
  established traceability mechanism; selector validation can be considered
  separately if it becomes a recurring source of broken links.
- **Root/global expectations.** Root Area reports may show root Area Findings,
  but they will not rank child Area Findings globally. The report text and
  closeout should avoid implying a global "top findings" view exists.

## Open questions

None for this design. Global finding synthesis and recommendations are explicit
deferred work.
