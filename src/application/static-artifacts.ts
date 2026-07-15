import * as Effect from "effect/Effect"

import specification from "../../SPECIFICATION.md" with { type: "text" }
import schemaAsset from "../../quality.schema.json" with { type: "text" }
import { commandResult } from "../domain/command-result.ts"

const schema = schemaAsset as unknown as string

export const specCommand = Effect.sync(() => commandResult(specification))

export const schemaCommand = Effect.sync(() => commandResult(schema))
