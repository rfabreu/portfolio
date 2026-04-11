# P7: Dark/Light Theme Toggle — Design Spec

## Overview

Add a theme toggle to the portfolio allowing visitors to switch between dark and light modes. Dark remains the default. The light theme uses a hybrid palette (warm base, crisp surfaces) that preserves brand identity across both modes.

## Toggle Placement

- **Desktop:** SVG icon button in the nav bar, positioned between the nav links and the CTA buttons (RESUME / LET'S TALK). Sun icon in dark mode, moon icon in light mode.
- **Mobile:** Same icon button appears at the top of the hamburger menu overlay (MobileNav island).
- **No emoji.** Use inline SVG for the sun/moon icons.

## Color Palette

### Dark Theme (existing, default)

| Token | Value |
|-------|-------|
| base | #050508 |
| surface | rgba(255, 255, 255, 0.03) |
| surface-border | rgba(255, 255, 255, 0.06) |
| text-primary | #FFFFFF |
| text-secondary | #A0A0B0 |
| text-muted | #666666 |
| accent-indigo | #6366F1 |
| accent-purple | #8B5CF6 |
| accent-pink | #EC4899 |
| accent-indigo-light | #A5B4FC |
| accent-purple-light | #C4B5FD |
| accent-pink-light | #F9A8D4 |

### Light Theme (new, hybrid palette)

| Token | Value |
|-------|-------|
| base | #F8F7F4 |
| surface | #FFFFFF |
| surface-border | rgba(0, 0, 0, 0.06) |
| text-primary | #111111 |
| text-secondary | #555555 |
| text-muted | #999999 |
| accent-indigo | #4F46E5 |
| accent-purple | #7C3AED |
| accent-pink | #DB2777 |
| accent-indigo-light | #4F46E5 |
| accent-purple-light | #7C3AED |
| accent-pink-light | #DB2777 |

Light accent-light tokens collapse to match their base accent values (pastel tints don't read well on light backgrounds).

## Implementation Approach

### CSS Custom Properties

The `@theme` block in `global.css` defines the dark defaults. Light overrides are applied via `[data-theme="light"]` on the `<html>` element:

```css
html[data-theme="light"] {
  --color-base: #F8F7F4;
  --color-surface: #FFFFFF;
  /* ... all token overrides */
}
```

No component code duplication — all components already use theme tokens via Tailwind classes (`bg-base`, `text-text-primary`, etc.).

### Hardcoded Colors

~20 instances across React islands use literal hex values instead of theme tokens. These must be converted:

- **ChatWidget.tsx** — `bg-[#0a0a0f]`, `border-white/10`, `bg-white/5`, `text-white`, `bg-[#6366f1]`, `text-gray-300`, `text-gray-400`, `text-gray-500`
- **MobileNav.tsx** — `bg-[#050508]`, `bg-white`, `text-white`, `text-[#6366f1]`, `border-[#6366f1]/40`
- **GameBoard.tsx** — `text-white` (Play Again button), `text-red-400` (error text)
- **Hero.astro** — `border-white/15`, `text-gray-300`
- **ProjectLayout.astro** — `<style>` block has 6 hardcoded colors for prose-custom styles: `#ffffff`, `#a0a0b0`, `#6366f1`, `rgba(255,255,255,0.05)`
- **Skills.astro** — `bg-white/5` on skill tags
- **ProjectCard.astro** — `bg-white/5` on tech tags
- **success.astro** — `text-white` on CTA button
- **GradientDivider.astro** — hardcoded rgba accent values in gradient (faint on light backgrounds)

Strategy: Replace hardcoded values with theme-aware Tailwind classes (`bg-base`, `text-text-primary`, `border-surface-border`, etc.) or CSS custom properties where Tailwind classes aren't sufficient. For the ProjectLayout prose styles, convert hardcoded hex to `var(--color-*)` references.

### ThemeToggle Island

New React island: `frontend/src/islands/ThemeToggle.tsx`

- Reads current theme from `document.documentElement.dataset.theme`
- Toggles between `"dark"` and `"light"`
- Persists choice to `localStorage` key `"theme"`
- Renders sun SVG (dark mode) or moon SVG (light mode)
- Styled as a small bordered icon button matching the nav aesthetic

### Nav Integration

- **Nav.astro:** Add `<ThemeToggle client:load />` between the nav links `<div>` and the CTA buttons `<div>`
- **MobileNav.tsx:** Import and render `ThemeToggle` at the top of the mobile overlay, above the nav links

### ParticleHero Adaptation

The WebGL particle canvas inverts its color scheme in light mode:

- **Dark mode:** Bright indigo/purple particles, light connecting lines on dark background
- **Light mode:** Dark indigo/purple particles, dark connecting lines on light background
- The island reads `data-theme` from `<html>` at init and listens for changes (via `MutationObserver` on the `data-theme` attribute)
- **Color changes:** Particle colors are vertex attributes sourced from a `PALETTE` constant; the line color is a `u_color` uniform. Both need a dark-mode and light-mode variant.
- **Blending mode change:** The current additive blend (`gl.SRC_ALPHA, gl.ONE`) makes particles glow on dark backgrounds but washes out on light ones. Light mode must switch to standard alpha blending (`gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA`).
- Animation logic and interaction are identical across both modes.

### Theme Persistence & Flash Prevention

1. **Inline `<script>` in `<head>`** (in `BaseLayout.astro`, before any CSS or body content):
   - Reads `localStorage.getItem("theme")`
   - Falls back to `window.matchMedia("(prefers-color-scheme: light)").matches`
   - Sets `document.documentElement.dataset.theme` immediately
   - This prevents the dark→light flash on page load for light-mode users

2. **`<meta name="theme-color">`** updates dynamically:
   - Dark: `#050508`
   - Light: `#F8F7F4`
   - Updated by the ThemeToggle island on toggle and by the head script on load

### Transition Animation

- On toggle, add `.theme-transitioning` class to `<html>`
- CSS rule: `html.theme-transitioning * { transition: background-color 200ms ease, color 200ms ease, border-color 200ms ease; }`
- Remove class after 200ms via `setTimeout`
- This provides a smooth visual transition without permanently adding transition overhead to every element

## Accessibility

- Toggle button has `aria-label` ("Switch to light mode" / "Switch to dark mode")
- Respects `prefers-color-scheme` as initial default when no stored preference exists
- All color combinations pass WCAG AA contrast requirements
- `prefers-reduced-motion`: transition animation is skipped (instant swap)

## Files Changed

| File | Change |
|------|--------|
| `src/styles/global.css` | Add `[data-theme="light"]` overrides, `.theme-transitioning` rule |
| `src/islands/ThemeToggle.tsx` | New island — toggle button with sun/moon SVG |
| `src/components/Nav.astro` | Add ThemeToggle island |
| `src/islands/MobileNav.tsx` | Add theme toggle to overlay |
| `src/layouts/BaseLayout.astro` | Add head script for flash prevention, update meta theme-color |
| `src/islands/ParticleHero.tsx` | Read theme, invert particle colors in light mode |
| `src/islands/ChatWidget.tsx` | Replace hardcoded colors with theme tokens |
| `src/islands/GameBoard.tsx` | Replace hardcoded colors with theme tokens |
| `src/layouts/ProjectLayout.astro` | Convert prose-custom `<style>` hardcoded colors to CSS vars, add head script |
| `src/components/Skills.astro` | Replace `bg-white/5` with theme-aware class |
| `src/components/ProjectCard.astro` | Replace `bg-white/5` with theme-aware class |
| `src/components/GradientDivider.astro` | Adjust gradient opacity for light mode |
| `src/pages/success.astro` | Replace `text-white` with theme-aware class |

## Out of Scope

- Custom theme beyond dark/light (e.g., high contrast, sepia)
- Per-section theme overrides
- Theme-aware project cover images
- Animated toggle icon (morph sun→moon)
