# Plan: Fix React Hydration Mismatch

File: docs/superpowers/plans/2026-11-18-hydration-mismatch.md
Related spec: docs/superpowers/specs/2026-11-18-hydration-mismatch-design.md
Related issues: none
Date: 2026-11-18
Branch: fix/hydration-mismatch

---

## Summary of Spec

`ThemeToggle` reads `document` during its first render, so a light-theme visitor's first client render differs from the server HTML (always dark) and React logs #418. Read the theme with `useSyncExternalStore` and a `'dark'` server snapshot. Add an E2E regression test.

## Tasks

### Task 1: Reproduce
- [x] Playwright script over 5 pages × 3 theme setups × 2 reduced-motion values, in dev and production builds; record the mismatching component.
- [x] Capture before screenshots (pre-hydration and hydrated, both themes) and the per-frame theme/background samples.

### Task 2: Failing test
- [x] Add `frontend/tests/e2e/hydration.spec.ts`.
- [x] Run it against a build of `main`; confirm the light-theme cases fail.

### Task 3: Fix
- [x] `frontend/src/islands/ThemeToggle.tsx`: `useSyncExternalStore(subscribe, getTheme, getServerTheme)`.
- [x] Rebuild; the repro script reports 0 mismatches; the E2E spec passes; screenshots and nav bars match the baseline.

### Task 4: Verify
- [x] In `frontend/`: `npm ci`, `npx astro check`, `npm run build`, `npm test`, `npm run build:pdf`, `npm run test:e2e`.

### Task 5: Review and PR
- [x] Adversarial review of the diff vs `origin/main` by a reviewer subagent; fix findings.
- [x] Push `fix/hydration-mismatch`, open the PR to `main`, check `gh pr checks`.
