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

// The toggle renders both icons and CSS picks one from html[data-theme], so the
// right icon shows from first paint, before the island hydrates.
function toggleIcons(page: Page) {
  const button = page.getByRole('button', { name: /Switch to (dark|light) mode/ }).first();
  return { sun: button.locator('.theme-icon-sun'), moon: button.locator('.theme-icon-moon') };
}

test.describe('theme toggle icon before hydration', () => {
  // Blocking the island bundles leaves the inline theme script running, so this is
  // the page a visitor sees before React hydrates the toggle.
  test.beforeEach(async ({ page }) => {
    await page.route('**/_astro/**/*.js', (route) => route.abort());
  });

  for (const theme of ['dark', 'light']) {
    test(`stored ${theme} theme shows the ${theme === 'dark' ? 'sun' : 'moon'} icon`, async ({ page }) => {
      await page.addInitScript((t) => localStorage.setItem('theme', t), theme);
      await page.goto('/');

      await expect(page.locator('html')).toHaveAttribute('data-theme', theme);
      await expect(page.locator('astro-island[ssr]').first()).toBeAttached();
      const { sun, moon } = toggleIcons(page);
      await expect(theme === 'dark' ? sun : moon).toBeVisible();
      await expect(theme === 'dark' ? moon : sun).toBeHidden();
    });
  }

  test.describe('no stored theme, OS prefers light', () => {
    test.use({ colorScheme: 'light' });

    test('shows the moon icon', async ({ page }) => {
      await page.goto('/');

      await expect(page.locator('html')).toHaveAttribute('data-theme', 'light');
      const { sun, moon } = toggleIcons(page);
      await expect(moon).toBeVisible();
      await expect(sun).toBeHidden();
    });
  });
});

test.describe('theme toggle icon with JavaScript disabled', () => {
  test.use({ javaScriptEnabled: false });

  // No script runs, so data-theme is never set and the page keeps the dark default.
  test('the default (no data-theme) shows the sun icon', async ({ page }) => {
    await page.goto('/');

    await expect(page.locator('html')).not.toHaveAttribute('data-theme', /.*/);
    const { sun, moon } = toggleIcons(page);
    await expect(sun).toBeVisible();
    await expect(moon).toBeHidden();
  });
});

test.describe('theme toggle icon after hydration', () => {
  for (const theme of ['dark', 'light']) {
    test(`${theme} theme: icon and label match, and toggling swaps both`, async ({ page }) => {
      await page.addInitScript((t) => localStorage.setItem('theme', t), theme);
      const errors = collectErrors(page);
      await page.goto('/');
      await waitForHydration(page);

      const isDark = theme === 'dark';
      const { sun, moon } = toggleIcons(page);
      const toggle = page.getByRole('button', { name: /Switch to (dark|light) mode/ }).first();
      await expect(toggle).toHaveAttribute('aria-label', isDark ? 'Switch to light mode' : 'Switch to dark mode');
      await expect(isDark ? sun : moon).toBeVisible();
      await expect(isDark ? moon : sun).toBeHidden();

      await toggle.click();
      await expect(page.locator('html')).toHaveAttribute('data-theme', isDark ? 'light' : 'dark');
      await expect(toggle).toHaveAttribute('aria-label', isDark ? 'Switch to dark mode' : 'Switch to light mode');
      await expect(isDark ? moon : sun).toBeVisible();
      await expect(isDark ? sun : moon).toBeHidden();

      expect(errors.filter((e) => HYDRATION_ERROR.test(e))).toEqual([]);
    });
  }
});
