---
step: 1
title: Task Index with Test Pipeline Generation
journey: test-generation
---

# Step 1: Task Index with Test Pipeline Generation

## Given
- A forge project with a feature that has tasks in tasks/ directory
- .forge/config.yaml with language profiles (e.g., go, javascript)
- Test cases document (test-cases.md) or profile manifest capabilities

## When
- `forge task index --feature <slug>` is executed

## Then
- Test pipeline tasks are generated in index.json
- Per-type gen-scripts tasks created based on profile capabilities
- Quick mode: gen-cases + per-type gen-and-run + graduate + verify + drift
- Breakdown mode: gen-cases + eval + per-type gen-scripts + run + graduate + verify
- Multi-profile: letter-suffixed task IDs (a, b) for each profile
- Dependency chain is correct (gen -> per-type -> graduate -> verify -> drift)

## Contract Dimensions
- **Actor**: CLI user executing `forge task index --feature <slug>`
- **Input**: Feature directory with tasks/, .forge/config.yaml, optional test-cases.md
- **Output**: Updated index.json with test pipeline tasks, generated task .md files
- **Side Effects**: index.json creation/update, task .md file generation
- **Error Cases**: Missing feature directory, invalid config.yaml
- **Invariants**: Idempotent re-runs produce same task set; shared tasks (gen-cases) not per-type
