---
title: "WRAPP"
description: "Real-time weather radar dashboard with precipitation overlay, optimized for satellite operations."
image: "/images/placeholder.svg"
tags: ["JavaScript", "Leaflet.js", "GitHub Actions"]
category: "Frontend"
featured: false
order: 6
github: "https://github.com/rfabreu/wrapp"
---

## Overview

A professional weather radar visualization dashboard built for satellite operations at Nextologies. Displays real-time precipitation radar overlays, current conditions, wind patterns, and 5-day forecasts centered on the operations area. Designed for 24/7 workstation environments with a dark theme and responsive layout for field devices.

## Problem

Satellite operations teams need at-a-glance weather awareness without switching away from their primary tools. Commercial weather apps are cluttered with ads and irrelevant data. The team needed a focused, always-on radar display with professional meteorological color coding.

## Approach

Built as a zero-build-step static site using vanilla JavaScript and Leaflet.js. Integrates RainViewer API for radar imagery (5-minute refresh cycles with dBZ reflectivity scale) and OpenWeatherMap for current conditions and forecasts. GitHub Actions handles deployment to GitHub Pages with secure API key injection at build time.

## Tech Stack

- **Vanilla JavaScript (ES6+)** — No framework, no build step, instant deployment
- **Leaflet.js** — Interactive maps with CARTO Voyager base tiles and canvas rendering
- **RainViewer API** — Real-time precipitation radar with historical timeline navigation
- **OpenWeatherMap API** — Current conditions, wind, and 5-day forecast
- **GitHub Actions** — CI/CD with secrets injection for API keys, auto-deploy to Pages
- **CSS3** — Glassmorphism panels (backdrop-filter blur) for operational readability

## Outcome

A focused operational tool used in a professional satellite operations environment. Responsive design works on both desktop workstations and mobile field devices, with automatic radar refresh and graceful offline degradation.
