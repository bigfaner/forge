---
feature: "forge-testing-optimization"
sources:
  - docs/proposals/forge-testing-optimization/proposal.md
generated: "2026-05-10"
---

# Test Cases: forge-testing-optimization

> **WARNING**: sitemap.json not found — Element set to `sitemap-missing`. Run `/gen-sitemap` for precise element references.

## Summary

| Type | Count |
|------|-------|
| CLI  | 7     |
| **Total** | **7** |

---

## CLI Test Cases

### TC-001: validate-specs detects E1-E4 ERROR rules

- **Source**: Success Criterion — "validate-specs.mjs 脚本能检测 E1-E4 四种 ERROR 和 W1-W4 四种 WARNING"
- **Type**: CLI
- **Target**: cli/validate-specs
- **Test ID**: cli/validate-specs/detect-e1-e4-error-rules
- **Pre-conditions**: A spec file containing deliberate violations of E1 (waitForTimeout/setTimeout), E2 (missing TC IDs), E3 (missing Traceability comments), and E4 (DOM parent traversal `locator('..')`)
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Create a fixture spec file with `page.waitForTimeout(5000)` call (E1 violation)
  2. Create a fixture spec file with `setTimeout(() => {}, 1000)` call (E1 violation)
  3. Create a fixture test-cases.md with TC-999, but fixture spec omits TC-999 (E2 violation)
  4. Create a fixture spec file with `test('foo', () => {})` lacking `// Traceability:` comment above or within (E3 violation)
  5. Create a fixture spec file with `locator('..')` DOM parent traversal (E4 violation)
  6. Run `node validate-specs.mjs` against all fixture files
  7. Verify exit code is non-zero (ERROR level blocks)
  8. Verify output contains ERROR entries for E1, E2, E3, and E4
- **Expected**: All four ERROR types (E1-E4) are detected and reported. Script exits with non-zero code indicating blocking errors.
- **Priority**: P0

### TC-002: validate-specs detects W1-W4 WARNING rules

- **Source**: Success Criterion — "validate-specs.mjs 脚本能检测 E1-E4 四种 ERROR 和 W1-W4 四种 WARNING"
- **Type**: CLI
- **Target**: cli/validate-specs
- **Test ID**: cli/validate-specs/detect-w1-w4-warning-rules
- **Pre-conditions**: A spec file containing deliberate violations of W1 (serial suite > 15 tests), W2 (serial suite without afterAll), W3 (beforeEach with login call), and W4 (CSS class selector)
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Create a fixture spec file with a `test.describe.serial('big suite', () => { ... })` containing 20+ `test()` calls (W1 violation)
  2. Create a fixture spec file with a serial describe that has no `afterAll()` call (W2 violation)
  3. Create a fixture spec file with `beforeEach(() => { login(...) })` pattern (W3 violation)
  4. Create a fixture spec file with `locator('.ant-btn')` CSS class selector (W4 violation)
  5. Run `node validate-specs.mjs` against all fixture files
  6. Verify output contains WARNING entries for W1, W2, W3, and W4
  7. Verify exit code is zero (WARNING does not block)
- **Expected**: All four WARNING types (W1-W4) are detected and reported as non-blocking warnings. Script exits with zero code.
- **Priority**: P0

### TC-003: ts-morph devDependency in package.json

- **Source**: Success Criterion — "ts-morph 在 tests/e2e/package.json 中作为 devDependency 存在"
- **Type**: CLI
- **Target**: cli/package-json
- **Test ID**: cli/package-json/ts-morph-dev-dependency-present
- **Pre-conditions**: gen-test-scripts templates package.json exists
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read `gen-test-scripts/templates/package.json`
  2. Check that `devDependencies` object contains a `ts-morph` key
  3. Verify `ts-morph` version string is a valid semver range (not empty, not "*")
- **Expected**: `ts-morph` is listed as a devDependency with a valid version in the templates package.json.
- **Priority**: P0

### TC-004: task validate-specs command executes and returns structured output

- **Source**: Success Criterion — "`task validate-specs` 命令能执行校验并返回结构化输出"
- **Type**: CLI
- **Target**: cli/task-validate-specs
- **Test ID**: cli/task-validate-specs/executes-and-returns-structured-output
- **Pre-conditions**: task-cli binary is built and installed; validate-specs.mjs script exists; a target spec file exists for validation
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Build task-cli binary (`go build` in task-cli directory)
  2. Create a minimal spec fixture file for validation
  3. Run `task validate-specs --path <fixture-dir>`
  4. Verify command exits without error when specs are clean
  5. Create a spec fixture with known E1 violation
  6. Run `task validate-specs --path <fixture-dir>` again
  7. Verify output is structured (JSON or key-value format with rule ID, severity, file, and message fields)
  8. Verify the E1 violation is reported in structured output
