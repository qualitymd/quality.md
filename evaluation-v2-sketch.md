---
title: Evaluation v2 sketch
status: sketch
---

# Evaluation v2 sketch

This sketch explores a possible v2 shape for QUALITY.md evaluation. It is not a
specification or implementation plan yet.

The main shift is to treat evaluation as an agent-orchestrated judgment protocol,
not as an executable algorithm. The protocol should name the moves an agent makes,
the inputs and outputs expected at each move, the bottom-up traversal rules, and
the points where the agent must stop rather than produce fake precision.

## Core idea

Evaluation should separate these concerns:

1. Frame the Requirement Evaluation before evidence judgment, including evidence
   targets and applied Rating Level criteria.
2. Assess evidence for each Requirement.
3. Rate each Requirement Assessment against the pre-framed criteria.
4. Analyze each Factor from direct Requirement Ratings and child Factor analyses.
5. Analyze each Area from local Factor analyses and child Area analyses.
6. Select recommendations from the resulting analysis graph.
7. Report conversationally, and optionally capture durable artifacts.

This makes the durable record a possible projection of the evaluation, not the
starting mental model for the user experience.

## Evaluation protocol

Frame records use a shared shape:

```text
<Routine>Frame:
  subject        # what the routine is about
  inputs         # prior records or model references the routine may consume
  derivedContext # rules, criteria, policies, stop conditions, and limits derived before the routine runs
```

A frame is the structured input context for an agent-run routine. It freezes what
the routine may consider, what rules it must follow, and when it must stop.

### Frame the evaluation

Inputs:

- selected `QUALITY.md`
- requested scope
- rigor or depth expectation
- Rating Scale
- Area tree
- Factor trees within each Area
- in-scope Requirements
- evaluation limits and safety rules

Output:

- evaluation frame
- in-scope Areas, Factors, and Requirements
- out-of-scope or deferred elements
- known limitations

Stop when:

- the model cannot be resolved
- the scope is ambiguous and cannot be safely inferred
- in-scope Requirements are absent
- source evidence cannot be resolved

Stub output:

```text
EvaluationFrame:
  subject:
    modelRef
  inputs:
    requestedScope
    ratingScaleRef
    areaTreeRef
    factorTreeRefs
  derivedContext:
    resolvedScope
    rigor
    evaluationPolicies
    expectedEvaluationLimits
```

### Frame Area Evaluations

For each in-scope Area, produce an Area Evaluation Frame before evaluating local
Requirements, local Factors, and child Areas. This frame defines the Area-local
evaluation boundary without duplicating lower-level Requirement or Factor frames.

Stub output:

```text
AreaEvaluationFrame:
  subject:
    areaRef
  inputs:
    sourceRefs
    localRequirementRefs
    rootFactorRefs
    childAreaRefs
  derivedContext:
    scope
    expectedEvaluationLimits
```

### Frame Requirement Evaluations

For each local Requirement in an Area, produce a Requirement Evaluation Frame
before assessing evidence. This frame sets the evidence target and rating bar for
the Requirement, so the agent does not adapt criteria after seeing the evidence.

Inputs:

- Requirement statement
- Requirement criteria
- connected Factors
- source scope
- evaluation frame
- Rating Scale
- Requirement-specific rating overrides, if any

Output:

- Requirement reference
- Factor references
- evidence targets
- applied Rating Level criteria
- stop conditions
- expected evaluation limits

Output shape:

```text
RequirementEvaluationFrame:
  subject:
    requirementRef                       # stable reference to the Requirement
    factorRefs                           # Factors this Requirement contributes to in this Area
  inputs:
    ratingScaleRef
    requirementCriteriaRef
    ratingOverrideRefs
  derivedContext:
    evidenceTargets: EvidenceTarget[]    # pre-assessment questions or inspection targets
    appliedRatingCriteria: AppliedRatingCriterion[]
                                          # Requirement-specific criteria for each assessable Rating Level
    stopConditions: StopCondition[]      # conditions that should stop assessment/rating
    expectedEvaluationLimits: EvaluationLimit[]
                                          # known limits that constrain claims but do not necessarily stop assessment

EvidenceTarget:
  id                                     # local identifier within the Requirement frame
  question                               # what the assessment needs to establish
  purpose                                # why this target matters for judgment
  sourceRefs                             # Area-owned source references or narrower locators to inspect
  required: true | false                 # whether this target is rating-critical

AppliedRatingCriterion:
  level                                  # Rating Level id; title is intentionally omitted
  criterion                              # criterion adapted to this Requirement before evidence judgment
  source: model_default | requirement_override
  sourceRef                              # model default or Requirement override used as source
  adaptationRationale                    # why this adaptation preserves the intended bar

StopCondition:
  id                                     # local identifier within the Requirement frame
  condition                              # condition under which the agent should stop
  reason                                 # why continuing would produce weak or unsafe judgment

EvaluationLimit:
  id                                     # local identifier within the frame or result
  description                            # known boundary on what the evaluation can claim
  impact                                 # how the limit affects confidence, coverage, or rating
```

