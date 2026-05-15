---
name: eval-consistency
description: Evaluate and fix cross-document consistency (PRD, Design, UI, Tasks). Detects inconsistencies via 1000-point scoring, then auto-fixes downstream docs to align with PRD as source of truth. Supports --scope docs|full. Uses doc-scorer and doc-reviser subagents.
---

# /eval-consistency

Evaluate and fix cross-document consistency with iterative adversarial scoring.

Delegate to the generic eval skill:

```
Skill(skill="forge:eval", args="--type consistency")
```
