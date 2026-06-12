#!/usr/bin/env node
// Builds the npm distribution for quality.md.
//
// For each supported platform it cross-compiles the Go binary into a
// per-platform package (@qualitymd/cli-<os>-<arch>) gated by npm `os`/`cpu`
// fields, then stamps the shared version across every package.json.
//
// Usage:
//   node scripts/build-npm.mjs <version> [--publish]
//
//   <version>   semver to stamp (e.g. 1.2.3); defaults to 0.0.0-dev
//   --publish   run `npm publish` for each platform package and the launcher
//
// Layout produced:
//   npm/quality.md/                 launcher package (committed)
//   npm/platforms/<os>-<arch>/      generated per-platform packages (gitignored)

import { execFileSync } from "node:child_process";
import { mkdirSync, writeFileSync, readFileSync, rmSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const root = join(dirname(fileURLToPath(import.meta.url)), "..");
const MODULE = "github.com/qualitymd/quality.md";

// node platform/arch -> Go GOOS/GOARCH
const TARGETS = [
  { os: "darwin", arch: "arm64", goos: "darwin", goarch: "arm64" },
  { os: "darwin", arch: "x64", goos: "darwin", goarch: "amd64" },
  { os: "linux", arch: "arm64", goos: "linux", goarch: "arm64" },
  { os: "linux", arch: "x64", goos: "linux", goarch: "amd64" },
  { os: "win32", arch: "arm64", goos: "windows", goarch: "arm64" },
  { os: "win32", arch: "x64", goos: "windows", goarch: "amd64" },
];

const args = process.argv.slice(2);
const publish = args.includes("--publish");
const version = args.find((a) => !a.startsWith("--")) ?? "0.0.0-dev";

function run(cmd, cmdArgs, { env, cwd } = {}) {
  execFileSync(cmd, cmdArgs, {
    stdio: "inherit",
    cwd,
    env: { ...process.env, ...env },
  });
}

console.log(`Building quality.md npm packages @ ${version}`);

for (const t of TARGETS) {
  const key = `${t.os}-${t.arch}`;
  const pkgName = `@qualitymd/cli-${key}`;
  const pkgDir = join(root, "npm", "platforms", key);
  const binDir = join(pkgDir, "bin");
  const ext = t.os === "win32" ? ".exe" : "";

  rmSync(pkgDir, { recursive: true, force: true });
  mkdirSync(binDir, { recursive: true });

  const ldflags = [
    "-s",
    "-w",
    `-X ${MODULE}/internal/cli.version=${version}`,
  ].join(" ");

  console.log(`  ${pkgName}  (${t.goos}/${t.goarch})`);
  run(
    "go",
    [
      "build",
      "-trimpath",
      "-ldflags",
      ldflags,
      "-o",
      join(binDir, `qualitymd${ext}`),
      "./cmd/qualitymd",
    ],
    { env: { GOOS: t.goos, GOARCH: t.goarch, CGO_ENABLED: "0" } },
  );

  writeFileSync(
    join(pkgDir, "package.json"),
    JSON.stringify(
      {
        name: pkgName,
        version,
        description: `quality.md native binary for ${key}`,
        license: "MIT",
        repository: { type: "git", url: "https://github.com/qualitymd/quality.md.git" },
        os: [t.os],
        cpu: [t.arch],
        files: ["bin"],
        publishConfig: { access: "public" },
      },
      null,
      2,
    ) + "\n",
  );

  if (publish) {
    run("npm", ["publish", "--access", "public"], { cwd: pkgDir });
  }
}

// Stamp the launcher package: its own version + every optionalDependency.
const launcherPath = join(root, "npm", "quality.md", "package.json");
const launcher = JSON.parse(readFileSync(launcherPath, "utf8"));
launcher.version = version;
for (const t of TARGETS) {
  launcher.optionalDependencies[`@qualitymd/cli-${t.os}-${t.arch}`] = version;
}
writeFileSync(launcherPath, JSON.stringify(launcher, null, 2) + "\n");
console.log(`Stamped launcher and ${TARGETS.length} platform packages @ ${version}`);

if (publish) {
  console.log("Publishing launcher package...");
  execFileSync("npm", ["publish", "--access", "public"], {
    stdio: "inherit",
    cwd: join(root, "npm", "quality.md"),
  });
}
