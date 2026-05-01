// frontend/scripts/generate-pdf.mjs
import { chromium } from 'playwright';
import { createServer } from 'http';
import handler from 'serve-handler';
import { resolve } from 'path';

const distPath = resolve('dist');
const outputPath = resolve('dist', 'Rafael_Abreu_Resume.pdf');

const server = createServer((req, res) => handler(req, res, { public: distPath }));
// Port 0 = OS-assigned free port. Avoids collision with Astro's default dev port (4321).
await new Promise((resolve) => server.listen(0, resolve));
const { port } = server.address();

const browser = await chromium.launch();
try {
  const page = await browser.newPage();
  await page.goto(`http://localhost:${port}/resume`, { waitUntil: 'load' });
  await page.emulateMedia({ media: 'print' });
  await page.pdf({
    path: outputPath,
    format: 'A4',
    printBackground: true,
    margin: { top: '0.5in', right: '0.5in', bottom: '0.5in', left: '0.5in' },
  });
  console.log(`PDF written to ${outputPath}`);
} finally {
  await browser.close();
  server.close();
}
