import { assert, describe, it } from "@effect/vitest"
import * as BunFileSystem from "@effect/platform-bun/BunFileSystem"
import * as BunPath from "@effect/platform-bun/BunPath"
import * as Deferred from "effect/Deferred"
import * as Effect from "effect/Effect"
import * as Fiber from "effect/Fiber"
import * as Layer from "effect/Layer"
import { mkdir, mkdtemp, readFile, readdir, rm, writeFile } from "node:fs/promises"
import { join } from "node:path"
import { fileURLToPath } from "node:url"

import { executeProviderRun } from "../../src/application/evaluation-provider.ts"
import type { EvaluationRequest } from "../../src/domain/evaluator/types.ts"
import evaluationExamples from "../../src/assets/evaluation-examples.json"
import type { EvaluatorService } from "../../src/services/evaluator.ts"
import { HostRuntime, type HostRuntimeService } from "../../src/services/host-runtime.ts"

const repositoryRoot = fileURLToPath(new URL("../..", import.meta.url))
const services = Layer.mergeAll(
  BunFileSystem.layer,
  BunPath.layer,
  Layer.succeed(HostRuntime, {
    cwd: repositoryRoot,
    environment: {},
    hardwareConcurrency: 4,
    currentTimeMillis: Effect.succeed(Date.UTC(2026, 6, 14)),
    randomBytes: (length) => Effect.succeed(new Uint8Array(length).fill(7)),
    readStdin: Effect.succeed(""),
    which: () => null,
    codexAuthenticated: () => false,
  } satisfies HostRuntimeService),
)

const clone = <A>(value: A): A => structuredClone(value)

const replaceReference = (value: unknown, from: string, to: string): unknown => {
  if (value === from) return to
  if (Array.isArray(value)) return value.map((entry) => replaceReference(entry, from, to))
  if (value !== null && typeof value === "object")
    return Object.fromEntries(
      Object.entries(value).map(([key, child]) => [key, replaceReference(child, from, to)]),
    )
  return value
}

const replaceEvidenceRefs = (value: unknown): unknown => {
  if (Array.isArray(value)) return value.map(replaceEvidenceRefs)
  if (value !== null && typeof value === "object")
    return Object.fromEntries(
      Object.entries(value).map(([key, child]) => [
        key,
        key === "sourceRef" ? "evidence[ev-001]" : replaceEvidenceRefs(child),
      ]),
    )
  return value
}

const payloadFor = (request: EvaluationRequest): Record<string, unknown> => {
  const requirement = "requirement:root::has-tests"
  if (request.kind === "assessRateRequirement") {
    const assessment = replaceEvidenceRefs(
      replaceReference(
        clone(evaluationExamples.RequirementAssessmentResult),
        requirement,
        request.subject,
      ),
    ) as Record<string, unknown>
    const rating = replaceReference(
      clone(evaluationExamples.RequirementRatingResult),
      requirement,
      request.subject,
    ) as Record<string, unknown>
    return {
      assessment,
      rating,
      evidence: {
        observations: [{ id: "ev-001", kind: "file", role: "evaluated", path: "evidence.txt" }],
        limits: [],
      },
    }
  }
  if (request.kind === "analyzeFactor") {
    const result = clone(evaluationExamples.FactorAnalysisResult) as Record<string, unknown>
    result.factorId = request.subject
    return result
  }
  if (request.kind === "analyzeArea") {
    const result = clone(evaluationExamples.AreaAnalysisResult) as Record<string, unknown>
    result.areaId = request.subject
    return result
  }
  if (request.kind === "rankFindings") {
    const result = replaceReference(
      clone(evaluationExamples.FindingRankingResult),
      requirement,
      "requirement:root::exists",
    ) as Record<string, unknown>
    result.orderedFindings = (result.orderedFindings as Array<Record<string, unknown>>).slice(0, 2)
    return result
  }
  if (request.kind === "recommend") {
    const recommendation = replaceReference(
      clone(evaluationExamples.RecommendationResult),
      requirement,
      "requirement:root::exists",
    )
    return { recommendations: [recommendation] }
  }
  if (request.kind === "rankRecommendations") {
    const result = replaceReference(
      clone(evaluationExamples.RecommendationRankingResult),
      requirement,
      "requirement:root::exists",
    ) as Record<string, unknown>
    return result
  }
  throw new Error(`unexpected evaluator work ${request.kind}`)
}

const capabilities = {
  structuredOutput: true,
  workspaceInspection: true,
  instructionIsolation: true,
  verification: false,
  networkAccess: "disabled",
  tools: true,
  concurrent: true,
  subagents: false,
  freshContext: true,
  cancellation: true,
  usage: true,
  maxTurns: "unsupported",
  tokenBudget: "supported",
  costBudget: "advisory",
  contextWindow: "reported",
  compaction: "opaque",
  sandbox: "provider",
  executableOverride: false,
} as const

