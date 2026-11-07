import { experimental_AstroContainer as AstroContainer } from 'astro/container';
import { describe, expect, it } from 'vitest';
import Services from '../src/components/Services.astro';
import Contact from '../src/components/Contact.astro';
import ProjectLayout from '../src/layouts/ProjectLayout.astro';
import reactRenderer from '@astrojs/react/server.js';
import fs from 'node:fs';
import path from 'node:path';

async function createContainerWithReact() {
  const container = await AstroContainer.create();
  container.addServerRenderer({
    name: '@astrojs/react',
    renderer: reactRenderer,
  });
  return container;
}

describe('P10 Phase 2: Services & Project Brief Form', () => {
  it('Services renders 4 core studio offerings and deliverables', async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(Services);

    expect(html).toContain('SERVICES');
    expect(html).toContain('What We Build');
    expect(html).toContain('Full-Stack Software Development');
    expect(html).toContain('Cloud Architecture &amp; DevOps');
    expect(html).toContain('AI Integration &amp; Automation');
    expect(html).toContain('Systems Modernization &amp; Optimization');
    expect(html).toContain('OrangeCat Labs');
    expect(html).toContain('href="#contact"');
  });

  it('Contact renders project brief form fields without exposing email addresses', async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(Contact);

    expect(html).toContain('INQUIRY');
    expect(html).toContain('Start a Project');
    expect(html).toContain('name="project_type"');
    expect(html).toContain('name="timeline"');
    expect(html).not.toContain('name="budget"');
    expect(html).toContain('Submit Project Brief →');
    expect(html).toContain('orangecatlabs.ca');
    // Ensure no mailto links exist
    expect(html).not.toContain('href="mailto:');
  });

  it('ProjectLayout renders conversion CTA linking to OrangeCat Labs', async () => {
    const container = await createContainerWithReact();
    const html = await container.renderToString(ProjectLayout, {
      props: {
        title: 'Test Project',
        description: 'Test Description',
        image: '/images/placeholder.svg',
        tags: ['Go', 'React'],
        category: 'Full-Stack',
      },
    });

    expect(html).toContain('Have a similar project in mind?');
    expect(html).toContain('Start a Project with Us →');
    expect(html).toContain('href="https://orangecatlabs.ca"');
  });

  it('Resume profile headline is aligned with studio co-founder position', () => {
    const profilePath = path.resolve(__dirname, '../src/content/resume/profile.md');
    const content = fs.readFileSync(profilePath, 'utf-8');
    expect(content).toContain('title: "Co-Founder, OrangeCat Labs | Full-Stack & Infrastructure Leader"');
  });
});
