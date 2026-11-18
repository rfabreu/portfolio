# CLAUDE.md Follow-ups from PRs #47–#50 — Design Spec

> Applies the "CLAUDE.md follow-ups" listed in the bodies of PRs #47, #48, #49 and #50 (merged 2026-11-18 UTC). Documentation only: no source, config, workflow, or content changes. No human was available for brainstorming; uncertainties are in section 5.

---

## 1. Overview & Goals

- **Accuracy:** each PR left CLAUDE.md edits for later. Applying them keeps CLAUDE.md true to `origin/main` (`5b3a4d1`).
- **Minimal:** only verified follow-ups, in the existing structure and voice.

## 2. Follow-ups and Evidence

| PR | Follow-up (from the PR body) | Verified against |
|---|---|---|
| #47 | CI line: audit gate in `Frontend`, pinned `govulncheck` in `API` | `.github/workflows/ci.yml`. The merged gate is `.github/scripts/npm-audit-gate.mjs` (allowlist with expiry), not the plain `npm audit --audit-level=high` the PR body describes, so CLAUDE.md describes the merged code. |
| #47 | "Dependency updates" note | `.github/dependabot.yml`: weekly, npm `/frontend`, gomod `/api`, github-actions `/`, minor/patch grouped, no auto-merge |
| #47 | Record GHSA-fc5q-39xc-vjpf | `ALLOWLIST` in the gate script, `expires: '2027-01-31'` |
| #47 | `.github/` in Project Structure | `ls .github` |
| #48 | None required in the body; the pins and the Playwright upgrade are recorded anyway as the task asked | `runs-on: ubuntu-24.04` in `ci.yml` (`frontend-e2e`) and `deploy-frontend.yml` (`deploy`); `node_modules/playwright` is 1.59.1 in `frontend/package-lock.json` |
| #49 | P10: analytics events done | `umami.track` and `data-umami-event` usages in `frontend/src`; `ProjectLayout`/`BaseLayout` script tags |
| #49 | Analytics note, `data-domains`, test that enforces tracking attributes | `frontend/src/layouts/*.astro`, `frontend/tests/analytics-funnel-events.test.ts` |
| #50 | API limits | `api/internal/chat/handler.go`, `provider.go`, `api/internal/game/handler.go`, `api/cmd/server/main.go`, `ChatWidget.tsx` (`MAX_INPUT_LENGTH`) |

## 3. Decisions

- Project Structure: the existing `.github/` line is widened to name the dependabot config and the audit gate script.
- Deployment gets the CI line rewrite plus three short bullets: runner pin, dependency updates, analytics.
- The API limits go under Build & Run → API (Go), after the env var line, as #50 suggested. Exact server timeout values stay in code (`newServer`), not in CLAUDE.md, so they cannot drift.
- A new "Known Follow-ups" section holds the Playwright upgrade and the allowlist expiry. It is placed before Docs Workflow.
- P10 in Phase 2 Status: "analytics events" moves from Open to done (PR #49).

## 4. Out of Scope

The `rafaelswe.netlify.app` site question, and the Services, page order and chatbot topics.

## 5. Open Questions

1. The PR #47 body says the audit step is `npm audit --audit-level=high`; the merged `ci.yml` uses the allowlist gate (commits `7df550c`, `49a3354`, `b24f157` landed in the PR before it merged). CLAUDE.md follows the merged code.
2. Whether the Playwright upgrade should be tracked as a GitHub issue instead of (or as well as) a CLAUDE.md line is the owner's call.
3. Event data keys, per-event `location` values and the preview-host behaviour are described in the analytics spec; CLAUDE.md links to it rather than duplicating the catalogue.
