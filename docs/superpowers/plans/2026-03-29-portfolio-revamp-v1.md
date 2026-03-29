# Portfolio Revamp V1 — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Revamp the portfolio from React/CRA to a modern Astro + Go API monorepo with AI chatbot, deployed via GitHub Actions to Netlify (frontend) and Fly.io (API).

**Architecture:** Astro static-first frontend with React islands for interactive components (hero animations, AI chatbot). Go microservice for the AI chatbot API using Google Gemini free tier. Monorepo with path-filtered CI/CD. All secrets managed via environment variables — never committed.

**Tech Stack:** Astro 5, React 19, Tailwind CSS v4, TypeScript, Go 1.22+, Google Gemini API, GitHub Actions, Netlify, Fly.io, Docker

---

## File Map

### Frontend (`frontend/`)

| File | Responsibility |
|---|---|
| `frontend/package.json` | Dependencies and scripts |
| `frontend/astro.config.mjs` | Astro configuration with React integration |
| `frontend/tsconfig.json` | TypeScript configuration |
| `frontend/tailwind.config.mjs` | Tailwind v4 theme with design tokens |
| `frontend/.env.example` | Template for frontend env vars (API URL) |
| `frontend/src/styles/global.css` | Tailwind directives + base styles + design tokens as CSS custom properties |
| `frontend/src/layouts/BaseLayout.astro` | HTML shell, meta tags, fonts, global styles |
| `frontend/src/layouts/ProjectLayout.astro` | Layout for project case study pages |
| `frontend/src/pages/index.astro` | Main page composing all sections |
| `frontend/src/pages/project/[slug].astro` | Dynamic project case study route |
| `frontend/src/components/Nav.astro` | Sticky navigation with scroll-aware states |
| `frontend/src/components/Hero.astro` | Hero section with gradient background and animated text |
| `frontend/src/components/About.astro` | About/bio section |
| `frontend/src/components/Projects.astro` | Two-tier project grid |
| `frontend/src/components/FeaturedProjectCard.astro` | Large card for featured projects |
| `frontend/src/components/ProjectCard.astro` | Smaller card for other projects |
| `frontend/src/components/Skills.astro` | Categorized tech stack section |
| `frontend/src/components/Contact.astro` | Contact form (Netlify Forms) + socials |
| `frontend/src/components/Footer.astro` | Footer with social links |
| `frontend/src/components/SectionHeader.astro` | Reusable numbered section header (`01 // ABOUT`) |
| `frontend/src/components/GradientDivider.astro` | Reusable gradient line divider |
| `frontend/src/components/TechTag.astro` | Reusable monospace tech tag |
| `frontend/src/islands/ChatWidget.tsx` | React island: AI chatbot floating widget |
| `frontend/src/islands/MobileNav.tsx` | React island: mobile hamburger menu with state |
| `frontend/src/content/config.ts` | Astro content collection schema for projects |
| `frontend/src/content/projects/fptv.md` | Project: FPTV |
| `frontend/src/content/projects/freckles-design.md` | Project: Freckles Design |
| `frontend/src/content/projects/charge-it-up.md` | Project: Charge It Up |
| `frontend/public/favicon.ico` | Favicon (migrate from current) |
| `frontend/public/images/` | Project images (migrate from current) |

### API (`api/`)

| File | Responsibility |
|---|---|
| `api/go.mod` | Go module definition |
| `api/go.sum` | Dependency checksums |
| `api/cmd/server/main.go` | Entry point: HTTP server setup, route registration |
| `api/internal/chat/handler.go` | HTTP handler for `POST /api/chat` |
| `api/internal/chat/handler_test.go` | Tests for chat handler |
| `api/internal/chat/provider.go` | LLM provider interface + Gemini implementation |
| `api/internal/chat/provider_test.go` | Tests for provider (with mock) |
| `api/internal/chat/context.go` | Resume/bio context for system prompt |
| `api/internal/middleware/cors.go` | CORS middleware |
| `api/internal/middleware/ratelimit.go` | Rate limiting middleware |
| `api/internal/middleware/ratelimit_test.go` | Tests for rate limiter |
| `api/Dockerfile` | Multi-stage Docker build |
| `api/.env.example` | Template for API env vars (GEMINI_API_KEY, etc.) |

### CI/CD (`.github/workflows/`)

| File | Responsibility |
|---|---|
| `.github/workflows/ci.yml` | PR validation: lint, test, build both frontend and API |
| `.github/workflows/deploy-frontend.yml` | Production frontend deploy to Netlify on push to main |
| `.github/workflows/deploy-api.yml` | Production API deploy to Fly.io on push to main |

### Root

| File | Responsibility |
|---|---|
| `.gitignore` | Updated for monorepo (Go, Astro, env files, .superpowers/) |
| `.env.example` | Root-level reference to both env files |
| `README.md` | Updated project documentation |

---

## Task 1: Project Scaffolding & Monorepo Setup

**Files:**
- Create: `.gitignore` (replace existing)
- Create: `frontend/package.json`
- Create: `frontend/astro.config.mjs`
- Create: `frontend/tsconfig.json`
- Create: `frontend/.env.example`
- Create: `frontend/.env`
- Create: `api/go.mod`
- Create: `api/.env.example`
- Create: `api/.env`

- [ ] **Step 1: Update `.gitignore` for monorepo**

Replace the root `.gitignore` with one covering both Astro and Go:

```gitignore
# Dependencies
node_modules/

# Build output
frontend/dist/
api/bin/

# Environment variables — NEVER commit these
.env
.env.local
.env.*.local
frontend/.env
api/.env

# OS
.DS_Store
Thumbs.db

# IDE
.vscode/
.idea/
*.swp
*.swo

# Go
api/tmp/

# Astro
.astro/

# Superpowers (brainstorm sessions)
.superpowers/

# Old CRA build
build/
```

- [ ] **Step 2: Scaffold Astro project**

Run from the repo root:

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio
npm create astro@latest frontend -- --template minimal --install --no-git --typescript strict
```

Expected: `frontend/` directory created with Astro boilerplate.

- [ ] **Step 3: Install frontend dependencies**

```bash
cd frontend
npx astro add react tailwind
npm install
```

Expected: `@astrojs/react`, `@astrojs/tailwind`, `react`, `react-dom` added to `package.json`.

- [ ] **Step 4: Create `frontend/astro.config.mjs`**

```javascript
import { defineConfig } from 'astro/config';
import react from '@astrojs/react';
import tailwind from '@astrojs/tailwind';

export default defineConfig({
  integrations: [
    react(),
    tailwind(),
  ],
  output: 'static',
});
```

- [ ] **Step 5: Create `frontend/tsconfig.json`**

```json
{
  "extends": "astro/tsconfigs/strict",
  "compilerOptions": {
    "jsx": "react-jsx",
    "jsxImportSource": "react",
    "baseUrl": ".",
    "paths": {
      "@components/*": ["src/components/*"],
      "@islands/*": ["src/islands/*"],
      "@layouts/*": ["src/layouts/*"],
      "@styles/*": ["src/styles/*"]
    }
  }
}
```

- [ ] **Step 6: Create frontend environment variable template**

Create `frontend/.env.example`:

```env
# URL of the Go API (e.g., http://localhost:8080 for local, https://your-api.fly.dev for production)
PUBLIC_API_URL=http://localhost:8080
```

Create `frontend/.env` (gitignored, local development):

```env
PUBLIC_API_URL=http://localhost:8080
```

- [ ] **Step 7: Initialize Go module**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/api
go mod init github.com/rfabreu/portfolio-api
```

Expected: `api/go.mod` created.

- [ ] **Step 8: Create API environment variable template**

Create `api/.env.example`:

```env
# Google Gemini API key — get one at https://aistudio.google.com/apikey
GEMINI_API_KEY=your-gemini-api-key-here

# Server configuration
PORT=8080
ALLOWED_ORIGINS=http://localhost:4321,https://your-site.netlify.app

# Rate limiting
RATE_LIMIT_RPM=30
```

Create `api/.env` (gitignored, local development):

```env
GEMINI_API_KEY=
PORT=8080
ALLOWED_ORIGINS=http://localhost:4321
RATE_LIMIT_RPM=30
```

- [ ] **Step 9: Migrate static assets**

```bash
mkdir -p frontend/public/images
cp /Users/rafaeldeabreugomes/projects/portfolio/public/images/* frontend/public/images/
cp /Users/rafaeldeabreugomes/projects/portfolio/public/favicon.ico frontend/public/
```

- [ ] **Step 10: Verify Astro builds**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/frontend
npm run build
```

Expected: Build completes, output in `frontend/dist/`.

- [ ] **Step 11: Commit**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio
git add .gitignore frontend/ api/go.mod api/.env.example
git commit -m "feat: scaffold monorepo with Astro frontend and Go API module"
```

---

## Task 2: Design System & Global Styles

**Files:**
- Create: `frontend/src/styles/global.css`
- Create: `frontend/tailwind.config.mjs`

- [ ] **Step 1: Create `frontend/tailwind.config.mjs` with design tokens**

```javascript
/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}'],
  theme: {
    extend: {
      colors: {
        base: '#050508',
        surface: 'rgba(255,255,255,0.03)',
        'surface-border': 'rgba(255,255,255,0.06)',
        'text-primary': '#ffffff',
        'text-secondary': '#a0a0b0',
        'text-muted': '#666666',
        accent: {
          indigo: '#6366f1',
          purple: '#8b5cf6',
          pink: '#ec4899',
          'indigo-light': '#a5b4fc',
          'purple-light': '#c4b5fd',
          'pink-light': '#f9a8d4',
        },
      },
      fontFamily: {
        mono: ["'Courier New'", 'monospace'],
      },
      letterSpacing: {
        'headline': '-0.06em',
        'tight-2': '-0.05em',
      },
      borderRadius: {
        card: '12px',
        tag: '5px',
        btn: '8px',
      },
    },
  },
  plugins: [],
};
```

