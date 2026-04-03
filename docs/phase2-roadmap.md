# Phase 2 Roadmap

> Phase 1 (V1) shipped March 2026: Astro frontend, Go API with AI chatbot, CI/CD pipeline.
> Phase 2 builds on this foundation with interactive features, infrastructure evolution, and content expansion.

## Priority Order

Features are ordered by impact and dependency. Each is independent and can be built as a standalone PR.

---

### P1: Interactive Particle Hero

**What:** Replace the static gradient hero background with a canvas-based particle mesh that reacts to cursor movement — interconnected dots forming a neural network pattern.

**Why:** The hero is the first thing visitors see. An interactive, technically impressive animation signals engineering craft in the first 3 seconds. This was designed during brainstorming as the "show-stopping moment."

**Scope:**
- Canvas animation rendered at ~60fps via `requestAnimationFrame`
- Particles connected by thin lines in indigo/purple palette
- Cursor proximity: nodes glow brighter and pull toward mouse
- Scroll transition: particles fade out as user scrolls past hero
- Mobile: gentle drift animation (no cursor on touch devices)
- `prefers-reduced-motion`: falls back to static gradient
- Implemented as a React island — only this component hydrates

**Estimated effort:** 1 session

---

### P2: AI Playground (`/playground`)

**What:** A dedicated page with an AI-powered strategy game where visitors play against the Go API. Primary candidate: Tic-Tac-Toe with difficulty scaling (minimax/alpha-beta pruning).

**Why:** Turns the Go backend into a playable demonstration of AI engineering, not just a chatbot proxy. Visitors literally compete against Rafael's code.

**Scope:**
- New Go endpoint: `POST /api/game` (game state management, AI move calculation)
- Frontend: React island rendering game board, move history, difficulty selector
- Astro page at `frontend/src/pages/playground.astro`
- Nav item linking to `/playground`
- Designed as a sandbox — future interactive experiments can be added to this page

**Estimated effort:** 1-2 sessions

---

### P3: Kubernetes Migration

**What:** Migrate the Go API from Fly.io (Firecracker VMs) to a Kubernetes cluster on GKE or EKS free tier.

**Why:** Rafael holds a Red Hat Kubernetes/OpenShift certificate and is actively focusing on K8s. Running the portfolio API on Kubernetes turns the infrastructure itself into a portfolio piece — demonstrating container orchestration, not just containerization. The existing Dockerfile transfers directly.

**Scope:**
- Kubernetes manifests (Deployment, Service, Ingress, ConfigMap for env vars, Secret for API keys)
- Helm chart or Kustomize overlays for environment management
- Update `deploy-api.yml` GitHub Actions workflow to deploy to K8s cluster
- Health checks and readiness probes configured
- Horizontal Pod Autoscaler (optional, demonstrates K8s knowledge)
- Document the infrastructure in README

**Estimated effort:** 2-3 sessions

**Prerequisites:** GKE/EKS free tier account with credit card on file

---

### P4: Expanded Project Case Studies

**What:** Add more projects as they become ready, with full case study detail.

**Why:** The portfolio is designed to grow. Content collections make this trivial — drop a `.md` file and it's live on next deploy.

**Scope:**
- Write case studies for new projects as they're completed
- Optionally add the remaining old projects (Jest-Another-RPG, README Generator, Portfolio Generator) with updated descriptions
- Add project cover images to `frontend/public/images/`

**Estimated effort:** Per-project, ~30 minutes each

---

### P5: Privacy-Respecting Analytics

**What:** Add Plausible or Umami analytics to understand visitor behavior without tracking individuals.

**Why:** Data on which projects get the most clicks, where visitors drop off, and how many use the chatbot informs future portfolio decisions. Privacy-respecting options align with engineering values.

**Scope:**
- Self-hosted Umami (free) or Plausible Cloud ($9/mo)
- Single script tag in `BaseLayout.astro`
- Dashboard for page views, referrers, project clicks, chatbot usage

**Estimated effort:** 1 hour

---

### P6: Blog / Writing Section

**What:** Add a `/blog` route with Markdown-based articles about engineering topics.

**Why:** Long-form writing demonstrates depth of understanding. Good for SEO and positions Rafael as a thought leader, not just a builder.

**Scope:**
- New content collection for blog posts (similar to projects)
- Blog listing page at `/blog`
- Individual post pages at `/blog/[slug]`
- RSS feed generation (Astro plugin)
- Nav item for Blog

**Estimated effort:** 1 session

---

### P7: Dark/Light Theme Toggle

**What:** Add a theme toggle allowing visitors to switch between dark and light modes.

**Why:** Accessibility and user preference. Some recruiters review portfolios in bright office environments where dark themes are harder to read.

**Scope:**
- Light theme color tokens in `global.css`
- Toggle button in Nav
- Persist preference in `localStorage`
- `prefers-color-scheme` media query for initial state
- React island for the toggle (needs client-side state)

**Estimated effort:** 1 session

---

### P8: Testimonials Section

**What:** Reactivate the testimonials section with real testimonial data.

**Why:** Social proof from colleagues, managers, or collaborators adds credibility that self-reported skills can't.

**Scope:**
- Content collection for testimonials (name, role, company, quote, image)
- Testimonials component on main page between Skills and Contact
- Carousel or grid layout

**Prerequisites:** Collect real testimonials

**Estimated effort:** 1 session (once content is ready)

---

## Dependency Map

```
P1 (Particle Hero)     — independent, no dependencies
P2 (Playground)        — independent, extends Go API
P3 (K8s Migration)     — independent, replaces Fly.io deploy
P4 (Case Studies)      — independent, content-only
P5 (Analytics)         — independent, one script tag
P6 (Blog)              — independent, new content collection
P7 (Theme Toggle)      — independent, CSS + React island
P8 (Testimonials)      — blocked on collecting testimonial content
```

All features are independent. Recommended order follows priority numbering, but any can be built at any time.

## V1 Cleanup (Optional)

Minor improvements to the current V1 that can be addressed alongside Phase 2 work:

- **Dependabot alerts** — GitHub reports 88 vulnerabilities on the old CRA dependencies at root level. These are from the legacy `node_modules/` and `package.json` that are no longer used. Consider removing old root-level CRA files (`src/`, `public/`, `package.json`, `postcss.config.js`, `tailwind.config.js`) once the new frontend is stable.
- **GitHub Pages** — the `pages-build-deployment` workflow fails because we use Netlify, not GitHub Pages. Consider disabling GitHub Pages in repo settings.
- **go.sum** — the Go API has no external dependencies yet, so no `go.sum` exists. CI cache warnings about missing `go.sum` are harmless but will resolve once an external dependency is added.
