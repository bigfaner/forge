---
id: "3"
title: "更新 guide.md hook 测试速查表"
priority: "P0"
estimated_time: "1h"
dependencies: [1]
type: "doc"
mainSession: false
---

# 3: 更新 guide.md hook 测试速查表

## Description
在 guide.md hook 新增 Testing section，包含 Surface → Test Type 映射表、e2e 术语约束、测试文件位置规则。这是冷启动覆盖的关键——用户初始项目不存在任何 convention 文件时，agent 从 guide.md 即可正确回答测试相关问题。

guide.md 是每次会话的固定 token 成本，增量必须控制在 20 行以内。

## Reference Files
- `docs/proposals/surface-first-testing/proposal.md` — Proposed Solution (核心变更 #2), Non-Functional Requirements, Success Criteria
- `plugins/forge/hooks/guide.md`: 现有 guide.md 内容，需新增 Testing section (ref: Proposed Solution)
- `plugins/forge/skills/test-guide/references/test-type-model.md`: 映射表和 e2e 约束的权威来源 (ref: Proposed Solution)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/hooks/guide.md` | 新增 Testing section（速查表 + 映射表 + e2e 约束） |

## Acceptance Criteria
- [ ] guide.md 新增 Testing section，增量 <= 20 行
- [ ] 包含 Surface → Test Type 映射表（cli/api/web/tui/mobile 5 种 Surface）
- [ ] 包含 e2e 术语约束说明（"e2e" 仅用于 Web/Mobile）
- [ ] 包含测试文件位置规则（`tests/{surface}/` 或 `tests/e2e/`）
- [ ] 包含引导运行 `/test-guide` 获取完整策略的提示

## Hard Rules
- 增量 <= 20 行，超标则必须压缩
- 映射表内容必须与 `references/test-type-model.md` 保持一致

## Implementation Notes
- guide.md 是 hook 文件，修改后影响所有用户的所有会话，需精简
- 映射表至少包含 Surface 名称和 Test Type 名称两列，建议包含执行模型的一句话描述
- 参考提案中"速查表如果只传达映射结果而不传达映射理由"的评审建议——优先保证映射表的可推理性
