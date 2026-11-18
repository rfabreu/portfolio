# P10 Phase 3: Analytics Funnel Events — Design Spec

> Adds Umami custom events for the business funnel: every click through to OrangeCat Labs, the "View Projects" CTA, LinkedIn outbound clicks, contact form submission attempts, and successful submissions. Covers the "Track CTA clicks, form submissions, and outbound clicks to orangecatlabs.ca in the existing analytics" item from Issue #30.

---

## 1. Overview & Objectives

The site exists to drive business to OrangeCat Labs, but the only custom events today are `chat-opened` and `game-started` (P5). There is no way to tell which placement sends visitors to `orangecatlabs.ca`, or how many visitors start and complete the contact form.

Goals:
1. Attribute every OrangeCat Labs click to the placement that produced it.
2. Track the secondary CTA (View Projects) and LinkedIn outbound clicks.
3. Track contact form attempts and successes so the funnel (visit → click / attempt → success) can be read from one dashboard.
4. Send no personal data, change nothing visible (copy, layout, styling, targets), and break nothing when Umami is absent. The one behavioural difference is the same-tab LinkedIn links on `/resume` (open question 2).

Non-goals: Playwright funnel E2E specs (separate follow-up PR), GitHub/demo outbound links, internal anchor links (`#contact`), new copy, new analytics vendors, other changes to the Umami script tags (the one exception, `data-domains`, is in open question 7).

---

## 2. Umami Behaviour (verified)

Source: https://umami.is/docs/track-events and the live tracker `https://cloud.umami.is/script.js`, read 2026-11-17.

- `data-umami-event="<name>"` on an element records an event when the element, or a descendant, is clicked. `data-umami-event-<prop>="<value>"` adds event data. All values arrive as strings. Event names are limited to 50 characters.
- The tracker registers one capture-phase `click` listener on `document` and resolves the target with `closest('[data-umami-event]')`. It therefore works for React-rendered markup (MobileNav) and for elements added after load, and it runs before the React `onClick` that closes the mobile menu.
- For an `<a>` with an `href`:
  - `target="_blank"`, or a Ctrl/Shift/Meta/middle click: the tracker does **not** call `preventDefault`. It sends the event and the browser opens the tab as normal. Every OrangeCat Labs link on the site opens in a new tab.
  - Any other click (same tab): the tracker calls `preventDefault`, sends the event, and only after the request settles assigns `location.href`. Navigation is delayed by one network round trip.
- Tracking is skipped entirely when the script has not loaded (`astro dev`, blockers). The site's script tags do not set `data-do-not-track`, so the tracker does not honour the browser's Do Not Track setting. No element needs a fallback.
- Netlify deploy previews run the same production build (`netlify.toml`, `[context.deploy-preview]`) and the script tags set no `data-domains`, so previews load the tracker and report to the production Umami site, tagged with the preview hostname. This was already true for pageviews; the new events add to it (open question 7).
- The script tag loads only when `import.meta.env.PROD`. This spec does not change that.

Consequences for the design:
- Declarative attributes are used wherever the link opens a new tab, because they cannot affect navigation.
- "View Projects" is a same-tab, same-page anchor (`#projects`). The declarative attribute would make the tracker intercept the click and delay the scroll, so it uses an imperative listener that never calls `preventDefault` (section 4.3).
- A declarative attribute on the form's submit `<button>` would fire on every click, including clicks that fail HTML5 validation, so it would not measure "submission attempts". Attempts use the form's `submit` event (section 4.4).

---

## 3. Event Catalogue

Naming: kebab-case event names, matching `chat-opened` and `game-started`. One event per destination or action, with the placement carried in a `location` property. Property names and values are fixed strings authored in source; none is derived from user input.

