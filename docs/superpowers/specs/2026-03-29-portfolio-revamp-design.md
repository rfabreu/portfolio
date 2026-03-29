# Portfolio Revamp — Design Specification

## Overview

Complete revamp of Rafael Abreu's developer portfolio from a React/CRA single-page app to a modern Astro + Go API architecture. The portfolio positions Rafael as a versatile, AI-forward software engineer who builds across the full stack. It must make a strong impression in 30 seconds for recruiters while holding up under deep technical scrutiny from engineering leads.

## Goals & Positioning

- **Primary positioning:** Versatile full-stack engineer who understands and builds with AI
- **Target audience:** Recruiters (30-second scan), technical leads (deep evaluation), collaborators, freelance clients
- **Differentiators:** The portfolio's own tech stack demonstrates engineering range (Astro + Go), an AI chatbot serves as a live demo of AI integration skills, and the design is crafted — not templated

## Architecture

### Monorepo Structure

```
portfolio/
├── frontend/                 ← Astro project
│   ├── src/
│   │   ├── layouts/          (base page layouts)
│   │   ├── pages/            (routes: /, /project/:slug, /playground)
│   │   ├── components/       (static Astro components)
│   │   ├── islands/          (React interactive components)
│   │   ├── content/          (Markdown/MDX project files)
│   │   ├── styles/           (Tailwind v4 + global CSS)
│   │   └── assets/           (images, fonts)
│   ├── public/               (static files, favicon, etc.)
│   ├── astro.config.mjs
│   └── package.json
│
├── api/                      ← Go service
│   ├── cmd/server/           (entry point)
│   ├── internal/
│   │   ├── chat/             (AI chatbot handler)
│   │   ├── game/             (AI game logic — Phase 2)
│   │   └── middleware/       (CORS, rate limiting)
│   ├── Dockerfile
│   └── go.mod
│
├── .github/workflows/        ← CI/CD
│   ├── ci.yml                (lint, test, build on PRs)
│   ├── deploy-frontend.yml   (Astro → Netlify on push to main)
│   └── deploy-api.yml        (Go → Fly.io on push to main)
│
└── README.md
```

### Key Architectural Decisions

- **Astro with React islands:** Static HTML by default. Only the hero animation, AI chatbot, and game hydrate as interactive React components. Zero JS shipped for static sections.
- **Go API:** Lightweight service with two responsibilities: AI chat endpoint (proxies to Gemini with resume context) and game state management (Phase 2). Deployed as a Docker container.
- **Content Collections:** Projects defined as Markdown/MDX files with typed frontmatter. Adding a project = dropping a `.md` file.
- **LLM Provider:** Google Gemini API (free tier — Gemini 2.0 Flash, 15 RPM, 1M tokens/day). Go API uses an interface-based LLM provider so it's swappable to another provider via config change.
- **Monorepo:** Single repository, two deployable units. Path-filtered CI/CD so changes to `frontend/` don't redeploy the Go service and vice versa.

## Branching Strategy

**Simplified flow replacing the current feature → develop → main model:**

```
feature/*  ──PR──►  main  ──auto-deploy──►  production
                │
                ├── CI runs (lint, test, build)
                ├── Preview deploy URL generated
                └── Merge when green
```

- No `develop` branch — PR-based preview deploys replace the staging role
- Branch protection on `main`: require CI to pass before merge
- Cleaner git history with fewer merge commits

## Visual Design

### Design Direction: "Noir Engineer + Electric Accents + Typographic Confidence"

A blend of three directions:
- **A's dark precision:** Deep black base, monospace touches for code/labels
- **C's electric accents:** Indigo (#6366f1) / Purple (#8b5cf6) / Pink (#ec4899) accent palette with subtle gradient orbs for depth
- **D's typographic boldness:** Large, bold, tight-tracked headlines that command attention

### Design System

| Token | Value |
|---|---|
| Background (base) | `#050508` |
| Surface | `rgba(255,255,255,0.03)` with `rgba(255,255,255,0.06)` border |
| Text primary | `#ffffff` |
| Text secondary | `#a0a0b0` |
| Text muted | `#666666` |
| Accent indigo | `#6366f1` |
| Accent purple | `#8b5cf6` |
| Accent pink | `#ec4899` |
| Accent indigo light | `#a5b4fc` |
| Accent purple light | `#c4b5fd` |
| Accent pink light | `#f9a8d4` |
| Font headline | System sans-serif, weight 800-900, letter-spacing -2px to -3px |
| Font body | System sans-serif, weight 400, line-height 1.6 |
| Font code/labels | `'Courier New', monospace`, used for tech tags, section numbers, and code references |
| Section numbering | `01 // ABOUT`, `02 // PROJECTS`, etc. — monospace, indigo color |
| Dividers | Gradient lines: `linear-gradient(90deg, transparent, rgba(99,102,241,0.3), transparent)` |
| Border radius | Cards: 12px, Tags: 5px, Buttons: 8px |
| Gradient orbs | Radial gradients of accent colors at low opacity (0.06-0.15) for atmospheric depth |

