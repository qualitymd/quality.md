import * as Effect from "effect/Effect"
import * as FileSystem from "effect/FileSystem"
import * as Option from "effect/Option"
import * as Path from "effect/Path"

import { FileSystemFailure } from "../domain/errors.ts"
import { commandResult, ExitCode, type CommandResult } from "../domain/command-result.ts"
import { jsonDocument } from "../domain/json.ts"
import { invalidDocumentResult, lintDocument } from "../domain/lint/lint.ts"
import type { LintResult, RepairRecord } from "../domain/lint/result.ts"
import {
  parseQualityDocument,
  QualityDocumentParseError,
  renderQualityDocument,
} from "../domain/model/document.ts"

export interface LintCommandInput {
  readonly path: string
  readonly json: boolean
  readonly fix: boolean
}

const failure = (operation: string, path: string, cause: unknown) =>
  new FileSystemFailure({
    detail: `${operation} ${path}: ${cause instanceof Error ? cause.message : String(cause)}`,
  })

const readDocument = (path: string) =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const raw = yield* fs
      .readFileString(path)
      .pipe(Effect.mapError((cause) => failure("read", path, cause)))
    return yield* Effect.try({
      try: () => parseQualityDocument(path, raw),
      catch: (cause) =>
        cause instanceof QualityDocumentParseError
          ? cause
          : new FileSystemFailure({ detail: String(cause) }),
    })
  })

const writeAtomic = (path: string, content: string) =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const paths = yield* Path.Path
    const info = yield* fs.stat(path)
    const temp = yield* fs.makeTempFile({
      directory: paths.dirname(path),
      prefix: `.${paths.basename(path)}.`,
    })
    yield* Effect.gen(function* () {
      yield* fs.writeFileString(temp, content, { mode: info.mode })
      yield* fs.chmod(temp, info.mode)
      yield* fs.rename(temp, path)
    }).pipe(Effect.onError(() => Effect.ignore(fs.remove(temp, { force: true }))))
  }).pipe(Effect.mapError((cause) => failure("write", path, cause)))

const check = (path: string) =>
  readDocument(path).pipe(
    Effect.map((document) => lintDocument(document).result),
    Effect.catchIf(
      (cause): cause is QualityDocumentParseError => cause instanceof QualityDocumentParseError,
      () => Effect.succeed(invalidDocumentResult(path)),
    ),
  )

const fix = (path: string) =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const document = yield* readDocument(path)
    const link = yield* fs.readLink(path).pipe(Effect.option)
    if (Option.isSome(link)) {
      return yield* Effect.fail(
        new FileSystemFailure({ detail: `${path} is a symbolic link; refusing to repair it` }),
      )
    }
    const original = lintDocument(document)
    const repairs = original.applyRepairs()
    if (repairs.length === 0) return original.result
    const rendered = renderQualityDocument(document)
    yield* writeAtomic(path, rendered)
    const repaired = yield* readDocument(path)
    return withRepairs(lintDocument(repaired).result, repairs)
  }).pipe(
    Effect.catchIf(
      (cause): cause is QualityDocumentParseError => cause instanceof QualityDocumentParseError,
      () => Effect.succeed(invalidDocumentResult(path)),
    ),
  )

const withRepairs = (result: LintResult, repairs: ReadonlyArray<RepairRecord>): LintResult => ({
  ...result,
  summary: { ...result.summary, fixed: repairs.length },
  repairs,
})

const renderHuman = (result: LintResult) => {
  let stdout = result.summary.fixed > 0 ? `Applied ${result.summary.fixed} repair(s).\n` : ""
  if (result.findings.length === 0) {
    stdout += `${result.path} is valid.\n`
  } else {
    for (const finding of result.findings) {
      stdout += `${finding.severity} ${finding.ruleId}: ${finding.message} (${finding.location.label})\n`
    }
    stdout += `\n${result.summary.errors} error(s), ${result.summary.warnings} warning(s).\n`
  }
  const action = result.nextActions.at(0)
  const stderr = action === undefined || action.command === "" ? "" : `\nNext: ${action.command}\n`
  return { stdout, stderr }
}

export const lintCommand = (
  input: LintCommandInput,
): Effect.Effect<CommandResult, FileSystemFailure, FileSystem.FileSystem | Path.Path> => {
  if (input.path === "-") {
    return Effect.succeed(
      commandResult("", {
        stderr: "qualitymd: lint does not read from stdin yet; pass a file path\n",
        exitCode: ExitCode.usage,
      }),
    )
  }
  return (input.fix ? fix(input.path) : check(input.path)).pipe(
    Effect.map((result) => {
      const exitCode = result.valid ? ExitCode.ok : ExitCode.problems
      if (input.json) return commandResult(jsonDocument(result), { exitCode })
      const rendered = renderHuman(result)
      return commandResult(rendered.stdout, { stderr: rendered.stderr, exitCode })
    }),
  )
}
