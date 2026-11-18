# API Rate Limiter — Implementation Plan

> **For agentic workers:** Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the Go API's rate limiter with a correct per-client, per-route fixed-window limiter, add a global Gemini call cap, and validate config.

**Architecture:** `middleware.RateLimit(RateLimitConfig)` per route in `newHandler`; `chat.CappedProvider` decorates the Gemini provider; `loadConfig` parses env with fallbacks.

**Spec:** `docs/superpowers/specs/2026-11-18-api-rate-limiter-design.md`
**Branch:** `fix/api-rate-limiter`

---

## Task 1: Client key

- [x] **Step 1:** Tests first (`clientkey_test.go`): trust on/off, invalid headers, `X-Forwarded-For`, IPv6 /64, IPv4-mapped, zone, no-port `RemoteAddr`. Confirm they fail (undefined `clientKey`).
- [x] **Step 2:** Implement `clientkey.go`.

## Task 2: Limiter

- [x] **Step 1:** Tests first (`ratelimit_test.go`, fake clock): limit, reset, rejects don't extend, exact boundary, JSON 429 with `Retry-After`, map cap and expiry, defaults, concurrency.
- [x] **Step 2:** Rewrite `ratelimit.go` (fixed window, bounded FIFO map, `WriteJSONError`).

## Task 3: Gemini cap

- [x] **Step 1:** Tests first (`cap_test.go`): trips at cap, reset, rejects don't extend, failed calls count, concurrency, handler 503.
- [x] **Step 2:** Implement `cap.go`; map `CapExceededError` in `handler.go`.

## Task 4: Config and wiring

- [x] **Step 1:** Tests first (`config_test.go`, `handler_test.go`): defaults, invalid values, legacy variable, `/health` unlimited, independent routes, cap, preflight, CORS on 429.
- [x] **Step 2:** `config.go`, `newHandler` in `main.go`.
- [x] **Step 3:** `fly.toml`, `.env.example`, `CLAUDE.md`.

## Task 5: Verify and open the PR

- [x] **Step 1:** `go vet ./... && go test -race ./... -v && go build ./cmd/server`; `go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...`.
- [x] **Step 2:** Run the server and exercise limits with curl (health, game, chat before the provider call, daily cap, `Retry-After`).
- [x] **Step 3:** Adversarial review (Opus subagent) of the diff against `origin/main`. Applied: FIFO eviction instead of rejecting new keys when full (lockout, scan cost); per-client daily chat limit; legacy `RATE_LIMIT_RPM` can only tighten chat; `NewCappedProvider` validates arguments; flaky test fixed. Documented as open: `Fly-Client-IP` trust, machine count, `hard_limit`, layering.
- [ ] **Step 4:** Push the branch, open the PR (deploy impact and rollback in the body), check `gh pr checks`.