- **Expected**: `task validate-specs` runs successfully, returns structured output with rule ID, severity, file path, and message for each finding. Clean specs produce success output; dirty specs report violations.
- **Priority**: P0

### TC-005: gen-test-scripts SKILL.md contains Step 4.5 structural validation

- **Source**: Success Criterion — "gen-test-scripts SKILL.md 包含 Step 4.5 结构校验步骤"
- **Type**: CLI
- **Target**: cli/gen-test-scripts-skill
- **Test ID**: cli/gen-test-scripts-skill/step-4-5-structural-validation
- **Pre-conditions**: gen-test-scripts SKILL.md exists
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read `plugins/forge/skills/gen-test-scripts/SKILL.md`
  2. Search for "Step 4.5" section heading
  3. Verify the section describes structural validation using `task validate-specs`
  4. Verify the section specifies that ERROR results block downstream flow (mark T-test-2 as blocked)
  5. Verify the section specifies that WARNING results are reported but non-blocking
- **Expected**: SKILL.md contains a Step 4.5 that describes calling `task validate-specs`, handling ERROR as blocking and WARNING as non-blocking.
- **Priority**: P1

### TC-006: gen-test-scripts aborts when eval-test-cases Step Actionability below 20

- **Source**: Success Criterion — "gen-test-scripts 在 eval-test-cases Step Actionability < 20 时中止"
- **Type**: CLI
- **Target**: cli/gen-test-scripts-skill
- **Test ID**: cli/gen-test-scripts-skill/abort-on-low-step-actionability
- **Pre-conditions**: gen-test-scripts SKILL.md exists; an eval-test-cases report exists with Step Actionability < 20
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read `plugins/forge/skills/gen-test-scripts/SKILL.md`
  2. Locate the Prerequisites section
  3. Verify there is a check for eval-test-cases report existence
  4. Verify the check reads the Step Actionability score from the report
  5. Verify the check specifies: if Step Actionability < 20, abort gen-test-scripts and prompt user to fix test-cases.md
- **Expected**: SKILL.md Prerequisites section includes a Step Actionability threshold check that aborts execution when the score is below 20.
- **Priority**: P1

### TC-007: gen-test-cases Element field marked as required

- **Source**: Success Criterion — "gen-test-cases SKILL.md 和模板中 Element 字段标记为必填"
- **Type**: CLI
- **Target**: cli/gen-test-cases-skill
- **Test ID**: cli/gen-test-cases-skill/element-field-required
- **Pre-conditions**: gen-test-cases SKILL.md and template exist
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read `plugins/forge/skills/gen-test-cases/SKILL.md`
  2. Search for "Element" in context of required/mandatory/必填
  3. Verify the SKILL.md states Element field is required (HARD-RULE or equivalent)
  4. Verify the template at `plugins/forge/skills/gen-test-cases/templates/test-cases.md` includes Element field
  5. Verify the behavior when sitemap.json is missing is defined (fallback to `sitemap-missing`)
- **Expected**: gen-test-cases SKILL.md and template both define Element as a required field with clear fallback behavior when sitemap is absent.
- **Priority**: P1

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | SC: "validate-specs.mjs 脚本能检测 E1-E4 四种 ERROR" | CLI | cli/validate-specs | P0 |
| TC-002 | SC: "validate-specs.mjs 脚本能检测 W1-W4 四种 WARNING" | CLI | cli/validate-specs | P0 |
| TC-003 | SC: "ts-morph 在 tests/e2e/package.json 中作为 devDependency 存在" | CLI | cli/package-json | P0 |
| TC-004 | SC: "`task validate-specs` 命令能执行校验并返回结构化输出" | CLI | cli/task-validate-specs | P0 |
| TC-005 | SC: "gen-test-scripts SKILL.md 包含 Step 4.5 结构校验步骤" | CLI | cli/gen-test-scripts-skill | P1 |
| TC-006 | SC: "gen-test-scripts 在 eval-test-cases Step Actionability < 20 时中止" | CLI | cli/gen-test-scripts-skill | P1 |
| TC-007 | SC: "gen-test-cases SKILL.md 和模板中 Element 字段标记为必填" | CLI | cli/gen-test-cases-skill | P1 |

---

## Route Validation

_Omitted — this feature is a forge plugin (CLI/SKILL definitions), not a web application. No route files to validate._
