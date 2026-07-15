import * as Effect from "effect/Effect"
import * as FileSystem from "effect/FileSystem"
import * as Path from "effect/Path"
import * as Result from "effect/Result"

import { runDirectoryNumber } from "../domain/evaluation/run.ts"

export interface EvaluationRunDirectory {
  readonly number: number
  readonly name: string
  readonly absolute: string
  readonly display: string
}

const readCandidate = Effect.fn("qualitymd.readEvaluationRunCandidate")(function* (path: string) {
  const fs = yield* FileSystem.FileSystem
  if (!(yield* fs.exists(path))) return undefined
  const result = yield* fs.readFileString(path).pipe(Effect.result)
  return Result.isSuccess(result) ? result.success : undefined
})

export const evaluationRunDirectories = Effect.fn("qualitymd.evaluationRunDirectories")(function* (
  absolute: string,
  relative: string,
) {
  const fs = yield* FileSystem.FileSystem
  const paths = yield* Path.Path
  if (!(yield* fs.exists(absolute))) return []
  const candidates = yield* Effect.forEach(yield* fs.readDirectory(absolute), (name) =>
    Effect.gen(function* () {
      const run = paths.join(absolute, name)
      if ((yield* fs.stat(run)).type !== "Directory") return undefined
      const number = runDirectoryNumber({
        name,
        evaluationArtifact: yield* readCandidate(paths.join(run, "evaluation.json")),
        historicalManifest: yield* readCandidate(paths.join(run, "data/evaluation-manifest.json")),
      })
      return number === undefined
        ? undefined
        : {
            number,
            name,
            absolute: run,
            display: `${relative}/${name}`,
          }
    }),
  )
  return candidates
    .filter((candidate): candidate is EvaluationRunDirectory => candidate !== undefined)
    .sort((left, right) => left.number - right.number || left.name.localeCompare(right.name))
})

export const nextEvaluationRunNumber = Effect.fn("qualitymd.nextEvaluationRunNumber")(function* (
  absolute: string,
) {
  const runs = yield* evaluationRunDirectories(absolute, absolute)
  return Math.max(0, ...runs.map((run) => run.number)) + 1
})