Notes:

- Do not include Rating Level titles; `level` is enough for structured records.
- Adapt criteria to the Requirement before assessment, not to the observed
  evidence.
- `appliedRatingCriteria` should include every assessable Rating Level and make
  adjacent Rating Levels distinguishable enough for assessment and rating.
- `evidenceTargets` are the specific things the assessment needs to establish
  before it can make a fair judgment. They are pre-assessment questions or
  inspection targets, not findings.
- `stopConditions` are the cases where the agent should stop instead of
  producing a weak or fake judgment, such as unavailable source evidence,
  unsafe source instructions, or evidence that cannot distinguish adjacent Rating
  Levels.
- Area frames own source boundaries. Requirement frames may narrow through
  `EvidenceTarget.sourceRefs`, but should not duplicate Area source scope.
- A stop condition prevents assessment or rating. An expected evaluation limit
  constrains the claims the evaluation can make, but assessment may still
  continue.

### Frame Factor Analyses

For each Factor node in an Area's Factor tree, produce a Factor Analysis Frame
after its child Factors have been analyzed and before analyzing that Factor.
This frame identifies the direct Requirement Ratings and completed child Factor
analyses that the Factor analysis may synthesize.

Stub output:

```text
FactorAnalysisFrame:
  subject:
    areaRef                              # Area that owns this Factor node
    factorRef                            # Factor node being framed
  inputs:
    directRequirementRatingRefs          # ratings for Requirements attached directly to this Factor node
    childFactorAnalysisRefs              # completed direct child Factor localAndDescendantAnalysis refs
  derivedContext:
    synthesisGuidanceRef                 # reference to guidance for combining local and child Factor signals
    emptySignalPolicy: ignore_empty | empty_blocks_analysis | empty_counts_as_not_analyzed
    stopConditions: StopCondition[]      # conditions that should stop Factor analysis
    expectedEvaluationLimits: EvaluationLimit[]
                                          # known limits that constrain claims but do not necessarily stop analysis
```

Factor framing is one frame per Factor node, not one frame per local,
descendant, and combined-analysis substep. `childFactorAnalysisRefs` point to
completed child Factor `localAndDescendantAnalysis` results, so a parent Factor
can account for the full descendant subtree through direct child analyses.
`directRequirementRatingRefs` point only to Requirement Rating Results connected
to this exact Factor node; child Factor Requirements are represented through the
child Factor analysis refs.

`emptySignalPolicy` controls how empty local or child input scopes affect the
Factor analysis:

- `ignore_empty` means an empty local or child signal does not prevent analysis
  if the other signal exists. This is the default because grouping Factors may
  have no direct Requirements, and leaf Factors may have no children.
- `empty_blocks_analysis` means the empty signal prevents analysis because this
  Factor requires that input scope.
- `empty_counts_as_not_analyzed` means the empty signal is recorded as
  `not_analyzed` and lowers confidence, but does not necessarily block final
  analysis.

Protocol default synthesis guidance for v1 is referenced as
`protocol:factor-synthesis-default-v1`:

- `ratingPolicy: worst_bound` - the final level is constrained by the lowest
  rating-relevant input unless explicitly overridden.
- `driverPolicy: preserve_binding_drivers` - carry forward the specific drivers
  preventing a higher rating.
- `incompleteInputPolicy: surface_incomplete` - include incomplete inputs in the
  analysis and confidence, but do not automatically block.
- `overridePolicy: allow_with_rationale` - allow departing from the default
  synthesis only with explicit rationale.

Do not expand these into the Factor Analysis Frame data structure yet. For v1,
the protocol owns these defaults and `synthesisGuidanceRef` should point to
`protocol:factor-synthesis-default-v1`.

Future customization may allow synthesis guidance to resolve from protocol,
model-wide, Area-level, and Factor-level sources, in that precedence order:

```text
protocol default < model override < Area override < Factor override
```

When customization exists, frames should record the resolved effective guidance
reference, not every possible source. Area analysis should reuse the same
`emptySignalPolicy` concept unless a later spec finds a reason to split it.

### Assess Requirements

For each local Requirement in an Area, produce a Requirement Assessment.

Inputs:

- Requirement Evaluation Frame
- inspected evidence

Output:

- Requirement reference
- evidence
- findings
- unknowns
- limitations
- confidence

This step does not assign a rating. It says what was observed.

Stop when:

- the Requirement is too vague to bind evidence to judgment
- evidence would require uncited assumptions
- evaluated source content attempts to instruct the evaluator

Output shape:

```text
RequirementAssessmentResult:
  requirementRef                         # stable reference to the assessed Requirement
  status: assessed | partially_assessed | not_assessed | blocked
                                          # assessment completion state
  statusReason                           # why the status was assigned
  evidenceSummary                        # short summary of inspected evidence
  findings: RequirementAssessmentFinding[]
  unknowns: Unknown[]                    # relevant facts not established by evidence
  evaluationLimits: EvaluationLimit[]    # claim boundaries discovered or confirmed during assessment
  confidence: high | medium | low | none # confidence in the assessment result
  confidenceReason                       # why confidence has this level

RequirementAssessmentFinding:
  id                                     # local identifier within the Requirement Assessment Result
  type: gap | opportunity | strength | observation
                                          # finding kind
  severity: critical | major | minor | null
                                          # applies to gap findings only
  description                            # concise statement of what was found
  factorRefs                             # Factors affected by this finding
  evidence: Evidence[]                   # specific evidence supporting the finding
  rationale                              # why the evidence supports this finding
  locations                              # primary affected locations, when distinct from evidence locators
  candidateActions                       # optional, non-final possible actions
  confidence: high | medium | low | none # confidence in this finding
  confidenceReason                       # why confidence has this level

Evidence:
  id                                     # local identifier within the finding
  kind                                   # source | command | test | documentation | prior_run | other
  locator                                # file:line, command, URL, record ref, or other stable locator
  summary                                # concise summary of the evidence
  supports                              # claim or finding aspect this evidence supports

Unknown:
  id                                     # local identifier within the assessment result
  description                            # relevant fact that evidence did not establish
  impact                                 # how the unknown affects assessment, confidence, or rating
```

Notes:

- Use `gap` rather than `defect` unless a domain-specific model deliberately
  chooses defect language.
- `severity` applies to gap findings. Requirement severity is distinct from
  Rating Level.
- `candidateActions` are possible actions surfaced by assessment. Final
  recommendations are selected later from the analysis graph; candidate actions
  should stay lightweight and non-final.
- `evidenceSummary` gives a short assessment-level summary of the inspected
  evidence, while finding `evidence` carries the specific evidence behind each
  finding.
- `confidenceReason` explains why confidence is high, medium, or low, especially
  when the evidence is partial or indirect.
- `confidence` is a coarse enum:
  - `high` means evidence is direct, current, sufficient for the claim, and
    independently checkable.
  - `medium` means evidence is relevant and plausible, but partial, indirect,
    sampled, or not fully verified.
  - `low` means evidence is thin, ambiguous, stale, inferred, or materially
    incomplete.
  - `none` means no confidence judgment was possible because there was no
    assessment-quality evidence.
- Finding IDs are local to the Requirement Assessment Result unless a later
  record needs to reference them by path.
- `unknowns` are relevant facts the assessment did not establish.
- `evaluationLimits` are boundaries on what the assessment attempted or can
  honestly claim.
- Status semantics:
  - `assessed` means required evidence targets were sufficiently addressed for
    assessment.
  - `partially_assessed` means some evidence targets were addressed, but material
    unknowns or limits remain.
  - `not_assessed` means assessment-quality evidence was not available or not
    inspected.
  - `blocked` means a stop condition or safety issue prevented assessment from
    proceeding.

### Rate Requirements

For each Requirement Assessment, produce a Requirement Rating.

Inputs:

- Requirement
- Requirement Evaluation Frame
- Requirement Assessment

Output:

- Requirement reference
- rating level
- rationale
- rating drivers
- criteria results
- missing evidence
- confidence

The rating rationale should explain why this rating was chosen and why the
rating is not higher.

Stop when:

- evidence cannot distinguish adjacent Rating Levels
- material claims are not tied to evidence
- rating overrides cannot be interpreted

Output shape:

```text
RequirementRatingResult:
  requirementRef                         # stable reference to the rated Requirement
  status: rated | not_rated | blocked    # rating completion state
  statusReason                           # why the status was assigned
  level                                  # selected Rating Level id; present only when status is rated
  rationale                              # why this level was selected
  ratingDrivers: RequirementRatingDriver[]
                                          # findings, unknowns, or limits that determine the rating
  criteriaResults: CriterionResult[]     # how assessment output maps to each applied criterion
  missingEvidence: MissingEvidence[]     # evidence targets still missing or unresolved
  evaluationLimits: EvaluationLimit[]    # claim boundaries carried into the rating
  confidence: high | medium | low | none # confidence in the rating result
  confidenceReason                       # why confidence has this level

CriterionResult:
  level                                  # Rating Level id from appliedRatingCriteria
  matched: true | false | partial | unknown
                                          # whether assessment output satisfies this level's criterion
  rationale                              # why this criterion result was assigned

RequirementRatingDriver:
  findingRefs                            # RequirementAssessmentFinding refs, when applicable
  unknownRefs                            # Unknown refs, when applicable
  evaluationLimitRefs                    # EvaluationLimit refs, when applicable
  description                            # concise explanation of the rating driver
  effect: supports_level | prevents_higher | blocks_rating
                                          # how this driver affects the rating

MissingEvidence:
  evidenceTargetRef                      # EvidenceTarget ref from the RequirementEvaluationFrame
  description                            # what evidence is missing
  impact                                 # how the missing evidence affects rating or confidence
```

