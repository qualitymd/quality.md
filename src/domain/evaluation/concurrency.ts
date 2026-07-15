import type { EvaluatorDispatchCapability } from "../evaluator/types.ts"

export type DispatchMode = "direct" | "delegated" | "sequential"

export interface ConcurrencyResolution {
  readonly source: "configured" | "automatic"
  readonly requested: number
  readonly automatic: number
  readonly maximum?: number
  readonly resolved: number
  readonly clamped: boolean
  readonly mode: DispatchMode
}

const positiveInteger = (value: unknown, label: string): number => {
  if (typeof value !== "number" || !Number.isInteger(value) || value <= 0)
    throw new Error(`${label} must be a positive integer`)
  return value
}

export const validateDispatchCapability = (
  dispatch: EvaluatorDispatchCapability,
): EvaluatorDispatchCapability => {
  const automaticConcurrency = positiveInteger(
    dispatch.automaticConcurrency,
    "evaluator automaticConcurrency",
  )
  const maxConcurrency =
    dispatch.maxConcurrency === undefined
      ? undefined
      : positiveInteger(dispatch.maxConcurrency, "evaluator maxConcurrency")
  if (maxConcurrency !== undefined && automaticConcurrency > maxConcurrency)
    throw new Error("evaluator automaticConcurrency must not exceed maxConcurrency")
  if (!dispatch.concurrentCalls && !dispatch.delegatedRequests) {
    if (automaticConcurrency !== 1 || maxConcurrency !== 1)
      throw new Error(
        "a sequential evaluator must declare automaticConcurrency 1 and maxConcurrency 1",
      )
  }
  return dispatch
}

export const dispatchMode = (dispatch: EvaluatorDispatchCapability): DispatchMode =>
  dispatch.concurrentCalls ? "direct" : dispatch.delegatedRequests ? "delegated" : "sequential"

export const resolveConcurrency = (
  configured: unknown,
  dispatch: EvaluatorDispatchCapability,
): ConcurrencyResolution => {
  validateDispatchCapability(dispatch)
  const source = configured === undefined ? "automatic" : "configured"
  const requested =
    configured === undefined
      ? dispatch.automaticConcurrency
      : positiveInteger(configured, "evaluation.concurrency")
  const maximum = dispatch.maxConcurrency
  const mode = dispatchMode(dispatch)
  const resolved =
    mode === "sequential" ? 1 : maximum === undefined ? requested : Math.min(requested, maximum)
  return {
    source,
    requested,
    automatic: dispatch.automaticConcurrency,
    ...(maximum === undefined ? {} : { maximum }),
    resolved,
    clamped: resolved !== requested,
    mode,
  }
}
