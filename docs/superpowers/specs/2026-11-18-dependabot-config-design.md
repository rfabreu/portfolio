# Dependabot Config Tuning — Design Spec

> Tunes `.github/dependabot.yml` after the first batch of Dependabot PRs (#51–#59). Config and docs only: nothing under `frontend/**` or `api/**` changes, so neither deploy workflow triggers. No human was available for brainstorming; uncertainties are in section 5.

---

## 1. Overview & Goals

The first Dependabot run produced eight open PRs (#51–#58; #51–#54, the four action PRs, have since merged) plus one closed (#59). Three problems:

- **Action bumps arrive one PR per action.** Major bumps are not grouped, so #51–#54 were four separate PRs. They touch `.github/` only.
- **TypeScript 7 PR (#59) cannot work.** `@astrojs/check` declares the peer range `typescript ^5.0.0 || ^6.0.0` (0.9.12, the latest, still does).
- **Playwright can split.** `playwright` and `@playwright/test` must resolve to the same version. Today the `npm-minor-patch` group covers both, but a major would arrive as two separate PRs.

Related, and recorded in CLAUDE.md: when several npm PRs conflict on `frontend/package-lock.json`, they are combined into one PR, because every `frontend/**` merge is a production deploy.

## 2. Changes

| File | Change |
|---|---|
| `.github/dependabot.yml` | github-actions: group renamed `actions-minor-patch` → `actions`, `update-types: [major, minor, patch]` |
| | npm: `ignore` entry `typescript` with `versions: [">=7.0.0"]` and a comment naming the revisit trigger |
| | npm: new `playwright` group (`patterns: [playwright, "@playwright/test"]`), defined before `npm-minor-patch`; `npm-minor-patch` gets `exclude-patterns` for the same two packages |
| `CLAUDE.md` | "Dependency updates" bullet describes the above, and states the combine-conflicting-PRs rule |

## 3. Decisions and Evidence

- **Two mechanisms for Playwright.** GitHub's options reference states: "If a dependency matches more than one rule, it's included in the first group that it matches." Config order is the natural reading of "first", but the docs do not define it, and reviewing Dependabot's source suggests a package can be placed in every group whose rules match (a group without `patterns`, like `npm-minor-patch`, loses no specificity contest). So the `playwright` group is listed first **and** `npm-minor-patch` excludes both packages, which makes the outcome independent of order.
- **No `update-types` on the `playwright` group.** The docs say a group includes all SemVer levels by default, so majors group too.
- **`versions`, not `update-types`, for the TypeScript ignore.** An `update-types: [version-update:semver-major]` ignore would also block TypeScript 6 (a major from 5.9). `versions: [">=7.0.0"]` blocks only 7+. The docs show `versions: [">=1.0.0"]`-style ranges for ignore. Ignore rules apply to security updates as well as version updates ("ignore those dependencies when it opens pull requests for version updates and security updates"), so a future TypeScript 7 security fix would not be proposed either; this is acceptable while TypeScript 7 is unsupported.
- **Group rename.** `actions-minor-patch` no longer describes the group. Dependabot branch names embed the group name, so the next run opens a new `actions` PR.
- **Unchanged:** schedules, `open-pull-requests-limit`, commit-message prefixes, the gomod ecosystem.

## 4. Verification

- The file parses and validates against the SchemaStore `dependabot-2.0.json` schema (Ajv, draft-07): every key is a documented key.
- Dependabot itself is the only real test of grouping behaviour. While PR #55 (`npm-minor-patch`, already containing both Playwright packages) is open, Dependabot treats those packages as handled by that PR, so the `playwright` group will not show anything until #55 is closed or superseded. The combined dependency PR (`chore/frontend-deps-batch`) supersedes #55; the first weekly run after it merges is the first real test.

## 5. Open Questions

1. **Existing action PRs.** #51–#54 (setup-go 6, setup-node 7, actions-netlify 4, checkout 7) were merged one by one before this change, so there are no single-action PRs left to duplicate; the next grouped `actions` PR will only carry later bumps.
2. **Deploy-only actions.** `nwtgck/actions-netlify` is used only in `deploy-frontend.yml`; no CI job runs it. A grouped major bump can pass all three required checks and still break the next production deploy; this has already happened in principle with `actions-netlify@v4` (#53), which first runs in the deploy that follows the combined frontend dependency PR. The owner decides whether to ignore majors of that action or review its bumps by reading the release notes.
3. **Ignore semantics.** If a TypeScript 7 PR appears after the next run, the `versions` form was not honoured.
4. **gomod** is intentionally unchanged (still minor and patch only).
