#!/usr/bin/env node
// Builds the npm distribution for QUALITY.md.
//
// For each supported platform it cross-compiles the Go binary into a
// per-platform package (@qualitymd/cli-<os>-<arch>) gated by npm `os`/`cpu`
// fields, then stamps the shared version across every package.json.
//
// Usage:
//   node scripts/build-npm.mjs <version> [--publish]
//
//   <version>   semver to stamp (e.g. 1.2.3); defaults to 0.0.0-dev
//   --publish          run `npm publish` for each platform package and the launcher
//   --skip-existing    with --publish, skip packages already present on npm
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
const skipExisting = args.includes("--skip-existing");
const version = args.find((a) => !a.startsWith("--")) ?? "0.0.0-dev";

function run(cmd, cmdArgs, { env, cwd } = {}) {
  execFileSync(cmd, cmdArgs, {
    stdio: "inherit",
    cwd,
    env: { ...process.env, ...env },
  });
}

function platformName(t) {
  const osName = {
    darwin: "macOS",
    linux: "Linux",
    win32: "Windows",
  }[t.os];

  return `${osName} ${t.arch}`;
}

function platformPackageDescription(t) {
  return `Native qualitymd binary for ${platformName(t)}. Installed automatically by quality.md.`;
}

function platformPackageReadme(pkgName, t) {
  return `# ${pkgName}

Native \`qualitymd\` binary package for ${platformName(t)}.

This package is installed automatically by the user-facing \`quality.md\` npm
package. Most users should not install it directly.

Install the CLI with:

\`\`\`sh
npm install -g quality.md
\`\`\`

or run it with:

\`\`\`sh
npx quality.md --version
\`\`\`

See the main package and project docs:

- npm: https://www.npmjs.com/package/quality.md
- GitHub: https://github.com/qualitymd/quality.md
`;
}

function packageVersionExists(pkgName) {
  try {
    const raw = execFileSync("npm", ["view", `${pkgName}@${version}`, "version", "--json"], {
      encoding: "utf8",
      stdio: ["ignore", "pipe", "pipe"],
    }).trim();
    return JSON.parse(raw) === version;
  } catch {
    return false;
  }
}

function publishPackage(pkgName, cwd) {
  if (skipExisting && packageVersionExists(pkgName)) {
    console.log(`  ${pkgName}@${version} already published; skipping`);
    return;
  }
  run("npm", ["publish", "--access", "public"], { cwd });
}

console.log(`Building QUALITY.md npm packages @ ${version}`);
run(process.execPath, [join(root, "scripts", "sync-npm-readme.mjs")], {
  cwd: root,
});

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
        description: platformPackageDescription(t),
        homepage: "https://getquality.md",
        keywords: ["qualitymd", "quality.md", "cli", "native-binary"],
        license: "MIT",
        repository: { type: "git", url: "https://github.com/qualitymd/quality.md.git" },
        bugs: { url: "https://github.com/qualitymd/quality.md/issues" },
        os: [t.os],
        cpu: [t.arch],
        files: ["bin"],
        publishConfig: { access: "public" },
      },
      null,
      2,
    ) + "\n",
  );
  writeFileSync(join(pkgDir, "README.md"), platformPackageReadme(pkgName, t));

  if (publish) {
    publishPackage(pkgName, pkgDir);
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
  publishPackage(launcher.name, join(root, "npm", "quality.md"));
}
