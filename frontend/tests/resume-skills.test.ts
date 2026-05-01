import { experimental_AstroContainer as AstroContainer } from 'astro/container';
import { expect, test } from 'vitest';
import Skills from '../src/components/resume/Skills.astro';

test('Skills renders each non-empty bucket with its label and items', async () => {
  const container = await AstroContainer.create();
  const html = await container.renderToString(Skills, {
    props: {
      skills: {
        leadership: ['Strategic Planning'],
        technical: ['Go', 'Python'],
        domain: [],
      },
    },
  });
  expect(html).toContain('Leadership');
  expect(html).toContain('Strategic Planning');
  expect(html).toContain('Technical');
  expect(html).toContain('Go');
  expect(html).toContain('Python');
  expect(html).not.toContain('Domain');  // empty bucket skipped
});
