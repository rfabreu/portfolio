# P9: Resume Feature — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the broken Resume button (404 on `/Rafael_Abreu_Resume.pdf`) with a markdown-driven resume that renders as a styled in-site modal, downloads as a build-time PDF, and grounds the chatbot's system prompt — all from a single `profile.md` source of truth.

**Architecture:** One markdown file at `frontend/src/content/resume/profile.md` feeds three deterministic renderers: (1) an Astro modal using the native HTML `<dialog>` element, (2) a Playwright-generated PDF written to `frontend/dist/Rafael_Abreu_Resume.pdf` in CI, (3) a Go-embedded system prompt assembled from the resume plus existing project case studies. No AI in the build or request path.

**Tech Stack:** Astro 6, native HTML `<dialog>`, Playwright (PDF + E2E), Vitest (unit/snapshot tests), Go 1.26+ with `//go:embed`, GitHub Actions CI.

**Spec:** `docs/superpowers/specs/2026-04-30-resume-feature-design.md`
**Issue:** #25
**Branch:** `feature/resume`

**Reference for initial content (DO NOT copy or commit):** `/Users/rafaeldeabreugomes/Downloads/rafael_abreu_linkedin_profile.md`

---

## Phase 1: Foundation — Schema & Initial Content

### Task 1.1: Add `resume` collection schema

**Files:**
- Modify: `frontend/src/content.config.ts`

- [ ] **Step 1: Replace the file contents with the resume collection added alongside `projects`:**

```typescript
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
```

- [ ] **Step 2: Verify schema parses (no profile.md yet, so collection will be empty):** Run `cd frontend && npx astro check`. Expected: no schema errors. (Warnings about empty collection are acceptable.)

- [ ] **Step 3: Commit:**

```bash
git add frontend/src/content.config.ts
git commit -m "feat(p9): add resume content collection schema

Defines the structured schema for profile.md per the design spec —
required fields (name, title, location, links.linkedin, summary,
work_history, education, skills) plus optional sections (certifications,
awards_recognition, speaking_writing, languages_spoken, selected_projects).
Uses reference('projects') for slug validation against the existing
projects collection."
```

### Task 1.2: Create `profile.md` from LinkedIn export

**Files:**
- Create: `frontend/src/content/resume/profile.md`
- Reference (read-only, never copied/committed): `/Users/rafaeldeabreugomes/Downloads/rafael_abreu_linkedin_profile.md`

- [ ] **Step 1: Read the LinkedIn export** to extract the source data:

```bash
cat /Users/rafaeldeabreugomes/Downloads/rafael_abreu_linkedin_profile.md
```

Extract: name, current and past job titles with dates, descriptions, education entries, certifications, languages.

- [ ] **Step 2: Verify the export file is NOT in the repo and won't be:**

```bash
git ls-files | grep -i linkedin
```

Expected: no output. The file lives only at `~/Downloads/`.

- [ ] **Step 3: Write `frontend/src/content/resume/profile.md`** with frontmatter populated from the LinkedIn data. Apply these rules:

  - **Tone:** exec-track — verbs like "Led", "Architected", "Drove", "Owned", "Established", "Scaled". Avoid "helped with", "contributed to", "learned".
  - **Accuracy:** dates, titles, employers, certifications come verbatim from the LinkedIn export. Bullets paraphrase the LinkedIn descriptions but never invent metrics, scope, or roles not in the source.
  - **`tech_stack` per work entry:** drawn from technologies mentioned in that job's LinkedIn description. Do not pad with technologies he uses generally.
  - **`leadership` skills bucket:** populated from leadership signals in LinkedIn job descriptions ("led", "managed", "owned"). Skip if not evidenced.
  - **`awards_recognition` and `speaking_writing`:** start as empty arrays. They will fill over time.
  - **`languages_spoken`:** "English (Native)", "Portuguese (Native)" at minimum; check LinkedIn for others.
  - **`selected_projects`:** reference 2–4 strongest case studies from `frontend/src/content/projects/`. Use slugs that exist in that directory.

Template structure (replace placeholder values with real LinkedIn-derived content):

```markdown
---
name: "Rafael Abreu"
title: "Software Engineer | Engineering Leadership Track"
location: "Toronto, ON, Canada"

links:
  linkedin:  "https://www.linkedin.com/in/rafael-abreu"
  github:    "https://github.com/rfabreu"
  portfolio: "https://rafaelabreu.netlify.app"

summary: |
  [Three-sentence executive-positioned summary. Who you are, what
  you've delivered, where you're heading. Drawn from LinkedIn's
  About section, sharpened for senior/director/C-level positioning.]

work_history:
  - company:    "[Verbatim from LinkedIn]"
    title:      "[Verbatim from LinkedIn]"
    dates:      "[Verbatim from LinkedIn, format: '2023 — Present']"
    location:   "[Verbatim from LinkedIn]"
    tech_stack: ["...", "..."]
    bullets:
      - "[Exec-tone paraphrase of LinkedIn bullet — never invents]"
      - "..."
  # repeat for each role

education:
  - school:  "[Verbatim]"
    program: "[Verbatim]"
    dates:   "[Verbatim]"

certifications:
  - name:   "[Verbatim certification name]"
    issuer: "[Verbatim issuer]"
    year:   "[Verbatim year]"

skills:
  leadership:
    - "[Skills evidenced in LinkedIn job descriptions]"
  technical:
    - "[Tech skills consolidated from LinkedIn 'Skills' section + job descriptions]"
  domain:
    - "[Industry/domain expertise from LinkedIn]"

awards_recognition: []
speaking_writing:   []

languages_spoken:
  - "English (Native)"
  - "Portuguese (Native)"

selected_projects:
  - slug:     "[existing project slug from frontend/src/content/projects/]"
    headline: "[Resume-specific one-line framing]"
---
```

- [ ] **Step 4: Validate:** `cd frontend && npx astro check`. Expected: no errors. If `selected_projects.slug` references a non-existent project, the build fails with a useful pointer — fix the slug before continuing.

- [ ] **Step 5: Commit:**

```bash
git add frontend/src/content/resume/profile.md
git commit -m "feat(p9): add initial profile.md derived from LinkedIn export

Source data lifted from local LinkedIn export (kept off-repo per
no-personal-data policy). Tone tuned for exec-track positioning;
all dates, titles, and employers verbatim from source."
```

---

## Phase 2: Section Components

These render the structured data into HTML. Each uses Tailwind v4 theme tokens (`text-text-primary`, `bg-surface`, etc.) so dark/light themes work without per-component overrides.

### Task 2.1: Add Vitest for snapshot testing

**Files:**
- Modify: `frontend/package.json`
- Create: `frontend/vitest.config.ts`
- Create: `frontend/tests/setup.ts`

- [ ] **Step 1: Install Vitest:**

```bash
cd frontend
npm install --save-dev vitest
```

Note: the Astro Container API used in component tests is exported from `astro/container` (already a dependency); no extra adapter package is required.

- [ ] **Step 2: Add `test` script to `package.json` "scripts":**

```json
"test": "vitest run",
"test:watch": "vitest"
```

- [ ] **Step 3: Create `frontend/vitest.config.ts`:**

```typescript
import { defineConfig } from 'vitest/config';
import { getViteConfig } from 'astro/config';

export default getViteConfig(
  defineConfig({
    test: {
      globals: true,
      environment: 'node',
      include: ['tests/**/*.test.ts'],
    },
  }),
);
```

- [ ] **Step 4: Verify the test runner starts:**

```bash
cd frontend && npm test
```

Expected: "No test files found" or similar — runner works, no tests yet. PASS.

- [ ] **Step 5: Commit:**

```bash
git add frontend/package.json frontend/package-lock.json frontend/vitest.config.ts
git commit -m "chore(p9): add Vitest for snapshot testing"
```

### Task 2.2: `Summary.astro` — summary paragraph

**Files:**
- Create: `frontend/src/components/resume/Summary.astro`
- Create: `frontend/tests/resume-summary.test.ts`

- [ ] **Step 1: Write the failing test:**

```typescript
// frontend/tests/resume-summary.test.ts
import { experimental_AstroContainer as AstroContainer } from 'astro/container';
import { expect, test } from 'vitest';
import Summary from '../src/components/resume/Summary.astro';

test('Summary renders the summary text in a section element', async () => {
  const container = await AstroContainer.create();
  const html = await container.renderToString(Summary, {
    props: { summary: 'Engineering leader building distributed systems.' },
  });
  expect(html).toContain('<section');
  expect(html).toContain('Engineering leader building distributed systems.');
});
```

- [ ] **Step 2: Run the test, expect failure:**

```bash
cd frontend && npm test -- resume-summary
```

Expected: FAIL — module not found.

