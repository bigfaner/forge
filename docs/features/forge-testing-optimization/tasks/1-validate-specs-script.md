---
id: "1"
title: "Create validate-specs.mjs validation script + update package.json"
priority: "P0"
estimated_time: "2h"
dependencies: []
status: pending
breaking: false
noTest: false
mainSession: false
---

# 1: Create validate-specs.mjs validation script + update package.json

## Description

Create a ts-morph-based AST validation script (`validate-specs.mjs`) that checks generated spec files against 8 structural rules (4 ERROR + 4 WARNING). This is the core programmatic validation artifact that replaces the current reliance on LLM compliance with SKILL.md rules.

The script parses spec files using ts-morph AST analysis and produces structured output (JSON) that downstream tooling can consume.

## Reference Files
- `docs/proposals/forge-testing-optimization/proposal.md` — Source proposal (Phase 2, Section 2.1)
- `plugins/forge/skills/gen-test-scripts/templates/helpers.ts` — Existing template for reference on coding patterns
- `plugins/forge/skills/gen-test-scripts/templates/playwright-ui.spec.ts` — Existing spec template (validation targets)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/gen-test-scripts/templates/validate-specs.mjs` | ts-morph AST validation script with 8 rules |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-test-scripts/templates/package.json` | Add `ts-morph` as devDependency |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] `validate-specs.mjs` detects E1: `waitForTimeout` / `setTimeout` usage via AST CallExpression analysis
- [ ] `validate-specs.mjs` detects E2: TC ID full coverage — all `TC-\d+` from test-cases.md must appear in spec files
- [ ] `validate-specs.mjs` detects E3: every `test()` has a Traceability comment (`// Traceability:`)
- [ ] `validate-specs.mjs` detects E4: no DOM parent traversal `locator('..')`
- [ ] `validate-specs.mjs` detects W1: serial suite with >15 `test()` calls
- [ ] `validate-specs.mjs` detects W2: serial suite without `afterAll` cleanup
- [ ] `validate-specs.mjs` detects W3: `beforeEach` containing login/loginViaUI calls
- [ ] `validate-specs.mjs` detects W4: CSS class selectors (`.xxx` pattern in locator strings)
- [ ] Script outputs structured JSON with `{ errors: [{id, rule, file, line, message}], warnings: [...] }`
- [ ] Exit code: 0 if no errors (warnings OK), 1 if any errors present
- [ ] `ts-morph` listed as devDependency in `package.json` template
- [ ] Script handles missing files gracefully (reports ERROR, doesn't crash)

## Implementation Notes

1. **ts-morph version**: Pin to a version compatible with the TypeScript version in `package.json` template. Use `ts-morph` ^21.x which supports TS 5.x
2. **Script location**: Place in `gen-test-scripts/templates/` so it gets copied alongside helpers.ts to `tests/e2e/`
3. **CLI interface**: Accept spec directory as first arg, optional `--test-cases` path for E2 TC ID coverage check
4. **E2 TC ID coverage**: Read test-cases.md, extract all `TC-\d+` patterns, grep spec files for matches. Report missing IDs
5. **Fallback on ts-morph failure**: If AST parsing fails (e.g., invalid TS), catch the error and report as WARNING rather than ERROR — per the risk mitigation in proposal
6. **Windows path handling**: Use `path.join()` for all file paths, avoid hardcoded separators
7. **Rule E3 Traceability**: Check for `// Traceability:` comment either as a leading comment above the `test()` call or as a string literal inside the test callback
