import * as Schema from "effect/Schema"

export type EvaluatorKind = "harness" | "codex" | "claude" | "openai" | "anthropic"

export interface EvaluatorCapabilities {
  readonly structuredOutput: boolean
  readonly sourceResolution: boolean
  readonly tools: boolean
  readonly concurrent: boolean
  readonly subagents: boolean
  readonly freshContext: boolean
  readonly cancellation: boolean
  readonly usage: boolean
  readonly maxTurns: "supported" | "unsupported"
  readonly tokenBudget: "supported" | "advisory" | "unsupported"
  readonly costBudget: "supported" | "advisory" | "unsupported"
  readonly contextWindow: "reported" | "configured" | "unknown"
  readonly compaction: "observable" | "configurable" | "opaque"
  readonly sandbox: "provider" | "host" | "unsupported"
  readonly executableOverride: boolean
}

export interface SourceFile {
  readonly path: string
  readonly content: string
  readonly sha256: string
  readonly truncated?: boolean
}

export interface SourceBundle {
  readonly files: ReadonlyArray<SourceFile>
  readonly hash: string
  readonly truncated: boolean
}

export interface AreaContext {
  readonly areaId: string
  readonly sourceBundleHash: string
  readonly frame: Readonly<Record<string, unknown>>
  readonly ratingCriteria: Readonly<Record<string, string>>
  readonly bodyGuidance: string
  readonly files: ReadonlyArray<SourceFile>
  readonly hash: string
}

export interface EvaluationRequest {
  readonly runId: string
  readonly workUnitId: string
  readonly kind: string
  readonly subject: string
  readonly instructions: string
  readonly areaContext: AreaContext
  readonly context: Readonly<Record<string, unknown>>
  readonly expectedSchema: Readonly<Record<string, unknown>>
  readonly workspaceRoot: string
  readonly timeoutMs: number
}

export interface EvaluationUsage {
  readonly inputTokens?: number
  readonly outputTokens?: number
  readonly cachedInputTokens?: number
  readonly costUsd?: number
}

export interface EvaluationResponse {
  readonly workUnitId: string
  readonly payload: Readonly<Record<string, unknown>>
  readonly evaluatorKind: EvaluatorKind
  readonly model?: string
  readonly contextMeta?: Readonly<Record<string, string>>
  readonly usage?: EvaluationUsage
}

const FailureTypeId = "qualitymd/evaluator/failure"

export class EvaluatorFailure extends Schema.TaggedErrorClass<EvaluatorFailure>(FailureTypeId)(
  "EvaluatorFailure",
  {
    category: Schema.Literals([
      "missing_evaluator",
      "evaluator_unauthenticated",
      "evaluator_incompatible",
      "missing_api_key",
      "rate_limited",
      "timeout",
      "invalid_evaluator_output",
      "schema_invalid_output",
      "unsafe_source_content",
      "insufficient_evidence",
      "source_unavailable",
      "selector_unsupported",
      "run_state_invalid",
      "cancelled",
      "report_build_failed",
      "internal_error",
    ]),
    detail: Schema.String,
  },
) {
  override get message() {
    return this.detail
  }
}