- [ ] **Step 3: Create the component:**

```astro
---
// frontend/src/components/resume/Summary.astro
interface Props {
  summary: string;
}
const { summary } = Astro.props;
---

<section class="resume-section resume-summary" aria-label="Summary">
  <p class="text-text-secondary leading-relaxed">{summary}</p>
</section>
```

- [ ] **Step 4: Run the test, expect pass:**

```bash
cd frontend && npm test -- resume-summary
```

Expected: PASS.

- [ ] **Step 5: Commit:**

```bash
git add frontend/src/components/resume/Summary.astro frontend/tests/resume-summary.test.ts
git commit -m "feat(p9): add Summary section component"
```

### Task 2.3: `Skills.astro` — bucketed skill lists

**Files:**
- Create: `frontend/src/components/resume/Skills.astro`
- Create: `frontend/tests/resume-skills.test.ts`

- [ ] **Step 1: Write the failing test:**

```typescript
import { experimental_AstroContainer as AstroContainer } from 'astro/container';
import { expect, test } from 'vitest';
import Skills from '../src/components/resume/Skills.astro';

test('Skills renders each non-empty bucket with its label and items', async () => {
  const container = await AstroContainer.create();
  const html = await container.renderToString(Skills, {
    props: {
      skills: {
        leadership: ['Strategic Planning'],
        technical: ['Go', 'Python'],
        domain: [],
      },
    },
  });
  expect(html).toContain('Leadership');
  expect(html).toContain('Strategic Planning');
  expect(html).toContain('Technical');
  expect(html).toContain('Go');
  expect(html).toContain('Python');
  expect(html).not.toContain('Domain');  // empty bucket skipped
});
```

- [ ] **Step 2: Run, expect FAIL** (`cd frontend && npm test -- resume-skills`).

- [ ] **Step 3: Create the component:**

```astro
---
// frontend/src/components/resume/Skills.astro
interface Props {
  skills: {
    leadership: string[];
    technical: string[];
    domain: string[];
  };
}
const { skills } = Astro.props;

const buckets: Array<{ key: keyof typeof skills; label: string }> = [
  { key: 'leadership', label: 'Leadership' },
  { key: 'technical', label: 'Technical' },
  { key: 'domain', label: 'Domain' },
];
---

<section class="resume-section resume-skills" aria-label="Skills">
  <h2 class="resume-heading">Skills</h2>
  {buckets.map(({ key, label }) => (
    skills[key].length > 0 && (
      <div class="resume-skill-bucket">
        <h3 class="resume-skill-label">{label}</h3>
        <ul class="resume-skill-list">
          {skills[key].map((skill) => <li>{skill}</li>)}
        </ul>
      </div>
    )
  ))}
</section>
```

- [ ] **Step 4: Run, expect PASS.**

- [ ] **Step 5: Commit:**

```bash
git add frontend/src/components/resume/Skills.astro frontend/tests/resume-skills.test.ts
git commit -m "feat(p9): add Skills section component"
```

### Task 2.4: `WorkHistory.astro` — jobs with bullets and tech_stack

**Files:**
- Create: `frontend/src/components/resume/WorkHistory.astro`
- Create: `frontend/tests/resume-work-history.test.ts`

- [ ] **Step 1: Failing test:**

```typescript
import { experimental_AstroContainer as AstroContainer } from 'astro/container';
import { expect, test } from 'vitest';
import WorkHistory from '../src/components/resume/WorkHistory.astro';

test('WorkHistory renders each job with title, company, dates, bullets, tech_stack', async () => {
  const container = await AstroContainer.create();
  const html = await container.renderToString(WorkHistory, {
    props: {
      entries: [
        {
          company: 'Nextologies',
          title: 'Software Engineer',
          dates: '2023 — Present',
          location: 'Toronto, ON',
          tech_stack: ['Go', 'Kubernetes'],
          bullets: ['Architected the platform.'],
        },
      ],
    },
  });
  expect(html).toContain('Nextologies');
  expect(html).toContain('Software Engineer');
  expect(html).toContain('2023 — Present');
  expect(html).toContain('Toronto, ON');
  expect(html).toContain('Architected the platform.');
  expect(html).toContain('Go');
  expect(html).toContain('Kubernetes');
});

test('WorkHistory handles empty bullets and missing location gracefully', async () => {
  const container = await AstroContainer.create();
  const html = await container.renderToString(WorkHistory, {
    props: {
      entries: [
        {
          company: 'Acme',
          title: 'Engineer',
          dates: '2020 — 2022',
          tech_stack: [],
          bullets: [],
        },
      ],
    },
  });
  expect(html).toContain('Acme');
  expect(html).not.toContain('<ul');  // no bullet list when empty
});
```

- [ ] **Step 2: Run, expect FAIL.**

- [ ] **Step 3: Create the component:**

```astro
---
// frontend/src/components/resume/WorkHistory.astro
interface Entry {
  company: string;
  title: string;
  dates: string;
  location?: string;
  tech_stack: string[];
  bullets: string[];
}

interface Props {
  entries: Entry[];
}

const { entries } = Astro.props;
---

<section class="resume-section resume-work-history" aria-label="Work history">
  <h2 class="resume-heading">Experience</h2>
  {entries.map((entry) => (
    <article class="resume-job">
      <header class="resume-job-header">
        <div class="resume-job-title-row">
          <span class="resume-job-title">{entry.title}</span>
          <span class="resume-job-company">{entry.company}</span>
        </div>
        <div class="resume-job-meta">
          <span class="resume-job-dates">{entry.dates}</span>
          {entry.location && <span class="resume-job-location"> · {entry.location}</span>}
        </div>
      </header>
      {entry.bullets.length > 0 && (
        <ul class="resume-job-bullets">
          {entry.bullets.map((b) => <li>{b}</li>)}
        </ul>
      )}
      {entry.tech_stack.length > 0 && (
        <div class="resume-job-tech">
          {entry.tech_stack.map((tech) => <span class="resume-tech-tag">{tech}</span>)}
        </div>
      )}
    </article>
  ))}
</section>
```

- [ ] **Step 4: Run, expect PASS.**

- [ ] **Step 5: Commit:**

```bash
git add frontend/src/components/resume/WorkHistory.astro frontend/tests/resume-work-history.test.ts
git commit -m "feat(p9): add WorkHistory section component"
```

### Task 2.5: `Education.astro` — education entries

**Files:**
- Create: `frontend/src/components/resume/Education.astro`
- Create: `frontend/tests/resume-education.test.ts`

- [ ] **Step 1: Failing test:**

```typescript
import { experimental_AstroContainer as AstroContainer } from 'astro/container';
import { expect, test } from 'vitest';
import Education from '../src/components/resume/Education.astro';

test('Education renders school, program, and dates per entry', async () => {
  const container = await AstroContainer.create();
  const html = await container.renderToString(Education, {
    props: {
      entries: [
        { school: 'University of Toronto', program: 'Coding Bootcamp', dates: '2022 — 2023' },
      ],
    },
  });
  expect(html).toContain('University of Toronto');
  expect(html).toContain('Coding Bootcamp');
  expect(html).toContain('2022 — 2023');
});
```

- [ ] **Step 2: Run, expect FAIL.**

- [ ] **Step 3: Create the component:**

```astro
---
// frontend/src/components/resume/Education.astro
interface Entry {
  school: string;
  program: string;
  dates: string;
}
interface Props {
  entries: Entry[];
}
const { entries } = Astro.props;
---

{entries.length > 0 && (
  <section class="resume-section resume-education" aria-label="Education">
    <h2 class="resume-heading">Education</h2>
    <ul class="resume-list">
      {entries.map((entry) => (
        <li>
          <span class="resume-edu-school">{entry.school}</span>
          <span class="resume-edu-program"> — {entry.program}</span>
          <span class="resume-edu-dates"> ({entry.dates})</span>
        </li>
      ))}
    </ul>
  </section>
)}
```

- [ ] **Step 4: Run, expect PASS.**

- [ ] **Step 5: Commit:**

```bash
git add frontend/src/components/resume/Education.astro frontend/tests/resume-education.test.ts
git commit -m "feat(p9): add Education section component"
```

### Task 2.6: `Certifications.astro`, `Awards.astro`, `Speaking.astro` — list-shaped sections

**Files:**
- Create: `frontend/src/components/resume/Certifications.astro`
- Create: `frontend/src/components/resume/Awards.astro`
- Create: `frontend/src/components/resume/Speaking.astro`
- Create: `frontend/tests/resume-list-sections.test.ts`

These three sections are structurally similar — each renders a list of records and skips itself when empty.

- [ ] **Step 1: Failing test (covers all three components):**

