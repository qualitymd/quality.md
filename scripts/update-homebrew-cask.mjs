#!/usr/bin/env node
// Updates qualitymd/homebrew-tap Casks/qualitymd.rb from release checksums.

import { readFileSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import {
  assertExpectedChecksums,
  downloadGitHubAssetText,
  githubJSON,
  parseChecksums,
  run,
  tapRepo,
  targets,
  tempDir,
  tokenEnv,
  versionFromTag,
} from "./lib/release.mjs";

try {
  const tag = process.argv[2];
  if (!tag) {
    throw new Error("usage: node scripts/update-homebrew-cask.mjs <tag>");
  }
  const version = versionFromTag(tag);
  const tapToken = tokenEnv("HOMEBREW_TAP_GITHUB_TOKEN");
  const githubToken = process.env.GITHUB_TOKEN || process.env.GH_TOKEN || "";
  const release = await githubJSON(`/repos/qualitymd/quality.md/releases/tags/${tag}`, {
    token: githubToken,
  });
  const checksumsAsset = release.assets.find((asset) => asset.name === "checksums.txt");
  if (!checksumsAsset) {
    throw new Error(`${tag} release is missing checksums.txt`);
  }
  const checksums = parseChecksums(await downloadGitHubAssetText(checksumsAsset, { token: githubToken }));
  assertExpectedChecksums(checksums);
  updateTap(version, tapToken, checksums);
} catch (error) {
  console.error(error.message);
  process.exit(1);
}

function updateTap(version, token, checksums) {
  const dir = tempDir("qualitymd-tap-");
  try {
    const checkout = join(dir.path, "homebrew-tap");
    const authHeader = `AUTHORIZATION: bearer ${token}`;
    run("git", ["-c", `http.https://github.com/.extraheader=${authHeader}`, "clone", `https://github.com/${tapRepo}.git`, checkout], {
      stdio: "ignore",
    });
    run("git", ["config", "user.name", "qualitymd-release"], { cwd: checkout, stdio: "ignore" });
    run("git", ["config", "user.email", "release@qualitymd.local"], { cwd: checkout, stdio: "ignore" });

    const caskPath = join(checkout, "Casks", "qualitymd.rb");
    let cask = readFileSync(caskPath, "utf8");
    cask = cask.replace(/version "[^"]+"/, `version "${version}"`);
    for (const target of targets.filter((target) => target.os !== "win32")) {
      const checksum = checksums.get(target.asset);
      const assetPattern = target.asset.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
      const pattern = new RegExp(
        `(sha256 ")[0-9a-f]{64}("\\n\\s+url "https://github\\.com/qualitymd/quality\\.md/releases/download/v#\\{version\\}/${assetPattern}")`,
      );
      if (!pattern.test(cask)) {
        throw new Error(`could not find cask stanza for ${target.asset}`);
      }
      cask = cask.replace(pattern, `$1${checksum}$2`);
    }
    writeFileSync(caskPath, cask);

    if (!runHasDiff(checkout)) {
      console.log(`Homebrew cask already current for v${version}`);
      return;
    }

    run("git", ["add", "Casks/qualitymd.rb"], { cwd: checkout });
    run("git", ["commit", "-m", `Update qualitymd cask to v${version}`], { cwd: checkout });
    run("git", ["-c", `http.https://github.com/.extraheader=${authHeader}`, "push", "origin", "HEAD:main"], {
      cwd: checkout,
    });
    console.log(`Homebrew cask updated for v${version}`);
  } finally {
    dir.remove();
  }
}

function runHasDiff(cwd) {
  try {
    run("git", ["diff", "--quiet"], { cwd, stdio: "ignore" });
    return false;
  } catch {
    return true;
  }
}
