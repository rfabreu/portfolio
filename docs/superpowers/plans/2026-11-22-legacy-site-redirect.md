# Redirect the legacy `rafaelswe` site — Implementation Plan

> **For agentic workers:** Steps use checkbox (`- [x]`) syntax for tracking.

**Goal:** 301 every path on `rafaelswe.netlify.app` to production; previews and production unaffected.

**Spec:** `docs/superpowers/specs/2026-11-22-legacy-site-redirect-design.md`
**Branch:** `fix/legacy-site-redirect`

## Task 1: Edit

- [x] **Step 1:** Root `netlify.toml`: add the host-scoped `[[redirects]]` rule.
- [x] **Step 2:** Add `frontend/tests/legacy-site-redirect.test.ts`.
- [x] **Step 3:** `CLAUDE.md`: Deployment bullet on `rafaelswe`.

## Task 2: Verify and open the PR

- [x] **Step 1:** `npm ci`, `astro check`, build, Vitest.
- [x] **Step 2:** Adversarial review of the diff; apply findings.
- [ ] **Step 3:** Push, open PR with the post-merge checklist; confirm the deploy preview serves 200 HTML.
