# P9: Resume Feature — Design Spec

> Replaces the currently broken Resume button (404 on `/Rafael_Abreu_Resume.pdf`) with a hand-edited markdown source rendered to a styled in-site modal, a downloadable PDF generated at build time, and a chatbot system prompt assembled from the same source plus the existing project case studies.

## Overview

The portfolio's Resume button currently 404s — no PDF exists in `frontend/public/`. Rather than dropping a static PDF file, this feature establishes a sustainable single-source-of-truth model: Rafael edits one markdown file, and the site renders three artifacts from it (modal, PDF, chatbot context).

The feature explicitly avoids the AI-driven pipeline considered during brainstorming (LinkedIn ingestion, GitHub-derived skill extraction, AI-generated bullets, accuracy validators). For an artifact that changes 4–6 times per year, that pipeline's maintenance cost outweighed its value. The simple version captures ~95% of the value with ~5% of the moving parts.

## Architecture

```
                ┌──────────────────────────────────────────┐
                │  frontend/src/content/resume/profile.md  │  ← single source of truth
                └──────────────────┬───────────────────────┘
                                   │
        ┌──────────────────────────┼──────────────────────────┐
        │                          │                          │
        ▼                          ▼                          ▼
   Astro build              Playwright (CI)             CI sync step
   (modal + page)           PDF generation              (copies markdown)
        │                          │                          │
        ▼                          ▼                          ▼
   ResumeModal.astro       Rafael_Abreu_              api/internal/chat/
   + /resume page          Resume.pdf                 embedded/
                           (frontend/dist)            ↓
                                                Go embed → SystemPrompt()
                                                (resume + projects + persona)
```

Three render targets, one source. All rendering is deterministic at build time; no AI in the request path. The chatbot is grounded in both the resume and the existing project case studies (`frontend/src/content/projects/*.md`), fixing a latent staleness bug where the current hardcoded `SystemPrompt()` drifts from the actual portfolio.

## Data Model

The resume lives at `frontend/src/content/resume/profile.md` inside a new Astro content collection. The frontmatter carries all structured data; the markdown body is reserved for free-form narrative, expected to remain empty in practice.

```yaml
---
name: "Rafael Abreu"
title: "Software Engineer | Engineering Leadership Track"
location: "Toronto, ON, Canada"

links:
  linkedin:  "https://www.linkedin.com/in/rafael-abreu"
  github:    "https://github.com/rfabreu"
  portfolio: "https://rafaelabreu.netlify.app"

summary: |
  One-paragraph executive-positioned summary. Roughly three sentences:
  who you are, what you've delivered, where you're heading.

work_history:
  - company:    "Nextologies"
    title:      "Software Engineer"
    dates:      "2023 — Present"
    location:   "Toronto, ON"
    tech_stack: ["Go", "Python", "Kubernetes", "Docker", "AWS"]
    bullets:
      - "Architected ..."
      - "Drove ..."

education:
  - school:  "University of Toronto"
    program: "Coding Bootcamp — Full-Stack Web Development"
    dates:   "2022 — 2023"

certifications:
  - name:   "Red Hat Certified Specialist in OpenShift Administration"
    issuer: "Red Hat"
    year:   "2024"

skills:
  leadership:
    - "Strategic Planning"
    - "Cross-functional Stakeholder Management"
    - "Hiring & Team Building"
  technical:
    - "Go, Python, TypeScript"
    - "Kubernetes, Docker, CI/CD"
    - "AI/ML Integration (Gemini, RAG)"
  domain:
    - "Broadcast / Live TV Operations"
    - "EdTech"

awards_recognition: []      # optional, fills over time
speaking_writing:   []      # optional, fills over time

languages_spoken:
  - "English (Native)"
  - "Portuguese (Native)"

selected_projects:
  - slug:     "fptv"        # references frontend/src/content/projects/fptv.md
    headline: "Glassmorphism dashboard for Toronto-based TV broadcaster"
  - slug:     "portfolio"
    headline: "Astro + Go monorepo with AI chatbot and game"
---
```

### Schema Validation

