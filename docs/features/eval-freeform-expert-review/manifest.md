---
feature: "eval-freeform-expert-review"
created: "2026-05-23"
status: completed
mode: quick
---

# Feature (Quick): eval-freeform-expert-review

<!-- Status flow: tasks -> in-progress -> completed -->

## Overview

为 eval-proposal 增加 `--freeform-expert` 参数，启用后在 rubric 循环之前插入 Phase 0 自由专家评审阶段。动态生成领域专家进行纯叙事评审，提取 key findings 并注入 rubric scorer，覆盖 rubric 维度之外的盲点。

## Documents

| Document | Path |
|----------|------|
| Proposal | ../../proposals/eval-freeform-expert-review/proposal.md |
| Test Cases | testing/test-cases.md |

## Tasks

| ID | Title | Status | File |
|----|-------|--------|------|
| 1 | Expert Profile Template & Inference Prompt | pending | tasks/1-expert-profile-template.md |
| 2 | Freeform Review Protocol & Agent Prompt | pending | tasks/2-freeform-review-protocol.md |
| 3 | Extraction Prompt Template & Injection Mechanism | pending | tasks/3-extraction-injection.md |
| 4 | Expert Persistence, Reuse & Deprecation | pending | tasks/4-expert-persistence.md |
| 5 | eval Skill Phase 0 Integration | pending | tasks/5-eval-integration.md |
