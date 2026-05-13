# Playwright Run Strategy

Profile-specific execution and result-parsing rules for the `run-e2e-tests` skill.

## Execution

| Item | Value |
|------|-------|
| Command | `npx playwright test` |
| Invoked via | `just test-e2e --feature <slug>` |
| Setup | `just e2e-setup` (npm deps + Playwright chromium) |
| Teardown | Kill tracked servers via PID files; browser instances auto-closed by test runner |
| Reporters | JSON (structured) + list (human-readable) |

## Config Behavior

`playwright.config.ts` sets `testIgnore: /^features\//` — default test discovery skips staging specs. The justfile recipe sets `E2E_FEATURE=1` to disable `testIgnore` for feature-specific runs.

## Timeouts

| Setting | Value | Notes |
|---------|-------|-------|
| `globalTimeout` | 300000ms | 5-min hard cap for entire run |
| `timeout` | 30000ms | Per-test max |
| `expect.timeout` | 10000ms | Per-assertion max |
| `retries` | 0 | Override via `E2E_RETRIES` env var |

## Result Format

Playwright JSON reporter output:

```json
{
  "config": { "rootDir": "...", "projects": [...] },
  "suites": [{
    "specs": [{
      "title": "TC-001: Login with valid credentials",
      "ok": true,
      "tests": [{
        "status": "expected",
        "results": [{
          "status": "passed",
          "duration": 1234,
          "errors": [{ "message": "..." }],
          "attachments": [{ "name": "screenshot", "path": "results/TC-001.png", "contentType": "image/png" }]
        }]
      }]
    }]
  }]
}
```

## Result Parsing Rules

### Suite Traversal

Recursively walk `suites[].specs[]`, descending into nested `suites[].suites[]`.

### Field Mapping

| Data | Source | Notes |
|------|--------|-------|
| TC ID | `spec.title` via regex `TC-\d+` | Extract from title string |
| Quick pass | `spec.ok` (boolean) | May be absent in some Playwright versions; derive from `tests[].results[].status` |
| Test-level status | `spec.tests[].status` | `"expected"` (passed), `"unexpected"` (failed), `"skipped"`, `"flaky"` |
| Result-level status | `spec.tests[].results[].status` | **Authoritative**: `"passed"`, `"failed"`, `"skipped"`, `"timedOut"`, `"interrupted"` |
| Duration | `spec.tests[].results[].duration` | Milliseconds |
| Errors | `spec.tests[].results[].errors[].message` | Per-result |
| Attachments | `spec.tests[].results[].attachments[]` | Fields: `name`, `path`, `contentType` |

## Test Type Classification

| Type | Indicators |
|------|------------|
| UI | Uses `page` fixture (`async ({ page })`), imports `screenshot`/`loginViaUI` |
| API | Uses `curl`/`fetch`/`authCurl`, imports `curl`, `apiBaseUrl`, `getApiToken` |
| CLI | Uses `runCli` or `child_process` |

## Screenshot Discovery

Two sources, both captured via single glob:

```
glob tests/e2e/results/**/*.png
```

| Source | Path pattern | Origin |
|--------|-------------|--------|
| Explicit screenshots | `tests/e2e/results/screenshots/TC-NNN.png` | `screenshot(page, tcId)` calls |
| Auto-captured failures | `tests/e2e/results/...` (nested via attachments) | Playwright failure screenshots |

## Error Handling

| Condition | Action |
|-----------|--------|
| Browser not installed | Run `just e2e-setup`, retry once |
| Server won't start | Report error, skip tests needing it |
| Health check fails | Report unreachable services, abort |
| Test timeout | Mark FAIL with timeout reason |
| node_modules missing | Run `just e2e-setup`, retry once |
| Spec compile error | Report TypeScript error, skip that spec |

## Failure Diagnosis

| Failure rate | Response |
|--------------|----------|
| >30% UI tests fail | App health problem — check screenshots for white/blank screen first |
| 10-30% | Spot-check 2-3 screenshots |
| <10% | Per-test fix tasks |

App health diagnostic order: failure screenshots -> dependency compatibility -> root component -> browser console.
