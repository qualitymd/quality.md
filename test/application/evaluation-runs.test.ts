import { assert, describe, it } from "@effect/vitest"
import * as BunFileSystem from "@effect/platform-bun/BunFileSystem"
import * as BunPath from "@effect/platform-bun/BunPath"
import * as Effect from "effect/Effect"
import * as Layer from "effect/Layer"
import { afterEach } from "vitest"
import { mkdir, mkdtemp, rm, writeFile } from "node:fs/promises"
import { tmpdir } from "node:os"
import { join } from "node:path"

import {
  evaluationRunDirectories,
  nextEvaluationRunNumber,
} from "../../src/application/evaluation-runs.ts"

const temporaryDirectories: string[] = []
const services = Layer.merge(BunFileSystem.layer, BunPath.layer)
const promise = <A>(evaluate: () => Promise<A>) => Effect.promise(evaluate)

afterEach(async () => {
  await Promise.all(
    temporaryDirectories.splice(0).map((directory) => rm(directory, { recursive: true })),
  )
})

describe("evaluation run enumeration", () => {
  it.effect("applies shared precedence and excludes non-directories", () =>
    Effect.gen(function* () {
      const root = yield* promise(async () => {
        const directory = await mkdtemp(join(tmpdir(), "qualitymd-runs-test-"))
        temporaryDirectories.push(directory)
        return directory
      })
      const current = join(root, "0001-current-eval")
      const historical = join(root, "0002-historical-eval")
      const unreadable = join(root, "0003-unreadable-eval")
      const quality = join(root, "0004-quality-eval")
      yield* promise(() => mkdir(current))
      yield* promise(() =>
        writeFile(
          join(current, "evaluation.json"),
          JSON.stringify({ manifest: { run: { number: 9 } } }),
        ),
      )
      yield* promise(() => mkdir(join(historical, "data"), { recursive: true }))
      yield* promise(() =>
        writeFile(
          join(historical, "data/evaluation-manifest.json"),
          JSON.stringify({ run: { number: 7 } }),
        ),
      )
      yield* promise(() => mkdir(join(unreadable, "evaluation.json"), { recursive: true }))
      yield* promise(() => mkdir(quality))
      yield* promise(() => writeFile(join(root, "0010-not-a-directory-eval"), "ignored"))

      const runs = yield* evaluationRunDirectories(root, ".quality/evaluations")
      assert.deepStrictEqual(
        runs.map(({ number, name, display }) => ({ number, name, display })),
        [
          {
            number: 3,
            name: "0003-unreadable-eval",
            display: ".quality/evaluations/0003-unreadable-eval",
          },
          {
            number: 4,
            name: "0004-quality-eval",
            display: ".quality/evaluations/0004-quality-eval",
          },
          {
            number: 7,
            name: "0002-historical-eval",
            display: ".quality/evaluations/0002-historical-eval",
          },
          {
            number: 9,
            name: "0001-current-eval",
            display: ".quality/evaluations/0001-current-eval",
          },
        ],
      )
      assert.strictEqual(yield* nextEvaluationRunNumber(root), 10)
    }).pipe(Effect.provide(services)),
  )
})
