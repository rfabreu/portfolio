# Portfolio — Astro + Go Monorepo

## Project Structure

```
frontend/          Astro 6 site (static output, islands architecture)
api/               Go API server (chat + game endpoints)
docs/              Roadmap, design specs, implementation plans
.github/           CI + deploy workflows, dependabot.yml, scripts/ (npm audit gate)
```

Legacy CRA dirs at root (`src/`, `public/`, `build/`, `node_modules/`) are gitignored — not tracked.

## Build & Run

### Frontend (Astro)
```bash
cd frontend
npm ci
npm run dev          # dev server on :4321
npm run build        # static build → frontend/dist/
npx astro check      # type checking
```
Requires: Node >= 22.12.0

Frontend env var: `PUBLIC_API_URL` (injected at build, see `frontend/.env.example`)

### API (Go)
```bash
cd api
go run ./cmd/server  # dev server on :8080
go test ./... -v     # run tests
go vet ./...         # lint
go build ./cmd/server
```
Requires: Go 1.26.1+

**First-time clone:** run `make sync-content` (or `make dev`) before any `go` command — `internal/chat/embedded/` is gitignored and populated from `frontend/src/content/` at build time. `make dev` and `make test` invoke `sync-content` automatically.

**API env vars:** `PORT`, `ALLOWED_ORIGINS`, `GEMINI_API_KEY`, `RATE_LIMIT_RPM`

