---
created: 2026-05-16
author: "faner"
status: Draft
---

# Proposal: Auto-Behavior Configuration

## Problem

Forge pipeline 的自动化行为（e2e 测试、consolidate-specs、git push、代码清理）全部硬编码在流程中，用户无法按需开关。用户无法跳过不需要的 e2e 测试生成，无法在流程结束后自动推送代码，无法在提交前自动清理代码。

### Evidence

- `forge task index` 总是根据 `test-profiles` 生成 T-test/T-quick 任务，无跳过选项
- T-test-5（consolidate-specs）始终生成，即使项目不需要 spec 同步
- `/run-tasks` 完成后需要手动 `git push`
- `/simplify` 作为独立 skill 存在，但不在 pipeline 中自动调用
- todo.txt 第 79 行明确记录了此需求

### Urgency

随着不同类型项目（CLI、库、服务）使用 Forge，一刀切的 pipeline 越来越不灵活。CLI 工具可能不需要 e2e 测试但需要 spec 同步；库项目可能需要 e2e 测试但不需要自动 push。缺少配置控制导致用户要么接受不必要的任务，要么手动跳过。

## Proposed Solution

在 `.forge/config.yaml` 中添加 `auto` 配置块，其中 `e2eTest`、`consolidateSpecs`、`cleanCode` 按 quick/full 模式分别配置，`push` 为全局开关：

```yaml
auto:
  e2eTest:
    quick: true
    full: true
  consolidateSpecs:
    quick: true
    full: true
  cleanCode:
    quick: false
    full: false
  gitPush: false
```

### 配置项说明

| 配置项 | 类型 | 默认值 | 控制内容 |
|--------|------|--------|----------|
| `auto.e2eTest.quick` | bool | true | quick 模式是否生成 e2e 测试任务（T-quick-1~4） |
| `auto.e2eTest.full` | bool | true | full 模式是否生成 e2e 测试任务（T-test-1~4.5） |
| `auto.consolidateSpecs.quick` | bool | true | quick 模式是否生成 spec 同步任务（T-quick-specs-1） |
| `auto.consolidateSpecs.full` | bool | true | full 模式是否生成 spec 同步任务（T-specs-1） |
| `auto.cleanCode.quick` | bool | false | quick 模式是否追加 T-clean-code-1 任务 |
| `auto.cleanCode.full` | bool | false | full 模式是否追加 T-clean-code-1 任务 |
| `auto.gitPush` | bool | false | all-completed hook 通过后是否自动 git push |

### 创新亮点

- **模式分离**：e2eTest、consolidateSpecs、cleanCode 按 quick/full 模式独立配置，适应不同粒度的工作流
- **简洁结构**：所有配置收纳在 `auto` 块下，YAML 层级清晰
- **分层默认值**：已有行为默认 true（向后兼容），新增行为默认 false（opt-in）
- **命名修正**：将 T-test-5 重命名为 T-specs-1，反映其真实职责（文档-代码一致性维护，非测试）
- **任务类型正交化**：测试任务与维护任务完全独立控制，消除概念混淆

## Requirements Analysis

### Key Scenarios

- **纯 CLI 项目**：e2eTest 全 false，consolidateSpecs 全 true，跳过 e2e 但保留 spec 同步
- **库项目快速迭代**：quick 模式全部 true（含 cleanCode），full 模式仅 e2eTest
- **文档密集型项目**：e2eTest 全 false，consolidateSpecs 全 false，跳过所有自动任务
- **代码质量优先**：cleanCode.quick=true，quick 模式自动清理代码后跑测试

### Non-Functional Requirements

- 向后兼容：所有现有 `.forge/config.yaml` 无需修改即可继续工作
- 配置校验：JSON Schema 约束，非法值在 `forge task index` 时报错
- 明确反馈：跳过任务时在 index.json 中标记 skipped 并注明原因
- 平铺兼容：旧格式 `autoPush: true`（flat bool）自动升级为新格式

### Constraints & Dependencies

- 依赖 `forge task index` CLI（Go 二进制）读取配置并决定任务生成
- `auto.gitPush` 依赖 `/run-tasks` skill 的 post-completion 流程
- `auto.cleanCode` 依赖 `/simplify` skill

## Alternatives & Industry Benchmarking

### Industry Solutions

