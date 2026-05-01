import { experimental_AstroContainer as AstroContainer } from 'astro/container';
import { expect, test } from 'vitest';
import Languages from '../src/components/resume/Languages.astro';

test('Languages renders each entry inline; omits when empty', async () => {
  const container = await AstroContainer.create();
  const withEntries = await container.renderToString(Languages, {
    props: { entries: ['English (Native)', 'Portuguese (Native)'] },
  });
  expect(withEntries).toContain('English (Native)');
  expect(withEntries).toContain('Portuguese (Native)');

  const empty = await container.renderToString(Languages, { props: { entries: [] } });
  expect(empty.trim()).toBe('');
});
