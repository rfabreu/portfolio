# Plan: Go API Hardening

File: docs/superpowers/plans/2026-11-17-api-hardening.md
Related spec: docs/superpowers/specs/2026-11-17-api-hardening-design.md
Date: 2026-11-17
Branch: fix/api-hardening

---

## Summary of Spec

Stop the Gemini API key from reaching logs (header instead of URL), add an outbound timeout, cap request sizes and validate chat roles, give the HTTP server timeouts, and check (but not change) how the rate limiter sees client IPs on Fly.io.

## Governance Note

CLAUDE.md is the governing document for this work. Its rules on feature branches, never pushing to `main`, adversarial review before merge, and the CI checks apply. If this plan conflicts with CLAUDE.md, CLAUDE.md wins.

## Approach

Test first for every fix: write the test, watch it fail for the expected reason, then implement.

### Step 0: Setup

```
git checkout -b fix/api-hardening origin/main
cd api && make sync-content && go test ./...
```

### Step 1: Key out of the URL (`internal/chat/provider.go`)

1. Add tests in `provider_test.go`. Each points `BaseURL` at an `httptest` server:
   - refused connection (server closed): `err` doesn't contain the key
   - HTTP 500: `err` doesn't contain the key; the server saw the key in `x-goog-api-key` and not in the query
   - 200: the text is returned and the path is `/models/<model>:generateContent`
2. Run: compile failure (no `BaseURL`). Add `BaseURL`/`HTTPClient` while keeping `?key=`. Run again: the leak test fails with the key visible in the `*url.Error`.
3. Move the key to `req.Header.Set("x-goog-api-key", ...)` and drop the query. Run: green.

### Step 2: Outbound timeout

1. Test: `NewGeminiProvider(...).HTTPClient.Timeout > 0`, and a hanging server with a 50ms client returns an error without the key.
2. Add `RequestTimeout = 30 * time.Second` and use `&http.Client{Timeout: RequestTimeout}` in `NewGeminiProvider`.

### Step 3: Chat limits (`internal/chat/handler.go`)

1. Tests: oversized body → 413; roles `system`, `model`, `""`, `User` → 400; user > 2000 runes → 400; assistant > 8000 runes → 400; 2000 two-byte runes → 200; the worst-case 20-message widget conversation (emoji user messages, `\x01` assistant messages) fits `MaxBodyBytes` and returns 200.
2. Add the exported constants and confirm the tests fail on behaviour. Then add `http.MaxBytesReader`, `errors.As(err, &*http.MaxBytesError)` → 413, and a role/length loop that counts with `utf8.RuneCountInString`.

### Step 4: Game body cap (`internal/game/handler.go`)

1. Test: a body padded past `MaxBodyBytes` → 413. Confirm it fails without the cap (it returns 200).
2. Add `MaxBodyBytes = 4 << 10` and the same `MaxBytesReader`/413 handling.

### Step 5: Server timeouts (`cmd/server/main.go`)

1. Test `main_test.go`: `newServer` sets all four timeouts and `WriteTimeout > chat.RequestTimeout`.
2. Add `newServer(addr, handler) *http.Server` (5s / 15s / `chat.RequestTimeout + 15s` / 120s) and call `ListenAndServe` on it.

### Step 6: Widget parity (`frontend/src/islands/ChatWidget.tsx`)

1. Add `const MAX_INPUT_LENGTH = 2000;` and `maxLength={MAX_INPUT_LENGTH}` on the input.
2. Add `frontend/tests/chat-input-limit.test.ts`, which reads both files and asserts the constants match. Confirm it fails when the widget value is changed to 1999.

### Step 7: Rate limiter research (no code change)

Read https://docs.fly.io/networking/request-headers and https://docs.fly.io/machines/runtime-environment. The docs don't say whether `Fly-Client-IP` can be spoofed, so record that under Open Questions and leave `ratelimit.go` unchanged.

### Step 8: Verify

- `cd api && go vet ./... && go test ./... -v && go build ./cmd/server`
- `cd frontend && npx astro check && npm run build && npm test`
- Local smoke test (`PORT=18080 GEMINI_API_KEY=<fake marker> ./server`):
  - `GET /health` → 200
  - 2 MiB chat body → 413
  - 2001-char user message → 400
  - role `system` → 400
  - 5 KB game body → 413
  - normal game move → 200
  - chat with the fake key → clean 500; `grep -c <fake marker>` on the log → 0
  - chat with no key set → clean 500

### Step 9: Adversarial review

Run the `pr-review-toolkit:code-reviewer` subagent on `git diff origin/main...HEAD`. Ask it to try to bypass the limits, leak the key, break ChatWidget traffic, and find timeouts too short for Gemini, and to check CLAUDE.md conventions. Fix valid findings and re-verify.

### Step 10: Commit, push the feature branch, open the PR

Commit message:
```
fix(api): keep Gemini key out of logs, add timeouts and request limits
```

The PR body should cover: spec and plan paths, changes, the Fly.io deploy note (merging deploys the API to Fly.io and the frontend to Netlify), Open Questions, the review summary, and CLAUDE.md follow-ups.

## Pre-commit Checklist (per CLAUDE.md)

- [x] Spec exists and is referenced in commit/PR
- [x] Plan exists and follows spec
- [x] `go vet`, `go test`, `go build` pass
- [x] `astro check`, `npm run build`, `npm test` pass
- [x] Local smoke test done; no key in logs
- [ ] Adversarial review complete
- [ ] CI green