**API limits** (PR #50): `/api/chat` body 1 MiB (413), roles `user`/`assistant` only, at most 20 messages, `user` content 2000 runes (must match `MAX_INPUT_LENGTH` in `ChatWidget.tsx`, guarded by `frontend/tests/chat-input-limit.test.ts`), `assistant` 8000 runes (400). `/api/game` body 4 KiB (413). Gemini calls time out after 30s (`chat.RequestTimeout`); the key is sent in the `x-goog-api-key` header, never the URL, and the client refuses redirects. Server timeouts are set in `newServer` in `api/cmd/server/main.go`.

## API Routes

- `GET  /health`    — health check
- `POST /api/chat`  — AI chat (Google Gemini via `internal/chat/`)
- `POST /api/game`  — Tic-Tac-Toe AI move (minimax, `internal/game/`)

## Frontend Architecture

- **Pages:** `src/pages/` — `index.astro`, `playground.astro`, `project/[slug].astro`, `resume.astro`, `success.astro`
- **Islands (React):** `src/islands/` — `ParticleHero.tsx`, `ChatWidget.tsx`, `GameBoard.tsx`, `MobileNav.tsx`, `ThemeToggle.tsx`
- **Content collections:** `src/content/projects/*.md`, `src/content/resume/`, `src/content/testimonials/` — schemas in `content.config.ts`
- **Path aliases:** `@components`, `@islands`, `@layouts`, `@styles` (defined in `tsconfig.json`)
- **Global styles/tokens:** `src/styles/global.css` (Tailwind v4 theme)

## Deployment

- **Frontend:** Netlify (static deploy via `nwtgck/actions-netlify@v3`). Triggered on push to `main` when `frontend/**` changes.
- **API:** Fly.io (Docker). Triggered on push to `main` when `api/**`, `frontend/src/content/resume/**`, or `frontend/src/content/projects/**` changes (the chatbot embeds that content). A resume or project content change therefore deploys both the frontend and the API.
- **CI:** Three jobs on every PR to `main`: `Frontend` (`astro check`, build, Vitest, then the npm audit gate), `Frontend E2E` (build, resume PDF, Playwright) and `API` (`make sync-content`, `go vet`, `go test`, `go build`, then pinned `govulncheck`). The audit gate runs its own unit tests, then `npm audit --json` checked by `.github/scripts/npm-audit-gate.mjs`: high and critical advisories fail unless allowlisted there with a reason and an expiry date. Currently allowlisted: GHSA-fc5q-39xc-vjpf (picomatch ReDoS; no patched release, 4.0.7 was still the latest on npm, checked 2026-11-18), until 2027-01-31. An expired entry fails the gate. `govulncheck` is pinned to `v1.9.0` in `ci.yml`; bump it by hand.
- **Playwright runner pin:** the `Frontend E2E` job and the `Deploy to Netlify` job (`deploy-frontend.yml`) run on `ubuntu-24.04`, not `ubuntu-latest` (Ubuntu 26.04), because the locked Playwright 1.59.1 cannot install Chromium there (PR #48). Remove both pins after upgrading Playwright to 1.61 or later.
- **Dependency updates:** `.github/dependabot.yml` opens weekly PRs for `frontend/` (npm), `api/` (Go modules) and GitHub Actions, with minor and patch updates grouped per ecosystem. Nothing auto-merges; each PR needs CI and owner approval.
- **Analytics:** Umami loads in production builds only, with `data-domains="rafaelabreu.netlify.app"`, so previews and localhost report nothing. Events: `chat-opened`, `game-started`, `orangecat-click`, `linkedin-click`, `view-projects-click`, `contact-form-submit`, `contact-form-success`. New OrangeCat Labs or Rafael LinkedIn links must carry `data-umami-event` (enforced by `frontend/tests/analytics-funnel-events.test.ts`). Catalogue: `docs/superpowers/specs/2026-11-17-analytics-funnel-events-design.md`.
- **Branch protection (checked 2026-11-17):** `main` requires the `Frontend`, `API` and `Frontend E2E` checks, with strict mode (branch must be up to date) on. Admin enforcement is off, so the rules bind non-admins only. Force pushes and deletions are disallowed. Required reviews: the protection endpoint returns no review rule, but the dedicated `required_pull_request_reviews` endpoint returns 1 approving review, so the owner should confirm which applies.
- **Netlify preview checks:** Issue #34 (failing `rafaelswe` previews) is closed (2026-11-07). PR #44 added a root `netlify.toml` (`base = "frontend"`); on PR #42 `netlify/rafaelswe/deploy-preview` passes (checked 2026-11-17).

## Phase 2 Status

Roadmap: `docs/phase2-roadmap.md`

- [x] P1: Interactive Particle Hero (PR #8)
- [x] P2: AI Playground (PR #9)
- [x] V1 Cleanup (PR #10)
- [ ] P3: Kubernetes Migration (on hold — needs $0 hosting option)
- [x] P4: Expanded Project Case Studies (PR #20, #22)
- [x] P5: Privacy Analytics (PR #19)
- [ ] P6: Blog/Writing Section (needs article content)
- [x] P7: Dark/Light Theme Toggle (PR #17)
- [x] P8: Testimonials architecture (PR #31) — section stays hidden until real testimonials are added
- [x] P9: Resume Feature (PR #26, #27)
- [ ] P10: Reposition for OrangeCat Labs (Issue #30) — done: Phase 1 (PR #32, og follow-up PR #35), About copy (PR #36), Contact copy (PR #38), Freckles Design credit (PR #41), analytics funnel events (PR #49), Services section, Project Brief contact fields (`project_type`, `timeline`; no budget field), project-page CTA and resume `title` (PR #42), GitHub repo description (now names OrangeCat Labs; checked 2026-11-17). Open: resume summary and the rest of the resume (via P11), chatbot (spec first), chat widget greeting, page order. **Owner to confirm:** Issue #30 says the standalone Services section was dropped from the plan, but PR #42 shipped `frontend/src/components/Services.astro` and a Services nav item; was reinstating it intended?
- [ ] P11: Resume for Executive and Technical Leadership Audiences (Issue #37) — open; PR #42 already changed the `profile.md` `title` to "Co-Founder, OrangeCat Labs | Full-Stack & Infrastructure Leader", the `summary` is unchanged
- [x] Dependabot triage (Issue #39): closed 2026-11-13; 0 open alerts (checked 2026-11-17)

## Known Follow-ups

- Upgrade Playwright to 1.61 or later, confirm Chromium installs on Ubuntu 26.04, then remove the `ubuntu-24.04` pins in `ci.yml` and `deploy-frontend.yml`.
- Review the GHSA-fc5q-39xc-vjpf allowlist entry before 2027-01-31: drop it if a patched picomatch exists, otherwise renew it with a fresh assessment.

## Docs Workflow

Every change follows: design spec → implementation plan → feature branch → PR. This includes copy-only and docs-only changes (owner decision, 2026-11-07).
- Specs: `docs/superpowers/specs/`
- Plans: `docs/superpowers/plans/`

## Production Safety Rules

**`main` is production. Pushing to `main` triggers live deployments.**

- **NEVER push to `main` unless the user explicitly says "push to main".** Committing locally is fine when asked; pushing is not.
- **ALL work happens on feature branches** (`feature/<name>`, `fix/<name>`). Merge to `main` only via PR.
- **No "safe" exceptions.** Docs, dependency patches, config changes — everything triggers CI/CD. There is no safe push to main.
- **Verify locally before any merge.** Run `npm run build` (frontend) and `go test ./...` (API) and visually confirm in the browser.
- **Adversarial review before merge.** Every PR gets an adversarial review for repo conventions, the positioning constraints below, and bugs. Merge only after that review and all CI checks pass, and only with the user's approval.
- **`.gitignore` pattern rules:** Patterns that could shadow real source directories (`src/`, `build/`, `public/`, `node_modules/`) **must be root-relative** (`/src/` not `src/`) to avoid blocking Tailwind v4's scanning of `frontend/src/`. Tool/IDE directories that legitimately appear at any depth (`.vscode/`, `.astro/`, `.netlify`, `.claude/`) **must stay bare** so they match recursively (e.g. `frontend/.vscode/`).

## Branching

- `main` is production — protected, deploy-on-push
- Feature branches: `feature/<name>`
- Bug fixes: `fix/<name>`
- Each branch → PR to `main`

## Key Conventions

- **Astro components** are `.astro` files; interactive islands use React (`.tsx`) with `client:load`
- **Islands pattern:** only interactive components hydrate (`client:load`); everything else is static Astro
- **Styling:** Tailwind CSS v4 with custom theme tokens in `frontend/src/styles/global.css`
- **Go API:** stdlib `net/http` with `http.NewServeMux()` — no frameworks
- **API architecture:** `internal/<domain>/handler.go` + `engine.go` or `provider.go`
- **Contact form:** Netlify Forms (`data-netlify="true"`) — not handled by Go API
- **No email exposure:** Contact is form-only, never display email addresses publicly
- **Accessibility:** respect `prefers-reduced-motion`; use semantic HTML and ARIA labels

## Positioning Constraints (P10)

The site drives business to OrangeCat Labs (https://orangecatlabs.ca). While P10 is in progress:

- No budget or pricing language anywhere: site copy, form fields, or chatbot.
- No Upwork link on the portfolio.
- Only Freckles Design is attributed to OrangeCat Labs; no other project is linked to the studio. The owner confirmed this attribution is accurate: Rafael built it under the company.
- Contact paths are the Netlify form and LinkedIn only.

## Canonical Facts

Confirmed by the owner on 2026-11-07:

- **Production URL:** https://rafaelabreu.netlify.app (`site` in `frontend/astro.config.mjs`). `rafaelabreu.dev` is not the production URL. `links.portfolio` in `frontend/src/content/resume/profile.md` and `demo` in `frontend/src/content/projects/portfolio.md` now use it (corrected in PR #41, checked on `main` 2026-11-17).
- **LinkedIn:** https://linkedin.com/in/rafael-a-gomes. `links.linkedin` in `profile.md` now matches (PR #41, checked on `main` 2026-11-17).

## Working Rules

- **Never assume.** Do not state a conclusion without evidence (code, tests, CI results, or GitHub data). If something cannot be verified, say so and ask the owner.
