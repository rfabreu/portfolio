# Plan: CI Supply-Chain Checks (Dependabot, govulncheck, npm audit)

File: docs/superpowers/plans/2026-11-17-ci-supply-chain.md
Related spec: docs/superpowers/specs/2026-11-17-ci-supply-chain-design.md
Related issues: none
Date: 2026-11-17
Branch: chore/ci-supply-chain

---

## Summary of Spec

Add `.github/dependabot.yml` (weekly, grouped minor/patch, no auto-merge) for npm, gomod and github-actions. Add a pinned `govulncheck` step to the `API` job after `make sync-content`, and an npm audit gate (`npm audit --json` checked by `.github/scripts/npm-audit-gate.mjs`, dev dependencies included, time-boxed allowlist) to the `Frontend` job. Required job names stay unchanged.

## Governance Note

CLAUDE.md is the governing document for this work. Its rules on branch protection (required checks `Frontend` and `API`; `Frontend E2E` not required but kept), no direct push to `main`, adversarial review and owner approval apply. If any conflict arises between this plan and CLAUDE.md, CLAUDE.md takes precedence. CLAUDE.md itself is not edited in this PR; follow-ups are listed in the spec and PR body.

## Approach

### Step 1: Baseline — run both checks before changing anything

```
cd api && make sync-content && go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...
cd frontend && npm ci && npm audit --audit-level=high
```

Always `cd` with an absolute path or from the repo root; a shell that is still inside another directory will run `npm audit` without a lockfile and give a misleading error.

Expected and observed: govulncheck exits 0 (`No vulnerabilities found.`); `npm audit` exits 1 with 13 high findings, all from GHSA-fc5q-39xc-vjpf (picomatch, no patched version). Per spec, this is reported, not silenced and not fixable by a small bump.

### Step 2: Create `.github/dependabot.yml`

- `version: 2`.
- Three `updates` entries: `npm` at `/frontend`, `gomod` at `/api`, `github-actions` at `/`.
- Each: `schedule.interval: weekly`, `open-pull-requests-limit: 5`, a `commit-message.prefix`, and one group of `update-types: [minor, patch]`.
- A header comment states that nothing auto-merges and why.
- No `auto-merge`, `reviewers` or `assignees` keys.

### Step 3: Edit `.github/workflows/ci.yml` — steps only

- `frontend-check` (name `Frontend`): append `npm audit --json` piped through `.github/scripts/npm-audit-gate.mjs` as the last step (orchestrator review replaced the plain `npm audit --audit-level=high`; see the spec's "Decision: time-boxed allowlist").
- `.github/scripts/npm-audit-gate.mjs`: new gate; GHSA-fc5q-39xc-vjpf allowlisted until 2027-01-31. `.github/scripts/npm-audit-gate.test.mjs`: unit tests, run first in the same step with `node --test ../.github/scripts/npm-audit-gate.test.mjs` (Node 22 needs the file path, not the directory).
- `api-check` (name `API`): append `go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...` as the last step, after `make sync-content`, `go vet`, `go test` and `go build`.
- Do not touch `name:` fields, job keys, triggers, or `frontend-e2e`.

### Step 4: Validate

- Parse both YAML files. Validate `dependabot.yml` against the SchemaStore Dependabot 2.0 schema and `ci.yml` against the SchemaStore GitHub Workflow schema. PyYAML reads the bare key `on:` as boolean `True`; map it back to `"on"` before validating, otherwise the schema check gives a false failure.
- Programmatically compare the job names to `Frontend`, `Frontend E2E` and `API`.
- Use `actionlint` if installed (it was not available locally).
- Run the gate's unit tests under Node 22 (the CI runtime): `npx -y node@22 --test .github/scripts/npm-audit-gate.test.mjs`.

### Step 5: Spec, plan, and adversarial review

- Write this plan and the spec.
- Start a reviewer subagent on `git diff origin/main...HEAD`. The planned type `pr-review-toolkit:code-reviewer` was not installed in this environment; the review ran as a `general-purpose` agent with the same reviewer brief (the closest installed reviewer type, `feature-dev:code-reviewer`, has no shell access to re-run the checks). Ask it to hunt for: broken required-check names, wrong directories, unpinned tools, steps that run before content sync, PR-noise problems (grouping, limits) and security problems (permissions, auto-merge, secrets).
- Fix valid findings; record the review outcome in the PR body.

### Step 6: Verification before completion

Do not claim success from reasoning alone. Re-run, in this state of the tree, and quote the results:

- `git diff origin/main --stat` matches "Files Changed" in the spec (no stray files; `api/internal/chat/embedded/` stays untracked/gitignored).
- Step 1 commands, again.
- Step 4 validation, again.
- After the PR is open: `gh pr checks <n>`. Compare each result to the spec's acceptance criteria. `Frontend` should pass, with the gate log showing GHSA-fc5q-39xc-vjpf allowed until 2027-01-31; any other blocking advisory is a real finding.

### Step 7: Commit, push feature branch, open PR

Branch: `chore/ci-supply-chain`

Commit message:
```
chore(ci): add dependabot, govulncheck and npm audit
```

Commit contents:
- `.github/dependabot.yml` (new)
- `.github/workflows/ci.yml`
- `.github/scripts/npm-audit-gate.mjs` (new)
- `.github/scripts/npm-audit-gate.test.mjs` (new)
- `docs/superpowers/specs/2026-11-17-ci-supply-chain-design.md` (new)
- `docs/superpowers/plans/2026-11-17-ci-supply-chain.md` (new)

PR body should:
- Reference spec and plan paths.
- Lead with the npm audit finding (no patched picomatch) and the time-boxed allowlist chosen in orchestrator review, for the owner to confirm at merge.
- Summarise the adversarial review.
- List Open questions and "CLAUDE.md follow-ups".
- State that nothing auto-merges and no deploy is triggered by this change.

Never push to `main`; never merge.

## Pre-commit Checklist (per CLAUDE.md)

- [x] Spec exists and is referenced in commit/PR
- [x] Plan exists and follows spec
- [x] Required job names unchanged (`Frontend`, `Frontend E2E`, `API`)
- [x] govulncheck pinned, runs after `make sync-content`
- [x] YAML valid; schema-checked
- [x] Both checks run locally; failures reported, not silenced
- [ ] Adversarial review complete
- [ ] PR CI results compared to acceptance criteria
- [ ] PR title/body follow conventions in CLAUDE.md
