# Design Spec: CI Supply-Chain Checks (Dependabot, govulncheck, npm audit)

File: docs/superpowers/specs/2026-11-17-ci-supply-chain-design.md
Related: none (new work; no tracking issue)
Created: 2026-11-17
Status: Ready for owner review. The picomatch finding is handled by a time-boxed allowlist (orchestrator decision, owner approves at merge); see "Decision: time-boxed allowlist".

---

## Context

The repo has no dependency automation and no vulnerability scanning in CI. Only GitHub's default security alerts run. `main` is production and every push deploys, so a vulnerable or compromised dependency should be caught on the PR, not after the deploy.

## Constraints

- Branch protection requires the job names `Frontend` and `API`. `Frontend E2E` is not required (CLAUDE.md, checked 2026-11-07), but its name is kept as well. No job may be renamed, split or removed. New checks are steps inside the existing jobs only.
- Nothing auto-merges. Every Dependabot PR runs CI and needs the owner's approval. No auto-merge configuration is added.
- Work follows spec → plan → feature branch → PR, adversarial review, green CI, owner approval (CLAUDE.md).
- Never push to `main`; feature branch → PR only.
- No failing check may be silenced to make this PR green (no `|| true`, `continue-on-error`, lowered threshold or ignore-list entry). A failure is reported with its evidence and the decision goes to the owner. This is a constraint set for this PR (2026-11-17), consistent with the "Never assume" working rule in CLAUDE.md.

## Scope

### 1. `.github/dependabot.yml` (new)

Three ecosystems, all `interval: weekly`:

| Ecosystem | Directory | Why this directory |
|-----------|-----------|--------------------|
| `npm` | `/frontend` | `package.json` and `package-lock.json` live there |
| `gomod` | `/api` | `go.mod` lives there |
| `github-actions` | `/` | Dependabot reads workflows from `/.github/workflows` when the directory is `/` |

Per ecosystem:

- One group (`npm-minor-patch`, `gomod-minor-patch`, `actions-minor-patch`) with `update-types: [minor, patch]`. All minor and patch bumps for an ecosystem arrive as one PR per week. Major bumps are not grouped, so each arrives as its own PR and gets its own review.
- `open-pull-requests-limit: 5`. The default is 5 for version updates; it is set explicitly so the limit is visible and easy to tune.
- `commit-message.prefix`: `chore(deps)` for npm and gomod, `chore(ci)` for actions, matching the repo's conventional-commit style (`fix(ci): ...`).
- No `auto-merge`, no `reviewers`, no `assignees`. The owner reviews from the normal PR list.

Security updates are separate from version updates and keep working through GitHub's alerts. Grouping applies to version updates.

### 2. `govulncheck` in the `API` job

New last step in `api-check`, after `make sync-content`, `go vet`, `go test` and `go build`:

```yaml
- name: Scan Go dependencies for known vulnerabilities
  run: go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...
```

- **Order.** It runs after `make sync-content` because `internal/chat/embedded/` is gitignored and `go:embed` fails to compile without it. govulncheck loads and type-checks the packages, so it needs the same precondition as `go vet`. It is last so that a failing scan never hides the vet, test or build result.
- **Pinned version.** `@v1.9.0`, not `@latest`. See "Open questions" for why it is v1.9.0 and not the newer v1.10.0.
- **Why `go run`, not `go install`.** No binary is kept and there is nothing to cache or add to `PATH`. `go run pkg@version` ignores the current module's `go.mod`, so govulncheck's own dependencies never enter `api/go.mod` or `api/go.sum`.
- **Exit behaviour.** govulncheck exits non-zero when it finds a vulnerability that the code actually reaches. It does not fail on vulnerabilities in modules that are required but whose vulnerable symbols are never called. That is the intended trade-off: a call-graph-aware gate produces far fewer false positives than a module-level audit.
- **Toolchain.** The job uses `actions/setup-go` with `go-version: '1.26'`, which resolves to the latest 1.26 patch. govulncheck also reports standard-library vulnerabilities fixed in later patches, so the scan stays current without a pin.

### 3. `npm audit` in the `Frontend` job

New last step in `frontend-check`:

```yaml
- name: Audit npm dependencies
  run: |
    npm audit --json > "$RUNNER_TEMP/npm-audit.json" || true
    node ../.github/scripts/npm-audit-gate.mjs "$RUNNER_TEMP/npm-audit.json"
```

