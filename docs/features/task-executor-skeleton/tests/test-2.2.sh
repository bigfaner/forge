#!/usr/bin/env bash
# Test script for task 2.2: Add Execution Workflow to quick templates + update schemas
# Validates acceptance criteria from task definition

set -uo pipefail

PASS=0
FAIL=0
TOTAL=0

pass() { ((PASS++)); ((TOTAL++)); echo "  PASS: $1"; }
fail() { ((FAIL++)); ((TOTAL++)); echo "  FAIL: $1"; }

QUICK_TEMPLATES_DIR="plugins/forge/skills/quick-tasks/templates"
BREAKDOWN_SCHEMA="plugins/forge/skills/breakdown-tasks/templates/index.schema.json"
QUICK_SCHEMA="plugins/forge/skills/quick-tasks/templates/index.schema.json"

TEMPLATES=(
  "$QUICK_TEMPLATES_DIR/task.md"
  "$QUICK_TEMPLATES_DIR/quick-test-cases.md"
  "$QUICK_TEMPLATES_DIR/quick-gen-scripts.md"
  "$QUICK_TEMPLATES_DIR/quick-run-tests.md"
  "$QUICK_TEMPLATES_DIR/quick-graduate.md"
  "$QUICK_TEMPLATES_DIR/quick-verify-regression.md"
)

echo "=== Task 2.2 Acceptance Tests ==="
echo ""

# AC: All 6 quick templates have ## Execution Workflow section with non-empty content
echo "AC: All 6 quick templates have ## Execution Workflow section..."
for tmpl in "${TEMPLATES[@]}"; do
  basename=$(basename "$tmpl")
  if grep -q "^## Execution Workflow" "$tmpl"; then
    pass "$basename has ## Execution Workflow heading"
  else
    fail "$basename missing ## Execution Workflow heading"
  fi

  # Check workflow has non-empty content (at least one numbered step)
  # Use sed to get content after ## Execution Workflow to end of file
  if sed -n '/^## Execution Workflow/,$ p' "$tmpl" | grep -qE "^\s*[0-9]+\."; then
    pass "$basename Execution Workflow has numbered steps"
  else
    fail "$basename Execution Workflow is empty or has no numbered steps"
  fi
done

# AC: All 6 quick templates have no noTest in frontmatter
echo ""
echo "AC: All 6 quick templates have no noTest in frontmatter..."
for tmpl in "${TEMPLATES[@]}"; do
  basename=$(basename "$tmpl")
  # Check frontmatter only (between --- markers)
  if awk '/^---$/{n++; next} n==1' "$tmpl" | grep -qi "noTest"; then
    fail "$basename still has noTest in frontmatter"
  else
    pass "$basename has no noTest in frontmatter"
  fi
done

# AC: Both index.schema.json files have no noTest field definition
echo ""
echo "AC: Both schemas have no noTest field definition..."
for schema in "$BREAKDOWN_SCHEMA" "$QUICK_SCHEMA"; do
  basename=$(basename "$schema")
  dirname=$(basename "$(dirname "$(dirname "$schema")")")
  label="$dirname/$basename"
  if grep -q '"noTest"' "$schema"; then
    fail "$label still has noTest field"
  else
    pass "$label has no noTest field"
  fi
done

# AC: grep -r "noTest" quick-tasks/templates/ in .md files -> zero matches
echo ""
echo "AC: grep -r noTest in quick-tasks .md templates -> zero matches..."
count=$(grep -rl "noTest" "$QUICK_TEMPLATES_DIR"/*.md 2>/dev/null | wc -l)
if [ "$count" -eq 0 ]; then
  pass "No noTest matches in quick-tasks .md templates"
else
  fail "Found noTest in quick-tasks .md templates ($count files)"
fi

# AC: Schemas are valid JSON
echo ""
echo "AC: Schema files are valid JSON..."
for schema in "$BREAKDOWN_SCHEMA" "$QUICK_SCHEMA"; do
  basename=$(basename "$schema")
  dirname=$(basename "$(dirname "$(dirname "$schema")")")
  label="$dirname/$basename"
  if python3 -c "import json; json.load(open('$schema'))" 2>/dev/null; then
    pass "$label is valid JSON"
  elif node -e "JSON.parse(require('fs').readFileSync('$schema','utf8'))" 2>/dev/null; then
    pass "$label is valid JSON"
  else
    fail "$label is NOT valid JSON"
  fi
done

# Summary
echo ""
echo "=== Results: $PASS passed, $FAIL failed (out of $TOTAL) ==="

if [ "$FAIL" -gt 0 ]; then
    exit 1
fi
