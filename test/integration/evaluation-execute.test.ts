import { assert, it } from "@effect/vitest"
import * as BunFileSystem from "@effect/platform-bun/BunFileSystem"
import * as BunPath from "@effect/platform-bun/BunPath"
import * as Effect from "effect/Effect"
import * as Layer from "effect/Layer"
import { createHash } from "node:crypto"
import { mkdir, mkdtemp, readFile, rm, writeFile } from "node:fs/promises"
import { join } from "node:path"

import { executeHarnessRun } from "../../src/application/evaluation-execute.ts"
import { HostRuntime, type HostRuntimeService } from "../../src/services/host-runtime.ts"

const digest = (value: string | Buffer) => createHash("sha256").update(value).digest("hex")

it.effect("preserves deterministic harness checkpoint bytes", () =>
  Effect.scoped(
    Effect.gen(function* () {
      const directory = yield* Effect.acquireRelease(
        Effect.promise(async () => {
          await mkdir(join(process.cwd(), "tmp"), { recursive: true })
          return mkdtemp(join(process.cwd(), "tmp/qualitymd-0202-byte-"))
        }),
        (path) => Effect.promise(() => rm(path, { recursive: true })),
      )
      const model = join(directory, "QUALITY.md")
      yield* Effect.promise(() => mkdir(join(directory, "src")))
      yield* Effect.promise(() =>
        writeFile(join(directory, "src/example.ts"), "export const ready = true\n"),
      )
      yield* Effect.promise(() =>
        writeFile(
          model,
          `---
title: Byte fixture
source: src
ratingScale:
  - level: target
    title: Target
    description: Meets the target.
    criterion: Meets the target.
  - level: unacceptable
    title: Unacceptable
    description: Does not meet the target.
    criterion: Does not meet the target.
factors:
  reliability:
    title: Reliability
    description: Reliable operation.
    requirements:
      ready:
        title: Ready
        assessment: Inspect readiness.
---
Body guidance.
`,
        ),
      )
      const services = Layer.mergeAll(
        BunFileSystem.layer,
        BunPath.layer,
        Layer.succeed(HostRuntime, {
          cwd: directory,
          environment: {},
          hardwareConcurrency: 2,
          currentTimeMillis: Effect.succeed(Date.UTC(2026, 6, 14)),
          randomBytes: (length) => Effect.succeed(new Uint8Array(length).fill(7)),
          readStdin: Effect.succeed(""),
          which: () => null,
          codexAuthenticated: () => false,
        } satisfies HostRuntimeService),
      )
      const result = yield* executeHarnessRun({
        model,
        evaluationDir: "",
        area: "",
        factors: [],
        evaluator: "harness",
        resume: "",
        evaluatorResult: "",
        dryRun: false,
        json: true,
      }).pipe(Effect.provide(services))
      assert.strictEqual(result.exitCode, 0)
      const receipt = JSON.parse(result.stdout) as { path: string }
      const run = join(directory, receipt.path)
      const files = yield* Effect.forEach(
        ["evaluation.json", "logs/events.jsonl", "logs/evaluator-calls.jsonl", "model-snapshot.md"],
        (path) => Effect.promise(() => readFile(join(run, path))),
      )
      assert.deepStrictEqual(
        {
          stdout: digest(result.stdout),
          stderr: digest(result.stderr),
          files: files.map(digest),
        },
        {
          stdout: "b63713fb24e9c5883d8f60e3cff60637e8833ae69ada58f31d3fe1826d699759",
          stderr: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
          files: [
            "2729e257c1d9bf0501b9cb1579eb035e7c398c0b6405d3b02d85020acda72e2e",
            "6c15831d65ed4af8143f1c69152a46c40248bb0bdf2ba669b02d879aacb89d0a",
            "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
            "6b247c30352ac60048741350d5d7326dc73f19cc3779d571ce5d05687a94fdc5",
          ],
        },
      )
    }),
  ),
)
