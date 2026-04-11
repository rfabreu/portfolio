# P5: Privacy-Respecting Analytics — Design Spec

## Overview

Add Umami Cloud analytics to the portfolio to track visitor behavior without cookies or personal data. Zero cost (free tier), zero maintenance (hosted), GDPR compliant.

## Script Tag

```html
<script is:inline defer src="https://cloud.umami.is/script.js" data-website-id="f17cb708-288b-471b-989f-d7115bfcbec4"></script>
```

Added to both `BaseLayout.astro` and `ProjectLayout.astro` `<head>` sections. Uses `is:inline` to prevent Astro bundling. Wrapped in `import.meta.env.PROD` check to skip in dev.

## Custom Event Tracking

Umami supports custom events via `umami.track('event-name')` calls. Track these interactions:

- **chat-opened** — fired when user opens the ChatWidget (clicks the FAB button to open, not close)
- **game-started** — fired on the first move in a GameBoard session (not on subsequent moves or replays)

Events use the global `umami.track()` function which is available after the script loads. Guard calls with `typeof umami !== 'undefined'` to prevent errors in dev or if the script fails to load.

## Environment Awareness

Only load the script in production builds. In Astro, use:

```astro
{import.meta.env.PROD && (
  <script is:inline defer src="..." data-website-id="..."></script>
)}
```

This prevents localhost traffic from polluting analytics data.

## Files Changed

| File | Change |
|------|--------|
| `src/layouts/BaseLayout.astro` | Add Umami script tag (prod only) |
| `src/layouts/ProjectLayout.astro` | Add Umami script tag (prod only) |
| `src/islands/ChatWidget.tsx` | Add `umami.track('chat-opened')` on open |
| `src/islands/GameBoard.tsx` | Add `umami.track('game-started')` on first move |

## Out of Scope

- Self-hosted Umami
- Custom dashboard embedding in the portfolio
- Server-side analytics
- Tracking project card clicks (Umami auto-tracks page navigations)
