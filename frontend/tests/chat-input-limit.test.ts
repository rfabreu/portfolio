import { readFileSync } from 'node:fs';
import { fileURLToPath } from 'node:url';
import { describe, expect, it } from 'vitest';

function read(relative: string) {
  return readFileSync(fileURLToPath(new URL(relative, import.meta.url)), 'utf8');
}

describe('Chat input limit', () => {
  it('matches the API per-message cap for user messages', () => {
    const widget = read('../src/islands/ChatWidget.tsx');
    const handler = read('../../api/internal/chat/handler.go');

    const widgetLimit = widget.match(/const MAX_INPUT_LENGTH = (\d+);/)?.[1];
    const apiLimit = handler.match(/MaxUserContentRunes = (\d+)/)?.[1];

    expect(widgetLimit).toBeDefined();
    expect(apiLimit).toBeDefined();
    expect(widgetLimit).toBe(apiLimit);
    expect(widget).toContain('maxLength={MAX_INPUT_LENGTH}');
  });
});
