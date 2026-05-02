# Merge Example: Profile API Tests

## Source (staging): `tests/e2e/features/user-profile/api.spec.ts`

```typescript
import { test, expect } from '@playwright/test';
import { curl, apiBaseUrl } from '../../helpers.js';

test.describe('Profile API', () => {
  // Traceability: TC-010 → Story 5 / AC-2
  test('TC-010: GET /api/profile returns current user profile', async () => {
    const res = await curl('GET', `${apiBaseUrl()}/api/profile`);
    expect(res.status).toBe(200);
    const data = JSON.parse(res.body);
    expect(data.username).toBeTruthy();
  });

  // Traceability: TC-011 → Story 5 / AC-3
  test('TC-011: PUT /api/profile updates display name', async () => {
    const res = await curl('PUT', `${apiBaseUrl()}/api/profile`, {
      body: JSON.stringify({ displayName: 'New Name' }),
    });
    expect(res.status).toBe(200);
  });
});
```

## Target (regression, already exists): `tests/e2e/profile/api.spec.ts`

```typescript
import { test, expect } from '@playwright/test';
import { curl, apiBaseUrl } from '../helpers.js';

test.describe('Profile API', () => {
  // Traceability: TC-001 → Story 3 / AC-1
  test('TC-001: GET /api/profile/:id returns user by ID', async () => {
    const res = await curl('GET', `${apiBaseUrl()}/api/profile/42`);
    expect(res.status).toBe(200);
    const data = JSON.parse(res.body);
    expect(data.id).toBe(42);
  });

  // Traceability: TC-002 → Story 3 / AC-2
  test('TC-002: DELETE /api/profile/:id removes user', async () => {
    const res = await curl('DELETE', `${apiBaseUrl()}/api/profile/42`);
    expect(res.status).toBe(204);
  });
});
```

## Merged result: `tests/e2e/profile/api.spec.ts`

```typescript
import { test, expect } from '@playwright/test';
import { curl, apiBaseUrl } from '../helpers.js';
//                          ^^^ import rewritten from ../../helpers.js to ../helpers.js

test.describe('Profile API', () => {
  // ── Existing tests (from prior graduation) ──

  // Traceability: TC-001 → Story 3 / AC-1
  test('TC-001: GET /api/profile/:id returns user by ID', async () => {
    const res = await curl('GET', `${apiBaseUrl()}/api/profile/42`);
    expect(res.status).toBe(200);
    const data = JSON.parse(res.body);
    expect(data.id).toBe(42);
  });

  // Traceability: TC-002 → Story 3 / AC-2
  test('TC-002: DELETE /api/profile/:id removes user', async () => {
    const res = await curl('DELETE', `${apiBaseUrl()}/api/profile/42`);
    expect(res.status).toBe(204);
  });

  // ── New tests (from user-profile graduation) ──

  // Traceability: TC-010 → Story 5 / AC-2
  test('TC-010: GET /api/profile returns current user profile', async () => {
    const res = await curl('GET', `${apiBaseUrl()}/api/profile`);
    expect(res.status).toBe(200);
    const data = JSON.parse(res.body);
    expect(data.username).toBeTruthy();
  });

  // Traceability: TC-011 → Story 5 / AC-3
  test('TC-011: PUT /api/profile updates display name', async () => {
    const res = await curl('PUT', `${apiBaseUrl()}/api/profile`, {
      body: JSON.stringify({ displayName: 'New Name' }),
    });
    expect(res.status).toBe(200);
  });
});
```

## What happened in this merge

| Step | Action |
|------|--------|
| Import dedup | Both files import `test, expect, curl, apiBaseUrl` → single combined import block |
| Describe match | Both have `test.describe('Profile API')` → merged into one block |
| Test dedup | No duplicate titles → all 4 tests kept (TC-001, TC-002, TC-010, TC-011) |
| Import path | `'../../helpers.js'` → `'../helpers.js'` (staging → regression) |
| Nesting | Both are flat → no nesting preservation needed |

## Manifest entry

```json
{
  "entries": [
    {
      "targetPath": "tests/e2e/profile/api.spec.ts",
      "wasExistingBeforeMerge": true,
      "status": "done"
    }
  ]
}
```

## Backup

```
tests/e2e/.graduated/.backup/user-profile/profile__api.spec.ts
```
Contains the pre-merge target file for rollback.
