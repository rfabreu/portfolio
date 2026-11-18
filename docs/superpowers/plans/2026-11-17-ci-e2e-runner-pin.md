# CI: Pin Playwright Jobs to Ubuntu 24.04 — Implementation Plan

> **For agentic workers:** Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make the required `Frontend E2E` check pass again, and keep the production frontend deploy working, by running both Playwright jobs on Ubuntu 24.04, where Playwright 1.59.1 is supported.

**Architecture:** One line each in `.github/workflows/ci.yml` and `.github/workflows/deploy-frontend.yml`, plus comments. No source, dependency, or content changes; no deploy (workflow files are outside both deploy path filters).

**Tech Stack:** GitHub Actions.

**Spec:** `docs/superpowers/specs/2026-11-17-ci-e2e-runner-pin-design.md`
**Branch:** `fix/ci-e2e-runner-pin`

---

## Task 1: Pin the runner

**Files:**
- Modify: `.github/workflows/ci.yml`
- Modify: `.github/workflows/deploy-frontend.yml`

- [ ] **Step 1:** In the `frontend-e2e` job, change `runs-on: ubuntu-latest` to `runs-on: ubuntu-24.04` and add a comment pointing at the Playwright upgrade follow-up.
- [ ] **Step 1b:** Same change for the `deploy` job in `deploy-frontend.yml`.
- [ ] **Step 2:** Confirm the job keys and `name:` values are unchanged (`Frontend`, `Frontend E2E`, `API`).
- [ ] **Step 3:** Validate the YAML parses.

## Task 2: Verify and open the PR

- [ ] **Step 1:** Adversarial review of the diff.
- [ ] **Step 2:** Push the branch, open a PR to `main`.
- [ ] **Step 3:** Confirm all three CI jobs pass on the PR, with `Frontend E2E` reaching and passing the Playwright tests.
