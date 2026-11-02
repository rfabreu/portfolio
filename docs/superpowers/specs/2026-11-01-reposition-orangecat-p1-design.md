# P10 Phase 1: Reposition Portfolio for OrangeCat Labs — Design Spec

> Shift the portfolio's first impression from "hire me as an employee" to "Rafael, Co-Founder of OrangeCat Labs — let's build your next product." Experience and credentials stay on the site as credibility, not as the pitch. Phase 1 covers positioning, copy, CTAs, nav, SEO, and footer elevation.

---

## 1. Positioning Strategy

### Before (Employee-seeking)
- Hero: "Software Engineer building intelligent systems"
- Primary CTA: "View Projects"
- Secondary CTA: "Ask my AI →"
- Nav CTA: "LET'S TALK"
- Page title: "Rafael Abreu | Software Engineer"

### After (Studio founder driving leads)
- Hero: Rafael as Co-Founder of OrangeCat Labs, outcome-focused tagline
- Primary CTA: "Work with OrangeCat Labs" → `https://orangecatlabs.ca`
- Secondary CTA: "View Projects" (demoted, stays as anchor)
- Nav CTA: "WORK WITH US" → `https://orangecatlabs.ca`
- Page title: "Rafael Abreu | Co-Founder, OrangeCat Labs"

**Key principle:** The visitor's next step is always clear: go to OrangeCat Labs to start a project.

---

## 2. Hero Section (`Hero.astro`)

### Current
```
// HELLO WORLD
RAFAEL
ABREU
Software Engineer building intelligent systems…
[Go] [Python] [React] [Kubernetes] [AI/ML]
[View Projects]  [Ask my AI →]
```

### Proposed
```
// CO-FOUNDER · ORANGECAT LABS
RAFAEL
ABREU
I build solutions that drive results — full-stack applications,
cloud infrastructure, and AI-powered systems for businesses
ready to scale.
[Go] [Python] [React] [Kubernetes] [AI/ML]
[Work with OrangeCat Labs]  [View Projects]
```

### Changes
| Element | Before | After |
|---------|--------|-------|
| Mono label | `// HELLO WORLD` | `// CO-FOUNDER · ORANGECAT LABS` |
| Subtitle | "Software Engineer building **intelligent systems**…" | "I build solutions that drive results — **full-stack applications**, cloud infrastructure, and AI-powered systems for businesses ready to scale." |
| Primary CTA | "View Projects" → `#projects` | "Work with OrangeCat Labs" → `https://orangecatlabs.ca` (external, `target="_blank"`) |
| Secondary CTA | "Ask my AI →" button | "View Projects" → `#projects` (ghost/outline style) |
| Tech tags | Unchanged | Unchanged |
| H1 name | Unchanged | Unchanged |

### CTA Styling
- **Primary:** Keeps existing `bg-accent-indigo` solid button styling. Adds OrangeCat Labs logo (16×16px inline) before text for brand recognition.
- **Secondary:** Uses existing outline/ghost style (`border border-surface-border`). No icon.
- **"Ask my AI" button removal:** The hero no longer has a direct chat trigger. The chat widget FAB (bottom-right) remains the entry point. The `hero-chat-trigger` element ID and its listener in `ChatWidget.tsx` can be safely removed.

---

## 3. Navigation (`Nav.astro` + `MobileNav.tsx`)

### Desktop Nav
| Element | Before | After |
|---------|--------|-------|
| Nav links | About · Projects · Skills · Playground · Contact | About · Projects · Skills · Playground · Contact (**unchanged**) |
| CTA right group | `[RESUME]` `[LET'S TALK]` | `[RESUME]` `[WORK WITH US]` |
| "LET'S TALK" target | `/#contact` | `https://orangecatlabs.ca` (external) |
| "LET'S TALK" label | "LET'S TALK" | "WORK WITH US" |

### Mobile Nav
| Element | Before | After |
|---------|--------|-------|
| Bottom CTAs | `[RESUME]` `[LET'S TALK]` | `[RESUME]` `[WORK WITH US]` |
| Target | `/#contact` | `https://orangecatlabs.ca` (external) |

### External link handling
The "WORK WITH US" link opens `https://orangecatlabs.ca` in a new tab (`target="_blank"`, `rel="noopener noreferrer"`). An optional subtle external-link icon (↗) is appended for clarity.

---

## 4. Footer Elevation (`Footer.astro`)

The footer already has an OrangeCat Labs credit with logo and orange hover. Phase 1 elevates it:

| Element | Before | After |
|---------|--------|-------|
| Credit text | "Built by **OrangeCat Labs**" | "Co-Founded by Rafael · **OrangeCat Labs**" |
| Logo | 20×20px, subtle | 24×24px, slightly more prominent |
| Position | Center column | Center column (unchanged layout) |

Minimal change — the footer already does its job well.

---

## 5. SEO & Meta Tags (`BaseLayout.astro`)

| Meta | Before | After |
|------|--------|-------|
| `<title>` default | "Rafael Abreu \| Software Engineer" | "Rafael Abreu \| Co-Founder, OrangeCat Labs" |
| `<meta name="description">` | "Software Engineer building intelligent systems across the full stack…" | "Rafael Abreu — Co-Founder of OrangeCat Labs. Full-stack engineering, cloud infrastructure, and AI solutions for businesses ready to scale." |
| Open Graph (new) | None | `og:title`, `og:description`, `og:type`, `og:url`, `og:image` (use existing portfolio screenshot or OrangeCat branding) |

---

## 6. Files Changed

| File | Change Type | Description |
|------|-------------|-------------|
| `frontend/src/components/Hero.astro` | Modified | Rewrite mono label, subtitle, CTA buttons; remove `hero-chat-trigger` |
| `frontend/src/components/Nav.astro` | Modified | Change "LET'S TALK" → "WORK WITH US" with external link |
| `frontend/src/islands/MobileNav.tsx` | Modified | Change "LET'S TALK" → "WORK WITH US" with external link |
| `frontend/src/islands/ChatWidget.tsx` | Modified | Remove `hero-chat-trigger` listener (dead code after hero change) |
| `frontend/src/components/Footer.astro` | Modified | Update credit text wording |
| `frontend/src/layouts/BaseLayout.astro` | Modified | Update default title, description, add OG tags |
| `frontend/tests/e2e/resume.spec.ts` | Reviewed | Existing tests target `RESUME` button which is unchanged; no breakage expected |

---

## 7. What Does NOT Change in Phase 1

- **About section** — stays as-is (personal credibility, not the pitch)
- **Projects section** — no reframing yet (Phase 2 scope)
- **Skills section** — unchanged
- **Contact form** — stays as-is (Phase 2 converts to project brief)
- **Chatbot grounding** — no prompt changes (Phase 3 scope)
- **Resume modal/page** — stays as-is (Phase 2 demotes it)
- **Analytics tracking** — no new events (Phase 3 scope)
- **Section numbering** — unchanged

---

## 8. Constraints (from Issue #30)

- No email addresses shown publicly.
- Contact path: Netlify form + LinkedIn only.
- All work on feature branches via PR.
- OrangeCat Labs Upwork URL is Phase 2 (blocked on URL).
- No project is labeled as OrangeCat Labs work until decided.

---

## 9. Out of Scope (Phase 1)

- Services section (Phase 2)
- Case study reframing of projects (Phase 2)
- Contact form → project brief conversion (Phase 2)
- Chatbot persona/content updates (Phase 3)
- Analytics event tracking for CTAs (Phase 3)
- Page section reordering (Phase 3)
