import * as Clock from "effect/Clock"
import * as Effect from "effect/Effect"
import * as FileSystem from "effect/FileSystem"
import * as Path from "effect/Path"

import { FileSystemFailure } from "../domain/errors.ts"
import { currentVersionInfo } from "./version.ts"
import { commandResult, ExitCode, type CommandResult } from "../domain/command-result.ts"
import { jsonDocument } from "../domain/json.ts"
import { UpdateRuntime, type UpdateRuntimeService } from "../services/update-runtime.ts"

type InstallMethod = "managed-standalone" | "npm" | "homebrew" | "archive" | "unknown"

interface LatestRelease {
  readonly version: string
  readonly ready: boolean
  readonly releaseNotesURL?: string
}

export interface UpdateResult {
  readonly schemaVersion: 1
  readonly currentVersion: string
  readonly commit?: string
  readonly specificationVersion: string
  readonly developmentBuild: boolean
  readonly installMethod: InstallMethod
  readonly latestVersion?: string
  readonly latestVersionReady: boolean
  readonly updateAvailable: boolean
  readonly applySupported: boolean
  readonly recommendedAction?: string
  readonly recommendedCommand?: string
  readonly applied: boolean
  readonly releaseNotesURL?: string
}

interface UpdateCache {
  readonly latestVersion?: string
  readonly releaseNotesURL?: string
  readonly ready: boolean
  readonly checkedAt: string
}

export interface UpdateInput {
  readonly check: boolean
  readonly json: boolean
}

const truthy = (value: string | undefined) =>
  ["1", "true", "yes", "on"].includes(value?.trim().toLowerCase() ?? "")

const qualitymdHome = (environment: Readonly<Record<string, string | undefined>>) =>
  environment.QUALITYMD_HOME || `${environment.HOME || environment.USERPROFILE || "."}/.qualitymd`

const cacheName = ".qualitymd-update-cache"
const managedMarker = ".qualitymd-managed-install"
const cacheTtlMs = 20 * 60 * 60 * 1000

const normalizeVersion = (input: string) => {
  let value = input
    .trim()
    .replace(/^qualitymd\s+/, "")
    .replace(/^v/, "")
  const suffix = value.search(/[ +(]/)
  if (suffix >= 0) value = value.slice(0, suffix)
  return value
}

const semver = (input: string) => {
  const normalized = normalizeVersion(input)
  const match = /^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-([^+]+))?(?:\+.+)?$/.exec(normalized)
  if (match === null) return undefined
  return {
    normalized,
    numbers: [Number(match[1]), Number(match[2]), Number(match[3])] as const,
    prerelease: match[4] ?? "",
  }
}

export const updateAvailable = (current: string, latest: string) => {
  const left = semver(current)
  const right = semver(latest)
  if (
    left === undefined ||
    right === undefined ||
    left.prerelease !== "" ||
    right.prerelease !== ""
  )
    return false
  for (let index = 0; index < 3; index += 1) {
    if (right.numbers[index]! > left.numbers[index]!) return true
    if (right.numbers[index]! < left.numbers[index]!) return false
  }
  return false
}

const sameRelease = (left: string, right: string) => {
  const a = semver(left)
  const b = semver(right)
  return a !== undefined && b !== undefined && a.normalized === b.normalized
}

export const platformArchiveName = (
  platform: NodeJS.Platform,
  hostArch: string,
  linuxIsMusl: boolean,
) => {
  const arch = hostArch === "x64" ? "amd64" : hostArch
  if (arch !== "amd64" && arch !== "arm64") return ""
  if (platform === "win32") return `qualitymd_windows_${arch}.zip`
  if (platform === "darwin") return `qualitymd_darwin_${arch}.tar.gz`
  if (platform === "linux") return `qualitymd_linux_${arch}${linuxIsMusl ? "_musl" : ""}.tar.gz`
  return ""
}

