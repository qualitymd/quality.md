import * as Cause from "effect/Cause"
import * as Effect from "effect/Effect"
import type * as FileSystem from "effect/FileSystem"
import * as Option from "effect/Option"
import type * as Path from "effect/Path"
import { Argument, Command, Flag } from "effect/unstable/cli"

import { lintCommand } from "../application/lint.ts"
import { initCommand } from "../application/init.ts"
import { evaluationCreateCommand } from "../application/evaluation-create.ts"
import {
  evaluationListCommand,
  evaluationStatusCommand,
} from "../application/evaluation-inspect.ts"
import { evaluationRunCommand } from "../application/evaluation-run.ts"
import { evaluationReportBuildCommand } from "../application/evaluation-report.ts"
import {
  evaluationDataExampleCommand,
  evaluationDataGetCommand,
  evaluationDataKindsCommand,
  evaluationDataListCommand,
  evaluationDataSchemaCommand,
  evaluationDataSetCommand,
  evaluationDataVerifyCommand,
} from "../application/evaluation-data.ts"
import { modelGetCommand, modelListCommand, modelTreeCommand } from "../application/model.ts"
import type { FileSystemFailure } from "../domain/errors.ts"
import { schemaCommand, specCommand } from "../application/static-artifacts.ts"
import { statusCommand } from "../application/status.ts"
import { versionCommand } from "../application/version.ts"
import { cachedUpdateNotice, updateCommand, updateRefreshCommand } from "../application/update.ts"
import { buildCommit, buildVersion } from "../build-info.ts"
import { commandResult, ExitCode, type CommandResult } from "../domain/command-result.ts"
import type { HostRuntime } from "../services/host-runtime.ts"
import { Output } from "../services/output.ts"

export interface CliCommand<Input, Error = never, Requirements = never> {
  readonly name: string
  readonly description: string
  readonly run: (input: Input) => Effect.Effect<CommandResult, Error, Requirements>
}

const emit = (result: CommandResult) =>
  Effect.gen(function* () {
    const output = yield* Output
    if (result.stdout !== "") yield* output.stdout(result.stdout)
    if (result.stderr !== "") yield* output.stderr(result.stderr)
    yield* output.setExitCode(result.exitCode)
  })

const execute = <Input, Error, Requirements>(
  descriptor: CliCommand<Input, Error, Requirements>,
  input: Input,
) =>
  descriptor.run(input).pipe(
    Effect.catchCause((cause) => {
      const squashed = Cause.squash(cause)
      return Effect.succeed(
        commandResult("", {
          stderr: `qualitymd: ${squashed instanceof Error ? squashed.message : String(squashed)}\n`,
          exitCode: ExitCode.internal,
        }),
      )
    }),
    Effect.flatMap(emit),
  )

const lintDescriptor: CliCommand<
  { readonly path: string; readonly json: boolean; readonly fix: boolean },
  import("../domain/errors.ts").FileSystemFailure,
  FileSystem.FileSystem | Path.Path
> = {
  name: "lint",
  description: "Validate a QUALITY.md file",
  run: lintCommand,
}

const lint = Command.make(
  lintDescriptor.name,
  {
    path: Argument.string("path").pipe(Argument.withDefault("QUALITY.md")),
    json: Flag.boolean("json").pipe(
      Flag.withDescription("Emit a machine-readable JSON lint result"),
    ),
    fix: Flag.boolean("fix").pipe(Flag.withDescription("Apply deterministic in-place repairs")),
  },
  (input) => execute(lintDescriptor, input),
).pipe(Command.withDescription(lintDescriptor.description))

const initDescriptor: CliCommand<
  {
    readonly path: string
    readonly force: boolean
    readonly json: boolean
    readonly minimal: boolean
    readonly noAgentInstructions: boolean
  },
  FileSystemFailure,
  FileSystem.FileSystem | Path.Path | HostRuntime
> = {
  name: "init",
  description: "Scaffold a starter QUALITY.md",
  run: initCommand,
}