### Responsive Breakpoints

- Mobile: < 768px — single column, stacked layout, hamburger menu
- Tablet: 768px-1024px — 2-column grids, condensed nav
- Desktop: > 1024px — full layout as designed

## Site Structure

### Main Page (`/`) — Single Scroll

```
Nav (sticky, scroll-aware active states)
├── Hero          — animated first impression
├── About         — brief bio, positioning statement
├── Projects      — curated grid (featured + other)
├── Skills        — categorized tech stack
├── Contact       — streamlined form + socials
└── Footer
```

### Breakout Routes

- **`/project/:slug`** — Individual project case study pages, generated from Markdown content collections
- **`/playground`** — AI game + interactive experiments (Phase 2)

### Navigation

- Sticky nav with scroll-aware active states on main page
- Smooth scroll for section links on main page
- Standard routing for breakout pages
- Mobile: hamburger with full-screen overlay
- Nav items: About, Projects, Skills, Playground, Contact
- Nav actions: Resume (outline button), Let's Talk (filled button)
- Contact form: Netlify Forms (free, serverless, works independently of Go API)

## Hero Section

### Concept

Particle/node mesh background that reacts to cursor movement — interconnected dots forming a loose neural network pattern, signaling "builds intelligent systems."

### MVP (V1)

Static gradient background with the approved typography layout:
- Monospace label: `// HELLO WORLD`
- Bold name: `RAFAEL ABREU` (900 weight, -3px tracking)
- Tagline: "Software Engineer building intelligent systems across the full stack"
- Tech tags: Go, Python, React, Kubernetes, AI/ML (color-coded indigo/purple/pink)
- CTAs: "View Projects" (filled indigo) + "Ask my AI →" (outline)
- Text animates in on load: name slides up, tagline fades in with delay, tags stagger in

### Phase 2

Canvas-based particle animation:
- ~60fps, `requestAnimationFrame` with visibility check
- Particles connected by thin lines in indigo/purple palette
- Cursor proximity: nodes glow brighter and pull slightly toward mouse
- Scroll transition: particles fade out as user scrolls past hero
- Mobile: gentle drift animation instead of cursor interaction (no hover on touch)
- `prefers-reduced-motion`: static gradient fallback
- Implemented as a React island — only this component hydrates

## Projects Section

### Two-Tier Display

**Tier 1 — Featured Projects (2-3 max):**
- Large cards (full-width or half-width)
- Project screenshot/visual
- Title + one-line description
- Tech stack tags (monospace, color-coded)
- Category badge (AI, Full-Stack, DevOps, etc.)
- Links to `/project/:slug` case study page

**Tier 2 — Other Projects:**
- Smaller grid cards (2-3 per row)
- Title + brief description
- Tech tags
- GitHub + live demo external links
- No case study page required

### Case Study Pages (`/project/:slug`)

Generated from Markdown/MDX. Structure:

```
Overview        — what is it, one paragraph
Problem         — what problem does it solve
Approach        — how you built it, key decisions
Architecture    — diagram or description of the system
Tech Stack      — detailed breakdown with reasoning
Outcome         — results, learnings, what you'd do differently
Links           — GitHub, live demo, related writing
```

### Project Frontmatter Schema

```yaml
---
title: "AI Chat Service"
description: "Go-powered conversational AI backend"
image: "./ai-chat-cover.png"
tags: ["Go", "gRPC", "Claude API"]
category: "AI"
featured: true
order: 1
---
```

### Scroll Animations

Cards reveal on scroll with subtle upward fade, staggered timing. CSS transitions + Intersection Observer — no heavy animation library needed.

## AI Features

### 1. AI Chatbot — "Ask Rafael's AI" (MVP)

**UX:**
- Floating button (bottom-right): gradient indigo/purple, 52px pill
- Click → chat panel slides up
- Greeting: "Hey! I'm Rafael's AI assistant. Ask me about his experience, projects, or tech stack."
- Dark-themed panel matching site aesthetic
- Streaming responses for real-time typing feel

