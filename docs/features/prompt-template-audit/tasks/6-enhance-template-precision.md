---
id: "6"
title: "P2: 增强模板指令精确度"
priority: "P2"
estimated_time: "45m"
dependencies: ["1", "2", "3", "4", "5"]
type: "doc"
mainSession: false
---

# 6: P2: 增强模板指令精确度

## Description
增强 3 个模板中模糊指令的精确度：(1) coding-feature.md Step 2 TDD 引导不足——添加从 Acceptance Criteria 提取需求的引导；(2) doc.md Step 2 Execute 过于模糊——添加按任务类型执行的引导；(3) 5 个 coding 模板的测试示例使用 Go 特定语法——改为语言无关描述。

## Reference Files
- `docs/proposals/prompt-template-audit/proposal.md` — Source proposal (Sections 2.1, 2.6, P2 #10/#11/#12)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/prompt/data/coding-feature.md` | Step 2 添加从 Acceptance Criteria 提取测试需求的引导 |
| `forge-cli/pkg/prompt/data/coding-enhancement.md` | Step 2 添加相同的 AC 提取引导；将 Go 测试示例改为语言无关描述 |
| `forge-cli/pkg/prompt/data/coding-cleanup.md` | 将 Go 测试示例改为语言无关描述 |
| `forge-cli/pkg/prompt/data/coding-fix.md` | 将 Go 测试示例改为语言无关描述 |
| `forge-cli/pkg/prompt/data/coding-refactor.md` | 将 Go 测试示例改为语言无关描述 |
| `forge-cli/pkg/prompt/data/doc.md` | Step 2 添加按任务类型（创建/修改/删除文档）执行的引导 |

## Acceptance Criteria
- [ ] coding-feature.md Step 2 包含"从 task 文件的 Acceptance Criteria 中提取每个测试需求"的引导
- [ ] coding-enhancement.md Step 2 包含相同引导
- [ ] 所有 5 个 coding 模板中 `go test -race -cover ./changed/package/...` 等 Go 特定示例替换为"运行受影响包/模块的测试命令"等语言无关描述
- [ ] doc.md Step 2 包含按文档任务类型执行的具体引导（识别类型 → 按类型执行）

## Implementation Notes
- Go 示例仅是说明性的——agent 实际使用 `just test` 命令，模板中不应硬编码特定语言的测试命令
- doc.md Step 2 可参考 doc-summary.md 的精确度标准
- 依赖前序任务完成（确保修改基于已修复的模板版本）
