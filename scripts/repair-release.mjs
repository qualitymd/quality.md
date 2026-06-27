#!/usr/bin/env node
// Repairs safe-to-rerun release channels for a tag, then verifies.

import { join } from "node:path";
import { output, run, tempDir, versionFromTag } from "./lib/release.mjs";

try {
  const tag = process.argv[2];
  if (!tag) {
    throw new Error("usage: node scripts/repair-release.mjs <tag>");
  }
  const version = versionFromTag(tag);
  const notes = tempDir("qualitymd-release-notes-");
  try {
    const notesPath = join(notes.path, "release-notes.md");
    const text = output("node", ["scripts/extract-release-notes.mjs", tag]);
    await import("node:fs").then(({ writeFileSync }) => writeFileSync(notesPath, text));
    run("gh", ["release", "edit", tag, "--notes-file", notesPath]);
  } finally {
    notes.remove();
  }

  run("node", ["scripts/update-homebrew-cask.mjs", tag]);
  run("node", ["scripts/build-npm.mjs", version, "--publish", "--skip-existing"]);
  run("node", ["scripts/release-verify.mjs", tag]);
  console.log(`release repair completed for ${tag}`);
} catch (error) {
  console.error(error.message);
  process.exit(1);
}
