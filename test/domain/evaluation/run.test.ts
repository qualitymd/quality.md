import { describe, expect, it } from "vitest"

import {
  evaluationRunArtifact,
  evaluationRunEvents,
  evaluationRunReceipt,
  renderHarnessAwaiting,
  runDirectoryNumber,
} from "../../../src/domain/evaluation/run.ts"

describe("evaluation run directory classification", () => {
  it("uses the current artifact before the historical manifest and folder prefix", () => {
    expect(
      runDirectoryNumber({
        name: "0003-quality-eval",
        evaluationArtifact: JSON.stringify({ manifest: { run: { number: 11 } } }),
        historicalManifest: JSON.stringify({ run: { number: 7 } }),
      }),
    ).toBe(11)
  })

  it("falls through invalid and unreadable documents in precedence order", () => {
    expect(
      runDirectoryNumber({
        name: "0003-quality-eval",
        evaluationArtifact: JSON.stringify({ manifest: { run: { number: 0 } } }),
        historicalManifest: JSON.stringify({ run: { number: 7 } }),
      }),
    ).toBe(7)
    expect(
      runDirectoryNumber({
        name: "0003-quality-eval",
        evaluationArtifact: "not json",
        historicalManifest: JSON.stringify({ run: { number: -1 } }),
      }),
    ).toBe(3)
  })

  it("recognizes quality slugs and rejects non-positive or unrelated folder names", () => {
    expect(runDirectoryNumber({ name: "0006-quality-eval" })).toBe(6)
    expect(runDirectoryNumber({ name: "0000-quality-eval" })).toBeUndefined()
    expect(runDirectoryNumber({ name: "quality-eval" })).toBeUndefined()
  })

  it("assembles stable harness artifacts, events, receipts, and human summaries", () => {
    const pending = {
      requestId: "req_1",
      workUnitId: "assessRateRequirement:requirement:root::ready",
      inputHash: "hash",
      correlationId: "eval#unit",
      attempt: 1,
    }
    const identity = { evaluationId: "eval", createdAt: "2026-07-14T00:00:00Z" }
    expect(
      evaluationRunArtifact({
        identity,
        model: "QUALITY.md",
        scope: {
          requestedScope: {},
          plannedScope: { areaId: "area:root", factorFilter: [] },
        },
        number: 1,
        label: "0001-full-eval",
        evaluator: {
          name: "harness",
          kind: "harness",
          capabilities: {
            structuredOutput: true,
            workspaceInspection: true,
            instructionIsolation: true,
            verification: false,
            networkAccess: "disabled",
            tools: true,
            dispatch: {
              concurrentCalls: false,
              delegatedRequests: true,
              automaticConcurrency: 4,
            },
            freshContext: true,
            cancellation: true,
            usage: true,
            maxTurns: "supported",
            tokenBudget: "supported",
            costBudget: "supported",
            contextWindow: "unknown",
            compaction: "opaque",
            sandbox: "host",
            executableOverride: false,
          },
        },
        concurrency: {
          source: "configured",
          requested: 2,
          automatic: 4,
          resolved: 2,
          clamped: false,
          mode: "delegated",
        },
        areaSources: { "area:root": { selector: ".", kind: "path" } },
        workUnits: { frameEvaluation: { status: "completed" } },
        pending: [pending],
        payloads: [{ workUnit: "frameEvaluation", payload: { kind: "EvaluationFrame" } }],
      }),
    ).toMatchObject({
      schemaVersion: 9,
      kind: "EvaluationRun",
      manifest: { ...identity, run: { number: 1, label: "0001-full-eval" } },
      state: { status: "awaiting_evaluator", pendingEvaluatorCalls: [pending] },
    })
    expect(
      evaluationRunEvents(
        identity.createdAt,
        "eval",
        {
          name: "harness",
          kind: "harness",
          capabilities: {
            structuredOutput: true,
            workspaceInspection: true,
            instructionIsolation: true,
            verification: false,
            networkAccess: "disabled",
            tools: true,
            dispatch: {
              concurrentCalls: false,
              delegatedRequests: true,
              automaticConcurrency: 4,
            },
            freshContext: true,
            cancellation: true,
            usage: true,
            maxTurns: "supported",
            tokenBudget: "supported",
            costBudget: "supported",
            contextWindow: "unknown",
            compaction: "opaque",
            sandbox: "host",
            executableOverride: false,
          },
        },
        {
          source: "configured",
          requested: 2,
          automatic: 4,
          resolved: 2,
          clamped: false,
          mode: "delegated",
        },
        1,
      ),
    ).toContain('"peakOutstanding":1')
    const request = { requestId: "req_1", workUnitId: pending.workUnitId, attempt: 1 }
    expect(
      evaluationRunReceipt({
        path: ".quality/evaluations/0001-full-eval",
        evaluator: "harness",
        evaluatorKind: "harness",
        concurrency: 2,
        total: 4,
        evaluatorUnits: 1,
        completed: 3,
        sources: [],
        requests: [request],
        dispatchMode: "delegated",
      }),
    ).toMatchObject({ status: "awaiting_evaluator", evaluatorRequests: [request] })
    expect(renderHarnessAwaiting([request], 2, "qualitymd evaluation run --resume run")).toContain(
      "req_1",
    )
  })
})
