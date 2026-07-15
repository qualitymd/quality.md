export type Severity = "error" | "warning" | "info"
export type PathSegment = string | number

export const RuleId = {
  invalidFrontmatter: "invalid-frontmatter",
  missingTitle: "missing-title",
  missingRatingScale: "missing-rating-scale",
  tooFewLevels: "too-few-levels",
  missingLevelName: "missing-level-name",
  invalidRatingLevelId: "invalid-rating-level-id",
  duplicateLevel: "duplicate-level",
  invalidConfig: "invalid-config",
  missingCriterion: "missing-criterion",
  missingLevelDescription: "missing-level-description",
  emptyModel: "empty-model",
  misplacedRootKey: "misplaced-root-key",
  invalidAreaName: "invalid-area-name",
  reservedAreaName: "reserved-area-name",
  emptyArea: "empty-area",
  invalidFactorName: "invalid-factor-name",
  duplicateFactorName: "duplicate-factor-name",
  invalidRequirementName: "invalid-requirement-name",
  duplicateRequirement: "duplicate-requirement",
  emptyFactor: "empty-factor",
  missingFactorDescription: "missing-factor-description",
  invalidAssessment: "invalid-assessment",
  unknownRatingKey: "unknown-rating-key",
  unknownFactor: "unknown-factor",
  missingFactorReference: "missing-factor-reference",
  emptyProperty: "empty-property",
  unknownKey: "unknown-key",
} as const

export type RuleId = (typeof RuleId)[keyof typeof RuleId]

export interface Rule {
  readonly id: RuleId
  readonly severity: Severity
  readonly fixable: boolean
  readonly description: string
}

const rule = (id: RuleId, severity: Severity, fixable: boolean, description: string): Rule => ({
  id,
  severity,
  fixable,
  description,
})

export const Rules: ReadonlyArray<Rule> = [
  rule(RuleId.invalidFrontmatter, "error", false, "The frontmatter is missing or invalid."),
  rule(RuleId.missingTitle, "error", false, "A model element declares no title."),
  rule(RuleId.missingRatingScale, "error", false, "The model root declares no rating scale."),
  rule(RuleId.tooFewLevels, "error", false, "The rating scale has fewer than two levels."),
  rule(RuleId.missingLevelName, "error", false, "A rating level declares no level name."),
  rule(RuleId.invalidRatingLevelId, "error", false, "A rating level ID is invalid."),
  rule(RuleId.duplicateLevel, "error", false, "A rating level name is duplicated."),
  rule(RuleId.invalidConfig, "error", false, "The root config path is invalid."),
  rule(RuleId.missingCriterion, "error", false, "A rating level declares no criterion."),
  rule(RuleId.missingLevelDescription, "warning", false, "A rating level has no description."),
  rule(RuleId.emptyModel, "error", false, "The model root supplies no model content."),
  rule(RuleId.misplacedRootKey, "error", false, "A root-only key appears on an area."),
  rule(RuleId.invalidAreaName, "error", false, "An area name is invalid."),
  rule(RuleId.reservedAreaName, "error", false, "An area name is reserved."),
  rule(RuleId.emptyArea, "warning", false, "An area reaches no requirements."),
  rule(RuleId.invalidFactorName, "error", false, "A factor name is invalid."),
  rule(RuleId.duplicateFactorName, "warning", false, "A factor name is repeated."),
  rule(RuleId.invalidRequirementName, "error", false, "A requirement name is invalid."),
  rule(RuleId.duplicateRequirement, "error", false, "A requirement name is duplicated."),
  rule(RuleId.emptyFactor, "warning", false, "A factor leads to no requirements."),
  rule(RuleId.missingFactorDescription, "warning", false, "A factor has no description."),
  rule(RuleId.invalidAssessment, "error", false, "A requirement has no scalar assessment."),
  rule(RuleId.unknownRatingKey, "error", false, "A rating override names an unknown level."),
  rule(RuleId.unknownFactor, "error", false, "A requirement references an unknown factor."),
  rule(RuleId.missingFactorReference, "error", false, "A direct requirement has no factor."),
  rule(RuleId.emptyProperty, "warning", true, "An optional property is empty."),
  rule(RuleId.unknownKey, "warning", false, "A key names no model property."),
]

export const RulesById = new Map(Rules.map((entry) => [entry.id, entry]))

export interface Action {
  readonly id: string
  readonly label: string
  readonly command: string
}

export interface Location {
  readonly path: string
  readonly modelPath: ReadonlyArray<PathSegment>
  readonly label: string
  readonly line?: number
  readonly column?: number
}

export interface Finding {
  readonly ruleId: RuleId
  readonly severity: Severity
  readonly message: string
  readonly location: Location
  readonly fixable: boolean
}

export interface RepairRecord {
  readonly ruleId: RuleId
  readonly message: string
  readonly location: Location
}

export interface LintResult {
  readonly schemaVersion: 1
  readonly path: string
  readonly valid: boolean
  readonly summary: {
    readonly errors: number
    readonly warnings: number
    readonly info: number
    readonly fixable: number
    readonly fixed: number
  }
  readonly findings: ReadonlyArray<Finding>
  readonly repairs: ReadonlyArray<RepairRecord>
  readonly nextActions: ReadonlyArray<Action>
}
