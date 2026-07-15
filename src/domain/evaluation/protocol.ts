import evaluationSchema from "../../assets/evaluation-data.schema.json"
import { hashJson, sha256 } from "../json.ts"
import type { SourceBundle } from "../evaluator/types.ts"
import type { QualityModel } from "../model/model.ts"
import type { WorkUnit } from "./graph.ts"
import type { EvaluationPlan, PlannedFactor } from "./plan.ts"

type JsonObject = Record<string, unknown>

export interface StoredPayload {
  readonly workUnit: string
  readonly payload: JsonObject
}

export interface ProtocolRequest {
  readonly workUnitId: string
  readonly kind: string
  readonly subject: string
  readonly instructions: string
  readonly sharedContext: JsonObject | null
  readonly context: JsonObject | null
  readonly source: SourceBundle["files"]
  readonly expectedSchema: JsonObject
  readonly expectedSchemaText: string
  readonly inputHash: string
  readonly correlationId: string
}

const schema = evaluationSchema as {
  readonly $schema: string
  readonly $defs: Readonly<Record<string, Readonly<JsonObject>>>
}

const kindSchema = (kind: string): JsonObject => ({
  $id: `https://getquality.md/evaluation-data.schema.json/${kind}`,
  $schema: schema.$schema,
  ...schema.$defs[kind],
})

const sourceSchema: JsonObject = {
  additionalProperties: false,
  properties: {
    files: {
      items: {
        additionalProperties: false,
        properties: {
          path: { minLength: 1, type: "string" },
        },
        required: ["path"],
        type: "object",
      },
      minItems: 1,
      type: "array",
    },
  },
  required: ["files"],
  type: "object",
}

const assessmentSchema: JsonObject = {
  additionalProperties: false,
  properties: {
    assessment: kindSchema("RequirementAssessmentResult"),
    rating: kindSchema("RequirementRatingResult"),
  },
  required: ["assessment", "rating"],
  type: "object",
}

const recommendationSchema: JsonObject = {
  additionalProperties: false,
  properties: {
    recommendations: {
      items: kindSchema("RecommendationResult"),
      minItems: 1,
      type: "array",
    },
  },
  required: ["recommendations"],
  type: "object",
}

export const expectedSchemaFor = (unit: WorkUnit): JsonObject => {
  if (unit.kind === "resolveSource") return sourceSchema
  if (unit.kind === "assessRateRequirement") return assessmentSchema
  if (unit.kind === "recommend") return recommendationSchema
  return kindSchema(unit.dataKind!)
}

const payloadFor = (payloads: ReadonlyArray<StoredPayload>, workUnit: string) =>
  payloads.find((entry) => entry.workUnit === workUnit)?.payload

const payloadsFor = (payloads: ReadonlyArray<StoredPayload>, workUnit: string) =>
  payloads.filter((entry) => entry.workUnit === workUnit).map((entry) => entry.payload)

const requirementPayload = (
  payloads: ReadonlyArray<StoredPayload>,
  requirement: string,
  kind: string,
) =>
  payloadsFor(payloads, `assessRateRequirement:${requirement}`).find(
    (payload) => payload.kind === kind,
  )

const routineRef = (kind: string, subject: JsonObject, selector = "") => ({
  kind,
  ...(selector === "" ? {} : { selector }),
  subject,
})

const factorChildren = (factor: PlannedFactor, plan: EvaluationPlan) =>
  plan.factors
    .filter(
      (candidate) =>
        candidate.areaId === factor.areaId &&
        candidate.path.length === factor.path.length + 1 &&
        candidate.path.slice(0, -1).join("/") === factor.path.join("/"),
    )
    .map((candidate) => candidate.ref)

const factorRequirements = (factor: PlannedFactor, plan: EvaluationPlan) =>
  plan.requirements
    .filter((requirement) => requirement.factorIds.includes(factor.ref))
    .map((requirement) => requirement.ref)

