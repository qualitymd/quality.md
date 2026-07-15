import * as Effect from "effect/Effect"
import * as FileSystem from "effect/FileSystem"
import * as Path from "effect/Path"

import { FileSystemFailure } from "../domain/errors.ts"
import { commandResult, ExitCode, type CommandResult } from "../domain/command-result.ts"
import { buildGraph, type WorkUnit } from "../domain/evaluation/graph.ts"
import { planEvaluation } from "../domain/evaluation/plan.ts"
import {
  buildProtocolRequest,
  deterministicPayload,
  type ProtocolRequest,
  type StoredPayload,
} from "../domain/evaluation/protocol.ts"
import { normalizeEvaluatorResult, ResultValidationError } from "../domain/evaluation/result.ts"
import { hashJson, jsonDocument } from "../domain/json.ts"
import { decodeModel } from "../domain/model/model.ts"
import { parseQualityDocument } from "../domain/model/document.ts"
import { captureSource, packageSource, type SourceBundle } from "../services/source.ts"
import { HostRuntime } from "../services/host-runtime.ts"
import type { EvaluationRunInput } from "./evaluation-run.ts"
import { buildReportsAtRun } from "./evaluation-report.ts"

type JsonObject = Record<string, unknown>

interface UnitState extends JsonObject {
  status: string
  attempts?: number
  inputHash?: string
  completedAt?: string
  failure?: { readonly category: string; readonly detail?: string }
}

interface PendingCall {
  requestId: string
  workUnitId: string
  inputHash: string
  correlationId: string
  attempt: number
}

interface Artifact {
  schemaVersion: number
  kind: string
  manifest: {
    evaluationId: string
    createdAt: string
    model: string
    requestedScope: JsonObject
    plannedScope: { areaId: string; factorFilter: ReadonlyArray<string> }
    run: { number: number; label: string }
    evaluator: string
    evaluatorKind: string
    evaluatorCapabilities: JsonObject
    concurrency: number
  }
  state: {
    status: string
    workUnits: Record<string, UnitState>
    startedAt: string
    updatedAt: string
    completedAt?: string
    pendingEvaluatorCalls?: Array<PendingCall>
    harnessIdentity?: { runtime: string }
    failure?: { category: string; detail?: string }
  }
  sources: Record<string, JsonObject>
  results: { payloads: Array<StoredPayload> }
  outputs?: JsonObject
}

interface Envelope {
  readonly requestId: string
  readonly inputHash: string
  readonly evaluator: { readonly runtime: string; readonly model?: string }
  readonly usage?: JsonObject
  readonly durationMs?: number
  readonly contextMeta?: JsonObject
  readonly payload?: unknown
  readonly failure?: string
  readonly detail?: string
}

const retryable = new Set([
  "rate_limited",
  "timeout",
  "invalid_evaluator_output",
  "schema_invalid_output",
])

const failureResult = (detail: string, exitCode: ExitCode = ExitCode.internal) =>
  commandResult("", { stderr: `qualitymd: ${detail}\n`, exitCode })

const readEnvelopes = (path: string) =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const runtime = yield* HostRuntime
    const raw = path === "-" ? yield* runtime.readStdin : yield* fs.readFileString(path)
    return yield* Effect.try({
      try: () => {
        const decoded = JSON.parse(raw) as unknown
        const values = Array.isArray(decoded) ? decoded : [decoded]
        if (values.length === 0)
          throw new Error("--evaluator-result must carry at least one result envelope")
        return values.map((value): Envelope => {
          if (value === null || Array.isArray(value) || typeof value !== "object")
            throw new Error("--evaluator-result must not carry a null envelope")
          const envelope = value as Partial<Envelope>
          if (typeof envelope.requestId !== "string" || envelope.requestId === "")
            throw new Error("--evaluator-result envelopes must carry the outstanding requestId")
          if (typeof envelope.inputHash !== "string" || envelope.inputHash === "")
            throw new Error("--evaluator-result envelopes must carry the outstanding inputHash")
          if (
            envelope.evaluator === undefined ||
            typeof envelope.evaluator.runtime !== "string" ||
            envelope.evaluator.runtime === ""
          )
            throw new Error(
              "--evaluator-result envelopes must identify the harness runtime in evaluator.runtime",
            )
          if (
            (envelope.failure === undefined || envelope.failure === "") &&
            envelope.payload === undefined
          )
            throw new Error(
              "--evaluator-result envelopes must carry a payload or a classified failure",
            )
          return envelope as Envelope
        })
      },
      catch: (cause) => (cause instanceof Error ? cause : new Error(String(cause))),
    })
  })