Notes:

- Rating maps `RequirementAssessmentResult` output to the pre-framed
  `RequirementEvaluationFrame.appliedRatingCriteria`. It should not re-assess
  evidence.
- The rationale should explain both why the selected level applies and why the
  rating is not higher.
- Missing evidence should not automatically become an unacceptable rating. It
  should produce `not_rated`, `blocked`, or a lower-confidence rating depending
  on the pre-framed criteria and assessment status.
- Status semantics:
  - `rated` means a Rating Level was assigned.
  - `not_rated` means no Rating Level was assigned because assessment evidence
    was insufficient, but no stop or safety condition blocked the routine.
  - `blocked` means a stop condition, unsafe input, invalid frame, or
    contradiction prevented rating.
- `level` should be omitted unless `status` is `rated`.
- `criteriaResults` should include one result for every
  `AppliedRatingCriterion` in the Requirement Evaluation Frame, even when the
  result is `unknown`.
- When `status` is `rated`, `level` should match one of the pre-framed
  `AppliedRatingCriterion.level` values.
- Each `RequirementRatingDriver` should have at least one of `findingRefs`,
  `unknownRefs`, or `evaluationLimitRefs`.
- A `RequirementRatingDriver` with `effect: blocks_rating` should produce
  `status: blocked`, not `status: rated`.
- Requirement rating uses the same confidence enum and semantics as Requirement
  assessment.
- Unmet criteria are represented by `criteriaResults`; do not add a separate
  `unmetCriteria` field unless a later report projection needs a denormalized
  convenience field.

### Analyze Factors

Walk each Area's Factor tree bottom-up.

For each Factor:

1. Analyze child Factors first.
2. Gather direct Requirement Ratings connected to this Factor.
3. Frame this Factor Analysis with direct Requirement Ratings and completed child
   Factor analyses.
4. Analyze local and local-and-descendant signal in one Factor Analysis Result.

Inputs:

- Area reference
- Factor reference
- direct Requirement Ratings
- child Factor Analyses
- Rating Scale
- synthesis guidance

Output:

- Area reference
- Factor reference
- local analysis
- local and descendant analysis
- rationale
- rating drivers
- incomplete direct Requirements
- incomplete child Factors
- confidence

Output shape:

```text
FactorAnalysisResult:
  areaRef                                # Area that owns this Factor node
  factorRef                              # Factor node being analyzed
  localAnalysis: FactorScopedAnalysis    # direct Requirement Ratings only
  localAndDescendantAnalysis: FactorScopedAnalysis
                                          # direct Requirement Ratings plus descendant Factor analyses

FactorScopedAnalysis:
  status: analyzed | empty | not_analyzed | blocked
                                          # completion state for this input scope
  statusReason                           # why the status was assigned
  level                                  # selected Rating Level id; present only when status is analyzed
  rationale                              # why this scoped analysis has this result
  inputRefs                              # Requirement Rating refs, child Factor analysis refs, or localAnalysis ref
  ratingDrivers: FactorRatingDriver[]    # inputs that determine the scoped analysis
  incompleteInputs: IncompleteInput[]    # inputs that were absent, incomplete, or unusable
  evaluationLimits: EvaluationLimit[]    # claim boundaries carried into this analysis
  confidence: high | medium | low | none # confidence in this scoped analysis
  confidenceReason                       # why confidence has this level

IncompleteInput:
  inputRef                               # reference to the missing or incomplete input
  reason                                 # why the input is incomplete or unusable
  impact                                 # how the incomplete input affects analysis or confidence

FactorRatingDriver:
  requirementRatingRefs                  # Requirement Rating refs, for localAnalysis drivers
  childFactorAnalysisRefs                # child Factor localAndDescendantAnalysis refs, when applicable
  localAnalysisRef                       # localAnalysis ref, for localAndDescendantAnalysis drivers
  evaluationLimitRefs                    # EvaluationLimit refs, when applicable
  description                            # concise explanation of the rating driver
  effect: supports_level | prevents_higher | blocks_analysis
                                          # how this driver affects the scoped analysis
```

