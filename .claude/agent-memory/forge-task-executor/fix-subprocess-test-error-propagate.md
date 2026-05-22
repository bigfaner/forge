---
name: fix-subprocess-test-error-propagate
description: Subprocess tests must propagate errors via Exit(), not ignore with _ = when functions return AIErrors instead of calling os.Exit directly
metadata:
  type: reference
---

After the AIError/state machine refactoring, functions like `runSubmit`, `runReopen`, `runForensicExtract`, `runCheckDeps`, `runVerifyTaskDone`, `saveIndexAndSignalCompletion` return errors instead of calling `os.Exit()` directly. Subprocess-based tests that call these functions directly (not through Cobra's RunE) must propagate the error with `Exit(err)` instead of ignoring it with `_ = func()`.

For functions that use `panic()` instead of `Exit()` (like `validateQualityGate`), wrap in `recover/defer` to convert the panic to a proper `Exit()` call.

**Why:** The refactoring moved from direct `os.Exit()` calls to returning `*AIError` via Cobra's `RunE`. When tests bypass Cobra and call functions directly, the error must still propagate as a non-zero exit for the parent test process to detect it.

**How to apply:** In subprocess test stanzas (the `if os.Getenv("TEST_FLAG") == "1" { ... }` pattern), replace `_ = func()` with `if err := func(); err != nil { Exit(err) }`. For `panic`-based functions, wrap in defer/recover.