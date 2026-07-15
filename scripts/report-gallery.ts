#!/usr/bin/env bun

import { mkdir, readdir, readFile, rm, writeFile } from "node:fs/promises"
import { dirname, join, relative } from "node:path"

import { buildReportTree, type ReportManifest } from "../src/domain/evaluation/report.ts"
import { planEvaluation } from "../src/domain/evaluation/plan.ts"
import { jsonDocument } from "../src/domain/json.ts"
import { parseQualityDocument } from "../src/domain/model/document.ts"
import { decodeModel } from "../src/domain/model/model.ts"

const root = new URL("..", import.meta.url).pathname
const gallery = join(root, "examples/report-gallery")

const files = async (directory: string): Promise<string[]> => {
  const output: string[] = []
  for (const entry of await readdir(directory, { withFileTypes: true })) {
    const path = join(directory, entry.name)
    if (entry.isDirectory()) output.push(...(await files(path)))
    else output.push(path)
  }
  return output.sort()
}

for (const manifestPath of (await files(gallery)).filter((path) =>
  path.endsWith("/data/evaluation-manifest.json"),
)) {
  const run = dirname(dirname(manifestPath))
  const outputPath = join(run, "data/evaluation-output-result.json")
  const previousOutput = JSON.parse(await readFile(outputPath, "utf8")) as {
    readonly reportOutputs?: ReadonlyArray<{ readonly path?: string }>
  }
  for (const output of previousOutput.reportOutputs ?? []) {
    if (output.path !== undefined) await rm(join(run, output.path), { force: true })
  }
  const snapshot = join(run, "model-snapshot.md")
  const model = decodeModel(parseQualityDocument(snapshot, await readFile(snapshot, "utf8")))
  const manifest = JSON.parse(await readFile(manifestPath, "utf8")) as ReportManifest
  const payloads = await Promise.all(
    (await files(join(run, "data")))
      .filter((path) => path.endsWith(".json") && !path.endsWith("evaluation-output-result.json"))
      .map(async (path) => JSON.parse(await readFile(path, "utf8")) as Record<string, unknown>),
  )
  const tree = buildReportTree({
    model,
    manifest,
    plan: planEvaluation(model, manifest.plannedScope),
    payloads,
    runRel: `.quality/evaluations/${manifest.run.label}`,
  })
  for (const report of tree.reports) {
    const path = join(run, report.path)
    await mkdir(dirname(path), { recursive: true })
    await writeFile(path, report.content)
  }
  await writeFile(outputPath, jsonDocument(tree.output))
  console.log(`generated ${relative(root, run)}`)
}