const findingIndex = (plan: EvaluationPlan, payloads: ReadonlyArray<StoredPayload>) => {
  const findings = plan.requirements.flatMap((requirement) => {
    const assessment = requirementPayload(payloads, requirement.ref, "RequirementAssessmentResult")
    if (assessment === undefined || !Array.isArray(assessment.findings)) return []
    return assessment.findings.flatMap((value) => {
      if (value === null || typeof value !== "object") return []
      const finding = value as JsonObject
      if (typeof finding.id !== "string" || finding.id === "") return []
      return [
        {
          requirementId: requirement.ref,
          findingId: finding.id,
          findingRef: routineRef(
            "RequirementAssessmentResult",
            { requirementId: requirement.ref },
            `findings[${finding.id}]`,
          ),
          ...(typeof finding.type === "string" ? { type: finding.type } : {}),
          ...(typeof finding.severity === "string" ? { severity: finding.severity } : {}),
          ...(typeof finding.confidence === "string" ? { confidence: finding.confidence } : {}),
          ...(typeof finding.statement === "string" ? { statement: finding.statement } : {}),
        },
      ]
    })
  })
  return findings.length === 0 ? null : findings
}

const instructions = {
  resolveSource:
    'Resolve this area\'s source selector and return one JSON object of the form {"files": [{"path": string}, ...]}.\n' +
    "- The selector describes a body of evidence; use read-only tools to locate exactly the workspace text files it names.\n" +
    "- path is the unique, non-empty workspace-relative path of each gathered file; do not return URLs, external identifiers, absolute paths, or paths that escape the workspace.\n" +
    "- Return paths only. The runner rereads, bounds, hashes, and persists the files; do not return file content, exploration transcripts, summaries, assessments, or ratings.\n" +
    "- Gather only what the selector describes; do not widen to adjacent material.\n" +
    "- If the material the selector describes does not exist — including when the selector reads like a filesystem path that names nothing — return the classified failure source_unavailable naming the selector instead of improvising or substituting evidence.",
  assessRateRequirement:
    'Assess this requirement against the packaged source evidence, then rate it from that assessment, and return one JSON object of the form {"assessment": RequirementAssessmentResult, "rating": RequirementRatingResult}.\n' +
    "- Set requirementId in both objects to the subject reference exactly.\n" +
    "Assessment:\n" +
    "- status is one of: assessed, partially_assessed, blocked, not_applicable.\n" +
    "- Record every finding with the full core shape (id, type, confidence, statement, condition, criteria, basis, effect, evidence). Gap and risk findings carry severity; strength and note findings must not.\n" +
    "- criteria entries reference this requirement and a rating level from the frame's appliedRatingCriteria.\n" +
    "- Cite evidence sourceRef values from the packaged source paths.\n" +
    "- If required evidence is unavailable, say so via status, unknowns, and evaluationLimits instead of guessing.\n" +
    "- Use finding ids like gap-001, strength-001, risk-001, note-001, unique within this assessment.\n" +
    "Rating:\n" +
    "- Judge only from your assessment and the frame's appliedRatingCriteria; do not rate on evidence the assessment does not record.\n" +
    "- status is one of: rated, not_rated, blocked, not_applicable. When rated, set ratingLevelId to the highest rating level whose criterion the assessed evidence satisfies and explain the rationale.\n" +
    "- Record criteriaResults for each rating level considered and ratingDrivers referencing the assessment.",
  analyzeFactor:
    "Synthesize this factor's analysis from its direct requirement ratings and child factor analyses, and return one FactorAnalysisResult JSON object.\n" +
    "- Set factorId to the subject reference exactly.\n" +
    "- Fill localAnalysis (direct inputs only) and localAndDescendantAnalysis (including child factors), each with status, ratingLevelId when analyzed, rationale, inputRefs, and ratingDrivers.\n" +
    "- Follow the frame's synthesis guidance: the roll-up rating is bounded by the worst contributing input (worst_bound); ignore empty inputs.\n" +
    "- Do not inspect new evidence; synthesize only the given inputs.",
  analyzeArea:
    "Synthesize this area's analysis from its factor analyses, child area analyses, and local requirement ratings, and return one AreaAnalysisResult JSON object.\n" +
    "- Set areaId to the subject reference exactly.\n" +
    "- Fill localAnalysis (local inputs only) and localAndDescendantAnalysis (including child areas), each with status, ratingLevelId when analyzed, rationale, inputRefs, and ratingDrivers.\n" +
    "- Follow the frame's synthesis guidance: the roll-up rating is bounded by the worst contributing input (worst_bound); ignore empty inputs.\n" +
    "- Do not inspect new evidence; synthesize only the given inputs.",
  rankFindings:
    "Rank every persisted finding across the evaluation scope and return one FindingRankingResult JSON object.\n" +
    "- orderedFindings must contain exactly one entry per finding in the findings context, no more, no fewer.\n" +
    "- Copy each entry's findingRef object verbatim from the findings context.\n" +
    "- rank is 1-based and unique; tier is one of P1, P2, P3, P4 (P1 = act first).\n" +
    "- Rank by severity, confidence, and breadth of effect; give each entry a one-sentence rationale.",
  recommend:
    'Propose quality-management recommendations from the ranked findings and analyses, and return one JSON object of the form {"recommendations": [RecommendationResult, ...]}.\n' +
    "- Each recommendation needs title, description, background, expectedValue, doneCriterion, impact (high, medium, or low), confidence, and non-empty traceRefs pointing at the findings or analyses it follows from.\n" +
    "- Omit the id field; the runner assigns canonical recommendation IDs.\n" +
    "- Do not include planning fields (effort, roi, quickWin, priority, score).\n" +
    "- Cover the highest-ranked findings first; propose the smallest set of recommendations that addresses the advice-driving findings (typically 2-6).",
  rankRecommendations:
    "Rank the recommendations and account for finding coverage, and return one RecommendationRankingResult JSON object.\n" +
    "- orderedRecommendations must contain exactly one entry per recommendation in the recommendations context, using each recommendation's id as recommendationRef, with 1-based unique ranks, impact, confidence, and rationale.\n" +
    "- findingCoverage must contain exactly one entry per finding in the findings context: copy findingRef verbatim, set disposition to addressed_by_recommendation (with recommendationRefs listing covering recommendation ids) or not_advice_driving (with a short rationale).",
} as const