- **Threshold.** The gate (`.github/scripts/npm-audit-gate.mjs`, plain Node, no dependencies) fails on any high or critical advisory that is not in its allowlist, and passes on low and moderate. It also fails if the audit report is missing, unreadable or carries an `error` (for example a registry outage).
- **Allowlist.** Each entry is keyed by GHSA id and must carry a reason and an expiry date. An expired entry fails the gate, so an exception cannot be forgotten. An entry whose advisory is no longer reported prints a notice to remove it.
- **Dev dependencies are included** (no `--omit=dev`). Justification:
  1. This is a static site. Almost nothing runs in production from `node_modules`; Astro, Vite, Tailwind and React are bundled at build time. The `dependencies` / `devDependencies` split in `package.json` therefore does not track what matters, which is what executes on the build machine.
  2. The packages that execute with the most trust are build and CI tooling: `astro`, `vite`, `esbuild`, `vitest`, `playwright`. `deploy-frontend.yml` runs `npm ci`, the build, Playwright and the PDF generator from the same dependency tree and then publishes `frontend/dist` to Netlify. A compromised build-time or dev dependency can therefore alter what ships to production, even though it is not a runtime dependency of the site.
  3. Measured today, the choice changes the count only slightly: 13 high findings with dev dependencies, 12 with `--omit=dev`. Both fail.
  4. The cost is more false positives from dev-only advisories (for example a test-runner ReDoS). That cost is accepted; if it becomes noise, the right response is a documented, time-boxed exception, not a weaker default.
- **Placement.** Last in the job, for the same reason as govulncheck.
- **Source of truth.** `npm audit` reads `package-lock.json`, which `npm ci` has already installed, so the audit covers exactly what CI built.

## Finding: npm audit fails today

Run locally on this branch's base (`aaebfa5`, equal to `origin/main`) with `npm ci && npm audit --audit-level=high` in `frontend/`:

- **Result: exit 1, 13 high-severity findings** (12 with `--omit=dev`).
- **Single root advisory:** GHSA-fc5q-39xc-vjpf / CVE-2026-112366, "picomatch vulnerable to ReDoS through single-branch repeated extglobs", CWE-1333, CVSS 7.5, published 2026-10-19.
- All other findings (`astro`, `vite`, `vitest`, `tinyglobby`, `fdir`, `unstorage`, `anymatch`, `@astrojs/*`) are transitive "depends on a vulnerable picomatch" entries. `npm audit --json` lists exactly one advisory source.
- **No fix exists upstream.** The advisory's vulnerable range is `<= 4.0.7` with `first_patched_version: null`. `4.0.7` is the newest picomatch on npm (`dist-tags.latest`), and `2.3.2`, used by `anymatch` under `unstorage`, is the newest 2.x. So the dependency-bump route (`npm update`, `overrides`) cannot resolve it. `npm audit fix --force` offers only breaking changes (it names `vitest@3.1.4` and `astro@5.3.0`, multi-major downgrades from `vitest` 4.x and `astro` 7.x). That is not a small bump, so it was not applied or evaluated here.
- **Practical exposure** (assessment, not verified exploitability): picomatch is used here by build-time file globbing in Astro, Vite and Vitest. The advisory describes a ReDoS triggered by crafted glob patterns, so exploitation would need an attacker-influenced pattern, and the shipped site is static output. This is an assessment from reading the advisory, not a tested conclusion, and it has not been checked against how each tool in the tree builds its patterns. It is the owner's call, not this PR's.

Consequence: a plain `npm audit --audit-level=high` step fails on every PR until upstream ships a patch, which would also hold back the Dependabot and govulncheck parts of this PR.

## Decision: time-boxed allowlist

Added in orchestrator review (2026-11-17), as option E under "Open questions"; the owner approves it at merge:

