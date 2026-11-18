# Design Spec: Go API Hardening

File: docs/superpowers/specs/2026-11-17-api-hardening-design.md
Related: none (findings from a code read of origin/main aaebfa5)
Created: 2026-11-17
Status: Approved for implementation (all decisions come from the task brief; uncertain items are under Open Questions)

---

## Context

A code read of `api/` on origin/main (aaebfa5) found five hardening gaps in the Go API that serves `POST /api/chat` and `POST /api/game` from Fly.io. This spec covers fixing them in one PR on `fix/api-hardening`.

## Constraints

- Stdlib `net/http` only, no frameworks (CLAUDE.md).
- Out of scope: chatbot persona and prompt content, CORS, new endpoints.
- No change may reject legitimate `ChatWidget.tsx` or `GameBoard.tsx` traffic.
- Work follows spec → plan → branch → PR, with adversarial review, green CI and owner approval (CLAUDE.md). Never push to `main`.
- Merging deploys the API to Fly.io (`deploy-api.yml` triggers on `api/**`). The `ChatWidget.tsx` change also deploys the frontend to Netlify.

## Scope

### 1. Gemini API key leaks into logs

`provider.go` built the URL as `...:generateContent?key=<key>`. When `http.DefaultClient.Do` fails (DNS, refused connection, timeout), the returned `*url.Error` includes the full URL, and `handler.go` logs it with `log.Printf("Chat error: %v", err)`. A failing test reproduced the leak:

```
gemini request: Post "http://127.0.0.1:59428/models/gemini-2.5-flash-lite:generateContent?key=test-secret-key-do-not-log": dial tcp ...: connect: connection refused
```

Fix:
- Send the key in the `x-goog-api-key` header. Never put it in the URL.
- Add `BaseURL` and `HTTPClient` fields on `GeminiProvider` for testability. `NewGeminiProvider` sets them to the real endpoint (`https://generativelanguage.googleapis.com/v1beta`) and a client with a timeout. A zero-value provider falls back to the same defaults, so it can't hit a nil client.
- The default client refuses redirects (`CheckRedirect` returns `http.ErrUseLastResponse`). Go forwards custom headers such as `x-goog-api-key` to a redirect target, even on another host. Only `Authorization`, `Cookie` and a few others are stripped. With the key in a header, a redirect could otherwise forward it. A redirect is now reported as a non-200 error.
- Responses are read through `io.LimitReader` (1 MiB) instead of an unbounded `io.ReadAll`.
- Tests: a refused connection, an HTTP 500 and a timeout each return an error that doesn't contain the key. The request carries the key in the header and not in the query. A 307 to another server doesn't forward the key.

Header support evidence:
- Google documents `x-goog-api-key` as the way to pass an API key on REST calls: https://docs.cloud.google.com/docs/authentication/api-keys-use
- Gemini's own REST examples use it on the same host (`generativelanguage.googleapis.com/v1beta/interactions`): https://aistudio.google.com/docs/api-key
- An empirical check against `.../models/gemini-2.5-flash-lite:generateContent` (2026-11-17): with no key, Google returns `403 ACCESS_DENIED`. With a fake key in `x-goog-api-key`, it returns `400 API_KEY_INVALID`. So this endpoint reads the header. It hasn't been tested with a real key (see Open Questions).

### 2. No timeout on outbound calls

`http.DefaultClient` has no timeout, so a hung Gemini call held the handler goroutine until the client disconnected.

Fix: `chat.RequestTimeout = 30 * time.Second`, used as the default `http.Client{Timeout: ...}`. The test asserts the default equals `RequestTimeout` and that a hanging server produces a `net.Error` with `Timeout() == true`. The request still uses `r.Context()`, so a client disconnect also cancels the call. 30s is well above typical `gemini-2.5-flash-lite` latency for a reply capped at 1024 output tokens, and the prompt asks for under 200 words.

### 3. No input size limits on POST /api/chat and POST /api/game

