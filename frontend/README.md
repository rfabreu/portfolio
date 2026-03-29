# Portfolio Frontend

Astro 6 static site with React islands for interactive components. Part of the [portfolio monorepo](../README.md).

## Setup

```bash
cp .env.example .env    # Configure API URL
npm install
npm run dev             # http://localhost:4321
```

## Commands

| Command | Action |
|---------|--------|
| `npm run dev` | Start dev server at `localhost:4321` |
| `npm run build` | Build production site to `./dist/` |
| `npm run preview` | Preview production build locally |
| `npx astro check` | Run TypeScript type checking |

## Project Structure

```
src/
├── components/     Static Astro components (Nav, Hero, About, etc.)
├── islands/        React components that hydrate client-side
│   ├── ChatWidget.tsx    AI chatbot floating widget
│   └── MobileNav.tsx     Mobile hamburger menu
├── content/        Markdown content collections
│   └── projects/   Project case studies (.md files)
├── layouts/        Page layouts (BaseLayout, ProjectLayout)
├── pages/          Routes
│   ├── index.astro           Main page
│   └── project/[slug].astro  Dynamic project pages
├── styles/         Tailwind v4 theme and global CSS
└── content.config.ts  Content collection schema
```

## Design System

Tailwind v4 CSS-first configuration in `src/styles/global.css` using `@theme` directives.

**Colors:** Dark base (`#050508`), indigo/purple/pink accents
**Typography:** System sans-serif for body, `Courier New` monospace for code/labels
**Animations:** `animate-fade-up`, `animate-fade-in`, `animate-stagger`, `.reveal` (scroll)

## Environment Variables

| Variable | Description |
|----------|-------------|
| `PUBLIC_API_URL` | Go API URL (e.g., `http://localhost:8080` local, `https://portfolio-api-rfabreu.fly.dev` production) |
