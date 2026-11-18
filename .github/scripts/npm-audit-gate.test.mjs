// Run with: node --test .github/scripts/
import { test } from 'node:test';
import assert from 'node:assert/strict';
import { evaluate, ALLOWLIST } from './npm-audit-gate.mjs';

const TODAY = '2026-11-17';
const ALLOW = { 'GHSA-aaaa-aaaa-aaaa': { reason: 'test exception', expires: '2027-01-31' } };

function advisory(ghsa, severity, name = 'pkg') {
  return { source: 1, name, title: 't', url: `https://github.com/advisories/${ghsa}`, severity };
}

function report(...vias) {
  const vulnerabilities = {};
  vias.forEach((via, i) => {
    vulnerabilities[`pkg${i}`] = { via: [via] };
  });
  return { vulnerabilities, metadata: { vulnerabilities: { high: vias.length } } };
}

function passes(result) {
  return result.errors.length === 0 && result.blocking.length === 0;
}

test('allowlisted advisory passes', () => {
  assert.ok(passes(evaluate(report(advisory('GHSA-aaaa-aaaa-aaaa', 'high')), ALLOW, TODAY)));
});

test('unlisted high or critical advisory fails, whatever the casing', () => {
  for (const severity of ['high', 'critical', 'HIGH', 'Critical']) {
    const r = report(advisory('GHSA-aaaa-aaaa-aaaa', 'high'), advisory('GHSA-bbbb-bbbb-bbbb', severity));
    assert.equal(evaluate(r, ALLOW, TODAY).blocking.length, 1, severity);
  }
});

test('advisory without a GHSA url still blocks', () => {
  const r = report({ source: 99, name: 'x', title: 't', severity: 'critical' });
  assert.deepEqual(evaluate(r, {}, TODAY).blocking.length, 1);
});

test('moderate only passes and flags a stale allowlist entry', () => {
  const r = { vulnerabilities: { x: { via: [advisory('GHSA-cccc-cccc-cccc', 'moderate')] } }, metadata: { vulnerabilities: { moderate: 1 } } };
  const result = evaluate(r, ALLOW, TODAY);
  assert.ok(passes(result));
  assert.equal(result.notices.length, 1);
});

test('expired entry fails', () => {
  const expired = { 'GHSA-aaaa-aaaa-aaaa': { reason: 'r', expires: '2026-01-01' } };
  assert.equal(evaluate(report(advisory('GHSA-aaaa-aaaa-aaaa', 'high')), expired, TODAY).blocking.length, 1);
});

test('entry without a valid expires date or reason fails', () => {
  const bad = [
    { reason: 'r' },
    { reason: 'r', expires: '2027-1-31' },
    { reason: 'r', expires: '2027-13-45' },
    { expires: '2027-01-31' },
    { reason: '  ', expires: '2027-01-31' },
  ];
  for (const entry of bad) {
    const result = evaluate(report(advisory('GHSA-aaaa-aaaa-aaaa', 'high')), { 'GHSA-aaaa-aaaa-aaaa': entry }, TODAY);
    assert.ok(result.errors.length > 0, JSON.stringify(entry));
  }
});

test('npm error payload and malformed reports fail', () => {
  for (const r of [null, 'text', {}, { error: { code: 'ENOAUDIT' } }, { vulnerabilities: { x: {} } }]) {
    assert.ok(evaluate(r, {}, TODAY).errors.length > 0, JSON.stringify(r));
  }
});

test('dependency on a package missing from the report fails closed', () => {
  const r = report(advisory('GHSA-aaaa-aaaa-aaaa', 'high'));
  r.vulnerabilities.foo = { severity: 'critical', via: ['ghostpkg'] };
  assert.ok(evaluate(r, ALLOW, TODAY).errors.length > 0);
});

test('metadata totals with no advisory fail closed', () => {
  const r = { vulnerabilities: {}, metadata: { vulnerabilities: { high: 3 } } };
  assert.ok(evaluate(r, {}, TODAY).errors.length > 0);
});

test('shipped allowlist entries are well formed', () => {
  assert.deepEqual(evaluate({ vulnerabilities: {} }, ALLOWLIST, TODAY).errors, []);
});
