export const ExitCode = {
  ok: 0,
  problems: 1,
  usage: 2,
  internal: 70,
} as const

export type ExitCode = (typeof ExitCode)[keyof typeof ExitCode]

export interface CommandResult {
  readonly exitCode: ExitCode
  readonly stdout: string
  readonly stderr: string
}

export const commandResult = (
  stdout = "",
  options: {
    readonly stderr?: string
    readonly exitCode?: ExitCode
  } = {},
): CommandResult => ({
  exitCode: options.exitCode ?? ExitCode.ok,
  stdout,
  stderr: options.stderr ?? "",
})
