---
created: "2026-05-30"
author: "faner"
status: Draft
intent: "new-feature"
---

# Proposal: Sitemap Dedup CLI

## Problem

gen-sitemap 生成的 sitemap.json 存在四种重复节点问题：布局元素在每个 page.elements 中重复出现、跨页面元素重复、State 元素重复、以及全局 ID 冲突。尽管 SKILL.md 和 merge-validation.md 中已定义了七层去重规则，但这些规则完全依赖 LLM agent 在执行过程中手动计算集合交集、过滤、合并——本质上是让非确定性的 LLM 执行确定性的集合操作，导致规则在执行层面失效。

### Evidence

- `docs/lessons/pattern-sitemap-shared-layout.md` 记录了一次典型故障：生成 351 个元素中约 120 个是 sidebar 重复（10 个 sidebar 元素 × 12 个认证页面）。根因分析指出 agent 获取了架构信息但未使用，机械地逐页转录。
- 去重规则分散在 SKILL.md（Step 2、Step 4、Step 5）和 merge-validation.md（Step 5a-5d）中，agent 难以一贯地遵守所有步骤。
- 四类重复（布局重复、页面元素重复、State 重复、ID 冲突）全面出现，说明问题不是某个单一规则的遗漏，而是执行机制的系统性失效。

### Urgency

sitemap.json 是多个下游 skill 的输入源（write-prd、gen-test-scripts、breakdown-tasks、eval）。重复节点导致这些 skill 产生冗余或矛盾的输出，影响整个测试 pipeline 的可靠性。每次手动清理 sitemap 是持续的时间消耗。

## Proposed Solution

将去重逻辑从自然语言规则转为 Forge Go CLI 中的确定性命令：

1. **新增 `forge sitemap dedup` 命令**：读取 sitemap.json，执行五阶段去重算法（跨页共享检测 → 单页内去重 → State 元素去重 → ID 重分配 → 结构验证），输出修复后的 sitemap.json 和变更报告。
2. **新增 `forge sitemap validate` 命令**：只验证不修复，输出详细的错误报告。
3. **修改 SKILL.md Step 5**：将自然语言去重规则替换为 `forge sitemap dedup` 调用指令，agent 审查变更报告处理边缘情况。
4. **简化 merge-validation.md**：移除已由 Go 代码处理的去重规则，仅保留需要 agent-browser 的任务（stale route detection）。

### Innovation Highlights

这是将 AI agent 的"软规则"转化为确定性代码的典型案例。核心洞察是：集合操作（交集、去重、ID 分配）天然适合程序化执行，而非 LLM 推理。这个模式可以推广到其他 skill 中存在类似问题的环节。

## Requirements Analysis

### Key Scenarios

- **Happy path**：Agent 生成原始 sitemap.json → 调用 `forge sitemap dedup` → 脚本输出变更报告（提升 5 个共享元素、移除 3 个重复、重分配 12 个 ID）→ agent 审查报告，确认无需额外修改
- **Layout 识别遗漏**：Agent 未在 Step 2 识别到共享 layout → 所有元素都在 page 级别 → `forge sitemap dedup` 检测到多个 wrapped 页面中的相同元素 → 自动提升到 layout.elements
- **ID 冲突**：Agent 分配了重复的 E-NNN ID → 脚本重新编号，修复所有引用
- **边缘情况**：脚本处理后仍有个别语义重复（不同 role+name 但实际是同一元素）→ agent 审查报告时手动修复
- **增量更新**：已有 sitemap.json → 新增页面探索后 → 调用 dedup → 保留已有 ID，仅处理新增部分的重复

### Non-Functional Requirements

- **性能**：sitemap.json 通常 < 1000 个元素，Go 处理应在 < 100ms 内完成
- **向后兼容**：不改变 sitemap.json 的 schema（字段名、ID 格式、嵌套结构不变）
- **幂等性**：对已去重的 sitemap.json 多次运行 `dedup`，结果不变

### Constraints & Dependencies

- 依赖 Forge Go CLI 的编译和分发机制
- `forge sitemap dedup` 需要读取用户项目中的 `docs/sitemap/sitemap.json`
- 阈值采用严格模式：只有出现在所有 layout-wrapped 页面的元素才被提升为 layout.elements

## Alternatives & Industry Benchmarking

### Industry Solutions

sitemap 去重本质上是集合操作问题（交集、唯一性、ID 分配），业界标准做法是用确定性代码处理，而非自然语言指令。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 无开发成本 | 重复问题持续，下游 skill 受影响 | Rejected: 问题影响整个测试 pipeline |
| 加强 SKILL.md 规则 | 当前实践 | 无新代码 | 仍然依赖 LLM 执行确定性逻辑，根本问题未解决 | Rejected: 根因是执行机制而非规则不够 |
| 独立验证 Agent | User proposed | 分离关注点 | 两层 LLM 不如一层确定性代码可靠，增加 token/时间成本 | Rejected: 同样的根本弱点 |
| **Go CLI 去重命令** | 确定性代码 | 可靠、可测试、幂等 | 需要写 Go 代码 | **Selected: 直接解决根因** |

