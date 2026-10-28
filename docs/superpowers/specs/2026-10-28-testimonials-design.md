# P8: Testimonials Section — Design Spec

> Adds a structured testimonials content collection and responsive endorsement cards to the portfolio between Skills and Contact, with a zero-JS architecture, dark/light theme support, and a dormant state until real quotes are gathered.

---

## 1. Overview & Goals

- **Social Proof:** Provide credibility through authentic recommendations from engineering managers, collaborators, clients, and colleagues.
- **Single Source of Truth:** Testimonials are maintained as Markdown files in `frontend/src/content/testimonials/` with a strictly validated Zod schema, matching the pattern established by `projects` and `resume`.
- **Zero-JS by Default:** Built with static Astro components using pure CSS for responsive grid layout and optional touch scroll-snap, avoiding unnecessary client bundle overhead.
- **Dormant / Safe State:** The component cleanly renders nothing when the content collection is empty, ensuring production remains unaffected until real testimonials are ready to publish.
- **Future AI Grounding:** Architecture designed to sync directly into the Go API's embedded context (`api/internal/chat/`) to ground the chatbot in peer feedback.

---

## 2. Architecture & Component Hierarchy

```
frontend/src/content/testimonials/*.md
          │
          ▼ (Astro Content Loader: glob)
frontend/src/content.config.ts (Collection definition & Zod schema)
          │
          ▼
frontend/src/pages/index.astro
          │
          └── frontend/src/components/Testimonials.astro
                ├── SectionHeader.astro (Number, label, title)
                └── frontend/src/components/TestimonialCard.astro (Individual quote card)
                      ├── Quote icon (SVG)
                      ├── Markdown quote body
                      ├── Author avatar / SVG initials monogram fallback
                      ├── Author name & LinkedIn verification link
                      └── Role & company attribution
```

---

## 3. Content Collection Schema

In `frontend/src/content.config.ts`, define the `testimonials` collection:

```typescript
const testimonials = defineCollection({
  loader: glob({ pattern: '**/*.md', base: './src/content/testimonials' }),
  schema: z.object({
    name: z.string(),
    role: z.string(),                     // e.g. "Director of Software Engineering"
    company: z.string(),                  // e.g. "Nextologies"
    relationship: z.string(),             // e.g. "Managed Rafael directly", "Engineering Colleague"
    avatar: z.string().optional(),        // e.g. "/images/testimonials/avatar.jpg"
    linkedin: z.string().url().optional(),// Direct LinkedIn profile URL for social verification
    order: z.number().default(99),        // Display order
    featured: z.boolean().default(false), // Highlighted in featured placement
  }),
});
```

### Content File Example (`frontend/src/content/testimonials/example.md`)

```markdown
---
name: "Alex Rivera"
role: "Engineering Manager"
company: "Nextologies"
relationship: "Managed Rafael directly"
linkedin: "https://linkedin.com/in/example"
order: 1
featured: true
---

Rafael consistently demonstrated deep technical competence across our microservice infrastructure. He tackled complex Kubernetes and Go challenges with diligence and built tools our team relies on daily.
```

The quote content is written as the Markdown body (`await render(entry)` or standard Markdown parsing), enabling multi-paragraph quotes and italicized emphasis.

---

## 4. UI / UX Design System

### A. Layout & Placement
- **Location:** On `frontend/src/pages/index.astro` between `<Skills />` and `<Contact />`.
- **Container:** Wrapped in `<section id="testimonials" class="py-20 px-6">` with `<div class="container mx-auto max-w-6xl reveal">`.
- **Grid:** 
  - Mobile (<768px): Single-column cards with smooth horizontal scroll-snap or vertical stack.
  - Tablet/Desktop (≥768px): 2-column or 3-column responsive grid (`grid md:grid-cols-2 lg:grid-cols-3 gap-6`).

