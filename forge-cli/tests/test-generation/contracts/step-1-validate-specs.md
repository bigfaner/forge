# Contract: test-generation / Step 1: Validate Specs

## Outcome "detects-e1-error-rules"
- Preconditions: "validate-specs.mjs exists in gen-test-scripts skill templates"
- Input: "Spec files with waitForTimeout and setTimeout (E1 violations)"
- Output: "validate-specs exits non-zero with JSON errors array containing rule E1"
- State: "no state changes"
- Side-effect: none

## Outcome "detects-e3-error-rules"
- Preconditions: "validate-specs.mjs exists in gen-test-scripts skill templates"
- Input: "Spec file without Traceability comment (E3 violation)"
- Output: "validate-specs exits non-zero with JSON errors array containing rule E3"
- State: "no state changes"
- Side-effect: none

## Outcome "detects-e4-error-rules"
- Preconditions: "validate-specs.mjs exists in gen-test-scripts skill templates"
- Input: "Spec file with DOM parent traversal locator('..') (E4 violation)"
- Output: "validate-specs exits non-zero with JSON errors array containing rule E4"
- State: "no state changes"
- Side-effect: none

## Outcome "detects-w1-w4-warning-rules"
- Preconditions: "validate-specs.mjs exists in gen-test-scripts skill templates"
- Input: "Spec files with serial suite > 15 tests (W1), no afterAll (W2), beforeEach login (W3), CSS class selector (W4)"
- Output: "validate-specs exits 0 with JSON warnings array containing rules W1-W4"
- State: "no state changes"
- Side-effect: none

## Outcome "structured-output-shape"
- Preconditions: "validate-specs.mjs exists in gen-test-scripts skill templates"
- Input: "Spec file with E1 violation"
- Output: "JSON output has 'errors' array and 'warnings' array, each entry has 'rule', 'file', 'message' fields"
- State: "no state changes"
- Side-effect: none

## Journey Invariants
- validate-specs.mjs is the authoritative spec validation script
- Error rules (E1-E4) cause non-zero exit; Warning rules (W1-W4) allow zero exit
- Output is structured JSON with errors and warnings arrays
