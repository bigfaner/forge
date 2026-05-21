---
step: 3
title: Type Reclassification and Migration
journey: task-type-system
---

# Step 3: Type Reclassification and Migration

## Given
- A forge project with tasks of various types
- Tasks may have type different from actual work performed

## When
- `forge task submit` is executed with typeReclassification data
- `forge prompt get-by-task-id` is executed for typed tasks
- `forge task migrate` is executed

## Then
- Submit with reclassification includes Type Reclassification section in record
- Submit without reclassification omits Type Reclassification section
- Prompt templates are selected based on task type (feature/cleanup/refactor)
- Migrate maps deprecated implementation type to feature
- Quality gate skips cleanup-only features (docs-only logic)
- Dynamic fix tasks get correct type based on failing step

## Contract Dimensions
- **Actor**: CLI user submitting records, getting prompts, or migrating types
- **Input**: task ID, record data with optional reclassification, config with feature
- **Output**: record markdown with optional reclassification, prompt template text, migrated index.json
- **Side Effects**: record files written, index.json types updated
- **Invariants**: reclassification only appears when original != actual type
