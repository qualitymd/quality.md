import { describe, expect, it } from "vitest"

import { initialFramePayloads } from "../../../src/domain/evaluation/frames.ts"
import { planEvaluation } from "../../../src/domain/evaluation/plan.ts"
import type { QualityModel } from "../../../src/domain/model/model.ts"

const model: QualityModel = {
  path: "QUALITY.md",
  title: "Example",
  source: "src/**",
  ratingScale: [
    { level: "partial", criterion: "Partially meets the requirement." },
    { level: "target", criterion: "Meets the requirement." },
  ],
  factors: {
    reliability: {
      title: "Reliability",
      requirements: {
        recovery: {
          title: "Recovery",
          assessment: "Inspect recovery behavior.",
          ratings: { target: "Recovers within the target." },
        },
      },
    },
  },
}

describe("initial evaluation frames", () => {
  it("derives stable ordered frame payloads from the model and plan", () => {
    const plan = planEvaluation(model, { areaId: "area:root", factorFilter: [] })

    expect(initialFramePayloads(model, plan, "QUALITY.md")).toEqual([
      {
        workUnit: "frameEvaluation",
        payload: {
          derivedContext: {
            evaluationPolicies: ["source-as-data", "secret-redaction"],
            rigor: "standard",
          },
          inputs: { ratingLevelIds: ["rating:partial", "rating:target"] },
          kind: "EvaluationFrame",
          schemaVersion: 3,
          subject: { modelLocator: "QUALITY.md" },
        },
      },
      {
        workUnit: "frameAreaEvaluation:area:root",
        payload: {
          inputs: {
            childAreaIds: [],
            localRequirementIds: ["requirement:root::recovery"],
            rootFactorIds: ["factor:root::reliability"],
            sourceRefs: ["src/**"],
          },
          kind: "AreaEvaluationFrame",
          schemaVersion: 3,
          subject: { areaId: "area:root" },
        },
      },
      {
        workUnit: "frameRequirementEvaluation:requirement:root::recovery",
        payload: {
          derivedContext: {
            appliedRatingCriteria: [
              {
                criterion: "Partially meets the requirement.",
                ratingLevelId: "rating:partial",
                source: "model_default",
              },
              {
                criterion: "Recovers within the target.",
                ratingLevelId: "rating:target",
                source: "requirement_override",
              },
            ],
          },
          inputs: {
            ratingLevelIds: ["rating:partial", "rating:target"],
            requirementAssessmentBasis: "Inspect recovery behavior.",
          },
          kind: "RequirementEvaluationFrame",
          schemaVersion: 3,
          subject: {
            factorIds: ["factor:root::reliability"],
            requirementId: "requirement:root::recovery",
          },
        },
      },
    ])
  })
})
