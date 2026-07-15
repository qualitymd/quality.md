import * as Effect from "effect/Effect"

import { hashJson } from "../domain/json.ts"

export const hashJsonEffect = Effect.fn("qualitymd.hashJson")((value: unknown) =>
  Effect.promise(() => hashJson(value)),
)

export const requestId = Effect.fn("qualitymd.evaluationRequestId")(function* (
  evaluationId: string,
  workUnit: string,
  inputHash: string,
  attempt: number,
) {
  const hash = yield* hashJsonEffect({ evaluationId, workUnit, inputHash, attempt })
  return `req_${hash.slice(0, 16)}`
})
