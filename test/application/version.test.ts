import { assert, describe, it } from "@effect/vitest"
import * as Effect from "effect/Effect"

import { currentVersionInfo, versionCommand } from "../../src/application/version.ts"
import { specificationVersion } from "../../src/build-info.ts"

describe("version", () => {
  it("reports the bundled specification version", () => {
    assert.strictEqual(specificationVersion, "0.12 (Draft)")
    assert.strictEqual(currentVersionInfo().specificationVersion, specificationVersion)
  })

  it.effect("emits stable JSON", () =>
    Effect.gen(function* () {
      const result = yield* versionCommand({ json: true })
      assert.deepStrictEqual(JSON.parse(result.stdout), currentVersionInfo())
      assert.strictEqual(result.exitCode, 0)
    }),
  )
})
