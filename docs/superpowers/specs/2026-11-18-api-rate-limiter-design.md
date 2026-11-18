# Design Spec: Go API Rate Limiter

File: docs/superpowers/specs/2026-11-18-api-rate-limiter-design.md
Related: `2026-11-17-api-hardening-design.md` (section 5 left the client key unchanged pending evidence)
Created: 2026-11-18
Status: Approved for implementation (decisions come from the task brief; uncertain items are under Open Questions)

---

## Context

`api/internal/middleware/ratelimit.go` (origin/main 1e9296d) has three problems:

1. **Counting bug.** A client's window only resets after 60 s with no request, and `lastSeen` is refreshed on every request, rejected ones included. Once over the limit, a client stays blocked for as long as it keeps sending at least one request per 60 s. A scaled-down copy of the code confirmed this.
2. **Probably one shared bucket.** The key is the host of `r.RemoteAddr`, and no proxy header is read. Behind Fly's `http_service` proxy, `RemoteAddr` is very likely the proxy, so all visitors share a bucket. With bug 1, one script sending a request every 59 s blocks chat and game for everyone. Not confirmed in production (Open Question 1).
3. **Wrong shape for the routes.** One limiter wraps `/health`, `/api/chat` and `/api/game` with a single default of 30/min. The game makes one request per AI move; chat spends the owner's Gemini key. Nothing caps total Gemini calls. An invalid `RATE_LIMIT_RPM` silently becomes 0 (`strconv.Atoi` error ignored), which blocks everyone after one request. The 429 body is JSON sent as `text/plain` with no `Retry-After`.

## Constraints

- Stdlib only, `internal/<domain>/` layout (CLAUDE.md).
- No frontend change unless needed. No change may reject legitimate `ChatWidget.tsx` or `GameBoard.tsx` traffic.
- No permanent IP logging (IPs are personal data).
- Work follows spec → plan → branch → PR with adversarial review. Never push to `main`. Merging deploys the API.

## Design

### 1. Client key (`middleware/clientkey.go`)

`clientKey(r, trustFly)`:

- With `trustFly` (env `TRUST_FLY_CLIENT_IP=1`, set in `fly.toml`): the `Fly-Client-IP` header is used when there is exactly one header line and `netip.ParseAddr` accepts it. `"1.2.3.4, 5.6.7.8"`, `"garbage"`, `""`, `"1.2.3.4:80"`, a CIDR and repeated header lines all fall back.
- Otherwise, or on fallback: the host of `RemoteAddr`.
- Normalisation: `Unmap()` (so `::ffff:1.2.3.4` equals `1.2.3.4`), drop the IPv6 zone, IPv6 becomes its `/64` prefix string (a host owns a whole /64, so per-address keys would let it rotate for free).
- `X-Forwarded-For`, `X-Real-IP` and `Forwarded` are never read. Fly documents `X-Forwarded-For` as spoofable.

Why trusting `Fly-Client-IP` is a flag and not automatic: https://fly.io/docs/networking/request-headers/ says request headers "are added by Fly Proxy" and defines `Fly-Client-IP` as "the IP address of the client from the perspective of Fly Proxy". It does not say whether the proxy overwrites a value the client sent (Open Question 2). The flag is off by default, so a local run or a non-Fly deploy never trusts a client-supplied header. The page also says the header is the better choice when no other reverse proxy sits in front of Fly. None does here: the browser calls the API origin directly.

### 2. Fixed window that resets on its own (`middleware/ratelimit.go`)

`RateLimit(RateLimitConfig{Limit, Window, MaxKeys, TrustFlyClientIP, Now})`; `Limit` is required and a value `<= 0` panics (programmer error; `main` validates first).

- Per key `{count, start}`. A window starts at the key's first request. The window is over when `now - start >= Window` (exact boundary: 1 ns before is still blocked), and the next request starts a new one.
- Over the limit (`count >= Limit`) means reject **without changing anything**, so retries never extend a block.
- A client can send up to 2×Limit requests across a window boundary. Accepted for abuse protection.
- `Now func() time.Time` is injectable. There is no background goroutine.

**Bounded map.** `MaxKeys` defaults to 10,000 (a few MB at most, per limiter). Entries sit in a queue ordered by window start, so the next to expire is always at the front: each request pops expired windows (amortised O(1), no background goroutine and no full scans). If a new key arrives while the map is still full of live windows, the client with the oldest window is dropped, which resets that client's count.

This fails open for one client instead of closed for everyone. The first draft rejected new keys with 429 while the map was full; adversarial review showed that an attacker who fills 10,000 keys and renews each once per window (about 167 requests/s) locks every new visitor out indefinitely, and a free IPv6 tunnel (a /48 holds 65,536 /64 prefixes) supplies the keys. Eviction costs the attacker nothing it did not already have: N keys already buy N×Limit requests. Memory stays bounded, and the daily Gemini cap and the concurrency limit are the backstops against a flood.

