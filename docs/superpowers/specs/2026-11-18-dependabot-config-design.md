# Dependabot Config Tuning — Design Spec

> Tunes `.github/dependabot.yml` after the first batch of Dependabot PRs (#51–#59). Config and docs only: nothing under `frontend/**` or `api/**` changes, so neither deploy workflow triggers. No human was available for brainstorming; uncertainties are in section 5.

---

## 1. Overview & Goals

The first Dependabot run produced eight open PRs (#51–#58) plus one closed (#59). Three problems:

- **Action bumps arrive one PR per action.** Major bumps are not grouped, so #51–#54 are four separate PRs. They touch `.github/` only and never deploy.
- **TypeScript 7 PR (#59) cannot work.** `@astrojs/check` declares the peer range `typescript ^5.0.0 || ^6.0.0`.
- **Playwright can split.** `playwright` and `@playwright/test` must resolve to the same version. A minor/patch group covers them today, but a major would arrive as two separate PRs, and a group change by Dependabot could separate them.

## 2. Changes

| File | Change |
|---|---|
| `.github/dependabot.yml` | github-actions: group renamed `actions-minor-patch` → `actions`, `update-types: [major, minor, patch]` |
| | npm: `ignore` entry `typescript` with `versions: [">=7.0.0"]` and a comment naming the revisit trigger |
| | npm: new `playwright` group (`patterns: [playwright, "@playwright/test"]`), defined **before** `npm-minor-patch` |
| `CLAUDE.md` | "Dependency updates" bullet describes the above |

## 3. Decisions and Evidence

- **Group order.** GitHub's Dependabot options reference states: "If a dependency matches more than one rule, it's included in the first group that it matches." Config order is the natural reading of "first" (the docs do not define it). `playwright` is therefore listed above `npm-minor-patch`; otherwise a minor Playwright bump would land in `npm-minor-patch` and the dedicated group would only catch majors.
- **No `update-types` on the `playwright` group.** The docs say a group includes all SemVer levels by default, so majors group too. Both packages move in one PR at every level.
- **`versions`, not `update-types`, for the TypeScript ignore.** An `update-types: [version-update:semver-major]` ignore would also block TypeScript 6 (a major from 5.9). `versions: [">=7.0.0"]` blocks only 7+. Ignore rules do not suppress security updates.
- **Group rename.** `actions-minor-patch` no longer describes the group. Dependabot branch names embed the group name, so the next run opens a new PR; the four existing action PRs (#51–#54) are single-dependency PRs and unaffected.
- **Unchanged:** schedules, `open-pull-requests-limit`, commit-message prefixes, the gomod ecosystem.

## 4. Verification

- The file parses and validates against the SchemaStore `dependabot-2.0.json` schema (Ajv, draft-07): every key is a documented key.
- Dependabot itself is the only real test of grouping behaviour; the first weekly run after merge shows it (see Open questions).

## 5. Open Questions

1. The Dependabot docs do not say whether "first group" means config order, and do not say whether patterns and `update-types` in one group are ANDed. The config relies on config order only. If a minor Playwright PR appears in `npm-minor-patch`, move the `playwright` group or add `exclude-patterns` to `npm-minor-patch`.
2. The `>=7.0.0` form is an npm semver range, which Dependabot's npm ignore handling accepts per "use the standard pattern for the package manager". The docs give no `>=` example; the first run confirms it (no TypeScript 7 PR should appear).
3. Whether the owner wants the same "all update types" grouping for `gomod` is left unchanged.