const latestVersion = (
  runtime: UpdateRuntimeService,
  method: InstallMethod,
): Effect.Effect<LatestRelease, Error> => {
  if (method === "npm") {
    return runtime.request("https://registry.npmjs.org/quality.md/latest", "application/json").pipe(
      Effect.map((body) => {
        const payload = JSON.parse(body) as { readonly version?: string }
        return payload.version === undefined
          ? { version: "", ready: false }
          : { version: `v${payload.version.replace(/^v/, "")}`, ready: true }
      }),
    )
  }
  if (method === "homebrew") {
    return runtime
      .request(
        "https://raw.githubusercontent.com/qualitymd/homebrew-tap/main/Casks/qualitymd.rb",
        "text/plain",
      )
      .pipe(
        Effect.map((body) => {
          const version = /^\s*version\s+"([^"]+)"/m.exec(body)?.[1]
          if (version === undefined) return { version: "", ready: false }
          const tag = `v${version.replace(/^v/, "")}`
          return {
            version: tag,
            ready: true,
            releaseNotesURL: `https://github.com/qualitymd/quality.md/releases/tag/${tag}`,
          }
        }),
      )
  }
  return runtime
    .request(
      "https://api.github.com/repos/qualitymd/quality.md/releases/latest",
      "application/json",
    )
    .pipe(
      Effect.map((body) => {
        const payload = JSON.parse(body) as {
          readonly tag_name?: string
          readonly html_url?: string
          readonly assets?: ReadonlyArray<{ readonly name?: string }>
        }
        const names = new Set((payload.assets ?? []).map((asset) => asset.name))
        const archive = platformArchiveName(runtime.platform, runtime.arch, runtime.linuxIsMusl)
        return {
          version: payload.tag_name ?? "",
          ready: archive !== "" && names.has(archive) && names.has("checksums.txt"),
          ...(payload.html_url === undefined ? {} : { releaseNotesURL: payload.html_url }),
        }
      }),
    )
}

const commandFor = (
  method: InstallMethod,
  platform: NodeJS.Platform,
): { readonly display: string; readonly argv: string[] } | undefined => {
  if (method === "npm")
    return {
      display: "npm install -g quality.md@latest",
      argv: ["npm", "install", "-g", "quality.md@latest"],
    }
  if (method === "homebrew")
    return {
      display: "brew upgrade qualitymd/tap/qualitymd",
      argv: ["brew", "upgrade", "qualitymd/tap/qualitymd"],
    }
  if (method === "managed-standalone") {
    if (platform === "win32") {
      const script =
        "$env:QUALITYMD_NO_INPUT='1'; iwr https://getquality.md/install.ps1 -UseB | iex"
      return {
        display: `powershell -NoProfile -ExecutionPolicy Bypass -Command ${JSON.stringify(script)}`,
        argv: ["powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", script],
      }
    }
    const script = "curl -fsSL https://getquality.md/install.sh | QUALITYMD_NO_INPUT=1 sh"
    return { display: script, argv: ["sh", "-c", script] }
  }
  return undefined
}

const recommendedAction = (method: InstallMethod, command?: string) => {
  if (command !== undefined) return command
  if (method === "archive")
    return "download and replace the archive-installed binary from the latest GitHub release"
  return "install or upgrade with a managed channel such as npm, Homebrew, or the hosted installer"
}

const detectInstallMethod = (
  fs: FileSystem.FileSystem,
  paths: Path.Path,
  runtime: UpdateRuntimeService,
) =>
  Effect.gen(function* () {
    const explicit = runtime.environment.QUALITYMD_INSTALL_METHOD?.toLowerCase()
    if (explicit === "npm" || explicit === "homebrew" || explicit === "managed-standalone")
      return explicit
    const executable = paths.resolve(runtime.execPath)
    const root = paths.resolve(qualitymdHome(runtime.environment))
    if (yield* fs.exists(paths.join(root, managedMarker))) {
      const relative = paths.relative(root, executable)
      if (relative !== ".." && !relative.startsWith(`..${paths.sep}`)) return "managed-standalone"
    }
    const lower = executable.replaceAll("\\", "/").toLowerCase()
    if (
      lower.includes("/cellar/qualitymd/") ||
      lower.includes("/caskroom/qualitymd/") ||
      (lower.includes("/homebrew/") && lower.includes("/qualitymd"))
    )
      return "homebrew"
    return /\/qualitymd(?:\.exe)?$/.test(lower) ? "archive" : "unknown"
  })

