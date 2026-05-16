---
name: eval-consistency
description: Evaluate and fix cross-document consistency (PRD, Design, UI, Tasks). Detects inconsistencies via 1000-point scoring, then auto-fixes downstream docs to align with PRD as source of truth. Supports --scope docs|full.
argument-hints:
  - name: target
    description: Target score threshold (default: 900).
    required: false
  - name: iterations
    description: Max adversarial iterations (default: 3).
    required: false
  - name: scope
    description: Evaluation scope: docs (cross-document only) or full (docs + code). Defaults to docs.
    required: false
---
Skill(skill="forge:eval", args="--type consistency [--target N] [--iterations N] [--scope docs|full]")
