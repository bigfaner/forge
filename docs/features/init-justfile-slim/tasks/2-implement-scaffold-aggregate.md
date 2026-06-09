---
id: "2"
title: "实现 scaffold aggregate 聚合模式"
priority: "P0"
estimated_time: "2h"
complexity: "high"
dependencies: [1]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.feature"
mainSession: false
---

# 2: 实现 scaffold aggregate 聚合模式

## Description

实现 `forge justfile scaffold --aggregate` 模式：读取 `forge surfaces` 输出获取全部 surface key，生成跨 surface 的聚合 recipe（install / ci / clean），以及多服务编排时的 test-setup 聚合 recipe。

## Reference Files
- `docs/proposals/init-justfile-slim/proposal.md` — 新增：`forge justfile scaffold` CLI 命令 (聚合 recipe 生成、多服务编排模式)
- `forge-cli/internal/cmd/scaffold/` — Task 1 创建的 scaffold 包，本任务在此基础上扩展

## Acceptance Criteria
- [ ] `forge justfile scaffold --aggregate` 生成 install、ci、clean 三个聚合 recipe，聚合 recipe 无 `# user-customized` 标记
- [ ] ci recipe 聚合所有 surface 的 `<key>-lint` + `<key>-compile` + `<key>-unit-test`，不包含 surface-level test recipe
- [ ] 当 `forge surfaces` 返回多个 service-type surface（api/web/mobile）时，额外生成 test-setup 聚合 recipe，按依赖顺序编排启动和 teardown
- [ ] 纯 cli/tui 组合不生成 test-setup 聚合 recipe
- [ ] 聚合 recipe 正确引用带前缀的 surface recipe 名（如 `backend-compile`、`frontend-lint`）

## Implementation Notes
- 聚合生成器约 100 行，读取 `forge surfaces` 输出（调用现有 `forgeconfig` 包），按依赖顺序排列
- ci recipe 排除 surface-level test 是因为 surface test 需要运行时环境，不属于 CI 流水线
- 多服务编排的启动顺序为 api → web → mobile，teardown 逆序