const writeCache = (
  fs: FileSystem.FileSystem,
  paths: Path.Path,
  runtime: UpdateRuntimeService,
  release: LatestRelease,
) =>
  Effect.gen(function* () {
    if (release.version === "") return
    const root = paths.resolve(qualitymdHome(runtime.environment))
    const now = yield* Clock.currentTimeMillis
    yield* fs.makeDirectory(root, { recursive: true, mode: 0o700 })
    const record: UpdateCache = {
      latestVersion: release.version,
      ...(release.releaseNotesURL === undefined
        ? {}
        : { releaseNotesURL: release.releaseNotesURL }),
      ready: release.ready,
      checkedAt: new Date(now).toISOString(),
    }
    yield* fs.writeFileString(paths.join(root, cacheName), jsonDocument(record), { mode: 0o600 })
  })

const renderHuman = (result: UpdateResult, check: boolean) => {
  let output =
    `Current version: ${result.currentVersion}\n` +
    `Latest version: ${result.latestVersion || "not recorded"}\n` +
    `Latest version ready: ${result.latestVersionReady}\n` +
    `Install method: ${result.installMethod}\n` +
    `Update available: ${result.updateAvailable}\n`
  if (result.releaseNotesURL) output += `Release notes: ${result.releaseNotesURL}\n`
  if (result.applied) return `${output}Update applied.\n`
  if (!check && !result.updateAvailable) return `${output}Already up to date.\n`
  if (result.applySupported) return `${output}Recommended command: ${result.recommendedCommand}\n`
  return `${output}Recommended action: ${result.recommendedAction}\n`
}

export const updateCommand = (
  input: UpdateInput,
): Effect.Effect<
  CommandResult,
  FileSystemFailure,
  FileSystem.FileSystem | Path.Path | UpdateRuntime
> =>
  Effect.gen(function* () {
    const fs = yield* FileSystem.FileSystem
    const paths = yield* Path.Path
    const runtime = yield* UpdateRuntime
    const info = currentVersionInfo()
    const installMethod = yield* detectInstallMethod(fs, paths, runtime)
    const command = commandFor(installMethod, runtime.platform)
    const disabled = truthy(runtime.environment.QUALITYMD_NO_UPDATE_CHECK)
    const latest = disabled
      ? ({ version: "", ready: false } satisfies LatestRelease)
      : yield* latestVersion(runtime, installMethod)
    const newer = updateAvailable(info.version, latest.version)
    const available = !info.developmentBuild && newer && latest.ready
    let result: UpdateResult = {
      schemaVersion: 1,
      currentVersion: info.version,
      ...(info.commit === undefined ? {} : { commit: info.commit }),
      specificationVersion: info.specificationVersion,
      developmentBuild: info.developmentBuild,
      installMethod,
      ...(latest.version === "" ? {} : { latestVersion: latest.version }),
      latestVersionReady: latest.ready,
      updateAvailable: available,
      applySupported: command !== undefined,
      recommendedAction: recommendedAction(installMethod, command?.display),
      ...(command === undefined ? {} : { recommendedCommand: command.display }),
      applied: false,
      ...(latest.releaseNotesURL === undefined ? {} : { releaseNotesURL: latest.releaseNotesURL }),
    }
    if (!input.check && !info.developmentBuild && newer && !latest.ready)
      return commandResult("", {
        stderr: `qualitymd: latest qualitymd release ${latest.version} is not yet available for ${installMethod} installs\n`,
        exitCode: ExitCode.problems,
      })
    if (!input.check && available) {
      if (command === undefined)
        return commandResult("", {
          stderr: `qualitymd: update apply is not supported for ${installMethod} installs; ${result.recommendedAction}\n`,
          exitCode: ExitCode.problems,
        })
      const applied = yield* runtime.run(command.argv).pipe(Effect.result)
      if (applied._tag === "Failure")
        return commandResult("", {
          stderr: `qualitymd: ${applied.failure.message}\n`,
          exitCode: ExitCode.problems,
        })
      const verified = yield* runtime.visibleVersion.pipe(Effect.result)
      if (verified._tag === "Failure")
        return commandResult("", {
          stderr: `qualitymd: update command ran, but qualitymd --version could not be verified: ${verified.failure.message}\n`,
          exitCode: ExitCode.problems,
        })
      if (!sameRelease(verified.success, latest.version))
        return commandResult("", {
          stderr: `qualitymd: update command ran, but visible qualitymd version is ${verified.success || "not recorded"}, not ${latest.version}\n`,
          exitCode: ExitCode.problems,
        })
      const { recommendedAction: _, recommendedCommand: __, ...withoutRecommendation } = result
      result = {
        ...withoutRecommendation,
        currentVersion: verified.success,
        updateAvailable: false,
        applied: true,
      }
    }
    if (!info.developmentBuild && latest.version !== "")
      yield* writeCache(fs, paths, runtime, latest).pipe(Effect.ignore)
    return commandResult(input.json ? jsonDocument(result) : renderHuman(result, input.check))
  }).pipe(
    Effect.mapError(
      (cause) =>
        new FileSystemFailure({ detail: cause instanceof Error ? cause.message : String(cause) }),
    ),
  )

