/**
 * Resolve a path against the configured site URL (Astro.site).
 * Returns the path unchanged when no site is configured, so callers
 * never throw during tests or local tooling that run without `site`.
 */
export function toAbsoluteUrl(path: string, site?: URL | string): string {
  if (!site) return path;
  return new URL(path, site).href;
}
