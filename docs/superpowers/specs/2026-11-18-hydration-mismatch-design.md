# Fix React Hydration Mismatch (#418) — Design Spec

> Production logs a minified `React error #418` on pages with the nav. The cause is `ThemeToggle`, and only for light-theme visitors. One-file island fix plus an E2E regression test. No human was available for brainstorming; uncertainties are in section 7.

---

## 1. Overview & Goals

Remove every server/client markup mismatch on the site's React islands, with no visible change: the theme must stay correct on first paint, the toggle must keep working, and `prefers-reduced-motion` and the analytics hooks must be untouched.

## 2. Root-cause evidence

Reproduced with a Playwright script (not committed) over 5 pages × 3 theme setups × 2 `prefers-reduced-motion` values = 30 loads, in the dev server (unminified message) and in `npm run build` + `astro preview` (minified).
Pages: `/`, `/playground/`, `/resume/`, `/project/freckles-design/`, `/success/`. Theme setups: `localStorage.theme=dark`, `localStorage.theme=light`, no stored theme with OS `prefers-color-scheme: light`.

| Result | Dev | Production |
|---|---|---|
| Dark theme (stored) | 0 mismatches | 0 mismatches |
| Light theme (stored or OS-derived) | 16 mismatches, all in `<ThemeToggle>` | 16 × `Minified React error #418` |
| `/resume/` (any theme) | 0 (the page has no `ThemeToggle`) | 0 |

The dev message names the component and the diff:

```
<ThemeToggle>
  <button ...
+   aria-label="Switch to dark mode"     (client)
-                                        (server)
```

Cause, in `frontend/src/islands/ThemeToggle.tsx` before the fix: `useState(getInitialTheme)` reads `document.documentElement.dataset.theme` when `document` exists. The server has no `document`, so it always renders the `'dark'` branch (sun icon, "Switch to light mode"). `client:load` hydrates against that HTML. The inline script in `BaseLayout.astro` has already set `data-theme="light"` for a light-theme visitor, so the first client render is the moon icon and "Switch to dark mode": different markup, so React discards the server tree and logs #418. Dark visitors match by coincidence, which is why the error only shows for some visitors.

Ruled out by the same runs (zero mismatches in dark theme, and the only reported component in light theme is `ThemeToggle`): `ParticleHero` (canvas only; reads theme, `matchMedia` and `Math.random` inside `useEffect`), `MobileNav` (closed on first render), `ChatWidget` (closed on first render), `GameBoard` (deterministic initial state). The only other console output was headless-GPU driver warnings from WebGL on `/`.

## 3. Design

`ThemeToggle` reads the theme with `useSyncExternalStore`:

- `subscribe`: a `MutationObserver` on `<html>`'s `data-theme` (the same observer the old `useEffect` ran).
- `getSnapshot`: `document.documentElement.dataset.theme || 'dark'` (client).
- `getServerSnapshot`: `'dark'`. React uses it for the server render **and** for the hydration pass, so the markup matches the server HTML, then re-renders with the real value straight after hydration.

This is the pattern React documents for external mutable state with SSR. Alternatives rejected:

- `useState('dark')` + `useEffect` sync: also correct and equally one extra render after hydration, but the subscription and the state copy are hand-written; `useSyncExternalStore` has both built in, with no state to keep in sync.
- `suppressHydrationWarning` on the button: hides the error but leaves the server markup wrong and only covers one attribute, not the swapped SVG child.
- `client:only="react"`: loses the server-rendered button (layout shift in the nav) for no benefit.

Behaviour is otherwise unchanged: `toggle()`, `localStorage`, `theme-color` meta and the `theme-transitioning` class are untouched. `ParticleHero` keeps its own `MutationObserver` for the theme.

## 4. Behaviour kept, and one known pre-existing quirk

Before and after the fix (production build, JS for islands blocked to capture the server-rendered state, then hydrated):

- Page theme on first frame: a single `data-theme | background-color` state per run, sampled every animation frame (dark `rgb(5, 5, 8)`, light `rgb(248, 247, 244)`). No flash, before or after.
- Nav bar screenshots: byte-identical before and after, in dark and light, pre-hydration and hydrated.
- Pre-existing and unchanged: before hydration, a light-theme visitor's toggle shows the server-rendered dark-state icon and label for a moment. Fixing it (CSS-driven icons) is a separate, visible change; see section 7.

## 5. Tests

- `frontend/tests/e2e/hydration.spec.ts` (Playwright, runs in the existing `Frontend E2E` job):
  - The five pages × stored dark/light theme × `reducedMotion` `no-preference`/`reduce` = 20 tests, plus the five pages with no stored theme and OS `prefers-color-scheme: light` = 5 tests. Each collects `console` errors and `pageerror`s and fails on `Minified React error #418`, `Hydration`, `did not match` or `didn't match`.
  - Waiting for hydration: Astro drops an island's `ssr` attribute when `hydrateRoot` is scheduled, and React reports a mismatch when that hydration commits, so the helper waits for no `astro-island[ssr]` and then for an idle callback, to avoid asserting before a late error arrives (found in review).
  - Two behaviour tests: a persisted light theme gets the "Switch to dark mode" label after hydration and clicking flips `data-theme` and the label, on desktop and in the 390px mobile menu overlay.
  - `/resume/` has no React islands (it does not render `Nav`), so its tests pass trivially; they are kept so a future island there is covered.
- **Proven to fail before the fix:** run against a build with `main`'s `ThemeToggle.tsx`, 12 of the 27 fail (every light-theme page except `/resume/`, stored or OS-derived) and 15 pass. After the fix, 27 pass.

## 6. Verification

In `frontend/`: `npm ci`, `npx astro check`, `npm run build`, `npm test`, `npm run build:pdf`, `npm run test:e2e`. Plus the 30-load repro script: 0 mismatches in dev and production after the fix.

## 7. Open questions

1. **Pre-hydration toggle icon in light theme.** As above, the toggle briefly shows the dark-state icon. Rendering both icons and selecting with CSS on `html[data-theme]` would remove it, but it changes markup and styling; left for the owner to request.
2. **Other mismatch sources not reachable by this script**: states behind interaction (open mobile menu, open chat, game moves) render client-side only after a click, so they cannot cause hydration mismatches; they were not separately exercised.
3. **Production confirmation.** After deploy, the #418 error should stop appearing for light-theme visitors; the owner should confirm in whichever tool collects the logs (not visible from this repo).
