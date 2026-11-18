import { test, expect, type Page } from '@playwright/test';

// React reports a server/client markup mismatch as "Minified React error #418"
// in production builds and as "Hydration failed ... didn't match" in dev.
const HYDRATION_ERROR = /Minified React error #418|Hydration|did not match|didn't match/i;

// /resume/ has no React islands; it is listed so a future island there is covered.
const PAGES = ['/', '/playground/', '/resume/', '/project/freckles-design/', '/success/'];

function collectErrors(page: Page): string[] {
  const errors: string[] = [];
  page.on('console', (msg) => {
    if (msg.type() === 'error') errors.push(msg.text());
  });
  page.on('pageerror', (err) => errors.push(err.message));
  return errors;
}

// Astro drops an island's `ssr` attribute as soon as hydrateRoot is scheduled, and
// React reports a mismatch only when that hydration commits, so also wait for idle.
async function waitForHydration(page: Page) {
  await page.waitForFunction(() => !document.querySelector('astro-island[ssr]'));
  await page.evaluate(
    () => new Promise<void>((resolve) => requestIdleCallback(() => resolve(), { timeout: 2000 })),
  );
}

for (const theme of ['dark', 'light']) {
  for (const reducedMotion of ['no-preference', 'reduce'] as const) {
    test.describe(`${theme} theme, prefers-reduced-motion: ${reducedMotion}`, () => {
      test.use({ reducedMotion });

      test.beforeEach(async ({ page }) => {
        await page.addInitScript((t) => localStorage.setItem('theme', t), theme);
      });

      for (const path of PAGES) {
        test(`${path} hydrates without a mismatch`, async ({ page }) => {
          const errors = collectErrors(page);

          await page.goto(path);
          await waitForHydration(page);

          expect(errors.filter((e) => HYDRATION_ERROR.test(e))).toEqual([]);
        });
      }
    });
  }
}

test.describe('no stored theme, OS prefers light', () => {
  test.use({ colorScheme: 'light' });

  for (const path of PAGES) {
    test(`${path} hydrates without a mismatch`, async ({ page }) => {
      const errors = collectErrors(page);

      await page.goto(path);
      await waitForHydration(page);

      expect(errors.filter((e) => HYDRATION_ERROR.test(e))).toEqual([]);
    });
  }
});

test.describe('theme toggle after hydration', () => {
  test('label follows the persisted light theme and toggling still works', async ({ page }) => {
    await page.addInitScript(() => localStorage.setItem('theme', 'light'));
    await page.goto('/');
    await waitForHydration(page);

    const toggle = page.getByRole('button', { name: /Switch to (dark|light) mode/ }).first();
    await expect(toggle).toHaveAttribute('aria-label', 'Switch to dark mode');

    await toggle.click();
    await expect(page.locator('html')).toHaveAttribute('data-theme', 'dark');
    await expect(toggle).toHaveAttribute('aria-label', 'Switch to light mode');
  });

  test('the mobile menu toggle matches the theme and toggles', async ({ page }) => {
    await page.setViewportSize({ width: 390, height: 800 });
    await page.addInitScript(() => localStorage.setItem('theme', 'light'));
    await page.goto('/');
    await waitForHydration(page);

    await page.getByRole('button', { name: 'Open menu' }).click();
    const toggle = page.getByRole('button', { name: /Switch to (dark|light) mode/ });
    await expect(toggle).toHaveAttribute('aria-label', 'Switch to dark mode');

    await toggle.click();
    await expect(page.locator('html')).toHaveAttribute('data-theme', 'dark');
    await expect(toggle).toHaveAttribute('aria-label', 'Switch to light mode');
  });
});
