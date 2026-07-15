import * as Effect from "effect/Effect"
import * as FileSystem from "effect/FileSystem"
import * as Path from "effect/Path"
import * as Result from "effect/Result"

import { FileSystemFailure } from "../domain/errors.ts"
import { commandResult, ExitCode, type CommandResult } from "../domain/command-result.ts"
import {
  assignRecommendationId,
  dataKinds,
  dataPathForPayload,
  dataPathForQuery,
  evaluationDataExample,
  evaluationDataSchema,
  resolveDataKind,
  validateDataPayload,
  type DataKind,
  type DataQuery,
} from "../domain/evaluation/data.ts"
import { jsonDocument } from "../domain/json.ts"
import { decodeModel } from "../domain/model/model.ts"
import { parseQualityDocument } from "../domain/model/document.ts"
import { HostRuntime } from "../services/host-runtime.ts"
import { resolveWorkspace } from "../services/workspace.ts"
import { evaluationRunDirectories } from "./evaluation-runs.ts"

type RunInput = {
  readonly run?: string
  readonly latest: boolean
  readonly model: string
  readonly evaluationDir: string
}

interface ResolvedRun {
  readonly absolute: string
  readonly display: string
}

const usage = (detail: string) =>
  commandResult("", { stderr: `qualitymd: ${detail}\n`, exitCode: ExitCode.usage })

const resolveRun = (input: RunInput) =>
  Effect.gen(function* () {
    const paths = yield* Path.Path
    if (input.run !== undefined && input.latest)
      throw new Error("pass a run path or --latest, not both")
    if (input.run === undefined && !input.latest) throw new Error("pass a run path or --latest")
    if (input.run !== undefined) {
      if (input.model === "")
        return { absolute: paths.resolve(input.run), display: input.run } satisfies ResolvedRun
      const workspace = yield* resolveWorkspace({ model: input.model })
      return {
        absolute: paths.isAbsolute(input.run)
          ? input.run
          : paths.resolve(workspace.workspaceRoot.abs, input.run),
        display: input.run,
      } satisfies ResolvedRun
    }
    const workspace = yield* resolveWorkspace({
      ...(input.model === "" ? {} : { model: input.model }),
      ...(input.evaluationDir === "" ? {} : { evaluationDir: input.evaluationDir }),
    })
    const latest = (yield* evaluationRunDirectories(
      workspace.evaluations.abs,
      workspace.evaluations.rel,
    )).at(-1)
    if (latest === undefined) throw new Error("no evaluation runs found")
    return {
      absolute: latest.absolute,
      display: latest.display,
    } satisfies ResolvedRun
  })

const mapFailure = <A, R>(
  effect: Effect.Effect<A, unknown, R>,
): Effect.Effect<A, FileSystemFailure, R> =>
  effect.pipe(
    Effect.mapError(
      (cause) =>
        new FileSystemFailure({ detail: cause instanceof Error ? cause.message : String(cause) }),
    ),
  )

const recursiveJson = (directory: string) =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const paths = yield* Path.Path
    const visit = (current: string): Effect.Effect<ReadonlyArray<string>, unknown> =>
      Effect.gen(function* () {
        const entries = yield* Effect.forEach(yield* fs.readDirectory(current), (name) =>
          Effect.gen(function* () {
            const path = paths.join(current, name)
            const info = yield* fs.stat(path)
            if (info.type === "Directory") return yield* visit(path)
            return info.type === "File" && name.endsWith(".json") ? [path] : []
          }),
        )
        return entries.flat()
      })
    return [...(yield* visit(directory))].sort()
  })

const runModel = (run: string) =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const paths = yield* Path.Path
    const snapshot = paths.join(run, "model-snapshot.md")
    return decodeModel(parseQualityDocument(snapshot, yield* fs.readFileString(snapshot)))
  })

