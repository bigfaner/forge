# E2E Test Report: forge-info-commands

**Date**: 2026-05-14
**Duration**: N/A (blocked)

## Summary

| Type  | Total | Pass | Fail | Skip |
|-------|-------|------|------|------|
| UI    | 0     | 0    | 0    | 0    |
| API   | 0     | 0    | 0    | 0    |
| CLI   | 0     | 0    | 0    | 0    |
| **All** | **0** | **0** | **0** | **0** |

**Result**: BLOCKED

---

## Blocked Reason

Tests could not be executed. Two prerequisites are missing:

1. **Justfile missing `e2e-setup` recipe**: The project `Justfile` does not contain an `e2e-setup` target. Run `/init-justfile` to scaffold the required targets (`e2e-setup`, `test-e2e`, `e2e-verify`).

2. **No test scripts generated**: The directory `tests/e2e/features/forge-info-commands/` does not exist. Run `/gen-test-scripts` first to generate executable test scripts.

---

## Failed Tests Detail

N/A

---

## Screenshots

N/A
