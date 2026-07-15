import { describe, expect, it } from "vitest"

import { buildGraph, readyUnits } from "../../../src/domain/evaluation/graph.ts"
import type { EvaluationPlan } from "../../../src/domain/evaluation/plan.ts"

const plan: EvaluationPlan = {
  areas: [
    {
      ref: "area:root",
      path: [],
      value: {} as EvaluationPlan["areas"][number]["value"],
      source: "the maintained entity",
      childAreaIds: [],
      rootFactorIds: ["factor:root::reliability"],
      localRequirementIds: [],
    },
  ],
  factors: [
    {
      ref: "factor:root::reliability",
      areaId: "area:root",
      path: ["reliability"],
      value: {} as EvaluationPlan["factors"][number]["value"],
    },
  ],
  requirements: [
    {
      ref: "requirement:root::recovers",
      areaId: "area:root",
      factorIds: ["factor:root::reliability"],
      value: {} as EvaluationPlan["requirements"][number]["value"],
    },
  ],
}

describe("evaluation graph", () => {
  it("starts requirement judgment from its frame and reports last", () => {
    const graph = buildGraph(plan)
    const requirement = graph.find((unit) => unit.kind === "assessRateRequirement")!
    const summary = graph.find((unit) => unit.kind === "summarizeEvaluation")!
    const report = graph.at(-1)!

    expect(requirement.dependsOn).toEqual(["frameRequirementEvaluation:requirement:root::recovers"])
    expect(graph.map((unit) => unit.kind)).not.toContain("resolveSource")
    expect(summary).toMatchObject({
      dataKind: "EvaluationSummaryResult",
      dependsOn: ["analyzeArea:area:root", "rankFindings", "rankRecommendations"],
      evaluatorBacked: true,
    })
    expect(report.kind).toBe("buildReports")
    expect(report.dependsOn).toEqual(graph.slice(0, -1).map((unit) => unit.id))
    expect(new Set(graph.map((unit) => unit.id)).size).toBe(graph.length)
  })

  it("selects evaluator work whose dependencies are complete without mutating the graph", () => {
    const graph = buildGraph(plan)
    const before = structuredClone(graph)

    expect(
      readyUnits(
        graph,
        new Set([
          "frameEvaluation",
          "frameAreaEvaluation:area:root",
          "frameRequirementEvaluation:requirement:root::recovers",
        ]),
        1,
      ).map((unit) => unit.id),
    ).toEqual(["assessRateRequirement:requirement:root::recovers"])
    expect(graph).toEqual(before)
  })
})
