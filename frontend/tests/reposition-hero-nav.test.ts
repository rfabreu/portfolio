import { experimental_AstroContainer as AstroContainer } from 'astro/container';
import { describe, expect, it } from 'vitest';
import Hero from '../src/components/Hero.astro';
import Nav from '../src/components/Nav.astro';
import Footer from '../src/components/Footer.astro';
import BaseLayout from '../src/layouts/BaseLayout.astro';
import reactRenderer from '@astrojs/react/server.js';

async function createContainerWithReact() {
  const container = await AstroContainer.create();
  container.addServerRenderer({
    name: '@astrojs/react',
    renderer: reactRenderer,
  });
  return container;
}

describe('P10 Phase 1: Repositioning Components', () => {
  it('Hero renders Co-Founder branding, outcome copy, and OrangeCat Labs CTA', async () => {
    const container = await createContainerWithReact();
    const html = await container.renderToString(Hero);

    expect(html).toContain('// CO-FOUNDER · ORANGECAT LABS');
    expect(html).toContain('I build solutions that drive results');
    expect(html).toContain('Work with OrangeCat Labs');
    expect(html).toContain('href="https://orangecatlabs.ca"');
    expect(html).toContain('target="_blank"');
    expect(html).toContain('href="#projects"');
    expect(html).toContain('View Projects');
    // hero-chat-trigger button was removed
    expect(html).not.toContain('id="hero-chat-trigger"');
  });

  it('Nav renders WORK WITH US linking to OrangeCat Labs', async () => {
    const container = await createContainerWithReact();
    const html = await container.renderToString(Nav);

    expect(html).toContain('WORK WITH US');
    expect(html).toContain('href="https://orangecatlabs.ca"');
    expect(html).toContain('target="_blank"');
  });

  it('Footer renders elevated Co-Founded credit', async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(Footer);

    expect(html).toContain('Co-Founded by Rafael');
    expect(html).toContain('OrangeCat Labs');
    expect(html).toContain('href="https://orangecatlabs.ca"');
  });

  it('BaseLayout renders new title and OpenGraph metadata', async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(BaseLayout, {
      props: {
        includeResumeModal: false,
      },
    });

    expect(html).toContain('<title>Rafael Abreu | Co-Founder, OrangeCat Labs</title>');
    expect(html).toContain('content="website"');
    expect(html).toContain('content="Rafael Abreu | Co-Founder, OrangeCat Labs"');
    expect(html).toContain('name="twitter:card"');
  });
});
