import { experimental_AstroContainer as AstroContainer } from 'astro/container';
import { readFileSync } from 'node:fs';
import { describe, expect, it } from 'vitest';
import Contact from '../src/components/Contact.astro';

// Netlify builds notification field names from <label> text, so label text
// must be clean. See docs/superpowers/specs/2026-11-18-contact-notification-labels-design.md.

async function renderForm() {
  const container = await AstroContainer.create();
  const html = await container.renderToString(Contact);
  return html.slice(html.indexOf('<form'), html.indexOf('</form>'));
}

const FIELDS = [
  { id: 'name', name: 'name', label: 'Name', required: true },
  { id: 'email', name: 'email', label: 'Email', required: true },
  { id: 'project_type', name: 'project_type', label: 'Project type', required: true },
  { id: 'timeline', name: 'timeline', label: 'Target timeline', required: false },
  { id: 'message', name: 'message', label: 'Project brief', required: true },
];

function labelFor(form: string, id: string) {
  const match = form.match(new RegExp(`<label[^>]*for="${id}"[^>]*>([\\s\\S]*?)</label>`));
  expect(match, `label for ${id}`).not.toBeNull();
  return { attrs: match![0].slice(0, match![0].indexOf('>')), text: match![1] };
}

function controlTag(form: string, id: string) {
  const match = form.match(new RegExp(`<(?:input|select|textarea)[^>]*id="${id}"[^>]*>`));
  expect(match, `control ${id}`).not.toBeNull();
  return match![0];
}

describe('Contact form notification labels', () => {
  for (const f of FIELDS) {
    it(`${f.id}: label text is clean, name attribute unchanged`, async () => {
      const form = await renderForm();
      const label = labelFor(form, f.id);

      expect(label.text.trim()).toBe(f.label);
      expect(label.text).not.toContain('*');
      expect(label.text).not.toMatch(/<[a-z]/i);
      expect(controlTag(form, f.id)).toContain(`name="${f.name}"`);
    });

    it(`${f.id}: required marker and required attribute`, async () => {
      const form = await renderForm();
      const label = labelFor(form, f.id);
      const control = controlTag(form, f.id);

      if (f.required) {
        expect(label.attrs).toContain('field-required');
        expect(control).toMatch(/\srequired[\s>=]/);
      } else {
        expect(label.attrs).not.toContain('field-required');
        expect(control).not.toMatch(/\srequired[\s>=]/);
      }
    });

    it(`${f.id}: label keeps uppercase look via CSS`, async () => {
      const label = labelFor(await renderForm(), f.id);
      expect(label.attrs).toContain('uppercase');
    });
  }

  it('leaves the honeypot label and field untouched', async () => {
    const form = await renderForm();
    expect(form).toContain("<label>Don't fill this out: <input name=\"bot-field\"");
  });

  it('draws the asterisk as a CSS pseudo-element with the old " *" spacing', () => {
    const css = readFileSync(new URL('../src/styles/global.css', import.meta.url), 'utf8');
    expect(css).toMatch(/\.field-required::after\s*\{\s*content:\s*" \*";\s*\}/);
  });
});