```typescript
import { experimental_AstroContainer as AstroContainer } from 'astro/container';
import { expect, test } from 'vitest';
import Certifications from '../src/components/resume/Certifications.astro';
import Awards from '../src/components/resume/Awards.astro';
import Speaking from '../src/components/resume/Speaking.astro';

test('Certifications renders entries with name/issuer/year, omits when empty', async () => {
  const container = await AstroContainer.create();
  const withEntries = await container.renderToString(Certifications, {
    props: { entries: [{ name: 'OpenShift Specialist', issuer: 'Red Hat', year: '2024' }] },
  });
  expect(withEntries).toContain('OpenShift Specialist');
  expect(withEntries).toContain('Red Hat');
  expect(withEntries).toContain('2024');

  const empty = await container.renderToString(Certifications, { props: { entries: [] } });
  expect(empty.trim()).toBe('');
});

test('Awards renders entries with name/org/year, omits when empty', async () => {
  const container = await AstroContainer.create();
  const withEntries = await container.renderToString(Awards, {
    props: { entries: [{ name: 'Hackathon Winner', org: 'TechFest', year: '2025' }] },
  });
  expect(withEntries).toContain('Hackathon Winner');
  expect(withEntries).toContain('TechFest');

  const empty = await container.renderToString(Awards, { props: { entries: [] } });
  expect(empty.trim()).toBe('');
});

test('Speaking renders entries with title/venue/year and link when present, omits when empty', async () => {
  const container = await AstroContainer.create();
  const withEntries = await container.renderToString(Speaking, {
    props: {
      entries: [
        { type: 'talk', title: 'Building with Go', venue: 'GoTO Meetup', year: '2026', url: 'https://example.com/talk' },
      ],
    },
  });
  expect(withEntries).toContain('Building with Go');
  expect(withEntries).toContain('GoTO Meetup');
  expect(withEntries).toContain('href="https://example.com/talk"');

  const empty = await container.renderToString(Speaking, { props: { entries: [] } });
  expect(empty.trim()).toBe('');
});
```

- [ ] **Step 2: Run, expect FAIL.**

- [ ] **Step 3: Create `Certifications.astro`:**

```astro
---
interface Entry { name: string; issuer: string; year: string; }
interface Props { entries: Entry[]; }
const { entries } = Astro.props;
---

{entries.length > 0 && (
  <section class="resume-section resume-certifications" aria-label="Certifications">
    <h2 class="resume-heading">Certifications</h2>
    <ul class="resume-list">
      {entries.map((entry) => (
        <li>
          <span class="resume-cert-name">{entry.name}</span>
          <span class="resume-cert-meta"> — {entry.issuer} ({entry.year})</span>
        </li>
      ))}
    </ul>
  </section>
)}
```

- [ ] **Step 4: Create `Awards.astro`:**

```astro
---
interface Entry { name: string; org: string; year: string; }
interface Props { entries: Entry[]; }
const { entries } = Astro.props;
---

{entries.length > 0 && (
  <section class="resume-section resume-awards" aria-label="Awards & Recognition">
    <h2 class="resume-heading">Awards &amp; Recognition</h2>
    <ul class="resume-list">
      {entries.map((entry) => (
        <li>
          <span class="resume-award-name">{entry.name}</span>
          <span class="resume-award-meta"> — {entry.org} ({entry.year})</span>
        </li>
      ))}
    </ul>
  </section>
)}
```

- [ ] **Step 5: Create `Speaking.astro`:**

```astro
---
interface Entry {
  type: 'talk' | 'article' | 'podcast' | 'panel';
  title: string;
  venue: string;
  year: string;
  url?: string;
}
interface Props { entries: Entry[]; }
const { entries } = Astro.props;
---

{entries.length > 0 && (
  <section class="resume-section resume-speaking" aria-label="Speaking & Writing">
    <h2 class="resume-heading">Speaking &amp; Writing</h2>
    <ul class="resume-list">
      {entries.map((entry) => (
        <li>
          {entry.url ? (
            <a href={entry.url} class="resume-speaking-link">{entry.title}</a>
          ) : (
            <span>{entry.title}</span>
          )}
          <span class="resume-speaking-meta"> — {entry.venue} ({entry.year})</span>
        </li>
      ))}
    </ul>
  </section>
)}
```

- [ ] **Step 6: Run, expect PASS.**

- [ ] **Step 7: Commit:**

```bash
git add frontend/src/components/resume/Certifications.astro frontend/src/components/resume/Awards.astro frontend/src/components/resume/Speaking.astro frontend/tests/resume-list-sections.test.ts
git commit -m "feat(p9): add Certifications, Awards, and Speaking section components"
```

### Task 2.7: `Languages.astro` — languages spoken

**Files:**
- Create: `frontend/src/components/resume/Languages.astro`
- Create: `frontend/tests/resume-languages.test.ts`

- [ ] **Step 1: Failing test:**

```typescript
import { experimental_AstroContainer as AstroContainer } from 'astro/container';
import { expect, test } from 'vitest';
import Languages from '../src/components/resume/Languages.astro';

test('Languages renders each entry inline; omits when empty', async () => {
  const container = await AstroContainer.create();
  const withEntries = await container.renderToString(Languages, {
    props: { entries: ['English (Native)', 'Portuguese (Native)'] },
  });
  expect(withEntries).toContain('English (Native)');
  expect(withEntries).toContain('Portuguese (Native)');

  const empty = await container.renderToString(Languages, { props: { entries: [] } });
  expect(empty.trim()).toBe('');
});
```

- [ ] **Step 2: Run, expect FAIL.**

- [ ] **Step 3: Create the component:**

```astro
---
interface Props { entries: string[]; }
const { entries } = Astro.props;
---

{entries.length > 0 && (
  <section class="resume-section resume-languages" aria-label="Languages">
    <h2 class="resume-heading">Languages</h2>
    <p class="resume-languages-list">{entries.join(' · ')}</p>
  </section>
)}
```

- [ ] **Step 4: Run, expect PASS.**

- [ ] **Step 5: Commit:**

```bash
git add frontend/src/components/resume/Languages.astro frontend/tests/resume-languages.test.ts
git commit -m "feat(p9): add Languages section component"
```

### Task 2.8: `Projects.astro` — selected projects with case-study links

**Files:**
- Create: `frontend/src/components/resume/Projects.astro`
- Create: `frontend/tests/resume-projects.test.ts`

- [ ] **Step 1: Failing test:**

```typescript
import { experimental_AstroContainer as AstroContainer } from 'astro/container';
import { expect, test } from 'vitest';
import Projects from '../src/components/resume/Projects.astro';

test('Projects renders headline and links to case study by slug', async () => {
  const container = await AstroContainer.create();
  const html = await container.renderToString(Projects, {
    props: {
      entries: [
        { slug: 'fptv', headline: 'Glassmorphism dashboard for Toronto TV broadcaster' },
      ],
    },
  });
  expect(html).toContain('Glassmorphism dashboard for Toronto TV broadcaster');
  expect(html).toContain('href="/project/fptv"');
});

test('Projects omits section when entries is empty', async () => {
  const container = await AstroContainer.create();
  const html = await container.renderToString(Projects, { props: { entries: [] } });
  expect(html.trim()).toBe('');
});
```

- [ ] **Step 2: Run, expect FAIL.**

- [ ] **Step 3: Create the component:**

```astro
---
interface Entry { slug: string; headline: string; }
interface Props { entries: Entry[]; }
const { entries } = Astro.props;
---

{entries.length > 0 && (
  <section class="resume-section resume-projects" aria-label="Selected projects">
    <h2 class="resume-heading">Selected Projects</h2>
    <ul class="resume-list">
      {entries.map((entry) => (
        <li>
          <a href={`/project/${entry.slug}`} class="resume-project-link">
            {entry.headline}
          </a>
        </li>
      ))}
    </ul>
  </section>
)}
```

Note: `entry.slug` here is already the resolved string from the `reference()` helper because we'll convert the `DataEntryMap` shape to plain props in `ResumeBody.astro` (Task 3.1). Astro's `reference()` validates at build time but the runtime value is `{ collection, id }`, so the assembler will call `.id` to extract the slug string before passing to this component.

- [ ] **Step 4: Run, expect PASS.**

- [ ] **Step 5: Commit:**

```bash
git add frontend/src/components/resume/Projects.astro frontend/tests/resume-projects.test.ts
git commit -m "feat(p9): add Projects section component"
```

---

## Phase 3: Wrappers — `ResumeBody.astro` and `/resume` page

### Task 3.1: `ResumeBody.astro` — assembles all sections

**Files:**
- Create: `frontend/src/components/resume/ResumeBody.astro`

- [ ] **Step 1: Create the assembler component:**

