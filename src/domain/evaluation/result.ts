import Ajv2020, { type ErrorObject } from "ajv/dist/2020.js"

import { expectedSchemaFor, type StoredPayload } from "./protocol.ts"
import type { WorkUnit } from "./graph.ts"

type JsonObject = Record<string, unknown>

export class ResultValidationError extends Error {
  readonly category: "invalid_evaluator_output" | "schema_invalid_output"

  constructor(category: "invalid_evaluator_output" | "schema_invalid_output", message: string) {
    super(message)
    this.category = category
  }
}

const describe = (errors: ReadonlyArray<ErrorObject> | null | undefined) =>
  (errors ?? [])
    .map((error) => `${error.instancePath || "payload"} ${error.message ?? "is invalid"}`)
    .join("; ")

const validate = (unit: WorkUnit, payload: JsonObject) => {
  const ajv = new Ajv2020({ allErrors: true, strict: false, validateFormats: false })
  const check = ajv.compile(expectedSchemaFor(unit))
  if (!check(payload)) {
    throw new ResultValidationError("schema_invalid_output", describe(check.errors))
  }
}

const object = (value: unknown, detail: string): JsonObject => {
  if (value === null || Array.isArray(value) || typeof value !== "object") {
    throw new ResultValidationError("invalid_evaluator_output", detail)
  }
  return { ...(value as JsonObject) }
}

const normalizeKind = (
  payload: JsonObject,
  kind: string,
  subjectField: string,
  subject: string,
) => {
  const output: JsonObject = { ...payload, schemaVersion: 3, kind }
  if (subjectField !== "") {
    const existing = output[subjectField]
    if (typeof existing === "string" && existing !== "" && existing !== subject) {
      throw new ResultValidationError(
        "invalid_evaluator_output",
        `payload ${subjectField} = ${JSON.stringify(existing)}, want the work unit subject ${JSON.stringify(subject)}`,
      )
    }
    output[subjectField] = subject
  }
  return output
}

export const normalizeEvaluatorResult = (
  unit: WorkUnit,
  value: unknown,
  recommendationToken: (index: number) => string,
): ReadonlyArray<StoredPayload> => {
  const payload = object(value, "evaluator returned no payload")
  if (unit.kind === "assessRateRequirement") {
    const assessment = normalizeKind(
      object(payload.assessment, "combined requirement judgment must carry an assessment object"),
      "RequirementAssessmentResult",
      "requirementId",
      unit.subject,
    )
    const rating = normalizeKind(
      object(payload.rating, "combined requirement judgment must carry a rating object"),
      "RequirementRatingResult",
      "requirementId",
      unit.subject,
    )
    validate(unit, { assessment, rating, evidence: payload.evidence })
    return [
      { workUnit: unit.id, payload: assessment },
      { workUnit: unit.id, payload: rating },
    ]
  }
  if (unit.kind === "recommend") {
    if (!Array.isArray(payload.recommendations) || payload.recommendations.length === 0) {
      throw new ResultValidationError(
        "invalid_evaluator_output",
        "recommend result must carry a non-empty recommendations array",
      )
    }
    const recommendationState = payload.recommendations.reduce<{
      readonly used: ReadonlyArray<string>
      readonly recommendations: ReadonlyArray<JsonObject>
    }>(
      (state, value, index) => {
        const recommendation = normalizeKind(
          object(value, `recommendations[${index}] must be an object`),
          "RecommendationResult",
          "",
          "",
        )
        let id = typeof recommendation.id === "string" ? recommendation.id : ""
        if (!/^qrec_[a-z0-9]+$/.test(id)) {
          do id = `qrec_${recommendationToken(index)}`
          while (state.used.includes(id))
          recommendation.id = id
        }
        if (state.used.includes(id)) {
          throw new ResultValidationError(
            "invalid_evaluator_output",
            `recommendations[${index}] duplicates id ${id}`,
          )
        }
        return {
          used: [...state.used, id],
          recommendations: [...state.recommendations, recommendation],
        }
      },
      { used: [], recommendations: [] },
    )
    const recommendations = recommendationState.recommendations
    validate(unit, { recommendations })
    return recommendations.map((recommendation) => ({ workUnit: unit.id, payload: recommendation }))
  }
  const subjectField =
    unit.kind === "analyzeFactor" ? "factorId" : unit.kind === "analyzeArea" ? "areaId" : ""
  const normalized = normalizeKind(payload, unit.dataKind!, subjectField, unit.subject)
  validate(unit, normalized)
  return [{ workUnit: unit.id, payload: normalized }]
}
