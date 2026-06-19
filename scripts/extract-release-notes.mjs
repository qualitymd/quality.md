#!/usr/bin/env node
// Prints the curated CHANGELOG.md section for a release tag.

import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import { extractReleaseSection, readChangelog } from "./lib/changelog.mjs";

const root = join(dirname(fileURLToPath(import.meta.url)), "..");

try {
  const tag = process.argv[2];

  if (!tag) {
    throw new Error("usage: node scripts/extract-release-notes.mjs <tag>");
  }

  const changelog = readChangelog(join(root, "CHANGELOG.md"));
  const { section } = extractReleaseSection(changelog, tag);

  process.stdout.write(`${section}\n`);
} catch (error) {
  console.error(error.message);
  process.exit(1);
}