- [ ] **Step 2: Create `frontend/src/styles/global.css`**

```css
@import 'tailwindcss';

@layer base {
  body {
    @apply bg-base text-text-secondary antialiased;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    line-height: 1.6;
  }

  h1, h2, h3, h4 {
    @apply text-text-primary font-extrabold;
    letter-spacing: -0.06em;
  }

  /* Scroll behavior */
  html {
    scroll-behavior: smooth;
  }

  @media (prefers-reduced-motion: reduce) {
    html {
      scroll-behavior: auto;
    }
    .animate-fade-up,
    .animate-fade-in,
    .animate-stagger {
      animation: none !important;
      opacity: 1 !important;
      transform: none !important;
    }
  }
}

@layer components {
  /* Gradient divider */
  .gradient-divider {
    height: 1px;
    background: linear-gradient(90deg, transparent, rgba(99,102,241,0.3), rgba(139,92,246,0.2), transparent);
  }

  /* Gradient orbs for atmospheric depth */
  .orb-indigo {
    background: radial-gradient(circle, rgba(99,102,241,0.12) 0%, transparent 70%);
  }
  .orb-purple {
    background: radial-gradient(circle, rgba(139,92,246,0.08) 0%, transparent 70%);
  }
  .orb-pink {
    background: radial-gradient(circle, rgba(236,72,153,0.06) 0%, transparent 70%);
  }
}

@layer utilities {
  /* Hero text entrance animations */
  @keyframes fade-up {
    from {
      opacity: 0;
      transform: translateY(20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  @keyframes fade-in {
    from { opacity: 0; }
    to { opacity: 1; }
  }

  .animate-fade-up {
    animation: fade-up 0.8s ease-out forwards;
    opacity: 0;
  }

  .animate-fade-in {
    animation: fade-in 0.6s ease-out forwards;
    opacity: 0;
  }

  .animate-stagger > * {
    opacity: 0;
    animation: fade-up 0.5s ease-out forwards;
  }

  .animate-stagger > *:nth-child(1) { animation-delay: 0.1s; }
  .animate-stagger > *:nth-child(2) { animation-delay: 0.2s; }
  .animate-stagger > *:nth-child(3) { animation-delay: 0.3s; }
  .animate-stagger > *:nth-child(4) { animation-delay: 0.4s; }
  .animate-stagger > *:nth-child(5) { animation-delay: 0.5s; }

  /* Scroll reveal */
  .reveal {
    opacity: 0;
    transform: translateY(24px);
    transition: opacity 0.6s ease-out, transform 0.6s ease-out;
  }
  .reveal.visible {
    opacity: 1;
    transform: translateY(0);
  }
}
```

- [ ] **Step 3: Verify Tailwind compiles with new config**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/frontend
npm run build
```

Expected: Build completes without errors.

- [ ] **Step 4: Commit**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio
git add frontend/tailwind.config.mjs frontend/src/styles/global.css
git commit -m "feat: add design system with Tailwind tokens and global styles"
```

---

## Task 3: Base Layout & Index Page Shell

**Files:**
- Create: `frontend/src/layouts/BaseLayout.astro`
- Create: `frontend/src/pages/index.astro`

- [ ] **Step 1: Create `frontend/src/layouts/BaseLayout.astro`**

```astro
---
interface Props {
  title?: string;
  description?: string;
}

const {
  title = 'Rafael Abreu | Software Engineer',
  description = 'Software Engineer building intelligent systems across the full stack. Go, Python, React, and the infrastructure that connects them.',
} = Astro.props;
---

<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="description" content={description} />
    <meta name="theme-color" content="#050508" />
    <link rel="icon" type="image/x-icon" href="/favicon.ico" />
    <title>{title}</title>
  </head>
  <body>
    <slot />

    <script>
      // Scroll reveal: observe .reveal elements and add .visible on intersection
      const observer = new IntersectionObserver(
        (entries) => {
          entries.forEach((entry) => {
            if (entry.isIntersecting) {
              entry.target.classList.add('visible');
              observer.unobserve(entry.target);
            }
          });
        },
        { threshold: 0.1 }
      );

      document.querySelectorAll('.reveal').forEach((el) => observer.observe(el));
    </script>
  </body>
</html>
```

- [ ] **Step 2: Create `frontend/src/pages/index.astro`**

```astro
---
import BaseLayout from '../layouts/BaseLayout.astro';
---

<BaseLayout>
  <main>
    <p class="text-text-primary text-center py-20 text-xl">Portfolio coming soon.</p>
  </main>
</BaseLayout>
```

- [ ] **Step 3: Verify page renders**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/frontend
npm run dev
```

Open `http://localhost:4321` — should show "Portfolio coming soon." on a dark background.

- [ ] **Step 4: Commit**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio
git add frontend/src/layouts/BaseLayout.astro frontend/src/pages/index.astro
git commit -m "feat: add base layout and index page shell"
```

---

## Task 4: Reusable UI Components

**Files:**
- Create: `frontend/src/components/SectionHeader.astro`
- Create: `frontend/src/components/GradientDivider.astro`
- Create: `frontend/src/components/TechTag.astro`

- [ ] **Step 1: Create `frontend/src/components/SectionHeader.astro`**

```astro
---
interface Props {
  number: string;
  label: string;
  title: string;
}

const { number, label, title } = Astro.props;
---

<div class="mb-8">
  <div class="font-mono text-accent-indigo text-xs tracking-[3px] mb-2">
    {number} // {label}
  </div>
  <h2 class="text-3xl md:text-4xl font-black tracking-headline">
    {title}
  </h2>
</div>
```

- [ ] **Step 2: Create `frontend/src/components/GradientDivider.astro`**

```astro
<div class="gradient-divider mx-6 md:mx-10" role="separator"></div>
```

- [ ] **Step 3: Create `frontend/src/components/TechTag.astro`**

```astro
---
interface Props {
  label: string;
  color?: 'indigo' | 'purple' | 'pink';
}

const { label, color = 'indigo' } = Astro.props;

const colorMap = {
  indigo: 'bg-accent-indigo/10 border-accent-indigo/25 text-accent-indigo-light',
  purple: 'bg-accent-purple/10 border-accent-purple/25 text-accent-purple-light',
  pink: 'bg-accent-pink/10 border-accent-pink/20 text-accent-pink-light',
};
---

<span class={`font-mono text-xs px-3 py-1 border rounded-tag ${colorMap[color]}`}>
  {label}
</span>
```

- [ ] **Step 4: Commit**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio
git add frontend/src/components/SectionHeader.astro frontend/src/components/GradientDivider.astro frontend/src/components/TechTag.astro
git commit -m "feat: add reusable UI components (SectionHeader, GradientDivider, TechTag)"
```

---

## Task 5: Navigation Component

**Files:**
- Create: `frontend/src/components/Nav.astro`
- Create: `frontend/src/islands/MobileNav.tsx`

- [ ] **Step 1: Create `frontend/src/components/Nav.astro`**

```astro
---
import MobileNav from '../islands/MobileNav.tsx';

const navItems = [
  { label: 'About', href: '#about' },
  { label: 'Projects', href: '#projects' },
  { label: 'Skills', href: '#skills' },
  { label: 'Contact', href: '#contact' },
];

const resumePath = '/Rafael_Abreu_Resume.pdf';
---

<header class="sticky top-0 z-50 bg-base/80 backdrop-blur-md border-b border-surface-border">
  <nav class="container mx-auto max-w-6xl px-6 py-4 flex items-center justify-between">
    <!-- Logo -->
    <a href="/" class="font-mono text-text-primary text-sm font-bold tracking-wide">
      RA<span class="text-accent-indigo">.</span>
    </a>

    <!-- Desktop nav -->
    <div class="hidden md:flex items-center gap-7">
      {navItems.map((item) => (
        <a
          href={item.href}
          class="text-text-muted text-xs tracking-[1.5px] uppercase hover:text-accent-indigo transition-colors duration-200"
        >
          {item.label}
        </a>
      ))}
    </div>

    <!-- Desktop actions -->
    <div class="hidden md:flex items-center gap-3">
      <a
        href={resumePath}
        target="_blank"
        rel="noopener noreferrer"
        class="px-4 py-1.5 border border-accent-indigo/40 rounded-btn text-accent-indigo text-xs tracking-wide hover:bg-accent-indigo/10 transition-colors duration-200"
      >
        RESUME
      </a>
      <a
        href="#contact"
        class="px-4 py-1.5 bg-accent-indigo rounded-btn text-white text-xs tracking-wide hover:bg-accent-indigo/90 transition-colors duration-200"
      >
        LET'S TALK
      </a>
    </div>

    <!-- Mobile hamburger -->
    <MobileNav client:load navItems={navItems} resumePath={resumePath} />
  </nav>
</header>
```

- [ ] **Step 2: Create `frontend/src/islands/MobileNav.tsx`**

```tsx
import { useState } from 'react';

interface NavItem {
  label: string;
  href: string;
}

interface Props {
  navItems: NavItem[];
  resumePath: string;
}

export default function MobileNav({ navItems, resumePath }: Props) {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div className="md:hidden">
      {/* Hamburger button */}
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="relative z-50 w-8 h-8 flex flex-col justify-center items-center gap-1.5"
        aria-label={isOpen ? 'Close menu' : 'Open menu'}
      >
        <span
          className={`block w-5 h-0.5 bg-white transition-all duration-300 ${
            isOpen ? 'rotate-45 translate-y-1' : ''
          }`}
        />
        <span
          className={`block w-5 h-0.5 bg-white transition-all duration-300 ${
            isOpen ? '-rotate-45 -translate-y-1' : ''
          }`}
        />
      </button>

      {/* Full-screen overlay */}
      {isOpen && (
        <div className="fixed inset-0 z-40 bg-base flex flex-col items-center justify-center gap-8">
          {navItems.map((item) => (
            <a
              key={item.label}
              href={item.href}
              onClick={() => setIsOpen(false)}
              className="text-text-primary text-2xl font-bold tracking-headline hover:text-accent-indigo transition-colors"
            >
              {item.label}
            </a>
          ))}
          <div className="flex gap-4 mt-4">
            <a
              href={resumePath}
              target="_blank"
              rel="noopener noreferrer"
              className="px-6 py-2 border border-accent-indigo/40 rounded-btn text-accent-indigo text-sm"
            >
              RESUME
            </a>
            <a
              href="#contact"
              onClick={() => setIsOpen(false)}
              className="px-6 py-2 bg-accent-indigo rounded-btn text-white text-sm"
            >
              LET'S TALK
            </a>
          </div>
        </div>
      )}
    </div>
  );
}
```

