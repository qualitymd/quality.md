import type { ParsedNode } from "yaml"

import {
  isEmptyNode,
  isMapNode,
  isScalarNode,
  isSequenceNode,
  mapEntries,
  mapEntry,
  nodeValue,
  type QualityDocument,
  type YamlMap,
  type YamlNode,
  type YamlPair,
} from "../model/document.ts"
import {
  Area,
  Factor,
  hasSchemaProperty,
  Model,
  ModelNamePattern,
  ModelNamePatternText,
  property,
  RatingLevel,
  Requirement,
  schemaProperty,
  type Property,
  type SchemaNode,
} from "../model/schema.ts"
import {
  RuleId,
  RulesById,
  type Finding,
  type LintResult,
  type Location,
  type PathSegment,
  type RepairRecord,
  type Severity,
} from "./result.ts"

interface AreaRef {
  readonly name: string
  readonly node: YamlMap
  readonly path: ReadonlyArray<PathSegment>
  readonly parent?: AreaRef
  factors: Array<FactorRef>
  requirements: Array<RequirementRef>
  areas: Array<AreaRef>
}

interface FactorRef {
  readonly name: string
  readonly node: YamlMap
  readonly path: ReadonlyArray<PathSegment>
  readonly area: AreaRef
  readonly parent?: FactorRef
  factors: Array<FactorRef>
  requirements: Array<RequirementRef>
}

interface RequirementRef {
  readonly name: string
  readonly node: YamlMap
  readonly path: ReadonlyArray<PathSegment>
  readonly area: AreaRef
  readonly factor?: FactorRef
}

interface RepairOperation {
  readonly record: RepairRecord
  readonly apply: () => boolean
}

interface SchemaContext {
  readonly areaName?: string
  readonly factorName?: string
  readonly requirement?: string
}

export interface LintRun {
  readonly result: LintResult
  readonly applyRepairs: () => ReadonlyArray<RepairRecord>
}

const appendPath = (
  path: ReadonlyArray<PathSegment>,
  ...parts: ReadonlyArray<PathSegment>
): Array<PathSegment> => [...path, ...parts]

const pathLabel = (path: ReadonlyArray<PathSegment>) =>
  path.length === 0 ? "frontmatter" : path.map(String).join(".")

const severityRank = (severity: Severity) =>
  severity === "error" ? 0 : severity === "warning" ? 1 : severity === "info" ? 2 : 3

const comparePath = (left: ReadonlyArray<PathSegment>, right: ReadonlyArray<PathSegment>) => {
  const length = Math.min(left.length, right.length)
  for (let index = 0; index < length; index += 1) {
    const a = left[index]
    const b = right[index]
    if (typeof a === "number" && typeof b === "number" && a !== b) return a - b
    const compared = String(a).localeCompare(String(b))
    if (compared !== 0) return compared
  }
  return left.length - right.length
}

const compareLocations = (left: Location, right: Location) => {
  if (left.line !== undefined && right.line !== undefined) {
    if (left.line !== right.line) return left.line - right.line
    if (left.column !== right.column) return (left.column ?? 0) - (right.column ?? 0)
  }
  return comparePath(left.modelPath, right.modelPath) || left.label.localeCompare(right.label)
}

const compareFindings = (left: Finding, right: Finding) =>
  compareLocations(left.location, right.location) ||
  severityRank(left.severity) - severityRank(right.severity) ||
  left.ruleId.localeCompare(right.ruleId)

const nextActions = (path: string, errors: number, fixable: number) => {
  if (fixable > 0) {
    return [
      {
        id: "fix",
        label: "Apply deterministic lint repairs",
        command: `qualitymd lint --fix ${path}`,
      },
    ]
  }
  if (errors > 0) {
    return [
      {
        id: "rerun-lint",
        label: "Re-run validation",
        command: `qualitymd lint ${path}`,
      },
    ]
  }
  return []
}

class RunState {
  private readonly findings: Array<Finding> = []
  private readonly repairs: Array<RepairOperation> = []
  private readonly levels = new Set<string>()
  private root!: AreaRef

  constructor(private readonly document: QualityDocument) {}

