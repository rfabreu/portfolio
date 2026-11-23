# Fix deprecations — plan

- [x] Record the `astro check` baseline (83 hints).
- [x] Change the `z` import to `astro/zod` and `z.string().url()` to `z.url()` in `frontend/src/content.config.ts`.
- [x] Change `React.FormEvent` to `React.SubmitEvent<HTMLFormElement>` in `ChatWidget.tsx`.
- [x] Remove the unused `tag` parameter in `analytics-funnel-events.test.ts`.
- [x] Verify `astro check` (0 hints), `npm run build`, `vitest run` and E2E.
