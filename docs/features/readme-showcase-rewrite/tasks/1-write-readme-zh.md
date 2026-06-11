---
id: "1"
title: "重写 README.md（中文版）"
priority: "P0"
estimated_time: "1.5h"
dependencies: []
type: "doc"
mainSession: false
---

# 1: 重写 README.md（中文版）

## Description
从零重写 README.md，采用 7 节精简结构（叙述体开头 → 一句话定位 → 竞品对比表 → 核心特性 → 5 分钟体验 → 安装 → 贡献 + License），将 README 从"命令手册"重新定位为"项目名片"。当前 README 约 474 行，以 CLI 命令表和任务类型枚举为主，首屏缺少情感共鸣和差异化定位。

## Reference Files
- `docs/proposals/readme-showcase-rewrite/proposal.md` — Problem, Proposed Solution, Competitor Feature Matrix, Success Criteria, Key Risks
- `README.md` — 现有安装步骤需复用

## Affected Files

### Create
| File | Description |
|------|-------------|
| _(无新建文件)_ | |

### Modify
| File | Changes |
|------|---------|
| `README.md` | 全文重写为 7 节精简结构 |

### Delete
| File | Reason |
|------|--------|
| _(无删除文件)_ | |

## Hard Rules
- 竞品功能矩阵中每个 ✓/✗ 标记必须可溯源：撰写时通过实时搜索竞品（Superpowers / Spec Kit / OpenSpec）的最新 README、文档和仓库，逐项验证功能是否支持，确保信息准确。对有争议的项在 Implementation Notes 中注明判断依据。

## Acceptance Criteria
- [ ] 首屏（前 30 行）包含核心定位 + 痛点共鸣叙述（场景化叙事，非空泛口号）
- [ ] 总长度 ≤ 200 行（当前 ~474 行）
- [ ] 包含功能矩阵对比表，覆盖 Superpowers / Spec Kit / OpenSpec 共 3 个竞品，≥ 6 个对比维度
- [ ] 4 大核心特性（质量门控 / 上下文持久化 / 知识沉淀 / Agent 自动编排）各有独立小节
- [ ] 无 CLI 命令速查表（命令信息由 `forge --help` 承担）
- [ ] 保留现有安装步骤的技术准确性（install 命令、前置依赖不变）

## Implementation Notes
- 痛点场景参考：方向漂移、质量失控、上下文丢失、知识不沉淀
- 竞品功能矩阵数据来自 proposal 的 Competitor Feature Matrix 章节，使用 ✓/✗ 标记
- 叙述体开头保持技术准确性，用具体场景而非空泛口号
- 5 分钟体验路径需极简，让用户快速上手
