---
id: "5"
title: "gen-test-scripts 重写 + 内置模板迁移"
priority: "P0"
estimated_time: "4h"
dependencies: ["2", "4"]
scope: "backend"
breaking: false
type: "feature"
mainSession: false
---

# 5: gen-test-scripts 重写 + 内置模板迁移

## Description

重写 gen-test-scripts skill：基于 Contract 规范 + 配置驱动的模板 + 代码侦察（Fact Table），生成可执行测试代码。将语义描述符转换为精确正则（基于 Fact Table）。同时将现有 6 个 language profile 迁移为可覆盖的内置模板。自动生成 Journey 烟测试。

来源：proposal Pipeline 第 3 步、Scope "gen-test-scripts skill"和"内置模板迁移"。

## Reference Files
- `docs/proposals/contract-journey-test-model/proposal.md` — Source proposal
- `plugins/forge/skills/gen-test-scripts/` — 现有 gen-test-scripts 实现
- `plugins/forge/skills/gen-test-scripts/references/` — 现有 6 个 language profile（go/javascript/python/java/rust/mobile）
- `forge-cli/pkg/e2e/` — Fact Table 机制

## Acceptance Criteria

- [ ] gen-test-scripts 将语义描述符转换为精确正则（基于 Fact Table），生成可编译的测试代码（`go test -c` / `pytest --collect-only` / `tsc --noEmit` 通过）
- [ ] 生成的测试带 `@feature` 标签（Go 使用 `//go:build feature`），直接放入 `tests/<journey>/` 目录（无 staging 中间目录）
- [ ] 至少包含 1 个 Journey 烟测试，烟测试端到端运行 Journey 的 happy path
- [ ] 烟测试输出与 Contract 中 "success" Outcome 声明的 Output/State 完全匹配
- [ ] 零配置时 gen-test-scripts 输出与现有 profile 输出 diff 为空；config 声明自定义模板路径时使用自定义模板
- [ ] 现有 6 个 language profile 的 generate.md/run.md 作为内置默认模板工作

## Hard Rules

- 测试数据安全：生成的测试代码中不硬编码真实 secret/token；敏感字段使用占位符（如 `token: <from-env>`）
- 分批生成：单次生成一个 Journey（含 happy path + 边缘场景）
- 标签以语言框架原生方式嵌入（Go `//go:build`、Python `@pytest.mark`、JS `describe("@feature")`）

## Implementation Notes

- 语义描述符 → regex 转换管道示例：`Output: "success confirmation containing feature-slug"` → Fact Table 查询 → `Feature\s+([\w-]+)\s+created successfully`
- verify 自身准确性通过 bootstrap 策略保证：Phase 1 结束时用 126+ 已知正确输出生成 Fact Table 快照
- 模板迁移是机械性工作：现有 generate.md/run.md → 可覆盖的默认模板文件
