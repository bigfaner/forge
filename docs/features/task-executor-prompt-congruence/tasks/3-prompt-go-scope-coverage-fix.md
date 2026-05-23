---
id: "3"
title: "prompt.go scope resolution 与 coverage 语言修复"
priority: "P1"
estimated_time: "2h"
dependencies: [2]
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 3: prompt.go scope resolution 与 coverage 语言修复

## Description

修复 `forge-cli/pkg/prompt/prompt.go` 中的 scope resolution 逻辑和 resolveCoverage() 函数，解决 3 处代码层面的问题（Issues 5-code, 7, 10）。

当前 prompt.go 存在：(1) scope 空值未做 project-type fallback；(2) resolveCoverage() 对 cleanup/refactor 类型注入矛盾指令；(3) resolveCoverage() 返回中文文本注入英文模板。

## Reference Files
- `docs/proposals/task-executor-prompt-congruence/proposal.md` — Source proposal
- `docs/conventions/forge-distribution.md` — Forge distribution model constraints

## Acceptance Criteria

- [ ] scope resolution 逻辑扩展：当 task scope 与 project-type 不匹配时（如 backend 项目的 task scope="frontend"），fallback 到无 scope 参数的默认命令（如 `just compile` 而非 `just compile frontend`）
- [ ] resolveCoverage() 对 coding-cleanup 和 coding-refactor 类型不注入 percentage 覆盖率指令（因为这两类任务的 "no new tests" 指令与 coverage 要求矛盾）
- [ ] resolveCoverage() 返回英文文本，模板内无语言混杂
- [ ] 现有测试通过：`forge-cli/pkg/prompt/prompt_test.go` 和 `forge-cli/tests/scope-resolution/scope_resolution_test.go`
- [ ] 新增测试覆盖 scope fallback 逻辑（backend project + frontend scope → default command）

## Hard Rules

- 不修改 markdown 模板文件（Task 2 负责）
- 不修改 forge-cli 核心逻辑，仅改 prompt.go 及其测试
- 遵循 `docs/conventions/forge-distribution.md` 中的路径解析约束
- 向后兼容：不影响已生成的 task 文件

## Implementation Notes

- Issue 7 (scope resolution): `renderTemplate` 已有 scope 空值处理，需扩展为完整 resolution —— 检查 `forge config get project-type`，若 project-type 为 backend 且 scope 为 frontend，则忽略 scope
- Issue 5-code (coverage): resolveCoverage() 当前对所有 coding.* 类型注入覆盖率指令，需要对 cleanup/refactor 特殊处理——跳过 percentage 策略注入，或改为 "maintain existing coverage" 文本
- Issue 10 (语言): resolveCoverage() 中的中文字符串需替换为英文。影响的是 prompt 模板内的注入文本
- 已有测试文件：`forge-cli/pkg/prompt/prompt_test.go`（resolveCoverage 单元测试）和 `forge-cli/tests/scope-resolution/scope_resolution_test.go`（scope resolution 集成测试）
- 风险：scope resolution 引入新 bug，但已有测试覆盖
