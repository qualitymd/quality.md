import * as Effect from "effect/Effect"
import * as FileSystem from "effect/FileSystem"
import * as Path from "effect/Path"

import { FileSystemFailure } from "../domain/errors.ts"
import { commandResult, ExitCode, type CommandResult } from "../domain/command-result.ts"
import { jsonDocument } from "../domain/json.ts"
import { decodeModel } from "../domain/model/model.ts"
import { parseQualityDocument } from "../domain/model/document.ts"
import { resolveWorkspace } from "../services/workspace.ts"
import { evaluationRunDirectories } from "./evaluation-runs.ts"

export interface RunFlags {
  readonly run?: string
  readonly latest: boolean
  readonly model: string
  readonly evaluationDir: string
}

export interface ResolvedRun {
  readonly absolute: string
  readonly display: string
  readonly model: string
}

interface ArtifactDocument {
  readonly manifest: {
    readonly requestedScope: Record<string, unknown>
    readonly plannedScope: { readonly areaId: string; readonly factorFilter: ReadonlyArray<string> }
    readonly run: { readonly number: number }
  }
  readonly state: {
    readonly status: string
    readonly workUnits?: Readonly<Record<string, { readonly status?: string }>>
    readonly pendingEvaluatorCalls?: ReadonlyArray<{
      readonly requestId: string
      readonly workUnitId: string
      readonly attempt: number
    }>
  }
  readonly results: {
    readonly payloads: ReadonlyArray<{
      readonly payload?: { readonly kind?: string; readonly areaId?: string }
    }>
  }
}

const usage = (detail: string) =>
  commandResult("", { stderr: `qualitymd: ${detail}\n`, exitCode: ExitCode.usage })

export const resolveRun = (flags: RunFlags) =>
  Effect.gen(function* () {
    const paths = yield* Path.Path
    if (flags.run !== undefined && flags.latest)
      return { error: "pass a run path or --latest, not both" } as const
    if (flags.run === undefined && !flags.latest)
      return { error: "pass a run path or --latest" } as const
    if (flags.run !== undefined) {
      if (flags.model === "") {
        return {
          run: {
            absolute: paths.resolve(flags.run),
            display: flags.run,
            model: "",
          } satisfies ResolvedRun,
        } as const
      }
      const workspace = yield* resolveWorkspace({ model: flags.model })
      const absolute = paths.isAbsolute(flags.run)
        ? flags.run
        : paths.resolve(workspace.workspaceRoot.abs, flags.run)
      return {
        run: { absolute, display: flags.run, model: flags.model } satisfies ResolvedRun,
      } as const
    }
    const workspace = yield* resolveWorkspace({
      ...(flags.model === "" ? {} : { model: flags.model }),
      ...(flags.evaluationDir === "" ? {} : { evaluationDir: flags.evaluationDir }),
    })
    const runs = yield* evaluationRunDirectories(
      workspace.evaluations.abs,
      workspace.evaluations.rel,
    )
    const latest = runs.at(-1)
    if (latest === undefined) return { error: "no evaluation runs found" } as const
    return {
      run: {
        absolute: latest.absolute,
        display: latest.display,
        model: flags.model,
      } satisfies ResolvedRun,
    } as const
  })

const areaDataPath = (areaId: string) => {
  const area = areaId.slice(5)
  return `data/areas/${area}/area-analysis-result.json`
}

