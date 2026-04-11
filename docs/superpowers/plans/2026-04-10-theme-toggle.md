# P7: Dark/Light Theme Toggle — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a dark/light theme toggle to the portfolio with a hybrid light palette, flash-free persistence, and theme-aware updates to all components.

**Architecture:** CSS custom properties define both palettes in `global.css`. A `data-theme` attribute on `<html>` drives the active palette. A head script prevents flash on load. A React island provides the toggle button. All hardcoded colors across components are replaced with theme tokens.

**Tech Stack:** Astro, React, Tailwind CSS v4, WebGL, CSS custom properties, localStorage

---

### Task 1: Light theme CSS tokens and transition rule

**Files:**
- Modify: `frontend/src/styles/global.css`

- [ ] **Step 1: Add light theme overrides after the `@theme` block**

Add this block after the closing `}` of the `@theme` block (after line 20) in `global.css`:

```css
html[data-theme="light"] {
  --color-base: #F8F7F4;
  --color-surface: #FFFFFF;
  --color-surface-border: rgba(0, 0, 0, 0.06);
  --color-text-primary: #111111;
  --color-text-secondary: #555555;
  --color-text-muted: #999999;
  --color-accent-indigo: #4F46E5;
  --color-accent-purple: #7C3AED;
  --color-accent-pink: #DB2777;
  --color-accent-indigo-light: #4F46E5;
  --color-accent-purple-light: #7C3AED;
  --color-accent-pink-light: #DB2777;
}
```

- [ ] **Step 2: Add transition animation rule**

Add after the light theme block:

```css
html.theme-transitioning * {
  transition: background-color 200ms ease, color 200ms ease, border-color 200ms ease !important;
}

@media (prefers-reduced-motion: reduce) {
  html.theme-transitioning * {
    transition: none !important;
  }
}
```

- [ ] **Step 3: Add light-mode gradient-divider override**

Update the existing `gradient-divider` utility to use CSS custom properties, then add a light-mode override. Replace the existing `@utility gradient-divider` block:

```css
@utility gradient-divider {
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(99, 102, 241, 0.3), rgba(139, 92, 246, 0.2), transparent);
}

html[data-theme="light"] .gradient-divider {
  background: linear-gradient(90deg, transparent, rgba(79, 70, 229, 0.4), rgba(124, 58, 237, 0.3), transparent);
}
```

- [ ] **Step 4: Verify dark theme is unchanged**

Run: `cd frontend && npm run build`
Expected: Build succeeds. No visual changes since no `data-theme` attribute is set yet (dark is the default).

- [ ] **Step 5: Commit**

```bash
git add frontend/src/styles/global.css
git commit -m "feat(theme): add light palette tokens and transition rule"
```

---

### Task 2: Flash prevention script in layouts

**Files:**
- Modify: `frontend/src/layouts/BaseLayout.astro`
- Modify: `frontend/src/layouts/ProjectLayout.astro`

- [ ] **Step 1: Add inline theme script to BaseLayout head**

In `frontend/src/layouts/BaseLayout.astro`, add this `<script is:inline>` block inside `<head>`, immediately before the closing `</head>` tag (after line 23):

```html
<script is:inline>
  (function() {
    var theme = localStorage.getItem('theme');
    if (!theme) {
      theme = window.matchMedia('(prefers-color-scheme: light)').matches ? 'light' : 'dark';
    }
    document.documentElement.dataset.theme = theme;
    var meta = document.querySelector('meta[name="theme-color"]');
    if (meta) meta.setAttribute('content', theme === 'light' ? '#F8F7F4' : '#050508');
  })();
</script>
```

Note: Using `is:inline` ensures Astro includes this script directly in the HTML without bundling. This must run before any rendering to prevent flash.

- [ ] **Step 2: Add the same script to ProjectLayout head**

In `frontend/src/layouts/ProjectLayout.astro`, add the identical `<script is:inline>` block inside `<head>`, immediately before `</head>` (after line 33).

- [ ] **Step 3: Verify build still passes**

Run: `cd frontend && npm run build`
Expected: Build succeeds, 6 pages built.

- [ ] **Step 4: Commit**

```bash
git add frontend/src/layouts/BaseLayout.astro frontend/src/layouts/ProjectLayout.astro
git commit -m "feat(theme): add flash prevention script to layouts"
```

---

### Task 3: ThemeToggle island

**Files:**
- Create: `frontend/src/islands/ThemeToggle.tsx`

- [ ] **Step 1: Create the ThemeToggle component**

