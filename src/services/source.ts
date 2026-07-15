import * as Effect from "effect/Effect"
import * as FileSystem from "effect/FileSystem"
import * as Path from "effect/Path"

import { hashJson, sha256 } from "../domain/json.ts"
import type { SourceKind } from "../domain/evaluator/types.ts"

type JsonObject = Record<string, unknown>

export interface AreaSource {
  readonly selector: string
  readonly kind: SourceKind
}

export interface SealedEvidenceManifest extends JsonObject {
  readonly requirementId: string
  readonly source: AreaSource
  readonly observations: ReadonlyArray<JsonObject>
  readonly limits: ReadonlyArray<string>
  readonly capturedAt: string
  readonly manifestHash: string
}

export class EvidenceValidationError extends Error {}

const validGlob = (value: string) => {
  for (let index = 0; index < value.length; index += 1) {
    if (value[index] === "\\") {
      index += 1
      if (index >= value.length) return false
      continue
    }
    if (value[index] !== "[") continue
    let cursor = index + 1
    if (value[cursor] === "!" || value[cursor] === "^") cursor += 1
    const first = cursor
    let closed = false
    for (; cursor < value.length; cursor += 1) {
      if (value[cursor] === "\\") {
        cursor += 1
        if (cursor >= value.length) return false
      } else if (value[cursor] === "]") {
        closed = true
        break
      }
    }
    if (!closed || cursor === first) return false
    index = cursor
  }
  return true
}

export const detectSourceKind = (workspaceRoot: string, selector: string) =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const paths = yield* Path.Path
    const value = selector.trim()
    const normalized = value.replaceAll("\\", "/")
    if (
      value === "" ||
      paths.isAbsolute(value) ||
      normalized === ".." ||
      normalized.startsWith("../")
    )
      return "path" as const
    if (/[*?[]/.test(value) && validGlob(value)) return "glob" as const
    if (yield* fs.exists(paths.join(workspaceRoot, value))) return "path" as const
    return "prose" as const
  })

export const validateSourceSelector = (
  workspaceRoot: string,
  source: AreaSource,
): Effect.Effect<void, EvidenceValidationError, Path.Path> =>
  Effect.gen(function* () {
    if (source.kind === "prose") return
    const paths = yield* Path.Path
    const selector = source.selector.trim().replaceAll("\\", "/")
    if (
      selector === "" ||
      paths.isAbsolute(selector) ||
      selector === ".." ||
      selector.startsWith("../") ||
      selector.split("/").includes("..")
    ) {
      yield* Effect.fail(
        new EvidenceValidationError(
          `source selector ${JSON.stringify(source.selector)} must be workspace-relative and contained`,
        ),
      )
    }
    const resolved = paths.resolve(workspaceRoot, selector)
    const relative = paths.relative(workspaceRoot, resolved)
    if (relative === ".." || relative.startsWith(`..${paths.sep}`)) {
      yield* Effect.fail(
        new EvidenceValidationError(
          `source selector ${JSON.stringify(source.selector)} escapes the workspace`,
        ),
      )
    }
  })

const object = (value: unknown, detail: string): JsonObject => {
  if (value === null || Array.isArray(value) || typeof value !== "object") {
    throw new EvidenceValidationError(detail)
  }
  return value as JsonObject
}

const exactKeys = (value: JsonObject, allowed: ReadonlyArray<string>, detail: string) => {
  const extra = Object.keys(value).filter((key) => !allowed.includes(key))
  if (extra.length > 0)
    throw new EvidenceValidationError(`${detail} carries unknown field ${JSON.stringify(extra[0])}`)
}