`frontend/src/content.config.ts` defines a Zod schema for the `resume` collection. Required fields: `name`, `title`, `location`, `links.linkedin`, `summary`, `work_history`, `education`, `skills`. Optional fields: `certifications`, `awards_recognition`, `speaking_writing`, `languages_spoken`, `selected_projects`. The `bullets` field on each `work_history` entry defaults to an empty array, allowing a job to be added before bullets are written. Invalid frontmatter fails the build with a useful error pointing at the offending field.

`selected_projects[].slug` uses Astro's `reference('projects')` Zod helper to validate at build time that the referenced slug exists in the `projects` collection. A typo or stale reference fails the build immediately, rather than silently rendering a broken link at runtime.

### Design Decisions

- **No phone or email in the schema.** Per the project's no-direct-contact-channel policy. Rendered footer in modal and PDF reads: *"Contact via the form on rafaelabreu.netlify.app or LinkedIn"*.
- **Skills are bucketed by category** (`leadership` / `technical` / `domain`) so the renderer can group without inference.
- **No skill levels** (`level: senior`). Modern senior resumes don't use them.
- **`selected_projects` references existing project slugs** rather than duplicating content. The renderer looks up the corresponding case study and links to it.
- **Freeform `dates` strings** ("2023 — Present"). Sortability isn't needed for ~5 jobs ever; structured dates would be friction without payoff.
- **`awards_recognition` and `speaking_writing` start empty.** Optional sections render nothing when empty. Strong exec-track resumes accumulate entries here over time; having the structure ready means there's no friction when content arrives.
- **`tech_stack` per work entry** keeps bullets outcome-focused. The renderer shows the stack as small tags below each role's bullets.

## Components

### Modal — `frontend/src/components/ResumeModal.astro`

Pure Astro component using the native HTML `<dialog>` element. No React island. The modal has only open/close behavior — that's all native to `<dialog>` (ESC handling, focus trap, click-outside via the `::backdrop` pseudo-element). Eliminates a hydrating bundle.

```astro
<dialog id="resume-modal" class="resume-modal">
  <header class="resume-modal-header">
    <div class="resume-modal-title">
      <span class="name">{resume.name}</span>
      <span class="title">{resume.title}</span>
    </div>
    <div class="resume-modal-actions">
      <a href="/Rafael_Abreu_Resume.pdf" download class="btn-download">Download PDF</a>
      <button type="button" class="btn-close" aria-label="Close">×</button>
    </div>
  </header>
  <div class="resume-modal-body">
    <ResumeBody resume={resume} />
  </div>
  <footer class="resume-modal-footer">
    <span>Contact via <a href="/#contact">contact form</a> or
    <a href={resume.links.linkedin}>LinkedIn →</a></span>
  </footer>
</dialog>

<script>
  const modal = document.getElementById('resume-modal') as HTMLDialogElement | null;

  // Event delegation — works for static Astro triggers AND React-hydrated triggers
  // (MobileNav button is rendered after page load, so per-element binding would miss it).
  document.addEventListener('click', (e) => {
    const trigger = (e.target as HTMLElement)?.closest('[data-resume-trigger]');
    if (trigger) {
      e.preventDefault();
      modal?.showModal();
    }
  });

  modal?.querySelector('.btn-close')?.addEventListener('click', () => modal.close());
  modal?.addEventListener('click', (e) => {
    if (e.target === modal) modal.close();
  });
</script>
```

The modal is included site-wide via `BaseLayout.astro`, **except on the `/resume` route** (which already renders the same content inline). The layout accepts an `includeResumeModal` prop that defaults to `true`; `resume.astro` sets it to `false`. Triggers across the site (Nav, MobileNav, hero CTA if added later) all use the `data-resume-trigger` attribute. Event delegation means the React-hydrated MobileNav button works the same as the static Astro Nav anchor — both add `data-resume-trigger` to their rendered button element.

Because the modal isn't on `/resume`, the Resume button in Nav and MobileNav must also be suppressed there — clicking a `data-resume-trigger` element on `/resume` would be a no-op without the modal listener installed. Both `Nav.astro` and `MobileNav.tsx` accept the current pathname (via `Astro.url.pathname` for the former, a prop for the latter) and conditionally hide the Resume button when `pathname === '/resume'`. The Nav remains otherwise unchanged — visitors on the resume page have already arrived at the destination the button would have taken them to.

### Print Page — `frontend/src/pages/resume.astro`

