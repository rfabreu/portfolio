# P1: Interactive Particle Hero — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the static gradient orbs in the hero section with a raw WebGL particle mesh that forms an interconnected neural network pattern, reacts to cursor movement, and fades on scroll.

**Architecture:** A single React island (`ParticleHero.tsx`) renders a full-bleed `<canvas>` behind the existing hero text. All WebGL logic (custom vertex/fragment shaders, particle simulation, mouse interaction) lives in this one file. Hero.astro is modified to swap the static orb divs for the island.

**Tech Stack:** React 19, raw WebGL (no libraries), TypeScript, Astro islands (`client:load`)

---

### Task 1: Scaffold ParticleHero island with canvas and WebGL context

**Files:**
- Create: `frontend/src/islands/ParticleHero.tsx`

This task creates the React island with a canvas element and WebGL context initialization. No particles yet — just a transparent canvas that proves WebGL is working.

- [ ] **Step 1: Create the ParticleHero component with canvas ref and WebGL context**

Create `frontend/src/islands/ParticleHero.tsx`:

```tsx
import { useRef, useEffect } from 'react';

export default function ParticleHero() {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const gl = canvas.getContext('webgl', { alpha: true, premultipliedAlpha: false });
    if (!gl) return;

    const dpr = Math.min(window.devicePixelRatio, 2);

    function resize() {
      if (!canvas || !gl) return;
      const rect = canvas.getBoundingClientRect();
      canvas.width = rect.width * dpr;
      canvas.height = rect.height * dpr;
      gl.viewport(0, 0, canvas.width, canvas.height);
    }

    resize();

    let resizeTimer: number;
    function onResize() {
      clearTimeout(resizeTimer);
      resizeTimer = window.setTimeout(resize, 250);
    }
    window.addEventListener('resize', onResize);

    // Clear to transparent to confirm WebGL works
    gl.clearColor(0, 0, 0, 0);
    gl.clear(gl.COLOR_BUFFER_BIT);

    return () => {
      window.removeEventListener('resize', onResize);
      clearTimeout(resizeTimer);
    };
  }, []);

  return (
    <canvas
      ref={canvasRef}
      style={{
        position: 'absolute',
        inset: 0,
        width: '100%',
        height: '100%',
        pointerEvents: 'none',
      }}
    />
  );
}
```

- [ ] **Step 2: Verify the file has no TypeScript errors**

Run: `cd frontend && npx astro check`
Expected: No errors in `ParticleHero.tsx`

- [ ] **Step 3: Commit**

```bash
git add frontend/src/islands/ParticleHero.tsx
git commit -m "feat: scaffold ParticleHero island with WebGL context"
```

---

### Task 2: Integrate ParticleHero into Hero.astro

**Files:**
- Modify: `frontend/src/components/Hero.astro`

Replace the three static orb `<div>`s with the `ParticleHero` island. The canvas renders behind the hero text.

- [ ] **Step 1: Modify Hero.astro**

In `frontend/src/components/Hero.astro`, add the import in the frontmatter:

```astro
---
import TechTag from './TechTag.astro';
import ParticleHero from '../islands/ParticleHero.tsx';

const tags: { label: string; color: 'indigo' | 'purple' | 'pink' }[] = [
  { label: 'Go', color: 'indigo' },
  { label: 'Python', color: 'purple' },
  { label: 'React', color: 'pink' },
  { label: 'Kubernetes', color: 'indigo' },
  { label: 'AI/ML', color: 'purple' },
];
---
```

Then replace the three orb divs with the island. Change the `<section>` body from:

```html
<section id="hero" class="relative min-h-[85vh] flex items-center overflow-hidden">
  <div class="absolute -top-16 right-20 w-72 h-72 orb-indigo rounded-full pointer-events-none"></div>
  <div class="absolute -bottom-10 left-48 w-52 h-52 orb-purple rounded-full pointer-events-none"></div>
  <div class="absolute top-10 -left-10 w-44 h-44 orb-pink rounded-full pointer-events-none"></div>

  <div class="container mx-auto max-w-6xl px-6 relative z-10">
```

To:

```html
<section id="hero" class="relative min-h-[85vh] flex items-center overflow-hidden">
  <ParticleHero client:load />

  <div class="container mx-auto max-w-6xl px-6 relative z-10">
```

- [ ] **Step 2: Build and verify**