| Limit | Value | Response |
|-------|-------|----------|
| Chat body (`http.MaxBytesReader`) | `MaxBodyBytes = 1 MiB` | 413 `{"error":"request body too large"}` |
| Chat message role | `user` or `assistant` only | 400 `{"error":"invalid message role"}` |
| Chat `user` message | `MaxUserContentRunes = 2000` runes | 400 `{"error":"message too long"}` |
| Chat `assistant` message | `MaxAssistantContentRunes = 8000` runes | 400 `{"error":"message too long"}` |
| Chat message count | 20 (unchanged), now enforced while decoding | 400 (unchanged message) |
| Game body (`http.MaxBytesReader`) | `MaxBodyBytes = 4 KiB` | 413 `{"error":"request body too large"}` |

Errors use the existing ``http.Error(w, `{"error":"..."}`, status)`` style.

Why these values:
- **User 2000.** Before this change, the widget input had no `maxLength`, so any server cap could reject a long paste with a generic "Something went wrong." This PR adds `maxLength={2000}` to the widget input (`MAX_INPUT_LENGTH`). The server counts runes, and `maxLength` counts UTF-16 code units. A string never has more runes than UTF-16 code units, so anything the input accepts passes the server cap. A Vitest test (`frontend/tests/chat-input-limit.test.ts`) fails if the two constants drift apart.
- **Assistant 8000.** The widget sends Gemini's replies back as `assistant` messages. Gemini output is capped at 1024 tokens (`maxOutputTokens`), and 8000 runes allows about 7.8 characters per token, well above typical tokenisation. Rejecting a genuine reply would break the conversation, so this cap is deliberately generous.
- **Body 1 MiB.** 20 messages × 8000 runes × 6 bytes (the worst-case JSON escape `\u00XX` that `JSON.stringify` emits for control characters) ≈ 960 KB, plus structure, which is under 1,048,576. So the body cap never rejects a request that passes the per-message caps. `TestChatHandler_AcceptsLargestWidgetConversation` checks this with 4-byte emoji user messages and fully escaped assistant messages.
- **Game 4 KiB.** `GameBoard.tsx` sends a 9-cell board plus three session counters (about 100 bytes).

`json.Decoder` only reads the first JSON value, so data after it is never decoded. The body cap still bounds anything read.

The message count is enforced while decoding. `ChatRequest.Messages` is an unexported `messageList` type whose `UnmarshalJSON` decodes one element at a time and stops at element 21. Before this, a 1 MiB body of `{},{},...` decoded about 349k structs before the count check. A measurement showed 55 MiB allocated for one request on a 256 MB VM. After the change it's 2 MiB. `TestChatHandler_RejectsManyMessagesCheaply` asserts under 16 MiB. A literal `[]chat.Message{...}` is still assignable to the field.

### 4. No server timeouts

`http.ListenAndServe` has no timeouts (slowloris exposure). Fix: `newServer(addr, handler)` returns an `http.Server` with:

| Field | Value | Reason |
|-------|-------|--------|
| `ReadHeaderTimeout` | 5s | Slow-header protection |
| `ReadTimeout` | 15s | Enough to upload a 1 MiB body |
| `WriteTimeout` | `ReadTimeout + chat.RequestTimeout + 15s` (60s) | The write deadline starts once headers are read, so it covers the body read plus the Gemini call. This way the Gemini error still reaches even a slow client |
| `IdleTimeout` | 120s | Keep-alive reuse without holding idle connections forever |

`WriteTimeout` is derived from `ReadTimeout` and `chat.RequestTimeout`. `TestNewServer_SetsTimeouts` asserts `WriteTimeout > ReadTimeout + RequestTimeout`. The request context isn't cancelled by `ReadTimeout`: the read deadline covers headers and body only (Go issue https://github.com/golang/go/issues/70834).

### 5. Rate limiter client IP behind Fly.io (verified, not changed)

`middleware/ratelimit.go` keys on `r.RemoteAddr`. Research:

- https://docs.fly.io/networking/request-headers lists `Fly-Client-IP` as a header "added by Fly Proxy": "The IP address of the client from the perspective of Fly Proxy."
- The same page says `X-Forwarded-For` "must be treated with caution to avoid spoofing attempts."
- The page does **not** say whether Fly Proxy overwrites a `Fly-Client-IP` header sent by the client.
- https://docs.fly.io/machines/runtime-environment documents `FLY_APP_NAME` and `FLY_MACHINE_ID` as set inside Machines, so the app could tell it's on Fly.

Decision: **no behaviour change.** If the proxy passes a client-supplied `Fly-Client-IP` through, trusting it would let anyone bypass the rate limit by rotating the header. Today's limiter can't be bypassed that way. It may instead be stricter than intended (see Open Questions). The brief says to change behaviour only if the docs confirm it, and they don't confirm the spoofing guarantee.

## Open Questions

1. **Does Fly Proxy overwrite a client-supplied `Fly-Client-IP`?** It isn't stated in https://docs.fly.io/networking/request-headers. Needs confirmation from Fly support or an empirical test on a deployed Machine (for example a temporary debug log of `Fly-Client-IP` and `RemoteAddr` for a request sent with a forged header). If confirmed, trust `Fly-Client-IP` only when `FLY_APP_NAME` is set.
2. **What does `r.RemoteAddr` hold on Fly?** It isn't stated in the docs read. If it's the proxy's address, all visitors share one 30 rpm bucket today. That's restrictive, not exploitable.
3. **Real-key confirmation of header auth.** `x-goog-api-key` is confirmed against `generateContent` with a fake key only (Google returns `API_KEY_INVALID`, not "missing key"). After deploy, the owner should send one chat message in production to confirm.
4. **Long pastes are now silently truncated at 2000 characters** by the widget's `maxLength`. There's no counter or hint. This is intended, and a visible counter could be a later UX tweak.
5. **Widget conversation length (pre-existing, out of scope).** `ChatWidget.tsx` has no 20-message cap, so the 11th user message gets the server's 400 and shows "Something went wrong." This PR doesn't change that.

## Acceptance Criteria

1. No error returned by `GeminiProvider` contains the API key (tests: refused connection, HTTP 500, timeout).
2. The key is sent only in `x-goog-api-key`. The request URL has no query string. Redirects are not followed.
3. Gemini calls time out after `chat.RequestTimeout`.
4. Chat: more than 20 messages → 400 without decoding them all; body > 1 MiB → 413; role other than `user`/`assistant` → 400; user content > 2000 runes or assistant > 8000 runes → 400; the largest widget conversation → 200.
5. Game: body > 4 KiB → 413; a normal move → 200.
6. The server sets all four timeouts, and `WriteTimeout > chat.RequestTimeout`.
7. The widget input has `maxLength` equal to the server user cap, guarded by a Vitest test.
8. `go vet`, `go test`, `go build`, `astro check`, `npm run build` and `npm test` all pass. CI is green.

## Files Changed (expected)

| File | Change |
|------|--------|
| `api/internal/chat/provider.go` | Key in header; `BaseURL`, `HTTPClient` fields with defaults; `RequestTimeout`; no redirects; bounded response read |
| `api/internal/chat/provider_test.go` | Key-leak, header, redirect, success, zero-value and timeout tests |
| `api/internal/chat/handler.go` | Body cap, streaming message-count check, role validation, per-message caps |
| `api/internal/chat/handler_test.go` | Limit tests, including the largest widget conversation |
| `api/internal/game/handler.go` | Body cap |
| `api/internal/game/handler_test.go` | Oversized body test |
| `api/cmd/server/main.go` | `newServer` with timeouts |
| `api/cmd/server/main_test.go` | Timeout test (new) |
| `frontend/src/islands/ChatWidget.tsx` | `maxLength={MAX_INPUT_LENGTH}` (2000) |
| `frontend/tests/chat-input-limit.test.ts` | Client/server cap parity test (new) |

Branch name: `fix/api-hardening`
