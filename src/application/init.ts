import * as Effect from "effect/Effect"
import * as FileSystem from "effect/FileSystem"
import * as Path from "effect/Path"

import skeleton from "../assets/skeleton.md" with { type: "text" }
import minimalSkeleton from "../assets/skeleton-minimal.md" with { type: "text" }
import { FileSystemFailure } from "../domain/errors.ts"
import { commandResult, ExitCode, type CommandResult } from "../domain/command-result.ts"
import { jsonDocument } from "../domain/json.ts"
import { HostRuntime } from "../services/host-runtime.ts"

export interface InitInput {
  readonly path: string
  readonly force: boolean
  readonly json: boolean
  readonly minimal: boolean
  readonly noAgentInstructions: boolean
}

interface AgentInstructionFile {
  readonly path: string
  readonly created: boolean
  readonly updated: boolean
}

const marker = "<!-- Added by qualitymd init. -->"
const pointerTrail = " for this project's quality model."

const errorDetail = (cause: unknown) => (cause instanceof Error ? cause.message : String(cause))

const updateAgentInstructions = (modelPath: string) =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const paths = yield* Path.Path
    const runtime = yield* HostRuntime
    const cwd = runtime.cwd
    const absoluteModel = paths.resolve(cwd, modelPath)
    const candidates = [
      { path: "AGENTS.md", create: true },
      { path: "CLAUDE.md", create: false },
      { path: "GEMINI.md", create: false },
    ]
    const seen = new Set<string>()
    const results: Array<AgentInstructionFile> = []
    for (const candidate of candidates) {
      const target = paths.join(cwd, candidate.path)
      const exists = yield* fs.exists(target)
      if (!exists && !candidate.create) continue
      const key = exists ? yield* fs.realPath(target) : target
      if (seen.has(key)) continue
      seen.add(key)
      const existing = exists ? yield* fs.readFileString(target) : ""
      if (
        existing.includes(marker) ||
        (existing.includes("See [QUALITY.md](") && existing.includes(pointerTrail))
      ) {
        continue
      }
      const relative = paths.relative(paths.dirname(target), absoluteModel).replaceAll("\\", "/")
      const block = `${marker}\nSee [QUALITY.md](${relative})${pointerTrail}`
      const content =
        existing === ""
          ? `${block}\n`
          : existing.endsWith("\n\n")
            ? `${existing}${block}\n`
            : existing.endsWith("\n")
              ? `${existing}\n${block}\n`
              : `${existing}\n\n${block}\n`
      yield* fs.writeFileString(target, content, { mode: 0o644 })
      results.push({ path: candidate.path, created: !exists, updated: exists })
    }
    return results
  }).pipe(Effect.mapError((cause) => new FileSystemFailure({ detail: errorDetail(cause) })))

const actions = (path: string) => [
  {
    id: "lint",
    label: "Validate the scaffolded file",
    command: `qualitymd lint ${path}`,
  },
]

export const initCommand = (
  input: InitInput,
): Effect.Effect<
  CommandResult,
  FileSystemFailure,
  FileSystem.FileSystem | Path.Path | HostRuntime
> => {
  const content = input.minimal ? minimalSkeleton : skeleton
  if (input.json && input.path === "-") {
    return Effect.succeed(
      commandResult("", {
        stderr: "qualitymd: --json cannot be combined with path -\n",
        exitCode: ExitCode.usage,
      }),
    )
  }
  if (input.path === "-") return Effect.succeed(commandResult(content))
  return Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const existed = yield* fs.exists(input.path)
    const created = !existed
    const written = yield* fs
      .writeFileString(input.path, content, {
        flag: input.force ? "w" : "wx",
        mode: 0o644,
      })
      .pipe(Effect.result)
    if (written._tag === "Failure") {
      const reason =
        existed && !input.force
          ? `${input.path} already exists; pass --force to overwrite`
          : errorDetail(written.failure)
      if (input.json) {
        return commandResult("", {
          stderr: jsonDocument({ schemaVersion: 1, path: input.path, reason }),
          exitCode: ExitCode.internal,
        })
      }
      return yield* Effect.fail(new FileSystemFailure({ detail: reason }))
    }
    const agentInstructionFiles = input.noAgentInstructions
      ? []
      : yield* updateAgentInstructions(input.path)
    const nextActions = actions(input.path)
    if (input.json) {
      return commandResult(
        jsonDocument({
          schemaVersion: 1,
          path: input.path,
          created,
          agentInstructionFiles,
          nextActions,
        }),
      )
    }
    const agentLine =
      agentInstructionFiles.length === 0
        ? ""
        : `Agent instructions: ${agentInstructionFiles.map((entry) => entry.path).join(", ")}\n`
    return commandResult("", {
      stderr: `Created ${input.path}\n${agentLine}\nNext: ${nextActions[0]!.command}\n`,
    })
  }).pipe(
    Effect.mapError((cause) =>
      cause instanceof FileSystemFailure
        ? cause
        : new FileSystemFailure({ detail: errorDetail(cause) }),
    ),
  )
}
