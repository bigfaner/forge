import { describe, test } from 'node:test';
import assert from 'node:assert/strict';
import { curl, baseUrl } from './helpers.js';

const apiUrl = process.env.E2E_API_URL ?? 'http://localhost:8080';

describe('API E2E Tests', () => {
  // Traceability: TC-010 → Spec Section 5.2
  test('TC-010: GET /api/health returns 200', async () => {
    const res = await curl('GET', `${apiUrl}/api/health`);
    assert.equal(res.status, 200);
  });
});
