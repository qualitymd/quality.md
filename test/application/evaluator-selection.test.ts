import { describe, expect, it } from "vitest"

import {
  reportEvaluatorSelection,
  selectEvaluator,
  type EvaluatorDiscovery,
} from "../../src/application/evaluation-run.ts"
import { commandResult } from "../../src/domain/command-result.ts"
import type { Workspace } from "../../src/services/workspace.ts"

const workspace = (evaluators: Workspace["evaluators"] = {}, selected?: string): Workspace =>
  ({
    evaluation: selected === undefined ? {} : { evaluator: selected },
    evaluators,
  }) as Workspace

const discovery = (
  commands: ReadonlyArray<string>,
  codexAuthenticated: boolean,
  claudeAuthenticated: boolean | null = null,
): EvaluatorDiscovery => ({
  which: (command) => (commands.includes(command) ? `/bin/${command}` : null),
  codexAuthenticated: () => codexAuthenticated,
  claudeAuthenticated: () => claudeAuthenticated,
})

describe("evaluator selection", () => {
  it("prefers an authenticated Codex agent runtime", () => {
    const selected = selectEvaluator("auto", workspace(), discovery(["codex", "claude"], true))
    expect(selected.name).toBe("codex")
    expect("candidates" in selected && selected.candidates).toHaveLength(2)
    expect("candidates" in selected && selected.candidates.map(({ name }) => name)).toEqual([
      "codex",
      "claude",
    ])
  })

  it("names deterministic ordering and each usable candidate not selected", () => {
    const selected = selectEvaluator(
      "auto",
      workspace(),
      discovery(["codex", "claude"], true, true),
    )
    expect(selected.name).toBe("codex")
    expect(selected.reason).toContain("deterministic discovery order decided")
    expect(selected.reason).toContain("usable but not selected: claude")
  })

  it("reports verified and assumed authentication bases as structured data", () => {
    const selected = selectEvaluator("auto", workspace(), discovery(["codex", "claude"], false))
    expect(selected.name).toBe("claude")
    expect("candidates" in selected && selected.candidates).toMatchObject([
      { name: "codex", authenticated: false, authenticationBasis: "verified", usable: false },
      { name: "claude", authenticated: true, authenticationBasis: "assumed", usable: true },
    ])
  })

  it("gates Claude usability on its verified authentication probe", () => {
    expect(() => selectEvaluator("auto", workspace(), discovery(["claude"], false, false))).toThrow(
      "no evaluator is available",
    )
    const selected = selectEvaluator("auto", workspace(), discovery(["claude"], false, true))
    expect(selected.name).toBe("claude")
    expect("candidates" in selected && selected.candidates[1]).toMatchObject({
      authenticated: true,
      authenticationBasis: "verified",
      usable: true,
    })
  })

  it("adds auto candidates and the selection reason to run receipts", () => {
    const selected = selectEvaluator(
      "auto",
      workspace(),
      discovery(["codex", "claude"], true, true),
    )
    const result = reportEvaluatorSelection(
      commandResult(`${JSON.stringify({ schemaVersion: 3, status: "completed" })}\n`),
      selected,
    )
    expect(JSON.parse(result.stdout)).toMatchObject({
      status: "completed",
      evaluatorReason: expect.stringContaining("usable but not selected: claude"),
      evaluatorCandidates: [
        { name: "codex", authenticationBasis: "verified", usable: true },
        { name: "claude", authenticationBasis: "verified", usable: true },
      ],
    })
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
