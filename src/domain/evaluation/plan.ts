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
  const areas: Array<PlannedArea> = []
  const factors: Array<PlannedFactor> = []
  const requirements: Array<PlannedRequirement> = []
  const filters = scope.factorFilter

  const factorSelected = (reference: string) =>
    filters.length === 0 ||
    filters.some((filter) => reference === filter || reference.startsWith(`${filter}/`))

  const walkFactor = (
    areaPath: ReadonlyArray<string>,
    path: ReadonlyArray<string>,
    value: Factor,
  ) => {
    const reference = factorReference(areaPath, path)
    if (!factorSelected(reference)) return
    factors.push({ ref: reference, areaId: areaRef(areaPath), path, value })
    for (const [name, child] of sorted(value.factors)) walkFactor(areaPath, [...path, name], child)
    for (const [name, requirement] of sorted(value.requirements)) {
      const ref = requirementReference(areaPath, name)
      const linked = (requirement.factors ?? []).map((factor) =>
        factor.startsWith("factor:") ? factor : factorReference(areaPath, factor.split("/")),
      )
      requirements.push({
        ref,
        areaId: areaRef(areaPath),
        factorIds: [...new Set([reference, ...linked.filter(factorSelected)])],
        value: requirement,
      })
    }
  }

  const walkArea = (path: ReadonlyArray<string>, value: Area) => {
    const requirementStart = requirements.length
    for (const [name, factor] of sorted(value.factors)) walkFactor(path, [name], factor)
    if (filters.length === 0) {
      for (const [name, requirement] of sorted(value.requirements)) {
        requirements.push({
          ref: requirementReference(path, name),
          areaId: areaRef(path),
          factorIds: (requirement.factors ?? []).map((factor) =>
            factor.startsWith("factor:") ? factor : factorReference(path, factor.split("/")),
          ),
          value: requirement,
        })
      }
    }
    const factorStart = factors.findIndex((factor) => factor.areaId === areaRef(path))
    const childAreaIds = sorted(value.areas).map(([name]) => areaRef([...path, name]))
    areas.push({
      ref: areaRef(path),
      path,
      value,
      source: effectiveSource(model, path).selector,
      childAreaIds,
      rootFactorIds:
        factorStart < 0
          ? []
          : factors
              .slice(factorStart)
              .filter((factor) => factor.areaId === areaRef(path) && factor.path.length === 1)
              .map((factor) => factor.ref),
      localRequirementIds: requirements
        .slice(requirementStart)
        .map((requirement) => requirement.ref),
    })
    for (const [name, child] of sorted(value.areas)) walkArea([...path, name], child)
  }

  walkArea(scopedPath, areaAt(model, scopedPath))
  return { areas, factors, requirements }
}
