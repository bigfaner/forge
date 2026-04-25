import { describe, test, before } from 'node:test';
import assert from 'node:assert/strict';
import { curl, getApiToken, createAuthCurl, defaultCreds, apiUrl } from './helpers.js';

describe('API E2E Tests', () => {
  let authCurl: ReturnType<typeof createAuthCurl>;

  before(async () => {
    // CONDITIONAL: Keep only if auth-required-test exists; remove if only public-test/login-test
    const token = await getApiToken(apiUrl);
    authCurl = createAuthCurl(apiUrl, token);
  });

  // ── Auth Tests (no shared auth) ─────────────────────────────────
  // CONDITIONAL: Remove this block if no login-test exists
  // Traceability: TC-010 → Spec Section 5.2
  test('TC-010: POST /api/auth/login returns 200 with valid credentials', async () => {
    const res = await curl('POST', `${apiUrl}/api/auth/login`, {
      body: JSON.stringify(defaultCreds),
    });
    assert.equal(res.status, 200);
    const data = JSON.parse(res.body);
    assert.ok(data.token ?? data.access_token, 'Response contains token');
  });

  // ── Authenticated Tests (use shared auth) ───────────────────────
  // Traceability: TC-011 → Spec Section 5.3
  test('TC-011: GET /api/users returns 200', async () => {
    const res = await authCurl('GET', '/api/users');
    assert.equal(res.status, 200);
  });
});
