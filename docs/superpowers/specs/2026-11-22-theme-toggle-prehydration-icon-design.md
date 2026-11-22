# Theme Toggle: Correct Icon From First Paint — Design Spec

> A light-theme visitor sees the dark-state toggle (sun icon) until React hydrates. This spec makes the icon correct on first paint using CSS alone, with no hydration mismatch. Owner approved this visible change (2026-11-22). It resolves open question 1 of `2026-11-18-hydration-mismatch-design.md`.

---

## 1. Overview & Goals

- The right icon (sun in dark theme, moon in light theme) is visible from first paint, before any island JavaScript runs.
- Server and client markup stay identical, so there are still zero hydration mismatches (guarded by `frontend/tests/e2e/hydration.spec.ts`).
- Toggling, `localStorage`, the `theme-color` meta, `theme-transitioning` and `prefers-reduced-motion` behaviour are untouched.

## 2. Root cause

`ThemeToggle.tsx` reads the theme with `useSyncExternalStore(subscribe, getTheme, getServerTheme)`. `getServerTheme` returns `'dark'` on purpose: React uses it for the server render and the hydration pass (PR #67). So the static HTML always holds the dark-state button (sun SVG, `aria-label="Switch to light mode"`).

The inline script in `frontend/src/layouts/BaseLayout.astro` (and the copy in `ProjectLayout.astro`) sets `html[data-theme]` before paint from `localStorage.theme`, else `prefers-color-scheme`. A light-theme visitor therefore gets a light page with a sun icon until `client:load` hydrates and re-renders (a visible flash of the wrong icon, longer on slow connections).

## 3. Design

### 3.1 Icons: render both, select with CSS

`ThemeToggle` always renders both SVGs, with the classes `theme-icon-sun` and `theme-icon-moon`, both `aria-hidden="true"` (the button's accessible name comes from `aria-label`). The markup no longer depends on the theme, so the icon part of the server HTML and the client render are identical in every theme.

`frontend/src/styles/global.css` shows one of them, keyed on `html[data-theme]`:

```css
html[data-theme='light'] .theme-icon-sun { display: none; }
html:not([data-theme='light']) .theme-icon-moon { display: none; }
```

- It lives in `global.css`, not a scoped style: the toggle is a React island, so Astro scoped styles do not apply to it. `global.css` already holds the same kind of rule (`html[data-theme='light'] a.oc-link`).
- **Default / unknown theme:** only the exact value `light` selects the moon. No `data-theme` (the inline script has not run, for example JavaScript disabled), `dark`, or any other value shows the sun. This matches `getTheme`'s `|| 'dark'` fallback and the token CSS (`html[data-theme="dark"], html:not([data-theme])`). `ThemeToggle` computes `isDark` as `theme !== 'light'`, so a stray `localStorage.theme` value cannot make the label and the icon disagree.
- `display: none` removes the hidden SVG from the accessibility tree too; `aria-hidden` makes that independent of CSS.

### 3.2 Accessible name: keep the state-driven `aria-label`

Options considered:

1. **Keep `aria-label={isDark ? 'Switch to light mode' : 'Switch to dark mode'}` (chosen).** The server and hydration pass render the dark label, then `useSyncExternalStore` re-renders with the real theme straight after hydration, so the label is correct once the island is live. React hydration does not mismatch (same snapshot as the server). Before hydration the button has no click handler, so it does nothing; a stale label on an inert control for a moment is not a functional problem, and it is the same window the previous code had.
2. A static label ("Toggle theme") would be correct from first paint and with no JS, but it drops the state-describing name, and `hydration.spec.ts` and screen-reader users rely on "Switch to dark/light mode". Rejected.
3. Two visually-hidden spans selected by CSS (`<span class="sr-only theme-label-sun">`) would be correct pre-hydration, but then the name depends on CSS and both strings exist in the DOM; more markup for an inert-control window. Rejected.

### 3.3 Other places

- `MobileNav.tsx` renders the same `ThemeToggle` component, inside the menu overlay that only mounts after a click (client-only), so it gets the change for free.
- No other theme toggle exists (`grep` for `dataset.theme` finds `ThemeToggle`, `ParticleHero` (read only) and the two layout scripts).

## 4. Finding: no-JavaScript visitors always get the dark theme

The theme is set by an inline script, so with JavaScript disabled `data-theme` is never set and the page is dark regardless of `prefers-color-scheme` or `localStorage`. The toggle then shows the sun, which is correct for that dark page. Honouring `prefers-color-scheme` in CSS for no-JS visitors is a separate change and out of scope.

So the test for "first paint, before hydration" blocks only the island bundles (`/_astro/*.js`) and leaves the inline script running: that is the real window a light visitor sees. A second test uses `javaScriptEnabled: false` and asserts the default (no `data-theme`) shows the sun.

## 5. Tests

`frontend/tests/e2e/hydration.spec.ts` gains a `theme toggle icon before hydration` block:

- Dark and light themes (stored in `localStorage`), island JS aborted: the expected icon is visible, the other is hidden, and `data-theme` matches. The no-stored-theme + OS light case is covered too.
- `javaScriptEnabled: false`: no `data-theme`, sun visible, moon hidden.
- After hydration, in both themes: the correct icon is visible, `aria-label` matches the theme, and clicking swaps both the icon and the label.
- The existing 27 hydration tests and behaviour tests keep running (0 hydration errors).

## 6. Verification

In `frontend/`: `npm ci`, `npx astro check`, `npm run build`, `npx vitest run`, then the E2E suite per CLAUDE.md. Plus screenshots of the nav toggle area (light and dark, island JS blocked, JS disabled).

## 7. Open questions

1. No-JS visitors always see the dark theme (section 4); the owner may want a CSS `prefers-color-scheme` fallback later.
