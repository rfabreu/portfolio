import { defineCollection, z, reference } from 'astro:content';
import { glob } from 'astro/loaders';

const projects = defineCollection({
  loader: glob({ pattern: '**/*.md', base: './src/content/projects' }),
  schema: z.object({
    title: z.string(),
    description: z.string(),
    image: z.string(),
    tags: z.array(z.string()),
    category: z.string(),
    featured: z.boolean().default(false),
    order: z.number(),
    github: z.string().url().optional(),
    demo: z.string().url().optional(),
  }),
});

const resume = defineCollection({
  loader: glob({ pattern: 'profile.md', base: './src/content/resume' }),
  schema: z.object({
    name: z.string(),
    title: z.string(),
    location: z.string(),
    links: z.object({
      linkedin: z.string().url(),
      github: z.string().url().optional(),
      portfolio: z.string().url().optional(),
    }),
    summary: z.string(),
    work_history: z.array(z.object({
      company: z.string(),
      title: z.string(),
      dates: z.string(),
      location: z.string().optional(),
      tech_stack: z.array(z.string()).default([]),
      bullets: z.array(z.string()).default([]),
    })),
    education: z.array(z.object({
      school: z.string(),
      program: z.string(),
      dates: z.string(),
    })),
    certifications: z.array(z.object({
      name: z.string(),
      issuer: z.string(),
      year: z.string(),
    })).default([]),
    skills: z.object({
      leadership: z.array(z.string()).default([]),
      technical: z.array(z.string()).default([]),
      domain: z.array(z.string()).default([]),
    }),
    awards_recognition: z.array(z.object({
      name: z.string(),
      org: z.string(),
      year: z.string(),
    })).default([]),
    speaking_writing: z.array(z.object({
      type: z.enum(['talk', 'article', 'podcast', 'panel']),
      title: z.string(),
      venue: z.string(),
      year: z.string(),
      url: z.string().url().optional(),
    })).default([]),
    languages_spoken: z.array(z.string()).default([]),
    selected_projects: z.array(z.object({
      slug: reference('projects'),
      headline: z.string(),
    })).default([]),
  }),
});

export const collections = { projects, resume };
