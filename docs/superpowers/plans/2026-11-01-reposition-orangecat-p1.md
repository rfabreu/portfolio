# P10 Phase 1: Reposition Portfolio for OrangeCat Labs — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Reposition the portfolio homepage from an individual employment pitch to presenting Rafael as Co-Founder of OrangeCat Labs, shifting the primary call to action to `https://orangecatlabs.ca` while keeping technical depth as credibility.

**Architecture:** Update static Astro components (`Hero.astro`, `Nav.astro`, `Footer.astro`, `BaseLayout.astro`) and React islands (`MobileNav.tsx`, `ChatWidget.tsx`) with outcome-driven positioning, updated meta/OG tags, and outbound CTAs to OrangeCat Labs.

**Tech Stack:** Astro 6, React, Tailwind CSS v4, Vitest, Playwright.

**Spec:** `docs/superpowers/specs/2026-11-01-reposition-orangecat-p1-design.md`  
**Issue:** #30  
**Branch:** `feature/reposition-orangecat-p1`

---

## Phase 1: Component & Island Updates

### Task 1: Update `Hero.astro` and clean up `ChatWidget.tsx`

**Files:**
- Modify: `frontend/src/components/Hero.astro`
- Modify: `frontend/src/islands/ChatWidget.tsx`

- [ ] **Step 1: Update `Hero.astro`**
  - Update mono header from `// HELLO WORLD` to `// CO-FOUNDER · ORANGECAT LABS`.
  - Update subtitle to: `I build solutions that drive results — <span class="text-accent-purple-light">full-stack applications</span>, cloud infrastructure, and AI-powered systems for businesses ready to scale.`
  - Replace primary CTA with external link to `https://orangecatlabs.ca`:
    - Label: `Work with OrangeCat Labs`
    - Class: `inline-flex items-center gap-2 px-6 py-2.5 bg-accent-indigo rounded-btn text-white text-sm font-semibold hover:bg-accent-indigo/90 transition-colors duration-200`
    - Icon: 18x18px inline OrangeCat logo with `rounded-sm object-contain p-0.5 bg-white/90`
    - Target: `target="_blank" rel="noopener noreferrer"`
  - Demote secondary CTA to "View Projects":
    - Target: `#projects`
    - Class: `px-6 py-2.5 border border-surface-border rounded-btn text-text-secondary text-sm font-medium hover:border-accent-indigo/40 hover:text-accent-indigo transition-colors duration-200`
  - Remove `#hero-chat-trigger` button.

- [ ] **Step 2: Clean up `ChatWidget.tsx`**
  - Remove the `useEffect` listening to `document.getElementById('hero-chat-trigger')`.

---

### Task 2: Update Navigation in `Nav.astro` and `MobileNav.tsx`

**Files:**
- Modify: `frontend/src/components/Nav.astro`
- Modify: `frontend/src/islands/MobileNav.tsx`

- [ ] **Step 1: Update desktop CTA in `Nav.astro`**
  - Change `LET'S TALK` button to:
    - Text: `WORK WITH US`
    - Href: `https://orangecatlabs.ca`
    - Target: `target="_blank" rel="noopener noreferrer"`
    - Styling: keep `px-4 py-1.5 bg-accent-indigo rounded-btn text-white text-xs tracking-wide hover:bg-accent-indigo/90 transition-colors duration-200`

- [ ] **Step 2: Update mobile CTA in `MobileNav.tsx`**
  - Change `LET'S TALK` button to:
    - Text: `WORK WITH US`
    - Href: `https://orangecatlabs.ca`
    - Target: `target="_blank" rel="noopener noreferrer"`
    - Styling: keep `px-6 py-2 bg-accent-indigo rounded-lg text-white text-sm`

---

### Task 3: Elevate Footer Credit in `Footer.astro`

**Files:**
- Modify: `frontend/src/components/Footer.astro`

- [ ] **Step 1: Update credit wording**
  - Change `<span>Built by <span class="font-semibold">OrangeCat Labs</span></span>` to `<span>Co-Founded by Rafael · <span class="font-semibold">OrangeCat Labs</span></span>`.

---

### Task 4: Enhance SEO & Open Graph in `BaseLayout.astro`

**Files:**
- Modify: `frontend/src/layouts/BaseLayout.astro`

- [ ] **Step 1: Update default title and description**
  - Title default: `Rafael Abreu | Co-Founder, OrangeCat Labs`
  - Description default: `Rafael Abreu — Co-Founder of OrangeCat Labs. Full-stack engineering, cloud infrastructure, and AI solutions for businesses ready to scale.`

- [ ] **Step 2: Add Open Graph and Twitter Card tags in `<head>`**
  - `og:type` -> `website`
  - `og:title` -> `{title}`
  - `og:description` -> `{description}`
  - `og:image` -> `/images/portfolio_website.png`
  - `twitter:card` -> `summary_large_image`

---

## Phase 2: Tests & Verification

### Task 5: Add Unit Tests & Verify E2E

**Files:**
- Create: `frontend/tests/reposition-hero-nav.test.ts`

- [ ] **Step 1: Write Vitest test suite**
  - Verify `Hero.astro` outputs `// CO-FOUNDER · ORANGECAT LABS`, the new tagline, and outbound link to `orangecatlabs.ca`.
  - Verify `Nav.astro` outputs `WORK WITH US` linking to `https://orangecatlabs.ca`.
  - Verify `Footer.astro` outputs `Co-Founded by Rafael`.

- [ ] **Step 2: Run test suite and check builds**
  - `npm test -- --run`
  - `npx astro check`
  - `npm run build`
  - `npm run test:e2e`
  - `cd ../api && make test`
