import * as Effect from "effect/Effect"
import * as FileSystem from "effect/FileSystem"
import * as Option from "effect/Option"
import * as Path from "effect/Path"

import { FileSystemFailure } from "../domain/errors.ts"
import type { SourceBundle, SourceFile } from "../domain/evaluator/types.ts"
export type { SourceBundle } from "../domain/evaluator/types.ts"

export class SourceCaptureError extends Error {}

const skipped = new Set([".git", ".quality", "node_modules", "vendor", "dist"])
const maxFileBytes = 64 * 1024
const maxBundleBytes = 512 * 1024

const validGlob = (value: string) => {
  for (let index = 0; index < value.length; index += 1) {
    if (value[index] === "\\") {
      index += 1
      if (index >= value.length) return false
      continue
    }
    if (value[index] !== "[") continue
    let cursor = index + 1
    if (value[cursor] === "!" || value[cursor] === "^") cursor += 1
    const first = cursor
    let closed = false
    for (; cursor < value.length; cursor += 1) {
      if (value[cursor] === "\\") {
        cursor += 1
        if (cursor >= value.length) return false
      } else if (value[cursor] === "]") {
        closed = true
        break
      }
    }
    if (!closed || cursor === first) return false
    index = cursor
  }
  return true
}

export const detectSourceKind = (workspaceRoot: string, selector: string) =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const paths = yield* Path.Path
    const value = selector.trim()
    const normalized = value.replaceAll("\\", "/")
    if (
      value === "" ||
      paths.isAbsolute(value) ||
      normalized === ".." ||
      normalized.startsWith("../")
    )
      return "path" as const
    if (/[*?[]/.test(value) && validGlob(value)) return "glob" as const
    if (yield* fs.exists(paths.join(workspaceRoot, value))) return "path" as const
    return "prose" as const
  })

const digest = (bytes: Uint8Array) =>
  Effect.tryPromise(async () => {
    const result = await crypto.subtle.digest("SHA-256", Uint8Array.from(bytes))
    return [...new Uint8Array(result)].map((byte) => byte.toString(16).padStart(2, "0")).join("")
  })

const walk = (root: string) =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const paths = yield* Path.Path
    const output: Array<string> = []
    const visit = (directory: string, selectedRoot: string): Effect.Effect<void, unknown> =>
      Effect.gen(function* () {
        for (const name of (yield* fs.readDirectory(directory)).sort()) {
          const absolute = paths.join(directory, name)
          if (Option.isSome(yield* fs.readLink(absolute).pipe(Effect.option))) continue
          const info = yield* fs.stat(absolute)
          if (info.type === "Directory") {
            if (absolute !== selectedRoot && skipped.has(name)) continue
            yield* visit(absolute, selectedRoot)
          } else if (info.type === "File") {
            output.push(absolute)
          }
        }
      })
    yield* visit(root, root)
    return output
  })

const resolveFiles = (workspaceRoot: string, selector: string, kind: "path" | "glob") =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const paths = yield* Path.Path
    if (kind === "glob") {
      const segments = selector.replaceAll("\\", "/").split("/")
      const firstMeta = segments.findIndex((part) => /[*?[]/.test(part))
      const literalParts = firstMeta < 0 ? segments : segments.slice(0, firstMeta)
      const found = yield* Effect.tryPromise(async () => {
        const matches: Array<string> = []
        for await (const relative of new Bun.Glob(selector).scan({
          cwd: workspaceRoot,
          absolute: false,
          dot: true,
          followSymlinks: false,
          onlyFiles: true,
        })) {
          if (
            relative
              .split("/")
              .some((part, index) => skipped.has(part) && literalParts[index] !== part)
          )
            continue
          matches.push(paths.join(workspaceRoot, relative))
        }
        return matches
      })
      return found.sort()
    }
    const absolute = paths.resolve(workspaceRoot, selector)
    const relative = paths.relative(workspaceRoot, absolute)
    if (relative === ".." || relative.startsWith(`..${paths.sep}`)) return []
    if (!(yield* fs.exists(absolute))) return []
    if (Option.isSome(yield* fs.readLink(absolute).pipe(Effect.option))) return []
    const info = yield* fs.stat(absolute)
    if (info.type === "Directory") return yield* walk(absolute)
    return info.type === "File" ? [absolute] : []
  })

export const packageSource = (
  workspaceRoot: string,
  selector: string,
  kind: "path" | "glob",
): Effect.Effect<SourceBundle, FileSystemFailure, FileSystem.FileSystem | Path.Path> =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const paths = yield* Path.Path
    const files: Array<SourceFile> = []
    let total = 0
    let bundleTruncated = false
    const normalized = selector.replaceAll("\\", "/")
    if (
      kind === "path" &&
      (paths.isAbsolute(selector) || normalized === ".." || normalized.startsWith("../"))
    )
      yield* Effect.fail(
        new Error(`source selector ${JSON.stringify(selector)} escapes the workspace`),
      )
    const selected = yield* resolveFiles(workspaceRoot, selector, kind)
    if (selected.length === 0)
      yield* Effect.fail(
        new Error(
          kind === "glob"
            ? `source glob ${JSON.stringify(selector)} matched no readable files`
            : `source path ${JSON.stringify(selector)} contains no readable files`,
        ),
      )
    for (const absolute of selected) {
      if (total >= maxBundleBytes) {
        bundleTruncated = true
        break
      }
      const raw = yield* fs.readFile(absolute)
      if (raw.slice(0, 8 * 1024).includes(0)) continue
      const sha256 = yield* digest(raw)
      let included = raw
      let truncated = false
      if (included.length > maxFileBytes) {
        included = included.slice(0, maxFileBytes)
        truncated = true
      }
      if (total + included.length > maxBundleBytes) {
        included = included.slice(0, maxBundleBytes - total)
        truncated = true
        bundleTruncated = true
      }
      total += included.length
      files.push({
        path: paths.relative(workspaceRoot, absolute).replaceAll("\\", "/"),
        content: new TextDecoder().decode(included),
        sha256,
        ...(truncated ? { truncated: true } : {}),
      })
    }
    if (files.length === 0)
      yield* Effect.fail(
        new Error(`source selector ${JSON.stringify(selector)} contains no readable text files`),
      )
    const hashInput = files.map((file) => `${file.path}\0${file.sha256}\0`).join("")
    return {
      files,
      hash: yield* digest(new TextEncoder().encode(hashInput)),
      truncated: bundleTruncated,
    }
  }).pipe(
    Effect.mapError((cause) => {
      const detail = cause instanceof Error ? cause.message : String(cause)
      return new FileSystemFailure({
        detail: detail.startsWith("source ")
          ? detail
          : `source ${kind} ${JSON.stringify(selector)} could not be packaged: ${detail}`,
      })
    }),
  )

