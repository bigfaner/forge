---
name: eval-consistency
description: Evaluate and fix cross-document consistency (PRD, Design, UI, Tasks). Detects inconsistencies via 1000-point scoring, then auto-fixes downstream docs to align with PRD as source of truth. Supports --scope docs|full.
argument-hint: "[--target 900] [--iterations 3] [--scope docs|full]"
---
Skill(skill="forge:eval", args="--type consistency [--target N] [--iterations N] [--scope docs|full]")
