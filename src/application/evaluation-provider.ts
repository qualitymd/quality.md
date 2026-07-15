import * as Clock from "effect/Clock"
import * as Effect from "effect/Effect"
import * as FileSystem from "effect/FileSystem"
import * as Path from "effect/Path"

import { FileSystemFailure } from "../domain/errors.ts"
import { type CommandResult, ExitCode } from "../domain/command-result.ts"
import { hashJson, jsonDocument, sha256 } from "../domain/json.ts"
import type { AreaContext, EvaluationRequest, SourceFile } from "../domain/evaluator/types.ts"
import type { EvaluatorService } from "../services/evaluator.ts"
import { HostRuntime } from "../services/host-runtime.ts"
import { resolveWorkspace, type Workspace } from "../services/workspace.ts"
import { executeHarnessRun } from "./evaluation-execute.ts"
import type { EvaluationRunInput } from "./evaluation-run.ts"
import { resumeHarnessRun } from "./evaluation-resume.ts"

type JsonObject = Record<string, unknown>

interface ReceiptRequest {
  readonly requestId: string
  readonly inputHash: string
  readonly workUnitId: string
  readonly kind: string
  readonly subject?: string
  readonly instructions: string
  readonly sharedContext?: JsonObject
  readonly context?: JsonObject
  readonly source?: ReadonlyArray<SourceFile>
  readonly expectedSchema: JsonObject
}

interface Receipt {
  readonly path: string
  readonly status: string
  readonly evaluatorRequests?: ReadonlyArray<ReceiptRequest>
}

const parseReceipt = (result: CommandResult): Receipt => {
  if (result.exitCode !== ExitCode.ok || result.stdout === "")
    throw new Error(result.stderr.trim() || "evaluation run failed")
  return JSON.parse(result.stdout) as Receipt
}

const ratingCriteria = (context: JsonObject | undefined) => {
  const frame = context?.requirementEvaluationFrame
  if (frame === null || typeof frame !== "object") return {}
  const derived = (frame as JsonObject).derivedContext
  if (derived === null || typeof derived !== "object") return {}
  const criteria = (derived as JsonObject).appliedRatingCriteria
  if (!Array.isArray(criteria)) return {}
  return Object.fromEntries(
    criteria.flatMap((value) => {
      if (value === null || typeof value !== "object") return []
      const item = value as JsonObject
      return typeof item.ratingLevelId === "string" && typeof item.criterion === "string"
        ? [[item.ratingLevelId, item.criterion] as const]
        : []
    }),
  )
}

const areaFrame = (request: ReceiptRequest) => {
  const frame = request.sharedContext?.areaEvaluationFrame
  return frame !== null && typeof frame === "object" ? (frame as JsonObject) : {}
}

const areaId = (request: ReceiptRequest, frame: JsonObject) => {
  const subject = frame.subject
  if (subject !== null && typeof subject === "object") {
    const value = (subject as JsonObject).areaId
    if (typeof value === "string") return value
  }
  return request.subject?.startsWith("area:") === true ? request.subject : "area:root"
}

const areaContextFor = async (
  request: ReceiptRequest,
  bodyGuidance: string,
): Promise<AreaContext> => {
  const files = request.source ?? []
  const sourceBundleHash = await sha256(
    files.map((file) => `${file.path}\0${file.sha256}\0`).join(""),
  )
  const frame = areaFrame(request)
  const input = {
    areaId: areaId(request, frame),
    sourceBundleHash,
    frame,
    ratingCriteria: ratingCriteria(request.context),
    bodyGuidance,
    files,
  }
  return { ...input, hash: await hashJson(input) }
}