Run: `cd frontend && npm run build`
Expected: Build succeeds, 5 pages built, no errors.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/Hero.astro
git commit -m "feat: integrate ParticleHero island into Hero section"
```

---

### Task 3: Add shader compilation helpers and particle shaders

**Files:**
- Modify: `frontend/src/islands/ParticleHero.tsx`

Add WebGL shader compilation utility functions and the particle vertex/fragment shader source strings. No rendering yet — just the shader programs compiled and ready.

- [ ] **Step 1: Add shader helpers and particle shader sources**

At the top of `ParticleHero.tsx`, above the component, add:

```tsx
function createShader(gl: WebGLRenderingContext, type: number, source: string): WebGLShader | null {
  const shader = gl.createShader(type);
  if (!shader) return null;
  gl.shaderSource(shader, source);
  gl.compileShader(shader);
  if (!gl.getShaderParameter(shader, gl.COMPILE_STATUS)) {
    console.error('Shader compile error:', gl.getShaderInfoLog(shader));
    gl.deleteShader(shader);
    return null;
  }
  return shader;
}

function createProgram(
  gl: WebGLRenderingContext,
  vertSrc: string,
  fragSrc: string,
): WebGLProgram | null {
  const vert = createShader(gl, gl.VERTEX_SHADER, vertSrc);
  const frag = createShader(gl, gl.FRAGMENT_SHADER, fragSrc);
  if (!vert || !frag) return null;
  const program = gl.createProgram();
  if (!program) return null;
  gl.attachShader(program, vert);
  gl.attachShader(program, frag);
  gl.linkProgram(program);
  if (!gl.getProgramParameter(program, gl.LINK_STATUS)) {
    console.error('Program link error:', gl.getProgramInfoLog(program));
    gl.deleteProgram(program);
    return null;
  }
  // Shaders are linked — detach and delete the individual shader objects
  gl.detachShader(program, vert);
  gl.deleteShader(vert);
  gl.detachShader(program, frag);
  gl.deleteShader(frag);
  return program;
}

const PARTICLE_VERT = `
  attribute vec2 a_position;
  attribute float a_size;
  attribute vec3 a_color;
  attribute float a_brightness;
  uniform vec2 u_resolution;
  varying vec3 v_color;
  varying float v_brightness;
  void main() {
    vec2 clip = (a_position / u_resolution) * 2.0 - 1.0;
    clip.y *= -1.0;
    gl_Position = vec4(clip, 0.0, 1.0);
    gl_PointSize = a_size;
    v_color = a_color;
    v_brightness = a_brightness;
  }
`;

const PARTICLE_FRAG = `
  precision mediump float;
  varying vec3 v_color;
  varying float v_brightness;
  void main() {
    float dist = distance(gl_PointCoord, vec2(0.5));
    if (dist > 0.5) discard;
    float alpha = smoothstep(0.5, 0.1, dist) * v_brightness;
    gl_FragColor = vec4(v_color * v_brightness, alpha);
  }
`;
```

- [ ] **Step 2: Verify no TypeScript errors**

Run: `cd frontend && npx astro check`
Expected: No errors.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/islands/ParticleHero.tsx
git commit -m "feat: add WebGL shader helpers and particle shader sources"
```

---

### Task 4: Add line shaders

**Files:**
- Modify: `frontend/src/islands/ParticleHero.tsx`

Add the vertex/fragment shader sources for connection lines. These are simpler — just position and alpha.

- [ ] **Step 1: Add line shader sources after the particle shaders**

Add these constants after `PARTICLE_FRAG`:

```tsx
const LINE_VERT = `
  attribute vec2 a_position;
  attribute float a_alpha;
  uniform vec2 u_resolution;
  varying float v_alpha;
  void main() {
    vec2 clip = (a_position / u_resolution) * 2.0 - 1.0;
    clip.y *= -1.0;
    gl_Position = vec4(clip, 0.0, 1.0);
    v_alpha = a_alpha;
  }
`;

const LINE_FRAG = `
  precision mediump float;
  uniform vec3 u_color;
  varying float v_alpha;
  void main() {
    gl_FragColor = vec4(u_color, v_alpha);
  }
`;
```

- [ ] **Step 2: Verify no TypeScript errors**