- [ ] **Step 3: Add Nav to index page**

Update `frontend/src/pages/index.astro`:

```astro
---
import BaseLayout from '../layouts/BaseLayout.astro';
import Nav from '../components/Nav.astro';
---

<BaseLayout>
  <Nav />
  <main>
    <p class="text-text-primary text-center py-20 text-xl">Sections coming next.</p>
  </main>
</BaseLayout>
```

- [ ] **Step 4: Verify nav renders on desktop and mobile**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/frontend
npm run dev
```

Open `http://localhost:4321`. Desktop: nav links + buttons visible. Resize to mobile: hamburger button, click opens full-screen overlay.

- [ ] **Step 5: Commit**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio
git add frontend/src/components/Nav.astro frontend/src/islands/MobileNav.tsx frontend/src/pages/index.astro
git commit -m "feat: add navigation with mobile hamburger menu"
```

---

## Task 6: Hero Section

**Files:**
- Create: `frontend/src/components/Hero.astro`

- [ ] **Step 1: Create `frontend/src/components/Hero.astro`**

```astro
---
import TechTag from './TechTag.astro';

const tags: { label: string; color: 'indigo' | 'purple' | 'pink' }[] = [
  { label: 'Go', color: 'indigo' },
  { label: 'Python', color: 'purple' },
  { label: 'React', color: 'pink' },
  { label: 'Kubernetes', color: 'indigo' },
  { label: 'AI/ML', color: 'purple' },
];
---

<section id="hero" class="relative min-h-[85vh] flex items-center overflow-hidden">
  <!-- Gradient orbs -->
  <div class="absolute -top-16 right-20 w-72 h-72 orb-indigo rounded-full pointer-events-none"></div>
  <div class="absolute -bottom-10 left-48 w-52 h-52 orb-purple rounded-full pointer-events-none"></div>
  <div class="absolute top-10 -left-10 w-44 h-44 orb-pink rounded-full pointer-events-none"></div>

  <div class="container mx-auto max-w-6xl px-6 relative z-10">
    <div class="font-mono text-accent-indigo text-xs tracking-[3px] mb-5 animate-fade-in" style="animation-delay: 0.1s;">
      // HELLO WORLD
    </div>

    <h1 class="text-5xl sm:text-7xl md:text-8xl font-black tracking-headline leading-none animate-fade-up" style="animation-delay: 0.2s;">
      RAFAEL<br />ABREU
    </h1>

    <p class="mt-5 text-text-secondary text-base md:text-lg max-w-lg leading-relaxed animate-fade-in" style="animation-delay: 0.5s;">
      Software Engineer building <span class="text-accent-purple-light">intelligent systems</span> across the full stack. Go, Python, React, and the infrastructure that connects them.
    </p>

    <div class="mt-6 flex flex-wrap gap-2.5 animate-stagger" style="animation-delay: 0.7s;">
      {tags.map((tag) => (
        <TechTag label={tag.label} color={tag.color} />
      ))}
    </div>

    <div class="mt-8 flex flex-wrap gap-4 animate-fade-in" style="animation-delay: 1s;">
      <a
        href="#projects"
        class="px-6 py-2.5 bg-accent-indigo rounded-btn text-white text-sm font-semibold hover:bg-accent-indigo/90 transition-colors duration-200"
      >
        View Projects
      </a>
      <button
        id="hero-chat-trigger"
        class="px-6 py-2.5 border border-white/15 rounded-btn text-gray-300 text-sm font-medium hover:border-accent-indigo/40 hover:text-accent-indigo transition-colors duration-200"
      >
        Ask my AI →
      </button>
    </div>
  </div>
</section>
```

- [ ] **Step 2: Add Hero to index page**

Update `frontend/src/pages/index.astro`:

```astro
---
import BaseLayout from '../layouts/BaseLayout.astro';
import Nav from '../components/Nav.astro';
import Hero from '../components/Hero.astro';
import GradientDivider from '../components/GradientDivider.astro';
---

<BaseLayout>
  <Nav />
  <main>
    <Hero />
    <GradientDivider />
  </main>
</BaseLayout>
```

- [ ] **Step 3: Verify hero renders with animations**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/frontend
npm run dev
```

Open `http://localhost:4321`. Hero should show gradient orbs, animated text entrance, tech tags, and two CTA buttons. Verify responsive: check mobile breakpoint for text sizing.

- [ ] **Step 4: Commit**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio
git add frontend/src/components/Hero.astro frontend/src/pages/index.astro
git commit -m "feat: add hero section with animated text and gradient orbs"
```

---

## Task 7: About Section

**Files:**
- Create: `frontend/src/components/About.astro`

- [ ] **Step 1: Create `frontend/src/components/About.astro`**

```astro
---
import SectionHeader from './SectionHeader.astro';
---

<section id="about" class="py-20 px-6">
  <div class="container mx-auto max-w-6xl reveal">
    <SectionHeader number="01" label="ABOUT" title="Who I Am" />

    <div class="grid md:grid-cols-2 gap-12 mt-8">
      <div class="space-y-5">
        <p class="text-text-secondary leading-relaxed">
          I'm a Software Engineer with experience across technology, television, and education.
          I build full-stack applications and the infrastructure that supports them — from React
          frontends to Go microservices to Kubernetes clusters.
        </p>
        <p class="text-text-secondary leading-relaxed">
          I completed the University of Toronto Coding Bootcamp and hold a Red Hat Kubernetes
          certificate (OpenShift). My current focus is on building intelligent systems that
          leverage AI to solve real problems.
        </p>
      </div>

      <div class="space-y-4">
        <div class="p-5 bg-surface border border-surface-border rounded-card">
          <div class="font-mono text-accent-indigo text-xs tracking-wider mb-2">FOCUS AREAS</div>
          <ul class="space-y-2 text-text-secondary text-sm">
            <li>→ Full-Stack Engineering (Go, Python, React)</li>
            <li>→ AI Integration & LLM Applications</li>
            <li>→ Cloud Infrastructure & DevOps</li>
            <li>→ System Design & Architecture</li>
          </ul>
        </div>
        <div class="p-5 bg-surface border border-surface-border rounded-card">
          <div class="font-mono text-accent-purple text-xs tracking-wider mb-2">EDUCATION</div>
          <ul class="space-y-2 text-text-secondary text-sm">
            <li>→ University of Toronto — Coding Bootcamp</li>
            <li>→ Red Hat — Kubernetes/OpenShift Certificate</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</section>
```

- [ ] **Step 2: Add About to index page**

Update `frontend/src/pages/index.astro` to add the About section after GradientDivider:

```astro
---
import BaseLayout from '../layouts/BaseLayout.astro';
import Nav from '../components/Nav.astro';
import Hero from '../components/Hero.astro';
import GradientDivider from '../components/GradientDivider.astro';
import About from '../components/About.astro';
---

<BaseLayout>
  <Nav />
  <main>
    <Hero />
    <GradientDivider />
    <About />
    <GradientDivider />
  </main>
</BaseLayout>
```

- [ ] **Step 3: Verify about section renders and scroll reveal works**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/frontend
npm run dev
```

Open `http://localhost:4321`. Scroll down past hero — About section should fade in. Check two-column layout on desktop, stacked on mobile.

- [ ] **Step 4: Commit**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio
git add frontend/src/components/About.astro frontend/src/pages/index.astro
git commit -m "feat: add about section with bio and focus areas"
```

---

## Task 8: Content Collections & Project Data

**Files:**
- Create: `frontend/src/content/config.ts`
- Create: `frontend/src/content/projects/fptv.md`
- Create: `frontend/src/content/projects/freckles-design.md`
- Create: `frontend/src/content/projects/charge-it-up.md`

- [ ] **Step 1: Create content collection schema**

Create `frontend/src/content/config.ts`:

```typescript
import { defineCollection, z } from 'astro:content';

const projects = defineCollection({
  type: 'content',
  schema: z.object({
    title: z.string(),
    description: z.string(),
    image: z.string(),
    tags: z.array(z.string()),
    category: z.string(),
    featured: z.boolean().default(false),
    order: z.number(),
    github: z.string().url().optional(),
    demo: z.string().url().optional(),
  }),
});

export const collections = { projects };
```

- [ ] **Step 2: Create `frontend/src/content/projects/fptv.md`**

```markdown
---
title: "FPTV"
description: "Toronto-based TV channel website with glassmorphism dashboard and live broadcasting features."
image: "/images/fptv_website.png"
tags: ["HTML", "CSS", "JavaScript", "Weebly"]
category: "Full-Stack"
featured: true
order: 1
github: "https://github.com/rfabreu/FPTV-WebFeatures"
demo: "https://fptv.ca"
---

## Overview

Developed the FPTV website, a Toronto-based TV channel broadcasting across Canada. The site provides seamless viewing with modern UI patterns including a glassmorphism dashboard for subscribers.

## Problem

FPTV needed a professional web presence that reflected their broadcasting quality while providing subscribers with an intuitive dashboard experience. The existing site lacked modern design patterns and subscriber tooling.

## Approach

Built custom frontend features using vanilla HTML, CSS, and JavaScript integrated with the Weebly (Square) platform. Focused on glassmorphism design for the subscriber dashboard and integrated the Environment Canada weather widget for local relevance.