const bundleFromRecord = (record: JsonObject): SourceBundle | undefined => {
  if (record.kind !== "prose" || typeof record.bundleHash !== "string") return undefined
  if (!Array.isArray(record.files)) return undefined
  return {
    hash: record.bundleHash,
    truncated: record.truncated === true,
    files: record.files.flatMap((value) => {
      if (value === null || typeof value !== "object") return []
      const file = value as JsonObject
      if (
        typeof file.path !== "string" ||
        typeof file.sha256 !== "string" ||
        typeof file.content !== "string"
      )
        return []
      return [
        {
          path: file.path,
          sha256: file.sha256,
          content: file.content,
          ...(file.truncated === true ? { truncated: true } : {}),
        },
      ]
    }),
  }
}

const requestReceipt = (
  protocol: ProtocolRequest,
  pending: PendingCall,
  lastFailure?: { readonly category: string; readonly detail?: string },
) => ({
  requestId: pending.requestId,
  workUnitId: protocol.workUnitId,
  kind: protocol.kind,
  ...(protocol.subject === "" ? {} : { subject: protocol.subject }),
  attempt: pending.attempt,
  instructions: protocol.instructions,
  ...(protocol.sharedContext === null || Object.keys(protocol.sharedContext).length === 0
    ? {}
    : { sharedContext: protocol.sharedContext }),
  ...(protocol.context === null || Object.keys(protocol.context).length === 0
    ? {}
    : { context: protocol.context }),
  ...(protocol.source.length === 0 ? {} : { source: protocol.source }),
  expectedSchema: protocol.expectedSchema,
  inputHash: protocol.inputHash,
  correlationId: protocol.correlationId,
  ...(lastFailure === undefined ? {} : { lastFailure }),
})

const randomRequestId = (
  evaluationId: string,
  workUnit: string,
  inputHash: string,
  attempt: number,
) =>
  hashJson({ evaluationId, workUnit, inputHash, attempt }).then(
    (hash) => `req_${hash.slice(0, 16)}`,
  )

const mergePayloads = (
  artifact: Artifact,
  graph: ReadonlyArray<WorkUnit>,
  workUnit: string,
  values: ReadonlyArray<StoredPayload>,
) => {
  const order = new Map(graph.map((unit, index) => [unit.id, index]))
  artifact.results.payloads = [
    ...artifact.results.payloads.filter((entry) => entry.workUnit !== workUnit),
    ...values,
  ].sort((left, right) => (order.get(left.workUnit) ?? 0) - (order.get(right.workUnit) ?? 0))
}

const appendCallLog = (fs: FileSystem.FileSystem, path: string, entry: JsonObject) =>
  fs.writeFileString(path, `${JSON.stringify(entry)}\n`, { flag: "a", mode: 0o600 })

export const resumeHarnessRun = (
  input: EvaluationRunInput,
): Effect.Effect<
  CommandResult,
  FileSystemFailure,
  FileSystem.FileSystem | Path.Path | HostRuntime
