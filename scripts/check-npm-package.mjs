#!/usr/bin/env node
// Verifies the launcher npm package publishes the repository README.

import { execFileSync } from "node:child_process";
import { existsSync, readFileSync, readdirSync, statSync } from "node:fs";
import { dirname, join, normalize, relative, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const root = join(dirname(fileURLToPath(import.meta.url)), "..");
const packageDir = join(root, "npm", "quality.md");
const rootReadmePath = join(root, "README.md");
const packageReadmePath = join(packageDir, "README.md");
const skillDir = join(root, "skills", "quality");

const output = execFileSync("npm", ["pack", "--dry-run", "--json", "--silent"], {
  cwd: packageDir,
  encoding: "utf8",
});
const [pack] = JSON.parse(output);
const readme = pack.files.find((file) => file.path === "README.md");
const expectedSize = statSync(rootReadmePath).size;

if (!readme) {
  throw new Error("npm package does not include README.md");
}

if (readme.size !== expectedSize) {
  throw new Error(
    `npm package README.md is ${readme.size} bytes, expected ${expectedSize}`,
  );
}

if (
  readFileSync(packageReadmePath, "utf8") !== readFileSync(rootReadmePath, "utf8")
) {
  throw new Error("npm package README.md does not match repository README.md");
}

console.log(
  `npm package README.md matches repository README.md (${expectedSize} bytes)`,
);

checkSkillRelativeLinks(skillDir);

function checkSkillRelativeLinks(dir) {
  const broken = [];
  for (const file of markdownFiles(dir)) {
    const body = readFileSync(file, "utf8");
    for (const target of relativeMarkdownLinks(body)) {
      const withoutAnchor = target.split("#", 1)[0];
      if (!withoutAnchor) {
        continue;
      }
      const resolved = resolve(dirname(file), withoutAnchor);
      const rel = relative(dir, resolved);
      if (rel.startsWith("..") || normalize(rel) === "..") {
        broken.push(`${relative(root, file)} -> ${target} escapes skill bundle`);
        continue;
      }
      if (!existsSync(resolved)) {
        broken.push(`${relative(root, file)} -> ${target} does not resolve`);
      }
    }
  }
  if (broken.length > 0) {
    throw new Error(`Broken skills/quality relative links:\n${broken.join("\n")}`);
  }
  console.log("skills/quality relative links resolve within the bundle");
}

function markdownFiles(dir) {
  const files = [];
  for (const entry of readdirSync(dir, { withFileTypes: true })) {
    const path = join(dir, entry.name);
    if (entry.isDirectory()) {
      files.push(...markdownFiles(path));
    } else if (entry.isFile() && entry.name.endsWith(".md")) {
      files.push(path);
    }
  }
  return files;
}

function relativeMarkdownLinks(body) {
  const links = [];
  const linkPattern = /!?\[[^\]]*\]\(([^)]+)\)/g;
  for (const match of body.matchAll(linkPattern)) {
    const raw = match[1].trim();
    const target = raw.replace(/^<|>$/g, "").split(/\s+/, 1)[0];
    if (
      !target ||
      target.startsWith("#") ||
      target.startsWith("http://") ||
      target.startsWith("https://") ||
      target.startsWith("mailto:") ||
      target.startsWith("/")
    ) {
      continue;
    }
    links.push(target);
  }
  return links;
}
