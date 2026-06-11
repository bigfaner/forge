---
step: 1
title: Task Type Constants and Validation
journey: task-type-system
---

# Step 1: Task Type Constants and Validation

## Given
- Forge CLI binary built from source

## When
- `forge task list-types` is executed
- `forge task validate-index` is executed on an index with new type values

## Then
- list-types displays feature, enhancement, cleanup, refactor types
- Deprecated implementation type is still shown with deprecation notice
- validate-index accepts index.json with coding.feature, coding.enhancement, coding.cleanup, coding.refactor types

## Contract Dimensions
- **Actor**: CLI user listing or validating task types
- **Input**: CLI args, optional index.json path
- **Output**: stdout with type list or validation results
- **Error Cases**: invalid type in index -> validation error
- **Invariants**: implementation type shown as deprecated; new types always listed
