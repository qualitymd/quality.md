import * as Cause from "effect/Cause"
import * as Effect from "effect/Effect"
import * as FileSystem from "effect/FileSystem"
import * as Path from "effect/Path"

import { FileSystemFailure } from "../domain/errors.ts"
import { commandResult, ExitCode, type CommandResult } from "../domain/command-result.ts"
import { resolveScope, scopeSlug } from "../domain/evaluation/run.ts"
import { jsonDocument } from "../domain/json.ts"
import {
  decodeModel,
  effectiveSource,
  findElement,
  flattenModel,
  projectModel,
  type QualityModel,
} from "../domain/model/model.ts"
import { parseQualityDocument } from "../domain/model/document.ts"
import { claudeEvaluator, codexEvaluator, harnessCapabilities } from "../adapters/evaluator.ts"
import type { EvaluatorService } from "../services/evaluator.ts"
import { HostRuntime, type HostRuntimeService } from "../services/host-runtime.ts"
import { resolveWorkspace, type Workspace } from "../services/workspace.ts"
import { detectSourceKind, validateSourceSelector } from "../services/source.ts"
import { executeHarnessRun } from "./evaluation-execute.ts"
import { executeProviderRun, resumeProviderRun } from "./evaluation-provider.ts"
import { resumeHarnessRun } from "./evaluation-resume.ts"

export interface EvaluationRunInput {
  readonly model: string
  readonly evaluationDir: string
  readonly area: string
  readonly factors: ReadonlyArray<string>
  readonly evaluator: string
  readonly resume: string
  readonly resumeDisplay?: string
  readonly evaluatorResult: string
  readonly dryRun: boolean
  readonly json: boolean
}

export interface EvaluatorDiscovery {
  readonly which: (command: string) => string | null
  readonly codexAuthenticated: () => boolean
}

const evaluatorDiscovery = (runtime: HostRuntimeService): EvaluatorDiscovery => ({
  which: runtime.which,
  codexAuthenticated: runtime.codexAuthenticated,
})

const usage = (detail: string) =>
  commandResult("", { stderr: `qualitymd: ${detail}\n`, exitCode: ExitCode.usage })

const resumeRun = (input: EvaluationRunInput) =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const paths = yield* Path.Path
    const runtime = yield* HostRuntime
    const runAbs = paths.resolve(input.resume)
    const artifactPath = paths.join(runAbs, "evaluation.json")
    if (!(yield* fs.exists(artifactPath))) return yield* resumeHarnessRun(input)
    const artifact = JSON.parse(yield* fs.readFileString(artifactPath)) as {
      readonly manifest?: {
        readonly model?: string
        readonly evaluator?: string
        readonly evaluatorKind?: string
      }
    }
    const evaluatorKind = artifact.manifest?.evaluatorKind ?? "harness"
    const evaluatorName = artifact.manifest?.evaluator ?? evaluatorKind
    if (input.evaluator !== "" && input.evaluator !== evaluatorName)
      return usage(
        `run evaluator is pinned to ${JSON.stringify(evaluatorName)}; omit --evaluator or start a new run`,
      )
    if (evaluatorKind === "harness") return yield* resumeHarnessRun(input)
    if (input.evaluatorResult !== "")
      return usage("--evaluator-result is only accepted for a harness evaluation run")
    const model = artifact.manifest?.model ?? "QUALITY.md"
    let directory = runAbs
    let modelPath = ""
    while (paths.dirname(directory) !== directory) {
      const candidate = paths.join(directory, model)
      if (yield* fs.exists(candidate)) {
        modelPath = candidate
        break
      }
      directory = paths.dirname(directory)
    }
    if (modelPath === "")
      return commandResult("", {
        stderr: `qualitymd: cannot locate ${model} for resumed run ${input.resume}\n`,
        exitCode: ExitCode.problems,
      })
    const workspace = yield* resolveWorkspace({ model: modelPath })
    const selected = selectEvaluator(evaluatorName, workspace, evaluatorDiscovery(runtime))
    if (selected.kind !== evaluatorKind)
      return commandResult("", {
        stderr: `qualitymd: configured evaluator ${JSON.stringify(evaluatorName)} changed kind from ${JSON.stringify(evaluatorKind)} to ${JSON.stringify(selected.kind)}; restore the original profile or start a new run\n`,
        exitCode: ExitCode.problems,
      })
    if (!("evaluate" in selected))
      return commandResult("", {
        stderr: `qualitymd: evaluator ${JSON.stringify(evaluatorName)} cannot service this provider run\n`,
        exitCode: ExitCode.problems,
      })
    return yield* resumeProviderRun(input, selected, workspace)
  }).pipe(
    Effect.mapError((cause) =>
      cause instanceof FileSystemFailure
        ? cause
        : new FileSystemFailure({ detail: cause instanceof Error ? cause.message : String(cause) }),
    ),
  )

