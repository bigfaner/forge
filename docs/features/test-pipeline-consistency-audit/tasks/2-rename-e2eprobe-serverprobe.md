---
id: "2"
title: "重命名 e2eprobe 包为 serverprobe"
priority: "P0"
estimated_time: "1h"
dependencies: [1]
surface-key: "."
surface-type: "cli"
breaking: true
type: "coding.refactor"
mainSession: false
---

# 2: 重命名 e2eprobe 包为 serverprobe

## Description
将 `pkg/e2eprobe` 包重命名为 `pkg/serverprobe`，更新所有 import 路径，并将探针配置路径从 `tests/e2e/config.yaml` 更新为 `tests/config.yaml`。此重命名消除了与特定测试类型绑定的包名，使探针功能适用于所有 server 类 surface。

## Reference Files
- `proposal.md#Layer-1-Go-代码层术语路径统一` — 第 3 项定义了 e2eprobe → serverprobe 重命名规则
- `proposal.md#Risks` — Go 包重命名导致外部引用断裂的风险及缓解措施

## Acceptance Criteria
- [ ] `pkg/e2eprobe/` 目录已重命名为 `pkg/serverprobe/`
- [ ] 所有 `import` 路径中的 `e2eprobe` 已更新为 `serverprobe`
- [ ] 探针配置路径使用 `GetTestConfigPath()` 而非硬编码
- [ ] `grep -rn "e2eprobe" forge-cli/ --include="*.go"` 返回 0 结果
- [ ] `go build ./...` 通过

## Implementation Notes
- 全面 grep 检查所有 import 路径，确保无遗漏
- 包名重命名后 CI 编译验证

### Integration Test Impact
- Affected test suite(s): `forge-cli/pkg/e2eprobe/`（重命名为 `forge-cli/pkg/serverprobe/`）
- Expected fixture changes: 无
- Risk level: medium
