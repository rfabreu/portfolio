import { experimental_AstroContainer as AstroContainer } from 'astro/container';
import { describe, expect, it } from 'vitest';
import fs from 'node:fs';
import path from 'node:path';
import reactRenderer from '@astrojs/react/server.js';
import Hero from '../src/components/Hero.astro';
import Nav from '../src/components/Nav.astro';
import Footer from '../src/components/Footer.astro';
import About from '../src/components/About.astro';
import Contact from '../src/components/Contact.astro';
import Services from '../src/components/Services.astro';
import ProjectLayout from '../src/layouts/ProjectLayout.astro';

const SRC = path.resolve(__dirname, '../src');

async function createContainerWithReact() {
  const container = await AstroContainer.create();
  container.addServerRenderer({
    name: '@astrojs/react',
    renderer: reactRenderer,
  });
  return container;
}

/** Opening `<a ...>` tags in a chunk of markup. `=>` is neutralised so JSX handlers don't end a tag early. */
function anchorTags(markup: string): string[] {
  return markup.replace(/=>/g, '=_').match(/<a\b[^>]*>/g) ?? [];
}

function anchorsWithHref(markup: string, href: string): string[] {
  return anchorTags(markup).filter((tag) => tag.includes(`href="${href}"`));
}

function attr(tag: string, name: string): string | null {
  return tag.match(new RegExp(`\\s${name}="([^"]*)"`))?.[1] ?? null;
}

function sourceFiles(dir = SRC): string[] {
  return fs.readdirSync(dir, { withFileTypes: true }).flatMap((entry) => {
    const full = path.join(dir, entry.name);
    if (entry.isDirectory()) return sourceFiles(full);
    return /\.(astro|tsx|md|mdx)$/.test(entry.name) ? [full] : [];
  });
}

// Deliberately loose: any scheme, `www.`, path, query or quoting still counts as a funnel link.
const ORANGECAT_HREF = /href=["']?[^\s>]*orangecatlabs\.ca|href=\{ORANGECAT_URL\}/;
const LINKEDIN_HREF = /href=["']?[^\s>]*linkedin\.com\/in\/rafael|href=\{resume\.data\.links\.linkedin\}/;

function funnelAnchors(pattern: RegExp) {
  return sourceFiles().flatMap((file) =>
    anchorTags(fs.readFileSync(file, 'utf-8'))
      .filter((tag) => pattern.test(tag))
      .map((tag) => ({ file: path.relative(SRC, file), tag })),
  );
}

describe('P10 Phase 3: orangecat-click attributes', () => {
  it('Hero primary CTA is tracked as hero', async () => {
    const container = await createContainerWithReact();
    const html = await container.renderToString(Hero);
    const [tag, ...rest] = anchorsWithHref(html, 'https://orangecatlabs.ca');

    expect(rest).toHaveLength(0);
    expect(attr(tag, 'data-umami-event')).toBe('orangecat-click');
    expect(attr(tag, 'data-umami-event-location')).toBe('hero');
  });

  it('Nav desktop CTA is tracked as nav', async () => {
    const container = await createContainerWithReact();
    const html = await container.renderToString(Nav);
    const [tag, ...rest] = anchorsWithHref(html, 'https://orangecatlabs.ca');

    expect(rest).toHaveLength(0);
    expect(attr(tag, 'data-umami-event')).toBe('orangecat-click');
    expect(attr(tag, 'data-umami-event-location')).toBe('nav');
  });

  it('Footer credit is tracked as footer', async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(Footer);
    const [tag, ...rest] = anchorsWithHref(html, 'https://orangecatlabs.ca');

    expect(rest).toHaveLength(0);
    expect(attr(tag, 'data-umami-event')).toBe('orangecat-click');
    expect(attr(tag, 'data-umami-event-location')).toBe('footer');
  });

  it('About inline link is tracked as about', async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(About);
    const [tag, ...rest] = anchorsWithHref(html, 'https://orangecatlabs.ca');

    expect(rest).toHaveLength(0);
    expect(attr(tag, 'data-umami-event')).toBe('orangecat-click');
    expect(attr(tag, 'data-umami-event-location')).toBe('about');
  });

  it('Contact studio link is tracked as contact', async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(Contact);
    const [tag, ...rest] = anchorsWithHref(html, 'https://orangecatlabs.ca');

    expect(rest).toHaveLength(0);
    expect(attr(tag, 'data-umami-event')).toBe('orangecat-click');
    expect(attr(tag, 'data-umami-event-location')).toBe('contact');
  });

  describe('ProjectLayout CTA', () => {
    const props = {
      title: 'Test Project',
      description: 'Test Description',
      image: '/images/placeholder.svg',
      tags: ['Go', 'React'],
      category: 'Full-Stack',
    };

    function projectCta(html: string) {
      return anchorsWithHref(html, 'https://orangecatlabs.ca').filter(
        (tag) => attr(tag, 'data-umami-event-location') === 'project-cta',
      );
    }

    it('is tracked as project-cta with the case study slug', async () => {
      const container = await createContainerWithReact();
      const html = await container.renderToString(ProjectLayout, {
        props: { ...props, slug: 'freckles-design' },
      });
      const [tag, ...rest] = projectCta(html);

      expect(rest).toHaveLength(0);
      expect(attr(tag, 'data-umami-event')).toBe('orangecat-click');
      expect(attr(tag, 'data-umami-event-project')).toBe('freckles-design');
    });

    it('omits the project property when no slug is supplied', async () => {
      const container = await createContainerWithReact();
      const html = await container.renderToString(ProjectLayout, { props });
      const [tag] = projectCta(html);

      expect(tag).toBeDefined();
      expect(tag).not.toContain('data-umami-event-project');
    });
  });
});

