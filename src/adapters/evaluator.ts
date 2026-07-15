import { query, type SDKResultMessage } from "@anthropic-ai/claude-agent-sdk"
import { Codex } from "@openai/codex-sdk"
import * as Effect from "effect/Effect"
import * as FileSystem from "effect/FileSystem"

import { renderEvaluationPrompt } from "../domain/evaluator/context.ts"
import {
  EvaluatorFailure,
  type EvaluationResponse,
  type EvaluatorCapabilities,
} from "../domain/evaluator/types.ts"
import type { EvaluatorService } from "../services/evaluator.ts"

const failure = (cause: unknown, category: EvaluatorFailure["category"] = "internal_error") =>
  new EvaluatorFailure({
    category,
    detail: cause instanceof Error ? cause.message : String(cause),
  })

export const strictProviderSchema = (value: unknown): unknown => {
  if (Array.isArray(value)) return value.map(strictProviderSchema)
  if (value === null || typeof value !== "object") return value
  const input = value as Record<string, unknown>
  const unsupported = new Set(["$schema", "allOf", "if", "then", "else", "not"])
  const schema = Object.fromEntries(
    Object.entries(input)
      .filter(([key]) => key !== "properties" && key !== "required" && !unsupported.has(key))
      .map(([key, child]) => [key, strictProviderSchema(child)]),
  ) as Record<string, unknown>
  if (
    input.properties !== null &&
    typeof input.properties === "object" &&
    !Array.isArray(input.properties)
  ) {
    const properties = input.properties as Record<string, unknown>
    const originallyRequired = new Set(Array.isArray(input.required) ? input.required : [])
    schema.properties = Object.fromEntries(
      Object.entries(properties).map(([name, child]) => {
        const adapted = strictProviderSchema(child)
        return [
          name,
          originallyRequired.has(name) ? adapted : { anyOf: [adapted, { type: "null" }] },
        ]
      }),
    )
    schema.additionalProperties = false
    schema.required = Object.keys(properties)
  }
  return schema
}

export const stripNullProperties = (value: unknown): unknown => {
  if (Array.isArray(value)) return value.map(stripNullProperties)
  if (value === null || typeof value !== "object") return value
  return Object.fromEntries(
    Object.entries(value)
      .filter(([, child]) => child !== null)
      .map(([key, child]) => [key, stripNullProperties(child)]),
  )
}

const freshDirectory = () =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    return yield* fs.makeTempDirectoryScoped({ prefix: "qualitymd-evaluator-" })
  })

const childEnvironment = (provider: "codex" | "claude") => {
  const common = [
    "HOME",
    "PATH",
    "USER",
    "LOGNAME",
    "SHELL",
    "TMPDIR",
    "TEMP",
    "TMP",
    "LANG",
    "LC_ALL",
    "TERM",
  ]
  const providerNames =
    provider === "codex"
      ? ["CODEX_HOME", "OPENAI_API_KEY", "OPENAI_BASE_URL"]
      : ["CLAUDE_CONFIG_DIR", "ANTHROPIC_API_KEY", "ANTHROPIC_BASE_URL"]
  return Object.fromEntries(
    [...common, ...providerNames].flatMap((name) => {
      const value = process.env[name]
      return value === undefined ? [] : [[name, value]]
    }),
  )
}

const codexCapabilities: EvaluatorCapabilities = {
  structuredOutput: true,
  workspaceInspection: true,
  instructionIsolation: true,
  verification: false,
  networkAccess: "disabled",
  tools: true,
  concurrent: true,
  subagents: true,
  freshContext: true,
  cancellation: true,
  usage: true,
  maxTurns: "unsupported",
  tokenBudget: "unsupported",
  costBudget: "unsupported",
  contextWindow: "unknown",
  compaction: "opaque",
  sandbox: "provider",
  executableOverride: true,
}

export const codexEvaluator = (
  options: {
    readonly name?: string
    readonly model?: string
    readonly command?: string
  } = {},
): EvaluatorService => ({
  name: options.name ?? "codex",
  kind: "codex",
  capabilities: codexCapabilities,
  evaluate: (request) =>
    Effect.scoped(
      Effect.gen(function* () {
        const workingDirectory = yield* freshDirectory()
        const prompt = renderEvaluationPrompt(request)
        const codex = new Codex({
          ...(options.command === undefined ? {} : { codexPathOverride: options.command }),
          config: {
            features: { multi_agent: false },
            shell_environment_policy: {
              inherit: "none",
              include_only: ["HOME", "PATH", "TMPDIR", "TEMP", "TMP", "LANG", "LC_ALL"],
            },
          },
          env: childEnvironment("codex"),
        })
        const result = yield* Effect.tryPromise({
          try: (signal) =>
            codex
              .startThread({
                ...(options.model === undefined ? {} : { model: options.model }),
                workingDirectory,
                skipGitRepoCheck: true,
                sandboxMode: "read-only",
                approvalPolicy: "never",
                networkAccessEnabled: false,
                webSearchMode: "disabled",
                ...(request.inspection === undefined
                  ? {}
                  : { additionalDirectories: [request.inspection.workspaceRoot] }),
              })
              .run(prompt, { outputSchema: strictProviderSchema(request.expectedSchema), signal }),
          catch: (cause) => failure(cause),
        }).pipe(
          Effect.timeoutOrElse({
            duration: request.timeoutMs,
            orElse: () => Effect.fail(failure("Codex evaluation timed out", "timeout")),
          }),
        )
        let payload: Readonly<Record<string, unknown>>
        try {
          payload = stripNullProperties(JSON.parse(result.finalResponse)) as Readonly<
            Record<string, unknown>
          >
        } catch (cause) {
          return yield* Effect.fail(failure(cause, "invalid_evaluator_output"))
        }
        return {
          workUnitId: request.workUnitId,
          payload,
          evaluatorKind: "codex",
          ...(options.model === undefined ? {} : { model: options.model }),
          ...(result.usage === null
            ? {}
            : {
                usage: {
                  inputTokens: result.usage.input_tokens,
                  outputTokens: result.usage.output_tokens,
                  cachedInputTokens: result.usage.cached_input_tokens,
                },
              }),
        } satisfies EvaluationResponse
      }),
    ).pipe(
      Effect.mapError((cause) => (cause instanceof EvaluatorFailure ? cause : failure(cause))),
    ),
})