**Backend (Go API):**
- `POST /api/chat` — receives user message + conversation history
- Server injects resume, project descriptions, and bio as system context
- Calls Google Gemini API (Gemini 2.0 Flash, free tier)
- Streams response back via Server-Sent Events
- Conversation history maintained in session (no database)

**Guardrails:**
- Rate limiting in Go middleware (prevent abuse)
- System prompt constrains responses to professional background only
- Conversation capped at ~20 messages per session
- Fallback: if API is down, shows "Leave a message" routing to contact form
- LLM provider interface — swappable via config (Gemini now, Claude/OpenAI later if needed)

### 2. AI Playground — `/playground` (Phase 2)

AI-powered strategy game where visitors play against a Go-built AI opponent. Options:
- AI Tic-Tac-Toe with difficulty scaling (minimax/alpha-beta pruning)
- AI Snake with competing AI agent (pathfinding demonstration)

Game logic runs in Go API, frontend renders game state. Visitors literally play against Rafael's Go code.

Page designed as a sandbox that grows — future interactive experiments can be added here.

### 3. AI-Enhanced Project Descriptions (Phase 2)

Dynamic project summaries that adapt based on visitor context. Stretch goal.

## CI/CD Pipeline

### Workflow 1: `ci.yml` — Pull Request Validation

```
Trigger: pull_request → main

Jobs (parallel):
├── frontend-check:
│   ├── npm ci (in frontend/)
│   ├── astro check (type checking)
│   ├── eslint
│   └── astro build (verify it compiles)
│
├── api-check:
│   ├── go vet
│   ├── golangci-lint
│   ├── go test ./...
│   └── go build ./cmd/server
│
└── preview-deploy:
    └── Deploy Astro to Netlify preview URL
        (auto-comment on PR with preview link)
```

### Workflow 2: `deploy-frontend.yml` — Production Frontend

```
Trigger: push to main (paths: frontend/**)

Steps:
├── npm ci && astro build
├── Deploy to Netlify (production)
└── Purge CDN cache
```

### Workflow 3: `deploy-api.yml` — Production Go API

```
Trigger: push to main (paths: api/**)

Steps:
├── go test ./...
├── Build Docker image
├── Push to GitHub Container Registry
├── Deploy to Fly.io
└── Health check (curl /health endpoint)
```

### Deployment Targets

| Component | Platform | Tier | Reason |
|---|---|---|---|
| Frontend | Netlify | Free | Great CDN, auto HTTPS, preview deploys on PRs |
| Go API | Fly.io | Free | Docker-based, auto-sleep when idle (~2s wake), scales to zero cost |
| LLM | Google Gemini API | Free | 15 RPM, 1M tokens/day — sufficient for portfolio traffic |
| Container Registry | GitHub Container Registry | Free | Integrated with GitHub Actions |

### Secrets (GitHub Actions)

- `NETLIFY_AUTH_TOKEN` — Netlify deployment
- `NETLIFY_SITE_ID` — Netlify site identifier
- `FLY_API_TOKEN` — Fly.io deployment
- `GEMINI_API_KEY` — Google Gemini API for chatbot

### Branch Protection on `main`

- Require CI workflow to pass before merge
- No direct pushes to `main`

## MVP vs Phase 2

### MVP (V1) — Ship Fast

| Feature | Details |
|---|---|
| Astro frontend | All sections: nav, hero (static gradient + text animation), about, projects grid, skills, contact, footer |
| Visual design | Full design system applied — dark base, indigo/purple/pink accents, bold typography, gradient orbs |
| Project markdown | Content collections, 2-3 projects with case study pages |
| Responsive design | Full mobile/tablet/desktop |
| Scroll animations | Fade-in reveals via Intersection Observer + CSS |
| Go API | `/health` endpoint + chatbot endpoint with Gemini Flash |
| AI Chatbot | Floating widget, streaming responses, resume context, rate limiting |
| CI/CD | All 3 GitHub Actions workflows, Netlify + Fly.io deployment |
| Branch protection | CI gate on PRs to `main` |

### Phase 2 — Layer In After V1

| Feature | Details |
|---|---|
| Particle hero | Canvas-based interactive particle mesh with cursor reactivity |
| Playground page | AI-powered game at `/playground` |
| Expanded case studies | Full detail for remaining projects |
| Testimonials | Reactivate section with real data |
| AI-enhanced descriptions | Dynamic project highlights based on visitor context |
| Analytics | Privacy-respecting (Plausible or Umami) |
| Dark/light toggle | Currently dark only |
| Blog/writing section | Engineering topics, if desired |
