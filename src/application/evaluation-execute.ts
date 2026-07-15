import * as Effect from "effect/Effect"
import * as FileSystem from "effect/FileSystem"
import * as Path from "effect/Path"

import { FileSystemFailure } from "../domain/errors.ts"
import { commandResult, type CommandResult } from "../domain/command-result.ts"
import { buildGraph } from "../domain/evaluation/graph.ts"
import { planEvaluation } from "../domain/evaluation/plan.ts"
import { buildProtocolRequest } from "../domain/evaluation/protocol.ts"
import { newEvaluationIdentity, resolveScope, scopeSlug } from "../domain/evaluation/run.ts"
import { hashJson, jsonDocument } from "../domain/json.ts"
import { decodeModel, type QualityModel } from "../domain/model/model.ts"
import { parseQualityDocument } from "../domain/model/document.ts"
import { harnessCapabilities } from "../adapters/evaluator.ts"
import { HostRuntime } from "../services/host-runtime.ts"
import { detectSourceKind, packageSource, type SourceBundle } from "../services/source.ts"
import { resolveWorkspace } from "../services/workspace.ts"
import type { EvaluationRunInput } from "./evaluation-run.ts"

const ratingIds = (model: QualityModel) => model.ratingScale.map((level) => `rating:${level.level}`)

const evaluationFrame = (model: QualityModel, modelPath: string) => ({
  derivedContext: { evaluationPolicies: ["source-as-data", "secret-redaction"], rigor: "standard" },
  inputs: { ratingLevelIds: ratingIds(model) },
  kind: "EvaluationFrame",
  schemaVersion: 3,
  subject: { modelLocator: modelPath },
})

const nextNumber = (directory: string) =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const paths = yield* Path.Path
    if (!(yield* fs.exists(directory))) return 1
    let maximum = 0
    for (const name of yield* fs.readDirectory(directory)) {
      let number: number | undefined
      for (const fileName of ["evaluation.json", "data/evaluation-manifest.json"]) {
        const file = paths.join(directory, name, fileName)
        if (!(yield* fs.exists(file))) continue
        try {
          const value = JSON.parse(yield* fs.readFileString(file)) as {
            readonly manifest?: { readonly run?: { readonly number?: number } }
            readonly run?: { readonly number?: number }
          }
          number = value.manifest?.run?.number ?? value.run?.number
        } catch {
          // Fall through to current run-folder parsing.
        }
        break
      }
      if (number === undefined) {
        const match = /^(\d{4})-([a-z0-9-]+)-eval$/.exec(name)
        if (match !== null && !match[2]!.split("-").includes("quality")) number = Number(match[1])
      }
      maximum = Math.max(maximum, number ?? 0)
    }
    return maximum + 1
  })

const requirementFrame = (
  model: QualityModel,
  requirement: ReturnType<typeof planEvaluation>["requirements"][number],
) => {
  const criteria = model.ratingScale.map((level) => ({
    criterion: requirement.value.ratings?.[level.level] ?? level.criterion,
    ratingLevelId: `rating:${level.level}`,
    source:
      requirement.value.ratings?.[level.level] === undefined
        ? "model_default"
        : "requirement_override",
  }))
  return {
    derivedContext: { appliedRatingCriteria: criteria },
    inputs: {
      ratingLevelIds: ratingIds(model),
      requirementAssessmentBasis: requirement.value.assessment,
    },
    kind: "RequirementEvaluationFrame",
    schemaVersion: 3,
    subject: { factorIds: requirement.factorIds, requirementId: requirement.ref },
  }
}

const atomicWrite = (path: string, content: string) =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const paths = yield* Path.Path
    const temp = yield* fs.makeTempFile({ directory: paths.dirname(path), prefix: ".evaluation." })
    yield* fs.writeFileString(temp, content, { mode: 0o644 })
    yield* fs.rename(temp, path)
  })

