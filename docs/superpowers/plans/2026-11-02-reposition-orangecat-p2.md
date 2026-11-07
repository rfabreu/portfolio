# P10 Phase 2: Conversion & Services — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement Phase 2 of the OrangeCat Labs repositioning (Issue #30) by adding a dedicated Services section, turning the contact form into a structured project brief, adding conversion CTAs to project case studies, and aligning the resume headline while demoting its nav styling.

**Architecture:** Astro 6 static components (`Services.astro`, `Contact.astro`, `ProjectLayout.astro`), markdown content updates (`profile.md`), and Vitest/Playwright test suites.

**Tech Stack:** Astro 6, Tailwind CSS v4, Netlify Forms, Vitest, Playwright.

**Spec:** `docs/superpowers/specs/2026-11-02-reposition-orangecat-p2-design.md`  
**Issue:** #30  
**Branch:** `feature/reposition-orangecat-p2`

---

## Phase 1: Services & Navigation

### Task 1: Build `Services.astro` Component

**Files:**
- Create: `frontend/src/components/Services.astro`

- [ ] **Step 1: Implement `Services.astro`**
  - Use `SectionHeader` with `number="02"`, `label="SERVICES"`, `title="What We Build"`.
  - Grid of 4 services:
    1. Full-Stack Software Development
    2. Cloud Architecture & DevOps
    3. AI Integration & Automation
    4. Systems Modernization & Optimization
  - Each card displays icon/tag, title, summary, key focus bullets, and an outbound inquiry link.

---

### Task 2: Mount `Services` and Update Navigation

**Files:**
- Modify: `frontend/src/pages/index.astro`
- Modify: `frontend/src/components/Nav.astro`
- Modify: `frontend/src/islands/MobileNav.tsx`

- [ ] **Step 1: Add `Services` to `Nav.astro` and `MobileNav.tsx`**
  - Add `{ label: 'Services', href: '/#services' }` to `navItems`.
  - Demote `RESUME` button in `Nav.astro` from a bordered button to a secondary text link with `data-resume-trigger`.

- [ ] **Step 2: Mount `<Services />` in `index.astro`**
  - Place `<Services />` with dividers right after `<Hero />` (or after `<About />`).

---

## Phase 2: Project Brief Form & Case Studies

### Task 3: Upgrade Contact Form to Project Brief

**Files:**
- Modify: `frontend/src/components/Contact.astro`

- [ ] **Step 1: Add project brief fields**
  - `name` (required text)
  - `email` (required email)
  - `project_type` (dropdown: Full-Stack Web App, Cloud & DevOps, AI & Automation, Modernization, Other)
  - `timeline` (dropdown: Immediate, 1–3 months, 3–6 months, Flexible)
  - `message` (required textarea)
  - Submit button: `Submit Project Brief →`
  - Update section header to `START A PROJECT` / `PROJECT BRIEF`.
  - No budget or pricing fields are included, per P10 constraints.

---

### Task 4: Add Conversion CTA to Case Studies & Align Resume

**Files:**
- Modify: `frontend/src/layouts/ProjectLayout.astro`
- Modify: `frontend/src/content/resume/profile.md`

- [ ] **Step 1: Add conversion block in `ProjectLayout.astro`**
  - Add callout at the end of the case study linking to `https://orangecatlabs.ca`.

- [ ] **Step 2: Update resume headline in `profile.md`**
  - Change `title: "Operations & Infrastructure Leader"` to `title: "Co-Founder, OrangeCat Labs | Full-Stack & Infrastructure Leader"`.

---

## Phase 3: Testing & Verification

### Task 5: Add Unit Tests & Run Verification Suite

**Files:**
- Create: `frontend/tests/reposition-services-brief.test.ts`

- [ ] **Step 1: Write Vitest unit tests**
  - Test `Services.astro` outputs all 4 service offerings.
  - Test `Contact.astro` contains project brief selectors (`project_type`, `timeline`) and does NOT contain a `budget` field.
  - Test `ProjectLayout.astro` renders conversion CTA.

- [ ] **Step 2: Run verification suite**
  - `npm test -- --run`
  - `npx astro check`
  - `npm run build`
  - `npm run test:e2e`
  - `cd ../api && make test`
