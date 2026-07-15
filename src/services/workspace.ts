import * as Effect from "effect/Effect"
import * as FileSystem from "effect/FileSystem"
import * as Path from "effect/Path"
import { parse as parseYaml } from "yaml"

import { FileSystemFailure } from "../domain/errors.ts"
import { mapEntry, nodeValue, parseQualityDocument } from "../domain/model/document.ts"

export interface PathRef {
  readonly abs: string
  readonly rel: string
  readonly repoRel: string
}

export interface EvaluatorProfile {
  readonly kind: string
  readonly model?: string
  readonly command?: string
}

export interface Workspace {
  readonly repoRoot: PathRef
  readonly workspaceRoot: PathRef
  readonly model: PathRef
  readonly config: PathRef
  readonly configPresent: boolean
  readonly dataDir: PathRef
  readonly evaluations: PathRef
  readonly changelog: PathRef
  readonly logs: PathRef
  readonly evaluation: { readonly evaluator?: string; readonly concurrency?: number }
  readonly evaluators: Readonly<Record<string, EvaluatorProfile>>
}

export interface WorkspaceOptions {
  readonly repoRoot?: string
  readonly model?: string
  readonly evaluationDir?: string
}

const slash = (value: string) => value.replaceAll("\\", "/")

const fail = (cause: unknown) =>
  new FileSystemFailure({ detail: cause instanceof Error ? cause.message : String(cause) })

const cleanModelRelative = (paths: Path.Path, value: string) => {
  if (value.trim() === "") throw new Error("path must be non-empty")
  if (paths.isAbsolute(value))
    throw new Error(`path ${JSON.stringify(value)} must be model-relative`)
  return slash(paths.normalize(value))
}

const repoRelative = (paths: Path.Path, root: string, absolute: string) => {
  const relative = paths.relative(root, absolute)
  if (relative === ".." || relative.startsWith(`..${paths.sep}`)) {
    throw new Error(`path ${JSON.stringify(absolute)} escapes the repository`)
  }
  return relative === "" ? "." : slash(relative)
}

const findRepoRoot = (start: string) =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const paths = yield* Path.Path
    let current = paths.resolve(start)
    if (!(yield* fs.exists(current)) || (yield* fs.stat(current)).type !== "Directory") {
      current = paths.dirname(current)
    }
    while (true) {
      if (yield* fs.exists(paths.join(current, ".git"))) return current
      const parent = paths.dirname(current)
      if (parent === current) throw new Error(`could not find repository root from ${start}`)
      current = parent
    }
  })

const resolveRef = (
  paths: Path.Path,
  repoRoot: string,
  workspaceRoot: string,
  value: string,
): PathRef => {
  const rel = cleanModelRelative(paths, value)
  const abs = paths.join(workspaceRoot, rel)
  return { abs, rel, repoRel: repoRelative(paths, repoRoot, abs) }
}

export const resolveWorkspace = (options: WorkspaceOptions = {}) =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const paths = yield* Path.Path
    const modelInput = options.model ?? "QUALITY.md"
    const modelAbs = paths.resolve(options.repoRoot ?? process.cwd(), modelInput)
    const repoRoot =
      options.repoRoot === undefined
        ? yield* findRepoRoot(modelAbs)
        : paths.resolve(options.repoRoot)
    const workspaceRoot = paths.dirname(modelAbs)
    const modelRaw = yield* fs.readFileString(modelAbs)
    const document = parseQualityDocument(modelAbs, modelRaw)
    const configPair = mapEntry(document.frontmatter, "config")
    const configValue =
      configPair === undefined ? ".quality/config.yaml" : nodeValue(configPair.value)
    const config = resolveRef(paths, repoRoot, workspaceRoot, configValue)
    const configPresent = yield* fs.exists(config.abs)
    const parsed = configPresent
      ? ((parseYaml(yield* fs.readFileString(config.abs)) ?? {}) as {
          readonly evaluationDir?: string
          readonly evaluation?: { readonly evaluator?: string; readonly concurrency?: number }
          readonly evaluators?: Readonly<Record<string, EvaluatorProfile>>
        })
      : {}
    const evaluationDir = options.evaluationDir ?? parsed.evaluationDir ?? ".quality/evaluations"
    return {
      repoRoot: { abs: repoRoot, rel: ".", repoRel: "." },
      workspaceRoot: {
        abs: workspaceRoot,
        rel: ".",
        repoRel: repoRelative(paths, repoRoot, workspaceRoot),
      },
      model: {
        abs: modelAbs,
        rel: "QUALITY.md",
        repoRel: repoRelative(paths, repoRoot, modelAbs),
      },
      config,
      configPresent,
      dataDir: resolveRef(paths, repoRoot, workspaceRoot, ".quality"),
      evaluations: resolveRef(paths, repoRoot, workspaceRoot, evaluationDir),
      changelog: resolveRef(paths, repoRoot, workspaceRoot, ".quality/changelog"),
      logs: resolveRef(paths, repoRoot, workspaceRoot, ".quality/logs"),
      evaluation: parsed.evaluation ?? {},
      evaluators: parsed.evaluators ?? {},
    } satisfies Workspace
  }).pipe(Effect.mapError(fail))
