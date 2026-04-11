---
title: "Loudness Analyzer"
description: "Automated EBU R128 loudness compliance system with Python analyzer, serverless API, and React dashboard."
image: "/images/placeholder.svg"
tags: ["Python", "React", "Supabase", "Netlify"]
category: "Full-Stack"
featured: false
order: 5
github: "https://github.com/rfabreu/ffmpeg-dialnorm-detector"
---

## Overview

An automated loudness analysis system for broadcast MPEG-TS streams. A Python analyzer runs concurrent ffmpeg processes measuring EBU R128 loudness per channel, classifies results against broadcast thresholds, and submits to a cloud database. A React dashboard displays a color-coded compliance matrix by date and channel for quick operational scanning.

## Problem

Broadcast operations teams need to monitor audio loudness across dozens of channels continuously. Manual spot-checking is slow and misses transient violations. The industry standard (EBU R128) defines strict thresholds, but existing tools don't integrate measurement with visualization in a single workflow.

## Approach

Split the system into three independent layers: a Python CLI analyzer that runs ffmpeg in parallel (ThreadPoolExecutor, batch size 10, 60s timeout per stream), serverless Netlify Functions as the API layer with Supabase PostgreSQL for storage, and a React frontend rendering the compliance matrix. The analyzer submits measurements via authenticated HTTP POST; the dashboard fetches and visualizes by date.

## Tech Stack

- **Python** — ffmpeg-python bindings, ThreadPoolExecutor for parallel stream analysis
- **ffmpeg** — EBU R128 loudness filter (`ebur128=peak=true:meter=9`) on multicast UDP streams
- **Netlify Functions** — Serverless Node.js API (submit, matrix, dates, streams endpoints)
- **Supabase PostgreSQL** — Stream deduplication via UPSERT, time-series measurements
- **React 19 + Vite** — Dashboard with color-coded compliance matrix (green/yellow/red)
- **Tailwind CSS 4 + Recharts** — Dark-theme visualization with charting ready for trends

## Outcome

A three-tier system that automates what was previously manual spot-checking. The color-coded matrix lets operators scan dozens of channels at a glance — green (acceptable), yellow (low), red (loud) — with drill-down to exact dB readings per time slot.
