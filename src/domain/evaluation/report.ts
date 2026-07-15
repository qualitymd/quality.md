import { dataPathForPayload, resolveDataKind } from "./data.ts"
import type { EvaluationPlan, PlannedArea, PlannedFactor, PlannedRequirement } from "./plan.ts"
import type { QualityModel } from "../model/model.ts"

type Json = Record<string, unknown>

export interface ReportManifest {
  readonly evaluationId: string
  readonly createdAt: string
  readonly model: string
  readonly requestedScope: Json
  readonly plannedScope: { readonly areaId: string; readonly factorFilter: ReadonlyArray<string> }
  readonly run: { readonly number: number; readonly label: string }
}

export interface RenderedReport {
  readonly kind: string
  readonly path: string
  readonly areaId?: string
  readonly factorId?: string
  readonly requirementId?: string
  readonly recommendationId?: string
  readonly content: string
}

export interface BuiltReportTree {
  readonly reports: ReadonlyArray<RenderedReport>
  readonly output: Json
  readonly rating: Json
  readonly payloadFiles: ReadonlyArray<{ readonly path: string; readonly payload: Json }>
}

const object = (value: unknown): Json =>
  value !== null && !Array.isArray(value) && typeof value === "object" ? (value as Json) : {}
const objects = (value: unknown): ReadonlyArray<Json> =>
  Array.isArray(value) ? value.map(object).filter((entry) => Object.keys(entry).length > 0) : []
const string = (value: unknown) => (typeof value === "string" ? value : "")
const value = (source: unknown, key: string) => object(source)[key]
const field = (source: unknown, key: string) => string(value(source, key))
const scoped = (source: unknown, key: string) => object(value(source, key))

const payloadIndex = (payloads: ReadonlyArray<Json>) => {
  const byKind = new Map<string, Array<Json>>()
  for (const payload of payloads) {
    const kind = field(payload, "kind")
    const items = byKind.get(kind) ?? []
    items.push(payload)
    byKind.set(kind, items)
  }
  return {
    all: payloads,
    kind: (kind: string) => byKind.get(kind) ?? [],
    one: (kind: string, key?: string, expected?: string) =>
      (byKind.get(kind) ?? []).find((payload) =>
        key === undefined ? true : field(payload, key) === expected,
      ),
  }
}

const splitArea = (reference: string) =>
  reference === "area:root" ? [] : reference.slice("area:".length).split("/")
const splitFactor = (reference: string) => {
  const [area = "root", path = ""] = reference.slice("factor:".length).split("::")
  return { areaId: `area:${area}`, path: path.split("/") }
}
const splitRequirement = (reference: string) => {
  const [area = "root", name = ""] = reference.slice("requirement:".length).split("::")
  return { areaId: `area:${area}`, name }
}
const areaDir = (areaId: string) => {
  const parts = splitArea(areaId)
  return parts.length === 0 ? "" : `areas/${parts.join("/")}/`
}
export const areaReportPath = (areaId: string) => {
  const parts = splitArea(areaId)
  return parts.length === 0 ? "root-area.md" : `areas/${parts.join("/")}/${parts.at(-1)}-area.md`
}
export const factorReportPath = (factorId: string) => {
  const factor = splitFactor(factorId)
  return `${areaDir(factor.areaId)}${factor.path.flatMap((part) => ["factors", part]).join("/")}/${factor.path.at(-1)}-factor.md`
}
export const requirementReportPath = (requirementId: string) => {
  const requirement = splitRequirement(requirementId)
  return `${areaDir(requirement.areaId)}requirements/${requirement.name}/${requirement.name}-requirement.md`
}

const dirname = (path: string) => (path.includes("/") ? path.slice(0, path.lastIndexOf("/")) : ".")
const normalizeParts = (parts: ReadonlyArray<string>) => {
  const result: Array<string> = []
  for (const part of parts) {
    if (part === "" || part === ".") continue
    if (part === "..") result.pop()
    else result.push(part)
  }
  return result
}
const relative = (fromFile: string, to: string) => {
  const from = normalizeParts(dirname(fromFile).split("/"))
  const target = normalizeParts(to.split("/"))
  let common = 0
  while (common < from.length && common < target.length && from[common] === target[common]) common++
  return [...from.slice(common).map(() => ".."), ...target.slice(common)].join("/") || "."
}
const link = (from: string, to: string, label: string) =>
  `[${escapeCell(label)}](${relative(from, to)})`
const escapeCell = (text: string) => text.replaceAll("|", "\\|").replaceAll("\n", " ")
const code = (text: string) => `\`${text.replaceAll("`", "\\`")}\``
const row = (...cells: ReadonlyArray<string>) => `| ${cells.map(escapeCell).join(" | ")} |\n`

const frontmatter = (type: string, title: string, rest: ReadonlyArray<[string, string]> = []) =>
  `---\ntype: ${type}\ntitle: ${/^[A-Za-z0-9 ._-]+$/.test(title) ? title : JSON.stringify(title)}\n${rest
    .map(([key, item]) => `${key}: ${item.includes(":") ? JSON.stringify(item) : item}\n`)
    .join("")}---\n\n`

const titleForArea = (area: PlannedArea) => area.value.title || area.path.at(-1) || "root"
const titleForFactor = (factor: PlannedFactor) =>
  factor.value.title || factor.path.at(-1) || "factor"
const titleForRequirement = (requirement: PlannedRequirement) =>
  requirement.value.title || splitRequirement(requirement.ref).name

const display = {
  analysis: (status: string) =>
    ({
      analyzed: "✅ Analyzed",
      empty: "⬜ Empty",
      not_analyzed: "⚪ Not Analyzed",
      blocked: "⛔ Blocked",
    })[status] || humanize(status),
  assessment: (status: string) =>
    ({
      assessed: "✅ Assessed",
      partially_assessed: "🟡 Partially Assessed",
      not_assessed: "⚪ Not Assessed",
      blocked: "⛔ Blocked",
    })[status] || humanize(status),
  confidence: (status: string) =>
    ({ high: "🟢 High", medium: "🔵 Medium", low: "🟡 Low", none: "⚪ None" })[status] ||
    humanize(status),
  impact: (status: string) =>
    ({ very_high: "⬥⬥ Very high", high: "⬥ High", medium: "● Medium", low: "○ Low" })[status] ||
    humanize(status),
  severity: (status: string) =>
    ({ critical: "🔴 Critical", high: "🔴 High", medium: "🟡 Medium", low: "🔵 Low" })[status] ||
    humanize(status),
  findingType: (status: string) =>
    ({ gap: "🚩 Gap", risk: "⚠️ Risk", strength: "💪 Strength", note: "ℹ️ Note" })[status] ||
    humanize(status),
  basis: (status: string) =>
    ({
      verified: "✅ Verified",
      plausible: "🟡 Plausible",
      not_assessed: "⚪ Not Assessed",
      not_applicable: "⬜ Not Applicable",
    })[status] || humanize(status),
  tier: (status: string) =>
    ({ P1: "🔴 P1 Highest", P2: "🟠 P2 High", P3: "🟡 P3 Medium", P4: "⚪ P4 Low" })[status] ||
    humanize(status),
}
const humanize = (text: string) =>
  text === ""
    ? ""
    : text.replaceAll(/[_-]+/g, " ").replaceAll(/\b\w/g, (letter) => letter.toUpperCase())

