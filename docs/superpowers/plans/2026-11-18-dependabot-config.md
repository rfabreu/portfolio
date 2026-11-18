# Dependabot Config Tuning — Implementation Plan

> **For agentic workers:** Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Group GitHub Actions majors, keep Playwright packages together, and ignore TypeScript 7+.

**Architecture:** Edits to `.github/dependabot.yml` and `CLAUDE.md` only, plus this plan and the spec. No deploy (neither deploy path filter matches).

**Spec:** `docs/superpowers/specs/2026-11-18-dependabot-config-design.md`
**Branch:** `chore/dependabot-config`

---

## Task 1: Edit config and docs

- [x] **Step 1:** github-actions: rename group to `actions`, add `major` to `update-types`.
- [x] **Step 2:** npm: add `ignore` for `typescript` `>=7.0.0` with revisit comment.
- [x] **Step 3:** npm: add `playwright` group above `npm-minor-patch`.
- [x] **Step 4:** Update the CLAUDE.md "Dependency updates" bullet.

## Task 2: Verify and open the PR

- [x] **Step 1:** Validate YAML against the SchemaStore dependabot schema.
- [x] **Step 2:** Adversarial review of the diff by a reviewer subagent (Opus); findings applied: `exclude-patterns` on `npm-minor-patch`, corrected spec claims, deploy-only action risk recorded.
- [ ] **Step 3:** Push the branch, open a PR to `main`, check `gh pr checks`.
