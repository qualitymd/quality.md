import { findElement, parseAreaReference, projectModel, type QualityModel } from "../model/model.ts"
import type { EvaluatorCapabilities, EvaluatorKind } from "../evaluator/types.ts"
import type { ConcurrencyResolution } from "./concurrency.ts"
import type { StoredPayload } from "./protocol.ts"

export const EvaluationSchemaVersion = 3

export interface PlannedScope {
  readonly areaId: string
  readonly factorFilter: ReadonlyArray<string>
}

export interface RequestedScope {
  readonly areaId?: string
  readonly factorFilter?: ReadonlyArray<string>
}

export interface EvaluationManifest {
  readonly schemaVersion: 3
  readonly kind: "EvaluationManifest"
  readonly evaluationId: string
  readonly createdAt: string
  readonly model: string
  readonly requestedScope: RequestedScope
  readonly plannedScope: PlannedScope
  readonly run: { readonly number: number; readonly label: string }
}

export interface RunDirectoryInput {
  readonly name: string
  readonly evaluationArtifact?: string | undefined
  readonly historicalManifest?: string | undefined
}

type JsonObject = Record<string, unknown>

export interface HarnessPendingCall {
  readonly requestId: string
  readonly workUnitId: string
  readonly inputHash: string
  readonly correlationId: string
  readonly attempt: number
}

export interface HarnessRequestReceipt extends JsonObject {
  readonly requestId: string
  readonly workUnitId: string
  readonly attempt: number
}

interface EvaluationRunOptions {
  readonly identity: { readonly evaluationId: string; readonly createdAt: string }
  readonly model: string
  readonly scope: {
    readonly requestedScope: RequestedScope
    readonly plannedScope: PlannedScope
  }
  readonly number: number
  readonly label: string
  readonly evaluator: {
    readonly name: string
    readonly kind: EvaluatorKind
    readonly capabilities: EvaluatorCapabilities
  }
  readonly concurrency: ConcurrencyResolution
  readonly areaSources: Readonly<
    Record<string, { readonly selector: string; readonly kind: string }>
  >
  readonly workUnits: Readonly<Record<string, JsonObject>>
  readonly pending: ReadonlyArray<HarnessPendingCall>
  readonly payloads: ReadonlyArray<StoredPayload>
}

export const evaluationRunArtifact = (options: EvaluationRunOptions) => ({
  schemaVersion: 10,
  kind: "EvaluationRun",
  manifest: {
    ...options.identity,
    model: options.model,
    ...options.scope,
    run: { number: options.number, label: options.label },
    evaluator: options.evaluator.name,
    evaluatorKind: options.evaluator.kind,
    evaluatorCapabilities: options.evaluator.capabilities,
    concurrency: options.concurrency.resolved,
    areaSources: options.areaSources,
  },
  state: {
    status: options.concurrency.mode === "delegated" ? "awaiting_evaluator" : "running",
    workUnits: options.workUnits,
    startedAt: options.identity.createdAt,
    updatedAt: options.identity.createdAt,
    pendingEvaluatorCalls: options.pending,
  },
  evidence: {},
  results: { payloads: options.payloads },
})

export const evaluationRunEvents = (
  timestamp: string,
  evaluationId: string,
  evaluator: EvaluationRunOptions["evaluator"],
  concurrency: ConcurrencyResolution,
  pending: number,
): string =>
  [
    {
      timestamp,
      event: "run_created",
      evaluationId,
      evaluator: evaluator.name,
      evaluatorKind: evaluator.kind,
      capabilities: evaluator.capabilities,
      dispatchMode: concurrency.mode,
      concurrencyResolution: {
        source: concurrency.source,
        requested: concurrency.requested,
        automatic: concurrency.automatic,
        ...(concurrency.maximum === undefined ? {} : { maximum: concurrency.maximum }),
        resolved: concurrency.resolved,
        clamped: concurrency.clamped,
      },
    },
    {
      timestamp,
      event: "run_status",
      status: concurrency.mode === "delegated" ? "awaiting_evaluator" : "running",
      ...(concurrency.mode === "delegated" ? { outstanding: pending } : { pending }),
      ...(concurrency.mode === "delegated" ? { peakOutstanding: pending } : {}),
    },
  ]
    .map((entry) => JSON.stringify(entry))
    .join("\n") + "\n"

export const evaluationRunReceipt = (options: {
  readonly path: string
  readonly evaluator: string
  readonly evaluatorKind: EvaluatorKind
  readonly concurrency: number
  readonly total: number
  readonly evaluatorUnits: number
  readonly completed: number
  readonly sources: ReadonlyArray<JsonObject>
  readonly requests: ReadonlyArray<HarnessRequestReceipt>
  readonly dispatchMode: "direct" | "delegated" | "sequential"
}) => ({
  schemaVersion: 3,
  path: options.path,
  status: options.dispatchMode === "delegated" ? "awaiting_evaluator" : "running",
  evaluator: options.evaluator,
  evaluatorKind: options.evaluatorKind,
  concurrency: options.concurrency,
  workUnits: {
    total: options.total,
    evaluatorUnits: options.evaluatorUnits,
    completed: options.completed,
  },
  sources: options.sources,
  evaluatorRequests: options.requests,
  nextActions: [
    {
      id: "evaluation-evaluator-result",
      label: "Submit harness judgment results for outstanding work requests",
      command: `qualitymd evaluation run --resume ${options.path} --evaluator-result - --json`,
    },
    {
      id: "evaluation-run-reemit",
      label: "Recover the outstanding work requests",
      command: `qualitymd evaluation run --resume ${options.path} --json`,
    },
  ],
})

