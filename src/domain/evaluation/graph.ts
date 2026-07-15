import type { EvaluationPlan, PlannedArea, PlannedFactor } from "./plan.ts"

export type WorkKind =
  | "frameEvaluation"
  | "frameAreaEvaluation"
  | "frameRequirementEvaluation"
  | "assessRateRequirement"
  | "frameFactorAnalysis"
  | "analyzeFactor"
  | "frameAreaAnalysis"
  | "analyzeArea"
  | "rankFindings"
  | "recommend"
  | "rankRecommendations"
  | "buildReports"

export interface WorkUnit {
  readonly id: string
  readonly kind: WorkKind
  readonly subject: string
  readonly dependsOn: ReadonlyArray<string>
  readonly evaluatorBacked: boolean
  readonly dataKind?: string
}

const unitId = (kind: WorkKind, subject: string) => (subject === "" ? kind : `${kind}:${subject}`)

const childrenOf = (factor: PlannedFactor, factors: ReadonlyArray<PlannedFactor>) =>
  factors
    .filter(
      (candidate) =>
        candidate.areaId === factor.areaId &&
        candidate.path.length === factor.path.length + 1 &&
        candidate.path.slice(0, -1).join("/") === factor.path.join("/"),
    )
    .map((candidate) => candidate.ref)

const requirementsOf = (factor: PlannedFactor, plan: EvaluationPlan) =>
  plan.requirements
    .filter((requirement) => requirement.factorIds.includes(factor.ref))
    .map((requirement) => requirement.ref)

const makeUnit = (
  kind: WorkKind,
  subject: string,
  dependsOn: ReadonlyArray<string>,
  evaluatorBacked: boolean,
  dataKind?: string,
): WorkUnit => ({
  id: unitId(kind, subject),
  kind,
  subject,
  dependsOn: [...new Set(dependsOn)],
  evaluatorBacked,
  ...(dataKind === undefined ? {} : { dataKind }),
})

const bottomUp = <A extends { readonly path: ReadonlyArray<string> }>(items: ReadonlyArray<A>) =>
  items
    .map((item, index) => ({ item, index }))
    .sort(
      (left, right) => right.item.path.length - left.item.path.length || left.index - right.index,
    )
    .map(({ item }) => item)

export const buildGraph = (plan: EvaluationPlan): ReadonlyArray<WorkUnit> => {
  const evaluationFrame = makeUnit("frameEvaluation", "", [], false, "EvaluationFrame")
  const areaEvaluationFrames = plan.areas.map((area) => {
    const parent =
      area.path.length === 0 ? "" : `area:${area.path.slice(0, -1).join("/") || "root"}`
    return makeUnit(
      "frameAreaEvaluation",
      area.ref,
      ["frameEvaluation", ...(parent === "" ? [] : [unitId("frameAreaEvaluation", parent)])],
      false,
      "AreaEvaluationFrame",
    )
  })
  const requirementUnits = plan.requirements.flatMap((requirement) => {
    const frame = makeUnit(
      "frameRequirementEvaluation",
      requirement.ref,
      [unitId("frameAreaEvaluation", requirement.areaId)],
      false,
      "RequirementEvaluationFrame",
    )
    return [frame, makeUnit("assessRateRequirement", requirement.ref, [frame.id], true)]
  })
  const factorUnits = bottomUp(plan.factors).flatMap((factor) => {
    const dependencies = [
      ...requirementsOf(factor, plan).map((ref) => unitId("assessRateRequirement", ref)),
      ...childrenOf(factor, plan.factors).map((ref) => unitId("analyzeFactor", ref)),
    ]
    const frame = makeUnit(
      "frameFactorAnalysis",
      factor.ref,
      dependencies,
      false,
      "FactorAnalysisFrame",
    )
    return [
      frame,
      makeUnit(
        "analyzeFactor",
        factor.ref,
        [frame.id, ...dependencies],
        true,
        "FactorAnalysisResult",
      ),
    ]
  })
  const areaAnalysisUnits = bottomUp(plan.areas).flatMap((area) => {
    const dependencies = [
      ...area.rootFactorIds.map((ref) => unitId("analyzeFactor", ref)),
      ...area.childAreaIds.map((ref) => unitId("analyzeArea", ref)),
      ...area.localRequirementIds.map((ref) => unitId("assessRateRequirement", ref)),
    ]
    const frame = makeUnit("frameAreaAnalysis", area.ref, dependencies, false, "AreaAnalysisFrame")
    return [
      frame,
      makeUnit("analyzeArea", area.ref, [frame.id, ...dependencies], true, "AreaAnalysisResult"),
    ]
  })
  const rankFindings = makeUnit(
    "rankFindings",
    "",
    plan.requirements.map((requirement) => unitId("assessRateRequirement", requirement.ref)),
    true,
    "FindingRankingResult",
  )
  const recommend = makeUnit(
    "recommend",
    "",
    [
      rankFindings.id,
      ...bottomUp(plan.factors).map((factor) => unitId("analyzeFactor", factor.ref)),
      ...bottomUp(plan.areas).map((area) => unitId("analyzeArea", area.ref)),
    ],
    true,
    "RecommendationResult",
  )
  const rankRecommendations = makeUnit(
    "rankRecommendations",
    "",
    [recommend.id, rankFindings.id],
    true,
    "RecommendationRankingResult",
  )
  const beforeReports = [
    evaluationFrame,
    ...areaEvaluationFrames,
    ...requirementUnits,
    ...factorUnits,
    ...areaAnalysisUnits,
    rankFindings,
    recommend,
    rankRecommendations,
  ]
  const buildReports = makeUnit(
    "buildReports",
    "",
    beforeReports.map((unit) => unit.id),
    false,
  )
  return [...beforeReports, buildReports]
}

export const readyUnits = (
  graph: ReadonlyArray<WorkUnit>,
  completed: ReadonlySet<string>,
  limit: number,
): ReadonlyArray<WorkUnit> =>
  graph
    .filter(
      (unit) =>
        unit.evaluatorBacked && unit.dependsOn.every((dependency) => completed.has(dependency)),
    )
    .slice(0, limit)

export const areaForUnit = (unit: WorkUnit, plan: EvaluationPlan): PlannedArea | undefined => {
  if (unit.subject.startsWith("area:")) return plan.areas.find((area) => area.ref === unit.subject)
  return plan.areas.find((area) =>
    plan.requirements.some(
      (requirement) => requirement.ref === unit.subject && requirement.areaId === area.ref,
    ),
  )
}
