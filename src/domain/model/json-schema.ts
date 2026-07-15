import {
  Area,
  Factor,
  Model,
  RatingLevel,
  Requirement,
  type Property,
  type SchemaNode,
} from "./schema.ts"

type JsonObject = Record<string, unknown>

const scalarValueSchema = (): JsonObject => ({
  anyOf: [{ type: "string", minLength: 1 }, { type: "number" }, { type: "boolean" }],
})

const elementSchema = (property: Property): JsonObject => {
  if (property.elementKind !== undefined) return { $ref: `#/$defs/${property.elementKind}` }
  if (property.elementShape === "scalar" || property.valueShape === "scalar")
    return scalarValueSchema()
  throw new Error(`schema property ${JSON.stringify(property.name)} has no element type`)
}

const propertySchema = (property: Property): JsonObject => {
  if (property.shape === "scalar")
    return property.pattern === undefined
      ? scalarValueSchema()
      : { type: "string", pattern: property.pattern }
  if (property.shape === "sequence")
    return {
      type: "array",
      items: elementSchema(property),
      ...(property.minItems === undefined ? {} : { minItems: property.minItems }),
    }
  const schema: JsonObject = {
    type: "object",
    additionalProperties: elementSchema(property),
  }
  if (property.keyPattern !== undefined) {
    schema.propertyNames =
      property.elementKind === "area"
        ? {
            allOf: [{ pattern: property.keyPattern }, { not: { enum: ["root"] } }],
          }
        : { pattern: property.keyPattern }
  }
  return schema
}

const nodeSchema = (node: SchemaNode): JsonObject => {
  const required = node.properties
    .filter((property) => property.presence === "required")
    .map((property) => property.name)
  const schema: JsonObject = {
    type: "object",
    properties: Object.fromEntries(
      node.properties.map((property) => [property.name, propertySchema(property)]),
    ),
    ...(required.length === 0 ? {} : { required }),
  }
  if (node.requiredAny?.length === 1) {
    schema.anyOf = node.requiredAny[0]!.properties.map((name) => ({ required: [name] }))
  } else if (node.requiredAny !== undefined && node.requiredAny.length > 1) {
    schema.allOf = node.requiredAny.map((group) => ({
      anyOf: group.properties.map((name) => ({ required: [name] })),
    }))
  }
  return schema
}

const sortKeys = (value: unknown): unknown => {
  if (Array.isArray(value)) return value.map(sortKeys)
  if (value === null || typeof value !== "object") return value
  return Object.fromEntries(
    Object.entries(value as JsonObject)
      .sort(([left], [right]) => left.localeCompare(right))
      .map(([key, child]) => [key, sortKeys(child)]),
  )
}

export const generateQualitySchema = () => {
  const root = {
    ...nodeSchema(Model),
    $schema: "https://json-schema.org/draft/2020-12/schema",
    $id: "https://getquality.md/quality.schema.json",
    title: "QUALITY.md frontmatter",
    description:
      "Structural JSON Schema for QUALITY.md frontmatter, derived from the qualitymd linter's schema.",
    $comment:
      "Non-normative and subordinate to SPECIFICATION.md (https://getquality.md). Structural-only: passing this schema does not imply full conformance. Semantic rules (factor-reference resolution, rating-override keys, the placement-dependent factor-connection rule, and rating-level uniqueness) are enforced by `qualitymd lint`, not here. Generated from src/domain/model/schema.ts and guarded against drift by a consistency test; do not edit by hand — run `mise run schema`.",
    $defs: Object.fromEntries(
      [Area, Factor, Requirement, RatingLevel].map((node) => [node.kind, nodeSchema(node)]),
    ),
  }
  return `${JSON.stringify(sortKeys(root), null, 2)}\n`
}
