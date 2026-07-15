import type { QualityModel } from "../model/model.ts"
import type { WorkUnit } from "./graph.ts"
import type { EvaluationPlan, PlannedFactor } from "./plan.ts"
import type { StoredPayload } from "./protocol.ts"

type JsonObject = Record<string, unknown>

const ratingIds = (model: QualityModel) => model.ratingScale.map((level) => `rating:${level.level}`)

const routineRef = (kind: string, subject: JsonObject, selector = "") => ({
  kind,
  ...(selector === "" ? {} : { selector }),
  subject,
})

const factorChildren = (factor: PlannedFactor, plan: EvaluationPlan) =>
  plan.factors
    .filter(
      (candidate) =>
        candidate.areaId === factor.areaId &&
        candidate.path.length === factor.path.length + 1 &&
        candidate.path.slice(0, -1).join("/") === factor.path.join("/"),
    )
    .map((candidate) => candidate.ref)

const factorRequirements = (factor: PlannedFactor, plan: EvaluationPlan) =>
  plan.requirements
    .filter((requirement) => requirement.factorIds.includes(factor.ref))
    .map((requirement) => requirement.ref)

const evaluationFrame = (model: QualityModel, modelPath: string): JsonObject => ({
  derivedContext: { evaluationPolicies: ["source-as-data", "secret-redaction"], rigor: "standard" },
  inputs: { ratingLevelIds: ratingIds(model) },
  kind: "EvaluationFrame",
  schemaVersion: 3,
  subject: { modelLocator: modelPath },
})

const areaEvaluationFrame = (area: EvaluationPlan["areas"][number]): JsonObject => ({
  inputs: {
    childAreaIds: area.childAreaIds,
    localRequirementIds: area.localRequirementIds,
    rootFactorIds: area.rootFactorIds,
    sourceRefs: [area.source],
  },
  kind: "AreaEvaluationFrame",
  schemaVersion: 3,
  subject: { areaId: area.ref },
})

const requirementEvaluationFrame = (
  model: QualityModel,
  requirement: EvaluationPlan["requirements"][number],
): JsonObject => ({
  derivedContext: {
    appliedRatingCriteria: model.ratingScale.map((level) => ({
      criterion: requirement.value.ratings?.[level.level] ?? level.criterion,
      ratingLevelId: `rating:${level.level}`,
      source:
        requirement.value.ratings?.[level.level] === undefined
          ? "model_default"
          : "requirement_override",
    })),
  },
  inputs: {
    ratingLevelIds: ratingIds(model),
    requirementAssessmentBasis: requirement.value.assessment,
  },
  kind: "RequirementEvaluationFrame",
  schemaVersion: 3,
  subject: { factorIds: requirement.factorIds, requirementId: requirement.ref },
})

export const initialFramePayloads = (
  model: QualityModel,
  plan: EvaluationPlan,
  modelPath: string,
): ReadonlyArray<StoredPayload> => [
  { workUnit: "frameEvaluation", payload: evaluationFrame(model, modelPath) },
  ...plan.areas.map((area) => ({
    workUnit: `frameAreaEvaluation:${area.ref}`,
    payload: areaEvaluationFrame(area),
  })),
  ...plan.requirements.map((requirement) => ({
    workUnit: `frameRequirementEvaluation:${requirement.ref}`,
    payload: requirementEvaluationFrame(model, requirement),
  })),
]

export const deterministicPayload = (
  unit: WorkUnit,
  model: QualityModel,
  plan: EvaluationPlan,
  modelPath: string,
): JsonObject => {
  if (unit.kind === "frameEvaluation") return evaluationFrame(model, modelPath)
  if (unit.kind === "frameAreaEvaluation") {
    const area = plan.areas.find((entry) => entry.ref === unit.subject)!
    return {
      inputs: {
        childAreaIds: area.childAreaIds,
        localRequirementIds: area.localRequirementIds,
        rootFactorIds: area.rootFactorIds,
        ...(area.source === "" ? {} : { sourceRefs: [area.source] }),
      },
      kind: "AreaEvaluationFrame",
      schemaVersion: 3,
      subject: { areaId: area.ref },
    }
  }
  if (unit.kind === "frameRequirementEvaluation") {
    const requirement = plan.requirements.find((entry) => entry.ref === unit.subject)!
    const criteria = model.ratingScale.flatMap((level) => {
      const criterion = requirement.value.ratings?.[level.level] ?? level.criterion
      if (criterion === "") return []
      return [
        {
          criterion,
          ratingLevelId: `rating:${level.level}`,
          source:
            requirement.value.ratings?.[level.level] === undefined
              ? "model_default"
              : "requirement_override",
        },
      ]
    })
    return {
      ...(criteria.length === 0 ? {} : { derivedContext: { appliedRatingCriteria: criteria } }),
      inputs: {
        ratingLevelIds: ratingIds(model),
        ...(requirement.value.assessment === ""
          ? {}
          : { requirementAssessmentBasis: requirement.value.assessment }),
      },
      kind: "RequirementEvaluationFrame",
      schemaVersion: 3,
      subject: {
        ...(requirement.factorIds.length === 0 ? {} : { factorIds: requirement.factorIds }),
        requirementId: requirement.ref,
      },
    }
  }
  if (unit.kind === "frameFactorAnalysis") {
    const factor = plan.factors.find((entry) => entry.ref === unit.subject)!
    return {
      derivedContext: {
        emptySignalPolicy: "ignore_empty",
        synthesisGuidanceRef: "protocol:factor-synthesis-default-v0",
      },
      inputs: {
        childFactorAnalysisRefs: factorChildren(factor, plan).map((ref) =>
          routineRef("FactorAnalysisResult", { factorId: ref }, "localAndDescendantAnalysis"),
        ),
        directRequirementRatingRefs: factorRequirements(factor, plan).map((ref) =>
          routineRef("RequirementRatingResult", { requirementId: ref }),
        ),
      },
      kind: "FactorAnalysisFrame",
      schemaVersion: 3,
      subject: { areaId: factor.areaId, factorId: factor.ref },
    }
  }
  if (unit.kind === "frameAreaAnalysis") {
    const area = plan.areas.find((entry) => entry.ref === unit.subject)!
    return {
      derivedContext: {
        emptySignalPolicy: "ignore_empty",
        synthesisGuidanceRef: "protocol:area-synthesis-default-v0",
      },
      inputs: {
        childAreaAnalysisRefs: area.childAreaIds.map((ref) =>
          routineRef("AreaAnalysisResult", { areaId: ref }, "localAndDescendantAnalysis"),
        ),
        factorAnalysisRefs: area.rootFactorIds.map((ref) =>
          routineRef("FactorAnalysisResult", { factorId: ref }, "localAndDescendantAnalysis"),
        ),
      },
      kind: "AreaAnalysisFrame",
      schemaVersion: 3,
      subject: { areaId: area.ref },
    }
  }
  throw new Error(`no deterministic payload for ${unit.id}`)
}
