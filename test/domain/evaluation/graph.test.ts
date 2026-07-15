import { describe, expect, it } from "vitest"

import { buildGraph } from "../../../src/domain/evaluation/graph.ts"
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
  it("orders source resolution before requirement judgment and reports last", () => {
    const graph = buildGraph(plan, { "area:root": "prose" })
    const source = graph.find((unit) => unit.kind === "resolveSource")!
    const requirement = graph.find((unit) => unit.kind === "assessRateRequirement")!
    const report = graph.at(-1)!

    expect(requirement.dependsOn).toContain(source.id)
    expect(report.kind).toBe("buildReports")
    expect(report.dependsOn).toEqual(graph.slice(0, -1).map((unit) => unit.id))
    expect(new Set(graph.map((unit) => unit.id)).size).toBe(graph.length)
  })
})
