import type { EvaluationPlan, PlannedArea, PlannedFactor } from "./plan.ts"

export type WorkKind =
  | "frameEvaluation"
  | "frameAreaEvaluation"
  | "resolveSource"
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

const add = (
  units: Array<WorkUnit>,
  kind: WorkKind,
  subject: string,
  dependsOn: ReadonlyArray<string>,
  evaluatorBacked: boolean,
  dataKind?: string,
) => {
  const unit: WorkUnit = {
    id: unitId(kind, subject),
    kind,
    subject,
    dependsOn: [...new Set(dependsOn)],
    evaluatorBacked,
    ...(dataKind === undefined ? {} : { dataKind }),
  }
  units.push(unit)
  return unit.id
}

const bottomUp = <A extends { readonly path: ReadonlyArray<string> }>(items: ReadonlyArray<A>) =>
  items
    .map((item, index) => ({ item, index }))
    .sort(
      (left, right) => right.item.path.length - left.item.path.length || left.index - right.index,
    )
    .map(({ item }) => item)

export const buildGraph = (
  plan: EvaluationPlan,
  sourceKinds: Readonly<Record<string, "path" | "glob" | "prose">>,
): ReadonlyArray<WorkUnit> => {
  const units: Array<WorkUnit> = []
  add(units, "frameEvaluation", "", [], false, "EvaluationFrame")
  for (const area of plan.areas) {
    const parent =
      area.path.length === 0 ? "" : `area:${area.path.slice(0, -1).join("/") || "root"}`
    add(
      units,
      "frameAreaEvaluation",
      area.ref,
      ["frameEvaluation", ...(parent === "" ? [] : [unitId("frameAreaEvaluation", parent)])],
      false,
      "AreaEvaluationFrame",
    )
  }
  for (const area of plan.areas) {
    if (sourceKinds[area.ref] === "prose") {
      add(units, "resolveSource", area.ref, [unitId("frameAreaEvaluation", area.ref)], true)
    }
  }
  for (const requirement of plan.requirements) {
    const frame = add(
      units,
      "frameRequirementEvaluation",
      requirement.ref,
      [unitId("frameAreaEvaluation", requirement.areaId)],
      false,
      "RequirementEvaluationFrame",
    )
    add(
      units,
      "assessRateRequirement",
      requirement.ref,
      [
        frame,
        ...(sourceKinds[requirement.areaId] === "prose"
          ? [unitId("resolveSource", requirement.areaId)]
          : []),
      ],
      true,
    )
  }
  const analyzedFactors = new Map<string, string>()
  for (const factor of bottomUp(plan.factors)) {
    const dependencies = [
      ...requirementsOf(factor, plan).map((ref) => unitId("assessRateRequirement", ref)),
      ...childrenOf(factor, plan.factors).map((ref) => unitId("analyzeFactor", ref)),
    ]
    const frame = add(
      units,
      "frameFactorAnalysis",
      factor.ref,
      dependencies,
      false,
      "FactorAnalysisFrame",
    )
    analyzedFactors.set(
      factor.ref,
      add(
        units,
        "analyzeFactor",
        factor.ref,
        [frame, ...dependencies],
        true,
        "FactorAnalysisResult",
      ),
    )
  }
  const analyzedAreas = new Map<string, string>()
  for (const area of bottomUp(plan.areas)) {
    const dependencies = [
      ...area.rootFactorIds.map((ref) => analyzedFactors.get(ref)!),
      ...area.childAreaIds.map((ref) => unitId("analyzeArea", ref)),
      ...area.localRequirementIds.map((ref) => unitId("assessRateRequirement", ref)),
    ]
    const frame = add(
      units,
      "frameAreaAnalysis",
      area.ref,
      dependencies,
      false,
      "AreaAnalysisFrame",
    )
    analyzedAreas.set(
      area.ref,
      add(units, "analyzeArea", area.ref, [frame, ...dependencies], true, "AreaAnalysisResult"),
    )
  }
  const rankFindings = add(
    units,
    "rankFindings",
    "",
    plan.requirements.map((requirement) => unitId("assessRateRequirement", requirement.ref)),
    true,
    "FindingRankingResult",
  )
  const recommend = add(
    units,
    "recommend",
    "",
    [rankFindings, ...analyzedFactors.values(), ...analyzedAreas.values()],
    true,
    "RecommendationResult",
  )
  add(
    units,
    "rankRecommendations",
    "",
    [recommend, rankFindings],
    true,
    "RecommendationRankingResult",
  )
  add(
    units,
    "buildReports",
    "",
    units.map((unit) => unit.id),
    false,
  )
  return units
}

export const areaForUnit = (unit: WorkUnit, plan: EvaluationPlan): PlannedArea | undefined => {
  if (unit.subject.startsWith("area:")) return plan.areas.find((area) => area.ref === unit.subject)
  return plan.areas.find((area) =>
    plan.requirements.some(
      (requirement) => requirement.ref === unit.subject && requirement.areaId === area.ref,
    ),
  )
}
