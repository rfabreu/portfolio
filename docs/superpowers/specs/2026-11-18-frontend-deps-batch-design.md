# Frontend Dependency Batch (#55 + #56 + #57 + #58) — Design Spec

> Combines four conflicting Dependabot PRs into one. Merging touches `frontend/**`, so it **deploys the frontend to production**. No human was available for brainstorming; uncertainties are in section 6.

---

## 1. Overview & Goals

Dependabot PRs #55 (npm-minor-patch group), #56 (`@astrojs/react` 5 → 7), #57 (`vitest` 4 → 5) and #58 (`@types/node` 25 → 26) each pass CI alone, but all four edit `frontend/package-lock.json`, so only one could merge before the others needed rebasing, and each merge is a production deploy. One combined PR means one deploy and one lockfile.

## 2. Changes

`frontend/package.json` ranges (taken from the four PR diffs):

| Package | From | To | PR |
|---|---|---|---|
| `@astrojs/check` | ^0.9.9 | ^0.9.12 | #55 |
| `@tailwindcss/vite` | ^4.2.2 | ^4.3.5 | #55 |
| `@types/react` / `@types/react-dom` | ^19.2.14 / ^19.2.3 | ^19.3.0 / ^19.3.0 | #55 |
| `astro` | ^7.4.7 | ^7.4.12 | #55 |
| `react` / `react-dom` | ^19.2.4 | ^19.3.1 | #55 |
| `@playwright/test` | ^1.59.1 | ^1.65.1 | #55 |
| `serve-handler` | ^6.1.7 | ^6.1.8 | #55 |
| `@astrojs/react` | ^5.0.2 | ^7.1.9 | #56 |
| `vitest` | ^4.1.5 | ^5.1.0 | #57 |
| `@types/node` | ^25.6.0 | ^26.6.6 | #58 |
| `playwright` | ^1.59.1 | ^1.65.1 | not in #55 (see below) |
| `engines.node` | >=22.12.0 | >=22.23.4 | new |

Also: `frontend/package-lock.json` regenerated with `npm install`; `CLAUDE.md` (Node requirement, "Astro 6" → "Astro 7", Playwright pin and follow-up wording).

## 3. Decisions

- **`playwright` range.** #55 moved the lockfile to 1.65.1 but left the `playwright` range at `^1.59.1` (the lock already satisfied it). `playwright` and `@playwright/test` must match (`generate-pdf.mjs` imports `playwright`; E2E uses `@playwright/test`), so both ranges now say `^1.65.1`.
- **`tailwindcss` range stays `^4.2.2`**, as in #55 (the lock resolves 4.3.5). The owner can raise it to match `@tailwindcss/vite`.
- **`engines.node` ≥ 22.23.4.** Astro 7.4.7 (already on `main`) and 7.4.12 both declare `>=22.23.4`; `main` said `>=22.12.0`, which was stale. CI and `deploy-frontend.yml` use `node-version: 22` (latest 22.x, 22.23.5 at time of writing).
- **Vite.** The lock on `main` carried two Vite majors (7.3.7 hoisted for `@astrojs/react`/vitest/`@tailwindcss/vite`, 8.3.6 under `astro`). With `@astrojs/react` 7.1.9 and vitest 5 the install resolves a single hoisted `vite@8.3.8`, as #56 does alone.
- **Out of scope:** removing the `ubuntu-24.04` runner pins (separate follow-up), the React #418 hydration warning that already exists on `main`, TypeScript 7 (closed #59), and any `.github/` change (see the Dependabot config PR).

## 4. Verification

On Node 26 (local) and Node 22.23.5: `npm ci`, `astro check` (0 errors, 83 hints as on `main`), `npm run build`, `npm test` (81/81, also under Node 22), `npm run build:pdf`, Playwright E2E (7/7, `CI=1`, preview on a free port), browser check of the islands against a build of `origin/main` for comparison, and the npm audit gate.

## 5. Security pass (lockfile)

See the PR body for numbers. Summary: all `resolved` URLs on `registry.npmjs.org`; no git/http/file sources; one fewer install script (esbuild is gone with Vite 8's switch to rolldown); six new packages, all from Astro/Rolldown/sindresorhus/lukeed publishers; six lockfile entries without `integrity` are the bundled `@tailwindcss/oxide-wasm32-wasi` dependencies, which have no registry tarball.

## 6. Open Questions

1. **Netlify's Node version** is not pinned anywhere in the repo (`netlify.toml` has no `NODE_VERSION`; there is no `.nvmrc`). The Netlify deploy previews of #55–#58 passed with Astro ≥ 7.4.7, which requires Node ≥ 22.23.4, which is evidence but not proof of the site setting. The production frontend deploy is built by GitHub Actions (`deploy-frontend.yml`, Node 22), not by Netlify's builder.
2. `vitest` resolves 5.1.1 here, not the 5.1.0 that #57 locks (both satisfy `^5.1.0`; 5.1.1 was newest at install).
3. vitest 5 declares `engines.node` `^22.12.0 || ^24.0.0 || >=26.0.0`, so Node 23 and 25 get an engine warning; the project targets 22.
4. Whether to also raise the `tailwindcss` range (see section 3).