export const selectEvaluator = (
  requested: string,
  workspace: Workspace,
  discovery: EvaluatorDiscovery,
) => {
  const name = requested || workspace.evaluation.evaluator || "auto"
  if (name === "harness") {
    return {
      name,
      kind: "harness" as const,
      capabilities: harnessCapabilities,
      reason:
        "requested harness evaluator: judgment is supplied by the invoking agent harness through checkpoints",
    }
  }
  const configuredEvaluator = (profileName: string) => {
    const profile = workspace.evaluators[profileName]
    if (profile === undefined) throw new Error(`unknown evaluator ${JSON.stringify(profileName)}`)
    const options = {
      name: profileName,
      ...(profile.model === undefined ? {} : { model: profile.model }),
      ...(profile.command === undefined ? {} : { command: profile.command }),
    }
    let evaluator: EvaluatorService
    if (profile.kind === "codex") evaluator = codexEvaluator(options)
    else if (profile.kind === "claude") evaluator = claudeEvaluator(options)
    else
      throw new Error(
        `evaluator profile ${JSON.stringify(name)} declares unsupported kind ${JSON.stringify(profile.kind)}; use codex or claude`,
      )
    return { ...evaluator, reason: "configured evaluator profile" }
  }
  if (name === "auto") {
    const candidates: Array<{
      readonly name: string
      readonly executable: boolean
      readonly structuredOutput: boolean
      readonly authenticated: boolean
      readonly usable: boolean
      readonly evidence: ReadonlyArray<string>
    }> = []
    const codexExecutable = discovery.which("codex") !== null
    const codexAuthenticated = codexExecutable && discovery.codexAuthenticated()
    candidates.push({
      name: "codex",
      executable: codexExecutable,
      structuredOutput: true,
      authenticated: codexAuthenticated,
      usable: codexExecutable && codexAuthenticated,
      evidence: [
        codexExecutable ? "agent runtime executable found" : "agent runtime executable not found",
        "non-interactive structured output available through the Codex SDK",
        codexAuthenticated ? "authenticated (codex login status)" : "not authenticated",
      ],
    })
    if (codexExecutable && codexAuthenticated) {
      const evaluator = codexEvaluator()
      return {
        ...evaluator,
        reason:
          "auto: codex agent runtime is installed, authenticated, and supports non-interactive structured output",
        candidates,
      }
    }
    const claudeExecutable = discovery.which("claude") !== null
    candidates.push({
      name: "claude",
      executable: claudeExecutable,
      structuredOutput: true,
      authenticated: claudeExecutable,
      usable: claudeExecutable,
      evidence: [
        claudeExecutable ? "agent runtime executable found" : "agent runtime executable not found",
        "non-interactive structured output available through the Claude Agent SDK",
        claudeExecutable
          ? "authentication assumed because the runtime exposes no non-interactive status probe"
          : "authentication not checked without an executable",
      ],
    })
    if (claudeExecutable) {
      const evaluator = claudeEvaluator()
      return {
        ...evaluator,
        reason:
          "auto: claude agent runtime is installed and supports non-interactive structured output; authentication is assumed",
        candidates,
      }
    }
    throw new Error(
      "no evaluator is available; install and authenticate codex or claude, select a configured codex/claude profile, or pass --evaluator harness from a capable invoking agent",
    )
  }
  if (name === "codex")
    return { ...codexEvaluator(), reason: "requested built-in SDK agent evaluator" }
  if (name === "claude")
    return { ...claudeEvaluator(), reason: "requested built-in SDK agent evaluator" }
  return configuredEvaluator(name)
}

