#!/usr/bin/env node

import { spawnSync } from "node:child_process";
import { createHash } from "node:crypto";
import { readdirSync, readFileSync, statSync } from "node:fs";
import { resolve } from "node:path";

const separator = process.argv.indexOf("--", 2);
if (separator < 3 || separator === process.argv.length - 1) {
  console.error(
    "usage: node scripts/check-generated.mjs <command> [args...] -- <output> [outputs...]",
  );
  process.exit(2);
}

const command = process.argv[2];
const arguments_ = process.argv.slice(3, separator);
const outputs = process.argv.slice(separator + 1).map((path) => resolve(path));

const files = (path) => {
  const info = statSync(path);
  if (info.isFile()) return [path];
  return readdirSync(path, { withFileTypes: true })
    .flatMap((entry) => files(resolve(path, entry.name)))
    .sort();
};

const snapshot = () =>
  new Map(
    outputs.flatMap((output) =>
      files(output).map((path) => [
        path,
        createHash("sha256").update(readFileSync(path)).digest("hex"),
      ]),
    ),
  );

const before = snapshot();
const result = spawnSync(command, arguments_, { stdio: "inherit" });
if (result.error) throw result.error;
if (result.status !== 0) process.exit(result.status ?? 1);
const after = snapshot();

const changed = new Set([...before.keys(), ...after.keys()]);
for (const path of changed) {
  if (before.get(path) !== after.get(path)) {
    console.error(`generated output is stale: ${path}`);
    process.exit(1);
  }
}
