import { describe, expect, it } from "vitest"

import {
  dataKinds,
  evaluationDataExample,
  validateDataPayload,
} from "../../../src/domain/evaluation/data.ts"

describe("evaluation data", () => {
  it("keeps every bundled example valid against schema and semantic rules", () => {
    for (const [kind, writable] of dataKinds) {
      const example = structuredClone(evaluationDataExample(kind))
      expect(validateDataPayload(example, undefined, !writable)).toBe(kind)
    }
  })

  it("rejects severity on a strength finding", () => {
    const example = structuredClone(evaluationDataExample("RequirementAssessmentResult"))
    const finding = (example.findings as Array<Record<string, unknown>>)[0]!
    finding.severity = "minor"
    expect(() => validateDataPayload(example)).toThrow(/severity|must NOT be valid/)
  })

  it("rejects planning fields on recommendations", () => {
    const example = structuredClone(evaluationDataExample("RecommendationResult"))
    example.priority = "high"
    expect(() => validateDataPayload(example)).toThrow(/priority|additional properties/)
  })

  it("requires three to five executive-summary key points", () => {
    const example = structuredClone(evaluationDataExample("EvaluationSummaryResult"))
    example.keyPoints = ["Only one"]
    expect(() => validateDataPayload(example)).toThrow(/keyPoints|fewer than 3/)
  })
})