export const validateCrossPayloads = (payloads: ReadonlyArray<Record<string, unknown>>) => {
  const object = (value: unknown): Record<string, unknown> =>
    value !== null && !Array.isArray(value) && typeof value === "object"
      ? (value as Record<string, unknown>)
      : {}
  const objects = (value: unknown) => (Array.isArray(value) ? value.map(object) : [])
  const payloadPaths = payloads.map((payload) => ({
    path: dataPathForPayload(resolveDataKind(String(payload.kind)), payload),
    payload,
  }))
  const duplicatePath = payloadPaths.find(
    (entry, index) =>
      payloadPaths.findIndex((candidate) => candidate.path === entry.path) !== index,
  )
  if (duplicatePath !== undefined)
    throw new Error(`effective payloads duplicate ${duplicatePath.path}`)
  const paths = new Map(payloadPaths.map(({ path, payload }) => [path, payload]))
  const assessments = new Map(
    payloads
      .filter((payload) => payload.kind === "RequirementAssessmentResult")
      .map((payload) => [payload.requirementId, payload]),
  )
  const recommendationList = payloads.filter((payload) => payload.kind === "RecommendationResult")
  const recommendations = new Set(recommendationList.map((payload) => payload.id))
  if (recommendations.size !== recommendationList.length)
    throw new Error("RecommendationResult.id must be unique across the run")
  const findingKeys = new Set(
    [...assessments].flatMap(([requirementId, assessment]) =>
      objects(assessment.findings).map(
        (finding) => `${String(requirementId)}#${String(finding.id)}`,
      ),
    ),
  )
  const refPath = (raw: unknown) => {
    const ref = object(raw)
    const subject = object(ref.subject)
    const kind = resolveDataKind(String(ref.kind))
    return dataPathForQuery({
      kind,
      area: String(subject.areaId ?? ""),
      factor: String(subject.factorId ?? ""),
      requirement: String(subject.requirementId ?? ""),
      selector: String(ref.selector ?? ""),
    })
  }
  const findingKey = (raw: unknown) => {
    const ref = object(raw)
    if (ref.kind !== "RequirementAssessmentResult")
      throw new Error("findingRef.kind must be RequirementAssessmentResult")
    const requirement = String(object(ref.subject).requirementId ?? "")
    const match = /^findings\[([^\]]+)]$/.exec(String(ref.selector ?? ""))
    if (match === null) throw new Error("findingRef.selector must name findings[<id>]")
    return `${requirement}#${match[1]}`
  }
  const validateRanks = (entries: ReadonlyArray<Record<string, unknown>>, name: string) => {
    for (const [index, entry] of entries.entries()) {
      const rank = entry.rank
      if (!Number.isInteger(rank) || Number(rank) < 1)
        throw new Error(`${name}[${index}].rank must be a positive integer`)
      if (entries.slice(0, index).some((candidate) => Number(candidate.rank) === Number(rank)))
        throw new Error(`${name}[${index}].rank duplicates ${rank}`)
    }
  }
  const validateDrivers = (owner: Record<string, unknown>, label: string) => {
    for (const [index, driver] of objects(owner.ratingDrivers).entries()) {
      const refs = objects(driver.inputRefs)
      if (refs.length === 0)
        throw new Error(`${label}.ratingDrivers[${index}] requires at least one inputRefs entry`)
      for (const [refIndex, ref] of refs.entries()) {
        const path = refPath(ref)
        if (!paths.has(path))
          throw new Error(
            `${label}.ratingDrivers[${index}].inputRefs[${refIndex}] does not resolve to ${path}`,
          )
      }
    }
  }
  for (const payload of payloads) {
    if (payload.kind === "RequirementRatingResult" && payload.status === "rated") {
      const assessment = assessments.get(payload.requirementId)
      if (
        assessment === undefined ||
        !["assessed", "partially_assessed"].includes(String(assessment.status)) ||
        !Array.isArray(assessment.findings) ||
        assessment.findings.length === 0
      )
        throw new Error(
          "rated requirement requires paired RequirementAssessmentResult with at least one finding",
        )
      if (!Array.isArray(payload.ratingDrivers) || payload.ratingDrivers.length === 0)
        throw new Error("rated requirement requires at least one ratingDrivers entry")
      validateDrivers(payload, "RequirementRatingResult")
    }
    if (["FactorAnalysisResult", "AreaAnalysisResult"].includes(String(payload.kind)))
      for (const scopeName of ["localAnalysis", "localAndDescendantAnalysis"]) {
        const scope = object(payload[scopeName])
        if (
          scope.status === "analyzed" &&
          typeof scope.ratingLevelId === "string" &&
          (!Array.isArray(scope.ratingDrivers) || scope.ratingDrivers.length === 0)
        )
          throw new Error(
            `${scopeName} with a ratingLevelId requires at least one ratingDrivers entry`,
          )
        if (Array.isArray(scope.ratingDrivers)) validateDrivers(scope, scopeName)
      }
    if (payload.kind === "FindingRankingResult") {
      const entries = objects(payload.orderedFindings)
      const orderedKeys = entries.map((entry) => findingKey(entry.findingRef))
      for (const [index, key] of orderedKeys.entries()) {
        if (!findingKeys.has(key))
          throw new Error(
            `orderedFindings[${index}].findingRef does not resolve to an in-scope finding`,
          )
        if (orderedKeys.indexOf(key) !== index)
          throw new Error(`orderedFindings[${index}].findingRef duplicates ${key}`)
      }
      for (const key of findingKeys)
        if (!orderedKeys.includes(key)) throw new Error(`orderedFindings missing ${key}`)
      validateRanks(entries, "orderedFindings")
    }
    if (payload.kind === "RecommendationResult")
      for (const [index, ref] of objects(payload.traceRefs).entries()) {
        const path = refPath(ref)
        if (!paths.has(path)) throw new Error(`traceRefs[${index}] does not resolve to ${path}`)
      }
    if (
      payload.kind === "RecommendationRankingResult" &&
      Array.isArray(payload.orderedRecommendations)
    ) {
      const entries = objects(payload.orderedRecommendations)
      const ranked = entries.map((value) => String(value.recommendationRef ?? ""))
      for (const [index, recommendationRef] of ranked.entries()) {
        if (!recommendations.has(recommendationRef))
          throw new Error(
            `orderedRecommendations[${index}].recommendationRef does not resolve to a RecommendationResult`,
          )
        if (ranked.indexOf(recommendationRef) !== index)
          throw new Error(
            `orderedRecommendations[${index}].recommendationRef duplicates ${recommendationRef}`,
          )
      }
      for (const id of recommendations)
        if (!ranked.includes(String(id)))
          throw new Error(`orderedRecommendations missing ${String(id)}`)
      validateRanks(entries, "orderedRecommendations")
      const coverageEntries = objects(payload.findingCoverage)
      const covered = coverageEntries.map((coverage) => findingKey(coverage.findingRef))
      for (const [index, coverage] of coverageEntries.entries()) {
        const key = covered[index]!
        if (!findingKeys.has(key))
          throw new Error(
            `findingCoverage[${index}].findingRef does not resolve to an in-scope finding`,
          )
        if (covered.indexOf(key) !== index)
          throw new Error(`findingCoverage[${index}].findingRef duplicates ${key}`)
        if (
          coverage.disposition === "not_advice_driving" &&
          String(coverage.rationale ?? "") === ""
        )
          throw new Error(`findingCoverage[${index}].rationale is required for not_advice_driving`)
        if (coverage.disposition === "addressed_by_recommendation") {
          const refs = Array.isArray(coverage.recommendationRefs)
            ? coverage.recommendationRefs.map(String)
            : []
          if (refs.length === 0)
            throw new Error(
              `findingCoverage[${index}].recommendationRefs requires at least one RecommendationResult id`,
            )
          for (const ref of refs)
            if (!recommendations.has(ref))
              throw new Error(
                `findingCoverage[${index}].recommendationRefs includes unknown RecommendationResult ${ref}`,
              )
        }
      }
      for (const key of findingKeys)
        if (!covered.includes(key)) throw new Error(`findingCoverage missing ${key}`)
    }
  }
}

