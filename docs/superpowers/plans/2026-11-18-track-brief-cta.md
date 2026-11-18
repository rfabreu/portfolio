# Plan: Track Brief CTA Clicks

File: docs/superpowers/plans/2026-11-18-track-brief-cta.md
Related spec: docs/superpowers/specs/2026-11-18-track-brief-cta-design.md
Related issues: none (follows up open question 5 of the analytics funnel spec, PR #49)
Date: 2026-11-18
Branch: feature/track-brief-cta

---

## Summary of Spec

One imperative event, `brief-cta-click`, with `location` ∈ `nav`, `mobile-nav`, `services`, `chat-error`, `resume-modal`, `resume-page-footer`. Never `preventDefault`, always behind `typeof umami !== 'undefined'`. No visible change.

## Tasks

### Task 1: Failing tests
- [x] Extend `frontend/tests/analytics-funnel-events.test.ts`: Services and Nav hooks, per-component listener scripts, resume modal and page, ChatWidget handler, `#contact` source scan, event list.
- [x] Extend `frontend/tests/analytics-funnel-mobile-nav.test.ts`: click tracking, menu still closes, no throw without `umami`, markup unchanged.
- [x] Run `npx vitest run tests/analytics`; confirm the new tests fail.

### Task 2: Implement
- [x] `Nav.astro`, `Services.astro`, `resume.astro`, `ResumeModal.astro`: `data-brief-cta` hook plus listener script.
- [x] `MobileNav.tsx`: `handleNavClick`.
- [x] `ChatWidget.tsx`: `trackBriefFromChat`.
- [x] `npx vitest run` passes.

### Task 3: Verify
- [x] In `frontend/`: `npm ci`, `npx astro check`, `npm run build`, `npm test`, `npm run build:pdf && npm run test:e2e`.
- [x] Inspect `dist/` for the hooks and handlers.

### Task 4: Review and PR
- [x] Adversarial review of the diff vs `origin/main` by a reviewer subagent; fix findings.
- [x] Push `feature/track-brief-cta`, open the PR to `main`, check `gh pr checks`.
- [x] CLAUDE.md follow-ups (listed in the PR body, not edited here): add `brief-cta-click` to the Analytics events list.
