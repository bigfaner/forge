---
name: eval-test-cases
description: Evaluate test-cases.md for downstream executability with 1000-point scoring, then run adversarial iterations until target score is met. Main session orchestrates doc-scorer and doc-reviser subagents.
---

# /eval-test-cases

Evaluate test cases for downstream executability with iterative adversarial scoring.

Delegate to the generic eval skill:

```
Skill(skill="forge:eval", args="--type test-cases")
```
