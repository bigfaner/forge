---
name: init-forge
description: Build and install the task-cli tool.
---

# /init-forge

Build and install the task-cli tool.

## Process

1. Detect operating system (Windows/Linux/macOS)
2. Locate task-cli path (task-cli/)
3. Run the corresponding install script
4. Prompt user to reopen terminal

## Steps

### Step 1: Locate task-cli

```bash
# task-cli is in the task-cli/ directory
TASK_CLI_DIR="${CLAUDE_PROJECT_ROOT}/task-cli"
if [ ! -d "$TASK_CLI_DIR" ]; then
  echo "ERROR: task-cli not found at $TASK_CLI_DIR"
  exit 1
fi
echo "Found task-cli at: $TASK_CLI_DIR"
```

### Step 2: Detect OS and Install

**Windows (PowerShell):**
```powershell
cd $TASK_CLI_DIR
powershell -ExecutionPolicy Bypass -File scripts/install-local.ps1
```

**Linux/macOS:**
```bash
cd "$TASK_CLI_DIR" && bash scripts/install-local.sh
```

### Step 3: Verify Installation

```bash
task --version
```

### Step 4: Prompt User

After installation, output:

```
╔════════════════════════════════════════════════════════════════╗
║  task-cli installed successfully                              ║
╠════════════════════════════════════════════════════════════════╣
║  Please reopen your terminal to refresh PATH, then run:       ║
║  task --version                                               ║
╚════════════════════════════════════════════════════════════════╝
```

## Error Handling

| Error | Solution |
|-------|----------|
| task-cli not found | Ensure task-cli directory exists in the project root |
| Build failed | Check Go environment |
| Permission denied | Check write permissions for install directory |
