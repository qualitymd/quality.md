#!/usr/bin/env node
// Keeps the npm package README generated from the repository README.

import { copyFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const root = join(dirname(fileURLToPath(import.meta.url)), "..");

copyFileSync(
  join(root, "README.md"),
  join(root, "npm", "quality.md", "README.md"),
);
