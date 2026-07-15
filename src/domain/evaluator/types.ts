import * as Schema from "effect/Schema"

export type EvaluatorKind = "harness" | "codex" | "claude"

export type SourceKind = "path" | "glob" | "prose"

export interface EvaluatorCapabilities {
  readonly structuredOutput: boolean
  readonly workspaceInspection: boolean
  readonly instructionIsolation: boolean
  readonly verification: boolean
  readonly networkAccess: "disabled"
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

export interface InspectionContext {
  readonly workspaceRoot: string
  readonly source: {
    readonly selector: string
    readonly kind: SourceKind
  }
  readonly policy: {
    readonly workspace: "read-only"
    readonly network: "disabled"
    readonly approvals: "never"
    readonly verification: "unavailable" | "sandboxed"
    readonly repositoryInstructions: "untrusted-data"
  }
}

export interface EvaluationRequest {
  readonly runId: string
  readonly workUnitId: string
  readonly kind: string
  readonly subject: string
  readonly instructions: string
  readonly sharedContext: Readonly<Record<string, unknown>>
  readonly context: Readonly<Record<string, unknown>>
  readonly bodyGuidance: string
  readonly inspection?: InspectionContext
  readonly expectedSchema: Readonly<Record<string, unknown>>
  readonly timeoutMs: number
}

export interface EvaluationUsage {
  readonly inputTokens?: number
  readonly outputTokens?: number
  readonly cachedInputTokens?: number
  readonly cacheWriteInputTokens?: number
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
      "rate_limited",
      "timeout",
      "invalid_evaluator_output",
      "schema_invalid_output",
      "evidence_invalid",
      "unsafe_source_content",
      "workspace_access_denied",
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
