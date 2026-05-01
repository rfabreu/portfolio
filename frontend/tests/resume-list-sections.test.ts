import { experimental_AstroContainer as AstroContainer } from 'astro/container';
import { expect, test } from 'vitest';
import Certifications from '../src/components/resume/Certifications.astro';
import Awards from '../src/components/resume/Awards.astro';
import Speaking from '../src/components/resume/Speaking.astro';

test('Certifications renders entries with name/issuer/year, omits when empty', async () => {
  const container = await AstroContainer.create();
  const withEntries = await container.renderToString(Certifications, {
    props: { entries: [{ name: 'OpenShift Specialist', issuer: 'Red Hat', year: '2024' }] },
  });
  expect(withEntries).toContain('OpenShift Specialist');
  expect(withEntries).toContain('Red Hat');
  expect(withEntries).toContain('2024');

  const empty = await container.renderToString(Certifications, { props: { entries: [] } });
  expect(empty.trim()).toBe('');
});

test('Awards renders entries with name/org/year, omits when empty', async () => {
  const container = await AstroContainer.create();
  const withEntries = await container.renderToString(Awards, {
    props: { entries: [{ name: 'Hackathon Winner', org: 'TechFest', year: '2025' }] },
  });
  expect(withEntries).toContain('Hackathon Winner');
  expect(withEntries).toContain('TechFest');

  const empty = await container.renderToString(Awards, { props: { entries: [] } });
  expect(empty.trim()).toBe('');
});

test('Speaking renders entries with title/venue/year and link when present, omits when empty', async () => {
  const container = await AstroContainer.create();
  const withEntries = await container.renderToString(Speaking, {
    props: {
      entries: [
        { type: 'talk', title: 'Building with Go', venue: 'GoTO Meetup', year: '2026', url: 'https://example.com/talk' },
      ],
    },
  });
  expect(withEntries).toContain('Building with Go');
  expect(withEntries).toContain('GoTO Meetup');
  expect(withEntries).toContain('href="https://example.com/talk"');

  const empty = await container.renderToString(Speaking, { props: { entries: [] } });
  expect(empty.trim()).toBe('');
});
