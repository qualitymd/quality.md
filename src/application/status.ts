import * as Effect from "effect/Effect"
import * as FileSystem from "effect/FileSystem"
import * as Path from "effect/Path"

import { FileSystemFailure } from "../domain/errors.ts"
import { commandResult, ExitCode, type CommandResult } from "../domain/command-result.ts"
import { jsonDocument } from "../domain/json.ts"
import { invalidDocumentResult, lintDocument } from "../domain/lint/lint.ts"
import type { Action, LintResult } from "../domain/lint/result.ts"
import { parseQualityDocument, QualityDocumentParseError } from "../domain/model/document.ts"
import {
  decodeModel,
  effectiveSource,
  flattenModel,
  projectModel,
  type Area,
  type QualityModel,
} from "../domain/model/model.ts"
import { resolveWorkspace, type Workspace } from "../services/workspace.ts"

type Readiness =
  | "missing-model"
  | "invalid-model"
  | "ready-to-evaluate"
  | "has-evaluation-history"
  | "needs-evaluation-reconciliation"

interface EvaluationRunSummary {
  readonly path: string
  readonly reportable: boolean
  readonly lifecycle?: string
  readonly stale: boolean
  readonly dataArtifacts: number
  readonly gaps: number
  readonly problem?: string
}

interface EvaluationHistory {
  readonly path?: string
  readonly runs: number
  readonly latest?: EvaluationRunSummary
  readonly items: ReadonlyArray<EvaluationRunSummary>
  readonly summary: {
    readonly reportable: number
    readonly incomplete: number
    readonly awaitingEvaluator?: number
    readonly stale: number
    readonly problems: number
  }
}

interface StatusSnapshot {
  readonly schemaVersion: 2
  readonly path: string
  readonly workspace?: {
    readonly root: string
    readonly model: string
    readonly config: string
    readonly configPresent: boolean
    readonly dataDir: string
    readonly evaluationDir: string
    readonly changelogDir: string
    readonly logDir: string
  }
  readonly readiness: Readiness
  readonly model: Record<string, unknown>
  readonly evaluations: EvaluationHistory
  readonly nextActions: ReadonlyArray<Action>
}

const emptyHistory = (path?: string): EvaluationHistory => ({
  ...(path === undefined ? {} : { path }),
  runs: 0,
  items: [],
  summary: { reportable: 0, incomplete: 0, stale: 0, problems: 0 },
})

const parseLint = (
  path: string,
  raw: string,
): { readonly result: LintResult; readonly model?: QualityModel } => {
  try {
    const document = parseQualityDocument(path, raw)
    const result = lintDocument(document).result
    return { result, ...(result.valid ? { model: decodeModel(document) } : {}) }
  } catch (cause) {
    if (cause instanceof QualityDocumentParseError) return { result: invalidDocumentResult(path) }
    throw cause
  }
}

const sourceCoverage = (model: QualityModel) => {
  const rows: Array<Record<string, unknown>> = []
  const visit = (area: Area, path: ReadonlyArray<string>, fallback: string) => {
    const resolved = effectiveSource(model, path)
    rows.push({
      areaPath: path,
      label: area.title || fallback,
      sourceState: resolved.state,
      ...(resolved.state === "default" ? {} : { source: resolved.selector }),
      factors: Object.keys(area.factors ?? {}).length,
      requirements: Object.keys(area.requirements ?? {}).length,
      childAreas: Object.keys(area.areas ?? {}).length,
    })
    for (const [name, child] of Object.entries(area.areas ?? {}).sort(([a], [b]) =>
      a.localeCompare(b),
    )) {
      visit(child, [...path, name], name)
    }
  }
  visit(model, [], "Model")
  return rows
}

const modelShape = (model: QualityModel) => {
  const shape = { areas: 0, factors: 0, requirements: 0, ratingLevels: model.ratingScale.length }
  for (const element of flattenModel(projectModel(model))) {
    if (element.kind === "area") shape.areas += 1
    else if (element.kind === "factor") shape.factors += 1
    else shape.requirements += 1
  }
  return shape
}

const inspectArtifact = (
  path: string,
  currentModel: string,
  snapshot: string,
  artifactRaw: string,
): EvaluationRunSummary => {
  try {
    const artifact = JSON.parse(artifactRaw) as {
      readonly state?: {
        readonly status?: string
        readonly workUnits?: Readonly<Record<string, { readonly status?: string }>>
        readonly pendingEvaluatorCalls?: ReadonlyArray<unknown>
      }
      readonly results?: { readonly payloads?: ReadonlyArray<unknown> }
    }
    const lifecycle = artifact.state?.status ?? ""
    const gaps = Object.values(artifact.state?.workUnits ?? {}).filter(
      (unit) => unit.status !== "completed",
    ).length
    return {
      path,
      reportable: gaps === 0 && lifecycle === "completed",
      ...(lifecycle === "" ? {} : { lifecycle }),
      stale: snapshot !== currentModel,
      dataArtifacts: artifact.results?.payloads?.length ?? 0,
      gaps,
    }
  } catch (cause) {
    return {
      path,
      reportable: false,
      stale: snapshot !== currentModel,
      dataArtifacts: 0,
      gaps: 0,
      problem: cause instanceof Error ? cause.message : String(cause),
    }
  }
}

