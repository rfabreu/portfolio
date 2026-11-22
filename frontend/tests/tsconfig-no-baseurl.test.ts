import { readFileSync } from 'node:fs';
import { fileURLToPath } from 'node:url';
import { describe, expect, it } from 'vitest';

function read(relative: string) {
  return readFileSync(fileURLToPath(new URL(relative, import.meta.url)), 'utf8');
}

describe('tsconfig.json', () => {
  const config = JSON.parse(read('../tsconfig.json'));

  it('does not set the deprecated baseUrl (TS5101, removed in TS 7)', () => {
    expect(config.compilerOptions).not.toHaveProperty('baseUrl');
  });

  it('points every path alias at a ./-prefixed target', () => {
    const targets: string[] = Object.values<string[]>(config.compilerOptions.paths).flat();
    expect(targets.length).toBeGreaterThan(0);
    for (const target of targets) expect(target.startsWith('./')).toBe(true);
  });
});