export const evaluationDataKindsCommand = (json: boolean): CommandResult => {
  const result = {
    schemaVersion: 3,
    kinds: dataKinds.map(([kind, agentWritable, description]) => ({
      kind,
      agentWritable,
      description,
    })),
  }
  if (json) return commandResult(jsonDocument(result))
  return commandResult(
    result.kinds
      .map(
        (kind) =>
          `${kind.kind}\t${kind.agentWritable ? "agent-writable" : "cli-owned"}\t${kind.description}\n`,
      )
      .join(""),
  )
}

export const evaluationDataSchemaCommand = (kind: string, jsonFlag: boolean): CommandResult => {
  if (jsonFlag)
    return usage("evaluation data schema already emits JSON on stdout; rerun without --json")
  try {
    return commandResult(
      jsonDocument(evaluationDataSchema(kind === "" ? undefined : resolveDataKind(kind))),
    )
  } catch (cause) {
    return usage(cause instanceof Error ? cause.message : String(cause))
  }
}

export const evaluationDataExampleCommand = (kind: string, jsonFlag: boolean): CommandResult => {
  if (jsonFlag)
    return usage("evaluation data example already emits JSON on stdout; rerun without --json")
  try {
    return commandResult(jsonDocument(evaluationDataExample(resolveDataKind(kind))))
  } catch (cause) {
    return usage(cause instanceof Error ? cause.message : String(cause))
  }
}

