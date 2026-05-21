---
step: 1
title: Removed Test Commands
journey: command-regression
---

# Step 1: Removed Test Commands

## Given
- Forge CLI binary built from source
- Profile-based test commands have been removed (forge test detect/get/interfaces/framework)

## When
- `forge test detect` is executed
- `forge test get` is executed
- `forge test interfaces` is executed
- `forge test framework` is executed

## Then
- Each removed command returns non-zero exit code
- Output contains "unknown command" or "command not found" or "removed" message

## Contract Dimensions
- **Actor**: CLI user attempting to use removed test commands
- **Input**: CLI args with removed subcommands
- **Output**: non-zero exit code, error message
- **Error Cases**: all invocations are error cases (commands removed)
- **Invariants**: removed commands never succeed; error message always identifies the issue