const evaluationHistory = (workspace: Workspace, currentModel: string) =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const paths = yield* Path.Path
    if (!(yield* fs.exists(workspace.evaluations.abs)))
      return emptyHistory(workspace.evaluations.rel)
    const runs: Array<{ readonly number: number; readonly name: string }> = []
    for (const name of yield* fs.readDirectory(workspace.evaluations.abs)) {
      const runAbs = paths.join(workspace.evaluations.abs, name)
      if ((yield* fs.stat(runAbs)).type !== "Directory") continue
      const artifactPath = paths.join(runAbs, "evaluation.json")
      if (yield* fs.exists(artifactPath)) {
        try {
          const artifact = JSON.parse(yield* fs.readFileString(artifactPath)) as {
            readonly manifest?: { readonly run?: { readonly number?: number } }
          }
          runs.push({ number: artifact.manifest?.run?.number ?? Number(name.slice(0, 4)), name })
        } catch {
          runs.push({ number: Number(name.slice(0, 4)), name })
        }
        continue
      }
      const match = /^(\d{4})-([a-z0-9-]+)-eval$/.exec(name)
      if (match === null || match[2]!.split("-").includes("quality")) continue
      runs.push({ number: Number(match[1]), name })
    }
    runs.sort((a, b) => a.number - b.number || a.name.localeCompare(b.name))
    const items: Array<EvaluationRunSummary> = []
    for (const run of runs) {
      const runAbs = paths.join(workspace.evaluations.abs, run.name)
      const path = `${workspace.evaluations.rel}/${run.name}`
      try {
        const snapshot = yield* fs.readFileString(paths.join(runAbs, "model-snapshot.md"))
        const artifactPath = paths.join(runAbs, "evaluation.json")
        if (yield* fs.exists(artifactPath)) {
          items.push(
            inspectArtifact(path, currentModel, snapshot, yield* fs.readFileString(artifactPath)),
          )
        } else {
          const data = paths.join(runAbs, "data")
          const artifacts = (yield* fs.exists(data))
            ? (yield* fs.readDirectory(data)).filter((name) => name.endsWith(".json")).length
            : 0
          items.push({
            path,
            reportable: false,
            stale: snapshot !== currentModel,
            dataArtifacts: artifacts,
            gaps: 1,
          })
        }
      } catch (cause) {
        items.push({
          path,
          reportable: false,
          stale: false,
          dataArtifacts: 0,
          gaps: 0,
          problem: cause instanceof Error ? cause.message : String(cause),
        })
      }
    }
    const summary = {
      reportable: items.filter((item) => item.reportable).length,
      incomplete: items.filter((item) => !item.reportable).length,
      awaitingEvaluator: items.filter((item) => item.lifecycle === "awaiting_evaluator").length,
      stale: items.filter((item) => item.stale).length,
      problems: items.filter((item) => item.problem !== undefined).length,
    }
    return {
      path: workspace.evaluations.rel,
      runs: items.length,
      ...(items.length === 0 ? {} : { latest: items.at(-1)! }),
      items,
      summary: {
        reportable: summary.reportable,
        incomplete: summary.incomplete,
        ...(summary.awaitingEvaluator === 0
          ? {}
          : { awaitingEvaluator: summary.awaitingEvaluator }),
        stale: summary.stale,
        problems: summary.problems,
      },
    } satisfies EvaluationHistory
  })

const actionsFor = (
  path: string,
  readiness: Readiness,
  history: EvaluationHistory,
): ReadonlyArray<Action> => {
  if (readiness === "ready-to-evaluate") {
    return [
      {
        id: "evaluation-create",
        label: "Create an evaluation run",
        command: `qualitymd evaluation create --model ${path}`,
      },
    ]
  }
  const latest = history.latest
  if (readiness === "needs-evaluation-reconciliation" && latest !== undefined) {
    if (latest.problem === undefined && latest.lifecycle === "awaiting_evaluator") {
      return [
        {
          id: "evaluation-run-reemit",
          label: "Resume the awaiting evaluation run to recover its pending work request",
          command: `qualitymd evaluation run --model ${path} --resume ${latest.path} --json`,
        },
      ]
    }
    if (latest.problem !== undefined || !latest.reportable) {
      return [
        {
          id: "evaluation-status-latest",
          label: "Inspect the latest evaluation run",
          command: `qualitymd evaluation status --model ${path} ${latest.path}`,
        },
      ]
    }
    if (latest.stale) {
      return [
        {
          id: "evaluation-create",
          label: "Create a fresh evaluation run",
          command: `qualitymd evaluation create --model ${path}`,
        },
      ]
    }
  }
  if (readiness === "has-evaluation-history" && latest?.reportable === true) {
    return [
      {
        id: "report-build",
        label: "Build the latest evaluation report",
        command: `qualitymd evaluation report build --model ${path} ${latest.path}`,
      },
    ]
  }
  return []
}

