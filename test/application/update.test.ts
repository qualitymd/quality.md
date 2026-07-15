import { assert, describe, it } from "@effect/vitest"
import * as BunFileSystem from "@effect/platform-bun/BunFileSystem"
import * as BunPath from "@effect/platform-bun/BunPath"
import * as Effect from "effect/Effect"
import * as Layer from "effect/Layer"

import {
  platformArchiveName,
  updateAvailable,
  updateCommand,
} from "../../src/application/update.ts"
import { UpdateRuntime, type UpdateRuntimeService } from "../../src/services/update-runtime.ts"

describe("updates", () => {
  it("compares only stable semantic versions", () => {
    assert.strictEqual(updateAvailable("v1.2.3", "v1.2.4"), true)
    assert.strictEqual(updateAvailable("1.3.0", "1.2.9"), false)
    assert.strictEqual(updateAvailable("1.2.3-dev", "1.2.4"), false)
    assert.strictEqual(updateAvailable("not-a-version", "1.2.4"), false)
  })

  it("selects release archives by platform, architecture, and libc", () => {
    assert.strictEqual(
      platformArchiveName("darwin", "arm64", false),
      "qualitymd_darwin_arm64.tar.gz",
    )
    assert.strictEqual(platformArchiveName("win32", "x64", false), "qualitymd_windows_amd64.zip")
    assert.strictEqual(platformArchiveName("linux", "x64", false), "qualitymd_linux_amd64.tar.gz")
    assert.strictEqual(
      platformArchiveName("linux", "x64", true),
      "qualitymd_linux_amd64_musl.tar.gz",
    )
    assert.strictEqual(platformArchiveName("freebsd", "x64", false), "")
  })

  it.effect("uses the injected runtime for channel readiness", () => {
    const runtime: UpdateRuntimeService = {
      platform: "linux",
      arch: "x64",
      execPath: "/tmp/qualitymd",
      environment: {
        QUALITYMD_INSTALL_METHOD: "npm",
        QUALITYMD_HOME: "/tmp/qualitymd-update-test",
      },
      stderrIsTTY: false,
      linuxIsMusl: false,
      request: () => Effect.succeed('{"version":"0.31.0"}'),
      run: () => Effect.void,
      visibleVersion: Effect.succeed("v0.31.0"),
      spawnRefresh: () => Effect.void,
    }
    const services = Layer.mergeAll(
      BunFileSystem.layer,
      BunPath.layer,
      Layer.succeed(UpdateRuntime, runtime),
    )
    return updateCommand({ check: true, json: true }).pipe(
      Effect.provide(services),
      Effect.map((result) => {
        const payload = JSON.parse(result.stdout) as {
          readonly installMethod: string
          readonly latestVersion: string
          readonly latestVersionReady: boolean
        }
        assert.strictEqual(payload.installMethod, "npm")
        assert.strictEqual(payload.latestVersion, "v0.31.0")
        assert.strictEqual(payload.latestVersionReady, true)
      }),
    )
  })
})
