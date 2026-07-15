import { readdir, readFile } from "node:fs/promises"
import { join } from "node:path"
import { fileURLToPath } from "node:url"
import { describe, expect, it } from "vitest"

const root = fileURLToPath(new URL("../..", import.meta.url))

const typescriptFiles = async (directory: string): Promise<ReadonlyArray<string>> => {
  const output: Array<string> = []
  for (const entry of await readdir(directory, { withFileTypes: true })) {
    const path = join(directory, entry.name)
    if (entry.isDirectory()) output.push(...(await typescriptFiles(path)))
    else if (entry.isFile() && path.endsWith(".ts")) output.push(path)
  }
  return output.sort()
}

const violations = async (directory: string, pattern: RegExp) => {
  const output: Array<string> = []
  for (const path of await typescriptFiles(join(root, "src", directory))) {
    const content = await readFile(path, "utf8")
    for (const [index, line] of content.split("\n").entries()) {
      if (pattern.test(line))
        output.push(`${path.slice(root.length + 1)}:${index + 1}: ${line.trim()}`)
    }
  }
  return output
}

describe("source architecture", () => {
  it("keeps domain modules independent of runtime layers", async () => {
    expect(
      await violations("domain", /from ["'][^"']*(?:application|services|adapters|cli)\//),
    ).toEqual([])
  })

  it("keeps services independent of application, adapter, and CLI layers", async () => {
    expect(await violations("services", /from ["'][^"']*(?:application|adapters|cli)\//)).toEqual(
      [],
    )
  })

  it("keeps direct host I/O out of domain and application modules", async () => {
    const hostAccess =
      /\b(?:process\.|Bun\.|navigator\.|performance\.|fetch\(|Date\.now\(|crypto\.getRandomValues)/
    expect([
      ...(await violations("domain", hostAccess)),
      ...(await violations("application", hostAccess)),
    ]).toEqual([])
  })
})