const rating = (model: QualityModel, analysis: Json) => {
  const status = field(analysis, "status")
  if (status !== "analyzed" && status !== "rated")
    return status === "blocked" ? "⛔ Blocked" : status === "empty" ? "⬜ Empty" : "⚪ Not Rated"
  const id = field(analysis, "ratingLevelId").replace(/^rating:/, "")
  const level = model.ratingScale.find((candidate) => candidate.level === id)
  return level?.title || humanize(id)
}
const confidence = (analysis: Json) => display.confidence(field(analysis, "confidence")) || "—"
const summary = (analysis: Json) =>
  field(analysis, "rationale") || "No analysis summary was recorded."
const pair = (overall: Json, local: Json, fn: (item: Json) => string) =>
  `${fn(overall) || "—"} / ${fn(local) || "—"}`

const requestedScope = (manifest: ReportManifest) => {
  const area = field(manifest.requestedScope, "areaId")
  const filters = Array.isArray(manifest.requestedScope.factorFilter)
    ? manifest.requestedScope.factorFilter.map(string).filter(Boolean)
    : []
  if (area === "" && filters.length === 0) return "full evaluation"
  if (filters.length > 0) return `factor-scoped evaluation of ${filters.join(", ")}`
  return `area-scoped evaluation of ${area}`
}

const glossaryTarget = (reportPath: string, runRel: string) => {
  const runDepth = normalizeParts(runRel.split("/")).length
  return `${"../".repeat(normalizeParts(dirname(reportPath).split("/")).length + runDepth)}glossary.md`
}
const evaluationLinks = (reportPath: string, runRel: string) =>
  `> **Evaluation links:** ${link(reportPath, "report.md", "report.md")} | ${link(reportPath, "findings.md", "findings.md")} | ${link(reportPath, "recommendations.md", "recommendations.md")} | [glossary.md](${glossaryTarget(reportPath, runRel)})\n\n`
const runLine = (reportPath: string, manifest: ReportManifest) => {
  const label =
    reportPath === "report.md"
      ? manifest.run.label
      : link(reportPath, "report.md", manifest.run.label)
  return `Run: ${label} - Evaluation ID: ${code(manifest.evaluationId)} - Created: ${manifest.createdAt} - Scope: ${requestedScope(manifest)}\n\n`
}
const contents = (items: ReadonlyArray<string>) =>
  items.length < 2
    ? ""
    : `## Contents\n\n${items
        .map(
          (item) =>
            `- [${item}](#${item
              .toLowerCase()
              .replaceAll(/[^a-z0-9 -]/g, "")
              .replaceAll(" ", "-")})`,
        )
        .join("\n")}\n\n`
const sourceData = (reportPath: string, paths: ReadonlyArray<string>) =>
  `## Primary source data\n\n${[...new Set(paths)].map((path) => `- ${link(reportPath, path, path)}`).join("\n")}\n`

const header = (input: {
  type: string
  heading: string
  reportPath: string
  runRel: string
  manifest: ReportManifest
  context?: ReadonlyArray<string>
  summaryHeaders?: ReadonlyArray<string>
  summaryCells?: ReadonlyArray<string>
  contentLinks: ReadonlyArray<string>
}) => {
  let output = frontmatter(input.type, input.heading) + `# ${input.heading}\n\n`
  output += evaluationLinks(input.reportPath, input.runRel)
  output += runLine(input.reportPath, input.manifest)
  for (const context of input.context ?? []) if (context !== "") output += `${context}\n\n`
  if ((input.summaryHeaders?.length ?? 0) > 0) {
    output += "## Key details\n\n"
    output +=
      row(...input.summaryHeaders!) +
      row(...input.summaryHeaders!.map(() => "---")) +
      row(...input.summaryCells!) +
      "\n"
  }
  return output + contents(input.contentLinks)
}

const areaContext = (plan: EvaluationPlan, areaId: string, from: string) => {
  const path = splitArea(areaId)
  const ids = ["area:root", ...path.map((_, index) => `area:${path.slice(0, index + 1).join("/")}`)]
  const parts = ids.flatMap((id) => {
    const area = plan.areas.find((candidate) => candidate.ref === id)
    return area === undefined ? [] : [link(from, areaReportPath(id), titleForArea(area))]
  })
  return parts.length === 0 ? "" : `Area: ${parts.join(" / ")}`
}

const factorContext = (plan: EvaluationPlan, factor: PlannedFactor, from: string) =>
  `Factor: ${factor.path
    .map((_, index) => {
      const path = factor.path.slice(0, index + 1)
      const id = `factor:${splitArea(factor.areaId).join("/") || "root"}::${path.join("/")}`
      const item = plan.factors.find((candidate) => candidate.ref === id)
      return link(
        from,
        factorReportPath(id),
        item === undefined ? path.at(-1)! : titleForFactor(item),
      )
    })
    .join(" / ")}`

interface RankedFinding {
  readonly rank: number
  readonly tier: string
  readonly rationale: string
  readonly findingId: string
  readonly finding: Json
  readonly requirement?: PlannedRequirement
}
const rankedFindings = (plan: EvaluationPlan, index: ReturnType<typeof payloadIndex>) => {
  const found = new Map<string, { finding: Json; requirement: PlannedRequirement }>()
  for (const requirement of plan.requirements) {
    const assessment = index.one("RequirementAssessmentResult", "requirementId", requirement.ref)
    for (const finding of objects(assessment?.findings)) {
      const id = field(finding, "id")
      if (id !== "") found.set(`${requirement.ref}#${id}`, { finding, requirement })
    }
  }
  const ranking = index.one("FindingRankingResult")
  return objects(ranking?.orderedFindings).flatMap((item, position): Array<RankedFinding> => {
    const findingRef = object(item.findingRef)
    const selector = field(findingRef, "selector")
    const findingId = /^findings\[([^\]]+)]$/.exec(selector)?.[1] ?? field(item, "id")
    const requirementId = field(object(findingRef.subject), "requirementId")
    const match =
      found.get(`${requirementId}#${findingId}`) ??
      [...found.values()].find((entry) => field(entry.finding, "id") === findingId)
    return match === undefined
      ? []
      : [
          {
            rank: Number(item.rank ?? position + 1),
            tier: field(item, "tier"),
            rationale: field(item, "rationale"),
            findingId,
            ...match,
          },
        ]
  })
}

