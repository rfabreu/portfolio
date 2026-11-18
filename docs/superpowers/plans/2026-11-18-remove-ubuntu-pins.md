# Remove ubuntu-24.04 Runner Pins — Implementation Plan

> **For agentic workers:** Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Run the E2E and deploy jobs on `ubuntu-latest` again.

**Architecture:** Edits to two workflow files and `CLAUDE.md`, plus this plan and the spec. No deploy (neither deploy path filter matches).

**Spec:** `docs/superpowers/specs/2026-11-18-remove-ubuntu-pins-design.md`
**Branch:** `chore/remove-ubuntu-pins`

---

## Task 1: Edit workflows and docs

- [x] **Step 1:** `ci.yml` and `deploy-frontend.yml`: `runs-on: ubuntu-latest`, drop the "Pinned" comments.
- [x] **Step 2:** CLAUDE.md: runner bullet, follow-up removal, `brief-cta-click`, local E2E note.

## Task 2: Verify and open the PR

- [ ] **Step 1:** Validate workflow YAML; grep for leftover `ubuntu-24.04` / 1.59 references.
- [ ] **Step 2:** Adversarial review of the diff by a reviewer subagent; apply findings.
- [ ] **Step 3:** Push the branch, open a PR, read the `Frontend E2E` log (runner image, install, `build:pdf`, `test:e2e`), put the lines in the PR body.
