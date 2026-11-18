import { createElement } from 'react';
import { renderToStaticMarkup } from 'react-dom/server';
import { afterEach, describe, expect, it, vi } from 'vitest';
import MobileNav from '../src/islands/MobileNav';

// The menu is closed in server-rendered output, so force MobileNav's `useState(false)` open.
const setIsOpen = vi.hoisted(() => vi.fn());

vi.mock('react', async (importOriginal) => {
  const actual = await importOriginal<typeof import('react')>();
  return {
    ...actual,
    useState: ((initial: unknown) =>
      initial === false ? [true, setIsOpen] : actual.useState(initial as never)) as unknown as typeof actual.useState,
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

describe('brief-cta-click: MobileNav Contact link', () => {
  type El = { type?: unknown; props?: Record<string, unknown> };
  const items = [
    { label: 'Services', href: '/#services' },
    { label: 'Contact', href: '/#contact' },
  ];

  function anchors(node: unknown, out: El[] = []): El[] {
    if (Array.isArray(node)) node.forEach((child) => anchors(child, out));
    else if (node && typeof node === 'object') {
      const el = node as El;
      if (el.type === 'a') out.push(el);
      anchors(el.props?.children, out);
    }
    return out;
  }

  function click(label: string) {
    const track = vi.fn();
    vi.stubGlobal('umami', { track });
    setIsOpen.mockClear();
    const tree = (MobileNav as unknown as (p: unknown) => unknown)({ navItems: items, isResumePage: false });
    const link = anchors(tree).find((a) => a.props?.children === label);
    (link?.props?.onClick as () => void)();
    return track;
  }

  afterEach(() => vi.unstubAllGlobals());

  it('tracks Contact as brief-cta-click / mobile-nav and still closes the menu', () => {
    const track = click('Contact');

    expect(track).toHaveBeenCalledTimes(1);
    expect(track).toHaveBeenCalledWith('brief-cta-click', { location: 'mobile-nav' });
    expect(setIsOpen).toHaveBeenCalledWith(false);
  });

  it('does not track other nav links but still closes the menu', () => {
    const track = click('Services');

    expect(track).not.toHaveBeenCalled();
    expect(setIsOpen).toHaveBeenCalledWith(false);
  });

  it('does nothing extra when Umami is absent', () => {
    vi.unstubAllGlobals();
    setIsOpen.mockClear();
    const tree = (MobileNav as unknown as (p: unknown) => unknown)({ navItems: items, isResumePage: false });
    const link = anchors(tree).find((a) => a.props?.children === 'Contact');

    expect(() => (link?.props?.onClick as () => void)()).not.toThrow();
    expect(setIsOpen).toHaveBeenCalledWith(false);
  });

  it('keeps the Contact link markup unchanged (no hook attributes)', () => {
    const html = renderToStaticMarkup(createElement(MobileNav, { navItems: items, isResumePage: false }));
    const tag = (html.match(/<a\b[^>]*>/g) ?? []).find((t) => t.includes('href="/#contact"')) ?? '';

    expect(tag).not.toContain('data-');
  });
});
