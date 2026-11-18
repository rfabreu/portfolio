# Docs Reconcile: State as of 2026-11-17 — Implementation Plan

> **For agentic workers:** Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make `CLAUDE.md` and `docs/phase2-roadmap.md` match the verified state in the spec, without changing any source, config, workflow, or content file.

**Architecture:** Two Markdown files edited in place plus this spec and plan. No code, no deploy.

**Tech Stack:** Markdown, `git`, `gh` (read-only).

**Spec:** `docs/superpowers/specs/2026-11-17-docs-reconcile-state-design.md`
**Issues:** #30 (P10), #37 (P11), #34, #39
**Branch:** `docs/reconcile-state-2026-11`

---

## Task 1: Update `CLAUDE.md`

**Files:**
- Modify: `CLAUDE.md`

- [ ] **Step 1: Branch protection** — new check date, three required checks, strict on, admin enforcement off.
- [ ] **Step 2: Netlify line** — #34 closed, PR #44, previews pass.
- [ ] **Step 3: Phase 2 Status** — P10 done and open lists, Services owner question, P11 note, Dependabot line.
- [ ] **Step 4: Canonical Facts** — `.dev` and LinkedIn notes now match (PR #41).

## Task 2: Update `docs/phase2-roadmap.md`

**Files:**
- Modify: `docs/phase2-roadmap.md`

- [ ] **Step 1: P10 Status** — rewrite from the spec's verified state.
- [ ] **Step 2: P11** — state paragraph and URL bullet.
- [ ] **Step 3: Dependency Map and V1 Cleanup** — P10 row, Dependabot bullet.

## Task 3: Verify and review

- [ ] **Step 1:** re-run every `gh`/`git` command behind spec section 2.
- [ ] **Step 2:** independent reviewer re-verifies the diff; fix findings.
- [ ] **Step 3:** `git diff --name-only origin/main` lists only `CLAUDE.md` and `docs/**`.
- [ ] **Step 4:** push the branch, open the PR, wait for CI.