const normalizedPath = (paths: Path.Path, value: unknown, detail: string) => {
  if (typeof value !== "string" || value.trim() === "")
    throw new EvidenceValidationError(`${detail} must be a non-empty workspace-relative path`)
  const normalized = value.trim().replaceAll("\\", "/")
  if (
    paths.isAbsolute(normalized) ||
    normalized === ".." ||
    normalized.startsWith("../") ||
    normalized.split("/").includes("..")
  )
    throw new EvidenceValidationError(`${detail} must be workspace-relative and contained`)
  return normalized.replace(/^\.\//, "")
}

const normalizeLocator = (value: unknown, content: string, detail: string) => {
  if (value === undefined) return undefined
  const locator = object(value, `${detail} must be an object`)
  exactKeys(locator, ["startLine", "endLine", "heading"], detail)
  const heading = locator.heading
  const hasHeading = typeof heading === "string" && heading.trim() !== ""
  const hasLines = locator.startLine !== undefined || locator.endLine !== undefined
  if (hasHeading && hasLines)
    throw new EvidenceValidationError(`${detail} must use a heading or line range, not both`)
  if (hasHeading) {
    const wanted = heading.trim()
    const found = content
      .split(/\r?\n/)
      .some((line) => /^#{1,6}\s+/.test(line) && line.replace(/^#{1,6}\s+/, "").trim() === wanted)
    if (!found)
      throw new EvidenceValidationError(`${detail}.heading ${JSON.stringify(wanted)} was not found`)
    return { heading: wanted }
  }
  if (!hasLines) throw new EvidenceValidationError(`${detail} must carry a heading or line range`)
  const start = locator.startLine
  const end = locator.endLine ?? start
  if (!Number.isInteger(start) || (start as number) < 1)
    throw new EvidenceValidationError(`${detail}.startLine must be a positive integer`)
  if (!Number.isInteger(end) || (end as number) < (start as number))
    throw new EvidenceValidationError(`${detail}.endLine must be at least startLine`)
  const lineCount = content.split(/\r?\n/).length
  if ((end as number) > lineCount)
    throw new EvidenceValidationError(`${detail}.endLine exceeds the file's ${lineCount} lines`)
  return { startLine: start as number, endLine: end as number }
}

const evaluatedBySource = (
  paths: Path.Path,
  workspaceRoot: string,
  source: AreaSource,
  evidencePath: string,
) => {
  if (source.kind === "prose") return true
  if (source.kind === "glob") {
    const selector = source.selector.replaceAll("\\", "/").replace(/^\.\//, "")
    return new Bun.Glob(selector).match(evidencePath)
  }
  const selected = paths.resolve(workspaceRoot, source.selector)
  const evidence = paths.resolve(workspaceRoot, evidencePath)
  const relative = paths.relative(selected, evidence)
  return relative === "" || (relative !== ".." && !relative.startsWith(`..${paths.sep}`))
}

const evidenceReferences = (assessment: unknown) => {
  const refs: Array<string> = []
  const visit = (value: unknown) => {
    if (Array.isArray(value)) {
      for (const entry of value) visit(entry)
      return
    }
    if (value === null || typeof value !== "object") return
    for (const [key, entry] of Object.entries(value as JsonObject)) {
      if (key === "sourceRef" && typeof entry === "string") refs.push(entry)
      else visit(entry)
    }
  }
  visit(assessment)
  return refs
}

export const sealEvidenceManifest = (options: {
  readonly workspaceRoot: string
  readonly requirementId: string
  readonly source: AreaSource
  readonly proposal: unknown
  readonly assessment: unknown
  readonly capturedAt: string
}): Effect.Effect<
  SealedEvidenceManifest,
  EvidenceValidationError,
  FileSystem.FileSystem | Path.Path
> =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const paths = yield* Path.Path
    yield* validateSourceSelector(options.workspaceRoot, options.source)
    const proposal = object(options.proposal, "evidence must be an object")
    exactKeys(proposal, ["observations", "limits"], "evidence")
    if (!Array.isArray(proposal.observations))
      throw new EvidenceValidationError("evidence.observations must be an array")
    if (!Array.isArray(proposal.limits))
      throw new EvidenceValidationError("evidence.limits must be an array")
    const limits = proposal.limits.map((value, index) => {
      if (typeof value !== "string" || value.trim() === "")
        throw new EvidenceValidationError(`evidence.limits[${index}] must be a non-empty string`)
      return value.trim()
    })
    const workspaceReal = yield* fs.realPath(options.workspaceRoot)
    const seen = new Set<string>()
    const observations: Array<JsonObject> = []
    for (const [index, raw] of proposal.observations.entries()) {
      const detail = `evidence.observations[${index}]`
      const observation = object(raw, `${detail} must be an object`)
      exactKeys(observation, ["id", "kind", "role", "path", "locator"], detail)
      if (typeof observation.id !== "string" || !/^ev-[a-z0-9][a-z0-9-]*$/.test(observation.id))
        throw new EvidenceValidationError(`${detail}.id must match ^ev-[a-z0-9][a-z0-9-]*$`)
      if (seen.has(observation.id))
        throw new EvidenceValidationError(
          `${detail}.id duplicates ${JSON.stringify(observation.id)}`,
        )
      seen.add(observation.id)
      if (observation.kind !== "file")
        throw new EvidenceValidationError(`${detail}.kind must be file`)
      if (observation.role !== "evaluated" && observation.role !== "supporting")
        throw new EvidenceValidationError(`${detail}.role must be evaluated or supporting`)
      const path = normalizedPath(paths, observation.path, `${detail}.path`)
      const absolute = paths.resolve(options.workspaceRoot, path)
      if (!(yield* fs.exists(absolute)))
        throw new EvidenceValidationError(`${detail}.path ${JSON.stringify(path)} does not exist`)
      const real = yield* fs.realPath(absolute)
      const relativeReal = paths.relative(workspaceReal, real)
      if (relativeReal === ".." || relativeReal.startsWith(`..${paths.sep}`))
        throw new EvidenceValidationError(
          `${detail}.path ${JSON.stringify(path)} escapes the workspace`,
        )
      if ((yield* fs.stat(real)).type !== "File")
        throw new EvidenceValidationError(
          `${detail}.path ${JSON.stringify(path)} must name a regular file`,
        )
      const bytes = yield* fs.readFile(real)
      let content: string
      try {
        content = new TextDecoder("utf-8", { fatal: true }).decode(bytes)
      } catch {
        throw new EvidenceValidationError(
          `${detail}.path ${JSON.stringify(path)} must be UTF-8 text`,
        )
      }
      if (
        observation.role === "evaluated" &&
        !evaluatedBySource(paths, options.workspaceRoot, options.source, path)
      )
        throw new EvidenceValidationError(
          `${detail}.path ${JSON.stringify(path)} is outside the evaluated source selector; classify it as supporting or correct the path`,
        )
      const locator = normalizeLocator(observation.locator, content, `${detail}.locator`)
      observations.push({
        id: observation.id,
        kind: "file",
        role: observation.role,
        path,
        ...(locator === undefined ? {} : { locator }),
        sha256: yield* Effect.promise(() => sha256(bytes)),
        bytes: bytes.length,
        capturedAt: options.capturedAt,
      })
    }
    for (const ref of evidenceReferences(options.assessment)) {
      const match = /^evidence\[([^\]]+)\]$/.exec(ref)
      if (match === null || !seen.has(match[1]!))
        throw new EvidenceValidationError(
          `assessment evidence sourceRef ${JSON.stringify(ref)} does not name an accepted evidence observation`,
        )
    }
    const manifest = {
      requirementId: options.requirementId,
      source: options.source,
      observations,
      limits,
      capturedAt: options.capturedAt,
    }
    return {
      ...manifest,
      manifestHash: yield* Effect.promise(() => hashJson(manifest)),
    } satisfies SealedEvidenceManifest
  }).pipe(
    Effect.catchDefect((defect) =>
      Effect.fail(
        defect instanceof EvidenceValidationError
          ? defect
          : new EvidenceValidationError(defect instanceof Error ? defect.message : String(defect)),
      ),
    ),
    Effect.mapError((cause) =>
      cause instanceof EvidenceValidationError
        ? cause
        : new EvidenceValidationError(cause instanceof Error ? cause.message : String(cause)),
    ),
  )