describe('P10 Phase 3: linkedin-click attributes', () => {
  function linkedinTags(html: string) {
    return anchorTags(html).filter((tag) => /href="https:\/\/(www\.)?linkedin\.com\/in\//.test(tag));
  }

  it('Footer LinkedIn link is tracked as footer', async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(Footer);
    const [tag, ...rest] = linkedinTags(html);

    expect(rest).toHaveLength(0);
    expect(attr(tag, 'data-umami-event')).toBe('linkedin-click');
    expect(attr(tag, 'data-umami-event-location')).toBe('footer');
  });

  it('Contact LinkedIn link is tracked as contact', async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(Contact);
    const [tag, ...rest] = linkedinTags(html);

    expect(rest).toHaveLength(0);
    expect(attr(tag, 'data-umami-event')).toBe('linkedin-click');
    expect(attr(tag, 'data-umami-event-location')).toBe('contact');
  });

  // ResumeModal and the resume page read the `resume` content collection, which the Astro
  // container does not load under Vitest, so they are asserted from source here and from the
  // built HTML during build verification.
  it('Resume modal LinkedIn link is tracked as resume-modal', () => {
    const source = fs.readFileSync(path.join(SRC, 'components/ResumeModal.astro'), 'utf-8');
    const [tag, ...rest] = anchorTags(source).filter((t) => LINKEDIN_HREF.test(t));

    expect(rest).toHaveLength(0);
    expect(attr(tag, 'data-umami-event')).toBe('linkedin-click');
    expect(attr(tag, 'data-umami-event-location')).toBe('resume-modal');
  });

  it('Resume page header and footer LinkedIn links are tracked separately', () => {
    const source = fs.readFileSync(path.join(SRC, 'pages/resume.astro'), 'utf-8');
    const tags = anchorTags(source).filter((t) => LINKEDIN_HREF.test(t));

    expect(tags.map((tag) => attr(tag, 'data-umami-event'))).toEqual(['linkedin-click', 'linkedin-click']);
    expect(tags.map((tag) => attr(tag, 'data-umami-event-location'))).toEqual([
      'resume-page-header',
      'resume-page-footer',
    ]);
  });

  it('leaves GitHub links untracked', async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(Contact);
    const github = anchorsWithHref(html, 'https://github.com/rfabreu');

    expect(github).toHaveLength(1);
    expect(github[0]).not.toContain('data-umami-event');
  });
});

describe('P10 Phase 3: view-projects-click', () => {
  it('Hero View Projects link carries the hook for the imperative listener and no declarative event', async () => {
    const container = await createContainerWithReact();
    const html = await container.renderToString(Hero);
    const [tag, ...rest] = anchorsWithHref(html, '#projects');

    expect(rest).toHaveLength(0);
    expect(tag).toContain('data-view-projects');
    // A declarative event would make Umami preventDefault and delay the in-page scroll.
    expect(tag).not.toContain('data-umami-event');
  });

  it('Hero script tracks the click without preventing default', () => {
    const hero = fs.readFileSync(path.join(SRC, 'components/Hero.astro'), 'utf-8');
    const script = hero.match(/<script>([\s\S]*?)<\/script>/)?.[1] ?? '';

    expect(script).toContain('[data-view-projects]');
    expect(script).toMatch(/umami\.track\('view-projects-click', \{ location: 'hero' \}\)/);
    expect(script).not.toContain('preventDefault');
  });
});