### 3. Per-route limits (`cmd/server/main.go`)

| Route | Limit (per client, per minute) | Env |
|-------|-------------------------------|-----|
| `/health` | none | |
| `POST /api/chat` | 10 per minute | `CHAT_RATE_LIMIT_RPM` |
| `POST /api/chat` | 50 per 24 h window | `CHAT_CLIENT_DAILY_LIMIT` |
| `POST /api/game` | 60 per minute | `GAME_RATE_LIMIT_RPM` |

The daily chat limit exists because of section 4: at 10 per minute a single client would spend the whole 500-call global budget in 50 minutes and turn chat off for everyone. A request must pass the per-minute limiter, then the daily one, then validation, then the global cap. Each limiter has its own map. A client that hits the daily limit gets a 429, which the widget words as "Rate limit reached. Try again in a minute."; that is inaccurate for a wait of hours, but only a client that sent 50 messages in a day sees it, and the "Leave a message instead" link sits beside it. The limiters wrap the route handlers inside the mux, so `/health`, CORS preflights (answered by the CORS middleware first) and 404/405 responses never touch them. The 429 still passes through CORS, so the widget can read its status.

Defaults: a visitor types a chat message every few seconds at best, so 10/min is generous and caps Gemini spend per client. A tic-tac-toe game makes at most 5 AI moves, so 60/min allows about a dozen games a minute. The old default was 30/min for everything. These numbers are judgement calls (Open Question 4).

`RATE_LIMIT_RPM` (the old variable) is not set in `fly.toml`, the workflows or any file in the repo; it appears only in `api/.env.example`. It may exist as a Fly secret (Open Question 5), so it stays supported as a deprecated fallback. Under the old code it was one bucket shared by every route and, behind Fly's proxy, probably every visitor, so its value says little about a sensible per-client limit. When set and valid, it is the default for the game route, and the chat default becomes the smaller of it and 10: it can tighten chat, never loosen it. The per-route variables override it. Startup logs a deprecation warning.

### 4. Global Gemini cap (`chat/cap.go`)

`CappedProvider` wraps the `Provider`: a process-wide fixed window (`chat.DailyCapWindow` = 24 h, starting at the first call), `CHAT_DAILY_CAP` calls (default 500). Rejects change nothing. Properties:

- It sits after per-client limiting and body validation, so 400s do not spend it. Failed Gemini calls do count (they may be billed or count against quota).
- The count is in memory and per machine: it resets when the machine stops (`auto_stop_machines`) or restarts, and each additional machine has its own budget. It is a cost brake, not an accounting system. With more than one Machine, each has its own budget, and the proxy may send one client to either, which also doubles that client's per-minute limits. The durable control is a Google Cloud budget or quota (Open Question 3).
- Once spent, `/api/chat` returns **503** `{"error":"chat is unavailable right now (daily limit reached)"}` with `Retry-After` (seconds until the window resets, up to 86400). 503 over 429 because the condition is not about this client's rate; the widget maps 429 to "Rate limit reached. Try again in a minute.", which would be wrong for a wait of hours. For any other non-OK status it shows "Something went wrong." followed by the "Leave a message instead" link to the contact form, which is the right fallback. **No frontend change is needed.**
- The first rejection in a window logs one `WARNING` line (no personal data) so the owner sees it in `fly logs`.

### 5. Config validation (`cmd/server/config.go`)

Invalid rate or cap values, including `CHAT_CLIENT_DAILY_LIMIT` (non-numeric, `<= 0`, overflow, whitespace) log a `WARNING` naming the variable and fall back to the default. Crash on startup was rejected: a bad setting would turn a deploy into a crash loop for a storefront, while every default is restrictive, so falling back is the safe direction. An invalid per-route value falls back to what that route would otherwise use (the legacy value if set, else the default). An unparseable `TRUST_FLY_CLIENT_IP` leaves trust off. Startup logs the effective limits.

### 6. Responses

429 (limiter) and 503 (cap) are written by `middleware.WriteJSONError`: `Content-Type: application/json`, `X-Content-Type-Options: nosniff`, `Retry-After` (whole seconds, rounded up, minimum 1), and `{"error":"..."}` like the rest of the API.

### 7. `fly.toml`

```toml
[env]
  TRUST_FLY_CLIENT_IP = '1'

[http_service.concurrency]
  type = 'requests'
  soft_limit = 20
  hard_limit = 25
```

Keys verified against https://fly.io/docs/reference/configuration/ (`[http_service.concurrency]` with `type`, `soft_limit`, `hard_limit`; `requests` is recommended for HTTP services). Per the docs, at `hard_limit` the proxy stops sending new traffic to that Machine; the docs do not say what happens to the excess when there is only one Machine (Open Question 6). `soft_limit` also feeds autostart/autostop. 25 concurrent requests is far above this site's traffic, so it bounds worst-case memory on the 256 MB VM (a chat request holds up to 1 MiB of body).

