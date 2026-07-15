import { describe, expect, it } from "vitest"

import {
  resolveConcurrency,
  validateDispatchCapability,
} from "../../../src/domain/evaluation/concurrency.ts"

const harness = {
  concurrentCalls: false,
  delegatedRequests: true,
  automaticConcurrency: 4,
} as const

const codex = {
  concurrentCalls: true,
  delegatedRequests: false,
  automaticConcurrency: 4,
} as const

const claude = {
  concurrentCalls: false,
  delegatedRequests: false,
  automaticConcurrency: 1,
  maxConcurrency: 1,
} as const

describe("evaluation concurrency resolution", () => {
  it.each([
    ["harness automatic", undefined, harness, "delegated", 4, false],
    ["harness configured", 2, harness, "delegated", 2, false],
    ["codex automatic", undefined, codex, "direct", 4, false],
    ["codex configured", 7, codex, "direct", 7, false],
    ["claude automatic", undefined, claude, "sequential", 1, false],
    ["claude configured clamp", 8, claude, "sequential", 1, true],
  ] as const)("resolves %s", (_label, configured, dispatch, mode, resolved, clamped) => {
    expect(resolveConcurrency(configured, dispatch)).toMatchObject({ mode, resolved, clamped })
  })

  it("records configured and automatic resolution sources", () => {
    expect(resolveConcurrency(undefined, codex)).toMatchObject({
      source: "automatic",
      requested: 4,
      automatic: 4,
    })
    expect(resolveConcurrency(2, codex)).toMatchObject({
      source: "configured",
      requested: 2,
      automatic: 4,
    })
  })

  it.each([0, -1, 1.5, "4", null])("rejects invalid configured cap %j", (configured) => {
    expect(() => resolveConcurrency(configured, codex)).toThrow(
      "evaluation.concurrency must be a positive integer",
    )
  })

  it("validates maximums and the sequential declaration invariant", () => {
    expect(() =>
      validateDispatchCapability({ ...codex, automaticConcurrency: 5, maxConcurrency: 4 }),
    ).toThrow("must not exceed")
    expect(() =>
      validateDispatchCapability({
        concurrentCalls: false,
        delegatedRequests: false,
        automaticConcurrency: 2,
      }),
    ).toThrow("sequential evaluator")
  })
})
