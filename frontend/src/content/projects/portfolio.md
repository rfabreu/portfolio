---
title: "Developer Portfolio"
description: "Full-stack monorepo with Astro frontend, Go API, AI chatbot, WebGL particle hero, and dark/light theming."
image: "/images/portfolio_website.png"
tags: ["Astro", "React", "Go", "Tailwind CSS"]
category: "Full-Stack"
featured: true
order: 1
github: "https://github.com/rfabreu/portfolio"
demo: "https://rafaelabreu.dev"
---

## Overview

This portfolio is itself a full-stack engineering project — an Astro 6 static site with React islands, a Go API server powering an AI chatbot and game engine, WebGL particle animations, and a dark/light theme system. Built as a monorepo with separate CI/CD pipelines for frontend (Netlify) and API (Fly.io).

## Problem

Most developer portfolios are static pages that don't demonstrate engineering depth. The goal was to build a portfolio that _is_ the portfolio piece — showcasing frontend craft, backend engineering, AI integration, and infrastructure in a single, cohesive project.

## Approach

Chose Astro for its islands architecture — static HTML by default, with React hydration only where interactivity is needed (chat widget, game board, particle canvas, theme toggle). The Go API handles AI chat via Google Gemini and a Tic-Tac-Toe engine with minimax/alpha-beta pruning. Tailwind CSS v4 with custom theme tokens powers a dark/light mode system that reaches every component.

## Tech Stack

- **Astro 6** — Static site generation with islands architecture
- **React** — Interactive islands (ChatWidget, GameBoard, ParticleHero, ThemeToggle)
- **Go** — API server with stdlib `net/http`, no frameworks
- **WebGL** — 60fps particle mesh animation with cursor interaction
- **Tailwind CSS v4** — Custom theme tokens with dark/light palette switching
- **Google Gemini** — AI chatbot backend via Go API
- **Netlify + Fly.io** — Separate deploy pipelines triggered by path-filtered GitHub Actions

## Outcome

A portfolio that demonstrates full-stack capability across frontend, backend, AI, animation, and DevOps — not just describes it. Every visitor interaction (chatting with the AI, playing Tic-Tac-Toe, toggling themes) is backed by real engineering.
