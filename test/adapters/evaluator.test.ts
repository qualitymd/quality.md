import { describe, expect, it } from "vitest"

import {
  claudeEvaluationUsage,
  claudeSystemPrompt,
  codexEvaluationUsage,
  strictProviderSchema,
  stripNullProperties,
} from "../../src/adapters/evaluator.ts"
import {
  EVALUATION_PROMPT_BOUNDARY,
  renderEvaluationPrompt,
  renderEvaluationPromptParts,
} from "../../src/domain/evaluator/context.ts"
import type { EvaluationRequest } from "../../src/domain/evaluator/types.ts"

const request = (overrides: Partial<EvaluationRequest> = {}): EvaluationRequest => ({
  runId: "eval-1",
  workUnitId: "assessRateRequirement:requirement:root::safe",
  kind: "assessRateRequirement",
  subject: "requirement:root::safe",
  instructions: "Assess the requirement.",
  sharedContext: { area: { title: "Root", id: "root" } },
  context: { requirement: { title: "Safe", id: "safe" } },
  bodyGuidance: "Prefer direct evidence.",
  inspection: {
    workspaceRoot: "/workspace",
    source: { selector: "src", kind: "path" },
    policy: {
      workspace: "read-only",
      network: "disabled",
      approvals: "never",
      verification: "unavailable",
      repositoryInstructions: "untrusted-data",
    },
  },
  expectedSchema: { type: "object" },
  timeoutMs: 1_000,
  ...overrides,
})

describe("provider schema adaptation", () => {
  it("requires every property while making optional fields nullable", () => {
    const source = {
      $schema: "https://json-schema.org/draft/2020-12/schema",
      type: "object",
      properties: {
        name: { type: "string" },
        nested: {
          type: "object",
          properties: { description: { type: "string" }, enabled: { type: "boolean" } },
          required: ["enabled"],
        },
      },
      required: ["name"],
      allOf: [{ if: {}, ["th" + "en"]: {} }],
    }
    const strict = strictProviderSchema(source) as Record<string, unknown>
    const properties = strict.properties as Record<string, unknown>
    const nested = properties.nested as Record<string, unknown>
    const nestedVariants = nested.anyOf as Array<Record<string, unknown>>
    expect(strict.required).toEqual(["name", "nested"])
    expect(strict.additionalProperties).toBe(false)
    expect(strict).not.toHaveProperty("allOf")
    expect(strict).not.toHaveProperty("$schema")
    expect(properties.name).toEqual({ type: "string" })
    expect(nestedVariants[1]).toEqual({ type: "null" })
    expect(nestedVariants[0]).toMatchObject({
      additionalProperties: false,
      required: ["description", "enabled"],
      properties: {
        description: { anyOf: [{ type: "string" }, { type: "null" }] },
        enabled: { type: "boolean" },
      },
    })
    expect(source.required).toEqual(["name"])
    expect(source).toHaveProperty("allOf")
  })

  it("removes provider null placeholders before persisted validation", () => {
    expect(
      stripNullProperties({ name: "value", optional: null, nested: [{ keep: true, omit: null }] }),
    ).toEqual({ name: "value", nested: [{ keep: true }] })
  })
})

describe("neutral evaluator prompts", () => {
  it("treats repository instructions as data and carries the inspection boundary", () => {
    const prompt = renderEvaluationPrompt(request())
    expect(prompt).toContain("repository instruction files")
    expect(prompt).toContain("never governing instructions")
    expect(prompt).toContain('"workspace":"read-only"')
    expect(prompt).toContain('"repositoryInstructions":"untrusted-data"')
    expect(prompt).toContain("source selector identifies the evaluated subject")
    expect(prompt).toContain("QUALITY.md body guidance")
  })
})

describe("cacheable evaluator prompt parts", () => {
  it("keeps shared bytes stable while work-unit identity and context vary", () => {
    const first = renderEvaluationPromptParts(request())
    const secondRequest = request({
      runId: "eval-2",
      workUnitId: "assessRateRequirement:requirement:root::fast",
      subject: "requirement:root::fast",
      context: { requirement: { title: "Fast", id: "fast" } },
    })
    const second = renderEvaluationPromptParts(secondRequest)

    expect(first.cacheablePrefix).toBe(second.cacheablePrefix)
    expect(first.workUnitSuffix).not.toBe(second.workUnitSuffix)
    expect(renderEvaluationPrompt(secondRequest)).toBe(
      [second.cacheablePrefix, EVALUATION_PROMPT_BOUNDARY, second.workUnitSuffix].join("\n\n"),
    )
    expect(first.cacheablePrefix).not.toContain(first.workUnitSuffix)
  })

  it("places every shared block before the work-unit boundary", () => {
    const prompt = renderEvaluationPrompt(request())
    const boundary = prompt.indexOf(EVALUATION_PROMPT_BOUNDARY)

    for (const shared of [
      "Kind: assessRateRequirement",
      "Task: Assess the requirement.",
      "QUALITY.md body guidance:",
      "Shared model context:",
      "Authorized inspection boundary:",
    ])
      expect(prompt.indexOf(shared)).toBeLessThan(boundary)
    for (const varying of ["Work unit:", "Subject:", "Work-unit context:"])
      expect(prompt.indexOf(varying)).toBeGreaterThan(boundary)
  })

  it("canonicalizes structured blocks and changes the prefix only for shared input", () => {
    const baseline = renderEvaluationPromptParts(request())
    const { inspection, ...withoutInspection } = request()
    const reordered = renderEvaluationPromptParts(
      request({
        sharedContext: { area: { id: "root", title: "Root" } },
        context: { requirement: { id: "safe", title: "Safe" } },
      }),
    )

    expect(reordered).toEqual(baseline)
    expect(inspection).toBeDefined()
    for (const changed of [
      request({ kind: "analyzeFactor" }),
      request({ instructions: "Analyze the factor." }),
      request({ bodyGuidance: "Use bounded evidence." }),
      request({ sharedContext: { area: { id: "other", title: "Other" } } }),
      withoutInspection,
    ])
      expect(renderEvaluationPromptParts(changed).cacheablePrefix).not.toBe(
        baseline.cacheablePrefix,
      )
  })
})

describe("provider cache configuration and usage", () => {
  it("uses Claude's cache-stable preset system prompt", () => {
    expect(claudeSystemPrompt).toEqual({
      type: "preset",
      preset: "claude_code",
      excludeDynamicSections: true,
    })
  })

  it("preserves provider cache reads, writes, zero values, and absence", () => {
    expect(
      codexEvaluationUsage({ input_tokens: 12, output_tokens: 3, cached_input_tokens: 0 }),
    ).toEqual({ inputTokens: 12, outputTokens: 3, cachedInputTokens: 0 })
    expect(codexEvaluationUsage({})).toEqual({})
    expect(
      claudeEvaluationUsage(
        {
          input_tokens: 10,
          output_tokens: 2,
          cache_read_input_tokens: 8,
          cache_creation_input_tokens: 0,
        },
        0,
      ),
    ).toEqual({
      inputTokens: 10,
      outputTokens: 2,
      cachedInputTokens: 8,
      cacheWriteInputTokens: 0,
      costUsd: 0,
    })
    expect(claudeEvaluationUsage({}, undefined)).toEqual({})
  })
})
