---
created: 2026-05-21
author: faner
status: Draft
---

# Proposal: Decouple run-tests from e2e/just Hardcoding

## Problem

`run-e2e-tests` skill 硬编码了 `just e2e-setup` / `just e2e-verify` / `just e2e-test` 三个命令，将 Forge 绑定到 `just` 任务运行器和 e2e 测试语义上。非 e2e 项目（Go unit tests、Vitest integration tests 等）无法使用此 skill。

### Evidence

- SKILL.md 中 3 处硬编码 `just e2e-*` 命令（第 87、99、109 行）
- 前置检查 `grep -q "^e2e-setup:" Justfile` 强制要求 justfile
- skill 名称 `run-e2e-tests` 在整个 pipeline 中传播 e2e 语义（gen-test-scripts 的 Related Skills 表、task chain 中的 T-test-3）
- `gen-test-scripts/templates/` 目录包含 27MB 死代码（validate-specs.mjs + node_modules），SKILL.md 零引用

### Urgency

v3.0.0 正在重构测试体系（test profile 系统、Convention-driven 模型）。保留 e2e/just 硬编码会与框架无关化方向矛盾，越晚改迁移成本越高。

## Proposed Solution

将 `run-e2e-tests` 重新定位为纯执行器 `run-tests`，执行命令从 `.forge/config.yaml` 的 `test.execution` 节点读取。skill 只负责三件事：执行配置中的命令 → 解析结果 → 生成报告。

### Innovation Highlights

这不是创新，是职责归位。Convention 管"怎么写测试"（framework、assertion、result format），Config 管"怎么跑测试"（执行命令模板）。两者已在各自的文件中，只是缺少 `test.execution` 这个桥接点。

## Requirements Analysis

### Key Scenarios

1. **Happy path**: 用户在 config.yaml 配置 `test.execution.run`，调用 `/run-tests`，skill 执行命令、解析结果、生成报告
2. **完整生命周期**: 用户配置 setup → run → teardown，skill 按顺序执行
3. **无配置**: config.yaml 中没有 test.execution → 报错并提示配置方法
4. **verify 失败**: 用户配置了 verify 步骤，检测到 `// VERIFY:` 标记 → 报错并提示回到 gen-test-scripts
5. **结果解析失败**: Convention Result Format 与实际输出不匹配 → 使用 text-verbose fallback

### Non-Functional Requirements

- 执行命令支持 `{slug}` 模板变量替换
- Convention Result Format 解析逻辑不受影响（已框架无关）
- 报告模板去除 e2e 专属文案，适用于任何测试类型

### Constraints & Dependencies

- 依赖 `.forge/config.yaml` 存在且包含 `test.execution` 配置
- 依赖 Convention 文件提供 Result Format（用于结果解析）
- 不依赖 justfile（但用户可以在 config 中使用 just 命令）

## Alternatives & Industry Benchmarking

### Industry Solutions

测试运行器的配置通常在项目配置文件中（package.json scripts、Justfile、Makefile）。本方案将这种模式提升为 Forge config 的标准字段，使 skill 不关心具体运行器。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零改动成本 | 绑定 just+e2e，非 TS/Playwright 项目无法使用 | Rejected: 与 v3.0 框架无关化方向矛盾 |
| Convention Execution 段落 | 自研 | Convention 已加载，无新文件 | 混淆"怎么写"和"怎么跑"的职责边界 | Rejected: 职责不清 |
| **Config test.execution** | 自研 | 职责清晰（Convention=写，Config=跑），支持任意执行器，集中管理 | 新增 config schema | **Selected: 最干净的职责分离** |
| CLI 参数传入 | CLI pattern | 零配置 | 无法复用，每次调用都要输入完整命令 | Rejected: 操作成本高 |

## Feasibility Assessment

### Technical Feasibility

完全可行。只需：(1) 读取 config.yaml 的 test.execution 字段、(2) 替换 {slug} 变量、(3) bash 执行。现有 Result Format 解析和报告生成逻辑不变。

### Resource & Timeline

改动量小：1 个 skill 重写 + 1 个 config schema 定义 + 若干引用更新。1-2 个 coding task 可完成。

### Dependency Readiness

`.forge/config.yaml` 已存在，forge CLI 已有 `forge config get` 命令。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "skill 需要管理测试生命周期（setup/verify/teardown）" | Occam's Razor | Refined: 生命周期由用户在 config 中定义，skill 只是按顺序执行。不内置任何生命周期逻辑 |
| "validate-specs.mjs 是 gen-test-scripts 的必要组件" | Evidence check (grep) | Overturned: SKILL.md 零引用，470 行 Playwright 专属代码 + 27MB node_modules 是死代码 |
| "justfile 是 Forge 的标准执行层" | Assumption Flip | Refined: justfile 是 Forge 推荐的执行层，但不是唯一选择。Config 中可以写 make、npm run、cargo test 等任何命令 |

## Scope

### In Scope

- 重命名 `run-e2e-tests` → `run-tests`，重构 SKILL.md 为纯执行器
- 定义 `.forge/config.yaml` 的 `test.execution` schema（run 必填，setup/verify/teardown 可选）
- 删除 `gen-test-scripts/templates/` 整个目录（validate-specs.mjs + test + fixtures + node_modules）
- 泛化报告模板（`e2e-report.md` → `test-report.md`，去除 e2e 专属文案）
- 更新上游引用（gen-test-scripts Related Skills 表、其他 skill 中的 run-e2e-tests 引用）

### Out of Scope

- `init-justfile` skill 的修改（justfile target 名称不变，用户在 config 中引用它们）
- `gen-test-scripts` skill 的核心逻辑修改（只更新 Related Skills 引用）
- Convention 文件格式的修改（Result Format 已框架无关，不动）
- 其他 skill 中 e2e 命名引用的全面清理（只处理直接关联的引用）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 已有项目缺少 test.execution 配置导致 skill 无法使用 | M | M | 报错信息明确：提示用户在 config.yaml 中添加 test.execution 配置 |
| 用户配置的执行命令语法错误 | L | L | 执行前不校验命令语法，依赖 bash 的自然错误输出 |
| teardown 在异常中断时未执行 | M | M | 使用 bash trap 机制确保 teardown 在任何退出情况下执行 |

## Success Criteria

- [ ] `run-tests` skill 不包含任何硬编码的 `just` 或 `e2e` 命令
- [ ] `test.execution.run` 配置缺失时，skill 报错并给出配置示例
- [ ] `{slug}` 变量在所有命令模板中正确替换为当前 feature slug
- [ ] `gen-test-scripts/templates/` 目录已被删除（减少 27MB）
- [ ] 报告模板中不含 "E2E" 字样，适用于任何测试类型
