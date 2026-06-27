import { execFileSync } from "node:child_process";
import { mkdirSync, mkdtempSync, readFileSync, rmSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import { assertReleaseTag } from "./changelog.mjs";

export const root = join(dirname(fileURLToPath(import.meta.url)), "../..");
export const mainRepo = "qualitymd/quality.md";
export const tapRepo = "qualitymd/homebrew-tap";

export const targets = [
  {
    os: "darwin",
    arch: "arm64",
    npm: "@qualitymd/cli-darwin-arm64",
    asset: "qualitymd_darwin_arm64.tar.gz",
  },
  {
    os: "darwin",
    arch: "x64",
    npm: "@qualitymd/cli-darwin-x64",
    asset: "qualitymd_darwin_amd64.tar.gz",
  },
  {
    os: "linux",
    arch: "arm64",
    npm: "@qualitymd/cli-linux-arm64",
    asset: "qualitymd_linux_arm64.tar.gz",
  },
  {
    os: "linux",
    arch: "x64",
    npm: "@qualitymd/cli-linux-x64",
    asset: "qualitymd_linux_amd64.tar.gz",
  },
  {
    os: "win32",
    arch: "arm64",
    npm: "@qualitymd/cli-win32-arm64",
    asset: "qualitymd_windows_arm64.zip",
  },
  {
    os: "win32",
    arch: "x64",
    npm: "@qualitymd/cli-win32-x64",
    asset: "qualitymd_windows_amd64.zip",
  },
];

export const launcherPackage = "quality.md";
export const expectedReleaseAssets = [
  "checksums.txt",
  ...targets.map((target) => target.asset),
].sort();

export function versionFromTag(tag) {
  assertReleaseTag(tag);
  return tag.replace(/^v/, "");
}

export function run(cmd, args, { cwd = root, env = {}, stdio = "inherit" } = {}) {
  return execFileSync(cmd, args, {
    cwd,
    env: { ...process.env, ...env },
    stdio,
  });
}

export function output(cmd, args, { cwd = root, env = {} } = {}) {
  return execFileSync(cmd, args, {
    cwd,
    env: { ...process.env, ...env },
    encoding: "utf8",
    stdio: ["ignore", "pipe", "pipe"],
  });
}

export function tempDir(prefix) {
  const path = mkdtempSync(join(tmpdir(), prefix));
  return {
    path,
    remove() {
      rmSync(path, { recursive: true, force: true });
    },
  };
}

export function readJSON(path) {
  return JSON.parse(readFileSync(path, "utf8"));
}

export function writeJSON(path, value) {
  writeFileSync(path, JSON.stringify(value, null, 2) + "\n");
}

export function tokenEnv(name) {
  const value = process.env[name];
  if (!value) {
    throw new Error(`${name} is required`);
  }
  return value;
}

export function githubToken() {
  return process.env.GITHUB_TOKEN || process.env.GH_TOKEN || "";
}

export async function githubJSON(path, { token = githubToken(), method = "GET", body } = {}) {
  const response = await fetch(`https://api.github.com${path}`, {
    method,
    headers: {
      Accept: "application/vnd.github+json",
      "X-GitHub-Api-Version": "2022-11-28",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    ...(body === undefined ? {} : { body: JSON.stringify(body) }),
  });
  const text = await response.text();
  let json = {};
  if (text) {
    try {
      json = JSON.parse(text);
    } catch {
      json = { message: text };
    }
  }
  if (!response.ok) {
    const message = json.message ? `: ${json.message}` : "";
    throw new Error(`GitHub ${method} ${path} failed with ${response.status}${message}`);
  }
  return json;
}

export async function githubMaybeJSON(path, opts = {}) {
  const response = await fetch(`https://api.github.com${path}`, {
    method: opts.method || "GET",
    headers: {
      Accept: "application/vnd.github+json",
      "X-GitHub-Api-Version": "2022-11-28",
      ...(opts.token ? { Authorization: `Bearer ${opts.token}` } : {}),
    },
    ...(opts.body === undefined ? {} : { body: JSON.stringify(opts.body) }),
  });
  const text = await response.text();
  let json = {};
  if (text) {
    try {
      json = JSON.parse(text);
    } catch {
      json = { message: text };
    }
  }
  return { ok: response.ok, status: response.status, json };
}

export async function downloadText(url, { token = githubToken() } = {}) {
  const response = await fetch(url, {
    headers: token ? { Authorization: `Bearer ${token}` } : {},
  });
  if (!response.ok) {
    throw new Error(`download failed with ${response.status}: ${url}`);
  }
  return response.text();
}

export async function downloadGitHubAssetText(asset, { token = githubToken() } = {}) {
  const response = await fetch(asset.url, {
    headers: {
      Accept: "application/octet-stream",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
  });
  if (!response.ok) {
    throw new Error(`GitHub asset download failed with ${response.status}: ${asset.name}`);
  }
  return response.text();
}

export async function downloadFile(url, path, { token = githubToken() } = {}) {
  const response = await fetch(url, {
    headers: token ? { Authorization: `Bearer ${token}` } : {},
  });
  if (!response.ok) {
    throw new Error(`download failed with ${response.status}: ${url}`);
  }
  const buffer = Buffer.from(await response.arrayBuffer());
  mkdirSync(dirname(path), { recursive: true });
  writeFileSync(path, buffer);
}

export async function downloadGitHubAssetFile(asset, path, { token = githubToken() } = {}) {
  const response = await fetch(asset.url, {
    headers: {
      Accept: "application/octet-stream",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
  });
  if (!response.ok) {
    throw new Error(`GitHub asset download failed with ${response.status}: ${asset.name}`);
  }
  const buffer = Buffer.from(await response.arrayBuffer());
  mkdirSync(dirname(path), { recursive: true });
  writeFileSync(path, buffer);
}

export function parseChecksums(text) {
  const checksums = new Map();
  for (const line of text.split(/\r?\n/)) {
    const match = line.match(/^([0-9a-f]{64})\s+(.+)$/);
    if (match) {
      checksums.set(match[2], match[1]);
    }
  }
  return checksums;
}

export function assertExpectedChecksums(checksums) {
  for (const asset of targets.map((target) => target.asset)) {
    if (!checksums.has(asset)) {
      throw new Error(`checksums.txt is missing ${asset}`);
    }
  }
}

export function npmUserConfig(token) {
  const dir = tempDir("qualitymd-npm-");
  const config = join(dir.path, ".npmrc");
  writeFileSync(config, `//registry.npmjs.org/:_authToken=${token}\n`);
  return {
    env: {
      NPM_CONFIG_USERCONFIG: config,
      NODE_AUTH_TOKEN: token,
    },
    cleanup: () => dir.remove(),
  };
}

export function packageVersionExists(pkg, version, env = {}) {
  try {
    const got = output("npm", ["view", `${pkg}@${version}`, "version", "--json"], { env }).trim();
    return JSON.parse(got) === version;
  } catch {
    return false;
  }
}

export function currentPlatformTarget() {
  const os = process.platform;
  const arch = process.arch;
  return targets.find((target) => target.os === os && target.arch === arch);
}
