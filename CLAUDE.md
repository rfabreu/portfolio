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

**API env vars:** `PORT`, `ALLOWED_ORIGINS`, `GEMINI_API_KEY`, `RATE_LIMIT_RPM`

## API Routes

- `GET  /health`    — health check
- `POST /api/chat`  — AI chat (Google Gemini via `internal/chat/`)
- `POST /api/game`  — Tic-Tac-Toe AI move (minimax, `internal/game/`)

## Frontend Architecture

- **Pages:** `src/pages/` — `index.astro`, `playground.astro`, `project/[slug].astro`, `success.astro`
- **Islands (React):** `src/islands/` — `ParticleHero.tsx`, `ChatWidget.tsx`, `GameBoard.tsx`, `MobileNav.tsx`
- **Content collections:** `src/content/projects/*.md` — schema in `content.config.ts`
- **Path aliases:** `@components`, `@islands`, `@layouts`, `@styles` (defined in `tsconfig.json`)
- **Global styles/tokens:** `src/styles/global.css` (Tailwind v4 theme)

## Deployment

- **Frontend:** Netlify (static deploy via `nwtgck/actions-netlify@v3`). Triggered on push to `main` when `frontend/**` changes.
- **API:** Fly.io (Docker). Triggered on push to `main` when `api/**` changes.
- **CI:** Runs `astro check` + `npm run build` (frontend) and `go vet` + `go test` + `go build` (API) on every PR to `main`.

## Phase 2 Status

Roadmap: `docs/phase2-roadmap.md`

- [x] P1: Interactive Particle Hero (PR #8)
- [x] P2: AI Playground (PR #9)
- [x] V1 Cleanup (PR #10)
- [ ] P3: Kubernetes Migration (needs cloud account)
- [ ] P4: Expanded Project Case Studies
- [ ] P5: Privacy Analytics
- [ ] P6: Blog/Writing Section
- [ ] P7: Dark/Light Theme Toggle
- [ ] P8: Testimonials (blocked on content)

## Docs Workflow

Each feature follows: design spec → implementation plan → feature branch → PR
- Specs: `docs/superpowers/specs/`
- Plans: `docs/superpowers/plans/`

## Production Safety Rules

**`main` is production. Pushing to `main` triggers live deployments.**

- **NEVER push to `main` unless the user explicitly says "push to main".** Committing locally is fine when asked; pushing is not.
- **ALL work happens on feature branches** (`feature/<name>`, `fix/<name>`). Merge to `main` only via PR.
- **No "safe" exceptions.** Docs, dependency patches, config changes — everything triggers CI/CD. There is no safe push to main.
- **Verify locally before any merge.** Run `npm run build` (frontend) and `go test ./...` (API) and visually confirm in the browser.
- **`.gitignore` patterns must be root-relative** (`/src/` not `src/`) to avoid blocking Tailwind v4's source scanning of `frontend/src/`. Bare patterns match recursively and will silently break utility class generation.

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