Notes:

- `localAnalysis` analyzes only direct Requirement Rating Results attached to
  this Factor node.
- `localAndDescendantAnalysis` is the final Factor analysis. It synthesizes
  `localAnalysis` with direct child Factor `localAndDescendantAnalysis` results.
- A parent Factor only needs direct child Factor analyses because each child
  analysis already accounts for its own descendants.
- `level` should be omitted unless the scoped analysis `status` is `analyzed`.
- Status semantics:
  - `analyzed` means inputs existed and were sufficient to produce this scoped
    analysis.
  - `empty` means no inputs existed for this scope.
  - `not_analyzed` means inputs existed, but the routine could not produce a
    defensible analysis.
  - `blocked` means a stop condition, invalid input, or safety issue prevented
    analysis.
- For `localAnalysis`, `empty` means no direct Requirement Rating Results.
- For `localAndDescendantAnalysis`, `empty` means no local signal and no child
  Factor signal.
- For a leaf Factor, `localAndDescendantAnalysis` usually matches
  `localAnalysis`, subject to empty-signal policy.
- For a grouping Factor with no direct Requirement Ratings, `localAnalysis` may
  be `empty` while `localAndDescendantAnalysis` is synthesized from child Factor
  analyses.
- `localAnalysis.ratingDrivers` should use `requirementRatingRefs` and
  `evaluationLimitRefs`.
- `localAndDescendantAnalysis.ratingDrivers` should use `localAnalysisRef`,
  `childFactorAnalysisRefs`, and `evaluationLimitRefs`.

### Analyze Areas

Walk the Area tree bottom-up.

For each Area:

1. Assess and rate local Requirements.
2. Analyze the Area's local Factor forest.
3. Analyze child Areas.
4. Frame this Area Analysis with Factor analyses and completed child Area
   analyses.
5. Analyze local and local-and-descendant signal in one Area Analysis Result.

Inputs:

- Area reference
- Factor Analyses
- child Area Analyses
- Rating Scale
- synthesis guidance

Frame shape:

```text
AreaAnalysisFrame:
  subject:
    areaRef                              # Area node being framed
  inputs:
    factorAnalysisRefs                   # this Area's root Factor localAndDescendantAnalysis refs
    childAreaAnalysisRefs                # completed direct child Area localAndDescendantAnalysis refs
  derivedContext:
    synthesisGuidanceRef                 # reference to guidance for combining local and child Area signals
    emptySignalPolicy: ignore_empty | empty_blocks_analysis | empty_counts_as_not_analyzed
    stopConditions: StopCondition[]      # conditions that should stop Area analysis
    expectedEvaluationLimits: EvaluationLimit[]
                                          # known limits that constrain claims but do not necessarily stop analysis
```

Output:

- Area reference
- local analysis
- local and descendant analysis
- rationale
- rating drivers
- incomplete Factors
- incomplete child Areas
- confidence

Output shape:

```text
AreaAnalysisResult:
  areaRef                                # Area node being analyzed
  localAnalysis: AreaScopedAnalysis      # this Area's Factor analyses only
  localAndDescendantAnalysis: AreaScopedAnalysis
                                          # this Area's Factor analyses plus descendant Area analyses

AreaScopedAnalysis:
  status: analyzed | empty | not_analyzed | blocked
                                          # completion state for this input scope
  statusReason                           # why the status was assigned
  level                                  # selected Rating Level id; present only when status is analyzed
  rationale                              # why this scoped analysis has this result
  inputRefs                              # Factor analysis refs, child Area analysis refs, or localAnalysis ref
  ratingDrivers: AreaRatingDriver[]      # inputs that determine the scoped analysis
  incompleteInputs: IncompleteInput[]    # inputs that were absent, incomplete, or unusable
  evaluationLimits: EvaluationLimit[]    # claim boundaries carried into this analysis
  confidence: high | medium | low | none # confidence in this scoped analysis
  confidenceReason                       # why confidence has this level

AreaRatingDriver:
  factorAnalysisRefs                     # Factor localAndDescendantAnalysis refs, for localAnalysis drivers
  childAreaAnalysisRefs                  # child Area localAndDescendantAnalysis refs, when applicable
  localAnalysisRef                       # localAnalysis ref, for localAndDescendantAnalysis drivers
  evaluationLimitRefs                    # EvaluationLimit refs, when applicable
  description                            # concise explanation of the rating driver
  effect: supports_level | prevents_higher | blocks_analysis
                                          # how this driver affects the scoped analysis
```

Notes:

- `localAnalysis` analyzes this Area's root Factor
  `localAndDescendantAnalysis` results.