interface RankedRecommendation {
  readonly rank: number
  readonly recommendation: Json
  readonly ranking: Json
}
const rankedRecommendations = (index: ReturnType<typeof payloadIndex>) => {
  const recommendations = new Map(
    index.kind("RecommendationResult").map((item) => [field(item, "id"), item]),
  )
  return objects(index.one("RecommendationRankingResult")?.orderedRecommendations).flatMap(
    (ranking, position): Array<RankedRecommendation> => {
      const recommendation = recommendations.get(field(ranking, "recommendationRef"))
      return recommendation === undefined
        ? []
        : [{ rank: Number(ranking.rank ?? position + 1), recommendation, ranking }]
    },
  )
}

const slug = (text: string) =>
  text
    .toLowerCase()
    .replaceAll(/[^a-z0-9]+/g, "-")
    .replaceAll(/^-|-$/g, "") || "recommendation"
const recommendationPath = (item: RankedRecommendation) =>
  `recommendations/${String(item.rank).padStart(3, "0")}-${slug(field(item.recommendation, "title"))}.md`

const factorsForRequirement = (requirement: PlannedRequirement, from: string) =>
  requirement.factorIds.length === 0
    ? "—"
    : requirement.factorIds
        .map((id) => {
          return link(from, factorReportPath(id), splitFactor(id).path.join("/"))
        })
        .join("; ")

const areaContains = (parent: string, child: string) =>
  parent === "area:root" || child === parent || child.startsWith(`${parent}/`)

const factorContains = (parent: string, child: string) =>
  child === parent || child.startsWith(`${parent}/`)

const traceContexts = (plan: EvaluationPlan, recommendation: Json) =>
  objects(recommendation.traceRefs).flatMap((ref) => {
    const requirementId = field(object(ref.subject), "requirementId")
    const requirement = plan.requirements.find((item) => item.ref === requirementId)
    return requirement === undefined
      ? []
      : [{ areaId: requirement.areaId, factorIds: requirement.factorIds }]
  })

const recommendationMatchesArea = (plan: EvaluationPlan, recommendation: Json, areaId: string) =>
  traceContexts(plan, recommendation).some(
    (context) =>
      areaContains(areaId, context.areaId) ||
      context.factorIds.some((factorId) => areaContains(areaId, splitFactor(factorId).areaId)),
  )

const recommendationMatchesFactor = (
  plan: EvaluationPlan,
  recommendation: Json,
  factorId: string,
) =>
  traceContexts(plan, recommendation).some((context) =>
    context.factorIds.some((candidate) => factorContains(factorId, candidate)),
  )

const findingMatchesArea = (finding: RankedFinding, areaId: string) =>
  finding.requirement !== undefined && areaContains(areaId, finding.requirement.areaId)

const findingMatchesFactor = (finding: RankedFinding, factorId: string) =>
  finding.requirement?.factorIds.some((candidate) => factorContains(factorId, candidate)) ?? false

const recommendationAreaFactorLinks = (
  plan: EvaluationPlan,
  recommendation: Json,
  reportPath: string,
) => {
  const groups = new Map<string, Array<string>>()
  for (const context of traceContexts(plan, recommendation)) {
    const items = groups.get(context.areaId) ?? []
    for (const factorId of context.factorIds) if (!items.includes(factorId)) items.push(factorId)
    groups.set(context.areaId, items)
  }
  if (groups.size === 0) return "—"
  return [...groups]
    .map(([areaId, factorIds]) => {
      const area = plan.areas.find((candidate) => candidate.ref === areaId)
      const areaLink = link(
        reportPath,
        areaReportPath(areaId),
        area === undefined ? areaId : titleForArea(area),
      )
      const factorLinks = factorIds
        .map((factorId) => {
          const factor = plan.factors.find((candidate) => candidate.ref === factorId)
          return link(
            reportPath,
            factorReportPath(factorId),
            factor === undefined
              ? (splitFactor(factorId).path.at(-1) ?? factorId)
              : titleForFactor(factor),
          )
        })
        .join(", ")
      return `${areaLink} / ${factorLinks || "—"}`
    })
    .join("; ")
}

const limitTable = (...scopes: ReadonlyArray<Json>) => {
  let output = row("Type", "Scope", "Impact") + row("---", "---", "---")
  let count = 0
  for (const scope of scopes) {
    for (const [key, label] of [
      ["incompleteInputs", "🧩 Incomplete Inputs"],
      ["evaluationLimits", "⚠️ Evaluation Limits"],
    ] as const) {
      for (const item of objects(scope[key])) {
        count++
        output += row(
          label,
          field(item, "scope") || field(item, "ref") || field(item, "id"),
          field(item, "impact") || field(item, "description") || field(item, "reason"),
        )
      }
    }
  }
  return output + (count === 0 ? row("(no limits or incomplete inputs)", "—", "—") : "") + "\n"
}

const findingEffectSummary = (finding: Json) => {
  const effect = object(finding.effect)
  return field(effect, "statement") || field(effect, "ratingEffect")
}

const findingBasisSummary = (finding: Json) => {
  const basis = object(finding.basis)
  const status = display.basis(field(basis, "status"))
  const statement = field(basis, "statement")
  return status === "" ? statement : statement === "" ? status : `${status}: ${statement}`
}

const findingSection = (level: number, title: string, body: string) =>
  `${"#".repeat(level)} ${title}\n\n${body || "(not recorded)"}\n\n`

const evidenceSection = (level: number, title: string, evidence: ReadonlyArray<Json>) => {
  let output = `${"#".repeat(level)} ${title}\n\n`
  if (evidence.length === 0) return output + "(none recorded)\n\n"
  for (const item of evidence) {
    output += `- ${code(field(item, "sourceRef"))}: ${field(item, "statement")}\n`
    if (field(item, "rationale") !== "") output += `  Rationale: ${field(item, "rationale")}\n`
  }
  return output + "\n"
}

const findingDetails = (
  finding: Json,
  ranked: RankedFinding | undefined,
  total: number,
  fallbackId: string,
) => {
  const id = field(finding, "id") || fallbackId
  const statement = field(finding, "statement")
  let output = `<a id="finding-${id}"></a>\n\n### ${id}${statement === "" ? "" : ` ${statement}`}\n\n`
  output += row("Advice rank", "Tier", "Ranking rationale") + row("---", "---", "---")
  output +=
    ranked === undefined
      ? row("(not ranked)", "—", "—") + "\n"
      : row(`${ranked.rank} / ${total}`, display.tier(ranked.tier), ranked.rationale) + "\n"
  output += findingSection(4, "Condition", field(finding, "condition"))
  output += "#### Criteria\n\n"
  const criteria = objects(finding.criteria)
  if (criteria.length === 0) output += "(none recorded)\n\n"
  else {
    for (const criterion of criteria) {
      const label = [field(criterion, "requirementId"), field(criterion, "ratingLevelId")]
        .filter(Boolean)
        .join(" / ")
      output += `- ${code(label)}: ${field(criterion, "criterion")}\n`
      if (field(criterion, "rationale") !== "")
        output += `  Rationale: ${field(criterion, "rationale")}\n`
    }
    output += "\n"
  }
  const basis = object(finding.basis)
  output += "#### Basis\n\n"
  if (Object.keys(basis).length === 0) output += "(not recorded)\n\n"
  else {
    if (field(basis, "status") !== "")
      output += `Status: ${display.basis(field(basis, "status"))}\n\n`
    if (field(basis, "statement") !== "") output += `${field(basis, "statement")}\n\n`
    if (field(basis, "rationale") !== "") output += `Rationale: ${field(basis, "rationale")}\n\n`
    output += evidenceSection(5, "Basis evidence", objects(basis.evidence))
  }
  const effect = object(finding.effect)
  output += "#### Effect\n\n"
  if (Object.keys(effect).length === 0) output += "(not recorded)\n\n"
  else {
    if (field(effect, "statement") !== "") output += `${field(effect, "statement")}\n\n`
    if (field(effect, "ratingEffect") !== "")
      output += `Rating effect: ${field(effect, "ratingEffect")}\n\n`
    if (field(effect, "rationale") !== "") output += `Rationale: ${field(effect, "rationale")}\n\n`
  }
  return output + evidenceSection(4, "Evidence", objects(finding.evidence))
}