## Out of scope

Persistent or cross-machine counters, a CAPTCHA, auth, changes to CORS, chat prompt content and the frontend error copy.

## Open Questions

1. **What do `RemoteAddr` and `Fly-Client-IP` hold in production?** Expected: `RemoteAddr` is a Fly proxy address (so today everyone shares one bucket) and `Fly-Client-IP` is the visitor. The owner can check with a one-off, **temporary** line such as `log.Printf("remote=%s fly=%s", r.RemoteAddr, r.Header.Get("Fly-Client-IP"))` in the limiter on a throwaway branch, deployed, read with `fly logs`, and removed straight away. Do not keep IP logging: IPs are personal data.
2. **Does Fly Proxy overwrite a client-sent `Fly-Client-IP`?** The docs do not say. If it passes the value through, anyone can rotate it and defeat the per-client limit (the daily cap, the map bound and the concurrency limit remain). Test: send `curl -H 'Fly-Client-IP: 203.0.113.7' https://<app>/…` with the temporary log above and check the logged value is the real one. If it is forged-through, set `TRUST_FLY_CLIENT_IP` to `0` in `fly.toml`.
3. **Does the Gemini key have a Google Cloud quota or budget alert?** The in-memory cap resets on machine restart and per machine. A budget alert or a quota on the key is the durable control. Owner action.
4. **Are the defaults right?** Chat 10/min, game 60/min, daily cap 500 (500 calls of at most 1024 output tokens on `gemini-2.5-flash-lite` is a small spend; the owner can confirm against the price list). Tune with the env vars; no code change.
5. **Is `RATE_LIMIT_RPM` set as a Fly secret?** `fly secrets list` shows names only. If it is set, it now drives both defaults (with a deprecation warning at startup); the owner can unset it and use the new variables.
6. **How does Fly treat requests beyond `hard_limit` with one Machine?** Not stated in the docs. Check `fly status` for the Machine count. If requests are rejected rather than queued, a burst above 25 concurrent requests would see 503s from the proxy, not from the app.
7. **`hard_limit = 25` is a new failure mode.** Before, no hard limit was enforced. Errors from Fly Proxy carry no CORS headers, so the widget would show a generic error. Slow-body requests hold a slot for up to the 15 s `ReadTimeout`, so a few dozen slow connections could occupy all 25. Accepted as a memory bound on a 256 MB VM; raise `hard_limit` if it bites.
8. **Shared addresses.** Everyone behind one IPv4 address (office, carrier-grade NAT) shares one bucket: 10 chat messages per minute and 50 per day.
9. **Layering.** `chat` imports `middleware` for `WriteJSONError`, and the other chat/game errors still use `http.Error` (`text/plain`). A neutral `internal/httpjson` package used by every handler is a reasonable follow-up; it is out of scope here.
10. **`Retry-After` and CORS.** `Retry-After` is not in `Access-Control-Expose-Headers`, so browser code cannot read it. The widgets do not use it today.

## Acceptance Criteria

1. Same `RemoteAddr`, different `Fly-Client-IP` → separate buckets when trusted; a forged header is ignored when not trusted; invalid headers fall back; `X-Forwarded-For` is never used.
2. `2001:db8::1` and `2001:db8::2` share a bucket, `2001:db8:0:1::1` does not; `::ffff:1.2.3.4` equals `1.2.3.4`.
3. A window resets after `Window` with no traffic needed, and rejected requests never extend it.
4. The map never exceeds `MaxKeys`; expired entries go first, then the oldest window, and a flood of new keys is never locked out.
5. `/health` is never limited; chat and game limits are independent.
6. The daily cap stops calls to the provider at `CHAT_DAILY_CAP`, answers 503 JSON with `Retry-After`, and is not spent by invalid requests.
7. Invalid config falls back to defaults with warnings; the deprecated variable behaves as documented.
8. 429 and 503 are `application/json` with `Retry-After`.
9. `go vet`, `go test -race ./...`, `go build`, `govulncheck` pass; CI is green.

## Files Changed

| File | Change |
|------|--------|
| `api/internal/middleware/clientkey.go`, `clientkey_test.go` | New: key derivation and tests |
| `api/internal/middleware/ratelimit.go`, `ratelimit_test.go`, `limiter_internal_test.go` | Rewrite: fixed window, bounded FIFO map, injectable clock, `WriteJSONError` |
| `api/internal/chat/cap.go`, `cap_test.go` | New: global Gemini call cap |
| `api/internal/chat/handler.go` | Map `CapExceededError` to 503 |
| `api/cmd/server/config.go`, `config_test.go` | New: env parsing and validation |
| `api/cmd/server/main.go`, `handler_test.go` | `newHandler` with per-route limiters; integration tests with a fake provider |
| `api/fly.toml`, `api/.env.example` | Env and concurrency settings |
| `CLAUDE.md` | API env vars and limits |

Branch name: `fix/api-rate-limiter`
