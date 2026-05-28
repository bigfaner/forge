---
id: "1"
title: "Create SC-Pre baseline measurement"
priority: "P0"
estimated_time: "1h"
complexity: "medium"
dependencies: []
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.feature"
mainSession: false
---

# 1: Create SC-Pre baseline measurement

## Description
Before any template modification, establish token/line count baselines and functional snapshot checklists for all In Scope files. These baselines are the reference points for SC6/SC8 (token reduction), SC1 (functional retention), and SC-FM-1 (frontmatter coverage) verification.

Create three categories of baseline artifacts: (1) per-file token and line counts for content slimming scope, (2) frontmatter structure records for all 41 templates, (3) per-template functional snapshot checklists capturing every semantic node.

## Reference Files
- forge-cli/pkg/prompt/templates/coding-*.md: Token/line baseline for 5 coding templates (source: proposal.md#In-Scope-内容精简)
- forge-cli/pkg/prompt/templates/gate.md: Token/line baseline (source: proposal.md#In-Scope-内容精简)
- forge-cli/pkg/prompt/templates/doc.md: Token/line baseline (source: proposal.md#In-Scope-内容精简)
- forge-cli/pkg/prompt/templates/test-*.md: Token/line baseline for test templates (source: proposal.md#In-Scope-内容精简)
- plugins/forge/agents/task-executor.md: Token/line baseline + Execution Protocol step count (source: proposal.md#In-Scope-内容精简)

## Acceptance Criteria
- [ ] `eval/baseline-token-counts.json` exists with per-file token and line counts for all prompt templates in content slimming scope + task-executor.md
- [ ] `eval/frontmatter-baseline.json` records current frontmatter structure (field count, variables list) for all 41 templates (21 prompt + 14 task + 6 record)
- [ ] Functional snapshot checklists created per template: JSON array `[{id, category, summary, sourceLine}]` using classification dictionary `{instruction, constraint, example, format, metadata}`
- [ ] All baseline files committed before any template modification begins

## Implementation Notes
- Token counting: use tiktoken or Claude-compatible tokenizer for accurate measurement
- Functional snapshots: node granularity is "semantic node" — one or more consecutive lines with independent functional purpose (e.g., one AC instruction, one CODING_PRINCIPLES principle, one Record Field declaration)
- Classification dictionary: `{instruction: 正面指令, constraint: 负面约束, example: 行为示范, format: 格式约定, metadata: 元数据声明}`