export const evaluationDataSetCommand = (
  input: RunInput & { readonly json: boolean; readonly dryRun: boolean },
): Effect.Effect<
  CommandResult,
  FileSystemFailure,
  FileSystem.FileSystem | Path.Path | HostRuntime
> =>
  mapFailure(
    Effect.gen(function* () {
      const fs = yield* FileSystem.FileSystem
      const paths = yield* Path.Path
      const runtime = yield* HostRuntime
      let run: ResolvedRun
      try {
        run = yield* resolveRun(input)
      } catch (cause) {
        return usage(cause instanceof Error ? cause.message : String(cause))
      }
      const raw = yield* runtime.readStdin
      let values: unknown
      try {
        values = JSON.parse(raw)
      } catch (cause) {
        return usage(
          `invalid JSON payload array: ${cause instanceof Error ? cause.message : String(cause)}`,
        )
      }
      if (!Array.isArray(values)) return usage("payload batch must contain one JSON array")
      if (values.length === 0) return usage("payload batch must contain at least one JSON object")
      const model = yield* runModel(run.absolute)
      const candidatesResult = yield* Effect.forEach(values, (value, index) =>
        Effect.gen(function* () {
          if (value === null || Array.isArray(value) || typeof value !== "object") {
            return yield* Effect.fail(new Error(`payload[${index}]: payload must be a JSON object`))
          }
          const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
          const token = [...(yield* runtime.randomBytes(10))]
            .map((byte) => alphabet[byte % alphabet.length])
            .join("")
          const payload = assignRecommendationId({ ...(value as Record<string, unknown>) }, token)
          return yield* Effect.try({
            try: () => {
              const kind = validateDataPayload(payload, model)
              return { index, kind, path: dataPathForPayload(kind, payload), payload }
            },
            catch: (cause) => new Error(cause instanceof Error ? cause.message : String(cause)),
          })
        }),
      ).pipe(Effect.result)
      if (Result.isFailure(candidatesResult)) return usage(candidatesResult.failure.message)
      const candidates = candidatesResult.success
      const validation = yield* Effect.try({
        try: () => {
          if (new Set(candidates.map((candidate) => candidate.path)).size !== candidates.length)
            throw new Error("payload batch contains duplicate artifact paths")
        },
        catch: (cause) => new Error(cause instanceof Error ? cause.message : String(cause)),
      }).pipe(Effect.result)
      if (Result.isFailure(validation)) return usage(validation.failure.message)
      const dataRoot = paths.join(run.absolute, "data")
      const existingPaths = (yield* fs.exists(dataRoot)) ? yield* recursiveJson(dataRoot) : []
      const existingResults = yield* Effect.forEach(existingPaths, (path) =>
        fs.readFileString(path).pipe(
          Effect.flatMap((raw) => Effect.try(() => JSON.parse(raw) as Record<string, unknown>)),
          Effect.result,
        ),
      )
      const existing = existingResults.flatMap((result) =>
        Result.isSuccess(result) ? [result.success] : [],
      )
      const crossValidation = yield* Effect.try({
        try: () => {
          const replaced = new Set(candidates.map((candidate) => candidate.path))
          validateCrossPayloads([
            ...existing.filter((payload) => {
              try {
                return !replaced.has(
                  dataPathForPayload(resolveDataKind(String(payload.kind)), payload),
                )
              } catch {
                return true
              }
            }),
            ...candidates.map((candidate) => candidate.payload),
          ])
        },
        catch: (cause) => new Error(cause instanceof Error ? cause.message : String(cause)),
      }).pipe(Effect.result)
      if (Result.isFailure(crossValidation)) return usage(crossValidation.failure.message)
      if (!input.dryRun) {
        yield* Effect.forEach(
          candidates,
          (candidate) =>
            Effect.gen(function* () {
              const absolute = paths.join(run.absolute, candidate.path)
              yield* fs.makeDirectory(paths.dirname(absolute), { recursive: true, mode: 0o755 })
              const temp = yield* fs.makeTempFile({
                directory: paths.dirname(absolute),
                prefix: ".data.",
              })
              yield* fs.writeFileString(temp, jsonDocument(candidate.payload), { mode: 0o644 })
              yield* fs.rename(temp, absolute)
            }),
          { discard: true },
        )
      }
      const result = {
        schemaVersion: 3,
        count: candidates.length,
        writes: candidates.map(({ index, kind, path }) => ({ index, kind, path })),
        ...(input.dryRun ? { dryRun: true } : {}),
        nextActions: [
          {
            id: "evaluation-status",
            label: "Inspect evaluation data status",
            command: `qualitymd evaluation status ${run.display}`,
          },
        ],
      }
      if (input.json) return commandResult(jsonDocument(result))
      return commandResult("", {
        stderr:
          `${input.dryRun ? "Would write" : "Wrote"} ${candidates.length} payloads\n` +
          candidates.map((candidate) => `- ${candidate.path}\n`).join(""),
      })
    }),
  )

