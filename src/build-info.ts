import specification from "../SPECIFICATION.md" with { type: "text" }

declare const __QUALITYMD_VERSION__: string | undefined
declare const __QUALITYMD_COMMIT__: string | undefined

const definedVersion = typeof __QUALITYMD_VERSION__ === "string" ? __QUALITYMD_VERSION__.trim() : ""
const definedCommit = typeof __QUALITYMD_COMMIT__ === "string" ? __QUALITYMD_COMMIT__.trim() : ""

export const buildVersion = definedVersion === "" ? "dev" : definedVersion
export const buildCommit =
  definedCommit === "" || definedCommit === "none" ? undefined : definedCommit

export const specificationVersion = (() => {
  const marker = "**Specification version:**"
  const line = specification.split("\n").find((candidate) => candidate.trim().startsWith(marker))
  return line?.trim().slice(marker.length).trim() ?? ""
})()

export const isDevelopmentVersion = (version: string): boolean => {
  const trimmed = version.trim()
  if (trimmed.startsWith("dev") || trimmed.includes("+dirty")) return true
  const normalized = `v${trimmed.replace(/^v/, "")}`
  return normalized.includes("-0.")
}
