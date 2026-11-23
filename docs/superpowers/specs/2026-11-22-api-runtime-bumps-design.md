# API Runtime Bumps — Design Spec

> The API image's runtime base (`alpine:3.19`) is past end of support, and the Go toolchain and `govulncheck` pin are a minor behind. Touches `api/Dockerfile`, `api/go.mod`, the Go lines of two workflows, `.github/dependabot.yml`, `CLAUDE.md` and docs. `api/**` changes, so **merging deploys the API to Fly.io**.

## 1. Problem

- `api/Dockerfile` runs on `alpine:3.19`. Alpine's release table lists 3.19 with EOL 2025-11-01, so it has had no security fixes for three weeks (today is 2026-11-22).
- `.github/dependabot.yml` has no `docker` ecosystem, which is how the image went stale unnoticed.
- The container runs as root. Port 8080 needs no privilege and the server writes no files (grep for `os.Create`, `WriteFile`, `OpenFile`, `MkdirAll`, `CreateTemp` in non-test code finds nothing).
- Go 1.27 is the current release; `go.mod`, the builder image and CI are on 1.26. `govulncheck` is pinned one release behind.

## 2. Versions and Sources (checked 2026-11-22)

| Item | Was | Now | Source |
|---|---|---|---|
| Alpine runtime | `alpine:3.19` (EOL 2025-11-01) | `alpine:3.24` (latest stable branch, EOL 2028-06-01; 3.24.4 is its current release) | https://alpinelinux.org/releases/ (read via `releases.json`); Docker Hub tag `alpine:3.24` updated 2026-11-08 |
| Go | `go 1.26.1`, builder `golang:1.26-alpine`, CI `1.26` | `go 1.27.4`, builder `golang:1.27-alpine`, CI `1.27` | https://go.dev/doc/devel/release (lists go1.27.4); `go.dev/dl/?mode=json` lists go1.27.4 as the latest stable (go1.26.11 is the latest 1.26); Docker Hub `golang:1.27-alpine` updated 2026-11-16 and reports `go1.27.4` |
| govulncheck | `v1.9.0` | `v1.10.0` (tagged 2026-11-17) | `proxy.golang.org/golang.org/x/vuln/@v/list` and `@v/v1.10.0.info`; `.github/dependabot.yml` validates against the SchemaStore `dependabot-2.0.json` schema (Ajv, draft-07) |

Alpine 3.20 (EOL 2026-04-01) and 3.21 (EOL 2026-11-01) are also past EOL; 3.22 and 3.23 are supported but older than 3.24, so 3.24 is the only choice that maximises remaining support.

## 3. Changes

| File | Change |
|---|---|
| `api/Dockerfile` | builder `golang:1.27-alpine`; runtime `alpine:3.24` (minor pinned, not `latest`); `ca-certificates` kept; `adduser -D -H -u 10001 app` and `USER app` before `CMD` |
| `api/go.mod` | `go 1.26.1` -> `go 1.27.4`; no `toolchain` line |
| `.github/workflows/ci.yml` | API job `go-version: '1.27'`; `govulncheck@v1.9.0` -> `@v1.10.0` (Go lines only: another PR edits this file's Node versions) |
| `.github/workflows/deploy-api.yml` | `go-version: '1.27'` |
| `.github/dependabot.yml` | new `docker` entry for `/api`: weekly, `chore(deps)` prefix, limit 5, group `docker-minor-patch` (minor and patch), with a comment |
| `CLAUDE.md` | Go requirement, runtime image and non-root user, `govulncheck` pin, Dependency updates bullet |

## 4. Decisions

- **Fixed uid 10001, no home directory.** `-D` (no password) `-H` (no home directory). A fixed high uid keeps the user stable across image rebuilds. The binary is built `0755` and owned by root, so `app` can read and execute it but not modify it; `WORKDIR /app` is root-owned, which is fine because nothing is written. `USER` comes after `COPY`, so the copy does not need `--chown`.
- **`go 1.27.4`, not `1.27`.** The existing directive is a full patch version (`1.26.1`), so the format stays consistent. The builder's `GOTOOLCHAIN` is `local` (checked in the image), so the build only works if `golang:1.27-alpine` is at or above 1.27.4; it is today (go1.27.4). If the image ever lags the directive the build fails loudly rather than downloading a toolchain.
- **The module has no dependencies**, so `go mod tidy` changes nothing beyond the directive and there is no `go.sum`.
- **Docker grouping mirrors gomod/npm** (`minor` and `patch`; majors arrive alone). Dependabot bumps the Dockerfile only; it does not bump `go.mod` or `go-version` in the workflows, so those must move by hand (stated in the config comment and in `CLAUDE.md`). A Docker bump merges into a Fly deploy.
- **`fly.toml` unchanged.** `internal_port = 8080` matches the unprivileged listener.

## 5. Verification

Run from `api/` after `make sync-content` with go1.27.4: `go vet ./...`, `go test ./...`, `go build ./cmd/server`, `go run golang.org/x/vuln/cmd/govulncheck@v1.10.0 ./...` ("No vulnerabilities found."). Docker Desktop was available locally: `docker build`, then `docker run -p 18080:8080`, `curl localhost:18080/health` -> 200 `{"status":"ok"}`, and `id` / `ps` inside the container show uid 10001 `app`, PID 1 `./server`. Raw output is in the PR body. The orchestrator verifies `/health` on Fly after merge, which is the first run on Fly's builder.

## 6. Risks and Open Questions

- **Deploy risk.** The deploy builds the Dockerfile on Fly's remote builder; the local build used the same Dockerfile and succeeded, but only a Fly deploy proves it there. Rollback is a revert of this PR (which redeploys) or redeploying the previous Fly release image.
- **Fly machine runs the process as a non-root user.** Nothing in the API needs root (no volumes, no file writes, port 8080). If a future feature writes to disk it must write somewhere `app` owns.
- **CI `go-version: '1.27'`** resolves to the latest 1.27 patch at run time, which can be newer than `go.mod`; that is allowed (`go.mod` states a minimum).
- **Go minors keep needing manual bumps** of `go.mod` and the workflows even with the Docker entry; Dependabot's gomod ecosystem only moves module requirements.
