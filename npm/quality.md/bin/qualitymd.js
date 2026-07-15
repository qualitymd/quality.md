#!/usr/bin/env node
"use strict";

// Launcher for the quality.md CLI. The actual binary ships in a per-platform
// optional dependency (e.g. @qualitymd/cli-darwin-arm64). npm installs only the
// package matching the host's `os`/`cpu`, so we resolve that one and exec it.

const { spawnSync } = require("node:child_process");

const isMusl =
  process.platform === "linux" &&
  !process.report?.getReport?.().header?.glibcVersionRuntime;
const platformKey = `${process.platform}-${process.arch}${isMusl ? "-musl" : ""}`;
const pkg = `@qualitymd/cli-${platformKey}`;
const ext = process.platform === "win32" ? ".exe" : "";

let binary;
try {
  binary = require.resolve(`${pkg}/bin/qualitymd${ext}`);
} catch {
  console.error(
    `quality.md: no prebuilt binary for ${platformKey}.\n` +
      `The optional dependency ${pkg} is missing. If you installed with ` +
      `--no-optional, reinstall without it, or use the managed installer:\n` +
      `  curl -fsSL https://getquality.md/install.sh | sh`,
  );
  process.exit(1);
}

const env = { ...process.env, QUALITYMD_INSTALL_METHOD: "npm" };
const result = spawnSync(binary, process.argv.slice(2), { stdio: "inherit", env });

if (result.error) {
  throw result.error;
}
process.exit(result.status === null ? 1 : result.status);
