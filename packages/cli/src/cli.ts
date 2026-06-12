import { Console } from "effect";
import { Command, Flag } from "effect/unstable/cli";

/**
 * The root `qualitymd` command.
 *
 * Add subcommands with `Command.withSubcommands` as the CLI grows.
 */
const root = Command.make(
  "qualitymd",
  {
    name: Flag.string("name").pipe(
      Flag.withAlias("n"),
      Flag.withDescription("Who to greet"),
      Flag.withDefault("world"),
    ),
  },
  ({ name }) => Console.log(`Hello, ${name}! quality.md is ready.`),
).pipe(Command.withDescription("quality.md command-line interface"));

/**
 * Parses the provided arguments and runs the root command. The resulting
 * `Effect` still requires the CLI `Environment`, which the entry point
 * (`bin.ts`) provides via the Node services layer.
 */
export const run = Command.run(root, { version: "0.0.0" });
