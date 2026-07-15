import * as Effect from "effect/Effect"
import * as FileSystem from "effect/FileSystem"
import * as Path from "effect/Path"

import { FileSystemFailure } from "../domain/errors.ts"
import { commandResult, type CommandResult } from "../domain/command-result.ts"
import { initialFramePayloads } from "../domain/evaluation/frames.ts"
import { resolveConcurrency } from "../domain/evaluation/concurrency.ts"
import { buildGraph, readyUnits } from "../domain/evaluation/graph.ts"
import { planEvaluation } from "../domain/evaluation/plan.ts"
import {
  buildProtocolRequest,
  completeProtocolRequest,
  protocolRequestReceipt,
} from "../domain/evaluation/protocol.ts"
import {
  evaluationRunArtifact,
  evaluationRunEvents,
  evaluationRunReceipt,
  renderHarnessAwaiting,
  newEvaluationIdentity,
  resolveScope,
  scopeSlug,
} from "../domain/evaluation/run.ts"
import { jsonDocument } from "../domain/json.ts"
import { decodeModel } from "../domain/model/model.ts"
import { parseQualityDocument } from "../domain/model/document.ts"
import { harnessCapabilities } from "../adapters/evaluator.ts"
import type { EvaluatorCapabilities, EvaluatorKind } from "../domain/evaluator/types.ts"
import { HostRuntime } from "../services/host-runtime.ts"
import { detectSourceKind, validateSourceSelector, type AreaSource } from "../services/source.ts"
import { resolveWorkspace } from "../services/workspace.ts"
import { atomicWriteFileString } from "../services/atomic-file.ts"
import { hashJsonEffect, requestId } from "./evaluation-hash.ts"
import { nextEvaluationRunNumber } from "./evaluation-runs.ts"
import type { EvaluationRunInput } from "./evaluation-run.ts"

export interface RunEvaluator {
  readonly name: string
  readonly kind: EvaluatorKind
  readonly capabilities: EvaluatorCapabilities
}