const inspect = (run: ResolvedRun) =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const paths = yield* Path.Path
    if (!(yield* fs.exists(run.absolute)) || (yield* fs.stat(run.absolute)).type !== "Directory") {
      throw new Error(`reading run ${run.display}: no such file or directory`)
    }
    const snapshotPath = paths.join(run.absolute, "model-snapshot.md")
    if (!(yield* fs.exists(snapshotPath)))
      throw new Error(`${run.display} is not an evaluation run folder: missing model-snapshot.md`)
    const model = decodeModel(
      parseQualityDocument(snapshotPath, yield* fs.readFileString(snapshotPath)),
    )
    const artifactPath = paths.join(run.absolute, "evaluation.json")
    if (!(yield* fs.exists(artifactPath))) {
      return {
        model,
        manifest: undefined,
        status: {
          schemaVersion: 3,
          path: run.display,
          reportable: false,
          data: { artifacts: 0 },
          gaps: [
            {
              kind: "missing-evaluation-data",
              ref: "data/evaluation-manifest.json",
              detail: "required evaluation manifest payload is missing",
            },
          ],
          nextActions: [
            {
              id: "evaluation-data-set",
              label: "Persist required evaluation data",
              command: `qualitymd evaluation data set ${run.display} \u003c payloads.json`,
            },
          ],
        },
      }
    }
    const artifact = JSON.parse(yield* fs.readFileString(artifactPath)) as ArtifactDocument
    const payloads = artifact.results?.payloads ?? []
    const areaId = artifact.manifest.plannedScope.areaId
    const hasScopedAnalysis = payloads.some(
      (entry) => entry.payload?.kind === "AreaAnalysisResult" && entry.payload.areaId === areaId,
    )
    const payloadKinds = new Set(payloads.map((entry) => entry.payload?.kind ?? ""))
    const gaps = [
      ...(hasScopedAnalysis
        ? []
        : [
            {
              kind: "missing-evaluation-data",
              ref: areaDataPath(areaId),
              detail: "required scoped area analysis payload is missing",
            },
          ]),
      ...(
        [
          ["FindingRankingResult", "data/advice/finding-ranking-result.json"],
          ["RecommendationResult", "data/advice/recommendations/"],
          ["RecommendationRankingResult", "data/advice/recommendation-ranking-result.json"],
          ["EvaluationSummaryResult", "data/advice/evaluation-summary-result.json"],
        ] as const
      ).flatMap(([kind, ref]) =>
        payloadKinds.has(kind)
          ? []
          : [
              {
                kind: "missing-evaluation-data",
                ref,
                detail: `required ${kind} payload is missing`,
              },
            ],
      ),
    ]
    const lifecycle = artifact.state.status
    const reportable = gaps.length === 0
    const awaiting =
      lifecycle === "awaiting_evaluator"
        ? (artifact.state.pendingEvaluatorCalls ?? []).map((call) => ({
            requestId: call.requestId,
            workUnitId: call.workUnitId,
            attempt: call.attempt,
          }))
        : []
    const nextActions =
      lifecycle === "awaiting_evaluator"
        ? [
            {
              id: "evaluation-run-reemit",
              label: "Recover the outstanding harness work requests",
              command: `qualitymd evaluation run --resume ${run.display} --json`,
            },
            {
              id: "evaluation-evaluator-result",
              label: "Submit harness judgment results",
              command: `qualitymd evaluation run --resume ${run.display} --evaluator-result - --json`,
            },
          ]
        : reportable
          ? [
              {
                id: "evaluation-report-build",
                label: "Build evaluation report",
                command: `qualitymd evaluation report build ${run.display}`,
              },
            ]
          : [
              {
                id: "evaluation-run-resume",
                label: "Resume the evaluation run",
                command: `qualitymd evaluation run --resume ${run.display}`,
              },
            ]
    return {
      model,
      manifest: artifact.manifest,
      status: {
        schemaVersion: 3,
        path: run.display,
        reportable,
        ...(lifecycle === "" ? {} : { lifecycle }),
        ...(awaiting.length === 0 ? {} : { awaitingEvaluator: awaiting }),
        data: { artifacts: payloads.length },
        gaps,
        nextActions,
      },
    }
  })