const claudeCapabilities: EvaluatorCapabilities = {
  structuredOutput: true,
  workspaceInspection: true,
  instructionIsolation: true,
  verification: false,
  networkAccess: "disabled",
  tools: true,
  concurrent: false,
  subagents: false,
  freshContext: true,
  cancellation: true,
  usage: true,
  maxTurns: "supported",
  tokenBudget: "supported",
  costBudget: "advisory",
  contextWindow: "reported",
  compaction: "observable",
  sandbox: "host",
  executableOverride: true,
}

export const claudeEvaluator = (
  options: {
    readonly name?: string
    readonly model?: string
    readonly command?: string
  } = {},
): EvaluatorService => ({
  name: options.name ?? "claude",
  kind: "claude",
  capabilities: claudeCapabilities,
  evaluate: (request) =>
    Effect.scoped(
      Effect.gen(function* () {
        const cwd = yield* freshDirectory()
        const prompt = renderEvaluationPrompt(request)
        const result = yield* Effect.tryPromise({
          try: async (signal) => {
            const abortController = new AbortController()
            signal.addEventListener("abort", () => abortController.abort(), { once: true })
            let stderr = ""
            let session: ReturnType<typeof query> | undefined
            try {
              session = query({
                prompt,
                options: {
                  abortController,
                  cwd,
                  ...(request.inspection === undefined
                    ? {}
                    : { additionalDirectories: [request.inspection.workspaceRoot] }),
                  ...(options.model === undefined ? {} : { model: options.model }),
                  ...(options.command === undefined
                    ? {}
                    : { pathToClaudeCodeExecutable: options.command }),
                  maxTurns: 8,
                  outputFormat: {
                    type: "json_schema",
                    schema: strictProviderSchema(request.expectedSchema) as Record<string, unknown>,
                  },
                  permissionMode: "dontAsk",
                  settingSources: [],
                  persistSession: false,
                  env: childEnvironment("claude"),
                  stderr: (data) => {
                    stderr = (stderr + data).slice(-2048)
                  },
                  tools: request.inspection === undefined ? [] : ["Read", "Glob", "Grep"],
                  disallowedTools: ["Write", "Edit", "Bash", "Agent", "Task"],
                },
              })
              let result: SDKResultMessage | undefined
              for await (const message of session) if (message.type === "result") result = message
              if (result === undefined) throw new Error("Claude returned no result message")
              return result
            } catch (cause) {
              const detail = cause instanceof Error ? cause.message : String(cause)
              throw new Error(stderr.trim() === "" ? detail : `${detail}: ${stderr.trim()}`)
            } finally {
              session?.close()
            }
          },
          catch: (cause) => failure(cause),
        }).pipe(
          Effect.timeoutOrElse({
            duration: request.timeoutMs,
            orElse: () => Effect.fail(failure("Claude evaluation timed out", "timeout")),
          }),
        )
        if (result.subtype !== "success") {
          return yield* Effect.fail(
            failure(result.errors.join("; ") || result.subtype, "invalid_evaluator_output"),
          )
        }
        const payload = stripNullProperties(
          result.structured_output ?? JSON.parse(result.result),
        ) as Readonly<Record<string, unknown>>
        return {
          workUnitId: request.workUnitId,
          payload,
          evaluatorKind: "claude",
          ...(options.model === undefined ? {} : { model: options.model }),
          contextMeta: { sessionId: result.session_id },
          usage: {
            inputTokens: result.usage.input_tokens,
            outputTokens: result.usage.output_tokens,
            costUsd: result.total_cost_usd,
          },
        } satisfies EvaluationResponse
      }),
    ).pipe(
      Effect.mapError((cause) => (cause instanceof EvaluatorFailure ? cause : failure(cause))),
    ),
})

export const harnessCapabilities: EvaluatorCapabilities = {
  structuredOutput: true,
  workspaceInspection: true,
  instructionIsolation: true,
  verification: false,
  networkAccess: "disabled",
  tools: true,
  concurrent: false,
  subagents: true,
  freshContext: true,
  cancellation: true,
  usage: true,
  maxTurns: "supported",
  tokenBudget: "supported",
  costBudget: "supported",
  contextWindow: "unknown",
  compaction: "opaque",
  sandbox: "host",
  executableOverride: false,
}
