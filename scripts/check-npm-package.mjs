#!/usr/bin/env node
// Verifies the launcher npm package publishes the repository README.

import { execFileSync } from "node:child_process";
import { readFileSync, statSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const root = join(dirname(fileURLToPath(import.meta.url)), "..");
const packageDir = join(root, "npm", "quality.md");
const rootReadmePath = join(root, "README.md");
const packageReadmePath = join(packageDir, "README.md");

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
