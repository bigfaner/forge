---
step: 2
title: Type-Driven Pipeline Generation
journey: task-type-system
---

# Step 2: Type-Driven Pipeline Generation

## Given
- A forge project with task-type-refinement feature
- Tasks with various type values in index.json

## When
- `forge task index --feature <slug>` is executed

## Then
- Feature/enhancement/fix tasks trigger T-quick-* test pipeline generation
- Cleanup-only features skip test pipeline
- Refactor-only features skip test pipeline
- Documentation-only features get T-eval-doc task instead of pipeline
- Mixed cleanup+refactor features get neither pipeline nor eval-doc

## Contract Dimensions
- **Actor**: CLI user building task index
- **Input**: feature slug, index.json with typed tasks
- **Output**: updated index.json with auto-generated test pipeline tasks
- **Side Effects**: index.json updated with T-quick-*, T-test-*, T-eval-doc tasks
- **Invariants**: cleanup and refactor types never trigger pipeline; documentation triggers eval only
