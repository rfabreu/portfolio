# Track Brief CTA Clicks — Design Spec

> Adds one Umami event, `brief-cta-click`, for in-page links that lead to the contact form (`#contact` or `/#contact`). Resolves open question 5 of `2026-11-17-analytics-funnel-events-design.md` (the untracked Services link) and reverses its "intentionally untracked `#contact` anchors" line for the placements below.

---

## 1. Overview & Objectives

The funnel is visit → click → `contact-form-submit` → `contact-form-success`. The step between "visit" and "submit" is invisible for visitors who reach the form through an in-page link. Tracking those links shows which placement sends people to the Project Brief form.

Goals:
1. One event, `brief-cta-click`, with a `location` property naming the placement.
2. Same-tab scrolling is never delayed.
3. Change nothing visible: no copy, layout, styling, `href`, `target` or `rel`. Send no PII.

Non-goals: new copy or links, Playwright funnel specs, the Hero "Work with OrangeCat Labs" CTA (external, already `orangecat-click`).

---

## 2. Links Found

`grep -rnE '#contact' frontend/src` finds these placements (the Nav item is defined once in a data array and rendered twice):

| `location` | Element | File | Tracked |
|-----------|---------|------|---------|
| `nav` | Desktop nav "Contact" (`/#contact`) | `src/components/Nav.astro` | Yes |
| `mobile-nav` | Mobile menu "Contact" (same `navItems` entry) | `src/islands/MobileNav.tsx` | Yes |
| `services` | "Have a project in mind? Submit a brief below ↓" | `src/components/Services.astro` | Yes |
| `chat-error` | "Leave a message instead" (shown when a chat request fails) | `src/islands/ChatWidget.tsx` | Yes |
| `resume-modal` | "contact form" in the resume modal footer | `src/components/ResumeModal.astro` | Yes |
| `resume-page-footer` | "contact form" in the `/resume` footer | `src/pages/resume.astro` | Yes |

Justification: each is an explicit invitation to submit a brief, so all are funnel CTAs. Nav "Contact" is the most-visible entry and is tracked separately for desktop and mobile so the two surfaces can be compared. `chat-error` is a recovery path: it shows how often a failed chat still converts. The resume links show how many resume readers move to the form. No `#contact` link was left out. Links to `#services`, `#projects` and so on are navigation, not funnel, and stay untracked.

Content files (`src/content/**`) contain no `#contact` links (grep checked).

---

## 3. Event

| Event | Properties | Fires when | Mechanism |
|-------|-----------|-----------|-----------|
| `brief-cta-click` | `location` ∈ `nav`, `mobile-nav`, `services`, `chat-error`, `resume-modal`, `resume-page-footer` | Any link in section 2 is clicked | Imperative |

Naming follows the existing kebab-case scheme (`view-projects-click`). `location` values are fixed strings in source.

---

## 4. Implementation

Same-tab links must not use `data-umami-event`: the tracker calls `preventDefault` and navigates after its request settles, which would delay the scroll (spec of 2026-11-17, section 2). All six use the `view-projects-click` pattern: a click listener that never calls `preventDefault`, behind `typeof umami !== 'undefined'`.

- **Astro links** get a hook attribute, `data-brief-cta="<location>"` (a non-visual attribute), and the owning component gets a `<script>` that adds a click listener to that element and calls `umami.track('brief-cta-click', { location: '<location>' })`. `Nav.astro` sets the hook only on the entry whose `href` is `/#contact`. `ResumeModal.astro` gets a second `<script>` block rather than extending the existing one, whose handler legitimately calls `preventDefault` for the modal trigger.
- **MobileNav** replaces the inline `() => setIsOpen(false)` with `handleNavClick(item.href)`, which tracks when `href === '/#contact'` and then always closes the menu, as before.
- **ChatWidget** replaces the inline close handler on the error link with `trackBriefFromChat()`, which tracks and then closes the chat, as before.

`/#contact` from `/resume`, `/playground`, `/success` and case studies is a full page load. Umami's tracker sends events with `fetch(..., { keepalive: true })`, so the event survives that navigation (same assumption as `contact-form-submit`; open question 2).

---

## 5. Privacy

Only the fixed `location` string is sent. No user input, URL, form value or message text is read. The existing scan that allows only `location` and `project` event data keys is unaffected, since this event uses imperative `umami.track` data rather than `data-umami-event-*` attributes.

---

## 6. Testing & Verification

Vitest, extending `tests/analytics-funnel-events.test.ts` and `tests/analytics-funnel-mobile-nav.test.ts`:
- Services and Nav are rendered through the Astro container: the link carries `data-brief-cta` with the expected location, no `data-umami-event`, and only the Contact nav link carries the hook.
- Each Astro owner's script contains the selector and the guarded `umami.track('brief-cta-click', { location })` call and no `preventDefault`.
- MobileNav is called with `useState` mocked: clicking Contact calls `umami.track('brief-cta-click', { location: 'mobile-nav' })` and `setIsOpen(false)`; other links only close the menu; no throw when `umami` is absent; markup has no extra attributes.
- ChatWidget's link calls `trackBriefFromChat` (source check; the function tracks `chat-error` then closes).
- Source scan: every `#contact` anchor in `src/` is one of the known files (a new untracked one fails the test); the global guard test now lists `brief-cta-click`.

Build verification (`frontend/`): `npm ci`, `npx astro check`, `npm run build`, `npm test`, `npm run build:pdf && npm run test:e2e`; inspect `dist/` for the hooks and bundled handlers.

Not verifiable here: production analytics (open question 1).

---

## 7. Open Questions

1. **Production verification.** After deploy, click each placement once and check Umami → Events → `brief-cta-click` → `location`. Previews report nothing (`data-domains`).
2. **Cross-page clicks.** On non-home pages `/#contact` navigates away. The event relies on the tracker's `keepalive` request; the owner should confirm `nav` and `resume-page-footer` events from non-home pages arrive.
3. **Hook attribute naming.** `data-brief-cta` is a new non-visual hook, like `data-view-projects`. A single shared script (delegation) was rejected to keep each component self-contained, at the cost of four small scripts.
4. **Not distinguished:** the Nav item is one entry, so `nav` counts clicks from every page it is on; add page context only if the owner wants it (it would add a property).
5. **`chat-error` volume** is expected to be low, since it only shows after a failed request.
6. **`resume-modal` on the home page.** The link is a same-page hash change while the `<dialog>` is still modal, so the form may stay hidden behind it. The click is counted regardless; read `resume-modal` counts as intent, not arrival. Not changed (would be a visible behaviour change).