const init = Command.make(
  initDescriptor.name,
  {
    path: Argument.string("path").pipe(Argument.withDefault("QUALITY.md")),
    force: Flag.boolean("force").pipe(Flag.withDescription("Overwrite an existing file")),
    json: Flag.boolean("json").pipe(
      Flag.withDescription("Emit a machine-readable JSON init receipt"),
    ),
    minimal: Flag.boolean("minimal").pipe(
      Flag.withDescription("Write a minimal valid skeleton without the guided template prose"),
    ),
    noAgentInstructions: Flag.boolean("no-agent-instructions").pipe(
      Flag.withDescription("Do not create or update agent instruction files"),
    ),
  },
  (input) => execute(initDescriptor, input),
).pipe(Command.withDescription(initDescriptor.description))

const modelTreeDescriptor: CliCommand<
  {
    readonly path: string
    readonly json: boolean
    readonly area: string
    readonly depth?: number
  },
  FileSystemFailure,
  FileSystem.FileSystem
> = {
  name: "tree",
  description: "Render the model as a containment hierarchy",
  run: modelTreeCommand,
}

const modelTree = Command.make(
  modelTreeDescriptor.name,
  {
    path: Argument.string("path").pipe(Argument.withDefault("QUALITY.md")),
    json: Flag.boolean("json").pipe(Flag.withDescription("Emit the tree as nested JSON")),
    area: Flag.string("area").pipe(
      Flag.withDefault(""),
      Flag.withDescription("Root the tree at a canonical area reference (area:<path>)"),
    ),
    depth: Flag.integer("depth").pipe(
      Flag.optional,
      Flag.withDescription("Limit nesting depth; 0 emits only the rooted node"),
    ),
  },
  ({ depth, ...input }) =>
    execute(modelTreeDescriptor, {
      ...input,
      ...(Option.isSome(depth) ? { depth: depth.value } : {}),
    }),
).pipe(Command.withDescription(modelTreeDescriptor.description))

const modelListDescriptor: CliCommand<
  {
    readonly path: string
    readonly json: boolean
    readonly area: string
    readonly types: ReadonlyArray<string>
  },
  FileSystemFailure,
  FileSystem.FileSystem
> = {
  name: "list",
  description: "Enumerate model elements with their canonical IDs",
  run: modelListCommand,
}

const modelList = Command.make(
  modelListDescriptor.name,
  {
    path: Argument.string("path").pipe(Argument.withDefault("QUALITY.md")),
    json: Flag.boolean("json").pipe(Flag.withDescription("Emit the enumeration as a JSON array")),
    area: Flag.string("area").pipe(
      Flag.withDefault(""),
      Flag.withDescription("Restrict output to one area's subtree (area:<path>)"),
    ),
    types: Flag.string("type").pipe(
      Flag.atMost(100),
      Flag.withDescription("Restrict output to kinds: area, factor, requirement"),
    ),
  },
  (input) => execute(modelListDescriptor, input),
).pipe(Command.withDescription(modelListDescriptor.description))

const modelGetDescriptor: CliCommand<
  { readonly id: string; readonly path: string; readonly json: boolean },
  FileSystemFailure,
  FileSystem.FileSystem
> = {
  name: "get",
  description: "Show one element's detail and immediate relations",
  run: modelGetCommand,
}

const modelGet = Command.make(
  modelGetDescriptor.name,
  {
    args: Argument.string("id [path]").pipe(Argument.between(1, 2)),
    json: Flag.boolean("json").pipe(
      Flag.withDescription("Emit the element detail as a JSON object"),
    ),
  },
  ({ args, json }) =>
    execute(modelGetDescriptor, { id: args[0]!, path: args[1] ?? "QUALITY.md", json }),
).pipe(Command.withDescription(modelGetDescriptor.description))

const model = Command.make("model", {}, () => Effect.void).pipe(
  Command.withDescription("Query a quality model's structure and canonical reference IDs"),
  Command.withSubcommands([modelTree, modelList, modelGet]),
)

const statusDescriptor: CliCommand<
  { readonly path: string; readonly json: boolean },
  FileSystemFailure,
  FileSystem.FileSystem | Path.Path
> = {
  name: "status",
  description: "Show a QUALITY.md workspace status snapshot",
  run: statusCommand,
}

