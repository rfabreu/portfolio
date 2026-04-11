---
title: "PRJTSKMNGR"
description: "Cross-department project coordination platform with role-based workflows, approval chains, and audit trails."
image: "/images/placeholder.svg"
tags: ["Go", "Next.js", "PostgreSQL", "Docker"]
category: "Full-Stack"
featured: true
order: 2
github: "https://github.com/rfabreu/PRJTSKMNGR"
---

## Overview

A workforce coordination platform designed for operations teams managing projects across multiple departments. Features hierarchical role-based access (Admin → Manager → Lead → Member), multi-level approval workflows, task dependencies with cycle detection, and an immutable audit trail for compliance.

## Problem

Generic project management tools lack the hierarchical approval workflows and department-scoped visibility that operations teams need. Tasks that cross department boundaries require coordination patterns that flat task boards can't express.

## Approach

Built a layered Go backend (handler → service → repository) with PostgreSQL using the LTREE extension for efficient hierarchical queries on departments and organizations. The Next.js frontend adapts its dashboard, navigation, and data access based on the authenticated user's role. Every action — creation, assignment, approval, completion — is logged to an append-only activity table.

## Tech Stack

- **Go 1.22** — Layered backend with 11 HTTP handlers, JWT auth (httpOnly cookies)
- **PostgreSQL 16** — LTREE extension for hierarchical departments, 12 repository modules
- **Next.js 16** — App Router with role-adaptive dashboard and admin views
- **React 19 + Tailwind CSS 4** — UI with Radix primitives and Recharts for metrics
- **Docker Compose** — Multi-stage builds (Go alpine + Node alpine) for local development

## Outcome

A platform that models real organizational complexity — scoped visibility, approval chains, cross-team dependencies, and performance metrics — built with a modern Go + Next.js stack. Currently in active development.