describe('P10 Phase 3: contact form events', () => {
  const contact = fs.readFileSync(path.join(SRC, 'components/Contact.astro'), 'utf-8');
  const script = contact.match(/<script>([\s\S]*?)<\/script>/)?.[1] ?? '';

  it('tracks validated submissions via the form submit event, with no event data', () => {
    expect(script).toContain(`form[name="contact"]`);
    expect(script).toContain(`addEventListener('submit'`);
    expect(script).toContain(`umami.track('contact-form-submit')`);
  });

  it('never reads form values', () => {
    expect(script).not.toMatch(/FormData|\.value|\.elements|\.checked|\.files|\.text(Content)?\b/);
  });

  it('puts no Umami event on the submit button (it would also count failed validation)', async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(Contact);
    const form = html.slice(html.indexOf('<form'), html.indexOf('</form>'));

    expect(form).not.toContain('data-umami-event');
  });

  it('success page tracks contact-form-success on load', () => {
    const success = fs.readFileSync(path.join(SRC, 'pages/success.astro'), 'utf-8');
    const successScript = success.match(/<script>([\s\S]*?)<\/script>/)?.[1] ?? '';

    expect(successScript).toContain(`addEventListener('load'`);
    expect(successScript).toContain(`umami.track('contact-form-success')`);
  });
});

describe('P10 Phase 3: source-level guards', () => {
  it('every OrangeCat Labs link in src carries orangecat-click and the full catalogue is present', () => {
    const found = funnelAnchors(ORANGECAT_HREF);

    for (const { file, tag } of found) {
      expect(tag, `${file} is missing orangecat-click`).toContain('data-umami-event="orangecat-click"');
    }
    expect(found.map(({ tag }) => tag.match(/data-umami-event-location="([^"]+)"/)?.[1]).sort()).toEqual(
      [
        'about',
        'contact',
        'footer',
        'freckles-credit',
        'hero',
        'mobile-nav',
        'nav',
        'project-cta',
      ].sort(),
    );
  });

  it('every Rafael LinkedIn link in src carries linkedin-click and the full catalogue is present', () => {
    const found = funnelAnchors(LINKEDIN_HREF);

    for (const { file, tag } of found) {
      expect(tag, `${file} is missing linkedin-click`).toContain('data-umami-event="linkedin-click"');
    }
    const locations = found.map(({ tag }) => tag.match(/data-umami-event-location="([^"]+)"/)?.[1]).sort();
    expect(locations).toEqual(
      ['contact', 'footer', 'resume-modal', 'resume-page-footer', 'resume-page-header'].sort(),
    );
  });

  it('has no Markdown-syntax links to OrangeCat Labs or Rafael\'s LinkedIn (they cannot carry tracking attributes)', () => {
    const offenders = sourceFiles().flatMap((file) =>
      [...fs.readFileSync(file, 'utf-8').matchAll(/\]\([^)]*(orangecatlabs\.ca|linkedin\.com\/in\/rafael)[^)]*\)/g)].map(
        () => path.relative(SRC, file),
      ),
    );

    expect(offenders).toEqual([]);
  });

  it('the Freckles Design credit is attributed as freckles-credit', () => {
    const md = fs.readFileSync(path.join(SRC, 'content/projects/freckles-design.md'), 'utf-8');
    const [tag] = anchorsWithHref(md, 'https://orangecatlabs.ca');

    expect(attr(tag, 'data-umami-event')).toBe('orangecat-click');
    expect(attr(tag, 'data-umami-event-location')).toBe('freckles-credit');
  });

  it('only sends the location and project event properties (no PII keys)', () => {
    const keys = new Set<string>();
    for (const file of sourceFiles()) {
      for (const match of fs.readFileSync(file, 'utf-8').matchAll(/data-umami-event-([\w-]+)/g)) {
        keys.add(match[1]);
      }
    }

    expect([...keys].sort()).toEqual(['location', 'project']);
  });

  it('every umami.track call is behind a typeof guard and the event list is the documented one', () => {
    const names = new Set<string>();
    for (const file of sourceFiles()) {
      const source = fs.readFileSync(file, 'utf-8');
      const calls = [...source.matchAll(/umami\.track\('([\w-]+)'/g)];
      const guardedCalls = source.match(/typeof umami !== 'undefined'\) umami\.track\('/g) ?? [];

      expect(guardedCalls.length, `${path.relative(SRC, file)} has an unguarded umami.track`).toBe(calls.length);
      calls.forEach((call) => names.add(call[1]));
    }

    expect([...names].sort()).toEqual(
      ['brief-cta-click', 'chat-opened', 'contact-form-submit', 'contact-form-success', 'game-started', 'view-projects-click'].sort(),
    );
  });
});