const status = Command.make(
  statusDescriptor.name,
  {
    path: Argument.string("path").pipe(Argument.withDefault("QUALITY.md")),
    json: Flag.boolean("json").pipe(
      Flag.withDescription("Emit a machine-readable workspace status snapshot"),
    ),
  },
  (input) => execute(statusDescriptor, input),
).pipe(Command.withDescription(statusDescriptor.description))

const evaluationCreateDescriptor: CliCommand<
  {
    readonly modelArgument?: string
    readonly model: string
    readonly area: string
    readonly factors: ReadonlyArray<string>
    readonly evaluationDir: string
    readonly json: boolean
  },
  FileSystemFailure,
  FileSystem.FileSystem | Path.Path | HostRuntime
> = {
  name: "create",
  description: "Create a numbered evaluation run folder",
  run: evaluationCreateCommand,
}

const evaluationCreate = Command.make(
  evaluationCreateDescriptor.name,
  {
    modelArgument: Argument.string("model").pipe(Argument.optional),
    model: Flag.string("model").pipe(
      Flag.withDefault(""),
      Flag.withDescription("QUALITY.md file to snapshot"),
    ),
    area: Flag.string("area").pipe(
      Flag.withDefault(""),
      Flag.withDescription("Canonical area reference for the evaluation scope"),
    ),
    factors: Flag.string("factor").pipe(
      Flag.atMost(100),
      Flag.withDescription("Canonical factor reference for a scoped evaluation; repeatable"),
    ),
    evaluationDir: Flag.string("evaluation-dir").pipe(
      Flag.withDefault(""),
      Flag.withDescription("Override the model-relative evaluation directory"),
    ),
    json: Flag.boolean("json").pipe(
      Flag.withDescription("Emit a machine-readable run creation receipt"),
    ),
  },
  ({ modelArgument, ...input }) =>
    execute(evaluationCreateDescriptor, {
      ...input,
      ...(Option.isSome(modelArgument) ? { modelArgument: modelArgument.value } : {}),
    }),
).pipe(Command.withDescription(evaluationCreateDescriptor.description))

const runFlags = () => ({
  run: Argument.string("run").pipe(Argument.optional),
  latest: Flag.boolean("latest").pipe(Flag.withDescription("Use the most recent evaluation run")),
  evaluationDir: Flag.string("evaluation-dir").pipe(
    Flag.withDefault(""),
    Flag.withDescription("Override the model-relative evaluation directory when using --latest"),
  ),
  model: Flag.string("model").pipe(
    Flag.withDefault(""),
    Flag.withDescription("QUALITY.md file that anchors model-relative run paths"),
  ),
})

const evaluationStatusDescriptor: CliCommand<
  {
    readonly run?: string
    readonly latest: boolean
    readonly model: string
    readonly evaluationDir: string
    readonly json: boolean
  },
  FileSystemFailure,
  FileSystem.FileSystem | Path.Path
> = {
  name: "status",
  description: "Show whether an evaluation run is reportable",
  run: evaluationStatusCommand,
}

const evaluationStatus = Command.make(
  evaluationStatusDescriptor.name,
  {
    ...runFlags(),
    json: Flag.boolean("json").pipe(
      Flag.withDescription("Emit a machine-readable status document"),
    ),
  },
  ({ run, ...input }) =>
    execute(evaluationStatusDescriptor, {
      ...input,
      ...(Option.isSome(run) ? { run: run.value } : {}),
    }),
).pipe(Command.withDescription(evaluationStatusDescriptor.description))

const evaluationListDescriptor: CliCommand<
  {
    readonly model: string
    readonly evaluationDir: string
    readonly state: string
    readonly json: boolean
  },
  FileSystemFailure,
  FileSystem.FileSystem | Path.Path
> = {
  name: "list",
  description: "List evaluation runs",
  run: evaluationListCommand,
}

const evaluationList = Command.make(
  evaluationListDescriptor.name,
  {
    model: Flag.string("model").pipe(
      Flag.withDefault(""),
      Flag.withDescription("QUALITY.md file that anchors evaluation history"),
    ),
    evaluationDir: Flag.string("evaluation-dir").pipe(
      Flag.withDefault(""),
      Flag.withDescription("Override the model-relative evaluation directory"),
    ),
    state: Flag.string("state").pipe(
      Flag.withDefault("all"),
      Flag.withDescription("Filter runs: all, reportable, incomplete, awaiting"),
    ),
    json: Flag.boolean("json").pipe(Flag.withDescription("Emit a machine-readable run list")),
  },
  (input) => execute(evaluationListDescriptor, input),
).pipe(Command.withDescription(evaluationListDescriptor.description))

