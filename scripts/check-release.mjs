#!/usr/bin/env node
// Runs pre-tag release checks for a prepared release commit.

import { execFileSync } from "node:child_process";
import { readFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import {
  assertReleaseTag,
  extractReleaseSection,
  hasUnreleasedSection,
  readChangelog,
} from "./lib/changelog.mjs";

const root = join(dirname(fileURLToPath(import.meta.url)), "..");

try {
  const args = process.argv.slice(2);
  const tag = args.find((arg) => !arg.startsWith("--"));
  const options = new Set(args.filter((arg) => arg.startsWith("--")));

  if (!tag) {
    throw new Error("usage: node scripts/check-release.mjs <tag>");
  }

  assertReleaseTag(tag);

  const skipClean = options.has("--skip-clean");
  const skipTagCheck = options.has("--skip-tag-check");
  const skipGates = options.has("--skip-gates");

  if (!skipClean) {
    assertCleanTree("before release checks");
  }

  if (!skipTagCheck) {
    assertTagAvailable(tag);
  }

  const changelog = readChangelog(join(root, "CHANGELOG.md"));
  const skillMetadata = readSkillMetadata(join(root, "skills/quality/SKILL.md"));
  assertSkillMetadata(tag, skillMetadata);
  if (!hasUnreleasedSection(changelog)) {
    throw new Error("CHANGELOG.md must contain a top-level ## Unreleased section");
  }

  const { body } = extractReleaseSection(changelog, tag);
  assertCompatibilityBlock(tag, body, skillMetadata);

  if (!skipGates) {
    run("mise", ["run", "fmt"]);
    run("mise", ["run", "test"]);
    run("mise", ["run", "vet"]);
    run("mise", ["run", "lint"]);
    run("mise", ["run", "snapshot"]);
    run("mise", ["run", "npm-build"]);
  }

  if (!skipClean) {
    assertCleanTree("after release checks");
  }

  console.log(`release checks passed for ${tag}`);
} catch (error) {
  console.error(error.message);
  process.exit(1);
}

function run(cmd, cmdArgs) {
  execFileSync(cmd, cmdArgs, { cwd: root, stdio: "inherit" });
}

function output(cmd, cmdArgs) {
  return execFileSync(cmd, cmdArgs, {
    cwd: root,
    encoding: "utf8",
    stdio: ["ignore", "pipe", "pipe"],
  });
}

function assertCleanTree(context) {
  const status = output("git", ["status", "--porcelain"]);
  if (status.trim()) {
    throw new Error(`working tree must be clean ${context}`);
  }
}

function assertTagAvailable(releaseTag) {
  const local = output("git", ["tag", "--list", releaseTag]).trim();
  if (local) {
    throw new Error(`local tag ${releaseTag} already exists`);
  }

  const remote = output("git", ["ls-remote", "--tags", "origin", releaseTag]).trim();
  if (remote) {
    throw new Error(`remote tag ${releaseTag} already exists on origin`);
  }
}

function readSkillMetadata(path) {
  const raw = readFileSync(path, "utf8");
  const match = raw.match(/^---\n([\s\S]*?)\n---\n/);
  if (!match) {
    throw new Error("skills/quality/SKILL.md does not contain YAML frontmatter");
  }
  const frontmatter = match[1];
  const metadata = {};
  let compatibility = "";
  let inMetadata = false;
  for (const line of frontmatter.split(/\r?\n/)) {
    if (/^\S/.test(line)) {
      inMetadata = false;
    }
    const compatibilityMatch = line.match(/^compatibility:\s*(.+)$/);
    if (compatibilityMatch) {
      compatibility = unquoteYAMLScalar(compatibilityMatch[1].trim());
      continue;
    }
    if (line.trim() === "metadata:") {
      inMetadata = true;
      continue;
    }
    if (!inMetadata) {
      continue;
    }
    const metadataMatch = line.match(/^\s+([A-Za-z0-9_.-]+):\s*(.+)$/);
    if (metadataMatch) {
      metadata[metadataMatch[1]] = unquoteYAMLScalar(metadataMatch[2].trim());
    }
  }
  return { compatibility, metadata };
}

function unquoteYAMLScalar(value) {
  if (
    (value.startsWith('"') && value.endsWith('"')) ||
    (value.startsWith("'") && value.endsWith("'"))
  ) {
    return value.slice(1, -1);
  }
  return value;
}

function assertSkillMetadata(releaseTag, skill) {
  const expectedVersion = releaseTag.replace(/^v/, "");
  const actualVersion = skill.metadata.version;
  if (actualVersion !== expectedVersion) {
    throw new Error(
      `skills/quality/SKILL.md metadata.version is ${actualVersion || "(missing)"}, expected ${expectedVersion}`,
    );
  }

  const range = skill.metadata["requires-qualitymd-cli"];
  if (!range) {
    throw new Error("skills/quality/SKILL.md metadata.requires-qualitymd-cli is missing");
  }
  if (!/^>=\d+\.\d+\.\d+(?:[-+][0-9A-Za-z.-]+)? <\d+\.\d+\.\d+(?:[-+][0-9A-Za-z.-]+)?$/.test(range)) {
    throw new Error(
      `skills/quality/SKILL.md metadata.requires-qualitymd-cli is not a supported SemVer range: ${range}`,
    );
  }

  if (!skill.compatibility.includes(range)) {
    throw new Error(
      "skills/quality/SKILL.md compatibility prose does not match metadata.requires-qualitymd-cli",
    );
  }
}

function assertCompatibilityBlock(releaseTag, body, skillMetadata) {
  const cliLine = body.match(/^- CLI:\s+`?([^`\n]+)`?$/m);
  if (cliLine && cliLine[1].trim() !== releaseTag) {
    throw new Error(
      `CHANGELOG.md CLI compatibility line is ${cliLine[1].trim()}, expected ${releaseTag}`,
    );
  }

  const specLine = body.match(/^- QUALITY\.md specification:\s+`?([^`\n]+)`?$/m);
  if (specLine) {
    const spec = readChangelog(join(root, "SPECIFICATION.md"));
    const version = spec.match(/\*\*Specification version:\*\*\s+(.+)/);
    if (!version) {
      throw new Error("SPECIFICATION.md does not declare a specification version");
    }

    if (!specLine[1].trim().includes(version[1].trim())) {
      throw new Error(
        `CHANGELOG.md specification compatibility line does not match SPECIFICATION.md (${version[1].trim()})`,
      );
    }
  }

  const skillLine = body.match(/^- \/quality skill:\s+(.+)$/m);
  if (skillLine) {
    const text = skillLine[1].trim();
    const skillVersion = skillMetadata.metadata.version;
    const cliRange = skillMetadata.metadata["requires-qualitymd-cli"];
    if (!text.includes(skillVersion) || !text.includes(cliRange)) {
      throw new Error(
        `CHANGELOG.md /quality skill compatibility line must include skill ${skillVersion} and qualitymd ${cliRange}`,
      );
    }
  }
}
