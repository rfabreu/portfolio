# Design Spec: P11 Resume — Executive and Technical Leadership Audiences

File: docs/superpowers/specs/2026-11-07-p11-resume-executive-design.md
Related: Issue #37 (P11), Issue #30 (P10)
Created: 2026-11-07
Status: Draft for review

---

## Context

P11 (Issue #37) reworks the resume content so it reads as credible to C-level, VP, and senior technical audiences, with business-oriented framing (scope, outcomes, risk, scale) alongside technical depth. The resume is the secondary link behind the OrangeCat Labs call to action (P10), so it should match the founder voice the rest of the site now has.

`frontend/src/content/resume/profile.md` is the single source for: the resume modal, the `/resume` page, the generated PDF, and the chatbot's grounding. Changes here trigger both the frontend and API deploys via `deploy-api.yml`.

## Constraints (from Issue #37 and P10)

- Only claims that are true and supportable; metrics must come from Rafael. No invented figures.
- No budget or pricing language anywhere.
- No email address on the resume.
- Only Freckles Design is linked to OrangeCat Labs (portfolio projects). How the resume describes the co-founder role is a decision for this issue.
- Work follows spec → plan → branch → PR, adversarial review, green CI, owner approval (CLAUDE.md).
- Never push to `main`.

## Decisions Needed

### 1. Headline / Title

Current: `"Operations & Infrastructure Leader"`
P2 (PR #42) changed it to: `"Co-Founder, OrangeCat Labs | Full-Stack & Infrastructure Leader"`

Options:
- A. Keep P2's title as-is — co-founder first, signals studio leadership
- B. `"Co-Founder, OrangeCat Labs"` — studio-first, minimal
- C. `"Full-Stack & Infrastructure Leader, Co-Founder OrangeCat Labs"` — technical-first

Recommendation: Keep P2's title (Option A) — it leads with the co-founder role while retaining technical credibility.

### 2. Summary Rewrite (Executive Voice)

Current summary leads with "Operations & Infrastructure Leader and hands-on software engineer" — too individual-contributor focused.

Rewrite to lead with scope, scale, and outcomes:
- Scope: 24/7 mission-critical environments, 16-person team, 5,000+ concurrent streams
- Scale: global clients, 24/7 high-availability
- Outcomes: uptime/SLA ownership, automation, cost reduction, capacity expansion
- Voice: executive, outcome-first, then technical means

### 3. Bullet Audit (Results / Scale)

Audit every work-history bullet:
- Add a result or scale figure where factual (existing strong examples: 92% backlog reduction, 80% viewership increase, 62% traffic increase, 16-person team, 5,000+ streams)
- Cut or merge activity-only bullets that lack a measurable result
- Lead each role with scope and business outcome, then technical means

### 4. OrangeCat Labs in Work History

Current: "Freelance Software Engineer" (Dec 2021 – Present) entry covers portfolio work.

Options:
- A. Add a separate "OrangeCat Labs — Co-Founder" entry (dates: co-founding to Present)
- B. Merge co-founder role into the Freelance entry (add subtitle/description)
- C. Keep Freelance entry, add co-founder as a separate company entry

Recommendation: Option A — separate co-founder entry makes the studio role explicit and searchable.

### 5. Certification Wording

Current: "Red Hat Student Education Programme Kubernetes Bootcamp" (2022)
About section says: "Red Hat Kubernetes certificate (OpenShift)"

Reconcile to a consistent, supportable wording.

### 6. selected_projects and Skills

Review `selected_projects` and skills groupings for executive + technical audiences. Ensure the portfolio project headline emphasizes engineering depth and business outcomes.

### 7. URLs (Already Fixed in PR #41)

- `links.portfolio`: `https://rafaelabreu.netlify.app` ✅ (was `.dev`)
- `links.linkedin`: `https://linkedin.com/in/rafael-a-gomes` ✅ (was `rfabreu/`)

## Scope

| File | Change |
|------|--------|
| `frontend/src/content/resume/profile.md` | Title, summary, work history bullets, add OrangeCat Labs co-founder entry, cert wording, selected_projects, skills |
| `docs/superpowers/specs/2026-11-07-p11-resume-executive-design.md` | This spec |
| `docs/superpowers/plans/2026-11-07-p11-resume-executive.md` | Implementation plan (to be drafted after spec approval) |
| `docs/phase2-roadmap.md` | Mark P11 status |
| `CLAUDE.md` | Phase 2 status list |

## Out of Scope (per Issue #37)

- AI-driven resume generation or LinkedIn ingestion
- Redesigning the resume layout or components (unless new content requires it)
- Chatbot updates (handled in P10 Phase 3, separate spec)

## Deploy Impact

- `deploy-api.yml` triggers on `frontend/src/content/resume/**` — both frontend and API deploy on merge
- Chatbot grounding changes with resume content — spot-check chatbot answers after deploy
- PDF regeneration required in CI (`npm run build:pdf`)

## Acceptance Criteria

1. Headline/title positions Rafael as Co-Founder, OrangeCat Labs with technical credibility
2. Summary rewritten in executive voice: scope, scale, outcomes first
3. Every work-history bullet has a result or scale figure where factual; activity-only bullets cut or merged
4. OrangeCat Labs co-founder role present in work history
5. Certification wording reconciled with About section
6. `links.portfolio` = `https://rafaelabreu.netlify.app`, `links.linkedin` = `https://linkedin.com/in/rafael-a-gomes`
7. No budget/pricing language, no public email
8. CI passes: `astro check`, `npm run build` (incl. PDF), `go vet`, `go test`, `go build`
9. Adversarial review complete: no invented metrics, no budget/pricing, correct links only

## Files Changed (expected)

- `frontend/src/content/resume/profile.md`
- `docs/superpowers/specs/2026-11-07-p11-resume-executive-design.md`
- `docs/superpowers/plans/2026-11-07-p11-resume-executive.md` (after spec approval)
- `docs/phase2-roadmap.md`
- `CLAUDE.md`

Branch name: `feature/p11-resume-executive`
