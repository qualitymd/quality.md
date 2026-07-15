import {
  isMap,
  isScalar,
  isSeq,
  LineCounter,
  parseDocument,
  type Document,
  type Node,
  type Pair,
  type ParsedNode,
  type Scalar,
  type YAMLMap,
  type YAMLSeq,
} from "yaml"

export type YamlNode = ParsedNode | null | undefined
export type YamlMap = YAMLMap.Parsed
export type YamlSeq = YAMLSeq.Parsed
export type YamlPair = Pair<ParsedNode, ParsedNode | null>

export interface QualityDocument {
  readonly path: string
  readonly yaml: Document.Parsed<ParsedNode>
  readonly frontmatter: YamlNode
  readonly body: string
  readonly lineCounter: LineCounter
}

export class QualityDocumentParseError extends Error {
  readonly path: string

  constructor(path: string, detail: string) {
    super(`${path}: invalid frontmatter: ${detail}`)
    this.name = "QualityDocumentParseError"
    this.path = path
  }
}

const splitFrontmatter = (raw: string): { readonly frontmatter: string; readonly body: string } => {
  if (!raw.startsWith("---")) {
    throw new Error("file must begin with a YAML frontmatter block")
  }
  const firstLineEnd = raw.indexOf("\n")
  if (firstLineEnd < 0) throw new Error('unterminated frontmatter: missing closing "---"')
  if (raw.slice(0, firstLineEnd).trim() !== "---") {
    throw new Error('opening frontmatter fence must be "---"')
  }
  let position = firstLineEnd + 1
  while (position <= raw.length) {
    const next = raw.indexOf("\n", position)
    const lineEnd = next < 0 ? raw.length : next
    if (raw.slice(position, lineEnd).replace(/\r$/, "").trim() === "---") {
      return {
        frontmatter: raw.slice(firstLineEnd + 1, position),
        body: raw.slice(next < 0 ? lineEnd : lineEnd + 1),
      }
    }
    if (next < 0) break
    position = next + 1
  }
  throw new Error("unexpected end of file")
}

export const parseQualityDocument = (path: string, raw: string): QualityDocument => {
  try {
    const { frontmatter, body } = splitFrontmatter(raw)
    const lineCounter = new LineCounter()
    const yaml = parseDocument(frontmatter, {
      keepSourceTokens: true,
      lineCounter,
      prettyErrors: false,
      strict: false,
      uniqueKeys: false,
    })
    if (yaml.errors.length > 0) throw yaml.errors[0]
    if (yaml.contents === null) throw new Error("frontmatter is empty")
    return {
      path,
      yaml: yaml as Document.Parsed<ParsedNode>,
      frontmatter: yaml.contents,
      body,
      lineCounter,
    }
  } catch (cause) {
    if (cause instanceof QualityDocumentParseError) throw cause
    throw new QualityDocumentParseError(
      path,
      cause instanceof Error ? cause.message : String(cause),
    )
  }
}

export const renderQualityDocument = (document: QualityDocument): string => {
  const yaml = document.yaml.toString({ indent: 2, lineWidth: 0 })
  return `---\n${yaml}---\n${document.body}`
}

export const mapEntries = (node: YamlNode): ReadonlyArray<YamlPair> =>
  isMap(node) ? (node.items as Array<YamlPair>) : []

export const mapEntry = (node: YamlNode, name: string): YamlPair | undefined =>
  mapEntries(node).find((pair) => nodeValue(pair.key) === name)

export const nodeValue = (node: Node | null | undefined): string => {
  if (isScalar(node)) return node.value === null ? "" : String(node.value)
  if (node === null || node === undefined) return ""
  return String(node.toJSON())
}

export const isEmptyNode = (node: YamlNode): boolean => {
  if (node === null) return true
  if (isScalar(node)) return node.value === null || String(node.value).trim() === ""
  if (isMap(node) || isSeq(node)) return node.items.length === 0
  return false
}

export const isScalarNode = (node: YamlNode): node is Scalar.Parsed => isScalar(node)
export const isMapNode = (node: YamlNode): node is YamlMap => isMap(node)
export const isSequenceNode = (node: YamlNode): node is YamlSeq => isSeq(node)
