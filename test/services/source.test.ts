import { assert, describe, it } from "@effect/vitest"
import * as BunFileSystem from "@effect/platform-bun/BunFileSystem"
import * as BunPath from "@effect/platform-bun/BunPath"
import * as Effect from "effect/Effect"
import * as Layer from "effect/Layer"
import { mkdtemp, mkdir, rm, writeFile } from "node:fs/promises"
import { join } from "node:path"
import { tmpdir } from "node:os"
import { afterEach } from "vitest"

import { captureSource, detectSourceKind, packageSource } from "../../src/services/source.ts"

const temporaryDirectories: string[] = []

const temporaryDirectory = async () => {
  const directory = await mkdtemp(join(tmpdir(), "qualitymd-source-test-"))
  temporaryDirectories.push(directory)
  return directory
}

const services = Layer.merge(BunFileSystem.layer, BunPath.layer)
const promise = <A>(evaluate: () => Promise<A>) => Effect.promise(evaluate)

afterEach(async () => {
  await Promise.all(
    temporaryDirectories.splice(0).map((directory) => rm(directory, { recursive: true })),
  )
})

describe("source selection and capture", () => {
  it.effect("distinguishes paths, valid globs, malformed globs, and prose", () =>
    Effect.gen(function* () {
      const root = yield* promise(temporaryDirectory)
      yield* promise(() => writeFile(join(root, "README.md"), "hello"))

      assert.strictEqual(yield* detectSourceKind(root, "README.md"), "path")
      assert.strictEqual(yield* detectSourceKind(root, "**/*.md"), "glob")
      assert.strictEqual(yield* detectSourceKind(root, "[unterminated"), "prose")
      assert.strictEqual(yield* detectSourceKind(root, "the current release process"), "prose")
    }).pipe(Effect.provide(services)),
  )

  it.effect("packages deterministic text files while skipping default vendor trees", () =>
    Effect.gen(function* () {
      const root = yield* promise(temporaryDirectory)
      yield* promise(() => mkdir(join(root, "src")))
      yield* promise(() => mkdir(join(root, "vendor")))
      yield* promise(() => writeFile(join(root, "src", "a.txt"), "alpha"))
      yield* promise(() => writeFile(join(root, "vendor", "b.txt"), "beta"))

      const rootBundle = yield* packageSource(root, ".", "path")
      assert.deepStrictEqual(
        rootBundle.files.map((file) => file.path),
        ["src/a.txt"],
      )
      const vendorBundle = yield* packageSource(root, "vendor", "path")
      assert.deepStrictEqual(
        vendorBundle.files.map((file) => file.path),
        ["vendor/b.txt"],
      )
    }).pipe(Effect.provide(services)),
  )

  it.effect("rejects escaping and empty selectors", () =>
    Effect.gen(function* () {
      const root = yield* promise(temporaryDirectory)
      const escaping = yield* packageSource(root, "../outside", "path").pipe(Effect.result)
      assert.strictEqual(escaping._tag, "Failure")
      if (escaping._tag === "Failure")
        assert.match(escaping.failure.detail, /escapes the workspace/)
      const missing = yield* packageSource(root, "missing", "path").pipe(Effect.result)
      assert.strictEqual(missing._tag, "Failure")
      if (missing._tag === "Failure")
        assert.match(missing.failure.detail, /contains no readable files/)
    }).pipe(Effect.provide(services)),
  )

  it.effect("validates resolved source containment and applies file bounds", () =>
    Effect.gen(function* () {
      const root = yield* promise(temporaryDirectory)
      yield* promise(() => writeFile(join(root, "large.txt"), "x".repeat(70 * 1024)))
      const escaping = yield* captureSource(root, { files: [{ path: "../secret" }] }).pipe(
        Effect.result,
      )
      assert.strictEqual(escaping._tag, "Failure")
      if (escaping._tag === "Failure")
        assert.match(escaping.failure.message, /workspace-relative and contained/)
      const injected = yield* captureSource(root, {
        files: [{ path: "large.txt", content: "injected" }],
      }).pipe(Effect.result)
      assert.strictEqual(injected._tag, "Failure")
      if (injected._tag === "Failure")
        assert.match(injected.failure.message, /must carry only path/)
      const bundle = yield* captureSource(root, { files: [{ path: "large.txt" }] })
      assert.strictEqual(bundle.files[0]?.truncated, true)
      assert.strictEqual(new TextEncoder().encode(bundle.files[0]?.content).length, 64 * 1024)
    }).pipe(Effect.provide(services)),
  )
})
