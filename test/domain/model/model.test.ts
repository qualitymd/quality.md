import { describe, expect, it } from "vitest"

import {
  effectiveSource,
  findElement,
  flattenModel,
  parseAreaReference,
  projectModel,
  truncateDepth,
  type QualityModel,
} from "../../../src/domain/model/model.ts"

const model: QualityModel = {
  path: "QUALITY.md",
  title: "Example",
  source: "root/**",
  ratingScale: [{ level: "target", criterion: "Meets the target." }],
  factors: {
    reliability: {
      title: "Reliability",
      factors: { recovery: { title: "Recovery" } },
      requirements: {
        restarts: { title: "Restarts", assessment: "Inspect recovery behavior." },
      },
    },
  },
  areas: {
    api: {
      title: "API",
      source: "api/**",
      areas: { routes: { title: "Routes" } },
      factors: { usability: { title: "Usability" } },
    },
  },
}

describe("quality model projection", () => {
  it("builds stable canonical IDs in deterministic order", () => {
    expect(flattenModel(projectModel(model)).map((entry) => entry.id)).toEqual([
      "area:root",
      "factor:root::reliability",
      "factor:root::reliability/recovery",
      "requirement:root::restarts",
      "area:api",
      "factor:api::usability",
      "area:api/routes",
    ])
  })

  it("finds elements and truncates without mutating the source tree", () => {
    const root = projectModel(model)
    expect(findElement(root, "factor:api::usability")?.label).toBe("Usability")
    expect(truncateDepth(root, 0)).not.toHaveProperty("children")
    expect(root.children).toHaveLength(2)
  })

  it("resolves valid areas and rejects invalid or unresolved references", () => {
    expect(parseAreaReference(model, "area:api/routes")).toEqual(["api", "routes"])
    expect(() => parseAreaReference(model, "api")).toThrow(/must start with area:/)
    expect(() => parseAreaReference(model, "area:missing")).toThrow(/does not resolve/)
  })

  it("distinguishes declared, inherited, and default sources", () => {
    expect(effectiveSource(model, [])).toEqual({ selector: "root/**", state: "declared" })
    expect(effectiveSource(model, ["api"])).toEqual({ selector: "api/**", state: "declared" })
    expect(effectiveSource(model, ["api", "routes"])).toEqual({
      selector: "api/**",
      state: "inherited",
    })
    const { source: _source, ...withoutSource } = model
    expect(effectiveSource(withoutSource, [])).toEqual({
      selector: ".",
      state: "default",
    })
  })
})
