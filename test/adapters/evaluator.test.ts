import { describe, expect, it } from "vitest"

import { strictProviderSchema, stripNullProperties } from "../../src/adapters/evaluator.ts"

describe("provider schema adaptation", () => {
  it("requires every property while making optional fields nullable", () => {
    const source = {
      $schema: "https://json-schema.org/draft/2020-12/schema",
      type: "object",
      properties: {
        name: { type: "string" },
        nested: {
          type: "object",
          properties: { description: { type: "string" }, enabled: { type: "boolean" } },
          required: ["enabled"],
        },
      },
      required: ["name"],
      allOf: [{ if: {}, ["th" + "en"]: {} }],
    }
    const strict = strictProviderSchema(source) as Record<string, unknown>
    const properties = strict.properties as Record<string, unknown>
    const nested = properties.nested as Record<string, unknown>
    const nestedVariants = nested.anyOf as Array<Record<string, unknown>>
    expect(strict.required).toEqual(["name", "nested"])
    expect(strict.additionalProperties).toBe(false)
    expect(strict).not.toHaveProperty("allOf")
    expect(strict).not.toHaveProperty("$schema")
    expect(properties.name).toEqual({ type: "string" })
    expect(nestedVariants[1]).toEqual({ type: "null" })
    expect(nestedVariants[0]).toMatchObject({
      additionalProperties: false,
      required: ["description", "enabled"],
      properties: {
        description: { anyOf: [{ type: "string" }, { type: "null" }] },
        enabled: { type: "boolean" },
      },
    })
    expect(source.required).toEqual(["name"])
    expect(source).toHaveProperty("allOf")
  })

  it("removes provider null placeholders before persisted validation", () => {
    expect(
      stripNullProperties({ name: "value", optional: null, nested: [{ keep: true, omit: null }] }),
    ).toEqual({ name: "value", nested: [{ keep: true }] })
  })
})
