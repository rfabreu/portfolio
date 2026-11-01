import { experimental_AstroContainer as AstroContainer } from 'astro/container';
import { describe, expect, it } from 'vitest';
import TestimonialCard from '../src/components/TestimonialCard.astro';
import Testimonials from '../src/components/Testimonials.astro';

describe('TestimonialCard', () => {
  it('renders author details, quote, and monogram initials fallback when no avatar provided', async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(TestimonialCard, {
      props: {
        name: 'Alex Rivera',
        role: 'Engineering Manager',
        company: 'Nextologies',
        relationship: 'Managed Rafael directly',
        quote: 'Rafael is an exceptional full-stack engineer.',
      },
    });

    expect(html).toContain('Alex Rivera');
    expect(html).toContain('Engineering Manager · Nextologies');
    expect(html).toContain('Managed Rafael directly');
    expect(html).toContain('Rafael is an exceptional full-stack engineer.');
    // Initials monogram fallback for "Alex Rivera" is "AR"
    expect(html).toContain('AR');
    expect(html).not.toContain('<img');
  });

  it('renders avatar image when avatar is provided', async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(TestimonialCard, {
      props: {
        name: 'Jane Doe',
        role: 'Product Lead',
        company: 'Tech Corp',
        relationship: 'Collaborated on projects',
        avatar: '/images/testimonials/jane.jpg',
        quote: 'Delivered features ahead of schedule.',
      },
    });

    expect(html).toContain('<img');
    expect(html).toContain('src="/images/testimonials/jane.jpg"');
    expect(html).toContain('alt="Jane Doe"');
  });

  it('renders linkedin link when linkedin url is provided', async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(TestimonialCard, {
      props: {
        name: 'Jane Doe',
        role: 'Product Lead',
        company: 'Tech Corp',
        relationship: 'Collaborated on projects',
        linkedin: 'https://linkedin.com/in/janedoe',
        quote: 'Great teammate.',
      },
    });

    expect(html).toContain('href="https://linkedin.com/in/janedoe"');
    expect(html).toContain('aria-label="Jane Doe\'s LinkedIn profile"');
  });
});

describe('Testimonials Section', () => {
  it('renders nothing when content collection is empty (dormant guard)', async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(Testimonials);

    // Section should be completely empty / whitespace only
    expect(html.trim()).toBe('');
    expect(html).not.toContain('<section id="testimonials"');
  });
});
