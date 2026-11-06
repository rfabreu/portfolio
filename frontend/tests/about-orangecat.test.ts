import { experimental_AstroContainer as AstroContainer } from 'astro/container';
import { describe, expect, it } from 'vitest';
import About from '../src/components/About.astro';

async function renderAbout() {
  const container = await AstroContainer.create();
  return container.renderToString(About);
}

describe('P10: About section copy', () => {
  it('introduces Rafael as co-founder with a single OrangeCat Labs link', async () => {
    const html = await renderAbout();

    expect(html).toContain('co-founder of');
    expect(html).toContain('href="https://orangecatlabs.ca"');
    expect(html).toContain('target="_blank"');
    expect(html).toContain('rel="noopener noreferrer"');
    expect(html).toContain('a Canadian software studio');
    expect(html.split('orangecatlabs.ca').length - 1).toBe(2);
  });

  it('tells assistive tech that the link opens in a new tab', async () => {
    const html = await renderAbout();

    expect(html).toContain('aria-label="OrangeCat Labs (opens orangecatlabs.ca in a new tab)"');
  });

  it('uses the approved founder-voice copy', async () => {
    const html = await renderAbout();

    expect(html).toContain('24/7 broadcast operations');
    expect(html).toContain('real-time infrastructure for global clients');
    expect(html).toContain('My current focus is on intelligent systems that use AI to solve real problems.');
  });

  it('removes the old job-seeker opener and credentials paragraph', async () => {
    const html = await renderAbout();

    expect(html).not.toContain('Software Engineer with experience across');
    expect(html).not.toContain('I completed the University of Toronto Coding Bootcamp');
  });

  it('keeps the focus areas and education cards', async () => {
    const html = await renderAbout();

    expect(html).toContain('FOCUS AREAS');
    expect(html).toContain('EDUCATION');
    expect(html).toContain('University of Toronto');
  });

  it('contains no pricing, budget, or Upwork language', async () => {
    const html = await renderAbout();

    expect(html).not.toMatch(/\b(price|pricing|budget|quote|upwork)\b/i);
  });
});
