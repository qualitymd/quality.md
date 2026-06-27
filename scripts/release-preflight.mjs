#!/usr/bin/env node
// Verifies release targets and publish credentials before tagging.

import {
  githubJSON,
  githubMaybeJSON,
  launcherPackage,
  mainRepo,
  npmUserConfig,
  output,
  packageVersionExists,
  run,
  tapRepo,
  targets,
  tokenEnv,
  versionFromTag,
} from "./lib/release.mjs";

try {
  const args = process.argv.slice(2);
  const tag = args.find((arg) => !arg.startsWith("--"));
  const credentialsOnly = args.includes("--credentials-only");
  if (!tag) {
    throw new Error("usage: node scripts/release-preflight.mjs <tag> [--credentials-only]");
  }
  const version = versionFromTag(tag);

  if (!credentialsOnly) {
    assertGitTargetAvailable(tag);
    await assertGitHubReleaseAbsent(tag);
  } else {
    await assertGitHubAccess();
  }
  await assertNPMReady(version, { credentialsOnly });
  await assertHomebrewTapWritable(version);

  console.log(`release preflight passed for ${tag}`);
} catch (error) {
  console.error(error.message);
  process.exit(1);
}

function assertGitTargetAvailable(tag) {
  if (output("git", ["tag", "--list", tag]).trim()) {
    throw new Error(`local tag ${tag} already exists`);
  }
  if (output("git", ["ls-remote", "--tags", "origin", tag]).trim()) {
    throw new Error(`remote tag ${tag} already exists on origin`);
  }
}

async function assertGitHubReleaseAbsent(tag) {
  const token = process.env.GITHUB_TOKEN || process.env.GH_TOKEN;
  if (!token) {
    throw new Error("GITHUB_TOKEN or GH_TOKEN is required for release preflight");
  }
  await githubJSON(`/repos/${mainRepo}`, { token });
  const release = await githubMaybeJSON(`/repos/${mainRepo}/releases/tags/${tag}`, { token });
  if (release.ok) {
    throw new Error(`GitHub release ${tag} already exists`);
  }
  if (release.status !== 404) {
    throw new Error(`could not verify GitHub release absence for ${tag}: ${release.status}`);
  }
}

async function assertGitHubAccess() {
  const token = process.env.GITHUB_TOKEN || process.env.GH_TOKEN;
  if (!token) {
    throw new Error("GITHUB_TOKEN or GH_TOKEN is required for release preflight");
  }
  await githubJSON(`/repos/${mainRepo}`, { token });
}

async function assertNPMReady(version, { credentialsOnly = false } = {}) {
  const token = tokenEnv("NPM_TOKEN");
  const npm = npmUserConfig(token);
  try {
    run("npm", ["whoami"], { env: npm.env, stdio: "ignore" });
    if (credentialsOnly) {
      return;
    }
    for (const pkg of [launcherPackage, ...targets.map((target) => target.npm)]) {
      if (packageVersionExists(pkg, version, npm.env)) {
        throw new Error(`npm package ${pkg}@${version} already exists`);
      }
    }
  } finally {
    npm.cleanup();
  }
}

async function assertHomebrewTapWritable(version) {
  const token = tokenEnv("HOMEBREW_TAP_GITHUB_TOKEN");
  const repo = await githubJSON(`/repos/${tapRepo}`, { token });
  const branch = repo.default_branch || "main";
  const ref = await githubJSON(`/repos/${tapRepo}/git/ref/heads/${branch}`, { token });
  const tempBranch = `release-preflight-${version}-${Date.now()}`;
  let created = false;
  try {
    await githubJSON(`/repos/${tapRepo}/git/refs`, {
      token,
      method: "POST",
      body: {
        ref: `refs/heads/${tempBranch}`,
        sha: ref.object.sha,
      },
    });
    created = true;
  } finally {
    if (created) {
      await githubJSON(`/repos/${tapRepo}/git/refs/heads/${tempBranch}`, {
        token,
        method: "DELETE",
      });
    }
  }
}