const snapshot = (path: string) =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    if (!(yield* fs.exists(path))) {
      return {
        schemaVersion: 2,
        path,
        readiness: "missing-model",
        model: { present: false },
        evaluations: emptyHistory(),
        nextActions: [
          { id: "init", label: "Create a starter QUALITY.md", command: `qualitymd init ${path}` },
        ],
      } satisfies StatusSnapshot
    }
    const raw = yield* fs.readFileString(path)
    const parsed = parseLint(path, raw)
    const lint = {
      summary: parsed.result.summary,
      ...(parsed.result.findings.length === 0 ? {} : { findings: parsed.result.findings }),
    }
    if (!parsed.result.valid || parsed.model === undefined) {
      return {
        schemaVersion: 2,
        path,
        readiness: "invalid-model",
        model: { present: true, valid: false, lint },
        evaluations: emptyHistory(),
        nextActions: [
          parsed.result.summary.fixable > 0
            ? {
                id: "fix",
                label: "Apply deterministic lint repairs",
                command: `qualitymd lint --fix ${path}`,
              }
            : { id: "lint", label: "Review lint findings", command: `qualitymd lint ${path}` },
        ],
      } satisfies StatusSnapshot
    }
    const workspace = yield* resolveWorkspace({ model: path })
    const history = yield* evaluationHistory(workspace, raw)
    const readiness: Readiness =
      history.runs === 0
        ? "ready-to-evaluate"
        : history.summary.incomplete > 0 ||
            history.summary.stale > 0 ||
            history.summary.problems > 0
          ? "needs-evaluation-reconciliation"
          : "has-evaluation-history"
    return {
      schemaVersion: 2,
      path,
      workspace: {
        root: workspace.workspaceRoot.repoRel,
        model: workspace.model.rel,
        config: workspace.config.rel,
        configPresent: workspace.configPresent,
        dataDir: workspace.dataDir.rel,
        evaluationDir: workspace.evaluations.rel,
        changelogDir: workspace.changelog.rel,
        logDir: workspace.logs.rel,
      },
      readiness,
      model: {
        present: true,
        valid: true,
        lint,
        shape: modelShape(parsed.model),
        sourceCoverage: sourceCoverage(parsed.model),
      },
      evaluations: history,
      nextActions: actionsFor(path, readiness, history),
    } satisfies StatusSnapshot
  }).pipe(
    Effect.mapError((cause) =>
      cause instanceof FileSystemFailure
        ? cause
        : new FileSystemFailure({ detail: cause instanceof Error ? cause.message : String(cause) }),
    ),
  )

const renderHuman = (value: StatusSnapshot) => {
  const lint = value.model.lint as { readonly summary: LintResult["summary"] } | undefined
  const present = value.model.present === true
  const valid = value.model.valid === true
  const model = !present
    ? "absent"
    : valid
      ? "present, valid"
      : `present, invalid (${lint?.summary.errors ?? 0} error(s), ${lint?.summary.warnings ?? 0} warning(s))`
  const shape = value.model.shape as
    | { readonly areas: number; readonly factors: number; readonly requirements: number }
    | undefined
  let output = "Workspace Status\n"
  if (value.workspace !== undefined) output += `- Workspace: ${value.workspace.root}\n`
  output += `- Model file: ${value.workspace?.model ?? value.path}: ${model}\n`
  if (shape !== undefined)
    output += `- Model: ${shape.areas} area(s), ${shape.factors} factor(s), ${shape.requirements} requirement(s)\n`
  output += `- Evaluation history: ${value.evaluations.runs} run(s), ${value.evaluations.summary.incomplete} incomplete, ${value.evaluations.summary.stale} stale\n`
  output += `- Readiness: ${value.readiness}\n`
  return output
}

export const statusCommand = (input: {
  readonly path: string
  readonly json: boolean
}): Effect.Effect<CommandResult, FileSystemFailure, FileSystem.FileSystem | Path.Path> => {
  if (input.path === "-") {
    return Effect.succeed(
      commandResult("", {
        stderr: "qualitymd: status does not read from stdin; pass a file path\n",
        exitCode: ExitCode.usage,
      }),
    )
  }
  return snapshot(input.path).pipe(
    Effect.map((value) => {
      const action = value.nextActions.at(0)
      return commandResult(input.json ? jsonDocument(value) : renderHuman(value), {
        stderr: input.json || action === undefined ? "" : `\nNext: ${action.command}\n`,
      })
    }),
  )
}
