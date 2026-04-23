#!/bin/bash
# validate-manifest.sh
# Validates that docs/decisions/manifest.md counts match actual row counts in type files

DECISIONS_DIR="docs/decisions"
MANIFEST="$DECISIONS_DIR/manifest.md"
ERRORS=0

CATEGORIES=(
  "Architecture"
  "Interface"
  "Data Model"
  "Dependencies"
  "Error Handling"
  "Testing"
  "Security"
  "Local Dev & Deployment"
)

FILES=(
  "architecture.md"
  "interface.md"
  "data-model.md"
  "dependencies.md"
  "error-handling.md"
  "testing.md"
  "security.md"
  "local-dev-deployment.md"
)

for i in "${!CATEGORIES[@]}"; do
  category="${CATEGORIES[$i]}"
  file="$DECISIONS_DIR/${FILES[$i]}"

  # count actual data rows (total | lines minus header and separator = -2)
  raw=$(grep -c "^|" "$file" 2>/dev/null || echo 0)
  actual=$((raw - 2))
  [ "$actual" -lt 0 ] && actual=0

  # read manifest count for this category
  manifest_count=$(grep "| $category |" "$MANIFEST" | awk -F'|' '{print $4}' | tr -d ' ')

  if [ "$actual" != "$manifest_count" ]; then
    echo "MISMATCH: $category — actual=$actual, manifest=$manifest_count"
    ERRORS=$((ERRORS + 1))
  fi
done

if [ "$ERRORS" -eq 0 ]; then
  echo "manifest OK"
  exit 0
else
  exit 1
fi
