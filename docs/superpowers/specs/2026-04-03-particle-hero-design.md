# P1: Interactive Particle Hero — Design Spec

> Replaces the static gradient orbs in the hero section with a WebGL particle mesh that reacts to cursor movement, forming an interconnected neural network pattern.

## Architecture

### Component Structure

```
Hero.astro (text content — name, tagline, tech tags, CTAs)
  └── ParticleHero.tsx (React island, client:load)
        ├── Vertex shader (particle positions → clip space)
        ├── Fragment shader (soft glow circles, palette colors)
        └── JS: particle simulation, mouse tracking, connection lines
```

### Files Changed

- **New:** `frontend/src/islands/ParticleHero.tsx` — React island owning a full-bleed `<canvas>` and all WebGL logic
- **Modified:** `frontend/src/components/Hero.astro` — remove the three static orb `<div>`s, add `<ParticleHero client:load />` positioned absolutely behind hero text

### Integration

The canvas renders behind the hero text via absolute positioning and `z-index`. Hero.astro continues to own all text content (static Astro). The React island manages animation state internally — no props from Astro required.

## Particle Simulation

### Init

- ~100 particles randomly positioned across the canvas
- Each particle assigned a random color from the existing palette: indigo (`#6366f1`), purple (`#8b5cf6`), pink (`#ec4899`)
- Particle count scales with canvas width: `Math.floor(width / 15)`, clamped to 60–120

### Movement

- Constant gentle drift via random velocity vectors (~0.2–0.5 px/frame)
- Bounce off canvas edges to stay in view

### Connection Lines

- Drawn between particles within a distance threshold (~150px, scaled with canvas size)
- Line opacity fades linearly with distance (close = visible, far = transparent)
- Creates the neural network mesh appearance

## Mouse Interaction

- Track cursor position via `mousemove` on the canvas
- Particles within ~120px radius of cursor are pulled toward it (force proportional to `1 / distance`, capped so particles don't snap or overlap the cursor)
- Particles near cursor glow brighter (increased alpha and size in fragment shader)
- When cursor leaves the canvas, particles resume natural drift

## Scroll Transition

- Scroll listener on `window` tracks scroll position relative to the hero section
- Canvas opacity fades to 0 as user scrolls past the hero
- Driven by inline `style.opacity` based on scroll percentage through the hero height

## Mobile Behavior

- No mouse tracking (no cursor on touch devices)
- Particles drift with slightly more pronounced movement to keep visual interest
- Touch events are not tracked (a background effect has no meaningful touch interaction)

## Accessibility: `prefers-reduced-motion`

- Detected via `window.matchMedia('(prefers-reduced-motion: reduce)')`
- Falls back to a single static render: particles and connections drawn once, animation loop does not run, mouse interaction disabled
- Visually equivalent to a still screenshot of the mesh

## WebGL Rendering

### Shaders

**Vertex shader:**
- Accepts particle position (vec2) and size (float) as attributes
- Outputs `gl_Position` (clip space) and `gl_PointSize`
- Passes color and brightness to fragment shader as varyings

**Fragment shader:**
- Renders each particle as a soft radial gradient circle using `distance()` from `gl_PointCoord` center
- Applies particle's palette color with brightness modulated by cursor proximity
- Uses `discard` for pixels outside the circle radius

### Connection Lines

- Separate draw pass using `gl.LINES` with a simple shader pair (position + alpha)
- Line alpha derived from distance between the two connected particles
- Same additive blending mode (`gl.blendFunc(gl.SRC_ALPHA, gl.ONE)`) as particles for visual consistency

### Render Loop

- `requestAnimationFrame` at display refresh rate
- Each frame: update particle positions in JS, upload position buffer to GPU, draw lines then particles (two draw calls)
- Particles rendered as `gl.POINTS` on top of lines

### Blending

- Additive blending (`SRC_ALPHA, ONE`) for glow effect on both particles and lines
- Depth testing disabled (2D scene, painter's order via draw call sequence)

## Performance

- **Particle count:** 60–120, scaled by canvas width (`Math.floor(width / 15)`, clamped)
- **Connection threshold:** scales with canvas size to maintain visual density
- **`devicePixelRatio`:** capped at 2 (retina supported, 3x skipped as wasteful)
- **Resize handling:** `window.resize` listener with 250ms debounce — recreates canvas dimensions and rescales particle positions
- **Tab visibility:** animation pauses via `document.visibilitychange` when tab is hidden
- **Cleanup:** React `useEffect` return function deletes WebGL program, buffers, and cancels `requestAnimationFrame` — no leaks on unmount or HMR

## Bundle Impact

- Estimated ~4–5KB (shader strings + simulation logic + WebGL boilerplate)
- Zero external dependencies
- Hydrates only this island (`client:load`) — does not affect other static components