const emptyBundle = async (): Promise<SourceBundle> => ({
  files: [],
  hash: await sha256(""),
  truncated: false,
})

export const buildProtocolRequest = async (options: {
  readonly unit: WorkUnit
  readonly plan: EvaluationPlan
  readonly payloads: ReadonlyArray<StoredPayload>
  readonly bundles: ReadonlyMap<string, SourceBundle>
  readonly sourceRecords: Readonly<Record<string, JsonObject>>
  readonly evaluationId: string
}): Promise<ProtocolRequest> => {
  const { unit, plan, payloads, bundles, sourceRecords, evaluationId } = options
  const expectedSchema = expectedSchemaFor(unit)
  const expectedSchemaText =
    JSON.stringify(expectedSchema, null, 2) +
    (unit.kind === "resolveSource" ||
    unit.kind === "assessRateRequirement" ||
    unit.kind === "recommend"
      ? ""
      : "\n")
  let requestInstructions = ""
  let sharedContext: JsonObject | null = null
  let context: JsonObject | null = null
  let bundle = await emptyBundle()

  if (unit.kind === "resolveSource") {
    requestInstructions = instructions.resolveSource
    sharedContext = {
      areaEvaluationFrame: payloadFor(payloads, `frameAreaEvaluation:${unit.subject}`)!,
    }
    const record = sourceRecords[unit.subject]!
    context = { sourceSelector: { selector: record.selector, kind: record.kind } }
  } else if (unit.kind === "assessRateRequirement") {
    requestInstructions = instructions.assessRateRequirement
    const requirement = plan.requirements.find((entry) => entry.ref === unit.subject)!
    bundle = bundles.get(requirement.areaId)!
    sharedContext = {
      areaEvaluationFrame: payloadFor(payloads, `frameAreaEvaluation:${requirement.areaId}`)!,
    }
    context = {
      requirement: {
        assessment: requirement.value.assessment,
        description: requirement.value.description ?? "",
        requirementId: requirement.ref,
        title: requirement.value.title,
      },
      requirementEvaluationFrame: payloadFor(
        payloads,
        `frameRequirementEvaluation:${requirement.ref}`,
      )!,
    }
  } else if (unit.kind === "analyzeFactor") {
    requestInstructions = instructions.analyzeFactor
    const factor = plan.factors.find((entry) => entry.ref === unit.subject)!
    context = {
      factorAnalysisFrame: payloadFor(payloads, `frameFactorAnalysis:${factor.ref}`)!,
      directRequirementRatings: Object.fromEntries(
        factorRequirements(factor, plan).map((ref) => [
          ref,
          requirementPayload(payloads, ref, "RequirementRatingResult"),
        ]),
      ),
      childFactorAnalyses: Object.fromEntries(
        factorChildren(factor, plan).map((ref) => [
          ref,
          payloadFor(payloads, `analyzeFactor:${ref}`),
        ]),
      ),
    }
  } else if (unit.kind === "analyzeArea") {
    requestInstructions = instructions.analyzeArea
    const area = plan.areas.find((entry) => entry.ref === unit.subject)!
    context = {
      areaAnalysisFrame: payloadFor(payloads, `frameAreaAnalysis:${area.ref}`)!,
      factorAnalyses: Object.fromEntries(
        area.rootFactorIds.map((ref) => [ref, payloadFor(payloads, `analyzeFactor:${ref}`)]),
      ),
      childAreaAnalyses: Object.fromEntries(
        area.childAreaIds.map((ref) => [ref, payloadFor(payloads, `analyzeArea:${ref}`)]),
      ),
      localRequirementRatings: Object.fromEntries(
        area.localRequirementIds.map((ref) => [
          ref,
          requirementPayload(payloads, ref, "RequirementRatingResult"),
        ]),
      ),
    }
  } else if (unit.kind === "rankFindings") {
    requestInstructions = instructions.rankFindings
    context = { findings: findingIndex(plan, payloads) }
  } else if (unit.kind === "recommend") {
    requestInstructions = instructions.recommend
    context = {
      areaAnalyses: Object.fromEntries(
        plan.areas.map((area) => [area.ref, payloadFor(payloads, `analyzeArea:${area.ref}`)]),
      ),
      findingRanking: payloadFor(payloads, "rankFindings"),
      findings: findingIndex(plan, payloads),
    }
  } else if (unit.kind === "rankRecommendations") {
    requestInstructions = instructions.rankRecommendations
    context = {
      recommendations: payloadsFor(payloads, "recommend"),
      findings: findingIndex(plan, payloads),
      findingRanking: payloadFor(payloads, "rankFindings"),
    }
  }
  const inputHash = await hashJson({
    instructions: requestInstructions,
    sharedContext,
    context,
    schema: expectedSchemaText,
    source: bundle.hash,
  })
  return {
    workUnitId: unit.id,
    kind: unit.kind,
    subject: unit.subject,
    instructions: requestInstructions,
    sharedContext,
    context,
    source: bundle.files,
    expectedSchema,
    expectedSchemaText,
    inputHash,
    correlationId: `${evaluationId}#${unit.id}`,
  }
}