const serviceProviderCheckpoints = (
  input: EvaluationRunInput,
  evaluator: EvaluatorService,
  runPath: string,
  workspace: Workspace,
  bodyGuidance: string,
): Effect.Effect<CommandResult, Error, FileSystem.FileSystem | Path.Path | HostRuntime> =>
  Effect.scoped(
    Effect.gen(function* () {
      const fs = yield* FileSystem.FileSystem
      const paths = yield* Path.Path
      const runAbs = paths.resolve(workspace.workspaceRoot.abs, runPath)
      const artifact = JSON.parse(
        yield* fs.readFileString(paths.join(runAbs, "evaluation.json")),
      ) as {
        readonly manifest: JsonObject
      }
      let checkpoint = yield* resumeHarnessRun({
        ...input,
        evaluator: evaluator.name,
        resume: runAbs,
        resumeDisplay: runPath,
        evaluatorResult: "",
        json: true,
      })
      for (let cycle = 0; cycle < 100; cycle += 1) {
        const receipt = parseReceipt(checkpoint)
        if (receipt.status !== "awaiting_evaluator") return checkpoint
        const requests = receipt.evaluatorRequests ?? []
        if (requests.length === 0)
          throw new Error("evaluation checkpoint has no evaluator requests")
        if (
          requests.some((request) => request.kind === "resolveSource") &&
          !evaluator.capabilities.sourceResolution
        )
          throw new Error(
            `evaluator ${JSON.stringify(evaluator.name)} cannot resolve inferred source selectors; choose codex, claude, or harness, or replace the prose selector with a path or glob`,
          )
        const envelopes = yield* Effect.forEach(
          requests,
          (request) =>
            Effect.gen(function* () {
              const started = yield* Clock.currentTimeMillis
              const areaContext = yield* Effect.promise(() => areaContextFor(request, bodyGuidance))
              const evaluationRequest: EvaluationRequest = {
                runId: String(artifact.manifest.evaluationId),
                workUnitId: request.workUnitId,
                kind: request.kind,
                subject: request.subject ?? "",
                instructions: request.instructions,
                areaContext,
                context: request.context ?? {},
                expectedSchema: request.expectedSchema,
                workspaceRoot: workspace.workspaceRoot.abs,
                timeoutMs: 10 * 60 * 1000,
              }
              return yield* evaluator.evaluate(evaluationRequest).pipe(
                Effect.flatMap((response) =>
                  Effect.gen(function* () {
                    const ended = yield* Clock.currentTimeMillis
                    return {
                      requestId: request.requestId,
                      inputHash: request.inputHash,
                      evaluator: {
                        runtime: evaluator.kind,
                        ...(response.model === undefined ? {} : { model: response.model }),
                      },
                      ...(response.usage === undefined ? {} : { usage: response.usage }),
                      ...(response.contextMeta === undefined
                        ? {}
                        : { contextMeta: response.contextMeta }),
                      durationMs: Math.max(0, ended - started),
                      payload: response.payload,
                    }
                  }),
                ),
                Effect.catch((failure) =>
                  Effect.gen(function* () {
                    const ended = yield* Clock.currentTimeMillis
                    return {
                      requestId: request.requestId,
                      inputHash: request.inputHash,
                      evaluator: { runtime: evaluator.kind },
                      durationMs: Math.max(0, ended - started),
                      failure: failure.category,
                      detail: failure.detail,
                    }
                  }),
                ),
              )
            }),
          { concurrency: artifact.manifest.concurrency as number },
        )
        const resultFile = yield* fs.makeTempFileScoped({ prefix: "qualitymd-evaluator-results-" })
        yield* fs.writeFileString(resultFile, jsonDocument(envelopes), { mode: 0o600 })
        checkpoint = yield* resumeHarnessRun({
          ...input,
          evaluator: evaluator.name,
          resume: runAbs,
          resumeDisplay: runPath,
          evaluatorResult: resultFile,
          json: true,
        })
      }
      throw new Error("evaluation exceeded 100 evaluator checkpoint cycles")
    }),
  )

export const executeProviderRun = (
  input: EvaluationRunInput,
  evaluator: EvaluatorService,
): Effect.Effect<
  CommandResult,
  FileSystemFailure,
  FileSystem.FileSystem | Path.Path | HostRuntime
> =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const paths = yield* Path.Path
    const workspace = yield* resolveWorkspace({
      ...(input.model === "" ? {} : { model: input.model }),
      ...(input.evaluationDir === "" ? {} : { evaluationDir: input.evaluationDir }),
    })
    const runtime = yield* HostRuntime
    const configuredConcurrency =
      workspace.evaluation.concurrency ?? runtime.hardwareConcurrency * 2
    const resolvedConcurrency =
      configuredConcurrency > 1 &&
      !evaluator.capabilities.concurrent &&
      !evaluator.capabilities.subagents
        ? 1
        : configuredConcurrency
    const document = yield* fs.readFileString(workspace.model.abs)
    const bodyStart = document.indexOf("\n---", 4)
    const bodyGuidance = bodyStart < 0 ? "" : document.slice(bodyStart + 4)
    const created = yield* executeHarnessRun({ ...input, evaluator: "harness", json: true })
    const initial = parseReceipt(created)
    const runAbs = paths.resolve(workspace.workspaceRoot.abs, initial.path)
    const artifactPath = paths.join(runAbs, "evaluation.json")
    const artifact = JSON.parse(yield* fs.readFileString(artifactPath)) as {
      readonly manifest: JsonObject
      readonly state: JsonObject
    }
    artifact.manifest.evaluator = evaluator.name
    artifact.manifest.evaluatorKind = evaluator.kind
    artifact.manifest.evaluatorCapabilities = evaluator.capabilities
    artifact.manifest.concurrency = resolvedConcurrency
    delete artifact.state.harnessIdentity
    yield* fs.writeFileString(artifactPath, jsonDocument(artifact), { mode: 0o644 })

    return yield* serviceProviderCheckpoints(
      input,
      evaluator,
      initial.path,
      workspace,
      bodyGuidance,
    )
  }).pipe(
    Effect.mapError(
      (cause) =>
        new FileSystemFailure({ detail: cause instanceof Error ? cause.message : String(cause) }),
    ),
  )

export const resumeProviderRun = (
  input: EvaluationRunInput,
  evaluator: EvaluatorService,
  workspace: Workspace,
): Effect.Effect<
  CommandResult,
  FileSystemFailure,
  FileSystem.FileSystem | Path.Path | HostRuntime
> =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const document = yield* fs.readFileString(workspace.model.abs)
    const bodyStart = document.indexOf("\n---", 4)
    const bodyGuidance = bodyStart < 0 ? "" : document.slice(bodyStart + 4)
    return yield* serviceProviderCheckpoints(
      input,
      evaluator,
      input.resume,
      workspace,
      bodyGuidance,
    )
  }).pipe(
    Effect.mapError(
      (cause) =>
        new FileSystemFailure({ detail: cause instanceof Error ? cause.message : String(cause) }),
    ),
  )