describe("provider evaluation", () => {
  it.effect("completes the work graph and builds reports with reconstructible requests", () =>
    Effect.gen(function* () {
      const directory = yield* Effect.acquireRelease(
        Effect.promise(async () => {
          await mkdir(join(repositoryRoot, "tmp"), { recursive: true })
          return mkdtemp(join(repositoryRoot, "tmp", "qualitymd-provider-test-"))
        }),
        (path) => Effect.promise(() => rm(path, { recursive: true })),
      )
      const model = join(directory, "QUALITY.md")
      yield* Effect.promise(() => writeFile(join(directory, "evidence.txt"), "bounded evidence"))
      yield* Effect.promise(() =>
        writeFile(
          model,
          `---
title: Test
source: evidence.txt
ratingScale:
  - level: target
    title: Target
    criterion: Satisfies the requirement.
  - level: unacceptable
    title: Unacceptable
    criterion: Does not satisfy the requirement.
factors:
  evidence:
    title: Evidence
    requirements:
      exists:
        title: Evidence exists
        assessment: Inspect the selected source.
---
`,
        ),
      )
      const evaluator: EvaluatorService = {
        name: "mock-provider",
        kind: "claude",
        capabilities: { ...capabilities, concurrent: false },
        evaluate: (request) =>
          Effect.succeed({
            workUnitId: request.workUnitId,
            evaluatorKind: "claude",
            model: "mock-model",
            payload: payloadFor(request),
            usage: { inputTokens: 10, outputTokens: 5 },
          }),
      }
      const result = yield* executeProviderRun(
        {
          model,
          evaluationDir: "",
          area: "",
          factors: [],
          evaluator: evaluator.name,
          resume: "",
          evaluatorResult: "",
          dryRun: false,
          json: true,
        },
        evaluator,
      ).pipe(Effect.provide(services))
      assert.strictEqual(result.exitCode, 0)
      const receipt = JSON.parse(result.stdout) as {
        path: string
        status: string
        reportMd: string
      }
      assert.strictEqual(receipt.status, "completed")
      const runPath = join(directory, receipt.path)
      assert.match(
        yield* Effect.promise(() => readFile(join(runPath, "report.md"), "utf8")),
        /# Quality evaluation/,
      )
      const artifact = JSON.parse(
        yield* Effect.promise(() => readFile(join(runPath, "evaluation.json"), "utf8")),
      ) as {
        schemaVersion: number
        manifest: { concurrency: number }
        state: { status: string }
        evidence: Record<
          string,
          { observations: ReadonlyArray<{ path: string; sha256: string }>; manifestHash: string }
        >
      }
      assert.strictEqual(artifact.schemaVersion, 8)
      assert.strictEqual(artifact.manifest.concurrency, 1)
      assert.strictEqual(artifact.state.status, "completed")
      const evidence = Object.values(artifact.evidence)[0]!
      assert.strictEqual(evidence.observations[0]?.path, "evidence.txt")
      assert.match(evidence.observations[0]?.sha256 ?? "", /^[a-f0-9]{64}$/)
      assert.match(evidence.manifestHash, /^[a-f0-9]{64}$/)
      assert.match(
        yield* Effect.promise(() => readFile(join(runPath, "logs", "events.jsonl"), "utf8")),
        /"status":"completed"/,
      )
    }),
  )

  it.effect("interrupts in-flight evaluator work while leaving a resumable artifact", () =>
    Effect.gen(function* () {
      const directory = yield* Effect.acquireRelease(
        Effect.promise(async () => {
          await mkdir(join(repositoryRoot, "tmp"), { recursive: true })
          return mkdtemp(join(repositoryRoot, "tmp", "qualitymd-cancel-test-"))
        }),
        (path) => Effect.promise(() => rm(path, { recursive: true })),
      )
      const model = join(directory, "QUALITY.md")
      yield* Effect.promise(() => writeFile(join(directory, "evidence.txt"), "bounded evidence"))
      yield* Effect.promise(() =>
        writeFile(
          model,
          `---
title: Test
source: evidence.txt
ratingScale:
  - level: target
    title: Target
    criterion: Satisfies the requirement.
  - level: unacceptable
    title: Unacceptable
    criterion: Does not satisfy the requirement.
factors:
  evidence:
    title: Evidence
    requirements:
      exists:
        title: Evidence exists
        assessment: Inspect the selected source.
---
`,
        ),
      )
      const started = yield* Deferred.make<void>()
      let finalized = false
      const evaluator: EvaluatorService = {
        name: "interruptible-provider",
        kind: "claude",
        capabilities: { ...capabilities, concurrent: false },
        evaluate: () =>
          Deferred.succeed(started, undefined).pipe(
            Effect.andThen(Effect.never),
            Effect.ensuring(Effect.sync(() => void (finalized = true))),
          ),
      }
      const fiber = yield* Effect.forkChild(
        executeProviderRun(
          {
            model,
            evaluationDir: "",
            area: "",
            factors: [],
            evaluator: evaluator.name,
            resume: "",
            evaluatorResult: "",
            dryRun: false,
            json: true,
          },
          evaluator,
        ).pipe(Effect.provide(services)),
      )
      yield* Deferred.await(started)
      yield* Fiber.interrupt(fiber)
      assert.strictEqual(finalized, true)

      const evaluationRoot = join(directory, ".quality", "evaluations")
      const runs = yield* Effect.promise(() => readdir(evaluationRoot))
      assert.strictEqual(runs.length, 1)
      const artifact = JSON.parse(
        yield* Effect.promise(() =>
          readFile(join(evaluationRoot, runs[0]!, "evaluation.json"), "utf8"),
        ),
      ) as { state: { status: string; pendingEvaluatorCalls?: ReadonlyArray<unknown> } }
      assert.strictEqual(artifact.state.status, "awaiting_evaluator")
      assert.ok((artifact.state.pendingEvaluatorCalls?.length ?? 0) > 0)
    }),
  )
})
