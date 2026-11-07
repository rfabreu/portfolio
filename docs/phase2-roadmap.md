# Phase 2 Roadmap

> Phase 1 (V1) shipped March 2026: Astro frontend, Go API with AI chatbot, CI/CD pipeline.
> Phase 2 builds on this foundation with interactive features, infrastructure evolution, and content expansion.

## Priority Order

Features are ordered by impact and dependency. P1–P9 are independent and can each be built as a standalone PR; P10 and P11 are sequenced (see the Dependency Map).

---

### ~~P1: Interactive Particle Hero~~ ✅ Done (PR #8)

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

### ~~P2: AI Playground (`/playground`)~~ ✅ Done (PR #9)

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

### ~~P4: Expanded Project Case Studies~~ ✅ Done (PR #20, #22)

**What:** Add more projects as they become ready, with full case study detail.

**Why:** The portfolio is designed to grow. Content collections make this trivial — drop a `.md` file and it's live on next deploy.

**Scope:**
- Write case studies for new projects as they're completed
- Optionally add the remaining old projects (Jest-Another-RPG, README Generator, Portfolio Generator) with updated descriptions
- Add project cover images to `frontend/public/images/`

**Estimated effort:** Per-project, ~30 minutes each

---

### ~~P5: Privacy-Respecting Analytics~~ ✅ Done (PR #19)

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

### ~~P7: Dark/Light Theme Toggle~~ ✅ Done (PR #17)

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

### ~~P8: Testimonials Section~~ ✅ Architecture done (PR #31)

**What:** Reactivate the testimonials section with real testimonial data.

**Why:** Social proof from colleagues, managers, or collaborators adds credibility that self-reported skills can't.

**Scope:**
- Content collection for testimonials (name, role, company, quote, image)
- Testimonials component on main page between Skills and Contact
- Carousel or grid layout

**Status:** The content collection, card and section components, and unit tests shipped in PR #31 (Issue #16 closed). The section renders nothing until at least one testimonial exists in `frontend/src/content/testimonials/`.

**Remaining:** Collect real testimonials and add them as content files.

---

### ~~P9: Resume Feature~~ ✅ Done (PR #26, #27)

**What:** A markdown-driven resume rendered as an in-site modal, downloadable PDF, and chatbot grounding source — replacing the previously broken Resume button.

**Why:** The Resume button on the portfolio was a 404. Rather than dropping a static PDF that goes stale, this establishes a single-source-of-truth model: edit one markdown file, three artifacts (modal, PDF, chatbot prompt) regenerate deterministically. Also fixes a latent staleness bug where the chatbot's project knowledge drifted from the actual portfolio.

**Scope:**
- New Astro content collection at `frontend/src/content/resume/profile.md` with full Zod schema
- 9 section components rendering structured data into HTML
- ResumeModal using the native `<dialog>` element (no React island)
- /resume standalone page (also serves as Playwright's PDF print source)
- Build-time PDF generation via Playwright in CI
- Chatbot system prompt assembled from resume + project markdown via `go:embed`
- E2E coverage for the modal flow

**Estimated effort:** 1-2 sessions

**Out of scope (deferred):** AI-driven LinkedIn ingestion, GitHub-based skill auto-extraction, AI-generated bullets. Considered during brainstorming and rejected as overkill for an artifact that changes 4–6 times per year.

---

### P10: Reposition for OrangeCat Labs 🚧 In progress (Issue #30)

**What:** Shift the portfolio from an individual employment pitch to presenting Rafael as a successful professional and Co-Founder of OrangeCat Labs, with the site driving qualified leads to the studio. Experience stays as credibility.

**Why:** The site is the most visible asset for the studio. Its first screen, calls to action, and contact path should send business to https://orangecatlabs.ca instead of recruiting an employee.

**Status** (checked against `main`, Issue #30, and owner decisions on 2026-11-07):
- Phase 1: positioning, CTAs, nav, footer credit, SEO and social tags — ✅ Done (PR #32); absolute `og:image` and `og:url` — ✅ Done (PR #35)
- Phase 2:
  - About copy in founder voice with one OrangeCat Labs link — ✅ Done (PR #36)
  - Contact section copy for business inquiries — ✅ Done (PR #38, copy only)
  - Services section — dropped: the studio's own site covers services, and the hero, nav, and footer already link to it
  - Contact form fields (company, project type; never a budget field) — deferred, undecided
  - Freckles Design case study credits OrangeCat Labs — planned (the owner confirmed the credit is accurate)
  - Resume headline alignment — planned, delivered through P11 (Issue #37)
  - `projects/portfolio.md` `demo` link corrected to the production URL — planned, no issue or PR yet; it can share a PR with the Freckles credit because both change projects content
- Phase 3: chatbot grounding (design spec first), chat widget greeting, analytics events (only `chat-opened` and `game-started` exist today), page section order — planned
- Follow-ups: GitHub repo description (repo settings, owner)

**Constraints:**
- No budget or pricing language anywhere (copy, form fields, chatbot)
- No Upwork link on the portfolio
- Only Freckles Design is attributed to OrangeCat Labs
- Contact paths remain the Netlify form and LinkedIn; no public email

See Issue #30 for the full checklist.

---

### P11: Resume for Executive and Technical Leadership Audiences 📋 Planned (Issue #37)

**What:** Rework the resume content (`frontend/src/content/resume/profile.md`) so it reads as credible to C-level, VP, and senior technical audiences, with business framing (scope, outcomes, risk, scale) alongside the technical depth.

**Why:** The resume is the secondary link behind the OrangeCat Labs call to action. Its headline and summary should match the co-founder voice the rest of the site now has.

**Scope** (see Issue #37 for the full checklist):
- Decide the target audiences and the headline title; rewrite the summary in executive voice
- Audit every bullet for a factual result or scale figure; metrics must come from Rafael, none invented
- Decide how OrangeCat Labs appears in the work history
- Reconcile the Red Hat certification wording with the About section
- Correct `links.linkedin` to `https://linkedin.com/in/rafael-a-gomes` and `links.portfolio` to `https://rafaelabreu.netlify.app` (both owner-confirmed 2026-11-07; Issue #37 still lists the portfolio URL as a question)
- Regenerate the PDF and check the modal and `/resume` page

**Constraints:** no email address on the resume; no budget or pricing language; how the resume describes the co-founder role is a decision for Issue #37.

**Impact:** `profile.md` feeds the resume modal, the `/resume` page, the PDF, and the chatbot's grounding. Merging a change under `frontend/src/content/resume/**` deploys both the frontend and the API.

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
P8 (Testimonials)      — architecture done; section hidden until testimonials are collected
P9 (Resume Feature)    — independent, replaces broken Resume button
P10 (OrangeCat Labs)   — Phase 1, About copy, and Contact copy done; the remaining Phase 2 and Phase 3 items build on it (chatbot work needs a spec first)
P11 (Resume)           — builds on P10; delivers P10's resume headline item
```

P1–P9 are independent. P10 and P11 are sequenced as described above. Recommended order follows priority numbering.

## V1 Cleanup (Optional)

Minor items from the V1 review that can be addressed alongside Phase 2 work (some are now resolved):

- **Dependabot alerts** — PR #10 removed the legacy CRA files that the original 88 alerts referred to. The alerts open now are on the current frontend dependencies and are triaged in Issue #39.
- **GitHub Pages** — disabled in repo settings (`has_pages` is false, checked 2026-11-07); nothing further to do.
- **go.sum** — the Go API has no external dependencies yet, so no `go.sum` exists. CI cache warnings about missing `go.sum` are harmless but will resolve once an external dependency is added.
