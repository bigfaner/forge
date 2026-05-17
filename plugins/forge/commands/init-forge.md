---
name: init-forge
description: Build and install the forge CLI tool.
allowed-tools: Bash Read
---

# /init-forge

Build and install the forge CLI tool.

## Process

1. Detect operating system (Windows/Linux/macOS)
2. Locate forge-cli path (forge-cli/)
3. Run the corresponding install script
4. Prompt user to reopen terminal

## Steps

### Step 1: Locate forge-cli

```bash
# forge-cli is in the forge-cli/ directory
FORGE_CLI_DIR="${CLAUDE_PROJECT_ROOT}/forge-cli"
if [ ! -d "$FORGE_CLI_DIR" ]; then
  echo "ERROR: forge-cli not found at $FORGE_CLI_DIR"
  exit 1
fi
echo "Found forge-cli at: $FORGE_CLI_DIR"
```

### Step 2: Detect OS and Install

**Windows (PowerShell):**
```powershell
cd $FORGE_CLI_DIR
powershell -ExecutionPolicy Bypass -File scripts/install-local.ps1
```

**Linux/macOS:**
```bash
cd "$FORGE_CLI_DIR" && bash scripts/install-local.sh
```

### Step 3: Verify Installation

```bash
forge --version
```

### Step 4: Prompt User

After installation, output:

```
╔════════════════════════════════════════════════════════════════╗
║  forge CLI installed successfully to ~/.forge/bin/            ║
╠════════════════════════════════════════════════════════════════╣
║  Run the following command to update PATH in current session: ║
║    source ~/.zshrc    (for zsh)                               ║
║    source ~/.bashrc   (for bash)                              ║
║  Then verify with: forge --version                            ║
╚════════════════════════════════════════════════════════════════╝
```

## Error Handling

| Error | Solution |
|-------|----------|
| forge-cli not found | Ensure forge-cli directory exists in the project root |
| Build failed | Check Go environment |
| Permission denied | Check write permissions for install directory |

<EXTREMELY-IMPORTANT>
- Do NOT modify any project source files — this command only builds and installs the forge CLI binary.
- Do NOT overwrite an existing forge CLI installation without user confirmation.
- If the build fails, report the error and stop. Do not attempt partial installs.
</EXTREMELY-IMPORTANT>
