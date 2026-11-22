# Plan: Theme Toggle Icon From First Paint

File: docs/superpowers/plans/2026-11-22-theme-toggle-prehydration-icon.md
Related spec: docs/superpowers/specs/2026-11-22-theme-toggle-prehydration-icon-design.md
Related issues: none
Date: 2026-11-22
Branch: fix/theme-toggle-prehydration-icon

---

## Summary of Spec

Render both toggle icons and pick one with CSS on `html[data-theme]`, so the static HTML is correct for light visitors before hydration and matches the client render. Keep the state-driven `aria-label`.

## Tasks

### Task 1: Failing tests
- [x] Extend `frontend/tests/e2e/hydration.spec.ts` with pre-hydration icon tests (island JS blocked, JS disabled) and post-hydration icon + label tests.
- [x] Run them against `main`'s `ThemeToggle`; confirm the light-theme pre-hydration test fails.

### Task 2: Fix
- [x] `frontend/src/islands/ThemeToggle.tsx`: render both SVGs (`aria-hidden`, `theme-icon-sun` / `theme-icon-moon`); `isDark = theme !== 'light'`.
- [x] `frontend/src/styles/global.css`: the two display rules.

### Task 3: Verify
- [ ] In `frontend/`: `npm ci`, `npx astro check`, `npm run build`, `npx vitest run`, E2E per CLAUDE.md.
- [x] Screenshots of the nav toggle area, light and dark, island JS blocked.

### Task 4: Review and PR
- [ ] Adversarial self-review of the diff; fix findings.
- [ ] Push `fix/theme-toggle-prehydration-icon`, open the PR to `main` (do not merge).
