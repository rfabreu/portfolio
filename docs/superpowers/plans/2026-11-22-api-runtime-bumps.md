# API Runtime Bumps — Implementation Plan

> **For agentic workers:** Steps use checkbox (`- [x]`) syntax for tracking.

**Goal:** Leave end-of-life Alpine, bump Go and `govulncheck`, run the container non-root, and cover Docker in Dependabot.

**Spec:** `docs/superpowers/specs/2026-11-22-api-runtime-bumps-design.md`
**Branch:** `chore/api-runtime-bumps`

## Task 1: Edit

- [x] **Step 1:** `api/Dockerfile`: `golang:1.27-alpine`, `alpine:3.24`, non-root `app` user (uid 10001).
- [x] **Step 2:** `api/go.mod`: `go 1.27.4`; `go mod tidy` (no change beyond the directive, no `go.sum`).
- [x] **Step 3:** `ci.yml` (Go lines only) and `deploy-api.yml`: `go-version: '1.27'`; `govulncheck@v1.10.0`.
- [x] **Step 4:** `dependabot.yml`: `docker` ecosystem for `/api`.
- [x] **Step 5:** `CLAUDE.md`: Go requirement, runtime image, pin, Dependency updates.

## Task 2: Verify and open the PR

- [x] **Step 1:** `make sync-content`; `go vet`, `go test`, `go build`, `govulncheck@v1.10.0`.
- [x] **Step 2:** `docker build`, run, `curl /health`, confirm non-root.
- [x] **Step 3:** Adversarial review of the diff; apply findings.
- [x] **Step 4:** Push, open PR with real output in the body. Do not merge.