```astro
---
// frontend/src/components/resume/ResumeBody.astro
import type { CollectionEntry } from 'astro:content';
import Summary from './Summary.astro';
import Skills from './Skills.astro';
import WorkHistory from './WorkHistory.astro';
import Education from './Education.astro';
import Certifications from './Certifications.astro';
import Awards from './Awards.astro';
import Speaking from './Speaking.astro';
import Languages from './Languages.astro';
import Projects from './Projects.astro';

interface Props {
  resume: CollectionEntry<'resume'>;
}
const { resume } = Astro.props;
const data = resume.data;

const projectEntries = data.selected_projects.map((p) => ({
  slug: p.slug.id,
  headline: p.headline,
}));
---

<div class="resume-body">
  <Summary summary={data.summary} />
  <Skills skills={data.skills} />
  <WorkHistory entries={data.work_history} />
  <Education entries={data.education} />
  <Certifications entries={data.certifications} />
  <Awards entries={data.awards_recognition} />
  <Speaking entries={data.speaking_writing} />
  <Languages entries={data.languages_spoken} />
  <Projects entries={projectEntries} />
</div>
```

- [ ] **Step 2: Verify type-check:** `cd frontend && npx astro check`. Expected: no errors.

- [ ] **Step 3: Commit:**

```bash
git add frontend/src/components/resume/ResumeBody.astro
git commit -m "feat(p9): add ResumeBody assembler"
```

### Task 3.2: `/resume` standalone page

**Files:**
- Create: `frontend/src/pages/resume.astro`

- [ ] **Step 1: Create the page:**

```astro
---
// frontend/src/pages/resume.astro
import { getEntry } from 'astro:content';
import BaseLayout from '../layouts/BaseLayout.astro';
import ResumeBody from '../components/resume/ResumeBody.astro';

const resume = await getEntry('resume', 'profile');
if (!resume) throw new Error('Resume entry not found at frontend/src/content/resume/profile.md');
---

<BaseLayout
  title={`${resume.data.name} — Resume`}
  description={`Resume for ${resume.data.name}, ${resume.data.title}.`}
  includeResumeModal={false}
>
  <main class="resume-page">
    <header class="resume-page-header">
      <h1 class="resume-page-name">{resume.data.name}</h1>
      <p class="resume-page-title">{resume.data.title}</p>
      <p class="resume-page-location">{resume.data.location}</p>
      <p class="resume-page-links">
        <a href={resume.data.links.linkedin}>LinkedIn</a>
        {resume.data.links.github && <> · <a href={resume.data.links.github}>GitHub</a></>}
        {resume.data.links.portfolio && <> · <a href={resume.data.links.portfolio}>Portfolio</a></>}
      </p>
    </header>
    <ResumeBody resume={resume} />
    <footer class="resume-page-footer">
      <p>Contact via <a href="/#contact">contact form</a> or <a href={resume.data.links.linkedin}>LinkedIn</a>.</p>
    </footer>
  </main>
</BaseLayout>
```

- [ ] **Step 2: Run dev server and visit `/resume`:**

```bash
cd frontend && npm run dev
```

Open `http://localhost:4321/resume` in browser. Confirm:
- Resume content renders without errors
- All section data appears correctly
- Stop the dev server with Ctrl+C

- [ ] **Step 3: Commit:**

```bash
git add frontend/src/pages/resume.astro
git commit -m "feat(p9): add /resume standalone page"
```

### Task 3.3: Add resume styling to `global.css`

**Files:**
- Modify: `frontend/src/styles/global.css`

- [ ] **Step 1: Read the current global.css** to understand the existing token system:

```bash
cat frontend/src/styles/global.css
```

- [ ] **Step 2: Append the resume style block** at the bottom of `global.css`:

```css
/* ===== Resume ===== */

/* Page (used standalone at /resume and as the print source for the PDF) */
.resume-page {
  max-width: 48rem;
  margin: 0 auto;
  padding: 3rem 1.5rem 4rem;
  color: var(--color-text-primary);
  font-family: ui-sans-serif, system-ui, sans-serif;
  line-height: 1.6;
}

.resume-page-header {
  margin-bottom: 2.5rem;
  padding-bottom: 1.5rem;
  border-bottom: 1px solid var(--color-surface-border);
}

.resume-page-name {
  font-size: 2rem;
  font-weight: 700;
  letter-spacing: -0.03em;
  margin-bottom: 0.25rem;
}

.resume-page-title {
  font-size: 1.125rem;
  color: var(--color-text-secondary);
  margin-bottom: 0.5rem;
}

.resume-page-location,
.resume-page-links {
  font-size: 0.875rem;
  color: var(--color-text-muted);
}

.resume-page-links a {
  color: var(--color-accent-indigo);
  text-decoration: none;
}

.resume-page-links a:hover { text-decoration: underline; }

.resume-page-footer {
  margin-top: 3rem;
  padding-top: 1.5rem;
  border-top: 1px solid var(--color-surface-border);
  font-size: 0.875rem;
  color: var(--color-text-muted);
}

.resume-page-footer a { color: var(--color-accent-indigo); }

/* Sections */
.resume-section { margin-top: 2rem; }
.resume-heading {
  font-size: 0.875rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: var(--color-accent-indigo);
  margin-bottom: 0.75rem;
}

/* Work history */
.resume-job { margin-bottom: 1.5rem; }
.resume-job-header { margin-bottom: 0.5rem; }
.resume-job-title-row {
  display: flex;
  align-items: baseline;
  gap: 0.5rem;
  flex-wrap: wrap;
}
.resume-job-title { font-weight: 600; }
.resume-job-company { color: var(--color-text-secondary); }
.resume-job-meta {
  font-size: 0.875rem;
  color: var(--color-text-muted);
}
.resume-job-bullets {
  list-style: disc;
  padding-left: 1.5rem;
  margin-top: 0.5rem;
}
.resume-job-bullets li { margin-bottom: 0.25rem; }
.resume-job-tech {
  display: flex;
  flex-wrap: wrap;
  gap: 0.375rem;
  margin-top: 0.5rem;
}
.resume-tech-tag {
  display: inline-block;
  font-size: 0.75rem;
  padding: 0.125rem 0.5rem;
  background: var(--color-surface);
  border: 1px solid var(--color-surface-border);
  border-radius: 0.25rem;
  color: var(--color-text-secondary);
}

/* Skills */
.resume-skill-bucket { margin-bottom: 0.75rem; }
.resume-skill-label {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 0.25rem;
}
.resume-skill-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.375rem 0.75rem;
  list-style: none;
  padding: 0;
}
.resume-skill-list li { color: var(--color-text-secondary); }
.resume-skill-list li:not(:last-child)::after {
  content: '·';
  margin-left: 0.75rem;
  color: var(--color-text-muted);
}

/* Generic list */
.resume-list { list-style: none; padding: 0; }
.resume-list li { margin-bottom: 0.375rem; }

/* Modal */
.resume-modal {
  background: var(--color-base);
  color: var(--color-text-primary);
  border: 1px solid var(--color-surface-border);
  border-radius: 0.75rem;
  padding: 0;
  max-width: 56rem;
  width: 90vw;
  max-height: 90vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
.resume-modal::backdrop {
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
}
.resume-modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid var(--color-surface-border);
  flex-shrink: 0;
}
.resume-modal-title { display: flex; flex-direction: column; }
.resume-modal-title .name { font-weight: 700; }
.resume-modal-title .title { font-size: 0.875rem; color: var(--color-text-muted); }
.resume-modal-actions { display: flex; gap: 0.5rem; align-items: center; }
.btn-download {
  padding: 0.375rem 0.875rem;
  background: var(--color-accent-indigo);
  color: white;
  border-radius: 0.375rem;
  font-size: 0.75rem;
  text-decoration: none;
}
.btn-download:hover { opacity: 0.9; }
.btn-close {
  background: transparent;
  border: none;
  color: var(--color-text-muted);
  font-size: 1.5rem;
  cursor: pointer;
  padding: 0.25rem 0.5rem;
  line-height: 1;
}
.btn-close:hover { color: var(--color-text-primary); }
.resume-modal-body {
  overflow-y: auto;
  padding: 1.5rem;
}
.resume-modal-footer {
  padding: 1rem 1.5rem;
  border-top: 1px solid var(--color-surface-border);
  font-size: 0.875rem;
  color: var(--color-text-muted);
  flex-shrink: 0;
}
.resume-modal-footer a { color: var(--color-accent-indigo); }

/* Print stylesheet — used by Playwright when generating the PDF */
@media print {
  /* Hide site chrome on the /resume page when printing */
  header, footer.site-footer, nav { display: none !important; }
  .resume-page { max-width: none; padding: 0; }
  .resume-page-header { border-bottom-color: #ccc; }
  .resume-section { page-break-inside: avoid; }
  .resume-heading { color: #4f46e5; }
  body { background: white; color: black; }
  a { color: #4f46e5; text-decoration: none; }
  .resume-tech-tag {
    background: #f5f5f5;
    border-color: #ddd;
    color: #444;
  }
}
```

