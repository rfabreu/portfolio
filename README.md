# Rafael Abreu | Software Engineer Portfolio

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Frontend](https://img.shields.io/badge/Frontend-Astro_6-ff5d01)](https://astro.build)
[![API](https://img.shields.io/badge/API-Go-00ADD8)](https://go.dev)
[![Deploy](https://img.shields.io/badge/Deploy-Netlify-00C7B7)](https://rafaelabreu.netlify.app)

**Live:** [rafaelabreu.netlify.app](https://rafaelabreu.netlify.app)

A modern developer portfolio built as a monorepo with an Astro static frontend and a Go API backend powering an AI chatbot. Designed to demonstrate full-stack engineering range — from frontend craft to backend services to container orchestration.

## Architecture

```
portfolio/
├── frontend/          Astro 6 + React islands + Tailwind v4
├── api/               Go microservice (AI chatbot via Gemini API)
├── .github/workflows/ CI/CD (GitHub Actions)
└── docs/              Design specs and implementation plans
```

| Component | Tech | Deployment |
|-----------|------|------------|
| Frontend | Astro 6, React 19, Tailwind CSS v4, TypeScript | Netlify |
| API | Go, Google Gemini API | Fly.io (Docker) |
| CI/CD | GitHub Actions (path-filtered) | Auto-deploy on push to `main` |

## Features

- **Static-first frontend** — Astro ships zero JS by default; only interactive components (chat widget, mobile nav) hydrate as React islands
- **AI chatbot** — floating widget powered by a Go backend proxying Google Gemini, constrained to professional Q&A
- **Content collections** — projects defined as Markdown files with typed frontmatter; add a project by dropping a `.md` file
- **Project case studies** — individual detail pages generated from Markdown (overview, problem, approach, tech stack, outcome)
- **Design system** — dark base with indigo/purple/pink accents, bold typography, scroll reveal animations, gradient orbs
- **Responsive** — mobile-first with hamburger nav overlay
- **CI/CD pipeline** — PR validation (lint, test, build), auto-deploy frontend to Netlify, auto-deploy API to Fly.io

## Local Development

### Prerequisites

- Node.js >= 22.12.0
- Go >= 1.22
- A [Google Gemini API key](https://aistudio.google.com/apikey) (free tier)

### Frontend

```bash
cd frontend
cp .env.example .env       # Set PUBLIC_API_URL
npm install
npm run dev                 # http://localhost:4321
```

### API

```bash
cd api
cp .env.example .env       # Set GEMINI_API_KEY
go run ./cmd/server         # http://localhost:8080
```

### Both together

Run the API in one terminal and the frontend in another. The chat widget connects to `PUBLIC_API_URL`.

## Adding a Project

Create a Markdown file in `frontend/src/content/projects/`:

```yaml
---
title: "Project Name"
description: "One-line description"
image: "/images/project-cover.png"
tags: ["Go", "Docker", "React"]
category: "Full-Stack"
featured: true
order: 1
github: "https://github.com/rfabreu/project"
demo: "https://project.example.com"
---

## Overview
What it is...

## Problem
What it solves...

## Approach
How you built it...

## Tech Stack
- **Go** — reason for choice
...

## Outcome
Results and learnings...
```

Set `featured: true` for large cards on the main page, `false` for the compact grid.

## CI/CD

| Workflow | Trigger | Action |
|----------|---------|--------|
| `ci.yml` | PR to `main` | Lint, test, build (frontend + API) |
| `deploy-frontend.yml` | Push to `main` (frontend changes) | Build Astro, deploy to Netlify |
| `deploy-api.yml` | Push to `main` (API changes) | Test, build Docker, deploy to Fly.io |

### Required GitHub Secrets

| Secret | Purpose |
|--------|---------|
| `NETLIFY_AUTH_TOKEN` | Netlify deploy authentication |
| `NETLIFY_SITE_ID` | Target Netlify site |
| `FLY_API_TOKEN` | Fly.io deploy authentication |
| `PUBLIC_API_URL` | Go API URL injected at build time |
| `GEMINI_API_KEY` | Google Gemini API for chatbot |

## Branching Strategy

```
feature/*  ──PR──►  main  ──auto-deploy──►  production
```

No `develop` branch. PRs to `main` trigger CI checks; merges trigger production deploys. Preview URLs generated on PRs via Netlify.

## License

This project is licensed under the [MIT License](LICENSE).