const evaluationRunDescriptor: CliCommand<
  {
    readonly model: string
    readonly evaluationDir: string
    readonly area: string
    readonly factors: ReadonlyArray<string>
    readonly evaluator: string
    readonly resume: string
    readonly evaluatorResult: string
    readonly dryRun: boolean
    readonly json: boolean
  },
  FileSystemFailure,
  FileSystem.FileSystem | Path.Path | HostRuntime
> = {
  name: "run",
  description: "Execute a complete evaluation run with the deterministic runner",
  run: evaluationRunCommand,
}

const evaluationRun = Command.make(
  evaluationRunDescriptor.name,
  {
    model: Flag.string("model").pipe(
      Flag.withDefault(""),
      Flag.withDescription("QUALITY.md file to evaluate"),
    ),
    evaluationDir: Flag.string("evaluation-dir").pipe(
      Flag.withDefault(""),
      Flag.withDescription("Override the model-relative evaluation directory"),
    ),
    area: Flag.string("area").pipe(
      Flag.withDefault(""),
      Flag.withDescription("Canonical area reference for the evaluation scope"),
    ),
    factors: Flag.string("factor").pipe(
      Flag.atMost(100),
      Flag.withDescription("Canonical factor reference for a scoped evaluation; repeatable"),
    ),
    evaluator: Flag.string("evaluator").pipe(
      Flag.withDefault(""),
      Flag.withDescription(
        "Evaluator to use: auto (default), a built-in name, or a configured profile",
      ),
    ),
    resume: Flag.string("resume").pipe(
      Flag.withDefault(""),
      Flag.withDescription("Resume an existing run from its evaluation.json"),
    ),
    evaluatorResult: Flag.string("evaluator-result").pipe(
      Flag.withDefault(""),
      Flag.withDescription("Submit harness result envelopes for outstanding work requests"),
    ),
    dryRun: Flag.boolean("dry-run").pipe(
      Flag.withAlias("n"),
      Flag.withDescription(
        "Preview the resolved run without invoking an evaluator or writing evaluation data",
      ),
    ),
    json: Flag.boolean("json").pipe(Flag.withDescription("Emit a machine-readable run receipt")),
  },
  (input) => execute(evaluationRunDescriptor, input),
).pipe(Command.withDescription(evaluationRunDescriptor.description))

const dataRunInput = <A extends { readonly run: Option.Option<string> }>(input: A) => {
  const { run, ...rest } = input
  return { ...rest, ...(Option.isSome(run) ? { run: run.value } : {}) }
}

const evaluationDataSet = Command.make(
  "set",
  {
    ...runFlags(),
    dryRun: Flag.boolean("dry-run").pipe(
      Flag.withAlias("n"),
      Flag.withDescription("Validate and report intended write without persisting"),
    ),
    json: Flag.boolean("json").pipe(Flag.withDescription("Emit a machine-readable write receipt")),
  },
  (input) =>
    execute(
      {
        name: "set",
        description: "Validate and persist a batch of evaluation JSON payloads",
        run: evaluationDataSetCommand,
      },
      dataRunInput(input),
    ),
).pipe(Command.withDescription("Validate and persist a batch of evaluation JSON payloads"))

const evaluationDataList = Command.make(
  "list",
  {
    ...runFlags(),
    kind: Flag.string("kind").pipe(
      Flag.withDefault(""),
      Flag.withDescription("Filter by evaluation data kind"),
    ),
    json: Flag.boolean("json").pipe(Flag.withDescription("Emit a machine-readable data list")),
  },
  (input) =>
    execute(
      {
        name: "list",
        description: "List stored evaluation JSON payloads",
        run: evaluationDataListCommand,
      },
      dataRunInput(input),
    ),
).pipe(Command.withDescription("List stored evaluation JSON payloads"))