## Tech Stack

- **HTML/CSS** — Semantic markup with modern CSS including backdrop-filter for glassmorphism
- **JavaScript** — Interactive dashboard components and widget integration
- **Weebly (Square)** — CMS platform for content management

## Outcome

Delivered a polished, professional website that serves FPTV's Canadian audience with an engaging subscriber dashboard and integrated broadcasting tools.
```

- [ ] **Step 3: Create `frontend/src/content/projects/freckles-design.md`**

```markdown
---
title: "Freckles Design"
description: "Artist portfolio and brand platform built in partnership with a fine artist and graphic designer."
image: "/images/freckles_design_website.png"
tags: ["React", "Bootstrap", "JavaScript", "Netlify"]
category: "Full-Stack"
featured: true
order: 2
github: "https://github.com/rfabreu/frecklesDesign"
demo: "https://frecklesdesign.ca/"
---

## Overview

A web platform developed in partnership with fine artist and graphic designer Carla Antunes to promote her brand and showcase her art. Built with React and deployed using CI/CD workflows.

## Problem

The artist needed a professional online presence to showcase their portfolio, attract clients, and establish their brand identity. The solution needed to be easy to update and maintain.

## Approach

Built with React and Bootstrap for rapid, responsive development. Deployed to Netlify with CI/CD for seamless updates. Worked closely with the artist to ensure the design reflected their creative vision.

## Tech Stack

- **React** — Component-based UI for gallery and portfolio views
- **Bootstrap** — Responsive grid and UI components
- **Netlify** — Hosting with CI/CD pipeline
- **Google Domains** — Custom domain management

## Outcome

Delivered a branded portfolio platform that the artist actively uses to showcase work and attract new clients. The CI/CD pipeline enables rapid content updates.
```

- [ ] **Step 4: Create `frontend/src/content/projects/charge-it-up.md`**

```markdown
---
title: "Charge It Up"
description: "EV charging station locator covering all of North America with real-time availability."
image: "/images/charge_it_up.jpg"
tags: ["HTML", "CSS", "Tailwind", "JavaScript"]
category: "Frontend"
featured: false
order: 3
github: "https://github.com/rfabreu/ev-mapper"
demo: "https://tarequem.github.io/ev-mapper/"
---

## Overview

Charge It Up is a web application that allows electric vehicle owners to quickly locate charging stations near their current location, covering all of North America.

## Problem

EV owners need a fast, reliable way to find nearby charging stations. Existing solutions were often slow or limited in coverage.

## Approach

Built a lightweight frontend application using vanilla HTML, CSS (Tailwind), and JavaScript. Integrated with mapping APIs to provide real-time station location data across North America.

## Tech Stack

- **HTML/CSS/Tailwind** — Responsive, utility-first styling
- **JavaScript** — Maps API integration and geolocation
- **Third-party APIs** — Charging station data and mapping

## Outcome

A fast, user-friendly tool that helps EV owners find charging stations anywhere in North America with minimal load time and an intuitive interface.
```

- [ ] **Step 5: Verify content collections load**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/frontend
npx astro check
```

Expected: No errors. Content collection schemas validated.

- [ ] **Step 6: Commit**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio
git add frontend/src/content/
git commit -m "feat: add content collections with project case studies"
```

---

## Task 9: Projects Section

**Files:**
- Create: `frontend/src/components/FeaturedProjectCard.astro`
- Create: `frontend/src/components/ProjectCard.astro`
- Create: `frontend/src/components/Projects.astro`

- [ ] **Step 1: Create `frontend/src/components/FeaturedProjectCard.astro`**

```astro
---
interface Props {
  title: string;
  description: string;
  image: string;
  tags: string[];
  category: string;
  slug: string;
}

const { title, description, image, tags, category, slug } = Astro.props;

const tagColors = ['indigo', 'purple', 'pink'] as const;
---

<a href={`/project/${slug}`} class="group block bg-surface border border-surface-border rounded-card overflow-hidden hover:border-accent-indigo/30 transition-colors duration-300">
  <div class="relative h-48 md:h-56 overflow-hidden">
    <img
      src={image}
      alt={title}
      class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
      loading="lazy"
    />
    <div class="absolute top-3 right-3 font-mono text-[10px] text-accent-indigo bg-accent-indigo/10 border border-accent-indigo/25 px-2 py-0.5 rounded-tag">
      {category}
    </div>
  </div>
  <div class="p-5">
    <h3 class="text-text-primary text-lg font-bold tracking-tight">{title}</h3>
    <p class="text-text-muted text-sm mt-2 line-clamp-2">{description}</p>
    <div class="mt-4 flex flex-wrap gap-1.5">
      {tags.map((tag, i) => (
        <span class={`font-mono text-[10px] px-2 py-0.5 rounded-tag bg-${['accent-indigo', 'accent-purple', 'accent-pink'][i % 3]}/10 text-${['accent-indigo', 'accent-purple', 'accent-pink'][i % 3]}-light border border-${['accent-indigo', 'accent-purple', 'accent-pink'][i % 3]}/20`}>
          {tag}
        </span>
      ))}
    </div>
  </div>
</a>
```

- [ ] **Step 2: Create `frontend/src/components/ProjectCard.astro`**

```astro
---
interface Props {
  title: string;
  description: string;
  tags: string[];
  github?: string;
  demo?: string;
}

const { title, description, tags, github, demo } = Astro.props;
---

<div class="bg-surface border border-surface-border rounded-card p-5 hover:border-accent-indigo/20 transition-colors duration-300">
  <h3 class="text-text-primary text-base font-bold tracking-tight">{title}</h3>
  <p class="text-text-muted text-sm mt-2 line-clamp-2">{description}</p>
  <div class="mt-3 flex flex-wrap gap-1.5">
    {tags.slice(0, 3).map((tag) => (
      <span class="font-mono text-[10px] px-2 py-0.5 rounded-tag bg-white/5 text-text-secondary">
        {tag}
      </span>
    ))}
  </div>
  <div class="mt-4 flex gap-3">
    {github && (
      <a
        href={github}
        target="_blank"
        rel="noopener noreferrer"
        class="text-text-muted text-xs hover:text-accent-indigo transition-colors"
      >
        GitHub →
      </a>
    )}
    {demo && demo !== '#' && (
      <a
        href={demo}
        target="_blank"
        rel="noopener noreferrer"
        class="text-text-muted text-xs hover:text-accent-pink transition-colors"
      >
        Live Demo →
      </a>
    )}
  </div>
</div>
```

- [ ] **Step 3: Create `frontend/src/components/Projects.astro`**

```astro
---
import { getCollection } from 'astro:content';
import SectionHeader from './SectionHeader.astro';
import FeaturedProjectCard from './FeaturedProjectCard.astro';
import ProjectCard from './ProjectCard.astro';

const allProjects = await getCollection('projects');
const sorted = allProjects.sort((a, b) => a.data.order - b.data.order);
const featured = sorted.filter((p) => p.data.featured);
const other = sorted.filter((p) => !p.data.featured);
---

<section id="projects" class="py-20 px-6">
  <div class="container mx-auto max-w-6xl reveal">
    <div class="flex justify-between items-baseline mb-8">
      <SectionHeader number="02" label="PROJECTS" title="Selected Work" />
    </div>

    <!-- Featured projects -->
    {featured.length > 0 && (
      <div class="grid md:grid-cols-2 gap-4 mb-6">
        {featured.map((project) => (
          <FeaturedProjectCard
            title={project.data.title}
            description={project.data.description}
            image={project.data.image}
            tags={project.data.tags}
            category={project.data.category}
            slug={project.slug}
          />
        ))}
      </div>
    )}

    <!-- Other projects -->
    {other.length > 0 && (
      <div class="grid sm:grid-cols-2 lg:grid-cols-3 gap-4">
        {other.map((project) => (
          <ProjectCard
            title={project.data.title}
            description={project.data.description}
            tags={project.data.tags}
            github={project.data.github}
            demo={project.data.demo}
          />
        ))}
      </div>
    )}
  </div>
</section>
```

- [ ] **Step 4: Add Projects to index page**

Update `frontend/src/pages/index.astro` to include Projects after About:

```astro
---
import BaseLayout from '../layouts/BaseLayout.astro';
import Nav from '../components/Nav.astro';
import Hero from '../components/Hero.astro';
import GradientDivider from '../components/GradientDivider.astro';
import About from '../components/About.astro';
import Projects from '../components/Projects.astro';
---

<BaseLayout>
  <Nav />
  <main>
    <Hero />
    <GradientDivider />
    <About />
    <GradientDivider />
    <Projects />
    <GradientDivider />
  </main>
</BaseLayout>
```

- [ ] **Step 5: Verify projects render with two tiers**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/frontend
npm run dev
```

Open `http://localhost:4321`. Scroll to Projects — should show 2 featured cards (FPTV, Freckles Design) as large cards and 1 other card (Charge It Up) in a smaller grid.

- [ ] **Step 6: Commit**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio
git add frontend/src/components/FeaturedProjectCard.astro frontend/src/components/ProjectCard.astro frontend/src/components/Projects.astro frontend/src/pages/index.astro
git commit -m "feat: add two-tier projects section with content collections"
```

---

## Task 10: Project Case Study Pages

**Files:**
- Create: `frontend/src/layouts/ProjectLayout.astro`
- Create: `frontend/src/pages/project/[slug].astro`

- [ ] **Step 1: Create `frontend/src/layouts/ProjectLayout.astro`**

```astro
---
import Nav from '../components/Nav.astro';
import Footer from '../components/Footer.astro';
import GradientDivider from '../components/GradientDivider.astro';

interface Props {
  title: string;
  description: string;
  image: string;
  tags: string[];
  category: string;
  github?: string;
  demo?: string;
}

const { title, description, image, tags, category, github, demo } = Astro.props;
---

