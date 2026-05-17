# E2E Test Report: task-record-immutability

**Date**: 2026-05-17
**Duration**: 1.765s

## Summary

| Type  | Total | Pass | Fail | Skip |
|-------|-------|------|------|------|
| UI    | 0     | 0    | 0    | 0    |
| API   | 0     | 0    | 0    | 0    |
| CLI   | 12    | 12   | 0    | 0    |
| **All** | **12** | **12** | **0** | **0** |

**Result**: PASS

---

## Results by Test Case

| TC ID | Test Name | Status | Duration |
|-------|-----------|--------|----------|
| TC-001 | Submit record succeeds when no record exists | PASS | 0.11s |
| TC-002 | Submit record blocked when record already exists | PASS | 0.08s |
| TC-003 | Submit record with --force overwrites existing record | PASS | 0.11s |
| TC-004 | Default query shows 4 fields unchanged | PASS | 0.08s |
| TC-005 | Verbose query displays all task fields | PASS | 0.07s |
| TC-006 | Verbose query shows RELATED_FIXES for tasks with fix records | PASS | 0.06s |
| TC-007 | Verbose query omits RELATED_FIXES when no fixes exist | PASS | 0.05s |
| TC-008 | Status command behavior unchanged | PASS | 0.07s |
| TC-009 | Verbose query with short flag -v | PASS | 0.12s |
| TC-010 | Verbose query omits SCOPE when empty | PASS | 0.07s |
| TC-011 | Verbose query omits BREAKING when false | PASS | 0.09s |
| TC-012 | Verbose query displays multi-line DEPENDENCIES | PASS | 0.08s |

---

## Failed Tests Detail

No failed tests.

---

## Screenshots

No screenshots (CLI tests only).
