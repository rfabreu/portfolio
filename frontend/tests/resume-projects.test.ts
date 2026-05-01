import { experimental_AstroContainer as AstroContainer } from 'astro/container';
import { expect, test } from 'vitest';
import Projects from '../src/components/resume/Projects.astro';

test('Projects renders headline and links to case study by slug', async () => {
  const container = await AstroContainer.create();
  const html = await container.renderToString(Projects, {
    props: {
      entries: [
        { slug: 'fptv', headline: 'Glassmorphism dashboard for Toronto TV broadcaster' },
      ],
    },
  });
  expect(html).toContain('Glassmorphism dashboard for Toronto TV broadcaster');
  expect(html).toContain('href="/project/fptv"');
});

test('Projects omits section when entries is empty', async () => {
  const container = await AstroContainer.create();
  const html = await container.renderToString(Projects, { props: { entries: [] } });
  expect(html.trim()).toBe('');
});
