---
name: eval-proposal
description: Evaluate a proposal document with 1000-point scoring, then run adversarial iterations until target score is met. Specify target score and max iterations.
argument-hints:
  - name: target
    description: Target score threshold (default: 900).
    required: false
  - name: iterations
    description: Max adversarial iterations (default: 3).
    required: false
---
Skill(skill="forge:eval", args="--type proposal [--target N] [--iterations N]")
