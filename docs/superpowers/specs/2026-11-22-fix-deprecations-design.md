# Fix deprecation diagnostics

## Problem
`npx astro check` reports 0 errors but 83 hints: 74x `ts(6385)` on `z` (imported from `astro:content`), 7x on `z.string().url()`, 1x on `React.FormEvent`, and 1x `ts(6133)` for an unused `tag` parameter in `frontend/tests/analytics-funnel-events.test.ts`.

## Design
- `content.config.ts`: import `z` from `astro/zod` (the replacement named in Astro's `client.d.ts` deprecation notice; `astro/zod` re-exports `zod/v4`). `defineCollection` and `reference` stay on `astro:content`.
- Replace `z.string().url()` with `z.url()` (zod 4.6.7). Both run the same `new URL()` check, so optional-ness and accepted schemes are unchanged.
- `ChatWidget.tsx`: `React.FormEvent` becomes `React.SubmitEvent<HTMLFormElement>`, as the `@types/react` 19.3 JSDoc suggests.
- Test: drop the unused `tag` parameter from the `.map` callback.

## Non-goals
No change to `src/content/`, `package.json`, the lockfile or `.github/`. No runtime behaviour change.
