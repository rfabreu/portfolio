# Plan: P11 Resume — Executive and Technical Leadership Audiences

File: docs/superpowers/plans/2026-11-07-p11-resume-executive.md
Related spec: docs/superpowers/specs/2026-11-07-p11-resume-executive-design.md
Related issue: #37 (P11)
Date: 2026-11-07
Branch: feature/p11-resume-executive

---

## Summary of Spec

Rework `frontend/src/content/resume/profile.md` for C-level, VP, and senior technical audiences: executive-voice summary, results-driven bullets, OrangeCat Labs co-founder work-history entry, reconciled certification wording, and reviewed selected_projects/skills groupings. Portfolio URL and LinkedIn already fixed in PR #41.

## Governance Note

CLAUDE.md is the governing document for this work. All references to agent workflow rules, CI check requirements (Frontend + API required; Frontend E2E not required), branch protection (no direct push to `main`), and the adversarial review process in CLAUDE.md apply. If any conflict arises between this plan and CLAUDE.md, CLAUDE.md takes precedence.

## Approach

### Phase 1: Content Decisions

Based on spec decisions:

**Title (keep P2's):** `"Co-Founder, OrangeCat Labs | Full-Stack & Infrastructure Leader"`

**Summary (executive voice):** Lead with scope (24/7 mission-critical, 16-person team, 5,000+ concurrent streams, global clients), scale, and outcomes (uptime/SLA, automation, capacity expansion), then technical means.

**Work history audit:** For each bullet, ensure it leads with scope/outcome and includes a result or scale figure where factual. Cut or merge activity-only bullets.

**OrangeCat Labs co-founder entry:** Add a separate work-history entry (company: OrangeCat Labs, title: Co-Founder, dates: co-founding to Present) — do not merge into Freelance entry.

**Certification wording:** Reconcile "Red Hat Student Education Programme Kubernetes Bootcamp" with About section's "Kubernetes/OpenShift Certificate" — use consistent, supportable wording.

**selected_projects and skills:** Review for executive + technical audiences; emphasize engineering depth and business outcomes.

### Phase 2: Implement profile.md Changes

Edit `frontend/src/content/resume/profile.md`:

1. Update `title` (keep P2's value)
2. Rewrite `summary` in executive voice
3. Audit and rewrite work-history bullets (results/scale first)
4. Add OrangeCat Labs co-founder work-history entry
5. Reconcile certification name/issuer
6. Review `selected_projects` headlines
7. Review `skills` groupings
8. Verify `links.portfolio` = `https://rafaelabreu.netlify.app` (already fixed in PR #41)
9. Verify `links.linkedin` = `https://linkedin.com/in/rafael-a-gomes` (already fixed in PR #41)

### Phase 3: Documentation Updates

- `docs/phase2-roadmap.md`: Mark P11 status
- `CLAUDE.md`: Update Phase 2 status list

### Phase 4: Verification

- `npx astro check` → 0 errors
- `npm run build` → complete (incl. PDF generation)
- `npm test -- --run` → all tests pass
- `go vet ./...`, `go test ./...`, `go build ./cmd/server` → all pass
- `grep -i 'budget\|pricing\|mailto:' frontend/src/content/resume/profile.md` → no matches
- `grep 'rafaelabreu\.dev' frontend/src/content/resume/profile.md` → no matches

### Phase 5: Commit, Push, PR

Branch: `feature/p11-resume-executive`

Commit message:
```
feat(p11): rework resume for executive and technical leadership audiences
```

PR body should:
- Reference spec and plan file paths
- List all changed files with descriptions
- Reference issue #37
- Note that `deploy-api.yml` triggers on resume content changes
- Request adversarial review before merge

## Pre-commit Checklist (per CLAUDE.md)

- [ ] Spec exists and is referenced in commit/PR
- [ ] Plan exists and follows spec
- [ ] All files pass `astro check` (frontend CI)
- [ ] API builds successfully (API CI — triggered by resume content change)
- [ ] No budget/pricing language, no public email
- [ ] No invented metrics — all figures supportable
- [ ] PR title/body follow conventions in CLAUDE.md