Create `frontend/src/islands/ThemeToggle.tsx`:

```tsx
import { useState, useEffect } from 'react';

function getInitialTheme(): string {
  if (typeof document !== 'undefined') {
    return document.documentElement.dataset.theme || 'dark';
  }
  return 'dark';
}

export default function ThemeToggle() {
  const [theme, setTheme] = useState(getInitialTheme);

  useEffect(() => {
    const observer = new MutationObserver(() => {
      const current = document.documentElement.dataset.theme || 'dark';
      setTheme(current);
    });
    observer.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ['data-theme'],
    });
    return () => observer.disconnect();
  }, []);

  function toggle() {
    const next = theme === 'dark' ? 'light' : 'dark';

    document.documentElement.classList.add('theme-transitioning');
    document.documentElement.dataset.theme = next;
    localStorage.setItem('theme', next);

    const meta = document.querySelector('meta[name="theme-color"]');
    if (meta) meta.setAttribute('content', next === 'light' ? '#F8F7F4' : '#050508');

    setTimeout(() => {
      document.documentElement.classList.remove('theme-transitioning');
    }, 200);
  }

  const isDark = theme === 'dark';

  return (
    <button
      onClick={toggle}
      className="w-8 h-8 rounded-btn border border-surface-border flex items-center justify-center text-text-muted hover:text-accent-indigo hover:border-accent-indigo/40 transition-colors duration-200"
      aria-label={isDark ? 'Switch to light mode' : 'Switch to dark mode'}
    >
      {isDark ? (
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <circle cx="12" cy="12" r="5" />
          <line x1="12" y1="1" x2="12" y2="3" />
          <line x1="12" y1="21" x2="12" y2="23" />
          <line x1="4.22" y1="4.22" x2="5.64" y2="5.64" />
          <line x1="18.36" y1="18.36" x2="19.78" y2="19.78" />
          <line x1="1" y1="12" x2="3" y2="12" />
          <line x1="21" y1="12" x2="23" y2="12" />
          <line x1="4.22" y1="19.78" x2="5.64" y2="18.36" />
          <line x1="18.36" y1="5.64" x2="19.78" y2="4.22" />
        </svg>
      ) : (
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
        </svg>
      )}
    </button>
  );
}
```

- [ ] **Step 2: Verify build**

Run: `cd frontend && npm run build`
Expected: Build succeeds. The island exists but isn't mounted yet.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/islands/ThemeToggle.tsx
git commit -m "feat(theme): add ThemeToggle island with sun/moon SVG"
```

---

### Task 4: Nav integration (desktop + mobile)

**Files:**
- Modify: `frontend/src/components/Nav.astro`
- Modify: `frontend/src/islands/MobileNav.tsx`

- [ ] **Step 1: Add ThemeToggle to desktop Nav**

In `frontend/src/components/Nav.astro`, add the import in the frontmatter:

```astro
import ThemeToggle from '../islands/ThemeToggle.tsx';
```

Then add the toggle between the nav links div and the CTA buttons div. Replace the existing `<div class="hidden md:flex items-center gap-3">` block (the CTA buttons wrapper, starting at line 32) with:

```html
    <div class="hidden md:flex items-center gap-3">
      <ThemeToggle client:load />
      <a
        href={resumePath}
        target="_blank"
        rel="noopener noreferrer"
        class="px-4 py-1.5 border border-accent-indigo/40 rounded-btn text-accent-indigo text-xs tracking-wide hover:bg-accent-indigo/10 transition-colors duration-200"
      >
        RESUME
      </a>
      <a
        href="#contact"
        class="px-4 py-1.5 bg-accent-indigo rounded-btn text-white text-xs tracking-wide hover:bg-accent-indigo/90 transition-colors duration-200"
      >
        LET'S TALK
      </a>
    </div>
