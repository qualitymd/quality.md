#!/usr/bin/env bun

import { chmod, mkdir, readdir, rm, stat, writeFile } from "node:fs/promises"
import { basename, join } from "node:path"

import { root, targets, versionFromTag } from "./lib/release.mjs"

const args = process.argv.slice(2)
const tag = args[0]
if (tag === undefined) {
  console.error("usage: bun scripts/build-release.ts <tag> [commit]")
  process.exit(2)
}
const version = versionFromTag(tag)
const commit = args[1] ?? process.env.GITHUB_SHA ?? "none"
const dist = join(root, "dist")
const staging = join(dist, ".build")
const maxExecutableBytes = 100 * 1024 * 1024
const maxArchiveBytes = 40 * 1024 * 1024

const run = async (command: string[]) => {
  const child = Bun.spawn(command, {
    cwd: root,
    stdin: "inherit",
    stdout: "inherit",
    stderr: "inherit",
  })
  const status = await child.exited
  if (status !== 0) throw new Error(`${command.join(" ")} exited with status ${status}`)
}

await rm(dist, { recursive: true, force: true })
await mkdir(staging, { recursive: true })

for (const target of targets) {
  const directory = join(staging, target.asset)
  await mkdir(directory, { recursive: true })
  const binary = join(directory, target.os === "win32" ? "qualitymd.exe" : "qualitymd")
  console.log(`building ${target.asset} (${target.bunTarget})`)
  await run([
    "bun",
    "build",
    "--compile",
    `--target=${target.bunTarget}`,
    "--minify",
    "--define",
    `__QUALITYMD_VERSION__=${JSON.stringify(version)}`,
    "--define",
    `__QUALITYMD_COMMIT__=${JSON.stringify(commit)}`,
    `--outfile=${binary}`,
    "src/main.ts",
  ])
  const executableBytes = (await stat(binary)).size
  if (executableBytes > maxExecutableBytes)
    throw new Error(
      `${target.asset} executable is ${executableBytes} bytes; budget is ${maxExecutableBytes}`,
    )
  if (target.os !== "win32") await chmod(binary, 0o755)
  const archive = join(dist, target.asset)
  if (target.asset.endsWith(".zip")) {
    await run(["zip", "-X", "-j", archive, binary])
  } else {
    await run(
      process.platform === "linux"
        ? [
            "tar",
            "--sort=name",
            "--mtime=@0",
            "--owner=0",
            "--group=0",
            "--numeric-owner",
            "-C",
            directory,
            "-czf",
            archive,
            "qualitymd",
          ]
        : ["tar", "-C", directory, "-czf", archive, "qualitymd"],
    )
  }
  const archiveBytes = (await stat(archive)).size
  if (archiveBytes > maxArchiveBytes)
    throw new Error(`${target.asset} is ${archiveBytes} bytes; budget is ${maxArchiveBytes}`)
}

const archives = (await readdir(dist)).filter((name) => /\.(?:tar\.gz|zip)$/.test(name)).sort()
const checksums: string[] = []
for (const name of archives) {
  const hash = new Bun.CryptoHasher("sha256")
    .update(await Bun.file(join(dist, name)).arrayBuffer())
    .digest("hex")
  checksums.push(`${hash}  ${name}`)
}
await writeFile(join(dist, "checksums.txt"), `${checksums.join("\n")}\n`)
await writeFile(
  join(dist, "metadata.json"),
  `${JSON.stringify(
    {
      schemaVersion: 1,
      version,
      commit,
      builder: `bun ${Bun.version}`,
      reproducibility:
        "Bun standalone executables are not byte-reproducible across independent compiles; release repair reuses checksum-verified published archives instead of rebuilding them.",
      budgets: {
        maxExecutableBytes,
        maxArchiveBytes,
      },
      targets: targets.map((target) => ({
        asset: target.asset,
        bunTarget: target.bunTarget,
        ...(target.libc === undefined ? {} : { libc: target.libc }),
      })),
    },
    null,
    2,
  )}\n`,
)
await writeFile(
  join(dist, "artifacts.json"),
  `${JSON.stringify(
    archives.map((name) => ({ name, path: join("dist", basename(name)) })),
    null,
    2,
  )}\n`,
)
await rm(staging, { recursive: true, force: true })
