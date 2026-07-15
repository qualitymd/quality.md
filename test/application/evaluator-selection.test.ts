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
): EvaluatorDiscovery => ({
  which: (command) => (commands.includes(command) ? `/bin/${command}` : null),
  codexAuthenticated: () => authenticated,
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

  it("does not turn configured agent profiles into an implicit auto fallback", () => {
    expect(() =>
      selectEvaluator(
        "auto",
        workspace({ team: { kind: "claude", command: "/opt/team-claude" } }),
        discovery([], false),
      ),
    ).toThrow("no evaluator is available")
  })

  it("supports explicitly configured agent-runtime profiles", () => {
    const selected = selectEvaluator(
      "team-agent",
      workspace({ "team-agent": { kind: "claude", command: "/opt/team-claude" } }),
      discovery([], false),
    )
    expect(selected.kind).toBe("claude")
    expect(selected.name).toBe("team-agent")
    expect(selected.capabilities.dispatch).toEqual({
      concurrentCalls: false,
      delegatedRequests: false,
      automaticConcurrency: 1,
      maxConcurrency: 1,
    })
  })

  it("rejects direct API profile kinds", () => {
    expect(() =>
      selectEvaluator(
        "legacy-api",
        workspace({ "legacy-api": { kind: "openai" } }),
        discovery([], false),
      ),
    ).toThrow("use codex or claude")
  })

  it("reports the default Claude adapter's sequential, no-subagent execution boundary", () => {
    const selected = selectEvaluator("claude", workspace(), discovery(["claude"], false))
    expect(selected.capabilities.dispatch).toEqual({
      concurrentCalls: false,
      delegatedRequests: false,
      automaticConcurrency: 1,
      maxConcurrency: 1,
    })
    expect(selected.capabilities.executableOverride).toBe(true)
    expect(selected.capabilities.workspaceInspection).toBe(true)
    expect(selected.capabilities.instructionIsolation).toBe(true)
    expect(selected.capabilities.verification).toBe(false)
  })

  it("gives actionable remedies when discovery finds nothing usable", () => {
    expect(() => selectEvaluator("auto", workspace(), discovery([], false))).toThrow(
      "--evaluator harness",
    )
  })
})
