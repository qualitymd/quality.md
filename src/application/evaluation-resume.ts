import * as Effect from "effect/Effect"
import * as FileSystem from "effect/FileSystem"
import * as Path from "effect/Path"
import * as Result from "effect/Result"

import { FileSystemFailure } from "../domain/errors.ts"
import { commandResult, ExitCode, type CommandResult } from "../domain/command-result.ts"
import { deterministicPayload } from "../domain/evaluation/frames.ts"
import { buildGraph, type WorkUnit } from "../domain/evaluation/graph.ts"
import { planEvaluation } from "../domain/evaluation/plan.ts"
import {
  buildProtocolRequest,
  completeProtocolRequest,
  protocolRequestReceipt,
  type StoredPayload,
} from "../domain/evaluation/protocol.ts"
import { normalizeEvaluatorResult, ResultValidationError } from "../domain/evaluation/result.ts"
import { jsonDocument } from "../domain/json.ts"
import { decodeModel } from "../domain/model/model.ts"
import { parseQualityDocument } from "../domain/model/document.ts"
import {
  EvidenceValidationError,
  sealEvidenceManifest,
  type AreaSource,
  type SealedEvidenceManifest,
} from "../services/source.ts"
import { HostRuntime } from "../services/host-runtime.ts"
import { hashJsonEffect, requestId } from "./evaluation-hash.ts"
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
    areaSources: Record<string, AreaSource>
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
  evidence: Record<string, SealedEvidenceManifest>
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
  "evidence_invalid",
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
    if (artifact.schemaVersion !== 8 || artifact.kind !== "EvaluationRun")
      return failureResult(
        `evaluation artifact schema ${artifact.schemaVersion} is incompatible with schema 8; start a new run`,
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
    const document = parseQualityDocument(snapshotPath, snapshot)
    const bodyGuidance = document.body.trim()
    const model = decodeModel(document)
    const plan = planEvaluation(model, artifact.manifest.plannedScope)
    const graph = buildGraph(plan)
    const now = new Date(yield* runtime.currentTimeMillis).toISOString().replace(/\.\d{3}Z$/, "Z")
    const initialStatus = artifact.state.status
    let events: ReadonlyArray<JsonObject> = []
    let retryFailures: Readonly<
      Record<string, { readonly category: string; readonly detail?: string }>
    > = {}
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
        const draft = buildProtocolRequest({
          unit,
          plan,
          payloads: artifact.results.payloads,
          areaSources: artifact.manifest.areaSources,
          bodyGuidance,
          evaluationId: artifact.manifest.evaluationId,
        })
        const protocol = completeProtocolRequest(draft, yield* hashJsonEffect(draft.hashInput))
        if (protocol.inputHash !== pending.inputHash)
          return failureResult(
            `the model or request context for ${unit.id} changed after the work request was emitted; a result must not be attached to a different request — start a new run`,
          )
        const state = (artifact.state.workUnits[unit.id] ??= { status: "pending" })
        state.attempts = (state.attempts ?? 0) + 1
        let category = envelope.failure ?? ""
        let detail = envelope.detail ?? ""
        let accepted: ReadonlyArray<StoredPayload> = []
        let sealedEvidence: SealedEvidenceManifest | undefined
        if (category === "") {
          try {
            const bytes = yield* runtime.randomBytes(100)
            const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
            accepted = normalizeEvaluatorResult(unit, envelope.payload, (index) =>
              Array.from(
                bytes.slice(index * 10, index * 10 + 10),
                (byte) => alphabet[byte % alphabet.length],
              ).join(""),
            )
            if (unit.kind === "assessRateRequirement") {
              if (
                envelope.payload === null ||
                Array.isArray(envelope.payload) ||
                typeof envelope.payload !== "object"
              )
                throw new EvidenceValidationError(
                  "combined requirement judgment must carry an evidence manifest",
                )
              const proposed = envelope.payload as JsonObject
              const assessment = accepted.find(
                (entry) => entry.payload.kind === "RequirementAssessmentResult",
              )?.payload
              const requirement = plan.requirements.find((entry) => entry.ref === unit.subject)!
              const evidenceResult = yield* sealEvidenceManifest({
                workspaceRoot,
                requirementId: unit.subject,
                source: artifact.manifest.areaSources[requirement.areaId]!,
                proposal: proposed.evidence,
                assessment,
                capturedAt: now,
              }).pipe(Effect.result)
              if (Result.isFailure(evidenceResult)) throw evidenceResult.failure
              sealedEvidence = evidenceResult.success
            }
          } catch (cause) {
            category =
              cause instanceof ResultValidationError
                ? cause.category
                : cause instanceof EvidenceValidationError
                  ? "evidence_invalid"
                  : "invalid_evaluator_output"
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
          ...(sealedEvidence === undefined
            ? {}
            : {
                evidence: {
                  observations: sealedEvidence.observations.length,
                  manifestHash: sealedEvidence.manifestHash,
                },
              }),
          ...(category === "" ? {} : { failure: { category, detail } }),
        })
        if (category === "") {
          if (artifact.manifest.evaluatorKind === "harness")
            artifact.state.harnessIdentity ??= { runtime: envelope.evaluator.runtime }
          artifact.state.pendingEvaluatorCalls = (
            artifact.state.pendingEvaluatorCalls ?? []
          ).filter((call) => call.requestId !== pending.requestId)
          mergePayloads(artifact, graph, unit.id, accepted)
          if (sealedEvidence !== undefined) artifact.evidence[unit.id] = sealedEvidence
          Object.assign(state, {
            status: "completed",
            inputHash: protocol.inputHash,
            completedAt: now,
            ...(sealedEvidence === undefined ? {} : { evidenceHash: sealedEvidence.manifestHash }),
          })
          delete state.failure
          events = [
            ...events,
            {
              timestamp: now,
              event: "work_unit_completed",
              workUnit: unit.id,
              attempt: state.attempts,
            },
          ]
        } else if (retryable.has(category) && state.attempts < 3) {
          pending.attempt = state.attempts + 1
          pending.requestId = yield* requestId(
            artifact.manifest.evaluationId,
            unit.id,
            pending.inputHash,
            pending.attempt,
          )
          retryFailures = {
            ...retryFailures,
            [pending.requestId]: { category, ...(detail === "" ? {} : { detail }) },
          }
          events = [
            ...events,
            {
              timestamp: now,
              event: "work_unit_retry",
              workUnit: unit.id,
              attempt: state.attempts,
              failure: { category, ...(detail === "" ? {} : { detail }) },
            },
          ]
        } else {
          artifact.state.pendingEvaluatorCalls = (
            artifact.state.pendingEvaluatorCalls ?? []
          ).filter((call) => call.requestId !== pending.requestId)
          state.status = "failed"
          state.failure = { category, ...(detail === "" ? {} : { detail }) }
          artifact.state.status = "failed"
          artifact.state.failure = state.failure
          artifact.state.completedAt = now
          events = [
            ...events,
            {
              timestamp: now,
              event: "work_unit_failed",
              workUnit: unit.id,
              attempt: state.attempts,
              failure: state.failure,
            },
          ]
        }
      }
      if (remaining.length > 0)
        return failureResult(
          `the submitted result ${remaining[0]!.requestId} does not correlate with an outstanding work request; resume without --evaluator-result to recover the outstanding requests`,
        )
    }
    let requests: ReadonlyArray<JsonObject> = []
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
            events = [
              ...events,
              { timestamp: now, event: "work_unit_completed", workUnit: unit.id },
            ]
            continue
          }
          const payload = deterministicPayload(unit, model, plan, artifact.manifest.model)
          const inputHash = yield* hashJsonEffect(payload)
          mergePayloads(artifact, graph, unit.id, [{ workUnit: unit.id, payload }])
          artifact.state.workUnits[unit.id] = {
            status: "completed",
            inputHash,
            completedAt: now,
          }
          events = [...events, { timestamp: now, event: "work_unit_completed", workUnit: unit.id }]
          continue
        }
        const draft = buildProtocolRequest({
          unit,
          plan,
          payloads: artifact.results.payloads,
          areaSources: artifact.manifest.areaSources,
          bodyGuidance,
          evaluationId: artifact.manifest.evaluationId,
        })
        const protocol = completeProtocolRequest(draft, yield* hashJsonEffect(draft.hashInput))
        let pending = (artifact.state.pendingEvaluatorCalls ?? []).find(
          (call) => call.workUnitId === unit.id,
        )
        if (pending === undefined) {
          if ((artifact.state.pendingEvaluatorCalls ?? []).length >= artifact.manifest.concurrency)
            continue
          const attempt = (state?.attempts ?? 0) + 1
          pending = {
            requestId: yield* requestId(
              artifact.manifest.evaluationId,
              unit.id,
              protocol.inputHash,
              attempt,
            ),
            workUnitId: unit.id,
            inputHash: protocol.inputHash,
            correlationId: protocol.correlationId,
            attempt,
          }
          artifact.state.pendingEvaluatorCalls = [
            ...(artifact.state.pendingEvaluatorCalls ?? []),
            pending,
          ]
          artifact.state.workUnits[unit.id] ??= { status: "pending" }
        } else if (pending.inputHash !== protocol.inputHash) {
          return failureResult(
            `the model or request context for ${unit.id} changed after the work request was emitted; a result must not be attached to a different request — start a new run`,
          )
        }
        requests = [
          ...requests,
          protocolRequestReceipt(protocol, pending, retryFailures[pending.requestId]),
        ]
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
      events = [
        ...events,
        {
          timestamp: now,
          event: "run_status",
          status: artifact.state.status,
          ...(artifact.state.failure === undefined ? {} : { failure: artifact.state.failure }),
        },
      ]
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
        selector: artifact.manifest.areaSources[area.ref]!.selector,
        kind: artifact.manifest.areaSources[area.ref]!.kind,
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