  run(): LintRun {
    if (!isMapNode(this.document.frontmatter)) {
      this.add(
        RuleId.invalidFrontmatter,
        "The frontmatter is not a model mapping; a QUALITY.md frontmatter block must be a map of model properties.",
        this.location(this.document.frontmatter, [], "frontmatter"),
      )
      return this.finish()
    }
    this.root = {
      name: "root",
      node: this.document.frontmatter,
      path: [],
      factors: [],
      requirements: [],
      areas: [],
    }
    this.checkRoot()
    this.checkRatingScale()
    this.walkModel()
    this.checkEmptyModel()
    this.checkAreas(this.root)
    this.checkFactors(this.root)
    this.checkRequirements(this.root)
    return this.finish()
  }

  private finish(): LintRun {
    this.findings.sort(compareFindings)
    this.repairs.sort((left, right) =>
      compareLocations(left.record.location, right.record.location),
    )
    const result = this.result([])
    return {
      result,
      applyRepairs: () => {
        const records: Array<RepairRecord> = []
        for (const repair of this.repairs) {
          if (repair.apply()) records.push(repair.record)
        }
        return records
      },
    }
  }

  result(repairs: ReadonlyArray<RepairRecord>): LintResult {
    let errors = 0
    let warnings = 0
    let info = 0
    let fixable = 0
    for (const finding of this.findings) {
      if (finding.severity === "error") errors += 1
      else if (finding.severity === "warning") warnings += 1
      else info += 1
      if (finding.fixable) fixable += 1
    }
    return {
      schemaVersion: 1,
      path: this.document.path,
      valid: errors === 0,
      summary: { errors, warnings, info, fixable, fixed: repairs.length },
      findings: this.findings,
      repairs,
      nextActions: nextActions(this.document.path, errors, fixable),
    }
  }

  private add(ruleId: RuleId, message: string, location: Location, repair?: RepairOperation) {
    const rule = RulesById.get(ruleId)
    if (rule === undefined) throw new Error(`unknown lint rule: ${ruleId}`)
    this.findings.push({
      ruleId,
      severity: rule.severity,
      message,
      location,
      fixable: rule.fixable,
    })
    if (repair !== undefined) this.repairs.push(repair)
  }

  private invalid(
    node: YamlNode,
    path: ReadonlyArray<PathSegment>,
    label: string,
    message: string,
  ) {
    this.add(RuleId.invalidFrontmatter, message, this.location(node, path, label))
  }

  private location(node: YamlNode, modelPath: ReadonlyArray<PathSegment>, label: string): Location {
    const range = node?.range
    if (range === undefined || range === null) {
      return { path: this.document.path, modelPath: [...modelPath], label }
    }
    const position = this.document.lineCounter.linePos(range[0])
    return {
      path: this.document.path,
      modelPath: [...modelPath],
      label,
      line: position.line + 1,
      column: position.col,
    }
  }

  private locationForMissing(modelPath: ReadonlyArray<PathSegment>, label: string): Location {
    return { path: this.document.path, modelPath: [...modelPath], label }
  }

  private locationForNodeOrMissing(
    node: YamlNode | undefined,
    modelPath: ReadonlyArray<PathSegment>,
    label: string,
  ) {
    return node === undefined
      ? this.locationForMissing(modelPath, label)
      : this.location(node, modelPath, label)
  }

  private checkRoot() {
    this.checkSchemaProperties(Model, this.document.frontmatter, [], {})
    this.checkRootConfig()
    this.checkRequiredTitle(
      this.document.frontmatter,
      [],
      true,
      "The model root declares no `title`; a model title is required for human-facing display.",
    )
  }

  private checkRootConfig() {
    const pair = mapEntry(this.document.frontmatter, "config")
    if (pair === undefined) return
    const value = pair.value
    const modelPath = ["config"]
    if (!isScalarNode(value) || nodeValue(value).trim() === "") {
      this.add(
        RuleId.invalidConfig,
        "The root `config` value must be a non-empty model-relative scalar path.",
        this.locationForNodeOrMissing(value ?? undefined, modelPath, "config"),
      )
      return
    }
    const path = nodeValue(value)
    const normalized = path.replaceAll("\\", "/")
    let detail: string | undefined
    if (/^(?:\/|[A-Za-z]:\/|\/\/)/.test(normalized)) {
      detail = `path ${JSON.stringify(path)} must be model-relative`
    } else if (
      normalized
        .split("/")
        .reduce(
          (depth, part) =>
            part === ".." ? depth - 1 : part === "." || part === "" ? depth : depth + 1,
          0,
        ) < 0 ||
      normalized === ".." ||
      normalized.startsWith("../")
    ) {
      detail = `path ${JSON.stringify(path)} escapes the repository`
    }
    if (detail !== undefined) {
      this.add(
        RuleId.invalidConfig,
        `The root \`config\` value is invalid: ${detail}.`,
        this.location(value, modelPath, "config"),
      )
    }
  }