<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="description" content={description} />
    <meta name="theme-color" content="#050508" />
    <link rel="icon" type="image/x-icon" href="/favicon.ico" />
    <title>{title} | Rafael Abreu</title>
  </head>
  <body class="bg-base text-text-secondary">
    <Nav />

    <article class="container mx-auto max-w-4xl px-6 py-16">
      <!-- Header -->
      <a href="/#projects" class="text-text-muted text-sm hover:text-accent-indigo transition-colors mb-8 inline-block">
        ← Back to Projects
      </a>

      <div class="font-mono text-accent-indigo text-xs tracking-[3px] mb-3">{category}</div>
      <h1 class="text-4xl md:text-5xl font-black tracking-headline mb-4">{title}</h1>
      <p class="text-text-secondary text-lg mb-6">{description}</p>

      <div class="flex flex-wrap gap-2 mb-6">
        {tags.map((tag, i) => {
          const colors = ['accent-indigo', 'accent-purple', 'accent-pink'];
          const c = colors[i % 3];
          return (
            <span class={`font-mono text-xs px-3 py-1 border rounded-tag bg-${c}/10 border-${c}/25 text-${c}-light`}>
              {tag}
            </span>
          );
        })}
      </div>

      <div class="flex gap-4 mb-10">
        {github && (
          <a href={github} target="_blank" rel="noopener noreferrer" class="px-5 py-2 border border-surface-border rounded-btn text-text-secondary text-sm hover:border-accent-indigo/40 hover:text-accent-indigo transition-colors">
            View Code →
          </a>
        )}
        {demo && demo !== '#' && (
          <a href={demo} target="_blank" rel="noopener noreferrer" class="px-5 py-2 bg-accent-indigo rounded-btn text-white text-sm hover:bg-accent-indigo/90 transition-colors">
            Live Demo →
          </a>
        )}
      </div>

      <!-- Project image -->
      <img src={image} alt={title} class="w-full rounded-card border border-surface-border mb-12" loading="lazy" />

      <GradientDivider />

      <!-- Markdown content -->
      <div class="prose-custom mt-12">
        <slot />
      </div>
    </article>

    <Footer />

    <style>
      .prose-custom :global(h2) {
        color: #ffffff;
        font-size: 1.5rem;
        font-weight: 800;
        letter-spacing: -0.04em;
        margin-top: 2.5rem;
        margin-bottom: 1rem;
      }
      .prose-custom :global(p) {
        color: #a0a0b0;
        line-height: 1.7;
        margin-bottom: 1rem;
      }
      .prose-custom :global(ul) {
        color: #a0a0b0;
        list-style: none;
        padding-left: 0;
      }
      .prose-custom :global(li) {
        padding-left: 1.25rem;
        position: relative;
        margin-bottom: 0.5rem;
      }
      .prose-custom :global(li::before) {
        content: '→';
        position: absolute;
        left: 0;
        color: #6366f1;
      }
      .prose-custom :global(strong) {
        color: #ffffff;
        font-weight: 600;
      }
      .prose-custom :global(code) {
        font-family: 'Courier New', monospace;
        background: rgba(255,255,255,0.05);
        padding: 0.15rem 0.4rem;
        border-radius: 4px;
        font-size: 0.875rem;
      }
    </style>
  </body>
</html>
```

- [ ] **Step 2: Create `frontend/src/pages/project/[slug].astro`**

```astro
---
import { getCollection } from 'astro:content';
import ProjectLayout from '../../layouts/ProjectLayout.astro';

export async function getStaticPaths() {
  const projects = await getCollection('projects');
  return projects.map((project) => ({
    params: { slug: project.slug },
    props: { project },
  }));
}

const { project } = Astro.props;
const { Content } = await project.render();
---

<ProjectLayout
  title={project.data.title}
  description={project.data.description}
  image={project.data.image}
  tags={project.data.tags}
  category={project.data.category}
  github={project.data.github}
  demo={project.data.demo}
>
  <Content />
</ProjectLayout>
```

- [ ] **Step 3: Verify case study pages render**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/frontend
npm run dev
```

Open `http://localhost:4321/project/fptv`. Should show the full case study with header, tags, image, and rendered markdown content. Click "Back to Projects" to verify navigation.

- [ ] **Step 4: Commit**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio
git add frontend/src/layouts/ProjectLayout.astro frontend/src/pages/project/
git commit -m "feat: add project case study pages with dynamic routing"
```

---

## Task 11: Skills Section

**Files:**
- Create: `frontend/src/components/Skills.astro`

- [ ] **Step 1: Create `frontend/src/components/Skills.astro`**

```astro
---
import SectionHeader from './SectionHeader.astro';

const categories = [
  {
    label: 'LANGUAGES',
    color: 'text-accent-indigo',
    skills: ['Go', 'Python', 'JavaScript', 'TypeScript', 'HTML', 'CSS'],
  },
  {
    label: 'FRAMEWORKS & TOOLS',
    color: 'text-accent-purple',
    skills: ['React', 'Astro', 'Node.js', 'Tailwind CSS', 'Bootstrap', 'Git'],
  },
  {
    label: 'INFRASTRUCTURE',
    color: 'text-accent-pink',
    skills: ['Kubernetes', 'Docker', 'CI/CD', 'Netlify', 'Fly.io'],
  },
  {
    label: 'AI & DATA',
    color: 'text-accent-indigo',
    skills: ['LLM Integration', 'Gemini API', 'RAG', 'MongoDB', 'MySQL'],
  },
  {
    label: 'PRACTICES',
    color: 'text-accent-purple',
    skills: ['TDD', 'OOP', 'MVC', 'System Design', 'Agile'],
  },
];
---

<section id="skills" class="py-20 px-6">
  <div class="container mx-auto max-w-6xl reveal">
    <SectionHeader number="03" label="EXPERTISE" title="Tech Stack" />

    <div class="grid sm:grid-cols-2 lg:grid-cols-3 gap-4 mt-8">
      {categories.map((cat) => (
        <div class="bg-surface border border-surface-border rounded-card p-5">
          <div class={`font-mono text-xs font-bold tracking-wider mb-3 ${cat.color}`}>
            {cat.label}
          </div>
          <div class="flex flex-wrap gap-2">
            {cat.skills.map((skill) => (
              <span class="font-mono text-xs text-text-secondary bg-white/5 px-2.5 py-1 rounded-tag">
                {skill}
              </span>
            ))}
          </div>
        </div>
      ))}
    </div>
  </div>
</section>
```

- [ ] **Step 2: Add Skills to index page**

Update `frontend/src/pages/index.astro`:

```astro
---
import BaseLayout from '../layouts/BaseLayout.astro';
import Nav from '../components/Nav.astro';
import Hero from '../components/Hero.astro';
import GradientDivider from '../components/GradientDivider.astro';
import About from '../components/About.astro';
import Projects from '../components/Projects.astro';
import Skills from '../components/Skills.astro';
---

<BaseLayout>
  <Nav />
  <main>
    <Hero />
    <GradientDivider />
    <About />
    <GradientDivider />
    <Projects />
    <GradientDivider />
    <Skills />
    <GradientDivider />
  </main>
</BaseLayout>
```

- [ ] **Step 3: Verify skills section renders**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/frontend
npm run dev
```

Open `http://localhost:4321`. Scroll to Skills — should show 5 categorized cards in a responsive grid.

- [ ] **Step 4: Commit**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio
git add frontend/src/components/Skills.astro frontend/src/pages/index.astro
git commit -m "feat: add categorized skills section"
```

---

## Task 12: Contact Section & Footer

**Files:**
- Create: `frontend/src/components/Contact.astro`
- Create: `frontend/src/components/Footer.astro`

- [ ] **Step 1: Create `frontend/src/components/Contact.astro`**

```astro
---
import SectionHeader from './SectionHeader.astro';
---

<section id="contact" class="py-20 px-6">
  <div class="container mx-auto max-w-6xl reveal">
    <SectionHeader number="04" label="CONTACT" title="Let's Connect" />

    <div class="grid md:grid-cols-2 gap-12 mt-8">
      <!-- Contact form (Netlify Forms) -->
      <form
        name="contact"
        method="POST"
        data-netlify="true"
        netlify-honeypot="bot-field"
        class="space-y-5"
      >
        <input type="hidden" name="form-name" value="contact" />
        <p class="hidden">
          <label>Don't fill this out: <input name="bot-field" /></label>
        </p>

        <div>
          <label for="name" class="font-mono text-xs text-text-muted tracking-wider block mb-2">NAME</label>
          <input
            type="text"
            id="name"
            name="name"
            required
            class="w-full bg-surface border border-surface-border rounded-btn px-4 py-3 text-text-primary text-sm focus:border-accent-indigo/60 focus:ring-1 focus:ring-accent-indigo/30 outline-none transition-colors"
            placeholder="Your name"
          />
        </div>

        <div>
          <label for="email" class="font-mono text-xs text-text-muted tracking-wider block mb-2">EMAIL</label>
          <input
            type="email"
            id="email"
            name="email"
            required
            class="w-full bg-surface border border-surface-border rounded-btn px-4 py-3 text-text-primary text-sm focus:border-accent-indigo/60 focus:ring-1 focus:ring-accent-indigo/30 outline-none transition-colors"
            placeholder="your@email.com"
          />
        </div>

        <div>
          <label for="message" class="font-mono text-xs text-text-muted tracking-wider block mb-2">MESSAGE</label>
          <textarea
            id="message"
            name="message"
            required
            rows="5"
            class="w-full bg-surface border border-surface-border rounded-btn px-4 py-3 text-text-primary text-sm focus:border-accent-indigo/60 focus:ring-1 focus:ring-accent-indigo/30 outline-none transition-colors resize-none"
            placeholder="Tell me about your project or just say hello..."
          ></textarea>
        </div>

        <button
          type="submit"
          class="px-8 py-3 bg-accent-indigo rounded-btn text-white text-sm font-semibold hover:bg-accent-indigo/90 transition-colors duration-200"
        >
          Send Message
        </button>
      </form>

      <!-- Contact info -->
      <div class="space-y-6">
        <div>
          <div class="font-mono text-xs text-accent-indigo tracking-wider mb-3">GET IN TOUCH</div>
          <p class="text-text-secondary text-sm leading-relaxed">
            I'm always interested in hearing about new opportunities, collaborations, or just connecting with fellow engineers. Drop me a message or find me on the platforms below.
          </p>
        </div>

        <div class="space-y-4">
          <a href="mailto:raabreugomes@gmail.com" class="flex items-center gap-3 text-text-secondary text-sm hover:text-accent-indigo transition-colors group">
            <span class="font-mono text-accent-indigo/60 group-hover:text-accent-indigo">→</span>
            raabreugomes@gmail.com
          </a>
          <a href="https://github.com/rfabreu" target="_blank" rel="noopener noreferrer" class="flex items-center gap-3 text-text-secondary text-sm hover:text-accent-purple transition-colors group">
            <span class="font-mono text-accent-purple/60 group-hover:text-accent-purple">→</span>
            github.com/rfabreu
          </a>
          <a href="https://linkedin.com/in/rafael-a-gomes" target="_blank" rel="noopener noreferrer" class="flex items-center gap-3 text-text-secondary text-sm hover:text-accent-pink transition-colors group">
            <span class="font-mono text-accent-pink/60 group-hover:text-accent-pink">→</span>
            linkedin.com/in/rafael-a-gomes
          </a>
        </div>

        <div class="p-5 bg-surface border border-surface-border rounded-card">
          <div class="font-mono text-xs text-text-muted tracking-wider mb-2">LOCATION</div>
          <p class="text-text-secondary text-sm">Toronto, Ontario, Canada</p>
        </div>
      </div>
    </div>
  </div>