const areaFactorTable = (
  model: QualityModel,
  plan: EvaluationPlan,
  index: ReturnType<typeof payloadIndex>,
  areaId: string,
  reportPath: string,
  findings: ReadonlyArray<RankedFinding>,
  recommendations: ReadonlyArray<RankedRecommendation>,
) => {
  let output =
    row("▦ Area / □ Factor", "Overall rating", "Local rating", "Findings", "Recommendations") +
    row("---", "---", "---", "---", "---")
  const area = plan.areas.find((candidate) => candidate.ref === areaId)
  if (area === undefined) return output
  const writeFactor = (factor: PlannedFactor, depth: number) => {
    const result = index.one("FactorAnalysisResult", "factorId", factor.ref) ?? {}
    output += row(
      `${"↳ ".repeat(depth)}${link(reportPath, factorReportPath(factor.ref), `□ ${titleForFactor(factor)}`)}`,
      rating(model, scoped(result, "localAndDescendantAnalysis")),
      rating(model, scoped(result, "localAnalysis")),
      String(findings.filter((finding) => findingMatchesFactor(finding, factor.ref)).length),
      String(
        recommendations.filter((item) =>
          recommendationMatchesFactor(plan, item.recommendation, factor.ref),
        ).length,
      ),
    )
    for (const child of plan.factors.filter(
      (candidate) =>
        candidate.areaId === factor.areaId &&
        candidate.path.length === factor.path.length + 1 &&
        candidate.path.slice(0, -1).join("/") === factor.path.join("/"),
    ))
      writeFactor(child, depth + 1)
  }
  const writeArea = (current: PlannedArea, depth: number, root: boolean) => {
    const result = index.one("AreaAnalysisResult", "areaId", current.ref) ?? {}
    const title = `${"↳ ".repeat(depth)}${link(reportPath, areaReportPath(current.ref), `▦ ${titleForArea(current)}`)}`
    output += row(
      root ? `**${title}**` : title,
      rating(model, scoped(result, "localAndDescendantAnalysis")),
      rating(model, scoped(result, "localAnalysis")),
      String(findings.filter((finding) => findingMatchesArea(finding, current.ref)).length),
      String(
        recommendations.filter((item) =>
          recommendationMatchesArea(plan, item.recommendation, current.ref),
        ).length,
      ),
    )
    for (const factor of plan.factors.filter(
      (candidate) => candidate.areaId === current.ref && candidate.path.length === 1,
    ))
      writeFactor(factor, depth + 1)
    for (const childId of current.childAreaIds) {
      const child = plan.areas.find((candidate) => candidate.ref === childId)
      if (child !== undefined) writeArea(child, depth + 1, false)
    }
  }
  output =
    row("▦ Area / □ Factor", "Overall rating", "Local rating", "Findings", "Recommendations") +
    row("---", "---", "---", "---", "---")
  writeArea(area, 0, true)
  return output + "\n"
}

const findingsTable = (
  findings: ReadonlyArray<RankedFinding>,
  plan: EvaluationPlan,
  reportPath: string,
) => {
  let output =
    row("Rank", "Finding", "Area", "Factors", "Type", "Severity") +
    row("---", "---", "---", "---", "---", "---")
  if (findings.length === 0)
    return output + row("(no ranked findings)", "—", "—", "—", "—", "—") + "\n"
  for (const item of findings) {
    const requirement = item.requirement!
    const area = plan.areas.find((candidate) => candidate.ref === requirement.areaId)
    output += row(
      String(item.rank),
      link(
        reportPath,
        `${requirementReportPath(requirement.ref)}#finding-${item.findingId}`,
        field(item.finding, "statement") ||
          field(item.finding, "title") ||
          field(item.finding, "id"),
      ),
      link(
        reportPath,
        areaReportPath(requirement.areaId),
        area === undefined ? requirement.areaId : titleForArea(area),
      ),
      requirement.factorIds
        .map((id) => {
          const factor = plan.factors.find((candidate) => candidate.ref === id)
          return link(
            reportPath,
            factorReportPath(id),
            factor === undefined ? (splitFactor(id).path.at(-1) ?? id) : titleForFactor(factor),
          )
        })
        .join(", ") || "—",
      display.findingType(field(item.finding, "type")),
      display.severity(field(item.finding, "severity")) || "—",
    )
  }
  return output + "\n"
}

const markedCount = (marker: string, count: number, label: string, plural = false) =>
  `${marker} ${count} ${label}${plural && count !== 1 ? "s" : ""}`

const findingCountSummary = (findings: ReadonlyArray<RankedFinding>) => {
  const parts: Array<string> = []
  for (const [type, marker, label] of [
    ["gap", "🚩", "Gap"],
    ["risk", "⚠️", "Risk"],
    ["strength", "💪", "Strength"],
    ["note", "ℹ️", "Note"],
  ] as const) {
    const matching = findings.filter((item) => field(item.finding, "type") === type)
    if (matching.length === 0) continue
    let part = markedCount(marker, matching.length, label, true)
    if (type === "gap" || type === "risk") {
      const severities = (
        [
          ["critical", "🔴", "Critical"],
          ["high", "🔴", "High"],
          ["medium", "🟡", "Medium"],
          ["low", "🔵", "Low"],
        ] as const
      ).flatMap(([severity, severityMarker, severityLabel]) => {
        const count = matching.filter((item) => field(item.finding, "severity") === severity).length
        return count === 0 ? [] : [markedCount(severityMarker, count, severityLabel)]
      })
      if (severities.length > 0) part += `: ${severities.join(", ")}`
    }
    parts.push(part)
  }
  return parts.join("; ")
}

