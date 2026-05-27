---
id: "6"
title: "更新 prompt/task 模板中的旧术语"
priority: "P1"
estimated_time: "1h"
dependencies: [4]
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 6: 更新 prompt/task 模板中的旧术语

## Description
将 `pkg/prompt/data/` 和 `pkg/task/data/` 中测试相关模板的旧术语替换为新模型术语："profile"/"active profile"/"profile resolution" → "Convention"/"surface"。

## Reference Files
- `proposal.md#Layer-1-Go-代码层术语路径统一` — 第 8 项定义了模板术语替换范围
- `proposal.md#Success-Criteria` — 验证条件：prompt/task data 中 grep "profile" 返回 0

## Acceptance Criteria
- [ ] `pkg/prompt/data/test-gen-scripts.md` 中旧术语已替换
- [ ] `pkg/prompt/data/test-run.md` 中旧术语已替换
- [ ] `pkg/task/data/test-gen-scripts.md` 中旧术语已替换
- [ ] `grep -rn "profile" forge-cli/pkg/prompt/data/ forge-cli/pkg/task/data/ --include="*.md"` 返回 0 结果

## Implementation Notes
- 仅替换 "profile"/"active profile"/"profile resolution" 相关术语，不改动模板结构
- 注意不要误改其他不相关的 "profile" 用法（如 CPU profile 等，但这些不太可能出现在测试模板中）