Run: `cd frontend && npx astro check`
Expected: No errors.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/islands/ParticleHero.tsx
git commit -m "feat: add connection line shader sources"
```

---

### Task 5: Implement particle data structures and initialization

**Files:**
- Modify: `frontend/src/islands/ParticleHero.tsx`

Add the Particle interface, palette colors, and initialization function. This creates the particle array with random positions, velocities, and colors.

- [ ] **Step 1: Add particle types and init function**

Add these after the shader constants, above the component:

```tsx
interface Particle {
  x: number;
  y: number;
  vx: number;
  vy: number;
  r: number;
  g: number;
  b: number;
  baseSize: number;
}

const PALETTE = [
  [0.389, 0.4, 0.945],   // #6366f1 indigo
  [0.545, 0.361, 0.965],  // #8b5cf6 purple
  [0.925, 0.286, 0.600],  // #ec4899 pink
] as const;

function initParticles(width: number, height: number, isMobile: boolean): Particle[] {
  const count = Math.max(60, Math.min(120, Math.floor(width / 15)));
  const speedMultiplier = isMobile ? 1.5 : 1.0;
  const particles: Particle[] = [];
  for (let i = 0; i < count; i++) {
    const color = PALETTE[Math.floor(Math.random() * PALETTE.length)];
    particles.push({
      x: Math.random() * width,
      y: Math.random() * height,
      vx: (Math.random() - 0.5) * 0.5 * speedMultiplier,
      vy: (Math.random() - 0.5) * 0.5 * speedMultiplier,
      r: color[0],
      g: color[1],
      b: color[2],
      baseSize: 2 + Math.random() * 2,
    });
  }
  return particles;
}
```

- [ ] **Step 2: Verify no TypeScript errors**

Run: `cd frontend && npx astro check`
Expected: No errors.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/islands/ParticleHero.tsx
git commit -m "feat: add particle data structures and initialization"
```

---

### Task 6: Implement the WebGL render loop with particles

**Files:**
- Modify: `frontend/src/islands/ParticleHero.tsx`

Replace the placeholder `useEffect` body with the full WebGL setup: compile shaders, create buffers, init particles, and run the animation loop rendering particles as `gl.POINTS`.

- [ ] **Step 1: Rewrite the useEffect to set up WebGL rendering**

Replace the entire `useEffect` in the component with:

```tsx
  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const gl = canvas.getContext('webgl', { alpha: true, premultipliedAlpha: false });
    if (!gl) return;

    const dpr = Math.min(window.devicePixelRatio, 2);
    const isMobile = 'ontouchstart' in window;
    const reducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches;

    // Compile shader programs
    const particleProg = createProgram(gl, PARTICLE_VERT, PARTICLE_FRAG);
    const lineProg = createProgram(gl, LINE_VERT, LINE_FRAG);
    if (!particleProg || !lineProg) return;

    // Particle program locations
    const pLoc = {
      position: gl.getAttribLocation(particleProg, 'a_position'),
      size: gl.getAttribLocation(particleProg, 'a_size'),
      color: gl.getAttribLocation(particleProg, 'a_color'),
      brightness: gl.getAttribLocation(particleProg, 'a_brightness'),
      resolution: gl.getUniformLocation(particleProg, 'u_resolution'),
    };

    // Line program locations
    const lLoc = {
      position: gl.getAttribLocation(lineProg, 'a_position'),
      alpha: gl.getAttribLocation(lineProg, 'a_alpha'),
      resolution: gl.getUniformLocation(lineProg, 'u_resolution'),
      color: gl.getUniformLocation(lineProg, 'u_color'),
    };

    // Buffers
    const particlePosBuf = gl.createBuffer();
    const particleSizeBuf = gl.createBuffer();
    const particleColorBuf = gl.createBuffer();
    const particleBrightBuf = gl.createBuffer();
    const linePosBuf = gl.createBuffer();
    const lineAlphaBuf = gl.createBuffer();

    // State
    let w = 0;
    let h = 0;
    let particles: Particle[] = [];
    let mouseX = -1000;
    let mouseY = -1000;
    let animId = 0;
    let paused = false;

    function resize() {
      if (!canvas || !gl) return;
      const rect = canvas.getBoundingClientRect();
      w = rect.width * dpr;
      h = rect.height * dpr;
      canvas.width = w;
      canvas.height = h;
      gl.viewport(0, 0, w, h);
      particles = initParticles(w, h, isMobile);
    }

    resize();

    let resizeTimer: number;
    function onResize() {
      clearTimeout(resizeTimer);
      resizeTimer = window.setTimeout(resize, 250);
    }

    // Mouse tracking (desktop only)
    function onMouseMove(e: MouseEvent) {
      if (isMobile) return;
      const rect = canvas!.getBoundingClientRect();
      mouseX = (e.clientX - rect.left) * dpr;
      mouseY = (e.clientY - rect.top) * dpr;
    }
    function onMouseLeave() {
      mouseX = -1000;
      mouseY = -1000;
    }

    // Scroll fade
    function onScroll() {
      if (!canvas) return;
      const rect = canvas.getBoundingClientRect();
      const progress = Math.max(0, Math.min(1, -rect.top / rect.height));
      canvas.style.opacity = String(1 - progress);
    }

    // Visibility
    function onVisibility() {
      paused = document.hidden;
      if (!paused && !reducedMotion) {
        animId = requestAnimationFrame(frame);
      }
    }

    const MOUSE_RADIUS = 120 * dpr;
    const CONNECTION_DIST = Math.min(150 * dpr, w * 0.12);

    function frame() {
      if (paused || !gl) return;

      const count = particles.length;

      // Update positions
      for (let i = 0; i < count; i++) {
        const p = particles[i];

        // Mouse attraction (desktop)
        if (mouseX > 0 && mouseY > 0) {
          const dx = mouseX - p.x;
          const dy = mouseY - p.y;
          const dist = Math.sqrt(dx * dx + dy * dy);
          if (dist < MOUSE_RADIUS && dist > 1) {
            const force = Math.min(0.8, 1 / dist) * 0.5;
            p.vx += dx * force * 0.01;
            p.vy += dy * force * 0.01;
          }
        }

        p.x += p.vx;
        p.y += p.vy;

        // Dampen velocity back toward base speed
        p.vx *= 0.99;
        p.vy *= 0.99;

        // Ensure minimum drift
        const speed = Math.sqrt(p.vx * p.vx + p.vy * p.vy);
        const minSpeed = isMobile ? 0.3 : 0.2;
        if (speed < minSpeed) {
          p.vx += (Math.random() - 0.5) * 0.1;
          p.vy += (Math.random() - 0.5) * 0.1;
        }

        // Bounce off edges
        if (p.x < 0) { p.x = 0; p.vx = Math.abs(p.vx); }
        if (p.x > w) { p.x = w; p.vx = -Math.abs(p.vx); }
        if (p.y < 0) { p.y = 0; p.vy = Math.abs(p.vy); }
        if (p.y > h) { p.y = h; p.vy = -Math.abs(p.vy); }
      }

      // Build particle buffers
      const positions = new Float32Array(count * 2);
      const sizes = new Float32Array(count);
      const colors = new Float32Array(count * 3);
      const brightnesses = new Float32Array(count);

      for (let i = 0; i < count; i++) {
        const p = particles[i];
        positions[i * 2] = p.x;
        positions[i * 2 + 1] = p.y;
        colors[i * 3] = p.r;
        colors[i * 3 + 1] = p.g;
        colors[i * 3 + 2] = p.b;

        let brightness = 0.5;
        let size = p.baseSize * dpr;
        if (mouseX > 0 && mouseY > 0) {
          const dx = mouseX - p.x;
          const dy = mouseY - p.y;
          const dist = Math.sqrt(dx * dx + dy * dy);
          if (dist < MOUSE_RADIUS) {
            const t = 1 - dist / MOUSE_RADIUS;
            brightness = 0.5 + t * 0.5;
            size += t * 3 * dpr;
          }
        }
        sizes[i] = size;
        brightnesses[i] = brightness;
      }

      // Build line buffers
      const linePositions: number[] = [];
      const lineAlphas: number[] = [];
      for (let i = 0; i < count; i++) {
        for (let j = i + 1; j < count; j++) {
          const dx = particles[i].x - particles[j].x;
          const dy = particles[i].y - particles[j].y;
          const dist = Math.sqrt(dx * dx + dy * dy);
          if (dist < CONNECTION_DIST) {
            const alpha = (1 - dist / CONNECTION_DIST) * 0.15;
            linePositions.push(
              particles[i].x, particles[i].y,
              particles[j].x, particles[j].y,
            );
            lineAlphas.push(alpha, alpha);
          }
        }
      }

      // Render
      gl.clearColor(0, 0, 0, 0);
      gl.clear(gl.COLOR_BUFFER_BIT);
      gl.enable(gl.BLEND);
      gl.blendFunc(gl.SRC_ALPHA, gl.ONE);
      gl.disable(gl.DEPTH_TEST);

      // Draw lines
      if (linePositions.length > 0) {
        gl.useProgram(lineProg);
        gl.uniform2f(lLoc.resolution, w, h);
        gl.uniform3f(lLoc.color, 0.389, 0.4, 0.945); // indigo for lines

        gl.bindBuffer(gl.ARRAY_BUFFER, linePosBuf);
        gl.bufferData(gl.ARRAY_BUFFER, new Float32Array(linePositions), gl.DYNAMIC_DRAW);
        gl.enableVertexAttribArray(lLoc.position);
        gl.vertexAttribPointer(lLoc.position, 2, gl.FLOAT, false, 0, 0);

        gl.bindBuffer(gl.ARRAY_BUFFER, lineAlphaBuf);
        gl.bufferData(gl.ARRAY_BUFFER, new Float32Array(lineAlphas), gl.DYNAMIC_DRAW);
        gl.enableVertexAttribArray(lLoc.alpha);
        gl.vertexAttribPointer(lLoc.alpha, 1, gl.FLOAT, false, 0, 0);

        gl.drawArrays(gl.LINES, 0, linePositions.length / 2);

        gl.disableVertexAttribArray(lLoc.position);
        gl.disableVertexAttribArray(lLoc.alpha);
      }

      // Draw particles
      gl.useProgram(particleProg);
      gl.uniform2f(pLoc.resolution, w, h);

      gl.bindBuffer(gl.ARRAY_BUFFER, particlePosBuf);
      gl.bufferData(gl.ARRAY_BUFFER, positions, gl.DYNAMIC_DRAW);
      gl.enableVertexAttribArray(pLoc.position);
      gl.vertexAttribPointer(pLoc.position, 2, gl.FLOAT, false, 0, 0);

      gl.bindBuffer(gl.ARRAY_BUFFER, particleSizeBuf);
      gl.bufferData(gl.ARRAY_BUFFER, sizes, gl.DYNAMIC_DRAW);
      gl.enableVertexAttribArray(pLoc.size);
      gl.vertexAttribPointer(pLoc.size, 1, gl.FLOAT, false, 0, 0);

      gl.bindBuffer(gl.ARRAY_BUFFER, particleColorBuf);
      gl.bufferData(gl.ARRAY_BUFFER, colors, gl.DYNAMIC_DRAW);
      gl.enableVertexAttribArray(pLoc.color);
      gl.vertexAttribPointer(pLoc.color, 3, gl.FLOAT, false, 0, 0);

      gl.bindBuffer(gl.ARRAY_BUFFER, particleBrightBuf);
      gl.bufferData(gl.ARRAY_BUFFER, brightnesses, gl.DYNAMIC_DRAW);
      gl.enableVertexAttribArray(pLoc.brightness);
      gl.vertexAttribPointer(pLoc.brightness, 1, gl.FLOAT, false, 0, 0);

      gl.drawArrays(gl.POINTS, 0, count);

      gl.disableVertexAttribArray(pLoc.position);
      gl.disableVertexAttribArray(pLoc.size);
      gl.disableVertexAttribArray(pLoc.color);
      gl.disableVertexAttribArray(pLoc.brightness);

      if (!reducedMotion) {
        animId = requestAnimationFrame(frame);
      }
    }

    // Event listeners
    window.addEventListener('resize', onResize);
    canvas.addEventListener('mousemove', onMouseMove);
    canvas.addEventListener('mouseleave', onMouseLeave);
    window.addEventListener('scroll', onScroll, { passive: true });
    document.addEventListener('visibilitychange', onVisibility);

    // Enable pointer events on canvas for mouse tracking
    canvas.style.pointerEvents = 'auto';

    // Start
    if (reducedMotion) {
      frame(); // Single static render
    } else {
      animId = requestAnimationFrame(frame);
    }

    return () => {
      cancelAnimationFrame(animId);
      clearTimeout(resizeTimer);
      window.removeEventListener('resize', onResize);
      canvas.removeEventListener('mousemove', onMouseMove);
      canvas.removeEventListener('mouseleave', onMouseLeave);
      window.removeEventListener('scroll', onScroll);
      document.removeEventListener('visibilitychange', onVisibility);
      // Clean up WebGL resources
      gl.deleteBuffer(particlePosBuf);
      gl.deleteBuffer(particleSizeBuf);
      gl.deleteBuffer(particleColorBuf);
      gl.deleteBuffer(particleBrightBuf);
      gl.deleteBuffer(linePosBuf);
      gl.deleteBuffer(lineAlphaBuf);
      gl.deleteProgram(particleProg);
      gl.deleteProgram(lineProg);
    };
  }, []);
```

