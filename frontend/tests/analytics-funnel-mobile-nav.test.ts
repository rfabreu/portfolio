import { createElement } from 'react';
import { renderToStaticMarkup } from 'react-dom/server';
import { describe, expect, it, vi } from 'vitest';
import MobileNav from '../src/islands/MobileNav';

// The menu is closed in server-rendered output, so force MobileNav's `useState(false)` open.
vi.mock('react', async (importOriginal) => {
  const actual = await importOriginal<typeof import('react')>();
  return {
    ...actual,
    useState: ((initial: unknown) =>
      initial === false ? [true, () => {}] : actual.useState(initial as never)) as unknown as typeof actual.useState,
  };
});

describe('P10 Phase 3: MobileNav tracking', () => {
  it('tracks the mobile WORK WITH US link as orangecat-click / mobile-nav without changing its markup', () => {
    const html = renderToStaticMarkup(
      createElement(MobileNav, {
        navItems: [{ label: 'Services', href: '/#services' }],
        isResumePage: false,
      }),
    );
    const tags = html.match(/<a\b[^>]*>/g) ?? [];
    const cta = tags.filter((tag) => tag.includes('href="https://orangecatlabs.ca"'));

    expect(cta).toHaveLength(1);
    expect(cta[0]).toContain('data-umami-event="orangecat-click"');
    expect(cta[0]).toContain('data-umami-event-location="mobile-nav"');
    expect(cta[0]).toContain('target="_blank"');
    expect(cta[0]).toContain('rel="noopener noreferrer"');
    // In-page navigation links stay untracked.
    expect(tags.filter((tag) => tag.includes('data-umami-event'))).toHaveLength(1);
  });
});