const ratingIds = (model: QualityModel) => model.ratingScale.map((level) => `rating:${level.level}`)

export const deterministicPayload = (
  unit: WorkUnit,
  model: QualityModel,
  plan: EvaluationPlan,
  modelPath: string,
): JsonObject => {
  if (unit.kind === "frameEvaluation") {
    return {
      derivedContext: {
        evaluationPolicies: ["source-as-data", "secret-redaction"],
        rigor: "standard",
      },
      inputs: { ratingLevelIds: ratingIds(model) },
      kind: "EvaluationFrame",
      schemaVersion: 3,
      subject: { modelLocator: modelPath },
    }
  }
  if (unit.kind === "frameAreaEvaluation") {
    const area = plan.areas.find((entry) => entry.ref === unit.subject)!
    return {
      inputs: {
        childAreaIds: area.childAreaIds,
        localRequirementIds: area.localRequirementIds,
        rootFactorIds: area.rootFactorIds,
        ...(area.source === "" ? {} : { sourceRefs: [area.source] }),
      },
      kind: "AreaEvaluationFrame",
      schemaVersion: 3,
      subject: { areaId: area.ref },
    }
  }
  if (unit.kind === "frameRequirementEvaluation") {
    const requirement = plan.requirements.find((entry) => entry.ref === unit.subject)!
    const criteria = model.ratingScale.flatMap((level) => {
      const criterion = requirement.value.ratings?.[level.level] ?? level.criterion
      if (criterion === "") return []
      return [
        {
          criterion,
          ratingLevelId: `rating:${level.level}`,
          source:
            requirement.value.ratings?.[level.level] === undefined
              ? "model_default"
              : "requirement_override",
        },
      ]
    })
    return {
      ...(criteria.length === 0 ? {} : { derivedContext: { appliedRatingCriteria: criteria } }),
      inputs: {
        ratingLevelIds: ratingIds(model),
        ...(requirement.value.assessment === ""
          ? {}
          : { requirementAssessmentBasis: requirement.value.assessment }),
      },
      kind: "RequirementEvaluationFrame",
      schemaVersion: 3,
      subject: {
        ...(requirement.factorIds.length === 0 ? {} : { factorIds: requirement.factorIds }),
        requirementId: requirement.ref,
      },
    }
  }
  if (unit.kind === "frameFactorAnalysis") {
    const factor = plan.factors.find((entry) => entry.ref === unit.subject)!
    return {
      derivedContext: {
        emptySignalPolicy: "ignore_empty",
        synthesisGuidanceRef: "protocol:factor-synthesis-default-v0",
      },
      inputs: {
        childFactorAnalysisRefs: factorChildren(factor, plan).map((ref) =>
          routineRef("FactorAnalysisResult", { factorId: ref }, "localAndDescendantAnalysis"),
        ),
        directRequirementRatingRefs: factorRequirements(factor, plan).map((ref) =>
          routineRef("RequirementRatingResult", { requirementId: ref }),
        ),
      },
      kind: "FactorAnalysisFrame",
      schemaVersion: 3,
      subject: { areaId: factor.areaId, factorId: factor.ref },
    }
  }
  if (unit.kind === "frameAreaAnalysis") {
    const area = plan.areas.find((entry) => entry.ref === unit.subject)!
    return {
      derivedContext: {
        emptySignalPolicy: "ignore_empty",
        synthesisGuidanceRef: "protocol:area-synthesis-default-v0",
      },
      inputs: {
        childAreaAnalysisRefs: area.childAreaIds.map((ref) =>
          routineRef("AreaAnalysisResult", { areaId: ref }, "localAndDescendantAnalysis"),
        ),
        factorAnalysisRefs: area.rootFactorIds.map((ref) =>
          routineRef("FactorAnalysisResult", { factorId: ref }, "localAndDescendantAnalysis"),
        ),
      },
      kind: "AreaAnalysisFrame",
      schemaVersion: 3,
      subject: { areaId: area.ref },
    }
  }
  throw new Error(`no deterministic payload for ${unit.id}`)
}
