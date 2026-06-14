

function isObject(val: unknown): val is Record<string, unknown> {
  return typeof val === 'object' && val !== null && !Array.isArray(val)
}

/** Convert PascalCase keys to camelCase (for telemetry API) */
export function pascalToCamel<T>(obj: T): T {
  if (Array.isArray(obj)) return obj.map(pascalToCamel) as unknown as T
  if (!isObject(obj)) return obj

  const result: Record<string, unknown> = {}
  for (const [key, value] of Object.entries(obj)) {
    const camelKey = key.charAt(0).toLowerCase() + key.slice(1)
    result[camelKey] = isObject(value) || Array.isArray(value) ? pascalToCamel(value) : value
  }
  return result as unknown as T
}

/** Convert snake_case keys to camelCase (for device API) */
export function snakeToCamel<T>(obj: T): T {
  if (Array.isArray(obj)) return obj.map(snakeToCamel) as unknown as T
  if (!isObject(obj)) return obj

  const result: Record<string, unknown> = {}
  for (const [key, value] of Object.entries(obj)) {
    const camelKey = key.replace(/_([a-z])/g, (_, c) => c.toUpperCase())
    result[camelKey] = isObject(value) || Array.isArray(value) ? snakeToCamel(value) : value
  }
  return result as unknown as T
}

/** Convert camelCase keys to snake_case (for sending to device API) */
export function camelToSnake<T>(obj: T): T {
  if (Array.isArray(obj)) return obj.map(camelToSnake) as unknown as T
  if (!isObject(obj)) return obj

  const result: Record<string, unknown> = {}
  for (const [key, value] of Object.entries(obj)) {
    const snakeKey = key.replace(/[A-Z]/g, (c) => '_' + c.toLowerCase())
    result[snakeKey] = isObject(value) || Array.isArray(value) ? camelToSnake(value) : value
  }
  return result as unknown as T
}
