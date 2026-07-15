import * as Effect from "effect/Effect"
import * as FileSystem from "effect/FileSystem"
import * as Path from "effect/Path"

import { FileSystemFailure } from "../domain/errors.ts"
import { commandResult, ExitCode, type CommandResult } from "../domain/command-result.ts"
import { jsonDocument } from "../domain/json.ts"
import {
  EvaluationSchemaVersion,
  newEvaluationIdentity,
  resolveScope,
  scopeSlug,
  type EvaluationManifest,
} from "../domain/evaluation/run.ts"
import { decodeModel } from "../domain/model/model.ts"
import { parseQualityDocument } from "../domain/model/document.ts"
import { HostRuntime } from "../services/host-runtime.ts"
import { resolveWorkspace } from "../services/workspace.ts"

export interface EvaluationCreateInput {
  readonly modelArgument?: string
  readonly model: string
  readonly area: string
  readonly factors: ReadonlyArray<string>
  readonly evaluationDir: string
  readonly json: boolean
}

const usage = (detail: string) =>
  commandResult("", { stderr: `qualitymd: ${detail}\n`, exitCode: ExitCode.usage })

const nextNumber = (directory: string) =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const paths = yield* Path.Path
    if (!(yield* fs.exists(directory))) return 1
    let maximum = 0
    for (const name of yield* fs.readDirectory(directory)) {
      const run = paths.join(directory, name)
      if ((yield* fs.stat(run)).type !== "Directory") continue
      let number = /^(\d{4})-[a-z0-9-]+-eval$/.exec(name)?.[1]
      for (const manifestPath of ["data/evaluation-manifest.json", "evaluation.json"]) {
        const file = paths.join(run, manifestPath)
        if (!(yield* fs.exists(file))) continue
        try {
          const parsed = JSON.parse(yield* fs.readFileString(file)) as {
            readonly run?: { readonly number?: number }
            readonly manifest?: { readonly run?: { readonly number?: number } }
          }
          number = String(parsed.run?.number ?? parsed.manifest?.run?.number ?? number)
        } catch {
          // A malformed run still contributes its parseable directory number.
        }
        break
      }
      maximum = Math.max(maximum, Number(number ?? 0))
    }
    return maximum + 1
  })

export const evaluationCreateCommand = (
  input: EvaluationCreateInput,
): Effect.Effect<
  CommandResult,
  FileSystemFailure,
  FileSystem.FileSystem | Path.Path | HostRuntime
> => {
  if (input.modelArgument !== undefined && input.model !== "") {
    return Effect.succeed(usage("pass a model argument or --model, not both"))
  }
  const modelPath = (input.modelArgument ?? input.model) || "QUALITY.md"
  if (modelPath === ".")
    return Effect.succeed(
      usage(`--model ${JSON.stringify(modelPath)} must name a QUALITY.md file, not a directory`),
    )
  return Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const paths = yield* Path.Path
    const runtime = yield* HostRuntime
    const workspace = yield* resolveWorkspace({
      model: modelPath,
      ...(input.evaluationDir === "" ? {} : { evaluationDir: input.evaluationDir }),
    })
    const modelRaw = yield* fs.readFileString(workspace.model.abs)
    const model = decodeModel(parseQualityDocument(workspace.model.abs, modelRaw))
    let scope: ReturnType<typeof resolveScope>
    try {
      scope = resolveScope(model, input.area, input.factors)
    } catch (cause) {
      return usage(cause instanceof Error ? cause.message : String(cause))
    }
    yield* fs.makeDirectory(workspace.evaluations.abs, { recursive: true, mode: 0o755 })
    const number = yield* nextNumber(workspace.evaluations.abs)
    const label = `${String(number).padStart(4, "0")}-${scopeSlug(scope.plannedScope)}-eval`
    const runAbs = paths.join(workspace.evaluations.abs, label)
    const runRel = `${workspace.evaluations.rel}/${label}`
    yield* fs.makeDirectory(runAbs, { mode: 0o755 })
    yield* fs.makeDirectory(paths.join(runAbs, "data"), { mode: 0o755 })
    yield* fs.writeFileString(paths.join(runAbs, "model-snapshot.md"), modelRaw, { mode: 0o644 })
    const identity = newEvaluationIdentity(
      new Date(yield* runtime.currentTimeMillis),
      yield* runtime.randomBytes(12),
    )
    const manifest: EvaluationManifest = {
      schemaVersion: EvaluationSchemaVersion,
      kind: "EvaluationManifest",
      ...identity,
      model: workspace.model.rel,
      ...scope,
      run: { number, label },
    }
    yield* fs.writeFileString(
      paths.join(runAbs, "data/evaluation-manifest.json"),
      jsonDocument(manifest),
      { mode: 0o644 },
    )
    const nextActions = [
      {
        id: "evaluation-data-set",
        label: "Record evaluation data",
        command: `qualitymd evaluation data set${modelPath === "QUALITY.md" ? "" : ` --model ${modelPath}`} ${runRel} < payloads.json`,
      },
    ]
    if (input.json) {
      return commandResult(jsonDocument({ schemaVersion: 3, path: runRel, nextActions }))
    }
    return commandResult("", {
      stderr: `Created ${runRel}\n\nNext: ${nextActions[0]!.command}\n`,
    })
  }).pipe(
    Effect.mapError((cause) =>
      cause instanceof FileSystemFailure
        ? cause
        : new FileSystemFailure({ detail: cause instanceof Error ? cause.message : String(cause) }),
    ),
  )
}
