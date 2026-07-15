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
})