> =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const paths = yield* Path.Path
    const runtime = yield* HostRuntime
    const displayRun = input.resumeDisplay ?? input.resume
    const runAbs = paths.resolve(input.resume)
    const artifactPath = paths.join(runAbs, "evaluation.json")
    if (!(yield* fs.exists(artifactPath)))
      return failureResult(`${displayRun} is not a resumable evaluation run`, ExitCode.usage)
    const artifact = JSON.parse(yield* fs.readFileString(artifactPath)) as Artifact
    if (artifact.schemaVersion !== 7 || artifact.kind !== "EvaluationRun")
      return failureResult(
        `evaluation artifact schema ${artifact.schemaVersion} is incompatible with schema 7`,
      )
    let workspaceRoot = ""
    if (paths.isAbsolute(artifact.manifest.model)) {
      if (yield* fs.exists(artifact.manifest.model))
        workspaceRoot = paths.dirname(artifact.manifest.model)
    } else {
      let candidate = runAbs
      while (true) {
        if (yield* fs.exists(paths.join(candidate, artifact.manifest.model))) {
          workspaceRoot = candidate
          break
        }
        const parent = paths.dirname(candidate)
        if (parent === candidate) break
        candidate = parent
      }
    }
    if (workspaceRoot === "")
      return failureResult(
        `cannot locate workspace model ${artifact.manifest.model} while resuming ${displayRun}`,
        ExitCode.usage,
      )
    const snapshotPath = paths.join(runAbs, "model-snapshot.md")
    const snapshot = yield* fs.readFileString(snapshotPath)
    const model = decodeModel(parseQualityDocument(snapshotPath, snapshot))
    const plan = planEvaluation(model, artifact.manifest.plannedScope)
    const sourceKinds = Object.fromEntries(
      Object.entries(artifact.sources).map(([area, record]) => [area, record.kind]),
    ) as Record<string, "path" | "glob" | "prose">
    const graph = buildGraph(plan, sourceKinds)
    const bundles = new Map<string, SourceBundle>()
    for (const area of plan.areas) {
      const record = artifact.sources[area.ref]!
      const captured = bundleFromRecord(record)
      if (captured !== undefined) bundles.set(area.ref, captured)
      else if (record.kind === "path" || record.kind === "glob") {
        bundles.set(
          area.ref,
          yield* packageSource(workspaceRoot, String(record.selector), record.kind),
        )
      }
    }
    const now = new Date(yield* runtime.currentTimeMillis).toISOString().replace(/\.\d{3}Z$/, "Z")
    const initialStatus = artifact.state.status
    const events: Array<JsonObject> = []
    const retryFailures = new Map<string, { readonly category: string; readonly detail?: string }>()
    if (input.evaluatorResult !== "") {
      const envelopes = yield* readEnvelopes(input.evaluatorResult).pipe(
        Effect.mapError((cause) => new FileSystemFailure({ detail: cause.message })),
      )
      const remaining = [...envelopes]
      const outstanding = [...(artifact.state.pendingEvaluatorCalls ?? [])]
      for (const pending of outstanding) {
        const match = remaining.findIndex(
          (envelope) =>
            envelope.requestId === pending.requestId && envelope.inputHash === pending.inputHash,
        )
        if (match < 0) continue
        const envelope = remaining.splice(match, 1)[0]!
        if (
          artifact.state.harnessIdentity !== undefined &&
          artifact.state.harnessIdentity.runtime !== envelope.evaluator.runtime
        ) {
          return failureResult(
            `run judgment is bound to harness runtime ${JSON.stringify(artifact.state.harnessIdentity.runtime)}; a result from ${JSON.stringify(envelope.evaluator.runtime)} would mix harnesses — start a new run`,
          )
        }
        const unit = graph.find((candidate) => candidate.id === pending.workUnitId)
        if (unit === undefined)
          return failureResult(
            `the pending work request ${pending.requestId} targets unknown work unit ${pending.workUnitId}; start a new run`,
          )
        const protocol = yield* Effect.promise(() =>
          buildProtocolRequest({
            unit,
            plan,
            payloads: artifact.results.payloads,
            bundles,
            sourceRecords: artifact.sources,
            evaluationId: artifact.manifest.evaluationId,
          }),
        )
        if (protocol.inputHash !== pending.inputHash)
          return failureResult(
            `the model or source for ${unit.id} changed after the work request was emitted; a result must not be attached to different evidence — start a new run`,
          )
        const state = (artifact.state.workUnits[unit.id] ??= { status: "pending" })
        state.attempts = (state.attempts ?? 0) + 1
        let category = envelope.failure ?? ""
        let detail = envelope.detail ?? ""
        let accepted: ReadonlyArray<StoredPayload> = []
        if (category === "") {
          try {
            if (unit.kind === "resolveSource") {
              const bundle = yield* captureSource(workspaceRoot, envelope.payload)
              bundles.set(unit.subject, bundle)
              const record = artifact.sources[unit.subject]!
              record.bundleHash = bundle.hash
              record.capturedAt = now
              record.harnessRuntime = envelope.evaluator.runtime
              if (bundle.truncated) record.truncated = true
              record.files = bundle.files.map((file) => ({
                path: file.path,
                sha256: file.sha256,
                content: file.content,
                ...(file.truncated === true ? { truncated: true } : {}),
              }))
            } else {
              const bytes = yield* runtime.randomBytes(100)
              const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
              accepted = normalizeEvaluatorResult(unit, envelope.payload, (index) =>
                Array.from(
                  bytes.slice(index * 10, index * 10 + 10),
                  (byte) => alphabet[byte % alphabet.length],
                ).join(""),
              )
            }
          } catch (cause) {
            category =
              cause instanceof ResultValidationError ? cause.category : "invalid_evaluator_output"
            detail = cause instanceof Error ? cause.message : String(cause)
          }
        }
        yield* appendCallLog(fs, paths.join(runAbs, "logs/evaluator-calls.jsonl"), {
          timestamp: now,
          evaluator: artifact.manifest.evaluator,
          evaluatorKind: artifact.manifest.evaluatorKind,
          capabilities: artifact.manifest.evaluatorCapabilities,
          harnessRuntime: envelope.evaluator.runtime,
          ...(envelope.evaluator.model === undefined ? {} : { model: envelope.evaluator.model }),
          workUnit: unit.id,
          attempt: state.attempts,
          durationMs: envelope.durationMs ?? 0,
          ...(envelope.contextMeta === undefined ? {} : { contextMeta: envelope.contextMeta }),
          ...(envelope.usage === undefined ? {} : { usage: envelope.usage }),
          ...(category === "" ? {} : { failure: { category, detail } }),
        })
        if (category === "") {
          if (artifact.manifest.evaluatorKind === "harness")
            artifact.state.harnessIdentity ??= { runtime: envelope.evaluator.runtime }
          artifact.state.pendingEvaluatorCalls = (
            artifact.state.pendingEvaluatorCalls ?? []
          ).filter((call) => call.requestId !== pending.requestId)
          mergePayloads(artifact, graph, unit.id, accepted)
          Object.assign(state, {
            status: "completed",
            inputHash: protocol.inputHash,
            completedAt: now,
          })
          delete state.failure
          events.push({
            timestamp: now,
            event: "work_unit_completed",
            workUnit: unit.id,
            attempt: state.attempts,
          })
        } else if (retryable.has(category) && state.attempts < 3) {
          pending.attempt = state.attempts + 1
          pending.requestId = yield* Effect.promise(() =>
            randomRequestId(
              artifact.manifest.evaluationId,
              unit.id,
              pending.inputHash,
              pending.attempt,
            ),
          )
          retryFailures.set(pending.requestId, { category, ...(detail === "" ? {} : { detail }) })
          events.push({
            timestamp: now,
            event: "work_unit_retry",
            workUnit: unit.id,
            attempt: state.attempts,
            failure: { category, ...(detail === "" ? {} : { detail }) },
          })
        } else {
          artifact.state.pendingEvaluatorCalls = (
            artifact.state.pendingEvaluatorCalls ?? []
          ).filter((call) => call.requestId !== pending.requestId)
          state.status = "failed"
          state.failure = { category, ...(detail === "" ? {} : { detail }) }
          artifact.state.status = "failed"
          artifact.state.failure = state.failure
          artifact.state.completedAt = now
          events.push({
            timestamp: now,
            event: "work_unit_failed",
            workUnit: unit.id,
            attempt: state.attempts,
            failure: state.failure,
          })
        }
      }
      if (remaining.length > 0)
        return failureResult(
          `the submitted result ${remaining[0]!.requestId} does not correlate with an outstanding work request; resume without --evaluator-result to recover the outstanding requests`,
        )
    }
    const requests: Array<JsonObject> = []
    if (artifact.state.status !== "failed") {
      for (const unit of graph) {
        const state = artifact.state.workUnits[unit.id]
        if (state?.status === "completed") continue
        if (
          unit.dependsOn.some(
            (dependency) => artifact.state.workUnits[dependency]?.status !== "completed",
          )
        )
          continue
        if (!unit.evaluatorBacked) {
          if (unit.kind === "buildReports") {
            const built = yield* buildReportsAtRun(runAbs, displayRun, artifact)
            artifact.outputs = built.outputs
            artifact.state.workUnits[unit.id] = {
              status: "completed",
              completedAt: now,
            }
            events.push({ timestamp: now, event: "work_unit_completed", workUnit: unit.id })
            continue
          }
          const payload = deterministicPayload(unit, model, plan, artifact.manifest.model)
          const inputHash = yield* Effect.promise(() => hashJson(payload))
          mergePayloads(artifact, graph, unit.id, [{ workUnit: unit.id, payload }])
          artifact.state.workUnits[unit.id] = {
            status: "completed",
            inputHash,
            completedAt: now,
          }
          events.push({ timestamp: now, event: "work_unit_completed", workUnit: unit.id })
          continue
        }
        const protocol = yield* Effect.promise(() =>
          buildProtocolRequest({
            unit,
            plan,
            payloads: artifact.results.payloads,
            bundles,
            sourceRecords: artifact.sources,
            evaluationId: artifact.manifest.evaluationId,
          }),
        )
        let pending = (artifact.state.pendingEvaluatorCalls ?? []).find(
          (call) => call.workUnitId === unit.id,
        )
        if (pending === undefined) {
          if ((artifact.state.pendingEvaluatorCalls ?? []).length >= artifact.manifest.concurrency)
            continue
          const attempt = (state?.attempts ?? 0) + 1
          pending = {
            requestId: yield* Effect.promise(() =>
              randomRequestId(artifact.manifest.evaluationId, unit.id, protocol.inputHash, attempt),
            ),
            workUnitId: unit.id,
            inputHash: protocol.inputHash,
            correlationId: protocol.correlationId,
            attempt,
          }
          ;(artifact.state.pendingEvaluatorCalls ??= []).push(pending)
          artifact.state.workUnits[unit.id] ??= { status: "pending" }
        } else if (pending.inputHash !== protocol.inputHash) {
          return failureResult(
            `the model or source for ${unit.id} changed after the work request was emitted; a result must not be attached to different evidence — start a new run`,
          )
        }
        requests.push(requestReceipt(protocol, pending, retryFailures.get(pending.requestId)))
      }
    }
    const pending = artifact.state.pendingEvaluatorCalls ?? []
    if (artifact.state.status !== "failed") {
      if (pending.length > 0) artifact.state.status = "awaiting_evaluator"
      else if (graph.every((unit) => artifact.state.workUnits[unit.id]?.status === "completed")) {
        artifact.state.status = "completed"
        artifact.state.completedAt = now
      }
    }
    artifact.state.updatedAt = now
    if (artifact.state.status !== initialStatus) {
      events.push({
        timestamp: now,
        event: "run_status",
        status: artifact.state.status,
        ...(artifact.state.failure === undefined ? {} : { failure: artifact.state.failure }),
      })
    }
    if (pending.length === 0) delete artifact.state.pendingEvaluatorCalls
    const temp = yield* fs.makeTempFile({ directory: runAbs, prefix: ".evaluation." })
    yield* fs.writeFileString(temp, jsonDocument(artifact), { mode: 0o644 })
    yield* fs.rename(temp, artifactPath)
    if (events.length > 0) {
      yield* fs.writeFileString(
        paths.join(runAbs, "logs/events.jsonl"),
        `${events.map((entry) => JSON.stringify(entry)).join("\n")}\n`,
        { flag: "a", mode: 0o600 },
      )
    }
    const completed = Object.values(artifact.state.workUnits).filter(
      (state) => state.status === "completed",
    ).length
    const result = {
      schemaVersion: 3,
      path: displayRun,
      status: artifact.state.status,
      evaluator: artifact.manifest.evaluator,
      evaluatorKind: artifact.manifest.evaluatorKind,
      concurrency: artifact.manifest.concurrency,
      workUnits: {
        total: graph.length,
        evaluatorUnits: graph.filter((unit) => unit.evaluatorBacked).length,
        completed,
      },
      sources: plan.areas.map((area) => ({
        area: area.ref,
        selector: artifact.sources[area.ref]!.selector,
        kind: artifact.sources[area.ref]!.kind,
        resolver: artifact.sources[area.ref]!.resolver,
      })),
      ...(requests.length === 0 ? {} : { evaluatorRequests: requests }),
      ...(artifact.state.status === "completed" && artifact.outputs !== undefined
        ? {
            reportMd: `${displayRun}/report.md`,
            ratingResult: artifact.outputs.rating,
          }
        : {}),
      nextActions:
        artifact.state.status === "awaiting_evaluator"
          ? [
              {
                id: "evaluation-evaluator-result",
                label: "Submit harness judgment results for outstanding work requests",
                command: `qualitymd evaluation run --resume ${displayRun} --evaluator-result - --json`,
              },
              {
                id: "evaluation-run-reemit",
                label: "Recover the outstanding work requests",
                command: `qualitymd evaluation run --resume ${displayRun} --json`,
              },
            ]
          : artifact.state.status === "completed"
            ? [
                {
                  id: "evaluation-report-read",
                  label: "Read the evaluation report",
                  command: `cat ${displayRun}/report.md`,
                },
              ]
            : [],
    }
    if (input.json) return commandResult(jsonDocument(result))
    if (artifact.state.status === "awaiting_evaluator") {
      return commandResult("", {
        stderr:
          `Awaiting harness judgment: ${requests.length} outstanding of up to ${artifact.manifest.concurrency} work requests\n` +
          requests
            .map(
              (request) =>
                `- ${request.workUnitId} (request ${request.requestId}, attempt ${request.attempt})\n`,
            )
            .join("") +
          "Run with --json to receive the bounded work requests.\n\n" +
          `Next: ${result.nextActions[0]!.command}\n`,
      })
    }
    return commandResult(`Evaluation ${artifact.state.status}: ${displayRun}\n`)
  }).pipe(
    Effect.mapError(
      (cause) =>
        new FileSystemFailure({
          detail: cause instanceof Error ? cause.message : String(cause),
        }),
    ),
  )
