# P10 Phase 3: Analytics Funnel Events — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add Umami custom events for the OrangeCat Labs funnel (Issue #30, P10 Phase 3): OrangeCat Labs clicks by placement, View Projects, LinkedIn clicks, contact form attempts and successes.

**Architecture:** Declarative `data-umami-event` / `data-umami-event-location` attributes on existing anchors; three tiny guarded `umami.track` scripts (Hero, Contact, success page) where declarative tracking is wrong. No visible change to copy, layout or styling (same-tab LinkedIn links on `/resume` wait for Umami's request; spec open question 2).

**Tech Stack:** Astro 6, React islands, Umami Cloud, Vitest.

**Spec:** `docs/superpowers/specs/2026-11-17-analytics-funnel-events-design.md`
**Issue:** #30
**Branch:** `feature/analytics-funnel-events`

---

## Task 1: Write failing tests (red)

**Files:**
- Create: `frontend/tests/analytics-funnel-events.test.ts`
- Create: `frontend/tests/analytics-funnel-mobile-nav.test.ts`

- [ ] **Step 1: Component attribute tests** (Astro container): Hero, Nav, Footer, About, Contact, ProjectLayout (with and without `slug`), ResumeModal, Resume page assert `data-umami-event="orangecat-click"` / `"linkedin-click"` with the expected `data-umami-event-location` on the right anchor.
- [ ] **Step 2: Source-scan tests**: every OrangeCat Labs or Rafael-LinkedIn `<a>` under `src/` carries `data-umami-event`; only `location` / `project` data keys exist; imperative scripts read no form values; every `umami.track` is guarded; Freckles credit markup has the attributes; `success.astro` fires `contact-form-success`.
- [ ] **Step 3: MobileNav test** (separate file so the `useState` mock is isolated): render with the menu forced open and assert `location="mobile-nav"`.
- [ ] **Step 4: Run `cd frontend && npm test` and confirm the new tests fail for the expected reason** (attributes missing), not for a typo in the test.

---

## Task 2: Declarative attributes (green)

**Files:**
- Modify: `frontend/src/components/Hero.astro`, `Nav.astro`, `Footer.astro`, `About.astro`, `Contact.astro`, `ResumeModal.astro`
- Modify: `frontend/src/islands/MobileNav.tsx`
- Modify: `frontend/src/layouts/ProjectLayout.astro`, `frontend/src/pages/project/[slug].astro`, `frontend/src/pages/resume.astro`
- Modify: `frontend/src/content/projects/freckles-design.md`

- [ ] **Step 1:** Add `data-umami-event="orangecat-click"` + `data-umami-event-location` to the eight OrangeCat Labs anchors (hero, nav, mobile-nav, footer, about, contact, project-cta, freckles-credit).
- [ ] **Step 2:** Add `data-umami-event="linkedin-click"` + `data-umami-event-location` to the five LinkedIn anchors (footer, contact, resume-modal, resume-page-header, resume-page-footer).
- [ ] **Step 3:** Add optional `slug` prop to `ProjectLayout.astro`, pass `slug={project.id}` from `[slug].astro`, render `data-umami-event-project={slug}` on the project CTA.
- [ ] **Step 4:** Do not touch classes, text, `href`, `target`, `rel` or `aria-label`.

---

## Task 3: Imperative events (green)

**Files:**
- Modify: `frontend/src/components/Hero.astro`, `Contact.astro`, `frontend/src/pages/success.astro`

- [ ] **Step 1: Hero** — hook attribute `data-view-projects` on the View Projects anchor; `<script>` adds a click listener calling `umami.track('view-projects-click', { location: 'hero' })` behind `typeof umami !== 'undefined'`. Never `preventDefault`.
- [ ] **Step 2: Contact** — `<script>` adds a `submit` listener on `form[name="contact"]` calling `umami.track('contact-form-submit')` with no data, behind the guard.
- [ ] **Step 3: Success** — `<script>` calls `umami.track('contact-form-success')` on `window.load`, behind the guard.
- [ ] **Step 4:** Run `cd frontend && npm test`; all new tests pass, existing tests unchanged.

---

## Task 4: Verify

- [ ] **Step 1:** `cd frontend && npm ci && npx astro check && npm run build && npm test` — all pass.
- [ ] **Step 2:** Inspect `frontend/dist/`: attributes on `index.html`, `project/freckles-design/index.html`, another `project/*/index.html` (has `data-umami-event-project`), `resume/index.html`, `success/index.html` (script present); `grep` the bundled JS for the MobileNav attributes and the three `umami.track` calls.
- [ ] **Step 3:** Confirm no visible change: `git diff` shows only attribute and script additions on existing elements.
- [ ] **Step 4:** Confirm no P10 violations: no pricing/budget/Upwork text added, no new OrangeCat Labs link, no email.

---

## Task 5: Review and PR

- [ ] **Step 1:** Adversarial review of the diff vs `origin/main` (missed CTAs, PII leaks, wrong Umami usage, visible UI changes, P10 violations, convention breaks). Fix valid findings.
- [ ] **Step 2:** Push `feature/analytics-funnel-events` only; open PR to `main` with the spec's open questions, review summary, "check events in the Umami dashboard after deploy", and CLAUDE.md follow-ups. Never merge.
- [ ] **Step 3:** Wait for CI (`Frontend`, `Frontend E2E`, `API`); fix failures.