  private checkSchemaProperties(
    schema: SchemaNode,
    owner: YamlNode,
    base: ReadonlyArray<PathSegment>,
    context: SchemaContext,
  ) {
    if (!isMapNode(owner)) return
    for (const pair of mapEntries(owner)) {
      const name = nodeValue(pair.key)
      const modelPath = appendPath(base, name)
      const label = pathLabel(modelPath)
      const schemaEntry = schemaProperty(schema, name)
      if (schemaEntry === undefined) {
        if (schema.kind === "model" && name === "config") continue
        if (schema.kind === "area" && hasSchemaProperty(Model, name)) {
          this.add(
            RuleId.misplacedRootKey,
            `The area \`${context.areaName ?? ""}\` declares \`${name}\`; \`${name}\` is only valid on the model root.`,
            this.location(pair.key as ParsedNode, modelPath, label),
          )
        } else {
          this.add(
            RuleId.unknownKey,
            unknownKeyMessage(schema, name, context),
            this.location(pair.key as ParsedNode, modelPath, label),
          )
        }
        continue
      }
      this.checkSchemaProperty(schema, schemaEntry, owner, pair, modelPath, label, context)
    }
  }

  private checkSchemaProperty(
    schema: SchemaNode,
    schemaEntry: Property,
    owner: YamlMap,
    pair: YamlPair,
    modelPath: ReadonlyArray<PathSegment>,
    label: string,
    context: SchemaContext,
  ) {
    const value = pair.value
    if (schema.kind === "requirement" && schemaEntry.name === property.assessment) {
      if (!isScalarNode(value) || isEmptyNode(value)) {
        this.add(
          RuleId.invalidAssessment,
          `The requirement \`${context.requirement ?? ""}\` has an invalid \`assessment\`; a requirement must declare exactly one non-empty scalar assessment.`,
          this.location(value, modelPath, label),
        )
      }
      return
    }
    if (schemaEntry.presence !== "required" && isEmptyNode(value)) {
      this.emptyProperty(owner, pair, modelPath, label)
      return
    }
    if (schemaEntry.shape === "scalar") {
      if (schemaEntry.presence === "required") {
        if (!isScalarNode(value) && !isEmptyNode(value)) {
          this.invalid(
            pair.key as ParsedNode,
            modelPath,
            label,
            `The \`${schemaEntry.name}\` property has the wrong YAML shape; it must be a scalar.`,
          )
        }
      } else if (!isScalarNode(value)) {
        this.invalid(
          pair.key as ParsedNode,
          modelPath,
          label,
          `The \`${schemaEntry.name}\` property has the wrong YAML shape; it must be a scalar.`,
        )
      }
      return
    }
    const valid =
      schemaEntry.shape === "map"
        ? isMapNode(value)
        : schemaEntry.shape === "sequence" && isSequenceNode(value)
    if (!isEmptyNode(value) && !valid) {
      this.invalid(
        pair.key as ParsedNode,
        modelPath,
        label,
        wrongShapeMessage(schemaEntry, context),
      )
      return
    }
    if (valid && schemaEntry.shape === "map" && schemaEntry.valueShape === "scalar") {
      for (const item of mapEntries(value)) {
        if (!isScalarNode(item.value) || isEmptyNode(item.value)) {
          const name = nodeValue(item.key)
          const itemPath = appendPath(modelPath, name)
          this.invalid(
            item.value,
            itemPath,
            pathLabel(itemPath),
            `The rating override \`${name}\` has the wrong YAML shape; override criteria must be non-empty scalars.`,
          )
        }
      }
    }
    if (
      isSequenceNode(value) &&
      schemaEntry.shape === "sequence" &&
      schemaEntry.elementShape === "scalar"
    ) {
      value.items.forEach((item, index) => {
        if (!isScalarNode(item) || isEmptyNode(item)) {
          const itemPath = appendPath(modelPath, index)
          this.invalid(
            item as YamlNode,
            itemPath,
            pathLabel(itemPath),
            `The requirement \`${context.requirement ?? ""}\` has a factor reference with the wrong YAML shape; each factor reference must be a non-empty scalar.`,
          )
        }
      })
    }
  }

