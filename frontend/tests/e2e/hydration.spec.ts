import { test, expect } from '@playwright/test';

// React reports a server/client markup mismatch as "Minified React error #418"
// in production builds and as "Hydration failed ... didn't match" in dev.
const HYDRATION_ERROR = /Minified React error #418|Hydration|did not match|didn't match/i;

const PAGES = ['/', '/playground/', '/resume/', '/project/freckles-design/', '/success/'];

for (const theme of ['dark', 'light']) {
  for (const reducedMotion of ['no-preference', 'reduce'] as const) {
    test.describe(`${theme} theme, prefers-reduced-motion: ${reducedMotion}`, () => {
      test.use({ reducedMotion });

      test.beforeEach(async ({ page }) => {
        await page.addInitScript((t) => localStorage.setItem('theme', t), theme);
      });

      for (const path of PAGES) {
        test(`${path} hydrates without a mismatch`, async ({ page }) => {
          const errors: string[] = [];
          page.on('console', (msg) => {
            if (msg.type() === 'error') errors.push(msg.text());
          });
          page.on('pageerror', (err) => errors.push(err.message));

          await page.goto(path);
          // Astro removes the `ssr` attribute from an island once it has hydrated.
          await page.waitForFunction(() => !document.querySelector('astro-island[ssr]'));

          expect(errors.filter((e) => HYDRATION_ERROR.test(e))).toEqual([]);
        });
      }
    });
  }
}

test.describe('theme toggle after hydration', () => {
  test('label follows the persisted light theme and toggling still works', async ({ page }) => {
    await page.addInitScript(() => localStorage.setItem('theme', 'light'));
    await page.goto('/');
    await page.waitForFunction(() => !document.querySelector('astro-island[ssr]'));

    const toggle = page.getByRole('button', { name: /Switch to (dark|light) mode/ }).first();
    await expect(toggle).toHaveAttribute('aria-label', 'Switch to dark mode');

    await toggle.click();
    await expect(page.locator('html')).toHaveAttribute('data-theme', 'dark');
    await expect(toggle).toHaveAttribute('aria-label', 'Switch to light mode');
  });
});