const recommendationsTable = (
  items: ReadonlyArray<RankedRecommendation>,
  plan: EvaluationPlan,
  reportPath: string,
  includeRationale: boolean,
) => {
  const headers = [
    "#",
    "Recommendation",
    "Area / factors",
    "Impact",
    "Confidence",
    "Reason",
    ...(includeRationale ? ["Ranking rationale"] : []),
  ]
  let output = row(...headers) + row(...headers.map(() => "---"))
  if (items.length === 0)
    return output + row("(no ranked recommendations)", ...headers.slice(1).map(() => "—")) + "\n"
  for (const item of items) {
    const cells = [
      String(item.rank),
      link(reportPath, recommendationPath(item), field(item.recommendation, "title")),
      recommendationAreaFactorLinks(plan, item.recommendation, reportPath),
      display.impact(field(item.recommendation, "impact")),
      display.confidence(
        field(includeRationale ? item.recommendation : item.ranking, "confidence"),
      ) || "—",
      field(item.recommendation, "expectedValue") || "—",
    ]
    if (includeRationale) cells.push(field(item.ranking, "rationale") || "—")
    output += row(...cells)
  }
  return output + "\n"
}

const renderRun = (
  model: QualityModel,
  manifest: ReportManifest,
  plan: EvaluationPlan,
  index: ReturnType<typeof payloadIndex>,
  runRel: string,
  findings: ReadonlyArray<RankedFinding>,
  recommendations: ReadonlyArray<RankedRecommendation>,
) => {
  const scopedArea =
    plan.areas.find((area) => area.ref === manifest.plannedScope.areaId) ?? plan.areas[0]!
  const analysis = index.one("AreaAnalysisResult", "areaId", scopedArea.ref) ?? {}
  const overall = scoped(analysis, "localAndDescendantAnalysis")
  const heading =
    manifest.requestedScope.areaId === undefined
      ? `Quality evaluation - ${titleForArea(scopedArea)}`
      : `Quality evaluation - ${titleForArea(scopedArea)}`
  let output = frontmatter("Evaluation Overview Report", heading, [
    ["evaluationId", manifest.evaluationId],
    ["created", manifest.createdAt],
    ["model", manifest.model],
    ["run", manifest.run.label],
  ])
  output += `# ${heading}\n\n${evaluationLinks("report.md", runRel)}`
  output += `## Summary\n\n${summary(overall)}\n\n## Key details\n\n`
  output +=
    row("Overall rating", "Confidence", "Scope", "Findings", "Recommendations") +
    row("---", "---", "---", "---", "---") +
    row(
      rating(model, overall),
      confidence(overall),
      manifest.requestedScope.areaId === undefined
        ? `Full evaluation of ${titleForArea(scopedArea)}`
        : requestedScope(manifest),
      `${findings.length} total`,
      `${recommendations.length} total`,
    ) +
    "\n"
  output += contents([
    "Model evaluation",
    "Top findings",
    "Top recommendations",
    "Primary source data",
  ])
  output += `## Model evaluation\n\n${areaFactorTable(model, plan, index, scopedArea.ref, "report.md", findings, recommendations)}`
  const findingSummary = findingCountSummary(findings)
  output += `## Top findings\n\n**Full findings report:** [findings.md](findings.md) (${findings.length} total${findingSummary === "" ? "" : `: ${findingSummary}`})\n\n${findingsTable(findings.slice(0, 10), plan, "report.md")}`
  const impactCounts = recommendations.reduce<Record<string, number>>((acc, item) => {
    const key = display.impact(field(item.recommendation, "impact"))
    acc[key] = (acc[key] ?? 0) + 1
    return acc
  }, {})
  const impact = Object.entries(impactCounts)
    .map(([key, count]) => {
      const [marker = "", ...label] = key.split(" ")
      return `${marker} ${count} ${label.join(" ")}`
    })
    .join(", ")
  output += `## Top recommendations\n\n**Full recommendations report:** [recommendations.md](recommendations.md) (${recommendations.length} total${impact === "" ? "" : `; impact: ${impact}`})\n\n${recommendationsTable(recommendations.slice(0, 10), plan, "report.md", false)}`
  return (
    output +
    sourceData("report.md", [
      "data/evaluation-manifest.json",
      dataPathForPayload("AreaAnalysisResult", analysis),
      "data/advice/finding-ranking-result.json",
      "data/advice/recommendation-ranking-result.json",
    ])
  )
}

const renderArea = (
  model: QualityModel,
  manifest: ReportManifest,
  plan: EvaluationPlan,
  area: PlannedArea,
  index: ReturnType<typeof payloadIndex>,
  runRel: string,
  recommendations: ReadonlyArray<RankedRecommendation>,
) => {
  const reportPath = areaReportPath(area.ref)
  const analysis = index.one("AreaAnalysisResult", "areaId", area.ref) ?? {}
  const local = scoped(analysis, "localAnalysis")
  const overall = scoped(analysis, "localAndDescendantAnalysis")
  let output = header({
    type: "Area Evaluation Report",
    heading: `Area: ${titleForArea(area)}`,
    reportPath,
    runRel,
    manifest,
    context: [areaContext(plan, area.ref, reportPath)],
    summaryHeaders: ["Overall rating", "Local rating", "Confidence"],
    summaryCells: [rating(model, overall), rating(model, local), pair(overall, local, confidence)],
    contentLinks: [
      "Summary",
      "Area / factor breakdown",
      "Requirements",
      "Limits and incomplete inputs",
      "Primary source data",
    ],
  })
  const findings = rankedFindings(plan, index)
  output += `## Summary\n\n${summary(overall)}\n\n## Area / factor breakdown\n\n${areaFactorTable(model, plan, index, area.ref, reportPath, findings, recommendations)}`
  output +=
    "## Requirements\n\n" +
    row("Requirement", "Rating", "Status", "Factors") +
    row("---", "---", "---", "---")
  const requirements = plan.requirements.filter((candidate) => candidate.areaId === area.ref)
  if (requirements.length === 0) output += row("(no local requirements)", "—", "—", "—")
  else
    for (const requirement of requirements) {
      const ratingResult =
        index.one("RequirementRatingResult", "requirementId", requirement.ref) ?? {}
      const assessment =
        index.one("RequirementAssessmentResult", "requirementId", requirement.ref) ?? {}
      output += row(
        link(reportPath, requirementReportPath(requirement.ref), titleForRequirement(requirement)),
        rating(model, ratingResult),
        display.assessment(field(assessment, "status")),
        factorsForRequirement(requirement, reportPath),
      )
    }
  output += "\n## Limits and incomplete inputs\n\n" + limitTable(local, overall)
  const sources = [
    "data/evaluation-manifest.json",
    dataPathForPayload("AreaAnalysisResult", analysis),
    "data/advice/finding-ranking-result.json",
    "data/advice/recommendation-ranking-result.json",
  ]
  for (const requirement of requirements) {
    const ratingResult = index.one("RequirementRatingResult", "requirementId", requirement.ref)
    const assessment = index.one("RequirementAssessmentResult", "requirementId", requirement.ref)
    if (ratingResult) sources.push(dataPathForPayload("RequirementRatingResult", ratingResult))
    if (assessment) sources.push(dataPathForPayload("RequirementAssessmentResult", assessment))
  }
  return output + sourceData(reportPath, sources)
}

