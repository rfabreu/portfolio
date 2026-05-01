import { test, expect } from '@playwright/test';

test('Resume button in nav opens the modal with expected content', async ({ page }) => {
  await page.goto('/');
  await page.locator('a:has-text("RESUME")').first().click();

  const dialog = page.locator('dialog#resume-modal');
  await expect(dialog).toBeVisible();
  await expect(dialog.locator('.name')).toBeVisible();
  await expect(dialog.locator('text=Download PDF')).toBeVisible();
  await expect(dialog.locator('text=LinkedIn')).toBeVisible();
});

test('Download PDF link points to the static PDF asset', async ({ page }) => {
  await page.goto('/');
  await page.locator('a:has-text("RESUME")').first().click();
  const downloadHref = await page.locator('.btn-download').getAttribute('href');
  expect(downloadHref).toBe('/Rafael_Abreu_Resume.pdf');
});

test('Resume button is hidden on /resume page', async ({ page }) => {
  await page.goto('/resume');
  await expect(page.locator('a:has-text("RESUME")')).toHaveCount(0);
});

test('/resume page renders resume content directly', async ({ page }) => {
  await page.goto('/resume');
  await expect(page.locator('h1.resume-page-name')).toBeVisible();
});
