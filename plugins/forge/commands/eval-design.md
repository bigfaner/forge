---
name: eval-design
description: Evaluate a tech design document with 1000-point scoring, then run adversarial iterations until target score is met. Main session orchestrates doc-scorer and doc-reviser subagents.
---

# /eval-design

Evaluate a tech design document with iterative adversarial scoring.

Delegate to the generic eval skill:

```
Skill(skill="forge:eval", args="--type design")
```
