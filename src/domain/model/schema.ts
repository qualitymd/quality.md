export const ModelNamePattern = /^[A-Za-z0-9](?:[A-Za-z0-9_-]*[A-Za-z0-9])?$/
export const ModelNamePatternText = "^[A-Za-z0-9](?:[A-Za-z0-9_-]*[A-Za-z0-9])?$"

export type NodeKind = "model" | "area" | "factor" | "requirement" | "ratingLevel"
export type Shape = "scalar" | "map" | "sequence"
export type Presence = "required" | "recommended" | "optional"

export interface Property {
  readonly name: string
  readonly shape: Shape
  readonly presence: Presence
  readonly elementKind?: NodeKind
  readonly elementShape?: Shape
  readonly valueShape?: Shape
  readonly minItems?: number
  readonly keyPattern?: string
  readonly pattern?: string
}

export interface SchemaNode {
  readonly kind: NodeKind
  readonly properties: ReadonlyArray<Property>
  readonly requiredAny?: ReadonlyArray<{
    readonly name: string
    readonly properties: ReadonlyArray<string>
  }>
}

export const property = {
  title: "title",
  ratingScale: "ratingScale",
  factors: "factors",
  requirements: "requirements",
  areas: "areas",
  source: "source",
  level: "level",
  description: "description",
  criterion: "criterion",
  assessment: "assessment",
  ratings: "ratings",
} as const

export const Model: SchemaNode = {
  kind: "model",
  properties: [
    { name: property.title, shape: "scalar", presence: "required" },
    { name: property.description, shape: "scalar", presence: "optional" },
    {
      name: property.ratingScale,
      shape: "sequence",
      presence: "required",
      elementKind: "ratingLevel",
      minItems: 2,
    },
    {
      name: property.factors,
      shape: "map",
      presence: "optional",
      elementKind: "factor",
      keyPattern: ModelNamePatternText,
    },
    {
      name: property.requirements,
      shape: "map",
      presence: "optional",
      elementKind: "requirement",
      keyPattern: ModelNamePatternText,
    },
    {
      name: property.areas,
      shape: "map",
      presence: "optional",
      elementKind: "area",
      keyPattern: ModelNamePatternText,
    },
    { name: property.source, shape: "scalar", presence: "optional" },
  ],
  requiredAny: [
    {
      name: "model-content",
      properties: [property.factors, property.requirements, property.areas],
    },
  ],
}

export const Area: SchemaNode = {
  kind: "area",
  properties: [
    { name: property.title, shape: "scalar", presence: "required" },
    { name: property.description, shape: "scalar", presence: "optional" },
    {
      name: property.factors,
      shape: "map",
      presence: "optional",
      elementKind: "factor",
      keyPattern: ModelNamePatternText,
    },
    {
      name: property.requirements,
      shape: "map",
      presence: "optional",
      elementKind: "requirement",
      keyPattern: ModelNamePatternText,
    },
    {
      name: property.areas,
      shape: "map",
      presence: "optional",
      elementKind: "area",
      keyPattern: ModelNamePatternText,
    },
    { name: property.source, shape: "scalar", presence: "optional" },
  ],
}

export const Factor: SchemaNode = {
  kind: "factor",
  properties: [
    { name: property.title, shape: "scalar", presence: "required" },
    { name: property.description, shape: "scalar", presence: "recommended" },
    {
      name: property.factors,
      shape: "map",
      presence: "optional",
      elementKind: "factor",
      keyPattern: ModelNamePatternText,
    },
    {
      name: property.requirements,
      shape: "map",
      presence: "optional",
      elementKind: "requirement",
      keyPattern: ModelNamePatternText,
    },
  ],
}

export const Requirement: SchemaNode = {
  kind: "requirement",
  properties: [
    { name: property.title, shape: "scalar", presence: "required" },
    { name: property.description, shape: "scalar", presence: "optional" },
    { name: property.assessment, shape: "scalar", presence: "required" },
    {
      name: property.factors,
      shape: "sequence",
      presence: "optional",
      elementShape: "scalar",
    },
    { name: property.ratings, shape: "map", presence: "optional", valueShape: "scalar" },
  ],
}

export const RatingLevel: SchemaNode = {
  kind: "ratingLevel",
  properties: [
    {
      name: property.level,
      shape: "scalar",
      presence: "required",
      pattern: ModelNamePatternText,
    },
    { name: property.title, shape: "scalar", presence: "required" },
    { name: property.description, shape: "scalar", presence: "recommended" },
    { name: property.criterion, shape: "scalar", presence: "required" },
  ],
}

export const schemaProperty = (node: SchemaNode, name: string) =>
  node.properties.find((candidate) => candidate.name === name)

export const hasSchemaProperty = (node: SchemaNode, name: string) =>
  schemaProperty(node, name) !== undefined
