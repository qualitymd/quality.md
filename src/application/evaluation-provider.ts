import * as Clock from "effect/Clock"
import * as Effect from "effect/Effect"
import * as FileSystem from "effect/FileSystem"
import * as Path from "effect/Path"

import { FileSystemFailure } from "../domain/errors.ts"
import { type CommandResult, ExitCode } from "../domain/command-result.ts"
import { jsonDocument } from "../domain/json.ts"
import type { EvaluationRequest, InspectionContext } from "../domain/evaluator/types.ts"
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
  readonly bodyGuidance?: string
  readonly inspection?: JsonObject
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

const serviceProviderCheckpoints = (
  input: EvaluationRunInput,
  evaluator: EvaluatorService,
  runPath: string,
  workspace: Workspace,
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
        const envelopes = yield* Effect.forEach(
          requests,
          (request) =>
            Effect.gen(function* () {
              const started = yield* Clock.currentTimeMillis
              const inspection =
                request.inspection === undefined
                  ? undefined
                  : ({
                      ...request.inspection,
                      workspaceRoot: workspace.workspaceRoot.abs,
                    } as unknown as InspectionContext)
              const evaluationRequest: EvaluationRequest = {
                runId: String(artifact.manifest.evaluationId),
                workUnitId: request.workUnitId,
                kind: request.kind,
                subject: request.subject ?? "",
                instructions: request.instructions,
                sharedContext: request.sharedContext ?? {},
                context: request.context ?? {},
                bodyGuidance: request.bodyGuidance ?? "",
                ...(inspection === undefined ? {} : { inspection }),
                expectedSchema: request.expectedSchema,
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

    return yield* serviceProviderCheckpoints(input, evaluator, initial.path, workspace)
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
    return yield* serviceProviderCheckpoints(input, evaluator, input.resume, workspace)
  }).pipe(
    Effect.mapError(
      (cause) =>
        new FileSystemFailure({ detail: cause instanceof Error ? cause.message : String(cause) }),
    ),
  )