export const evaluationStatusCommand = (
  input: RunFlags & { readonly json: boolean },
): Effect.Effect<CommandResult, FileSystemFailure, FileSystem.FileSystem | Path.Path> =>
  Effect.gen(function* () {
    const resolved = yield* resolveRun(input)
    if ("error" in resolved) return usage(resolved.error)
    const inspected = yield* inspect(resolved.run)
    const status = inspected.status
    if (input.json) return commandResult(jsonDocument(status))
    let stdout = `Run: ${status.path}\nReportable: ${String(status.reportable)}\nData artifacts: ${status.data.artifacts}\n`
    if ("lifecycle" in status) stdout += `Lifecycle: ${status.lifecycle}\n`
    if ("awaitingEvaluator" in status) {
      for (const awaiting of status.awaitingEvaluator)
        stdout += `Awaiting harness judgment: ${awaiting.workUnitId} (request ${awaiting.requestId}, attempt ${awaiting.attempt})\n`
    }
    for (const gap of status.gaps) stdout += `- ${gap.kind} ${gap.ref}: ${gap.detail}\n`
    const action = status.nextActions.at(0)
    return commandResult(stdout, {
      stderr: action === undefined ? "" : `\nNext: ${action.command}\n`,
    })
  }).pipe(
    Effect.mapError((cause) =>
      cause instanceof FileSystemFailure
        ? cause
        : new FileSystemFailure({ detail: cause instanceof Error ? cause.message : String(cause) }),
    ),
  )

export const evaluationListCommand = (input: {
  readonly model: string
  readonly evaluationDir: string
  readonly state: string
  readonly json: boolean
}): Effect.Effect<CommandResult, FileSystemFailure, FileSystem.FileSystem | Path.Path> => {
  const states = ["", "all", "reportable", "complete", "incomplete", "awaiting"]
  if (!states.includes(input.state))
    return Effect.succeed(
      usage("--state must be one of: all, complete, reportable, incomplete, awaiting"),
    )
  return Effect.gen(function* () {
    const workspace = yield* resolveWorkspace({
      ...(input.model === "" ? {} : { model: input.model }),
      ...(input.evaluationDir === "" ? {} : { evaluationDir: input.evaluationDir }),
    })
    const directories = yield* evaluationRunDirectories(
      workspace.evaluations.abs,
      workspace.evaluations.rel,
    )
    const inspectedRuns = yield* Effect.forEach(directories, (directory) =>
      inspect({
        absolute: directory.absolute,
        display: directory.display,
        model: input.model,
      }).pipe(Effect.map((inspected) => ({ directory, inspected }))),
    )
    const runs = inspectedRuns.flatMap(({ directory, inspected }) => {
      if (inspected.manifest === undefined) return []
      const status = inspected.status
      if ((input.state === "reportable" || input.state === "complete") && !status.reportable)
        return []
      if (input.state === "incomplete" && status.reportable) return []
      if (
        input.state === "awaiting" &&
        (!("lifecycle" in status) || status.lifecycle !== "awaiting_evaluator")
      )
        return []
      return [
        {
          path: directory.display,
          rootArea: inspected.model.title || "",
          requestedScope: inspected.manifest.requestedScope,
          plannedScope: inspected.manifest.plannedScope,
          dataArtifacts: status.data.artifacts,
          reportable: status.reportable,
          ...("lifecycle" in status ? { lifecycle: status.lifecycle } : {}),
          gaps: status.gaps.length,
        },
      ]
    })
    const value = { schemaVersion: 3, runs }
    if (input.json) return commandResult(jsonDocument(value))
    return commandResult(
      runs.length === 0
        ? "No evaluation runs found.\n"
        : `${runs.map((run) => `${run.path}\t${run.rootArea}\treportable=${String(run.reportable)}${"lifecycle" in run ? `\tlifecycle=${run.lifecycle}` : ""}\tdata=${run.dataArtifacts}\tgaps=${run.gaps}`).join("\n")}\n`,
    )
  }).pipe(
    Effect.mapError((cause) =>
      cause instanceof FileSystemFailure
        ? cause
        : new FileSystemFailure({ detail: cause instanceof Error ? cause.message : String(cause) }),
    ),
  )
}
