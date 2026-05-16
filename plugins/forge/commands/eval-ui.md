---
name: eval-ui
description: Evaluate a UI design document with 1000-point scoring from four stakeholder perspectives (User/Designer/Developer/PM), then run adversarial iterations until target score is met.
argument-hints:
  - name: target
    description: Target score threshold (default: 950).
    required: false
  - name: iterations
    description: Max adversarial iterations (default: 3).
    required: false
---
Skill(skill="forge:eval", args="--type ui [--target N] [--iterations N]")
