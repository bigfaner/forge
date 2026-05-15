---
name: eval-proposal
description: Evaluate a proposal document with 1000-point scoring, then run adversarial iterations until target score is met. Main session orchestrates doc-scorer and doc-reviser subagents. Specify target score and max iterations.
---

# /eval-proposal

Evaluate a proposal document with iterative adversarial scoring.

Delegate to the generic eval skill:

```
Skill(skill="forge:eval", args="--type proposal")
```
