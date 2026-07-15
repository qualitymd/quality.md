import * as Clock from "effect/Clock"
import * as Effect from "effect/Effect"
import * as FileSystem from "effect/FileSystem"
import * as Path from "effect/Path"
import * as Queue from "effect/Queue"

import { type CommandResult, ExitCode } from "../domain/command-result.ts"
import { FileSystemFailure } from "../domain/errors.ts"
import type { EvaluationRequest, InspectionContext } from "../domain/evaluator/types.ts"
import { jsonDocument } from "../domain/json.ts"
import type { EvaluatorService } from "../services/evaluator.ts"
import { HostRuntime } from "../services/host-runtime.ts"
import { resolveWorkspace, type Workspace } from "../services/workspace.ts"
import { executeEvaluationRun } from "./evaluation-execute.ts"
import type { EvaluationRunInput } from "./evaluation-run.ts"
import { resumeEvaluationRun } from "./evaluation-resume.ts"

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
  readonly concurrency: number
  readonly evaluatorRequests?: ReadonlyArray<ReceiptRequest>
}

interface ResultEnvelope extends JsonObject {
  readonly requestId: string
  readonly inputHash: string
}

interface Completion {
  readonly requestId: string
  readonly envelope: ResultEnvelope
}

const dispatching = (receipt: Receipt) =>
  receipt.status === "running" || receipt.status === "awaiting_evaluator"

const recordCancellation = Effect.fn("qualitymd.recordEvaluationCancellation")(function* (
  runAbs: string,
) {
  const fs = yield* FileSystem.FileSystem
  const paths = yield* Path.Path
  const artifactPath = paths.join(runAbs, "evaluation.json")
  if (!(yield* fs.exists(artifactPath))) return
  const artifact = JSON.parse(yield* fs.readFileString(artifactPath)) as {
    state: {
      status: string
      updatedAt: string
      cancelled?: boolean
    }
  }
  if (artifact.state.status === "completed" || artifact.state.status === "failed") return
  const timestamp = new Date(yield* Clock.currentTimeMillis).toISOString().replace(/\.\d{3}Z$/, "Z")
  artifact.state.status = "cancelled"
  artifact.state.cancelled = true
  artifact.state.updatedAt = timestamp
  const temp = yield* fs.makeTempFile({ directory: runAbs, prefix: ".evaluation." })
  yield* fs.writeFileString(temp, jsonDocument(artifact), { mode: 0o644 })
  yield* fs.rename(temp, artifactPath)
  yield* fs.writeFileString(
    paths.join(runAbs, "logs/events.jsonl"),
    `${JSON.stringify({ timestamp, event: "run_status", status: "cancelled" })}\n`,
    { flag: "a", mode: 0o600 },
  )
})

const parseReceipt = (result: CommandResult): Receipt => {
  if (result.exitCode !== ExitCode.ok || result.stdout === "")
    throw new Error(result.stderr.trim() || "evaluation run failed")
  return JSON.parse(result.stdout) as Receipt
}

const evaluateRequest = Effect.fn("qualitymd.evaluateProviderRequest")(function* (
  request: ReceiptRequest,
  evaluator: EvaluatorService,
  workspace: Workspace,
  runId: string,
) {
  const started = yield* Clock.currentTimeMillis
  const inspection =
    request.inspection === undefined
      ? undefined
      : ({
          ...request.inspection,
          workspaceRoot: workspace.workspaceRoot.abs,
        } as unknown as InspectionContext)
  const evaluationRequest: EvaluationRequest = {
    runId,
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
          ...(response.contextMeta === undefined ? {} : { contextMeta: response.contextMeta }),
          durationMs: Math.max(0, ended - started),
          payload: response.payload,
        } satisfies ResultEnvelope
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
        } satisfies ResultEnvelope
      }),
    ),
  )
})

const driveProviderRun = (
  input: EvaluationRunInput,
  evaluator: EvaluatorService,
  workspace: Workspace,
  initial: CommandResult,
): Effect.Effect<CommandResult, Error, FileSystem.FileSystem | Path.Path | HostRuntime> =>
  Effect.scoped(
    Effect.gen(function* () {
      const fs = yield* FileSystem.FileSystem
      const paths = yield* Path.Path
      const queue = yield* Queue.unbounded<Completion>()
      let receipt = parseReceipt(initial)
      const runAbs = paths.resolve(workspace.workspaceRoot.abs, receipt.path)
      const artifact = JSON.parse(
        yield* fs.readFileString(paths.join(runAbs, "evaluation.json")),
      ) as { readonly manifest: { readonly evaluationId: string } }
      const runId = artifact.manifest.evaluationId
      const active = new Set<string>()
      let peakActive = 0

      return yield* Effect.gen(function* () {
        for (let cycle = 0; cycle < 10_000; cycle += 1) {
          if (!dispatching(receipt)) return initial
          const requests = receipt.evaluatorRequests ?? []
          for (const request of requests) {
            if (active.has(request.requestId)) continue
            if (active.size >= receipt.concurrency) break
            active.add(request.requestId)
            yield* evaluateRequest(request, evaluator, workspace, runId).pipe(
              Effect.flatMap((envelope) =>
                Queue.offer(queue, { requestId: request.requestId, envelope }),
              ),
              Effect.forkScoped,
            )
          }
          if (active.size === 0)
            throw new Error("evaluation checkpoint has no provider requests to dispatch")
          if (active.size > peakActive) {
            peakActive = active.size
            const timestamp = new Date(yield* Clock.currentTimeMillis)
              .toISOString()
              .replace(/\.\d{3}Z$/, "Z")
            yield* fs.writeFileString(
              paths.join(runAbs, "logs/events.jsonl"),
              `${JSON.stringify({
                timestamp,
                event: "dispatch_activity",
                dispatchMode: evaluator.capabilities.dispatch.concurrentCalls
                  ? "direct"
                  : "sequential",
                active: active.size,
                peakActive,
              })}\n`,
              { flag: "a", mode: 0o600 },
            )
          }
          const completion = yield* Queue.take(queue)
          const resultFile = yield* fs.makeTempFileScoped({
            prefix: "qualitymd-evaluator-result-",
          })
          yield* fs.writeFileString(resultFile, jsonDocument(completion.envelope), { mode: 0o600 })
          const checkpoint = yield* resumeEvaluationRun({
            ...input,
            evaluator: evaluator.name,
            resume: runAbs,
            resumeDisplay: receipt.path,
            evaluatorResult: resultFile,
            json: true,
          })
          receipt = parseReceipt(checkpoint)
          active.delete(completion.requestId)
          if (!dispatching(receipt)) return checkpoint
        }
        throw new Error("evaluation exceeded 10000 provider completion cycles")
      }).pipe(Effect.onInterrupt(() => recordCancellation(runAbs)))
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
    const workspace = yield* resolveWorkspace({
      ...(input.model === "" ? {} : { model: input.model }),
      ...(input.evaluationDir === "" ? {} : { evaluationDir: input.evaluationDir }),
    })
    const initial = yield* executeEvaluationRun({ ...input, json: true }, evaluator)
    return yield* driveProviderRun(input, evaluator, workspace, initial)
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
    const checkpoint = yield* resumeEvaluationRun({ ...input, evaluatorResult: "", json: true })
    return yield* driveProviderRun(input, evaluator, workspace, checkpoint)
  }).pipe(
    Effect.mapError(
      (cause) =>
        new FileSystemFailure({ detail: cause instanceof Error ? cause.message : String(cause) }),
    ),
  )
