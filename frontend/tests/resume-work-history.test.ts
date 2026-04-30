import { experimental_AstroContainer as AstroContainer } from 'astro/container';
import { expect, test } from 'vitest';
import WorkHistory from '../src/components/resume/WorkHistory.astro';

test('WorkHistory renders each job with title, company, dates, bullets, tech_stack', async () => {
  const container = await AstroContainer.create();
  const html = await container.renderToString(WorkHistory, {
    props: {
      entries: [
        {
          company: 'Nextologies',
          title: 'Software Engineer',
          dates: '2023 — Present',
          location: 'Toronto, ON',
          tech_stack: ['Go', 'Kubernetes'],
          bullets: ['Architected the platform.'],
        },
      ],
    },
  });
  expect(html).toContain('Nextologies');
  expect(html).toContain('Software Engineer');
  expect(html).toContain('2023 — Present');
  expect(html).toContain('Toronto, ON');
  expect(html).toContain('Architected the platform.');
  expect(html).toContain('Go');
  expect(html).toContain('Kubernetes');
});

test('WorkHistory handles empty bullets and missing location gracefully', async () => {
  const container = await AstroContainer.create();
  const html = await container.renderToString(WorkHistory, {
    props: {
      entries: [
        {
          company: 'Acme',
          title: 'Engineer',
          dates: '2020 — 2022',
          tech_stack: [],
          bullets: [],
        },
      ],
    },
  });
  expect(html).toContain('Acme');
  expect(html).not.toContain('<ul');  // no bullet list when empty
});