  private emptyProperty(
    owner: YamlMap,
    pair: YamlPair,
    modelPath: ReadonlyArray<PathSegment>,
    label: string,
  ) {
    const name = nodeValue(pair.key)
    const location = this.location(pair.key as ParsedNode, modelPath, label)
    this.add(
      RuleId.emptyProperty,
      `The optional property \`${name}\` is empty; empty optional properties should be omitted.`,
      location,
      {
        record: {
          ruleId: RuleId.emptyProperty,
          message: `Removed empty optional property \`${name}\`.`,
          location,
        },
        apply: () => {
          const index = owner.items.indexOf(pair)
          if (index < 0) return false
          owner.items.splice(index, 1)
          return true
        },
      },
    )
  }

  private checkRatingScale() {
    const pair = mapEntry(this.document.frontmatter, property.ratingScale)
    const scale = pair?.value
    if (scale === undefined || isEmptyNode(scale)) {
      this.add(
        RuleId.missingRatingScale,
        "The model root declares no `ratingScale`; a QUALITY.md model requires one rating scale.",
        this.locationForMissing([property.ratingScale], property.ratingScale),
      )
      return
    }
    if (!isSequenceNode(scale)) return
    if (scale.items.length < 2) {
      this.add(
        RuleId.tooFewLevels,
        "The `ratingScale` has fewer than two levels; a QUALITY.md rating scale requires at least two levels.",
        this.location(pair?.key as ParsedNode, [property.ratingScale], property.ratingScale),
      )
    }
    const seen = new Map<string, Location>()
    scale.items.forEach((level, index) => {
      const modelPath = [property.ratingScale, index]
      if (!isMapNode(level)) {
        this.invalid(
          level as YamlNode,
          modelPath,
          `ratingScale[${index}]`,
          "A rating level has the wrong YAML shape; each `ratingScale` item must be a map.",
        )
        return
      }
      this.checkRatingLevel(level, modelPath, index, seen)
    })
  }

  private checkRatingLevel(
    level: YamlMap,
    modelPath: ReadonlyArray<PathSegment>,
    index: number,
    seen: Map<string, Location>,
  ) {
    this.checkSchemaProperties(RatingLevel, level, modelPath, {})
    const levelPair = mapEntry(level, property.level)
    const levelPath = appendPath(modelPath, property.level)
    const levelLabel = `ratingScale[${index}].level`
    if (levelPair === undefined || !isScalarNode(levelPair.value) || isEmptyNode(levelPair.value)) {
      this.add(
        RuleId.missingLevelName,
        "A rating level declares no `level` name; each rating level requires a non-empty `level`.",
        this.locationForNodeOrMissing(levelPair?.value ?? undefined, levelPath, levelLabel),
      )
    } else {
      const name = nodeValue(levelPair.value)
      const location = this.location(levelPair.value, levelPath, levelLabel)
      if (!ModelNamePattern.test(name)) {
        this.add(
          RuleId.invalidRatingLevelId,
          `The rating level ID \`${name}\` is invalid; area names, factor names, requirement names, and rating level IDs must match ${ModelNamePatternText}.`,
          location,
        )
      }
      const prior = seen.get(name)
      if (prior !== undefined) {
        const message = `The rating level \`${name}\` is duplicated; each \`level\` name must be unique within \`ratingScale\`.`
        this.add(RuleId.duplicateLevel, message, location)
        this.add(RuleId.duplicateLevel, message, prior)
      } else {
        seen.set(name, location)
        this.levels.add(name)
      }
    }
    const criterion = mapEntry(level, property.criterion)?.value
    if (!isScalarNode(criterion) || isEmptyNode(criterion)) {
      this.add(
        RuleId.missingCriterion,
        "A rating level declares no `criterion`; each rating level requires a non-empty criterion.",
        this.locationForNodeOrMissing(
          criterion ?? undefined,
          appendPath(modelPath, property.criterion),
          `ratingScale[${index}].criterion`,
        ),
      )
    }
    this.checkRequiredTitle(
      level,
      modelPath,
      false,
      `The rating level at \`ratingScale[${index}]\` declares no \`title\`; each rating level requires a human-facing title.`,
    )
    if (mapEntry(level, property.description) === undefined) {
      this.add(
        RuleId.missingLevelDescription,
        "A rating level declares no `description`; a description is recommended for each level.",
        this.locationForMissing(
          appendPath(modelPath, property.description),
          `ratingScale[${index}].description`,
        ),
      )
    }
  }

