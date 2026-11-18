# Remove ubuntu-24.04 Runner Pins — Design Spec

> Reverts the PR #48 runner pins now that Playwright is 1.65.1. Touches `.github/workflows/`, `CLAUDE.md` and docs only: no `frontend/**`, so no deploy on merge. No human was available for brainstorming; uncertainties are in section 5.

---

## 1. Overview & Goals

PR #48 pinned the `Frontend E2E` job (`ci.yml`) and the `Deploy to Netlify` job (`deploy-frontend.yml`) to `ubuntu-24.04`, because Playwright 1.59.1 could not install Chromium on Ubuntu 26.04 (`ubuntu-latest`). PR #63 upgraded Playwright to 1.65.1. Goal: return both jobs to `ubuntu-latest`, matching every other job, and prove it in CI.

## 2. Changes

| File | Change |
|---|---|
| `.github/workflows/ci.yml` | `frontend-e2e`: `runs-on: ubuntu-latest`, "Pinned" comment removed |
| `.github/workflows/deploy-frontend.yml` | `deploy`: same |
| `CLAUDE.md` | Runner-pin bullet replaced; Known Follow-ups item removed; `brief-cta-click` added to Analytics events; note on running E2E locally as an agent |

## 3. Evidence

- `frontend/package-lock.json` resolves `playwright` and `@playwright/test` to 1.65.1 (and `playwright-core` 1.65.1).
- The installed `playwright-core` 1.65.1 (`lib/coreBundle.js`) detects Ubuntu 26.04 as an officially supported platform and has `ubuntu26.04-x64` entries for the Chromium downloads and for the `--with-deps` package lists.
- Playwright 1.61.0 release notes list Ubuntu 26.04 support (from the triage agent; not re-fetched).
- Local checks for the CLAUDE.md E2E note: `isRunByAgent()` is in `node_modules/astro/dist/cli/agent.js` and `cli/preview/index.js` backgrounds the server when it returns true (unless `--ignore-lock` or `ASTRO_PREVIEW_BACKGROUND` is set); `playwright.config.ts` has `reuseExistingServer: !process.env.CI`.

## 4. Verification

The `Frontend E2E` job on this PR must show the 26.04 runner image, a successful `npx playwright install --with-deps chromium`, and passing `build:pdf` and `test:e2e`; the PR body quotes the lines. If Chromium fails to install, revert the pins and report. `deploy-frontend.yml` cannot run from a PR; it runs the same install step on the same image, and the next production deploy is the proof.

## 5. Open questions

- Whether `ubuntu-latest` stays 26.04 (GitHub may change the label); the CI log is the only evidence.
- The Playwright 1.61.0 release-note claim was relayed, not re-read.