export const renderHarnessAwaiting = (
  requests: ReadonlyArray<HarnessRequestReceipt>,
  concurrency: number,
  nextCommand: string,
): string =>
  `Awaiting harness judgment: ${requests.length} outstanding of up to ${concurrency} work requests\n` +
  requests
    .map(
      (request) =>
        `- ${request.workUnitId} (request ${request.requestId}, attempt ${request.attempt})\n`,
    )
    .join("") +
  "Run with --json to receive the bounded work requests.\n\n" +
  `Next: ${nextCommand}\n`

const alphabet = "0123456789abcdefghjkmnpqrstvwxyz"

const positiveInteger = (value: unknown): number | undefined =>
  typeof value === "number" && Number.isInteger(value) && value > 0 ? value : undefined

const documentNumber = (raw: string | undefined, kind: "artifact" | "manifest") => {
  if (raw === undefined) return undefined
  try {
    const parsed = JSON.parse(raw) as {
      readonly manifest?: { readonly run?: { readonly number?: unknown } }
      readonly run?: { readonly number?: unknown }
    }
    return positiveInteger(kind === "artifact" ? parsed.manifest?.run?.number : parsed.run?.number)
  } catch {
    return undefined
  }
}

export const runDirectoryNumber = (input: RunDirectoryInput): number | undefined => {
  const artifactNumber = documentNumber(input.evaluationArtifact, "artifact")
  if (artifactNumber !== undefined) return artifactNumber
  const manifestNumber = documentNumber(input.historicalManifest, "manifest")
  if (manifestNumber !== undefined) return manifestNumber
  const match = /^(\d{4})-[a-z0-9-]+-eval$/.exec(input.name)
  return match === null ? undefined : positiveInteger(Number(match[1]))
}

export const newEvaluationIdentity = (now: Date, bytes: Uint8Array) => {
  const tail = [...bytes].map((value) => alphabet[value & 31]).join("")
  const createdAt = `${now.toISOString().slice(0, 19)}Z`
  return {
    evaluationId: `${createdAt.replaceAll("-", "").replaceAll(":", "").replace(".000", "")}-${tail}`,
    createdAt,
  }
}

const factorArea = (reference: string) => {
  if (!reference.startsWith("factor:") || !reference.includes("::")) {
    throw new Error(`factor model reference ${JSON.stringify(reference)} is invalid`)
  }
  return `area:${reference.slice(7).split("::", 1)[0]}`
}

export const resolveScope = (
  model: QualityModel,
  area: string,
  factors: ReadonlyArray<string>,
): { readonly requestedScope: RequestedScope; readonly plannedScope: PlannedScope } => {
  const root = projectModel(model)
  if (area !== "") parseAreaReference(model, area)
  const initialArea = area === "" ? undefined : area
  const selection = factors.reduce(
    (current, reference) => {
      const requestedArea = current.requestedArea ?? factorArea(reference)
      const plannedArea =
        current.plannedArea === "area:root" && area === "" && current.filter.length === 0
          ? factorArea(reference)
          : current.plannedArea
      const element = findElement(root, reference)
      if (element?.kind !== "factor") {
        throw new Error(
          `--factor: factor model reference ${JSON.stringify(reference)} does not resolve in the model`,
        )
      }
      const owner = factorArea(reference)
      if (requestedArea !== owner || plannedArea !== owner) {
        throw new Error(`--factor ${reference} does not belong to --area ${plannedArea}`)
      }
      return { requestedArea, plannedArea, filter: [...current.filter, reference] }
    },
    {
      requestedArea: initialArea,
      plannedArea: area === "" ? "area:root" : area,
      filter: [] as ReadonlyArray<string>,
    },
  )
  return {
    requestedScope: {
      ...(selection.requestedArea === undefined ? {} : { areaId: selection.requestedArea }),
      ...(selection.filter.length === 0 ? {} : { factorFilter: selection.filter }),
    },
    plannedScope: { areaId: selection.plannedArea, factorFilter: selection.filter },
  }
}

export const scopeSlug = (scope: PlannedScope) => {
  if (scope.areaId === "area:root" && scope.factorFilter.length === 0) return "full"
  const area = scope.areaId.slice(5)
  const parts = [
    ...(area === "root" ? ["root"] : area.split("/")),
    ...scope.factorFilter.flatMap((factor) => factor.split("::")[1]!.split("/")),
  ]
  return parts.join("-")
}
