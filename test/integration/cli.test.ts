import { spawnSync } from "node:child_process"
import { mkdir, mkdtemp, readFile, rm, writeFile } from "node:fs/promises"
import { join } from "node:path"
import { fileURLToPath } from "node:url"
import { describe, expect, it } from "vitest"

import evaluationExamples from "../../src/assets/evaluation-examples.json"

const repositoryRoot = fileURLToPath(new URL("../..", import.meta.url))

const runFrom = (cwd: string, arguments_: ReadonlyArray<string>, input?: string) =>
  spawnSync("bun", ["src/main.ts", ...arguments_], {
    cwd,
    encoding: "utf8",
    env: { ...process.env, QUALITYMD_NO_UPDATE_CHECK: "1", NO_COLOR: "1" },
    ...(input === undefined ? {} : { input }),
  })

const run = (...arguments_: ReadonlyArray<string>) => runFrom(repositoryRoot, arguments_)

const replaceReference = (value: unknown, from: string, to: string): unknown => {
  if (value === from) return to
  if (Array.isArray(value)) return value.map((entry) => replaceReference(entry, from, to))
  if (value !== null && typeof value === "object")
    return Object.fromEntries(
      Object.entries(value).map(([key, child]) => [key, replaceReference(child, from, to)]),
    )
  return value
}

const replaceEvidenceRefs = (value: unknown): unknown => {
  if (Array.isArray(value)) return value.map(replaceEvidenceRefs)
  if (value !== null && typeof value === "object")
    return Object.fromEntries(
      Object.entries(value).map(([key, child]) => [
        key,
        key === "sourceRef" ? "evidence[ev-001]" : replaceEvidenceRefs(child),
      ]),
    )
  return value
}

const requirementPayload = (subject: string) => ({
  assessment: replaceEvidenceRefs(
    replaceReference(
      structuredClone(evaluationExamples.RequirementAssessmentResult),
      "requirement:root::has-tests",
      subject,
    ),
  ),
  rating: replaceReference(
    structuredClone(evaluationExamples.RequirementRatingResult),
    "requirement:root::has-tests",
    subject,
  ),
  evidence: {
    observations: [{ id: "ev-001", kind: "file", role: "evaluated", path: "evidence.txt" }],
    limits: [],
  },
})

