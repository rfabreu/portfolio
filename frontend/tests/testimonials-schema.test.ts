import { describe, expect, it } from 'vitest';
import { testimonialSchema } from '../src/content.config';

describe('testimonialSchema', () => {
  it('validates a complete valid testimonial payload', () => {
    const input = {
      name: 'Alex Rivera',
      role: 'Engineering Manager',
      company: 'Nextologies',
      relationship: 'Managed Rafael directly',
      avatar: '/images/testimonials/alex.jpg',
      linkedin: 'https://www.linkedin.com/in/alex-rivera',
      order: 1,
      featured: true,
    };

    const parsed = testimonialSchema.parse(input);
    expect(parsed.name).toBe('Alex Rivera');
    expect(parsed.role).toBe('Engineering Manager');
    expect(parsed.company).toBe('Nextologies');
    expect(parsed.relationship).toBe('Managed Rafael directly');
    expect(parsed.avatar).toBe('/images/testimonials/alex.jpg');
    expect(parsed.linkedin).toBe('https://www.linkedin.com/in/alex-rivera');
    expect(parsed.order).toBe(1);
    expect(parsed.featured).toBe(true);
  });

  it('applies defaults for order and featured when omitted', () => {
    const input = {
      name: 'Jane Smith',
      role: 'Staff Engineer',
      company: 'Acme Corp',
      relationship: 'Colleague',
    };

    const parsed = testimonialSchema.parse(input);
    expect(parsed.order).toBe(99);
    expect(parsed.featured).toBe(false);
    expect(parsed.avatar).toBeUndefined();
    expect(parsed.linkedin).toBeUndefined();
  });

  it('rejects payload when required fields are missing', () => {
    const missingName = {
      role: 'Engineer',
      company: 'Tech Co',
      relationship: 'Peer',
    };
    expect(() => testimonialSchema.parse(missingName)).toThrow();

    const missingRole = {
      name: 'Jane Smith',
      company: 'Tech Co',
      relationship: 'Peer',
    };
    expect(() => testimonialSchema.parse(missingRole)).toThrow();

    const missingCompany = {
      name: 'Jane Smith',
      role: 'Engineer',
      relationship: 'Peer',
    };
    expect(() => testimonialSchema.parse(missingCompany)).toThrow();

    const missingRelationship = {
      name: 'Jane Smith',
      role: 'Engineer',
      company: 'Tech Co',
    };
    expect(() => testimonialSchema.parse(missingRelationship)).toThrow();
  });

  it('rejects invalid linkedin URLs', () => {
    const invalidUrl = {
      name: 'Jane Smith',
      role: 'Engineer',
      company: 'Tech Co',
      relationship: 'Peer',
      linkedin: 'not-a-valid-url',
    };
    expect(() => testimonialSchema.parse(invalidUrl)).toThrow();
  });
});
