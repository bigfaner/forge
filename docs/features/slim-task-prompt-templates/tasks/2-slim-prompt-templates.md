---
id: "2"
title: "Slim all prompt templates"
priority: "P0"
estimated_time: "2h"
complexity: "high"
dependencies: [1]
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 2: Slim all prompt templates

## Description
Apply 6 content slimming techniques to all prompt template files in `forge-cli/pkg/prompt/templates/` that are in the content slimming scope (coding-*, gate, doc, test-*, code-quality-simplify, validation-*). The techniques are: (1) delete HTML comments, (2) delete Step 2 explanatory descriptions, (3) convert role descriptions from natural language to imperative sentences, (4) compress AC verification blocks from ~12 lines to ~4 lines, (5) streamline Record Fields descriptions, (6) compress CODING_PRINCIPLES by removing examples and explanations while retaining 1 representative example per principle.

The "prompt is instruction, not documentation" principle governs all changes: delete everything that doesn't directly guide agent behavior.

## Reference Files
- forge-cli/pkg/prompt/templates/coding-feature.md: Role + AC block + CODING_PRINCIPLES + Record Fields compression (source: proposal.md#Key-Scenarios-1)
- forge-cli/pkg/prompt/templates/coding-enhancement.md: Same pattern as coding-feature (source: proposal.md#Key-Scenarios-1)
- forge-cli/pkg/prompt/templates/gate.md: Role + AC block compression (source: proposal.md#Key-Scenarios-2)
- forge-cli/pkg/prompt/templates/test-run.md: Step 2 description + role compression (source: proposal.md#Key-Scenarios-3)
- forge-cli/pkg/prompt/templates/code-quality-simplify.md: Role compression only (source: proposal.md#Key-Scenarios-4)

## Acceptance Criteria
- [ ] SC1: 100% functional node retention rate — every instruction/constraint/format node from functional snapshot checklist is present in modified templates
- [ ] SC3: CODING_PRINCIPLES core constraint instructions preserved in all 5 coding-* templates (1 instruction line + 1 boundary summary + 1 representative example per principle)
- [ ] SC4: Record Fields field names and value structures preserved in all applicable templates
- [ ] SC5: All Step 2 explanatory descriptions fully deleted — grep confirms zero residuals
- [ ] Role descriptions converted to imperative sentences (no natural language character descriptions remain)

## Hard Rules
- **Line format invariance**: TASK_FILE, TASK_ID, SURFACE_KEY lines must retain `KEY: {{.Value}}` format and position — prompt.go uses strings.Replace on these lines
- **CRITICAL blocks preserved**: Do not modify Spec Authority Enforcement, Hard Rules, or any `**CRITICAL**` marked sections
- **AC:REQUIRED vs AC:STRONGLY distinction preserved**: These have different compliance strength semantics for the agent

## Implementation Notes
- Use coding-feature.md as the reference implementation — other coding-* templates follow the same pattern
- gate/doc templates: AC verification blocks but no CODING_PRINCIPLES
- test-* templates: mainly Step 2 description removal + role slimming
- code-quality/validation templates: simplest — role description slimming only (~22 lines total across 3 files)
- Instruction classification for review: A (positive instructions → keep), B (negative constraints → keep), C (behavioral examples → keep 1 per principle)
- Target: ~1200-1500 tokens reduction for coding.* templates, ~800-2200 tokens total
