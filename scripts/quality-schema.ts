#!/usr/bin/env bun

import { writeFile } from "node:fs/promises"

import { generateQualitySchema } from "../src/domain/model/json-schema.ts"

await writeFile(new URL("../quality.schema.json", import.meta.url), generateQualitySchema())