export const captureSource = (
  workspaceRoot: string,
  value: unknown,
): Effect.Effect<SourceBundle, SourceCaptureError, FileSystem.FileSystem | Path.Path> =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const paths = yield* Path.Path
    if (
      value === null ||
      typeof value !== "object" ||
      !Array.isArray((value as { files?: unknown }).files)
    ) {
      yield* Effect.fail(
        new SourceCaptureError("source resolution result must carry a non-empty files array"),
      )
    }
    const result = value as { readonly files: ReadonlyArray<unknown> }
    if (Object.keys(result).some((key) => key !== "files"))
      yield* Effect.fail(new SourceCaptureError("source resolution result must carry only files"))
    const items = result.files
    if (items.length === 0)
      yield* Effect.fail(
        new SourceCaptureError("source resolution result must carry a non-empty files array"),
      )
    const files: Array<SourceFile> = []
    const seen = new Set<string>()
    const decoder = new TextDecoder()
    let total = 0
    let bundleTruncated = false
    for (const [index, value] of items.entries()) {
      if (value === null || Array.isArray(value) || typeof value !== "object")
        yield* Effect.fail(new SourceCaptureError(`files[${index}] must be an object`))
      const item = value as { readonly path?: unknown }
      if (Object.keys(item).some((key) => key !== "path"))
        yield* Effect.fail(new SourceCaptureError(`files[${index}] must carry only path`))
      if (typeof item.path !== "string" || item.path.trim() === "")
        yield* Effect.fail(new SourceCaptureError(`files[${index}] must carry a non-empty path`))
      const normalized = (item.path as string).trim().replaceAll("\\", "/")
      if (
        paths.isAbsolute(normalized) ||
        normalized === ".." ||
        normalized.startsWith("../") ||
        normalized.split("/").includes("..")
      )
        yield* Effect.fail(
          new SourceCaptureError(`files[${index}].path must be workspace-relative and contained`),
        )
      if (seen.has(normalized))
        yield* Effect.fail(
          new SourceCaptureError(`files[${index}] duplicates path ${JSON.stringify(normalized)}`),
        )
      seen.add(normalized)
      const absolute = paths.resolve(workspaceRoot, normalized)
      const relative = paths.relative(workspaceRoot, absolute)
      if (relative === ".." || relative.startsWith(`..${paths.sep}`))
        yield* Effect.fail(
          new SourceCaptureError(`files[${index}].path must be workspace-relative and contained`),
        )
      if (!(yield* fs.exists(absolute)))
        yield* Effect.fail(
          new SourceCaptureError(
            `files[${index}].path ${JSON.stringify(normalized)} does not exist`,
          ),
        )
      if (Option.isSome(yield* fs.readLink(absolute).pipe(Effect.option)))
        yield* Effect.fail(
          new SourceCaptureError(
            `files[${index}].path ${JSON.stringify(normalized)} must not be a symlink`,
          ),
        )
      if ((yield* fs.stat(absolute)).type !== "File")
        yield* Effect.fail(
          new SourceCaptureError(
            `files[${index}].path ${JSON.stringify(normalized)} must name a regular file`,
          ),
        )
      const raw = yield* fs.readFile(absolute)
      if (raw.slice(0, 8 * 1024).includes(0))
        yield* Effect.fail(
          new SourceCaptureError(
            `files[${index}].path ${JSON.stringify(normalized)} must name a readable text file`,
          ),
        )
      const sha256 = yield* digest(raw)
      let included = raw
      let truncated = false
      if (included.length > maxFileBytes) {
        included = included.slice(0, maxFileBytes)
        truncated = true
      }
      if (total >= maxBundleBytes) {
        bundleTruncated = true
        break
      }
      if (total + included.length > maxBundleBytes) {
        included = included.slice(0, maxBundleBytes - total)
        truncated = true
        bundleTruncated = true
      }
      total += included.length
      files.push({
        path: normalized,
        content: decoder.decode(included),
        sha256,
        ...(truncated ? { truncated: true } : {}),
      })
    }
    return {
      files,
      hash: yield* digest(
        new TextEncoder().encode(files.map((file) => `${file.path}\0${file.sha256}\0`).join("")),
      ),
      truncated: bundleTruncated,
    }
  }).pipe(
    Effect.mapError((cause) =>
      cause instanceof SourceCaptureError
        ? cause
        : new SourceCaptureError(cause instanceof Error ? cause.message : String(cause)),
    ),
  )