describe('brief-cta-click: in-page links to the contact form', () => {
  const read = (file: string) => fs.readFileSync(path.join(SRC, file), 'utf-8');
  // The script block that owns the hook, so unrelated scripts in the same file (e.g. the resume modal opener) don't count.
  const scriptOf = (file: string, location: string) =>
    [...read(file).matchAll(/<script>([\s\S]*?)<\/script>/g)].map((m) => m[1]).find((body) => body.includes(`"${location}"`)) ?? '';

  /** Astro links carry a `data-brief-cta` hook and a listener that never prevents default. */
  function expectListener(file: string, location: string) {
    const script = scriptOf(file, location);
    expect(script, `${file} listener selector`).toContain(`[data-brief-cta="${location}"]`);
    expect(script, `${file} track call`).toMatch(
      new RegExp(`typeof umami !== 'undefined'\\) umami\\.track\\('brief-cta-click', \\{ location: '${location}' \\}\\)`),
    );
    expect(script, `${file} must not delay navigation`).not.toContain('preventDefault');
  }

  it('Services "Submit a brief" link is tracked as services', async () => {
    const container = await createContainerWithReact();
    const html = await container.renderToString(Services);
    const [tag, ...rest] = anchorsWithHref(html, '#contact');

    expect(rest).toHaveLength(0);
    expect(attr(tag, 'data-brief-cta')).toBe('services');
    expect(tag).not.toContain('data-umami-event');
    expectListener('components/Services.astro', 'services');
  });

  it('desktop nav Contact is tracked as nav; no other nav link carries the hook', async () => {
    const container = await createContainerWithReact();
    const html = await container.renderToString(Nav);
    const hooked = anchorTags(html).filter((tag) => tag.includes('data-brief-cta'));
    const [tag, ...rest] = anchorsWithHref(html, '/#contact');

    expect(rest).toHaveLength(0);
    expect(hooked).toEqual([tag]);
    expect(attr(tag, 'data-brief-cta')).toBe('nav');
    expect(tag).not.toContain('data-umami-event');
    expectListener('components/Nav.astro', 'nav');
  });

  it('resume modal and resume page footer "contact form" links are tracked', () => {
    for (const [file, location] of [
      ['components/ResumeModal.astro', 'resume-modal'],
      ['pages/resume.astro', 'resume-page-footer'],
    ]) {
      const [tag, ...rest] = anchorsWithHref(read(file), '/#contact');
      expect(rest).toHaveLength(0);
      expect(attr(tag, 'data-brief-cta')).toBe(location);
      expect(tag).not.toContain('data-umami-event');
      expectListener(file, location);
    }
  });

  it('ChatWidget "Leave a message instead" link tracks chat-error and still closes the chat', () => {
    const [tag, ...rest] = anchorTags(read('islands/ChatWidget.tsx')).filter((t) => t.includes('href="#contact"'));

    expect(rest).toHaveLength(0);
    expect(tag).toContain('onClick={trackBriefFromChat}');
    expect(read('islands/ChatWidget.tsx')).toMatch(
      /function trackBriefFromChat\(\) \{\s*if \(typeof umami !== 'undefined'\) umami\.track\('brief-cta-click', \{ location: 'chat-error' \}\);\s*setIsOpen\(false\);\s*\}/,
    );
  });

  it('every in-page link to #contact in src is a known tracked placement', () => {
    const found = sourceFiles().flatMap((file) =>
      anchorTags(fs.readFileSync(file, 'utf-8'))
        .filter((tag) => /href=["']?\/?#contact/.test(tag))
        .map((tag) => path.relative(SRC, file)),
    );
    // Nav's `/#contact` lives in a data array, and MobileNav reuses it; both are covered above and below.
    expect(found.sort()).toEqual(
      ['components/ResumeModal.astro', 'components/Services.astro', 'islands/ChatWidget.tsx', 'pages/resume.astro'].sort(),
    );
    expect(read('components/Nav.astro').match(/href: '\/?#contact'/g)).toHaveLength(1);
  });
});

describe('Umami reporting scope', () => {
  it('every Umami script tag reports only from the production host', () => {
    const tags = ['layouts/BaseLayout.astro', 'layouts/ProjectLayout.astro'].flatMap(
      (file) => fs.readFileSync(path.join(SRC, file), 'utf-8').match(/<script[^>]*cloud\.umami\.is[^>]*>/g) ?? [],
    );
    expect(tags).toHaveLength(2);
    for (const tag of tags) expect(tag).toContain('data-domains="rafaelabreu.netlify.app"');
  });
});