- `GHSA-fc5q-39xc-vjpf` is allowlisted in the gate script until **2027-01-31**, with the exposure reasoning above as its recorded reason.
- Every other high or critical advisory still fails `Frontend`, so the gate keeps its value for new findings.
- When picomatch ships a fix, Dependabot's grouped PR (or `npm update picomatch`) picks it up; the gate then prints a notice to remove the entry. If no fix exists by the expiry date, the gate fails and forces a re-review.
- The gate's logic is an exported `evaluate()` function covered by `.github/scripts/npm-audit-gate.test.mjs` (`node --test`, run in the `Frontend` step before the audit). Cases: allowlisted advisory, unlisted high/critical in any casing, advisory without a GHSA url, moderate-only with a stale-entry notice, expired entry, missing or malformed `reason`/`expires`, npm `error` payload, malformed reports, a dependency on a package missing from the report, and metadata totals with no advisory.
- The CLI was also run against the real report (passes, picomatch allowed) and against an adversarial reviewer's 12 crafted reports; all fail except a forged advisory URL, which would need npm itself to emit fabricated data and is not defended against.
- **Exposure evidence** (assessed, not tested): the only glob patterns in the repo are literals in `src/content.config.ts`; there is no `import.meta.glob` in `src`, and `astro.config.mjs` sets no glob options. picomatch is also used inside Astro's cache and remote-image pattern handling, both driven by config, not visitor input; those paths were not traced.
- **Expiry cost:** from 2027-02-01 every PR, including Dependabot's, fails `Frontend` until the entry is re-reviewed. There is no scheduled warning beforehand.

govulncheck, run locally from `api/` after `make sync-content`, with both v1.9.0 and v1.10.0: `No vulnerabilities found.` (exit 0).

## Verified Facts

| Fact | Source |
|------|--------|
| Required job names are `Frontend`, `Frontend E2E` and `API`; they are unchanged | `.github/workflows/ci.yml`, parsed and compared after the edit |
| `internal/chat/embedded/` is gitignored and created by `make sync-content` | CLAUDE.md, `api/Makefile` |
| govulncheck `v1.10.0` is the newest version on the Go module proxy (published 2026-11-17T20:25:28Z); `v1.9.0` was published 2026-10-11T03:39:53Z | `proxy.golang.org/golang.org/x/vuln/@v/<version>.info` |
| Both versions declare `go 1.26.0` | `proxy.golang.org/golang.org/x/vuln/@v/<version>.mod` |
| GitHub Releases for `golang/vuln` lag behind (latest release `v1.1.4`, 2025-01-13); tags and the module proxy are the real release signal | `gh api repos/golang/vuln/releases` and `/tags` |
| `dependabot.yml` is valid against the SchemaStore Dependabot 2.0 schema; `ci.yml` is valid against the SchemaStore GitHub Workflow schema | local validation with `jsonschema` |
| Neither deploy workflow triggers on `.github/**` (`deploy-frontend.yml`: `frontend/**`; `deploy-api.yml`: `api/**` and two content paths) | `.github/workflows/deploy-*.yml` |
| `npm audit` fails today on GHSA-fc5q-39xc-vjpf with no patched version | `npm audit`, `gh api advisories/GHSA-fc5q-39xc-vjpf`, `npm view picomatch dist-tags` |

Not asserted:

- That Dependabot accepts the config on GitHub. Only schema validity was checked. The first scheduled run, or the repo's Insights → Dependency graph → Dependabot page, is the real confirmation.
- That `govulncheck@v1.9.0` and `npm audit` behave identically on the CI runner (Node 22, Go 1.26.x) as locally (Node 26.13, npm 11, Go 1.27.4). The CI run on this PR is the evidence.
- Branch protection settings beyond what CLAUDE.md records.

## Deploy Impact

None. Only `.github/` and `docs/` change. No deploy workflow path filter matches, so merging deploys neither the frontend nor the API.

CI impact: the `API` job gains a govulncheck download, compile and run (about 13 seconds on a warm local machine; CI timing not yet measured). The `Frontend` job gains one registry request. Both depend on external services (`proxy.golang.org`, the npm registry and the Go vulnerability database). An outage there fails the check even though the code is fine.

## Acceptance Criteria

1. `.github/dependabot.yml` covers `npm` at `/frontend`, `gomod` at `/api` and `github-actions` at `/`, weekly, grouped minor/patch, limit 5, with no auto-merge.
2. `API` job runs govulncheck at a pinned version after `make sync-content`, and it passes.
3. `Frontend` job runs `npm audit` including dev dependencies, gated by `npm-audit-gate.mjs` (fails on unallowlisted high/critical advisories, expired entries and audit errors).
4. Job names `Frontend`, `Frontend E2E` and `API` are unchanged and no job was added or removed.
5. Findings are reported, not silenced: the npm audit failure is documented here and in the PR body.
6. Adversarial review covers required-check names, directories, unpinned tools, step order, and PR-noise and security problems.
7. CI: `Frontend` and `API` pass, with the gate log showing GHSA-fc5q-39xc-vjpf allowed until 2027-01-31. `Frontend E2E` passes once this branch includes the runner pin from PR #48.

