import * as Option from "effect/Option"
import type { Command, HelpDoc } from "effect/unstable/cli"

import { cliRoot } from "../src/cli/app.ts"

const outputPath = new URL("../mintlify/cli.mdx", import.meta.url)

const frontmatter = `---
title: CLI reference
description: The qualitymd command-line interface that the /quality skill builds on.
---

{/* Generated from the Effect CLI command tree by scripts/cli-docs.ts. Do not edit directly. Run \`mise run cli-docs\`. */}

The \`/quality\` skill drives the \`qualitymd\` CLI for you. This reference is for
when you want to run commands directly. It is generated from the CLI's own
command definitions, so it always matches the installed binary.
`

type AnyCommand = Command.Command.Any

type IntrospectableCommand = AnyCommand & {
  readonly buildHelpDoc: (path: ReadonlyArray<string>) => HelpDoc.HelpDoc
}

const escapeMdx = (value: string) =>
  value.replaceAll("<", "&lt;").replaceAll("{", "&#123;").replaceAll("}", "&#125;")

const cell = (value: string) => escapeMdx(value.replaceAll("\n", " ").replaceAll("|", "\\|"))

const visible = (command: AnyCommand) =>
  command.subcommands.flatMap((group) => group.commands).filter((child) => !child.hidden)

const commandPath = (path: ReadonlyArray<string>) => path.join(" ")

const renderFlags = (doc: HelpDoc.HelpDoc) => {
  if (doc.flags.length === 0) return ""
  let output = "\n**Flags**\n\n| Flag | Description |\n| --- | --- |\n"
  for (const flag of doc.flags) {
    const names = [flag.name, ...flag.aliases]
      .map((name) => `\`${name.startsWith("-") ? name : `--${name}`}\``)
      .join(", ")
    output += `| ${names} | ${cell(Option.getOrElse(flag.description, () => ""))} |\n`
  }
  return output
}

const renderCommand = (command: AnyCommand, path: ReadonlyArray<string>, level: number): string => {
  const current = command as IntrospectableCommand
  const doc = current.buildHelpDoc(path)
  const children = visible(command)
  let output = `\n${"#".repeat(Math.min(level, 6))} \`${commandPath(path)}\`\n\n`
  if (doc.description !== "") output += `${escapeMdx(doc.description)}\n\n`
  output += `\`\`\`bash\n${escapeMdx(doc.usage)}\n\`\`\`\n`
  if (children.length > 0) {
    output += "\n**Subcommands**\n\n| Command | Description |\n| --- | --- |\n"
    for (const child of children)
      output += `| \`${commandPath([...path, child.name])}\` | ${cell(child.description ?? "")} |\n`
  }
  output += renderFlags(doc)
  if ((doc.examples?.length ?? 0) > 0) {
    output += "\n**Examples**\n\n```bash\n"
    output += doc.examples!.map((example) => example.command).join("\n")
    output += "\n```\n"
  }
  for (const child of children) output += renderCommand(child, [...path, child.name], level + 1)
  return output
}

const top = visible(cliRoot)
let output = frontmatter + "\n## Commands\n\n| Command | Description |\n| --- | --- |\n"
for (const command of top)
  output += `| \`qualitymd ${command.name}\` | ${cell(command.description ?? "")} |\n`
for (const command of top) output += renderCommand(command, ["qualitymd", command.name], 2)

await Bun.write(outputPath, output)
