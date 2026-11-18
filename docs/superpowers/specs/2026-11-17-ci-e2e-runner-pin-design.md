# CI: Pin Playwright Jobs to Ubuntu 24.04 — Design Spec

> Unblocks the required `Frontend E2E` check, which fails on every PR since GitHub's `ubuntu-latest` label moved to Ubuntu 26.04, and stops the same failure from breaking the next production frontend deploy. Two-line workflow change; no source, dependency, or content changes.

---

## 1. Problem

`Frontend E2E` fails at the "Install Playwright Chromium" step on every PR (first seen on PR #46, 2026-11-18 04:32 UTC; a re-run failed the same way):

```
BEWARE: your OS is not officially supported by Playwright; installing dependencies for ubuntu26.04-x64 as a fallback.
Cannot install dependencies for ubuntu26.04-x64 with Playwright 1.59.1!
Error: ERROR: Playwright does not support chromium on ubuntu26.04-x64
```

`Frontend E2E` is a required check on `main`, so no PR can merge until this is fixed.

`.github/workflows/deploy-frontend.yml` runs the same sequence (`ubuntu-latest`, `npx playwright install --with-deps chromium`, `npm run build:pdf`). Its last run (2026-11-07) predates the label change, so the next push to `main` touching `frontend/**` would fail before deploying.

## 2. Verified Facts

| Fact | Source |
|---|---|
| The lockfile resolves `playwright`, `playwright-core`, and `@playwright/test` 1.59.1 (`package.json`: `^1.59.1`). | `frontend/package-lock.json`, CI log |
| `ubuntu-latest` now resolves to Ubuntu 26.04 (`Image: ubuntu-26.04` in the `Frontend` job log, 2026-11-18). | CI log, PR #46 |
| Playwright 1.59.1 and 1.60.0 have no `ubuntu26.04` platform entry; 1.61.0 and later (checked 1.61.0, 1.62.0, 1.65.1) do. | `playwright-core` package contents |
| The E2E job passed on `ubuntu-24.04` (image 20261101.360.1) on 2026-11-07, PR #45, with the same lockfile and steps. | job 128564690454 log |
| `Frontend` and `API` pass on Ubuntu 26.04. | PR #46 checks |
| `deploy-api.yml` does not install Playwright. | workflow file |

## 3. Decision

Pin the two jobs that install Playwright, `frontend-e2e` in `ci.yml` and `deploy` in `deploy-frontend.yml`, to `runs-on: ubuntu-24.04`, each with a comment explaining why.

- **Why pin, not upgrade Playwright now:** pinning returns both jobs to the environment they last passed on, with zero test-behaviour risk. Upgrading changes `frontend/package-lock.json`, which triggers a production frontend deploy, and deserves its own PR and review.
- **Why only these jobs:** the others pass on 26.04. Pinning them would hide real drift without fixing anything.
- **Job names unchanged:** `Frontend`, `Frontend E2E`, `API` are required by branch protection.
- **No deploy from this PR:** workflow files are outside both deploy path filters.

## 4. Follow-up

Upgrade Playwright to 1.61 or later (`npm update @playwright/test playwright` stays within the existing `^1.59.1` range; latest is 1.65.1), confirm Chromium installs on 26.04, then remove both pins. The Dependabot config proposed in PR #47 will surface the upgrade.

## 5. Open Questions

- How long GitHub keeps the `ubuntu-24.04` label. No deprecation notice was found; the follow-up removes the dependency on it.
- Strict required checks: open PRs (#45, #46, #47) must update from `main` after this merges before `Frontend E2E` reflects the fix.