Standalone full-page route, A4 layout, print-stylesheet-aware. Same data as the modal, different shell. Two purposes:

1. Playwright loads `http://localhost:<port>/resume` during CI to print the PDF.
2. Anyone with the URL can view the resume directly without going through the home page.

Uses the `BaseLayout` shell with a `print-mode` class that suppresses nav/footer in `@media print` and forces the PDF-friendly typography.

### Shared Section Components — `frontend/src/components/resume/`

Each frontmatter section has its own Astro component. Both `ResumeModal.astro` and `resume.astro` import these via a `<ResumeBody />` wrapper:

| Component | Renders |
|-----------|---------|
| `Summary.astro` | `summary` paragraph |
| `Skills.astro` | `skills` grouped by category |
| `WorkHistory.astro` | `work_history` entries with bullets and `tech_stack` tags |
| `Education.astro` | `education` entries |
| `Certifications.astro` | `certifications` list |
| `Awards.astro` | `awards_recognition` (skipped if empty) |
| `Speaking.astro` | `speaking_writing` (skipped if empty) |
| `Languages.astro` | `languages_spoken` list |
| `Projects.astro` | `selected_projects` with slug-resolved links to case studies |

Single source of layout per section, two render contexts. Empty optional sections render nothing — no visual penalty for blanks.

### PDF Generator — `frontend/scripts/generate-pdf.mjs`

Node script using Playwright. Runs in CI after `npm run build`:

```javascript
import { chromium } from 'playwright';
import { createServer } from 'http';
import handler from 'serve-handler';

const server = createServer((req, res) => handler(req, res, { public: 'dist' }));
// Port 0 = OS-assigned free port. Avoids collision with Astro's default dev port (4321).
await new Promise((resolve) => server.listen(0, resolve));
const { port } = server.address();

const browser = await chromium.launch();
const page = await browser.newPage();
await page.goto(`http://localhost:${port}/resume`, { waitUntil: 'load' });
await page.emulateMedia({ media: 'print' });
await page.pdf({
  path: 'dist/Rafael_Abreu_Resume.pdf',
  format: 'A4',
  printBackground: true,
  margin: { top: '0.5in', right: '0.5in', bottom: '0.5in', left: '0.5in' },
});

await browser.close();
server.close();
```

The PDF is written into `frontend/dist/`, which Netlify deploys. Resume button (already pointing at `/Rafael_Abreu_Resume.pdf`) starts working without a code change.

### Chatbot Integration — `api/internal/chat/`

Replace the hardcoded content in `SystemPrompt()` with embedded markdown files plus a hardcoded persona/rules block.

**New directory:** `api/internal/chat/embedded/` (gitignored, populated by CI).

**Modified `context.go`:**

```go
package chat

import (
    "embed"
    "fmt"
    "io/fs"
    "sort"
    "strings"
)

//go:embed embedded/resume.md
var resumeMD string

//go:embed embedded/projects/*.md
var projectsFS embed.FS

func SystemPrompt() string {
    var b strings.Builder
    b.WriteString(personaBlock())
    b.WriteString("\n\n## RESUME\n\n")
    b.WriteString(resumeMD)
    b.WriteString("\n\n## PROJECT CASE STUDIES\n\n")
    b.WriteString(assembleProjects())
    b.WriteString("\n\n")
    b.WriteString(rulesBlock())
    return b.String()
}