const evaluationDataGet = Command.make(
  "get",
  {
    ...runFlags(),
    kind: Flag.string("kind").pipe(
      Flag.withDefault(""),
      Flag.withDescription("Evaluation data kind"),
    ),
    area: Flag.string("area").pipe(Flag.withDefault(""), Flag.withDescription("Area ref")),
    factor: Flag.string("factor").pipe(Flag.withDefault(""), Flag.withDescription("Factor ref")),
    requirement: Flag.string("requirement").pipe(
      Flag.withDefault(""),
      Flag.withDescription("Requirement ref"),
    ),
    selector: Flag.string("selector").pipe(
      Flag.withDefault(""),
      Flag.withDescription("Optional sub-result selector"),
    ),
    json: Flag.boolean("json").pipe(
      Flag.withDescription("Not supported: data get already emits JSON"),
    ),
  },
  (input) =>
    execute(
      {
        name: "get",
        description: "Print one stored evaluation JSON payload",
        run: evaluationDataGetCommand,
      },
      dataRunInput(input),
    ),
).pipe(Command.withDescription("Print one stored evaluation JSON payload"))

const evaluationDataKinds = Command.make(
  "kinds",
  {
    json: Flag.boolean("json").pipe(Flag.withDescription("Emit a machine-readable data kind list")),
  },
  ({ json }) => emit(evaluationDataKindsCommand(json)),
).pipe(Command.withDescription("List evaluation data kinds"))

const evaluationDataExample = Command.make(
  "example",
  {
    kind: Argument.string("kind"),
    json: Flag.boolean("json").pipe(
      Flag.withDescription("Not supported: data example already emits JSON"),
    ),
  },
  ({ kind, json }) => emit(evaluationDataExampleCommand(kind, json)),
).pipe(Command.withDescription("Print a complete evaluation example JSON payload"))

const evaluationDataSchema = Command.make(
  "schema",
  {
    kind: Argument.string("kind").pipe(Argument.optional),
    json: Flag.boolean("json").pipe(
      Flag.withDescription("Not supported: data schema already emits JSON"),
    ),
  },
  ({ kind, json }) =>
    emit(evaluationDataSchemaCommand(Option.isSome(kind) ? kind.value : "", json)),
).pipe(Command.withDescription("Print the evaluation structured data JSON Schema"))

const evaluationDataVerify = Command.make(
  "verify",
  {
    ...runFlags(),
    json: Flag.boolean("json").pipe(
      Flag.withDescription("Emit a machine-readable verification receipt"),
    ),
  },
  (input) =>
    execute(
      {
        name: "verify",
        description: "Validate persisted evaluation JSON payloads",
        run: evaluationDataVerifyCommand,
      },
      dataRunInput(input),
    ),
).pipe(Command.withDescription("Validate persisted evaluation JSON payloads"))

const evaluationData = Command.make("data", {}, () => Effect.void).pipe(
  Command.withDescription("Work with evaluation structured data"),
  Command.withSubcommands([
    evaluationDataSet,
    evaluationDataList,
    evaluationDataGet,
    evaluationDataKinds,
    evaluationDataExample,
    evaluationDataSchema,
    evaluationDataVerify,
  ]),
)

const evaluationReportBuild = Command.make(
  "build",
  {
    ...runFlags(),
    json: Flag.boolean("json").pipe(
      Flag.withDescription("Emit a machine-readable report build receipt"),
    ),
  },
  (input) =>
    execute(
      {
        name: "build",
        description: "Build the complete Markdown report tree for an evaluation run",
        run: evaluationReportBuildCommand,
      },
      dataRunInput(input),
    ),
).pipe(Command.withDescription("Build the complete Markdown report tree for an evaluation run"))

const evaluationReport = Command.make("report", {}, () => Effect.void).pipe(
  Command.withDescription("Build evaluation reports"),
  Command.withSubcommands([evaluationReportBuild]),
)

const evaluation = Command.make("evaluation", {}, () => Effect.void).pipe(
  Command.withDescription("Work with QUALITY.md evaluation runs"),
  Command.withSubcommands([
    evaluationRun,
    evaluationCreate,
    evaluationList,
    evaluationStatus,
    evaluationData,
    evaluationReport,
  ]),
)

const versionDescriptor: CliCommand<{ readonly json: boolean }> = {
  name: "version",
  description: "Show structured qualitymd version metadata",
  run: versionCommand,
}