const execute = Effect.fn("qualitymd.executeEvaluationRun")(function* (
  input: EvaluationRunInput,
  evaluator: RunEvaluator,
) {
  const fs = yield* FileSystem.FileSystem
  const paths = yield* Path.Path
  const runtime = yield* HostRuntime
  const workspace = yield* resolveWorkspace({
    ...(input.model === "" ? {} : { model: input.model }),
    ...(input.evaluationDir === "" ? {} : { evaluationDir: input.evaluationDir }),
  })
  const raw = yield* fs.readFileString(workspace.model.abs)
  const document = parseQualityDocument(workspace.model.abs, raw)
  const bodyGuidance = document.body.trim()
  const model = decodeModel(document)
  const scope = resolveScope(model, input.area, input.factors)
  const plan = planEvaluation(model, scope.plannedScope)
  const concurrency = resolveConcurrency(
    workspace.evaluation.concurrency,
    evaluator.capabilities.dispatch,
  )
  const number = yield* nextEvaluationRunNumber(workspace.evaluations.abs)
  const label = `${String(number).padStart(4, "0")}-${scopeSlug(scope.plannedScope)}-eval`
  const runAbs = paths.join(workspace.evaluations.abs, label)
  const runRel = `${workspace.evaluations.rel}/${label}`
  const identity = newEvaluationIdentity(
    new Date(yield* runtime.currentTimeMillis),
    yield* runtime.randomBytes(12),
  )
  const timestamp = identity.createdAt
  const completedFrames = yield* Effect.forEach(
    initialFramePayloads(model, plan, workspace.model.rel),
    (entry) =>
      Effect.gen(function* () {
        return { ...entry, inputHash: yield* hashJsonEffect(entry.payload) }
      }),
  )
  const payloads = completedFrames.map(({ workUnit, payload }) => ({ workUnit, payload }))
  const completedWorkUnits = Object.fromEntries(
    completedFrames.map(({ workUnit, inputHash }) => [
      workUnit,
      { status: "completed", inputHash, completedAt: timestamp },
    ]),
  )
  const sourceEntries = yield* Effect.forEach(plan.areas, (area) =>
    Effect.gen(function* () {
      const kind = yield* detectSourceKind(workspace.workspaceRoot.abs, area.source)
      const source = { selector: area.source, kind } satisfies AreaSource
      yield* validateSourceSelector(workspace.workspaceRoot.abs, source)
      return { area: area.ref, source }
    }),
  )
  const areaSources = Object.fromEntries(sourceEntries.map(({ area, source }) => [area, source]))
  const sourcePlans = sourceEntries.map(({ area, source }) => ({ area, ...source }))
  const graph = buildGraph(plan)
  const initialRequests = yield* Effect.forEach(
    readyUnits(graph, new Set(Object.keys(completedWorkUnits)), concurrency.resolved),
    (unit) =>
      Effect.gen(function* () {
        const draft = buildProtocolRequest({
          unit,
          plan,
          payloads,
          areaSources,
          bodyGuidance,
          evaluationId: identity.evaluationId,
        })
        const protocol = completeProtocolRequest(draft, yield* hashJsonEffect(draft.hashInput))
        const pending = {
          requestId: yield* requestId(identity.evaluationId, unit.id, protocol.inputHash, 1),
          workUnitId: unit.id,
          inputHash: protocol.inputHash,
          correlationId: protocol.correlationId,
          attempt: 1,
        }
        return { pending, receipt: protocolRequestReceipt(protocol, pending) }
      }),
  )
  const pending = initialRequests.map((entry) => entry.pending)
  const requests = initialRequests.map((entry) => entry.receipt)
  const workUnits = {
    ...completedWorkUnits,
    ...Object.fromEntries(pending.map((entry) => [entry.workUnitId, { status: "pending" }])),
  }
  const artifact = evaluationRunArtifact({
    identity,
    model: workspace.model.rel,
    scope,
    number,
    label,
    evaluator,
    concurrency,
    areaSources,
    workUnits,
    pending,
    payloads,
  })
  yield* fs.makeDirectory(workspace.evaluations.abs, { recursive: true, mode: 0o755 })
  yield* fs.makeDirectory(runAbs, { mode: 0o755 })
  yield* fs.makeDirectory(paths.join(runAbs, "logs"), { mode: 0o755 })
  yield* fs.writeFileString(paths.join(runAbs, "model-snapshot.md"), raw, { mode: 0o644 })
  yield* atomicWriteFileString(paths.join(runAbs, "evaluation.json"), jsonDocument(artifact), {
    mode: 0o644,
  })
  yield* fs.writeFileString(paths.join(runAbs, "logs/evaluator-calls.jsonl"), "", {
    mode: 0o600,
  })
  yield* fs.writeFileString(
    paths.join(runAbs, "logs/events.jsonl"),
    evaluationRunEvents(timestamp, identity.evaluationId, evaluator, concurrency, pending.length),
    { mode: 0o600 },
  )
  const result = evaluationRunReceipt({
    path: runRel,
    evaluator: evaluator.name,
    evaluatorKind: evaluator.kind,
    concurrency: concurrency.resolved,
    total: graph.length,
    evaluatorUnits: graph.filter((unit) => unit.evaluatorBacked).length,
    completed: Object.keys(workUnits).length,
    sources: sourcePlans,
    requests,
    dispatchMode: concurrency.mode,
  })
  if (input.json) return commandResult(jsonDocument(result))
  return commandResult("", {
    stderr: renderHarnessAwaiting(requests, concurrency.resolved, result.nextActions[0]!.command),
  })
})

export const executeEvaluationRun = (
  input: EvaluationRunInput,
  evaluator: RunEvaluator,
): Effect.Effect<
  CommandResult,
  FileSystemFailure,
  FileSystem.FileSystem | Path.Path | HostRuntime
> =>
  execute(input, evaluator).pipe(
    Effect.mapError(
      (cause) =>
        new FileSystemFailure({
          detail: cause instanceof Error ? cause.message : String(cause),
        }),
    ),
  )

export const executeHarnessRun = (input: EvaluationRunInput) =>
  executeEvaluationRun(input, {
    name: "harness",
    kind: "harness",
    capabilities: harnessCapabilities,
  })
