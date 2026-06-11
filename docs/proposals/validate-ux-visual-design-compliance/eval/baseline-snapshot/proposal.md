---
created: "2026-06-08"
author: fanhuifeng
status: Draft
intent: "enhancement"
---

# Proposal: validate-ux Visual Design Compliance

## Problem

validate-ux 的 eval 管线验证代码结构是否符合规格说明，但不验证渲染后的视觉输出是否符合 ui-design.md 的设计要求，导致视觉问题上线后才发现。

### Evidence

pm-work-tracker 项目的 milestone-map 功能通过完整 pipeline（prd → design → tasks → code → eval → test → validate-code → validate-ux）后，人工 QA 发现了 6 个 validate-ux 未捕获的问题，其中 4 个是纯视觉问题：

1. 列表页和详情页使用不同风格的筛选组件（样式不一致）
2. 刷新按钮位置与设计稿不符（布局错误）
3. 日历弹出框无法显示（交互组件失效）
4. validate-ux 报告"全部 7 个 UI 函数符合设计规格"，实际存在明显视觉偏差（视觉回归未检测）

详细记录见：`docs/lessons/gotcha-validation-ux-misses-visual-gaps.md`（pm-work-tracker 项目）

### Urgency

每次 web 功能迭代都可能产生同样的视觉盲区。validate-ux 报告"全部通过"给用户虚假安全感，实际上渲染结果可能与设计文档有显著偏差。随着 forge 承担更多 web 项目的自动化开发，这个问题的影响频率会持续增加。

## Proposed Solution

在现有 validate-ux eval 管线的 Phase 1（预处理）中新增「Design Compliance」采集维度：

1. **读取 ui-design.md**：提取每个页面/路由的 UI 要求（组件类型、布局规范、样式属性、交互元素）
2. **启动 web 服务，逐路由采集**：对每个路由截图 + 提取 accessibility tree（复用现有 agent-browser 基础设施）
3. **渲染结果 vs 设计文档逐条对比**：将 accessibility tree 中的实际渲染结构与 ui-design.md 的要求逐条映射，标记 pass/fail
4. **生成 design-compliance 报告**：作为 ux-snapshot.md 的新增 section，供 Phase 2 评分

Phase 2 rubric 新增「Design Spec Compliance」评分维度，覆盖组件存在性、布局结构、样式属性、交互元素可用性四个子维度。

### Innovation Highlights

不是像素级截图对比（如 Percy/Chromatic），而是**结构化语义对比**——将渲染结果的 accessibility tree 与设计文档的结构化描述进行匹配：

- 不依赖 golden image 维护
- 能检测组件缺失、错误变体、布局偏差
- 生成结构化的 pass/fail 报告，而非模糊的"相似度百分比"
- 借鉴 web accessibility 领域思路：accessibility tree 本身就是页面的结构化语义表示

## Requirements Analysis

### Key Scenarios

- **Happy path**：ui-design.md 描述了表格页的筛选器、数据表格、分页组件。validate-ux 启动服务，导航到该路由，截图并提取 accessibility tree，发现所有组件均存在且布局正确，design compliance 全部 pass。
- **组件缺失**：ui-design.md 要求某页面有日期选择器组件，但渲染结果中对应字段仅为普通文本输入框。design compliance 标记该条目为 fail，附截图和 accessibility tree 差异。
- **组件变体不一致**：ui-design.md 要求列表页和详情页使用相同的筛选器组件，但渲染结果显示两个页面的筛选器 DOM 结构不同。design compliance 标记为 fail。
- **ui-design.md 不存在**：validate-ux 跳过 design compliance 步骤，仅执行现有的 PRD flow 验证（向后兼容）。

### Non-Functional Requirements

- **兼容性**：对没有 ui-design.md 的项目（CLI、TUI、无 UI 设计的 web 项目），行为完全不变
- **性能**：design compliance 采集不应使 validate-ux 的执行时间增加超过 50%
- **可观测性**：design compliance 报告中每一条 pass/fail 判断必须附有截图路径和具体的 accessibility tree 节点引用

### Constraints & Dependencies

- 依赖 ui-design.md 的存在和质量（由 /ui-design skill 生成）
- 依赖 sitemap.json 提供路由信息（web surface 项目）
- 依赖 agent-browser（或等效浏览器自动化）进行截图和 accessibility tree 提取
- 仅适用于 web surface 类型的项目

## Alternatives & Industry Benchmarking

### Industry Solutions

视觉回归测试在行业中有成熟方案：Percy、Chromatic、Playwright screenshot comparison 等，均为像素级截图对比，需要维护 baseline image。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | gotcha 证明现有 validate-ux 漏检视觉问题 | Rejected: 问题已实证存在 |
| 像素级截图对比 | Percy/Chromatic | 精确到像素，成熟工具 | 需维护 baseline image，任何细微变化触发 diff，误报率高 | Rejected: 维护成本过高，与 forge 单次验证场景不匹配 |
| 跨页面一致性审计 | 自研 | 无需设计文档即可发现不一致 | 不知道哪个是"正确"的，高误报率 | Rejected: 不解决"是否符合设计"的核心问题 |
| **Design-spec 结构化对比** | Accessibility tree 语义匹配 | 精准对标设计文档，结构化 pass/fail，无需 baseline | 依赖 ui-design.md 质量 | **Selected: 直接对标 gotcha 缺口，复用现有基础设施** |

