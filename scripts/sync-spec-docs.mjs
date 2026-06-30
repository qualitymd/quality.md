#!/usr/bin/env node
// Generates the Mintlify specification page from SPECIFICATION.md.
//
// SPECIFICATION.md is the single source of truth. This script adapts it into
// mintlify/specification.mdx: it adds Mintlify frontmatter, drops the leading
// H1 (Mintlify renders the title from frontmatter), and rewrites repo-relative
// links to resolvable targets. Run via `mise run sync-spec-docs`; the pre-commit
// hook and `mise run check` keep the generated page in sync.
//
// Do not edit mintlify/specification.mdx by hand. Edit SPECIFICATION.md instead.

import { mkdirSync, readFileSync, writeFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const root = join(dirname(fileURLToPath(import.meta.url)), "..");
const SOURCE = join(root, "SPECIFICATION.md");
const OUTPUT = join(root, "mintlify", "specification.mdx");

const REPO_BLOB_BASE = "https://github.com/qualitymd/quality.md/blob/main";

// Repo-relative paths with a better canonical home than a blob URL. Anchors are
// dropped because the canonical target is the document itself.
const LINK_OVERRIDES = {
  "docs/reference/rfc2119.md": "https://www.rfc-editor.org/rfc/rfc2119",
  "docs/reference/rfc8174.md": "https://www.rfc-editor.org/rfc/rfc8174",
};

const FRONTMATTER = `---
title: Specification
description: The QUALITY.md format specification — document structure, model vocabulary, frontmatter schema, and evaluation semantics.
---

{/* Generated from SPECIFICATION.md by scripts/sync-spec-docs.mjs. Do not edit directly. */}
`;

function rewriteLink(dest) {
  if (dest.startsWith("#") || /^[a-z][a-z0-9+.-]*:/i.test(dest)) {
    // In-page anchor or already-absolute (http:, https:, mailto:, …).
    return dest;
  }

  const hashIndex = dest.indexOf("#");
  const path = (hashIndex === -1 ? dest : dest.slice(0, hashIndex)).replace(
    /^\.\//,
    "",
  );
  const anchor = hashIndex === -1 ? "" : dest.slice(hashIndex);

  if (Object.prototype.hasOwnProperty.call(LINK_OVERRIDES, path)) {
    return LINK_OVERRIDES[path];
  }

  return `${REPO_BLOB_BASE}/${path}${anchor}`;
}

function rewriteLinks(markdown) {
  // Inline markdown links: [text](dest). Destinations here never carry titles.
  return markdown.replace(/(\]\()([^)\s]+)(\))/g, (_, open, dest, close) => {
    return `${open}${rewriteLink(dest)}${close}`;
  });
}

// Rewrite links only in prose. Code fences and inline code are literal in MDX
// and can contain `](` sequences (for example the name grammar regex) that are
// not links, so they must be left untouched.
function rewriteLinksOutsideCode(markdown) {
  const codeSpan = /(```[\s\S]*?```|`[^`\n]*`)/g;
  let result = "";
  let lastIndex = 0;
  let match;
  while ((match = codeSpan.exec(markdown)) !== null) {
    result += rewriteLinks(markdown.slice(lastIndex, match.index));
    result += match[0];
    lastIndex = match.index + match[0].length;
  }
  result += rewriteLinks(markdown.slice(lastIndex));
  return result;
}

function stripLeadingH1(markdown) {
  if (!markdown.startsWith("# ")) {
    throw new Error("SPECIFICATION.md must start with an H1 heading.");
  }
  const newlineIndex = markdown.indexOf("\n");
  return markdown.slice(newlineIndex + 1).replace(/^\n+/, "");
}

// Fail loudly if the prose contains characters MDX would parse as JSX, rather
// than emit a page that breaks the Mintlify build. Code fences and inline code
// are literal in MDX, so they are removed before the scan.
function assertMdxSafe(markdown) {
  const withoutCode = markdown
    .replace(/```[\s\S]*?```/g, "")
    .replace(/`[^`\n]*`/g, "");

  const lines = withoutCode.split("\n");
  const problems = [];
  lines.forEach((line, index) => {
    // `<` that looks like a tag, or `{` that opens a JSX expression.
    if (/<[A-Za-z/!]/.test(line) || line.includes("{")) {
      problems.push(`  line ${index + 1}: ${line.trim()}`);
    }
  });

  if (problems.length > 0) {
    throw new Error(
      "SPECIFICATION.md contains characters MDX parses as JSX outside code spans.\n" +
        "Wrap them in backticks or a code fence:\n" +
        problems.join("\n"),
    );
  }
}

const source = readFileSync(SOURCE, "utf8");
const body = rewriteLinksOutsideCode(stripLeadingH1(source));
assertMdxSafe(body);

const output = `${FRONTMATTER}\n${body}`;

mkdirSync(dirname(OUTPUT), { recursive: true });
writeFileSync(OUTPUT, output);
