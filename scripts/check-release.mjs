#!/usr/bin/env node
// Runs pre-tag release checks for a prepared release commit.

import { execFileSync } from "node:child_process";
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
  if (!hasUnreleasedSection(changelog)) {
    throw new Error("CHANGELOG.md must contain a top-level ## Unreleased section");
  }

  const { body } = extractReleaseSection(changelog, tag);
  assertCompatibilityBlock(tag, body);

  if (!skipGates) {
    run("mise", ["run", "fmt"]);
    run("mise", ["run", "test"]);
    run("mise", ["run", "vet"]);
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

function assertCompatibilityBlock(releaseTag, body) {
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
}
