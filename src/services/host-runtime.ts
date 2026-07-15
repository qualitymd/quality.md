import * as Context from "effect/Context"
import * as Effect from "effect/Effect"
import * as Layer from "effect/Layer"

export interface HostRuntimeService {
  readonly cwd: string
  readonly environment: Readonly<Record<string, string | undefined>>
  readonly currentTimeMillis: Effect.Effect<number>
  readonly randomBytes: (length: number) => Effect.Effect<Uint8Array>
  readonly readStdin: Effect.Effect<string>
  readonly which: (command: string) => string | null
  readonly codexAuthenticated: () => boolean
  readonly claudeAuthenticated: () => boolean | null
}

export class HostRuntime extends Context.Service<HostRuntime, HostRuntimeService>()(
  "qualitymd/HostRuntime",
) {}

export const HostRuntimeLive = Layer.succeed(HostRuntime, {
  cwd: process.cwd(),
  environment: process.env,
  currentTimeMillis: Effect.sync(() => Date.now()),
  randomBytes: (length) => Effect.sync(() => crypto.getRandomValues(new Uint8Array(length))),
  readStdin: Effect.promise(() => Bun.stdin.text()),
  which: (command) => Bun.which(command),
  codexAuthenticated: () => {
    const result = Bun.spawnSync(["codex", "login", "status"], {
      stdin: "ignore",
      stdout: "ignore",
      stderr: "ignore",
      env: process.env,
    })
    return result.exitCode === 0
  },
  claudeAuthenticated: () => {
    try {
      const result = Bun.spawnSync(["claude", "auth", "status", "--json"], {
        stdin: "ignore",
        stdout: "pipe",
        stderr: "ignore",
        env: process.env,
      })
      const status = JSON.parse(new TextDecoder().decode(result.stdout)) as {
        readonly loggedIn?: unknown
      }
      return typeof status.loggedIn === "boolean" ? status.loggedIn : null
    } catch {
      return null
    }
  },
} satisfies HostRuntimeService)