export const updateRefreshCommand = (): Effect.Effect<
  CommandResult,
  never,
  FileSystem.FileSystem | Path.Path | UpdateRuntime
> =>
  Effect.gen(function* () {
    const runtime = yield* UpdateRuntime
    if (
      truthy(runtime.environment.QUALITYMD_NO_UPDATE_CHECK) ||
      currentVersionInfo().developmentBuild
    )
      return commandResult()
    const fs = yield* FileSystem.FileSystem
    const paths = yield* Path.Path
    const method = yield* detectInstallMethod(fs, paths, runtime)
    const latest = yield* latestVersion(runtime, method).pipe(Effect.result)
    if (latest._tag === "Success")
      yield* writeCache(fs, paths, runtime, latest.success).pipe(Effect.ignore)
    return commandResult()
  }).pipe(Effect.catch(() => Effect.succeed(commandResult())))

export const cachedUpdateNotice = (
  args: ReadonlyArray<string>,
): Effect.Effect<string, never, FileSystem.FileSystem | Path.Path | UpdateRuntime> =>
  Effect.gen(function* () {
    const runtime = yield* UpdateRuntime
    const info = currentVersionInfo()
    if (
      info.developmentBuild ||
      truthy(runtime.environment.QUALITYMD_NO_UPDATE_CHECK) ||
      runtime.environment.CI !== undefined ||
      !runtime.stderrIsTTY ||
      args.includes("--json") ||
      args.includes("update") ||
      runtime.environment.QUALITYMD_UPDATE_REFRESH === "1"
    )
      return ""
    const fs = yield* FileSystem.FileSystem
    const paths = yield* Path.Path
    const path = paths.join(paths.resolve(qualitymdHome(runtime.environment)), cacheName)
    let stale = true
    if (yield* fs.exists(path)) {
      const raw = yield* fs.readFileString(path)
      const parsed = yield* Effect.try(() => JSON.parse(raw)).pipe(Effect.result)
      if (parsed._tag === "Success") {
        const record = parsed.success as UpdateCache
        const now = yield* Clock.currentTimeMillis
        stale = now - Date.parse(record.checkedAt) >= cacheTtlMs
        if (record.ready && updateAvailable(info.version, record.latestVersion ?? ""))
          return `update available: ${info.version} -> ${record.latestVersion} (run \`qualitymd update\`)${record.releaseNotesURL ? ` ${record.releaseNotesURL}` : ""}\n`
      }
    }
    if (stale) {
      yield* runtime.spawnRefresh({
        ...Object.fromEntries(
          Object.entries(runtime.environment).filter(
            (entry): entry is [string, string] => entry[1] !== undefined,
          ),
        ),
        QUALITYMD_UPDATE_REFRESH: "1",
      })
    }
    return ""
  }).pipe(Effect.catch(() => Effect.succeed("")))
