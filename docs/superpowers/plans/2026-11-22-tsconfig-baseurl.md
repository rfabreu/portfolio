# Remove tsconfig `baseUrl` — Implementation Plan

> **For agentic workers:** Steps use checkbox (`- [x]`) syntax for tracking.

**Goal:** Clear TS5101 and keep it cleared.

**Spec:** `docs/superpowers/specs/2026-11-22-tsconfig-baseurl-design.md`
**Branch:** `chore/tsconfig-baseurl`

## Task 1: Edit

- [x] **Step 1:** `tsconfig.json`: drop `baseUrl`, `./`-prefix `paths` targets.
- [x] **Step 2:** Add `frontend/tests/tsconfig-no-baseurl.test.ts`.
- [x] **Step 3:** `CLAUDE.md`: Runners bullet and tsconfig note.

## Task 2: Verify and open the PR

- [x] **Step 1:** `npm ci`, `astro check`, build, Vitest, `tsc` grep for TS5101.
- [x] **Step 2:** Adversarial review of the diff; apply findings.
- [x] **Step 3:** Push, open PR with real output in the body.