  private walkModel() {
    this.root.factors = this.walkFactors(this.root, undefined, this.root.node, [])
    this.root.requirements = this.walkRequirements(this.root, undefined, this.root.node, [])
    this.root.areas = this.walkAreas(this.root, this.root.node, [])
  }

  private walkAreas(
    parent: AreaRef,
    owner: YamlMap,
    base: ReadonlyArray<PathSegment>,
  ): Array<AreaRef> {
    const areas = mapEntry(owner, property.areas)?.value
    if (!isMapNode(areas)) return []
    const output: Array<AreaRef> = []
    for (const pair of mapEntries(areas)) {
      const name = nodeValue(pair.key)
      const modelPath = appendPath(base, property.areas, name)
      if (!isScalarNode(pair.key) || !ModelNamePattern.test(name)) {
        this.add(
          RuleId.invalidAreaName,
          `The area name \`${name}\` is invalid; area names must match ${ModelNamePatternText}.`,
          this.location(pair.key as ParsedNode, modelPath, pathLabel(modelPath)),
        )
      }
      if (name === "root") {
        this.add(
          RuleId.reservedAreaName,
          "The area name `root` is reserved for the root area reference and cannot be used as a child area name.",
          this.location(pair.key as ParsedNode, modelPath, pathLabel(modelPath)),
        )
      }
      if (!isMapNode(pair.value)) {
        this.invalid(
          pair.key as ParsedNode,
          modelPath,
          pathLabel(modelPath),
          `The area \`${name}\` has the wrong YAML shape; each area must be a map.`,
        )
        continue
      }
      const area: AreaRef = {
        name,
        node: pair.value,
        path: modelPath,
        parent,
        factors: [],
        requirements: [],
        areas: [],
      }
      this.checkAreaShape(area)
      area.factors = this.walkFactors(area, undefined, area.node, modelPath)
      area.requirements = this.walkRequirements(area, undefined, area.node, modelPath)
      area.areas = this.walkAreas(area, area.node, modelPath)
      output.push(area)
    }
    return output
  }

  private walkFactors(
    area: AreaRef,
    parent: FactorRef | undefined,
    owner: YamlMap,
    base: ReadonlyArray<PathSegment>,
  ): Array<FactorRef> {
    const factors = mapEntry(owner, property.factors)?.value
    if (!isMapNode(factors)) return []
    const output: Array<FactorRef> = []
    for (const pair of mapEntries(factors)) {
      const name = nodeValue(pair.key)
      const modelPath = appendPath(base, property.factors, name)
      if (!isScalarNode(pair.key) || !ModelNamePattern.test(name)) {
        this.add(
          RuleId.invalidFactorName,
          `The factor name \`${name}\` is invalid; factor names must match ${ModelNamePatternText}.`,
          this.location(pair.key as ParsedNode, modelPath, pathLabel(modelPath)),
        )
      }
      if (!isMapNode(pair.value)) {
        this.invalid(
          pair.key as ParsedNode,
          modelPath,
          pathLabel(modelPath),
          `The factor \`${name}\` has the wrong YAML shape; each factor must be a map.`,
        )
        continue
      }
      const factor: FactorRef = {
        name,
        node: pair.value,
        path: modelPath,
        area,
        ...(parent === undefined ? {} : { parent }),
        factors: [],
        requirements: [],
      }
      this.checkFactorShape(factor)
      factor.factors = this.walkFactors(area, factor, factor.node, modelPath)
      factor.requirements = this.walkRequirements(area, factor, factor.node, modelPath)
      output.push(factor)
    }
    return output
  }