func assembleProjects() string {
    entries, _ := fs.ReadDir(projectsFS, "embedded/projects")
    sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
    var b strings.Builder
    for _, entry := range entries {
        content, err := fs.ReadFile(projectsFS, fmt.Sprintf("embedded/projects/%s", entry.Name()))
        if err != nil { continue }
        b.WriteString(string(content))
        b.WriteString("\n\n---\n\n")
    }
    return b.String()
}
```

`personaBlock()` and `rulesBlock()` hold the existing "You are Rafael Abreu's AI assistant…" and "Only answer questions about Rafael, be concise…" content respectively, lifted out of the current monolithic string.

**CI sync step** (in `deploy-api.yml`, before `go build`):

```yaml
- name: Sync content for embedding
  run: |
    rm -rf api/internal/chat/embedded
    mkdir -p api/internal/chat/embedded/projects
    cp frontend/src/content/resume/profile.md api/internal/chat/embedded/resume.md
    cp frontend/src/content/projects/*.md api/internal/chat/embedded/projects/
```

**Local dev:** new `api/Makefile` target `sync-content` runs the same copy. `make dev` runs `sync-content && go run ./cmd/server`.

**First-time clone:** because `api/internal/chat/embedded/` is gitignored, a fresh clone won't compile the API until `make sync-content` (or `make dev`) runs once — `go:embed` requires the referenced paths to exist. README and `api/Makefile` should both document this. Alternative considered: commit the `embedded/` directory. Rejected because (a) it duplicates content already version-controlled in `frontend/`, and (b) it creates merge-conflict noise on every resume edit.

## Data Flow

```
Edit profile.md
       │
       ▼
push → feature branch → PR → review → merge to main
       │
       ├── deploy-frontend.yml
       │     ├── npm ci
       │     ├── npm run build
       │     ├── npx playwright install chromium
       │     ├── node scripts/generate-pdf.mjs   (writes dist/Rafael_Abreu_Resume.pdf)
       │     └── netlify deploy
       │
       └── deploy-api.yml
             ├── sync embedded markdown
             ├── go vet ./...
             ├── go test ./...
             ├── go build
             └── flyctl deploy
```

Both deploys complete within a few minutes of the push. Eventually consistent — there's a brief window where the frontend is live with new resume content but the API still has the old prompt. Acceptable for this surface (chatbot will catch up, no user-facing inconsistency).

## Error Handling

| Failure point | Behavior | Reasoning |
|---|---|---|
| Malformed `profile.md` frontmatter | Astro build fails with Zod error | Don't deploy a broken resume |
| Missing required field | Build fails with field path | Same |
| Playwright PDF generation fails | Frontend deploy fails | Don't ship the site without a working PDF — keeps the Resume button reliable |
| `embedded/` empty at API build | `go:embed` fails compile | Don't ship API with empty grounding |
| Modal runtime exception | Native `<dialog>` is bulletproof — closes via ESC/click-outside regardless | No JS framework runtime to fail |
| Resume PDF URL 404 in browser | Cannot happen post-fix; PDF is rebuilt every deploy | Fixes the existing bug |
| Project markdown referenced by `selected_projects` slug doesn't exist | Build fails with a Zod error pointing at the bad slug | Hard fail via `reference('projects')` — consistent with the 100%-accuracy stance; a typo silently rendering a broken link is worse than a build break that takes ten seconds to fix |

The pattern: **fail loudly at build time, silently never at runtime.** Bad data never reaches users.

## Testing

### Frontend

- **Schema validation:** Zod schema in `content.config.ts` validates `profile.md` on every build. Build is the test.
- **Snapshot tests:** rendered HTML for `ResumeModal.astro` and `/resume` page. Catches layout regressions.
- **E2E (Playwright, new):** click Resume nav button → modal opens → expected text present → "Download PDF" link resolves to `/Rafael_Abreu_Resume.pdf`.

### PDF

- **Build assertion:** `frontend/dist/Rafael_Abreu_Resume.pdf` exists, size > 10KB after `generate-pdf.mjs`.
- **Manual visual check** on first build. No automated visual regression — overkill for a single PDF rendered from a deterministic source.

### API

- Existing chat tests continue to pass.
- **New: `TestSystemPromptIncludesResumeContent`** — assembled prompt contains a known string from `profile.md` (e.g., the canonical name).
- **New: `TestSystemPromptIncludesAllProjects`** — for each file in `embedded/projects/`, assert its content appears in the prompt.
- **New: `TestSystemPromptStructure`** — assert persona block precedes resume, resume precedes projects, rules block is last.

### CI Gate

Frontend job + API job + new E2E job all green before merge. PR cannot merge with any failing.

## Files Changed

| File | Change |
|------|--------|
| `frontend/src/content/resume/profile.md` | New — single source of truth |
| `frontend/src/content.config.ts` | Add `resume` collection with Zod schema |
| `frontend/src/components/ResumeModal.astro` | New — `<dialog>`-based modal |
| `frontend/src/components/resume/Summary.astro` | New — section component |
| `frontend/src/components/resume/Skills.astro` | New — section component |
| `frontend/src/components/resume/WorkHistory.astro` | New — section component |
| `frontend/src/components/resume/Education.astro` | New — section component |
| `frontend/src/components/resume/Certifications.astro` | New — section component |
| `frontend/src/components/resume/Awards.astro` | New — section component |
| `frontend/src/components/resume/Speaking.astro` | New — section component |
| `frontend/src/components/resume/Languages.astro` | New — section component |
| `frontend/src/components/resume/Projects.astro` | New — section component |
| `frontend/src/components/resume/ResumeBody.astro` | New — assembles all section components |
| `frontend/src/pages/resume.astro` | New — print-styled standalone page |
| `frontend/src/layouts/BaseLayout.astro` | Include `<ResumeModal />` once site-wide; accept `includeResumeModal` prop (defaults `true`) for opt-out on `/resume` |
| `frontend/src/components/Nav.astro` | Resume button → `data-resume-trigger`, no longer href to PDF; conditionally hidden when `Astro.url.pathname === '/resume'` |
| `frontend/src/islands/MobileNav.tsx` | Resume button → `data-resume-trigger`; receives current pathname as prop, conditionally hidden on `/resume` |
| `frontend/src/styles/global.css` | Add resume modal + print styles, theme-token-aware |
| `frontend/scripts/generate-pdf.mjs` | New — Playwright PDF generator |
| `frontend/package.json` | Add `playwright` and `serve-handler` as devDependencies |
| `frontend/tests/resume.spec.ts` | New — Playwright E2E: nav click → modal opens → expected text → PDF link resolves |
| `frontend/tests/__snapshots__/resume-modal.snap` | New — snapshot for `ResumeModal.astro` rendered HTML |
| `frontend/tests/__snapshots__/resume-page.snap` | New — snapshot for `/resume` page rendered HTML |
| `api/internal/chat/context.go` | Refactor — split into persona/resume/projects/rules; use `go:embed` |
| `api/internal/chat/context_test.go` | New tests — `TestSystemPromptIncludesResumeContent`, `TestSystemPromptIncludesAllProjects`, `TestSystemPromptStructure` |
| `api/internal/chat/embedded/` | New gitignored directory; populated by CI sync |
| `api/Makefile` | New — `sync-content` and `dev` targets |
| `api/.gitignore` | Add `internal/chat/embedded/` |
| `.github/workflows/deploy-frontend.yml` | Add Playwright install + PDF generation step |
| `.github/workflows/deploy-api.yml` | Add content sync step before `go build` |
| `.github/workflows/ci.yml` | Add resume schema validation + E2E job |
| `docs/phase2-roadmap.md` | Add P9: Resume Feature entry |

## Out of Scope

The following were considered during brainstorming and explicitly excluded:

- **AI-driven resume generation pipeline** — LinkedIn ingestion, GitHub-derived skill extraction, AI-generated bullets, accuracy validators, PR-bot suggestions. The maintenance and noise outweigh the benefit at portfolio scale (resumes change 4–6 times/year). Reconsider only if manual updating becomes a friction point.
- **LinkedIn API integration** — official API requires Partner approval; scrapers violate ToS or break on layout changes. The realistic path (LinkedIn data export, manually refreshed) was not worth the ingestion plumbing for this version.
- **Auto-update from GitHub webhooks or portfolio commits** — same reasoning. New repos rarely warrant resume changes.
- **A separate AI editor / "rephrase this bullet" tool** — Rafael uses any LLM externally when editing the markdown. Not a build-pipeline feature.
- **Phone or email on the resume** — per the no-direct-contact-channel policy. Contact form + LinkedIn cover all reach-out paths.
- **Photo / avatar** — North American convention; introduces bias risk.
- **Skill levels (`level: senior`)** — modern senior resumes don't use them.
- **Multi-variant resumes** (e.g., backend-focused vs full-stack) — single canonical version for this iteration. Variants can be added later if a clear use case emerges.
- **Visual regression testing for the PDF** — overkill for a single deterministic artifact; manual eyeball on first build suffices.
- **"Site = resume" radical consolidation** — making the entire portfolio render as a resume in print mode. Larger lift, considered for future P-item if/when the markdown-resume model proves limiting.
