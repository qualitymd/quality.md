import {
  effectiveSource,
  factorReference,
  requirementReference,
  type Area,
  type Factor,
  type QualityModel,
  type Requirement,
} from "../model/model.ts"
import type { PlannedScope } from "./run.ts"

export interface PlannedRequirement {
  readonly ref: string
  readonly areaId: string
  readonly factorIds: ReadonlyArray<string>
  readonly value: Requirement
}

export interface PlannedFactor {
  readonly ref: string
  readonly areaId: string
  readonly path: ReadonlyArray<string>
  readonly value: Factor
}

export interface PlannedArea {
  readonly ref: string
  readonly path: ReadonlyArray<string>
  readonly value: Area
  readonly source: string
  readonly childAreaIds: ReadonlyArray<string>
  readonly rootFactorIds: ReadonlyArray<string>
  readonly localRequirementIds: ReadonlyArray<string>
}

export interface EvaluationPlan {
  readonly areas: ReadonlyArray<PlannedArea>
  readonly factors: ReadonlyArray<PlannedFactor>
  readonly requirements: ReadonlyArray<PlannedRequirement>
}

const sorted = <A>(record: Readonly<Record<string, A>> | undefined) =>
  Object.entries(record ?? {}).sort(([left], [right]) => left.localeCompare(right))

const areaRef = (path: ReadonlyArray<string>) =>
  `area:${path.length === 0 ? "root" : path.join("/")}`

const areaAt = (model: QualityModel, path: ReadonlyArray<string>): Area => {
  let area: Area = model
  for (const name of path) area = area.areas![name]!
  return area
}

export const planEvaluation = (model: QualityModel, scope: PlannedScope): EvaluationPlan => {
  const scopedPath = scope.areaId === "area:root" ? [] : scope.areaId.slice(5).split("/")
  const filters = scope.factorFilter

  const factorSelected = (reference: string) =>
    filters.length === 0 ||
    filters.some((filter) => reference === filter || reference.startsWith(`${filter}/`))

  const walkFactor = (
    areaPath: ReadonlyArray<string>,
    path: ReadonlyArray<string>,
    value: Factor,
  ): Pick<EvaluationPlan, "factors" | "requirements"> => {
    const reference = factorReference(areaPath, path)
    if (!factorSelected(reference)) return { factors: [], requirements: [] }
    const children = sorted(value.factors).map(([name, child]) =>
      walkFactor(areaPath, [...path, name], child),
    )
    const directRequirements = sorted(value.requirements).map(([name, requirement]) => {
      const linked = (requirement.factors ?? []).map((factor) =>
        factor.startsWith("factor:") ? factor : factorReference(areaPath, factor.split("/")),
      )
      return {
        ref: requirementReference(areaPath, name),
        areaId: areaRef(areaPath),
        factorIds: [...new Set([reference, ...linked.filter(factorSelected)])],
        value: requirement,
      }
    })
    return {
      factors: [
        { ref: reference, areaId: areaRef(areaPath), path, value },
        ...children.flatMap((child) => child.factors),
      ],
      requirements: [...children.flatMap((child) => child.requirements), ...directRequirements],
    }
  }

  const walkArea = (path: ReadonlyArray<string>, value: Area): EvaluationPlan => {
    const factorPlans = sorted(value.factors).map(([name, factor]) =>
      walkFactor(path, [name], factor),
    )
    const factors = factorPlans.flatMap((plan) => plan.factors)
    const factorRequirements = factorPlans.flatMap((plan) => plan.requirements)
    const areaRequirements =
      filters.length === 0
        ? sorted(value.requirements).map(([name, requirement]) => ({
            ref: requirementReference(path, name),
            areaId: areaRef(path),
            factorIds: (requirement.factors ?? []).map((factor) =>
              factor.startsWith("factor:") ? factor : factorReference(path, factor.split("/")),
            ),
            value: requirement,
          }))
        : []
    const requirements = [...factorRequirements, ...areaRequirements]
    const childAreaIds = sorted(value.areas).map(([name]) => areaRef([...path, name]))
    const area: PlannedArea = {
      ref: areaRef(path),
      path,
      value,
      source: effectiveSource(model, path).selector,
      childAreaIds,
      rootFactorIds: factors
        .filter((factor) => factor.areaId === areaRef(path) && factor.path.length === 1)
        .map((factor) => factor.ref),
      localRequirementIds: requirements.map((requirement) => requirement.ref),
    }
    const children = sorted(value.areas).map(([name, child]) => walkArea([...path, name], child))
    return {
      areas: [area, ...children.flatMap((child) => child.areas)],
      factors: [...factors, ...children.flatMap((child) => child.factors)],
      requirements: [...requirements, ...children.flatMap((child) => child.requirements)],
    }
  }

  return walkArea(scopedPath, areaAt(model, scopedPath))
}