### B. Card Anatomy (`TestimonialCard.astro`)
1. **Header / Quote Accent:**
   - Subtle decorative quotation mark or icon in `text-accent-indigo/30` or `text-accent-purple/30`.
2. **Body:**
   - Quote text styled in `text-text-secondary text-sm md:text-base leading-relaxed italic`.
3. **Divider / Spacer:**
   - Subtle top border or padding before author attribution.
4. **Author Attribution:**
   - **Avatar:** Circular 48x48px image if `avatar` exists. If omitted or broken, renders a clean SVG monogram containing the person's initials with an accent gradient border (`bg-surface-border text-accent-indigo font-bold`).
   - **Name:** `text-text-primary font-semibold text-sm` with optional LinkedIn icon link.
   - **Role & Company:** `text-text-muted text-xs`.
   - **Relationship Badge:** Subtle pill badge (`text-[10px] font-mono uppercase tracking-wider text-accent-indigo bg-accent-indigo/10 px-2 py-0.5 rounded-tag`).

### C. Color & Theme Tokens
Utilizes Tailwind CSS v4 custom theme tokens defined in `global.css`:
- Background: `bg-surface`
- Border: `border-surface-border`
- Text: `text-text-primary`, `text-text-secondary`, `text-text-muted`
- Accents: `text-accent-indigo`, `text-accent-purple`
- Transitions: Seamless contrast across dark and light modes.

---

## 5. Dormant & Publishing Strategy

### Safe Dormant State (Pre-collection)
To allow code preparation without displaying an empty section on the live portfolio:

```astro
---
import { getCollection } from 'astro:content';

const allTestimonials = await getCollection('testimonials');
const testimonials = allTestimonials.sort((a, b) => a.data.order - b.data.order);

// Dormant Guard: Render nothing if no testimonials exist
if (testimonials.length === 0) {
  return null;
}
---
```

### Section Numbering Harmony
- **When testimonials exist:**
  - `01 // ABOUT`
  - `02 // PROJECTS`
  - `03 // EXPERTISE`
  - `04 // RECOMMENDATIONS` (Testimonials)
  - `05 // CONTACT`
- **When testimonials are dormant (`length === 0`):**
  - Contact remains dynamically set or labeled as `04 // CONTACT` so there is no skipped number.

---

## 6. AI Chatbot Integration (Future Phase)

Following the portfolio's single-source-of-truth grounding architecture:
1. `api/Makefile`:
   ```makefile
   sync-content:
       ...
       @mkdir -p internal/chat/embedded/testimonials
       @cp -r ../frontend/src/content/testimonials/*.md internal/chat/embedded/testimonials/ 2>/dev/null || true
   ```
2. `api/internal/chat/context.go`:
   - Add `assembleTestimonials()` to `SystemPrompt()`.
   - Allows Gemini to accurately answer queries like *"What have colleagues and supervisors praised Rafael for?"*

---

## 7. Files Changed / Added

| File | Change | Description |
|------|--------|-------------|
| `frontend/src/content.config.ts` | Modified | Add `testimonials` collection definition & Zod schema |
| `frontend/src/content/testimonials/` | New Dir | Markdown storage for testimonial entries |
| `frontend/src/components/TestimonialCard.astro` | New | Card component with avatar fallback & quote layout |
| `frontend/src/components/Testimonials.astro` | New | Section component with empty collection guard |
| `frontend/src/pages/index.astro` | Modified | Import and position `<Testimonials />` between Skills and Contact |
| `frontend/tests/testimonials.test.ts` | New | Vitest unit tests verifying schema, sorting, and fallback rendering |

---

## 8. Out of Scope

- Automated LinkedIn scraping or OAuth testimonial ingestion (manual review ensures quality).
- Third-party widget embeds (e.g. Senja, Testimonial.to) which introduce tracking scripts and latency.
- Star ratings (corporate and engineering testimonials rely on written recommendations, not consumer star ratings).
