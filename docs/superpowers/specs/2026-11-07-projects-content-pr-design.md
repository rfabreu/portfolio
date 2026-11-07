# Design Spec: Projects Content — Freckles Credit + URL Fixes

File: docs/superpowers/specs/2026-11-07-projects-content-pr-design.md
Related: Issue #30 (P10), Issue #37 (P11)
Created: 2026-11-07
Status: Approved for implementation

---

## Context

P10 (Issue #30) requires crediting OrangeCat Labs on the Freckles Design case study and fixing all `.dev` URL references. The `demo` URLs in `frontend/src/content/projects/` content are the only location where project links appear on the site, and they currently point to incorrect or stale URLs.

## Constraints (from P10)

- Only Freckles Design is attributed to OrangeCat Labs (owner decision). No other project may be credited to the studio.
- No budget or pricing language anywhere — not in copy, form fields, or the chatbot.
• Contact remains the Netlify form and LinkedIn only (https://linkedin.com/in/rafael-a-gomes). No public email.
- Work follows spec → plan → branch → PR, adversarial review, green CI, owner approval (CLAUDE.md).
- Never push to `main`; feature branch → PR only.

## Scope

### 1. Add "Built by OrangeCat Labs" credit to Freckles Design case study

File: `frontend/src/content/projects/freckles-design.md`

Add a short credit line immediately after the YAML frontmatter (before `## Overview`):

**Built by <a href="https://orangecatlabs.ca" target="_blank" rel="noopener noreferrer" class="oc-link">OrangeCat Labs</a>**

- Keeps it short (5 words) — the case study stays focused on what was done; the studio's own messaging lives on orangecatlabs.ca and in chatbot grounding (P10 Phase 3).
- Link uses `target="_blank" rel="noopener noreferrer"` to match the site's external-link convention (Footer.astro, Nav.astro).
- Styled via class `oc-link` (see Scope item 3) so it glows OrangeCat orange on hover, matching the footer's `oc-credit` pattern.

### 2. Fix `demo` URL in projects/portfolio.md

File: `frontend/src/content/projects/portfolio.md`

Current (incorrect):
```
demo: "https://rafaelabreu.dev"
```

Update to:
```
demo: "https://rafaelabreu.netlify.app"
```

### 3. Add OrangeCat link hover-glow CSS

File: `frontend/src/styles/global.css`

Add an `a.oc-link` rule following the exact Footer.astro `oc-credit` pattern (brand orange sampled from the OrangeCat Labs logo):

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

### 4. Audit and fix all `.dev` URL references in projects/ content files

Files checked (all 7 under `frontend/src/content/projects/`):
- `prjtskmngr.md`, `freckles-design.md`, `fptv.md`, `loudness-analyzer.md`, `wrapp.md`, `charge-it-up.md`, `portfolio.md`

Verification (`grep -rn 'rafaelabreu\.dev' frontend/src/content/`) found exactly two matches, both now fixed:
- `frontend/src/content/projects/portfolio.md` — `demo`
- `frontend/src/content/resume/profile.md` — `links.portfolio`

Freckles Design's own `demo: "https://frecklesdesign.ca/"` is the client's real domain and is left unchanged.

### 5. Fix profile.md LinkedIn and portfolio URL mismatches (P11 prep)

File: `frontend/src/content/resume/profile.md`

- `links.portfolio`: `https://rafaelabreu.dev` → `https://rafaelabreu.netlify.app`
- `links.linkedin`: `https://www.linkedin.com/in/rfabreu/` → `https://linkedin.com/in/rafael-a-gomes`

These are the known mismatches listed in CLAUDE.md Canonical Facts and issue #37. Fixing them here addresses both #30 and #37 in one deploy cycle (content under `resume/` deploys frontend + API).

## Verified Facts

| Fact | Source |
|------|--------|
| Only Freckles Design is attributed to OrangeCat Labs | Owner decision; CLAUDE.md Canonical Facts, commit 3edf2c3 |
| Production URL is https://rafaelabreu.netlify.app | CLAUDE.md, astro.config.mj `site`, README |
| `rafaelabreu.dev` is not the production URL | CLAUDE.md, docs/superpowers/specs/2026-11-07-docs-sync-p10-p11-design.md |
| `deploy-api.yml` triggers on `frontend/src/content/resume/**` and `frontend/src/content/projects/**` | `.github/workflows/deploy-api.yml` |
| Required checks: `Frontend` and `API`; `Frontend E2E` not required | CLAUDE.md, GitHub branch protection |
| Footer orange hover pattern: `--oc-orange: #f65203` (dark) / `#d94400` (light), `transition: color 300ms`, hover `color: var(--oc-orange)` | `frontend/src/components/Footer.astro` |

Not asserted (verified locally via grep — only the two `.dev` matches above exist in `frontend/src/content/`).

## Deploy Impact

- Frontend: all changed files are under `frontend/src/content/` plus `frontend/src/styles/global.css`, so this triggers the frontend deploy.
- API: Go binary rebuilt because `resume/profile.md` changes feed chatbot grounding via `make sync-content`.
- No source code or component changes — only markdown content files and one CSS rule.

## Acceptance Criteria

1. Freckles Design page displays "Built by OrangeCat Labs" with a clickable link to orangecatlabs.ca that glows orange on hover.
2. All `demo:` URLs in `projects/` content point to valid `.netlify.app` URLs.
3. `profile.md` `links.portfolio` and `links.linkedin` point to correct URLs.
4. No `.dev` URL references remain in `frontend/src/content/`.
5. CI passes: `astro check`, frontend `npm run build`, `go vet`, `go test`, `go build`.
6. Adversarial review confirms: no budget/pricing language, correct links only, no project other than Freckles credited to OrangeCat Labs.

## Files Changed (expected)

| File | Change |
|------|--------|
| `frontend/src/content/projects/freckles-design.md` | Add "Built by OrangeCat Labs" credit line |
| `frontend/src/content/projects/portfolio.md` | Fix `demo` URL (`.dev` → `.netlify.app`) |
| `frontend/src/content/projects/*.md` (all 7) | Audit + fix any `.dev` references (only portfolio.md had one) |
| `frontend/src/content/resume/profile.md` | Fix `links.portfolio` and `links.linkedin` URLs |
| `frontend/src/styles/global.css` | Add `a.oc-link` hover-glow rule |

Branch name: `feature/projects-content-freckles-credit-and-urls`