- `localAndDescendantAnalysis` is the final Area analysis. It synthesizes
  `localAnalysis` with direct child Area `localAndDescendantAnalysis` results.
- A parent Area only needs direct child Area analyses because each child analysis
  already accounts for its own descendants.
- The root Area's `localAndDescendantAnalysis.level` is the overall evaluation
  rating.
- `level` should be omitted unless the scoped analysis `status` is `analyzed`.
- Status semantics match `FactorScopedAnalysis`:
  - `analyzed` means inputs existed and were sufficient to produce this scoped
    analysis.
  - `empty` means no inputs existed for this scope.
  - `not_analyzed` means inputs existed, but the routine could not produce a
    defensible analysis.
  - `blocked` means a stop condition, invalid input, or safety issue prevented
    analysis.
- For a leaf Area, `localAndDescendantAnalysis` usually matches
  `localAnalysis`, subject to empty-signal policy.
- For a grouping Area with no local Factors, `localAnalysis` may be `empty`
  while `localAndDescendantAnalysis` is synthesized from child Area analyses.
- `localAnalysis.ratingDrivers` should use `factorAnalysisRefs` and
  `evaluationLimitRefs`.
- `localAndDescendantAnalysis.ratingDrivers` should use `localAnalysisRef`,
  `childAreaAnalysisRefs`, and `evaluationLimitRefs`.

Protocol default synthesis guidance for v1 is referenced as
`protocol:area-synthesis-default-v1` and uses the same default policies as
`protocol:factor-synthesis-default-v1` unless a later spec finds a reason to
split them.

Open question:

- Whether Area analysis needs different synthesis defaults from Factor analysis.

### Frame Recommendation Selection

Before selecting recommendations, produce a Recommendation Selection Frame over
the completed analysis graph. This frame should identify candidate action
surfaces and selection constraints without deciding the final recommendation set.

Stub output:

```text
RecommendationSelectionFrame:
  subject:
    rootAreaAnalysisRef
  inputs:
    ratingDriverRefs
    candidateActionRefs
  derivedContext:
    selectionConstraints
```

### Frame Report Projection

Before generating human-readable reports, produce a Report Projection Frame. This
frame should identify which structured results project into `report.md` and
`areas/**/report.md`.

Stub output:

```text
ReportProjectionFrame:
  subject:
    rootAreaAnalysisRef
  inputs:
    recommendationRefs
  derivedContext:
    rootReportPath
    areaReportPaths
    projectionRules
```

## Pseudocode shape

This pseudocode names the protocol moves, but the real implementation target is a
Markdown instruction file for an agent.

```text
evaluateModel(model, evaluationFrame):
  modelContext = frameEvaluation(model, evaluationFrame)
  rootAreaAnalysis = evaluateArea(model.rootArea, modelContext)
  recommendationFrame = frameRecommendationSelection(rootAreaAnalysis, modelContext)
  recommendations = selectRecommendations(recommendationFrame, rootAreaAnalysis)
  reportFrame = frameReportProjection(
    rootAreaAnalysis,
    recommendations,
    modelContext
  )
  return assembleEvaluationResult(
    modelContext,
    rootAreaAnalysis,
    recommendations,
    reportFrame
  )
```

```text
evaluateArea(area, modelContext):
  areaFrame = frameAreaEvaluation(area, modelContext)
  areaContext = prepareAreaContext(area, areaFrame, modelContext)

  requirementFrames = frameAreaRequirementEvaluations(areaFrame, areaContext)
  requirementAssessments = assessAreaRequirements(requirementFrames, areaContext)
  requirementRatings = rateRequirementAssessments(
    requirementFrames,
    requirementAssessments,
    areaContext
  )

  factorAnalyses = analyzeAreaFactorForest(
    area,
    area.factorForest,
    requirementRatings,
    areaContext
  )

  childAreaAnalyses = []
  for childArea in area.children:
    childAreaAnalyses.append(evaluateArea(childArea, modelContext))

  areaAnalysisFrame = frameAreaAnalysis(
    areaFrame,
    factorAnalyses,
    childAreaAnalyses,
    areaContext
  )
  return analyzeArea(areaAnalysisFrame, areaContext)
```

```text
analyzeAreaFactorForest(area, factorForest, requirementRatings, areaContext):
  requirementIndex = indexRequirementRatingsByFactor(requirementRatings)

  factorAnalyses = []
  for factor in factorForest.roots:
    factorAnalyses.append(
      analyzeFactorNode(area, factor, requirementIndex, areaContext)
    )

  return factorAnalyses
```

