# P8: Testimonials Section — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Lay the complete architectural foundation for the testimonials section (Issue #16) using Astro 6 content collections and responsive Astro components, while maintaining a dormant state that introduces zero visual change to the live site until authentic quotes are ready.

**Architecture:** A new `testimonials` Astro content collection in `frontend/src/content.config.ts` loads Markdown files from `frontend/src/content/testimonials/`. Two new Astro components (`Testimonials.astro` and `TestimonialCard.astro`) render the section with a clean responsive grid and avatar monogram fallback. If no testimonials are present in the collection, the component safely evaluates to `null`, ensuring production integrity.

**Tech Stack:** Astro 6, Tailwind CSS v4, Zod, Vitest.

**Spec:** `docs/superpowers/specs/2026-10-28-testimonials-design.md`  
**Issue:** #16  
**Branch:** `feature/testimonials-architecture`

---

## Phase 1: Content Collection & Schema

### Task 1.1: Add `testimonials` collection to `frontend/src/content.config.ts`

**Files:**
- Modify: `frontend/src/content.config.ts`

- [ ] **Step 1: Add the `testimonials` collection definition and export:**

```typescript
const testimonials = defineCollection({
  loader: glob({ pattern: '**/*.md', base: './src/content/testimonials' }),
  schema: z.object({
    name: z.string(),
    role: z.string(),
    company: z.string(),
    relationship: z.string(),
    avatar: z.string().optional(),
    linkedin: z.string().url().optional(),
    order: z.number().default(99),
    featured: z.boolean().default(false),
  }),
});

export const collections = { projects, resume, testimonials };
```

- [ ] **Step 2: Create directory `frontend/src/content/testimonials/` and `.gitkeep`:**

Ensure git tracks the empty directory without failing the glob loader.

---

### Task 1.2: Add unit tests for schema validation

**Files:**
- Create: `frontend/tests/testimonials-schema.test.ts`

- [ ] **Step 1: Write Vitest tests verifying the testimonials Zod schema:**
  - Validates full payload with name, role, company, relationship, avatar, and LinkedIn.
  - Validates defaults for `order` (99) and `featured` (false).
  - Validates rejection of missing required fields (e.g. name or role).
  - Validates rejection of invalid URLs in `linkedin`.

- [ ] **Step 2: Run test suite:**
  `cd frontend && npm test -- --run tests/testimonials-schema.test.ts`

---

## Phase 2: Component Architecture

### Task 2.1: Build `TestimonialCard.astro`

**Files:**
- Create: `frontend/src/components/TestimonialCard.astro`

- [ ] **Step 1: Implement card component:**
  - Props: `name`, `role`, `company`, `relationship`, `avatar` (optional), `linkedin` (optional).
  - Body slot: `<slot />` for the rendered quote content.
  - SVG Monogram Generator: Extracts initials from `name` (e.g., "Carla Antunes" → "CA") and renders a styled SVG/CSS badge when `avatar` is missing.
  - Social link: Optional external link to LinkedIn profile with accessible `aria-label`.
  - Responsive styling matching Tailwind v4 surface tokens (`bg-surface`, `border-surface-border`, `rounded-card`, `p-6`).

---

### Task 2.2: Build `Testimonials.astro`

**Files:**
- Create: `frontend/src/components/Testimonials.astro`

- [ ] **Step 1: Implement section wrapper with dormant guard:**
  - Queries `getCollection('testimonials')`.
  - **Dormant Guard:** `if (testimonials.length === 0) return null;`
  - Sorts entries by `data.order`.
  - Uses `SectionHeader.astro` with `number="04"`, `label="RECOMMENDATIONS"`, `title="What Colleagues Say"`.
  - Renders responsive grid: `grid md:grid-cols-2 lg:grid-cols-3 gap-6`.

---

## Phase 3: Page Integration & Verification

### Task 3.1: Mount in `index.astro`

**Files:**
- Modify: `frontend/src/pages/index.astro`

- [ ] **Step 1: Import and place `Testimonials` between `Skills` and `Contact`:**

```astro
    <Skills />
    <GradientDivider />
    <Testimonials />
    <Contact />
```

Note: Because `Testimonials.astro` returns `null` when 0 entries exist, the extra `GradientDivider` is wrapped conditionally inside `Testimonials.astro` or kept clean.

- [ ] **Step 2: Run full validation suite:**
  - `cd frontend && npm test -- --run`
  - `cd frontend && npx astro check`
  - `cd frontend && npm run build`

---

## Phase 4: Activation Checklist (When Quotes are Collected)

Once authentic quotes from colleagues or managers are ready:
1. Drop `.md` files into `frontend/src/content/testimonials/<author-slug>.md`.
2. Add optional square avatar images to `frontend/public/images/testimonials/`.
3. (Optional) Run `cd api && make sync-content` to incorporate quotes into the AI chatbot prompt.
4. Run `npm run build` and visually confirm in browser.
