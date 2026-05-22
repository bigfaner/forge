---
id: "fix.8"
title: "Fix TestAutoConfigWithDefaults after WithDefaults() fix"
priority: "P1"
dependencies:
  - "2.gate"
status: 
type: "coding.fix"
---

# fix.8: Fix TestAutoConfigWithDefaults after WithDefaults() fix

Test TestAutoConfigWithDefaults/partial_preserves_set_values fails because it tests the OLD WithDefaults() behavior. After the fix, WithDefaults() only handles all-zero configs. Update the test expectation: partial configs now return unchanged.
