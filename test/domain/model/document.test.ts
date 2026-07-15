import { describe, expect, it } from "vitest"

import {
  mapEntries,
  nodeValue,
  parseQualityDocument,
  QualityDocumentParseError,
  renderQualityDocument,
} from "../../../src/domain/model/document.ts"

describe("QUALITY.md document", () => {
  it.each([
    ["missing frontmatter", "title: Missing", /must begin with a YAML frontmatter block/],
    ["unterminated frontmatter", "---\ntitle: Missing", /unexpected end of file/],
    ["empty frontmatter", "---\n---\n", /frontmatter is empty/],
  ])("rejects %s", (_name, raw, detail) => {
    expect(() => parseQualityDocument("QUALITY.md", raw)).toThrow(QualityDocumentParseError)
    expect(() => parseQualityDocument("QUALITY.md", raw)).toThrow(detail)
  })

  it("preserves the Markdown body while rendering YAML edits", () => {
    const raw = "---\ntitle: Example\nratingScale: []\n---\n# Body\n\nKeep me.\n"
    const document = parseQualityDocument("QUALITY.md", raw)
    expect(document.body).toBe("# Body\n\nKeep me.\n")
    expect(renderQualityDocument(document)).toBe(raw)
  })

  it("exposes ordered mapping entries", () => {
    const document = parseQualityDocument("QUALITY.md", "---\nfirst: one\nsecond: two\n---\n")
    expect(mapEntries(document.frontmatter).map((entry) => nodeValue(entry.key))).toEqual([
      "first",
      "second",
    ])
  })
})
