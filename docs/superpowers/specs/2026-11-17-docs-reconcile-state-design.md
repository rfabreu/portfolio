# Docs Reconcile: Branch Protection, Netlify Checks, P10/P11 State — Design Spec

> Brings `CLAUDE.md` and `docs/phase2-roadmap.md` in line with the repository and GitHub state read on 2026-11-17. Documentation only: no source, config, workflow, or content changes. No human was available for brainstorming; uncertainties are in section 5.

---

## 1. Overview & Goals

- **Accuracy:** several statements written on 2026-11-07 are now stale because PRs #41, #42, #43 and #44 merged afterwards.
- **No assumptions:** every statement added is backed by `git` or `gh` output read on 2026-11-17 (section 2).
- **Minimal edits:** existing structure and voice are kept; only the stale lines change.

---

## 2. Verified State (2026-11-17)

| Fact | Source |
|---|---|
| `main` is at `aaebfa5` (merge of PR #42). | `git log origin/main` |
| Required checks on `main`: `Frontend`, `API`, `Frontend E2E`; `strict: true`; `enforce_admins: false`. | `gh api repos/rfabreu/portfolio/branches/main/protection` |
| Issue #34 is closed (2026-11-07, completed). PR #44 (merged) added the root `netlify.toml`. | `gh issue view 34`, `gh pr view 44` |
| On PR #42, `netlify/rafaelswe/deploy-preview` is `pass`; the other three Netlify rows are `skipping`. `netlify.toml` sets `base = "frontend"`. No source states why the earlier previews failed. | `gh pr checks 42`, `git show origin/main:netlify.toml` |
| PR #41 changed `freckles-design.md` (OrangeCat credit), `portfolio.md` (`demo`) and `profile.md` (`links.*`). On `main`, `links.linkedin` is `https://linkedin.com/in/rafael-a-gomes` and both `links.portfolio` and `demo` are `https://rafaelabreu.netlify.app`. | `gh pr view 41`, `git show origin/main:<file>` |
| PR #42 added `Services.astro` (mounted after Hero in `index.astro`, Services nav item), Contact fields `project_type` and `timeline` (no budget, no company field), a project-page CTA block, and set `profile.md` `title` to "Co-Founder, OrangeCat Labs \| Full-Stack & Infrastructure Leader". The `summary` is unchanged. | `gh pr view 42`, `git grep` |
| Issue #30 text says the Services section "was dropped from the plan". It is open and was last updated 2026-11-07. | `gh issue view 30` |
| Issues #37 and #30 are open; #39 is closed (2026-11-13); open Dependabot alerts: 0 (185 fixed). | `gh issue view`, `gh api .../dependabot/alerts` |
| Repo description now reads "…Co-Founder of OrangeCat Labs… Driving business to OrangeCat Labs."; topic `orangecat-labs` exists. | `gh repo view --json` |
| Analytics events present in `frontend/src`: `chat-opened`, `game-started` only. | `git grep umami.track` |
| `origin/docs/sync-p10-p11-state` (c25523a, 971b746) was the head of merged PR #40 (squash commit 3edf2c3). `CLAUDE.md`, `README.md`, `docs/phase2-roadmap.md` and its spec and plan are byte-identical on the branch tip and on `main` before this PR. | `gh pr view 40`, `git diff --stat` |

---

## 3. Changes by File

### `CLAUDE.md`
- **Deployment:** branch-protection line (new check date, `Frontend E2E` required, strict on); the Netlify known-failing line becomes a short note that #34 is closed, PR #44 added `netlify.toml`, and previews pass.
- **Phase 2 Status:** P10 done and open lists rewritten from section 2, plus an explicit owner question on Services; P11 notes the `title` change; Dependabot triage line added.
- **Canonical Facts:** the `.dev` and LinkedIn mismatch notes become "now match" with PR #41. The canonical values are unchanged.

### `docs/phase2-roadmap.md`
- P10 Status, P11 state and scope bullet, Dependency Map row, V1 Cleanup Dependabot bullet, synced with the same facts.

---

## 4. Out of Scope

- `README.md` line 128 said Netlify previews report failures (Issue #34); corrected by the orchestrator during review.
- Editing Issue #30 (a proposed body is in the PR description for the owner), repo settings, branch deletion.
- Unmerged work by other agents (API hardening, CI supply-chain checks, analytics events, PR #45).

---

## 5. Open Questions

- Was reinstating the Services section (PR #42) intended, given Issue #30's "dropped from the plan" text? Only the owner can say; the docs state the facts and the question.
- Repo description names OrangeCat Labs, but Issue #30 still has it unticked; owner to tick it if final.
- Branch protection: `GET .../protection` has no `required_pull_request_reviews` key, while `GET .../protection/required_pull_request_reviews` returns 1 approving review. Owner to confirm which applies. Force pushes and deletions are disallowed.
- The `superpowers:verification-before-completion` skill is not installed in this environment; verification was done manually (every claim re-run with `gh`/`git`, plus an independent reviewer).

---

## 6. Deploy Impact and Verification

- Touches only `CLAUDE.md` and `docs/**`; neither deploy workflow's path filter matches.
- CI still runs all three jobs and must pass.
- An independent reviewer re-verifies every claim in section 2 with `gh` and `git`.
