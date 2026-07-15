import { describe, expect, it } from "vitest"

import { selectEvaluator, type EvaluatorDiscovery } from "../../src/application/evaluation-run.ts"
import type { Workspace } from "../../src/services/workspace.ts"

const workspace = (evaluators: Workspace["evaluators"] = {}, selected?: string): Workspace =>
  ({
    evaluation: selected === undefined ? {} : { evaluator: selected },
    evaluators,
  }) as Workspace

const discovery = (
  commands: ReadonlyArray<string>,
  authenticated: boolean,
  environment: Readonly<Record<string, string | undefined>> = {},
): EvaluatorDiscovery => ({
  which: (command) => (commands.includes(command) ? `/bin/${command}` : null),
  codexAuthenticated: () => authenticated,
  environment,
})

describe("evaluator selection", () => {
  it("prefers an authenticated Codex agent runtime", () => {
    const selected = selectEvaluator("auto", workspace(), discovery(["codex", "claude"], true))
    expect(selected.name).toBe("codex")
    expect("candidates" in selected && selected.candidates).toHaveLength(1)
  })

  it("skips an unauthenticated Codex runtime and records Claude's auth assumption", () => {
    const selected = selectEvaluator("auto", workspace(), discovery(["codex", "claude"], false))
    expect(selected.name).toBe("claude")
    expect("candidates" in selected && selected.candidates[0]?.usable).toBe(false)
    expect("candidates" in selected && selected.candidates[1]?.evidence.join(" ")).toContain(
      "assumed",
    )
  })

  it("selects ready API profiles in alphabetical order without exposing keys", () => {
    const selected = selectEvaluator(
      "auto",
      workspace({
        zeta: { kind: "openai", apiKeyEnv: "ZETA_KEY" },
        alpha: { kind: "anthropic", apiKeyEnv: "ALPHA_KEY" },
      }),
      discovery([], false, { ALPHA_KEY: "secret", ZETA_KEY: "secret" }),
    )
    expect(selected.name).toBe("alpha")
    expect(JSON.stringify(selected)).not.toContain("secret")
  })

  it("does not allow configured profiles to shadow reserved names", () => {
    const selected = selectEvaluator(
      "openai",
      workspace({ openai: { kind: "anthropic" } }),
      discovery([], false),
    )
    expect(selected.kind).toBe("openai")
  })

  it("reports the default Claude adapter's sequential, no-subagent execution boundary", () => {
    const selected = selectEvaluator("claude", workspace(), discovery(["claude"], false))
    expect(selected.capabilities.concurrent).toBe(false)
    expect(selected.capabilities.subagents).toBe(false)
    expect(selected.capabilities.executableOverride).toBe(true)
  })

  it("gives actionable remedies when discovery finds nothing usable", () => {
    expect(() => selectEvaluator("auto", workspace(), discovery([], false))).toThrow(
      "--evaluator harness",
    )
  })
})
