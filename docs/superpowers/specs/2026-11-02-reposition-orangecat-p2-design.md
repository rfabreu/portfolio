# P10 Phase 2: Conversion & Services — Design Spec

> Expands the OrangeCat Labs repositioning by adding a Services section, converting the contact form into a high-intent Project Brief, adding case study conversion CTAs to project pages, and demoting the resume link while aligning its headline with the studio co-founder positioning.

---

## 1. Overview & Objectives

Phase 1 established studio positioning and outbound CTAs to `https://orangecatlabs.ca`.
Phase 2 turns the portfolio into a conversion engine for prospective clients:
1. **Showcase Offerings:** Add a dedicated **Services** section articulating OrangeCat Labs' 4 core capabilities.
2. **High-Intent Lead Capture:** Transform the generic contact form into a structured **Project Brief** with project type and timeline selectors, using privacy-preserving Netlify Forms. No budget or pricing fields are collected.
3. **Case Study Conversion:** Add a conversion block at the footer of project case studies directing clients to work with OrangeCat Labs.
4. **Resume Demotion & Alignment:** Demote the `/resume` trigger to a secondary link style and update its profile headline to `"Co-Founder, OrangeCat Labs | Full-Stack & Infrastructure Leader"`.
5. **Trust Links:** LinkedIn only as the secondary public contact channel.

---

## 2. Services Section (`frontend/src/components/Services.astro`)

### Structure & Placement
- **Location:** Placed in `frontend/src/pages/index.astro` directly after `Hero` (or after `About`).
- **Section ID:** `#services`
- **Header:** `SectionHeader` with `label="SERVICES"`, `title="What We Build"`.
- **Layout:** Responsive 2x2 grid (`grid md:grid-cols-2 gap-6`).

### 4 Core Capabilities
1. **Full-Stack Software Development**
   - *Headline:* Web & Custom Software
   - *Description:* High-performance web applications, specialized internal tools, and scalable APIs engineered with Go, Python, React, and modern TypeScript.
   - *Key Focus:* Responsive UI/UX, microservices, secure authentication, type-safe architectures.
2. **Cloud Architecture & DevOps**
   - *Headline:* Infrastructure & Scalability
   - *Description:* Resilient cloud environments, Kubernetes orchestration, CI/CD pipeline automation, and containerized deployments engineered for 24/7 reliability.
   - *Key Focus:* Kubernetes/OpenShift, Docker, automated testing, high-availability monitoring.
3. **AI Integration & Automation**
   - *Headline:* Intelligent Systems
   - *Description:* Practical LLM integrations, retrieval-augmented generation (RAG), and custom workflow automations that eliminate manual bottlenecks.
   - *Key Focus:* Gemini & OpenAI APIs, vector search, automated pipelines, privacy-first processing.
4. **Systems Modernization & Optimization**
   - *Headline:* Modernization & Performance
   - *Description:* Refactoring legacy codebases, database performance tuning, reducing technical debt, and migrating monolithic systems to scalable services.
   - *Key Focus:* PostgreSQL/MySQL tuning, latency reduction, code refactoring, SLA stabilization.

### Card Styling
- Card container with `bg-surface border border-surface-border rounded-card p-6`.
- Accent indicators (`text-accent-indigo`, `text-accent-purple`, `text-accent-pink`).
- Action link at bottom: *"Discuss a project →"* scrolling to `#contact` or linking to `https://orangecatlabs.ca`.

---

## 3. Project Brief Form (`frontend/src/components/Contact.astro`)

### Conversion-Oriented Fields
Replace standard generic inputs with a structured project intake brief:
- **`name`** (text, required) — Client name.
- **`email`** (email, required) — Client contact email for response.
- **`project_type`** (select, required):
  - "Full-Stack Web / Software App"
  - "Cloud Infrastructure / DevOps"
  - "AI Integration / Automation"
  - "Legacy Modernization / Refactoring"
  - "Other / Advisory"
- **`timeline`** (select, optional):
  - "Immediate (< 1 month)"
  - "1 – 3 months"
  - "3 – 6 months"
  - "Flexible"
- **`message`** (textarea, required) — "Project details, technical challenges, or goals..."
- **CTA Button:** `"Submit Project Brief →"`

### Privacy & Form Backend
- Retains Netlify Forms integration (`data-netlify="true"`, `netlify-honeypot="bot-field"`).
- Action redirects to `/success`.
- **Zero email addresses shown publicly** (strict constraint preserved).

---

## 4. Case Study Conversion CTA (`frontend/src/layouts/ProjectLayout.astro`)

At the bottom of each project case study article (above the footer), add a high-conversion callout:
```html
<div class="mt-16 p-8 bg-surface border border-accent-indigo/30 rounded-card flex flex-col md:flex-row items-center justify-between gap-6">
  <div>
    <h3 class="text-xl font-bold text-text-primary">Have a similar project in mind?</h3>
    <p class="text-text-secondary text-sm mt-1">Let's discuss how OrangeCat Labs can design, build, or scale your software.</p>
  </div>
  <a href="https://orangecatlabs.ca" target="_blank" rel="noopener noreferrer" class="px-6 py-3 bg-accent-indigo rounded-btn text-white text-sm font-semibold hover:bg-accent-indigo/90 transition-colors whitespace-nowrap">
    Start a Project with Us →
  </a>
</div>
```

---

## 5. Resume Demotion & Headline Alignment

### Demotion from Header Primary Action
- In `Nav.astro` and `MobileNav.tsx`:
  - `WORK WITH US` is the sole prominent CTA button.
  - `RESUME` remains accessible as a secondary text link (maintaining `data-resume-trigger` and `"RESUME"` text for accessibility and Playwright E2E coverage).

### Headline Alignment
- In `frontend/src/content/resume/profile.md`:
  - Update `title` from `"Operations & Infrastructure Leader"` to `"Co-Founder, OrangeCat Labs | Full-Stack & Infrastructure Leader"`.
  - Aligns the resume PDF, modal, and `/resume` page with the studio leadership position.

---

## 6. Files Changed / Added

| File | Change | Description |
|---|---|---|
| `frontend/src/components/Services.astro` | New | Services section highlighting 4 studio offerings |
| `frontend/src/components/Contact.astro` | Modified | Upgrade form to Project Brief with type, timeline |
| `frontend/src/pages/index.astro` | Modified | Mount `<Services />` and updated section dividers |
| `frontend/src/components/Nav.astro` | Modified | Demote `RESUME` to secondary link; add `Services` nav item |
| `frontend/src/islands/MobileNav.tsx` | Modified | Add `Services` to mobile navigation |
| `frontend/src/layouts/ProjectLayout.astro` | Modified | Add project case study conversion callout |
| `frontend/src/content/resume/profile.md` | Modified | Update title to Co-Founder, OrangeCat Labs |
| `frontend/tests/reposition-services-brief.test.ts` | New | Unit tests verifying Services, Brief form, and CTA |
