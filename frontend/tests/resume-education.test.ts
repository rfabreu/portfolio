import { experimental_AstroContainer as AstroContainer } from 'astro/container';
import { expect, test } from 'vitest';
import Education from '../src/components/resume/Education.astro';

test('Education renders school, program, and dates per entry', async () => {
  const container = await AstroContainer.create();
  const html = await container.renderToString(Education, {
    props: {
      entries: [
        { school: 'University of Toronto', program: 'Coding Bootcamp', dates: '2022 — 2023' },
      ],
    },
  });
  expect(html).toContain('University of Toronto');
  expect(html).toContain('Coding Bootcamp');
  expect(html).toContain('2022 — 2023');
});
