// frontend/scripts/generate-pdf.mjs
import { chromium } from 'playwright';
import { createServer } from 'http';
import handler from 'serve-handler';
import { resolve } from 'path';

const distPath = resolve('dist');
const outputPath = resolve('dist', 'Rafael_Abreu_Resume.pdf');

const server = createServer((req, res) => handler(req, res, { public: distPath }));
// Port 0 = OS-assigned free port. Avoids collision with Astro's default dev port (4321)
// and guarantees availability even if a previous build leaked sockets.
await new Promise((resolve, reject) => {
  server.on('error', reject);
  server.listen(0, resolve);
});
const { port } = server.address();

const browser = await chromium.launch();
try {
  const page = await browser.newPage();
  // page.goto resolves on navigation completion regardless of HTTP status —
  // explicit response.ok() check prevents PDFing a 404 page.
  const response = await page.goto(`http://localhost:${port}/resume`, { waitUntil: 'load' });
  if (!response || !response.ok()) {
    throw new Error(`Failed to load /resume for PDF generation: HTTP ${response?.status() ?? 'no response'}`);
  }
  await page.emulateMedia({ media: 'print' });
  await page.pdf({
    path: outputPath,
    format: 'A4',
    printBackground: true,
    margin: { top: '0.5in', right: '0.5in', bottom: '0.5in', left: '0.5in' },
  });
  console.log(`PDF written to ${outputPath}`);
} finally {
  await browser.close().catch((err) => console.error('browser.close failed:', err));
  server.close();
}
