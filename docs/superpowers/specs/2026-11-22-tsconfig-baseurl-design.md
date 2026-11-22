# Remove tsconfig `baseUrl` — Design Spec

> TS 6.0.3 is on `main` (PR #69). It deprecates `compilerOptions.baseUrl`; TS 7 drops it. Touches `frontend/tsconfig.json`, a Vitest guard, `CLAUDE.md` and docs. `frontend/**` changes, so merging deploys the frontend (static output unchanged).

## 1. Problem

`npx tsc --noEmit -p .` in `frontend/` reports `TS5101: Option 'baseUrl' is deprecated`. `astro check` does not report it, so CI misses it.

## 2. Changes

| File | Change |
|---|---|
| `frontend/tsconfig.json` | Remove `"baseUrl": "."`; `paths` targets become `./src/<dir>/*` |
| `frontend/tests/tsconfig-no-baseurl.test.ts` | Guard: no `baseUrl`, every target `./`-prefixed |
| `CLAUDE.md` | Runners bullet: `Deploy to Netlify` confirmed (PR #69 deploy, run 44553727930, 2026-11-19, `ubuntu-26.04`); note on `paths` without `baseUrl` and the `tsc` check |

## 3. Evidence

- The four aliases (`@components`, `@islands`, `@layouts`, `@styles`) have no imports in `frontend/src` or `frontend/tests` (grep), so resolution is not exercised; build and `astro check` must still pass.
- A trial of this change gave `astro check` 0 errors / 0 warnings / 83 hints, same as before, and a passing build.

## 4. Verification

`npm ci`, `astro check`, `npm run build`, `vitest run`, and `npx tsc --noEmit -p . | grep TS5101` empty. The known TS2307 errors for `.astro` imports in `tests/` are pre-existing and out of scope.