</section>
```

- [ ] **Step 2: Create `frontend/src/components/Footer.astro`**

```astro
---
const currentYear = new Date().getFullYear();
---

<footer class="border-t border-surface-border py-8 px-6">
  <div class="container mx-auto max-w-6xl flex flex-col sm:flex-row justify-between items-center gap-4">
    <div class="font-mono text-text-muted text-xs">
      &copy; {currentYear} Rafael Abreu
    </div>
    <div class="flex gap-6">
      <a href="https://github.com/rfabreu" target="_blank" rel="noopener noreferrer" class="text-text-muted text-xs hover:text-accent-indigo transition-colors">
        GitHub
      </a>
      <a href="https://linkedin.com/in/rafael-a-gomes" target="_blank" rel="noopener noreferrer" class="text-text-muted text-xs hover:text-accent-purple transition-colors">
        LinkedIn
      </a>
      <a href="mailto:raabreugomes@gmail.com" class="text-text-muted text-xs hover:text-accent-pink transition-colors">
        Email
      </a>
    </div>
  </div>
</footer>
```

- [ ] **Step 3: Add Contact and Footer to index page**

Update `frontend/src/pages/index.astro`:

```astro
---
import BaseLayout from '../layouts/BaseLayout.astro';
import Nav from '../components/Nav.astro';
import Hero from '../components/Hero.astro';
import GradientDivider from '../components/GradientDivider.astro';
import About from '../components/About.astro';
import Projects from '../components/Projects.astro';
import Skills from '../components/Skills.astro';
import Contact from '../components/Contact.astro';
import Footer from '../components/Footer.astro';
---

<BaseLayout>
  <Nav />
  <main>
    <Hero />
    <GradientDivider />
    <About />
    <GradientDivider />
    <Projects />
    <GradientDivider />
    <Skills />
    <GradientDivider />
    <Contact />
  </main>
  <Footer />
</BaseLayout>
```

- [ ] **Step 4: Verify complete page renders end-to-end**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/frontend
npm run dev
```

Open `http://localhost:4321`. Scroll through entire page: Hero → About → Projects → Skills → Contact → Footer. Verify all sections, scroll reveals, links, and responsive layout.

- [ ] **Step 5: Commit**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio
git add frontend/src/components/Contact.astro frontend/src/components/Footer.astro frontend/src/pages/index.astro
git commit -m "feat: add contact section with Netlify Forms and footer"
```

---

## Task 13: Go API — Server Foundation

**Files:**
- Create: `api/cmd/server/main.go`
- Create: `api/internal/middleware/cors.go`
- Create: `api/internal/middleware/ratelimit.go`
- Create: `api/internal/middleware/ratelimit_test.go`

- [ ] **Step 1: Write the rate limiter test**

Create `api/internal/middleware/ratelimit_test.go`:

```go
package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rfabreu/portfolio-api/internal/middleware"
)

func TestRateLimiter_AllowsUnderLimit(t *testing.T) {
	handler := middleware.RateLimit(2)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/chat", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("request %d: expected 200, got %d", i, rec.Code)
		}
	}
}

