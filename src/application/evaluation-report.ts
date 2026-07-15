import * as Effect from "effect/Effect"
import * as FileSystem from "effect/FileSystem"
import * as Path from "effect/Path"

import { FileSystemFailure } from "../domain/errors.ts"
import { commandResult, ExitCode, type CommandResult } from "../domain/command-result.ts"
import { buildReportTree, type ReportManifest } from "../domain/evaluation/report.ts"
import { planEvaluation } from "../domain/evaluation/plan.ts"
import { jsonDocument } from "../domain/json.ts"
import { parseQualityDocument } from "../domain/model/document.ts"
import { decodeModel } from "../domain/model/model.ts"
import { atomicWriteFileString } from "../services/atomic-file.ts"
import { resolveRun, type RunFlags } from "./evaluation-inspect.ts"

type Json = Record<string, unknown>

export interface ReportArtifact {
  readonly schemaVersion: number
  readonly kind: string
  readonly manifest: ReportManifest & {
    readonly evaluator?: string
    readonly evaluatorKind?: string
    readonly concurrency?: number
    readonly evaluatorCapabilities?: Json
  }
  readonly results: { readonly payloads: ReadonlyArray<{ readonly payload: Json }> }
  outputs?: Json
}

export interface ReportBuildResult {
  readonly reportMd: string
  readonly evaluationOutputResult: string
  readonly ratingResult: Json
  readonly outputs: Json
}

const failure = (detail: string, exitCode: ExitCode = ExitCode.problems) =>
  commandResult("", { stderr: `qualitymd: ${detail}\n`, exitCode })

const write = (
  fs: FileSystem.FileSystem,
  paths: Path.Path,
  root: string,
  relative: string,
  content: string,
) =>
  Effect.gen(function* () {
    const target = paths.join(root, relative)
    yield* fs.makeDirectory(paths.dirname(target), { recursive: true })
    yield* atomicWriteFileString(target, content, { mode: 0o644 })
  })

export const buildReportsAtRun = (
  runAbs: string,
  display: string,
  artifact?: ReportArtifact,
): Effect.Effect<ReportBuildResult, Error, FileSystem.FileSystem | Path.Path> =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const paths = yield* Path.Path
    if (!(yield* fs.exists(runAbs)) || (yield* fs.stat(runAbs)).type !== "Directory")
      throw new Error(`reading run ${display}: no such file or directory`)
    const snapshotPath = paths.join(runAbs, "model-snapshot.md")
    if (!(yield* fs.exists(snapshotPath)))
      throw new Error(`${display} is not an evaluation run folder: missing model-snapshot.md`)
    const current =
      artifact ??
      (JSON.parse(
        yield* fs.readFileString(paths.join(runAbs, "evaluation.json")),
      ) as ReportArtifact)
    if (current.schemaVersion !== 9 || current.kind !== "EvaluationRun")
      throw new Error(
        `evaluation artifact schema ${current.schemaVersion} is incompatible with schema 9`,
      )
    const model = decodeModel(
      parseQualityDocument(snapshotPath, yield* fs.readFileString(snapshotPath)),
    )
    const plan = planEvaluation(model, current.manifest.plannedScope)
    const manifestPayload: Json = {
      createdAt: current.manifest.createdAt,
      evaluationId: current.manifest.evaluationId,
      kind: "EvaluationManifest",
      model: current.manifest.model,
      plannedScope: current.manifest.plannedScope,
      requestedScope: current.manifest.requestedScope,
      run: current.manifest.run,
      schemaVersion: 3,
    }
    const payloads = [manifestPayload, ...current.results.payloads.map((entry) => entry.payload)]
    let workspace = runAbs
    let runRel = `.quality/evaluations/${paths.basename(runAbs)}`
    while (paths.dirname(workspace) !== workspace) {
      if (yield* fs.exists(paths.join(workspace, current.manifest.model))) {
        runRel = paths.relative(workspace, runAbs).replaceAll("\\", "/")
        break
      }
      workspace = paths.dirname(workspace)
    }
    const tree = buildReportTree({ model, manifest: current.manifest, plan, payloads, runRel })
    for (const entry of tree.payloadFiles)
      yield* write(fs, paths, runAbs, entry.path, jsonDocument(entry.payload))
    yield* write(fs, paths, runAbs, "data/evaluation-output-result.json", jsonDocument(tree.output))
    for (const report of tree.reports) yield* write(fs, paths, runAbs, report.path, report.content)
    const reportMd = `${display}/report.md`
    const evaluationOutputResult = `${display}/data/evaluation-output-result.json`
    const outputs = {
      reportMd: `${runRel}/report.md`,
      evaluationOutput: tree.output,
      rating: tree.rating,
    }
    return { reportMd, evaluationOutputResult, ratingResult: tree.rating, outputs }
  })

export interface EvaluationReportBuildInput extends RunFlags {
  readonly json: boolean
}

export const evaluationReportBuildCommand = (
  input: EvaluationReportBuildInput,
): Effect.Effect<CommandResult, FileSystemFailure, FileSystem.FileSystem | Path.Path> =>
  Effect.gen(function* () {
    const resolved = yield* resolveRun(input)
    if ("error" in resolved) return failure(resolved.error, ExitCode.usage)
    const result = yield* buildReportsAtRun(resolved.run.absolute, resolved.run.display).pipe(
      Effect.result,
    )
    if (result._tag === "Failure") return failure(result.failure.message)
    const receipt = {
      schemaVersion: 3,
      path: resolved.run.display,
      reportMd: result.success.reportMd,
      evaluationOutputResult: result.success.evaluationOutputResult,
      ratingResult: result.success.ratingResult,
    }
    if (input.json) return commandResult(jsonDocument(receipt))
    return commandResult(`Built evaluation report: ${result.success.reportMd}\n`)
  }).pipe(
    Effect.mapError(
      (cause) =>
        new FileSystemFailure({
          detail: cause instanceof Error ? cause.message : String(cause),
        }),
    ),
  )
