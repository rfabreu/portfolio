# Portfolio — Astro + Go Monorepo

## Project Structure

```
frontend/          Astro 6 site (static output)
api/               Go API server (AI chatbot)
docs/              Roadmap and design specs
.github/workflows/ CI + deploy pipelines
```

Legacy CRA files at root (`src/`, `public/`, `build/`, `package.json`, `node_modules/`) are unused — the active frontend is in `frontend/`.

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

## Deployment

- **Frontend:** Netlify (static deploy via `nwtgck/actions-netlify@v3`). Triggered on push to `main` when `frontend/**` changes.
- **API:** Fly.io (Docker). Triggered on push to `main` when `api/**` changes.
- **CI:** Runs `astro check` + `npm run build` (frontend) and `go vet` + `go test` + `go build` (API) on every PR to `main`.

## Branching

- `main` is production
- Feature branches: `feature/<name>`
- Bug fixes: `fix/<name>`
- Each branch → PR to `main`

## Key Conventions

- **Astro components** are `.astro` files; interactive islands use React (`.tsx`) with `client:load`
- **Styling:** Tailwind CSS v4 with custom theme tokens in `frontend/src/styles/global.css`
- **Go API:** stdlib `net/http` with `http.NewServeMux()` — no frameworks
- **Contact form:** Netlify Forms (`data-netlify="true"`) — not handled by Go API
- **No email exposure:** Contact is form-only, never display email addresses publicly
