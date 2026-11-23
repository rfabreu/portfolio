# Drop the picomatch Allowlist Entry — Design Spec

> Removes the only `ALLOWLIST` entry from `.github/scripts/npm-audit-gate.mjs`. No frontend or API changes, so merging deploys nothing.

## Why
GitHub withdrew GHSA-fc5q-39xc-vjpf (picomatch ReDoS) on 2026-11-18T15:38:56Z (`gh api /advisories/GHSA-fc5q-39xc-vjpf --jq .withdrawn_at`). `npm audit --json` in `frontend/` reports 0 vulnerabilities and the gate prints a "no longer reported" notice.

## Changes
- `ALLOWLIST` becomes `{}`; reason/expiry validation and notices stay for future use.
- Pass message reads "passed (N allowlisted high/critical advisories)" so 0 reads sensibly.
- Existing unit tests use local fixtures; the "shipped allowlist entries are well formed" test already covers the real `ALLOWLIST`.
- CLAUDE.md: CI bullet states the allowlist is empty; the Known Follow-ups bullet is removed ("None.").
- Older specs/plans mentioning the entry are history and stay as written.
