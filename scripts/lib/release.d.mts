export interface ReleaseTarget {
  readonly os: "darwin" | "linux" | "win32"
  readonly arch: "arm64" | "x64"
  readonly npm: string
  readonly asset: string
  readonly bunTarget: string
  readonly libc?: "glibc" | "musl"
  readonly homebrew: boolean
}

export const root: string
export const targets: ReadonlyArray<ReleaseTarget>
export const versionFromTag: (tag: string) => string