## Feasibility Assessment

### Technical Feasibility

validate-ux 管线已有 web surface 支持（agent-browser 截图、sitemap.json 路由发现、accessibility tree 提取）。新增 design compliance 步骤是在已有流程上增加对比维度，不需要新的基础设施。ui-design.md 由 /ui-design skill 生成，已有结构化格式。

### Resource & Timeline

涉及 4-5 个文件的修改（pipeline 规则、rubric、expert、ux-snapshot 格式、可能需要更新 prompt template），属于中等规模改动，不需要新的外部依赖。

### Dependency Readiness

- agent-browser: 已在 validate-ux pipeline 中集成
- ui-design.md: 由现有 /ui-design skill 生成
- sitemap.json: web surface 项目已有
- accessibility tree 提取: agent-browser 已支持

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| validate-ux 已经验证了视觉合规性 | 5 Whys | Overturned: validate-ux 只验证代码结构（组件名、类名），不验证渲染结果。gotcha 中的 4 个视觉问题全部通过了 validate-ux |
| 需要像素级截图对比才能检测视觉问题 | Occam's Razor | Refined: accessibility tree 提供了结构化语义信息，足以检测组件缺失、错误变体、布局偏差。像素级对比对 forge 场景过重 |
| 应该创建新的 eval 类型 | XY Detection | Overturned: 用户真正需要的是"validate-ux 能检测视觉问题"，增强现有 validate-ux 更简单且复用基础设施 |

## Scope

### In Scope

1. validate-ux Phase 1 pipeline 新增「Design Compliance」采集步骤：读取 ui-design.md → 提取 UI 要求 → 逐路由对比渲染结果 → 生成 design-compliance 报告
2. validate-ux rubric 新增「Design Spec Compliance」评分维度（覆盖组件存在性、布局结构、样式属性、交互元素可用性）
3. 更新 validate-ux-pipeline.md 规则文件，增加 design compliance 采集流程定义
4. 更新 ux-snapshot.md 格式定义，增加 design-compliance section
5. 更新 ux-auditor expert，增加视觉设计合规评估专长

### Out of Scope

1. 交互行为测试（点击按钮验证 toast、表单提交反馈）——仅做渲染对比
2. API URL 模式校验——属于 validate-code 或 E2E 范围
3. 跨页面一致性审计（无设计文档参照的对比）——未来增强方向
4. 像素级视觉回归测试——超出当前范围
5. CLI/TUI surface 类型的视觉对比——本次聚焦 web surface
6. task 层面的 validation.ux 改动——仅增强 eval 层面的 validate-ux

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| ui-design.md 质量不足导致对比基准不可靠 | M | H | design compliance 报告中标注"无法对比"的条目，不产生 pass/fail 误判；在 report summary 中提示设计文档覆盖度 |
| accessibility tree 无法反映某些视觉属性（颜色渐变、阴影、动画） | M | M | 在 rubric 中明确 design compliance 的能力边界——验证结构/语义合规，不验证视觉细节 |
| agent-browser 截图在 CI 环境中不稳定 | L | M | 复用现有 validate-ux 的 dev→probe→test 编排，已有重试机制 |
| 新增步骤使 validate-ux 执行时间显著增加 | L | L | design compliance 仅在有 ui-design.md 时执行；无 ui-design.md 时完全跳过，零额外耗时 |

## Success Criteria

- [ ] SC-1: 对有 ui-design.md 的 web surface 项目，validate-ux Phase 1 为 sitemap.json 中 100% 的路由生成 design compliance 采集数据（截图 + accessibility tree）
- [ ] SC-2: ux-snapshot.md 中新增 design-compliance section，每条 ui-design.md UI 要求映射到 pass/fail 判断，附截图路径和 accessibility tree 节点引用
- [ ] SC-3: rubric 新增「Design Spec Compliance」评分维度，包含至少 4 个子维度（组件存在性、布局结构、样式属性、交互元素可用性），每个子维度有明确的评分标准和扣分规则
- [ ] SC-4: 无 ui-design.md 的项目执行 validate-ux 时，行为与当前完全一致（无 design compliance 步骤，无额外输出，无额外耗时）
- [ ] SC-5: design compliance 报告中的每一条 fail 判断必须包含：失败原因、期望值（来自 ui-design.md）、实际值（来自 accessibility tree）、截图路径

consistency_check_result:
  status: pass
  pairs_checked: 10
  conflicts_found: 0

## Next Steps

- Proceed to `/write-prd` to formalize requirements