- [ ] **Step 3: Visit `http://localhost:4321/resume`** in `npm run dev` to confirm the styles apply correctly. Use browser print preview (Cmd+P) to confirm the print stylesheet looks good.

- [ ] **Step 4: Commit:**

```bash
git add frontend/src/styles/global.css
git commit -m "feat(p9): add resume styles for page, modal, and print"
```

---

## Phase 4: Modal & Nav Integration

### Task 4.1: `ResumeModal.astro`

**Files:**
- Create: `frontend/src/components/ResumeModal.astro`

- [ ] **Step 1: Create the modal component:**

```astro
---
// frontend/src/components/ResumeModal.astro
import { getEntry } from 'astro:content';
import ResumeBody from './resume/ResumeBody.astro';

const resume = await getEntry('resume', 'profile');
if (!resume) throw new Error('Resume entry not found at frontend/src/content/resume/profile.md');
---

<dialog id="resume-modal" class="resume-modal">
  <header class="resume-modal-header">
    <div class="resume-modal-title">
      <span class="name">{resume.data.name}</span>
      <span class="title">{resume.data.title}</span>
    </div>
    <div class="resume-modal-actions">
      <a href="/Rafael_Abreu_Resume.pdf" download class="btn-download">Download PDF</a>
      <button type="button" class="btn-close" aria-label="Close resume">×</button>
    </div>
  </header>
  <div class="resume-modal-body">
    <ResumeBody resume={resume} />
  </div>
  <footer class="resume-modal-footer">
    Contact via <a href="/#contact">contact form</a> or
    <a href={resume.data.links.linkedin} target="_blank" rel="noopener noreferrer">LinkedIn →</a>
  </footer>
</dialog>

<script>
  const modal = document.getElementById('resume-modal') as HTMLDialogElement | null;

  // Event delegation handles both static Astro Nav anchors and React-hydrated MobileNav buttons.
  document.addEventListener('click', (e) => {
    const trigger = (e.target as HTMLElement)?.closest('[data-resume-trigger]');
    if (trigger) {
      e.preventDefault();
      modal?.showModal();
    }
  });

  modal?.querySelector('.btn-close')?.addEventListener('click', () => modal.close());
  modal?.addEventListener('click', (e) => {
    if (e.target === modal) modal.close();
  });
</script>
```

- [ ] **Step 2: Verify type-check:** `cd frontend && npx astro check`. Expected: no errors.

- [ ] **Step 3: Commit:**

```bash
git add frontend/src/components/ResumeModal.astro
git commit -m "feat(p9): add ResumeModal component"
```

### Task 4.2: BaseLayout integration with opt-out prop

**Files:**
- Modify: `frontend/src/layouts/BaseLayout.astro`

- [ ] **Step 1: Update the layout to include the modal site-wide with an opt-out prop:**

Replace the current `interface Props` and `Astro.props` destructure, and add the modal include before `</body>`:

```astro
---
import '../styles/global.css';
import ResumeModal from '../components/ResumeModal.astro';

interface Props {
  title?: string;
  description?: string;
  includeResumeModal?: boolean;
}

const {
  title = 'Rafael Abreu | Software Engineer',
  description = 'Software Engineer building intelligent systems across the full stack. Go, Python, React, and the infrastructure that connects them.',
  includeResumeModal = true,
} = Astro.props;
---

<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="description" content={description} />
    <meta name="theme-color" content="#050508" />
    <link rel="icon" type="image/x-icon" href="/favicon.ico" />
    <title>{title}</title>
    <script is:inline>
      (function() {
        var theme = localStorage.getItem('theme');
        if (!theme) {
          theme = window.matchMedia('(prefers-color-scheme: light)').matches ? 'light' : 'dark';
        }
        document.documentElement.dataset.theme = theme;
        var meta = document.querySelector('meta[name="theme-color"]');
        if (meta) meta.setAttribute('content', theme === 'light' ? '#F8F7F4' : '#050508');
      })();
    </script>
    {import.meta.env.PROD && <script is:inline async src="https://cloud.umami.is/script.js" data-website-id="f17cb708-288b-471b-989f-d7115bfcbec4"></script>}
  </head>
  <body>
    <slot />

    {includeResumeModal && <ResumeModal />}

    <script>
      const observer = new IntersectionObserver(
        (entries) => {
          entries.forEach((entry) => {
            if (entry.isIntersecting) {
              entry.target.classList.add('visible');
              observer.unobserve(entry.target);
            }
          });
        },
        { threshold: 0.1 }
      );
      document.querySelectorAll('.reveal').forEach((el) => observer.observe(el));
    </script>
  </body>
</html>
```

