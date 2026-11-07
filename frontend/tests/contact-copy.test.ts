import { experimental_AstroContainer as AstroContainer } from 'astro/container';
import { describe, expect, it } from 'vitest';
import Contact from '../src/components/Contact.astro';

async function renderContact() {
  const container = await AstroContainer.create();
  return container.renderToString(Contact);
}

function formMarkup(html: string) {
  const start = html.indexOf('<form');
  const end = html.indexOf('</form>');
  expect(start).toBeGreaterThan(-1);
  expect(end).toBeGreaterThan(start);
  return html.slice(start, end);
}

function count(haystack: string, needle: string) {
  return haystack.split(needle).length - 1;
}

const APPROVED_INTRO =
  "Whether you're planning a new system or trying to stabilize an existing one, tell me what you're working on. Send a message through the form, or connect with me on LinkedIn.";
const APPROVED_PLACEHOLDER =
  "Tell me about your project, your team, or the problem you're trying to solve...";

describe('P10: Contact section copy', () => {
  it('uses the approved business-inquiry intro', async () => {
    const html = await renderContact();

    expect(html).toContain(APPROVED_INTRO);
    expect(html).not.toContain('form below');
  });

  it('uses the approved message placeholder', async () => {
    const html = await renderContact();

    expect(html).toContain(`placeholder="${APPROVED_PLACEHOLDER}"`);
  });

  it('removes the old job-seeking copy', async () => {
    const html = await renderContact();

    expect(html).not.toContain('new opportunities');
    expect(html).not.toContain('fellow engineers');
    expect(html).not.toContain('just say hello');
  });

  it('keeps the Netlify form wiring and fields intact', async () => {
    const html = await renderContact();

    expect(html).toContain('name="contact"');
    expect(html).toContain('data-netlify="true"');
    expect(html).toContain('netlify-honeypot="bot-field"');
    expect(html).toContain('name="form-name" value="contact"');
    expect(html).toContain('action="/success"');
    expect(html).toContain('id="name"');
    expect(html).toContain('id="email"');
    expect(html).toContain('id="message"');
    expect(html).toContain('name="bot-field"');
  });

  it('keeps exactly the three visible fields plus the two hidden inputs', async () => {
    const form = formMarkup(await renderContact());

    // form-name (hidden), bot-field (honeypot), name, email
    expect(count(form, '<input')).toBe(4);
    expect(count(form, '<textarea')).toBe(1);
    expect(count(form, '<select')).toBe(0);
  });

  it('keeps the GitHub and LinkedIn links', async () => {
    const html = await renderContact();

    expect(html).toContain('href="https://github.com/rfabreu"');
    expect(html).toContain('href="https://linkedin.com/in/rafael-a-gomes"');
  });

  it('contains no pricing, budget, Upwork, or email address', async () => {
    const html = await renderContact();

    expect(html).not.toMatch(/\b(price|pricing|budget|quote|upwork)\b/i);
    expect(html).not.toContain('mailto:');
  });
});
