---
name: eval-test-cases
description: Evaluate test-cases.md for downstream executability with 1000-point scoring, then run adversarial iterations until target score is met.
argument-hints:
  - name: target
    description: Target score threshold (default: 900).
    required: false
  - name: iterations
    description: Max adversarial iterations (default: 6).
    required: false
---
Skill(skill="forge:eval", args="--type test-cases [--target N] [--iterations N]")
