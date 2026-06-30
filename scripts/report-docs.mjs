#!/usr/bin/env node
// Generates the Mintlify pages for the example quality evaluation from the
// report gallery.
//
// The report gallery is the single source of truth. This script renders every
// Markdown page of the generated LedgerLite evaluation into Mintlify pages under
// mintlify/examples/software-service/, rewrites cross-links to internal docs
// routes (and out-of-site targets to GitHub), and registers a single
// "Example quality evaluation" entry in the site navigation. Run via
// `mise run report-docs`; the pre-commit hook and `mise run check` keep the
// generated pages in sync.
//
// Do not edit mintlify/examples/software-service/** by hand. Regenerate the
// gallery (`mise run report-gallery`), then run `mise run report-docs`.

import {
  existsSync,
  mkdirSync,
  readdirSync,
  readFileSync,
  rmSync,
  writeFileSync,
} from "node:fs";
import { dirname, join, relative, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const root = join(dirname(fileURLToPath(import.meta.url)), "..");

const SRC_ROOT = join(
  root,
  "examples/report-gallery/software-service/.quality/evaluations/0001-full-eval",
);
const OUT_ROOT = join(root, "mintlify", "examples", "software-service");
const DOCS_JSON = join(root, "mintlify", "docs.json");

// Mintlify route for a page is its path under mintlify/ without the extension.
const ROUTE_BASE = "/examples/software-service";
const NAV_PAGE = "examples/software-service/report";
const NAV_GROUP = "Examples";
const SIDEBAR_TITLE = "Example quality evaluation";

const REPO_BLOB_BASE = "https://github.com/qualitymd/quality.md/blob/main";

// The gallery's relative `glossary.md` target does not resolve (the gallery ships
// no glossary); point every glossary link at the repository-root glossary.
const GLOSSARY_URL = `${REPO_BLOB_BASE}/glossary.md`;

// A short note prepended to the entry (report) page, marking the evaluation as a
// synthetic, illustrative example rather than a real assessment.
const REPORT_NOTE =
  "> **Synthetic example.** This is a generated, illustrative quality" +
  " evaluation for a fictional service. Ratings, findings, recommendations, and" +
  " evidence are synthetic and demonstrate report structure only.\n\n";

function main() {
  const files = collectMarkdown(SRC_ROOT);
  const pageSet = new Set(files); // absolute paths of every published page

  rmSync(OUT_ROOT, { recursive: true, force: true });

  for (const file of files) {
    const rel = relative(SRC_ROOT, file); // e.g. areas/api/api-area.md
    const out = join(OUT_ROOT, rel).replace(/\.md$/, ".mdx");
    const page = renderPage(file, pageSet);
    mkdirSync(dirname(out), { recursive: true });
    writeFileSync(out, page);
  }

  updateNavigation();
}

function collectMarkdown(dir) {
  const out = [];
  for (const entry of readdirSync(dir, { withFileTypes: true })) {
    const full = join(dir, entry.name);
    if (entry.isDirectory()) {
      out.push(...collectMarkdown(full));
    } else if (entry.isFile() && entry.name.endsWith(".md")) {
      out.push(full);
    }
  }
  return out.sort();
}

function renderPage(file, pageSet) {
  const source = readFileSync(file, "utf8");
  const { frontmatter, body } = splitFrontmatter(source);
  const title = frontmatter.title ?? leadingH1(body);
  if (!title) {
    throw new Error(`report-docs: no title for ${relative(root, file)}`);
  }

  const rewritten = rewriteLinksOutsideCode(
    stripLeadingH1(body),
    dirname(file),
    pageSet,
  );
  assertMdxSafe(rewritten, file);

  const isReport = relative(SRC_ROOT, file) === "report.md";
  const lines = [`title: ${yamlString(title)}`];
  if (isReport) {
    lines.push(`sidebarTitle: ${yamlString(SIDEBAR_TITLE)}`);
  }

  const head =
    `---\n${lines.join("\n")}\n---\n\n` +
    "{/* Generated from the report gallery by scripts/report-docs.mjs. Do not edit directly. */}\n\n";

  return head + (isReport ? REPORT_NOTE : "") + rewritten.replace(/^\n+/, "");
}

// splitFrontmatter parses a leading `--- ... ---` YAML block. Only the `title`
// key is needed, so the parse is intentionally minimal (flat scalar keys).
function splitFrontmatter(markdown) {
  if (!markdown.startsWith("---\n")) {
    return { frontmatter: {}, body: markdown };
  }
  const end = markdown.indexOf("\n---", 3);
  if (end === -1) {
    return { frontmatter: {}, body: markdown };
  }
  const block = markdown.slice(4, end);
  const body = markdown.slice(end + 4).replace(/^\n+/, "");
  const frontmatter = {};
  for (const line of block.split("\n")) {
    const match = /^([A-Za-z0-9_]+):\s*(.*)$/.exec(line);
    if (match) {
      frontmatter[match[1]] = unquote(match[2].trim());
    }
  }
  return { frontmatter, body };
}

function unquote(value) {
  if (
    (value.startsWith('"') && value.endsWith('"')) ||
    (value.startsWith("'") && value.endsWith("'"))
  ) {
    return value.slice(1, -1);
  }
  return value;
}

function leadingH1(body) {
  const match = /^#\s+(.+)$/m.exec(body);
  return match ? match[1].trim() : "";
}

function stripLeadingH1(body) {
  const trimmed = body.replace(/^\n+/, "");
  if (!trimmed.startsWith("# ")) {
    throw new Error("report-docs: page does not start with an H1 heading.");
  }
  const newline = trimmed.indexOf("\n");
  return trimmed.slice(newline + 1).replace(/^\n+/, "");
}

// yamlString quotes a frontmatter value when needed (the report titles contain a
// colon, which bare YAML scalars cannot carry).
function yamlString(value) {
  if (/^[^"':#{}[\],&*!|>%@`]+$/.test(value) && !/:\s|:$/.test(value)) {
    return value;
  }
  return `"${value.replace(/"/g, '\\"')}"`;
}

function rewriteLink(dest, fileDir, pageSet) {
  if (dest.startsWith("#") || /^[a-z][a-z0-9+.-]*:/i.test(dest)) {
    return dest; // in-page anchor or already-absolute URL
  }

  const hash = dest.indexOf("#");
  const path = hash === -1 ? dest : dest.slice(0, hash);
  const anchor = hash === -1 ? "" : dest.slice(hash);

  if (/(^|\/)glossary\.md$/.test(path)) {
    return GLOSSARY_URL;
  }

  const target = resolve(fileDir, path);

  if (path.endsWith(".md") && pageSet.has(target)) {
    const rel = relative(SRC_ROOT, target).replace(/\.md$/, "");
    return `${ROUTE_BASE}/${toPosix(rel)}${anchor}`;
  }

  // Anything else (data/**.json, a missing target) lives in the repo; link to it
  // on GitHub.
  return `${REPO_BLOB_BASE}/${toPosix(relative(root, target))}${anchor}`;
}

function toPosix(p) {
  return p.split(/[\\/]/).join("/");
}

function rewriteLinks(markdown, fileDir, pageSet) {
  // Inline links: [text](dest). Destinations here never carry titles or spaces.
  return markdown.replace(
    /(\]\()([^)\s]+)(\))/g,
    (_, open, dest, close) =>
      `${open}${rewriteLink(dest, fileDir, pageSet)}${close}`,
  );
}