## Files Changed (expected)

| File | Change |
|------|--------|
| `.github/dependabot.yml` | New: weekly grouped updates for npm, gomod, github-actions |
| `.github/workflows/ci.yml` | New last step in `Frontend` (npm audit + gate) and `API` (govulncheck) |
| `.github/scripts/npm-audit-gate.mjs` | New: allowlist-aware gate over `npm audit --json` |
| `.github/scripts/npm-audit-gate.test.mjs` | New: `node --test` unit tests for the gate |
| `docs/superpowers/specs/2026-11-17-ci-supply-chain-design.md` | New: this spec |
| `docs/superpowers/plans/2026-11-17-ci-supply-chain.md` | New: implementation plan |

Branch name: `chore/ci-supply-chain`

## Open questions

1. **How to handle the unfixable picomatch advisory. Resolved: option E** (see "Decision: time-boxed allowlist"); the owner confirms at merge. Options considered:
   - A (the original recommendation): keep the gate as is, hold this PR until a patched picomatch ships, then re-run CI. Dependabot's grouped PR or `npm update picomatch` will pick the patch up. Preserves the safety property. Cost: the PR waits for an upstream release of unknown date.
   - B: merge now with an admin override. Branch protection enforcement applies to non-admins only, so this is possible, but then every later PR, including Dependabot's, fails `Frontend` until the patch lands.
   - C: temporarily set `--audit-level=critical` in this PR. This weakens the gate. It is not done here because this PR's constraints rule out silencing a failing check; it is the owner's decision to take.
   - D: split into two PRs (Dependabot + govulncheck now, npm audit later). Not done: this work was scoped as a single PR.
   - **E (chosen in orchestrator review): allowlist only this advisory, with a reason and an expiry date**, in a small gate script. See "Decision: time-boxed allowlist". Option A remains available if the owner prefers to hold.
2. **govulncheck pinned to v1.9.0, not the newer v1.10.0.** v1.10.0 was published hours before this spec was written (2026-11-17). Pinning a brand-new release of a security tool defeats the point of a pin, so the previous release, five weeks old, was chosen. Both give the same result on this repo today. If the owner prefers current, change one version string. Dependabot does not update a version embedded in a `run:` command, so bumps are manual. An alternative that Dependabot's `gomod` ecosystem would track is a `tool` directive in `api/go.mod` (`go get -tool`), at the cost of adding govulncheck's dependency graph to `go.mod`/`go.sum`. Rejected for now to keep this PR out of the Go module graph.
3. **Dependabot config not run against GitHub.** See "Not asserted".
4. **`npm audit` depends on the npm registry's audit endpoint being reachable.** An outage fails the step. Accepted for now. If it proves flaky, a retry wrapper is a follow-up.
5. **No scheduled scan.** CI only runs on pull requests, so a new advisory shows up on the next PR, not on its own. GitHub's default security alerts and the weekly Dependabot PRs cover the gap. A scheduled workflow is possible but would add a new workflow and is out of scope.
6. **0.x minor bumps land in the grouped PR.** Under semver, a `0.y` → `0.(y+1)` bump is a minor change, so for example `@astrojs/check` (`^0.9.9`) moving to 0.10 would arrive inside the weekly `npm-minor-patch` PR even though 0.x minors can be breaking. Nothing auto-merges and CI runs on the group, so it is reviewed, but a breaking change could hide in a routine PR. Not changed here; a separate group or an `ignore` rule for 0.x packages is a follow-up if it becomes a problem.
7. **Out-of-scope hardening noticed, not changed:** `ci.yml` has no top-level `permissions:` block (default token scope), and actions are pinned by tag (`@v4`), not commit SHA. Both are possible follow-ups. The `github-actions` Dependabot entry will keep the tags current meanwhile.

## CLAUDE.md follow-ups (not edited in this PR)

- CI section: the CI line lists the jobs' checks; add the npm audit gate (`npm audit --json` checked by `.github/scripts/npm-audit-gate.mjs`, with its unit tests) to `Frontend` and `govulncheck` (pinned) to `API`.
- Add a "Dependency updates" note: `.github/dependabot.yml`, weekly, grouped, never auto-merged.
- Record the allowlisted picomatch advisory (GHSA-fc5q-39xc-vjpf, expires 2027-01-31) and where the allowlist lives (`.github/scripts/npm-audit-gate.mjs`).
- Add `.github/` to the Project Structure block.
