import { BunRuntime, BunServices } from "@effect/platform-bun"
import * as Effect from "effect/Effect"
import * as Layer from "effect/Layer"

import { runCli } from "./cli/app.ts"
import { HostRuntimeLive } from "./services/host-runtime.ts"
import { OutputLive } from "./services/output.ts"
import { UpdateRuntimeLive } from "./services/update-runtime.ts"

runCli(process.argv.slice(2)).pipe(
  Effect.provide(Layer.mergeAll(BunServices.layer, HostRuntimeLive, OutputLive, UpdateRuntimeLive)),
  BunRuntime.runMain({ disableErrorReporting: true }),
)
