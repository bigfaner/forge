---
step: 1
title: Feature Set
journey: feature-management
---

# Step 1: Feature Set

## Given
- A forge project directory (with or without .forge/state.json)
- Optional existing feature directories under docs/features/

## When
- `forge feature set <slug>` is executed
- `forge feature <slug>` (positional, legacy) is executed
- `forge feature` (query current) is executed
- `forge feature -v` (verbose query) is executed

## Then
- `set` creates feature directory structure and writes state.json
- Empty/whitespace slug is rejected
- Query returns current feature from state.json, falling back to features-dir scan
- Verbose query shows feature source (state.json or features-dir)
- Idempotent on repeated calls
- Positional arg does NOT write state.json

## Contract Dimensions
- **Actor**: CLI user setting or querying the current feature
- **Input**: slug string, optional state.json, optional feature directories
- **Output**: CLI stdout with FEATURE: <slug> or FEATURE: (none)
- **Side Effects**: .forge/state.json creation, docs/features/<slug>/ directory creation
- **Error Cases**: empty slug -> non-zero exit, corrupt state.json -> silent fallback
- **Invariants**: state.json takes priority over features-dir; nonexistent dir in state.json -> fallback
