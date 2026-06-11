---
id: "1"
title: "撰写 Journey-Contract 测试模型 Convention 文件"
priority: "P0"
estimated_time: "1-2h"
dependencies: []
type: "doc"
mainSession: true
---

# 1: 撰写 Journey-Contract 测试模型 Convention 文件

## Description

将 Journey-Contract 测试模型的核心定义和旧→新迁移指南整合为 `docs/conventions/testing-journey-contract.md`，作为 Forge 的默认项目级知识。Agent 加载该文件后应能理解模型概念并自主执行旧项目迁移。

源材料：
- 模型定义来自 `docs/features/contract-journey-test-model/design/model-and-directory-spec.md`
- 6 维规则来自 `plugins/forge/skills/gen-contracts/rules/dimension-rules.md`
- Fact Table 流程来自 `plugins/forge/skills/gen-contracts/rules/code-reconnaissance.md`
- 验证规则来自 `plugins/forge/skills/gen-contracts/rules/validation.md`
- TUI 异步语义来自 `plugins/forge/skills/gen-contracts/rules/tui-async.md`
- 旧模型信息来自 `docs/proposals/unified-workflow-test-model/proposal.md` 和 `plugins/forge/skills/gen-test-cases/SKILL.md`

## Reference Files
- `docs/proposals/journey-contract-default-knowledge/proposal.md` — Source proposal
- `docs/features/contract-journey-test-model/design/model-and-directory-spec.md` — 模型定义权威参考
- `plugins/forge/skills/gen-contracts/rules/dimension-rules.md` — 6 维声明规则
- `plugins/forge/skills/gen-contracts/rules/code-reconnaissance.md` — Fact Table 构建流程
- `plugins/forge/skills/gen-contracts/rules/validation.md` — Contract 验证规则
- `plugins/forge/skills/gen-contracts/rules/tui-async.md` — TUI 异步语义
- `docs/proposals/unified-workflow-test-model/proposal.md` — 旧模型参考
- `docs/conventions/testing-conventions.md` — 现有 convention 格式参考（避免重复）

## Affected Files

### Create
| File | Description |
|------|-------------|
| `docs/conventions/testing-journey-contract.md` | Journey-Contract 模型定义 + 迁移指南 |

### Modify
| File | Changes |
|------|---------|
| _(无)_ | — |

### Delete
| File | Reason |
|------|--------|
| _(无)_ | — |

## Acceptance Criteria

- [ ] `docs/conventions/testing-journey-contract.md` 存在
- [ ] frontmatter 包含 `title` 和 `domains: [testing, journey, contract, e2e, migration]`
- [ ] **模型定义部分**覆盖：
  - Journey、Step、Contract、Outcome 核心概念及属性
  - 6 维声明规则（4 必选 + 2 可选）
  - 语义描述符规则（自然语言，不含 regex）
  - Preconditions 互斥性规则
  - Contract 文件格式（Outcome block 模板）
  - State 验证级别（full/partial/deferred）
  - 目录约定（`tests/<journey>/_contracts/` 结构）
  - Tag-Based Promotion（`@feature` → `@regression`）
  - TUI 异步 Cmd await 语义
- [ ] **迁移指南部分**覆盖：
  - 旧→新概念映射表（TC format → Contract 6 维、interface-type → Journey-driven、staging/graduation → Tag-Based Promotion）
  - 旧目录结构→新目录结构重组步骤
  - TC 单步格式→Contract 多 Outcome 格式转换规则
  - 迁移检查清单
- [ ] 不与 `testing-conventions.md`（Convention 文件格式规范）和 `testing-isolation.md`（测试隔离规则）内容重复
- [ ] 文件行数 < 600 行
- [ ] 内容精炼，去除源文件中的配置 schema、CLI 实现细节、pipeline 流程描述等不属于模型核心的内容

## Hard Rules

- 语义描述符规则中必须明确禁止 regex 语法
- 迁移指南必须基于实际的旧模型差异（`unified-workflow-test-model` → `contract-journey-test-model`），不凭空设计
- 不修改任何源文件（feature 目录、skills/rules/ 保持不动）

## Implementation Notes

- 参考 `testing-conventions.md` 的 frontmatter 格式和 domains 用法
- 源文件 `model-and-directory-spec.md` 有 455 行，需要精炼提取核心内容（去除 Config Schema 4.x 节、Pipeline Integration 7.x 节、Framework-to-Tag Mapping 4.4 节等不属于模型核心的部分）
- 旧模型的关键特征：按接口类型分类（CLI/API/TUI/UI/Mobile）、单步 TC 格式、staging+graduation 生命周期、6 个硬编码语言 profile
- 新模型的关键特征：按用户工作流组织（Journey）、6 维 Contract + 语义描述符、Tag-Based Promotion、Convention 驱动
- 此任务为 mainSession（在主会话中执行），因为需要大量源文件交叉参考
