# Portfolio — Astro + Go Monorepo

## Project Structure

```
frontend/          Astro 6 site (static output, islands architecture)
api/               Go API server (chat + game endpoints)
docs/              Roadmap, design specs, implementation plans
.github/workflows/ CI + deploy pipelines
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
Requires: Go 1.26+

**First-time clone:** run `make sync-content` (or `make dev`) before any `go` command — `internal/chat/embedded/` is gitignored and populated from `frontend/src/content/` at build time. `make dev` and `make test` invoke `sync-content` automatically.

**API env vars:** `PORT`, `ALLOWED_ORIGINS`, `GEMINI_API_KEY`, `RATE_LIMIT_RPM`

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
- **CI:** Three jobs on every PR to `main`: `Frontend` (`astro check`, build, Vitest), `Frontend E2E` (build, resume PDF, Playwright) and `API` (`make sync-content`, `go vet`, `go test`, `go build`).
- **Branch protection (checked 2026-11-07):** `main` requires the `Frontend` and `API` checks. `Frontend E2E` is not required, and enforcement applies to non-admins only.
- **Known failing checks (not required):** three Netlify preview checks for the project `rafaelswe` ("Deploy failed", cause not yet known). Tracked in Issue #34.

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
- [ ] P10: Reposition for OrangeCat Labs (Issue #30) — done: Phase 1 (PR #32, og follow-up PR #35), About copy (PR #36), Contact copy (PR #38). Open: Freckles Design credit, resume headline (via P11), optional contact fields (deferred), chatbot (spec first), analytics events, page order, GitHub repo description
- [ ] P11: Resume for Executive and Technical Leadership Audiences (Issue #37)

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

- **Production URL:** https://rafaelabreu.netlify.app (`site` in `frontend/astro.config.mjs`). `rafaelabreu.dev` is not the production URL. The references to it in `frontend/src/content/resume/profile.md` (`links.portfolio`) and `frontend/src/content/projects/portfolio.md` (`demo`) are known mismatches awaiting correction (P11 and the pending projects-content PR).
- **LinkedIn:** https://linkedin.com/in/rafael-a-gomes. `links.linkedin` in `profile.md` currently differs; tracked in Issue #37.

## Working Rules

- **Never assume.** Do not state a conclusion without evidence (code, tests, CI results, or GitHub data). If something cannot be verified, say so and ask the owner.
