import { experimental_AstroContainer as AstroContainer } from 'astro/container';
import { expect, test } from 'vitest';
import Summary from '../src/components/resume/Summary.astro';

test('Summary renders the summary text in a section element', async () => {
  const container = await AstroContainer.create();
  const html = await container.renderToString(Summary, {
    props: { summary: 'Engineering leader building distributed systems.' },
  });
  expect(html).toContain('<section');
  expect(html).toContain('Engineering leader building distributed systems.');
});
