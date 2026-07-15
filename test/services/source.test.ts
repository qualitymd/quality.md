import { assert, describe, it } from "@effect/vitest"
import * as BunFileSystem from "@effect/platform-bun/BunFileSystem"
import * as BunPath from "@effect/platform-bun/BunPath"
import * as Effect from "effect/Effect"
import * as Layer from "effect/Layer"
import { afterEach } from "vitest"
import { mkdir, mkdtemp, rm, symlink, writeFile } from "node:fs/promises"
import { tmpdir } from "node:os"
import { join } from "node:path"

import {
  detectSourceKind,
  sealEvidenceManifest,
  validateSourceSelector,
} from "../../src/services/source.ts"

const temporaryDirectories: string[] = []

const temporaryDirectory = async () => {
  const directory = await mkdtemp(join(tmpdir(), "qualitymd-evidence-test-"))
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

const assessment = (sourceRef: string) => ({
  findings: [{ evidence: [{ sourceRef, statement: "Observed evidence." }] }],
})

describe("source selection and evidence sealing", () => {
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

  it.effect("seals requirement-specific evaluated and supporting evidence", () =>
    Effect.gen(function* () {
      const root = yield* promise(temporaryDirectory)
      yield* promise(() => mkdir(join(root, "src")))
      yield* promise(() => mkdir(join(root, "docs")))
      yield* promise(() =>
        writeFile(join(root, "src", "target.ts"), "# Target\n\nexport const ready = true\n"),
      )
      yield* promise(() => writeFile(join(root, "docs", "design.md"), "# Design\n\nReady.\n"))

      const manifest = yield* sealEvidenceManifest({
        workspaceRoot: root,
        requirementId: "requirement:root::ready",
        source: { selector: "src", kind: "path" },
        proposal: {
          observations: [
            {
              id: "ev-001",
              kind: "file",
              role: "evaluated",
              path: "src/target.ts",
              locator: { startLine: 3, endLine: 3 },
            },
            {
              id: "ev-002",
              kind: "file",
              role: "supporting",
              path: "docs/design.md",
              locator: { heading: "Design" },
            },
          ],
          limits: [],
        },
        assessment: assessment("evidence[ev-001]"),
        capturedAt: "2026-07-14T00:00:00Z",
      })

      assert.strictEqual(manifest.observations.length, 2)
      assert.strictEqual(manifest.observations[0]?.path, "src/target.ts")
      assert.match(String(manifest.observations[0]?.sha256), /^[a-f0-9]{64}$/)
      assert.match(manifest.manifestHash, /^[a-f0-9]{64}$/)
      assert.strictEqual("content" in (manifest.observations[0] ?? {}), false)
    }).pipe(Effect.provide(services)),
  )

  it.effect("rejects escaping selectors and symlink escapes", () =>
    Effect.gen(function* () {
      const root = yield* promise(temporaryDirectory)
      const outside = yield* promise(temporaryDirectory)
      yield* promise(() => writeFile(join(outside, "secret.txt"), "secret"))
      yield* promise(() => symlink(join(outside, "secret.txt"), join(root, "secret-link")))

      const selector = yield* validateSourceSelector(root, {
        selector: "../outside",
        kind: "path",
      }).pipe(Effect.result)
      assert.strictEqual(selector._tag, "Failure")

      const sealed = yield* sealEvidenceManifest({
        workspaceRoot: root,
        requirementId: "requirement:root::safe",
        source: { selector: "the selected subject", kind: "prose" },
        proposal: {
          observations: [{ id: "ev-001", kind: "file", role: "evaluated", path: "secret-link" }],
          limits: [],
        },
        assessment: assessment("evidence[ev-001]"),
        capturedAt: "2026-07-14T00:00:00Z",
      }).pipe(Effect.result)
      assert.strictEqual(sealed._tag, "Failure")
      if (sealed._tag === "Failure") assert.match(sealed.failure.message, /escapes the workspace/)
    }).pipe(Effect.provide(services)),
  )

  it.effect("rejects misclassified target evidence and unbound finding references", () =>
    Effect.gen(function* () {
      const root = yield* promise(temporaryDirectory)
      yield* promise(() => mkdir(join(root, "src")))
      yield* promise(() => mkdir(join(root, "docs")))
      yield* promise(() => writeFile(join(root, "src", "target.ts"), "target"))
      yield* promise(() => writeFile(join(root, "docs", "design.md"), "# Design\n"))

      const outsideTarget = yield* sealEvidenceManifest({
        workspaceRoot: root,
        requirementId: "requirement:root::ready",
        source: { selector: "src", kind: "path" },
        proposal: {
          observations: [{ id: "ev-001", kind: "file", role: "evaluated", path: "docs/design.md" }],
          limits: [],
        },
        assessment: assessment("evidence[ev-001]"),
        capturedAt: "2026-07-14T00:00:00Z",
      }).pipe(Effect.result)
      assert.strictEqual(outsideTarget._tag, "Failure")
      if (outsideTarget._tag === "Failure")
        assert.match(outsideTarget.failure.message, /outside the evaluated source selector/)

      const unbound = yield* sealEvidenceManifest({
        workspaceRoot: root,
        requirementId: "requirement:root::ready",
        source: { selector: "src", kind: "path" },
        proposal: {
          observations: [{ id: "ev-001", kind: "file", role: "evaluated", path: "src/target.ts" }],
          limits: [],
        },
        assessment: assessment("evidence[ev-missing]"),
        capturedAt: "2026-07-14T00:00:00Z",
      }).pipe(Effect.result)
      assert.strictEqual(unbound._tag, "Failure")
      if (unbound._tag === "Failure")
        assert.match(unbound.failure.message, /does not name an accepted evidence observation/)
    }).pipe(Effect.provide(services)),
  )
})
