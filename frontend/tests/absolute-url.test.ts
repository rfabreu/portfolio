import { describe, expect, it } from 'vitest';
import { toAbsoluteUrl } from '../src/utils/absoluteUrl';

const SITE = 'https://rafaelabreu.netlify.app';

describe('toAbsoluteUrl', () => {
  it('resolves a root-relative path against the site URL', () => {
    expect(toAbsoluteUrl('/images/portfolio_website.png', SITE)).toBe(
      'https://rafaelabreu.netlify.app/images/portfolio_website.png',
    );
  });

  it('accepts a URL object for the site', () => {
    expect(toAbsoluteUrl('/resume', new URL(SITE))).toBe('https://rafaelabreu.netlify.app/resume');
  });

  it('keeps the site root for the home page path', () => {
    expect(toAbsoluteUrl('/', SITE)).toBe('https://rafaelabreu.netlify.app/');
  });

  it('returns the path unchanged when no site is configured', () => {
    expect(toAbsoluteUrl('/resume', undefined)).toBe('/resume');
  });

  it('leaves an already absolute URL untouched', () => {
    expect(toAbsoluteUrl('https://example.com/a.png', SITE)).toBe('https://example.com/a.png');
  });
});
