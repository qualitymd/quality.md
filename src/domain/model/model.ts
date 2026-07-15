import type { QualityDocument } from "./document.ts"

export interface RatingLevel {
  readonly level: string
  readonly title?: string
  readonly description?: string
  readonly criterion: string
}

export interface Requirement {
  readonly title: string
  readonly description?: string
  readonly assessment: string
  readonly factors?: ReadonlyArray<string>
  readonly ratings?: Readonly<Record<string, string>>
}

export interface Factor {
  readonly title?: string
  readonly description?: string
  readonly factors?: Readonly<Record<string, Factor>>
  readonly requirements?: Readonly<Record<string, Requirement>>
}

export interface Area {
  readonly title?: string
  readonly factors?: Readonly<Record<string, Factor>>
  readonly requirements?: Readonly<Record<string, Requirement>>
  readonly areas?: Readonly<Record<string, Area>>
  readonly source?: string
}

export interface QualityModel extends Area {
  readonly ratingScale: ReadonlyArray<RatingLevel>
  readonly path: string
}

export const decodeModel = (document: QualityDocument): QualityModel => {
  const value = document.yaml.toJS() as Omit<QualityModel, "path">
  return { ...value, path: document.path }
}

export type ElementKind = "area" | "factor" | "requirement"

export interface ModelElement {
  readonly id: string
  readonly kind: ElementKind
  readonly label: string
  readonly parentId?: string
  readonly children?: ReadonlyArray<ModelElement>
}

const entries = <A>(record: Readonly<Record<string, A>> | undefined) =>
  Object.entries(record ?? {}).sort(([left], [right]) => left.localeCompare(right))

const areaReference = (path: ReadonlyArray<string>) =>
  `area:${path.length === 0 ? "root" : path.join("/")}`

export const factorReference = (
  areaPath: ReadonlyArray<string>,
  factorPath: ReadonlyArray<string>,
) => `factor:${areaPath.length === 0 ? "root" : areaPath.join("/")}::${factorPath.join("/")}`

export const requirementReference = (areaPath: ReadonlyArray<string>, name: string) =>
  `requirement:${areaPath.length === 0 ? "root" : areaPath.join("/")}::${name}`

const buildFactor = (
  parentId: string,
  areaPath: ReadonlyArray<string>,
  factorPath: ReadonlyArray<string>,
  factor: Factor,
): ModelElement => {
  const id = factorReference(areaPath, factorPath)
  const children = [
    ...entries(factor.factors).map(([name, child]) =>
      buildFactor(id, areaPath, [...factorPath, name], child),
    ),
    ...entries(factor.requirements).map(([name, requirement]) => ({
      id: requirementReference(areaPath, name),
      kind: "requirement" as const,
      label: requirement.title || name,
      parentId: id,
    })),
  ]
  return {
    id,
    kind: "factor",
    label: factor.title || factorPath.at(-1) || "root",
    parentId,
    ...(children.length === 0 ? {} : { children }),
  }
}

const buildArea = (
  parentId: string | undefined,
  areaPath: ReadonlyArray<string>,
  area: Area,
  fallback: string,
): ModelElement => {
  const id = areaReference(areaPath)
  const children = [
    ...entries(area.factors).map(([name, factor]) => buildFactor(id, areaPath, [name], factor)),
    ...entries(area.requirements).map(([name, requirement]) => ({
      id: requirementReference(areaPath, name),
      kind: "requirement" as const,
      label: requirement.title || name,
      parentId: id,
    })),
    ...entries(area.areas).map(([name, child]) => buildArea(id, [...areaPath, name], child, name)),
  ]
  return {
    id,
    kind: "area",
    label: area.title || fallback,
    ...(parentId === undefined ? {} : { parentId }),
    ...(children.length === 0 ? {} : { children }),
  }
}

export const projectModel = (model: QualityModel): ModelElement =>
  buildArea(undefined, [], model, "root")

export const flattenModel = (root: ModelElement): ReadonlyArray<ModelElement> => [
  root,
  ...(root.children ?? []).flatMap(flattenModel),
]

export const findElement = (root: ModelElement, id: string) =>
  flattenModel(root).find((entry) => entry.id === id)

export const truncateDepth = (element: ModelElement, depth: number): ModelElement => {
  if (depth === 0 || element.children === undefined) {
    const { children: _, ...withoutChildren } = element
    return withoutChildren
  }
  return {
    ...element,
    children: element.children.map((child) => truncateDepth(child, depth - 1)),
  }
}

export const parseAreaReference = (
  model: QualityModel,
  reference: string,
): ReadonlyArray<string> => {
  if (!reference.startsWith("area:")) {
    throw new Error(
      `area model reference ${JSON.stringify(reference)} must start with area: prefix`,
    )
  }
  const raw = reference.slice(5)
  const path = raw === "root" ? [] : raw.split("/")
  if (raw === "" || path.some((part) => !/^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(part))) {
    throw new Error(`area model reference ${JSON.stringify(reference)} is invalid`)
  }
  if (findElement(projectModel(model), areaReference(path)) === undefined) {
    throw new Error(
      `area model reference ${JSON.stringify(reference)} does not resolve in the model`,
    )
  }
  return path
}

export const areaId = areaReference

export type SourceState = "declared" | "inherited" | "default"

export const effectiveSource = (
  model: QualityModel,
  path: ReadonlyArray<string>,
): { readonly selector: string; readonly state: SourceState } => {
  let selector = model.source ?? ""
  let declaredAt = selector === "" ? -1 : 0
  let areas = model.areas ?? {}
  for (const [index, name] of path.entries()) {
    const area = areas[name]
    if (area === undefined) break
    if (area.source !== undefined && area.source !== "") {
      selector = area.source
      declaredAt = index + 1
    }
    areas = area.areas ?? {}
  }
  if (declaredAt === path.length) return { selector, state: "declared" }
  if (declaredAt >= 0) return { selector, state: "inherited" }
  return { selector: ".", state: "default" }
}
