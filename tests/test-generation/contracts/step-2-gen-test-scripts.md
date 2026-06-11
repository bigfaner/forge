---
step: 2
title: Generate Test Scripts with Convention Loading
journey: test-generation
---

# Step 2: Generate Test Scripts with Convention Loading

## Given
- A forge project with Convention files in docs/conventions/
- Convention files declare framework, assertion, tags, and result format
- A valid Journey directory structure with test types

## When
- `forge gen-test-scripts` is executed

## Then
- Convention files are loaded based on domain matching
- Framework selection respects Convention declarations over reconnaissance
- Missing/empty Convention sections trigger warnings with fallback to defaults
- Cold start (no Convention files) provides hints about test-guide
- Overlapping domain Conventions are merged with last-loaded-wins semantics
- Non-loadable Convention files (missing domains) are skipped with warning
- Unreadable/invalid-encoding files are skipped with warning

## Contract Dimensions
- **Actor**: CLI user or skill executing `forge gen-test-scripts`
- **Input**: Convention files in docs/conventions/, project structure, Journey definitions
- **Output**: Generated test script files, Convention loading logs/warnings
- **Side Effects**: Test file generation in appropriate Journey directories
- **Error Cases**: Missing Convention files, invalid Convention content, encoding errors, permission errors
- **Invariants**: Convention always wins over reconnaissance; generation completes within time budget (120s per Journey)
