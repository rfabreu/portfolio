# Plan: Projects Content — Freckles Credit + URL Fixes

File: docs/superpowers/plans/2026-11-07-projects-content-pr.md
Related spec: docs/superpowers/specs/2026-11-07-projects-content-pr-design.md
Related issues: #30 (P10), #37 (P11)
Date: 2026-11-07
Branch: feature/projects-content-freckles-credit-and-urls

---

## Summary of Spec

Add a short "Built by OrangeCat Labs" credit to the Freckles Design case study, fix all `.dev` URL references in projects content and resume profile to use rafaelabreu.netlify.app, fix the LinkedIn URL in profile.md, and add an orange hover-glow CSS rule matching the footer's `oc-credit` pattern.

## Governance Note

CLAUDE.md is the governing document for this work. All references to agent workflow rules, CI check requirements (Frontend + API required; Frontend E2E not required), branch protection (no direct push to `main`), and the adversarial review process in CLAUDE.md apply. If any conflict arises between this plan and CLAUDE.md, CLAUDE.md takes precedence.

## Approach

### Step 1: Inventory — grep for all `.dev` references

Run:
```
grep -rn 'rafaelabreu\.dev' frontend/src/content/ 2>&1 || true
```
This confirms the complete list of files needing fixes. Expected: `portfolio.md` (`demo`) and `profile.md` (`links.portfolio`).

### Step 2: Add Freckles Design credit — short form only

Edit `frontend/src/content/projects/freckles-design.md`. Insert after the YAML frontmatter closing `---` and before `## Overview`:

```
**Built by <a href="https://orangecatlabs.ca" target="_blank" rel="noopener noreferrer" class="oc-link">OrangeCat Labs</a>**
```

Blank line after the credit for visual separation from the body.

### Step 3: Add orange glow CSS for the OrangeCat link

Edit `frontend/src/styles/global.css`. Append the `a.oc-link` rule following the Footer.astro `oc-credit` pattern exactly:

```css
a.oc-link {
  --oc-orange: #f65203;
  color: inherit;
  text-decoration: none;
  transition: color 300ms ease;
}
html[data-theme='light'] a.oc-link {
  --oc-orange: #d94400;
}
a.oc-link:hover,
a.oc-link:focus-visible {
  color: var(--oc-orange);
}
```

### Step 4: Fix URLs — all affected files

- `frontend/src/content/projects/portfolio.md`: `demo: "https://rafaelabreu.dev"` → `demo: "https://rafaelabreu.netlify.app"`
- `frontend/src/content/resume/profile.md`:
  - `links.portfolio`: `https://rafaelabreu.dev` → `https://rafaelabreu.netlify.app`
  - `links.linkedin`: `https://www.linkedin.com/in/rfabreu/` → `https://linkedin.com/in/rafael-a-gomes`

### Step 5: Verify no `.dev` references remain

```
grep -rn 'rafaelabreu\.dev' frontend/src/content/ 2>&1 || true
```
Expected output: empty. If any hits remain, fix them before committing.

### Step 6: Adversarial review (bug hunt)

Run the adversarial checklist from CLAUDE.md on every changed file:
- No stale `.dev` references (typos, case variants).
- No budget/pricing language introduced.
- Markdown formatting valid; YAML frontmatter unchanged except intended fields.
- The Freckles credit mentions only Freckles Design as OrangeCat work.
- Links use `target="_blank" rel="noopener noreferrer"` and correct protocols.
- `deploy-api.yml` path filters will trigger (both `resume/` and `projects/` are in the glob) — confirmed.

### Step 7: CI verification

- `cd frontend && npx astro check` — 0 errors.
- `cd frontend && npm run build` — completes.
- `cd api && go vet ./... && go test ./... && go build ./cmd/server` — all pass.

### Step 8: Commit, push feature branch, open PR

Branch: `feature/projects-content-freckles-credit-and-urls`

Commit message:
```
feat(projects): add OrangeCat credit to Freckles case study and fix URL references
```

Commit contents:
- `docs/superpowers/specs/2026-11-07-projects-content-pr-design.md` (new)
- `docs/superpowers/plans/2026-11-07-projects-content-pr.md` (new)
- `frontend/src/content/projects/freckles-design.md`
- `frontend/src/content/projects/portfolio.md`
- `frontend/src/content/resume/profile.md`
- `frontend/src/styles/global.css`

PR body should:
- Reference spec and plan file paths.
- List all changed files with brief descriptions.
- Note no source code or component changes — only markdown content files and one CSS rule.
- Reference issues #30 and #37.
- Request adversarial review before merge.

## Pre-commit Checklist (per CLAUDE.md)

- [x] Spec exists and is referenced in commit/PR
- [x] Plan exists and follows spec
- [x] All files pass `astro check` (frontend CI)
- [x] API builds successfully (API CI — triggered by resume content change)
- [x] No `.dev` references remain in `frontend/src/content/`
- [ ] Adversarial review complete — no budget/pricing language, correct links only
- [ ] PR title/body follow conventions in CLAUDE.md