- [ ] **Step 2: Run dev server and verify the modal renders** (it'll be hidden until triggered): `cd frontend && npm run dev`. Open the homepage, run in console: `document.getElementById('resume-modal').showModal()`. Modal should appear with content. Stop dev server.

- [ ] **Step 3: Commit:**

```bash
git add frontend/src/layouts/BaseLayout.astro
git commit -m "feat(p9): include ResumeModal in BaseLayout with opt-out prop"
```

### Task 4.3: Update `Nav.astro`

**Files:**
- Modify: `frontend/src/components/Nav.astro`

- [ ] **Step 1: Replace the contents:**

```astro
---
import MobileNav from '../islands/MobileNav.tsx';
import ThemeToggle from '../islands/ThemeToggle.tsx';

const navItems = [
  { label: 'About', href: '/#about' },
  { label: 'Projects', href: '/#projects' },
  { label: 'Skills', href: '/#skills' },
  { label: 'Playground', href: '/playground' },
  { label: 'Contact', href: '/#contact' },
];

const isResumePage = Astro.url.pathname === '/resume';
---

<header class="sticky top-0 z-50 bg-base/80 backdrop-blur-md border-b border-surface-border">
  <nav class="container mx-auto max-w-6xl px-6 py-4 flex items-center justify-between">
    <a href="/" class="font-mono text-text-primary text-sm font-bold tracking-wide">
      RA<span class="text-accent-indigo">.</span>
    </a>

    <div class="hidden md:flex items-center gap-7">
      {navItems.map((item) => (
        <a
          href={item.href}
          class="text-text-muted text-xs tracking-[1.5px] uppercase hover:text-accent-indigo transition-colors duration-200"
        >
          {item.label}
        </a>
      ))}
    </div>

    <div class="hidden md:flex items-center gap-3">
      <ThemeToggle client:load />
      {!isResumePage && (
        <a
          href="/resume"
          data-resume-trigger
          class="px-4 py-1.5 border border-accent-indigo/40 rounded-btn text-accent-indigo text-xs tracking-wide hover:bg-accent-indigo/10 transition-colors duration-200"
        >
          RESUME
        </a>
      )}
      <a
        href="/#contact"
        class="px-4 py-1.5 bg-accent-indigo rounded-btn text-white text-xs tracking-wide hover:bg-accent-indigo/90 transition-colors duration-200"
      >
        LET'S TALK
      </a>
    </div>

    <MobileNav client:load navItems={navItems} isResumePage={isResumePage} />
  </nav>
</header>
```

Notes:
- The Resume link uses `href="/resume"` so users who somehow click it without the modal listener installed (e.g., on `/resume` itself, or with JS disabled) navigate to the standalone page. The `data-resume-trigger` attribute is what the modal script catches to override that navigation.
- `Astro.url.pathname` lets the Nav suppress the Resume button on `/resume`.

- [ ] **Step 2: Verify dev server still renders the home page correctly** and clicking RESUME opens the modal: `cd frontend && npm run dev`, visit `http://localhost:4321/`, click RESUME button. Modal should open. Visit `http://localhost:4321/resume` directly — RESUME button should be hidden.

- [ ] **Step 3: Commit:**

```bash
git add frontend/src/components/Nav.astro
git commit -m "feat(p9): wire Nav RESUME button to modal trigger; hide on /resume"
```

### Task 4.4: Update `MobileNav.tsx`

**Files:**
- Modify: `frontend/src/islands/MobileNav.tsx`

- [ ] **Step 1: Replace the file contents:**

```tsx
import { useState } from 'react';
import ThemeToggle from './ThemeToggle';

interface NavItem {
  label: string;
  href: string;
}

interface Props {
  navItems: NavItem[];
  isResumePage: boolean;
}

export default function MobileNav({ navItems, isResumePage }: Props) {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div className="md:hidden">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="relative z-50 w-8 h-8 flex flex-col justify-center items-center gap-1.5"
        aria-label={isOpen ? 'Close menu' : 'Open menu'}
      >
        <span
          className={`block w-5 h-0.5 bg-text-primary transition-all duration-300 ${
            isOpen ? 'rotate-45 translate-y-1' : ''
          }`}
        />
        <span
          className={`block w-5 h-0.5 bg-text-primary transition-all duration-300 ${
            isOpen ? '-rotate-45 -translate-y-1' : ''
          }`}
        />
      </button>

      {isOpen && (
        <div className="fixed inset-0 z-40 bg-base flex flex-col items-center justify-center gap-8">
          <div className="absolute top-5 right-14">
            <ThemeToggle />
          </div>
          {navItems.map((item) => (
            <a
              key={item.label}
              href={item.href}
              onClick={() => setIsOpen(false)}
              className="text-text-primary text-2xl font-bold hover:text-accent-indigo transition-colors"
              style={{ letterSpacing: '-0.06em' }}
            >
              {item.label}
            </a>
          ))}
          <div className="flex gap-4 mt-4">
            {!isResumePage && (
              <a
                href="/resume"
                data-resume-trigger
                onClick={() => setIsOpen(false)}
                className="px-6 py-2 border border-accent-indigo/40 rounded-lg text-accent-indigo text-sm"
              >
                RESUME
              </a>
            )}
            <a
              href="/#contact"
              onClick={() => setIsOpen(false)}
              className="px-6 py-2 bg-accent-indigo rounded-lg text-white text-sm"
            >
              LET'S TALK
            </a>
          </div>
        </div>
      )}
    </div>
  );
}
```

- [ ] **Step 2: Verify on a mobile viewport** in the dev server (`cd frontend && npm run dev`, then resize browser narrow or use device-emulation). Tap the hamburger, tap RESUME — modal should open.

- [ ] **Step 3: Commit:**

```bash
git add frontend/src/islands/MobileNav.tsx
git commit -m "feat(p9): wire MobileNav RESUME button to modal trigger; hide on /resume"
```

---

## Phase 5: PDF Generation in CI

### Task 5.1: Add Playwright and serve-handler

**Files:**
- Modify: `frontend/package.json`

- [ ] **Step 1: Install:**

```bash
cd frontend
npm install --save-dev playwright serve-handler
```

- [ ] **Step 2: Install Chromium for Playwright** (one-time, used locally; in CI we install fresh):

```bash
npx playwright install chromium
```

- [ ] **Step 3: Commit:**

```bash
git add frontend/package.json frontend/package-lock.json
git commit -m "chore(p9): add Playwright and serve-handler for PDF generation"
```

### Task 5.2: Write the PDF generator script

**Files:**
- Create: `frontend/scripts/generate-pdf.mjs`

- [ ] **Step 1: Create the script:**

```javascript
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
```

- [ ] **Step 2: Add npm script to `package.json`** under "scripts":

```json
"build:pdf": "node scripts/generate-pdf.mjs"
```

- [ ] **Step 3: Test locally:**

```bash
cd frontend
npm run build       # produce dist/
npm run build:pdf   # produce dist/Rafael_Abreu_Resume.pdf
ls -la dist/Rafael_Abreu_Resume.pdf
```

Expected: file exists, > 10KB. Open it manually to eyeball the layout.

- [ ] **Step 4: Commit:**

```bash
git add frontend/scripts/generate-pdf.mjs frontend/package.json
git commit -m "feat(p9): add Playwright-based PDF generator script"
```

### Task 5.3: Update frontend deploy workflow

**Files:**
- Modify: `.github/workflows/deploy-frontend.yml`

- [ ] **Step 1: Replace the file contents:**

```yaml
name: Deploy Frontend

on:
  push:
    branches: [main]
    paths:
      - 'frontend/**'

jobs:
  deploy:
    name: Deploy to Netlify
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: frontend
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-node@v4
        with:
          node-version: 22
          cache: npm
          cache-dependency-path: frontend/package-lock.json

      - run: npm ci

      - name: Build
        run: npm run build
        env:
          PUBLIC_API_URL: ${{ secrets.PUBLIC_API_URL }}

      - name: Install Playwright Chromium
        run: npx playwright install --with-deps chromium

      - name: Generate Resume PDF
        run: npm run build:pdf

      - name: Verify PDF was generated
        run: |
          test -f dist/Rafael_Abreu_Resume.pdf
          test "$(stat -c%s dist/Rafael_Abreu_Resume.pdf)" -gt 10240

      - name: Deploy to Netlify
        uses: nwtgck/actions-netlify@v3
        with:
          publish-dir: frontend/dist
          production-deploy: true
        env:
          NETLIFY_AUTH_TOKEN: ${{ secrets.NETLIFY_AUTH_TOKEN }}
          NETLIFY_SITE_ID: ${{ secrets.NETLIFY_SITE_ID }}
```

- [ ] **Step 2: Lint the YAML:**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio
yq eval '.github/workflows/deploy-frontend.yml' > /dev/null  # if yq is available
```

(If `yq` isn't installed, skip — the actual validation happens when GitHub parses it on next run.)

- [ ] **Step 3: Commit:**

```bash
git add .github/workflows/deploy-frontend.yml
git commit -m "ci(p9): generate resume PDF in deploy-frontend pipeline"
```

---

## Phase 6: Chatbot Integration

### Task 6.1: Add API gitignore for embedded directory

**Files:**
- Create: `api/.gitignore`

- [ ] **Step 1: Create the file:**

```
# Build artifacts mirroring frontend content (populated by `make sync-content`)
internal/chat/embedded/
```

- [ ] **Step 2: Commit:**

```bash
git add api/.gitignore
git commit -m "chore(p9): add api/.gitignore for embedded build artifacts"
```

### Task 6.2: Add API Makefile

**Files:**
- Create: `api/Makefile`

- [ ] **Step 1: Create the file:**

```makefile
.PHONY: sync-content dev test build clean

# Copy frontend content into embedded/ for go:embed.
# Required before `go build` / `go run` / `go test` because go:embed needs
# the referenced paths to exist at compile time.
sync-content:
	@rm -rf internal/chat/embedded
	@mkdir -p internal/chat/embedded/projects
	@cp ../frontend/src/content/resume/profile.md internal/chat/embedded/resume.md
	@cp ../frontend/src/content/projects/*.md internal/chat/embedded/projects/
	@echo "Content synced into internal/chat/embedded/"

dev: sync-content
	go run ./cmd/server

test: sync-content
	go test ./... -v

build: sync-content
	go build -o bin/server ./cmd/server

clean:
	@rm -rf internal/chat/embedded bin
```

- [ ] **Step 2: Run `make sync-content` and verify:**

```bash
cd api && make sync-content
ls internal/chat/embedded/
# Expected: resume.md (file), projects/ (directory with .md files)
```

- [ ] **Step 3: Commit:**

```bash
git add api/Makefile
git commit -m "chore(p9): add API Makefile with sync-content target"
```

### Task 6.3: Refactor `context.go` — split persona/rules helpers (TDD)

**Files:**
- Modify: `api/internal/chat/context.go`
- Create: `api/internal/chat/context_test.go`

This refactor decomposes the monolithic `SystemPrompt()` string into helper functions, **before** introducing the embed step. This isolates the structural change from the content-source change.

- [ ] **Step 1: Read existing `context.go`:**

```bash
cat api/internal/chat/context.go
```

- [ ] **Step 2: Write failing tests:**

```go
// api/internal/chat/context_test.go
package chat

import (
	"strings"
	"testing"
)

func TestSystemPromptIncludesPersona(t *testing.T) {
	p := SystemPrompt()
	if !strings.Contains(p, "Rafael Abreu's AI assistant") {
		t.Error("SystemPrompt should include the persona marker phrase")
	}
}

func TestSystemPromptIncludesRules(t *testing.T) {
	p := SystemPrompt()
	if !strings.Contains(p, "Only answer questions about Rafael") {
		t.Error("SystemPrompt should include the rules marker phrase")
	}
}

func TestSystemPromptStructure(t *testing.T) {
	p := SystemPrompt()
	personaIdx := strings.Index(p, "Rafael Abreu's AI assistant")
	rulesIdx := strings.Index(p, "Only answer questions about Rafael")
	if personaIdx == -1 || rulesIdx == -1 {
		t.Fatal("required markers missing")
	}
	if !(personaIdx < rulesIdx) {
		t.Error("persona block must precede rules block")
	}
}
```

- [ ] **Step 3: Run tests, expect PASS** (the existing `SystemPrompt()` already has both phrases):

```bash
cd api && go test ./internal/chat/ -v
```

Expected: 3 tests pass against the existing implementation.

- [ ] **Step 4: Refactor `context.go` to extract helpers** (preserves the prompt content exactly, just splits the structure):

```go
package chat

import "strings"

func SystemPrompt() string {
	var b strings.Builder
	b.WriteString(personaBlock())
	b.WriteString("\n\n")
	b.WriteString(rulesBlock())
	return b.String()
}

func personaBlock() string {
	return `You are Rafael Abreu's AI assistant on his portfolio website. You answer questions about Rafael's professional background, skills, projects, and experience.

ABOUT RAFAEL:
- Software Engineer based in Toronto, Ontario, Canada
- Full-stack developer with experience across technology, television, and education industries
- Currently focused on Go, Python, and AI/ML engineering
- Completed the University of Toronto Coding Bootcamp
- Holds a Red Hat Kubernetes certificate (OpenShift training)

TECHNICAL SKILLS:
- Languages: Go, Python, JavaScript, TypeScript, HTML, CSS
- Frameworks: React, Astro, Node.js, Tailwind CSS, Bootstrap
- Infrastructure: Kubernetes, Docker, CI/CD, Netlify, Fly.io
- AI & Data: LLM Integration, Gemini API, RAG, MongoDB, MySQL
- Practices: TDD, OOP, MVC, System Design, Agile

NOTABLE PROJECTS:
- FPTV: Developed the website for a Toronto-based TV channel broadcasting across Canada. Features glassmorphism dashboard. (HTML, CSS, JavaScript, Weebly)
- Freckles Design: Artist portfolio platform built in partnership with a fine artist. (React, Bootstrap, Netlify)
- Charge It Up: EV charging station locator for North America. (HTML, CSS, Tailwind, JavaScript)
- This portfolio itself: Astro frontend + Go API backend with AI chatbot, demonstrating modern architecture.`
}

func rulesBlock() string {
	return `RULES:
- Only answer questions about Rafael's professional background, skills, projects, and experience.
- Be concise and helpful. Keep responses under 200 words.
- If asked something unrelated to Rafael's professional life, politely redirect: "I'm here to help with questions about Rafael's experience and work. What would you like to know?"
- Be friendly and professional in tone.
- You may suggest that visitors check out specific projects or sections of the portfolio when relevant.`
}
```

- [ ] **Step 5: Re-run tests, expect PASS:**

```bash
cd api && go test ./internal/chat/ -v
```

- [ ] **Step 6: Commit:**

```bash
git add api/internal/chat/context.go api/internal/chat/context_test.go
git commit -m "refactor(p9): split SystemPrompt into persona and rules helpers

No content change — the assembled string is byte-identical to the
prior monolithic version. Sets up clean insertion points for the
upcoming embedded-markdown integration."
```

### Task 6.4: Switch `context.go` to embedded markdown content (TDD)

**Files:**
- Modify: `api/internal/chat/context.go`
- Modify: `api/internal/chat/context_test.go`

- [ ] **Step 1: Add failing tests for the new behavior:**

Append to `context_test.go`:

```go
func TestSystemPromptIncludesResumeContent(t *testing.T) {
	p := SystemPrompt()
	if !strings.Contains(p, "Rafael Abreu") {
		t.Error("SystemPrompt should include resume content (e.g., the canonical name)")
	}
	if !strings.Contains(p, "## RESUME") {
		t.Error("SystemPrompt should include the RESUME section header")
	}
}

func TestSystemPromptIncludesAllProjects(t *testing.T) {
	p := SystemPrompt()
	if !strings.Contains(p, "## PROJECT CASE STUDIES") {
		t.Error("SystemPrompt should include the PROJECT CASE STUDIES section header")
	}
	// Each project markdown's frontmatter title should appear.
	// Read embedded directory to know the count.
	entries, err := projectsFS.ReadDir("embedded/projects")
	if err != nil {
		t.Fatalf("failed to read embedded projects: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("no embedded projects found — sync step must run before tests")
	}
	// Spot-check: expect at least one project's content to be in the prompt.
	for _, entry := range entries {
		content, err := projectsFS.ReadFile("embedded/projects/" + entry.Name())
		if err != nil {
			t.Errorf("failed to read %s: %v", entry.Name(), err)
			continue
		}
		// First 100 chars should appear if the file was concatenated.
		snippet := string(content)
		if len(snippet) > 100 {
			snippet = snippet[:100]
		}
		if !strings.Contains(p, snippet) {
			t.Errorf("project file %s content not found in SystemPrompt", entry.Name())
		}
	}
}

func TestSystemPromptOrderingResumeBeforeProjects(t *testing.T) {
	p := SystemPrompt()
	resumeIdx := strings.Index(p, "## RESUME")
	projectsIdx := strings.Index(p, "## PROJECT CASE STUDIES")
	if resumeIdx == -1 || projectsIdx == -1 {
		t.Fatal("required section headers missing")
	}
	if !(resumeIdx < projectsIdx) {
		t.Error("RESUME section must precede PROJECT CASE STUDIES")
	}
}
```

- [ ] **Step 2: Run tests, expect FAIL** (the embed-based fields don't exist yet):

```bash
cd api && make sync-content && go test ./internal/chat/ -v
```

Expected: compilation error or test failures referencing `projectsFS`.

- [ ] **Step 3: Update `context.go` to use embedded files:**

```go
package chat

import (
	"embed"
	"io/fs"
	"sort"
	"strings"
)

//go:embed embedded/resume.md
var resumeMD string

//go:embed embedded/projects/*.md
var projectsFS embed.FS

func SystemPrompt() string {
	var b strings.Builder
	b.WriteString(personaBlock())
	b.WriteString("\n\n## RESUME\n\n")
	b.WriteString(resumeMD)
	b.WriteString("\n\n## PROJECT CASE STUDIES\n\n")
	b.WriteString(assembleProjects())
	b.WriteString("\n\n")
	b.WriteString(rulesBlock())
	return b.String()
}

func assembleProjects() string {
	entries, err := fs.ReadDir(projectsFS, "embedded/projects")
	if err != nil {
		return ""
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })

	var b strings.Builder
	for _, entry := range entries {
		content, err := fs.ReadFile(projectsFS, "embedded/projects/"+entry.Name())
		if err != nil {
			continue
		}
		b.Write(content)
		b.WriteString("\n\n---\n\n")
	}
	return b.String()
}

func personaBlock() string {
	// Persona is generic — Rafael's specific bio/skills now come from resume.md.
	return `You are Rafael Abreu's AI assistant on his portfolio website. You answer questions about Rafael's professional background, skills, projects, and experience.

The sections below contain Rafael's current resume and project case studies. Use them as your primary source of truth. The information here is authoritative; ignore any prior knowledge that conflicts.`
}

func rulesBlock() string {
	return `RULES:
- Only answer questions about Rafael's professional background, skills, projects, and experience.
- Be concise and helpful. Keep responses under 200 words.
- If asked something unrelated to Rafael's professional life, politely redirect: "I'm here to help with questions about Rafael's experience and work. What would you like to know?"
- Be friendly and professional in tone.
- You may suggest that visitors check out specific projects or sections of the portfolio when relevant.`
}
```

- [ ] **Step 4: Run tests, expect PASS:**

```bash
cd api && make test
```

- [ ] **Step 5: Commit:**

```bash
git add api/internal/chat/context.go api/internal/chat/context_test.go
git commit -m "feat(p9): replace hardcoded SystemPrompt content with embedded markdown

context.go now embeds frontend/src/content/resume/profile.md and
frontend/src/content/projects/*.md (copied via make sync-content) and
assembles them into the system prompt at startup. The prior bio/skills/
projects block in personaBlock() is removed in favor of the resume —
single source of truth."
```

### Task 6.5: Update API deploy workflow

**Files:**
- Modify: `.github/workflows/deploy-api.yml`

- [ ] **Step 1: Replace the file contents:**

```yaml
name: Deploy API

on:
  push:
    branches: [main]
    paths:
      - 'api/**'
      - 'frontend/src/content/resume/**'
      - 'frontend/src/content/projects/**'

jobs:
  deploy:
    name: Deploy to Fly.io
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: api
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.26'

      - name: Sync content for embedding
        run: make sync-content

      - name: Run tests
        run: go test ./... -v

      - name: Set up Fly CLI
        uses: superfly/flyctl-actions/setup-flyctl@master

      - name: Deploy to Fly.io
        run: flyctl deploy --remote-only
        env:
          FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
```

Note the additional `paths:` triggers — when `frontend/src/content/resume/**` or `frontend/src/content/projects/**` change on `main`, the API redeploys so the chatbot picks up new content.

- [ ] **Step 2: Commit:**

```bash
git add .github/workflows/deploy-api.yml
git commit -m "ci(p9): sync content before API build; trigger on resume/project changes"
```

---

## Phase 7: E2E Tests, Roadmap, README

### Task 7.1: E2E Playwright test for the modal flow

**Files:**
- Create: `frontend/playwright.config.ts`
- Create: `frontend/tests/e2e/resume.spec.ts`
- Modify: `frontend/package.json`

- [ ] **Step 1: Create `frontend/playwright.config.ts`:**

```typescript
import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: './tests/e2e',
  webServer: {
    command: 'npm run preview',
    url: 'http://localhost:4321',
    reuseExistingServer: !process.env.CI,
    timeout: 120000,
  },
  use: {
    baseURL: 'http://localhost:4321',
  },
});
```

- [ ] **Step 2: Install Playwright test runner:**

```bash
cd frontend && npm install --save-dev @playwright/test
```

- [ ] **Step 3: Create `frontend/tests/e2e/resume.spec.ts`:**

```typescript
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
```

- [ ] **Step 4: Add `test:e2e` script to `package.json`:**

```json
"test:e2e": "playwright test"
```

- [ ] **Step 5: Run the E2E test against the built site:**

```bash
cd frontend
npm run build
npm run build:pdf
npm run test:e2e
```

Expected: all 4 tests pass.

- [ ] **Step 6: Commit:**

```bash
git add frontend/playwright.config.ts frontend/tests/e2e/resume.spec.ts frontend/package.json frontend/package-lock.json
git commit -m "test(p9): add E2E coverage for resume modal and /resume page"
```

### Task 7.2: Update phase 2 roadmap

**Files:**
- Modify: `docs/phase2-roadmap.md`

- [ ] **Step 1: Read current roadmap to find P8 entry:**

```bash
cat docs/phase2-roadmap.md | tail -40
```

- [ ] **Step 2: Append the P9 entry** before the "Dependency Map" section:

```markdown
### ~~P9: Resume Feature~~ ✅ Done (PR #__)

**What:** A markdown-driven resume rendered as an in-site modal, downloadable PDF, and chatbot grounding source — replacing the previously broken Resume button.

**Why:** The Resume button on the portfolio was a 404. Rather than dropping a static PDF that goes stale, this establishes a single-source-of-truth model: edit one markdown file, three artifacts (modal, PDF, chatbot prompt) regenerate deterministically. Also fixes a latent staleness bug where the chatbot's project knowledge drifted from the actual portfolio.

**Scope:**
- New Astro content collection at `frontend/src/content/resume/profile.md` with full Zod schema
- 9 section components rendering structured data into HTML
- ResumeModal using the native `<dialog>` element (no React island)
- /resume standalone page (also serves as Playwright's PDF print source)
- Build-time PDF generation via Playwright in CI
- Chatbot system prompt assembled from resume + project markdown via `go:embed`
- E2E coverage for the modal flow

**Estimated effort:** 1-2 sessions

**Out of scope (deferred):** AI-driven LinkedIn ingestion, GitHub-based skill auto-extraction, AI-generated bullets. Considered during brainstorming and rejected as overkill for an artifact that changes 4–6 times per year.

---
```

Also update the Dependency Map block to include P9:

```
P9 (Resume Feature)    — independent, replaces broken Resume button
```

- [ ] **Step 3: Commit:**

```bash
git add docs/phase2-roadmap.md
git commit -m "docs(p9): add P9 resume feature entry to phase 2 roadmap"
```

### Task 7.3: Update root `CLAUDE.md` with first-time-clone note

**Files:**
- Modify: `CLAUDE.md`

- [ ] **Step 1: Add a sentence to the API "Build & Run" block** in `CLAUDE.md`:

Find the section that reads:
````
### API (Go)
```bash
cd api
go run ./cmd/server  # dev server on :8080
go test ./... -v     # run tests
go vet ./...         # lint
go build ./cmd/server
```
Requires: Go 1.26+
````

Insert a note after the code block:

```markdown
**First-time clone:** run `make sync-content` (or `make dev`) before any `go` command — `internal/chat/embedded/` is gitignored and populated from `frontend/src/content/` at build time. `make dev` and `make test` invoke `sync-content` automatically.
```

- [ ] **Step 2: Commit:**

```bash
git add CLAUDE.md
git commit -m "docs(p9): document first-time-clone setup for embedded chat content"
```

### Task 7.4: Final smoke test — full stack runs locally

**Files:** none modified.

- [ ] **Step 1: Frontend build pipeline succeeds end-to-end:**

```bash
cd frontend
npm ci
npm run build
npm run build:pdf
ls -la dist/Rafael_Abreu_Resume.pdf
```

Expected: PDF exists, > 10KB.

- [ ] **Step 2: API builds and tests pass with embedded content:**

```bash
cd ../api
make sync-content
make test
make build
```

Expected: tests green, binary built.

- [ ] **Step 3: Run the full app locally and verify chatbot picks up new content:**

In one terminal:

```bash
cd api && make dev
```

In another:

```bash
cd frontend && npm run dev
```

Open `http://localhost:4321/`:
- Click RESUME → modal opens with the resume content from `profile.md`
- Click Download PDF → PDF downloads (will only work after `npm run build && npm run build:pdf`)
- Open the chat widget → ask "What is Rafael's most recent role?" → response should reflect content from `profile.md`

- [ ] **Step 4: Stop both servers.** No commit needed (no files changed).

### Task 7.5: Open the pull request

**Files:** none.

- [ ] **Step 1: Push the branch:**

```bash
git push -u origin feature/resume
```

(Per the project's no-push-without-consent rule, only do this step when the user has explicitly asked to push.)

- [ ] **Step 2: Open the PR:**

```bash
gh pr create --title "feat(p9): resume feature — modal, PDF, chatbot grounding" --body "$(cat <<'EOF'
## Summary

Implements P9 per spec `docs/superpowers/specs/2026-04-30-resume-feature-design.md`. Closes #25.

- **Single source of truth:** `frontend/src/content/resume/profile.md`
- **In-site modal:** native `<dialog>`, no React hydration
- **Build-time PDF:** Playwright in CI, `frontend/dist/Rafael_Abreu_Resume.pdf`
- **Chatbot grounding:** `SystemPrompt()` now embeds resume + project markdown via `go:embed`

## Test plan
- [ ] CI green
- [ ] Visit `https://<deploy-preview>/` → click RESUME → modal renders correctly in dark and light themes
- [ ] Click "Download PDF" → file downloads, opens cleanly
- [ ] Visit `https://<deploy-preview>/resume` directly → page renders, RESUME nav button hidden
- [ ] Open chat widget → ask about a recent role → response reflects new content

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

- [ ] **Step 3:** Return the PR URL to the user.

---

## Self-Review Checklist (run after writing the plan)

**Spec coverage:**
- ✅ Schema (Task 1.1) covers every field in the spec's data model
- ✅ All 9 section components (Tasks 2.2–2.8) match the spec's component table
- ✅ ResumeBody assembler (Task 3.1) and standalone page (Task 3.2)
- ✅ Modal (Task 4.1) with event delegation per spec
- ✅ BaseLayout opt-out prop (Task 4.2)
- ✅ Nav and MobileNav conditional hiding (Tasks 4.3, 4.4)
- ✅ PDF generator with random port (Task 5.2)
- ✅ Frontend CI updates (Task 5.3)
- ✅ Chatbot embedding refactor in two steps — structural split first (Task 6.3), content swap second (Task 6.4)
- ✅ API gitignore + Makefile (Tasks 6.1, 6.2)
- ✅ API CI updates with content sync + path triggers (Task 6.5)
- ✅ E2E tests (Task 7.1)
- ✅ Roadmap and CLAUDE.md updates (Tasks 7.2, 7.3)
- ✅ Final smoke test (Task 7.4)
- ✅ PR opening (Task 7.5)

**Placeholder scan:** No "TBD"/"TODO"/"implement later" — every step has concrete code or commands.

**Type/path consistency:**
- `resume` collection always referenced as `'resume'` and entry `'profile'`
- Path `frontend/src/content/resume/profile.md` consistent
- `data-resume-trigger` attribute used identically in modal script, Nav, MobileNav
- `includeResumeModal` prop used identically in BaseLayout and resume.astro
- `assembleProjects()` and `projectsFS` referenced consistently
- The MobileNav prop change (`resumePath` → `isResumePage`) is reflected in both the Nav.astro caller and the MobileNav.tsx interface

**Notes for executor:**
- The plan assumes the executor has the LinkedIn export accessible at `/Users/rafaeldeabreugomes/Downloads/rafael_abreu_linkedin_profile.md`. If running on a different machine, the user must provide their own equivalent.
- The plan does **not** push to `main` or open the PR without explicit user consent (per project safety rules). Only the local branch commits run automatically.
