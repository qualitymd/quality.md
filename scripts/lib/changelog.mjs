import { readFileSync } from "node:fs";

export function assertReleaseTag(tag) {
  if (!/^v\d+\.\d+\.\d+(?:[-+][0-9A-Za-z.-]+)?$/.test(tag)) {
    throw new Error(`release tag must look like v1.2.3, got ${JSON.stringify(tag)}`);
  }
}

export function readChangelog(path) {
  return readFileSync(path, "utf8");
}

export function extractReleaseSection(changelog, tag) {
  assertReleaseTag(tag);

  const matches = [...changelog.matchAll(/^## (.+)$/gm)];
  for (let i = 0; i < matches.length; i++) {
    const match = matches[i];
    const heading = match[1].trim();
    const expected = new RegExp(`^${escapeRegExp(tag)} - \\d{4}-\\d{2}-\\d{2}$`);
    if (!expected.test(heading)) {
      continue;
    }

    const next = matches[i + 1];
    const start = match.index;
    const end = next ? next.index : changelog.length;
    const section = changelog.slice(start, end).trimEnd();
    const body = changelog.slice(match.index + match[0].length, end).trim();

    assertSubstantiveReleaseNotes(tag, body);
    return { heading, section, body };
  }

  throw new Error(`CHANGELOG.md does not contain a section for ${tag}`);
}

export function hasUnreleasedSection(changelog) {
  return /^## Unreleased$/m.test(changelog);
}

export function assertSubstantiveReleaseNotes(tag, body) {
  const substantiveLines = body
    .split(/\r?\n/)
    .map((line) => line.trim())
    .filter((line) => {
      if (!line) return false;
      if (/^#{3,6} /.test(line)) return false;
      if (/^<!--.*-->$/.test(line)) return false;
      return true;
    });

  if (substantiveLines.length === 0) {
    throw new Error(`${tag} release notes contain only headings or placeholders`);
  }
}

function escapeRegExp(value) {
  return value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
}
