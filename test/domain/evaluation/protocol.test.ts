import { describe, expect, it } from "vitest"

import {
  buildProtocolRequest,
  completeProtocolRequest,
  protocolRequestReceipt,
} from "../../../src/domain/evaluation/protocol.ts"
import type { WorkUnit } from "../../../src/domain/evaluation/graph.ts"
import type { EvaluationPlan } from "../../../src/domain/evaluation/plan.ts"

const unit: WorkUnit = {
  id: "assessRateRequirement:requirement:root::ready",
  kind: "assessRateRequirement",
  subject: "requirement:root::ready",
  dependsOn: ["frameRequirementEvaluation:requirement:root::ready"],
  evaluatorBacked: true,
}

const plan: EvaluationPlan = {
  areas: [
    {
      ref: "area:root",
      path: [],
      value: {} as EvaluationPlan["areas"][number]["value"],
      source: ".",
      childAreaIds: [],
      rootFactorIds: ["factor:root::reliability"],
      localRequirementIds: ["requirement:root::ready"],
    },
  ],
  factors: [],
  requirements: [
    {
      ref: "requirement:root::ready",
      areaId: "area:root",
      factorIds: ["factor:root::reliability"],
      value: { title: "Ready", assessment: "Inspect readiness." },
    },
  ],
}

describe("evaluation protocol assembly", () => {
  it("keeps optional request fields absent while retaining hash sentinels", () => {
    const draft = buildProtocolRequest({
      unit,
      plan,
      payloads: [
        {
          workUnit: "frameAreaEvaluation:area:root",
          payload: { kind: "AreaEvaluationFrame" },
        },
        {
          workUnit: "frameRequirementEvaluation:requirement:root::ready",
          payload: { kind: "RequirementEvaluationFrame" },
        },
      ],
      areaSources: { "area:root": { selector: ".", kind: "path" } },
      bodyGuidance: "",
      ratingScale: [{ level: "target", title: "Target", criterion: "Meets the target." }],
      evaluationId: "eval",
    })

    expect(draft).not.toHaveProperty("bodyGuidance")
    expect(draft.hashInput).toMatchObject({ bodyGuidance: "" })
    const protocol = completeProtocolRequest(draft, "input-hash")
    const receipt = protocolRequestReceipt(protocol, {
      requestId: "req_1",
      workUnitId: unit.id,
      inputHash: "input-hash",
      correlationId: protocol.correlationId,
      attempt: 1,
    })
    expect(receipt).not.toHaveProperty("bodyGuidance")
    expect(receipt).toMatchObject({
      requestId: "req_1",
      workUnitId: unit.id,
      subject: unit.subject,
      inputHash: "input-hash",
      correlationId: `eval#${unit.id}`,
    })
  })

  it("assembles an inputs-only stakeholder summary request", () => {
    const summaryUnit: WorkUnit = {
      id: "summarizeEvaluation",
      kind: "summarizeEvaluation",
      subject: "",
      dependsOn: ["analyzeArea:area:root", "rankFindings", "rankRecommendations"],
      evaluatorBacked: true,
      dataKind: "EvaluationSummaryResult",
    }
    const draft = buildProtocolRequest({
      unit: summaryUnit,
      plan,
      payloads: [
        { workUnit: "frameEvaluation", payload: { kind: "EvaluationFrame" } },
        { workUnit: "analyzeArea:area:root", payload: { kind: "AreaAnalysisResult" } },
        { workUnit: "rankFindings", payload: { kind: "FindingRankingResult" } },
        { workUnit: "recommend", payload: { kind: "RecommendationResult" } },
        {
          workUnit: "rankRecommendations",
          payload: { kind: "RecommendationRankingResult" },
        },
      ],
      areaSources: { "area:root": { selector: ".", kind: "path" } },
      bodyGuidance: "",
      ratingScale: [{ level: "target", title: "Target", criterion: "Meets the target." }],
      evaluationId: "eval",
    })

    expect(draft).not.toHaveProperty("inspection")
    expect(draft.instructions).toContain("stakeholder-facing executive summary")
    expect(draft.context).toMatchObject({
      overallAnalysis: { kind: "AreaAnalysisResult" },
      evaluationFrame: { kind: "EvaluationFrame" },
      ratingScale: [{ level: "target", title: "Target" }],
      findingRanking: { kind: "FindingRankingResult" },
      recommendations: [{ kind: "RecommendationResult" }],
      recommendationRanking: { kind: "RecommendationRankingResult" },
    })
    expect(draft.expectedSchema).toMatchObject({
      properties: {
        kind: { const: "EvaluationSummaryResult" },
        keyPoints: { minItems: 3, maxItems: 5 },
      },
    })
  })
})