export const evaluationDataListCommand = (
  input: RunInput & { readonly kind: string; readonly json: boolean },
): Effect.Effect<CommandResult, FileSystemFailure, FileSystem.FileSystem | Path.Path> =>
  mapFailure(
    Effect.gen(function* () {
      const fs = yield* FileSystem.FileSystem
      const paths = yield* Path.Path
      let run: ResolvedRun
      let filter: DataKind | undefined
      try {
        run = yield* resolveRun(input)
        filter = input.kind === "" ? undefined : resolveDataKind(input.kind)
      } catch (cause) {
        return usage(cause instanceof Error ? cause.message : String(cause))
      }
      const dataRoot = paths.join(run.absolute, "data")
      const dataPaths = (yield* fs.exists(dataRoot)) ? yield* recursiveJson(dataRoot) : []
      const inspected = yield* Effect.forEach(dataPaths, (path) =>
        fs.readFileString(path).pipe(
          Effect.flatMap((raw) =>
            Effect.try(() => {
              const payload = JSON.parse(raw) as Record<string, unknown>
              const kind = resolveDataKind(String(payload.kind))
              return filter === undefined || filter === kind
                ? {
                    kind,
                    path: paths.relative(run.absolute, path).replaceAll("\\", "/"),
                  }
                : undefined
            }),
          ),
          Effect.result,
        ),
      )
      const artifacts = inspected
        .flatMap((result) =>
          Result.isSuccess(result) && result.success !== undefined ? [result.success] : [],
        )
        .sort((left, right) => left.path.localeCompare(right.path))
      const result = { schemaVersion: 3, path: run.display, artifacts }
      if (input.json) return commandResult(jsonDocument(result))
      return commandResult(artifacts.map((item) => `${item.kind}\t${item.path}\n`).join(""))
    }),
  )

