---
status: "completed"
started: "2026-05-28 00:22"
completed: "2026-05-28 00:42"
time_spent: "~20m"
---

# Task Record: 1 EvalConfig 扁平化为 bool 并简化默认值

## Summary
EvalConfig flattened from ModeToggle to bool with old-format backward compatibility

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/forgeconfig/config.go
- forge-cli/pkg/forgeconfig/config_test.go

### Key Decisions
- EvalConfig.UnmarshalYAML handles old ModeToggle map format by extracting 'full' sub-key as bool value
- scanMappingNode extended to track scalar (bool) eval sub-fields alongside ModeToggle nodes
- Removed omitempty from Auto field to prevent yaml.v3 IsZero() from dropping explicit-false auto blocks
- formatValue YAML unmarshaler check scoped to non-struct types only, allowing EvalConfig struct summary

## Test Results
- **Tests Executed**: Yes
- **Passed**: 210
- **Failed**: 0
- **Coverage**: 84.6%

## Acceptance Criteria
- [x] EvalConfig 4 fields (Proposal, Prd, UiDesign, TechDesign) type is bool, not ModeToggle
- [x] AutoConfigDefaults() eval: proposal:true, prd:false, uiDesign:true, techDesign:false
- [x] forge config get auto.eval.proposal returns true or false (bool format)
- [x] forge config set auto.eval.prd true writes and persists to config.yaml
- [x] Old ModeToggle format (e.g. proposal: {quick: true, full: true}) compat: extracts 'full' sub-key

## Notes
Also removed yaml omitempty from Config.Auto field to fix yaml.v3 IsZero() dropping explicit-false eval values during writeConfig roundtrip