| Event | Properties | Fires when | Mechanism |
|-------|-----------|-----------|-----------|
| `orangecat-click` | `location`; `project` (only for `project-cta`) | A link to `https://orangecatlabs.ca` is clicked | Declarative |
| `linkedin-click` | `location` | A link to Rafael's LinkedIn profile is clicked | Declarative |
| `view-projects-click` | `location` = `hero` | The Hero "View Projects" link is clicked | Imperative |
| `contact-form-submit` | none | The contact form passes validation and is submitted | Imperative |
| `contact-form-success` | none | The `/success` page finishes loading | Imperative |

### `orangecat-click` locations

| `location` | Element | File |
|-----------|---------|------|
| `hero` | "Work with OrangeCat Labs" primary CTA | `src/components/Hero.astro` |
| `nav` | Desktop "WORK WITH US" | `src/components/Nav.astro` |
| `mobile-nav` | Mobile menu "WORK WITH US" | `src/islands/MobileNav.tsx` |
| `footer` | "Co-Founded by Rafael · OrangeCat Labs" credit | `src/components/Footer.astro` |
| `about` | "OrangeCat Labs" inline link | `src/components/About.astro` |
| `contact` | "orangecatlabs.ca (Studio Site)" link | `src/components/Contact.astro` |
| `project-cta` | "Start a Project with Us →" (every case study) | `src/layouts/ProjectLayout.astro` |
| `freckles-credit` | "Built by OrangeCat Labs" credit | `src/content/projects/freckles-design.md` |

`project-cta` also sends `project` = the case study slug (`project.id`), passed from `src/pages/project/[slug].astro`, so conversions can be attributed to a case study. It is omitted when no slug is supplied.

### `linkedin-click` locations

| `location` | Element | File |
|-----------|---------|------|
| `footer` | "LinkedIn" | `src/components/Footer.astro` |
| `contact` | "linkedin.com/in/rafael-a-gomes" | `src/components/Contact.astro` |
| `resume-modal` | "LinkedIn →" in the modal footer | `src/components/ResumeModal.astro` |
| `resume-page-header` | "LinkedIn" in the page header | `src/pages/resume.astro` |
| `resume-page-footer` | "LinkedIn" in the page footer | `src/pages/resume.astro` |

