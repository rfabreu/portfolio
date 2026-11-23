import { readFileSync } from 'node:fs';
import { fileURLToPath } from 'node:url';
import { describe, expect, it } from 'vitest';

const toml = readFileSync(fileURLToPath(new URL('../../netlify.toml', import.meta.url)), 'utf8');

// No TOML parser is a dependency and the file is flat: split into [[redirects]] blocks.
const blocks = toml
  .split(/^\[\[redirects\]\]\s*$/m)
  .slice(1)
  .map((block) => block.split(/^\[/m)[0]);

function value(block: string, key: string) {
  const match = block.match(new RegExp(`^\\s*${key}\\s*=\\s*(.+?)\\s*$`, 'm'));
  return match?.[1];
}

describe('netlify.toml legacy site redirect', () => {
  it('has exactly one redirect rule', () => {
    expect(blocks).toHaveLength(1);
  });

  const block = blocks[0] ?? '';

  it('matches only the exact legacy host, not a bare path (previews keep serving)', () => {
    expect(value(block, 'from')).toBe('"https://rafaelswe.netlify.app/*"');
  });

  it('keeps the path on the production host via :splat', () => {
    expect(value(block, 'to')).toBe('"https://rafaelabreu.netlify.app/:splat"');
  });

  it('is a permanent redirect that overrides existing files', () => {
    expect(value(block, 'status')).toBe('301');
    expect(value(block, 'force')).toBe('true');
  });
});
