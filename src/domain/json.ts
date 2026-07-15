export const jsonStringify = (value: unknown, space = 2): string =>
  JSON.stringify(value, null, space)
    .replaceAll("<", "\\u003c")
    .replaceAll(">", "\\u003e")
    .replaceAll("&", "\\u0026")
    .replaceAll("\u2028", "\\u2028")
    .replaceAll("\u2029", "\\u2029")

export const jsonDocument = (value: unknown): string => `${jsonStringify(value)}\n`

const sortValue = (value: unknown): unknown => {
  if (Array.isArray(value)) return value.map(sortValue)
  if (value !== null && typeof value === "object") {
    return Object.fromEntries(
      Object.entries(value as Record<string, unknown>)
        .sort(([left], [right]) => (left < right ? -1 : left > right ? 1 : 0))
        .map(([key, entry]) => [key, sortValue(entry)]),
    )
  }
  return value
}

export const canonicalJson = (value: unknown): string => jsonStringify(sortValue(value), 0)

export const sha256 = async (value: string | Uint8Array): Promise<string> => {
  const bytes = typeof value === "string" ? new TextEncoder().encode(value) : Uint8Array.from(value)
  const digest = await crypto.subtle.digest("SHA-256", bytes)
  return [...new Uint8Array(digest)].map((byte) => byte.toString(16).padStart(2, "0")).join("")
}

export const hashJson = (value: unknown): Promise<string> => sha256(canonicalJson(value))