CI/CD 系统（GitHub Actions, GitLab CI）通过 yaml 配置控制 pipeline 步骤的启用/禁用是标准做法。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 无开发成本 | 无法按需定制 pipeline | Rejected |
| 环境变量控制 | CI 系统惯例 | 简单直接 | 不可版本化，团队无法共享 | Rejected |
| CLI 参数控制 | CLI 工具惯例 | 灵活 | 每次运行都需指定 | Rejected |
| 平铺 boolean 配置 | — | 简单 | 无法区分 quick/full 模式 | Rejected: 不够灵活 |
| **嵌套 auto 配置块** | CI/CD 系统 | 简洁、模式分离、声明式 | 需要修改 schema 和 CLI | **Selected** |

## Feasibility Assessment

### Technical Feasibility

- Config 读取：`forge task index` 已读取 `.forge/config.yaml`，扩展 `auto` 字段无障碍
- Schema 更新：`forge-config.schema.json` 需添加 `auto` 对象定义
- 任务生成逻辑：在 `forge task index` 中根据模式和配置跳过任务生成
- auto.gitPush：在 `/run-tasks` 的 all-completed hook 后新增 push 步骤
- auto.cleanCode：在 `forge task index` 中追加 T-clean-code-1 任务模板

### Resource & Timeline

7 个配置项（3×2 模式 + 1 全局）+ 1 重命名 + 2 新任务类型，scope 明确，适合 quick mode。

### Dependency Readiness

所有依赖已就绪：`/simplify` skill 已存在，`forge task index` CLI 已成熟，JSON Schema 已有验证机制。

## Scope

### In Scope

- 在 `forge-config.schema.json` 中添加 `auto` 配置块（含 7 个字段）
- 更新 `forge-config.example.yaml` 文档
- 修改 `forge task index` CLI：根据配置和模式决定是否生成测试/维护/clean 任务
- 重命名 T-test-5 → T-specs-1（标准模式）
- 拆分 T-quick-5 → T-quick-4 + T-quick-specs-1（quick 模式）
- 新增 T-clean-code-1 任务类型（调用 /simplify）
- 在 `/run-tasks` 中添加 auto.gitPush 步骤（all-completed hook 后）
- 更新所有引用 T-test-5 和 T-quick-5 的 skill 文档

### Out of Scope

- T-test-1~4 内部的更细粒度控制（如只跳过 T-test-3 run-e2e-tests）
- auto.gitPush 的 PR 自动创建（只做 git push）
- T-clean-code-1 的自定义参数（固定调用 /simplify）
- 配置项的运行时覆盖（如 CLI 参数 `--no-e2e`）
- 全局（非项目级）配置

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| `additionalProperties: false` 被移除后其他非法字段不再报错 | M | L | 添加具体字段定义而非放宽为 true |
| auto.gitPush 在无权限时失败 | M | M | 捕获 git push 错误，给出明确提示而非 crash |
| T-test-5 重命名导致现有 index.json 不兼容 | M | H | 在 CHANGELOG 中标注 breaking change，提供迁移说明 |
| auto.cleanCode 的 /simplify 修改了测试依赖的代码 | L | M | T-clean-code-1 在测试任务之前执行，测试验证 clean 后的代码 |
| 嵌套配置增加解析复杂度 | L | L | 支持 flat bool 快捷写法（`push: false` 等价于不需要模式分离的配置） |

## Success Criteria

- [ ] `.forge/config.yaml` 中 `auto.e2eTest.quick=false` 时，quick 模式不生成任何 T-quick 测试任务
- [ ] `auto.e2eTest.full=false` 时，full 模式不生成任何 T-test 测试任务
- [ ] `auto.consolidateSpecs.quick=false` 时，不生成 T-quick-specs-1
- [ ] `auto.consolidateSpecs.full=false` 时，不生成 T-specs-1
- [ ] `auto.cleanCode.quick=true` 时，quick 模式在业务任务和测试任务之间自动生成 T-clean-code-1
- [ ] `auto.cleanCode.full=true` 时，full 模式在业务任务和测试任务之间自动生成 T-clean-code-1
- [ ] `auto.gitPush=true` 时，all-completed hook 通过后自动执行 git push
- [ ] 无 `auto` 配置的现有项目行为与改动前完全一致
- [ ] 所有引用 T-test-5 / T-quick-5 的文档和模板已更新为新命名

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
