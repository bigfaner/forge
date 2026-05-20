---
created: 2026-05-20
author: faner
status: Draft
---

# Proposal: 提取 Journey-Contract 测试模型为 Forge 默认知识

## Problem

已有项目使用旧测试模型（按接口类型分类、单步 TC 格式、staging+graduation 生命周期），当用户要求 agent 迁移到 Journey-Contract 模型时，agent 不知道迁移流程和映射规则，因为模型定义和迁移知识分散在 feature 目录和 skill 规则文件中，不在 conventions 层。

### Evidence

- `docs/features/contract-journey-test-model/design/model-and-directory-spec.md`（455 行权威参考）位于 feature 目录，agent 不会自动加载
- 6 维规则、代码侦察、验证规则、TUI 异步语义分散在 `plugins/forge/skills/gen-contracts/rules/*.md` 的 4 个文件中
- 没有任何 convention 文件包含旧模型→新模型的迁移映射规则
- 现有 `tests/e2e/features/` 下仍有 6 个 feature 子目录使用旧结构

### Urgency

v3.0.0 正在重构测试体系。多个项目需要迁移，每次迁移都需要人工解释旧→新映射，效率极低。

## Proposed Solution

创建 `docs/conventions/testing-journey-contract.md`，整合两部分知识：

1. **模型定义**：Journey/Step/Contract/Outcome 概念、6 维规则、语义描述符、Fact Table、验证规则、TUI 异步语义、目录约定
2. **迁移指南**：旧模型（接口类型分类 + TC 格式 + staging/graduation）→ 新模型（Journey 驱动 + Contract 6 维 + Tag-Based Promotion）的映射规则和重组步骤

### Innovation Highlights

知识固化 + 迁移路径。将已完成 feature 的设计文档从"历史记录"升级为"活跃知识"，同时补齐迁移指南让 agent 能自主执行迁移。

## Requirements Analysis

### Key Scenarios

- **迁移场景（核心）**：用户说"把项目的测试迁移到 Journey-Contract 模型"，agent 自动加载 convention 文件后知道映射规则和重组步骤
- Agent 执行 gen-journeys/gen-contracts 任务时自动加载模型定义
- 新会话的 agent 无需翻找 feature 目录即可理解 Journey-Contract 模型
- `/consolidate-specs` 可将该文件纳入漂移检测范围

### Non-Functional Requirements

- 文件大小合理（< 600 行），agent 单次加载不会占用过多上下文
- domains frontmatter 精确匹配，避免不相关任务加载测试知识

### Constraints & Dependencies

- 不修改原始 feature 目录和 skills/rules/ 文件（保留为历史记录）
- 内容必须与现有 convention 文件（testing-conventions.md、testing-isolation.md）无重复
- 迁移指南基于实际的旧→新模型差异，不凭空设计

## Alternatives & Industry Benchmarking

### Industry Solutions

测试框架通常通过"测试策略文档"或"测试章程"固化团队测试知识。Forge 的 convention 目录就是这种模式的实现。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | Agent 无法执行迁移，每次需人工解释 | Rejected: 核心知识不应该是隐性的 |
| /consolidate-specs 自动提取 | Forge 内置 skill | 自动化程度高 | 输出格式不可控，不含迁移指南 | Rejected: 用户需要迁移知识，不仅是模型定义 |
| **手动提取为 convention 文件** | 本方案 | 精确控制内容，含迁移指南，agent 自动加载 | 需要与源文件保持同步 | **Selected: 最直接满足需求** |

## Feasibility Assessment

### Technical Feasibility

纯文档工作，无代码改动。convention 文件格式已有成熟模板。迁移映射可从旧模型和新模型的实际差异推导。

### Resource & Timeline

1 个 doc 类型任务：撰写 convention 文件。

### Dependency Readiness

源文件全部就绪且已稳定（feature 状态为 completed）。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 仅搬运模型定义即满足需求 | XY Detection | Overturned: 真实需求是迁移知识，不是纯搬运。用户需要 agent 能自主执行旧→新迁移 |
| 不会与现有 convention 文件冲突 | Codebase Analysis | Confirmed: testing-conventions.md 描述文件格式规范，testing-isolation.md 描述隔离规则，本文件描述模型概念+迁移指南，三者互补 |
| 含迁移指南后文件仍然合理 | Stress Test | Refined: 模型定义 ~400 行 + 迁移指南 ~100 行 = ~500 行，可接受 |

## Scope

### In Scope

**模型定义部分：**
- Journey、Step、Contract、Outcome 核心概念
- 6 维声明规则（4 必选 + 2 可选）
- 语义描述符规则（自然语言，不含 regex）
- Fact Table 构建流程
- 验证规则和错误处理
- TUI 异步 Cmd await 语义、状态验证级别
- 目录约定（tests/ 结构、Tag-Based Promotion）
- Contract 文件格式

**迁移指南部分：**
- 旧模型→新模型的概念映射表
- 旧目录结构→新目录结构的重组步骤
- TC 格式→Contract 6 维格式的转换规则
- Staging+Graduation→Tag-Based Promotion 的迁移步骤
- 迁移检查清单

### Out of Scope

- CLI 实现细节（`forge test` 命令）
- Config schema 详细定义
- 框架特定的测试约定（已在 testing-go.md 等）
- gen-test-scripts 的代码生成规则
- 自动化迁移工具开发

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 源文件更新后 convention 文件漂移 | M | M | /consolidate-specs 的漂移检测会覆盖 convention 文件 |
| 迁移指南不够具体，agent 仍无法执行 | M | M | 基于实际旧→新模型差异推导，包含具体映射表和步骤 |
| 与 testing-conventions.md 内容边界模糊 | L | M | 明确分工：本文件 = 模型概念+规则+迁移指南，testing-conventions.md = Convention 文件格式规范 |
| 文件过长导致 agent 加载成本高 | L | L | 精炼提取，目标 < 600 行 |

## Success Criteria

- [ ] `docs/conventions/testing-journey-contract.md` 存在且包含完整的模型定义和迁移指南
- [ ] frontmatter 的 `domains` 字段包含 `[testing, journey, contract, e2e, migration]`
- [ ] 模型定义覆盖：Journey/Step/Contract/Outcome、6 维规则、语义描述符、Fact Table、验证规则、TUI 异步、目录约定
- [ ] 迁移指南包含：概念映射表、目录重组步骤、TC→Contract 转换规则、生命周期迁移步骤、检查清单
- [ ] 不与 `testing-conventions.md` 和 `testing-isolation.md` 内容重复
- [ ] Agent 加载该文件后能回答"如何从旧测试模型迁移到 Journey-Contract 模型"

## Next Steps

- Proceed to `/quick-tasks` to generate and execute tasks