```

- [ ] **Step 2: Add ThemeToggle to MobileNav**

In `frontend/src/islands/MobileNav.tsx`, add the import at the top:

```tsx
import ThemeToggle from './ThemeToggle';
```

Then add the toggle at the top of the mobile overlay. Replace the overlay `<div>` content (the `div` with `className="fixed inset-0 z-40 bg-[#050508]..."`) with:

```tsx
        <div className="fixed inset-0 z-40 bg-base flex flex-col items-center justify-center gap-8">
          <div className="absolute top-5 right-14">
            <ThemeToggle />
          </div>
          {navItems.map((item) => (
            <a
              key={item.label}
              href={item.href}
              onClick={() => setIsOpen(false)}
              className="text-text-primary text-2xl font-bold hover:text-accent-indigo transition-colors"
              style={{ letterSpacing: '-0.06em' }}
            >
              {item.label}
            </a>
          ))}
          <div className="flex gap-4 mt-4">
            <a
              href={resumePath}
              target="_blank"
              rel="noopener noreferrer"
              className="px-6 py-2 border border-accent-indigo/40 rounded-lg text-accent-indigo text-sm"
            >
              RESUME
            </a>
            <a
              href="#contact"
              onClick={() => setIsOpen(false)}
              className="px-6 py-2 bg-accent-indigo rounded-lg text-white text-sm"
            >
              LET'S TALK
            </a>
          </div>
        </div>