- [ ] **Step 2: Verify no TypeScript errors**

Run: `cd frontend && npx astro check`
Expected: No errors.

- [ ] **Step 3: Build the project**

Run: `cd frontend && npm run build`
Expected: Build succeeds with 5 pages.

- [ ] **Step 4: Commit**

```bash
git add frontend/src/islands/ParticleHero.tsx
git commit -m "feat: implement WebGL particle render loop with mouse interaction"
```

---

### Task 7: Visual verification and tuning

**Files:**
- Possibly modify: `frontend/src/islands/ParticleHero.tsx`

Run the dev server and verify the particle hero looks and behaves correctly across all states.

- [ ] **Step 1: Start the dev server**

Run: `cd frontend && npm run dev`

Open `http://localhost:4321` in a browser.

- [ ] **Step 2: Verify desktop behavior**

Check all of:
- Particles are visible behind the hero text (indigo, purple, pink colors)
- Connection lines appear between nearby particles
- Moving the mouse causes nearby particles to pull toward cursor
- Particles near cursor glow brighter and grow slightly
- Scrolling past the hero fades the canvas to transparent
- Moving the mouse away from canvas causes particles to resume drifting

- [ ] **Step 3: Verify mobile behavior**

Open browser DevTools, toggle device toolbar (mobile viewport):
- Particles drift without mouse interaction
- Drift is slightly more pronounced than desktop
- No console errors related to touch events