func TestRateLimiter_BlocksOverLimit(t *testing.T) {
	handler := middleware.RateLimit(1)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// First request: allowed
	req := httptest.NewRequest(http.MethodPost, "/api/chat", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("first request: expected 200, got %d", rec.Code)
	}

	// Second request: blocked
	req = httptest.NewRequest(http.MethodPost, "/api/chat", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("second request: expected 429, got %d", rec.Code)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/api
go test ./internal/middleware/ -v
```

Expected: FAIL — package/function not found.

- [ ] **Step 3: Implement rate limiter**

Create `api/internal/middleware/ratelimit.go`:

```go
package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type visitor struct {
	count    int
	lastSeen time.Time
}

// RateLimit returns middleware that limits requests per IP per minute.
func RateLimit(rpm int) func(http.Handler) http.Handler {
	var mu sync.Mutex
	visitors := make(map[string]*visitor)

	// Cleanup stale entries every 3 minutes
	go func() {
		for {
			time.Sleep(3 * time.Minute)
			mu.Lock()
			for ip, v := range visitors {
				if time.Since(v.lastSeen) > time.Minute {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			if ip == "" {
				ip = r.RemoteAddr
			}

			mu.Lock()
			v, exists := visitors[ip]
			if !exists || time.Since(v.lastSeen) > time.Minute {
				visitors[ip] = &visitor{count: 1, lastSeen: time.Now()}
				mu.Unlock()
				next.ServeHTTP(w, r)
				return
			}

			v.count++
			v.lastSeen = time.Now()

			if v.count > rpm {
				mu.Unlock()
				http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
				return
			}
			mu.Unlock()
			next.ServeHTTP(w, r)
		})
	}
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/api
go test ./internal/middleware/ -v
```

Expected: PASS — both tests pass.

- [ ] **Step 5: Create CORS middleware**

Create `api/internal/middleware/cors.go`:

```go
package middleware

import (
	"net/http"
	"strings"
)

// CORS returns middleware that sets CORS headers based on allowed origins.
func CORS(allowedOrigins string) func(http.Handler) http.Handler {
	origins := strings.Split(allowedOrigins, ",")
	originSet := make(map[string]bool, len(origins))
	for _, o := range origins {
		originSet[strings.TrimSpace(o)] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if originSet[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
				w.Header().Set("Access-Control-Max-Age", "86400")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
```

- [ ] **Step 6: Create server entry point**

Create `api/cmd/server/main.go`:

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/rfabreu/portfolio-api/internal/middleware"
)

func main() {
	port := getEnv("PORT", "8080")
	allowedOrigins := getEnv("ALLOWED_ORIGINS", "http://localhost:4321")
	rpm, _ := strconv.Atoi(getEnv("RATE_LIMIT_RPM", "30"))

	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	// Apply middleware
	handler := middleware.CORS(allowedOrigins)(middleware.RateLimit(rpm)(mux))

	log.Printf("API server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
```

- [ ] **Step 7: Install dependencies and verify server starts**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/api
go mod tidy
go build ./cmd/server
```

Expected: Binary compiles without errors.

- [ ] **Step 8: Test health endpoint**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/api
go run ./cmd/server &
sleep 1
curl -s http://localhost:8080/health
kill %1
```

Expected: `{"status":"ok"}`

- [ ] **Step 9: Commit**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio
git add api/
git commit -m "feat: add Go API server with health endpoint, CORS, and rate limiting"
```

---

## Task 14: Go API — AI Chatbot

**Files:**
- Create: `api/internal/chat/provider.go`
- Create: `api/internal/chat/provider_test.go`
- Create: `api/internal/chat/context.go`
- Create: `api/internal/chat/handler.go`
- Create: `api/internal/chat/handler_test.go`
- Modify: `api/cmd/server/main.go`

- [ ] **Step 1: Write the provider interface and mock test**

Create `api/internal/chat/provider_test.go`:

```go
package chat_test

import (
	"context"
	"testing"

	"github.com/rfabreu/portfolio-api/internal/chat"
)

type mockProvider struct {
	response string
	err      error
}

func (m *mockProvider) GenerateResponse(ctx context.Context, systemPrompt string, messages []chat.Message) (string, error) {
	return m.response, m.err
}

func TestMockProvider_ReturnsResponse(t *testing.T) {
	provider := &mockProvider{response: "Hello! I'm Rafael's AI assistant."}
	resp, err := provider.GenerateResponse(context.Background(), "system", []chat.Message{
		{Role: "user", Content: "Hi"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "Hello! I'm Rafael's AI assistant." {
		t.Errorf("expected greeting, got: %s", resp)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/api
go test ./internal/chat/ -v
```

Expected: FAIL — package/types not found.

- [ ] **Step 3: Implement provider interface**

Create `api/internal/chat/provider.go`:

```go
package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Message represents a chat message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Provider is the interface for LLM backends. Swap implementations to change providers.
type Provider interface {
	GenerateResponse(ctx context.Context, systemPrompt string, messages []Message) (string, error)
}

// GeminiProvider calls the Google Gemini API.
type GeminiProvider struct {
	APIKey string
	Model  string
}

// NewGeminiProvider creates a provider using the given API key.
// Model defaults to "gemini-2.0-flash" if empty.
func NewGeminiProvider(apiKey, model string) *GeminiProvider {
	if model == "" {
		model = "gemini-2.0-flash"
	}
	return &GeminiProvider{APIKey: apiKey, Model: model}
}

func (g *GeminiProvider) GenerateResponse(ctx context.Context, systemPrompt string, messages []Message) (string, error) {
	// Build Gemini API request
	var contents []map[string]any

	// System instruction is passed separately in Gemini
	for _, msg := range messages {
		role := msg.Role
		if role == "assistant" {
			role = "model"
		}
		contents = append(contents, map[string]any{
			"role":  role,
			"parts": []map[string]string{{"text": msg.Content}},
		})
	}

	body := map[string]any{
		"contents": contents,
		"systemInstruction": map[string]any{
			"parts": []map[string]string{{"text": systemPrompt}},
		},
		"generationConfig": map[string]any{
			"maxOutputTokens": 1024,
			"temperature":     0.7,
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", g.Model, g.APIKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("gemini request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gemini API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	// Parse Gemini response
	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from Gemini")
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/api
go test ./internal/chat/ -v
```

Expected: PASS.

- [ ] **Step 5: Create resume context**

Create `api/internal/chat/context.go`:

```go
package chat

// SystemPrompt returns the system prompt with Rafael's professional context.
func SystemPrompt() string {
	return `You are Rafael Abreu's AI assistant on his portfolio website. You answer questions about Rafael's professional background, skills, projects, and experience.

ABOUT RAFAEL:
- Software Engineer based in Toronto, Ontario, Canada
- Full-stack developer with experience across technology, television, and education industries
- Currently focused on Go, Python, and AI/ML engineering
- Completed the University of Toronto Coding Bootcamp
- Holds a Red Hat Kubernetes certificate (OpenShift training)

TECHNICAL SKILLS:
- Languages: Go, Python, JavaScript, TypeScript, HTML, CSS
- Frameworks: React, Astro, Node.js, Tailwind CSS, Bootstrap
- Infrastructure: Kubernetes, Docker, CI/CD, Netlify, Fly.io
- AI & Data: LLM Integration, Gemini API, RAG, MongoDB, MySQL
- Practices: TDD, OOP, MVC, System Design, Agile

NOTABLE PROJECTS:
- FPTV: Developed the website for a Toronto-based TV channel broadcasting across Canada. Features glassmorphism dashboard. (HTML, CSS, JavaScript, Weebly)
- Freckles Design: Artist portfolio platform built in partnership with a fine artist. (React, Bootstrap, Netlify)
- Charge It Up: EV charging station locator for North America. (HTML, CSS, Tailwind, JavaScript)
- This portfolio itself: Astro frontend + Go API backend with AI chatbot, demonstrating modern architecture.

RULES:
- Only answer questions about Rafael's professional background, skills, projects, and experience.
- Be concise and helpful. Keep responses under 200 words.
- If asked something unrelated to Rafael's professional life, politely redirect: "I'm here to help with questions about Rafael's experience and work. What would you like to know?"
- Be friendly and professional in tone.
- You may suggest that visitors check out specific projects or sections of the portfolio when relevant.`
}
```

- [ ] **Step 6: Write the handler test**

Create `api/internal/chat/handler_test.go`:

```go
package chat_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rfabreu/portfolio-api/internal/chat"
)

func TestChatHandler_ReturnsResponse(t *testing.T) {
	provider := &mockProvider{response: "Rafael is a Software Engineer based in Toronto."}
	handler := chat.NewHandler(provider)

	reqBody := chat.ChatRequest{
		Messages: []chat.Message{{Role: "user", Content: "Who is Rafael?"}},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp chat.ChatResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Message == "" {
		t.Error("expected non-empty message in response")
	}
}

func TestChatHandler_RejectsEmptyMessages(t *testing.T) {
	provider := &mockProvider{response: "test"}
	handler := chat.NewHandler(provider)

	reqBody := chat.ChatRequest{Messages: []chat.Message{}}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestChatHandler_EnforcesMessageLimit(t *testing.T) {
	provider := &mockProvider{response: "test"}
	handler := chat.NewHandler(provider)

	// Create 21 messages (over the 20 limit)
	msgs := make([]chat.Message, 21)
	for i := range msgs {
		msgs[i] = chat.Message{Role: "user", Content: "hello"}
	}

	reqBody := chat.ChatRequest{Messages: msgs}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for too many messages, got %d", rec.Code)
	}
}
```

- [ ] **Step 7: Run tests to verify they fail**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/api
go test ./internal/chat/ -v
```

Expected: FAIL — `ChatRequest`, `ChatResponse`, `NewHandler` not found.

- [ ] **Step 8: Implement chat handler**

Create `api/internal/chat/handler.go`:

```go
package chat

import (
	"encoding/json"
	"net/http"
)

const maxMessages = 20

// ChatRequest is the JSON body for POST /api/chat.
type ChatRequest struct {
	Messages []Message `json:"messages"`
}

// ChatResponse is the JSON response from the chat endpoint.
type ChatResponse struct {
	Message string `json:"message"`
}

// NewHandler returns an http.Handler for the chat endpoint.
func NewHandler(provider Provider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}

		if len(req.Messages) == 0 {
			http.Error(w, `{"error":"messages required"}`, http.StatusBadRequest)
			return
		}

		if len(req.Messages) > maxMessages {
			http.Error(w, `{"error":"conversation limit exceeded (max 20 messages)"}`, http.StatusBadRequest)
			return
		}

		response, err := provider.GenerateResponse(r.Context(), SystemPrompt(), req.Messages)
		if err != nil {
			http.Error(w, `{"error":"failed to generate response"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ChatResponse{Message: response})
	})
}
```

- [ ] **Step 9: Run tests to verify they pass**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/api
go test ./internal/chat/ -v
```

Expected: PASS — all 4 tests pass.

- [ ] **Step 10: Register chat route in server**

Update `api/cmd/server/main.go` — add the chat handler import and route:

Add import:
```go
import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/rfabreu/portfolio-api/internal/chat"
	"github.com/rfabreu/portfolio-api/internal/middleware"
)
```

Add route registration after the health check in `main()`:

```go
	// Chat endpoint
	geminiKey := os.Getenv("GEMINI_API_KEY")
	if geminiKey == "" {
		log.Println("WARNING: GEMINI_API_KEY not set — chat endpoint will fail")
	}
	provider := chat.NewGeminiProvider(geminiKey, "")
	mux.Handle("POST /api/chat", chat.NewHandler(provider))
```

- [ ] **Step 11: Verify server compiles with chat route**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/api
go mod tidy
go build ./cmd/server
```

Expected: Compiles without errors.

- [ ] **Step 12: Run all API tests**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/api
go test ./... -v
```

Expected: All tests pass.

- [ ] **Step 13: Commit**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio
git add api/
git commit -m "feat: add AI chatbot with Gemini provider, handler, and tests"
```

---

## Task 15: ChatWidget React Island (Frontend)

**Files:**
- Create: `frontend/src/islands/ChatWidget.tsx`
- Modify: `frontend/src/pages/index.astro`

- [ ] **Step 1: Create `frontend/src/islands/ChatWidget.tsx`**

```tsx
import { useState, useRef, useEffect } from 'react';

interface Message {
  role: 'user' | 'assistant';
  content: string;
}

export default function ChatWidget() {
  const [isOpen, setIsOpen] = useState(false);
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const API_URL = import.meta.env.PUBLIC_API_URL || 'http://localhost:8080';

  useEffect(() => {
    if (isOpen && inputRef.current) {
      inputRef.current.focus();
    }
  }, [isOpen]);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  // Listen for hero CTA click
  useEffect(() => {
    const trigger = document.getElementById('hero-chat-trigger');
    if (trigger) {
      const handler = () => setIsOpen(true);
      trigger.addEventListener('click', handler);
      return () => trigger.removeEventListener('click', handler);
    }
  }, []);

  async function sendMessage(e: React.FormEvent) {
    e.preventDefault();
    const trimmed = input.trim();
    if (!trimmed || isLoading) return;

    const userMessage: Message = { role: 'user', content: trimmed };
    const updatedMessages = [...messages, userMessage];
    setMessages(updatedMessages);
    setInput('');
    setIsLoading(true);
    setError(null);

    try {
      const resp = await fetch(`${API_URL}/api/chat`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ messages: updatedMessages }),
      });

      if (!resp.ok) {
        throw new Error(resp.status === 429 ? 'Rate limit reached. Try again in a minute.' : 'Something went wrong.');
      }

      const data = await resp.json();
      setMessages([...updatedMessages, { role: 'assistant', content: data.message }]);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to connect. Try again later.');
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <>
      {/* Floating button */}
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="fixed bottom-5 right-5 z-50 w-[52px] h-[52px] rounded-[14px] flex items-center justify-center text-white text-xl shadow-lg transition-transform hover:scale-105"
        style={{ background: 'linear-gradient(135deg, #6366f1, #8b5cf6)' }}
        aria-label={isOpen ? 'Close chat' : 'Open chat'}
      >
        {isOpen ? '✕' : '💬'}
      </button>

      {/* Chat panel */}
      {isOpen && (
        <div className="fixed bottom-20 right-5 z-50 w-[360px] max-h-[500px] bg-[#0a0a0f] border border-white/10 rounded-xl shadow-2xl flex flex-col overflow-hidden">
          {/* Header */}
          <div className="px-4 py-3 border-b border-white/10">
            <div className="text-white text-sm font-bold">Ask Rafael's AI</div>
            <div className="text-gray-500 text-xs mt-0.5">Ask about experience, projects, or skills</div>
          </div>

          {/* Messages */}
          <div className="flex-1 overflow-y-auto p-4 space-y-3 min-h-[200px] max-h-[340px]">
            {messages.length === 0 && (
              <div className="text-gray-500 text-sm">
                Hey! I'm Rafael's AI assistant. Ask me about his experience, projects, or tech stack.
              </div>
            )}

            {messages.map((msg, i) => (
              <div key={i} className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}>
                <div
                  className={`max-w-[85%] px-3 py-2 rounded-lg text-sm ${
                    msg.role === 'user'
                      ? 'bg-[#6366f1] text-white'
                      : 'bg-white/5 text-gray-300'
                  }`}
                >
                  {msg.content}
                </div>
              </div>
            ))}

            {isLoading && (
              <div className="flex justify-start">
                <div className="bg-white/5 text-gray-400 px-3 py-2 rounded-lg text-sm">
                  Thinking...
                </div>
              </div>
            )}

            {error && (
              <div className="text-red-400 text-xs text-center">
                {error}{' '}
                <a href="#contact" onClick={() => setIsOpen(false)} className="underline">
                  Leave a message instead
                </a>
              </div>
            )}

            <div ref={messagesEndRef} />
          </div>

          {/* Input */}
          <form onSubmit={sendMessage} className="p-3 border-t border-white/10 flex gap-2">
            <input
              ref={inputRef}
              type="text"
              value={input}
              onChange={(e) => setInput(e.target.value)}
              placeholder="Ask something..."
              className="flex-1 bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm text-white placeholder-gray-500 outline-none focus:border-[#6366f1]/50"
              disabled={isLoading}
            />
            <button
              type="submit"
              disabled={isLoading || !input.trim()}
              className="px-3 py-2 bg-[#6366f1] rounded-lg text-white text-sm disabled:opacity-40 hover:bg-[#6366f1]/90 transition-colors"
            >
              →
            </button>
          </form>
        </div>
      )}
    </>
  );
}
```

- [ ] **Step 2: Add ChatWidget to index page**

Update `frontend/src/pages/index.astro` — add the import and component before closing `</BaseLayout>`:

Add to imports:
```astro
import ChatWidget from '../islands/ChatWidget.tsx';
```

Add before `</BaseLayout>`:
```astro
  <ChatWidget client:load />
```

- [ ] **Step 3: Verify chat widget renders**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/frontend
npm run dev
```

Open `http://localhost:4321`. Floating chat button should appear bottom-right. Click to open panel. Type a message — it will fail to connect (API not running), which is expected. Verify the error fallback message appears with "Leave a message instead" link.

- [ ] **Step 4: Verify hero CTA opens chat**

Click "Ask my AI →" in the hero section. Chat panel should open.

- [ ] **Step 5: Commit**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio
git add frontend/src/islands/ChatWidget.tsx frontend/src/pages/index.astro
git commit -m "feat: add AI chat widget React island"
```

---

## Task 16: Dockerfile for Go API

**Files:**
- Create: `api/Dockerfile`
- Create: `api/.dockerignore`

- [ ] **Step 1: Create `api/Dockerfile`**

```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Run stage
FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

- [ ] **Step 2: Create `api/.dockerignore`**

```
.env
.env.*
tmp/
*.test
```

- [ ] **Step 3: Verify Docker build**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/api
docker build -t portfolio-api .
```

Expected: Image builds successfully.

- [ ] **Step 4: Verify Docker container runs**

```bash
docker run --rm -p 8080:8080 portfolio-api &
sleep 2
curl -s http://localhost:8080/health
docker stop $(docker ps -q --filter ancestor=portfolio-api)
```

Expected: `{"status":"ok"}`

- [ ] **Step 5: Commit**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio
git add api/Dockerfile api/.dockerignore
git commit -m "feat: add Dockerfile for Go API"
```

---

## Task 17: CI/CD — Pull Request Validation

**Files:**
- Create: `.github/workflows/ci.yml`

- [ ] **Step 1: Create `.github/workflows/ci.yml`**

```yaml
name: CI

on:
  pull_request:
    branches: [main]

jobs:
  frontend-check:
    name: Frontend
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: frontend
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-node@v4
        with:
          node-version: 20
          cache: npm
          cache-dependency-path: frontend/package-lock.json

      - run: npm ci
      - run: npx astro check
      - run: npm run build

  api-check:
    name: API
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: api
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache-dependency-path: api/go.sum

      - run: go vet ./...
      - run: go test ./... -v
      - run: go build ./cmd/server
```

- [ ] **Step 2: Commit**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio
mkdir -p .github/workflows
git add .github/workflows/ci.yml
git commit -m "ci: add PR validation workflow for frontend and API"
```

---

## Task 18: CI/CD — Frontend Production Deploy

**Files:**
- Create: `.github/workflows/deploy-frontend.yml`

- [ ] **Step 1: Create `.github/workflows/deploy-frontend.yml`**

```yaml
name: Deploy Frontend

on:
  push:
    branches: [main]
    paths:
      - 'frontend/**'

jobs:
  deploy:
    name: Deploy to Netlify
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: frontend
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-node@v4
        with:
          node-version: 20
          cache: npm
          cache-dependency-path: frontend/package-lock.json

      - run: npm ci

      - name: Build
        run: npm run build
        env:
          PUBLIC_API_URL: ${{ secrets.PUBLIC_API_URL }}

      - name: Deploy to Netlify
        uses: nwtgck/actions-netlify@v3
        with:
          publish-dir: frontend/dist
          production-deploy: true
        env:
          NETLIFY_AUTH_TOKEN: ${{ secrets.NETLIFY_AUTH_TOKEN }}
          NETLIFY_SITE_ID: ${{ secrets.NETLIFY_SITE_ID }}
```

- [ ] **Step 2: Commit**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio
git add .github/workflows/deploy-frontend.yml
git commit -m "ci: add frontend production deploy to Netlify"
```

---

## Task 19: CI/CD — API Production Deploy

**Files:**
- Create: `.github/workflows/deploy-api.yml`

- [ ] **Step 1: Create `.github/workflows/deploy-api.yml`**

```yaml
name: Deploy API

on:
  push:
    branches: [main]
    paths:
      - 'api/**'

jobs:
  deploy:
    name: Deploy to Fly.io
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: api
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache-dependency-path: api/go.sum

      - name: Run tests
        run: go test ./... -v

      - name: Set up Fly CLI
        uses: superfly/flyctl-actions/setup-flyctl@master

      - name: Deploy to Fly.io
        run: flyctl deploy --remote-only
        env:
          FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
```

- [ ] **Step 2: Commit**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio
git add .github/workflows/deploy-api.yml
git commit -m "ci: add API production deploy to Fly.io"
```

---

## Task 20: Final Integration & Verification

**Files:**
- Modify: `.gitignore` (verify .env files are excluded)
- Verify: full build pipeline works end-to-end

- [ ] **Step 1: Verify no secrets in tracked files**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio
git grep -i "api.key\|apikey\|api_key\|secret\|password\|token" -- ':!*.example' ':!*.md' ':!*.yml'
```

Expected: No matches. If any match, remove the secret and add the file to `.gitignore`.

- [ ] **Step 2: Verify .env files are gitignored**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio
echo "test" > frontend/.env
echo "test" > api/.env
git status
```

Expected: Neither `.env` file appears in git status (both ignored).

- [ ] **Step 3: Clean up test .env files**

```bash
rm -f frontend/.env api/.env
```

- [ ] **Step 4: Build frontend end-to-end**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/frontend
npm run build
```

Expected: Build completes, static output in `frontend/dist/`.

- [ ] **Step 5: Build API end-to-end**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/api
go test ./... -v
go build ./cmd/server
```

Expected: All tests pass, binary compiles.

- [ ] **Step 6: Test full stack locally**

Terminal 1 — start API:
```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/api
GEMINI_API_KEY=your-key-here PORT=8080 ALLOWED_ORIGINS=http://localhost:4321 go run ./cmd/server
```

Terminal 2 — start frontend:
```bash
cd /Users/rafaeldeabreugomes/projects/portfolio/frontend
PUBLIC_API_URL=http://localhost:8080 npm run dev
```

Open `http://localhost:4321`. Verify:
- All sections render correctly
- Nav smooth scrolling works
- Mobile hamburger menu works
- Project cards link to case study pages
- Chat widget opens and sends messages (if GEMINI_API_KEY is set)
- Chat widget shows error fallback gracefully (if key not set)
- Contact form is present with Netlify attributes
- Footer links work

- [ ] **Step 7: Final commit**

```bash
cd /Users/rafaeldeabreugomes/projects/portfolio
git add -A
git status
git commit -m "chore: final integration verification and cleanup"
```

---

## Secrets Reference

**These secrets must be configured — NEVER committed to the repository:**

| Secret | Where | Purpose |
|---|---|---|
| `GEMINI_API_KEY` | `api/.env` (local) / Fly.io env vars (prod) / GitHub Secrets (CI) | Google Gemini API authentication |
| `NETLIFY_AUTH_TOKEN` | GitHub Secrets | Netlify deploy from GitHub Actions |
| `NETLIFY_SITE_ID` | GitHub Secrets | Target Netlify site for deploy |
| `FLY_API_TOKEN` | GitHub Secrets | Fly.io deploy from GitHub Actions |
| `PUBLIC_API_URL` | `frontend/.env` (local) / GitHub Secrets (CI) | Go API URL for chat widget |

**Setup checklist (manual, after code is deployed):**
1. Create Netlify site → copy auth token and site ID → add to GitHub Secrets
2. Create Fly.io app → `flyctl auth token` → add to GitHub Secrets
3. Get Gemini API key at https://aistudio.google.com/apikey → add to Fly.io secrets (`flyctl secrets set GEMINI_API_KEY=...`) and GitHub Secrets
4. Set `PUBLIC_API_URL` in GitHub Secrets to your Fly.io app URL (e.g., `https://portfolio-api.fly.dev`)
5. Enable branch protection on `main` in GitHub repo settings
