---
title: "TypeScript Vitest Testing Convention"
---

# TypeScript Vitest Testing Convention

Convention for generating TypeScript/JavaScript test code using the Vitest framework.

## framework

- **name**: Vitest
- **version**: vitest 0.34+
- **language**: TypeScript
- **runner_command**: `vitest run --reporter=verbose`

## discovery

- **test_dir**: `tests/e2e/`
- **file_pattern**: `*.test.ts`, `*.spec.ts`
- **exclude_pattern**: `node_modules/`, `dist/`

## structure

- **suite_pattern**: `describe('...', () => { ... })` — BDD-style test container
- **case_pattern**: `it('should ...', () => { ... })` — individual test case within describe
- **hook_pattern**: `beforeAll` / `afterAll` / `beforeEach` / `afterEach`

### Test Structure

Use `describe` / `it` (BDD style):

```typescript
describe('Feature: Task Lifecycle', () => {
  let projectDir: string

  beforeAll(() => {
    projectDir = setupTestProject()
  })

  afterAll(() => {
    rmSync(projectDir, { recursive: true, force: true })
  })

  describe('Task claiming', () => {
    it('should claim a task successfully', () => {
      const result = runCLI('task', 'claim')
      expect(result.stdout).toContain('claimed task')
      expect(result.exitCode).toBe(0)
    })

    it('should fail when no tasks available', () => {
      const result = runCLI('task', 'claim')
      expect(result.exitCode).toBe(1)
      expect(result.stderr).toContain('no tasks available')
    })
  })
})
```

### Table-Driven Tests

Use `describe.each` / `it.each`:

```typescript
describe.each([
  { input: 'hello', expected: 'HELLO' },
  { input: '', expected: '' },
])('uppercase: $input', ({ input, expected }) => {
  it('should transform correctly', () => {
    expect(input.toUpperCase()).toBe(expected)
  })
})
```

Or `it.each`:

```typescript
it.each([
  { input: 'hello', expected: 'HELLO' },
  { input: '', expected: '' },
])('should transform $input', ({ input, expected }) => {
  expect(input.toUpperCase()).toBe(expected)
})
```

### CLI Testing

Use `child_process` for CLI invocation:

```typescript
function runCLI(...args: string[]): { stdout: string; stderr: string; exitCode: number } {
  const result = execSync(`forge ${args.join(' ')}`, {
    encoding: 'utf-8',
    env: { ...process.env, CLAUDE_PROJECT_DIR: projectDir },
  })
  return { stdout: result, stderr: '', exitCode: 0 }
}
```

### API Testing

Use built-in `fetch` (Node 18+) or `node-fetch`:

```typescript
it('should return OK from API endpoint', async () => {
  const response = await fetch('http://localhost:8080/api/resource')
  expect(response.status).toBe(200)
  const body = await response.json()
  expect(body.data).toBeDefined()
})
```

### Traceability

Each test should include a traceability comment:

```typescript
it('should login with valid credentials', () => {
  // Traceability: TC-001 -> PRD User Auth section
})
```

## assertions

- **style**: expect
- **library**: Vitest built-in `expect` (compatible with Jest matchers)
- **custom_matchers**: none

### Key Functions

- `expect(actual).toBe(expected)` — strict equality
- `expect(actual).toEqual(expected)` — deep equality
- `expect(actual).toContain(substr)` — substring/item in collection
- `expect(actual).toBeNull()` — null check
- `expect(actual).toBeDefined()` — defined check
- `expect(actual).toBeUndefined()` — undefined check
- `expect(actual).toBeTruthy()` — truthy check
- `expect(actual).toBeFalsy()` — falsy check
- `expect(fn).toThrow(error?)` — exception check
- `expect(actual).rejects.toThrow(error?)` — async rejection
- `expect(actual).resolves.toBe(expected)` — async resolution
- `expect(actual).toHaveLength(n)` — length check
- `expect(actual).toMatch(regex)` — regex match
- `expect(actual).toMatchSnapshot()` — snapshot match