- [ ] **Step 4: Verify reduced motion**

In browser DevTools, enable "Prefers reduced motion" in Rendering tab:
- Particles and connections render once as a static image
- No animation loop running (check via Performance tab — no recurring frames)

- [ ] **Step 5: Tune values if needed**

If the visual result needs adjustment, tune these constants in `ParticleHero.tsx`:
- `MOUSE_RADIUS` — how far the cursor influence reaches
- `CONNECTION_DIST` — how far particles can be and still connect
- `0.15` in line alpha calculation — line opacity intensity
- `0.5` base brightness — default particle brightness
- `p.baseSize` range (`2 + Math.random() * 2`) — particle size range

- [ ] **Step 6: Commit any tuning changes**

```bash
git add frontend/src/islands/ParticleHero.tsx
git commit -m "fix: tune particle hero visual parameters"
```

(Skip this commit if no changes were needed.)

---

### Task 8: Final build verification and commit

**Files:**
- No new changes expected

Final validation that everything builds clean and the branch is ready for PR.

- [ ] **Step 1: Run astro check**

Run: `cd frontend && npx astro check`
Expected: No errors.

- [ ] **Step 2: Run production build**

Run: `cd frontend && npm run build`
Expected: Build succeeds, 5 pages built.

- [ ] **Step 3: Run Go API checks (ensure no regressions)**

Run: `cd api && go vet ./... && go test ./... -v`
Expected: All pass (API is untouched, but confirm CI will pass).

- [ ] **Step 4: Verify git status is clean**

Run: `git status`
Expected: Nothing unexpected in the working tree. Only committed changes on `feature/particle-hero`.

- [ ] **Step 5: Push and create PR**

```bash
git push -u origin feature/particle-hero
gh pr create --base main --head feature/particle-hero \
  --title "feat: interactive WebGL particle hero background" \
  --body "## Summary
- Replace static gradient orbs with raw WebGL particle mesh
- Neural network pattern with ~100 interconnected particles
- Mouse interaction: particles attracted to cursor with glow effect
- Scroll fade: canvas opacity transitions to 0 as user scrolls past
- Mobile: gentle drift animation, no mouse tracking
- prefers-reduced-motion: single static render, no animation loop
- Zero dependencies, ~4-5KB bundle addition

## Test plan
- [ ] Desktop: particles visible, mouse interaction works, scroll fade works
- [ ] Mobile viewport: particles drift, no console errors
- [ ] Reduced motion: static render, no animation loop
- [ ] Build passes (astro check + npm run build)
- [ ] No API regressions (go vet + go test)"
```