const areaPath = (id: string) => (id === "area:root" ? [] : id.slice(5).split("/"))

const scopedElements = (
  model: QualityModel,
  areaId: string,
  factorFilters: ReadonlyArray<string>,
) => {
  const root = projectModel(model)
  const scoped = findElement(root, areaId)!
  const elements = flattenModel(scoped)
  if (factorFilters.length === 0) return elements
  const selectedFactors = new Set(
    elements
      .filter(
        (element) =>
          element.kind === "factor" &&
          factorFilters.some(
            (filter) => element.id === filter || element.id.startsWith(`${filter}/`),
          ),
      )
      .map((element) => element.id),
  )
  return elements.filter(
    (element) =>
      element.kind === "area" ||
      selectedFactors.has(element.id) ||
      (element.kind === "requirement" &&
        element.parentId !== undefined &&
        selectedFactors.has(element.parentId)),
  )
}

const nextRunNumber = (directory: string) =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const paths = yield* Path.Path
    if (!(yield* fs.exists(directory))) return 1
    let maximum = 0
    for (const name of yield* fs.readDirectory(directory)) {
      let number: number | undefined
      for (const manifestPath of ["evaluation.json", "data/evaluation-manifest.json"]) {
        const file = paths.join(directory, name, manifestPath)
        if (!(yield* fs.exists(file))) continue
        try {
          const parsed = JSON.parse(yield* fs.readFileString(file)) as {
            readonly manifest?: { readonly run?: { readonly number?: number } }
            readonly run?: { readonly number?: number }
          }
          number = parsed.manifest?.run?.number ?? parsed.run?.number
        } catch {
          // Fall back to a current-format folder name.
        }
        break
      }
      if (number === undefined) {
        const match = /^(\d{4})-([a-z0-9-]+)-eval$/.exec(name)
        if (match !== null && !match[2]!.split("-").includes("quality")) {
          number = Number(match[1])
        }
      }
      maximum = Math.max(maximum, number ?? 0)
    }
    return maximum + 1
  })

const dryRun = (input: EvaluationRunInput) =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const runtime = yield* HostRuntime
    const workspace = yield* resolveWorkspace({
      ...(input.model === "" ? {} : { model: input.model }),
      ...(input.evaluationDir === "" ? {} : { evaluationDir: input.evaluationDir }),
    })
    const raw = yield* fs.readFileString(workspace.model.abs)
    const model = decodeModel(parseQualityDocument(workspace.model.abs, raw))
    const scope = resolveScope(model, input.area, input.factors)
    const selected = selectEvaluator(input.evaluator, workspace, evaluatorDiscovery(runtime))
    const elements = scopedElements(
      model,
      scope.plannedScope.areaId,
      scope.plannedScope.factorFilter,
    )
    const areas = elements.filter((entry) => entry.kind === "area")
    const factors = elements.filter((entry) => entry.kind === "factor")
    const requirements = elements.filter((entry) => entry.kind === "requirement")
    const sources = []
    for (const area of areas) {
      const selector = effectiveSource(model, areaPath(area.id)).selector
      const kind = yield* detectSourceKind(workspace.workspaceRoot.abs, selector)
      yield* validateSourceSelector(workspace.workspaceRoot.abs, { selector, kind })
      sources.push({
        area: area.id,
        selector,
        kind,
      })
    }
    const total = 1 + areas.length * 3 + factors.length * 2 + requirements.length * 2 + 4
    const evaluatorUnits = areas.length + factors.length + requirements.length + 3
    const configured = workspace.evaluation.concurrency
    const concurrency = configured ?? runtime.hardwareConcurrency * 2
    if (concurrency < 1) throw new Error("evaluation.concurrency must be a positive integer")
    const resolvedConcurrency =
      concurrency > 1 && !selected.capabilities.concurrent && !selected.capabilities.subagents
        ? 1
        : concurrency
    const number = yield* nextRunNumber(workspace.evaluations.abs)
    const label = `${String(number).padStart(4, "0")}-${scopeSlug(scope.plannedScope)}-eval`
    let command = "qualitymd evaluation run"
    if (input.model !== "") command += ` --model ${input.model}`
    if (input.area !== "") command += ` --area ${input.area}`
    for (const factor of input.factors) command += ` --factor ${factor}`
    if (input.evaluator !== "") command += ` --evaluator ${input.evaluator}`
    return {
      schemaVersion: 3,
      model: workspace.model.rel,
      requestedScope: scope.requestedScope,
      plannedScope: scope.plannedScope,
      evaluator: selected.name,
      evaluatorKind: selected.kind,
      evaluatorReason: selected.reason,
      evaluatorCapabilities: selected.capabilities,
      ...("candidates" in selected ? { evaluatorCandidates: selected.candidates } : {}),
      inspectionPolicy: {
        workspace: "read-only",
        network: "disabled",
        approvals: "never",
        verification: "unavailable",
        repositoryInstructions: "untrusted-data",
      },
      concurrency: resolvedConcurrency,
      workUnits: { total, evaluatorUnits, completed: 0 },
      sources,
      expectedRunPath: `${workspace.evaluations.rel}/${label}`,
      nextActions: [{ id: "evaluation-run", label: "Execute the evaluation run", command }],
    }
  })