**Rule**: Use built-in `expect`, do not import external assertion libraries.

## Tags

- **Format**: Vitest `describe`/`it` with tag metadata

```typescript
describe('@feature', () => {
  it('should perform action', { tags: ['@feature'] }, () => {
    // test code
  })
})
```

- **Tag annotation**: Use the `tags` option in `it()` or `describe()` for categorization
- **CLI filtering**: `vitest run --reporter=verbose --testNamePattern="pattern"` or `vitest run -t "pattern"`

## Result Format

- **Output flags**: `--reporter=verbose` or `--reporter=json`
- **Format type**: `json-report` (structured JSON output)

### JSON Report Structure

```json
{
  "testResults": [
    {
      "name": "tests/e2e/feature.test.ts",
      "status": "passed",
      "startTime": 1716000000000,
      "endTime": 1716000001000,
      "assertionResults": [
        {
          "fullName": "Feature > should perform action",
          "status": "passed",
          "duration": 100,
          "failureMessages": []
        }
      ]
    }
  ]
}
```

## Import Patterns

Standard imports for Vitest e2e tests:

```typescript
import { describe, it, expect, beforeAll, afterAll, beforeEach, afterEach } from 'vitest'
import { execSync, exec } from 'child_process'
import { readFileSync, writeFileSync, mkdirSync, rmSync } from 'fs'
import { join } from 'path'
```

- HTTP tests add: `import fetch from 'node-fetch'` or built-in `fetch` (Node 18+)
- API tests add: `import { setupServer } from 'msw/node'` (mock service worker)

## Anti-patterns (Forbidden)

| Forbidden | Replacement |
|-----------|-------------|
| `setTimeout` / `sleep` for synchronization | `waitFor()` from `@testing-library/dom` or async retry |
| Jest imports (`@jest/globals`) | Vitest imports (`from 'vitest'`) |
| Hardcoded ports | Dynamic port allocation or environment variables |
| Real secrets/tokens in code | `process.env.E2E_API_TOKEN` |
| `test()` instead of `it()` | Use `it()` for consistency with BDD style |
| Mixed test runners | Use only Vitest, never Jest or Mocha |
| Unconditional `it.skip` | Implement properly or don't generate |
| `done` callback parameter | Use async/await pattern |

## Helpers

### runCLI helper

```typescript
import { execSync } from 'child_process'

interface CLIResult {
  stdout: string
  stderr: string
  exitCode: number
}

function runCLI(...args: string[]): CLIResult {
  try {
    const stdout = execSync(`forge ${args.join(' ')}`, {
      encoding: 'utf-8',
      env: { ...process.env, CLAUDE_PROJECT_DIR: projectDir },
    })
    return { stdout, stderr: '', exitCode: 0 }
  } catch (error: any) {
    return {
      stdout: error.stdout?.toString() ?? '',
      stderr: error.stderr?.toString() ?? '',
      exitCode: error.status ?? 1,
    }
  }
}
```

### waitFor helper

```typescript
async function waitFor(
  condition: () => boolean,
  { timeout = 5000, interval = 100 }: { timeout?: number; interval?: number } = {}
): Promise<void> {
  const start = Date.now()
  while (Date.now() - start < timeout) {
    if (condition()) return
    await new Promise(resolve => setTimeout(resolve, interval))
  }
  throw new Error(`waitFor timed out after ${timeout}ms`)
}
```

### setupTestProject helper

```typescript
import { mkdtempSync, mkdirSync, writeFileSync } from 'fs'
import { join } from 'path'
import { tmpdir } from 'os'

function setupTestProject(): string {
  const dir = mkdtempSync(join(tmpdir(), 'forge-e2e-'))
  mkdirSync(join(dir, '.forge'), { recursive: true })
  writeFileSync(join(dir, '.forge', 'config.yaml'), '{}')
  return dir
}
```
