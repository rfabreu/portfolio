// Fails CI on high or critical npm advisories, except allowlisted ones.
// Usage (from frontend/): npm audit --json > audit.json || true
//                         node ../.github/scripts/npm-audit-gate.mjs audit.json
// Tests: node --test ../.github/scripts/
//
// Allowlist entries must name a reason and an expiry date (YYYY-MM-DD). An
// expired or malformed entry fails the gate so the exception is reviewed
// rather than forgotten.
import { readFileSync } from 'node:fs';
import { fileURLToPath } from 'node:url';

export const ALLOWLIST = {
  'GHSA-fc5q-39xc-vjpf': {
    reason:
      'picomatch ReDoS; no patched release (<= 4.0.7 affected, 4.0.7 is latest). ' +
      'Assessed, not tested: in-repo glob patterns are literals in content.config.ts ' +
      'and the deployed site is static output.',
    expires: '2027-01-31',
  },
};

const BLOCKING = new Set(['high', 'critical']);
const ISO_DATE = /^\d{4}-\d{2}-\d{2}$/;

/** Returns { errors, blocking, allowed, notices }; the gate passes only when errors and blocking are empty. */
export function evaluate(report, allowlist, today) {
  const errors = [];
  const blocking = [];
  const allowed = [];
  const notices = [];

  for (const [id, entry] of Object.entries(allowlist)) {
    if (typeof entry?.reason !== 'string' || !entry.reason.trim()) {
      errors.push(`allowlist entry ${id} has no reason`);
    }
    if (!ISO_DATE.test(entry?.expires ?? '') || Number.isNaN(Date.parse(entry.expires))) {
      errors.push(`allowlist entry ${id} has no valid expires date (YYYY-MM-DD)`);
    }
  }

  if (!report || typeof report !== 'object') return { errors: [...errors, 'audit report is not an object'], blocking, allowed, notices };
  if (report.error) return { errors: [...errors, `npm audit error: ${JSON.stringify(report.error)}`], blocking, allowed, notices };
  if (!report.vulnerabilities || typeof report.vulnerabilities !== 'object') {
    return { errors: [...errors, 'audit report has no "vulnerabilities" field'], blocking, allowed, notices };
  }

  // Only the directly affected package carries the advisory object in `via`;
  // dependents list package names (strings) and are covered by that advisory.
  const advisories = new Map();
  for (const [name, vuln] of Object.entries(report.vulnerabilities)) {
    if (!Array.isArray(vuln?.via)) {
      errors.push(`audit entry ${name} has no "via" list`);
      continue;
    }
    for (const via of vuln.via) {
      if (typeof via === 'string') {
        if (!(via in report.vulnerabilities)) errors.push(`audit entry ${name} depends on ${via}, which the report does not list`);
        continue;
      }
      if (typeof via !== 'object' || !BLOCKING.has(String(via.severity).toLowerCase())) continue;
      const id = via.url?.match(/GHSA-[\w-]+/)?.[0] ?? `npm:${via.source}`;
      advisories.set(id, via);
    }
  }

  // Fail closed if npm's own totals report high/critical findings that no advisory accounts for.
  const counts = report.metadata?.vulnerabilities ?? {};
  if ((counts.high ?? 0) + (counts.critical ?? 0) > 0 && advisories.size === 0) {
    errors.push('audit metadata reports high/critical vulnerabilities but no advisory was found');
  }

  for (const [id, via] of advisories) {
    const entry = allowlist[id];
    if (!entry) {
      blocking.push(`${id} (${via.severity}) ${via.name}: ${via.title}`);
    } else if (entry.expires < today) {
      blocking.push(`${id} allowlist entry expired on ${entry.expires}; re-review it`);
    } else {
      allowed.push(`allowed until ${entry.expires}: ${id} ${via.name}: ${entry.reason}`);
    }
  }

  for (const id of Object.keys(allowlist)) {
    if (!advisories.has(id)) notices.push(`notice: ${id} is no longer reported; remove it from the allowlist`);
  }

  return { errors, blocking, allowed, notices };
}

function main() {
  const path = process.argv[2];
  if (!path) {
    console.error('npm-audit-gate: usage: npm-audit-gate.mjs <npm-audit.json>');
    process.exit(1);
  }

  let report;
  try {
    report = JSON.parse(readFileSync(path, 'utf8'));
  } catch (err) {
    console.error(`npm-audit-gate: could not read audit report: ${err.message}`);
    process.exit(1);
  }

  const today = new Date().toISOString().slice(0, 10);
  const { errors, blocking, allowed, notices } = evaluate(report, ALLOWLIST, today);
  for (const line of [...allowed, ...notices]) console.log(line);

  const problems = [...errors, ...blocking];
  if (problems.length) {
    console.error(`npm-audit-gate: failed:\n  ${problems.join('\n  ')}`);
    process.exit(1);
  }
  console.log(`npm-audit-gate: passed (${allowed.length} high/critical advisories, all allowlisted)`);
}

if (process.argv[1] === fileURLToPath(import.meta.url)) main();