describe("CLI compatibility", () => {
  it("reports version and root help on stdout", () => {
    const version = run("--version")
    expect(version.status).toBe(0)
    expect(version.stdout).toMatch(/^qualitymd version /)
    expect(version.stderr).toBe("")

    const help = run("--help")
    expect(help.status).toBe(0)
    expect(help.stdout).toContain("USAGE\n\n  qualitymd [command] [--flags]")
    expect(help.stdout).toContain("evaluation")
  })

  it("uses the stable usage exit category for unknown commands", () => {
    const result = run("definitely-unknown")
    expect(result.status).toBe(2)
    expect(result.stdout).toBe("")
    expect(result.stderr).toContain('Unknown command "definitely-unknown"')

    const nested = run("lint", "--definitely-unknown")
    expect(nested.status).toBe(2)
    expect(nested.stderr).toContain("Unrecognized flag")
  })

  it("lints the project model as machine-readable output", () => {
    const result = run("lint", "QUALITY.md", "--json")
    expect(result.status).toBe(0)
    expect(JSON.parse(result.stdout)).toMatchObject({ valid: true })
    expect(result.stderr).toBe("")
  })

  it("reports a valid workspace status snapshot without treating absent runs as failure", () => {
    const result = run("status", "QUALITY.md", "--json")
    expect(result.status).toBe(0)
    expect(JSON.parse(result.stdout)).toMatchObject({
      schemaVersion: 2,
      model: { valid: true },
    })
    expect(result.stderr).toBe("")
  })

  it("emits bundled spec and schema artifacts", () => {
    const specification = run("spec")
    expect(specification.status).toBe(0)
    expect(specification.stdout).toContain("# QUALITY.md specification")

    const schema = run("schema")
    expect(schema.status).toBe(0)
    expect(JSON.parse(schema.stdout)).toHaveProperty("$schema")

    const schemaArgument = run("schema", "unexpected")
    expect(schemaArgument.status).toBe(2)
    expect(schemaArgument.stderr).toContain("accepts no arguments")
  })

  it("projects model tree, list, and detail JSON with canonical IDs", () => {
    const tree = run("model", "tree", "QUALITY.md", "--depth", "0", "--json")
    expect(tree.status).toBe(0)
    expect(JSON.parse(tree.stdout)).toMatchObject({ id: "area:root", kind: "area" })
    expect(JSON.parse(tree.stdout)).not.toHaveProperty("children")

    const list = run("model", "list", "QUALITY.md", "--type", "factor", "--json")
    expect(list.status).toBe(0)
    expect(
      JSON.parse(list.stdout).every((entry: { kind: string }) => entry.kind === "factor"),
    ).toBe(true)

    const first = JSON.parse(list.stdout)[0] as { id: string }
    const detail = run("model", "get", first.id, "QUALITY.md", "--json")
    expect(detail.status).toBe(0)
    expect(JSON.parse(detail.stdout).id).toBe(first.id)
  })

  it("keeps removed command aliases rejected", () => {
    for (const arguments_ of [
      ["upgrade"],
      ["evaluation", "record"],
      ["evaluation", "report", "generate"],
    ]) {
      const result = run(...arguments_)
      expect(result.status).toBe(2)
    }
  })

  it("scaffolds default, minimal, JSON, force, and stdout modes", async () => {
    await mkdir(join(repositoryRoot, "tmp"), { recursive: true })
    const directory = await mkdtemp(join(repositoryRoot, "tmp", "qualitymd-init-test-"))
    try {
      const stdout = runFrom(repositoryRoot, ["init", "-", "--minimal"])
      expect(stdout.status).toBe(0)
      expect(stdout.stdout).toContain("ratingScale:")

      const executable = join(repositoryRoot, "src/main.ts")
      const runLocal = (...arguments_: ReadonlyArray<string>) =>
        spawnSync("bun", [executable, ...arguments_], {
          cwd: directory,
          encoding: "utf8",
          env: { ...process.env, QUALITYMD_NO_UPDATE_CHECK: "1", NO_COLOR: "1" },
        })
      const initialized = runLocal("init", "--json", "--minimal", "--no-agent-instructions")
      expect(initialized.status).toBe(0)
      expect(JSON.parse(initialized.stdout)).toMatchObject({
        path: "QUALITY.md",
        created: true,
        agentInstructionFiles: [],
      })
      expect(await readFile(join(directory, "QUALITY.md"), "utf8")).not.toContain(
        "## Purpose and scope",
      )

      const refused = runLocal("init")
      expect(refused.status).toBe(70)
      expect(refused.stderr).toContain("already exists")
      const overwritten = runLocal("init", "--force", "--no-agent-instructions")
      expect(overwritten.status).toBe(0)
      expect(await readFile(join(directory, "QUALITY.md"), "utf8")).toContain("## Overview")
    } finally {
      await rm(directory, { recursive: true })
    }
  })

  it("reports invalid models without corrupting stdout/stderr categories", async () => {
    await mkdir(join(repositoryRoot, "tmp"), { recursive: true })
    const directory = await mkdtemp(join(repositoryRoot, "tmp", "qualitymd-lint-test-"))
    const path = join(directory, "QUALITY.md")
    try {
      await writeFile(path, "not frontmatter\n")
      const json = run("lint", path, "--json")
      expect(json.status).toBe(1)
      expect(JSON.parse(json.stdout)).toMatchObject({ valid: false })
      expect(json.stderr).toBe("")

      const human = run("lint", path)
      expect(human.status).toBe(1)
      expect(human.stdout).toContain("invalid-frontmatter")
      expect(human.stderr).toContain("Next:")
    } finally {
      await rm(directory, { recursive: true })
    }
  })

  it("applies deterministic lint repairs atomically", async () => {
    await mkdir(join(repositoryRoot, "tmp"), { recursive: true })
    const directory = await mkdtemp(join(repositoryRoot, "tmp", "qualitymd-fix-test-"))
    const path = join(directory, "QUALITY.md")
    try {
      await writeFile(
        path,
        `---
title: Test
description: ""
ratingScale:
  - level: target
    title: Target
    criterion: Satisfies the requirement.
  - level: unacceptable
    title: Unacceptable
    criterion: Does not satisfy the requirement.
factors:
  evidence:
    title: Evidence
    requirements:
      exists:
        title: Evidence exists
        assessment: Inspect the source.
---
`,
      )
      const before = run("lint", path, "--json")
      expect(JSON.parse(before.stdout).summary.fixable).toBeGreaterThan(0)
      const fixed = run("lint", path, "--fix", "--json")
      expect(fixed.status).toBe(0)
      expect(JSON.parse(fixed.stdout).summary.fixed).toBeGreaterThan(0)
      expect(await readFile(path, "utf8")).not.toContain('description: ""')
      const repeated = run("lint", path, "--fix", "--json")
      expect(JSON.parse(repeated.stdout).summary.fixed).toBe(0)
    } finally {
      await rm(directory, { recursive: true })
    }
  })

  it("passes unmatched source selectors to requirement inspection without prepackaging", async () => {
    await mkdir(join(repositoryRoot, "tmp"), { recursive: true })
    const directory = await mkdtemp(join(repositoryRoot, "tmp", "qualitymd-cli-test-"))
    const model = join(directory, "QUALITY.md")
    await writeFile(
      model,
      `---
title: Test
source: "*.missing"
ratingScale:
  - level: target
    title: Target
    criterion: Satisfies the requirement.
  - level: unacceptable
    title: Unacceptable
    criterion: Does not satisfy the requirement.
factors:
  evidence:
    title: Evidence
    requirements:
      exists:
        title: Evidence exists
        assessment: Inspect the selected source.
---
`,
    )
    try {
      const preview = run(
        "evaluation",
        "run",
        "--model",
        model,
        "--evaluator",
        "harness",
        "--dry-run",
        "--json",
      )
      expect(JSON.parse(preview.stdout)).toMatchObject({
        inspectionPolicy: {
          workspace: "read-only",
          network: "disabled",
          approvals: "never",
          verification: "unavailable",
          repositoryInstructions: "untrusted-data",
        },
        sources: [{ selector: "*.missing", kind: "glob" }],
      })
      const result = run("evaluation", "run", "--model", model, "--evaluator", "harness", "--json")
      expect(result.status).toBe(0)
      const receipt = JSON.parse(result.stdout) as {
        path: string
        status: string
        sources: Array<{ selector: string; kind: string }>
        evaluatorRequests: Array<Record<string, unknown>>
      }
      expect(receipt).toMatchObject({
        status: "awaiting_evaluator",
        sources: [{ selector: "*.missing", kind: "glob" }],
      })
      expect(receipt.evaluatorRequests[0]).toHaveProperty("inspection")
      expect(receipt.evaluatorRequests[0]).not.toHaveProperty("source")
      const artifact = JSON.parse(
        await readFile(join(directory, receipt.path, "evaluation.json"), "utf8"),
      )
      expect(artifact).toMatchObject({ schemaVersion: 8, evidence: {} })
    } finally {
      await rm(directory, { recursive: true })
    }
  })

  it("records lifecycle events, retries, and evaluator capabilities", async () => {
    await mkdir(join(repositoryRoot, "tmp"), { recursive: true })
    const directory = await mkdtemp(join(repositoryRoot, "tmp", "qualitymd-events-test-"))
    const model = join(directory, "QUALITY.md")
    await writeFile(join(directory, "evidence.txt"), "bounded evidence")
    await writeFile(
      model,
      `---
title: Test
source: evidence.txt
ratingScale:
  - level: target
    title: Target
    criterion: Satisfies the requirement.
  - level: unacceptable
    title: Unacceptable
    criterion: Does not satisfy the requirement.
factors:
  evidence:
    title: Evidence
    requirements:
      exists:
        title: Evidence exists
        assessment: Inspect the selected source.
---
`,
    )
    try {
      const created = run("evaluation", "run", "--model", model, "--evaluator", "harness", "--json")
      expect(created.status).toBe(0)
      const receipt = JSON.parse(created.stdout) as {
        path: string
        evaluatorRequests: Array<{
          requestId: string
          inputHash: string
          inspection: { source: { selector: string; kind: string } }
        }>
      }
      expect(receipt.evaluatorRequests[0]?.inspection.source).toEqual({
        selector: "evidence.txt",
        kind: "path",
      })
      const runPath = join(directory, receipt.path)
      const repeated = run("evaluation", "run", "--resume", runPath, "--json")
      expect(repeated.status).toBe(0)
      expect(
        (JSON.parse(repeated.stdout).evaluatorRequests as Array<{ requestId: string }>).map(
          (request) => request.requestId,
        ),
      ).toEqual(receipt.evaluatorRequests.map((request) => request.requestId))
      const initialEvents = await readFile(join(runPath, "logs", "events.jsonl"), "utf8")
      expect(initialEvents).toContain('"event":"run_created"')
      expect(initialEvents).toContain('"structuredOutput":true')

      const resultPath = join(directory, "retry.json")
      await writeFile(
        resultPath,
        JSON.stringify(
          receipt.evaluatorRequests.map((request) => ({
            requestId: request.requestId,
            inputHash: request.inputHash,
            evaluator: { runtime: "test-harness" },
            failure: "rate_limited",
            detail: "retry later",
          })),
        ),
      )
      const resumed = run(
        "evaluation",
        "run",
        "--resume",
        runPath,
        "--evaluator-result",
        resultPath,
        "--json",
      )
      expect(resumed.status).toBe(0)
      const retried = JSON.parse(resumed.stdout) as {
        status: string
        evaluatorRequests: Array<{
          requestId: string
          inputHash: string
          subject: string
        }>
      }
      expect(retried).toMatchObject({ status: "awaiting_evaluator" })
      const evidenceRetry = retried.evaluatorRequests[0]!
      const invalidEvidence = requirementPayload(evidenceRetry.subject)
      invalidEvidence.evidence.observations[0]!.path = "../outside.txt"
      await writeFile(
        resultPath,
        JSON.stringify({
          requestId: evidenceRetry.requestId,
          inputHash: evidenceRetry.inputHash,
          evaluator: { runtime: "test-harness" },
          payload: invalidEvidence,
        }),
      )
      const evidenceRejected = run(
        "evaluation",
        "run",
        "--resume",
        runPath,
        "--evaluator-result",
        resultPath,
        "--json",
      )
      expect(evidenceRejected.status, evidenceRejected.stderr + evidenceRejected.stdout).toBe(0)
      expect(JSON.parse(evidenceRejected.stdout)).toMatchObject({
        status: "awaiting_evaluator",
        evaluatorRequests: [{ lastFailure: { category: "evidence_invalid" } }],
      })
      const events = await readFile(join(runPath, "logs", "events.jsonl"), "utf8")
      expect(events).toContain('"event":"work_unit_retry"')
      const calls = await readFile(join(runPath, "logs", "evaluator-calls.jsonl"), "utf8")
      expect(calls).toContain('"failure":{"category":"rate_limited"')
      expect(calls).toContain('"failure":{"category":"evidence_invalid"')
      expect(calls).toContain('"capabilities":{"structuredOutput":true')
    } finally {
      await rm(directory, { recursive: true })
    }
  })

  it("persists accepted harness results across a partial submission and resume", async () => {
    await mkdir(join(repositoryRoot, "tmp"), { recursive: true })
    const directory = await mkdtemp(join(repositoryRoot, "tmp", "qualitymd-partial-test-"))
    const model = join(directory, "QUALITY.md")
    await writeFile(join(directory, "evidence.txt"), "bounded evidence")
    await writeFile(
      model,
      `---
title: Test
source: evidence.txt
ratingScale:
  - level: target
    title: Target
    criterion: Satisfies the requirement.
  - level: unacceptable
    title: Unacceptable
    criterion: Does not satisfy the requirement.
factors:
  evidence:
    title: Evidence
    requirements:
      first:
        title: First evidence exists
        assessment: Inspect the selected source.
      second:
        title: Second evidence exists
        assessment: Inspect the selected source.
---
`,
    )
    try {
      const created = run("evaluation", "run", "--model", model, "--evaluator", "harness", "--json")
      expect(created.status).toBe(0)
      const receipt = JSON.parse(created.stdout) as {
        path: string
        evaluatorRequests: Array<{
          requestId: string
          inputHash: string
          workUnitId: string
          kind: string
          subject: string
        }>
      }
      expect(receipt.evaluatorRequests).toHaveLength(2)
      const accepted = receipt.evaluatorRequests[0]!
      expect(accepted.kind).toBe("assessRateRequirement")
      const resultPath = join(directory, "partial.json")
      await writeFile(
        resultPath,
        JSON.stringify([
          {
            requestId: accepted.requestId,
            inputHash: accepted.inputHash,
            evaluator: { runtime: "test-harness" },
            payload: requirementPayload(accepted.subject),
          },
        ]),
      )
      const runPath = join(directory, receipt.path)
      const partial = run(
        "evaluation",
        "run",
        "--resume",
        runPath,
        "--evaluator-result",
        resultPath,
        "--json",
      )
      expect(partial.status).toBe(0)
      expect(JSON.parse(partial.stdout)).toMatchObject({ status: "awaiting_evaluator" })

      const artifactPath = join(runPath, "evaluation.json")
      const afterPartial = JSON.parse(await readFile(artifactPath, "utf8")) as {
        state: { workUnits: Record<string, { status: string }> }
        results: { payloads: Array<{ workUnit: string }> }
        evidence: Record<string, { observations: ReadonlyArray<{ path: string }> }>
      }
      expect(afterPartial.state.workUnits[accepted.workUnitId]?.status).toBe("completed")
      expect(
        afterPartial.results.payloads.filter((entry) => entry.workUnit === accepted.workUnitId),
      ).toHaveLength(2)
      expect(afterPartial.evidence[accepted.workUnitId]?.observations[0]?.path).toBe("evidence.txt")

      const resumed = run("evaluation", "run", "--resume", runPath, "--json")
      expect(resumed.status).toBe(0)
      const afterResume = JSON.parse(await readFile(artifactPath, "utf8")) as typeof afterPartial
      expect(
        afterResume.results.payloads.filter((entry) => entry.workUnit === accepted.workUnitId),
      ).toHaveLength(2)
    } finally {
      await rm(directory, { recursive: true })
    }
  })
})