## Feasibility Assessment

### Technical Feasibility

Go 标准库 `encoding/json` 完全支持所需的 JSON 操作。集合操作（分组、交集、去重）是 Go 的基本能力。Forge CLI 已有类似的命令组结构（`forge task`、`forge prompt`），可复用模式。

### Resource & Timeline

单个 `sitemap` 命令组（dedup + validate），预计 2-3 个文件，约 300-500 行 Go 代码。工作量可控。

### Dependency Readiness

Forge CLI 编译链已就绪，无需额外依赖。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 更详细的自然语言规则能解决重复问题 | 5 Whys | Overturned: 根因是 LLM 执行确定性操作的不可靠性，不是规则不够详细。追问 5 次后发现问题始终回到"LLM 不可靠地执行集合操作" |
| 独立验证 agent 可以解决去重问题 | XY Detection | Overturned: 用户想要的是"可靠的 sitemap 输出"（Y），独立验证 agent（X）是用另一个 LLM 解决 LLM 的问题。确定性代码更直接 |
| ≥2 页面出现即提升为 layout 元素 | Stress Test | Refined: 严格模式下改为"所有 layout-wrapped 页面都出现才提升"，避免误提升页面特有元素 |
| 去重逻辑应该放在 skill 脚本中 | Assumption Flip | Refined: 内置到 Forge CLI 更可靠——编译型代码不受运行时环境波动影响，且可被其他 skill 复用 |

## Scope

### In Scope

- Go CLI `sitemap` 命令组（`forge-cli/internal/cmd/sitemap/`）
  - `forge sitemap dedup [path]`：五阶段去重算法
  - `forge sitemap validate [path]`：结构验证（不修复）
  - 变更报告输出（JSON 或人类可读格式）
- 修改 `plugins/forge/skills/gen-sitemap/SKILL.md` Step 5：用 `forge sitemap dedup` 调用替代自然语言去重规则
- 简化 `plugins/forge/skills/gen-sitemap/rules/merge-validation.md`：移除已由 Go 代码处理的规则，保留 stale route detection
- 更新 `docs/lessons/pattern-sitemap-shared-layout.md`：记录解决方案

### Out of Scope

- sitemap.json schema 变更（字段名、ID 格式、嵌套结构保持不变）
- 下游消费者变更（write-prd、gen-test-scripts、breakdown-tasks 等无需修改）
- agent-browser 集成变更
- 路由发现或页面探索逻辑变更
- `forge sitemap generate` 命令（当前 gen-sitemap 通过 skill 实现，不在此次范围内）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 严格阈值（所有页面）导致某些共享元素未被提升 | M | M | 在变更报告中明确列出"出现在多页但未提升"的元素，agent 可在审查时手动处理 |
| Go 代码的元素匹配逻辑与 agent 的匹配逻辑不一致 | L | H | Go 代码使用与 SKILL.md 相同的匹配键（role + name），并在变更报告中输出匹配详情供 agent 验证 |
| 增量更新场景下 Go 代码误删已有元素 | L | H | Go 代码只做去重和 ID 重分配，不删除页面或元素；删除操作仍由 agent 在 Step 5c（stale route detection）中处理 |
| `forge sitemap dedup` 成为单点故障（CLI 不可用时整个 skill 失败） | L | M | SKILL.md 中提供 fallback：CLI 不可用时回退到自然语言规则（并标记为降级模式） |

## Success Criteria

- [ ] `forge sitemap dedup` 对含重复节点的 sitemap.json 执行后，输出中零重复元素（无相同 role+name 出现在 page.elements 中且同时存在于 layout.elements）
- [ ] `forge sitemap dedup` 对同一文件连续运行两次，输出完全一致（幂等性）
- [ ] `forge sitemap validate` 正确检测并报告所有四类重复问题（布局重复、页面重复、State 重复、ID 冲突）
- [ ] ID 重分配后所有 E-NNN 连续无间隔，所有 L-NNN 连续无间隔，所有 state.trigger 引用指向有效 ID
- [ ] 处理 500 元素的 sitemap.json 耗时 < 100ms
- [ ] SKILL.md Step 5 的自然语言去重规则已被 `forge sitemap dedup` 调用指令替代
- [ ] merge-validation.md 中仅保留 stale route detection 相关规则

## Next Steps

- Proceed to `/tech-design` to define Go code structure and algorithm implementation details