const version = Command.make(
  versionDescriptor.name,
  {
    json: Flag.boolean("json").pipe(Flag.withDescription("Emit machine-readable version metadata")),
  },
  (input) => execute(versionDescriptor, input),
).pipe(Command.withDescription(versionDescriptor.description))

const update = Command.make(
  "update",
  {
    check: Flag.boolean("check").pipe(
      Flag.withDescription("Check for a newer release without applying it"),
    ),
    json: Flag.boolean("json").pipe(Flag.withDescription("Emit a machine-readable update result")),
  },
  (input) =>
    execute(
      {
        name: "update",
        description: "Update the qualitymd CLI through its owning install channel",
        run: updateCommand,
      },
      input,
    ),
).pipe(Command.withDescription("Update the qualitymd CLI through its owning install channel"))

const updateRefresh = Command.make("__update-refresh", {}, () =>
  execute(
    {
      name: "__update-refresh",
      description: "Refresh the local update cache",
      run: updateRefreshCommand,
    },
    {},
  ),
).pipe(Command.withHidden)

const specDescriptor: CliCommand<Record<string, never>> = {
  name: "spec",
  description: "Emit the QUALITY.md format specification",
  run: () => specCommand,
}

const noArguments = Argument.string("argument").pipe(Argument.variadic)

const rejectArguments = (command: string, arguments_: ReadonlyArray<unknown>) =>
  arguments_.length === 0
    ? undefined
    : emit(
        commandResult("", {
          stderr: `qualitymd: ${command} accepts no arguments\n`,
          exitCode: ExitCode.usage,
        }),
      )

const spec = Command.make(specDescriptor.name, { arguments_: noArguments }, ({ arguments_ }) => {
  const rejected = rejectArguments("spec", arguments_)
  return rejected ?? execute(specDescriptor, {})
}).pipe(Command.withDescription(specDescriptor.description))

const schemaDescriptor: CliCommand<Record<string, never>> = {
  name: "schema",
  description: "Emit the companion JSON Schema for QUALITY.md frontmatter",
  run: () => schemaCommand,
}

const schema = Command.make(
  schemaDescriptor.name,
  { arguments_: noArguments },
  ({ arguments_ }) => {
    const rejected = rejectArguments("schema", arguments_)
    return rejected ?? execute(schemaDescriptor, {})
  },
).pipe(Command.withDescription(schemaDescriptor.description))

const rootWelcome =
  "QUALITY.md\n\n" +
  "The companion CLI for the QUALITY.md file format for evaluating and improving\n" +
  "the quality of AI assistant projects and harnesses.\n\n" +
  "Designed to be used with the companion agent skill:\n" +
  "  npx skills add qualitymd/quality.md\n\n" +
  "Start:\n" +
  "  qualitymd init\n" +
  "  qualitymd lint QUALITY.md\n\n" +
  "Continue:\n" +
  "  qualitymd status\n" +
  "  qualitymd evaluation create\n\n" +
  "More:\n" +
  "  qualitymd --help\n" +
  "  docs: https://getquality.md\n" +
  "  report issues: https://github.com/qualitymd/quality.md/issues\n"

export const cliRoot = Command.make("qualitymd", {}, () => emit(commandResult(rootWelcome))).pipe(
  Command.withDescription("Evaluate and improve AI assistant projects with QUALITY.md"),
  Command.withSubcommands([
    evaluation,
    init,
    lint,
    model,
    spec,
    schema,
    status,
    update,
    version,
    updateRefresh,
  ]),
)

const run = Command.runWith(cliRoot, { version: buildVersion })

const rootHelp = `qualitymd is the companion CLI for the QUALITY.md file format for evaluating and improving the quality of AI assistant projects and harnesses.

Designed to be used with the companion agent skill.

Learn more at https://getquality.md
Report issues at https://github.com/qualitymd/quality.md/issues

USAGE

  qualitymd [command] [--flags]

EXAMPLES

  npx skills add qualitymd/quality.md

COMMON TASKS

  init [path] [--flags]    Scaffold a starter QUALITY.md
  lint [path] [--flags]    Validate a QUALITY.md file
  model [command]          Query a quality model's structure and canonical reference IDs
  spec                     Emit the QUALITY.md format specification
  schema                   Emit the companion JSON Schema for QUALITY.md frontmatter
  evaluation [command]     Work with QUALITY.md evaluation runs

MANAGE

  version [--flags]        Show structured qualitymd version metadata
  update [--flags]         Update the qualitymd CLI through its owning install channel
  status [path] [--flags]  Show a QUALITY.md workspace status snapshot
  help [command]           Help about any command
  completion [command]     Generate the autocompletion script for the specified shell

FLAGS

  -h --help                Help for qualitymd
  -v --version             Version for qualitymd
`

