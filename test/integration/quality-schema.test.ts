import { readFile } from "node:fs/promises"
import { describe, expect, it } from "vitest"

import { generateQualitySchema } from "../../src/domain/model/json-schema.ts"

describe("companion JSON Schema", () => {
  it("matches the TypeScript structural schema generator", async () => {
    const committed = await readFile(new URL("../../quality.schema.json", import.meta.url), "utf8")
    expect(committed).toBe(generateQualitySchema())
  })
})
