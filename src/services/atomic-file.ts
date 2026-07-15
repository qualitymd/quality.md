import * as Effect from "effect/Effect"
import * as FileSystem from "effect/FileSystem"
import * as Path from "effect/Path"

// Atomically replace `path` with `content` by writing a sibling temp file in the
// same directory and renaming it over the destination. Same-directory rename is
// atomic on POSIX, so readers never observe a partial write.
//
// The temp file is a plain sibling (`.<name>.<token>.tmp`), not `fs.makeTempFile`:
// `makeTempFile` creates a *directory* via `mkdtemp` and returns a file inside it,
// so renaming the inner file out leaves the container directory orphaned. A sibling
// file creates nothing extra to leak. On failure the temp file is removed; a
// successful rename consumes it.
export const atomicWriteFileString = Effect.fn("qualitymd.atomicWriteFileString")(function* (
  path: string,
  content: string,
  options: { readonly mode: number },
) {
  const fs = yield* FileSystem.FileSystem
  const paths = yield* Path.Path
  const token = yield* Effect.sync(() =>
    Array.from(crypto.getRandomValues(new Uint8Array(6)), (byte) =>
      byte.toString(16).padStart(2, "0"),
    ).join(""),
  )
  const temp = paths.join(paths.dirname(path), `.${paths.basename(path)}.${token}.tmp`)
  yield* Effect.gen(function* () {
    yield* fs.writeFileString(temp, content, { mode: options.mode })
    yield* fs.chmod(temp, options.mode)
    yield* fs.rename(temp, path)
  }).pipe(Effect.onError(() => Effect.ignore(fs.remove(temp, { force: true }))))
})
