import { describe, expect, it } from "vitest"

import {
  harnessRunArtifact,
  harnessRunEvents,
  harnessRunReceipt,
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
      harnessRunArtifact({
        identity,
        model: "QUALITY.md",
        scope: {
          requestedScope: {},
          plannedScope: { areaId: "area:root", factorFilter: [] },
        },
        number: 1,
        label: "0001-full-eval",
        capabilities: { concurrent: true },
        concurrency: 2,
        areaSources: { "area:root": { selector: ".", kind: "path" } },
        workUnits: { frameEvaluation: { status: "completed" } },
        pending: [pending],
        payloads: [{ workUnit: "frameEvaluation", payload: { kind: "EvaluationFrame" } }],
      }),
    ).toMatchObject({
      schemaVersion: 8,
      kind: "EvaluationRun",
      manifest: { ...identity, run: { number: 1, label: "0001-full-eval" } },
      state: { status: "awaiting_evaluator", pendingEvaluatorCalls: [pending] },
    })
    expect(harnessRunEvents(identity.createdAt, "eval", { concurrent: true }, 1)).toBe(
      '{"timestamp":"2026-07-14T00:00:00Z","event":"run_created","evaluationId":"eval","evaluator":"harness","evaluatorKind":"harness","capabilities":{"concurrent":true}}\n' +
        '{"timestamp":"2026-07-14T00:00:00Z","event":"run_status","status":"awaiting_evaluator","outstanding":1}\n',
    )
    const request = { requestId: "req_1", workUnitId: pending.workUnitId, attempt: 1 }
    expect(
      harnessRunReceipt({
        path: ".quality/evaluations/0001-full-eval",
        concurrency: 2,
        total: 4,
        evaluatorUnits: 1,
        completed: 3,
        sources: [],
        requests: [request],
      }),
    ).toMatchObject({ status: "awaiting_evaluator", evaluatorRequests: [request] })
    expect(renderHarnessAwaiting([request], 2, "qualitymd evaluation run --resume run")).toContain(
      "req_1",
    )
  })
})
