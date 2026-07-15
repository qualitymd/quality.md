import * as Effect from "effect/Effect"

import {
  buildCommit,
  buildVersion,
  isDevelopmentVersion,
  specificationVersion,
} from "../build-info.ts"
import { commandResult } from "../domain/command-result.ts"
import { jsonDocument } from "../domain/json.ts"

export interface VersionInput {
  readonly json: boolean
}

export interface VersionInfo {
  readonly schemaVersion: 1
  readonly version: string
  readonly commit?: string
  readonly developmentBuild: boolean
  readonly specificationVersion: string
}

export const currentVersionInfo = (): VersionInfo => ({
  schemaVersion: 1,
  version: buildVersion,
  ...(buildCommit === undefined ? {} : { commit: buildCommit }),
  developmentBuild: isDevelopmentVersion(buildVersion),
  specificationVersion,
})

const emptyDisplay = (value: string | undefined) =>
  value === undefined || value === "" ? "not recorded" : value

export const versionCommand = ({ json }: VersionInput) =>
  Effect.sync(() => {
    const info = currentVersionInfo()
    if (json) return commandResult(jsonDocument(info))
    return commandResult(
      `Version: ${info.version}\n` +
        `Commit: ${emptyDisplay(info.commit)}\n` +
        `Development build: ${String(info.developmentBuild)}\n` +
        `Specification version: ${emptyDisplay(info.specificationVersion)}\n`,
    )
  })
