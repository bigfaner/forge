---
name: eval-prd
description: Evaluate a PRD document with 1000-point scoring, then run adversarial iterations until target score is met. Main session orchestrates doc-scorer and doc-reviser subagents.
---

# /eval-prd

Evaluate a PRD document with iterative adversarial scoring.

Delegate to the generic eval skill:

```
Skill(skill="forge:eval", args="--type prd")
```
