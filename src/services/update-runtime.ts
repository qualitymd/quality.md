import * as Context from "effect/Context"
import * as Effect from "effect/Effect"
import * as Layer from "effect/Layer"

export interface UpdateRuntimeService {
  readonly platform: NodeJS.Platform
  readonly arch: string
  readonly execPath: string
  readonly environment: Readonly<Record<string, string | undefined>>
  readonly stderrIsTTY: boolean
  readonly linuxIsMusl: boolean
  readonly request: (url: string, accept: string) => Effect.Effect<string, Error>
  readonly run: (argv: ReadonlyArray<string>) => Effect.Effect<void, Error>
  readonly visibleVersion: Effect.Effect<string, Error>
  readonly spawnRefresh: (environment: Readonly<Record<string, string>>) => Effect.Effect<void>
}

export class UpdateRuntime extends Context.Service<UpdateRuntime, UpdateRuntimeService>()(
  "qualitymd/UpdateRuntime",
) {}

const failure = (cause: unknown) => (cause instanceof Error ? cause : new Error(String(cause)))

const linuxIsMusl = () => {
  if (process.platform !== "linux") return false
  const report = process.report?.getReport?.() as
    | { readonly header?: { readonly glibcVersionRuntime?: string } }
    | undefined
  return report?.header?.glibcVersionRuntime === undefined
}

export const UpdateRuntimeLive = Layer.succeed(UpdateRuntime, {
  platform: process.platform,
  arch: process.arch,
  execPath: process.execPath,
  environment: process.env,
  stderrIsTTY: process.stderr.isTTY,
  linuxIsMusl: linuxIsMusl(),
  request: (url, accept) =>
    Effect.tryPromise({
      try: async (signal) => {
        const response = await fetch(url, {
          headers: { accept, "user-agent": "qualitymd" },
          signal: AbortSignal.any([signal, AbortSignal.timeout(10_000)]),
        })
        if (!response.ok)
          throw new Error(`latest version check failed: ${response.status} ${response.statusText}`)
        return response.text()
      },
      catch: failure,
    }),
  run: (argv) =>
    Effect.tryPromise({
      try: async (signal) => {
        const subprocess = Bun.spawn([...argv], {
          stdin: "inherit",
          stdout: "inherit",
          stderr: "inherit",
        })
        signal.addEventListener("abort", () => subprocess.kill(), { once: true })
        const code = await subprocess.exited
        if (code !== 0) throw new Error(`${argv[0]} exited with status ${code}`)
      },
      catch: failure,
    }),
  visibleVersion: Effect.try({
    try: () => {
      const command = Bun.which("qualitymd") ?? process.execPath
      const result = Bun.spawnSync([command, "--version"], {
        stdout: "pipe",
        stderr: "pipe",
      })
      if (result.exitCode !== 0) throw new Error(new TextDecoder().decode(result.stderr).trim())
      return new TextDecoder().decode(result.stdout).trim()
    },
    catch: failure,
  }),
  spawnRefresh: (environment) =>
    Effect.sync(() => {
      Bun.spawn([process.execPath, "__update-refresh"], {
        env: environment,
        stdin: "ignore",
        stdout: "ignore",
        stderr: "ignore",
      }).unref()
    }),
} satisfies UpdateRuntimeService)