  private walkRequirements(
    area: AreaRef,
    factor: FactorRef | undefined,
    owner: YamlMap,
    base: ReadonlyArray<PathSegment>,
  ): Array<RequirementRef> {
    const requirements = mapEntry(owner, property.requirements)?.value
    if (!isMapNode(requirements)) return []
    const output: Array<RequirementRef> = []
    for (const pair of mapEntries(requirements)) {
      const name = nodeValue(pair.key)
      const modelPath = appendPath(base, property.requirements, name)
      if (!isScalarNode(pair.key) || !ModelNamePattern.test(name)) {
        this.add(
          RuleId.invalidRequirementName,
          `The requirement name \`${name}\` is invalid; requirement names must match ${ModelNamePatternText}.`,
          this.location(pair.key as ParsedNode, modelPath, pathLabel(modelPath)),
        )
      }
      if (!isMapNode(pair.value)) {
        this.invalid(
          pair.key as ParsedNode,
          modelPath,
          pathLabel(modelPath),
          `The requirement \`${name}\` has the wrong YAML shape; each requirement must be a map.`,
        )
        continue
      }
      const requirement: RequirementRef = {
        name,
        node: pair.value,
        path: modelPath,
        area,
        ...(factor === undefined ? {} : { factor }),
      }
      this.checkRequirementShape(requirement)
      output.push(requirement)
    }
    return output
  }

  private checkRequiredTitle(
    owner: YamlNode,
    base: ReadonlyArray<PathSegment>,
    root: boolean,
    message: string,
  ) {
    const value = mapEntry(owner, property.title)?.value
    const modelPath = appendPath(base, property.title)
    if (!isScalarNode(value) || isEmptyNode(value)) {
      this.add(
        RuleId.missingTitle,
        message,
        this.locationForNodeOrMissing(
          value ?? undefined,
          modelPath,
          root ? "title" : pathLabel(modelPath),
        ),
      )
    }
  }

  private checkAreaShape(area: AreaRef) {
    this.checkSchemaProperties(Area, area.node, area.path, { areaName: area.name })
    this.checkRequiredTitle(
      area.node,
      area.path,
      false,
      `The area \`${area.name}\` declares no \`title\`; each area requires a human-facing title.`,
    )
  }

  private checkFactorShape(factor: FactorRef) {
    this.checkSchemaProperties(Factor, factor.node, factor.path, { factorName: factor.name })
    this.checkRequiredTitle(
      factor.node,
      factor.path,
      false,
      `The factor \`${factor.name}\` declares no \`title\`; each factor requires a human-facing title.`,
    )
    if (mapEntry(factor.node, property.description) === undefined) {
      const modelPath = appendPath(factor.path, property.description)
      this.add(
        RuleId.missingFactorDescription,
        `The factor \`${factor.name}\` declares no \`description\`; a description is recommended for each factor.`,
        this.locationForMissing(modelPath, pathLabel(modelPath)),
      )
    }
  }

  private checkRequirementShape(requirement: RequirementRef) {
    this.checkSchemaProperties(Requirement, requirement.node, requirement.path, {
      requirement: requirement.name,
    })
    this.checkRequiredTitle(
      requirement.node,
      requirement.path,
      false,
      `The requirement \`${requirement.name}\` declares no \`title\`; each requirement requires a human-facing title.`,
    )
    if (mapEntry(requirement.node, property.assessment) === undefined) {
      const modelPath = appendPath(requirement.path, property.assessment)
      this.add(
        RuleId.invalidAssessment,
        `The requirement \`${requirement.name}\` has no \`assessment\`; a requirement must declare exactly one non-empty scalar assessment.`,
        this.locationForMissing(modelPath, pathLabel(modelPath)),
      )
    }
  }

  private checkEmptyModel() {
    if (
      this.root.factors.length === 0 &&
      this.root.requirements.length === 0 &&
      this.root.areas.length === 0
    ) {
      this.add(
        RuleId.emptyModel,
        "The model root supplies no entries under `factors`, `requirements`, or `areas`; a QUALITY.md model requires model content.",
        this.location(this.root.node, [], "frontmatter"),
      )
    }
  }

  private checkAreas(area: AreaRef): boolean {
    let hasRequirements = area.requirements.length > 0 || area.factors.some(factorHasRequirements)
    for (const child of area.areas) if (this.checkAreas(child)) hasRequirements = true
    if (area.parent !== undefined && !hasRequirements) {
      this.add(
        RuleId.emptyArea,
        `The area \`${area.name}\` reaches no requirements in its subtree; each area should lead to at least one requirement.`,
        this.location(area.node, area.path, pathLabel(area.path)),
      )
    }
    return hasRequirements
  }