export const executeHarnessRun = (
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
    const workspace = yield* resolveWorkspace({
      ...(input.model === "" ? {} : { model: input.model }),
      ...(input.evaluationDir === "" ? {} : { evaluationDir: input.evaluationDir }),
    })
    const raw = yield* fs.readFileString(workspace.model.abs)
    const document = parseQualityDocument(workspace.model.abs, raw)
    const model = decodeModel(document)
    const scope = resolveScope(model, input.area, input.factors)
    const plan = planEvaluation(model, scope.plannedScope)
    const concurrency = workspace.evaluation.concurrency ?? runtime.hardwareConcurrency * 2
    const resolvedConcurrency =
      concurrency > 1 && !harnessCapabilities.concurrent && !harnessCapabilities.subagents
        ? 1
        : concurrency
    const number = yield* nextNumber(workspace.evaluations.abs)
    const label = `${String(number).padStart(4, "0")}-${scopeSlug(scope.plannedScope)}-eval`
    const runAbs = paths.join(workspace.evaluations.abs, label)
    const runRel = `${workspace.evaluations.rel}/${label}`
    const identity = newEvaluationIdentity(
      new Date(yield* runtime.currentTimeMillis),
      yield* runtime.randomBytes(12),
    )
    const timestamp = identity.createdAt
    const payloads: Array<{
      readonly workUnit: string
      readonly payload: Record<string, unknown>
    }> = []
    const workUnits: Record<string, Record<string, unknown>> = {}
    const complete = async (workUnit: string, payload: Record<string, unknown>) => {
      payloads.push({ workUnit, payload })
      workUnits[workUnit] = {
        status: "completed",
        inputHash: await hashJson(payload),
        completedAt: timestamp,
      }
    }
    yield* Effect.promise(() =>
      complete("frameEvaluation", evaluationFrame(model, workspace.model.rel)),
    )
    const areaFrames = new Map<string, Record<string, unknown>>()
    for (const area of plan.areas) {
      const frame = {
        inputs: {
          childAreaIds: area.childAreaIds,
          localRequirementIds: area.localRequirementIds,
          rootFactorIds: area.rootFactorIds,
          sourceRefs: [area.source],
        },
        kind: "AreaEvaluationFrame",
        schemaVersion: 3,
        subject: { areaId: area.ref },
      }
      areaFrames.set(area.ref, frame)
      yield* Effect.promise(() => complete(`frameAreaEvaluation:${area.ref}`, frame))
    }
    const requirementFrames = new Map<string, Record<string, unknown>>()
    for (const requirement of plan.requirements) {
      const frame = requirementFrame(model, requirement)
      requirementFrames.set(requirement.ref, frame)
      yield* Effect.promise(() => complete(`frameRequirementEvaluation:${requirement.ref}`, frame))
    }
    const sourceRecords: Record<string, Record<string, unknown>> = {}
    const bundles = new Map<string, SourceBundle>()
    const sourcePlans = []
    for (const area of plan.areas) {
      const kind = yield* detectSourceKind(workspace.workspaceRoot.abs, area.source)
      const resolver = kind === "prose" ? "harness" : "walk"
      sourceRecords[area.ref] = { selector: area.source, kind, resolver }
      sourcePlans.push({ area: area.ref, selector: area.source, kind, resolver })
      if (area.localRequirementIds.length === 0) continue
      if (kind === "prose") continue
      const bundle = yield* packageSource(workspace.workspaceRoot.abs, area.source, kind)
      bundles.set(area.ref, bundle)
      sourceRecords[area.ref] = {
        ...sourceRecords[area.ref],
        bundleHash: bundle.hash,
        capturedAt: timestamp,
        ...(bundle.truncated ? { truncated: true } : {}),
        files: bundle.files.map(({ path, sha256, truncated }) => ({
          path,
          sha256,
          ...(truncated === true ? { truncated: true } : {}),
        })),
      }
    }
    const sourceKinds = Object.fromEntries(
      Object.entries(sourceRecords).map(([area, record]) => [area, record.kind]),
    ) as Record<string, "path" | "glob" | "prose">
    const graph = buildGraph(plan, sourceKinds)
    const requests = []
    const pending = []
    const completed = new Set(Object.keys(workUnits))
    for (const unit of graph) {
      if (
        !unit.evaluatorBacked ||
        pending.length >= resolvedConcurrency ||
        unit.dependsOn.some((dependency) => !completed.has(dependency))
      )
        continue
      const protocol = yield* Effect.promise(() =>
        buildProtocolRequest({
          unit,
          plan,
          payloads,
          bundles,
          sourceRecords,
          evaluationId: identity.evaluationId,
        }),
      )
      const workUnitId = unit.id
      workUnits[workUnitId] = { status: "pending" }
      const inputHash = protocol.inputHash
      const correlationId = `${identity.evaluationId}#${workUnitId}`
      const requestId = `req_${(yield* Effect.promise(() =>
        hashJson({
          evaluationId: identity.evaluationId,
          workUnit: workUnitId,
          inputHash,
          attempt: 1,
        }),
      )).slice(0, 16)}`
      pending.push({ requestId, workUnitId, inputHash, correlationId, attempt: 1 })
      requests.push({
        requestId,
        workUnitId,
        kind: protocol.kind,
        ...(protocol.subject === "" ? {} : { subject: protocol.subject }),
        attempt: 1,
        instructions: protocol.instructions,
        ...(protocol.sharedContext === null ? {} : { sharedContext: protocol.sharedContext }),
        ...(protocol.context === null ? {} : { context: protocol.context }),
        ...(protocol.source.length === 0 ? {} : { source: protocol.source }),
        expectedSchema: protocol.expectedSchema,
        inputHash,
        correlationId,
      })
    }
    const evaluatorUnits = graph.filter((unit) => unit.evaluatorBacked).length
    const total = graph.length
    const artifact = {
      schemaVersion: 7,
      kind: "EvaluationRun",
      manifest: {
        ...identity,
        model: workspace.model.rel,
        ...scope,
        run: { number, label },
        evaluator: "harness",
        evaluatorKind: "harness",
        evaluatorCapabilities: harnessCapabilities,
        concurrency: resolvedConcurrency,
      },
      state: {
        status: "awaiting_evaluator",
        workUnits,
        startedAt: timestamp,
        updatedAt: timestamp,
        pendingEvaluatorCalls: pending,
      },
      sources: sourceRecords,
      results: { payloads },
    }
    yield* fs.makeDirectory(workspace.evaluations.abs, { recursive: true, mode: 0o755 })
    yield* fs.makeDirectory(runAbs, { mode: 0o755 })
    yield* fs.makeDirectory(paths.join(runAbs, "logs"), { mode: 0o755 })
    yield* fs.writeFileString(paths.join(runAbs, "model-snapshot.md"), raw, { mode: 0o644 })
    yield* atomicWrite(paths.join(runAbs, "evaluation.json"), jsonDocument(artifact))
    yield* fs.writeFileString(paths.join(runAbs, "logs/evaluator-calls.jsonl"), "", {
      mode: 0o600,
    })
    yield* fs.writeFileString(
      paths.join(runAbs, "logs/events.jsonl"),
      [
        {
          timestamp,
          event: "run_created",
          evaluationId: identity.evaluationId,
          evaluator: "harness",
          evaluatorKind: "harness",
          capabilities: harnessCapabilities,
        },
        {
          timestamp,
          event: "run_status",
          status: "awaiting_evaluator",
          outstanding: pending.length,
        },
      ]
        .map((entry) => JSON.stringify(entry))
        .join("\n") + "\n",
      { mode: 0o600 },
    )
    const result = {
      schemaVersion: 3,
      path: runRel,
      status: "awaiting_evaluator",
      evaluator: "harness",
      evaluatorKind: "harness",
      concurrency: resolvedConcurrency,
      workUnits: { total, evaluatorUnits, completed: payloads.length },
      sources: sourcePlans,
      evaluatorRequests: requests,
      nextActions: [
        {
          id: "evaluation-evaluator-result",
          label: "Submit harness judgment results for outstanding work requests",
          command: `qualitymd evaluation run --resume ${runRel} --evaluator-result - --json`,
        },
        {
          id: "evaluation-run-reemit",
          label: "Recover the outstanding work requests",
          command: `qualitymd evaluation run --resume ${runRel} --json`,
        },
      ],
    }
    if (input.json) return commandResult(jsonDocument(result))
    return commandResult("", {
      stderr:
        `Awaiting harness judgment: ${requests.length} outstanding of up to ${resolvedConcurrency} work requests\n` +
        requests
          .map(
            (request) =>
              `- ${request.workUnitId} (request ${request.requestId}, attempt ${request.attempt})\n`,
          )
          .join("") +
        "Run with --json to receive the bounded work requests.\n\n" +
        `Next: ${result.nextActions[0]!.command}\n`,
    })
  }).pipe(
    Effect.mapError(
      (cause) =>
        new FileSystemFailure({
          detail: cause instanceof Error ? cause.message : String(cause),
        }),
    ),
  )
