---
step: 2
title: Task Submit & Record
journey: task-lifecycle
---

# Step 2: Task Submit & Record

## Given
- A forge project with tasks in index.json
- Some tasks have existing records, some do not

## When
- `forge task submit <id> --data <file>` is executed
- `forge task query <id>` is executed (default or --verbose)
- `forge task status <id>` is executed

## Then
- Submit creates record file, blocked if record exists (unless --force)
- Query shows TASK_ID, STATUS, SCOPE in default mode
- Query --verbose shows all fields (KEY, TITLE, PRIORITY, TYPE, DEPENDENCIES, TASK_FILE, RECORD_FILE, RELATED_FIXES)
- Status command shows TASK_ID and STATUS fields
- -v is equivalent to --verbose

## Contract Dimensions
- **Actor**: CLI user managing task records and querying task state
- **Input**: index.json, optional record files, CLI flags (--force, --verbose, -v)
- **Output**: CLI stdout with structured --- block format
- **Side Effects**: record file creation in docs/features/<slug>/tasks/records/
- **Error Cases**: existing record without --force -> exit 1 with hint message
- **Invariants**: SCOPE omitted when empty, BREAKING omitted when false, RELATED_FIXES omitted when none