export const evaluationDataGetCommand = (
  input: RunInput & Omit<DataQuery, "kind"> & { readonly kind: string; readonly json: boolean },
): Effect.Effect<CommandResult, FileSystemFailure, FileSystem.FileSystem | Path.Path> =>
  mapFailure(
    Effect.gen(function* () {
      const fs = yield* FileSystem.FileSystem
      const paths = yield* Path.Path
      if (input.json)
        return usage("evaluation data get already emits JSON on stdout; rerun without --json")
      let run: ResolvedRun
      let relative: string
      try {
        run = yield* resolveRun(input)
        relative = dataPathForQuery({ ...input, kind: resolveDataKind(input.kind) })
      } catch (cause) {
        return usage(cause instanceof Error ? cause.message : String(cause))
      }
      return commandResult(yield* fs.readFileString(paths.join(run.absolute, relative)))
    }),
  )

export const evaluationDataVerifyCommand = (
  input: RunInput & { readonly json: boolean },
): Effect.Effect<CommandResult, FileSystemFailure, FileSystem.FileSystem | Path.Path> =>
  mapFailure(
    Effect.gen(function* () {
      const fs = yield* FileSystem.FileSystem
      const paths = yield* Path.Path
      let run: ResolvedRun
      try {
        run = yield* resolveRun(input)
      } catch (cause) {
        return usage(cause instanceof Error ? cause.message : String(cause))
      }
      const model = yield* runModel(run.absolute)
      const dataRoot = paths.join(run.absolute, "data")
      const dataPaths = (yield* fs.exists(dataRoot)) ? yield* recursiveJson(dataRoot) : []
      const inspected = yield* Effect.forEach(dataPaths, (path) =>
        Effect.gen(function* () {
          const relative = paths.relative(run.absolute, path).replaceAll("\\", "/")
          const raw = yield* fs.readFileString(path).pipe(Effect.result)
          if (Result.isFailure(raw))
            return { checked: false, failure: { path: relative, reason: String(raw.failure) } }
          const parsed = yield* Effect.try({
            try: () => JSON.parse(raw.success) as Record<string, unknown>,
            catch: (cause) => (cause instanceof Error ? cause.message : String(cause)),
          }).pipe(Effect.result)
          if (Result.isFailure(parsed))
            return { checked: false, failure: { path: relative, reason: parsed.failure } }
          const validated = yield* Effect.try({
            try: () => {
              validateDataPayload(parsed.success, model, true)
              return parsed.success
            },
            catch: (cause) => (cause instanceof Error ? cause.message : String(cause)),
          }).pipe(Effect.result)
          return Result.isSuccess(validated)
            ? { checked: true, payload: validated.success }
            : { checked: true, failure: { path: relative, reason: validated.failure } }
        }),
      )
      const validPayloads = inspected.flatMap((entry) =>
        "payload" in entry ? [entry.payload] : [],
      )
      const fileFailures = inspected.flatMap((entry) => ("failure" in entry ? [entry.failure] : []))
      const crossFailure =
        fileFailures.length === 0
          ? yield* Effect.try({
              try: () => {
                validateCrossPayloads(validPayloads)
                return undefined
              },
              catch: (cause) => ({
                path: "data",
                reason: cause instanceof Error ? cause.message : String(cause),
              }),
            }).pipe(Effect.match({ onFailure: (failure) => failure, onSuccess: () => undefined }))
          : undefined
      const failures = [...fileFailures, ...(crossFailure === undefined ? [] : [crossFailure])]
      const checked = inspected.filter((entry) => entry.checked).length
      const result = {
        schemaVersion: 3,
        path: run.display,
        valid: failures.length === 0,
        checked,
        ...(failures.length === 0 ? {} : { failures }),
      }
      if (input.json)
        return commandResult(jsonDocument(result), {
          exitCode: failures.length === 0 ? ExitCode.ok : ExitCode.problems,
        })
      return commandResult("", {
        stderr:
          failures.length === 0
            ? `Verified ${checked} payloads\n`
            : failures.map((failure) => `${failure.path}: ${failure.reason}\n`).join(""),
        exitCode: failures.length === 0 ? ExitCode.ok : ExitCode.problems,
      })
    }),
  )