const rootCommands =
  "evaluation init lint model spec schema status update version help completion".split(" ")

const completion = (shell: string) => {
  const words = rootCommands.filter((word) => word !== "help" && word !== "completion").join(" ")
  switch (shell) {
    case "bash":
      return `_qualitymd() { COMPREPLY=( $(compgen -W "${words}" -- "\${COMP_WORDS[COMP_CWORD]}") ); }\ncomplete -F _qualitymd qualitymd\n`
    case "zsh":
      return `#compdef qualitymd\n_arguments '1:command:(${words})'\n`
    case "fish":
      return (
        words
          .split(" ")
          .map((word) => `complete -c qualitymd -n '__fish_use_subcommand' -a ${word}`)
          .join("\n") + "\n"
      )
    case "powershell":
      return `Register-ArgumentCompleter -Native -CommandName qualitymd -ScriptBlock { param($wordToComplete) '${words}'.Split(' ') | Where-Object { $_ -like "$wordToComplete*" } }\n`
    default:
      return ""
  }
}

const completionHelp = `Generate the autocompletion script for qualitymd for the specified shell.

USAGE

  qualitymd completion [command] [--flags]

COMMANDS

  bash        Generate the autocompletion script for bash
  zsh         Generate the autocompletion script for zsh
  fish        Generate the autocompletion script for fish
  powershell  Generate the autocompletion script for powershell

FLAGS

  -h --help   Help for completion
`

export const runCli = (args: ReadonlyArray<string>) => {
  if (args.length === 1 && (args[0] === "--help" || args[0] === "-h" || args[0] === "help"))
    return emit(commandResult(rootHelp))
  if (args.length === 1 && (args[0] === "--version" || args[0] === "-v"))
    return emit(
      commandResult(
        `qualitymd version ${buildVersion}${buildCommit === undefined ? "" : ` (${buildCommit})`}\n`,
      ),
    )
  if (args[0] === "help" && args.length > 1) args = [...args.slice(1), "--help"]
  if (args[0] === "completion") {
    const shell = args[1] ?? ""
    if (shell === "" || shell === "--help" || shell === "-h")
      return emit(commandResult(completionHelp))
    const script = completion(shell)
    return emit(
      script === ""
        ? commandResult("", {
            stderr: `qualitymd: unknown completion shell ${shell}\n`,
            exitCode: ExitCode.usage,
          })
        : commandResult(script),
    )
  }
  const first = args[0]
  if (first !== undefined && !rootCommands.includes(first)) {
    const detail = first.startsWith("-")
      ? `Unknown flag: ${first}.`
      : `Unknown command "${first}" for "qualitymd".`
    return emit(
      commandResult("", {
        stderr: `qualitymd: ${detail}\n\nTry --help for usage.\n`,
        exitCode: ExitCode.usage,
      }),
    )
  }
  return run([...args]).pipe(
    Effect.flatMap(() =>
      Effect.gen(function* () {
        const notice = yield* cachedUpdateNotice(args)
        if (notice !== "") {
          const output = yield* Output
          yield* output.stderr(notice)
        }
      }),
    ),
    Effect.catchCause((cause) =>
      Effect.gen(function* () {
        const output = yield* Output
        const rendered = String(cause)
        if (rendered.includes("CliError/ShowHelp")) {
          yield* output.setExitCode(ExitCode.usage)
        } else {
          if (process.env.QUALITYMD_DEBUG !== undefined) {
            yield* output.stderr(`${String(cause)}\n`)
          } else {
            yield* output.stderr("qualitymd: internal error\n")
          }
          yield* output.setExitCode(ExitCode.internal)
        }
      }),
    ),
  )
}
