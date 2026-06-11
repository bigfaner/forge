---
id: "3"
title: "移除 test.gen-and-run 废弃代码 + 更新测试文件"
priority: "P2"
estimated_time: "2h"
dependencies: ["1", "2"]
scope: "backend"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 3: 移除 test.gen-and-run 废弃代码

## Description

从生产代码、测试文件和活跃文档中彻底移除 `test.gen-and-run` / `T-quick-gen-and-run` 的所有引用。Quick 模式已使用 staged pipeline（gen-journeys → gen-contracts → gen-scripts → run → verify-regression），gen-and-run 为僵尸代码。

## Reference Files
- `proposal.md#P2-—-gen-and-run-废弃代码移除` — 完整的按文件+行号移除清单
- `proposal.md#Key-Risks` — 编译失败风险和 Synthesize() file-not-found 风险
- `proposal.md#Success-Criteria` — grep 零结果验证标准

## Acceptance Criteria

- [ ] `grep -r "gen-and-run\|quick-gen-and-run\|T-quick-gen" forge-cli/ --exclude-dir=docs/proposals` 返回零结果
- [ ] `grep -r "gen-and-run\|quick-gen-and-run\|T-quick-gen" plugins/forge/ --exclude-dir=docs/proposals` 返回零结果
- [ ] validate_index.go 对引用 `test.gen-and-run` 的旧 index.json 返回迁移指引错误信息
- [ ] Synthesize() ReadFile 失败时对 `gen-and-run` 文件名输出迁移指引
- [ ] 所有现有测试通过

## Hard Rules

- **移除顺序必须严格遵循**（消费者先于常量定义）：infer.go → prompt.go → validate_index.go → build.go → types.go
- 每步后执行 `go build ./...` 验证增量编译通过
- `prompt/data/test-gen-and-run.md` 和 `task/data/test-gen-and-run.md` 必须删除

## Implementation Notes

**生产代码移除清单（按顺序）**：

1. `infer.go:32-33` — 移除 gen-and-run 推断分支
2. `prompt.go:293-305` — 移除 genScriptBases 中 `T-quick-gen-and-run` 条目
3. `validate_index.go:224-226` — 移除前缀检查，替换为迁移错误提示
4. `build.go:484,492-494` — 清理 gen-and-run 注释（findFirstTestTaskIdx 的 P1 修改由 Task 2 完成）
5. `types.go:55,88,114,134` — 移除 TypeTestGenAndRun 常量及所有注册表条目

**测试文件更新（13 个文件）**：按 `proposal.md#P2` 中列出的文件和行号逐一清理。

**活跃文档**：OVERVIEW.md、task-lifecycle（WORKFLOW.md）中的 gen-and-run 引用。

**validate_index.go 迁移错误**：在 `ValidTypes` 检查之前添加 `test.gen-and-run` 专用检查，返回 "test.gen-and-run is deprecated; use staged test pipeline types"。

**Synthesize() 迁移指引**：在 ReadFile 失败分支中，检查文件名包含 `gen-and-run` 时输出迁移指引。
