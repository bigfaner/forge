# Test Code Quality Gates

## Antipattern Guard (Post-Compile)

Verify each generated test function does not match any forbidden pattern:

| # | Forbidden Pattern | Instead |
|---|-------------------|--------|
| 1 | Recursive test invocation | Recursion guard (env var) or `-run` flag |
| 2 | Unconditional `t.Skip` | Implement with fixture or don't generate |
| 3 | Vacuous assertions (conditional assert without else fail) | Every assertion reachable on every code path |
| 4 | Environment-dependent skip without own fixture | `t.TempDir()` + own project structure |
| 5 | Duplicate test function names across packages | Scan for collisions; unique names with journey slug |
| 6 | Static-file text grep (assert on source file content) | Test runtime behavior only |

## Duplicate Name Check

Before writing, scan existing test files in the module for matching function names. If a collision is found, use a unique name that includes the journey slug.

## Error Handling

| Situation | Action |
|-----------|--------|
| No Convention files found | Proceed with LLM defaults + Code Reconnaissance. Output hint: "No test Convention files found in `docs/conventions/testing/`. Generation will use LLM defaults. Run `/forge:test-guide` to create one." |
| Convention file missing required sections | Proceed with LLM defaults for missing sections. Log warning listing missing sections. |
| Convention file unreadable | Skip file, log warning with file path and error. |
| Convention file has no required sections (`framework`, `discovery`, `structure`, `assertions`) | Skip file, log warning. |
| Convention vs Reconnaissance conflict | Convention wins, log conflict for user awareness. |
| Contract files not found | Abort with prompt to run `/gen-contracts` |
| Fact Table lookup fails for a descriptor | Keep `// VERIFY:` marker, do not fabricate regex |
| `just compile` recipe missing | Block generation. Output actionable error with recovery instructions. |
| Compile gate failed (all retries) | Block task. Output error + file path + recovery actions. Preserve generated files. |
| No test files generated | Abort with clear diagnostic message |
| Custom template path not found | Fall back to Convention file patterns with WARNING |
| Syntax validation failed (attempt 1) | Auto-retry: regenerate the failing file with error context |
| Syntax validation failed (attempt 2) | Mark file as `gen-failed`, skip in subsequent steps |
| Import path resolution failed | Same as syntax validation: retry once, then `gen-failed` |
| Surface type not in config.yaml | Auto-detect from code signals (Step 0.5.2), or ask user |
