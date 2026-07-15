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

import { executeProviderRun, resumeProviderRun } from "../../src/application/evaluation-provider.ts"
import { EvaluatorFailure, type EvaluationRequest } from "../../src/domain/evaluator/types.ts"
import evaluationExamples from "../../src/assets/evaluation-examples.json"
import type { EvaluatorService } from "../../src/services/evaluator.ts"
import { HostRuntime, type HostRuntimeService } from "../../src/services/host-runtime.ts"
import { resolveWorkspace } from "../../src/services/workspace.ts"

const repositoryRoot = fileURLToPath(new URL("../..", import.meta.url))
const services = Layer.mergeAll(
  BunFileSystem.layer,
  BunPath.layer,
  Layer.succeed(HostRuntime, {
    cwd: repositoryRoot,
    environment: {},
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
  dispatch: {
    concurrentCalls: true,
    delegatedRequests: false,
    automaticConcurrency: 4,
  },
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
  it.effect("tops up a bounded direct pool after one durable completion", () =>
    Effect.gen(function* () {
      const directory = yield* Effect.acquireRelease(
        Effect.promise(async () => {
          await mkdir(join(repositoryRoot, "tmp"), { recursive: true })
          return mkdtemp(join(repositoryRoot, "tmp", "qualitymd-concurrency-test-"))
        }),
        (path) => Effect.promise(() => rm(path, { recursive: true })),
      )
      const model = join(directory, "QUALITY.md")
      yield* Effect.promise(() => writeFile(join(directory, "evidence.txt"), "bounded evidence"))
      yield* Effect.promise(() =>
        writeFile(
          model,
          `---
title: Concurrent test
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
      one:
        title: One
        assessment: Inspect the selected source.
      two:
        title: Two
        assessment: Inspect the selected source.
      three:
        title: Three
        assessment: Inspect the selected source.
---
`,
        ),
      )
      const gates = [
        yield* Deferred.make<void>(),
        yield* Deferred.make<void>(),
        yield* Deferred.make<void>(),
      ]
      const firstPairStarted = yield* Deferred.make<void>()
      const thirdStarted = yield* Deferred.make<void>()
      const started: Array<string> = []
      let active = 0
      let peak = 0
      const evaluator: EvaluatorService = {
        name: "concurrent-provider",
        kind: "codex",
        capabilities: {
          ...capabilities,
          dispatch: {
            concurrentCalls: true,
            delegatedRequests: false,
            automaticConcurrency: 2,
          },
        },
        evaluate: (request) => {
          if (request.kind !== "assessRateRequirement")
            return Effect.succeed({
              workUnitId: request.workUnitId,
              evaluatorKind: "codex",
              payload: payloadFor(request),
            })
          return Effect.suspend(() => {
            const index = started.length
            const gate = gates[index]!
            return Effect.gen(function* () {
              started.push(request.subject)
              active += 1
              peak = Math.max(peak, active)
              if (started.length === 2) yield* Deferred.succeed(firstPairStarted, undefined)
              if (started.length === 3) yield* Deferred.succeed(thirdStarted, undefined)
              yield* Deferred.await(gate)
              return {
                workUnitId: request.workUnitId,
                evaluatorKind: "codex" as const,
                payload: payloadFor(request),
              }
            }).pipe(Effect.ensuring(Effect.sync(() => void (active -= 1))))
          })
        },
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
      yield* Deferred.await(firstPairStarted)
      assert.strictEqual(active, 2)
      assert.strictEqual(peak, 2)

      const evaluationRoot = join(directory, ".quality", "evaluations")
      const runs = yield* Effect.promise(() => readdir(evaluationRoot))
      const runPath = join(evaluationRoot, runs[0]!)
      const initialArtifact = JSON.parse(
        yield* Effect.promise(() => readFile(join(runPath, "evaluation.json"), "utf8")),
      ) as {
        manifest: {
          evaluator: string
          evaluatorKind: string
          concurrency: number
          evaluatorCapabilities: {
            dispatch: {
              concurrentCalls: boolean
              delegatedRequests: boolean
              automaticConcurrency: number
              maxConcurrency?: number
            }
          }
        }
        state: { status: string; pendingEvaluatorCalls: ReadonlyArray<unknown> }
      }
      assert.strictEqual(initialArtifact.manifest.evaluator, evaluator.name)
      assert.strictEqual(initialArtifact.manifest.evaluatorKind, evaluator.kind)
      assert.strictEqual(initialArtifact.manifest.concurrency, 2)
      assert.deepStrictEqual(
        initialArtifact.manifest.evaluatorCapabilities.dispatch,
        evaluator.capabilities.dispatch,
      )
      assert.strictEqual(initialArtifact.state.status, "running")
      assert.strictEqual(initialArtifact.state.pendingEvaluatorCalls.length, 2)

      yield* Deferred.succeed(gates[1]!, undefined)
      yield* Deferred.await(thirdStarted)
      assert.strictEqual(active, 2)

      const midRun = JSON.parse(
        yield* Effect.promise(() => readFile(join(runPath, "evaluation.json"), "utf8")),
      ) as { state: { workUnits: Record<string, { status: string }> } }
      assert.strictEqual(
        midRun.state.workUnits[`assessRateRequirement:${started[1]}`]?.status,
        "completed",
      )

      yield* Deferred.succeed(gates[0]!, undefined)
      yield* Deferred.succeed(gates[2]!, undefined)
      const result = yield* Fiber.join(fiber)
      assert.strictEqual(result.exitCode, 0)
      const receipt = JSON.parse(result.stdout) as { status: string }
      assert.strictEqual(receipt.status, "completed")
      const artifact = JSON.parse(
        yield* Effect.promise(() => readFile(join(runPath, "evaluation.json"), "utf8")),
      ) as {
        manifest: { concurrency: number }
        results: { payloads: ReadonlyArray<{ workUnit: string }> }
      }
      assert.strictEqual(artifact.manifest.concurrency, 2)
      assert.deepStrictEqual(
        [
          ...new Set(
            artifact.results.payloads
              .map((entry) => entry.workUnit)
              .filter((workUnit) => workUnit.startsWith("assessRateRequirement:")),
          ),
        ],
        [
          "assessRateRequirement:requirement:root::one",
          "assessRateRequirement:requirement:root::three",
          "assessRateRequirement:requirement:root::two",
        ],
      )
      assert.match(
        yield* Effect.promise(() => readFile(join(runPath, "logs", "events.jsonl"), "utf8")),
        /"peakActive":2/,
      )
    }),
  )

  it.effect("retries only the failed direct work unit", () =>
    Effect.gen(function* () {
      const directory = yield* Effect.acquireRelease(
        Effect.promise(async () => {
          await mkdir(join(repositoryRoot, "tmp"), { recursive: true })
          return mkdtemp(join(repositoryRoot, "tmp", "qualitymd-retry-test-"))
        }),
        (path) => Effect.promise(() => rm(path, { recursive: true })),
      )
      const model = join(directory, "QUALITY.md")
      yield* Effect.promise(() => writeFile(join(directory, "evidence.txt"), "bounded evidence"))
      yield* Effect.promise(() =>
        writeFile(
          model,
          `---
title: Retry test
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
      retry:
        title: Retry once
        assessment: Inspect the selected source.
      sibling:
        title: Sibling succeeds
        assessment: Inspect the selected source.
---
`,
        ),
      )
      const attempts = new Map<string, number>()
      const evaluator: EvaluatorService = {
        name: "retry-provider",
        kind: "codex",
        capabilities: {
          ...capabilities,
          dispatch: {
            concurrentCalls: true,
            delegatedRequests: false,
            automaticConcurrency: 2,
          },
        },
        evaluate: (request) => {
          const count = (attempts.get(request.workUnitId) ?? 0) + 1
          attempts.set(request.workUnitId, count)
          if (request.subject.endsWith("::retry") && count === 1)
            return Effect.fail(
              new EvaluatorFailure({ category: "rate_limited", detail: "retry this unit" }),
            )
          return Effect.succeed({
            workUnitId: request.workUnitId,
            evaluatorKind: "codex",
            payload: payloadFor(request),
          })
        },
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
      const receipt = JSON.parse(result.stdout) as { path: string; status: string }
      assert.strictEqual(receipt.status, "completed")
      assert.strictEqual(attempts.get("assessRateRequirement:requirement:root::retry"), 2)
      assert.strictEqual(attempts.get("assessRateRequirement:requirement:root::sibling"), 1)
      const artifact = JSON.parse(
        yield* Effect.promise(() =>
          readFile(join(directory, receipt.path, "evaluation.json"), "utf8"),
        ),
      ) as { state: { workUnits: Record<string, { attempts?: number }> } }
      assert.strictEqual(
        artifact.state.workUnits["assessRateRequirement:requirement:root::retry"]?.attempts,
        2,
      )
      assert.strictEqual(
        artifact.state.workUnits["assessRateRequirement:requirement:root::sibling"]?.attempts,
        1,
      )
    }),
  )

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
        capabilities: {
          ...capabilities,
          dispatch: {
            concurrentCalls: false,
            delegatedRequests: false,
            automaticConcurrency: 1,
            maxConcurrency: 1,
          },
        },
        evaluate: (request) =>
          Effect.succeed({
            workUnitId: request.workUnitId,
            evaluatorKind: "claude",
            model: "mock-model",
            payload: payloadFor(request),
            usage: {
              inputTokens: 10,
              outputTokens: 5,
              cachedInputTokens: 4,
              cacheWriteInputTokens: 3,
            },
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
      assert.strictEqual(artifact.schemaVersion, 9)
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
      const calls = (yield* Effect.promise(() =>
        readFile(join(runPath, "logs", "evaluator-calls.jsonl"), "utf8"),
      ))
        .trim()
        .split("\n")
        .map((line) => JSON.parse(line) as { usage?: Record<string, number> })
      assert.ok(calls.length > 0)
      for (const call of calls)
        assert.deepStrictEqual(call.usage, {
          inputTokens: 10,
          outputTokens: 5,
          cachedInputTokens: 4,
          cacheWriteInputTokens: 3,
        })
      const persistedArtifact = JSON.stringify(artifact)
      assert.notMatch(persistedArtifact, /cachedInputTokens/)
      assert.notMatch(persistedArtifact, /cacheWriteInputTokens/)
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
      accepted:
        title: Accepted evidence exists
        assessment: Inspect the selected source.
      blocked:
        title: Blocked evidence exists
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
        capabilities: {
          ...capabilities,
          dispatch: {
            concurrentCalls: false,
            delegatedRequests: false,
            automaticConcurrency: 1,
            maxConcurrency: 1,
          },
        },
        evaluate: (request) =>
          request.subject.endsWith("::accepted")
            ? Effect.succeed({
                workUnitId: request.workUnitId,
                evaluatorKind: "claude",
                payload: payloadFor(request),
              })
            : Deferred.succeed(started, undefined).pipe(
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
      ) as {
        manifest: {
          evaluator: string
          evaluatorKind: string
          concurrency: number
          evaluatorCapabilities: { dispatch: { maxConcurrency?: number } }
        }
        state: {
          status: string
          workUnits: Record<string, { status: string }>
          pendingEvaluatorCalls?: ReadonlyArray<unknown>
        }
      }
      assert.strictEqual(artifact.manifest.evaluator, evaluator.name)
      assert.strictEqual(artifact.manifest.evaluatorKind, "claude")
      assert.strictEqual(artifact.manifest.concurrency, 1)
      assert.strictEqual(artifact.manifest.evaluatorCapabilities.dispatch.maxConcurrency, 1)
      assert.strictEqual(artifact.state.status, "cancelled")
      assert.strictEqual(artifact.state.pendingEvaluatorCalls?.length, 1)
      assert.strictEqual(
        artifact.state.workUnits["assessRateRequirement:requirement:root::accepted"]?.status,
        "completed",
      )
      const recoveredSubjects: Array<string> = []
      const recoveryEvaluator: EvaluatorService = {
        ...evaluator,
        evaluate: (request) => {
          recoveredSubjects.push(request.subject)
          return Effect.succeed({
            workUnitId: request.workUnitId,
            evaluatorKind: "claude",
            payload: payloadFor(request),
          })
        },
      }
      const runPath = join(evaluationRoot, runs[0]!)
      const workspace = yield* resolveWorkspace({ model }).pipe(Effect.provide(services))
      const resumed = yield* resumeProviderRun(
        {
          model,
          evaluationDir: "",
          area: "",
          factors: [],
          evaluator: evaluator.name,
          resume: runPath,
          evaluatorResult: "",
          dryRun: false,
          json: true,
        },
        recoveryEvaluator,
        workspace,
      ).pipe(Effect.provide(services))
      assert.strictEqual(resumed.exitCode, 0)
      assert.strictEqual((JSON.parse(resumed.stdout) as { status: string }).status, "completed")
      assert.ok(!recoveredSubjects.includes("requirement:root::accepted"))
    }),
  )
})
