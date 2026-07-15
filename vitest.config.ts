import { defineConfig } from "vitest/config"

export default defineConfig({
  plugins: [
    {
      name: "qualitymd-markdown-text",
      transform(source, id) {
        if (!id.endsWith(".md")) return undefined
        return { code: `export default ${JSON.stringify(source)}`, map: null }
      },
    },
  ],
  test: {
    include: ["test/**/*.test.ts"],
    testTimeout: 10_000,
  },
})
