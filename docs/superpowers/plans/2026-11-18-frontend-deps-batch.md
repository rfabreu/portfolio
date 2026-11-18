# Frontend Dependency Batch — Implementation Plan

> **For agentic workers:** Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** One PR with the dependency state of Dependabot #55, #56, #57 and #58, plus the Node `engines` bump.

**Architecture:** `frontend/package.json` ranges, regenerated `frontend/package-lock.json`, `CLAUDE.md`, spec and plan. Merging deploys the frontend.

**Spec:** `docs/superpowers/specs/2026-11-18-frontend-deps-batch-design.md`
**Branch:** `chore/frontend-deps-batch`

---

## Task 1: Combine the dependency state

- [x] **Step 1:** Read `gh pr diff` for #55–#58 and apply every `package.json` range change.
- [x] **Step 2:** Raise the `playwright` range to `^1.65.1` and `engines.node` to `>=22.23.4`.
- [x] **Step 3:** `npm install` in `frontend/`; compare resolved versions with each PR's lockfile.
- [x] **Step 4:** Update CLAUDE.md (Node requirement, Astro 7, Playwright pin and follow-up).

## Task 2: Verify

- [x] **Step 1:** `npm ci`, `astro check`, `npm run build`, `npm test`, `npm run build:pdf`.
- [x] **Step 2:** E2E with `CI=1` on a free port (port 4321 was held by another process).
- [x] **Step 3:** Browser check of islands vs a build of `origin/main`.
- [x] **Step 4:** Vitest, build and check under Node 22.23.5; npm audit gate.
- [x] **Step 5:** Lockfile security pass.

## Task 3: Review and open the PR

- [ ] **Step 1:** Adversarial review of the diff by a reviewer subagent; fix findings.
- [ ] **Step 2:** Push the branch, open a PR to `main`, check `gh pr checks`.
