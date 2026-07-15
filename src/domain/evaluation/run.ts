import { findElement, parseAreaReference, projectModel, type QualityModel } from "../model/model.ts"

export const EvaluationSchemaVersion = 3

export interface PlannedScope {
  readonly areaId: string
  readonly factorFilter: ReadonlyArray<string>
}

export interface RequestedScope {
  readonly areaId?: string
  readonly factorFilter?: ReadonlyArray<string>
}

export interface EvaluationManifest {
  readonly schemaVersion: 3
  readonly kind: "EvaluationManifest"
  readonly evaluationId: string
  readonly createdAt: string
  readonly model: string
  readonly requestedScope: RequestedScope
  readonly plannedScope: PlannedScope
  readonly run: { readonly number: number; readonly label: string }
}

const alphabet = "0123456789abcdefghjkmnpqrstvwxyz"

export const newEvaluationIdentity = (now: Date, bytes: Uint8Array) => {
  const tail = [...bytes].map((value) => alphabet[value & 31]).join("")
  const createdAt = `${now.toISOString().slice(0, 19)}Z`
  return {
    evaluationId: `${createdAt.replaceAll("-", "").replaceAll(":", "").replace(".000", "")}-${tail}`,
    createdAt,
  }
}

const factorArea = (reference: string) => {
  if (!reference.startsWith("factor:") || !reference.includes("::")) {
    throw new Error(`factor model reference ${JSON.stringify(reference)} is invalid`)
  }
  return `area:${reference.slice(7).split("::", 1)[0]}`
}

export const resolveScope = (
  model: QualityModel,
  area: string,
  factors: ReadonlyArray<string>,
): { readonly requestedScope: RequestedScope; readonly plannedScope: PlannedScope } => {
  const root = projectModel(model)
  let plannedArea = "area:root"
  const requested: { areaId?: string; factorFilter?: ReadonlyArray<string> } = {}
  if (area !== "") {
    parseAreaReference(model, area)
    requested.areaId = area
    plannedArea = area
  }
  const filter: Array<string> = []
  for (const reference of factors) {
    const element = findElement(root, reference)
    if (element?.kind !== "factor") {
      throw new Error(
        `--factor: factor model reference ${JSON.stringify(reference)} does not resolve in the model`,
      )
    }
    const owner = factorArea(reference)
    if (requested.areaId === undefined) requested.areaId = owner
    if (plannedArea === "area:root" && area === "" && filter.length === 0) plannedArea = owner
    if (requested.areaId !== owner || plannedArea !== owner) {
      throw new Error(`--factor ${reference} does not belong to --area ${plannedArea}`)
    }
    filter.push(reference)
  }
  if (filter.length > 0) requested.factorFilter = filter
  return {
    requestedScope: requested,
    plannedScope: { areaId: plannedArea, factorFilter: filter },
  }
}

export const scopeSlug = (scope: PlannedScope) => {
  if (scope.areaId === "area:root" && scope.factorFilter.length === 0) return "full"
  const area = scope.areaId.slice(5)
  const parts = area === "root" ? ["root"] : area.split("/")
  for (const factor of scope.factorFilter) parts.push(...factor.split("::")[1]!.split("/"))
  return parts.join("-")
}
