import * as Context from "effect/Context"
import type * as Effect from "effect/Effect"
import type * as FileSystem from "effect/FileSystem"

import type {
  EvaluationRequest,
  EvaluationResponse,
  EvaluatorFailure,
  EvaluatorCapabilities,
  EvaluatorKind,
} from "../domain/evaluator/types.ts"

export interface EvaluatorService {
  readonly name: string
  readonly kind: EvaluatorKind
  readonly capabilities: EvaluatorCapabilities
  readonly evaluate: (
    request: EvaluationRequest,
  ) => Effect.Effect<EvaluationResponse, EvaluatorFailure, FileSystem.FileSystem>
}

export class Evaluator extends Context.Service<Evaluator, EvaluatorService>()(
  "qualitymd/Evaluator",
) {}