const renderFactor = (
  model: QualityModel,
  manifest: ReportManifest,
  plan: EvaluationPlan,
  factor: PlannedFactor,
  index: ReturnType<typeof payloadIndex>,
  runRel: string,
) => {
  const reportPath = factorReportPath(factor.ref)
  const analysis = index.one("FactorAnalysisResult", "factorId", factor.ref) ?? {}
  const local = scoped(analysis, "localAnalysis")
  const overall = scoped(analysis, "localAndDescendantAnalysis")
  let output = header({
    type: "Factor Evaluation Report",
    heading: `Factor: ${titleForFactor(factor)}`,
    reportPath,
    runRel,
    manifest,
    context: [
      areaContext(plan, factor.areaId, reportPath),
      factorContext(plan, factor, reportPath),
    ],
    summaryHeaders: ["Overall rating", "Local rating", "Status", "Confidence"],
    summaryCells: [
      rating(model, overall),
      rating(model, local),
      pair(overall, local, (item) => display.analysis(field(item, "status"))),
      pair(overall, local, confidence),
    ],
    contentLinks: [
      "Summary",
      "Requirements",
      "Sub-factors",
      "Limits and incomplete inputs",
      "Primary source data",
    ],
  })
  output +=
    `## Summary\n\n${summary(overall)}\n\n## Requirements\n\n` +
    row("Requirement", "Rating", "Status") +
    row("---", "---", "---")
  const requirements = plan.requirements.filter((candidate) =>
    candidate.factorIds.includes(factor.ref),
  )
  if (requirements.length === 0) output += row("(no direct requirements)", "—", "—")
  else
    for (const requirement of requirements) {
      const ratingResult =
        index.one("RequirementRatingResult", "requirementId", requirement.ref) ?? {}
      const assessment =
        index.one("RequirementAssessmentResult", "requirementId", requirement.ref) ?? {}
      output += row(
        link(reportPath, requirementReportPath(requirement.ref), titleForRequirement(requirement)),
        rating(model, ratingResult),
        display.assessment(field(assessment, "status")),
      )
    }
  output +=
    "\n## Sub-factors\n\n" +
    row("Factor", "Path", "Local rating", "+ Sub-factors rating") +
    row("---", "---", "---", "---")
  const children = plan.factors.filter(
    (candidate) =>
      candidate.areaId === factor.areaId &&
      candidate.path.length === factor.path.length + 1 &&
      candidate.path.slice(0, -1).join("/") === factor.path.join("/"),
  )
  if (children.length === 0) output += row("(no sub-factors)", "—", "—", "—")
  else
    for (const child of children) {
      const result = index.one("FactorAnalysisResult", "factorId", child.ref) ?? {}
      output += row(
        link(reportPath, factorReportPath(child.ref), titleForFactor(child)),
        code(
          `${splitArea(child.areaId).join("/")}${child.areaId === "area:root" ? "" : "::"}${child.path.join("/")}`,
        ),
        rating(model, scoped(result, "localAnalysis")),
        plan.factors.some(
          (candidate) =>
            candidate.areaId === child.areaId &&
            candidate.path.length === child.path.length + 1 &&
            candidate.path.slice(0, -1).join("/") === child.path.join("/"),
        )
          ? rating(model, scoped(result, "localAndDescendantAnalysis"))
          : "—",
      )
    }
  output += "\n## Limits and incomplete inputs\n\n" + limitTable(local, overall)
  const sources = [
    "data/evaluation-manifest.json",
    dataPathForPayload("FactorAnalysisResult", analysis),
  ]
  for (const requirement of requirements) {
    const rr = index.one("RequirementRatingResult", "requirementId", requirement.ref)
    const ra = index.one("RequirementAssessmentResult", "requirementId", requirement.ref)
    if (rr) sources.push(dataPathForPayload("RequirementRatingResult", rr))
    if (ra) sources.push(dataPathForPayload("RequirementAssessmentResult", ra))
  }
  return output + sourceData(reportPath, sources)
}

const renderRequirement = (
  model: QualityModel,
  manifest: ReportManifest,
  plan: EvaluationPlan,
  requirement: PlannedRequirement,
  index: ReturnType<typeof payloadIndex>,
  runRel: string,
) => {
  const reportPath = requirementReportPath(requirement.ref)
  const assessment =
    index.one("RequirementAssessmentResult", "requirementId", requirement.ref) ?? {}
  const ratingResult = index.one("RequirementRatingResult", "requirementId", requirement.ref) ?? {}
  let output = header({
    type: "Requirement Evaluation Report",
    heading: `Requirement: ${titleForRequirement(requirement)}`,
    reportPath,
    runRel,
    manifest,
    context: [
      areaContext(plan, requirement.areaId, reportPath),
      `Factors: ${factorsForRequirement(requirement, reportPath)}`,
    ],
    summaryHeaders: ["Rating", "Assessment", "Confidence"],
    summaryCells: [
      rating(model, ratingResult),
      display.assessment(field(assessment, "status")),
      pair(ratingResult, assessment, confidence),
    ],
    contentLinks: [
      "Summary",
      "Findings summary",
      "Finding details",
      "Unknowns and missing evidence",
      "Primary source data",
    ],
  })
  output += `## Summary\n\n${field(assessment, "evidenceSummary") || field(ratingResult, "rationale") || "No assessment summary was recorded."}\n\n## Findings summary\n\n`
  output +=
    row("ID", "Statement", "Type", "Severity", "Confidence", "Effect", "Basis") +
    row("---", "---", "---", "---", "---", "---", "---")
  const findings = objects(assessment.findings)
  if (findings.length === 0) output += row("(no findings)", "—", "—", "—", "—", "—", "—")
  else
    for (const finding of findings)
      output += row(
        code(field(finding, "id")),
        field(finding, "statement"),
        display.findingType(field(finding, "type")),
        display.severity(field(finding, "severity")) || "—",
        display.confidence(field(finding, "confidence")) || "—",
        findingEffectSummary(finding),
        findingBasisSummary(finding),
      )
  output += "\n## Finding details\n\n"
  if (findings.length === 0) output += "(no finding details)\n\n"
  else {
    const ranked = rankedFindings(plan, index)
    for (const [position, finding] of findings.entries()) {
      const id = field(finding, "id") || `finding-${position + 1}`
      output += findingDetails(
        finding,
        ranked.find((item) => item.requirement?.ref === requirement.ref && item.findingId === id),
        ranked.length,
        id,
      )
    }
  }
  output += "## Unknowns and missing evidence\n\n" + row("Type", "Detail") + row("---", "---")
  const unknowns = [
    ...objects(assessment.unknowns).map((item) => [
      "❓ Unknowns",
      field(item, "description") ||
        field(item, "reason") ||
        field(item, "ref") ||
        field(item, "id"),
    ]),
    ...objects(ratingResult.missingEvidence).map((item) => [
      "🔎 Missing Evidence",
      field(item, "description") ||
        field(item, "reason") ||
        field(item, "ref") ||
        field(item, "id"),
    ]),
  ]
  if (unknowns.length === 0) output += row("(none recorded)", "—")
  else for (const item of unknowns) output += row(item[0]!, item[1]!)
  return (
    output +
    "\n" +
    sourceData(reportPath, [
      "data/evaluation-manifest.json",
      dataPathForPayload("RequirementAssessmentResult", assessment),
      dataPathForPayload("RequirementRatingResult", ratingResult),
      "data/advice/finding-ranking-result.json",
    ])
  )
}

