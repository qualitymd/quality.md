import * as Effect from "effect/Effect"
import * as FileSystem from "effect/FileSystem"

import { FileSystemFailure } from "../domain/errors.ts"
import { commandResult, ExitCode } from "../domain/command-result.ts"
import { jsonDocument } from "../domain/json.ts"
import {
  decodeModel,
  findElement,
  flattenModel,
  parseAreaReference,
  projectModel,
  truncateDepth,
  type ElementKind,
  type ModelElement,
} from "../domain/model/model.ts"
import { parseQualityDocument } from "../domain/model/document.ts"

interface CommonInput {
  readonly path: string
  readonly json: boolean
}

export interface ModelTreeInput extends CommonInput {
  readonly area: string
  readonly depth?: number
}

export interface ModelListInput extends CommonInput {
  readonly area: string
  readonly types: ReadonlyArray<string>
}

export interface ModelGetInput extends CommonInput {
  readonly id: string
}

const readModel = (path: string) =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const raw = yield* fs.readFileString(path)
    return decodeModel(parseQualityDocument(path, raw))
  }).pipe(
    Effect.mapError(
      (cause) =>
        new FileSystemFailure({
          detail: `${cause instanceof Error ? cause.message : String(cause)}\n\nmodel needs a parseable model; run \`qualitymd lint ${path}\` for diagnostics`,
        }),
    ),
  )

const usage = (detail: string) =>
  commandResult("", { stderr: `qualitymd: ${detail}\n`, exitCode: ExitCode.usage })

const json = jsonDocument

const rootAt = (model: ReturnType<typeof decodeModel>, area: string) => {
  const root = projectModel(model)
  if (area === "") return root
  const path = parseAreaReference(model, area)
  return findElement(root, `area:${path.length === 0 ? "root" : path.join("/")}`)!
}

const renderTree = (root: ModelElement): string => {
  const lines: Array<string> = []
  const walk = (element: ModelElement, depth: number) => {
    lines.push(`${"  ".repeat(depth)}${element.id}  ${element.label}`)
    for (const child of element.children ?? []) walk(child, depth + 1)
  }
  walk(root, 0)
  return `${lines.join("\n")}\n`
}

export const modelTreeCommand = (input: ModelTreeInput) => {
  if (input.path === "-")
    return Effect.succeed(usage("model does not read from stdin; pass a file path"))
  if (input.depth !== undefined && input.depth < 0)
    return Effect.succeed(usage("--depth must be zero or greater"))
  return readModel(input.path).pipe(
    Effect.map((model) => {
      try {
        const root = truncateDepth(rootAt(model, input.area), input.depth ?? -1)
        return commandResult(input.json ? json(root) : renderTree(root))
      } catch (cause) {
        return usage(cause instanceof Error ? cause.message : String(cause))
      }
    }),
  )
}

const parseKinds = (types: ReadonlyArray<string>): Set<ElementKind> | string => {
  const kinds = new Set<ElementKind>()
  for (const raw of types.flatMap((entry) => entry.split(","))) {
    if (raw !== "area" && raw !== "factor" && raw !== "requirement") {
      return `--type ${JSON.stringify(raw)} is not one of: area, factor, requirement`
    }
    kinds.add(raw)
  }
  return kinds
}

export const modelListCommand = (input: ModelListInput) => {
  if (input.path === "-")
    return Effect.succeed(usage("model does not read from stdin; pass a file path"))
  const kinds = parseKinds(input.types)
  if (typeof kinds === "string") return Effect.succeed(usage(kinds))
  return readModel(input.path).pipe(
    Effect.map((model) => {
      try {
        const rows = flattenModel(rootAt(model, input.area))
          .filter((entry) => kinds.size === 0 || kinds.has(entry.kind))
          .map(({ id, kind, label, parentId }) => ({
            id,
            kind,
            label,
            ...(parentId === undefined ? {} : { parentId }),
          }))
        const human =
          rows.length === 0
            ? "No elements.\n"
            : `${rows.map((row) => `${row.id}  ${row.label}`).join("\n")}\n`
        return commandResult(input.json ? json(rows) : human)
      } catch (cause) {
        return usage(cause instanceof Error ? cause.message : String(cause))
      }
    }),
  )
}

const detail = (element: ModelElement) => {
  const grouped: Record<ElementKind, Array<string>> = { area: [], factor: [], requirement: [] }
  for (const child of element.children ?? []) grouped[child.kind].push(child.id)
  return {
    id: element.id,
    kind: element.kind,
    label: element.label,
    ...(element.parentId === undefined ? {} : { parentId: element.parentId }),
    ...(grouped.factor.length === 0 ? {} : { factors: grouped.factor }),
    ...(grouped.requirement.length === 0 ? {} : { requirements: grouped.requirement }),
    ...(grouped.area.length === 0 ? {} : { areas: grouped.area }),
  }
}

const renderDetail = (value: ReturnType<typeof detail>) => {
  let output = `${value.id}\n  kind:   ${value.kind}\n  label:  ${value.label}\n`
  if (value.parentId !== undefined) output += `  parent: ${value.parentId}\n`
  for (const [label, ids] of [
    ["factors", value.factors],
    ["requirements", value.requirements],
    ["areas", value.areas],
  ] as const) {
    if (ids === undefined) continue
    output += `  ${label}:\n${ids.map((id) => `    ${id}`).join("\n")}\n`
  }
  return output
}

export const modelGetCommand = (input: ModelGetInput) => {
  if (input.path === "-")
    return Effect.succeed(usage("model does not read from stdin; pass a file path"))
  return readModel(input.path).pipe(
    Effect.map((model) => {
      const root = projectModel(model)
      const element = findElement(root, input.id)
      if (element === undefined)
        return usage(`no element in the model has id ${JSON.stringify(input.id)}`)
      const value = detail(element)
      return commandResult(input.json ? json(value) : renderDetail(value))
    }),
  )
}
