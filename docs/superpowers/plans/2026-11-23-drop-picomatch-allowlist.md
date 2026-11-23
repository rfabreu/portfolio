# Drop the picomatch Allowlist Entry — Plan

Spec: `docs/superpowers/specs/2026-11-23-drop-picomatch-allowlist-design.md`

1. Branch `chore/drop-picomatch-allowlist` off `origin/main`.
2. Empty `ALLOWLIST`; reword the pass message.
3. Update CLAUDE.md (CI bullet, Known Follow-ups).
4. Run the gate unit tests and the gate as `ci.yml` does (`npm ci`, `npm audit`, gate).
5. PR to `main`; no merge without owner approval.