  private checkFactors(area: AreaRef) {
    const byName = new Map<string, Array<FactorRef>>()
    const collect = (factors: ReadonlyArray<FactorRef>) => {
      for (const factor of factors) {
        const entries = byName.get(factor.name) ?? []
        entries.push(factor)
        byName.set(factor.name, entries)
        collect(factor.factors)
      }
    }
    collect(area.factors)
    for (const [name, factors] of byName) {
      if (factors.length < 2) continue
      for (const factor of factors) {
        this.add(
          RuleId.duplicateFactorName,
          `The factor name \`${name}\` appears more than once in area \`${area.parent === undefined ? "root" : area.name}\`; factor references use names, so repeated names in one area can be ambiguous.`,
          this.location(factor.node, factor.path, pathLabel(factor.path)),
        )
      }
    }
    for (const factor of area.factors) this.checkFactor(factor)
    for (const child of area.areas) this.checkFactors(child)
  }

  private checkFactor(factor: FactorRef): boolean {
    let hasRequirements = factor.requirements.length > 0
    for (const child of factor.factors) if (this.checkFactor(child)) hasRequirements = true
    for (const requirement of allRequirements(factor.area)) {
      if (this.referencedFactors(requirement).includes(factor)) hasRequirements = true
    }
    if (!hasRequirements) {
      this.add(
        RuleId.emptyFactor,
        `The factor \`${factor.name}\` leads to no requirements; each factor should be tied to at least one requirement.`,
        this.location(factor.node, factor.path, pathLabel(factor.path)),
      )
    }
    return hasRequirements
  }

  private checkRequirements(area: AreaRef) {
    const seen = new Map<string, Location>()
    for (const requirement of localRequirements(area)) {
      const location = this.locationForMissing(requirement.path, pathLabel(requirement.path))
      const prior = seen.get(requirement.name)
      if (prior === undefined) seen.set(requirement.name, location)
      else {
        const message = `The requirement \`${requirement.name}\` is duplicated; requirement names must be unique within their declaring area.`
        this.add(RuleId.duplicateRequirement, message, location)
        this.add(RuleId.duplicateRequirement, message, prior)
      }
    }
    for (const requirement of area.requirements) this.checkRequirementReferences(requirement)
    for (const factor of area.factors) this.checkFactorRequirementReferences(factor)
    for (const child of area.areas) this.checkRequirements(child)
  }

  private checkFactorRequirementReferences(factor: FactorRef) {
    for (const requirement of factor.requirements) this.checkRequirementReferences(requirement)
    for (const child of factor.factors) this.checkFactorRequirementReferences(child)
  }

  private checkRequirementReferences(requirement: RequirementRef) {
    const factorItems = mapEntry(requirement.node, property.factors)?.value
    const referencesFactor =
      isSequenceNode(factorItems) &&
      factorItems.items.some((item) => isScalarNode(item) && !isEmptyNode(item))
    if (requirement.factor === undefined && !referencesFactor) {
      const modelPath = appendPath(requirement.path, property.factors)
      this.add(
        RuleId.missingFactorReference,
        `The requirement \`${requirement.name}\` references no quality factor; place it under a factor or add one or more factor references under \`factors\`.`,
        this.locationForMissing(modelPath, pathLabel(modelPath)),
      )
    }
    const ratings = mapEntry(requirement.node, property.ratings)?.value
    if (isMapNode(ratings)) {
      for (const pair of mapEntries(ratings)) {
        const name = nodeValue(pair.key)
        if (this.levels.has(name)) continue
        const modelPath = appendPath(requirement.path, property.ratings, name)
        this.add(
          RuleId.unknownRatingKey,
          `The requirement \`${requirement.name}\` has a \`ratings\` override for unknown level \`${name}\`; override keys must name a rating-scale level.`,
          this.location(pair.key as ParsedNode, modelPath, pathLabel(modelPath)),
        )
      }
    }
    if (isSequenceNode(factorItems)) {
      factorItems.items.forEach((item, index) => {
        if (!isScalarNode(item) || isEmptyNode(item)) return
        const name = nodeValue(item)
        if (findFactor(requirement.area.factors, name) !== undefined) return
        const modelPath = appendPath(requirement.path, property.factors, index)
        this.add(
          RuleId.unknownFactor,
          `The requirement \`${requirement.name}\` references unknown factor \`${name}\`; factor references must resolve within the declaring area.`,
          this.location(item, modelPath, pathLabel(modelPath)),
        )
      })
    }
  }

