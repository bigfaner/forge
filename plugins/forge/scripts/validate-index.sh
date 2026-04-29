#!/bin/bash
# validate-index.sh
# PostToolUse hook to validate index.json files after Edit/Write

# Hook receives JSON input via stdin with structure:
# {"tool_name": "...", "tool_input": {"file_path": "..."}}
INPUT=$(cat)

# Extract file path from tool_input
FILE_PATH=$(echo "$INPUT" | jq -r '.tool_input.file_path // .tool_input.path // empty' 2>/dev/null)

# Skip if not an index.json file
if [[ -z "$FILE_PATH" ]] || [[ ! "$FILE_PATH" =~ index\.json$ ]]; then
    exit 0
fi

# Check if task CLI is available
if ! command -v task &> /dev/null; then
    echo "Warning: task CLI not found. Run /init-forge to install."
    exit 0
fi

# Run task validate
task validate "$FILE_PATH" 2>&1 || true

exit 0
