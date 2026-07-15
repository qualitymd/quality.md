#!/usr/bin/env node
// Verifies all public release channels for a release tag.

import { chmodSync, mkdtempSync, readFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { extractReleaseSection, readChangelog } from "./lib/changelog.mjs";
import {
  assertExpectedChecksums,
  currentPlatformTarget,
  downloadGitHubAssetFile,
  downloadGitHubAssetText,
  downloadText,
  expectedReleaseAssets,
  githubReleaseByTag,
  launcherPackage,
  output,
  packageVersionExists,
  parseChecksums,
  root,
  run,
  targets,
  versionFromTag,
} from "./lib/release.mjs";

try {
  const args = process.argv.slice(2);
  const tag = args.find((arg) => !arg.startsWith("--"));
  const allowDraft = args.includes("--allow-draft");
  if (!tag) {
    throw new Error("usage: node scripts/release-verify.mjs <tag> [--allow-draft]");
  }
  const version = versionFromTag(tag);

  const release = await verifyGitHubRelease(tag, allowDraft);
  const checksums = await verifyReleaseAssets(release);
  await verifyReleaseNotes(tag, release);
  await verifyNPM(version);
  await verifyHomebrew(version, checksums);
  await verifyNativeBinary(version, release);

  console.log(`release verification passed for ${tag}`);
} catch (error) {
  console.error(error.message);
  process.exit(1);
}

async function verifyGitHubRelease(tag, allowDraft) {
  const token = process.env.GITHUB_TOKEN || process.env.GH_TOKEN || "";
  const release = await githubReleaseByTag(tag, { token });
  if (release.draft && !allowDraft) {
    throw new Error(`${tag} GitHub release is still draft`);
  }
  if (release.prerelease) {
    throw new Error(`${tag} GitHub release is marked prerelease`);
  }
  return release;
}

async function verifyReleaseAssets(release) {
  const names = release.assets.map((asset) => asset.name).sort();
  const got = JSON.stringify(names);
  const want = JSON.stringify(expectedReleaseAssets);
  if (got !== want) {
    throw new Error(`release assets mismatch: got ${got}, want ${want}`);
  }
  const checksumsAsset = release.assets.find((asset) => asset.name === "checksums.txt");
  if (!checksumsAsset) {
    throw new Error("release is missing checksums.txt");
  }
  const checksums = parseChecksums(await downloadGitHubAssetText(checksumsAsset));
  assertExpectedChecksums(checksums);
  return checksums;
}

async function verifyReleaseNotes(tag, release) {
  const changelog = readChangelog(join(root, "CHANGELOG.md"));
  const { section } = extractReleaseSection(changelog, tag);
  if (trimFinalNewline(release.body || "") !== trimFinalNewline(section)) {
    throw new Error(`${tag} GitHub release notes do not match CHANGELOG.md`);
  }
}

async function verifyNPM(version) {
  for (const pkg of [launcherPackage, ...targets.map((target) => target.npm)]) {
    await waitForPackageVersion(pkg, version);
  }
}

async function waitForPackageVersion(pkg, version) {
  for (let attempt = 1; attempt <= 6; attempt += 1) {
    if (packageVersionExists(pkg, version)) {
      return;
    }
    if (attempt < 6) {
      await new Promise((resolve) => setTimeout(resolve, 5000));
    }
  }
  throw new Error(`npm package ${pkg}@${version} is not published`);
}

async function verifyHomebrew(version, checksums) {
  const content = await downloadText(
    "https://raw.githubusercontent.com/qualitymd/homebrew-tap/main/Casks/qualitymd.rb",
    { token: "" },
  );
  if (!content.includes(`version "${version}"`)) {
    throw new Error(`Homebrew cask does not declare version ${version}`);
  }
  for (const target of targets.filter((target) => target.homebrew)) {
    const checksum = checksums.get(target.asset);
    if (!content.includes(`sha256 "${checksum}"`)) {
      throw new Error(`Homebrew cask checksum mismatch for ${target.asset}`);
    }
  }
}

async function verifyNativeBinary(version, release) {
  const target = currentPlatformTarget();
  if (!target) {
    console.log(`skipping native binary check for unsupported ${process.platform}/${process.arch}`);
    return;
  }
  const asset = release.assets.find((candidate) => candidate.name === target.asset);
  if (!asset) {
    throw new Error(`missing native asset ${target.asset}`);
  }
  const dir = mkdtempSync(join(tmpdir(), "qualitymd-release-bin-"));
  const archive = join(dir, target.asset);
  await downloadGitHubAssetFile(asset, archive);
  if (target.asset.endsWith(".zip")) {
    run("unzip", ["-q", archive, "-d", dir]);
  } else {
    run("tar", ["-xzf", archive, "-C", dir]);
  }
  const binary = join(dir, process.platform === "win32" ? "qualitymd.exe" : "qualitymd");
  chmodSync(binary, 0o755);
  const info = JSON.parse(output(binary, ["version", "--json"], { cwd: dir }));
  if (info.version !== version) {
    throw new Error(`downloaded binary reports ${info.version}, want ${version}`);
  }
  const spec = output(binary, ["spec"], { cwd: dir });
  const rootSpec = readFileSync(join(root, "SPECIFICATION.md"), "utf8");
  const versionLine = rootSpec.match(/\*\*Specification version:\*\*\s+(.+)/);
  if (versionLine && !spec.includes(`**Specification version:** ${versionLine[1].trim()}`)) {
    throw new Error("downloaded binary spec output does not match SPECIFICATION.md version");
  }
}

function trimFinalNewline(value) {
  return value.replace(/\n+$/, "");
}
