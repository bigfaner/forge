---
created: "2026-06-11"
author: "fanhuifeng"
status: Draft
intent: "doc"
---

# Proposal: README 重写——从命令手册到项目名片

## Problem

当前 README.md 是一份内部命令速查手册（60%+ 篇幅是 CLI 命令表和任务类型枚举），无法在 5 秒内回答用户最核心的问题："为什么选 Forge 而不是别的？"。缺少竞品对比、没有痛点共鸣、缺少吸引贡献者的社区引导。

### Evidence

- 竞品 Superpowers（~224K stars）的 README 采用叙述体讲故事，情感感染力强
- OpenSpec（~54K stars）有明确的竞品对比表（vs Spec Kit / vs Kiro / vs nothing）
- GitHub Spec Kit（~111K stars）靠官方品牌背书 + 详尽的分步教程吸引用户
- 当前 README 把最有价值的信息（4 大核心优势）埋在了中间，被命令表淹没

### Urgency

README 是开源项目的"门面"。在 GitHub 搜索和社交媒体分享场景下，用户只会给首屏 5 秒。当前 README 首屏是痛点-解法对照表，缺少情感共鸣和差异化定位，直接影响用户转化率和社区吸引力。

## Proposed Solution

从零重写 README.md（中文）和 README.en.md（英文），采用 7 节精简结构：

1. **开头叙述** — 用场景化叙事引发 AI 编程用户痛点共鸣
2. **一句话定位** — 概括 Forge 的独特价值
3. **竞品对比表** — 功能矩阵表对标 Superpowers / Spec Kit / OpenSpec
4. **核心特性** — 4 大卖点各用一个视觉化小节展示
5. **5 分钟体验** — 极简上手路径
6. **安装** — 保留现有安装步骤
7. **贡献 + License**

### Innovation Highlights

- **叙述体开头**：借鉴 Superpowers 的讲故事手法，但不做空泛的愿景描述，而是用具体的痛点场景（方向漂移、质量失控、上下文丢失、知识不沉淀）引发共鸣
- **竞品功能矩阵表**：借鉴 OpenSpec 的对比手法，用 ✓/✗ 矩阵让 Forge 的差异化一目了然
- **去命令表化**：砍掉所有 CLI 命令罗列，把 README 定位为"项目名片"而非"速查手册"

## Requirements Analysis

### Key Scenarios

- **GitHub 浏览者**：在搜索结果或推荐流中看到 Forge，5 秒内判断"这东西值得试"
- **技术选型决策者**：对比 Superpowers / Spec Kit / OpenSpec / Forge，需要明确的差异点
- **潜在贡献者**：想了解项目架构和贡献方式
- **已有用户**：快速找到安装命令或更新信息

### Non-Functional Requirements

- 首屏（前 30 行）必须包含核心定位和吸引力钩子
- 总长度控制在 200 行以内（当前 README 约 475 行）
- 中英文内容对应，非直译，各自符合母语表达习惯

### Constraints & Dependencies

- 需要竞品信息准确（基于调研数据）
- 中文版为 README.md（主文件），英文版为 README.en.md
- 保留现有安装步骤的技术准确性

## Alternatives & Industry Benchmarking

### Industry Solutions

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 无法吸引新用户，与竞品差距拉大 | Rejected: README 是项目第一印象，不改进等于放弃增长 |
| 渐进改造 | — | 保留已有内容骨架 | 新旧内容割裂、语气不一致 | Rejected: 当前结构根深蒂固，修补式改造效果有限 |
| **全新重写** | — | 语气统一、结构清晰、定位精准 | 需投入撰写时间 | **Selected: 从零重写，确保首屏吸引力和整体一致性** |

### Competitor Feature Matrix

| 维度 | Forge | Superpowers | Spec Kit | OpenSpec |
|------|-------|-------------|----------|----------|
| 结构化流程 | ✓ | ✓ | ✓ | ✓ |
| 质量门控 (Quality Gate) | ✓ | ✗ | ✗ | ✗ |
| 上下文持久化 | ✓ | ✗ | ✗ | ✗ |
| Agent 自动编排 | ✓ | ✓ | ✗ | ✗ |
| 知识沉淀 (/learn) | ✓ | ✗ | ✗ | ✗ |
| 跨会话连续性 | ✓ | ✗ | ✗ | ✗ |
| 多 Agent 支持 | ✗ | ✓ | ✗ | ✗ |
| 跨 IDE/Agent 平台 | ✗ | ✓ | ✓ | ✓ |

## Feasibility Assessment

### Technical Feasibility

纯文档工作，无技术依赖。Markdown 编写即可。

### Resource & Timeline

单人完成，预计 2-3 小时（中文版 1.5h + 英文版 1h + 审校调整 0.5h）。

### Dependency Readiness

竞品调研已完成，数据可用。安装步骤直接复用现有内容。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "README 需要包含完整命令参考" | XY Detection | Overturned: 用户想要的是"为什么用"，不是"怎么用每个命令"。"怎么用"应该由 `forge --help` 和 wiki 承担 |
| "竞品对比会显得不自信" | Assumption Flip | Overturned: OpenSpec 和 GitHub Spec Kit 都有明确的竞品对比，这体现的是自信和透明 |
| "中文项目应该用英文 README" | 5 Whys | Refined: 用户选择以中文为主，因为目标社区以中文开发者为主 |

## Scope

### In Scope

- 全新 README.md（中文版，叙述体 + 7 节结构）
- 全新 README.en.md（英文版，对应翻译）
- 竞品功能矩阵对比表
- 4 大核心特性的视觉化展示
- 极简 5 分钟体验路径
- 贡献者引导

### Out of Scope

- CLI 命令速查表（移除，由 `forge --help` 承担）
- 架构详细文档（已有 docs/ARCHITECTURE.md）
- Skills 详细列表（已有多处文档）
- 网站 / landing page 设计
- 文档索引表

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 竞品信息不准确（功能矩阵的 ✓/✗ 有误） | M | H | 基于实际代码调研，标注信息来源 |
| 叙述体开头过于煽情，开发者不买账 | M | M | 保持技术准确性，用具体场景而非空泛口号 |
| 砍掉命令表后老用户找不到信息 | L | L | 命令信息已在 `forge --help` 和 wiki 中覆盖 |
| 中英文版本内容不一致 | L | M | 英文版独立撰写非直译，各自符合母语习惯 |

## Success Criteria

- [ ] README.md 首屏（前 30 行）包含核心定位 + 痛点共鸣叙述
- [ ] README 总长度 ≤ 200 行（当前 ~475 行）
- [ ] 包含功能矩阵对比表，覆盖 ≥ 3 个竞品，≥ 6 个维度
- [ ] 4 大核心特性各有独立小节
- [ ] 中英文版本内容对应（非直译），各自符合母语表达习惯
- [ ] 无 CLI 命令速查表（命令信息由 `forge --help` 承担）
- [ ] 保留安装步骤的技术准确性

## Next Steps

- Proceed to `/quick-tasks` to generate tasks directly from this proposal (doc intent, no PRD/design needed)