// Rewrite links only in prose. Code fences and inline code are literal in MDX and
// can contain `](` sequences that are not links, so they are left untouched.
function rewriteLinksOutsideCode(markdown, fileDir, pageSet) {
  const codeSpan = /(```[\s\S]*?```|`[^`\n]*`)/g;
  let result = "";
  let last = 0;
  let match;
  while ((match = codeSpan.exec(markdown)) !== null) {
    result += rewriteLinks(markdown.slice(last, match.index), fileDir, pageSet);
    result += match[0];
    last = match.index + match[0].length;
  }
  result += rewriteLinks(markdown.slice(last), fileDir, pageSet);
  return result;
}

// Fail loudly on characters MDX parses as JSX outside code spans, rather than
// emit a page that breaks the Mintlify build. Code spans are literal in MDX, and
// the `<a id="…"></a>` finding anchors are intentional valid JSX, so both are
// removed before the scan.
function assertMdxSafe(markdown, file) {
  const scrubbed = markdown
    .replace(/```[\s\S]*?```/g, "")
    .replace(/`[^`\n]*`/g, "")
    .replace(/<a id="[^"]*"><\/a>/g, "");

  const problems = [];
  scrubbed.split("\n").forEach((line, index) => {
    if (/<[A-Za-z/!]/.test(line) || line.includes("{")) {
      problems.push(`  line ${index + 1}: ${line.trim()}`);
    }
  });

  if (problems.length > 0) {
    throw new Error(
      `report-docs: ${relative(root, file)} contains characters MDX parses ` +
        "as JSX outside code spans:\n" +
        problems.join("\n"),
    );
  }
}

// updateNavigation appends a single-page `Examples` group to the navigation in
// docs.json. The insert is textual so the rest of the hand-maintained file —
// including its compact `pages` arrays — is preserved byte for byte; a
// JSON round-trip would reformat every array. The group is a fixed constant, so
// once present the file is already current and is left untouched (idempotent).
function updateNavigation() {
  const raw = readFileSync(DOCS_JSON, "utf8");

  // Validate structure (and fail loudly on a malformed file) before editing text.
  const groups = JSON.parse(raw).navigation?.groups;
  if (!Array.isArray(groups)) {
    throw new Error("report-docs: docs.json has no navigation.groups array.");
  }
  if (groups.some((g) => g.group === NAV_GROUP)) {
    return;
  }

  const open = raw.indexOf("[", raw.indexOf('"groups"'));
  const close = matchingBracket(raw, open);
  const lastBrace = raw.lastIndexOf("}", close);

  const block =
    `      {\n` +
    `        "group": ${JSON.stringify(NAV_GROUP)},\n` +
    `        "pages": [${JSON.stringify(NAV_PAGE)}]\n` +
    `      }`;

  const next = `${raw.slice(0, lastBrace + 1)},\n${block}${raw.slice(lastBrace + 1)}`;
  writeFileSync(DOCS_JSON, next);
}

// matchingBracket returns the index of the `]` that closes the `[` at openIndex.
function matchingBracket(text, openIndex) {
  let depth = 0;
  for (let i = openIndex; i < text.length; i++) {
    if (text[i] === "[") depth++;
    else if (text[i] === "]" && --depth === 0) return i;
  }
  throw new Error("report-docs: unbalanced brackets in docs.json navigation.");
}

if (!existsSync(SRC_ROOT)) {
  throw new Error(
    `report-docs: source evaluation not found at ${relative(root, SRC_ROOT)}; ` +
      "run `mise run report-gallery` first.",
  );
}

main();
