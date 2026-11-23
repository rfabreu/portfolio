# Frontend runtime bumps (Node 26, astro 7.4.13) — Design Spec

> Node 26 "Lithium" became Active LTS on 2026-10-28; Node 22 entered maintenance on 2025-10-21 (end of life 2027-04-30). Touches the three workflow `node-version` entries, `astro` and `CLAUDE.md`. `frontend/**` changes, so merging deploys the frontend (static output unchanged).

## 1. Evidence (checked 2026-11-22)

- Release schedule (`nodejs/Release` `schedule.json`): v26 LTS from 2026-10-28, maintenance 2027-10-20, EOL 2029-04-30. v24 went to maintenance on 2026-10-20. `nodejs.org/dist/index.json` marks v26.13.1 as `Lithium`.
- `astro@7.4.13` `engines.node` is `>=22.23.4`; `vitest@5.1.1` is `^22.12.0 || ^24.0.0 || >=26.0.0`; `@playwright/test@1.65.1` is `>=20`. All admit Node 26.
- No `.nvmrc` or `.node-version` exists. `netlify.toml` does not pin Node; the Netlify build machine is unchanged by this PR.
- `@astrojs/check@0.9.12` (latest) still peers `typescript ^5.0.0 || ^6.0.0`, so TypeScript 7.0.2 stays out and the dependabot ignore stays.

## 2. Changes

| File | Change |
|---|---|
| `.github/workflows/ci.yml`, `deploy-frontend.yml` | `node-version: 22` → `26` (3 entries) |
| `frontend/package.json`, lockfile | `astro` `^7.4.12` → `^7.4.13` (pulls `@astrojs/compiler-binding` 0.5.3) |
| `CLAUDE.md` | Node line: CI and deploy use Node 26; TS 7 re-check date |

`engines` stays `>=22.23.4`: the floor is Astro's own, and raising it would needlessly exclude Node 22/24 users.

## 3. Verification

`npm ci`, `astro check`, `npm run build`, `npm test`, the audit gate and E2E, all run locally on Node 26.13.1; the PR's CI jobs then run on 26 for the first time.