const renderFindings = (
  manifest: ReportManifest,
  plan: EvaluationPlan,
  runRel: string,
  findings: ReadonlyArray<RankedFinding>,
) => {
  const reportPath = "findings.md"
  let output = header({
    type: "Finding Index Report",
    heading: "Findings",
    reportPath,
    runRel,
    manifest,
    summaryHeaders: ["Findings", "Highest concern severity"],
    summaryCells: [
      `${findings.length} findings`,
      findings.length === 0 ? "—" : display.severity(field(findings[0]!.finding, "severity")),
    ],
    contentLinks: ["Ranked findings", "Primary source data"],
  })
  output += `## Ranked findings\n\n${findingsTable(findings, plan, reportPath)}`
  return (
    output +
    sourceData(reportPath, [
      "data/evaluation-manifest.json",
      "data/advice/finding-ranking-result.json",
      ...findings.flatMap((item) =>
        item.requirement === undefined
          ? []
          : [
              dataPathForPayload("RequirementAssessmentResult", {
                kind: "RequirementAssessmentResult",
                requirementId: item.requirement.ref,
              }),
            ],
      ),
    ])
  )
}

const renderRecommendationIndex = (
  manifest: ReportManifest,
  plan: EvaluationPlan,
  index: ReturnType<typeof payloadIndex>,
  runRel: string,
  recommendations: ReadonlyArray<RankedRecommendation>,
) => {
  const reportPath = "recommendations.md"
  const highest =
    recommendations.length === 0
      ? "—"
      : display.impact(field(recommendations[0]!.recommendation, "impact"))
  const coverage = objects(index.one("RecommendationRankingResult")?.findingCoverage)
  const addressed = coverage.filter(
    (item) => field(item, "disposition") === "addressed_by_recommendation",
  ).length
  const notDriving = coverage.filter(
    (item) => field(item, "disposition") === "not_advice_driving",
  ).length
  let output = header({
    type: "Recommendation Index Report",
    heading: "Recommendations",
    reportPath,
    runRel,
    manifest,
    summaryHeaders: ["Recommendations", "Highest impact", "Coverage"],
    summaryCells: [
      `${recommendations.length} recommendations`,
      highest,
      `✅ Addressed by Recommendation: ${addressed} / ⬜ Not Advice Driving: ${notDriving}`,
    ],
    contentLinks: ["Ranked recommendations", "Coverage", "Primary source data"],
  })
  output += `## Ranked recommendations\n\n${recommendationsTable(recommendations, plan, reportPath, true)}## Coverage\n\n- ✅ Addressed by Recommendation: ${addressed}\n- ⬜ Not Advice Driving: ${notDriving}\n\n`
  return (
    output +
    sourceData(reportPath, [
      "data/evaluation-manifest.json",
      "data/advice/recommendation-ranking-result.json",
      ...recommendations.map((item) =>
        dataPathForPayload("RecommendationResult", item.recommendation),
      ),
    ])
  )
}

const renderRecommendation = (
  manifest: ReportManifest,
  item: RankedRecommendation,
  runRel: string,
) => {
  const recommendation = item.recommendation
  const reportPath = recommendationPath(item)
  const id = field(recommendation, "id")
  let output = header({
    type: "Recommendation Report",
    heading: `Recommendation: ${field(recommendation, "title")}`,
    reportPath,
    runRel,
    manifest,
    summaryHeaders: ["#", "ID", "Impact", "Confidence", "Reference"],
    summaryCells: [
      String(item.rank),
      code(id),
      display.impact(field(recommendation, "impact")),
      display.confidence(field(recommendation, "confidence")),
      code(`evaluation:${manifest.evaluationId}/recommendation/${id}`),
    ],
    contentLinks: [
      "Description",
      "Background",
      "Expected value",
      "Done criterion",
      "Ranking rationale",
      "Trace",
      "Primary source data",
    ],
  })
  for (const [heading, key] of [
    ["Description", "description"],
    ["Background", "background"],
    ["Expected value", "expectedValue"],
    ["Done criterion", "doneCriterion"],
  ] as const)
    output += `## ${heading}\n\n${field(recommendation, key) || "(not recorded)"}\n\n`
  output += `## Ranking rationale\n\n${field(item.ranking, "rationale") || "(not recorded)"}\n\n## Trace\n\n`
  const refs = objects(recommendation.traceRefs)
  output +=
    refs.length === 0
      ? "(none recorded)\n\n"
      : refs.map((ref) => `- ${code(JSON.stringify(ref))}`).join("\n") + "\n\n"
  return (
    output +
    sourceData(reportPath, [
      "data/evaluation-manifest.json",
      dataPathForPayload("RecommendationResult", recommendation),
      "data/advice/recommendation-ranking-result.json",
    ])
  )
}

const ref = (kind: string, path: string, extra: Json = {}) => ({
  ...(extra.areaId === undefined ? {} : { areaId: extra.areaId }),
  ...(extra.factorId === undefined ? {} : { factorId: extra.factorId }),
  kind,
  path,
  ...(extra.recommendationId === undefined ? {} : { recommendationId: extra.recommendationId }),
  ...(extra.requirementId === undefined ? {} : { requirementId: extra.requirementId }),
})
const routineRef = (kind: string, subject: Json, selector = "") => ({
  kind,
  ...(selector === "" ? {} : { selector }),
  subject,
})

