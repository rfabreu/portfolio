# Frontend runtime bumps — Implementation Plan

> **For agentic workers:** Steps use checkbox (`- [x]`) syntax for tracking.

**Goal:** Move CI and deploy to the Active LTS Node and take the astro patch.

**Spec:** `docs/superpowers/specs/2026-11-22-frontend-runtime-bumps-design.md`
**Branch:** `chore/frontend-runtime-bumps`

## Task 1: Edit

- [x] **Step 1:** Workflows: `node-version: 26` in `ci.yml` (2) and `deploy-frontend.yml` (1).
- [x] **Step 2:** `npm install astro@7.4.13` in `frontend/`.
- [x] **Step 3:** `CLAUDE.md`: Node line and Dependency updates bullet.

## Task 2: Verify and open the PR

- [x] **Step 1:** `npm ci`, `astro check`, build, tests, audit gate, E2E on Node 26.
- [x] **Step 2:** Adversarial review of the diff; apply findings.
- [x] **Step 3:** Push, open PR with real output in the body.
