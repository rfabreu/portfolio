# Redirect the legacy `rafaelswe` Netlify site — Design Spec

> Two Netlify sites build this repo. `rafaelabreu` is production; `rafaelswe` is the old site. Touches the root `netlify.toml`, a Vitest guard, `CLAUDE.md` and docs. Owner approved the change. The new test is under `frontend/**`, so merging also runs the production frontend deploy (static output unchanged). Separately, the `rafaelswe` Git integration rebuilds on every push to `main`, which is what activates the rule.

## 1. Problem

Verified read-only with the Netlify CLI:

- `rafaelabreu` (https://rafaelabreu.netlify.app) is production, deployed only by `.github/workflows/deploy-frontend.yml` (`nwtgck/actions-netlify@v4`, `publish-dir: frontend/dist`).
- `rafaelswe` (https://rafaelswe.netlify.app, created July 2025) builds every push to `main` through its own Git integration using the root `netlify.toml` (`base = "frontend"`). It also builds the PR deploy previews (`deploy-preview-N--rafaelswe.netlify.app`, status `netlify/rafaelswe/deploy-preview`). `getSite` shows `custom_domain: null` and `domain_aliases: []`.
- On `rafaelswe` the chatbot is broken (no `PUBLIC_API_URL`, calls localhost) and the contact form is not registered, so messages are lost. It is also duplicate content for search engines.

## 2. Change

Add a host-scoped permanent redirect to the root `netlify.toml`:

```toml
[[redirects]]
  from = "https://rafaelswe.netlify.app/*"
  to = "https://rafaelabreu.netlify.app/:splat"
  status = 301
  force = true
```

- `from` names the exact host, so only requests to `rafaelswe.netlify.app` match. `deploy-preview-N--rafaelswe.netlify.app` and branch-deploy hosts are different hosts.
- `:splat` keeps the path. `force = true` is needed because real files exist at those paths (`/`, `/resume/`, assets), and a non-forced rule never shadows existing content.
- Query strings: Netlify passes them through on 301 (see section 3).

## 3. Evidence from Netlify docs (fetched 2026-11-22)

- `force` ("Whether to override any existing content in the path or not; `false` by default") and "By default, you can't shadow a URL that actually exists within the site" — https://docs.netlify.com/routing/redirects/ and https://docs.netlify.com/routing/redirects/rewrites-proxies/.
- Full-URL (domain-level) `from` is shown only as `http://blog.yoursite.com/* https://www.yoursite.com/blog/:splat 301!`, with: "For Netlify to redirect from a domain that is not the site's main custom domain, the domain must be assigned to the site" — https://docs.netlify.com/routing/redirects/redirect-options/. The same page: "http and https need separate rules if you don't force_ssl!".
- Query strings: "our service automatically passes on all query string parameters to destination paths for redirects with `200`, `301`, and `302` HTTP status" — redirect-options page.
- Rules are processed top to bottom, first match wins; this is the only rule in the file.

**Not documented (stated plainly):** the docs do not say whether a host-scoped `from` works for a site's default `*.netlify.app` host, nor whether it matches `deploy-preview-N--` hosts. They do not say it fails either, so the stop condition did not trigger. The `rafaelswe` host is the site's primary (`default_domain`) host, so the "assigned to the site" condition is plausibly met. Both points are therefore verified empirically, not from docs:

- Previews not matched: this PR's own deploy preview must return 200 with real HTML (checked before merge).
- Redirect works: only observable after merge, since `main` is the only `rafaelswe` production build. Post-merge checklist in the PR body. If it does not fire, the fallback is a manual setting in the Netlify UI, which this change does not touch.
- `http://rafaelswe.netlify.app` is expected to be upgraded to HTTPS by Netlify first; the post-merge checklist checks it, and a second `http://` rule is added only if it fails.

## 4. Production is unaffected

- The workflow does not set the action's `netlify-config-path` input. The action only forwards that input as `configPath` (`src/main.ts` lines 86 and 104 in `nwtgck/actions-netlify`); whether the client then reads the root `netlify.toml` was not verified.
- It does not matter: the rule matches only requests whose host is `rafaelswe.netlify.app`. Requests to `rafaelabreu.netlify.app` never match, so even if production read the file the rule is a no-op there. There is no redirect loop: the target host has no rule sending it elsewhere.

## 5. Tests

`frontend/tests/legacy-site-redirect.test.ts` reads the root `netlify.toml`. `frontend/package.json` has no TOML parser and the file is flat, so it splits the text into `[[redirects]]` blocks and asserts there is exactly one, with the exact host in `from`, `to` ending `/:splat` on the production host, `status = 301` and `force = true`. It also asserts the `from` is not a bare path (a path-only rule would redirect previews and production too).

## 6. Verification

`npm ci`, `npx astro check`, `npm run build`, `npx vitest run` in `frontend/`; PR deploy preview returns 200 HTML; post-merge checklist run by the orchestrator.
