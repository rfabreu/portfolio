import { experimental_AstroContainer as AstroContainer } from 'astro/container';
import { describe, expect, it } from 'vitest';
import Contact from '../src/components/Contact.astro';

async function renderContact() {
  const container = await AstroContainer.create();
  return container.renderToString(Contact);
}

describe('P10: Contact section copy', () => {
  it('uses the approved business-inquiry intro', async () => {
    const html = await renderContact();

    expect(html).toContain('planning a new system or trying to stabilize an existing one');
    expect(html).toContain('Send a message through the form below, or connect with me on LinkedIn.');
  });

  it('uses the approved message placeholder', async () => {
    const html = await renderContact();

    expect(html).toContain('Tell me about your project, your team, or the problem');
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