```text
analyzeFactorNode(area, factor, requirementIndex, areaContext):
  directRequirementRatings = requirementIndex.directRatingsFor(factor)

  childFactorAnalyses = []
  for childFactor in factor.children:
    childFactorAnalyses.append(
      analyzeFactorNode(area, childFactor, requirementIndex, areaContext)
    )

  factorFrame = frameFactorAnalysis(
    area,
    factor,
    directRequirementRatings,
    childFactorAnalyses,
    areaContext
  )

  return analyzeFactor(factorFrame, areaContext)
```

## Routine list

Framing and validation:

- `frameEvaluation`
- `frameAreaEvaluation`
- `prepareAreaContext`
- `validateEvaluationFrame`
- `validateAreaContext`

Requirement work:

- `frameRequirementEvaluation`
- `frameAreaRequirementEvaluations`
- `collectRequirementEvidence`
- `assessRequirement`
- `assessAreaRequirements`
- `rateRequirement`
- `rateRequirementAssessments`

Factor work:

- `indexRequirementRatingsByFactor`
- `analyzeAreaFactorForest`
- `analyzeFactorNode`
- `frameFactorAnalysis`
- `analyzeFactor`

Area work:

- `frameAreaAnalysis`
- `analyzeArea`

Recommendation and reporting work:

- `frameRecommendationSelection`
- `selectRecommendations`
- `frameReportProjection`
- `assembleEvaluationResult`
- `reportEvaluationResult`
- `captureEvaluationArtifacts`

## Routine prompt contracts

Because the evaluation algorithm is agent-orchestrated from Markdown
instructions, routines should be specified as prompt contracts rather than only
as program functions.

Each agent-run routine should define:

```text
RoutinePrompt:
  role
  task
  inputs
  requiredOutput
  constraints
  stopRules
  selfCheck
```

The prompt contract should make clear what the instruction file asks the agent to
do, what inputs the agent may use, what output shape is required, and when the
agent must stop rather than inventing precision.

### frameRequirementEvaluation prompt contract

Task:

Create a `RequirementEvaluationFrame` before inspecting assessment evidence.

Inputs:

- Requirement statement
- Requirement criteria
- connected Factors
- Area/source context
- Rating Scale
- Requirement-specific rating overrides, if any

Instructions:

- Identify the evidence targets needed to assess the Requirement.
- Adapt Rating Level criteria to this Requirement before evidence judgment.
- Use Requirement overrides when present; otherwise adapt model defaults.
- Do not use observed evidence to make criteria easier or harder.
- Define stop conditions for missing, ambiguous, unsafe, or insufficient
  evidence.
- Define expected evaluation limits known before assessment.

Required output:

- `RequirementEvaluationFrame` JSON.

Self-check:

- Are adjacent Rating Levels distinguishable?
- Are evidence targets sufficient to judge the Requirement?
- Are criteria adapted to the Requirement, not observed evidence?
- Are stop conditions specific enough to prevent fake precision?

## Spec and runtime organization

Evaluation v2 should have a parent durable spec folder that owns the whole v2
concept and splits detail by independently reviewable contract.

Proposed durable spec shape:

```text
specs/evaluation-v2/
  index.md
  evaluation-v2.md
  protocol.md
  routines/
    index.md
    routine-contract.md
    frame-evaluation.md
    frame-area-evaluation.md
    frame-requirement-evaluation.md
    assess-requirement.md
    rate-requirement.md
    frame-factor-analysis.md
    analyze-factor.md
    frame-area-analysis.md
    analyze-area.md
    frame-recommendation-selection.md
    select-recommendations.md
  records/
    index.md
    data-folder.md
    area-data-tree.md
    routine-record.md
    requirement-evaluation-frame-json.md
    requirement-assessment-result-json.md
    requirement-rating-result-json.md
    factor-analysis-frame-json.md
    factor-analysis-result-json.md
    area-analysis-frame-json.md
    area-analysis-result-json.md
    recommendation-selection-json.md
  reports/
    index.md
    report-projection.md
    report-md.md
    area-report-md.md
```

The parent `evaluation-v2.md` should hold shared invariants: bottom-up
traversal, frames-before-judgment, source-as-data, evidence locators,
confidence, rating-driver preservation, structured data as source of truth, and
the separation between evaluation state and report projection.

Reports are orthogonal to the core evaluation protocol. The evaluation source of
truth is `data/**`; `report.md` and `areas/**/report.md` are generated
human-readable projections that consume completed structured evaluation data.

The skill runtime should mirror the spec lightly:

```text
skills/quality/workflows/evaluation-v2.md
skills/quality/workflows/evaluation-v2/
  protocol.md
  routines/
    ...
```

The spec defines the durable contract. The skill files are runtime instructions
that implement that contract for the agent.

## Rating drivers

Use `rating drivers` for the specific evidence-backed findings or constraints
that determine a rating. They answer:

- why this rating?
- why not higher?
- what would need to change for the rating to move?

Rating drivers should survive roll-up. A parent analysis should not hide the
lower-level issue that determines its rating.

## Recommendation selection

Recommendations should probably be selected after the analysis graph exists,
rather than emitted independently at every layer.

Reasoning:

- The same issue may appear at Requirement, Factor, and Area levels.
- Recommendations should target the lowest useful level.
- Cross-cutting issues may need a Factor-level or Area-level recommendation.
- The final recommendation set should focus on rating movement and user action,
  not complete restatement of every finding.

Possible routine:

```text
selectRecommendations(analysisGraph):
  identify rating drivers
  identify the lowest useful action surface for each driver
  merge duplicates
  prioritize recommendations by expected rating movement and importance
  return recommendation candidates
```

## Persistence shape

The agent should orchestrate the protocol and make judgment calls. The CLI should
own persistence, retrieval, validation, numbering, and report projection.

Within an evaluation run folder, structured routine data should live under a
`data/` subfolder. Human-readable run artifacts can remain at the run root.

Possible run shape:

```text
NNNN-evaluation/
  report.md                  # root Area report
  areas/
    <area>/
      report.md              # focused Area report
      areas/
        <child-area>/
          report.md
  data/
    frame/
      evaluation-frame.json
      report-projection-frame.json
    areas/
      root/
        area-evaluation-frame.json
        area-analysis-frame.json
        area-analysis-result.json
        requirements/
          <requirement>/
            requirement-evaluation-frame.json
            assessment-result.json
            rating-result.json
        factors/
          <factor>/
            factor-analysis-frame.json
            factor-analysis-result.json
            factors/
              <child-factor>/
                factor-analysis-frame.json
                factor-analysis-result.json
        areas/
          <child-area>/
            ...
    recommendations/
      selection-frame.json
      selection-result.json
```

The exact filenames and record format are still open. The important boundary is:

- `report.md` is the human-readable report for the root Area
- `areas/**/report.md` is the human-readable report projection for non-root
  Areas, arranged in an OKF-style Area tree
- `data/` is for CLI-managed structured records
- `data/areas/` mirrors the evaluated Area tree; each Area owns its local
  `requirements/`, local `factors/`, and child `areas/`
- `factors/` lives inside an Area because Factor meaning is scoped by the Area;
  child Factors recurse through a nested `factors/` folder
- routine outputs should usually be single JSON files, not folders:
  `*-frame.json`, `assessment-result.json`, `rating-result.json`, and
  `*-analysis-result.json`
- use a folder for a routine output only if that output truly needs multiple
  files or attachments
- records should reference their routine inputs so the CLI can support resume,
  stale-record detection, QC, and report projection

Possible record concept:

```text
RoutineRecord:
  id
  kind
  inputRefs
  output
  status
  createdAt
  supersedes
  qc
```

Do not persist every intermediate thought. Persist routine outputs that help with
resume, audit, QC, or reporting.

## Future QC layer

This is not part of v1, but the protocol should leave room for a quality-control
step around substantive judgment moves.

QC should challenge routine output; it should not become a second full
evaluation.

Possible QC result:

```text
qcResult:
  status: accepted | needs_revision | blocked
  issues
  requiredRevisions
  confidence
```

Possible pause-and-check rules:

- Is every material claim tied to evidence?
- Does the rating match the stated criterion?
- Does the rationale explain why the rating is not higher?
- Are missing evidence and uncertainty visible?
- Did lower-level rating drivers survive the roll-up?
- Did aggregate analysis avoid averaging away an unacceptable child signal?
- Are model weakness, evidence weakness, and evaluated-source weakness kept
  distinct?

Use validation for structural checks and QC for judgment checks.

Validation examples:

- required fields exist
- references resolve
- Factor links are valid
- impossible states are rejected

QC examples:

- evidence supports findings
- findings support ratings
- roll-up rationale preserves rating drivers
- uncertainty is represented honestly

## Open questions

- How should unlinked Requirements contribute to Area analysis?
- How should one Requirement connected to multiple Factors contribute to each
  Factor?
- Can Requirement-to-Factor links carry role, strength, or weight?
- What is the right representation for empty local signal?
- What is the right representation for empty descendant signal?
- What roll-up policy should be default: worst-bound, weighted synthesis, apex
  rules, veto rules, or agent judgment with a proposed default?
- How should protocol synthesis guidance be customized later: model-wide,
  Area-level, Factor-level, or only through explicit future extensions?
- Should recommendations be selected only at the end, or can routines emit
  candidate recommendations for final selection?
- What durable artifacts should be created by default, if any?
- Which parts of this protocol belong in the skill instruction file, and which
  belong in CLI-supported record shapes?
