import * as Context from "effect/Context"
import * as Effect from "effect/Effect"
import * as Layer from "effect/Layer"

import { InternalFailure } from "../domain/errors.ts"

export interface OutputService {
  readonly stdout: (content: string) => Effect.Effect<void, InternalFailure>
  readonly stderr: (content: string) => Effect.Effect<void, InternalFailure>
  readonly setExitCode: (code: number) => Effect.Effect<void>
}

export class Output extends Context.Service<Output, OutputService>()("qualitymd/Output") {}

const write = (stream: NodeJS.WriteStream, content: string) =>
  Effect.try({
    try: () => {
      stream.write(content)
    },
    catch: (cause) =>
      new InternalFailure({
        detail: cause instanceof Error ? cause.message : String(cause),
      }),
  })

export const OutputLive = Layer.succeed(Output, {
  stdout: (content) => write(process.stdout, content),
  stderr: (content) => write(process.stderr, content),
  setExitCode: (code) => Effect.sync(() => void (process.exitCode = code)),
} satisfies OutputService)
