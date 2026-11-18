# CLAUDE.md Follow-ups from PRs #47–#50 — Implementation Plan

> **For agentic workers:** Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Apply the verified CLAUDE.md follow-ups from PRs #47, #48, #49 and #50.

**Architecture:** Edits to `CLAUDE.md` only, plus this plan and the spec. No source, workflow or content changes, so no deploy (neither deploy path filter matches `CLAUDE.md` or `docs/**`).

**Spec:** `docs/superpowers/specs/2026-11-18-claude-md-followups-design.md`
**Branch:** `docs/claude-md-followups-2026-11`

---

## Task 1: Edit CLAUDE.md

**Files:**
- Modify: `CLAUDE.md`

- [x] **Step 1:** Add `.github/` to Project Structure.
- [x] **Step 2:** Add the API limits paragraph under Build & Run → API (Go).
- [x] **Step 3:** Rewrite the Deployment CI line; add runner pin, dependency updates and analytics bullets.
- [x] **Step 4:** Update P10 in Phase 2 Status (analytics events done, PR #49).
- [x] **Step 5:** Add a "Known Follow-ups" section.

## Task 2: Verify and open the PR

- [ ] **Step 1:** Re-check every claim against `origin/main` and `gh`.
- [ ] **Step 2:** Adversarial review of the diff by a reviewer subagent; fix findings.
- [ ] **Step 3:** Push the branch, open a PR to `main`, check `gh pr checks`.