  private referencedFactors(requirement: RequirementRef) {
    const items = mapEntry(requirement.node, property.factors)?.value
    if (!isSequenceNode(items)) return []
    return items.items.flatMap((item) => {
      if (!isScalarNode(item) || isEmptyNode(item)) return []
      const factor = findFactor(requirement.area.factors, nodeValue(item))
      return factor === undefined ? [] : [factor]
    })
  }
}

const unknownKeyMessage = (schema: SchemaNode, key: string, context: SchemaContext) => {
  const advice =
    "the format permits extension properties, but check it is not a typo of a model property."
  if (schema.kind === "model")
    return `The frontmatter declares \`${key}\`, which is not a model root property; ${advice}`
  if (schema.kind === "area")
    return `The area \`${context.areaName ?? ""}\` declares \`${key}\`, which is not an area property; ${advice}`
  if (schema.kind === "factor")
    return `The factor \`${context.factorName ?? ""}\` declares \`${key}\`, which is not a factor property; ${advice}`
  if (schema.kind === "requirement")
    return `The requirement \`${context.requirement ?? ""}\` declares \`${key}\`, which is not a requirement property; ${advice}`
  if (schema.kind === "ratingLevel") {
    const names = schema.properties.map((entry) => `\`${entry.name}\``)
    const formatted = `${names.slice(0, -1).join(", ")}, and ${names.at(-1)}`
    return `A rating level declares \`${key}\`, which is not a rating level property (${formatted}); ${advice}`
  }
  return `The key \`${key}\` is not a model property; ${advice}`
}

const wrongShapeMessage = (entry: Property, context: SchemaContext) => {
  if (entry.name === property.ratingScale)
    return "The `ratingScale` property has the wrong YAML shape; it must be a list of rating levels."
  if (entry.name === property.factors && context.requirement !== undefined)
    return `The requirement \`${context.requirement}\` has the wrong \`factors\` shape; factor references must be a sequence.`
  if (entry.name === property.ratings)
    return `The requirement \`${context.requirement ?? ""}\` has the wrong \`ratings\` shape; ratings overrides must be a map.`
  return `The \`${entry.name}\` property has the wrong YAML shape; it must be a ${entry.shape === "sequence" ? "list" : entry.shape}.`
}

const factorHasRequirements = (factor: FactorRef): boolean =>
  factor.requirements.length > 0 || factor.factors.some(factorHasRequirements)

const findFactor = (factors: ReadonlyArray<FactorRef>, name: string): FactorRef | undefined => {
  for (const factor of factors) {
    if (factor.name === name) return factor
    const child = findFactor(factor.factors, name)
    if (child !== undefined) return child
  }
  return undefined
}

const allRequirements = (area: AreaRef): Array<RequirementRef> => [
  ...localRequirements(area),
  ...area.areas.flatMap(allRequirements),
]

const localRequirements = (area: AreaRef): Array<RequirementRef> => {
  const factorRequirements = (factor: FactorRef): Array<RequirementRef> => [
    ...factor.requirements,
    ...factor.factors.flatMap(factorRequirements),
  ]
  return [...area.requirements, ...area.factors.flatMap(factorRequirements)]
}

export const lintDocument = (document: QualityDocument): LintRun => new RunState(document).run()

export const invalidDocumentResult = (path: string): LintResult => ({
  schemaVersion: 1,
  path,
  valid: false,
  summary: { errors: 1, warnings: 0, info: 0, fixable: 0, fixed: 0 },
  findings: [
    {
      ruleId: RuleId.invalidFrontmatter,
      severity: "error",
      message:
        "The file does not begin with valid QUALITY.md frontmatter; a QUALITY.md file requires a YAML frontmatter block matching the model shape.",
      location: { path, modelPath: [], label: "frontmatter" },
      fixable: false,
    },
  ],
  repairs: [],
  nextActions: nextActions(path, 1, 0),
})
