# Plan: Contact Notification Labels

File: docs/superpowers/plans/2026-11-18-contact-notification-labels.md
Related spec: docs/superpowers/specs/2026-11-18-contact-notification-labels-design.md
Related issues: none
Date: 2026-11-18
Branch: fix/contact-notification-labels

---

## Summary of Spec

Move the required asterisk out of the label text (CSS `::after` on `.field-required`), keep uppercase via `uppercase`, use sentence-case label text. No change to `name`, `required`, copy, layout or styling.

## Tasks

### Task 1: Failing tests
- [x] Add `frontend/tests/contact-labels.test.ts`.
- [x] Run it and confirm it fails.

### Task 2: Implement
- [x] `Contact.astro`: label text and classes.
- [x] `global.css`: `.field-required::after`.
- [x] `npx vitest run` passes.

### Task 3: Verify
- [x] In `frontend/`: `npm ci`, `npx astro check`, `npm run build`, `npm test`, `npm run build:pdf`, `npm run test:e2e`.
- [x] Inspect `dist/index.html` labels.
- [x] Before/after screenshots, both themes, desktop and mobile.

### Task 4: Review and PR
- [x] Adversarial review of the diff vs `origin/main`; fix findings.
- [x] Push `fix/contact-notification-labels`, open the PR, check `gh pr checks`.
- [ ] Owner: submit the live form after deploy and confirm the notification.