```

Note: This also fixes the hardcoded colors in MobileNav — `bg-[#050508]` → `bg-base`, `text-white` → `text-text-primary`, `text-[#6366f1]` → `text-accent-indigo`, `border-[#6366f1]/40` → `border-accent-indigo/40`.

- [ ] **Step 3: Fix MobileNav hamburger icon colors**

In `MobileNav.tsx`, update the two hamburger line spans to use theme tokens. Replace both instances of `bg-white` with `bg-text-primary`:

```tsx
        <span
          className={`block w-5 h-0.5 bg-text-primary transition-all duration-300 ${
            isOpen ? 'rotate-45 translate-y-1' : ''
          }`}
        />
        <span
          className={`block w-5 h-0.5 bg-text-primary transition-all duration-300 ${
            isOpen ? '-rotate-45 -translate-y-1' : ''
          }`}
        />
```

- [ ] **Step 4: Verify dev server**

Run: `cd frontend && npm run dev`
Open `http://localhost:4321` in a browser. The toggle should appear in the nav bar. Clicking it should switch themes. Verify:
- Toggle icon changes between sun and moon
- Page colors change
- Refreshing the page preserves the chosen theme
- Mobile hamburger menu shows the toggle

- [ ] **Step 5: Commit**

```bash
git add frontend/src/components/Nav.astro frontend/src/islands/MobileNav.tsx
git commit -m "feat(theme): integrate toggle into desktop and mobile nav"
```

---

### Task 5: Fix hardcoded colors in Hero and Astro components

**Files:**
- Modify: `frontend/src/components/Hero.astro`
- Modify: `frontend/src/components/Skills.astro`
- Modify: `frontend/src/components/ProjectCard.astro`
- Modify: `frontend/src/pages/success.astro`

- [ ] **Step 1: Fix Hero.astro**

In `frontend/src/components/Hero.astro`, update the "Ask my AI" button (line 44-45). Replace:

```html
        class="px-6 py-2.5 border border-white/15 rounded-btn text-gray-300 text-sm font-medium hover:border-accent-indigo/40 hover:text-accent-indigo transition-colors duration-200"
```

With:

```html
        class="px-6 py-2.5 border border-surface-border rounded-btn text-text-secondary text-sm font-medium hover:border-accent-indigo/40 hover:text-accent-indigo transition-colors duration-200"
```

- [ ] **Step 2: Fix Skills.astro**

In `frontend/src/components/Skills.astro`, replace `bg-white/5` on line 45. Change:

```html
              <span class="font-mono text-xs text-text-secondary bg-white/5 px-2.5 py-1 rounded-tag">
```

To:

```html
              <span class="font-mono text-xs text-text-secondary bg-surface px-2.5 py-1 rounded-tag">
```

- [ ] **Step 3: Fix ProjectCard.astro**

In `frontend/src/components/ProjectCard.astro`, replace `bg-white/5` on line 18. Change:

```html
      <span class="font-mono text-[10px] px-2 py-0.5 rounded-tag bg-white/5 text-text-secondary">
```

To:

```html
      <span class="font-mono text-[10px] px-2 py-0.5 rounded-tag bg-surface text-text-secondary">
```

- [ ] **Step 4: Fix success.astro**

In `frontend/src/pages/success.astro`, the `text-white` on the CTA button (line 18) can stay — it's on a `bg-accent-indigo` button where white text is correct in both themes. No change needed.

- [ ] **Step 5: Verify build**

Run: `cd frontend && npm run build`
Expected: Build succeeds.

- [ ] **Step 6: Commit**

```bash
git add frontend/src/components/Hero.astro frontend/src/components/Skills.astro frontend/src/components/ProjectCard.astro
git commit -m "fix(theme): replace hardcoded colors in Hero, Skills, ProjectCard"
```

---

### Task 6: Fix hardcoded colors in ChatWidget

**Files:**
- Modify: `frontend/src/islands/ChatWidget.tsx`

- [ ] **Step 1: Update ChatWidget with theme-aware classes**

Replace the entire return statement in `ChatWidget.tsx` (lines 70-150) with:

```tsx
  return (
    <>
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="fixed bottom-5 right-5 z-50 w-[52px] h-[52px] rounded-[14px] flex items-center justify-center text-white text-xl shadow-lg transition-transform hover:scale-105"
        style={{ background: 'linear-gradient(135deg, #6366f1, #8b5cf6)' }}
        aria-label={isOpen ? 'Close chat' : 'Open chat'}
      >
        {isOpen ? '✕' : '💬'}
      </button>

      {isOpen && (
        <div className="fixed bottom-20 right-5 z-50 w-[360px] max-h-[500px] bg-base border border-surface-border rounded-xl shadow-2xl flex flex-col overflow-hidden">
          <div className="px-4 py-3 border-b border-surface-border">
            <div className="text-text-primary text-sm font-bold">Ask Rafael's AI</div>
            <div className="text-text-muted text-xs mt-0.5">Ask about experience, projects, or skills</div>
          </div>

          <div className="flex-1 overflow-y-auto p-4 space-y-3 min-h-[200px] max-h-[340px]">
            {messages.length === 0 && (
              <div className="text-text-muted text-sm">
                Hey! I'm Rafael's AI assistant. Ask me about his experience, projects, or tech stack.
              </div>
            )}

            {messages.map((msg, i) => (
              <div key={i} className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}>
                <div
                  className={`max-w-[85%] px-3 py-2 rounded-lg text-sm ${
                    msg.role === 'user'
                      ? 'bg-accent-indigo text-white'
                      : 'bg-surface text-text-secondary'
                  }`}
                >
                  {msg.content}
                </div>
              </div>
            ))}

            {isLoading && (
              <div className="flex justify-start">
                <div className="bg-surface text-text-muted px-3 py-2 rounded-lg text-sm">
                  Thinking...
                </div>
              </div>
            )}

            {error && (
              <div className="text-red-400 text-xs text-center">
                {error}{' '}
                <a href="#contact" onClick={() => setIsOpen(false)} className="underline">
                  Leave a message instead
                </a>
              </div>
            )}

            <div ref={messagesEndRef} />
          </div>

          <form onSubmit={sendMessage} className="p-3 border-t border-surface-border flex gap-2">
            <input
              ref={inputRef}
              type="text"
              value={input}
              onChange={(e) => setInput(e.target.value)}
              placeholder="Ask something..."
              className="flex-1 bg-surface border border-surface-border rounded-lg px-3 py-2 text-sm text-text-primary placeholder-text-muted outline-none focus:border-accent-indigo/50"
              disabled={isLoading}
            />
            <button
              type="submit"
              disabled={isLoading || !input.trim()}
              className="px-3 py-2 bg-accent-indigo rounded-lg text-white text-sm disabled:opacity-40 hover:bg-accent-indigo/90 transition-colors"
            >
              →
            </button>
          </form>
        </div>
      )}
    </>
  );
```

Key changes:
- `bg-[#0a0a0f]` → `bg-base`
- `border-white/10` → `border-surface-border`
- `text-white` (header) → `text-text-primary`
- `text-gray-500` → `text-text-muted`
- `bg-[#6366f1]` → `bg-accent-indigo`
- `bg-white/5` → `bg-surface`
- `text-gray-300` → `text-text-secondary`
- `text-gray-400` → `text-text-muted`
- `border-[#6366f1]/50` → `border-accent-indigo/50`
- `placeholder-gray-500` → `placeholder-text-muted`
- Kept `text-white` on the gradient chat FAB button and user message bubble (intentional — white text on colored backgrounds works in both themes)
- Kept `text-white` on send button (on accent-indigo background)

- [ ] **Step 2: Verify build**

Run: `cd frontend && npm run build`
Expected: Build succeeds.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/islands/ChatWidget.tsx
git commit -m "fix(theme): replace hardcoded colors in ChatWidget"
```

---

### Task 7: Verify GameBoard (no changes needed)

**Files:**
- Verify: `frontend/src/islands/GameBoard.tsx`

- [ ] **Step 1: Confirm GameBoard is already theme-aware**

GameBoard already uses theme tokens throughout: `bg-surface`, `border-surface-border`, `text-text-primary`, `text-accent-indigo`, `text-text-muted`. The `text-white` on the "Play Again" button is correct — white text on `bg-accent-indigo` works in both themes. The `text-red-400` error text is acceptable in both themes.

No changes or commits needed — skip to Task 8.

---

### Task 8: Fix hardcoded colors in ProjectLayout

**Files:**
- Modify: `frontend/src/layouts/ProjectLayout.astro`

- [ ] **Step 1: Convert prose-custom styles to CSS custom properties**

In `frontend/src/layouts/ProjectLayout.astro`, replace the entire `<style>` block (lines 80-121) with:

```html
    <style>
      .prose-custom :global(h2) {
        color: var(--color-text-primary);
        font-size: 1.5rem;
        font-weight: 800;
        letter-spacing: -0.04em;
        margin-top: 2.5rem;
        margin-bottom: 1rem;
      }
      .prose-custom :global(p) {
        color: var(--color-text-secondary);
        line-height: 1.7;
        margin-bottom: 1rem;
      }
      .prose-custom :global(ul) {
        color: var(--color-text-secondary);
        list-style: none;
        padding-left: 0;
      }
      .prose-custom :global(li) {
        padding-left: 1.25rem;
        position: relative;
        margin-bottom: 0.5rem;
      }
      .prose-custom :global(li::before) {
        content: '→';
        position: absolute;
        left: 0;
        color: var(--color-accent-indigo);
      }
      .prose-custom :global(strong) {
        color: var(--color-text-primary);
        font-weight: 600;
      }
      .prose-custom :global(code) {
        font-family: 'Courier New', monospace;
        background: var(--color-surface);
        padding: 0.15rem 0.4rem;
        border-radius: 4px;
        font-size: 0.875rem;
      }
    </style>
```

- [ ] **Step 2: Also fix the "Live Demo" button**

On line 63, `text-white` on `bg-accent-indigo` is correct — keep as-is. No change needed.

- [ ] **Step 3: Verify build**

Run: `cd frontend && npm run build`
Expected: Build succeeds.

- [ ] **Step 4: Commit**

```bash
git add frontend/src/layouts/ProjectLayout.astro
git commit -m "fix(theme): convert ProjectLayout prose styles to CSS vars"
```

---

### Task 9: ParticleHero theme adaptation

**Files:**
- Modify: `frontend/src/islands/ParticleHero.tsx`

- [ ] **Step 1: Add dark and light palette constants**

In `frontend/src/islands/ParticleHero.tsx`, replace the existing `PALETTE` constant (lines 104-108) with:

```tsx
const DARK_PALETTE = [
  [0.389, 0.4, 0.945],   // #6366f1 indigo
  [0.545, 0.361, 0.965],  // #8b5cf6 purple
  [0.925, 0.286, 0.600],  // #ec4899 pink
] as const;

const LIGHT_PALETTE = [
  [0.310, 0.275, 0.898],  // #4F46E5 indigo (deeper)
  [0.486, 0.227, 0.929],  // #7C3AED purple (deeper)
  [0.859, 0.153, 0.467],  // #DB2777 pink (deeper)
] as const;

const DARK_LINE_COLOR: [number, number, number] = [0.389, 0.4, 0.945];
const LIGHT_LINE_COLOR: [number, number, number] = [0.310, 0.275, 0.898];

function isLightTheme(): boolean {
  return document.documentElement.dataset.theme === 'light';
}
```

- [ ] **Step 2: Update initParticles to accept a palette**

Replace the `initParticles` function signature and its palette reference:

```tsx
function initParticles(width: number, height: number, isMobile: boolean, palette: typeof DARK_PALETTE): Particle[] {
  const count = Math.max(60, Math.min(120, Math.floor(width / 15)));
  const speedMultiplier = isMobile ? 1.5 : 1.0;
  const particles: Particle[] = [];
  for (let i = 0; i < count; i++) {
    const color = palette[Math.floor(Math.random() * palette.length)];
    particles.push({
      x: Math.random() * width,
      y: Math.random() * height,
      vx: (Math.random() - 0.5) * 0.8 * speedMultiplier,
      vy: (Math.random() - 0.5) * 0.8 * speedMultiplier,
      r: color[0],
      g: color[1],
      b: color[2],
      baseSize: 2.5 + Math.random() * 2.5,
    });
  }
  return particles;
}
```

- [ ] **Step 3: Add theme change listener and update the resize/frame functions**

Inside the `useEffect`, after the line `let paused = false;` (line 182), add:

```tsx
    let currentLight = isLightTheme();
```

Update the `resize` function to use the current theme palette:

```tsx
    function resize() {
      if (!canvas || !gl) return;
      const rect = canvas.getBoundingClientRect();
      w = rect.width * dpr;
      h = rect.height * dpr;
      canvas.width = w;
      canvas.height = h;
      gl.viewport(0, 0, w, h);
      currentLight = isLightTheme();
      particles = initParticles(w, h, isMobile, currentLight ? LIGHT_PALETTE : DARK_PALETTE);
    }
```

- [ ] **Step 4: Update the frame function for theme-aware blending and line color**

In the `frame` function, update the blend mode and line color. Replace the render section starting at `// Render` (line 325):

```tsx
      // Render
      gl.clearColor(0, 0, 0, 0);
      gl.clear(gl.COLOR_BUFFER_BIT);
      gl.enable(gl.BLEND);
      if (currentLight) {
        gl.blendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA);
      } else {
        gl.blendFunc(gl.SRC_ALPHA, gl.ONE);
      }
      gl.disable(gl.DEPTH_TEST);
```

Then update the line color uniform. Replace the line `gl.uniform3f(lLoc.color, 0.389, 0.4, 0.945);` with:

```tsx
        const lineColor = currentLight ? LIGHT_LINE_COLOR : DARK_LINE_COLOR;
        gl.uniform3f(lLoc.color, lineColor[0], lineColor[1], lineColor[2]);
```

- [ ] **Step 5: Add MutationObserver for theme changes**

After the event listener setup (after line 403), add:

```tsx
    // Listen for theme changes
    const themeObserver = new MutationObserver(() => {
      const nowLight = isLightTheme();
      if (nowLight !== currentLight) {
        currentLight = nowLight;
        const palette = currentLight ? LIGHT_PALETTE : DARK_PALETTE;
        for (let i = 0; i < particles.length; i++) {
          const color = palette[Math.floor(Math.random() * palette.length)];
          particles[i].r = color[0];
          particles[i].g = color[1];
          particles[i].b = color[2];
        }
      }
    });
    themeObserver.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ['data-theme'],
    });
```

And add cleanup in the return function, before the WebGL resource cleanup:

```tsx
      themeObserver.disconnect();
```

- [ ] **Step 6: Verify with dev server**

Run: `cd frontend && npm run dev`
Open `http://localhost:4321`. Toggle themes and verify:
- Particles change color when switching themes
- Lines change color
- No washed-out particles in light mode (blending mode switch)
- Animation continues smoothly during toggle

- [ ] **Step 7: Commit**

```bash
git add frontend/src/islands/ParticleHero.tsx
git commit -m "feat(theme): adapt ParticleHero for light mode"
```

---

### Task 10: Full integration test

**Files:** None (verification only)

- [ ] **Step 1: Build and verify**

Run: `cd frontend && npx astro check && npm run build`
Expected: Type checking passes, build succeeds with 6 pages.

- [ ] **Step 2: Run Go API tests**

Run: `cd api && go vet ./... && go test ./... -v`
Expected: All tests pass. (API code is unchanged but verify nothing is broken.)

- [ ] **Step 3: Manual verification checklist**

Run: `cd frontend && npm run dev`

Test in browser at `http://localhost:4321`:

1. **Default state:** Page loads in dark mode (or matches system preference if set to light)
2. **Toggle:** Click the sun icon → page smoothly transitions to light mode. Icon changes to moon.
3. **Persistence:** Refresh the page → stays in light mode. No flash of dark.
4. **Toggle back:** Click moon icon → transitions back to dark mode. Refresh → stays dark.
5. **Nav (desktop):** Toggle button visible between nav links and CTA buttons
6. **Nav (mobile):** Open hamburger → toggle appears at top of overlay
7. **Hero:** Particles visible and colored appropriately in both themes
8. **Chat widget:** Open chat → colors match current theme
9. **Skills section:** Skill tags have proper background in both themes
10. **Contact form:** Input fields and borders look correct in both themes
11. **Project cards:** Tags and borders correct in both themes
12. **Playground page:** Game board correct in both themes
13. **Project detail page:** Navigate to a project → prose content colors correct in both themes
14. **GradientDivider:** Visible in both themes

- [ ] **Step 4: Commit all remaining changes (if any fixups needed)**

```bash
git add -A
git commit -m "fix(theme): integration fixups"
```

(Skip this step if no fixups were needed.)