export const evaluationRunCommand = (
  input: EvaluationRunInput,
): Effect.Effect<CommandResult, never, FileSystem.FileSystem | Path.Path | HostRuntime> => {
  let operation: Effect.Effect<
    CommandResult,
    FileSystemFailure,
    FileSystem.FileSystem | Path.Path | HostRuntime
  >
  if (!input.dryRun) {
    if (input.evaluatorResult !== "" && input.resume === "") {
      return Effect.succeed(usage("--evaluator-result requires --resume"))
    }
    operation =
      input.resume !== ""
        ? resumeRun(input)
        : resolveWorkspace({
            ...(input.model === "" ? {} : { model: input.model }),
            ...(input.evaluationDir === "" ? {} : { evaluationDir: input.evaluationDir }),
          }).pipe(
            Effect.flatMap((workspace) => {
              return Effect.gen(function* () {
                const runtime = yield* HostRuntime
                const selected = selectEvaluator(
                  input.evaluator,
                  workspace,
                  evaluatorDiscovery(runtime),
                )
                return yield* selected.kind === "harness"
                  ? executeHarnessRun(input)
                  : executeProviderRun(input, selected)
              })
            }),
          )
  } else if (input.resume !== "" || input.evaluatorResult !== "") {
    return Effect.succeed(usage("--dry-run cannot be combined with --resume or --evaluator-result"))
  } else {
    operation = dryRun(input).pipe(
      Effect.map((preview) => {
        if (input.json) return commandResult(jsonDocument(preview))
        return commandResult("", {
          stderr:
            `Would evaluate ${preview.model} at ${preview.expectedRunPath}\n` +
            `Evaluator: ${preview.evaluator} (${preview.evaluatorReason}); concurrency: ${preview.concurrency}; work units: ${preview.workUnits.total} (${preview.workUnits.evaluatorUnits} evaluator-backed)\n\n` +
            `Next: ${preview.nextActions[0]!.command}\n`,
        })
      }),
      Effect.mapError((cause) =>
        cause instanceof FileSystemFailure
          ? cause
          : new FileSystemFailure({
              detail: cause instanceof Error ? cause.message : String(cause),
            }),
      ),
    )
  }
  return operation.pipe(
    Effect.catchCause((cause) => {
      const squashed = Cause.squash(cause)
      const detail = squashed instanceof Error ? squashed.message : String(squashed)
      const category = detail.includes("no evaluator is available")
        ? "missing_evaluator"
        : "internal_error"
      const failure = { category, detail }
      return Effect.succeed(
        input.json
          ? commandResult(jsonDocument({ schemaVersion: 3, status: "failed", failure }), {
              exitCode: ExitCode.problems,
            })
          : commandResult("", {
              stderr: `qualitymd: ${category}: ${detail}\n`,
              exitCode: ExitCode.problems,
            }),
      )
    }),
  )
}