export const buildReportTree = (input: {
  readonly model: QualityModel
  readonly manifest: ReportManifest
  readonly plan: EvaluationPlan
  readonly payloads: ReadonlyArray<Json>
  readonly runRel: string
}): BuiltReportTree => {
  const index = payloadIndex(input.payloads)
  const reportPlan: EvaluationPlan = {
    ...input.plan,
    requirements: input.plan.requirements
      .map((requirement) => {
        const frame = index
          .kind("RequirementEvaluationFrame")
          .find(
            (candidate) => field(object(candidate.subject), "requirementId") === requirement.ref,
          )
        const rawFrameFactors = object(frame?.subject).factorIds
        const frameFactors = Array.isArray(rawFrameFactors)
          ? rawFrameFactors.map(string).filter(Boolean)
          : []
        const assessment = index.one(
          "RequirementAssessmentResult",
          "requirementId",
          requirement.ref,
        )
        const assessmentFactors = Array.isArray(assessment?.factors)
          ? assessment.factors.map(string).filter(Boolean)
          : []
        return {
          ...requirement,
          factorIds:
            frameFactors.length > 0
              ? frameFactors
              : assessmentFactors.length > 0
                ? assessmentFactors
                : requirement.factorIds,
        }
      })
      .sort((left, right) => {
        const leftArea = input.plan.areas.findIndex((area) => area.ref === left.areaId)
        const rightArea = input.plan.areas.findIndex((area) => area.ref === right.areaId)
        return leftArea === rightArea ? left.ref.localeCompare(right.ref) : leftArea - rightArea
      }),
  }
  const scopedAnalysis = index.one(
    "AreaAnalysisResult",
    "areaId",
    input.manifest.plannedScope.areaId,
  )
  if (scopedAnalysis === undefined)
    throw new Error(
      `run is not reportable: missing-evaluation-data ${input.manifest.plannedScope.areaId}: required scoped area analysis payload is missing`,
    )
  for (const kind of ["EvaluationFrame", "FindingRankingResult", "RecommendationRankingResult"])
    if (index.one(kind) === undefined)
      throw new Error(
        `run is not reportable: missing-evaluation-data ${kind}: required evaluation payload is missing`,
      )
  const findings = rankedFindings(reportPlan, index)
  const recommendations = rankedRecommendations(index)
  const reports: Array<RenderedReport> = [
    {
      kind: "run",
      path: "report.md",
      content: renderRun(
        input.model,
        input.manifest,
        reportPlan,
        index,
        input.runRel,
        findings,
        recommendations,
      ),
    },
  ]
  for (const area of reportPlan.areas)
    if (index.one("AreaAnalysisResult", "areaId", area.ref))
      reports.push({
        kind: "area",
        path: areaReportPath(area.ref),
        areaId: area.ref,
        content: renderArea(
          input.model,
          input.manifest,
          reportPlan,
          area,
          index,
          input.runRel,
          recommendations,
        ),
      })
  for (const factor of reportPlan.factors)
    if (index.one("FactorAnalysisResult", "factorId", factor.ref))
      reports.push({
        kind: "factor",
        path: factorReportPath(factor.ref),
        areaId: factor.areaId,
        factorId: factor.ref,
        content: renderFactor(input.model, input.manifest, reportPlan, factor, index, input.runRel),
      })
  for (const requirement of reportPlan.requirements)
    if (
      index.one("RequirementAssessmentResult", "requirementId", requirement.ref) ||
      index.one("RequirementRatingResult", "requirementId", requirement.ref)
    )
      reports.push({
        kind: "requirement",
        path: requirementReportPath(requirement.ref),
        areaId: requirement.areaId,
        requirementId: requirement.ref,
        content: renderRequirement(
          input.model,
          input.manifest,
          reportPlan,
          requirement,
          index,
          input.runRel,
        ),
      })
  reports.push({
    kind: "findings",
    path: "findings.md",
    content: renderFindings(input.manifest, reportPlan, input.runRel, findings),
  })
  reports.push({
    kind: "recommendations",
    path: "recommendations.md",
    content: renderRecommendationIndex(
      input.manifest,
      reportPlan,
      index,
      input.runRel,
      recommendations,
    ),
  })
  for (const item of recommendations)
    reports.push({
      kind: "recommendation",
      path: recommendationPath(item),
      recommendationId: field(item.recommendation, "id"),
      content: renderRecommendation(input.manifest, item, input.runRel),
    })
  const reportRefs: Array<Json> = reports.map((report) =>
    ref(report.kind, report.path, {
      ...(report.areaId ? { areaId: report.areaId } : {}),
      ...(report.factorId ? { factorId: report.factorId } : {}),
      ...(report.requirementId ? { requirementId: report.requirementId } : {}),
      ...(report.recommendationId ? { recommendationId: report.recommendationId } : {}),
    }),
  )
  const areaOutputs = reportPlan.areas.flatMap((area) => {
    if (!index.one("AreaAnalysisResult", "areaId", area.ref)) return []
    const factors = reportPlan.factors.filter(
      (factor) =>
        factor.areaId === area.ref &&
        factor.path.length === 1 &&
        index.one("FactorAnalysisResult", "factorId", factor.ref),
    )
    const requirements = reportPlan.requirements.filter(
      (requirement) => requirement.areaId === area.ref,
    )
    return [
      {
        areaAnalysisResultRef: routineRef("AreaAnalysisResult", { areaId: area.ref }),
        areaEvaluationFrameRef: routineRef("AreaEvaluationFrame", { areaId: area.ref }),
        areaId: area.ref,
        factorAnalysisRefs: factors.map((factor) =>
          routineRef(
            "FactorAnalysisResult",
            { factorId: factor.ref },
            "localAndDescendantAnalysis",
          ),
        ),
        reportRefs: reportRefs.filter((report) => report.areaId === area.ref),
        requirementAssessmentRefs: requirements
          .filter((requirement) =>
            index.one("RequirementAssessmentResult", "requirementId", requirement.ref),
          )
          .map((requirement) =>
            routineRef("RequirementAssessmentResult", { requirementId: requirement.ref }),
          ),
        requirementRatingRefs: requirements
          .filter((requirement) =>
            index.one("RequirementRatingResult", "requirementId", requirement.ref),
          )
          .map((requirement) =>
            routineRef("RequirementRatingResult", { requirementId: requirement.ref }),
          ),
      },
    ]
  })
  const output: Json = {
    areaOutputs,
    kind: "EvaluationOutputResult",
    reportOutputs: reportRefs,
    ...(index.one("AreaAnalysisResult", "areaId", "area:root")
      ? {
          rootAreaAnalysisRef: routineRef(
            "AreaAnalysisResult",
            { areaId: "area:root" },
            "localAndDescendantAnalysis",
          ),
        }
      : {}),
    runReportRef: ref("run", "report.md"),
    schemaVersion: 3,
    scopedAreaAnalysisRef: routineRef(
      "AreaAnalysisResult",
      { areaId: input.manifest.plannedScope.areaId },
      "localAndDescendantAnalysis",
    ),
  }
  const scopedResult = scoped(scopedAnalysis, "localAndDescendantAnalysis")
  const ratingResult: Json =
    field(scopedResult, "status") === "analyzed"
      ? {
          kind: "rated",
          level: field(scopedResult, "ratingLevelId").replace(/^rating:/, ""),
          rationale: field(scopedResult, "rationale"),
        }
      : { kind: "not_assessed", rationale: field(scopedResult, "rationale") }
  const payloadFiles = input.payloads.map((payload) => ({
    path: dataPathForPayload(resolveDataKind(field(payload, "kind")), payload),
    payload,
  }))
  return { reports, output, rating: ratingResult, payloadFiles }
}