Both resume surfaces read the URL from `profile.md`; the attribute is tied to the element, not the URL, so it survives the pending LinkedIn URL correction (Issue #37).

### Intentionally untracked
- GitHub and live-demo links.
- Third-party LinkedIn links in `TestimonialCard.astro` (section is hidden; those are the testimonial authors' profiles, not Rafael's).
- `#contact` anchors (Services card, chat widget). Internal navigation, covered by pageviews.

---

## 4. Implementation

### 4.1 Declarative attributes
Add `data-umami-event` and `data-umami-event-location` to each anchor in the tables above. No class, text, `href`, `target`, `rel` or `aria-label` changes, so rendered visuals are identical.

### 4.2 Case study credit
`freckles-design.md` keeps its raw HTML anchor and gains the two attributes. Astro's Markdown pipeline passes raw HTML attributes through unchanged (verified in the build output, section 6).

### 4.3 View Projects
`Hero.astro` gets a small `<script>` that adds a `click` listener to the anchor (selected by a `data-view-projects` hook attribute) and calls `umami.track('view-projects-click', { location: 'hero' })` behind `typeof umami !== 'undefined'`. The listener never prevents default, so scrolling is exactly as before.

### 4.4 Contact form
`Contact.astro` gets a small `<script>` that listens for the `submit` event on `form[name="contact"]` and calls `umami.track('contact-form-submit')` with **no data**, behind the `typeof umami` guard. The `submit` event only fires after native validation passes, so the count is real attempts. The tracker's request uses `fetch(..., { keepalive: true })`, so it survives the navigation to `/success`. The script never reads `FormData`, field values or `.value`.

### 4.5 Success page
`success.astro` gets a small `<script>` that calls `umami.track('contact-form-success')` on `window.load`, behind the `typeof umami` guard. The Umami script tag is `async`; async scripts delay `window.load`, so waiting for `load` avoids running before `umami` is defined. The site's form posts to `/success` (Netlify Forms redirects the submit there), so this page load is the client-side success signal. Umami's automatic pageview already records `/success` visits; the named event makes the success count directly readable in Events and carries the same over-counting caveat (open question 3).

---

## 5. Privacy

- No form field value, name, email, project type, timeline or message text is read or sent, by any event.
- Event data keys are limited to `location` and `project`. A test scans all of `src/` and fails if any other `data-umami-event-*` key is introduced.
- `project` holds a public case study slug. `location` holds a fixed placement label.
- Umami Cloud is already cookieless (P5); no consent or script changes are needed.
- The positioning constraints (P10) are unaffected: no pricing language, no Upwork link, no new contact paths. Only Freckles Design is attributed to OrangeCat Labs, and no attribute adds a new OrangeCat Labs link.

---

## 6. Testing & Verification

Vitest, matching `tests/reposition-hero-nav.test.ts` and `tests/contact-copy.test.ts`:
- Each tracked element, rendered through the Astro container, carries the expected `data-umami-event` and `data-umami-event-location`.
- MobileNav is rendered with its menu open (React `useState` mocked in its own test file) to assert its attributes, since the menu is not in the server-rendered HTML.
- Source scan: every `<a>` under `src/**` (`.astro`, `.tsx`, `.md`) whose `href` points at `orangecatlabs.ca` (any scheme, `www.`, path or quoting, or `ORANGECAT_URL`) or at Rafael's LinkedIn profile carries a `data-umami-event`, and the set of locations equals the catalogue. Markdown-syntax links to either domain are rejected, since they cannot carry attributes. Limitation: the tag matcher stops at a bare `>` inside an attribute expression (`=>` is handled), so such an anchor would need a manual check.
- Source scan: only `location` and `project` data keys exist; the imperative scripts reference no form values; every `umami.track` call sits directly behind its own `typeof umami !== 'undefined'` guard.

Build verification (`frontend/`): `npm ci`, `npx astro check`, `npm run build`, `npm test`; then inspect `frontend/dist/` for the attributes on the home page, the Freckles Design case study, another case study (project slug), `/resume`, `/success`, and the bundled MobileNav chunk.

Not verifiable here: production analytics. The owner must open the Umami dashboard after deploy and confirm each event arrives (section 7, item 1).

---

## 7. Open Questions

1. **Production verification.** Events are only produced when the Umami script is loaded (production). After deploy, the owner should click each tracked element once and confirm it in Umami → Events (and Properties for `location`). Behaviour above was read from the live tracker on 2026-11-17; Umami can change the script.
2. **Same-tab LinkedIn links on `/resume`.** The two LinkedIn anchors in `resume.astro` have no `target="_blank"`, so Umami intercepts the click and navigates after the event request settles (a short delay; none if Umami is blocked). Adding `target="_blank"` would avoid it but is a visible behaviour change, so it is not done here. Decide whether the delay is acceptable.
3. **`/success` over-counting.** The event fires on every load of `/success`, including refreshes and direct visits. No de-duplication (e.g. `sessionStorage`) is added; compare against Netlify Forms submissions for ground truth.
4. **Undercounting.** Visitors using an ad/tracker blocker send nothing, so counts are a lower bound.
5. **Services "Discuss a project →" link** points to `#contact`, not `orangecatlabs.ca`, so it is untracked. Say if an internal-CTA event is wanted.
6. **Testimonial and GitHub links** are untracked by design (section 3). Say if either should be added.
7. **Deploy previews and local test traffic report to the production Umami site.** Previews are production builds without `data-domains`, so clicks on a preview (including reviewing this PR's preview) create events tagged with the `*.netlify.app` hostname. **Resolved in orchestrator review:** both layouts' script tags now set `data-domains="rafaelabreu.netlify.app"`, so the tracker only reports on the production hostname (`site` in `astro.config.mjs`). Previews, `localhost` and any other host load the script but send nothing. Every page with a tracked element uses `BaseLayout` or `ProjectLayout`, so coverage is unchanged in production. Consequence: events can no longer be checked on a deploy preview; verify them in production after deploy (item 1).
