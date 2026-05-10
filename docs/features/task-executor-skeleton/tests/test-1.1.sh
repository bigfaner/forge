#!/usr/bin/env bash
# Test script for task 1.1: Rewrite task-executor.md as workflow skeleton
# Validates acceptance criteria from task definition

set -uo pipefail

PASS=0
FAIL=0
TOTAL=0

pass() { ((PASS++)); ((TOTAL++)); echo "  PASS: $1"; }
fail() { ((FAIL++)); ((TOTAL++)); echo "  FAIL: $1"; }

AGENT_FILE="plugins/forge/agents/task-executor.md"
CMD_FILE="plugins/forge/commands/execute-task.md"

echo "=== Task 1.1 Acceptance Tests ==="
echo ""

# AC: Step 2 contains zero TDD-specific instructions
echo "AC: No TDD-specific instructions in Step 2..."
if grep -qiE "RED.*GREEN.*REFACTOR|TDD cycle|TDD Implementation" "$AGENT_FILE"; then
    fail "task-executor.md Step 2 still contains TDD keywords"
else
    pass "task-executor.md Step 2 has no TDD keywords"
fi

if grep -qiE "RED.*GREEN.*REFACTOR|TDD cycle|TDD Implementation" "$CMD_FILE"; then
    fail "execute-task.md Step 2 still contains TDD keywords"
else
    pass "execute-task.md Step 2 has no TDD keywords"
fi

# AC: Step 2 has 3-case dispatch (Case A, Case B, Case C)
echo ""
echo "AC: 3-case dispatch present in Step 2..."
for case_label in "CASE A" "CASE B" "CASE C"; do
    if grep -qF "$case_label" "$AGENT_FILE"; then
        pass "task-executor.md Step 2 has $case_label"
    else
        fail "task-executor.md Step 2 missing $case_label"
    fi
done

for case_label in "CASE A" "CASE B" "CASE C"; do
    if grep -qF "$case_label" "$CMD_FILE"; then
        pass "execute-task.md Step 2 has $case_label"
    else
        fail "execute-task.md Step 2 missing $case_label"
    fi
done

# AC: Case B references default template
echo ""
echo "AC: Case B references default template..."
if grep -qF "plugins/forge/skills/breakdown-tasks/templates/task.md" "$AGENT_FILE"; then
    pass "task-executor.md Case B references default template"
else
    fail "task-executor.md Case B does NOT reference default template"
fi

if grep -qF "plugins/forge/skills/breakdown-tasks/templates/task.md" "$CMD_FILE"; then
    pass "execute-task.md Case B references default template"
else
    fail "execute-task.md Case B does NOT reference default template"
fi

# AC: NO_TEST removed from Inputs section
echo ""
echo "AC: NO_TEST removed..."
if grep -qiE "NO_TEST|noTest" "$AGENT_FILE"; then
    fail "task-executor.md still has NO_TEST/noTest references"
else
    pass "task-executor.md has no NO_TEST/noTest references"
fi

if grep -qiE "NO_TEST|noTest" "$CMD_FILE"; then
    fail "execute-task.md still has NO_TEST/noTest references"
else
    pass "execute-task.md has no NO_TEST/noTest references"
fi

# AC: Steps renumbered (5 steps: 0-4)
echo ""
echo "AC: Steps renumbered to 5 steps (0-4)..."
if grep -qE "Step [45]/5" "$AGENT_FILE"; then
    fail "task-executor.md still has old Step 4/5 or 5/5 numbering"
else
    pass "task-executor.md has no old Step 4/5 or 5/5 numbering"
fi

if grep -qE "Step [45]/5" "$CMD_FILE"; then
    fail "execute-task.md still has old Step 4/5 or 5/5 numbering"
else
    pass "execute-task.md has no old Step 4/5 or 5/5 numbering"
fi

# AC: Step 2 heading says "Execute Workflow" not "TDD Implementation"
echo ""
echo "AC: Step 2 titled 'Execute Workflow'..."
if grep -qE "Step 2.*Execute Workflow" "$AGENT_FILE"; then
    pass "task-executor.md Step 2 is titled Execute Workflow"
else
    fail "task-executor.md Step 2 is NOT titled Execute Workflow"
fi

# AC: Step 3 is Record (old Step 4), Step 4 is Commit (old Step 5)
echo ""
echo "AC: New Step 3 = Record, Step 4 = Commit..."
if grep -qE "Step 3.*Record" "$AGENT_FILE"; then
    pass "task-executor.md Step 3 is Record"
else
    fail "task-executor.md Step 3 is NOT Record"
fi

if grep -qE "Step 4.*Commit" "$AGENT_FILE"; then
    pass "task-executor.md Step 4 is Commit"
else
    fail "task-executor.md Step 4 is NOT Commit"
fi

# AC: Quality Gate section removed from task-executor.md
echo ""
echo "AC: No standalone Quality Gate step..."
if grep -qE "Step 3.*Quality Gate|Step 3.*Full Verification" "$AGENT_FILE"; then
    fail "task-executor.md still has Quality Gate as Step 3"
else
    pass "task-executor.md has no standalone Quality Gate step"
fi

# Summary
echo ""
echo "=== Results: $PASS passed, $FAIL failed (out of $TOTAL) ==="

if [ "$FAIL" -gt 0 ]; then
    exit 1
fi
