# Docs Sync: P10 Progress, P11, and Verified Repo Facts — Design Spec

> Brings `CLAUDE.md`, `docs/phase2-roadmap.md`, and `README.md` in line with what is actually on `main` and with the owner's decisions of 2026-11-07. Documentation only: no source, config, workflow, or content changes.

---

## 1. Overview & Goals

- **Accuracy:** the three docs still describe P10 as "Phase 1 done", still list a Services section that was dropped, and omit P11. They also contain several statements that do not match the repository.
- **Decisions on record:** the owner's answers (canonical LinkedIn URL, production domain, Freckles Design attribution, "everything follows spec") are written down so later sessions do not re-ask or guess.
- **No assumptions:** every statement added to the docs is backed by repository contents or GitHub data read on 2026-11-07. Anything that could not be verified is left out and listed in section 6.

---

## 2. Verified State (2026-11-07)

| Fact | Source |
|---|---|
| `main` is at `9010b48` (PR #38). P10 Phase 1 (#32), og follow-up (#35), About copy (#36), and Contact copy (#38) are merged. | `git log`, GitHub PR list |
| Open PRs: none. Open issues: #11, #14, #30, #34, #37, #39. | GitHub issues API |
| `deploy-api.yml` triggers on `api/**`, `frontend/src/content/resume/**`, and `frontend/src/content/projects/**`. `deploy-frontend.yml` triggers on `frontend/**`. | workflow files |
| `ci.yml` runs three jobs on every PR to `main`: Frontend (`astro check`, build, Vitest), Frontend E2E (build, resume PDF, Playwright), API (`make sync-content`, `go vet`, `go test`, `go build`). | `ci.yml` |
| `main` is protected. Required status checks: `Frontend` and `API` (not `Frontend E2E`). Enforcement level: `non_admins`. | GitHub branches API |
| The three `rafaelswe` Netlify checks fail on PR heads with "Deploy failed" from Netlify's own preview build. Cause not readable from here. | GitHub check-runs API; tracked in #34 |
| `src/pages/resume.astro` and `src/islands/ThemeToggle.tsx` exist; neither is listed in `CLAUDE.md`. | `ls` |
| GitHub Pages is disabled (`has_pages: false`). | GitHub repo API |
| No workflow references `GEMINI_API_KEY`; it is read by the API at runtime (`api/.env.example`). | `grep`, `.env.example` |
| `go.mod` requires Go 1.26.1; CI and the Dockerfile use 1.26; the README says Go >= 1.22. | `go.mod`, `ci.yml`, `Dockerfile`, README |
| The 88 legacy Dependabot alerts were removed with the CRA files (PR #10). Issue #39 tracks the current frontend alerts. | PR #10 title, issue #39 |

---

## 3. Decisions Recorded (owner, 2026-11-07)

1. LinkedIn: `https://linkedin.com/in/rafael-a-gomes` is correct.
2. Production URL: `https://rafaelabreu.netlify.app`. `rafaelabreu.dev` is not the production URL.
3. Crediting OrangeCat Labs on the Freckles Design case study is accurate: the owner built it under the company.
4. Extra contact form fields (company, project type): undecided; stays deferred.
5. Every change follows the spec workflow (design spec, implementation plan, feature branch, PR), including copy-only and docs-only changes.
6. Standing rule from the earlier planning session: never assume; do not state conclusions without evidence (code, tests, CI, or GitHub data); ask the owner when unsure.

---

## 4. Changes by File

### `CLAUDE.md`
- **Frontend Architecture:** add `resume.astro` to Pages and `ThemeToggle.tsx` to Islands.
- **Deployment:** list the three API trigger paths and note that a resume or projects content change deploys both frontend and API. Describe all three CI jobs. Record the branch-protection facts and the known-failing Netlify checks (#34).
- **Phase 2 Status:** update the P10 entry (done: #32, #35, #36, #38; open: Freckles credit, resume headline via P11, chatbot, analytics events, page order, repo description) and add P11 (#37).
- **Docs Workflow:** every change follows spec, plan, branch, PR.
- **New "Canonical Facts" section:** production URL, LinkedIn URL, and the note that `rafaelabreu.dev` references in `profile.md` and `projects/portfolio.md` are known mismatches awaiting correction.
- **Positioning Constraints:** note that the Freckles Design attribution is owner-confirmed as accurate.
- **New "Working Rules" section:** the never-assume rule.

### `docs/phase2-roadmap.md`
- **P10 section:** rewrite Status to match section 2; mark the Services section as dropped and why; contact fields as deferred; resume headline as delivered through P11.
- **P11 section (new):** summary, scope, and impact notes taken from issue #37, plus the two URL corrections decided on 2026-11-07.
- **Dependency Map:** add P11; replace "All features are independent" with an accurate statement.
- **V1 Cleanup:** replace the stale Dependabot bullet with a pointer to #39; mark GitHub Pages as already disabled.

### `README.md`
- Go prerequisite: `>= 1.26`.
- CI/CD table: correct the `deploy-api.yml` trigger and name the `frontend/**` trigger.
- Required GitHub Secrets: remove `GEMINI_API_KEY` (no workflow reads it) and add one line saying it is an API runtime variable, not read by any workflow.

---

## 5. Out of Scope

- Source, config, workflow, content, or test changes.
- Repo description and topics (GitHub settings; proposed separately for owner approval).
- Making `Frontend E2E` a required check; branch deletion; Netlify checks (#34); Dependabot alerts (#39).
- README title and intro positioning (owner decision).
- Infrastructure cost notes in the roadmap ("$0 hosting", Plausible price): the P10 constraint covers site copy, form fields, and the chatbot, not internal docs.

---

## 6. Open Items (not asserted in the docs)

- Where `GEMINI_API_KEY` is stored (GitHub or Fly.io secret): not verifiable from the repo or API; owner to confirm with `fly secrets list` (names only).
- The cause of the Netlify preview failures: requires the Netlify deploy log.
- Whether the pending Freckles and `portfolio.md` `demo` URL fixes ship in one projects-content PR (they deploy frontend and API together).

---

## 7. Deploy Impact and Verification

- Touches only `CLAUDE.md`, `README.md`, and `docs/**`. Neither deploy workflow's path filter matches, so merging deploys nothing.
- CI still runs all three jobs on the PR (`ci.yml` has no path filter) and must pass.
- Verification: every file path named in the edited docs exists; every PR and issue number cited exists on GitHub; an independent adversarial review checks the diff against section 2 and the P10 constraints.
